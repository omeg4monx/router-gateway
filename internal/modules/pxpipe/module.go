package pxpipe

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName          = "pxpipe"
	interceptorPriority = 35
	maxJSONBodyBytes    = 1 << 20
	maxLogEntries       = 200
	eventTopicStatus    = "pxpipe.status"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Status is returned by GET /api/pxpipe/status and start/stop/restart/config endpoints.
type Status struct {
	Running   bool   `json:"running"`
	Mode      string `json:"mode"`
	StartedAt string `json:"started_at"`
	UptimeMs  int64  `json:"uptime_ms"`
	Config    Config `json:"config"`
}

// HealthCheck describes one step of the live pxpipe self-test.
type HealthCheck struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail"`
}

// HealthReport is returned by GET /api/pxpipe/health.
type HealthReport struct {
	Healthy bool          `json:"healthy"`
	Checks  []HealthCheck `json:"checks"`
	Error   string        `json:"error"`
}

// Stats aggregates live transform pipeline counters.
type Stats struct {
	RequestsSeen         int64 `json:"requests_seen"`
	RequestsTransformed  int64 `json:"requests_transformed"`
	BytesIn              int64 `json:"bytes_in"`
	BytesOut             int64 `json:"bytes_out"`
	BytesSaved           int64 `json:"bytes_saved"`
	EstimatedTokensSaved int64 `json:"estimated_tokens_saved"`
	Errors               int64 `json:"errors"`
}

// LogEntry records one redacted transform summary in the in-memory ring.
type LogEntry struct {
	Timestamp     string   `json:"timestamp"`
	RequestID     string   `json:"request_id"`
	Model         string   `json:"model"`
	BytesBefore   int      `json:"bytes_before"`
	BytesAfter    int      `json:"bytes_after"`
	BytesSaved    int64    `json:"bytes_saved"`
	StagesApplied []string `json:"stages_applied"`
}

// LogListResponse is returned by GET /api/pxpipe/logs.
type LogListResponse struct {
	Logs []LogEntry `json:"logs"`
}

// PreviewResponse is returned by POST /api/pxpipe/preview.
type PreviewResponse struct {
	OriginalBytes        int            `json:"original_bytes"`
	TransformedBytes     int            `json:"transformed_bytes"`
	BytesSaved           int64          `json:"bytes_saved"`
	EstimatedTokensSaved int64          `json:"estimated_tokens_saved"`
	StagesApplied        []string       `json:"stages_applied"`
	Body                 map[string]any `json:"body"`
}

// Options configures a pxpipe Module instance.
type Options struct {
	DB     *db.DB
	Bus    *eventbus.Bus
	Logger *slog.Logger
}

// Module implements lifecycle.Module and contracts.GatewayInterceptor.
type Module struct {
	repo   *repository
	bus    *eventbus.Bus
	logger *slog.Logger

	configMu  sync.Mutex
	mu        sync.RWMutex
	cfg       Config
	running   bool
	startedAt time.Time
	logs      []LogEntry

	reqSeen        atomic.Int64
	reqTransformed atomic.Int64
	bytesIn        atomic.Int64
	bytesOut       atomic.Int64
	bytesSaved     atomic.Int64
	tokensSaved    atomic.Int64
	errCount       atomic.Int64
}

var (
	_ lifecycle.Module             = (*Module)(nil)
	_ contracts.GatewayInterceptor = (*Module)(nil)
)

// New constructs a pxpipe Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Module{
		repo:   newRepository(opts.DB),
		bus:    opts.Bus,
		logger: logger,
		logs:   make([]LogEntry, 0, 32),
	}
}

func (m *Module) Name() string           { return moduleName }
func (m *Module) Dependencies() []string { return []string{"system"} }
func (m *Module) Priority() int          { return interceptorPriority }

func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("pxpipe: sub migrations fs: %v", err))
	}
	return sub
}

func (m *Module) Start(ctx context.Context) error {
	if err := m.repo.ensureDefaults(ctx); err != nil {
		return err
	}
	cfg, err := m.repo.getConfig(ctx)
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.cfg = cfg
	if cfg.AutoStart {
		m.running = true
		m.startedAt = time.Now().UTC()
	}
	m.mu.Unlock()
	return nil
}

func (m *Module) Stop(context.Context) error {
	m.mu.Lock()
	m.running = false
	m.startedAt = time.Time{}
	m.mu.Unlock()
	return nil
}

func (m *Module) snapshot() (cfg Config, running bool, startedAt time.Time) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cfg, m.running, m.startedAt
}

// OnRequest runs the configured in-process transform pipeline when pxpipe is active.
func (m *Module) OnRequest(_ context.Context, req *contracts.CanonicalRequest) error {
	if req == nil || len(req.Body) == 0 {
		return nil
	}
	cfg, running, _ := m.snapshot()
	if !running || !cfg.Enabled {
		return nil
	}

	m.reqSeen.Add(1)
	before := len(req.Body)
	m.bytesIn.Add(int64(before))

	res, err := runTransformPipeline(req.Body, cfg)
	if err != nil {
		m.errCount.Add(1)
		m.bytesOut.Add(int64(before))
		return nil
	}

	if !res.changed {
		m.bytesOut.Add(int64(before))
		return nil
	}
	after := len(res.out)
	m.bytesOut.Add(int64(after))
	req.Body = res.out
	m.reqTransformed.Add(1)
	saved := max(int64(before-after), 0)
	if saved > 0 {
		m.bytesSaved.Add(saved)
		m.tokensSaved.Add((saved + 3) / 4)
	}
	m.appendLog(LogEntry{
		Timestamp:     time.Now().UTC().Format(timestampLayout),
		RequestID:     req.RequestID,
		Model:         req.Model,
		BytesBefore:   before,
		BytesAfter:    after,
		BytesSaved:    saved,
		StagesApplied: res.stagesApplied,
	})
	return nil
}

func (m *Module) OnResponseHeaders(context.Context, *contracts.ResponseMeta) error { return nil }
func (m *Module) OnStreamChunk(context.Context, *contracts.StreamEvent) error      { return nil }

func (m *Module) appendLog(entry LogEntry) {
	if entry.StagesApplied == nil {
		entry.StagesApplied = []string{}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.logs) >= maxLogEntries {
		copy(m.logs, m.logs[1:])
		m.logs[len(m.logs)-1] = entry
		return
	}
	m.logs = append(m.logs, entry)
}

func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Route("/pxpipe", func(r chi.Router) {
		r.Get("/status", m.handleGetStatus)
		r.Patch("/config", m.handlePatchConfig)
		r.Post("/start", m.handleStart)
		r.Post("/stop", m.handleStop)
		r.Post("/restart", m.handleRestart)
		r.Get("/health", m.handleHealth)
		r.Get("/stats", m.handleGetStats)
		r.Delete("/stats", m.handleResetStats)
		r.Get("/logs", m.handleGetLogs)
		r.Delete("/logs", m.handleClearLogs)
		r.Post("/preview", m.handlePreview)
	})
}

func (m *Module) buildStatus() Status {
	cfg, running, startedAt := m.snapshot()
	var startedStr string
	var uptimeMs int64
	if running && !startedAt.IsZero() {
		startedStr = startedAt.Format(time.RFC3339)
		uptimeMs = time.Since(startedAt).Milliseconds()
	}
	return Status{
		Running:   running,
		Mode:      "in-process",
		StartedAt: startedStr,
		UptimeMs:  uptimeMs,
		Config:    cfg,
	}
}

func (m *Module) handleGetStatus(w http.ResponseWriter, _ *http.Request) {
	httpserver.WriteJSON(w, http.StatusOK, m.buildStatus())
}

type patchConfigRequest struct {
	Enabled             *bool `json:"enabled"`
	AutoStart           *bool `json:"auto_start"`
	StripNullFields     *bool `json:"strip_null_fields"`
	CollapseWhitespace  *bool `json:"collapse_whitespace"`
	DedupeSystemPrompts *bool `json:"dedupe_system_prompts"`
	MaxHistoryMessages  *int  `json:"max_history_messages"`
}

func (m *Module) handlePatchConfig(w http.ResponseWriter, r *http.Request) {
	var req patchConfigRequest
	if err := decodeStrictJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	m.configMu.Lock()
	defer m.configMu.Unlock()

	cfg, _, _ := m.snapshot()
	if req.Enabled != nil {
		cfg.Enabled = *req.Enabled
	}
	if req.AutoStart != nil {
		cfg.AutoStart = *req.AutoStart
	}
	if req.StripNullFields != nil {
		cfg.StripNullFields = *req.StripNullFields
	}
	if req.CollapseWhitespace != nil {
		cfg.CollapseWhitespace = *req.CollapseWhitespace
	}
	if req.DedupeSystemPrompts != nil {
		cfg.DedupeSystemPrompts = *req.DedupeSystemPrompts
	}
	if req.MaxHistoryMessages != nil {
		if *req.MaxHistoryMessages < 0 || *req.MaxHistoryMessages > 500 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_max_history", "max_history_messages must be between 0 and 500")
			return
		}
		cfg.MaxHistoryMessages = *req.MaxHistoryMessages
	}
	updated, err := m.repo.updateConfig(r.Context(), cfg)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "pxpipe_config_error", err.Error())
		return
	}
	m.mu.Lock()
	m.cfg = updated
	m.mu.Unlock()
	httpserver.WriteJSON(w, http.StatusOK, m.buildStatus())
}

func (m *Module) handleStart(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	if !m.running {
		m.running = true
		m.startedAt = time.Now().UTC()
	}
	m.mu.Unlock()
	if m.bus != nil {
		m.bus.Publish(eventTopicStatus, map[string]any{"running": true})
	}
	httpserver.WriteJSON(w, http.StatusOK, m.buildStatus())
}

func (m *Module) handleStop(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	m.running = false
	m.startedAt = time.Time{}
	m.mu.Unlock()
	if m.bus != nil {
		m.bus.Publish(eventTopicStatus, map[string]any{"running": false})
	}
	httpserver.WriteJSON(w, http.StatusOK, m.buildStatus())
}

func (m *Module) handleRestart(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	m.running = true
	m.startedAt = time.Now().UTC()
	m.mu.Unlock()
	if m.bus != nil {
		m.bus.Publish(eventTopicStatus, map[string]any{"running": true, "restarted": true})
	}
	httpserver.WriteJSON(w, http.StatusOK, m.buildStatus())
}

func (m *Module) handleHealth(w http.ResponseWriter, _ *http.Request) {
	cfg, running, _ := m.snapshot()
	checks := make([]HealthCheck, 0, 3)

	checks = append(checks, HealthCheck{
		ID:     "config_loaded",
		Label:  "Pipeline configuration loaded",
		OK:     cfg.UpdatedAt != "",
		Detail: "in-process Go engine",
	})

	activeOK := running && cfg.Enabled
	activeDetail := "active"
	if !running {
		activeDetail = "pipeline stopped"
	} else if !cfg.Enabled {
		activeDetail = "pipeline disabled in config"
	}
	checks = append(checks, HealthCheck{
		ID:     "pipeline_active",
		Label:  "Transform pipeline active",
		OK:     activeOK,
		Detail: activeDetail,
	})

	start := time.Now()
	sample := []byte(`{"model":"selftest","unused":null,"messages":[{"role":"user","content":"hello    world"}]}`)
	res, err := runTransformPipeline(sample, Config{
		Enabled:             true,
		StripNullFields:     true,
		CollapseWhitespace:  true,
		DedupeSystemPrompts: true,
	})
	dur := time.Since(start)
	selftestOK := err == nil && res.changed && len(res.out) < len(sample)
	selftestDetail := fmt.Sprintf("%dus (%d -> %d bytes)", dur.Microseconds(), len(sample), len(res.out))
	if err != nil {
		selftestDetail = err.Error()
	}
	checks = append(checks, HealthCheck{
		ID:     "transform_selftest",
		Label:  "Self-test payload transform",
		OK:     selftestOK,
		Detail: selftestDetail,
	})

	healthy := true
	errMsg := ""
	for _, c := range checks {
		if !c.OK {
			healthy = false
			if errMsg == "" {
				errMsg = c.Label + ": " + c.Detail
			}
		}
	}

	httpserver.WriteJSON(w, http.StatusOK, HealthReport{
		Healthy: healthy,
		Checks:  checks,
		Error:   errMsg,
	})
}

func (m *Module) handleGetStats(w http.ResponseWriter, _ *http.Request) {
	httpserver.WriteJSON(w, http.StatusOK, Stats{
		RequestsSeen:         m.reqSeen.Load(),
		RequestsTransformed:  m.reqTransformed.Load(),
		BytesIn:              m.bytesIn.Load(),
		BytesOut:             m.bytesOut.Load(),
		BytesSaved:           m.bytesSaved.Load(),
		EstimatedTokensSaved: m.tokensSaved.Load(),
		Errors:               m.errCount.Load(),
	})
}

func (m *Module) handleResetStats(w http.ResponseWriter, _ *http.Request) {
	m.reqSeen.Store(0)
	m.reqTransformed.Store(0)
	m.bytesIn.Store(0)
	m.bytesOut.Store(0)
	m.bytesSaved.Store(0)
	m.tokensSaved.Store(0)
	m.errCount.Store(0)
	httpserver.WriteJSON(w, http.StatusOK, Stats{})
}

func (m *Module) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	m.mu.RLock()
	n := len(m.logs)
	count := min(n, limit)
	out := make([]LogEntry, 0, count)
	for i := n - 1; i >= n-count; i-- {
		out = append(out, m.logs[i])
	}
	m.mu.RUnlock()
	httpserver.WriteJSON(w, http.StatusOK, LogListResponse{Logs: out})
}

func (m *Module) handleClearLogs(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	deleted := len(m.logs)
	m.logs = m.logs[:0]
	m.mu.Unlock()
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": deleted})
}

type previewRequest struct {
	Payload map[string]any `json:"payload"`
}

func (m *Module) handlePreview(w http.ResponseWriter, r *http.Request) {
	var req previewRequest
	if err := decodeStrictJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if req.Payload == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", "payload object is required")
		return
	}
	raw, err := json.Marshal(req.Payload)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	cfg, _, _ := m.snapshot()
	cfg.Enabled = true
	res, err := runTransformPipeline(raw, cfg)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "pxpipe_preview_error", err.Error())
		return
	}
	saved := max(int64(len(raw)-len(res.out)), 0)
	stages := res.stagesApplied
	if stages == nil {
		stages = []string{}
	}
	body := res.body
	if body == nil {
		body = map[string]any{}
	}
	httpserver.WriteJSON(w, http.StatusOK, PreviewResponse{
		OriginalBytes:        len(raw),
		TransformedBytes:     len(res.out),
		BytesSaved:           saved,
		EstimatedTokensSaved: (saved + 3) / 4,
		StagesApplied:        stages,
		Body:                 body,
	})
}

func decodeStrictJSON(r *http.Request, dst any) error {
	defer func() { _ = r.Body.Close() }()
	dec := json.NewDecoder(io.LimitReader(r.Body, maxJSONBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}
	return nil
}

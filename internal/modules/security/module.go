// Package security implements the security posture summary and redacted
// security audit trail module with automatic eventbus ingestion.
package security

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

const (
	moduleName   = "security"
	maxBodyBytes = 1 << 20 // 1 MiB
	subBufSize   = 64

	// TopicAudit allows any module or platform bridge to emit structured
	// security audit events asynchronously via eventbus.Bus.
	TopicAudit = "security.audit"
)

var watchedTopics = []string{
	TopicAudit,
	"auth.policy_changed",
	"keys.changed",
	"providers.changed",
	"settings.updated",
	"combos.changed",
	"proxypools.changed",
	"promptinjection.changed",
	"modeleditor.changed",
	"tokensaver.changed",
}

// Options configures a security Module instance.
type Options struct {
	DB            *db.DB
	Bus           *eventbus.Bus
	Logger        *slog.Logger
	Config        config.Config
	Store         contracts.AdminCredentialStore
	ActiveModules func() []string
}

// PostureReport summarizes the active security posture of the router instance.
type PostureReport struct {
	BindAddr            string   `json:"bind_addr"`
	IsLoopback          bool     `json:"is_loopback"`
	RequireLogin        bool     `json:"require_login"`
	AllowInsecureRemote bool     `json:"allow_insecure_remote"`
	MasterKeySource     string   `json:"master_key_source"`
	DisabledModules     []string `json:"disabled_modules"`
	ActiveModuleCount   int      `json:"active_module_count"`
	AuditEventCount     int      `json:"audit_event_count"`
	Warnings            []string `json:"warnings"`
	CheckedAt           string   `json:"checked_at"`
}

// Module implements lifecycle.Module for security posture and audit logging.
type Module struct {
	repo          *repo
	bus           *eventbus.Bus
	logger        *slog.Logger
	cfg           config.Config
	store         contracts.AdminCredentialStore
	activeModules func() []string

	mu      sync.Mutex
	started bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a security Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Module{
		repo:          newRepo(opts.DB),
		bus:           opts.Bus,
		logger:        logger,
		cfg:           opts.Config,
		store:         opts.Store,
		activeModules: opts.ActiveModules,
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares modules required before security initializes.
func (m *Module) Dependencies() []string { return []string{"system"} }

// Migrations returns the embedded SQLite migrations for security.
func (m *Module) Migrations() fs.FS { return migrationsFS }

// RegisterRoutes mounts security posture and audit trail routes onto AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Get("/security/posture", m.handlePosture)
	mounts.AdminAPI.Get("/security/audit", m.handleListAudit)
	mounts.AdminAPI.Post("/security/audit", m.handleCreateAudit)
	mounts.AdminAPI.Delete("/security/audit", m.handleClearAudit)
	mounts.AdminAPI.Get("/security/audit/export", m.handleExportAudit)
}

// Start subscribes to control-plane eventbus topics for automatic audit logging.
func (m *Module) Start(ctx context.Context) error {
	if m.bus == nil {
		return nil
	}

	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return nil
	}
	runCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	m.cancel = cancel
	m.started = true

	for _, topic := range watchedTopics {
		ch, unsub := m.bus.Subscribe(topic, subBufSize)
		m.wg.Add(1)
		go m.consumeTopic(runCtx, topic, ch, unsub)
	}
	m.mu.Unlock()

	return nil
}

// Stop cancels eventbus audit listeners and waits for clean shutdown.
func (m *Module) Stop(ctx context.Context) error {
	m.mu.Lock()
	if !m.started {
		m.mu.Unlock()
		return nil
	}
	cancel := m.cancel
	m.started = false
	m.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// PruneExpired removes expired audit events beyond the 90-day TTL and caps
// retained rows at maxRetainedRecords.
func (m *Module) PruneExpired(ctx context.Context) error {
	cutoff := time.Now().UTC().Add(-maxAuditRetention)
	_, err := m.repo.prune(ctx, cutoff, maxRetainedRecords)
	return err
}

func (m *Module) consumeTopic(ctx context.Context, topic string, ch <-chan eventbus.Event, unsub func()) {
	defer m.wg.Done()
	defer unsub()

	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			input := mapBusEventToAudit(topic, evt.Payload)
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			if _, err := m.repo.insert(writeCtx, input); err != nil && ctx.Err() == nil {
				m.logger.Warn("failed to record eventbus audit event", "topic", topic, "error", err)
			}
			cancel()
		}
	}
}

func mapBusEventToAudit(topic string, payload any) AuditInput {
	if topic == TopicAudit {
		switch v := payload.(type) {
		case AuditInput:
			if strings.TrimSpace(v.Category) == "" {
				v.Category = "security"
			}
			if strings.TrimSpace(v.Action) == "" {
				v.Action = "event"
			}
			if strings.TrimSpace(v.Summary) == "" {
				v.Summary = v.Category + " " + v.Action
			}
			return v
		case *AuditInput:
			if v != nil {
				return mapBusEventToAudit(topic, *v)
			}
		}
	}

	meta := payloadToMap(payload)
	switch topic {
	case "auth.policy_changed":
		return AuditInput{
			Category: "auth",
			Action:   "policy.updated",
			Severity: "warning",
			Actor:    "admin",
			Summary:  "Authentication enforcement policy updated",
			Metadata: meta,
		}
	case "keys.changed":
		return AuditInput{
			Category: "keys",
			Action:   "keys.changed",
			Severity: "info",
			Actor:    "admin",
			Summary:  "Gateway API key configuration updated",
			Metadata: meta,
		}
	case "providers.changed":
		return AuditInput{
			Category: "providers",
			Action:   "providers.changed",
			Severity: "info",
			Actor:    "admin",
			Summary:  "Upstream provider connections or nodes updated",
			Metadata: meta,
		}
	case "settings.updated":
		return AuditInput{
			Category: "settings",
			Action:   "settings.updated",
			Severity: "info",
			Actor:    "admin",
			Summary:  "Control-plane runtime settings updated",
			Metadata: meta,
		}
	case "combos.changed":
		return AuditInput{
			Category: "combos",
			Action:   "combos.changed",
			Severity: "info",
			Actor:    "admin",
			Summary:  "Routing combos, aliases, or disabled models updated",
			Metadata: meta,
		}
	case "proxypools.changed":
		return AuditInput{
			Category: "proxypools",
			Action:   "proxypools.changed",
			Severity: "info",
			Actor:    "admin",
			Summary:  "Outbound proxy pool configuration updated",
			Metadata: meta,
		}
	case "promptinjection.changed":
		return AuditInput{
			Category: "promptinjection",
			Action:   "config.changed",
			Severity: "warning",
			Actor:    "admin",
			Summary:  "Prompt injection guard configuration or rules updated",
			Metadata: meta,
		}
	case "modeleditor.changed":
		return AuditInput{
			Category: "modeleditor",
			Action:   "rules.changed",
			Severity: "info",
			Actor:    "admin",
			Summary:  "Model parameter overrides or bypass rules updated",
			Metadata: meta,
		}
	case "tokensaver.changed":
		return AuditInput{
			Category: "tokensaver",
			Action:   "config.changed",
			Severity: "info",
			Actor:    "admin",
			Summary:  "Token saver and Headroom cache configuration updated",
			Metadata: meta,
		}
	default:
		return AuditInput{
			Category: "system",
			Action:   topic,
			Severity: "info",
			Actor:    "system",
			Summary:  fmt.Sprintf("Event emitted on topic %s", topic),
			Metadata: meta,
		}
	}
}

func payloadToMap(payload any) map[string]any {
	if payload == nil {
		return map[string]any{}
	}
	if m, ok := payload.(map[string]any); ok {
		return m
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return map[string]any{"value": fmt.Sprint(payload)}
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return map[string]any{"value": string(raw)}
	}
	return out
}

func (m *Module) handlePosture(w http.ResponseWriter, r *http.Request) {
	isLoopback := m.cfg.IsLoopback()
	requireLogin := true
	if m.store != nil {
		if req, err := m.store.RequireLogin(r.Context()); err == nil {
			requireLogin = req
		}
	}

	masterKeySource := "file"
	if strings.TrimSpace(m.cfg.MasterKeyHex) != "" {
		masterKeySource = "env"
	}

	disabled := make([]string, 0, len(m.cfg.DisabledModules))
	for name, off := range m.cfg.DisabledModules {
		if off {
			disabled = append(disabled, name)
		}
	}
	slices.Sort(disabled)

	activeCount := 0
	if m.activeModules != nil {
		activeCount = len(m.activeModules())
	}

	auditCount, err := m.repo.countAll(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "audit_count_failed", "failed to query audit count")
		return
	}

	warnings := make([]string, 0)
	if !isLoopback && m.cfg.AllowInsecureRemote {
		warnings = append(warnings, "OMEG4MON_ALLOW_INSECURE_REMOTE is enabled on a non-loopback bind address.")
	}
	if !requireLogin {
		if !isLoopback {
			warnings = append(warnings, "Dashboard login is disabled while bound to a non-loopback address.")
		} else {
			warnings = append(warnings, "Dashboard login enforcement is disabled on loopback.")
		}
	}
	if masterKeySource == "file" {
		warnings = append(warnings, "Vault master key is loaded from data directory file; consider setting OMEG4MON_MASTER_KEY via secrets manager in production.")
	}

	httpserver.WriteJSON(w, http.StatusOK, PostureReport{
		BindAddr:            m.cfg.Addr(),
		IsLoopback:          isLoopback,
		RequireLogin:        requireLogin,
		AllowInsecureRemote: m.cfg.AllowInsecureRemote,
		MasterKeySource:     masterKeySource,
		DisabledModules:     disabled,
		ActiveModuleCount:   activeCount,
		AuditEventCount:     auditCount,
		Warnings:            warnings,
		CheckedAt:           formatUTC(time.Now()),
	})
}

func (m *Module) handleListAudit(w http.ResponseWriter, r *http.Request) {
	limit := defaultListLimit
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_limit", "limit must be a positive integer")
			return
		}
		limit = parsed
	}

	items, total, err := m.repo.list(r.Context(), ListFilter{
		Category: r.URL.Query().Get("category"),
		Severity: r.URL.Query().Get("severity"),
		Query:    r.URL.Query().Get("q"),
		Limit:    limit,
	})
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "audit_list_failed", "failed to list audit events")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": total,
	})
}

func (m *Module) handleCreateAudit(w http.ResponseWriter, r *http.Request) {
	var in AuditInput
	if err := decodeStrictJSON(w, r, &in); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	in.Category = strings.TrimSpace(in.Category)
	in.Action = strings.TrimSpace(in.Action)
	in.Summary = strings.TrimSpace(in.Summary)
	if in.Category == "" || in.Action == "" || in.Summary == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_audit_event", "category, action, and summary are required")
		return
	}
	if in.Severity != "" && !isValidSeverity(in.Severity) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_severity", "severity must be info, warning, or critical")
		return
	}
	if strings.TrimSpace(in.Actor) == "" {
		in.Actor = "admin"
	}
	if strings.TrimSpace(in.IPAddress) == "" {
		in.IPAddress = clientIP(r)
	}

	ev, err := m.repo.insert(r.Context(), in)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "audit_create_failed", "failed to create audit event")
		return
	}
	httpserver.WriteJSON(w, http.StatusCreated, ev)
}

func (m *Module) handleClearAudit(w http.ResponseWriter, r *http.Request) {
	deleted, marker, err := m.repo.clearAndRecordMarker(r.Context(), "admin", clientIP(r))
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "audit_clear_failed", "failed to clear audit trail")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"cleared": deleted,
		"marker":  marker,
	})
}

func (m *Module) handleExportAudit(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.listAllForExport(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "audit_export_failed", "failed to export audit trail")
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="security-audit.json"`)
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"exported_at": formatUTC(time.Now()),
		"total":       len(items),
		"items":       items,
	})
}

func isValidSeverity(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "info", "warning", "critical":
		return true
	default:
		return false
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func decodeStrictJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON payload: %w", err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("unexpected trailing data after JSON payload")
	}
	return nil
}

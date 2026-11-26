// Package notifications implements the system notifications inbox, live SSE
// notification stream, and background job observability/trigger endpoints.
package notifications

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
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

const (
	moduleName         = "notifications"
	maxBodyBytes       = 1 << 20 // 1 MiB
	busSubBufSize      = 64
	sseSubBufSize      = 32
	sseKeepalivePeriod = 15 * time.Second

	// TopicEmit allows modules or platform bridges to emit system notifications.
	TopicEmit = "notifications.emit"

	// TopicJobFailed is emitted by the app job observer when a periodic job errors or panics.
	TopicJobFailed = "jobs.failed"
)

var watchedTopics = []string{
	TopicEmit,
	TopicJobFailed,
	"auth.policy_changed",
}

// Options configures a notifications Module instance.
type Options struct {
	DB     *db.DB
	Bus    *eventbus.Bus
	Logger *slog.Logger
	Jobs   contracts.JobInspector
}

// JobStatusDTO is the JSON wire representation of a background job status.
type JobStatusDTO struct {
	Name           string `json:"name"`
	IntervalMs     int64  `json:"interval_ms"`
	Running        bool   `json:"running"`
	RunCount       int64  `json:"run_count"`
	ErrorCount     int64  `json:"error_count"`
	LastStatus     string `json:"last_status"`
	LastError      string `json:"last_error"`
	LastDurationMs int64  `json:"last_duration_ms"`
	LastStartedAt  string `json:"last_started_at"`
	LastFinishedAt string `json:"last_finished_at"`
	NextRunAt      string `json:"next_run_at"`
}

// Module implements lifecycle.Module for system notifications and background jobs UI.
type Module struct {
	repo   *repo
	bus    *eventbus.Bus
	logger *slog.Logger
	jobs   contracts.JobInspector

	mu      sync.Mutex
	started bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup

	streamMu   sync.RWMutex
	streamSubs map[uint64]chan Notification
	nextSubID  uint64
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a notifications Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Module{
		repo:       newRepo(opts.DB),
		bus:        opts.Bus,
		logger:     logger,
		jobs:       opts.Jobs,
		streamSubs: make(map[uint64]chan Notification),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares modules required before notifications initializes.
func (m *Module) Dependencies() []string { return []string{"system"} }

// Migrations returns the embedded SQLite migrations for notifications.
func (m *Module) Migrations() fs.FS { return migrationsFS }

// RegisterRoutes mounts notifications inbox, jobs control, and SSE stream routes.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Get("/notifications", m.handleList)
		mounts.AdminAPI.Post("/notifications", m.handleCreate)
		mounts.AdminAPI.Post("/notifications/read-all", m.handleMarkAllRead)
		mounts.AdminAPI.Post("/notifications/{id}/read", m.handleMarkRead)
		mounts.AdminAPI.Delete("/notifications/{id}", m.handleDelete)
		mounts.AdminAPI.Delete("/notifications", m.handleClear)

		mounts.AdminAPI.Get("/jobs", m.handleListJobs)
		mounts.AdminAPI.Post("/jobs/{name}/trigger", m.handleTriggerJob)
	}

	streamMount := mounts.AdminStreamAPI
	if streamMount == nil {
		streamMount = mounts.AdminAPI
	}
	if streamMount != nil {
		streamMount.Get("/notifications/stream", m.handleStream)
	}
}

// Start subscribes to eventbus topics for automatic system notifications.
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
		ch, unsub := m.bus.Subscribe(topic, busSubBufSize)
		m.wg.Add(1)
		go m.consumeTopic(runCtx, topic, ch, unsub)
	}
	m.mu.Unlock()

	return nil
}

// Stop cancels eventbus listeners and waits for clean shutdown.
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

// PruneExpired removes notifications older than 30 days and caps total rows.
func (m *Module) PruneExpired(ctx context.Context) error {
	cutoff := time.Now().UTC().Add(-maxNotificationAge)
	_, err := m.repo.prune(ctx, cutoff, maxRetainedNotification)
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
			payload := mapBusEventToNotification(topic, evt.Payload)
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			item, err := m.repo.insert(writeCtx, payload)
			cancel()
			if err != nil {
				if ctx.Err() == nil {
					m.logger.Warn("failed to persist eventbus notification", "topic", topic, "error", err)
				}
				continue
			}
			m.broadcast(item)
		}
	}
}

func mapBusEventToNotification(topic string, payload any) EmitPayload {
	switch v := payload.(type) {
	case EmitPayload:
		if strings.TrimSpace(v.Title) == "" {
			v.Title = "System notification"
		}
		if strings.TrimSpace(v.Message) == "" {
			v.Message = v.Title
		}
		return v
	case *EmitPayload:
		if v != nil {
			return mapBusEventToNotification(topic, *v)
		}
	}

	switch topic {
	case "auth.policy_changed":
		return EmitPayload{
			Title:    "Authentication policy changed",
			Message:  "Dashboard login enforcement settings were updated.",
			Severity: "warning",
			Category: "security",
			Source:   "auth",
		}
	case TopicJobFailed:
		return EmitPayload{
			Title:    "Background job failed",
			Message:  fmt.Sprintf("A background job reported an error: %v", payload),
			Severity: "error",
			Category: "jobs",
			Source:   "jobs.runner",
		}
	default:
		return EmitPayload{
			Title:    "System event",
			Message:  fmt.Sprintf("Event emitted on topic %s", topic),
			Severity: "info",
			Category: "system",
			Source:   topic,
		}
	}
}

func (m *Module) subscribeStream() (ch <-chan Notification, unsubscribe func()) {
	c := make(chan Notification, sseSubBufSize)
	m.streamMu.Lock()
	m.nextSubID++
	id := m.nextSubID
	m.streamSubs[id] = c
	m.streamMu.Unlock()

	var once sync.Once
	unsubscribe = func() {
		once.Do(func() {
			m.streamMu.Lock()
			if sub, ok := m.streamSubs[id]; ok {
				delete(m.streamSubs, id)
				close(sub)
			}
			m.streamMu.Unlock()
		})
	}
	return c, unsubscribe
}

func (m *Module) broadcast(item Notification) {
	m.streamMu.RLock()
	defer m.streamMu.RUnlock()
	for _, ch := range m.streamSubs {
		select {
		case ch <- item:
		default:
		}
	}
}

func (m *Module) handleList(w http.ResponseWriter, r *http.Request) {
	limit := defaultListLimit
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_limit", "limit must be a positive integer")
			return
		}
		limit = parsed
	}

	unreadOnly := false
	if raw := strings.TrimSpace(r.URL.Query().Get("unread_only")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_unread_only", "unread_only must be a boolean")
			return
		}
		unreadOnly = parsed
	}

	items, total, unreadCount, err := m.repo.list(r.Context(), ListFilter{
		UnreadOnly: unreadOnly,
		Severity:   r.URL.Query().Get("severity"),
		Limit:      limit,
	})
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "notifications_list_failed", "failed to list notifications")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"items":        items,
		"total":        total,
		"unread_count": unreadCount,
	})
}

func (m *Module) handleCreate(w http.ResponseWriter, r *http.Request) {
	var in EmitPayload
	if err := decodeStrictJSON(w, r, &in); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	in.Title = strings.TrimSpace(in.Title)
	in.Message = strings.TrimSpace(in.Message)
	if in.Title == "" || in.Message == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_notification", "title and message are required")
		return
	}
	if in.Severity != "" && !isValidSeverity(in.Severity) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_severity", "severity must be info, success, warning, or error")
		return
	}
	if strings.TrimSpace(in.Source) == "" {
		in.Source = "admin"
	}

	item, err := m.repo.insert(r.Context(), in)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "notification_create_failed", "failed to create notification")
		return
	}
	m.broadcast(item)
	httpserver.WriteJSON(w, http.StatusCreated, item)
}

func (m *Module) handleMarkRead(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_id", "notification id is required")
		return
	}

	item, err := m.repo.markRead(r.Context(), id)
	if err != nil {
		if errors.Is(err, errNotificationNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "notification_not_found", "notification not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "notification_update_failed", "failed to mark notification read")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, item)
}

func (m *Module) handleMarkAllRead(w http.ResponseWriter, r *http.Request) {
	updated, err := m.repo.markAllRead(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "notifications_update_failed", "failed to mark notifications read")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"updated": updated,
	})
}

func (m *Module) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_id", "notification id is required")
		return
	}

	if err := m.repo.deleteByID(r.Context(), id); err != nil {
		if errors.Is(err, errNotificationNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "notification_not_found", "notification not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "notification_delete_failed", "failed to delete notification")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": true})
}

func (m *Module) handleClear(w http.ResponseWriter, r *http.Request) {
	cleared, err := m.repo.clearAll(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "notifications_clear_failed", "failed to clear notifications")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"cleared": cleared})
}

func (m *Module) handleStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	rc := http.NewResponseController(w)
	ch, unsub := m.subscribeStream()
	defer unsub()

	if _, err := io.WriteString(w, ": connected\n\n"); err != nil {
		return
	}
	_ = rc.Flush()

	keepalive := time.NewTicker(sseKeepalivePeriod)
	defer keepalive.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-keepalive.C:
			if _, err := io.WriteString(w, ": keepalive\n\n"); err != nil {
				return
			}
			_ = rc.Flush()
		case item, ok := <-ch:
			if !ok {
				return
			}
			raw, err := json.Marshal(item)
			if err != nil {
				continue
			}
			if _, err := fmt.Fprintf(w, "event: notification\ndata: %s\n\n", raw); err != nil {
				return
			}
			_ = rc.Flush()
		}
	}
}

func (m *Module) handleListJobs(w http.ResponseWriter, r *http.Request) {
	if m.jobs == nil {
		httpserver.WriteJSON(w, http.StatusOK, map[string]any{"jobs": []JobStatusDTO{}})
		return
	}
	snap := m.jobs.Snapshot()
	out := make([]JobStatusDTO, 0, len(snap))
	for _, st := range snap {
		out = append(out, toJobDTO(st))
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"jobs": out})
}

func (m *Module) handleTriggerJob(w http.ResponseWriter, r *http.Request) {
	if m.jobs == nil {
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "jobs_unavailable", "background job inspector is not configured")
		return
	}

	name := strings.TrimSpace(chi.URLParam(r, "name"))
	if name == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_job_name", "job name is required")
		return
	}

	err := m.jobs.Trigger(r.Context(), name)
	if err != nil {
		switch {
		case errors.Is(err, contracts.ErrJobNotFound):
			httpserver.WriteError(w, r, http.StatusNotFound, "job_not_found", "background job not found")
			return
		case errors.Is(err, contracts.ErrJobAlreadyRunning):
			httpserver.WriteError(w, r, http.StatusConflict, "job_already_running", "background job is already running")
			return
		}
	}

	for _, st := range m.jobs.Snapshot() {
		if st.Name == name {
			httpserver.WriteJSON(w, http.StatusOK, map[string]any{"job": toJobDTO(st)})
			return
		}
	}

	httpserver.WriteError(w, r, http.StatusNotFound, "job_not_found", "background job not found")
}

func toJobDTO(st contracts.JobStatus) JobStatusDTO {
	dto := JobStatusDTO{
		Name:           st.Name,
		IntervalMs:     st.IntervalMs,
		Running:        st.Running,
		RunCount:       st.RunCount,
		ErrorCount:     st.ErrorCount,
		LastStatus:     st.LastStatus,
		LastError:      st.LastError,
		LastDurationMs: st.LastDurationMs,
	}
	if st.LastStartedAt != nil {
		dto.LastStartedAt = formatUTC(*st.LastStartedAt)
	}
	if st.LastFinishedAt != nil {
		dto.LastFinishedAt = formatUTC(*st.LastFinishedAt)
	}
	if st.NextRunAt != nil {
		dto.NextRunAt = formatUTC(*st.NextRunAt)
	}
	return dto
}

func isValidSeverity(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "info", "success", "warning", "error":
		return true
	default:
		return false
	}
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

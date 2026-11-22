// Package settings implements the typed control-plane settings store with
// built-in schema validation, revision history, atomic rollback, and event bus
// notifications.
package settings

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName = "settings"

	// TopicSettingsUpdated is published on the event bus whenever one or more
	// settings are updated or rolled back.
	TopicSettingsUpdated = "settings.updated"

	defaultAdminSubject  = "admin"
	defaultHistoryLimit  = 50
	maxHistoryLimit      = 200
	maxSettingsBodyBytes = 64 << 10 // 64 KiB
	maxReasonRunes       = 512
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// UpdatedEvent is published on TopicSettingsUpdated after committing a
// settings patch or rollback.
type UpdatedEvent struct {
	Keys      []string       `json:"keys"`
	Values    map[string]any `json:"values"`
	ChangedBy string         `json:"changed_by"`
	Reason    string         `json:"reason,omitempty"`
	Rollback  bool           `json:"rollback,omitempty"`
	UpdatedAt string         `json:"updated_at"`
}

// Options configures a new settings Module instance.
type Options struct {
	DB  *db.DB
	Bus *eventbus.Bus
}

// Module implements lifecycle.Module for typed runtime settings management.
type Module struct {
	db   *db.DB
	bus  *eventbus.Bus
	repo *repository
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a settings Module with the provided database and event bus.
func New(opts Options) *Module {
	return &Module{
		db:   opts.DB,
		bus:  opts.Bus,
		repo: newRepository(opts.DB),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string {
	return moduleName
}

// Dependencies declares that settings depends on the system module.
func (m *Module) Dependencies() []string {
	return []string{"system"}
}

// Migrations returns the embedded SQL migration filesystem rooted at migrations/.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("settings: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts the administrative settings endpoints onto AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Get("/settings", m.handleGetSettings)
		mounts.AdminAPI.Patch("/settings", m.handlePatchSettings)
		mounts.AdminAPI.Get("/settings/history", m.handleGetHistory)
		mounts.AdminAPI.Post("/settings/rollback", m.handleRollback)
	}
}

// Start validates that required dependencies are wired.
func (m *Module) Start(context.Context) error {
	if m.db == nil {
		return errors.New("settings: database is required")
	}
	return nil
}

// Stop completes graceful shutdown for the settings module.
func (m *Module) Stop(context.Context) error {
	return nil
}

type settingsResponse struct {
	Settings []SettingItem `json:"settings"`
}

type patchSettingsRequest struct {
	Updates map[string]json.RawMessage `json:"updates"`
	Reason  string                     `json:"reason"`
}

type historyResponse struct {
	Items []HistoryItem `json:"items"`
}

type rollbackRequest struct {
	HistoryID *int64 `json:"history_id"`
	Reason    string `json:"reason"`
}

func (m *Module) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.List(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load settings")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, settingsResponse{
		Settings: items,
	})
}

func (m *Module) handlePatchSettings(w http.ResponseWriter, r *http.Request) {
	var req patchSettingsRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	if len(req.Updates) == 0 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "updates must contain at least one setting")
		return
	}
	if utf8.RuneCountInString(req.Reason) > maxReasonRunes {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", fmt.Sprintf("reason must not exceed %d characters", maxReasonRunes))
		return
	}

	res, err := m.repo.ApplyUpdates(r.Context(), req.Updates, actorFromContext(r.Context()), req.Reason)
	if err != nil {
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			httpserver.WriteError(w, r, http.StatusBadRequest, valErr.Code, valErr.Message)
			return
		}
		if errors.Is(err, ErrEmptyUpdates) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "updates must contain at least one setting")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to update settings")
		return
	}

	if m.bus != nil {
		m.bus.Publish(TopicSettingsUpdated, res.Event)
	}

	httpserver.WriteJSON(w, http.StatusOK, settingsResponse{
		Settings: res.Settings,
	})
}

func (m *Module) handleGetHistory(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(r.URL.Query().Get("key"))
	limit := defaultHistoryLimit

	if rawLimit := strings.TrimSpace(r.URL.Query().Get("limit")); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed <= 0 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "limit must be a positive integer")
			return
		}
		if parsed > maxHistoryLimit {
			parsed = maxHistoryLimit
		}
		limit = parsed
	}

	items, err := m.repo.ListHistory(r.Context(), key, limit)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load settings history")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, historyResponse{
		Items: items,
	})
}

func (m *Module) handleRollback(w http.ResponseWriter, r *http.Request) {
	var req rollbackRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	if req.HistoryID == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "history_id is required")
		return
	}
	if utf8.RuneCountInString(req.Reason) > maxReasonRunes {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", fmt.Sprintf("reason must not exceed %d characters", maxReasonRunes))
		return
	}

	res, err := m.repo.Rollback(r.Context(), *req.HistoryID, actorFromContext(r.Context()), req.Reason)
	if err != nil {
		if errors.Is(err, ErrRevisionNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "revision_not_found", "settings history revision not found")
			return
		}
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			httpserver.WriteError(w, r, http.StatusBadRequest, valErr.Code, valErr.Message)
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to rollback setting")
		return
	}

	if m.bus != nil {
		m.bus.Publish(TopicSettingsUpdated, res.Event)
	}

	httpserver.WriteJSON(w, http.StatusOK, res.Item)
}

func actorFromContext(ctx context.Context) string {
	if subj := strings.TrimSpace(httpserver.AdminSubject(ctx)); subj != "" {
		return subj
	}
	return defaultAdminSubject
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxSettingsBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			httpserver.WriteError(w, r, http.StatusRequestEntityTooLarge, "request_too_large", "request body exceeds maximum allowed size")
			return false
		}
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
		return false
	}
	if dec.More() {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}
	return true
}

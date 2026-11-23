// Package plugins implements configurable gateway request-enhancement plugins
// (via contracts.GatewayInterceptor) and the reusable agent/IDE Skills catalog.
package plugins

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
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName = "plugins"

	// TopicPluginsChanged is published when a plugin configuration or skill changes.
	TopicPluginsChanged = "plugins.changed"

	maxBodyBytes = 128 << 10 // 128 KiB
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var slugRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{1,63}$`)

// Options configures a new plugins Module instance.
type Options struct {
	DB     *db.DB
	Bus    *eventbus.Bus
	Logger *slog.Logger
	Now    func() time.Time
}

// Module implements lifecycle.Module and contracts.GatewayInterceptor.
type Module struct {
	db     *db.DB
	bus    *eventbus.Bus
	logger *slog.Logger
	repo   *repository

	cacheMu sync.RWMutex
	cached  []PluginState
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a plugins Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return &Module{
		db:     opts.DB,
		bus:    opts.Bus,
		logger: logger,
		repo:   newRepository(opts.DB, opts.Now),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares module startup dependencies.
func (m *Module) Dependencies() []string { return []string{"system"} }

// Migrations returns the embedded SQL migrations for plugins.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("plugins: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts the plugins and skills endpoints on AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Get("/plugins", m.handleListPlugins)
	mounts.AdminAPI.Put("/plugins/{key}", m.handleUpdatePlugin)
	mounts.AdminAPI.Post("/plugins/preview", m.handlePreviewPlugins)

	mounts.AdminAPI.Get("/skills", m.handleListSkills)
	mounts.AdminAPI.Post("/skills", m.handleCreateSkill)
	mounts.AdminAPI.Get("/skills/{id}", m.handleGetSkill)
	mounts.AdminAPI.Patch("/skills/{id}", m.handleUpdateSkill)
	mounts.AdminAPI.Delete("/skills/{id}", m.handleDeleteSkill)
	mounts.AdminAPI.Get("/skills/{id}/export", m.handleExportSkill)
}

// Start seeds built-in skills and warms the in-memory plugin cache.
func (m *Module) Start(ctx context.Context) error {
	if m.db == nil {
		return errors.New("plugins: database is required")
	}
	if err := m.repo.SeedBuiltins(ctx); err != nil {
		return err
	}
	_, err := m.refreshCache(ctx)
	return err
}

// Stop completes graceful shutdown.
func (m *Module) Stop(context.Context) error { return nil }

func (m *Module) snapshotStates(ctx context.Context) ([]PluginState, error) {
	m.cacheMu.RLock()
	if m.cached != nil {
		out := make([]PluginState, len(m.cached))
		copy(out, m.cached)
		m.cacheMu.RUnlock()
		return out, nil
	}
	m.cacheMu.RUnlock()
	return m.refreshCache(ctx)
}

func (m *Module) refreshCache(ctx context.Context) ([]PluginState, error) {
	states, err := m.repo.ListPlugins(ctx)
	if err != nil {
		return nil, err
	}
	m.cacheMu.Lock()
	m.cached = make([]PluginState, len(states))
	copy(m.cached, states)
	m.cacheMu.Unlock()
	return states, nil
}

func (m *Module) handleListPlugins(w http.ResponseWriter, r *http.Request) {
	states, err := m.refreshCache(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "plugins_list_failed", "Failed to load plugins")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"plugins": states})
}

type updatePluginRequest struct {
	Enabled bool           `json:"enabled"`
	Models  []string       `json:"models"`
	Config  map[string]any `json:"config"`
}

func (m *Module) handleUpdatePlugin(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(chi.URLParam(r, "key"))
	if key == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_plugin_key", "Plugin key is required")
		return
	}
	var req updatePluginRequest
	if err := decodeJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	updated, err := m.repo.UpsertPlugin(r.Context(), key, req.Enabled, req.Models, req.Config)
	if err != nil {
		if errors.Is(err, ErrPluginNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "plugin_not_found", "Unknown plugin key")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "plugin_update_failed", "Failed to update plugin")
		return
	}
	_, _ = m.refreshCache(r.Context())
	if m.bus != nil {
		m.bus.Publish(TopicPluginsChanged, map[string]any{"action": "plugin.updated", "key": key})
	}
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

type previewPluginsRequest struct {
	Model   string          `json:"model"`
	Payload json.RawMessage `json:"payload"`
}

func (m *Module) handlePreviewPlugins(w http.ResponseWriter, r *http.Request) {
	var req previewPluginsRequest
	if err := decodeJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if len(req.Payload) == 0 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "payload_required", "Payload JSON object is required")
		return
	}
	var check map[string]any
	if err := json.Unmarshal(req.Payload, &check); err != nil || check == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", "Payload must be a JSON object")
		return
	}
	states, err := m.snapshotStates(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "plugins_load_failed", "Failed to load plugins")
		return
	}
	mutatedBytes, applied, changed := applyPluginsToBody(req.Model, req.Payload, states)
	if applied == nil {
		applied = []string{}
	}
	var outPayload any
	_ = json.Unmarshal(mutatedBytes, &outPayload)
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"changed":         changed,
		"applied_plugins": applied,
		"payload":         outPayload,
	})
}

func (m *Module) handleListSkills(w http.ResponseWriter, r *http.Request) {
	skills, err := m.repo.ListSkills(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "skills_list_failed", "Failed to list skills")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"skills": skills})
}

type createSkillRequest struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Content     string `json:"content"`
	Enabled     *bool  `json:"enabled"`
}

func (m *Module) handleCreateSkill(w http.ResponseWriter, r *http.Request) {
	var req createSkillRequest
	if err := decodeJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	slug := strings.ToLower(strings.TrimSpace(req.Slug))
	name := strings.TrimSpace(req.Name)
	content := strings.TrimSpace(req.Content)
	if !slugRe.MatchString(slug) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_slug", "Slug must match ^[a-z0-9][a-z0-9_-]{1,63}$")
		return
	}
	if name == "" || content == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_skill", "Name and content are required")
		return
	}
	category := strings.TrimSpace(req.Category)
	if category == "" {
		category = "general"
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	created, err := m.repo.CreateSkill(r.Context(), slug, name, strings.TrimSpace(req.Description), category, req.Content, enabled)
	if err != nil {
		if errors.Is(err, ErrSkillSlugExists) {
			httpserver.WriteError(w, r, http.StatusConflict, "skill_slug_exists", "A skill with this slug already exists")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "skill_create_failed", "Failed to create skill")
		return
	}
	if m.bus != nil {
		m.bus.Publish(TopicPluginsChanged, map[string]any{"action": "skill.created", "id": created.ID})
	}
	httpserver.WriteJSON(w, http.StatusCreated, created)
}

func (m *Module) handleGetSkill(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	skill, err := m.repo.GetSkill(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrSkillNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "skill_not_found", "Skill not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "skill_get_failed", "Failed to load skill")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, skill)
}

type updateSkillRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Content     *string `json:"content"`
	Enabled     *bool   `json:"enabled"`
}

func (m *Module) handleUpdateSkill(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	var req updateSkillRequest
	if err := decodeJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_name", "Skill name cannot be empty")
		return
	}
	if req.Content != nil && strings.TrimSpace(*req.Content) == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_content", "Skill content cannot be empty")
		return
	}
	updated, err := m.repo.UpdateSkill(r.Context(), id, updateSkillPatch(req))
	if err != nil {
		if errors.Is(err, ErrSkillNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "skill_not_found", "Skill not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "skill_update_failed", "Failed to update skill")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteSkill(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	err := m.repo.DeleteSkill(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrSkillNotFound):
			httpserver.WriteError(w, r, http.StatusNotFound, "skill_not_found", "Skill not found")
		case errors.Is(err, ErrBuiltinSkillProtected):
			httpserver.WriteError(w, r, http.StatusBadRequest, "builtin_skill_protected", "Built-in skills cannot be deleted")
		default:
			httpserver.WriteError(w, r, http.StatusInternalServerError, "skill_delete_failed", "Failed to delete skill")
		}
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id})
}

func (m *Module) handleExportSkill(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	skill, err := m.repo.GetSkill(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrSkillNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "skill_not_found", "Skill not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "skill_export_failed", "Failed to export skill")
		return
	}
	desc := strings.ReplaceAll(skill.Description, "\n", " ")
	doc := fmt.Sprintf("---\nname: %s\ndescription: %q\n---\n\n%s\n", skill.Slug, desc, strings.TrimSpace(skill.Content))
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", skill.Slug+"-SKILL.md"))
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, doc)
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	limited := io.LimitReader(r.Body, maxBodyBytes+1)
	raw, err := io.ReadAll(limited)
	if err != nil {
		return errors.New("failed to read request body")
	}
	if len(raw) > maxBodyBytes {
		return errors.New("request body exceeds size limit")
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return errors.New("invalid JSON payload")
	}
	return nil
}

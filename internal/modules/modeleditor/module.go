package modeleditor

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
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

const (
	moduleName          = "modeleditor"
	interceptorPriority = 20
	maxBodyBytes        = 1 << 20 // 1 MiB
	eventRulesChanged   = "modeleditor.changed"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Options configures a modeleditor Module instance.
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

	mu          sync.RWMutex
	overrides   []ModelOverride
	bypassRules []BypassRule
}

var (
	_ lifecycle.Module             = (*Module)(nil)
	_ contracts.GatewayInterceptor = (*Module)(nil)
)

// New constructs a modeleditor Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Module{
		repo:   newRepository(opts.DB),
		bus:    opts.Bus,
		logger: logger,
	}
}

func (m *Module) Name() string           { return moduleName }
func (m *Module) Dependencies() []string { return []string{"system"} }
func (m *Module) Priority() int          { return interceptorPriority }

func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("modeleditor: sub migrations fs: %v", err))
	}
	return sub
}

func (m *Module) Start(ctx context.Context) error {
	return m.reloadCache(ctx)
}

func (m *Module) Stop(context.Context) error { return nil }

func (m *Module) reloadCache(ctx context.Context) error {
	overrides, err := m.repo.listOverrides(ctx)
	if err != nil {
		return err
	}
	bypassRules, err := m.repo.listBypassRules(ctx)
	if err != nil {
		return err
	}
	for i := range overrides {
		overrides[i].compiledPattern = compileGlobPattern(overrides[i].ModelPattern)
	}
	for i := range bypassRules {
		bypassRules[i].compiledPattern = compileGlobPattern(bypassRules[i].ModelPattern)
	}
	m.mu.Lock()
	m.overrides = overrides
	m.bypassRules = bypassRules
	m.mu.Unlock()
	return nil
}

func (m *Module) reloadCacheAfterWrite(ctx context.Context) {
	reloadCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	if err := m.reloadCache(reloadCtx); err != nil {
		m.logger.Warn("modeleditor: reload cache failed after write", "err", err)
	}
}

func (m *Module) snapshot() ([]BypassRule, []ModelOverride) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.bypassRules, m.overrides
}

// ListVirtualModels returns enabled concrete model overrides that define a new
// virtual model name (non-wildcard model_pattern with a target_model) so they
// appear in /api/arena/models and /v1/models alongside original models.
func (m *Module) ListVirtualModels() []contracts.ModelDescriptor {
	_, overrides := m.snapshot()
	if len(overrides) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(overrides))
	out := make([]contracts.ModelDescriptor, 0, len(overrides))
	for _, ov := range overrides {
		if !ov.Enabled {
			continue
		}
		alias := strings.TrimSpace(ov.ModelPattern)
		target := strings.TrimSpace(ov.TargetModel)
		if alias == "" || target == "" || strings.ContainsAny(alias, "*?[") {
			continue
		}
		if _, dup := seen[alias]; dup {
			continue
		}
		seen[alias] = struct{}{}
		out = append(out, contracts.ModelDescriptor{
			ID:      alias,
			OwnedBy: "model-editor",
			Kind:    "override",
		})
	}
	return out
}

// OnRequest evaluates model bypass rules and applies matching model parameter overrides.
func (m *Module) OnRequest(_ context.Context, req *contracts.CanonicalRequest) error {
	if req == nil {
		return nil
	}
	bypassRules, overrides := m.snapshot()
	if len(bypassRules) == 0 && len(overrides) == 0 {
		return nil
	}
	if req.Metadata == nil {
		req.Metadata = make(map[string]any)
	}
	outModel, outBody, _, _ := applyPipeline(
		req.Protocol, req.Model, req.KeyID, req.Headers, req.Body, req.Metadata, bypassRules, overrides,
	)
	req.Model = outModel
	req.Body = outBody
	return nil
}

// OnResponseHeaders attaches observability headers when bypass or override rules applied.
func (m *Module) OnResponseHeaders(_ context.Context, meta *contracts.ResponseMeta) error {
	if meta == nil || meta.Headers == nil || meta.Metadata == nil {
		return nil
	}
	if v, ok := meta.Metadata["model_bypass_applied"].(string); ok && v != "" {
		meta.Headers.Set("X-Omeg4mon-Model-Bypass", v)
	}
	if v, ok := meta.Metadata["model_override_applied"].(string); ok && v != "" {
		meta.Headers.Set("X-Omeg4mon-Model-Override", v)
	}
	return nil
}

// OnStreamChunk is a no-op for modeleditor.
func (m *Module) OnStreamChunk(context.Context, *contracts.StreamEvent) error { return nil }

func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Route("/model-editor", func(r chi.Router) {
		r.Get("/overrides", m.handleListOverrides)
		r.Post("/overrides", m.handleCreateOverride)
		r.Patch("/overrides/{id}", m.handlePatchOverride)
		r.Delete("/overrides/{id}", m.handleDeleteOverride)

		r.Get("/bypass", m.handleListBypass)
		r.Post("/bypass", m.handleCreateBypass)
		r.Patch("/bypass/{id}", m.handlePatchBypass)
		r.Delete("/bypass/{id}", m.handleDeleteBypass)

		r.Post("/preview", m.handlePreview)
	})
}

func (m *Module) handleListOverrides(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.listOverrides(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to list model overrides")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"overrides": items})
}

func (m *Module) handleCreateOverride(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name               string         `json:"name"`
		ModelPattern       string         `json:"model_pattern"`
		TargetModel        string         `json:"target_model"`
		SystemPromptPrefix string         `json:"system_prompt_prefix"`
		SystemPromptSuffix string         `json:"system_prompt_suffix"`
		Temperature        *float64       `json:"temperature"`
		MaxTokens          *int           `json:"max_tokens"`
		TopP               *float64       `json:"top_p"`
		ExtraParams        map[string]any `json:"extra_params"`
		StripParams        []string       `json:"strip_params"`
		Priority           *int           `json:"priority"`
		Enabled            *bool          `json:"enabled"`
	}
	if !decodeJSONBody(w, r, &body) {
		return
	}
	priority := 100
	if body.Priority != nil {
		priority = *body.Priority
	}
	enabled := true
	if body.Enabled != nil {
		enabled = *body.Enabled
	}
	if body.ExtraParams == nil {
		body.ExtraParams = map[string]any{}
	}
	if body.StripParams == nil {
		body.StripParams = []string{}
	}

	item := ModelOverride{
		ID:                 newID("mov"),
		Name:               strings.TrimSpace(body.Name),
		ModelPattern:       strings.TrimSpace(body.ModelPattern),
		TargetModel:        strings.TrimSpace(body.TargetModel),
		SystemPromptPrefix: body.SystemPromptPrefix,
		SystemPromptSuffix: body.SystemPromptSuffix,
		Temperature:        body.Temperature,
		MaxTokens:          body.MaxTokens,
		TopP:               body.TopP,
		ExtraParams:        body.ExtraParams,
		StripParams:        body.StripParams,
		Priority:           priority,
		Enabled:            enabled,
	}
	if err := validateOverride(item); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_override", err.Error())
		return
	}
	created, err := m.repo.createOverride(r.Context(), item)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to create model override")
		return
	}
	m.reloadCacheAfterWrite(r.Context())
	m.publishChanged()
	httpserver.WriteJSON(w, http.StatusCreated, created)
}

func (m *Module) handlePatchOverride(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := m.repo.getOverride(r.Context(), id)
	if errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "model override not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to load model override")
		return
	}

	var patch struct {
		Name               *string        `json:"name"`
		ModelPattern       *string        `json:"model_pattern"`
		TargetModel        *string        `json:"target_model"`
		SystemPromptPrefix *string        `json:"system_prompt_prefix"`
		SystemPromptSuffix *string        `json:"system_prompt_suffix"`
		Temperature        *float64       `json:"temperature"`
		ClearTemperature   bool           `json:"clear_temperature"`
		MaxTokens          *int           `json:"max_tokens"`
		ClearMaxTokens     bool           `json:"clear_max_tokens"`
		TopP               *float64       `json:"top_p"`
		ClearTopP          bool           `json:"clear_top_p"`
		ExtraParams        map[string]any `json:"extra_params"`
		StripParams        []string       `json:"strip_params"`
		Priority           *int           `json:"priority"`
		Enabled            *bool          `json:"enabled"`
	}
	if !decodeJSONBody(w, r, &patch) {
		return
	}
	if patch.Name != nil {
		existing.Name = strings.TrimSpace(*patch.Name)
	}
	if patch.ModelPattern != nil {
		existing.ModelPattern = strings.TrimSpace(*patch.ModelPattern)
	}
	if patch.TargetModel != nil {
		existing.TargetModel = strings.TrimSpace(*patch.TargetModel)
	}
	if patch.SystemPromptPrefix != nil {
		existing.SystemPromptPrefix = *patch.SystemPromptPrefix
	}
	if patch.SystemPromptSuffix != nil {
		existing.SystemPromptSuffix = *patch.SystemPromptSuffix
	}
	if patch.ClearTemperature {
		existing.Temperature = nil
	} else if patch.Temperature != nil {
		existing.Temperature = patch.Temperature
	}
	if patch.ClearMaxTokens {
		existing.MaxTokens = nil
	} else if patch.MaxTokens != nil {
		existing.MaxTokens = patch.MaxTokens
	}
	if patch.ClearTopP {
		existing.TopP = nil
	} else if patch.TopP != nil {
		existing.TopP = patch.TopP
	}
	if patch.ExtraParams != nil {
		existing.ExtraParams = patch.ExtraParams
	}
	if patch.StripParams != nil {
		existing.StripParams = patch.StripParams
	}
	if patch.Priority != nil {
		existing.Priority = *patch.Priority
	}
	if patch.Enabled != nil {
		existing.Enabled = *patch.Enabled
	}
	if err := validateOverride(existing); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_override", err.Error())
		return
	}

	updated, err := m.repo.updateOverride(r.Context(), existing)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to update model override")
		return
	}
	m.reloadCacheAfterWrite(r.Context())
	m.publishChanged()
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteOverride(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := m.repo.deleteOverride(r.Context(), id)
	if errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "model override not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to delete model override")
		return
	}
	m.reloadCacheAfterWrite(r.Context())
	m.publishChanged()
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handleListBypass(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.listBypassRules(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to list model bypass rules")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"rules": items})
}

func (m *Module) handleCreateBypass(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name         string `json:"name"`
		ModelPattern string `json:"model_pattern"`
		KeyID        string `json:"key_id"`
		Protocol     string `json:"protocol"`
		HeaderName   string `json:"header_name"`
		HeaderValue  string `json:"header_value"`
		Action       string `json:"action"`
		TargetModel  string `json:"target_model"`
		Priority     *int   `json:"priority"`
		Enabled      *bool  `json:"enabled"`
	}
	if !decodeJSONBody(w, r, &body) {
		return
	}
	priority := 100
	if body.Priority != nil {
		priority = *body.Priority
	}
	enabled := true
	if body.Enabled != nil {
		enabled = *body.Enabled
	}
	rule := BypassRule{
		ID:           newID("mbr"),
		Name:         strings.TrimSpace(body.Name),
		ModelPattern: strings.TrimSpace(body.ModelPattern),
		KeyID:        strings.TrimSpace(body.KeyID),
		Protocol:     strings.TrimSpace(body.Protocol),
		HeaderName:   strings.TrimSpace(body.HeaderName),
		HeaderValue:  strings.TrimSpace(body.HeaderValue),
		Action:       strings.TrimSpace(body.Action),
		TargetModel:  strings.TrimSpace(body.TargetModel),
		Priority:     priority,
		Enabled:      enabled,
	}
	if err := validateBypassRule(rule); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_bypass_rule", err.Error())
		return
	}
	created, err := m.repo.createBypassRule(r.Context(), rule)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to create bypass rule")
		return
	}
	m.reloadCacheAfterWrite(r.Context())
	m.publishChanged()
	httpserver.WriteJSON(w, http.StatusCreated, created)
}

func (m *Module) handlePatchBypass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := m.repo.getBypassRule(r.Context(), id)
	if errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "bypass rule not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to load bypass rule")
		return
	}
	var patch struct {
		Name         *string `json:"name"`
		ModelPattern *string `json:"model_pattern"`
		KeyID        *string `json:"key_id"`
		Protocol     *string `json:"protocol"`
		HeaderName   *string `json:"header_name"`
		HeaderValue  *string `json:"header_value"`
		Action       *string `json:"action"`
		TargetModel  *string `json:"target_model"`
		Priority     *int    `json:"priority"`
		Enabled      *bool   `json:"enabled"`
	}
	if !decodeJSONBody(w, r, &patch) {
		return
	}
	if patch.Name != nil {
		existing.Name = strings.TrimSpace(*patch.Name)
	}
	if patch.ModelPattern != nil {
		existing.ModelPattern = strings.TrimSpace(*patch.ModelPattern)
	}
	if patch.KeyID != nil {
		existing.KeyID = strings.TrimSpace(*patch.KeyID)
	}
	if patch.Protocol != nil {
		existing.Protocol = strings.TrimSpace(*patch.Protocol)
	}
	if patch.HeaderName != nil {
		existing.HeaderName = strings.TrimSpace(*patch.HeaderName)
	}
	if patch.HeaderValue != nil {
		existing.HeaderValue = strings.TrimSpace(*patch.HeaderValue)
	}
	if patch.Action != nil {
		existing.Action = strings.TrimSpace(*patch.Action)
	}
	if patch.TargetModel != nil {
		existing.TargetModel = strings.TrimSpace(*patch.TargetModel)
	}
	if patch.Priority != nil {
		existing.Priority = *patch.Priority
	}
	if patch.Enabled != nil {
		existing.Enabled = *patch.Enabled
	}
	if err := validateBypassRule(existing); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_bypass_rule", err.Error())
		return
	}
	updated, err := m.repo.updateBypassRule(r.Context(), existing)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to update bypass rule")
		return
	}
	m.reloadCacheAfterWrite(r.Context())
	m.publishChanged()
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteBypass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := m.repo.deleteBypassRule(r.Context(), id)
	if errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "bypass rule not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to delete bypass rule")
		return
	}
	m.reloadCacheAfterWrite(r.Context())
	m.publishChanged()
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handlePreview(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Protocol string            `json:"protocol"`
		Model    string            `json:"model"`
		KeyID    string            `json:"key_id"`
		Headers  map[string]string `json:"headers"`
		Body     map[string]any    `json:"body"`
	}
	if !decodeJSONBody(w, r, &body) {
		return
	}
	if strings.TrimSpace(body.Model) == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_input", "model is required")
		return
	}
	proto := strings.TrimSpace(body.Protocol)
	if proto == "" {
		proto = "openai-chat"
	}
	if body.Body == nil {
		body.Body = map[string]any{"model": body.Model}
	}
	rawBody, err := json.Marshal(body.Body)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_body", "failed to encode preview body")
		return
	}
	hdr := make(http.Header)
	for k, v := range body.Headers {
		hdr.Set(k, v)
	}
	meta := make(map[string]any)
	bypassRules, overrides := m.snapshot()
	outModel, outRaw, appliedBypass, appliedOverrides := applyPipeline(
		proto, strings.TrimSpace(body.Model), strings.TrimSpace(body.KeyID), hdr, rawBody, meta, bypassRules, overrides,
	)
	outMap := make(map[string]any)
	_ = json.Unmarshal(outRaw, &outMap)
	httpserver.WriteJSON(w, http.StatusOK, PreviewResult{
		Model:              outModel,
		Body:               outMap,
		AppliedOverrides:   appliedOverrides,
		AppliedBypassRules: appliedBypass,
		Metadata:           meta,
	})
}

func (m *Module) publishChanged() {
	if m.bus != nil {
		m.bus.Publish(eventRulesChanged, nil)
	}
}

func validateOverride(ov ModelOverride) error {
	if ov.Name == "" || len(ov.Name) > 128 {
		return errors.New("name must be 1..128 characters")
	}
	if err := validateModelPattern(ov.ModelPattern); err != nil {
		return err
	}
	if ov.Temperature != nil && (*ov.Temperature < 0 || *ov.Temperature > 2) {
		return errors.New("temperature must be between 0 and 2")
	}
	if ov.MaxTokens != nil && (*ov.MaxTokens <= 0 || *ov.MaxTokens > 2000000) {
		return errors.New("max_tokens must be between 1 and 2000000")
	}
	if ov.TopP != nil && (*ov.TopP < 0 || *ov.TopP > 1) {
		return errors.New("top_p must be between 0 and 1")
	}
	return nil
}

func validateBypassRule(br BypassRule) error {
	if br.Name == "" || len(br.Name) > 128 {
		return errors.New("name must be 1..128 characters")
	}
	if err := validateModelPattern(br.ModelPattern); err != nil {
		return err
	}
	switch br.Action {
	case ActionRewriteModel:
		if br.TargetModel == "" {
			return errors.New("target_model is required when action is rewrite_model")
		}
	case ActionBypassOverrides, ActionBypassTokenSaver, ActionBypassHeadroom:
	default:
		return errors.New("action must be rewrite_model, bypass_overrides, bypass_tokensaver, or bypass_headroom")
	}
	return nil
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	limited := io.LimitReader(r.Body, maxBodyBytes+1)
	raw, err := io.ReadAll(limited)
	if err != nil || int64(len(raw)) > maxBodyBytes {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_body", "request body unreadable or exceeds 1 MiB")
		return false
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", "malformed JSON payload")
		return false
	}
	return true
}

func newID(prefix string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return prefix + "_" + hex.EncodeToString(buf)
}

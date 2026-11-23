package promptinjection

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
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName          = "promptinjection"
	interceptorPriority = 10
	maxBodyBytes        = 1 << 20 // 1 MiB
	eventConfigChanged  = "promptinjection.changed"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type httpStatusError struct {
	code int
	msg  string
}

func (e httpStatusError) Error() string   { return e.msg }
func (e httpStatusError) StatusCode() int { return e.code }

// Options configures a promptinjection Module instance.
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

	mu       sync.RWMutex
	cfg      Config
	compiled []compiledRule
}

var (
	_ lifecycle.Module             = (*Module)(nil)
	_ contracts.GatewayInterceptor = (*Module)(nil)
)

// New constructs a promptinjection Module.
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
		panic(fmt.Sprintf("promptinjection: sub migrations fs: %v", err))
	}
	return sub
}

func (m *Module) Start(ctx context.Context) error {
	if err := m.repo.ensureDefaults(ctx); err != nil {
		return err
	}
	return m.reloadCache(ctx)
}

func (m *Module) Stop(context.Context) error { return nil }

func (m *Module) reloadCache(ctx context.Context) error {
	cfg, err := m.repo.getConfig(ctx)
	if err != nil {
		return err
	}
	rules, err := m.repo.listRules(ctx)
	if err != nil {
		return err
	}
	compiled := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		cr, compErr := compileRule(r)
		if compErr != nil {
			m.logger.Warn("skipping invalid prompt guard rule", "rule_id", r.ID, "err", compErr)
			continue
		}
		compiled = append(compiled, cr)
	}
	m.mu.Lock()
	m.cfg = cfg
	m.compiled = compiled
	m.mu.Unlock()
	return nil
}

func (m *Module) snapshot() (Config, []compiledRule) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cfg, m.compiled
}

// OnRequest inspects inbound gateway requests according to the active prompt guard policy.
func (m *Module) OnRequest(ctx context.Context, req *contracts.CanonicalRequest) error {
	if req == nil || len(req.Body) == 0 {
		return nil
	}
	cfg, rules := m.snapshot()
	if !cfg.Enabled || len(rules) == 0 {
		return nil
	}

	sanitize := cfg.Mode == ModeSanitize
	outBody, hits, snippet := inspectAndMaybeSanitizeBody(req.Body, cfg, rules, sanitize)
	if len(hits) == 0 {
		if req.Metadata != nil {
			req.Metadata["prompt_guard_action"] = "pass"
		}
		return nil
	}

	topHit := hits[0]
	action := cfg.Mode
	if req.Metadata != nil {
		req.Metadata["prompt_guard_action"] = action
		req.Metadata["prompt_guard_rule"] = topHit.RuleName
	}

	_ = m.repo.recordEvent(ctx, GuardEvent{
		ID:        newID("pge"),
		RequestID: req.RequestID,
		Model:     req.Model,
		KeyID:     req.KeyID,
		RuleID:    topHit.RuleID,
		RuleName:  topHit.RuleName,
		Action:    action,
		Severity:  topHit.Severity,
		Snippet:   snippet,
	})

	switch action {
	case ModeBlock:
		code := cfg.BlockStatusCode
		if code < 400 || code >= 600 {
			code = http.StatusBadRequest
		}
		return httpStatusError{
			code: code,
			msg:  fmt.Sprintf("prompt guard blocked request (rule: %s)", topHit.RuleName),
		}
	case ModeSanitize:
		req.Body = outBody
		return nil
	default:
		return nil
	}
}

// OnResponseHeaders attaches prompt guard observability header when active.
func (m *Module) OnResponseHeaders(_ context.Context, meta *contracts.ResponseMeta) error {
	if meta == nil || meta.Headers == nil || meta.Metadata == nil {
		return nil
	}
	if action, ok := meta.Metadata["prompt_guard_action"].(string); ok && action != "" {
		meta.Headers.Set("X-Omeg4mon-Prompt-Guard", action)
	}
	return nil
}

// OnStreamChunk is a no-op for prompt guard.
func (m *Module) OnStreamChunk(context.Context, *contracts.StreamEvent) error { return nil }

func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Route("/prompt-guard", func(r chi.Router) {
		r.Get("/config", m.handleGetConfig)
		r.Patch("/config", m.handlePatchConfig)
		r.Get("/rules", m.handleListRules)
		r.Post("/rules", m.handleCreateRule)
		r.Patch("/rules/{id}", m.handlePatchRule)
		r.Delete("/rules/{id}", m.handleDeleteRule)
		r.Post("/test", m.handleTest)
		r.Get("/events", m.handleListEvents)
		r.Delete("/events", m.handleClearEvents)
	})
}

func (m *Module) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := m.repo.getConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to load prompt guard config")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, cfg)
}

func (m *Module) handlePatchConfig(w http.ResponseWriter, r *http.Request) {
	var patch struct {
		Enabled         *bool   `json:"enabled"`
		Mode            *string `json:"mode"`
		Sensitivity     *string `json:"sensitivity"`
		BlockStatusCode *int    `json:"block_status_code"`
	}
	if !decodeJSONBody(w, r, &patch) {
		return
	}

	cfg, err := m.repo.getConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to load prompt guard config")
		return
	}
	if patch.Enabled != nil {
		cfg.Enabled = *patch.Enabled
	}
	if patch.Mode != nil {
		mode := strings.TrimSpace(*patch.Mode)
		if !isValidMode(mode) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_mode", "mode must be block, sanitize, or audit")
			return
		}
		cfg.Mode = mode
	}
	if patch.Sensitivity != nil {
		sens := strings.TrimSpace(*patch.Sensitivity)
		if !isValidSensitivity(sens) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_sensitivity", "sensitivity must be low, medium, or high")
			return
		}
		cfg.Sensitivity = sens
	}
	if patch.BlockStatusCode != nil {
		code := *patch.BlockStatusCode
		if code != http.StatusBadRequest && code != http.StatusForbidden && code != http.StatusUnprocessableEntity {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_status_code", "block_status_code must be 400, 403, or 422")
			return
		}
		cfg.BlockStatusCode = code
	}

	saved, err := m.repo.saveConfig(r.Context(), cfg)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to save prompt guard config")
		return
	}
	_ = m.reloadCache(r.Context())
	m.publishChanged()
	httpserver.WriteJSON(w, http.StatusOK, saved)
}

func (m *Module) handleListRules(w http.ResponseWriter, r *http.Request) {
	rules, err := m.repo.listRules(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to list prompt guard rules")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"rules": rules})
}

func (m *Module) handleCreateRule(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		Pattern     string `json:"pattern"`
		MatchType   string `json:"match_type"`
		Action      string `json:"action"`
		Severity    string `json:"severity"`
		Enabled     *bool  `json:"enabled"`
		Description string `json:"description"`
	}
	if !decodeJSONBody(w, r, &body) {
		return
	}
	enabled := true
	if body.Enabled != nil {
		enabled = *body.Enabled
	}
	rule := Rule{
		ID:          newID("pgr"),
		Name:        strings.TrimSpace(body.Name),
		Pattern:     body.Pattern,
		MatchType:   strings.TrimSpace(body.MatchType),
		Action:      strings.TrimSpace(body.Action),
		Severity:    strings.TrimSpace(body.Severity),
		Enabled:     enabled,
		Builtin:     false,
		Description: strings.TrimSpace(body.Description),
	}
	if err := validateRule(rule); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_rule", err.Error())
		return
	}

	created, err := m.repo.createRule(r.Context(), rule)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to create prompt guard rule")
		return
	}
	_ = m.reloadCache(r.Context())
	m.publishChanged()
	httpserver.WriteJSON(w, http.StatusCreated, created)
}

func (m *Module) handlePatchRule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := m.repo.getRule(r.Context(), id)
	if errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "prompt guard rule not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to load prompt guard rule")
		return
	}

	var patch struct {
		Name        *string `json:"name"`
		Pattern     *string `json:"pattern"`
		MatchType   *string `json:"match_type"`
		Action      *string `json:"action"`
		Severity    *string `json:"severity"`
		Enabled     *bool   `json:"enabled"`
		Description *string `json:"description"`
	}
	if !decodeJSONBody(w, r, &patch) {
		return
	}
	if patch.Enabled != nil {
		existing.Enabled = *patch.Enabled
	}
	if !existing.Builtin {
		if patch.Name != nil {
			existing.Name = strings.TrimSpace(*patch.Name)
		}
		if patch.Pattern != nil {
			existing.Pattern = *patch.Pattern
		}
		if patch.MatchType != nil {
			existing.MatchType = strings.TrimSpace(*patch.MatchType)
		}
		if patch.Action != nil {
			existing.Action = strings.TrimSpace(*patch.Action)
		}
		if patch.Severity != nil {
			existing.Severity = strings.TrimSpace(*patch.Severity)
		}
		if patch.Description != nil {
			existing.Description = strings.TrimSpace(*patch.Description)
		}
	}
	if err := validateRule(existing); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_rule", err.Error())
		return
	}

	updated, err := m.repo.updateRule(r.Context(), existing)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to update prompt guard rule")
		return
	}
	_ = m.reloadCache(r.Context())
	m.publishChanged()
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteRule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := m.repo.deleteRule(r.Context(), id)
	if errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "custom prompt guard rule not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to delete prompt guard rule")
		return
	}
	_ = m.reloadCache(r.Context())
	m.publishChanged()
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handleTest(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Text        string `json:"text"`
		Sensitivity string `json:"sensitivity"`
		Mode        string `json:"mode"`
	}
	if !decodeJSONBody(w, r, &body) {
		return
	}
	if strings.TrimSpace(body.Text) == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_input", "text is required")
		return
	}

	cfg, rules := m.snapshot()
	if isValidSensitivity(body.Sensitivity) {
		cfg.Sensitivity = body.Sensitivity
	}
	if isValidMode(body.Mode) {
		cfg.Mode = body.Mode
	}

	sanitized, hits := evaluateText(body.Text, cfg, rules)
	action := "pass"
	if len(hits) > 0 {
		action = cfg.Mode
	}
	httpserver.WriteJSON(w, http.StatusOK, TestResult{
		Matched:       len(hits) > 0,
		Action:        action,
		Hits:          hits,
		SanitizedText: sanitized,
	})
}

func (m *Module) handleListEvents(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			limit = n
		}
	}
	events, err := m.repo.listEvents(r.Context(), limit)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to list prompt guard events")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"events": events})
}

func (m *Module) handleClearEvents(w http.ResponseWriter, r *http.Request) {
	if err := m.repo.clearEvents(r.Context()); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to clear prompt guard events")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) publishChanged() {
	if m.bus != nil {
		m.bus.Publish(eventConfigChanged, nil)
	}
}

func validateRule(r Rule) error {
	if r.Name == "" || len(r.Name) > 128 {
		return errors.New("name must be 1..128 characters")
	}
	if r.MatchType != MatchSubstring && r.MatchType != MatchRegex {
		return errors.New("match_type must be substring or regex")
	}
	if !isValidMode(r.Action) {
		return errors.New("action must be block, sanitize, or audit")
	}
	if r.Severity != SeverityLow && r.Severity != SeverityMedium && r.Severity != SeverityHigh {
		return errors.New("severity must be low, medium, or high")
	}
	if _, err := compileRule(r); err != nil {
		return err
	}
	return nil
}

func isValidMode(mode string) bool {
	return mode == ModeBlock || mode == ModeSanitize || mode == ModeAudit
}

func isValidSensitivity(s string) bool {
	return s == SensitivityLow || s == SensitivityMedium || s == SensitivityHigh
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

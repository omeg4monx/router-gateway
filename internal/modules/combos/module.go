// Package combos implements fallback combos, model aliases, disabled model
// controls, the 4-tier deterministic ProviderResolver, ModelLister, and the
// thread-safe AccountHealthTracker cooldown state machine.
package combos

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
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName = "combos"

	// TopicCombosChanged is published on the event bus whenever a combo,
	// model alias, disabled model entry, or manual cooldown override changes.
	TopicCombosChanged = "combos.changed"

	strategyOrdered = "ordered"

	defaultAdminSubject = "admin"
	probeUsageKeyID     = "admin-probe"
	defaultComboPrefix  = "combo_"
	defaultProbePrefix  = "probe_"
	idRandomBytes       = 12
	maxCombosBodyBytes  = 64 << 10 // 64 KiB
	maxNameRunes        = 128
	maxDescriptionRunes = 512
	maxTargetsCount     = 64
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// ChangedEvent is published on TopicCombosChanged after committing a routing
// configuration mutation.
type ChangedEvent struct {
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	ChangedBy    string `json:"changed_by"`
	UpdatedAt    string `json:"updated_at"`
}

// Options configures a new combos Module instance.
type Options struct {
	DB         *db.DB
	Providers  contracts.ProviderStore
	HTTPClient *http.Client
	Bus        *eventbus.Bus
	Usage      contracts.UsageRecorder
	Now        func() time.Time
}

// Module implements lifecycle.Module, contracts.ProviderResolver,
// contracts.AccountHealthTracker, and contracts.ModelLister.
type Module struct {
	db         *db.DB
	providers  contracts.ProviderStore
	httpClient *http.Client
	bus        *eventbus.Bus
	usage      contracts.UsageRecorder
	now        func() time.Time
	repo       *repository
	cooldown   *cooldownTracker
	stratState *strategyState
}

var (
	_ lifecycle.Module               = (*Module)(nil)
	_ contracts.ProviderResolver     = (*Module)(nil)
	_ contracts.AccountHealthTracker = (*Module)(nil)
	_ contracts.ModelLister          = (*Module)(nil)
)

// New constructs a combos Module with optional injectable clock for deterministic
// cooldown and timestamp behavior.
func New(opts Options) *Module {
	clock := opts.Now
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	return &Module{
		db:         opts.DB,
		providers:  opts.Providers,
		httpClient: opts.HTTPClient,
		bus:        opts.Bus,
		usage:      opts.Usage,
		now:        clock,
		repo:       newRepository(opts.DB, clock),
		cooldown:   newCooldownTracker(clock),
		stratState: newStrategyState(),
	}
}

func (m *Module) emitUsage(ctx context.Context, rec contracts.UsageRecord) {
	if m.usage == nil {
		return
	}
	if strings.TrimSpace(rec.RequestID) == "" {
		if ctx != nil {
			rec.RequestID = strings.TrimSpace(middleware.GetReqID(ctx))
		}
		if rec.RequestID == "" {
			if id, err := generateID(defaultProbePrefix); err == nil {
				rec.RequestID = id
			} else {
				rec.RequestID = fmt.Sprintf("%s%d", defaultProbePrefix, m.now().UTC().UnixNano())
			}
		}
	}
	if rec.Timestamp.IsZero() {
		rec.Timestamp = m.now().UTC()
	}
	m.usage.Record(ctx, rec)
}

// Name returns the canonical module name.
func (m *Module) Name() string {
	return moduleName
}

// Dependencies declares that combos depends on system and providers.
func (m *Module) Dependencies() []string {
	return []string{"system", "providers"}
}

// Migrations returns the embedded SQL migration filesystem rooted at migrations/.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("combos: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts administrative combo, alias, disabled model, and
// cooldown management routes onto AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}

	mounts.AdminAPI.Get("/combos", m.handleListCombos)
	mounts.AdminAPI.Post("/combos", m.handleCreateCombo)
	mounts.AdminAPI.Post("/combos/test", m.handleTestCombo)
	mounts.AdminAPI.Get("/combos/{id}", m.handleGetCombo)
	mounts.AdminAPI.Patch("/combos/{id}", m.handlePatchCombo)
	mounts.AdminAPI.Delete("/combos/{id}", m.handleDeleteCombo)

	mounts.AdminAPI.Get("/models/catalog", m.handleModelCatalog)
	mounts.AdminAPI.Post("/models/test", m.handleTestModel)
	mounts.AdminAPI.Post("/models/identify", m.handleIdentifyModel)
	mounts.AdminAPI.Post("/models/test-all", m.handleTestAllModels)

	mounts.AdminAPI.Get("/models/aliases", m.handleListAliases)
	mounts.AdminAPI.Put("/models/aliases/{alias}", m.handleUpsertAlias)
	mounts.AdminAPI.Delete("/models/aliases/{alias}", m.handleDeleteAlias)

	mounts.AdminAPI.Get("/models/disabled", m.handleListDisabledModels)
	mounts.AdminAPI.Put("/models/disabled", m.handleUpsertDisabledModel)
	mounts.AdminAPI.Delete("/models/disabled", m.handleDeleteDisabledModel)

	mounts.AdminAPI.Get("/models/cooldowns", m.handleListCooldowns)
	mounts.AdminAPI.Delete("/models/cooldowns/{connectionId}", m.handleDeleteCooldown)

	mounts.AdminAPI.Get("/providers/health-matrix", m.handleProviderHealthMatrix)
}

// Start validates that required dependencies are wired.
func (m *Module) Start(context.Context) error {
	if m.db == nil {
		return errors.New("combos: database is required")
	}
	if m.providers == nil {
		return errors.New("combos: provider store is required")
	}
	return nil
}

// Stop performs graceful module teardown.
func (m *Module) Stop(context.Context) error {
	return nil
}

// IsInCooldown reports whether connectionID is actively cooling down at now.
func (m *Module) IsInCooldown(connectionID string, now time.Time) bool {
	return m.cooldown.IsInCooldown(connectionID, now)
}

// RecordFailure records an upstream failure for connectionID and advances its
// exponential cooldown window.
func (m *Module) RecordFailure(connectionID string, statusCode int, retryAfter time.Duration) {
	m.cooldown.RecordFailure(connectionID, statusCode, retryAfter)
}

// RecordSuccess clears any active cooldown and resets failure backoff immediately.
func (m *Module) RecordSuccess(connectionID string) {
	m.cooldown.RecordSuccess(connectionID)
}

// ListCooldowns returns all currently active connection cooldown entries.
func (m *Module) ListCooldowns(now time.Time) []CooldownEntry {
	return m.cooldown.ListCooldowns(now)
}

// ClearCooldown manually clears any active cooldown state for connectionID.
func (m *Module) ClearCooldown(connectionID string) bool {
	return m.cooldown.ClearCooldown(connectionID)
}

type listCombosResponse struct {
	Combos []Combo `json:"combos"`
}

type createComboRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Targets     []string       `json:"targets"`
	Steps       []ComboStep    `json:"steps"`
	Strategy    string         `json:"strategy"`
	Config      map[string]any `json:"config"`
	Enabled     *bool          `json:"enabled"`
}

type patchComboRequest struct {
	Name        *string         `json:"name"`
	Description *string         `json:"description"`
	Targets     *[]string       `json:"targets"`
	Steps       *[]ComboStep    `json:"steps"`
	Strategy    *string         `json:"strategy"`
	Config      *map[string]any `json:"config"`
	Enabled     *bool           `json:"enabled"`
}

type listAliasesResponse struct {
	Aliases []ModelAlias `json:"aliases"`
}

type upsertAliasRequest struct {
	TargetModel string `json:"target_model"`
}

type listDisabledModelsResponse struct {
	DisabledModels []DisabledModel `json:"disabled_models"`
}

type upsertDisabledModelRequest struct {
	Model  string `json:"model"`
	Reason string `json:"reason"`
}

type deleteDisabledModelRequest struct {
	Model string `json:"model"`
}

type listCooldownsResponse struct {
	Cooldowns []CooldownEntry `json:"cooldowns"`
}

func (m *Module) handleListCombos(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.ListCombos(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to list combos")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, listCombosResponse{Combos: items})
}

func (m *Module) handleCreateCombo(w http.ResponseWriter, r *http.Request) {
	var req createComboRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	name, err := validateModelIdentifier("name", req.Name)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	desc, err := validateDescription("description", req.Description)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	targets, steps, err := normalizeAndValidateSteps(name, req.Targets, req.Steps)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	strategy, err := validateStrategy(req.Strategy)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	id, err := generateID(defaultComboPrefix)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate combo identifier")
		return
	}

	created, err := m.repo.CreateCombo(r.Context(), createComboParams{
		ID:          id,
		Name:        name,
		Description: desc,
		Targets:     targets,
		Steps:       steps,
		Strategy:    strategy,
		Config:      req.Config,
		Enabled:     enabled,
	})
	if err != nil {
		if errors.Is(err, ErrComboAlreadyExists) {
			httpserver.WriteError(w, r, http.StatusConflict, "combo_conflict", "combo name already exists")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to create combo")
		return
	}

	m.publishChanged("created", "combo", created.ID, actorFromContext(r.Context()), created.UpdatedAt)
	httpserver.WriteJSON(w, http.StatusCreated, created)
}

func (m *Module) handleGetCombo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "combo id is required")
		return
	}

	combo, err := m.repo.GetCombo(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrComboNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "combo_not_found", "combo not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load combo")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, combo)
}

func (m *Module) handlePatchCombo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "combo id is required")
		return
	}

	var req patchComboRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	if req.Name == nil &&
		req.Description == nil &&
		req.Targets == nil &&
		req.Steps == nil &&
		req.Strategy == nil &&
		req.Config == nil &&
		req.Enabled == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "at least one field must be provided")
		return
	}

	var (
		cleanName     string
		cleanDesc     string
		cleanStrategy string
		err           error
	)
	if req.Name != nil {
		cleanName, err = validateModelIdentifier("name", *req.Name)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}
	if req.Description != nil {
		cleanDesc, err = validateDescription("description", *req.Description)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}
	if req.Strategy != nil {
		cleanStrategy, err = validateStrategy(*req.Strategy)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}

	var validationErr error
	updated, err := m.repo.UpdateCombo(r.Context(), id, func(cur Combo) (updateComboValues, error) {
		vals := updateComboValues{
			Name:        cur.Name,
			Description: cur.Description,
			Targets:     cur.Targets,
			Steps:       cur.Steps,
			Strategy:    cur.Strategy,
			Config:      cur.Config,
			Enabled:     cur.Enabled,
		}
		if req.Name != nil {
			vals.Name = cleanName
		}
		if req.Description != nil {
			vals.Description = cleanDesc
		}
		if req.Strategy != nil {
			vals.Strategy = cleanStrategy
		}
		if req.Config != nil {
			vals.Config = *req.Config
		}
		if req.Enabled != nil {
			vals.Enabled = *req.Enabled
		}
		var (
			rawTargets []string
			rawSteps   []ComboStep
		)
		switch {
		case req.Steps != nil:
			rawSteps = *req.Steps
		case req.Targets != nil:
			rawTargets = *req.Targets
		case len(cur.Steps) > 0:
			rawSteps = cur.Steps
		default:
			rawTargets = cur.Targets
		}
		cleanTargets, cleanSteps, tErr := normalizeAndValidateSteps(vals.Name, rawTargets, rawSteps)
		if tErr != nil {
			validationErr = tErr
			return updateComboValues{}, tErr
		}
		vals.Targets = cleanTargets
		vals.Steps = cleanSteps
		return vals, nil
	})
	if err != nil {
		if validationErr != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", validationErr.Error())
			return
		}
		if errors.Is(err, ErrComboNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "combo_not_found", "combo not found")
			return
		}
		if errors.Is(err, ErrComboAlreadyExists) {
			httpserver.WriteError(w, r, http.StatusConflict, "combo_conflict", "combo name already exists")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to update combo")
		return
	}

	m.publishChanged("updated", "combo", updated.ID, actorFromContext(r.Context()), updated.UpdatedAt)
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteCombo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "combo id is required")
		return
	}

	deletedAt, err := m.repo.DeleteCombo(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrComboNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "combo_not_found", "combo not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to delete combo")
		return
	}

	m.publishChanged("deleted", "combo", id, actorFromContext(r.Context()), deletedAt)
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handleListAliases(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.ListAliases(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to list model aliases")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, listAliasesResponse{Aliases: items})
}

func (m *Module) handleUpsertAlias(w http.ResponseWriter, r *http.Request) {
	alias, err := decodePathModelParam(chi.URLParam(r, "alias"), "alias")
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	var req upsertAliasRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	targetModel, err := validateModelIdentifier("target_model", req.TargetModel)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if alias == targetModel {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "alias must not reference itself")
		return
	}

	saved, err := m.repo.UpsertAlias(r.Context(), alias, targetModel)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to save model alias")
		return
	}

	m.publishChanged("upserted", "alias", saved.Alias, actorFromContext(r.Context()), saved.UpdatedAt)
	httpserver.WriteJSON(w, http.StatusOK, saved)
}

func (m *Module) handleDeleteAlias(w http.ResponseWriter, r *http.Request) {
	alias, err := decodePathModelParam(chi.URLParam(r, "alias"), "alias")
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	deletedAt, err := m.repo.DeleteAlias(r.Context(), alias)
	if err != nil {
		if errors.Is(err, ErrAliasNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "alias_not_found", "model alias not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to delete model alias")
		return
	}

	m.publishChanged("deleted", "alias", alias, actorFromContext(r.Context()), deletedAt)
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handleListDisabledModels(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.ListDisabledModels(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to list disabled models")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, listDisabledModelsResponse{DisabledModels: items})
}

func (m *Module) handleUpsertDisabledModel(w http.ResponseWriter, r *http.Request) {
	var req upsertDisabledModelRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	model, err := validateModelIdentifier("model", req.Model)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	reason, err := validateDescription("reason", req.Reason)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	saved, err := m.repo.UpsertDisabledModel(r.Context(), model, reason)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to disable model")
		return
	}

	m.publishChanged("upserted", "disabled_model", saved.Model, actorFromContext(r.Context()), m.now().UTC().Format(rfc3339NanoFixed))
	httpserver.WriteJSON(w, http.StatusOK, saved)
}

func (m *Module) handleDeleteDisabledModel(w http.ResponseWriter, r *http.Request) {
	model := strings.TrimSpace(r.URL.Query().Get("model"))
	if model == "" {
		if r.Body == nil || r.Body == http.NoBody {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "model query parameter or JSON body is required")
			return
		}
		var req deleteDisabledModelRequest
		if !decodeJSONBody(w, r, &req) {
			return
		}
		model = strings.TrimSpace(req.Model)
	}

	cleanModel, err := validateModelIdentifier("model", model)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	deletedAt, err := m.repo.DeleteDisabledModel(r.Context(), cleanModel)
	if err != nil {
		if errors.Is(err, ErrDisabledModelNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "disabled_model_not_found", "disabled model not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to enable model")
		return
	}

	m.publishChanged("deleted", "disabled_model", cleanModel, actorFromContext(r.Context()), deletedAt)
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handleListCooldowns(w http.ResponseWriter, r *http.Request) {
	entries := m.ListCooldowns(m.now().UTC())
	httpserver.WriteJSON(w, http.StatusOK, listCooldownsResponse{Cooldowns: entries})
}

func (m *Module) handleDeleteCooldown(w http.ResponseWriter, r *http.Request) {
	connID := strings.TrimSpace(chi.URLParam(r, "connectionId"))
	if connID == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "connectionId is required")
		return
	}

	if m.ClearCooldown(connID) {
		m.publishChanged("cleared", "cooldown", connID, actorFromContext(r.Context()), m.now().UTC().Format(rfc3339NanoFixed))
	}
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) publishChanged(action, resourceType, resourceID, changedBy, updatedAt string) {
	if m.bus == nil {
		return
	}
	m.bus.Publish(TopicCombosChanged, ChangedEvent{
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ChangedBy:    changedBy,
		UpdatedAt:    updatedAt,
	})
}

func actorFromContext(ctx context.Context) string {
	if subj := strings.TrimSpace(httpserver.AdminSubject(ctx)); subj != "" {
		return subj
	}
	return defaultAdminSubject
}

func generateID(prefix string) (string, error) {
	buf := make([]byte, idRandomBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("combos: generate id: %w", err)
	}
	return prefix + hex.EncodeToString(buf), nil
}

func decodePathModelParam(raw, field string) (string, error) {
	unescaped, err := url.PathUnescape(raw)
	if err != nil {
		return "", fmt.Errorf("%s contains invalid URL escape sequence", field)
	}
	return validateModelIdentifier(field, unescaped)
}

func validateModelIdentifier(field, raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("%s must not be empty", field)
	}
	if trimmed == "*" {
		return "", fmt.Errorf("%s must not be wildcard '*'", field)
	}
	if strings.ContainsAny(trimmed, "\r\n\t") {
		return "", fmt.Errorf("%s must not contain control whitespace", field)
	}
	if utf8.RuneCountInString(trimmed) > maxNameRunes {
		return "", fmt.Errorf("%s must not exceed %d characters", field, maxNameRunes)
	}
	return trimmed, nil
}

func validateDescription(field, raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if utf8.RuneCountInString(trimmed) > maxDescriptionRunes {
		return "", fmt.Errorf("%s must not exceed %d characters", field, maxDescriptionRunes)
	}
	return trimmed, nil
}

func validateComboTargets(comboName string, rawTargets []string) ([]string, error) {
	if len(rawTargets) == 0 {
		return nil, errors.New("targets must contain at least one model")
	}
	if len(rawTargets) > maxTargetsCount {
		return nil, fmt.Errorf("targets must not exceed %d items", maxTargetsCount)
	}
	seen := make(map[string]struct{}, len(rawTargets))
	out := make([]string, 0, len(rawTargets))
	for _, raw := range rawTargets {
		clean, err := validateModelIdentifier("targets", raw)
		if err != nil {
			return nil, err
		}
		if clean == comboName {
			return nil, errors.New("targets must not reference the combo itself")
		}
		if _, dup := seen[clean]; dup {
			return nil, fmt.Errorf("targets contains duplicate entry %q", clean)
		}
		seen[clean] = struct{}{}
		out = append(out, clean)
	}
	return out, nil
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxCombosBodyBytes)
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

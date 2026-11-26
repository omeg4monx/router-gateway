// Package keys implements gateway API key lifecycle management (CRUD, CSPRNG
// secret generation, SHA-256 hash persistence, atomic rotation) and the
// read-through cached contracts.KeyValidator implementation.
package keys

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName = "keys"

	// TopicKeysChanged is published on the event bus whenever a gateway API key
	// is created, updated, rotated, or deleted.
	TopicKeysChanged = "keys.changed"

	defaultAdminSubject = "admin"
	maxKeysBodyBytes    = 64 << 10 // 64 KiB
	maxNameRunes        = 128
	maxSliceItems       = 256
	maxSliceItemRunes   = 128
)

var (
	// ErrKeyInvalid indicates a raw API key is empty, malformed, or does not match any stored key.
	ErrKeyInvalid = errors.New("keys: invalid API key")
	// ErrKeyDisabled indicates the matched API key has been administratively disabled.
	ErrKeyDisabled = errors.New("keys: API key is disabled")
	// ErrKeyExpired indicates the matched API key has passed its expires_at timestamp.
	ErrKeyExpired = errors.New("keys: API key has expired")
	// ErrKeyNotFound indicates no API key exists with the requested ID.
	ErrKeyNotFound = errors.New("keys: API key not found")
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// ChangedEvent is published on TopicKeysChanged after committing a key mutation.
type ChangedEvent struct {
	Action    string `json:"action"`
	KeyID     string `json:"key_id"`
	ChangedBy string `json:"changed_by"`
	UpdatedAt string `json:"updated_at"`
}

// Options configures a new keys Module instance.
type Options struct {
	DB  *db.DB
	Bus *eventbus.Bus
}

// Module implements lifecycle.Module, contracts.KeyValidator,
// contracts.KeyProvisioner, and contracts.GatewayInterceptor for gateway API keys.
type Module struct {
	db      *db.DB
	bus     *eventbus.Bus
	repo    *repository
	cache   *keyCache
	limiter *keyLimiter

	lifeMu  sync.Mutex
	started bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

var (
	_ lifecycle.Module             = (*Module)(nil)
	_ contracts.KeyValidator       = (*Module)(nil)
	_ contracts.KeyProvisioner     = (*Module)(nil)
	_ contracts.GatewayInterceptor = (*Module)(nil)
)

// New constructs a keys Module with the provided database and event bus.
func New(opts Options) *Module {
	return &Module{
		db:      opts.DB,
		bus:     opts.Bus,
		repo:    newRepository(opts.DB),
		cache:   newKeyCache(),
		limiter: newKeyLimiter(),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string {
	return moduleName
}

// Dependencies declares that keys depends on the system module.
func (m *Module) Dependencies() []string {
	return []string{"system"}
}

// Migrations returns the embedded SQL migration filesystem rooted at migrations/.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("keys: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts the administrative API key endpoints onto AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Get("/keys", m.handleListKeys)
		mounts.AdminAPI.Post("/keys", m.handleCreateKey)
		mounts.AdminAPI.Get("/keys/{id}", m.handleGetKey)
		mounts.AdminAPI.Patch("/keys/{id}", m.handlePatchKey)
		mounts.AdminAPI.Post("/keys/{id}/rotate", m.handleRotateKey)
		mounts.AdminAPI.Delete("/keys/{id}", m.handleDeleteKey)
	}
}

// Start validates that required dependencies are wired and subscribes to
// usage.recorded events to enforce real-time per-key USD quotas.
func (m *Module) Start(ctx context.Context) error {
	if m.db == nil {
		return errors.New("keys: database is required")
	}
	if m.bus == nil {
		return nil
	}

	m.lifeMu.Lock()
	if m.started {
		m.lifeMu.Unlock()
		return nil
	}
	runCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	m.cancel = cancel
	m.started = true
	ch, unsub := m.bus.Subscribe("usage.recorded", usageSubBufSize)
	m.wg.Add(1)
	go m.consumeUsage(runCtx, ch, unsub)
	m.lifeMu.Unlock()
	return nil
}

// Stop stops the usage.recorded subscriber and clears in-memory caches during graceful shutdown.
func (m *Module) Stop(ctx context.Context) error {
	m.lifeMu.Lock()
	wasStarted := m.started
	cancel := m.cancel
	m.started = false
	m.cancel = nil
	m.lifeMu.Unlock()

	if wasStarted && cancel != nil {
		cancel()
		done := make(chan struct{})
		go func() {
			m.wg.Wait()
			close(done)
		}()
		if ctx == nil {
			<-done
		} else {
			select {
			case <-done:
			case <-ctx.Done():
				m.cache.Clear()
				m.limiter.clear()
				return ctx.Err()
			}
		}
	}

	m.cache.Clear()
	m.limiter.clear()
	return nil
}

// ProvisionKey mints a new gateway API key programmatically on behalf of
// authorized modules (such as importer) and emits keys.changed on the bus.
func (m *Module) ProvisionKey(ctx context.Context, req contracts.ProvisionKeyRequest) (contracts.ProvisionedKey, error) {
	name, err := validateName(req.Name)
	if err != nil {
		return contracts.ProvisionedKey{}, err
	}
	if req.RateLimitRPM < 0 {
		return contracts.ProvisionedKey{}, errors.New("rate_limit_rpm must be non-negative")
	}
	if req.QuotaUSD < 0 || math.IsNaN(req.QuotaUSD) || math.IsInf(req.QuotaUSD, 0) {
		return contracts.ProvisionedKey{}, errors.New("quota_usd must be a non-negative finite number")
	}

	allowedModels, err := validateStringSlice("allowed_models", req.AllowedModels)
	if err != nil {
		return contracts.ProvisionedKey{}, err
	}
	scopes := req.Scopes
	if len(scopes) == 0 {
		scopes = []string{"gateway"}
	} else {
		scopes, err = validateStringSlice("scopes", scopes)
		if err != nil {
			return contracts.ProvisionedKey{}, err
		}
	}
	tags, err := validateStringSlice("tags", req.Tags)
	if err != nil {
		return contracts.ProvisionedKey{}, err
	}
	expiresAt, err := normalizeExpiresAt(req.ExpiresAt)
	if err != nil {
		return contracts.ProvisionedKey{}, err
	}

	keyID, err := generateKeyID()
	if err != nil {
		return contracts.ProvisionedKey{}, fmt.Errorf("generate key id: %w", err)
	}
	rawSecret, keyHash, keyPrefix, err := generateToken()
	if err != nil {
		return contracts.ProvisionedKey{}, fmt.Errorf("generate key secret: %w", err)
	}

	summary, err := m.repo.Create(ctx, createParams{
		ID:            keyID,
		Name:          name,
		KeyHash:       keyHash,
		KeyPrefix:     keyPrefix,
		AllowedModels: allowedModels,
		Scopes:        scopes,
		Tags:          tags,
		RateLimitRPM:  req.RateLimitRPM,
		QuotaUSD:      req.QuotaUSD,
		ExpiresAt:     expiresAt,
	})
	if err != nil {
		return contracts.ProvisionedKey{}, fmt.Errorf("create provisioned key: %w", err)
	}

	actor := strings.TrimSpace(req.Actor)
	if actor == "" {
		actor = "system"
	}
	m.publishChanged("created", summary.ID, actor, summary.UpdatedAt)

	return contracts.ProvisionedKey{
		ID:        summary.ID,
		Name:      summary.Name,
		KeyPrefix: summary.KeyPrefix,
		Secret:    rawSecret,
	}, nil
}

// RevokeKey deletes a provisioned API key by ID and invalidates the validation cache.
func (m *Module) RevokeKey(ctx context.Context, keyID string) error {
	id := strings.TrimSpace(keyID)
	if id == "" {
		return ErrKeyNotFound
	}
	res, err := m.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	m.cache.Delete(res.OldHash)
	m.limiter.purgeKey(id)
	m.publishChanged("deleted", id, "system", res.DeletedAt)
	return nil
}

// ValidateKey validates rawKey by SHA-256 hash using an in-memory read-through
// cache backed by SQLite, rejecting invalid, disabled, or expired keys.
func (m *Module) ValidateKey(ctx context.Context, rawKey string) (contracts.KeyInfo, error) {
	if !isValidTokenFormat(rawKey) {
		return contracts.KeyInfo{}, ErrKeyInvalid
	}
	hash := hashToken(rawKey)

	entry, ok := m.cache.Get(hash)
	if !ok {
		ver := m.cache.Version()
		loaded, err := m.repo.LookupByHash(ctx, hash)
		if err != nil {
			return contracts.KeyInfo{}, err
		}
		m.cache.SetIfVersion(hash, loaded, ver)
		entry = cloneCachedEntry(loaded)
	}

	if entry.info.Disabled {
		return contracts.KeyInfo{}, ErrKeyDisabled
	}
	if entry.expiresAt != nil && !time.Now().UTC().Before(*entry.expiresAt) {
		return contracts.KeyInfo{}, ErrKeyExpired
	}

	m.limiter.ensureLimits(entry.info.ID, keyLimitSnapshot{
		name:         entry.info.Name,
		rateLimitRPM: entry.info.RateLimitRPM,
		quotaUSD:     entry.info.QuotaUSD,
	})
	m.recordLastUsed(ctx, entry.info.ID)

	return cloneKeyInfo(entry.info), nil
}

type listKeysResponse struct {
	Keys []APIKeySummary `json:"keys"`
}

type keySecretResponse struct {
	Key    APIKeySummary `json:"key"`
	Secret string        `json:"secret"`
}

type createKeyRequest struct {
	Name          string   `json:"name"`
	AllowedModels []string `json:"allowed_models"`
	Scopes        []string `json:"scopes"`
	Tags          []string `json:"tags"`
	RateLimitRPM  int      `json:"rate_limit_rpm"`
	QuotaUSD      float64  `json:"quota_usd"`
	ExpiresAt     *string  `json:"expires_at"`
}

type nullableString struct {
	Set   bool
	Value *string
}

func (n *nullableString) UnmarshalJSON(data []byte) error {
	n.Set = true
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	n.Value = &s
	return nil
}

type patchKeyRequest struct {
	Name          *string        `json:"name"`
	AllowedModels *[]string      `json:"allowed_models"`
	Scopes        *[]string      `json:"scopes"`
	Tags          *[]string      `json:"tags"`
	RateLimitRPM  *int           `json:"rate_limit_rpm"`
	QuotaUSD      *float64       `json:"quota_usd"`
	Disabled      *bool          `json:"disabled"`
	ExpiresAt     nullableString `json:"expires_at"`
}

func (m *Module) handleListKeys(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.List(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to list API keys")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, listKeysResponse{
		Keys: items,
	})
}

func (m *Module) handleCreateKey(w http.ResponseWriter, r *http.Request) {
	var req createKeyRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	name, err := validateName(req.Name)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if req.RateLimitRPM < 0 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "rate_limit_rpm must be non-negative")
		return
	}
	if req.QuotaUSD < 0 || math.IsNaN(req.QuotaUSD) || math.IsInf(req.QuotaUSD, 0) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "quota_usd must be a non-negative finite number")
		return
	}

	allowedModels, err := validateStringSlice("allowed_models", req.AllowedModels)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	var scopes []string
	if req.Scopes == nil {
		scopes = []string{"gateway"}
	} else {
		scopes, err = validateStringSlice("scopes", req.Scopes)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}
	tags, err := validateStringSlice("tags", req.Tags)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	expiresAt, err := normalizeExpiresAt(req.ExpiresAt)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	keyID, err := generateKeyID()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate key identifier")
		return
	}
	rawSecret, keyHash, keyPrefix, err := generateToken()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate API key secret")
		return
	}

	summary, err := m.repo.Create(r.Context(), createParams{
		ID:            keyID,
		Name:          name,
		KeyHash:       keyHash,
		KeyPrefix:     keyPrefix,
		AllowedModels: allowedModels,
		Scopes:        scopes,
		Tags:          tags,
		RateLimitRPM:  req.RateLimitRPM,
		QuotaUSD:      req.QuotaUSD,
		ExpiresAt:     expiresAt,
	})
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to create API key")
		return
	}

	m.publishChanged("created", summary.ID, actorFromContext(r.Context()), summary.UpdatedAt)

	httpserver.WriteJSON(w, http.StatusCreated, keySecretResponse{
		Key:    summary,
		Secret: rawSecret,
	})
}

func (m *Module) handleGetKey(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "key id is required")
		return
	}

	summary, err := m.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "key_not_found", "API key not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load API key")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, summary)
}

func (m *Module) handlePatchKey(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "key id is required")
		return
	}

	var req patchKeyRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	if req.Name == nil &&
		req.AllowedModels == nil &&
		req.Scopes == nil &&
		req.Tags == nil &&
		req.RateLimitRPM == nil &&
		req.QuotaUSD == nil &&
		req.Disabled == nil &&
		!req.ExpiresAt.Set {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "at least one field must be provided")
		return
	}

	params := updateParams{
		Disabled:     req.Disabled,
		ExpiresAtSet: req.ExpiresAt.Set,
	}
	if req.Name != nil {
		cleanName, err := validateName(*req.Name)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		params.Name = &cleanName
	}
	if req.AllowedModels != nil {
		cleanModels, err := validateStringSlice("allowed_models", *req.AllowedModels)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		params.AllowedModels = &cleanModels
	}
	if req.Scopes != nil {
		cleanScopes, err := validateStringSlice("scopes", *req.Scopes)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		params.Scopes = &cleanScopes
	}
	if req.Tags != nil {
		cleanTags, err := validateStringSlice("tags", *req.Tags)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		params.Tags = &cleanTags
	}
	if req.RateLimitRPM != nil {
		if *req.RateLimitRPM < 0 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "rate_limit_rpm must be non-negative")
			return
		}
		rpm := *req.RateLimitRPM
		params.RateLimitRPM = &rpm
	}
	if req.QuotaUSD != nil {
		q := *req.QuotaUSD
		if q < 0 || math.IsNaN(q) || math.IsInf(q, 0) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "quota_usd must be a non-negative finite number")
			return
		}
		params.QuotaUSD = &q
	}
	if req.ExpiresAt.Set {
		normExp, err := normalizeExpiresAt(req.ExpiresAt.Value)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		params.ExpiresAt = normExp
	}

	res, err := m.repo.Update(r.Context(), id, params)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "key_not_found", "API key not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to update API key")
		return
	}

	// Immediately invalidate cache before emitting event bus notifications.
	m.cache.Delete(res.KeyHash)
	m.limiter.invalidateLimits(res.Summary.ID)
	m.publishChanged("updated", res.Summary.ID, actorFromContext(r.Context()), res.Summary.UpdatedAt)

	httpserver.WriteJSON(w, http.StatusOK, res.Summary)
}

func (m *Module) handleRotateKey(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "key id is required")
		return
	}

	newRaw, newHash, newPrefix, err := generateToken()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate API key secret")
		return
	}

	res, err := m.repo.Rotate(r.Context(), id, newHash, newPrefix)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "key_not_found", "API key not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to rotate API key")
		return
	}

	// Immediately invalidate both old and new hashes before publishing event.
	m.cache.Delete(res.OldHash, res.NewHash)
	m.publishChanged("rotated", res.Summary.ID, actorFromContext(r.Context()), res.Summary.UpdatedAt)

	httpserver.WriteJSON(w, http.StatusOK, keySecretResponse{
		Key:    res.Summary,
		Secret: newRaw,
	})
}

func (m *Module) handleDeleteKey(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "key id is required")
		return
	}

	res, err := m.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "key_not_found", "API key not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to delete API key")
		return
	}

	// Immediately invalidate deleted hash before publishing event.
	m.cache.Delete(res.OldHash)
	m.limiter.purgeKey(id)
	m.publishChanged("deleted", id, actorFromContext(r.Context()), res.DeletedAt)

	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) publishChanged(action, keyID, changedBy, updatedAt string) {
	if m.bus == nil {
		return
	}
	m.bus.Publish(TopicKeysChanged, ChangedEvent{
		Action:    action,
		KeyID:     keyID,
		ChangedBy: changedBy,
		UpdatedAt: updatedAt,
	})
}

func actorFromContext(ctx context.Context) string {
	if subj := strings.TrimSpace(httpserver.AdminSubject(ctx)); subj != "" {
		return subj
	}
	return defaultAdminSubject
}

func validateName(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("name must not be empty")
	}
	if utf8.RuneCountInString(trimmed) > maxNameRunes {
		return "", fmt.Errorf("name must not exceed %d characters", maxNameRunes)
	}
	return trimmed, nil
}

func validateStringSlice(field string, items []string) ([]string, error) {
	if len(items) == 0 {
		return []string{}, nil
	}
	if len(items) > maxSliceItems {
		return nil, fmt.Errorf("%s must not exceed %d items", field, maxSliceItems)
	}
	out := make([]string, 0, len(items))
	for _, raw := range items {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			return nil, fmt.Errorf("%s entries must not be empty", field)
		}
		if utf8.RuneCountInString(trimmed) > maxSliceItemRunes {
			return nil, fmt.Errorf("%s entry must not exceed %d characters", field, maxSliceItemRunes)
		}
		out = append(out, trimmed)
	}
	return out, nil
}

func normalizeExpiresAt(raw *string) (*string, error) {
	if raw == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil, errors.New("expires_at must be a valid RFC3339 timestamp or null")
	}
	parsed, err := time.Parse(time.RFC3339Nano, trimmed)
	if err != nil {
		return nil, errors.New("expires_at must be a valid RFC3339 timestamp")
	}
	formatted := parsed.UTC().Format(rfc3339NanoFixed)
	return &formatted, nil
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxKeysBodyBytes)
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

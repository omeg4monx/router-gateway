package proxypools

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
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName          = "proxypools"
	defaultProbeTimeout = 5 * time.Second
	maxBodyBytes        = 1 << 20 // 1 MiB
	eventPoolsChanged   = "proxypools.changed"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Options configures a proxypools Module instance.
type Options struct {
	DB           *db.DB
	Vault        *crypto.Vault
	Bus          *eventbus.Bus
	Logger       *slog.Logger
	ProbeTimeout time.Duration
}

// Module implements lifecycle.Module and contracts.ProxyDialerResolver.
type Module struct {
	repo         *repository
	vault        *crypto.Vault
	bus          *eventbus.Bus
	logger       *slog.Logger
	probeTimeout time.Duration

	mu          sync.RWMutex
	cacheLoaded bool
	boundURLs   map[string]*url.URL
}

var (
	_ lifecycle.Module              = (*Module)(nil)
	_ contracts.ProxyDialerResolver = (*Module)(nil)
)

// New constructs a proxypools Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	timeout := opts.ProbeTimeout
	if timeout <= 0 {
		timeout = defaultProbeTimeout
	}
	return &Module{
		repo:         newRepository(opts.DB),
		vault:        opts.Vault,
		bus:          opts.Bus,
		logger:       logger,
		probeTimeout: timeout,
		boundURLs:    make(map[string]*url.URL),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares the upstream modules required by proxypools.
func (m *Module) Dependencies() []string { return []string{"system", "providers"} }

// Migrations returns the embedded SQL migration filesystem rooted at migrations/.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("proxypools: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts proxy pool management endpoints onto mounts.AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Get("/proxypools", m.handleList)
	mounts.AdminAPI.Post("/proxypools", m.handleCreate)
	mounts.AdminAPI.Patch("/proxypools/{id}", m.handleUpdate)
	mounts.AdminAPI.Delete("/proxypools/{id}", m.handleDelete)
	mounts.AdminAPI.Post("/proxypools/{id}/test", m.handleTest)
	mounts.AdminAPI.Post("/proxypools/scan", m.handleScan)
	mounts.AdminAPI.Put("/proxypools/bindings/{providerId}", m.handleUpsertBinding)
	mounts.AdminAPI.Delete("/proxypools/bindings/{providerId}", m.handleDeleteBinding)
}

// Start pre-warms the in-memory provider proxy binding cache.
func (m *Module) Start(ctx context.Context) error {
	return m.reloadCache(ctx)
}

// Stop releases in-memory state cleanly.
func (m *Module) Stop(_ context.Context) error {
	m.mu.Lock()
	m.boundURLs = make(map[string]*url.URL)
	m.cacheLoaded = false
	m.mu.Unlock()
	return nil
}

// ProxyURLForProvider resolves the outbound proxy URL bound to providerID, or
// returns (nil, nil) if unbound or bound to a disabled pool.
func (m *Module) ProxyURLForProvider(ctx context.Context, providerID string) (*url.URL, error) {
	pid := strings.TrimSpace(providerID)
	if pid == "" {
		return nil, nil
	}

	m.mu.RLock()
	loaded := m.cacheLoaded
	u := m.boundURLs[pid]
	m.mu.RUnlock()

	if !loaded {
		if err := m.reloadCache(ctx); err != nil {
			return nil, err
		}
		m.mu.RLock()
		u = m.boundURLs[pid]
		m.mu.RUnlock()
	}
	if u == nil {
		return nil, nil
	}
	cp := *u
	return &cp, nil
}

func (m *Module) reloadCache(ctx context.Context) error {
	pools, err := m.repo.listPools(ctx)
	if err != nil {
		return err
	}
	bindings, err := m.repo.listBindings(ctx)
	if err != nil {
		return err
	}

	enabledURLs := make(map[string]*url.URL, len(pools))
	for _, p := range pools {
		if !p.Enabled {
			continue
		}
		raw, decErr := m.decryptURL(p.URLEncrypted)
		if decErr != nil {
			m.logger.Warn("proxypools: decrypt pool url failed", "pool_id", p.ID, "err", decErr)
			continue
		}
		parsed, parseErr := url.Parse(raw)
		if parseErr != nil {
			continue
		}
		enabledURLs[p.ID] = parsed
	}

	next := make(map[string]*url.URL, len(bindings))
	for _, b := range bindings {
		if u, ok := enabledURLs[b.PoolID]; ok {
			next[b.ProviderID] = u
		}
	}

	m.mu.Lock()
	m.boundURLs = next
	m.cacheLoaded = true
	m.mu.Unlock()
	return nil
}

func (m *Module) invalidateAndPublish(ctx context.Context) {
	if err := m.reloadCache(ctx); err != nil {
		m.logger.Warn("proxypools: reload cache failed", "err", err)
	}
	if m.bus != nil {
		m.bus.Publish(eventPoolsChanged, map[string]any{"updated_at": formatUTC(time.Now())})
	}
}

func (m *Module) encryptURL(raw string) (string, error) {
	if m.vault == nil {
		return "", errors.New("proxypools: secret vault unavailable")
	}
	return m.vault.Encrypt([]byte(raw))
}

func (m *Module) decryptURL(enc string) (string, error) {
	if m.vault == nil {
		return "", errors.New("proxypools: secret vault unavailable")
	}
	plain, err := m.vault.Decrypt(enc)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func (m *Module) handleList(w http.ResponseWriter, r *http.Request) {
	pools, err := m.repo.listPools(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "proxypools_list_failed", err.Error())
		return
	}
	bindings, err := m.repo.listBindings(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "proxypools_bindings_failed", err.Error())
		return
	}
	summaries := make([]PoolSummary, 0, len(pools))
	for _, p := range pools {
		summaries = append(summaries, p.PoolSummary)
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"pools":    summaries,
		"bindings": bindings,
	})
}

func (m *Module) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		Enabled *bool  `json:"enabled"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "name must not be empty")
		return
	}
	parsed, err := parseAndValidateProxyURL(req.URL)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_proxy_url", err.Error())
		return
	}

	enc, err := m.encryptURL(parsed.Raw)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_error", err.Error())
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	now := formatUTC(time.Now())
	item := storedPool{
		PoolSummary: PoolSummary{
			ID:             newPoolID(),
			Name:           name,
			Scheme:         parsed.Scheme,
			DisplayURL:     parsed.DisplayURL,
			HasCredentials: parsed.HasCredentials,
			Enabled:        enabled,
			Status:         "unknown",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		URLEncrypted: enc,
	}

	if err := m.repo.createPool(r.Context(), item); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "proxypools_create_failed", err.Error())
		return
	}
	m.invalidateAndPublish(r.Context())
	httpserver.WriteJSON(w, http.StatusCreated, item.PoolSummary)
}

func (m *Module) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	existing, err := m.repo.getPool(r.Context(), id)
	if errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "proxy pool not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "proxypools_get_failed", err.Error())
		return
	}

	var req struct {
		Name    *string `json:"name"`
		URL     *string `json:"url"`
		Enabled *bool   `json:"enabled"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "name must not be empty")
			return
		}
		existing.Name = name
	}
	if req.URL != nil {
		parsed, parseErr := parseAndValidateProxyURL(*req.URL)
		if parseErr != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_proxy_url", parseErr.Error())
			return
		}
		enc, encErr := m.encryptURL(parsed.Raw)
		if encErr != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_error", encErr.Error())
			return
		}
		existing.Scheme = parsed.Scheme
		existing.DisplayURL = parsed.DisplayURL
		existing.HasCredentials = parsed.HasCredentials
		existing.URLEncrypted = enc
	}
	if req.Enabled != nil {
		existing.Enabled = *req.Enabled
	}
	existing.UpdatedAt = formatUTC(time.Now())

	if err := m.repo.updatePool(r.Context(), existing); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "proxypools_update_failed", err.Error())
		return
	}
	m.invalidateAndPublish(r.Context())
	httpserver.WriteJSON(w, http.StatusOK, existing.PoolSummary)
}

func (m *Module) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	err := m.repo.deletePool(r.Context(), id)
	if errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "proxy pool not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "proxypools_delete_failed", err.Error())
		return
	}
	m.invalidateAndPublish(r.Context())
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handleTest(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	pool, err := m.repo.getPool(r.Context(), id)
	if errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "proxy pool not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "proxypools_get_failed", err.Error())
		return
	}

	var req struct {
		TargetURL string `json:"target_url"`
		TimeoutMS int    `json:"timeout_ms"`
	}
	if r.Body != nil {
		if decErr := decodeOptionalJSON(w, r, &req); decErr != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", decErr.Error())
			return
		}
	}

	raw, err := m.decryptURL(pool.URLEncrypted)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_error", err.Error())
		return
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "invalid_stored_url", err.Error())
		return
	}

	timeout := m.probeTimeout
	if req.TimeoutMS > 0 {
		timeout = min(time.Duration(req.TimeoutMS)*time.Millisecond, 30*time.Second)
	}

	lat, probeErr := probeProxy(r.Context(), parsed, strings.TrimSpace(req.TargetURL), timeout)
	testedAt := formatUTC(time.Now())
	status := "healthy"
	errMsg := ""
	if probeErr != nil {
		status = "unhealthy"
		errMsg = probeErr.Error()
	}

	if updateErr := m.repo.updatePoolHealth(r.Context(), pool.ID, status, lat, errMsg, testedAt); updateErr != nil {
		m.logger.Warn("proxypools: update health failed", "pool_id", pool.ID, "err", updateErr)
	}

	httpserver.WriteJSON(w, http.StatusOK, PoolTestResponse{
		ID:        pool.ID,
		Healthy:   probeErr == nil,
		Status:    status,
		LatencyMS: lat,
		Error:     errMsg,
		TestedAt:  testedAt,
	})
}

func (m *Module) handleScan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Candidates []string `json:"candidates"`
		AutoAdd    bool     `json:"auto_add"`
		NamePrefix string   `json:"name_prefix"`
		TimeoutMS  int      `json:"timeout_ms"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if len(req.Candidates) == 0 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "candidates must not be empty")
		return
	}
	if len(req.Candidates) > maxScanCandidates {
		httpserver.WriteError(w, r, http.StatusBadRequest, "too_many_candidates", fmt.Sprintf("maximum %d candidates per scan", maxScanCandidates))
		return
	}

	timeout := m.probeTimeout
	if req.TimeoutMS > 0 {
		timeout = min(time.Duration(req.TimeoutMS)*time.Millisecond, 15*time.Second)
	}
	prefix := strings.TrimSpace(req.NamePrefix)
	if prefix == "" {
		prefix = "Scanned Proxy"
	}

	var addedAny bool
	var addMu sync.Mutex
	var onHealthy func(idx int, parsed parsedProxyURL, latencyMS int64) (string, error)
	if req.AutoAdd {
		onHealthy = func(idx int, parsed parsedProxyURL, latencyMS int64) (string, error) {
			enc, encErr := m.encryptURL(parsed.Raw)
			if encErr != nil {
				return "", encErr
			}
			now := formatUTC(time.Now())
			item := storedPool{
				PoolSummary: PoolSummary{
					ID:             newPoolID(),
					Name:           fmt.Sprintf("%s #%d", prefix, idx+1),
					Scheme:         parsed.Scheme,
					DisplayURL:     parsed.DisplayURL,
					HasCredentials: parsed.HasCredentials,
					Enabled:        true,
					Status:         "healthy",
					LastLatencyMS:  latencyMS,
					LastTestedAt:   now,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				URLEncrypted: enc,
			}
			if err := m.repo.createPool(r.Context(), item); err != nil {
				return "", err
			}
			addMu.Lock()
			addedAny = true
			addMu.Unlock()
			return item.ID, nil
		}
	}

	results := runConcurrentScan(r.Context(), req.Candidates, timeout, onHealthy)
	if addedAny {
		m.invalidateAndPublish(r.Context())
	}
	httpserver.WriteJSON(w, http.StatusOK, ScanResponse{Results: results})
}

func (m *Module) handleUpsertBinding(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "providerId"))
	if providerID == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "providerId must not be empty")
		return
	}

	var req struct {
		PoolID string `json:"pool_id"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	poolID := strings.TrimSpace(req.PoolID)
	if poolID == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "pool_id must not be empty")
		return
	}

	if _, err := m.repo.getPool(r.Context(), poolID); errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "proxy pool not found")
		return
	} else if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "proxypools_get_failed", err.Error())
		return
	}

	now := formatUTC(time.Now())
	if err := m.repo.upsertBinding(r.Context(), providerID, poolID, now); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "proxypools_bind_failed", err.Error())
		return
	}
	m.invalidateAndPublish(r.Context())
	httpserver.WriteJSON(w, http.StatusOK, BindingSummary{
		ProviderID: providerID,
		PoolID:     poolID,
		UpdatedAt:  now,
	})
}

func (m *Module) handleDeleteBinding(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "providerId"))
	err := m.repo.deleteBinding(r.Context(), providerID)
	if errors.Is(err, errNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "not_found", "proxy binding not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "proxypools_unbind_failed", err.Error())
		return
	}
	m.invalidateAndPublish(r.Context())
	w.WriteHeader(http.StatusNoContent)
}

func decodeStrictJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("malformed JSON: %w", err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("unexpected trailing JSON content")
	}
	return nil
}

func decodeOptionalJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("malformed JSON: %w", err)
	}
	return nil
}

func newPoolID() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err == nil {
		return "pxp_" + hex.EncodeToString(buf[:])
	}
	return fmt.Sprintf("pxp_%d", time.Now().UnixNano())
}

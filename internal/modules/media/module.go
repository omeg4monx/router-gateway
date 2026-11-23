// Package media implements gateway media endpoints (/v1/audio/*, /v1/images/*,
// /v1/videos/*, /v1/search, /v1/web/fetch), encrypted media provider management,
// ordered media fallback combos, and an SSRF-safe built-in web fetcher.
package media

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

const (
	moduleName         = "media"
	maxAdminBodyBytes  = 1 << 20  // 1 MiB
	maxMediaBodyBytes  = 32 << 20 // 32 MiB
	maxMediaRespBytes  = 64 << 20 // 64 MiB
	maxPriority        = 10000
	defaultPriority    = 100
	TopicMediaChanged  = "media.changed"
	CapabilityAudio    = "audio"
	CapabilityImage    = "image"
	CapabilityVideo    = "video"
	CapabilitySearch   = "search"
	CapabilityFetch    = "fetch"
	DriverOpenAICompat = "openai_compatible"
	DriverGenericHTTP  = "generic_http"
	DriverBuiltinFetch = "builtin_fetch"
)

// Options configures a media Module instance.
type Options struct {
	DB                *db.DB
	Vault             *crypto.Vault
	HTTPClient        *http.Client
	Keys              contracts.KeyValidator
	Usage             contracts.UsageRecorder
	Bus               *eventbus.Bus
	Logger            *slog.Logger
	AllowPrivateFetch bool // test-only override for loopback httptest servers
}

// Module implements lifecycle.Module for media gateway proxying and admin CRUD.
type Module struct {
	repo              *repo
	vault             *crypto.Vault
	httpClient        *http.Client
	keys              contracts.KeyValidator
	usage             contracts.UsageRecorder
	bus               *eventbus.Bus
	logger            *slog.Logger
	allowPrivateFetch bool
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a media Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 120 * time.Second}
	}
	return &Module{
		repo:              newRepo(opts.DB),
		vault:             opts.Vault,
		httpClient:        client,
		keys:              opts.Keys,
		usage:             opts.Usage,
		bus:               opts.Bus,
		logger:            logger,
		allowPrivateFetch: opts.AllowPrivateFetch,
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares modules required before media initializes.
func (m *Module) Dependencies() []string { return []string{"system", "keys"} }

// Migrations returns the embedded SQLite migrations for media.
func (m *Module) Migrations() fs.FS { return migrationsFS }

// Start validates runtime readiness.
func (m *Module) Start(context.Context) error { return nil }

// Stop completes graceful shutdown.
func (m *Module) Stop(context.Context) error { return nil }

// RegisterRoutes mounts AdminAPI management routes and Gateway media routes.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Get("/media/providers", m.handleListProviders)
		mounts.AdminAPI.Post("/media/providers", m.handleCreateProvider)
		mounts.AdminAPI.Patch("/media/providers/{id}", m.handlePatchProvider)
		mounts.AdminAPI.Delete("/media/providers/{id}", m.handleDeleteProvider)

		mounts.AdminAPI.Get("/media/combos", m.handleListCombos)
		mounts.AdminAPI.Post("/media/combos", m.handleCreateCombo)
		mounts.AdminAPI.Patch("/media/combos/{id}", m.handlePatchCombo)
		mounts.AdminAPI.Delete("/media/combos/{id}", m.handleDeleteCombo)
	}
	if mounts.Gateway != nil {
		mounts.Gateway.Post("/v1/audio/transcriptions", m.proxyMediaHandler(CapabilityAudio, "/v1/audio/transcriptions"))
		mounts.Gateway.Post("/v1/audio/translations", m.proxyMediaHandler(CapabilityAudio, "/v1/audio/translations"))
		mounts.Gateway.Post("/v1/audio/speech", m.proxyMediaHandler(CapabilityAudio, "/v1/audio/speech"))

		mounts.Gateway.Post("/v1/images/generations", m.proxyMediaHandler(CapabilityImage, "/v1/images/generations"))
		mounts.Gateway.Post("/v1/images/edits", m.proxyMediaHandler(CapabilityImage, "/v1/images/edits"))
		mounts.Gateway.Post("/v1/images/variations", m.proxyMediaHandler(CapabilityImage, "/v1/images/variations"))

		mounts.Gateway.Post("/v1/videos/generations", m.proxyMediaHandler(CapabilityVideo, "/v1/videos/generations"))
		mounts.Gateway.Get("/v1/videos/{id}", m.handleVideoStatus)

		mounts.Gateway.Post("/v1/search", m.proxyMediaHandler(CapabilitySearch, "/search"))
		mounts.Gateway.Post("/v1/web/fetch", m.handleWebFetch)
	}
}

func isValidCapability(capName string) bool {
	switch capName {
	case CapabilityAudio, CapabilityImage, CapabilityVideo, CapabilitySearch, CapabilityFetch:
		return true
	default:
		return false
	}
}

func isValidDriver(driver string) bool {
	switch driver {
	case DriverOpenAICompat, DriverGenericHTTP, DriverBuiltinFetch:
		return true
	default:
		return false
	}
}

func validateBaseURL(driver, raw string) (string, error) {
	trimmed := strings.TrimRight(strings.TrimSpace(raw), "/")
	if driver == DriverBuiltinFetch {
		return trimmed, nil
	}
	if trimmed == "" {
		return "", errors.New("base_url is required for remote media drivers")
	}
	u, err := url.Parse(trimmed)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || strings.TrimSpace(u.Host) == "" {
		return "", errors.New("base_url must be a valid http:// or https:// URL")
	}
	return trimmed, nil
}

type createProviderReq struct {
	Name         string            `json:"name"`
	Capability   string            `json:"capability"`
	Driver       string            `json:"driver"`
	BaseURL      string            `json:"base_url"`
	APIKey       string            `json:"api_key"`
	DefaultModel string            `json:"default_model"`
	ExtraHeaders map[string]string `json:"extra_headers"`
	Priority     *int              `json:"priority"`
	Enabled      *bool             `json:"enabled"`
}

type updateProviderReq struct {
	Name         *string           `json:"name"`
	Capability   *string           `json:"capability"`
	Driver       *string           `json:"driver"`
	BaseURL      *string           `json:"base_url"`
	APIKey       *string           `json:"api_key"`
	ClearAPIKey  bool              `json:"clear_api_key"`
	DefaultModel *string           `json:"default_model"`
	ExtraHeaders map[string]string `json:"extra_headers"`
	Priority     *int              `json:"priority"`
	Enabled      *bool             `json:"enabled"`
}

func (m *Module) handleListProviders(w http.ResponseWriter, r *http.Request) {
	recs, err := m.repo.listProviders(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "media_list_failed", "failed to list media providers")
		return
	}
	out := make([]ProviderSummary, 0, len(recs))
	for _, rec := range recs {
		out = append(out, rec.ProviderSummary)
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"providers": out})
}

func (m *Module) handleCreateProvider(w http.ResponseWriter, r *http.Request) {
	var req createProviderReq
	if err := decodeStrictJSON(r, maxAdminBodyBytes, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 128 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_name", "name must be 1..128 characters")
		return
	}
	if !isValidCapability(req.Capability) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_capability", "capability must be one of audio, image, video, search, fetch")
		return
	}
	if !isValidDriver(req.Driver) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_driver", "driver must be one of openai_compatible, generic_http, builtin_fetch")
		return
	}
	if req.Driver == DriverBuiltinFetch && req.Capability != CapabilityFetch {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_driver", "builtin_fetch driver is only valid for fetch capability")
		return
	}
	baseURL, err := validateBaseURL(req.Driver, req.BaseURL)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_base_url", err.Error())
		return
	}
	priority := defaultPriority
	if req.Priority != nil {
		if *req.Priority < 0 || *req.Priority > maxPriority {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_priority", "priority must be between 0 and 10000")
			return
		}
		priority = *req.Priority
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	var encKey string
	if trimmedKey := strings.TrimSpace(req.APIKey); trimmedKey != "" {
		if m.vault == nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_unavailable", "secret vault is not configured")
			return
		}
		encKey, err = m.vault.Encrypt([]byte(trimmedKey))
		if err != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_encrypt_failed", "failed to encrypt api key")
			return
		}
	}

	rec, err := m.repo.createProvider(r.Context(), providerRecord{
		ProviderSummary: ProviderSummary{
			Name:         name,
			Capability:   req.Capability,
			Driver:       req.Driver,
			BaseURL:      baseURL,
			DefaultModel: strings.TrimSpace(req.DefaultModel),
			ExtraHeaders: sanitizeHeaders(req.ExtraHeaders),
			Priority:     priority,
			Enabled:      enabled,
		},
		APIKeyEnc: encKey,
	})
	if err != nil {
		if errors.Is(err, errAlreadyExists) {
			httpserver.WriteError(w, r, http.StatusConflict, "media_provider_exists", "a media provider with that name already exists")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "media_create_failed", "failed to create media provider")
		return
	}
	m.publishChange("provider.created", rec.ID)
	httpserver.WriteJSON(w, http.StatusCreated, rec.ProviderSummary)
}

func (m *Module) handlePatchProvider(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	existing, err := m.repo.getProvider(r.Context(), id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "media_provider_not_found", "media provider not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "media_lookup_failed", "failed to load media provider")
		return
	}

	var req updateProviderReq
	if err := decodeStrictJSON(r, maxAdminBodyBytes, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" || len(name) > 128 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_name", "name must be 1..128 characters")
			return
		}
		existing.Name = name
	}
	if req.Capability != nil {
		if !isValidCapability(*req.Capability) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_capability", "invalid media capability")
			return
		}
		existing.Capability = *req.Capability
	}
	if req.Driver != nil {
		if !isValidDriver(*req.Driver) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_driver", "invalid media driver")
			return
		}
		existing.Driver = *req.Driver
	}
	if existing.Driver == DriverBuiltinFetch && existing.Capability != CapabilityFetch {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_driver", "builtin_fetch driver is only valid for fetch capability")
		return
	}
	if req.BaseURL != nil {
		existing.BaseURL = *req.BaseURL
	}
	validatedURL, err := validateBaseURL(existing.Driver, existing.BaseURL)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_base_url", err.Error())
		return
	}
	existing.BaseURL = validatedURL

	if req.DefaultModel != nil {
		existing.DefaultModel = strings.TrimSpace(*req.DefaultModel)
	}
	if req.ExtraHeaders != nil {
		existing.ExtraHeaders = sanitizeHeaders(req.ExtraHeaders)
	}
	if req.Priority != nil {
		if *req.Priority < 0 || *req.Priority > maxPriority {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_priority", "priority must be between 0 and 10000")
			return
		}
		existing.Priority = *req.Priority
	}
	if req.Enabled != nil {
		existing.Enabled = *req.Enabled
	}
	if req.ClearAPIKey {
		existing.APIKeyEnc = ""
	} else if req.APIKey != nil && strings.TrimSpace(*req.APIKey) != "" {
		if m.vault == nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_unavailable", "secret vault is not configured")
			return
		}
		enc, encErr := m.vault.Encrypt([]byte(strings.TrimSpace(*req.APIKey)))
		if encErr != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_encrypt_failed", "failed to encrypt api key")
			return
		}
		existing.APIKeyEnc = enc
	}

	updated, err := m.repo.updateProvider(r.Context(), existing)
	if err != nil {
		if errors.Is(err, errAlreadyExists) {
			httpserver.WriteError(w, r, http.StatusConflict, "media_provider_exists", "a media provider with that name already exists")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "media_update_failed", "failed to update media provider")
		return
	}
	m.publishChange("provider.updated", updated.ID)
	httpserver.WriteJSON(w, http.StatusOK, updated.ProviderSummary)
}

func (m *Module) handleDeleteProvider(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if err := m.repo.deleteProvider(r.Context(), id); err != nil {
		if errors.Is(err, errNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "media_provider_not_found", "media provider not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "media_delete_failed", "failed to delete media provider")
		return
	}
	m.publishChange("provider.deleted", id)
	httpserver.WriteJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

type createComboReq struct {
	Name       string   `json:"name"`
	Capability string   `json:"capability"`
	Targets    []string `json:"targets"`
	Enabled    *bool    `json:"enabled"`
}

type updateComboReq struct {
	Name       *string   `json:"name"`
	Capability *string   `json:"capability"`
	Targets    *[]string `json:"targets"`
	Enabled    *bool     `json:"enabled"`
}

func normalizeTargets(raw []string) ([]string, error) {
	seen := make(map[string]struct{}, len(raw))
	out := make([]string, 0, len(raw))
	for _, t := range raw {
		trimmed := strings.TrimSpace(t)
		if trimmed == "" {
			continue
		}
		if _, dup := seen[trimmed]; dup {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	if len(out) == 0 {
		return nil, errors.New("targets must contain at least one provider id or name")
	}
	return out, nil
}

func (m *Module) handleListCombos(w http.ResponseWriter, r *http.Request) {
	combos, err := m.repo.listCombos(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "media_combo_list_failed", "failed to list media combos")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"combos": combos})
}

func (m *Module) handleCreateCombo(w http.ResponseWriter, r *http.Request) {
	var req createComboReq
	if err := decodeStrictJSON(r, maxAdminBodyBytes, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 128 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_name", "name must be 1..128 characters")
		return
	}
	if !isValidCapability(req.Capability) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_capability", "invalid media capability")
		return
	}
	targets, err := normalizeTargets(req.Targets)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_targets", err.Error())
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	created, err := m.repo.createCombo(r.Context(), Combo{
		Name:       name,
		Capability: req.Capability,
		Targets:    targets,
		Enabled:    enabled,
	})
	if err != nil {
		if errors.Is(err, errAlreadyExists) {
			httpserver.WriteError(w, r, http.StatusConflict, "media_combo_exists", "a media combo with that name already exists")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "media_combo_create_failed", "failed to create media combo")
		return
	}
	m.publishChange("combo.created", created.ID)
	httpserver.WriteJSON(w, http.StatusCreated, created)
}

func (m *Module) handlePatchCombo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	existing, err := m.repo.getCombo(r.Context(), id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "media_combo_not_found", "media combo not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "media_combo_lookup_failed", "failed to load media combo")
		return
	}
	var req updateComboReq
	if err := decodeStrictJSON(r, maxAdminBodyBytes, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" || len(name) > 128 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_name", "name must be 1..128 characters")
			return
		}
		existing.Name = name
	}
	if req.Capability != nil {
		if !isValidCapability(*req.Capability) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_capability", "invalid media capability")
			return
		}
		existing.Capability = *req.Capability
	}
	if req.Targets != nil {
		targets, err := normalizeTargets(*req.Targets)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_targets", err.Error())
			return
		}
		existing.Targets = targets
	}
	if req.Enabled != nil {
		existing.Enabled = *req.Enabled
	}
	updated, err := m.repo.updateCombo(r.Context(), existing)
	if err != nil {
		if errors.Is(err, errAlreadyExists) {
			httpserver.WriteError(w, r, http.StatusConflict, "media_combo_exists", "a media combo with that name already exists")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "media_combo_update_failed", "failed to update media combo")
		return
	}
	m.publishChange("combo.updated", updated.ID)
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteCombo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if err := m.repo.deleteCombo(r.Context(), id); err != nil {
		if errors.Is(err, errNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "media_combo_not_found", "media combo not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "media_combo_delete_failed", "failed to delete media combo")
		return
	}
	m.publishChange("combo.deleted", id)
	httpserver.WriteJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

func (m *Module) authenticateGateway(w http.ResponseWriter, r *http.Request) (contracts.KeyInfo, bool) {
	rawKey := extractGatewayAPIKey(r)
	if rawKey == "" {
		writeOpenAIError(w, http.StatusUnauthorized, "authentication_error", "missing_api_key", "missing API key")
		return contracts.KeyInfo{}, false
	}
	if m.keys == nil {
		writeOpenAIError(w, http.StatusInternalServerError, "server_error", "auth_unavailable", "key validator unavailable")
		return contracts.KeyInfo{}, false
	}
	info, err := m.keys.ValidateKey(r.Context(), rawKey)
	if err != nil || info.Disabled {
		writeOpenAIError(w, http.StatusUnauthorized, "authentication_error", "invalid_api_key", "invalid or disabled API key")
		return contracts.KeyInfo{}, false
	}
	return info, true
}

func extractGatewayAPIKey(r *http.Request) string {
	if auth := strings.TrimSpace(r.Header.Get("Authorization")); len(auth) > 7 && strings.EqualFold(auth[:7], "Bearer ") {
		if tok := strings.TrimSpace(auth[7:]); tok != "" {
			return tok
		}
	}
	if key := strings.TrimSpace(r.Header.Get("x-api-key")); key != "" {
		return key
	}
	if key := strings.TrimSpace(r.Header.Get("x-goog-api-key")); key != "" {
		return key
	}
	if r.URL != nil {
		if key := strings.TrimSpace(r.URL.Query().Get("key")); key != "" {
			return key
		}
	}
	return ""
}

func writeOpenAIError(w http.ResponseWriter, status int, errType, code, msg string) {
	httpserver.WriteJSON(w, status, map[string]any{
		"error": map[string]any{
			"message": msg,
			"type":    errType,
			"code":    code,
		},
	})
}

func (m *Module) resolveCandidates(ctx context.Context, capability, selector string) ([]providerRecord, error) {
	providers, err := m.repo.listProviders(ctx)
	if err != nil {
		return nil, err
	}
	byID := make(map[string]providerRecord, len(providers))
	byName := make(map[string]providerRecord, len(providers))
	capEnabled := make([]providerRecord, 0)
	for _, p := range providers {
		if !p.Enabled || p.Capability != capability {
			continue
		}
		byID[p.ID] = p
		byName[strings.ToLower(p.Name)] = p
		capEnabled = append(capEnabled, p)
	}

	sel := strings.TrimSpace(selector)
	if sel != "" {
		combos, err := m.repo.listCombos(ctx)
		if err != nil {
			return nil, err
		}
		for _, c := range combos {
			if !c.Enabled || c.Capability != capability {
				continue
			}
			if strings.EqualFold(c.Name, sel) || c.ID == sel {
				resolved := make([]providerRecord, 0, len(c.Targets))
				for _, t := range c.Targets {
					if p, ok := byID[t]; ok {
						resolved = append(resolved, p)
					} else if p, ok := byName[strings.ToLower(t)]; ok {
						resolved = append(resolved, p)
					}
				}
				return resolved, nil
			}
		}
		if p, ok := byID[sel]; ok {
			return []providerRecord{p}, nil
		}
		if p, ok := byName[strings.ToLower(sel)]; ok {
			return []providerRecord{p}, nil
		}
	}
	return capEnabled, nil
}

func (m *Module) proxyMediaHandler(capability, upstreamPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.handleMediaProxy(w, r, capability, upstreamPath)
	}
}

func (m *Module) handleVideoStatus(w http.ResponseWriter, r *http.Request) {
	videoID := strings.TrimSpace(chi.URLParam(r, "id"))
	if videoID == "" {
		writeOpenAIError(w, http.StatusBadRequest, "invalid_request_error", "missing_video_id", "video id is required")
		return
	}
	m.handleMediaProxy(w, r, CapabilityVideo, "/v1/videos/"+url.PathEscape(videoID))
}

func (m *Module) handleMediaProxy(w http.ResponseWriter, r *http.Request, capability, upstreamPath string) {
	keyInfo, ok := m.authenticateGateway(w, r)
	if !ok {
		return
	}

	start := time.Now()
	var bodyBytes []byte
	if r.Body != nil && r.Method != http.MethodGet {
		limited := io.LimitReader(r.Body, maxMediaBodyBytes+1)
		raw, err := io.ReadAll(limited)
		if err != nil {
			writeOpenAIError(w, http.StatusBadRequest, "invalid_request_error", "read_body_failed", "failed to read request body")
			return
		}
		if len(raw) > maxMediaBodyBytes {
			writeOpenAIError(w, http.StatusRequestEntityTooLarge, "invalid_request_error", "payload_too_large", "request payload exceeds 32 MiB limit")
			return
		}
		bodyBytes = raw
	}

	selector := extractMediaSelector(r, bodyBytes)
	candidates, err := m.resolveCandidates(r.Context(), capability, selector)
	if err != nil {
		writeOpenAIError(w, http.StatusInternalServerError, "server_error", "media_resolve_failed", "failed to resolve media provider")
		return
	}
	if len(candidates) == 0 {
		writeOpenAIError(w, http.StatusNotFound, "invalid_request_error", "media_provider_not_found", fmt.Sprintf("no enabled %s media provider configured", capability))
		return
	}

	var (
		lastStatus int
		lastBody   []byte
		lastCT     string
		lastErr    error
	)

	for _, cand := range candidates {
		targetURL := buildUpstreamMediaURL(cand, upstreamPath)
		var reqBody io.Reader
		if len(bodyBytes) > 0 {
			reqBody = bytes.NewReader(bodyBytes)
		}
		upReq, reqErr := http.NewRequestWithContext(r.Context(), r.Method, targetURL, reqBody)
		if reqErr != nil {
			lastErr = reqErr
			continue
		}
		if ct := r.Header.Get("Content-Type"); ct != "" {
			upReq.Header.Set("Content-Type", ct)
		}
		if accept := r.Header.Get("Accept"); accept != "" {
			upReq.Header.Set("Accept", accept)
		}
		for hk, hv := range cand.ExtraHeaders {
			upReq.Header.Set(hk, hv)
		}
		if cand.APIKeyEnc != "" && m.vault != nil {
			plain, decErr := m.vault.Decrypt(cand.APIKeyEnc)
			if decErr != nil {
				m.logger.Error("media: decrypt provider api key failed", "provider_id", cand.ID, "error", decErr)
				lastErr = decErr
				continue
			}
			upReq.Header.Set("Authorization", "Bearer "+string(plain))
		}

		//nolint:gosec // G704: upstream URL is admin-configured and scheme-validated on create/patch.
		resp, doErr := m.httpClient.Do(upReq)
		if doErr != nil {
			lastErr = doErr
			continue
		}
		respBody, readErr := io.ReadAll(io.LimitReader(resp.Body, maxMediaRespBytes))
		_ = resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}

		lastStatus = resp.StatusCode
		lastBody = respBody
		lastCT = resp.Header.Get("Content-Type")

		// Failover only on provider/account errors (401/403/429/5xx).
		// Client errors (400/404/422 etc.) return immediately without burning fallback providers.
		if isRetryableUpstreamStatus(resp.StatusCode) {
			continue
		}

		m.recordUsage(r.Context(), keyInfo.ID, selector, cand, capability, resp.StatusCode, time.Since(start), "")
		if lastCT != "" {
			w.Header().Set("Content-Type", lastCT)
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(respBody)
		return
	}

	if lastStatus > 0 {
		m.recordUsage(r.Context(), keyInfo.ID, selector, candidates[len(candidates)-1], capability, lastStatus, time.Since(start), "upstream media error")
		if lastCT != "" {
			w.Header().Set("Content-Type", lastCT)
		}
		w.WriteHeader(lastStatus)
		_, _ = w.Write(lastBody)
		return
	}

	errSummary := "all media providers failed"
	if lastErr != nil {
		errSummary = lastErr.Error()
	}
	m.recordUsage(r.Context(), keyInfo.ID, selector, candidates[0], capability, http.StatusBadGateway, time.Since(start), errSummary)
	writeOpenAIError(w, http.StatusBadGateway, "server_error", "media_upstream_failed", errSummary)
}

func (m *Module) handleWebFetch(w http.ResponseWriter, r *http.Request) {
	keyInfo, ok := m.authenticateGateway(w, r)
	if !ok {
		return
	}
	start := time.Now()
	var req WebFetchRequest
	if err := decodeStrictJSON(r, maxAdminBodyBytes, &req); err != nil {
		writeOpenAIError(w, http.StatusBadRequest, "invalid_request_error", "invalid_payload", err.Error())
		return
	}

	// If a remote fetch provider is explicitly selected or configured (non-builtin), proxy to it.
	if candidates, err := m.resolveCandidates(r.Context(), CapabilityFetch, req.Provider); err == nil {
		for _, c := range candidates {
			if c.Driver != DriverBuiltinFetch && c.BaseURL != "" {
				raw, _ := json.Marshal(req)
				r.Body = io.NopCloser(bytes.NewReader(raw))
				m.handleMediaProxy(w, r, CapabilityFetch, "/web/fetch")
				return
			}
		}
	}

	res, err := m.executeBuiltinFetch(r.Context(), req)
	if err != nil {
		status := http.StatusBadGateway
		code := "web_fetch_failed"
		if errors.Is(err, errSSRFBlocked) {
			status = http.StatusForbidden
			code = "ssrf_blocked"
		} else if errors.Is(err, errInvalidFetchURL) {
			status = http.StatusBadRequest
			code = "invalid_url"
		}
		m.recordUsage(r.Context(), keyInfo.ID, "builtin_fetch", providerRecord{ProviderSummary: ProviderSummary{ID: "builtin_fetch", Name: "builtin_fetch"}}, CapabilityFetch, status, time.Since(start), err.Error())
		writeOpenAIError(w, status, "invalid_request_error", code, err.Error())
		return
	}

	m.recordUsage(r.Context(), keyInfo.ID, "builtin_fetch", providerRecord{ProviderSummary: ProviderSummary{ID: "builtin_fetch", Name: "builtin_fetch"}}, CapabilityFetch, http.StatusOK, time.Since(start), "")
	httpserver.WriteJSON(w, http.StatusOK, res)
}

func isRetryableUpstreamStatus(status int) bool {
	return status == http.StatusUnauthorized ||
		status == http.StatusForbidden ||
		status == http.StatusTooManyRequests ||
		status >= 500
}

func buildUpstreamMediaURL(cand providerRecord, upstreamPath string) string {
	base := strings.TrimRight(cand.BaseURL, "/")
	if cand.Driver == DriverGenericHTTP {
		if u, err := url.Parse(base); err == nil && u.Path != "" && u.Path != "/" {
			return base
		}
	}
	if strings.HasSuffix(base, "/v1") && strings.HasPrefix(upstreamPath, "/v1/") {
		return base + strings.TrimPrefix(upstreamPath, "/v1")
	}
	return base + upstreamPath
}

func extractMediaSelector(r *http.Request, body []byte) string {
	if r.URL != nil {
		for _, q := range []string{"combo", "provider", "model"} {
			if v := strings.TrimSpace(r.URL.Query().Get(q)); v != "" {
				return v
			}
		}
	}
	if len(body) == 0 {
		return ""
	}
	ct := r.Header.Get("Content-Type")
	mediaType, params, _ := mime.ParseMediaType(ct)
	if mediaType == "application/json" || (mediaType == "" && bytes.HasPrefix(bytes.TrimSpace(body), []byte("{"))) {
		var payload struct {
			Model    string `json:"model"`
			Combo    string `json:"combo"`
			Provider string `json:"provider"`
		}
		if json.Unmarshal(body, &payload) == nil {
			if v := strings.TrimSpace(payload.Combo); v != "" {
				return v
			}
			if v := strings.TrimSpace(payload.Provider); v != "" {
				return v
			}
			if v := strings.TrimSpace(payload.Model); v != "" {
				return v
			}
		}
	}
	if strings.HasPrefix(mediaType, "multipart/") && params["boundary"] != "" {
		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		for {
			part, err := mr.NextPart()
			if err != nil {
				break
			}
			name := part.FormName()
			if name == "model" || name == "combo" || name == "provider" {
				val, _ := io.ReadAll(io.LimitReader(part, 256))
				_ = part.Close()
				if s := strings.TrimSpace(string(val)); s != "" {
					return s
				}
				continue
			}
			_ = part.Close()
		}
	}
	return ""
}

func (m *Module) recordUsage(ctx context.Context, keyID, model string, cand providerRecord, capability string, status int, latency time.Duration, errSummary string) {
	if m.usage == nil {
		return
	}
	if model == "" {
		model = cand.DefaultModel
	}
	if model == "" {
		model = "media/" + capability
	}
	m.usage.Record(ctx, contracts.UsageRecord{
		RequestID:    newID("mreq"),
		KeyID:        keyID,
		Model:        model,
		ProviderID:   cand.Name,
		ConnectionID: cand.ID,
		Protocol:     "media-" + capability,
		Latency:      latency,
		StatusCode:   status,
		ErrorSummary: errSummary,
		Timestamp:    time.Now().UTC(),
	})
}

func (m *Module) publishChange(action, id string) {
	if m.bus == nil {
		return
	}
	m.bus.Publish(TopicMediaChanged, map[string]any{
		"action": action,
		"id":     id,
	})
}

func sanitizeHeaders(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		tk := strings.TrimSpace(k)
		if tk == "" {
			continue
		}
		out[tk] = strings.TrimSpace(v)
	}
	return out
}

func decodeStrictJSON(r *http.Request, limit int64, dst any) error {
	if r.Body == nil {
		return errors.New("request body is required")
	}
	dec := json.NewDecoder(io.LimitReader(r.Body, limit))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}

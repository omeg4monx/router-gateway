// Package vouchers implements redeemable voucher lifecycle management, atomic
// claim provisioning of gateway API keys, optional team assignment, and
// anti-abuse per-IP / per-actor enforcement.
package vouchers

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"math"
	"net"
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

//go:embed migrations/*.sql
var migrationsFS embed.FS

const (
	moduleName            = "vouchers"
	maxBodyBytes          = 1 << 20
	defaultMaxClaimsPerIP = 10
	rateLimitWindow       = time.Minute
	maxTrackedIPs         = 2048

	// TopicVouchersChanged is published when vouchers are mutated or redeemed.
	TopicVouchersChanged = "vouchers.changed"
)

// Options configures a vouchers Module instance.
type Options struct {
	DB             *db.DB
	Bus            *eventbus.Bus
	Logger         *slog.Logger
	KeyProvisioner contracts.KeyProvisioner
	TeamAssigner   contracts.TeamAssigner
	MaxClaimsPerIP int
}

// Module implements lifecycle.Module for voucher administration and claims.
type Module struct {
	repo           *repo
	bus            *eventbus.Bus
	logger         *slog.Logger
	keys           contracts.KeyProvisioner
	teams          contracts.TeamAssigner
	maxClaimsPerIP int

	ipMu      sync.Mutex
	ipWindows map[string][]time.Time
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a vouchers Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	maxIP := opts.MaxClaimsPerIP
	if maxIP <= 0 {
		maxIP = defaultMaxClaimsPerIP
	}
	return &Module{
		repo:           newRepo(opts.DB),
		bus:            opts.Bus,
		logger:         logger,
		keys:           opts.KeyProvisioner,
		teams:          opts.TeamAssigner,
		maxClaimsPerIP: maxIP,
		ipWindows:      make(map[string][]time.Time),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares modules required before vouchers initializes.
func (m *Module) Dependencies() []string { return []string{"system", "keys"} }

// Migrations returns the embedded SQLite migrations for vouchers.
func (m *Module) Migrations() fs.FS { return migrationsFS }

// RegisterRoutes mounts admin voucher CRUD on AdminAPI and self-service claim on PublicAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Get("/vouchers", m.handleList)
		mounts.AdminAPI.Post("/vouchers", m.handleCreate)
		mounts.AdminAPI.Patch("/vouchers/{id}", m.handlePatch)
		mounts.AdminAPI.Delete("/vouchers/{id}", m.handleDelete)
		mounts.AdminAPI.Get("/vouchers/{id}/redemptions", m.handleListRedemptions)
	}
	if mounts.PublicAPI != nil {
		mounts.PublicAPI.Post("/vouchers/claim", m.handleClaim)
	} else if mounts.AdminAPI != nil {
		mounts.AdminAPI.Post("/vouchers/claim", m.handleClaim)
	}
}

// Start is a no-op for vouchers.
func (m *Module) Start(context.Context) error { return nil }

// Stop is a no-op for vouchers.
func (m *Module) Stop(context.Context) error { return nil }

func (m *Module) allowIP(ip string) bool {
	now := time.Now()
	cutoff := now.Add(-rateLimitWindow)

	m.ipMu.Lock()
	defer m.ipMu.Unlock()

	if len(m.ipWindows) > maxTrackedIPs {
		for k, times := range m.ipWindows {
			if len(times) == 0 || !times[len(times)-1].After(cutoff) {
				delete(m.ipWindows, k)
			}
		}
		for len(m.ipWindows) > maxTrackedIPs {
			var oldestKey string
			var oldestTime time.Time
			for k, times := range m.ipWindows {
				last := time.Time{}
				if len(times) > 0 {
					last = times[len(times)-1]
				}
				if oldestKey == "" || last.Before(oldestTime) {
					oldestKey = k
					oldestTime = last
				}
			}
			if oldestKey == "" {
				break
			}
			delete(m.ipWindows, oldestKey)
		}
	}

	window := m.ipWindows[ip]
	pruned := window[:0]
	for _, ts := range window {
		if ts.After(cutoff) {
			pruned = append(pruned, ts)
		}
	}
	if len(pruned) >= m.maxClaimsPerIP {
		m.ipWindows[ip] = pruned
		return false
	}
	m.ipWindows[ip] = append(pruned, now)
	return true
}

func (m *Module) publishChanged(action, voucherID string) {
	if m.bus != nil {
		m.bus.Publish(TopicVouchersChanged, map[string]any{
			"action":     action,
			"voucher_id": voucherID,
		})
	}
}

func (m *Module) handleList(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.list(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "vouchers_list_failed", "failed to list vouchers")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"vouchers": items})
}

func (m *Module) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code           string   `json:"code"`
		Name           string   `json:"name"`
		Description    string   `json:"description"`
		QuotaUSD       float64  `json:"quota_usd"`
		RateLimitRPM   int      `json:"rate_limit_rpm"`
		AllowedModels  []string `json:"allowed_models"`
		TeamID         string   `json:"team_id"`
		MaxRedemptions *int     `json:"max_redemptions"`
		Enabled        *bool    `json:"enabled"`
		ExpiresAt      string   `json:"expires_at"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || utf8.RuneCountInString(req.Name) > 128 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_name", "name must be 1..128 characters")
		return
	}
	if utf8.RuneCountInString(strings.TrimSpace(req.Code)) > 64 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_code", "code must be at most 64 characters")
		return
	}
	if req.QuotaUSD < 0 || math.IsNaN(req.QuotaUSD) || math.IsInf(req.QuotaUSD, 0) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_quota", "quota_usd must be a non-negative finite number")
		return
	}
	if req.RateLimitRPM < 0 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_rate_limit", "rate_limit_rpm must be non-negative")
		return
	}
	if len(req.AllowedModels) > 256 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_allowed_models", "allowed_models must have at most 256 entries")
		return
	}
	for _, model := range req.AllowedModels {
		trimmed := strings.TrimSpace(model)
		if trimmed == "" || utf8.RuneCountInString(trimmed) > 128 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_allowed_models", "each allowed_models entry must be 1..128 characters")
			return
		}
	}
	maxRed := 1
	if req.MaxRedemptions != nil {
		if *req.MaxRedemptions <= 0 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_max_redemptions", "max_redemptions must be >= 1")
			return
		}
		maxRed = *req.MaxRedemptions
	}
	if req.ExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(req.ExpiresAt))
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_expires_at", "expires_at must be valid RFC3339 timestamp")
			return
		}
		req.ExpiresAt = parsed.UTC().Format(time.RFC3339)
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	created, err := m.repo.create(r.Context(), createVoucherInput{
		Code:           req.Code,
		Name:           req.Name,
		Description:    req.Description,
		QuotaUSD:       req.QuotaUSD,
		RateLimitRPM:   req.RateLimitRPM,
		AllowedModels:  req.AllowedModels,
		TeamID:         req.TeamID,
		MaxRedemptions: maxRed,
		Enabled:        enabled,
		ExpiresAt:      req.ExpiresAt,
	})
	if err != nil {
		if errors.Is(err, errVoucherCodeExists) {
			httpserver.WriteError(w, r, http.StatusConflict, "voucher_code_exists", "voucher code already exists")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "voucher_create_failed", "failed to create voucher")
		return
	}
	m.publishChanged("created", created.ID)
	httpserver.WriteJSON(w, http.StatusCreated, created)
}

func (m *Module) handlePatch(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	var req struct {
		Name           *string `json:"name"`
		Description    *string `json:"description"`
		MaxRedemptions *int    `json:"max_redemptions"`
		Enabled        *bool   `json:"enabled"`
		ExpiresAt      *string `json:"expires_at"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" || utf8.RuneCountInString(trimmed) > 128 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_name", "name must be 1..128 characters")
			return
		}
	}
	if req.MaxRedemptions != nil && *req.MaxRedemptions <= 0 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_max_redemptions", "max_redemptions must be >= 1")
		return
	}
	if req.ExpiresAt != nil && strings.TrimSpace(*req.ExpiresAt) != "" {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.ExpiresAt))
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_expires_at", "expires_at must be valid RFC3339 timestamp")
			return
		}
		norm := parsed.UTC().Format(time.RFC3339)
		req.ExpiresAt = &norm
	}

	updated, err := m.repo.patch(r.Context(), id, patchVoucherInput{
		Name:           req.Name,
		Description:    req.Description,
		MaxRedemptions: req.MaxRedemptions,
		Enabled:        req.Enabled,
		ExpiresAt:      req.ExpiresAt,
	})
	if err != nil {
		if errors.Is(err, errVoucherNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "voucher_not_found", "voucher not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "voucher_update_failed", "failed to update voucher")
		return
	}
	m.publishChanged("updated", updated.ID)
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

func (m *Module) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if err := m.repo.deleteByID(r.Context(), id); err != nil {
		if errors.Is(err, errVoucherNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "voucher_not_found", "voucher not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "voucher_delete_failed", "failed to delete voucher")
		return
	}
	m.publishChanged("deleted", id)
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": true})
}

func (m *Module) handleListRedemptions(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	items, err := m.repo.listRedemptions(r.Context(), id)
	if err != nil {
		if errors.Is(err, errVoucherNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "voucher_not_found", "voucher not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "redemptions_list_failed", "failed to list voucher redemptions")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"redemptions": items})
}

func (m *Module) handleClaim(w http.ResponseWriter, r *http.Request) {
	ip := clientIP(r)
	if !m.allowIP(ip) {
		httpserver.WriteError(w, r, http.StatusTooManyRequests, "claim_rate_limited", "too many voucher claim attempts; please wait and retry")
		return
	}

	var req struct {
		Code      string `json:"code"`
		ClaimerID string `json:"claimer_id"`
		KeyName   string `json:"key_name"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	code := normalizeCode(req.Code)
	if code == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_code", "voucher code is required")
		return
	}
	keyName := strings.TrimSpace(req.KeyName)
	if utf8.RuneCountInString(keyName) > 128 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_key_name", "key_name must be at most 128 characters")
		return
	}
	claimerID := strings.TrimSpace(req.ClaimerID)
	if claimerID == "" {
		claimerID = "ip:" + ip
	}

	if m.keys == nil {
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "key_provisioner_unavailable", "API key provisioner is not configured")
		return
	}

	voucher, red, err := m.repo.reserveClaim(r.Context(), code, claimerID, ip)
	if err != nil {
		switch {
		case errors.Is(err, errVoucherNotFound):
			httpserver.WriteError(w, r, http.StatusNotFound, "voucher_not_found", "voucher code not found")
		case errors.Is(err, errVoucherDisabled):
			httpserver.WriteError(w, r, http.StatusForbidden, "voucher_disabled", "voucher is disabled")
		case errors.Is(err, errVoucherExpired):
			httpserver.WriteError(w, r, http.StatusGone, "voucher_expired", "voucher has expired")
		case errors.Is(err, errVoucherExhausted):
			httpserver.WriteError(w, r, http.StatusConflict, "voucher_exhausted", "voucher has reached its maximum redemptions")
		case errors.Is(err, errAlreadyClaimed):
			httpserver.WriteError(w, r, http.StatusConflict, "already_claimed", "this voucher has already been redeemed by this claimer")
		default:
			httpserver.WriteError(w, r, http.StatusInternalServerError, "voucher_claim_failed", "failed to redeem voucher")
		}
		return
	}

	compensate := func(keyID string) {
		if keyID != "" {
			revCtx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), 5*time.Second)
			if revErr := m.keys.RevokeKey(revCtx, keyID); revErr != nil {
				m.logger.Error("failed to revoke provisioned voucher key during rollback", "key_id", keyID, "error", revErr)
			}
			cancel()
		}
		if rbErr := m.repo.rollbackClaim(r.Context(), voucher.ID, red.ID); rbErr != nil {
			m.logger.Error("failed to rollback voucher claim reservation", "voucher_id", voucher.ID, "redemption_id", red.ID, "error", rbErr)
		}
	}

	if voucher.TeamID != "" && m.teams == nil {
		compensate("")
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "teams_unavailable", "teams module is disabled; cannot redeem team-bound voucher")
		return
	}

	if keyName == "" {
		keyName = truncateRunes(fmt.Sprintf("voucher-%s-%s", strings.ToLower(voucher.Code), claimerID), 120)
	}

	var expiresPtr *string
	if voucher.ExpiresAt != "" {
		exp := voucher.ExpiresAt
		expiresPtr = &exp
	}

	provisioned, err := m.keys.ProvisionKey(r.Context(), contracts.ProvisionKeyRequest{
		Name:          keyName,
		AllowedModels: voucher.AllowedModels,
		Scopes:        []string{"gateway"},
		Tags:          []string{"voucher:" + voucher.Code},
		RateLimitRPM:  voucher.RateLimitRPM,
		QuotaUSD:      voucher.QuotaUSD,
		ExpiresAt:     expiresPtr,
		Actor:         "voucher:" + voucher.Code,
	})
	if err != nil {
		compensate("")
		httpserver.WriteError(w, r, http.StatusInternalServerError, "key_provision_failed", "failed to provision API key for voucher")
		return
	}

	if err := m.repo.attachRedemptionKey(r.Context(), red.ID, provisioned.ID); err != nil {
		compensate(provisioned.ID)
		httpserver.WriteError(w, r, http.StatusInternalServerError, "voucher_claim_failed", "failed to record provisioned key on voucher redemption")
		return
	}
	if voucher.TeamID != "" {
		if assignErr := m.teams.AssignKeyToTeam(r.Context(), voucher.TeamID, provisioned.ID); assignErr != nil {
			compensate(provisioned.ID)
			httpserver.WriteError(w, r, http.StatusConflict, "team_assignment_failed", "failed to attach voucher key to target team")
			return
		}
	}

	m.publishChanged("redeemed", voucher.ID)
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"voucher_id":  voucher.ID,
		"code":        voucher.Code,
		"key_id":      provisioned.ID,
		"key_prefix":  provisioned.KeyPrefix,
		"secret":      provisioned.Secret,
		"team_id":     voucher.TeamID,
		"quota_usd":   voucher.QuotaUSD,
		"redeemed_at": red.RedeemedAt,
	})
}

func truncateRunes(s string, maxRunes int) string {
	if maxRunes <= 0 || utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	var b strings.Builder
	count := 0
	for _, r := range s {
		if count >= maxRunes {
			break
		}
		b.WriteRune(r)
		count++
	}
	return b.String()
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

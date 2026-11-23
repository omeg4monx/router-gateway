// Package teams implements API-key grouping with shared USD budgets and shared
// RPM quotas enforced via GatewayInterceptor under the single-admin model.
package teams

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

//go:embed migrations/*.sql
var migrationsFS embed.FS

const (
	moduleName          = "teams"
	interceptorPriority = 5
	maxBodyBytes        = 1 << 20
	usageSubBufSize     = 128

	// TopicTeamsChanged is published whenever team definitions or memberships change.
	TopicTeamsChanged = "teams.changed"
)

type statusError struct {
	code int
	msg  string
}

func (e statusError) Error() string   { return e.msg }
func (e statusError) StatusCode() int { return e.code }

// Options configures a teams Module instance.
type Options struct {
	DB     *db.DB
	Bus    *eventbus.Bus
	Logger *slog.Logger
}

// Module implements lifecycle.Module, contracts.GatewayInterceptor, and
// contracts.TeamAssigner.
type Module struct {
	repo   *repo
	bus    *eventbus.Bus
	logger *slog.Logger

	mu         sync.RWMutex
	teamsByID  map[string]Team
	keyToTeam  map[string]string
	rpmWindows map[string][]time.Time

	lifeMu  sync.Mutex
	started bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

var (
	_ lifecycle.Module             = (*Module)(nil)
	_ contracts.GatewayInterceptor = (*Module)(nil)
	_ contracts.TeamAssigner       = (*Module)(nil)
)

// New constructs a teams Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Module{
		repo:       newRepo(opts.DB),
		bus:        opts.Bus,
		logger:     logger,
		teamsByID:  make(map[string]Team),
		keyToTeam:  make(map[string]string),
		rpmWindows: make(map[string][]time.Time),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares modules required before teams initializes.
func (m *Module) Dependencies() []string { return []string{"system", "keys"} }

// Priority runs team budget/RPM checks first in the interceptor chain.
func (m *Module) Priority() int { return interceptorPriority }

// Migrations returns the embedded SQLite migrations for teams.
func (m *Module) Migrations() fs.FS { return migrationsFS }

// RegisterRoutes mounts teams CRUD and membership routes on AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Get("/teams", m.handleList)
	mounts.AdminAPI.Post("/teams", m.handleCreate)
	mounts.AdminAPI.Get("/teams/{id}", m.handleGet)
	mounts.AdminAPI.Patch("/teams/{id}", m.handlePatch)
	mounts.AdminAPI.Delete("/teams/{id}", m.handleDelete)
	mounts.AdminAPI.Post("/teams/{id}/keys", m.handleAddKey)
	mounts.AdminAPI.Delete("/teams/{id}/keys/{keyId}", m.handleRemoveKey)
	mounts.AdminAPI.Post("/teams/{id}/reset-spend", m.handleResetSpend)
}

// Start warms the in-memory team cache and subscribes to usage.recorded events.
func (m *Module) Start(ctx context.Context) error {
	if err := m.refreshCache(ctx); err != nil {
		return err
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

// Stop stops the usage.recorded subscriber cleanly.
func (m *Module) Stop(ctx context.Context) error {
	m.lifeMu.Lock()
	if !m.started {
		m.lifeMu.Unlock()
		return nil
	}
	cancel := m.cancel
	m.started = false
	m.lifeMu.Unlock()

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

// AssignKeyToTeam implements contracts.TeamAssigner.
func (m *Module) AssignKeyToTeam(ctx context.Context, teamID, keyID string) error {
	if strings.TrimSpace(teamID) == "" || strings.TrimSpace(keyID) == "" {
		return errors.New("team_id and key_id are required")
	}
	if _, err := m.repo.assignKey(ctx, teamID, keyID); err != nil {
		return err
	}
	_ = m.refreshCache(ctx)
	m.publishChanged("key_assigned", teamID)
	return nil
}

// OnRequest enforces shared team enabled state, USD budget, and sliding-window RPM.
func (m *Module) OnRequest(_ context.Context, req *contracts.CanonicalRequest) error {
	if req == nil || strings.TrimSpace(req.KeyID) == "" {
		return nil
	}

	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()

	teamID, ok := m.keyToTeam[req.KeyID]
	if !ok {
		return nil
	}
	team, exists := m.teamsByID[teamID]
	if !exists {
		return nil
	}

	if !team.Enabled {
		return statusError{
			code: http.StatusForbidden,
			msg:  fmt.Sprintf("team %q is disabled", team.Name),
		}
	}
	if team.BudgetUSD > 0 && team.SpentUSD >= team.BudgetUSD {
		return statusError{
			code: http.StatusTooManyRequests,
			msg:  fmt.Sprintf("team %q USD budget exhausted (%.4f / %.4f)", team.Name, team.SpentUSD, team.BudgetUSD),
		}
	}
	if team.RateLimitRPM > 0 {
		cutoff := now.Add(-time.Minute)
		window := m.rpmWindows[teamID]
		pruned := window[:0]
		for _, ts := range window {
			if ts.After(cutoff) {
				pruned = append(pruned, ts)
			}
		}
		if len(pruned) >= team.RateLimitRPM {
			m.rpmWindows[teamID] = pruned
			return statusError{
				code: http.StatusTooManyRequests,
				msg:  fmt.Sprintf("team %q rate limit exceeded (%d RPM)", team.Name, team.RateLimitRPM),
			}
		}
		m.rpmWindows[teamID] = append(pruned, now)
	}

	if req.Metadata != nil {
		req.Metadata["team_id"] = team.ID
		req.Metadata["team_name"] = team.Name
	}
	return nil
}

// OnResponseHeaders attaches team header when a request was billed to a team.
func (m *Module) OnResponseHeaders(_ context.Context, meta *contracts.ResponseMeta) error {
	if meta == nil || meta.Headers == nil || meta.Metadata == nil {
		return nil
	}
	if name, ok := meta.Metadata["team_name"].(string); ok && name != "" {
		meta.Headers.Set("X-Omeg4mon-Team", name)
	}
	return nil
}

// OnStreamChunk is a no-op for teams.
func (m *Module) OnStreamChunk(_ context.Context, _ *contracts.StreamEvent) error {
	return nil
}

func (m *Module) refreshCache(ctx context.Context) error {
	list, err := m.repo.list(ctx)
	if err != nil {
		return err
	}

	teamsByID := make(map[string]Team, len(list))
	keyToTeam := make(map[string]string)
	for _, t := range list {
		teamsByID[t.ID] = t
		for _, kid := range t.KeyIDs {
			keyToTeam[kid] = t.ID
		}
	}

	m.mu.Lock()
	m.teamsByID = teamsByID
	m.keyToTeam = keyToTeam
	m.mu.Unlock()
	return nil
}

func (m *Module) consumeUsage(ctx context.Context, ch <-chan eventbus.Event, unsub func()) {
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
			keyID, costUSD := extractUsageSpend(evt.Payload)
			if keyID == "" || costUSD <= 0 {
				continue
			}
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			teamID, err := m.repo.addSpendForKey(writeCtx, keyID, costUSD)
			cancel()
			if err != nil {
				if ctx.Err() == nil {
					m.logger.Warn("failed to accumulate team spend", "key_id", keyID, "error", err)
				}
				continue
			}
			if teamID != "" {
				m.mu.Lock()
				if t, exists := m.teamsByID[teamID]; exists {
					t.SpentUSD += costUSD
					m.teamsByID[teamID] = t
				}
				m.mu.Unlock()
			}
		}
	}
}

func extractUsageSpend(payload any) (keyID string, costUSD float64) {
	if payload == nil {
		return "", 0
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", 0
	}
	var parsed struct {
		KeyID   string  `json:"key_id"`
		CostUSD float64 `json:"cost_usd"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", 0
	}
	if math.IsNaN(parsed.CostUSD) || math.IsInf(parsed.CostUSD, 0) || parsed.CostUSD <= 0 {
		return "", 0
	}
	return strings.TrimSpace(parsed.KeyID), parsed.CostUSD
}

func (m *Module) publishChanged(action, teamID string) {
	if m.bus != nil {
		m.bus.Publish(TopicTeamsChanged, map[string]any{
			"action":  action,
			"team_id": teamID,
		})
	}
}

func (m *Module) handleList(w http.ResponseWriter, r *http.Request) {
	items, err := m.repo.list(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "teams_list_failed", "failed to list teams")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"teams": items})
}

func (m *Module) handleGet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	team, err := m.repo.getByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, errTeamNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "team_not_found", "team not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "team_get_failed", "failed to fetch team")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, team)
}

func (m *Module) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name         string  `json:"name"`
		Description  string  `json:"description"`
		BudgetUSD    float64 `json:"budget_usd"`
		RateLimitRPM int     `json:"rate_limit_rpm"`
		Enabled      *bool   `json:"enabled"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || len(req.Name) > 128 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_name", "name must be 1..128 characters")
		return
	}
	if req.BudgetUSD < 0 || math.IsNaN(req.BudgetUSD) || math.IsInf(req.BudgetUSD, 0) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_budget", "budget_usd must be a non-negative finite number")
		return
	}
	if req.RateLimitRPM < 0 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_rate_limit", "rate_limit_rpm must be non-negative")
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	created, err := m.repo.create(r.Context(), createTeamInput{
		Name:         req.Name,
		Description:  req.Description,
		BudgetUSD:    req.BudgetUSD,
		RateLimitRPM: req.RateLimitRPM,
		Enabled:      enabled,
	})
	if err != nil {
		if errors.Is(err, errTeamExists) {
			httpserver.WriteError(w, r, http.StatusConflict, "team_exists", "team name already exists")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "team_create_failed", "failed to create team")
		return
	}
	_ = m.refreshCache(r.Context())
	m.publishChanged("created", created.ID)
	httpserver.WriteJSON(w, http.StatusCreated, created)
}

func (m *Module) handlePatch(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	var req struct {
		Name         *string  `json:"name"`
		Description  *string  `json:"description"`
		BudgetUSD    *float64 `json:"budget_usd"`
		RateLimitRPM *int     `json:"rate_limit_rpm"`
		Enabled      *bool    `json:"enabled"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if req.Name != nil && (strings.TrimSpace(*req.Name) == "" || len(strings.TrimSpace(*req.Name)) > 128) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_name", "name must be 1..128 characters")
		return
	}
	if req.BudgetUSD != nil && (*req.BudgetUSD < 0 || math.IsNaN(*req.BudgetUSD) || math.IsInf(*req.BudgetUSD, 0)) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_budget", "budget_usd must be a non-negative finite number")
		return
	}
	if req.RateLimitRPM != nil && *req.RateLimitRPM < 0 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_rate_limit", "rate_limit_rpm must be non-negative")
		return
	}

	updated, err := m.repo.patch(r.Context(), id, patchTeamInput{
		Name:         req.Name,
		Description:  req.Description,
		BudgetUSD:    req.BudgetUSD,
		RateLimitRPM: req.RateLimitRPM,
		Enabled:      req.Enabled,
	})
	if err != nil {
		switch {
		case errors.Is(err, errTeamNotFound):
			httpserver.WriteError(w, r, http.StatusNotFound, "team_not_found", "team not found")
		case errors.Is(err, errTeamExists):
			httpserver.WriteError(w, r, http.StatusConflict, "team_exists", "team name already exists")
		default:
			httpserver.WriteError(w, r, http.StatusInternalServerError, "team_update_failed", "failed to update team")
		}
		return
	}
	_ = m.refreshCache(r.Context())
	m.publishChanged("updated", updated.ID)
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

func (m *Module) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if err := m.repo.deleteByID(r.Context(), id); err != nil {
		if errors.Is(err, errTeamNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "team_not_found", "team not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "team_delete_failed", "failed to delete team")
		return
	}
	_ = m.refreshCache(r.Context())
	m.publishChanged("deleted", id)
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": true})
}

func (m *Module) handleAddKey(w http.ResponseWriter, r *http.Request) {
	teamID := strings.TrimSpace(chi.URLParam(r, "id"))
	var req struct {
		KeyID string `json:"key_id"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	req.KeyID = strings.TrimSpace(req.KeyID)
	if req.KeyID == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_key_id", "key_id is required")
		return
	}
	updated, err := m.repo.assignKey(r.Context(), teamID, req.KeyID)
	if err != nil {
		if errors.Is(err, errTeamNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "team_not_found", "team not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "team_assign_key_failed", "failed to assign key to team")
		return
	}
	_ = m.refreshCache(r.Context())
	m.publishChanged("key_assigned", teamID)
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

func (m *Module) handleRemoveKey(w http.ResponseWriter, r *http.Request) {
	teamID := strings.TrimSpace(chi.URLParam(r, "id"))
	keyID := strings.TrimSpace(chi.URLParam(r, "keyId"))
	updated, err := m.repo.removeKey(r.Context(), teamID, keyID)
	if err != nil {
		if errors.Is(err, errTeamNotFound) || errors.Is(err, errKeyNotInTeam) {
			httpserver.WriteError(w, r, http.StatusNotFound, "not_found", err.Error())
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "team_remove_key_failed", "failed to remove key from team")
		return
	}
	_ = m.refreshCache(r.Context())
	m.publishChanged("key_removed", teamID)
	httpserver.WriteJSON(w, http.StatusOK, updated)
}

func (m *Module) handleResetSpend(w http.ResponseWriter, r *http.Request) {
	teamID := strings.TrimSpace(chi.URLParam(r, "id"))
	updated, err := m.repo.resetSpend(r.Context(), teamID)
	if err != nil {
		if errors.Is(err, errTeamNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "team_not_found", "team not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "team_reset_failed", "failed to reset team spend")
		return
	}
	_ = m.refreshCache(r.Context())
	m.publishChanged("spend_reset", teamID)
	httpserver.WriteJSON(w, http.StatusOK, updated)
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

package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

const (
	quotaCacheTTL           = 15 * time.Second
	quotaErrorCacheTTL      = 5 * time.Second
	maxQuotaCacheEntries    = 512
	maxQuotaConcurrency     = 6
	quotaFetchTimeout       = 10 * time.Second
	maxQuotaBodyBytes       = 256 << 10 // 256 KiB
	quotaLowRemainingPct    = 20.0
	harborDefaultBaseOrigin = "https://tokenharbor.ai"
	harborQuotaEndpointPath = "/api/cli/free-tier"
)

var errQuotaUnauthorized = errors.New("upstream rejected credentials")

// ProviderQuotaBucket represents a single normalized quota window (for example
// a paid subscription pass or a free-tier rolling allowance) for an upstream
// provider connection.
type ProviderQuotaBucket struct {
	ID              string  `json:"id"`
	Label           string  `json:"label"`
	UsedPct         float64 `json:"used_pct"`
	RemainingPct    float64 `json:"remaining_pct"`
	Exhausted       bool    `json:"exhausted"`
	PeriodStartedAt *string `json:"period_started_at"`
	PeriodEndsAt    *string `json:"period_ends_at"`
	ResetAt         *string `json:"reset_at"`
}

// ProviderConnectionQuota represents the normalized upstream quota snapshot for
// a single configured provider connection.
type ProviderConnectionQuota struct {
	ConnectionID   string                `json:"connection_id"`
	ProviderID     string                `json:"provider_id"`
	ProviderName   string                `json:"provider_name"`
	FetcherID      string                `json:"fetcher_id"`
	ConnectionName string                `json:"connection_name"`
	AuthType       string                `json:"auth_type"`
	Enabled        bool                  `json:"enabled"`
	Status         string                `json:"status"`
	PlanName       *string               `json:"plan_name"`
	OnPlan         bool                  `json:"on_plan"`
	Activated      bool                  `json:"activated"`
	Buckets        []ProviderQuotaBucket `json:"buckets"`
	Models         []string              `json:"models"`
	FetchedAt      string                `json:"fetched_at"`
	Cached         bool                  `json:"cached"`
	Error          string                `json:"error"`
}

// ProviderQuotasResponse is the response envelope for GET /api/providers/quotas.
type ProviderQuotasResponse struct {
	Quotas             []ProviderConnectionQuota `json:"quotas"`
	SupportedProviders []string                  `json:"supported_providers"`
	FetchedAt          string                    `json:"fetched_at"`
}

type quotaFetchResult struct {
	Status    string
	PlanName  *string
	OnPlan    bool
	Activated bool
	Buckets   []ProviderQuotaBucket
	Models    []string
}

type quotaFetcher interface {
	ID() string
	Matches(conn ConnectionSummary) bool
	Fetch(ctx context.Context, client *http.Client, conn ConnectionSummary, creds contracts.UpstreamCredentials) (quotaFetchResult, error)
}

func defaultQuotaFetchers() []quotaFetcher {
	return []quotaFetcher{
		tokenHarborQuotaFetcher{},
		codebuddyQuotaFetcher{providerID: providerIDCodeBuddyCN},
		codebuddyQuotaFetcher{providerID: providerIDCodeBuddyIntl},
		jetbrainsQuotaFetcher{},
		opencodeQuotaFetcher{providerID: "opencode-zen", planName: "OpenCode Zen", usageURL: "https://opencode.ai/zen/v1/usage"},
		opencodeQuotaFetcher{providerID: "opencode-go", planName: "OpenCode Go", usageURL: "https://opencode.ai/zen/go/v1/usage"},
	}
}

type quotaCacheEntry struct {
	quota     ProviderConnectionQuota
	expiresAt time.Time
}

type quotaFlightCall struct {
	done   chan struct{}
	result ProviderConnectionQuota
}

type quotaCache struct {
	mu       sync.Mutex
	ttl      time.Duration
	entries  map[string]quotaCacheEntry
	inflight map[string]*quotaFlightCall
}

func newQuotaCache(ttl time.Duration) *quotaCache {
	if ttl <= 0 {
		ttl = quotaCacheTTL
	}
	return &quotaCache{
		ttl:      ttl,
		entries:  make(map[string]quotaCacheEntry),
		inflight: make(map[string]*quotaFlightCall),
	}
}

func (c *quotaCache) Get(connectionID string, now time.Time) (ProviderConnectionQuota, bool) {
	if c == nil || connectionID == "" {
		return ProviderConnectionQuota{}, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[connectionID]
	if !ok {
		return ProviderConnectionQuota{}, false
	}
	if !now.Before(entry.expiresAt) {
		delete(c.entries, connectionID)
		return ProviderConnectionQuota{}, false
	}
	return cloneConnectionQuota(entry.quota), true
}

func (c *quotaCache) Set(connectionID string, quota ProviderConnectionQuota, now time.Time) {
	c.SetWithTTL(connectionID, quota, now, c.ttl)
}

func (c *quotaCache) SetWithTTL(connectionID string, quota ProviderConnectionQuota, now time.Time, ttl time.Duration) {
	if c == nil || connectionID == "" {
		return
	}
	if ttl <= 0 {
		ttl = c.ttl
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.entries) >= maxQuotaCacheEntries {
		for k, v := range c.entries {
			if !now.Before(v.expiresAt) {
				delete(c.entries, k)
			}
		}
		if len(c.entries) >= maxQuotaCacheEntries {
			var oldestKey string
			var oldestExp time.Time
			for k, v := range c.entries {
				if oldestKey == "" || v.expiresAt.Before(oldestExp) {
					oldestKey = k
					oldestExp = v.expiresAt
				}
			}
			if oldestKey != "" {
				delete(c.entries, oldestKey)
			}
		}
	}

	c.entries[connectionID] = quotaCacheEntry{
		quota:     cloneConnectionQuota(quota),
		expiresAt: now.Add(ttl),
	}
}

func (c *quotaCache) Invalidate(connectionID string) {
	if c == nil || connectionID == "" {
		return
	}
	c.mu.Lock()
	delete(c.entries, connectionID)
	c.mu.Unlock()
}

// AcquireFlight joins an existing in-flight fetch for connectionID or registers
// a new one. When owner is true, the caller must call FinishFlight(connectionID, call, result).
func (c *quotaCache) AcquireFlight(connectionID string) (call *quotaFlightCall, owner bool) {
	if c == nil || connectionID == "" {
		return nil, true
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if existing, ok := c.inflight[connectionID]; ok {
		return existing, false
	}
	call = &quotaFlightCall{done: make(chan struct{})}
	c.inflight[connectionID] = call
	return call, true
}

func (c *quotaCache) FinishFlight(connectionID string, call *quotaFlightCall, result ProviderConnectionQuota) {
	if c == nil || call == nil {
		return
	}
	call.result = cloneConnectionQuota(result)
	c.mu.Lock()
	if current, ok := c.inflight[connectionID]; ok && current == call {
		delete(c.inflight, connectionID)
	}
	c.mu.Unlock()
	close(call.done)
}

func cloneConnectionQuota(src ProviderConnectionQuota) ProviderConnectionQuota {
	cp := src
	if src.PlanName != nil {
		v := *src.PlanName
		cp.PlanName = &v
	}
	cp.Models = cloneStringSlice(src.Models)
	if len(src.Buckets) == 0 {
		cp.Buckets = []ProviderQuotaBucket{}
	} else {
		cp.Buckets = make([]ProviderQuotaBucket, len(src.Buckets))
		for i, b := range src.Buckets {
			nb := b
			if b.PeriodStartedAt != nil {
				v := *b.PeriodStartedAt
				nb.PeriodStartedAt = &v
			}
			if b.PeriodEndsAt != nil {
				v := *b.PeriodEndsAt
				nb.PeriodEndsAt = &v
			}
			if b.ResetAt != nil {
				v := *b.ResetAt
				nb.ResetAt = &v
			}
			cp.Buckets[i] = nb
		}
	}
	return cp
}

// tokenHarborQuotaFetcher queries Token Harbor's read-only /api/cli/free-tier
// endpoint for both built-in tokenharbor connections and custom/compatible
// connections targeting tokenharbor.ai.
type tokenHarborQuotaFetcher struct{}

type tokenHarborFreeTierPayload struct {
	PlanName        *string  `json:"plan_name"`
	OnPlan          bool     `json:"on_plan"`
	PlanUsedPct     *float64 `json:"plan_used_pct"`
	PlanExhausted   bool     `json:"plan_exhausted"`
	UsedPct         *float64 `json:"used_pct"`
	Exhausted       bool     `json:"exhausted"`
	Activated       bool     `json:"activated"`
	PeriodStartedAt *string  `json:"period_started_at"`
	PeriodEndsAt    *string  `json:"period_ends_at"`
	ResetAt         *string  `json:"reset_at"`
	Models          []string `json:"models"`
}

func (tokenHarborQuotaFetcher) ID() string {
	return "tokenharbor"
}

func (tokenHarborQuotaFetcher) Matches(conn ConnectionSummary) bool {
	if strings.EqualFold(strings.TrimSpace(conn.ProviderID), "tokenharbor") {
		return true
	}
	return isTokenHarborHost(conn.BaseURL)
}

func isTokenHarborHost(rawURL string) bool {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return false
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
	}
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	return host == "tokenharbor.ai" || strings.HasSuffix(host, ".tokenharbor.ai")
}

func resolveTokenHarborQuotaURL(baseURL string) string {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed != "" {
		if parsed, err := url.Parse(trimmed); err == nil &&
			(strings.EqualFold(parsed.Scheme, "http") || strings.EqualFold(parsed.Scheme, "https")) &&
			parsed.Host != "" {
			return parsed.Scheme + "://" + parsed.Host + harborQuotaEndpointPath
		}
	}
	return harborDefaultBaseOrigin + harborQuotaEndpointPath
}

func (tokenHarborQuotaFetcher) Fetch(
	ctx context.Context,
	client *http.Client,
	conn ConnectionSummary,
	creds contracts.UpstreamCredentials,
) (quotaFetchResult, error) {
	token := strings.TrimSpace(creds.APIKey)
	if token == "" {
		token = strings.TrimSpace(creds.AccessToken)
	}
	if token == "" {
		return quotaFetchResult{}, errors.New("missing API key or access token")
	}

	baseURL := strings.TrimSpace(creds.BaseURL)
	if baseURL == "" {
		baseURL = conn.BaseURL
	}
	endpoint := resolveTokenHarborQuotaURL(baseURL)

	reqCtx, cancel := context.WithTimeout(httpclient.WithProviderID(ctx, conn.ProviderID), quotaFetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return quotaFetchResult{}, fmt.Errorf("build quota request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return quotaFetchResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxQuotaBodyBytes))
	if err != nil {
		return quotaFetchResult{}, fmt.Errorf("read quota response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return quotaFetchResult{}, fmt.Errorf("%w (HTTP %d)", errQuotaUnauthorized, resp.StatusCode)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		snippet := strings.TrimSpace(string(body))
		if snippet != "" {
			return quotaFetchResult{}, fmt.Errorf("upstream returned HTTP %d: %s", resp.StatusCode, snippet)
		}
		return quotaFetchResult{}, fmt.Errorf("upstream returned HTTP %d", resp.StatusCode)
	}

	var payload tokenHarborFreeTierPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return quotaFetchResult{}, fmt.Errorf("decode Token Harbor quota response: %w", err)
	}

	return mapTokenHarborQuota(payload), nil
}

func mapTokenHarborQuota(payload tokenHarborFreeTierPayload) quotaFetchResult {
	planName := normalizeOptionalString(payload.PlanName)
	periodStart := normalizeOptionalString(payload.PeriodStartedAt)
	periodEnd := normalizeOptionalString(payload.PeriodEndsAt)
	resetAt := normalizeOptionalString(payload.ResetAt)
	if resetAt == nil {
		resetAt = periodEnd
	}

	buckets := make([]ProviderQuotaBucket, 0, 2)
	if payload.OnPlan {
		usedPct := clampPercent(derefFloat(payload.PlanUsedPct))
		exhausted := payload.PlanExhausted || usedPct >= 100
		if exhausted {
			usedPct = 100
		}
		remainingPct := clampPercent(100 - usedPct)
		label := "Plan Pass"
		if planName != nil {
			label = *planName
		}
		buckets = append(buckets, ProviderQuotaBucket{
			ID:              "plan",
			Label:           label,
			UsedPct:         usedPct,
			RemainingPct:    remainingPct,
			Exhausted:       exhausted,
			PeriodStartedAt: cloneOptionalString(periodStart),
			PeriodEndsAt:    cloneOptionalString(periodEnd),
			ResetAt:         cloneOptionalString(resetAt),
		})
	}

	freeUsedPct := clampPercent(derefFloat(payload.UsedPct))
	freeExhausted := payload.Exhausted || freeUsedPct >= 100
	if freeExhausted {
		freeUsedPct = 100
	}
	freeRemainingPct := clampPercent(100 - freeUsedPct)
	buckets = append(buckets, ProviderQuotaBucket{
		ID:              "free_tier",
		Label:           "Free Tier",
		UsedPct:         freeUsedPct,
		RemainingPct:    freeRemainingPct,
		Exhausted:       freeExhausted,
		PeriodStartedAt: cloneOptionalString(periodStart),
		PeriodEndsAt:    cloneOptionalString(periodEnd),
		ResetAt:         cloneOptionalString(resetAt),
	})

	status := deriveTokenHarborStatus(payload, buckets)
	models := cloneStringSlice(payload.Models)

	return quotaFetchResult{
		Status:    status,
		PlanName:  planName,
		OnPlan:    payload.OnPlan,
		Activated: payload.Activated,
		Buckets:   buckets,
		Models:    models,
	}
}

func deriveTokenHarborStatus(payload tokenHarborFreeTierPayload, buckets []ProviderQuotaBucket) string {
	if !payload.OnPlan && !payload.Activated {
		return "not_activated"
	}

	activeBuckets := make([]ProviderQuotaBucket, 0, len(buckets))
	for _, b := range buckets {
		if b.ID == "free_tier" && !payload.Activated {
			continue
		}
		activeBuckets = append(activeBuckets, b)
	}
	if len(activeBuckets) == 0 {
		return "not_activated"
	}

	allExhausted := true
	anyLow := false
	for _, b := range activeBuckets {
		if !b.Exhausted {
			allExhausted = false
		}
		if b.Exhausted || b.RemainingPct <= quotaLowRemainingPct {
			anyLow = true
		}
	}
	if allExhausted {
		return "exhausted"
	}
	if anyLow {
		return "low"
	}
	return "ok"
}

func derefFloat(v *float64) float64 {
	if v == nil || math.IsNaN(*v) || math.IsInf(*v, 0) {
		return 0
	}
	return *v
}

func clampPercent(v float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return math.Round(v*100) / 100
}

func normalizeOptionalString(v *string) *string {
	if v == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*v)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func cloneOptionalString(v *string) *string {
	if v == nil {
		return nil
	}
	cp := *v
	return &cp
}

func (m *Module) matchQuotaFetcher(conn ConnectionSummary) quotaFetcher {
	for _, f := range m.quotaFetchers {
		if f.Matches(conn) {
			return f
		}
	}
	return nil
}

func (m *Module) supportedQuotaProviders() []string {
	out := make([]string, 0, len(m.quotaFetchers))
	for _, f := range m.quotaFetchers {
		out = append(out, f.ID())
	}
	return out
}

func (m *Module) resolveProviderDisplayNames(ctx context.Context) map[string]string {
	names := make(map[string]string)
	for _, adapter := range m.ListAdapters() {
		names[adapter.ID()] = adapter.DisplayName()
	}
	if nodes, err := m.repo.ListNodes(ctx); err == nil {
		for _, node := range nodes {
			if strings.TrimSpace(node.Name) != "" {
				names[node.ID] = node.Name
			}
		}
	}
	return names
}

func (m *Module) fetchConnectionQuota(
	ctx context.Context,
	conn ConnectionSummary,
	fetcher quotaFetcher,
	providerName string,
	forceRefresh bool,
) ProviderConnectionQuota {
	now := time.Now().UTC()
	if !forceRefresh {
		if cached, ok := m.quotaCache.Get(conn.ID, now); ok {
			cached.Cached = true
			cached.ConnectionName = conn.Name
			cached.Enabled = conn.Enabled
			cached.ProviderName = providerName
			return cached
		}
	}

	flight, owner := m.quotaCache.AcquireFlight(conn.ID)
	if !owner {
		select {
		case <-flight.done:
			joined := cloneConnectionQuota(flight.result)
			joined.ConnectionName = conn.Name
			joined.Enabled = conn.Enabled
			joined.ProviderName = providerName
			return joined
		case <-ctx.Done():
			return ProviderConnectionQuota{
				ConnectionID:   conn.ID,
				ProviderID:     conn.ProviderID,
				ProviderName:   providerName,
				FetcherID:      fetcher.ID(),
				ConnectionName: conn.Name,
				AuthType:       conn.AuthType,
				Enabled:        conn.Enabled,
				Status:         "error",
				Buckets:        []ProviderQuotaBucket{},
				Models:         []string{},
				FetchedAt:      time.Now().UTC().Format(rfc3339NanoFixed),
				Error:          "request canceled",
			}
		}
	}

	// Re-check cache after acquiring flight ownership in case a preceding
	// leader just populated it (when forceRefresh is false).
	if !forceRefresh {
		if cached, ok := m.quotaCache.Get(conn.ID, time.Now().UTC()); ok {
			cached.Cached = true
			cached.ConnectionName = conn.Name
			cached.Enabled = conn.Enabled
			cached.ProviderName = providerName
			m.quotaCache.FinishFlight(conn.ID, flight, cached)
			return cached
		}
	}

	result := ProviderConnectionQuota{
		ConnectionID:   conn.ID,
		ProviderID:     conn.ProviderID,
		ProviderName:   providerName,
		FetcherID:      fetcher.ID(),
		ConnectionName: conn.Name,
		AuthType:       conn.AuthType,
		Enabled:        conn.Enabled,
		Status:         "error",
		Buckets:        []ProviderQuotaBucket{},
		Models:         []string{},
		FetchedAt:      now.Format(rfc3339NanoFixed),
		Cached:         false,
	}
	defer func() {
		m.quotaCache.FinishFlight(conn.ID, flight, result)
	}()

	// Detach caller cancellation with a bounded timeout so a single canceled
	// client does not abort a coalesced upstream fetch shared with other callers.
	fetchCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), quotaFetchTimeout)
	defer cancel()

	creds, err := m.ResolveCredentials(fetchCtx, conn.ID)
	if err != nil {
		result.Error = truncateError("failed to resolve credentials")
		m.quotaCache.SetWithTTL(conn.ID, result, time.Now().UTC(), quotaErrorCacheTTL)
		return result
	}

	fetched, err := fetcher.Fetch(fetchCtx, m.httpClient, conn, creds)
	if err != nil && conn.AuthType == authTypeOAuth && errors.Is(err, errQuotaUnauthorized) {
		if refreshedCreds, refreshErr := m.RefreshOAuthToken(fetchCtx, conn.ID); refreshErr == nil {
			creds = refreshedCreds
			fetched, err = fetcher.Fetch(fetchCtx, m.httpClient, conn, creds)
		}
	}
	if err != nil {
		result.Error = truncateError(sanitizeUpstreamError(err.Error(), creds))
		if !errors.Is(err, context.Canceled) {
			m.quotaCache.SetWithTTL(conn.ID, result, time.Now().UTC(), quotaErrorCacheTTL)
		}
		return result
	}

	result.Status = fetched.Status
	result.PlanName = fetched.PlanName
	result.OnPlan = fetched.OnPlan
	result.Activated = fetched.Activated
	result.Buckets = fetched.Buckets
	result.Models = fetched.Models
	result.FetchedAt = time.Now().UTC().Format(rfc3339NanoFixed)
	result.Error = ""

	m.quotaCache.Set(conn.ID, result, time.Now().UTC())
	return result
}

func (m *Module) handleListProviderQuotas(w http.ResponseWriter, r *http.Request) {
	connections, err := m.repo.ListConnections(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load provider connections")
		return
	}

	providerFilter := strings.TrimSpace(r.URL.Query().Get("provider_id"))
	forceRefresh := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("refresh")), "true") ||
		strings.TrimSpace(r.URL.Query().Get("refresh")) == "1"

	type matchedTarget struct {
		conn    ConnectionSummary
		fetcher quotaFetcher
	}

	targets := make([]matchedTarget, 0, len(connections))
	for _, conn := range connections {
		fetcher := m.matchQuotaFetcher(conn)
		if fetcher == nil {
			continue
		}
		if providerFilter != "" &&
			!strings.EqualFold(conn.ProviderID, providerFilter) &&
			!strings.EqualFold(fetcher.ID(), providerFilter) {
			continue
		}
		targets = append(targets, matchedTarget{conn: conn, fetcher: fetcher})
	}

	providerNames := m.resolveProviderDisplayNames(r.Context())
	quotas := make([]ProviderConnectionQuota, len(targets))

	if len(targets) > 0 {
		sem := make(chan struct{}, maxQuotaConcurrency)
		var wg sync.WaitGroup

		for i, target := range targets {
			providerName := providerNames[target.conn.ProviderID]
			if providerName == "" {
				providerName = target.conn.ProviderID
			}

			wg.Add(1)
			go func(idx int, t matchedTarget, displayName string) {
				defer wg.Done()
				select {
				case sem <- struct{}{}:
					defer func() { <-sem }()
				case <-r.Context().Done():
					quotas[idx] = ProviderConnectionQuota{
						ConnectionID:   t.conn.ID,
						ProviderID:     t.conn.ProviderID,
						ProviderName:   displayName,
						FetcherID:      t.fetcher.ID(),
						ConnectionName: t.conn.Name,
						AuthType:       t.conn.AuthType,
						Enabled:        t.conn.Enabled,
						Status:         "error",
						Buckets:        []ProviderQuotaBucket{},
						Models:         []string{},
						FetchedAt:      time.Now().UTC().Format(rfc3339NanoFixed),
						Error:          "request canceled",
					}
					return
				}
				quotas[idx] = m.fetchConnectionQuota(r.Context(), t.conn, t.fetcher, displayName, forceRefresh)
			}(i, target, providerName)
		}
		wg.Wait()
	}

	httpserver.WriteJSON(w, http.StatusOK, ProviderQuotasResponse{
		Quotas:             quotas,
		SupportedProviders: m.supportedQuotaProviders(),
		FetchedAt:          time.Now().UTC().Format(rfc3339NanoFixed),
	})
}

func (m *Module) handleRefreshConnectionQuota(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "connection id is required")
		return
	}

	conn, err := m.repo.GetConnectionSummary(r.Context(), id)
	if err != nil {
		if errors.Is(err, contracts.ErrConnectionNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "connection_not_found", "provider connection not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load provider connection")
		return
	}

	fetcher := m.matchQuotaFetcher(conn)
	if fetcher == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "quota_unsupported", "provider does not support live quota inspection")
		return
	}

	providerNames := m.resolveProviderDisplayNames(r.Context())
	providerName := providerNames[conn.ProviderID]
	if providerName == "" {
		providerName = conn.ProviderID
	}

	quota := m.fetchConnectionQuota(r.Context(), conn, fetcher, providerName, true)
	httpserver.WriteJSON(w, http.StatusOK, quota)
}

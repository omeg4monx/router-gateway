package providers_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/modules/providers"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestHarnessWithClient(t *testing.T, auth *stubAdminAuth, client *http.Client) *testHarness {
	t.Helper()

	cfg := config.Config{
		Host:       "127.0.0.1",
		Port:       20128,
		APITimeout: 30 * time.Second,
	}

	dbPath := filepath.Join(t.TempDir(), "router.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})

	masterKey := make([]byte, crypto.MasterKeySize)
	if _, err := rand.Read(masterKey); err != nil {
		t.Fatalf("rand.Read: %v", err)
	}
	vault, err := crypto.NewVault(masterKey)
	if err != nil {
		t.Fatalf("crypto.NewVault: %v", err)
	}

	bus := eventbus.New()
	reg := lifecycle.NewRegistry(nil)
	mod := providers.New(providers.Options{
		DB:         database,
		Vault:      vault,
		HTTPClient: client,
		Bus:        bus,
	})

	if err := reg.Register(stubSystemModule{}); err != nil {
		t.Fatalf("reg.Register(system): %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("reg.Register(providers): %v", err)
	}

	ctx := context.Background()
	if err := reg.MigrateAll(ctx, database); err != nil {
		t.Fatalf("reg.MigrateAll: %v", err)
	}

	srv := httpserver.New(httpserver.Options{
		Config: cfg,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Auth:   auth,
	})
	if err := reg.MountAll(srv.Mounts()); err != nil {
		t.Fatalf("reg.MountAll: %v", err)
	}
	if err := reg.StartAll(ctx); err != nil {
		t.Fatalf("reg.StartAll: %v", err)
	}
	t.Cleanup(func() {
		_ = reg.StopAll(context.Background())
	})

	return &testHarness{
		db:     database,
		vault:  vault,
		bus:    bus,
		reg:    reg,
		server: srv,
		mod:    mod,
		auth:   auth,
	}
}

type quotaBucketDTO struct {
	ID              string  `json:"id"`
	Label           string  `json:"label"`
	UsedPct         float64 `json:"used_pct"`
	RemainingPct    float64 `json:"remaining_pct"`
	Exhausted       bool    `json:"exhausted"`
	PeriodStartedAt *string `json:"period_started_at"`
	PeriodEndsAt    *string `json:"period_ends_at"`
	ResetAt         *string `json:"reset_at"`
}

type connectionQuotaDTO struct {
	ConnectionID   string           `json:"connection_id"`
	ProviderID     string           `json:"provider_id"`
	ProviderName   string           `json:"provider_name"`
	FetcherID      string           `json:"fetcher_id"`
	ConnectionName string           `json:"connection_name"`
	AuthType       string           `json:"auth_type"`
	Enabled        bool             `json:"enabled"`
	Status         string           `json:"status"`
	PlanName       *string          `json:"plan_name"`
	OnPlan         bool             `json:"on_plan"`
	Activated      bool             `json:"activated"`
	Buckets        []quotaBucketDTO `json:"buckets"`
	Models         []string         `json:"models"`
	FetchedAt      string           `json:"fetched_at"`
	Cached         bool             `json:"cached"`
	Error          string           `json:"error"`
}

type quotasResponseDTO struct {
	Quotas             []connectionQuotaDTO `json:"quotas"`
	SupportedProviders []string             `json:"supported_providers"`
	FetchedAt          string               `json:"fetched_at"`
}

func TestQuotaBuiltinTokenHarborAdapterInCatalog(t *testing.T) {
	h := newTestHarness(t, &stubAdminAuth{subject: "admin"}, nil)

	rec := h.doRequest(t, http.MethodGet, "/api/providers/catalog", "", true)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/providers/catalog status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var catalog struct {
		Adapters []struct {
			ID             string   `json:"id"`
			DisplayName    string   `json:"display_name"`
			Protocol       string   `json:"protocol"`
			AuthType       string   `json:"auth_type"`
			DefaultBaseURL string   `json:"default_base_url"`
			DefaultModels  []string `json:"default_models"`
		} `json:"adapters"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &catalog); err != nil {
		t.Fatalf("decode catalog: %v", err)
	}

	found := false
	for _, a := range catalog.Adapters {
		if a.ID == "tokenharbor" {
			found = true
			if a.DisplayName != "Token Harbor" {
				t.Errorf("display_name = %q, want Token Harbor", a.DisplayName)
			}
			if a.Protocol != "openai-chat" {
				t.Errorf("protocol = %q, want openai-chat", a.Protocol)
			}
			if a.AuthType != "api_key" {
				t.Errorf("auth_type = %q, want api_key", a.AuthType)
			}
			if a.DefaultBaseURL != "https://tokenharbor.ai/v1" {
				t.Errorf("default_base_url = %q, want https://tokenharbor.ai/v1", a.DefaultBaseURL)
			}
		}
	}
	if !found {
		t.Fatalf("built-in adapter 'tokenharbor' not found in catalog")
	}
}

func TestQuotaTokenHarborListAndStatuses(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cli/free-tier" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Header.Get("Authorization") {
		case "Bearer thk_live_plan_ok":
			_, _ = w.Write([]byte(`{
				"plan_name": "Office Pass",
				"on_plan": true,
				"plan_used_pct": 37,
				"plan_exhausted": false,
				"used_pct": 12,
				"exhausted": false,
				"activated": true,
				"period_started_at": "2026-11-05T09:12:44.512+00:00",
				"period_ends_at": "2026-11-12T09:12:44.512+00:00",
				"reset_at": "2026-11-12T09:12:44.512+00:00",
				"models": ["deepseek-v4.1-flash:free", "mimo-v2.6-flash:free"]
			}`))
		case "Bearer thk_live_not_activated":
			_, _ = w.Write([]byte(`{
				"plan_name": null,
				"on_plan": false,
				"plan_used_pct": 0,
				"plan_exhausted": false,
				"used_pct": 0,
				"exhausted": false,
				"activated": false,
				"period_started_at": null,
				"period_ends_at": null,
				"reset_at": null,
				"models": []
			}`))
		case "Bearer thk_live_low":
			_, _ = w.Write([]byte(`{
				"plan_name": null,
				"on_plan": false,
				"plan_used_pct": 0,
				"plan_exhausted": false,
				"used_pct": 85,
				"exhausted": false,
				"activated": true,
				"reset_at": "2026-11-30T00:00:00Z",
				"models": ["deepseek-v4.1-flash:free"]
			}`))
		case "Bearer thk_live_exhausted":
			_, _ = w.Write([]byte(`{
				"plan_name": null,
				"on_plan": false,
				"plan_used_pct": 0,
				"plan_exhausted": false,
				"used_pct": 100,
				"exhausted": true,
				"activated": true,
				"reset_at": "2026-11-30T00:00:00Z",
				"models": ["deepseek-v4.1-flash:free"]
			}`))
		case "Bearer thk_live_secret_bad":
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"invalid api key thk_live_secret_bad"}`))
		default:
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"unknown"}`))
		}
	}))
	t.Cleanup(upstream.Close)

	h := newTestHarness(t, &stubAdminAuth{subject: "admin"}, nil)

	// Unauthenticated must be rejected with 401.
	if rec := h.doRequest(t, http.MethodGet, "/api/providers/quotas", "", false); rec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated GET /api/providers/quotas = %d, want 401", rec.Code)
	}

	createConn := func(providerID, name, baseURL, apiKey string) string {
		t.Helper()
		payload := fmt.Sprintf(`{"provider_id":%q,"name":%q,"base_url":%q,"api_key":%q}`, providerID, name, baseURL, apiKey)
		rec := h.doRequest(t, http.MethodPost, "/api/providers/connections", payload, true)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create connection %s status = %d: %s", name, rec.Code, rec.Body.String())
		}
		var out struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatalf("decode connection: %v", err)
		}
		return out.ID
	}

	baseURL := upstream.URL + "/v1"
	idPlan := createConn("tokenharbor", "Harbor Plan", baseURL, "thk_live_plan_ok")
	idInactive := createConn("tokenharbor", "Harbor Inactive", baseURL, "thk_live_not_activated")
	idLow := createConn("tokenharbor", "Harbor Low", baseURL, "thk_live_low")
	idExhausted := createConn("tokenharbor", "Harbor Exhausted", baseURL, "thk_live_exhausted")
	idBad := createConn("tokenharbor", "Harbor Bad", baseURL, "thk_live_secret_bad")
	_ = createConn("openai", "OpenAI Unrelated", "https://api.openai.com/v1", "sk-openai-test")

	rec := h.doRequest(t, http.MethodGet, "/api/providers/quotas", "", true)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/providers/quotas status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var resp quotasResponseDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode quotas response: %v", err)
	}
	if len(resp.Quotas) != 5 {
		t.Fatalf("got %d quota entries, want 5", len(resp.Quotas))
	}
	if len(resp.SupportedProviders) == 0 || resp.SupportedProviders[0] != "tokenharbor" {
		t.Errorf("supported_providers = %v, want tokenharbor", resp.SupportedProviders)
	}

	byID := make(map[string]connectionQuotaDTO, len(resp.Quotas))
	for _, q := range resp.Quotas {
		byID[q.ConnectionID] = q
	}

	qPlan := byID[idPlan]
	if qPlan.Status != "ok" || !qPlan.OnPlan || !qPlan.Activated || qPlan.PlanName == nil || *qPlan.PlanName != "Office Pass" {
		t.Errorf("unexpected plan quota summary: %+v", qPlan)
	}
	if len(qPlan.Buckets) != 2 {
		t.Fatalf("plan connection buckets = %d, want 2", len(qPlan.Buckets))
	}
	if qPlan.Buckets[0].ID != "plan" || qPlan.Buckets[0].UsedPct != 37 || qPlan.Buckets[0].RemainingPct != 63 {
		t.Errorf("unexpected plan bucket: %+v", qPlan.Buckets[0])
	}
	if qPlan.Buckets[1].ID != "free_tier" || qPlan.Buckets[1].UsedPct != 12 || qPlan.Buckets[1].RemainingPct != 88 {
		t.Errorf("unexpected free_tier bucket: %+v", qPlan.Buckets[1])
	}
	if len(qPlan.Models) != 2 {
		t.Errorf("expected 2 models, got %v", qPlan.Models)
	}

	if qInactive := byID[idInactive]; qInactive.Status != "not_activated" || len(qInactive.Buckets) != 1 {
		t.Errorf("unexpected not_activated entry: %+v", qInactive)
	}
	if qLow := byID[idLow]; qLow.Status != "low" || len(qLow.Buckets) != 1 || qLow.Buckets[0].RemainingPct != 15 {
		t.Errorf("unexpected low entry: %+v", qLow)
	}
	if qExh := byID[idExhausted]; qExh.Status != "exhausted" || len(qExh.Buckets) != 1 || !qExh.Buckets[0].Exhausted {
		t.Errorf("unexpected exhausted entry: %+v", qExh)
	}
	qBad := byID[idBad]
	if qBad.Status != "error" || qBad.Error == "" {
		t.Errorf("expected error entry for bad key, got %+v", qBad)
	}
	if strings.Contains(qBad.Error, "thk_live_secret_bad") {
		t.Errorf("error message leaked raw API key: %q", qBad.Error)
	}
}

func TestQuotaCacheRefreshAndPatchInvalidation(t *testing.T) {
	var hits atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := hits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"plan_name":null,"on_plan":false,"plan_used_pct":0,"plan_exhausted":false,"used_pct":%d,"exhausted":false,"activated":true,"models":["deepseek-v4.1-flash:free"]}`, n*10)
	}))
	t.Cleanup(upstream.Close)

	h := newTestHarness(t, &stubAdminAuth{subject: "admin"}, nil)

	recCreate := h.doRequest(t, http.MethodPost, "/api/providers/connections", fmt.Sprintf(
		`{"provider_id":"tokenharbor","name":"Cache Account","base_url":%q,"api_key":"thk_live_cache"}`,
		upstream.URL+"/v1",
	), true)
	if recCreate.Code != http.StatusCreated {
		t.Fatalf("create connection status = %d: %s", recCreate.Code, recCreate.Body.String())
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(recCreate.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created connection: %v", err)
	}

	// Create an unsupported OpenAI connection to test 400 on single refresh.
	recOAI := h.doRequest(t, http.MethodPost, "/api/providers/connections", `{"provider_id":"openai","name":"OpenAI","api_key":"sk-test"}`, true)
	if recOAI.Code != http.StatusCreated {
		t.Fatalf("create openai status = %d", recOAI.Code)
	}
	var oaiCreated struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(recOAI.Body.Bytes(), &oaiCreated)

	// 1st call -> live fetch
	rec1 := h.doRequest(t, http.MethodGet, "/api/providers/quotas", "", true)
	var resp1 quotasResponseDTO
	_ = json.Unmarshal(rec1.Body.Bytes(), &resp1)
	if hits.Load() != 1 || len(resp1.Quotas) != 1 || resp1.Quotas[0].Cached || resp1.Quotas[0].Buckets[0].UsedPct != 10 {
		t.Fatalf("1st fetch unexpected: hits=%d resp=%+v", hits.Load(), resp1)
	}

	// 2nd call -> cache hit
	rec2 := h.doRequest(t, http.MethodGet, "/api/providers/quotas", "", true)
	var resp2 quotasResponseDTO
	_ = json.Unmarshal(rec2.Body.Bytes(), &resp2)
	if hits.Load() != 1 || !resp2.Quotas[0].Cached || resp2.Quotas[0].Buckets[0].UsedPct != 10 {
		t.Fatalf("2nd fetch should hit cache: hits=%d resp=%+v", hits.Load(), resp2)
	}

	// 3rd call with refresh=true -> bypasses cache
	rec3 := h.doRequest(t, http.MethodGet, "/api/providers/quotas?refresh=true", "", true)
	var resp3 quotasResponseDTO
	_ = json.Unmarshal(rec3.Body.Bytes(), &resp3)
	if hits.Load() != 2 || resp3.Quotas[0].Cached || resp3.Quotas[0].Buckets[0].UsedPct != 20 {
		t.Fatalf("refresh=true should bypass cache: hits=%d resp=%+v", hits.Load(), resp3)
	}

	// Single-connection force refresh
	recSingle := h.doRequest(t, http.MethodPost, "/api/providers/connections/"+created.ID+"/quota/refresh", "", true)
	if recSingle.Code != http.StatusOK {
		t.Fatalf("single refresh status = %d: %s", recSingle.Code, recSingle.Body.String())
	}
	var single connectionQuotaDTO
	_ = json.Unmarshal(recSingle.Body.Bytes(), &single)
	if hits.Load() != 3 || single.Cached || single.Buckets[0].UsedPct != 30 {
		t.Fatalf("single refresh unexpected: hits=%d single=%+v", hits.Load(), single)
	}

	// PATCH connection invalidates cached quota
	recPatch := h.doRequest(t, http.MethodPatch, "/api/providers/connections/"+created.ID, `{"api_key":"thk_live_updated"}`, true)
	if recPatch.Code != http.StatusOK {
		t.Fatalf("patch connection status = %d: %s", recPatch.Code, recPatch.Body.String())
	}
	rec4 := h.doRequest(t, http.MethodGet, "/api/providers/quotas", "", true)
	var resp4 quotasResponseDTO
	_ = json.Unmarshal(rec4.Body.Bytes(), &resp4)
	if hits.Load() != 4 || resp4.Quotas[0].Cached {
		t.Fatalf("patch should invalidate cache: hits=%d resp=%+v", hits.Load(), resp4)
	}

	// 404 on unknown connection, 400 on unsupported provider
	if rec404 := h.doRequest(t, http.MethodPost, "/api/providers/connections/conn_missing/quota/refresh", "", true); rec404.Code != http.StatusNotFound {
		t.Fatalf("unknown connection refresh status = %d, want 404", rec404.Code)
	}
	if rec400 := h.doRequest(t, http.MethodPost, "/api/providers/connections/"+oaiCreated.ID+"/quota/refresh", "", true); rec400.Code != http.StatusBadRequest {
		t.Fatalf("unsupported provider refresh status = %d, want 400", rec400.Code)
	}
}

func TestQuotaMatchesCustomTokenHarborHost(t *testing.T) {
	var requestedURL string
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requestedURL = req.URL.String()
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body: io.NopCloser(strings.NewReader(`{
					"plan_name": "Agent Pass",
					"on_plan": true,
					"plan_used_pct": 25,
					"plan_exhausted": false,
					"used_pct": 5,
					"exhausted": false,
					"activated": true,
					"reset_at": "2026-12-01T00:00:00Z",
					"models": ["deepseek-v4.1-flash:free"]
				}`)),
			}, nil
		}),
	}

	h := newTestHarnessWithClient(t, &stubAdminAuth{subject: "admin"}, client)

	recNode := h.doRequest(t, http.MethodPost, "/api/providers/nodes", `{"id":"my-harbor","name":"My Token Harbor","protocol":"openai-chat","base_url":"https://tokenharbor.ai/v1","models":["deepseek-v4.1-flash:free"]}`, true)
	if recNode.Code != http.StatusCreated {
		t.Fatalf("create custom node status = %d: %s", recNode.Code, recNode.Body.String())
	}
	recConn := h.doRequest(t, http.MethodPost, "/api/providers/connections", `{"provider_id":"my-harbor","name":"Custom Harbor Key","api_key":"thk_live_custom"}`, true)
	if recConn.Code != http.StatusCreated {
		t.Fatalf("create custom connection status = %d: %s", recConn.Code, recConn.Body.String())
	}

	rec := h.doRequest(t, http.MethodGet, "/api/providers/quotas", "", true)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/providers/quotas status = %d: %s", rec.Code, rec.Body.String())
	}
	if requestedURL != "https://tokenharbor.ai/api/cli/free-tier" {
		t.Errorf("requested URL = %q, want https://tokenharbor.ai/api/cli/free-tier", requestedURL)
	}

	var resp quotasResponseDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Quotas) != 1 {
		t.Fatalf("got %d quotas, want 1", len(resp.Quotas))
	}
	if resp.Quotas[0].ProviderID != "my-harbor" || resp.Quotas[0].FetcherID != "tokenharbor" || resp.Quotas[0].ProviderName != "My Token Harbor" {
		t.Errorf("unexpected custom node quota metadata: %+v", resp.Quotas[0])
	}
}

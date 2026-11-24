package combos_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/combos"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type catalogFakeProviderStore struct {
	connections []contracts.ProviderConnection
	creds       map[string]contracts.UpstreamCredentials
	refreshFn   func(connectionID string) (contracts.UpstreamCredentials, error)
}

func (f *catalogFakeProviderStore) ListEnabledConnections(_ context.Context) ([]contracts.ProviderConnection, error) {
	return append([]contracts.ProviderConnection(nil), f.connections...), nil
}

func (f *catalogFakeProviderStore) ResolveCredentials(_ context.Context, connectionID string) (contracts.UpstreamCredentials, error) {
	c, ok := f.creds[connectionID]
	if !ok {
		return contracts.UpstreamCredentials{}, contracts.ErrConnectionNotFound
	}
	return c, nil
}

func (f *catalogFakeProviderStore) RefreshOAuthToken(_ context.Context, connectionID string) (contracts.UpstreamCredentials, error) {
	if f.refreshFn != nil {
		return f.refreshFn(connectionID)
	}
	return contracts.UpstreamCredentials{}, contracts.ErrOAuthRefreshUnsupported
}

func setupCatalogTestModule(t *testing.T, store contracts.ProviderStore, client *http.Client) (*combos.Module, chi.Router) {
	t.Helper()
	database, err := db.Open(t.TempDir() + "/catalog-test.sqlite")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	mod := combos.New(combos.Options{
		DB:         database,
		Providers:  store,
		HTTPClient: client,
		Now:        func() time.Time { return time.Date(2026, 11, 24, 12, 0, 0, 0, time.UTC) },
	})
	reg := lifecycle.NewRegistry(nil)
	if err := reg.Register(stubModule{name: "system"}); err != nil {
		t.Fatalf("register system: %v", err)
	}
	if err := reg.Register(stubModule{name: "providers", deps: []string{"system"}}); err != nil {
		t.Fatalf("register providers: %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("register combos: %v", err)
	}
	if err := reg.MigrateAll(context.Background(), database); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("start module: %v", err)
	}

	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		mod.RegisterRoutes(lifecycle.RouteMounts{AdminAPI: r})
	})
	return mod, router
}

func TestModelCatalogAndHealthMatrix(t *testing.T) {
	store := &catalogFakeProviderStore{
		connections: []contracts.ProviderConnection{
			{
				ID:         "conn_openai_1",
				ProviderID: "openai",
				Name:       "OpenAI Primary",
				Protocol:   "openai",
				BaseURL:    "https://api.openai.com/v1",
				Priority:   10,
				Enabled:    true,
				Models:     []string{"gpt-4o", "o3-mini"},
			},
			{
				ID:         "conn_openai_2",
				ProviderID: "openai",
				Name:       "OpenAI Backup",
				Protocol:   "openai",
				BaseURL:    "https://api.openai.com/v1",
				Priority:   20,
				Enabled:    true,
				Models:     []string{"gpt-4o"},
			},
			{
				ID:         "conn_anthropic_1",
				ProviderID: "anthropic",
				Name:       "Claude Primary",
				Protocol:   "anthropic",
				BaseURL:    "https://api.anthropic.com/v1",
				Priority:   10,
				Enabled:    true,
				Models:     []string{"claude-sonnet-4-5"},
			},
		},
	}

	mod, router := setupCatalogTestModule(t, store, http.DefaultClient)

	// Disable openai/o3-mini so it becomes hidden in the catalog.
	disableBody := bytes.NewBufferString(`{"model":"openai/o3-mini","reason":"hidden from catalog"}`)
	disableReq := httptest.NewRequest(http.MethodPut, "/api/models/disabled", disableBody)
	disableRec := httptest.NewRecorder()
	router.ServeHTTP(disableRec, disableReq)
	if disableRec.Code != http.StatusOK {
		t.Fatalf("expected 200 disabling model, got %d: %s", disableRec.Code, disableRec.Body.String())
	}

	// GET /api/models/catalog?includeHidden=1 should include both visible and hidden models.
	req := httptest.NewRequest(http.MethodGet, "/api/models/catalog?includeHidden=1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for catalog, got %d: %s", rec.Code, rec.Body.String())
	}

	var catResp struct {
		Catalog map[string]struct {
			Provider string `json:"provider"`
			Models   []struct {
				ID            string          `json:"id"`
				Type          string          `json:"type"`
				Subtype       string          `json:"subtype"`
				Hidden        bool            `json:"hidden"`
				ContextLength int             `json:"context_length"`
				Capabilities  map[string]bool `json:"capabilities"`
			} `json:"models"`
		} `json:"catalog"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &catResp); err != nil {
		t.Fatalf("unmarshal catalog: %v", err)
	}

	openaiGroup, ok := catResp.Catalog["openai"]
	if !ok {
		t.Fatalf("expected openai group in catalog, got %+v", catResp.Catalog)
	}
	if len(openaiGroup.Models) != 2 {
		t.Fatalf("expected 2 deduplicated openai models with includeHidden=1, got %d", len(openaiGroup.Models))
	}

	var foundHidden bool
	for _, m := range openaiGroup.Models {
		if m.ID == "o3-mini" {
			foundHidden = m.Hidden
			if !m.Capabilities["reasoning"] {
				t.Errorf("expected o3-mini to infer reasoning capability")
			}
		}
		if m.ID == "gpt-4o" && m.ContextLength == 0 {
			t.Errorf("expected gpt-4o to infer non-zero context_length")
		}
	}
	if !foundHidden {
		t.Errorf("expected openai/o3-mini to be marked hidden=true")
	}

	// Without includeHidden=1, hidden models should be excluded.
	reqVisible := httptest.NewRequest(http.MethodGet, "/api/models/catalog", nil)
	recVisible := httptest.NewRecorder()
	router.ServeHTTP(recVisible, reqVisible)
	if recVisible.Code != http.StatusOK {
		t.Fatalf("expected 200 for visible catalog, got %d", recVisible.Code)
	}
	if strings.Contains(recVisible.Body.String(), `"o3-mini"`) {
		t.Errorf("expected o3-mini to be filtered out when includeHidden is not set")
	}

	// Put one openai connection in cooldown and check health-matrix -> degraded, anthropic -> healthy.
	mod.RecordFailure("conn_openai_1", http.StatusTooManyRequests, 5*time.Minute)
	hmReq := httptest.NewRequest(http.MethodGet, "/api/providers/health-matrix", nil)
	hmRec := httptest.NewRecorder()
	router.ServeHTTP(hmRec, hmReq)
	if hmRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for health-matrix, got %d: %s", hmRec.Code, hmRec.Body.String())
	}

	var hmResp struct {
		Providers []struct {
			Provider               string `json:"provider"`
			State                  string `json:"state"`
			TotalConnections       int    `json:"totalConnections"`
			HealthyConnections     int    `json:"healthyConnections"`
			CoolingDownConnections int    `json:"coolingDownConnections"`
		} `json:"providers"`
	}
	if err := json.Unmarshal(hmRec.Body.Bytes(), &hmResp); err != nil {
		t.Fatalf("unmarshal health-matrix: %v", err)
	}
	states := map[string]string{}
	for _, p := range hmResp.Providers {
		states[p.Provider] = p.State
	}
	if states["openai"] != "degraded" {
		t.Errorf("expected openai state degraded, got %q", states["openai"])
	}
	if states["anthropic"] != "healthy" {
		t.Errorf("expected anthropic state healthy, got %q", states["anthropic"])
	}
}

func TestModelAndComboLiveProbeEndpoints(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Model string `json:"model"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Model == "bad-model" {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"message":"rate limit exceeded for sk-secret-123"}}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"ok","choices":[]}`))
	}))
	defer upstream.Close()

	store := &catalogFakeProviderStore{
		connections: []contracts.ProviderConnection{
			{
				ID:         "conn_1",
				ProviderID: "openai",
				Name:       "Local Mock",
				Protocol:   "openai",
				BaseURL:    upstream.URL + "/v1",
				Priority:   1,
				Enabled:    true,
				Models:     []string{"gpt-4o", "bad-model"},
			},
		},
		creds: map[string]contracts.UpstreamCredentials{
			"conn_1": {
				ConnectionID: "conn_1",
				ProviderID:   "openai",
				Protocol:     "openai",
				BaseURL:      upstream.URL + "/v1",
				AuthType:     "api_key",
				APIKey:       "sk-secret-123",
			},
		},
	}

	_, router := setupCatalogTestModule(t, store, upstream.Client())

	// 1. POST /api/models/test (success)
	testReq := httptest.NewRequest(http.MethodPost, "/api/models/test", bytes.NewBufferString(`{"providerId":"openai","modelId":"gpt-4o"}`))
	testRec := httptest.NewRecorder()
	router.ServeHTTP(testRec, testReq)
	if testRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /api/models/test, got %d: %s", testRec.Code, testRec.Body.String())
	}
	var singleResp struct {
		OK         bool   `json:"ok"`
		Status     string `json:"status"`
		StatusCode int    `json:"statusCode"`
	}
	if err := json.Unmarshal(testRec.Body.Bytes(), &singleResp); err != nil {
		t.Fatalf("unmarshal single test response: %v", err)
	}
	if !singleResp.OK || singleResp.Status != "ok" || singleResp.StatusCode != 200 {
		t.Errorf("unexpected single test response: %+v", singleResp)
	}

	// 2. POST /api/models/test (rate-limited + secret redaction)
	failReq := httptest.NewRequest(http.MethodPost, "/api/models/test", bytes.NewBufferString(`{"providerId":"openai","modelId":"bad-model"}`))
	failRec := httptest.NewRecorder()
	router.ServeHTTP(failRec, failReq)
	if failRec.Code != http.StatusOK {
		t.Fatalf("expected 200 probe envelope for upstream 429, got %d", failRec.Code)
	}
	if strings.Contains(failRec.Body.String(), "sk-secret-123") {
		t.Fatalf("secret leaked in model test error response: %s", failRec.Body.String())
	}
	var failResp struct {
		OK          bool   `json:"ok"`
		Status      string `json:"status"`
		RateLimited bool   `json:"rateLimited"`
	}
	_ = json.Unmarshal(failRec.Body.Bytes(), &failResp)
	if failResp.OK || failResp.Status != "error" || !failResp.RateLimited {
		t.Errorf("expected error + rateLimited=true, got %+v", failResp)
	}

	// 3. POST /api/models/test-all
	allReq := httptest.NewRequest(http.MethodPost, "/api/models/test-all", bytes.NewBufferString(`{"providerId":"openai","modelIds":["gpt-4o","bad-model"],"respectRateLimit":true}`))
	allRec := httptest.NewRecorder()
	router.ServeHTTP(allRec, allReq)
	if allRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /api/models/test-all, got %d: %s", allRec.Code, allRec.Body.String())
	}
	var allResp struct {
		Results map[string]struct {
			Status string `json:"status"`
		} `json:"results"`
	}
	_ = json.Unmarshal(allRec.Body.Bytes(), &allResp)
	if allResp.Results["gpt-4o"].Status != "ok" || allResp.Results["bad-model"].Status != "error" {
		t.Errorf("unexpected test-all results: %+v", allResp.Results)
	}

	// 4. Create a combo and test POST /api/combos/test (fallback from bad-model -> gpt-4o)
	createCombo := httptest.NewRequest(http.MethodPost, "/api/combos", bytes.NewBufferString(`{"name":"fallback-chain","targets":["openai/bad-model","openai/gpt-4o"]}`))
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createCombo)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating combo, got %d: %s", createRec.Code, createRec.Body.String())
	}

	comboTestReq := httptest.NewRequest(http.MethodPost, "/api/combos/test", bytes.NewBufferString(`{"comboName":"fallback-chain"}`))
	comboTestRec := httptest.NewRecorder()
	router.ServeHTTP(comboTestRec, comboTestReq)
	if comboTestRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /api/combos/test, got %d: %s", comboTestRec.Code, comboTestRec.Body.String())
	}
	var comboTestResp struct {
		ResolvedBy string `json:"resolvedBy"`
		Results    []struct {
			Target string `json:"target"`
			Status string `json:"status"`
		} `json:"results"`
	}
	_ = json.Unmarshal(comboTestRec.Body.Bytes(), &comboTestResp)
	if comboTestResp.ResolvedBy != "openai/gpt-4o" || len(comboTestResp.Results) != 2 {
		t.Errorf("expected resolvedBy=openai/gpt-4o with 2 results, got %+v", comboTestResp)
	}
}

func TestModelProbeOpenAIResponsesAndOAuthRefresh(t *testing.T) {
	var seenResponsesPath string
	var seenResponsesPayload map[string]any
	var oauthAttempts int

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/responses"):
			seenResponsesPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&seenResponsesPayload)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"resp_ok","output":[]}`))
		case strings.HasSuffix(r.URL.Path, "/chat/completions"):
			oauthAttempts++
			if r.Header.Get("Authorization") != "Bearer fresh-token" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":{"message":"expired token"}}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"chat_ok","choices":[]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":{"message":"unexpected path ` + r.URL.Path + `"}}`))
		}
	}))
	defer upstream.Close()

	refreshCalls := 0
	store := &catalogFakeProviderStore{
		connections: []contracts.ProviderConnection{
			{
				ID:         "conn_responses",
				ProviderID: "codex",
				Name:       "Codex Responses",
				Protocol:   "openai-responses",
				BaseURL:    upstream.URL + "/v1",
				Priority:   1,
				Enabled:    true,
			},
			{
				ID:         "conn_oauth",
				ProviderID: "qwen",
				Name:       "Qwen OAuth",
				Protocol:   "openai-chat",
				BaseURL:    upstream.URL + "/v1",
				Priority:   1,
				Enabled:    true,
			},
		},
		creds: map[string]contracts.UpstreamCredentials{
			"conn_responses": {
				ConnectionID: "conn_responses",
				ProviderID:   "codex",
				Protocol:     "openai-responses",
				BaseURL:      upstream.URL + "/v1",
				AuthType:     "api_key",
				APIKey:       "sk-responses",
			},
			"conn_oauth": {
				ConnectionID: "conn_oauth",
				ProviderID:   "qwen",
				Protocol:     "openai-chat",
				BaseURL:      upstream.URL + "/v1",
				AuthType:     "oauth",
				AccessToken:  "stale-token",
			},
		},
		refreshFn: func(connectionID string) (contracts.UpstreamCredentials, error) {
			refreshCalls++
			return contracts.UpstreamCredentials{
				ConnectionID: connectionID,
				ProviderID:   "qwen",
				Protocol:     "openai-chat",
				BaseURL:      upstream.URL + "/v1",
				AuthType:     "oauth",
				AccessToken:  "fresh-token",
			}, nil
		},
	}

	_, router := setupCatalogTestModule(t, store, upstream.Client())

	// 1. openai-responses protocol should target /v1/responses with "input" payload.
	req1 := httptest.NewRequest(http.MethodPost, "/api/models/test", bytes.NewBufferString(`{"providerId":"codex","modelId":"gpt-5-codex"}`))
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec1.Code, rec1.Body.String())
	}
	var resp1 struct {
		OK         bool   `json:"ok"`
		StatusCode int    `json:"statusCode"`
		Error      string `json:"error"`
	}
	_ = json.Unmarshal(rec1.Body.Bytes(), &resp1)
	if !resp1.OK || resp1.StatusCode != http.StatusOK {
		t.Fatalf("expected openai-responses probe to succeed, got %+v", resp1)
	}
	if seenResponsesPath != "/v1/responses" {
		t.Errorf("expected path /v1/responses, got %q", seenResponsesPath)
	}
	if seenResponsesPayload["model"] != "gpt-5-codex" || seenResponsesPayload["input"] == nil {
		t.Errorf("expected model + input payload for openai-responses, got %+v", seenResponsesPayload)
	}

	// 2. OAuth connection getting 401 should refresh once and succeed on retry.
	req2 := httptest.NewRequest(http.MethodPost, "/api/models/test", bytes.NewBufferString(`{"providerId":"qwen","modelId":"qwen3-coder-plus"}`))
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec2.Code, rec2.Body.String())
	}
	var resp2 struct {
		OK         bool   `json:"ok"`
		StatusCode int    `json:"statusCode"`
		Error      string `json:"error"`
	}
	_ = json.Unmarshal(rec2.Body.Bytes(), &resp2)
	if !resp2.OK || resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected OAuth retry probe to succeed, got %+v", resp2)
	}
	if refreshCalls != 1 || oauthAttempts != 2 {
		t.Errorf("expected 1 refresh and 2 upstream attempts, got refresh=%d attempts=%d", refreshCalls, oauthAttempts)
	}
}

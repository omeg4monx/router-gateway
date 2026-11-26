package combos_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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

func setupCatalogTestModule(t testing.TB, store contracts.ProviderStore, client *http.Client) (*combos.Module, chi.Router) {
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

type stubUsageRecorder struct {
	records []contracts.UsageRecord
}

func (s *stubUsageRecorder) Record(_ context.Context, rec contracts.UsageRecord) {
	s.records = append(s.records, rec)
}

func TestProbeAndIdentifyRecordUsageTelemetry(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/v1/chat/completions":
			var reqBody struct {
				Model string `json:"model"`
			}
			_ = json.NewDecoder(r.Body).Decode(&reqBody)
			if reqBody.Model == "rate-limited-model" {
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":{"message":"rate limit exceeded"}}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id":"chatcmpl-1",
				"model":"gpt-4o",
				"choices":[{"message":{"role":"assistant","content":"I am GPT-4o by OpenAI."}}],
				"usage":{"prompt_tokens":18,"completion_tokens":7,"prompt_tokens_details":{"cached_tokens":4}}
			}`))
		case r.URL.Path == "/v1/responses":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id":"resp_1",
				"model":"gpt-5-codex",
				"output":[{"type":"message","content":[{"type":"output_text","text":"ok"}]}],
				"usage":{"input_tokens":12,"output_tokens":5,"input_tokens_details":{"cached_tokens":2}}
			}`))
		case r.URL.Path == "/v1/messages":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id":"msg_1",
				"model":"claude-sonnet-4",
				"content":[{"type":"text","text":"ok"}],
				"usage":{"input_tokens":22,"output_tokens":9,"cache_read_input_tokens":6}
			}`))
		case strings.Contains(r.URL.Path, ":generateContent"):
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"candidates":[{"content":{"parts":[{"text":"ok"}]}}],
				"usageMetadata":{"promptTokenCount":15,"candidatesTokenCount":3,"cachedContentTokenCount":5}
			}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer upstream.Close()

	store := &catalogFakeProviderStore{
		connections: []contracts.ProviderConnection{
			{ID: "conn_oai", ProviderID: "omeg4mon", Protocol: "openai-chat", BaseURL: upstream.URL + "/v1", Priority: 1, Enabled: true},
			{ID: "conn_resp", ProviderID: "codex", Protocol: "openai-responses", BaseURL: upstream.URL + "/v1", Priority: 1, Enabled: true},
			{ID: "conn_ant", ProviderID: "anthropic", Protocol: "anthropic", BaseURL: upstream.URL + "/v1", Priority: 1, Enabled: true},
			{ID: "conn_gem", ProviderID: "gemini", Protocol: "gemini", BaseURL: upstream.URL + "/v1beta", Priority: 1, Enabled: true},
		},
		creds: map[string]contracts.UpstreamCredentials{
			"conn_oai":  {ConnectionID: "conn_oai", ProviderID: "omeg4mon", Protocol: "openai-chat", BaseURL: upstream.URL + "/v1", AuthType: "api_key", APIKey: "sk-oai"},
			"conn_resp": {ConnectionID: "conn_resp", ProviderID: "codex", Protocol: "openai-responses", BaseURL: upstream.URL + "/v1", AuthType: "api_key", APIKey: "sk-resp"},
			"conn_ant":  {ConnectionID: "conn_ant", ProviderID: "anthropic", Protocol: "anthropic", BaseURL: upstream.URL + "/v1", AuthType: "api_key", APIKey: "sk-ant"},
			"conn_gem":  {ConnectionID: "conn_gem", ProviderID: "gemini", Protocol: "gemini", BaseURL: upstream.URL + "/v1beta", AuthType: "api_key", APIKey: "sk-gem"},
		},
	}

	database, err := db.Open(t.TempDir() + "/usage-probe-test.sqlite")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	recorder := &stubUsageRecorder{}
	mod := combos.New(combos.Options{
		DB:         database,
		Providers:  store,
		HTTPClient: upstream.Client(),
		Usage:      recorder,
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

	calls := []struct {
		path string
		body string
	}{
		{"/api/models/test", `{"providerId":"omeg4mon","modelId":"gpt-4o"}`},
		{"/api/models/test", `{"providerId":"codex","modelId":"gpt-5-codex"}`},
		{"/api/models/test", `{"providerId":"anthropic","modelId":"claude-sonnet-4"}`},
		{"/api/models/test", `{"providerId":"gemini","modelId":"gemini-2.5-flash"}`},
		{"/api/models/test", `{"providerId":"omeg4mon","modelId":"rate-limited-model"}`},
		{"/api/models/identify", `{"providerId":"omeg4mon","modelId":"gpt-4o"}`},
	}

	for _, c := range calls {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, c.path, bytes.NewBufferString(c.body)))
		if rec.Code != http.StatusOK {
			t.Fatalf("POST %s returned %d: %s", c.path, rec.Code, rec.Body.String())
		}
	}

	if len(recorder.records) != len(calls) {
		t.Fatalf("recorded %d usage entries, want %d", len(recorder.records), len(calls))
	}

	// OpenAI Chat probe
	if r := recorder.records[0]; r.ProviderID != "omeg4mon" || r.Model != "gpt-4o" || r.InputTokens != 18 || r.OutputTokens != 7 || r.CachedTokens != 4 || r.StatusCode != 200 || r.KeyID != "admin-probe" || r.RequestID == "" {
		t.Errorf("unexpected openai-chat probe usage record: %+v", r)
	}
	// OpenAI Responses probe
	if r := recorder.records[1]; r.ProviderID != "codex" || r.Model != "gpt-5-codex" || r.InputTokens != 12 || r.OutputTokens != 5 || r.CachedTokens != 2 || r.StatusCode != 200 {
		t.Errorf("unexpected openai-responses probe usage record: %+v", r)
	}
	// Anthropic probe
	if r := recorder.records[2]; r.ProviderID != "anthropic" || r.Model != "claude-sonnet-4" || r.InputTokens != 22 || r.OutputTokens != 9 || r.CachedTokens != 6 || r.StatusCode != 200 {
		t.Errorf("unexpected anthropic probe usage record: %+v", r)
	}
	// Gemini probe
	if r := recorder.records[3]; r.ProviderID != "gemini" || r.Model != "gemini-2.5-flash" || r.InputTokens != 15 || r.OutputTokens != 3 || r.CachedTokens != 5 || r.StatusCode != 200 {
		t.Errorf("unexpected gemini probe usage record: %+v", r)
	}
	// Rate-limited error probe
	if r := recorder.records[4]; r.ProviderID != "omeg4mon" || r.Model != "rate-limited-model" || r.StatusCode != 429 || r.ErrorSummary == "" {
		t.Errorf("unexpected error probe usage record: %+v", r)
	}
	// Identify probe
	if r := recorder.records[5]; r.ProviderID != "omeg4mon" || r.Model != "gpt-4o" || r.InputTokens != 18 || r.OutputTokens != 7 || r.CachedTokens != 4 || r.StatusCode != 200 || r.KeyID != "admin-probe" {
		t.Errorf("unexpected identify probe usage record: %+v", r)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestKeylessFreeProviderModelProbeAndIdentify(t *testing.T) {
	var capturedURLs []string
	var capturedAuthHeaders []string

	mockClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			capturedURLs = append(capturedURLs, req.URL.String())
			capturedAuthHeaders = append(capturedAuthHeaders, req.Header.Get("Authorization"))
			body := `{"id":"chatcmpl-keyless","model":"venice-uncensored-1-2","choices":[{"message":{"role":"assistant","content":"I am venice-uncensored-1-2."}}],"usage":{"prompt_tokens":5,"completion_tokens":3}}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		}),
	}

	database, err := db.Open(t.TempDir() + "/keyless-probe-test.sqlite")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	emptyStore := &catalogFakeProviderStore{}
	mod := combos.New(combos.Options{
		DB:         database,
		Providers:  emptyStore,
		HTTPClient: mockClient,
		Now:        func() time.Time { return time.Date(2026, 11, 26, 0, 0, 0, 0, time.UTC) },
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

	cases := []struct {
		providerID string
		modelID    string
		wantURL    string
	}{
		{"venice", "venice-uncensored-1-2", "https://outerface.venice.ai/api/inference/chat"},
		{"opencode", "muse-wonder-1.0", "https://opencode.ai/zen/v1/chat/completions"},
		{"ovh", "Qwen3.6-27B", "https://oai.endpoints.kepler.ai.cloud.ovh.net/v1/chat/completions"},
	}

	for _, tc := range cases {
		rec := httptest.NewRecorder()
		payload := `{"providerId":"` + tc.providerID + `","modelId":"` + tc.modelID + `"}`
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/models/test", bytes.NewBufferString(payload)))
		if rec.Code != http.StatusOK {
			t.Fatalf("POST /api/models/test (%s) status = %d: %s", tc.providerID, rec.Code, rec.Body.String())
		}
		var res combos.ModelTestResult
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("decode test result (%s): %v", tc.providerID, err)
		}
		if !res.OK || res.Status != "ok" {
			t.Fatalf("expected ok probe for keyless provider %q, got %+v", tc.providerID, res)
		}
	}

	// Identify probe must also work for keyless providers without a saved connection.
	idRec := httptest.NewRecorder()
	router.ServeHTTP(idRec, httptest.NewRequest(http.MethodPost, "/api/models/identify", bytes.NewBufferString(`{"providerId":"venice","modelId":"venice-uncensored-1-2"}`)))
	if idRec.Code != http.StatusOK {
		t.Fatalf("POST /api/models/identify (venice) status = %d: %s", idRec.Code, idRec.Body.String())
	}
	var idRes combos.ModelIdentityResult
	if err := json.Unmarshal(idRec.Body.Bytes(), &idRes); err != nil {
		t.Fatalf("decode identify result: %v", err)
	}
	if !idRes.OK || idRes.Verdict == "error" {
		t.Fatalf("expected ok identify probe for keyless provider venice, got %+v", idRes)
	}

	wantURLs := []string{
		"https://outerface.venice.ai/api/inference/chat",
		"https://opencode.ai/zen/v1/chat/completions",
		"https://oai.endpoints.kepler.ai.cloud.ovh.net/v1/chat/completions",
		"https://outerface.venice.ai/api/inference/chat",
	}
	if len(capturedURLs) != len(wantURLs) {
		t.Fatalf("captured %d upstream requests, want %d (%v)", len(capturedURLs), len(wantURLs), capturedURLs)
	}
	for i, want := range wantURLs {
		if capturedURLs[i] != want {
			t.Errorf("request[%d] URL = %q, want %q", i, capturedURLs[i], want)
		}
		if capturedAuthHeaders[i] != "" {
			t.Errorf("request[%d] unexpected Authorization header = %q", i, capturedAuthHeaders[i])
		}
	}

	// Non-keyless provider with zero connections must still report no enabled connection.
	nonKeylessRec := httptest.NewRecorder()
	router.ServeHTTP(nonKeylessRec, httptest.NewRequest(http.MethodPost, "/api/models/test", bytes.NewBufferString(`{"providerId":"openai","modelId":"gpt-4o"}`)))
	var nonKeylessRes combos.ModelTestResult
	if err := json.Unmarshal(nonKeylessRec.Body.Bytes(), &nonKeylessRes); err != nil {
		t.Fatalf("decode non-keyless result: %v", err)
	}
	if nonKeylessRes.OK || !strings.Contains(nonKeylessRes.Error, "no enabled connection") {
		t.Fatalf("expected no enabled connection error for openai without connections, got %+v", nonKeylessRes)
	}
}

package providers_test

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/providers"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

var (
	_ lifecycle.Module        = (*providers.Module)(nil)
	_ contracts.ProviderStore = (*providers.Module)(nil)
)

type stubSystemModule struct{}

func (stubSystemModule) Name() string                         { return "system" }
func (stubSystemModule) Dependencies() []string               { return nil }
func (stubSystemModule) Migrations() fs.FS                    { return nil }
func (stubSystemModule) RegisterRoutes(lifecycle.RouteMounts) {}
func (stubSystemModule) Start(context.Context) error          { return nil }
func (stubSystemModule) Stop(context.Context) error           { return nil }

type stubAdminAuth struct {
	subject string
}

func (a *stubAdminAuth) AuthenticateAdmin(r *http.Request) (string, error) {
	if r.Header.Get("Authorization") == "Bearer test-admin-token" {
		return a.subject, nil
	}
	return "", errors.New("unauthorized")
}

func (a *stubAdminAuth) RequireLogin(context.Context) bool {
	return true
}

type customTestAdapter struct{}

func (customTestAdapter) ID() string              { return "custom-plugin" }
func (customTestAdapter) DisplayName() string     { return "Custom Plugin" }
func (customTestAdapter) Category() string        { return "api_key" }
func (customTestAdapter) Protocol() string        { return "openai-chat" }
func (customTestAdapter) AuthType() string        { return "api_key" }
func (customTestAdapter) DefaultBaseURL() string  { return "https://plugin.example.test/v1" }
func (customTestAdapter) DefaultModels() []string { return []string{"plugin-model-1"} }
func (customTestAdapter) OAuthConfig() *providers.OAuthConfig {
	return nil
}
func (customTestAdapter) BuildTestRequest(ctx context.Context, creds contracts.UpstreamCredentials) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(creds.BaseURL, "/")+"/models", http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+creds.APIKey)
	return req, nil
}

type testHarness struct {
	db     *db.DB
	vault  *crypto.Vault
	bus    *eventbus.Bus
	reg    *lifecycle.Registry
	server *httpserver.Server
	mod    *providers.Module
	auth   *stubAdminAuth
}

func newTestHarness(t *testing.T, auth *stubAdminAuth, overrides map[string]providers.OAuthConfig) *testHarness {
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
		DB:             database,
		Vault:          vault,
		HTTPClient:     &http.Client{Timeout: 5 * time.Second},
		Bus:            bus,
		OAuthOverrides: overrides,
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

	srvOpts := httpserver.Options{
		Config: cfg,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	if auth != nil {
		srvOpts.Auth = auth
	}
	srv := httpserver.New(srvOpts)
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

func (h *testHarness) doRequest(t *testing.T, method, path, body string, authed bool) *httptest.ResponseRecorder {
	t.Helper()
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, bodyReader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if authed {
		req.Header.Set("Authorization", "Bearer test-admin-token")
	}
	rec := httptest.NewRecorder()
	h.server.Handler().ServeHTTP(rec, req)
	return rec
}

func TestCatalogAndAdapterRegistry(t *testing.T) {
	h := newTestHarness(t, &stubAdminAuth{subject: "admin"}, nil)

	// Unauthenticated admin route must return 401.
	unauthRec := h.doRequest(t, http.MethodGet, "/api/providers/catalog", "", false)
	if unauthRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unauthenticated catalog request, got %d: %s", unauthRec.Code, unauthRec.Body.String())
	}

	// Register a custom pluggable adapter at runtime.
	h.mod.RegisterAdapter(customTestAdapter{})

	rec := h.doRequest(t, http.MethodGet, "/api/providers/catalog", "", true)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/providers/catalog status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var catalog struct {
		Adapters []struct {
			ID             string   `json:"id"`
			DisplayName    string   `json:"display_name"`
			Category       string   `json:"category"`
			Protocol       string   `json:"protocol"`
			AuthType       string   `json:"auth_type"`
			DefaultBaseURL string   `json:"default_base_url"`
			DefaultModels  []string `json:"default_models"`
			OAuth          *struct {
				AuthorizeURL string   `json:"authorize_url"`
				TokenURL     string   `json:"token_url"`
				Scopes       []string `json:"scopes"`
			} `json:"oauth"`
		} `json:"adapters"`
		CustomNodes []providers.ProviderNode `json:"custom_nodes"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &catalog); err != nil {
		t.Fatalf("decode catalog: %v", err)
	}

	expectedIDs := []string{
		"openai",
		"anthropic",
		"openai-responses",
		"gemini",
		"openrouter",
		"openai-compatible",
		"anthropic-compatible",
		"claude",
		"codex",
		"gemini-oauth",
		"custom-plugin",
	}
	byID := make(map[string]struct {
		Protocol string
		AuthType string
		OAuthURL string
	})
	for _, a := range catalog.Adapters {
		oauthURL := ""
		if a.OAuth != nil {
			oauthURL = a.OAuth.AuthorizeURL
		}
		byID[a.ID] = struct {
			Protocol string
			AuthType string
			OAuthURL string
		}{Protocol: a.Protocol, AuthType: a.AuthType, OAuthURL: oauthURL}
	}

	for _, id := range expectedIDs {
		if _, ok := byID[id]; !ok {
			t.Errorf("catalog missing expected adapter %q", id)
		}
	}

	if byID["claude"].Protocol != "anthropic" || byID["claude"].AuthType != "oauth" || byID["claude"].OAuthURL != "https://claude.ai/oauth/authorize" {
		t.Errorf("unexpected claude adapter metadata: %+v", byID["claude"])
	}
	if byID["codex"].Protocol != "openai-responses" || byID["codex"].AuthType != "oauth" || byID["codex"].OAuthURL != "https://auth.openai.com/oauth/authorize" {
		t.Errorf("unexpected codex adapter metadata: %+v", byID["codex"])
	}
	if byID["gemini-oauth"].Protocol != "gemini" || byID["gemini-oauth"].AuthType != "oauth" || byID["gemini-oauth"].OAuthURL != "https://accounts.google.com/o/oauth2/v2/auth" {
		t.Errorf("unexpected gemini-oauth adapter metadata: %+v", byID["gemini-oauth"])
	}
}

func TestCustomProviderNodesCRUD(t *testing.T) {
	h := newTestHarness(t, nil, nil)
	events, unsub := h.bus.Subscribe(providers.TopicProvidersChanged, 16)
	defer unsub()

	// Create custom node.
	createRec := h.doRequest(t, http.MethodPost, "/api/providers/nodes", `{
		"id": "vllm-local",
		"name": "Local vLLM Cluster",
		"protocol": "openai-chat",
		"base_url": "http://127.0.0.1:8000/v1",
		"models": ["llama-3.3-70b", "qwen-2.5-coder"],
		"headers": {"X-Cluster": "home-lab"}
	}`, false)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("POST /api/providers/nodes status = %d, want 201: %s", createRec.Code, createRec.Body.String())
	}

	var created providers.ProviderNode
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created node: %v", err)
	}
	if created.ID != "vllm-local" || created.Protocol != "openai-chat" || len(created.Models) != 2 {
		t.Fatalf("unexpected created node: %+v", created)
	}

	// Duplicate ID must return 409.
	dupRec := h.doRequest(t, http.MethodPost, "/api/providers/nodes", `{
		"id": "vllm-local",
		"name": "Duplicate",
		"protocol": "openai-chat",
		"base_url": "http://127.0.0.1:8000/v1"
	}`, false)
	if dupRec.Code != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate node id, got %d: %s", dupRec.Code, dupRec.Body.String())
	}

	// List nodes & check catalog includes it.
	listRec := h.doRequest(t, http.MethodGet, "/api/providers/nodes", "", false)
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /api/providers/nodes status = %d: %s", listRec.Code, listRec.Body.String())
	}
	var listResp struct {
		Nodes []providers.ProviderNode `json:"nodes"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("decode nodes list: %v", err)
	}
	if len(listResp.Nodes) != 1 || listResp.Nodes[0].ID != "vllm-local" {
		t.Fatalf("unexpected nodes list: %+v", listResp.Nodes)
	}

	// Create a connection backed by this custom provider node.
	connRec := h.doRequest(t, http.MethodPost, "/api/providers/connections", `{
		"provider_id": "vllm-local",
		"name": "vLLM Primary",
		"api_key": "sk-vllm-super-secret-token"
	}`, false)
	if connRec.Code != http.StatusCreated {
		t.Fatalf("POST /api/providers/connections for custom node status = %d: %s", connRec.Code, connRec.Body.String())
	}

	// Delete node.
	delRec := h.doRequest(t, http.MethodDelete, "/api/providers/nodes/vllm-local", "", false)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /api/providers/nodes/vllm-local status = %d: %s", delRec.Code, delRec.Body.String())
	}

	// Deleting again returns 404.
	delAgain := h.doRequest(t, http.MethodDelete, "/api/providers/nodes/vllm-local", "", false)
	if delAgain.Code != http.StatusNotFound {
		t.Fatalf("expected 404 deleting missing node, got %d", delAgain.Code)
	}

	// Verify events published without any secrets.
	received := 0
	for i := 0; i < 3; i++ {
		select {
		case ev := <-events:
			received++
			raw, _ := json.Marshal(ev.Payload)
			if strings.Contains(string(raw), "sk-vllm-super-secret-token") {
				t.Fatalf("eventbus payload leaked plaintext API key: %s", string(raw))
			}
		case <-time.After(time.Second):
			t.Fatalf("expected 3 providers.changed events, got %d", received)
		}
	}
}

func TestConnectionsCRUDAndSecretEncryptionHygiene(t *testing.T) {
	h := newTestHarness(t, nil, nil)
	events, unsub := h.bus.Subscribe(providers.TopicProvidersChanged, 16)
	defer unsub()

	const (
		secretKey1 = "sk-live-openai-very-secret-0001"
		secretKey2 = "sk-live-openai-rotated-secret-0002"
	)

	// Create connection 1 (priority 20).
	rec1 := h.doRequest(t, http.MethodPost, "/api/providers/connections", `{
		"provider_id": "openai",
		"name": "OpenAI Primary",
		"api_key": "`+secretKey1+`",
		"priority": 20,
		"models": ["gpt-4o"],
		"headers": {"OpenAI-Organization": "org-test"}
	}`, false)
	if rec1.Code != http.StatusCreated {
		t.Fatalf("POST connection status = %d, want 201: %s", rec1.Code, rec1.Body.String())
	}
	if strings.Contains(rec1.Body.String(), secretKey1) || strings.Contains(rec1.Body.String(), "v1:") {
		t.Fatalf("POST response leaked secret or ciphertext: %s", rec1.Body.String())
	}

	var conn1 providers.ConnectionSummary
	if err := json.Unmarshal(rec1.Body.Bytes(), &conn1); err != nil {
		t.Fatalf("decode conn1: %v", err)
	}
	if conn1.ID == "" || conn1.ProviderID != "openai" || conn1.Protocol != "openai-chat" || conn1.BaseURL != "https://api.openai.com/v1" || !conn1.Enabled || !conn1.HasCredentials {
		t.Fatalf("unexpected conn1 summary: %+v", conn1)
	}

	// Verify SQLite stores AES-256-GCM ciphertext with "v1:" prefix and NEVER plaintext.
	var enc string
	if err := h.db.Reader().QueryRow(`SELECT credentials_enc FROM provider_connections WHERE id = ?`, conn1.ID).Scan(&enc); err != nil {
		t.Fatalf("query credentials_enc: %v", err)
	}
	if !strings.HasPrefix(enc, "v1:") {
		t.Fatalf("credentials_enc must start with v1:, got %q", enc)
	}
	if strings.Contains(enc, secretKey1) {
		t.Fatalf("credentials_enc must not contain plaintext api_key")
	}
	plain, err := h.vault.Decrypt(enc)
	if err != nil {
		t.Fatalf("vault.Decrypt(credentials_enc): %v", err)
	}
	if !strings.Contains(string(plain), secretKey1) {
		t.Fatalf("decrypted credentials missing expected api_key")
	}

	// Create connection 2 (priority 10, enabled) and connection 3 (disabled).
	rec2 := h.doRequest(t, http.MethodPost, "/api/providers/connections", `{
		"provider_id": "anthropic",
		"name": "Anthropic High Priority",
		"api_key": "sk-ant-secret-key-999",
		"priority": 10
	}`, false)
	if rec2.Code != http.StatusCreated {
		t.Fatalf("POST conn2 status = %d: %s", rec2.Code, rec2.Body.String())
	}
	var conn2 providers.ConnectionSummary
	_ = json.Unmarshal(rec2.Body.Bytes(), &conn2)

	rec3 := h.doRequest(t, http.MethodPost, "/api/providers/connections", `{
		"provider_id": "gemini",
		"name": "Gemini Disabled",
		"api_key": "AIza-secret-disabled",
		"enabled": false
	}`, false)
	if rec3.Code != http.StatusCreated {
		t.Fatalf("POST conn3 status = %d: %s", rec3.Code, rec3.Body.String())
	}

	// ListEnabledConnections must exclude disabled and sort by priority ascending.
	enabledConns, err := h.mod.ListEnabledConnections(context.Background())
	if err != nil {
		t.Fatalf("ListEnabledConnections: %v", err)
	}
	if len(enabledConns) != 2 {
		t.Fatalf("expected 2 enabled connections, got %d", len(enabledConns))
	}
	if enabledConns[0].ID != conn2.ID || enabledConns[1].ID != conn1.ID {
		t.Fatalf("expected priority order [%s, %s], got [%s, %s]", conn2.ID, conn1.ID, enabledConns[0].ID, enabledConns[1].ID)
	}

	// GET list and GET single must never leak plaintext secrets or ciphertext.
	listRec := h.doRequest(t, http.MethodGet, "/api/providers/connections", "", false)
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET connections status = %d", listRec.Code)
	}
	for _, forbidden := range []string{secretKey1, "sk-ant-secret-key-999", "AIza-secret-disabled", "v1:"} {
		if strings.Contains(listRec.Body.String(), forbidden) {
			t.Fatalf("GET /api/providers/connections leaked %q", forbidden)
		}
	}

	getRec := h.doRequest(t, http.MethodGet, "/api/providers/connections/"+conn1.ID, "", false)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET connection status = %d", getRec.Code)
	}
	if strings.Contains(getRec.Body.String(), secretKey1) || strings.Contains(getRec.Body.String(), "v1:") {
		t.Fatalf("GET single connection leaked secret: %s", getRec.Body.String())
	}

	// PATCH connection rotates api_key and updates priority.
	patchRec := h.doRequest(t, http.MethodPatch, "/api/providers/connections/"+conn1.ID, `{
		"name": "OpenAI Rotated",
		"api_key": "`+secretKey2+`",
		"priority": 5
	}`, false)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("PATCH connection status = %d: %s", patchRec.Code, patchRec.Body.String())
	}
	if strings.Contains(patchRec.Body.String(), secretKey2) || strings.Contains(patchRec.Body.String(), "v1:") {
		t.Fatalf("PATCH response leaked secret: %s", patchRec.Body.String())
	}

	// ResolveCredentials returns decrypted rotated secret.
	resolved, err := h.mod.ResolveCredentials(context.Background(), conn1.ID)
	if err != nil {
		t.Fatalf("ResolveCredentials: %v", err)
	}
	if resolved.APIKey != secretKey2 || resolved.ExtraHeaders["OpenAI-Organization"] != "org-test" {
		t.Fatalf("unexpected resolved credentials: %+v", resolved)
	}

	// ResolveCredentials for missing ID returns contracts.ErrConnectionNotFound.
	if _, err := h.mod.ResolveCredentials(context.Background(), "conn_missing"); !errors.Is(err, contracts.ErrConnectionNotFound) {
		t.Fatalf("expected contracts.ErrConnectionNotFound, got %v", err)
	}

	// RefreshOAuthToken on an API-key connection returns contracts.ErrOAuthRefreshUnsupported.
	if _, err := h.mod.RefreshOAuthToken(context.Background(), conn1.ID); !errors.Is(err, contracts.ErrOAuthRefreshUnsupported) {
		t.Fatalf("expected contracts.ErrOAuthRefreshUnsupported, got %v", err)
	}

	// DELETE connection.
	delRec := h.doRequest(t, http.MethodDelete, "/api/providers/connections/"+conn1.ID, "", false)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("DELETE connection status = %d", delRec.Code)
	}
	if _, err := h.mod.ResolveCredentials(context.Background(), conn1.ID); !errors.Is(err, contracts.ErrConnectionNotFound) {
		t.Fatalf("expected deleted connection to be not found, got %v", err)
	}

	// Verify none of the emitted events contain plaintext secrets or ciphertext.
	for i := 0; i < 5; i++ {
		select {
		case ev := <-events:
			raw, _ := json.Marshal(ev.Payload)
			for _, forbidden := range []string{secretKey1, secretKey2, "sk-ant-secret-key-999", "AIza-secret-disabled", "v1:"} {
				if strings.Contains(string(raw), forbidden) {
					t.Fatalf("eventbus payload leaked %q: %s", forbidden, string(raw))
				}
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for event %d", i)
		}
	}
}

func TestConnectionQuickTest(t *testing.T) {
	var seenAuth atomic.Value
	var seenVersion atomic.Value
	var seenExtra atomic.Value
	var returnStatus atomic.Int32
	returnStatus.Store(http.StatusOK)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenAuth.Store(r.Header.Get("x-api-key"))
		seenVersion.Store(r.Header.Get("anthropic-version"))
		seenExtra.Store(r.Header.Get("X-Custom-Header"))
		code := int(returnStatus.Load())
		w.WriteHeader(code)
		if code == http.StatusOK {
			_, _ = w.Write([]byte(`{"data":[{"id":"claude-sonnet-4-20250514"}]}`))
		} else {
			_, _ = w.Write([]byte(`{"error":{"message":"invalid api key"}}`))
		}
	}))
	defer upstream.Close()

	h := newTestHarness(t, nil, nil)
	const apiKey = "sk-ant-quicktest-secret"

	createRec := h.doRequest(t, http.MethodPost, "/api/providers/connections", `{
		"provider_id": "anthropic",
		"name": "Anthropic Test",
		"base_url": "`+upstream.URL+`",
		"api_key": "`+apiKey+`",
		"headers": {"X-Custom-Header": "ping-ok"}
	}`, false)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create conn status = %d: %s", createRec.Code, createRec.Body.String())
	}
	var conn providers.ConnectionSummary
	_ = json.Unmarshal(createRec.Body.Bytes(), &conn)

	// Successful quick test.
	testRec := h.doRequest(t, http.MethodPost, "/api/providers/connections/"+conn.ID+"/test", "", false)
	if testRec.Code != http.StatusOK {
		t.Fatalf("POST /test status = %d: %s", testRec.Code, testRec.Body.String())
	}
	if strings.Contains(testRec.Body.String(), apiKey) {
		t.Fatalf("POST /test leaked api_key")
	}

	var testResp struct {
		OK         bool   `json:"ok"`
		StatusCode int    `json:"status_code"`
		LatencyMS  int64  `json:"latency_ms"`
		Error      string `json:"error"`
	}
	if err := json.Unmarshal(testRec.Body.Bytes(), &testResp); err != nil {
		t.Fatalf("decode testResp: %v", err)
	}
	if !testResp.OK || testResp.StatusCode != http.StatusOK || testResp.Error != "" {
		t.Fatalf("unexpected success testResp: %+v", testResp)
	}
	if seenAuth.Load() != apiKey || seenVersion.Load() != "2023-06-01" || seenExtra.Load() != "ping-ok" {
		t.Fatalf("unexpected upstream headers: auth=%v version=%v extra=%v", seenAuth.Load(), seenVersion.Load(), seenExtra.Load())
	}

	// Verify persisted last_tested_at and last_test_status.
	getRec := h.doRequest(t, http.MethodGet, "/api/providers/connections/"+conn.ID, "", false)
	var updated providers.ConnectionSummary
	_ = json.Unmarshal(getRec.Body.Bytes(), &updated)
	if updated.LastTestedAt == nil || updated.LastTestStatus != "ok" || updated.LastTestError != "" {
		t.Fatalf("unexpected persisted test metadata: %+v", updated)
	}

	// Upstream 401 failure test.
	returnStatus.Store(http.StatusUnauthorized)
	failRec := h.doRequest(t, http.MethodPost, "/api/providers/connections/"+conn.ID+"/test", "", false)
	if failRec.Code != http.StatusOK {
		t.Fatalf("POST /test (upstream 401) status = %d: %s", failRec.Code, failRec.Body.String())
	}
	var failResp struct {
		OK         bool   `json:"ok"`
		StatusCode int    `json:"status_code"`
		Error      string `json:"error"`
	}
	_ = json.Unmarshal(failRec.Body.Bytes(), &failResp)
	if failResp.OK || failResp.StatusCode != http.StatusUnauthorized || failResp.Error == "" {
		t.Fatalf("unexpected failure testResp: %+v", failResp)
	}
}

func TestHeadlessOAuthFlowAndCallbackXSSProtection(t *testing.T) {
	const (
		expectedClientID     = "test-claude-client-id"
		expectedClientSecret = "test-claude-client-secret-xyz"
		expectedAuthCode     = "auth-code-987654321"
		issuedAccessToken    = "oauth-access-token-secret-111"
		issuedRefreshToken   = "oauth-refresh-token-secret-222"
	)

	var (
		capturedCodeChallenge string
		mu                    sync.Mutex
		tokenRequestCount     int
	)

	mockTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}
		mu.Lock()
		tokenRequestCount++
		wantChallenge := capturedCodeChallenge
		mu.Unlock()

		if r.FormValue("grant_type") != "authorization_code" {
			http.Error(w, "invalid grant_type", http.StatusBadRequest)
			return
		}
		if r.FormValue("code") != expectedAuthCode {
			http.Error(w, "invalid code", http.StatusBadRequest)
			return
		}
		if r.FormValue("client_id") != expectedClientID {
			http.Error(w, "invalid client_id", http.StatusUnauthorized)
			return
		}
		if r.FormValue("client_secret") != expectedClientSecret {
			http.Error(w, "invalid client_secret", http.StatusUnauthorized)
			return
		}
		verifier := r.FormValue("code_verifier")
		sum := sha256.Sum256([]byte(verifier))
		gotChallenge := base64.RawURLEncoding.EncodeToString(sum[:])
		if gotChallenge != wantChallenge {
			http.Error(w, "pkce mismatch", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  issuedAccessToken,
			"refresh_token": issuedRefreshToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
			"scope":         "user:inference",
		})
	}))
	defer mockTokenServer.Close()

	overrides := map[string]providers.OAuthConfig{
		"claude": {
			AuthorizeURL:       "https://claude.ai/oauth/authorize",
			TokenURL:           mockTokenServer.URL + "/token",
			ClientID:           expectedClientID,
			ClientSecret:       expectedClientSecret,
			DefaultRedirectURI: "http://localhost:20128/api/oauth/callback",
		},
	}

	h := newTestHarness(t, &stubAdminAuth{subject: "admin"}, overrides)
	events, unsub := h.bus.Subscribe(providers.TopicProvidersChanged, 16)
	defer unsub()

	// Step 1: Start OAuth session.
	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/claude/start", `{
		"connection_name": "Claude Work OAuth"
	}`, true)
	if startRec.Code != http.StatusOK {
		t.Fatalf("POST /oauth/claude/start status = %d: %s", startRec.Code, startRec.Body.String())
	}
	if strings.Contains(startRec.Body.String(), expectedClientSecret) {
		t.Fatalf("start response leaked client_secret")
	}

	var startResp struct {
		State        string `json:"state"`
		AuthorizeURL string `json:"authorize_url"`
		RedirectURI  string `json:"redirect_uri"`
		ExpiresAt    string `json:"expires_at"`
	}
	if err := json.Unmarshal(startRec.Body.Bytes(), &startResp); err != nil {
		t.Fatalf("decode startResp: %v", err)
	}
	if startResp.State == "" || startResp.AuthorizeURL == "" || startResp.RedirectURI == "" || startResp.ExpiresAt == "" {
		t.Fatalf("incomplete startResp: %+v", startResp)
	}

	parsedAuthURL, err := url.Parse(startResp.AuthorizeURL)
	if err != nil {
		t.Fatalf("parse authorize_url: %v", err)
	}
	q := parsedAuthURL.Query()
	if q.Get("code_challenge_method") != "S256" || q.Get("code_challenge") == "" || q.Get("state") != startResp.State || q.Get("client_id") != expectedClientID {
		t.Fatalf("unexpected authorize_url query params: %v", q)
	}
	mu.Lock()
	capturedCodeChallenge = q.Get("code_challenge")
	mu.Unlock()

	// Verify code_verifier and client_secret are encrypted in oauth_sessions.
	var verifierEnc, secretEnc string
	if err := h.db.Reader().QueryRow(`SELECT code_verifier_enc, client_secret_enc FROM oauth_sessions WHERE state = ?`, startResp.State).Scan(&verifierEnc, &secretEnc); err != nil {
		t.Fatalf("query oauth_sessions: %v", err)
	}
	if !strings.HasPrefix(verifierEnc, "v1:") || !strings.HasPrefix(secretEnc, "v1:") || strings.Contains(secretEnc, expectedClientSecret) {
		t.Fatalf("oauth_sessions secrets must be vault-encrypted, got verifier=%q secret=%q", verifierEnc, secretEnc)
	}

	// Step 2: Public GET /api/oauth/callback must work without admin auth and prevent reflected XSS.
	xssRec := h.doRequest(t, http.MethodGet, `/api/oauth/callback?code=%3Cscript%3Ealert(1)%3C/script%3E&state=%22%3E%3Cimg%20src=x%20onerror=alert(1)%3E`, "", false)
	if xssRec.Code != http.StatusOK {
		t.Fatalf("GET /api/oauth/callback status = %d", xssRec.Code)
	}
	htmlBody := xssRec.Body.String()
	if strings.Contains(htmlBody, "<script>alert(1)</script>") || strings.Contains(htmlBody, `<img src=x onerror=alert(1)>`) {
		t.Fatalf("GET /api/oauth/callback reflected unescaped HTML/XSS: %s", htmlBody)
	}
	if !strings.Contains(htmlBody, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Fatalf("expected HTML-escaped payload in callback output: %s", htmlBody)
	}

	// Step 3: Headless-safe exchange via pasted callback_url.
	pastedCallback := "http://localhost:20128/api/oauth/callback?code=" + url.QueryEscape(expectedAuthCode) + "&state=" + url.QueryEscape(startResp.State)
	exchangeBody, _ := json.Marshal(map[string]string{
		"callback_url": pastedCallback,
	})
	exchRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/claude/exchange", string(exchangeBody), true)
	if exchRec.Code != http.StatusCreated {
		t.Fatalf("POST /oauth/claude/exchange status = %d, want 201: %s", exchRec.Code, exchRec.Body.String())
	}
	for _, forbidden := range []string{issuedAccessToken, issuedRefreshToken, expectedClientSecret, "v1:"} {
		if strings.Contains(exchRec.Body.String(), forbidden) {
			t.Fatalf("exchange response leaked %q: %s", forbidden, exchRec.Body.String())
		}
	}

	var conn providers.ConnectionSummary
	if err := json.Unmarshal(exchRec.Body.Bytes(), &conn); err != nil {
		t.Fatalf("decode exchanged connection: %v", err)
	}
	if conn.ProviderID != "claude" || conn.AuthType != "oauth" || conn.Name != "Claude Work OAuth" || conn.TokenExpiresAt == nil {
		t.Fatalf("unexpected exchanged connection: %+v", conn)
	}

	// Verify SQLite stores encrypted credentials and deleted consumed oauth_sessions row.
	var connEnc string
	if err := h.db.Reader().QueryRow(`SELECT credentials_enc FROM provider_connections WHERE id = ?`, conn.ID).Scan(&connEnc); err != nil {
		t.Fatalf("query oauth connection credentials_enc: %v", err)
	}
	if !strings.HasPrefix(connEnc, "v1:") || strings.Contains(connEnc, issuedAccessToken) || strings.Contains(connEnc, issuedRefreshToken) {
		t.Fatalf("oauth credentials_enc not properly encrypted")
	}
	var remainingSessions int
	if err := h.db.Reader().QueryRow(`SELECT COUNT(*) FROM oauth_sessions WHERE state = ?`, startResp.State).Scan(&remainingSessions); err != nil {
		t.Fatalf("count oauth_sessions: %v", err)
	}
	if remainingSessions != 0 {
		t.Fatalf("expected consumed oauth_session to be deleted, found %d", remainingSessions)
	}

	// Replaying the same callback_url must fail with 400.
	replayRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/claude/exchange", string(exchangeBody), true)
	if replayRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on replayed oauth exchange, got %d: %s", replayRec.Code, replayRec.Body.String())
	}

	// ResolveCredentials returns the decrypted OAuth access token.
	creds, err := h.mod.ResolveCredentials(context.Background(), conn.ID)
	if err != nil {
		t.Fatalf("ResolveCredentials(oauth): %v", err)
	}
	if creds.AccessToken != issuedAccessToken || creds.AuthType != "oauth" {
		t.Fatalf("unexpected oauth creds: %+v", creds)
	}

	// Verify providers.changed event emitted without leaking tokens.
	select {
	case ev := <-events:
		raw, _ := json.Marshal(ev.Payload)
		for _, forbidden := range []string{issuedAccessToken, issuedRefreshToken, expectedClientSecret, "v1:"} {
			if strings.Contains(string(raw), forbidden) {
				t.Fatalf("oauth eventbus payload leaked %q", forbidden)
			}
		}
	case <-time.After(time.Second):
		t.Fatalf("expected providers.changed event after oauth exchange")
	}

	mu.Lock()
	if tokenRequestCount != 1 {
		t.Fatalf("expected 1 token exchange call, got %d", tokenRequestCount)
	}
	mu.Unlock()
}

func TestRefreshOAuthTokenSingleflightCoalescing(t *testing.T) {
	var refreshCalls atomic.Int32
	const (
		initialAccess  = "oauth-initial-access-001"
		initialRefresh = "oauth-initial-refresh-001"
		newAccess      = "oauth-refreshed-access-002"
		newRefresh     = "oauth-refreshed-refresh-002"
	)

	mockTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}
		switch r.FormValue("grant_type") {
		case "authorization_code":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token":  initialAccess,
				"refresh_token": initialRefresh,
				"token_type":    "Bearer",
				"expires_in":    60,
			})
		case "refresh_token":
			refreshCalls.Add(1)
			// Hold request briefly so all 8 concurrent goroutines overlap inside singleflight.
			time.Sleep(60 * time.Millisecond)
			if r.FormValue("refresh_token") != initialRefresh {
				http.Error(w, "unexpected refresh token", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token":  newAccess,
				"refresh_token": newRefresh,
				"token_type":    "Bearer",
				"expires_in":    7200,
			})
		default:
			http.Error(w, "unsupported grant_type", http.StatusBadRequest)
		}
	}))
	defer mockTokenServer.Close()

	overrides := map[string]providers.OAuthConfig{
		"codex": {
			AuthorizeURL:       "https://auth.openai.com/oauth/authorize",
			TokenURL:           mockTokenServer.URL + "/token",
			ClientID:           "codex-client-test",
			DefaultRedirectURI: "http://localhost:20128/api/oauth/callback",
		},
	}

	h := newTestHarness(t, nil, overrides)

	// Bootstrap OAuth connection via direct code + state exchange.
	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/codex/start", `{"connection_name":"Codex Account"}`, false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("start codex oauth status = %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		State string `json:"state"`
	}
	_ = json.Unmarshal(startRec.Body.Bytes(), &startResp)

	exchRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/codex/exchange", `{
		"code": "direct-code-123",
		"state": "`+startResp.State+`"
	}`, false)
	if exchRec.Code != http.StatusCreated {
		t.Fatalf("exchange codex oauth status = %d: %s", exchRec.Code, exchRec.Body.String())
	}
	var conn providers.ConnectionSummary
	_ = json.Unmarshal(exchRec.Body.Bytes(), &conn)

	// Fan out 8 concurrent RefreshOAuthToken calls for the same connectionID.
	const callers = 8
	startGate := make(chan struct{})
	var wg sync.WaitGroup
	results := make([]contracts.UpstreamCredentials, callers)
	errs := make([]error, callers)

	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			<-startGate
			creds, err := h.mod.RefreshOAuthToken(context.Background(), conn.ID)
			results[idx] = creds
			errs[idx] = err
		}(i)
	}
	close(startGate)
	wg.Wait()

	if got := refreshCalls.Load(); got != 1 {
		t.Fatalf("expected exactly 1 upstream refresh POST across %d concurrent callers, got %d", callers, got)
	}

	for i := 0; i < callers; i++ {
		if errs[i] != nil {
			t.Fatalf("caller %d RefreshOAuthToken error: %v", i, errs[i])
		}
		if results[i].AccessToken != newAccess {
			t.Fatalf("caller %d got AccessToken %q, want %q", i, results[i].AccessToken, newAccess)
		}
	}

	// Subsequent ResolveCredentials must see persisted refreshed token.
	persisted, err := h.mod.ResolveCredentials(context.Background(), conn.ID)
	if err != nil {
		t.Fatalf("ResolveCredentials after refresh: %v", err)
	}
	if persisted.AccessToken != newAccess {
		t.Fatalf("persisted AccessToken = %q, want %q", persisted.AccessToken, newAccess)
	}
}

func TestOAuthEnvResolutionAndCallbackJSONAndValidation(t *testing.T) {
	t.Setenv("OMEG4MON_OAUTH_GEMINI_OAUTH_CLIENT_ID", "env-gemini-client-id")
	t.Setenv("OMEG4MON_OAUTH_GEMINI_OAUTH_SECRET", "env-gemini-client-secret")

	h := newTestHarness(t, nil, nil)

	// Start gemini-oauth with empty body; client_id/secret should come from env.
	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/gemini-oauth/start", "", false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("POST /oauth/gemini-oauth/start status = %d: %s", startRec.Code, startRec.Body.String())
	}
	if strings.Contains(startRec.Body.String(), "env-gemini-client-secret") {
		t.Fatalf("start response leaked env client secret")
	}

	var startResp struct {
		State        string `json:"state"`
		AuthorizeURL string `json:"authorize_url"`
	}
	if err := json.Unmarshal(startRec.Body.Bytes(), &startResp); err != nil {
		t.Fatalf("decode startResp: %v", err)
	}
	parsed, err := url.Parse(startResp.AuthorizeURL)
	if err != nil {
		t.Fatalf("parse authorize_url: %v", err)
	}
	if parsed.Query().Get("client_id") != "env-gemini-client-id" || parsed.Query().Get("access_type") != "offline" {
		t.Fatalf("unexpected gemini-oauth authorize_url: %s", startResp.AuthorizeURL)
	}

	// Callback JSON handoff verifies active session state without consuming it.
	req := httptest.NewRequest(http.MethodGet, "/api/oauth/callback?code=cb-code&state="+url.QueryEscape(startResp.State), http.NoBody)
	req.Header.Set("Accept", "application/json")
	cbRec := httptest.NewRecorder()
	h.server.Handler().ServeHTTP(cbRec, req)
	if cbRec.Code != http.StatusOK {
		t.Fatalf("GET /api/oauth/callback JSON status = %d", cbRec.Code)
	}
	var cbJSON map[string]string
	if err := json.Unmarshal(cbRec.Body.Bytes(), &cbJSON); err != nil {
		t.Fatalf("decode callback JSON: %v", err)
	}
	if cbJSON["status"] != "verified" || cbJSON["provider_id"] != "gemini-oauth" || cbJSON["code"] != "cb-code" {
		t.Fatalf("unexpected callback JSON: %+v", cbJSON)
	}

	// Unknown field in JSON body must be rejected with 400.
	badFieldRec := h.doRequest(t, http.MethodPost, "/api/providers/nodes", `{"name":"X","protocol":"openai-chat","base_url":"https://example.test","unknown":true}`, false)
	if badFieldRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for unknown JSON field, got %d", badFieldRec.Code)
	}

	// Empty PATCH must return 400.
	emptyPatch := h.doRequest(t, http.MethodPatch, "/api/providers/connections/conn_any", `{}`, false)
	if emptyPatch.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty PATCH, got %d", emptyPatch.Code)
	}

	// Start() rejects missing DB or Vault.
	if err := providers.New(providers.Options{}).Start(context.Background()); err == nil {
		t.Fatalf("expected error starting module without DB/Vault")
	}
}

func TestOAuthExchangePreValidationAndConnectionCompatibility(t *testing.T) {
	var tokenCalls atomic.Int32
	mockTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "reauth-access-token-new",
			"refresh_token": "reauth-refresh-token-new",
			"token_type":    "Bearer",
			"expires_in":    3600,
		})
	}))
	defer mockTokenServer.Close()

	overrides := map[string]providers.OAuthConfig{
		"claude": {
			AuthorizeURL:       "https://claude.ai/oauth/authorize",
			TokenURL:           mockTokenServer.URL + "/token",
			ClientID:           "claude-client-id",
			DefaultRedirectURI: "http://localhost:20128/api/oauth/callback",
		},
		"codex": {
			AuthorizeURL:       "https://auth.openai.com/oauth/authorize",
			TokenURL:           mockTokenServer.URL + "/token",
			ClientID:           "codex-client-id",
			DefaultRedirectURI: "http://localhost:20128/api/oauth/callback",
		},
	}

	h := newTestHarness(t, nil, overrides)

	// Seed an API-key connection (openai), an OAuth connection for another provider (codex),
	// and a matching OAuth connection (claude).
	apiConnRec := h.doRequest(t, http.MethodPost, "/api/providers/connections", `{
		"provider_id": "openai",
		"name": "OpenAI Key Conn",
		"api_key": "sk-original-api-key"
	}`, false)
	if apiConnRec.Code != http.StatusCreated {
		t.Fatalf("seed api key conn status = %d: %s", apiConnRec.Code, apiConnRec.Body.String())
	}
	var apiConn providers.ConnectionSummary
	_ = json.Unmarshal(apiConnRec.Body.Bytes(), &apiConn)

	codexConnRec := h.doRequest(t, http.MethodPost, "/api/providers/connections", `{
		"provider_id": "codex",
		"name": "Codex OAuth Conn",
		"auth_type": "oauth",
		"access_token": "codex-original-access"
	}`, false)
	if codexConnRec.Code != http.StatusCreated {
		t.Fatalf("seed codex oauth conn status = %d: %s", codexConnRec.Code, codexConnRec.Body.String())
	}
	var codexConn providers.ConnectionSummary
	_ = json.Unmarshal(codexConnRec.Body.Bytes(), &codexConn)

	claudeConnRec := h.doRequest(t, http.MethodPost, "/api/providers/connections", `{
		"provider_id": "claude",
		"name": "Claude OAuth Old",
		"auth_type": "oauth",
		"access_token": "claude-original-access",
		"priority": 50
	}`, false)
	if claudeConnRec.Code != http.StatusCreated {
		t.Fatalf("seed claude oauth conn status = %d: %s", claudeConnRec.Code, claudeConnRec.Body.String())
	}
	var claudeConn providers.ConnectionSummary
	_ = json.Unmarshal(claudeConnRec.Body.Bytes(), &claudeConn)

	// Start a single claude OAuth session; all pre-validation failures must fail fast
	// WITHOUT calling the mock token endpoint or consuming the session.
	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/claude/start", `{}`, false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("start claude oauth status = %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		State string `json:"state"`
	}
	_ = json.Unmarshal(startRec.Body.Bytes(), &startResp)

	// 1. Target connection_id is an API-key connection -> 400, tokenCalls == 0, credentials untouched.
	badAuthType := h.doRequest(t, http.MethodPost, "/api/providers/oauth/claude/exchange", `{
		"code": "code-1",
		"state": "`+startResp.State+`",
		"connection_id": "`+apiConn.ID+`"
	}`, false)
	if badAuthType.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when targeting api_key connection, got %d: %s", badAuthType.Code, badAuthType.Body.String())
	}
	if tokenCalls.Load() != 0 {
		t.Fatalf("expected 0 token endpoint calls on api_key mismatch, got %d", tokenCalls.Load())
	}

	// 2. Target connection_id belongs to a different provider (codex) -> 400, tokenCalls == 0.
	badProvider := h.doRequest(t, http.MethodPost, "/api/providers/oauth/claude/exchange", `{
		"code": "code-1",
		"state": "`+startResp.State+`",
		"connection_id": "`+codexConn.ID+`"
	}`, false)
	if badProvider.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when targeting another provider's oauth connection, got %d: %s", badProvider.Code, badProvider.Body.String())
	}
	if tokenCalls.Load() != 0 {
		t.Fatalf("expected 0 token endpoint calls on provider mismatch, got %d", tokenCalls.Load())
	}

	// 3. Target connection_id does not exist -> 404, tokenCalls == 0.
	missingConn := h.doRequest(t, http.MethodPost, "/api/providers/oauth/claude/exchange", `{
		"code": "code-1",
		"state": "`+startResp.State+`",
		"connection_id": "conn_does_not_exist"
	}`, false)
	if missingConn.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when targeting missing connection_id, got %d: %s", missingConn.Code, missingConn.Body.String())
	}
	if tokenCalls.Load() != 0 {
		t.Fatalf("expected 0 token endpoint calls on missing connection_id, got %d", tokenCalls.Load())
	}

	// 4. Negative priority -> 400, tokenCalls == 0.
	badPriority := h.doRequest(t, http.MethodPost, "/api/providers/oauth/claude/exchange", `{
		"code": "code-1",
		"state": "`+startResp.State+`",
		"priority": -5
	}`, false)
	if badPriority.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for negative priority, got %d: %s", badPriority.Code, badPriority.Body.String())
	}
	if tokenCalls.Load() != 0 {
		t.Fatalf("expected 0 token endpoint calls on negative priority, got %d", tokenCalls.Load())
	}

	// 5. Invalid headers (CRLF injection) -> 400, tokenCalls == 0.
	badHeaders := h.doRequest(t, http.MethodPost, "/api/providers/oauth/claude/exchange", `{
		"code": "code-1",
		"state": "`+startResp.State+`",
		"headers": {"X-Bad": "val\r\nInjected: 1"}
	}`, false)
	if badHeaders.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid headers, got %d: %s", badHeaders.Code, badHeaders.Body.String())
	}
	if tokenCalls.Load() != 0 {
		t.Fatalf("expected 0 token endpoint calls on invalid headers, got %d", tokenCalls.Load())
	}

	// Verify existing mismatched connections remain completely untouched.
	apiCreds, err := h.mod.ResolveCredentials(context.Background(), apiConn.ID)
	if err != nil || apiCreds.APIKey != "sk-original-api-key" {
		t.Fatalf("api_key connection was unexpectedly altered: creds=%+v err=%v", apiCreds, err)
	}
	codexCreds, err := h.mod.ResolveCredentials(context.Background(), codexConn.ID)
	if err != nil || codexCreds.AccessToken != "codex-original-access" {
		t.Fatalf("codex oauth connection was unexpectedly altered: creds=%+v err=%v", codexCreds, err)
	}

	// 6. Valid exchange updating the matching claude OAuth connection succeeds with the same unburned session!
	okUpdate := h.doRequest(t, http.MethodPost, "/api/providers/oauth/claude/exchange", `{
		"code": "code-valid",
		"state": "`+startResp.State+`",
		"connection_id": "`+claudeConn.ID+`",
		"connection_name": "Claude OAuth Reauthed",
		"priority": 15
	}`, false)
	if okUpdate.Code != http.StatusCreated {
		t.Fatalf("expected 201 updating matching oauth connection, got %d: %s", okUpdate.Code, okUpdate.Body.String())
	}
	if tokenCalls.Load() != 1 {
		t.Fatalf("expected exactly 1 token endpoint call after valid exchange, got %d", tokenCalls.Load())
	}

	var updated providers.ConnectionSummary
	if err := json.Unmarshal(okUpdate.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode updated connection: %v", err)
	}
	if updated.ID != claudeConn.ID || updated.Name != "Claude OAuth Reauthed" || updated.Priority != 15 || updated.TokenExpiresAt == nil {
		t.Fatalf("unexpected updated connection summary: %+v", updated)
	}

	newCreds, err := h.mod.ResolveCredentials(context.Background(), claudeConn.ID)
	if err != nil || newCreds.AccessToken != "reauth-access-token-new" {
		t.Fatalf("expected updated claude oauth access token, got creds=%+v err=%v", newCreds, err)
	}
}

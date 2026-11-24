package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

func TestParseDiscoveredModelsFormats(t *testing.T) {
	t.Run("OpenAI data array", func(t *testing.T) {
		raw := []byte(`{"object":"list","data":[{"id":"gpt-5.4","name":"GPT-5.4"},{"id":"o4-mini"}]}`)
		got := parseDiscoveredModels(raw)
		if len(got) != 2 || got[0].ID != "gpt-5.4" || got[0].Name != "GPT-5.4" || got[1].ID != "o4-mini" {
			t.Fatalf("unexpected openai models: %+v", got)
		}
	})

	t.Run("Gemini models array strips models/ prefix", func(t *testing.T) {
		raw := []byte(`{"models":[{"name":"models/gemini-2.5-pro","displayName":"Gemini 2.5 Pro"},{"name":"models/gemini-2.5-flash"}]}`)
		got := parseDiscoveredModels(raw)
		if len(got) != 2 || got[0].ID != "gemini-2.5-pro" || got[0].Name != "Gemini 2.5 Pro" || got[1].ID != "gemini-2.5-flash" {
			t.Fatalf("unexpected gemini models: %+v", got)
		}
	})

	t.Run("Antigravity object map skips internal models", func(t *testing.T) {
		raw := []byte(`{"models":{"gemini-4-carbon":{"displayName":"Gemini 4 Carbon"},"internal-hidden":{"displayName":"Hidden","isInternal":true}}}`)
		got := parseDiscoveredModels(raw)
		if len(got) != 1 || got[0].ID != "gemini-4-carbon" || got[0].Name != "Gemini 4 Carbon" {
			t.Fatalf("unexpected antigravity models: %+v", got)
		}
	})

	t.Run("Codex slug list", func(t *testing.T) {
		raw := []byte(`{"models":[{"slug":"gpt-6.1-sol","display_name":"GPT-6.1 Sol"}]}`)
		got := parseDiscoveredModels(raw)
		if len(got) != 1 || got[0].ID != "gpt-6.1-sol" || got[0].Name != "GPT-6.1 Sol" {
			t.Fatalf("unexpected codex models: %+v", got)
		}
	})
}

type stubDiscoverySystemModule struct{}

func (stubDiscoverySystemModule) Name() string                         { return "system" }
func (stubDiscoverySystemModule) Dependencies() []string               { return nil }
func (stubDiscoverySystemModule) Migrations() fs.FS                    { return nil }
func (stubDiscoverySystemModule) RegisterRoutes(lifecycle.RouteMounts) {}
func (stubDiscoverySystemModule) Start(context.Context) error          { return nil }
func (stubDiscoverySystemModule) Stop(context.Context) error           { return nil }

func TestFetchConnectionModelsEndpoint(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") != "Bearer sk-live-discovery" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"live-model-alpha","name":"Live Alpha"},{"id":"live-model-beta"}]}`))
	}))
	defer upstream.Close()

	dbPath := filepath.Join(t.TempDir(), "router.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	masterKey := make([]byte, crypto.MasterKeySize)
	for i := range masterKey {
		masterKey[i] = byte(i + 1)
	}
	vault, err := crypto.NewVault(masterKey)
	if err != nil {
		t.Fatalf("crypto.NewVault: %v", err)
	}

	mod := New(Options{
		DB:         database,
		Vault:      vault,
		HTTPClient: upstream.Client(),
		Bus:        eventbus.New(),
	})

	ctx := context.Background()
	reg := lifecycle.NewRegistry(nil)
	if err := reg.Register(stubDiscoverySystemModule{}); err != nil {
		t.Fatalf("reg.Register(system): %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("reg.Register: %v", err)
	}
	if err := reg.MigrateAll(ctx, database); err != nil {
		t.Fatalf("MigrateAll: %v", err)
	}

	srv := httpserver.New(httpserver.Options{
		Config: config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: 10 * time.Second},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	mod.RegisterRoutes(srv.Mounts())

	doReq := func(method, path, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		return rec
	}

	createBody := `{"provider_id":"openai","name":"Discovery Key","protocol":"openai-chat","auth_type":"api_key","base_url":"` + upstream.URL + `/v1","api_key":"sk-live-discovery","models":["gpt-5.4"],"priority":10,"enabled":true}`
	createRec := doReq(http.MethodPost, "/api/providers/connections", createBody)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create connection failed: %d %s", createRec.Code, createRec.Body.String())
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal created: %v", err)
	}

	rec := doReq(http.MethodGet, "/api/providers/connections/"+created.ID+"/models", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}
	var res fetchConnectionModelsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if res.Source != "upstream" {
		t.Fatalf("expected source=upstream, got %q (warning=%q)", res.Source, res.Warning)
	}
	if len(res.Models) != 2 || res.Models[0].ID != "live-model-alpha" || res.Models[1].ID != "live-model-beta" {
		t.Fatalf("unexpected discovered models: %+v", res.Models)
	}
	_ = lifecycle.RouteMounts{}
}

func TestAutoDetectModelsDiscoveryCandidates(t *testing.T) {
	cases := []struct {
		name     string
		protocol string
		baseURL  string
		want0    string
	}{
		{
			name:     "strips /chat/completions suffix",
			protocol: protocolOpenAIChat,
			baseURL:  "https://api.example.com/v1/chat/completions",
			want0:    "https://api.example.com/v1/models",
		},
		{
			name:     "strips /v1/messages suffix for anthropic",
			protocol: protocolAnthropic,
			baseURL:  "https://api.anthropic.com/v1/messages",
			want0:    "https://api.anthropic.com/v1/models",
		},
		{
			name:     "avoids double /models/models",
			protocol: protocolOpenAIChat,
			baseURL:  "https://api.example.com/v1/models",
			want0:    "https://api.example.com/v1/models",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildDiscoveryCandidateURLs("", tc.protocol, tc.baseURL)
			if len(got) == 0 || got[0] != tc.want0 {
				t.Fatalf("expected first candidate %q, got %v", tc.want0, got)
			}
		})
	}
}

func TestOAuthCallbackEmitsAutoCompletionBridge(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "router.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	masterKey := make([]byte, crypto.MasterKeySize)
	for i := range masterKey {
		masterKey[i] = byte(i + 2)
	}
	vault, err := crypto.NewVault(masterKey)
	if err != nil {
		t.Fatalf("crypto.NewVault: %v", err)
	}

	mod := New(Options{DB: database, Vault: vault, Bus: eventbus.New()})
	reg := lifecycle.NewRegistry(nil)
	if err := reg.Register(stubDiscoverySystemModule{}); err != nil {
		t.Fatalf("register system: %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("register providers: %v", err)
	}
	if err := reg.MigrateAll(context.Background(), database); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	srv := httpserver.New(httpserver.Options{
		Config: config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: 10 * time.Second},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	mod.RegisterRoutes(srv.Mounts())

	req := httptest.NewRequest(http.MethodGet, "/api/providers/oauth/callback?code=abc123&state=st456", http.NoBody)
	req.Header.Set("Accept", "text/html")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"oauth_callback"`) || !strings.Contains(body, "BroadcastChannel") || !strings.Contains(body, "postMessage") {
		t.Fatalf("expected OAuth callback HTML to emit postMessage/BroadcastChannel bridge, got: %s", body)
	}
	csp := rec.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "script-src 'unsafe-inline'") {
		t.Fatalf("expected CSP to allow inline bridge script, got: %q", csp)
	}
}

func TestFetchConnectionModelsPersistsAndMergesExistingModels(t *testing.T) {
	modelCount := 300
	upstreamModels := make([]map[string]string, 0, modelCount)
	for i := 0; i < modelCount; i++ {
		upstreamModels = append(upstreamModels, map[string]string{
			"id": fmt.Sprintf("upstream-model-%03d", i),
		})
	}
	// Duplicate entry and one matching an existing manual model to verify deduplication.
	upstreamModels = append(upstreamModels,
		map[string]string{"id": "manual-custom-model"},
		map[string]string{"id": "upstream-model-000"},
	)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": upstreamModels})
	}))
	defer upstream.Close()

	dbPath := filepath.Join(t.TempDir(), "router.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	masterKey := make([]byte, crypto.MasterKeySize)
	for i := range masterKey {
		masterKey[i] = byte(i + 7)
	}
	vault, err := crypto.NewVault(masterKey)
	if err != nil {
		t.Fatalf("crypto.NewVault: %v", err)
	}

	bus := eventbus.New()
	eventsCh, unsubscribe := bus.Subscribe("providers.changed", 16)
	defer unsubscribe()
	drainEvents := func() int {
		count := 0
		for {
			select {
			case <-eventsCh:
				count++
			default:
				return count
			}
		}
	}

	mod := New(Options{
		DB:         database,
		Vault:      vault,
		Bus:        bus,
		HTTPClient: upstream.Client(),
	})

	reg := lifecycle.NewRegistry(nil)
	if err := reg.Register(stubDiscoverySystemModule{}); err != nil {
		t.Fatalf("register system: %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("register providers: %v", err)
	}
	if err := reg.MigrateAll(context.Background(), database); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	srv := httpserver.New(httpserver.Options{
		Config: config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: 10 * time.Second},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	mod.RegisterRoutes(srv.Mounts())

	// Seed a connection with one pre-existing manual model.
	createBody := `{"provider_id":"openai","name":"OpenAI Proxy","base_url":"` + upstream.URL + `/v1","api_key":"sk-persist-test","models":["manual-custom-model"]}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/providers/connections", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create connection status=%d body=%s", createRec.Code, createRec.Body.String())
	}

	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created connection: %v", err)
	}
	_ = drainEvents()

	// Fetch models — should persist merged upstream + manual models automatically.
	fetchReq := httptest.NewRequest(http.MethodGet, "/api/providers/connections/"+created.ID+"/models", http.NoBody)
	fetchRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(fetchRec, fetchReq)
	if fetchRec.Code != http.StatusOK {
		t.Fatalf("fetch models status=%d body=%s", fetchRec.Code, fetchRec.Body.String())
	}

	// Verify connection persisted all 301 unique models (1 manual + 300 upstream) via store contract.
	conns, err := mod.ListEnabledConnections(context.Background())
	if err != nil {
		t.Fatalf("ListEnabledConnections: %v", err)
	}
	if len(conns) != 1 {
		t.Fatalf("expected 1 connection, got %d", len(conns))
	}
	gotModels := conns[0].Models
	wantCount := modelCount + 1
	if len(gotModels) != wantCount {
		t.Fatalf("expected %d persisted models (manual + upstream), got %d", wantCount, len(gotModels))
	}
	if gotModels[0] != "manual-custom-model" {
		t.Fatalf("expected pre-existing manual model preserved first, got %q", gotModels[0])
	}
	if drainEvents() == 0 {
		t.Fatalf("expected providers.changed event after persisting discovered models")
	}

	// Verify PATCH with 302 models (adding another manual model on top of 301) succeeds without hitting 256 limit.
	nextModels := append(append([]string(nil), gotModels...), "manual-extra-model")
	patchPayload, err := json.Marshal(map[string]any{"models": nextModels})
	if err != nil {
		t.Fatalf("marshal patch payload: %v", err)
	}
	patchReq := httptest.NewRequest(http.MethodPatch, "/api/providers/connections/"+created.ID, strings.NewReader(string(patchPayload)))
	patchReq.Header.Set("Content-Type", "application/json")
	patchRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(patchRec, patchReq)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("patch connection with %d models failed: status=%d body=%s", len(nextModels), patchRec.Code, patchRec.Body.String())
	}

	// Re-fetching identical upstream models must preserve the newly added manual model and emit no redundant event.
	_ = drainEvents()
	fetchReq2 := httptest.NewRequest(http.MethodGet, "/api/providers/connections/"+created.ID+"/models", http.NoBody)
	fetchRec2 := httptest.NewRecorder()
	srv.Handler().ServeHTTP(fetchRec2, fetchReq2)
	if fetchRec2.Code != http.StatusOK {
		t.Fatalf("second fetch status=%d body=%s", fetchRec2.Code, fetchRec2.Body.String())
	}
	conns2, err := mod.ListEnabledConnections(context.Background())
	if err != nil {
		t.Fatalf("ListEnabledConnections second: %v", err)
	}
	if len(conns2[0].Models) != modelCount+2 {
		t.Fatalf("expected manual-extra-model preserved across re-fetch (%d models), got %d", modelCount+2, len(conns2[0].Models))
	}
	if delta := drainEvents(); delta != 0 {
		t.Fatalf("expected no redundant providers.changed event when models are unchanged, got delta %d", delta)
	}
}

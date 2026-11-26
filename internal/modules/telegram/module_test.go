package telegram_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/telegram"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubSystem struct{}

func (stubSystem) Name() string                         { return "system" }
func (stubSystem) Dependencies() []string               { return nil }
func (stubSystem) Migrations() fs.FS                    { return nil }
func (stubSystem) RegisterRoutes(lifecycle.RouteMounts) {}
func (stubSystem) Start(context.Context) error          { return nil }
func (stubSystem) Stop(context.Context) error           { return nil }

type stubProviderStore struct {
	conns []contracts.ProviderConnection
}

func (s *stubProviderStore) ListEnabledConnections(context.Context) ([]contracts.ProviderConnection, error) {
	return s.conns, nil
}

func (s *stubProviderStore) ResolveCredentials(context.Context, string) (contracts.UpstreamCredentials, error) {
	return contracts.UpstreamCredentials{}, contracts.ErrConnectionNotFound
}

func (s *stubProviderStore) RefreshOAuthToken(context.Context, string) (contracts.UpstreamCredentials, error) {
	return contracts.UpstreamCredentials{}, contracts.ErrOAuthRefreshUnsupported
}

type stubModelLister struct {
	models []contracts.ModelDescriptor
}

func (s *stubModelLister) ListModels(context.Context, []string) ([]contracts.ModelDescriptor, error) {
	return s.models, nil
}

type stubJobInspector struct {
	jobs []contracts.JobStatus
}

func (s *stubJobInspector) Snapshot() []contracts.JobStatus { return s.jobs }
func (s *stubJobInspector) Trigger(context.Context, string) error {
	return nil
}

type tgAPICall struct {
	Path string
	Body map[string]any
}

type tgTestEnv struct {
	mod       *telegram.Module
	bus       *eventbus.Bus
	adminSrv  *httptest.Server
	publicSrv *httptest.Server
	apiSrv    *httptest.Server

	mu       sync.Mutex
	apiCalls []tgAPICall
}

func (e *tgTestEnv) calls() []tgAPICall {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]tgAPICall(nil), e.apiCalls...)
}

func setupTelegramTest(t *testing.T) *tgTestEnv {
	t.Helper()

	env := &tgTestEnv{
		bus: eventbus.New(),
	}

	env.apiSrv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var parsed map[string]any
		_ = json.Unmarshal(raw, &parsed)
		env.mu.Lock()
		env.apiCalls = append(env.apiCalls, tgAPICall{Path: r.URL.Path, Body: parsed})
		env.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":42}}`))
	}))
	t.Cleanup(env.apiSrv.Close)

	database, err := db.Open(filepath.Join(t.TempDir(), "telegram.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	key := bytes.Repeat([]byte{0x42}, 32)
	vault, err := crypto.NewVault(key)
	if err != nil {
		t.Fatalf("NewVault: %v", err)
	}

	mod := telegram.New(telegram.Options{
		DB:         database,
		Vault:      vault,
		HTTPClient: env.apiSrv.Client(),
		Bus:        env.bus,
		Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
		Providers: &stubProviderStore{conns: []contracts.ProviderConnection{
			{ID: "c1", ProviderID: "openai", Name: "OpenAI Primary", Enabled: true},
		}},
		Models: &stubModelLister{models: []contracts.ModelDescriptor{
			{ID: "gpt-4o", OwnedBy: "openai", Kind: "model"},
			{ID: "fast-combo", OwnedBy: "combo", Kind: "combo"},
		}},
		Jobs: &stubJobInspector{jobs: []contracts.JobStatus{
			{Name: "usage-prune", LastStatus: "ok", RunCount: 3},
		}},
		ActiveModules: func() []string { return []string{"system", "gateway", "telegram"} },
		APIBaseURL:    env.apiSrv.URL,
	})
	env.mod = mod

	reg := lifecycle.NewRegistry(nil)
	if err := reg.Register(stubSystem{}); err != nil {
		t.Fatalf("register system: %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("register telegram: %v", err)
	}
	if err := reg.MigrateAll(context.Background(), database); err != nil {
		t.Fatalf("MigrateAll: %v", err)
	}

	adminRouter := chi.NewRouter()
	publicRouter := chi.NewRouter()
	adminRouter.Route("/api", func(r chi.Router) {
		publicRouter.Route("/api", func(pr chi.Router) {
			mod.RegisterRoutes(lifecycle.RouteMounts{
				AdminAPI:  r,
				PublicAPI: pr,
			})
		})
	})

	if err := reg.StartAll(context.Background()); err != nil {
		t.Fatalf("StartAll: %v", err)
	}
	t.Cleanup(func() { _ = reg.StopAll(context.Background()) })

	env.adminSrv = httptest.NewServer(adminRouter)
	env.publicSrv = httptest.NewServer(publicRouter)
	t.Cleanup(env.adminSrv.Close)
	t.Cleanup(env.publicSrv.Close)
	return env
}

func doJSON(t *testing.T, method, url string, body any, headers map[string]string) (*http.Response, []byte) {
	t.Helper()
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, url, rdr)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do %s %s: %v", method, url, err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	return resp, raw
}

func TestTelegramConfigVaultIsolationAndWebhookSecurity(t *testing.T) {
	env := setupTelegramTest(t)

	// 1. Default config is disabled and has no secrets.
	resp, raw := doJSON(t, http.MethodGet, env.adminSrv.URL+"/api/telegram/config", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET config status=%d body=%s", resp.StatusCode, raw)
	}
	var cfg map[string]any
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("unmarshal cfg: %v", err)
	}
	if cfg["enabled"] != false || cfg["bot_token_configured"] != false {
		t.Fatalf("unexpected default cfg: %+v", cfg)
	}

	// 2. Webhook fails closed while disabled.
	resp, _ = doJSON(t, http.MethodPost, env.publicSrv.URL+"/api/telegram/webhook", map[string]any{}, nil)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 when disabled, got %d", resp.StatusCode)
	}

	// 3. Configure bot token + chat_id + enabled (without user_id yet).
	rawToken := "123456789:AAH-secret-bot-token-xyz987"
	resp, raw = doJSON(t, http.MethodPatch, env.adminSrv.URL+"/api/telegram/config", map[string]any{
		"enabled":     true,
		"bot_token":   rawToken,
		"chat_id":     "10001",
		"webhook_url": "https://router.example.com/api/telegram/webhook",
	}, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PATCH config status=%d body=%s", resp.StatusCode, raw)
	}
	if strings.Contains(string(raw), rawToken) {
		t.Fatalf("raw bot token leaked in PATCH response: %s", raw)
	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("unmarshal updated cfg: %v", err)
	}
	if cfg["bot_token_configured"] != true || cfg["webhook_secret_configured"] != true {
		t.Fatalf("expected configured flags true, got %+v", cfg)
	}

	// 4. Register webhook to capture generated secret_token sent to Telegram Bot API.
	resp, raw = doJSON(t, http.MethodPost, env.adminSrv.URL+"/api/telegram/webhook/registration", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("register webhook status=%d body=%s", resp.StatusCode, raw)
	}
	calls := env.calls()
	if len(calls) == 0 || !strings.HasSuffix(calls[len(calls)-1].Path, "/setWebhook") {
		t.Fatalf("expected setWebhook call, got %+v", calls)
	}
	secretToken, _ := calls[len(calls)-1].Body["secret_token"].(string)
	if len(secretToken) < 16 {
		t.Fatalf("expected strong secret_token in setWebhook, got %q", secretToken)
	}

	// 5. Missing/wrong secret header -> 401.
	resp, _ = doJSON(t, http.MethodPost, env.publicSrv.URL+"/api/telegram/webhook", map[string]any{}, map[string]string{
		"X-Telegram-Bot-Api-Secret-Token": "wrong-secret",
	})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong secret, got %d", resp.StatusCode)
	}

	// 6. Valid secret but empty configured user_id -> 403 fail-closed.
	updateBody := map[string]any{
		"update_id": 1,
		"message": map[string]any{
			"message_id": 10,
			"from":       map[string]any{"id": 777001},
			"chat":       map[string]any{"id": 10001},
			"text":       "/status",
		},
	}
	resp, _ = doJSON(t, http.MethodPost, env.publicSrv.URL+"/api/telegram/webhook", updateBody, map[string]string{
		"X-Telegram-Bot-Api-Secret-Token": secretToken,
	})
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 when user_id unset, got %d", resp.StatusCode)
	}

	// 7. Set user_id=777001; mismatched sender -> 403; matching sender -> 200 + sendMessage reply.
	resp, _ = doJSON(t, http.MethodPatch, env.adminSrv.URL+"/api/telegram/config", map[string]any{
		"user_id": "777001",
	}, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PATCH user_id status=%d", resp.StatusCode)
	}

	wrongUserUpdate := map[string]any{
		"update_id": 2,
		"message": map[string]any{
			"message_id": 11,
			"from":       map[string]any{"id": 999999},
			"chat":       map[string]any{"id": 10001},
			"text":       "/status",
		},
	}
	resp, _ = doJSON(t, http.MethodPost, env.publicSrv.URL+"/api/telegram/webhook", wrongUserUpdate, map[string]string{
		"X-Telegram-Bot-Api-Secret-Token": secretToken,
	})
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for wrong user_id, got %d", resp.StatusCode)
	}

	resp, raw = doJSON(t, http.MethodPost, env.publicSrv.URL+"/api/telegram/webhook", updateBody, map[string]string{
		"X-Telegram-Bot-Api-Secret-Token": secretToken,
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for authorized webhook, got %d body=%s", resp.StatusCode, raw)
	}

	// 8. Rotate secret -> old secret returns 401.
	resp, raw = doJSON(t, http.MethodPost, env.adminSrv.URL+"/api/telegram/secret/rotate", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("rotate secret status=%d body=%s", resp.StatusCode, raw)
	}
	resp, _ = doJSON(t, http.MethodPost, env.publicSrv.URL+"/api/telegram/webhook", updateBody, map[string]string{
		"X-Telegram-Bot-Api-Secret-Token": secretToken,
	})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 after secret rotation, got %d", resp.StatusCode)
	}
}

func TestTelegramCommandsTestMessageAndEventBusForwarding(t *testing.T) {
	env := setupTelegramTest(t)

	resp, _ := doJSON(t, http.MethodPatch, env.adminSrv.URL+"/api/telegram/config", map[string]any{
		"enabled":               true,
		"bot_token":             "987654:AAH-test-token-123456",
		"chat_id":               "555001",
		"user_id":               "555001",
		"notify_on_job_failure": true,
	}, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PATCH config status=%d", resp.StatusCode)
	}

	// Command simulator tests.
	for _, cmd := range []string{"/help", "/ping", "/status", "/providers", "/models", "/jobs"} {
		resp, raw := doJSON(t, http.MethodPost, env.adminSrv.URL+"/api/telegram/command", map[string]any{
			"command": cmd,
		}, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("command %s status=%d body=%s", cmd, resp.StatusCode, raw)
		}
		var out struct {
			Command string `json:"command"`
			Reply   string `json:"reply"`
		}
		if err := json.Unmarshal(raw, &out); err != nil || strings.TrimSpace(out.Reply) == "" {
			t.Fatalf("command %s invalid output: %s", cmd, raw)
		}
	}

	// Send test message.
	resp, raw := doJSON(t, http.MethodPost, env.adminSrv.URL+"/api/telegram/test", map[string]any{
		"message": "Wave 3D probe",
	}, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST test status=%d body=%s", resp.StatusCode, raw)
	}

	// Publish jobs.failed on eventbus and wait for async sendMessage.
	beforeCalls := len(env.calls())
	env.bus.Publish("jobs.failed", map[string]any{
		"title":   "Background job usage-prune failed",
		"message": "disk busy",
	})

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if len(env.calls()) > beforeCalls {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if len(env.calls()) <= beforeCalls {
		t.Fatalf("expected eventbus jobs.failed to trigger sendMessage")
	}

	// Events list and clear.
	resp, raw = doJSON(t, http.MethodGet, env.adminSrv.URL+"/api/telegram/events", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET events status=%d body=%s", resp.StatusCode, raw)
	}
	var evList struct {
		Events []map[string]any `json:"events"`
	}
	if err := json.Unmarshal(raw, &evList); err != nil || len(evList.Events) == 0 {
		t.Fatalf("expected non-empty telegram events, got %s", raw)
	}

	resp, _ = doJSON(t, http.MethodDelete, env.adminSrv.URL+"/api/telegram/events", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE events status=%d", resp.StatusCode)
	}
}

func TestTelegramNetworkErrorRedactsBotToken(t *testing.T) {
	// Create a closed server to force a *url.Error on sendMessage.
	closedSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	closedURL := closedSrv.URL
	closedSrv.Close()

	database, err := db.Open(filepath.Join(t.TempDir(), "tg-redact.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	vault, err := crypto.NewVault(bytes.Repeat([]byte{0x55}, 32))
	if err != nil {
		t.Fatalf("NewVault: %v", err)
	}

	mod := telegram.New(telegram.Options{
		DB:         database,
		Vault:      vault,
		HTTPClient: &http.Client{Timeout: 2 * time.Second},
		Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
		APIBaseURL: closedURL,
	})
	reg := lifecycle.NewRegistry(nil)
	_ = reg.Register(stubSystem{})
	_ = reg.Register(mod)
	if err := reg.MigrateAll(context.Background(), database); err != nil {
		t.Fatalf("MigrateAll: %v", err)
	}
	r := chi.NewRouter()
	r.Route("/api", func(ar chi.Router) {
		mod.RegisterRoutes(lifecycle.RouteMounts{AdminAPI: ar})
	})
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	secretToken := "999888777:AAH-SUPER-SECRET-TOKEN-DO-NOT-LEAK"
	resp, _ := doJSON(t, http.MethodPatch, srv.URL+"/api/telegram/config", map[string]any{
		"enabled":   true,
		"bot_token": secretToken,
		"chat_id":   "12345",
	}, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PATCH config status=%d", resp.StatusCode)
	}

	resp, raw := doJSON(t, http.MethodPost, srv.URL+"/api/telegram/test", map[string]any{"message": "ping"}, nil)
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected 502 on closed Telegram API, got %d body=%s", resp.StatusCode, raw)
	}
	if strings.Contains(string(raw), secretToken) {
		t.Fatalf("bot token leaked in 502 error response: %s", raw)
	}

	_, evRaw := doJSON(t, http.MethodGet, srv.URL+"/api/telegram/events", nil, nil)
	if strings.Contains(string(evRaw), secretToken) {
		t.Fatalf("bot token leaked in telegram_events: %s", evRaw)
	}
}

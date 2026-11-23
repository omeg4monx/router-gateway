package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/modules/system"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/jobs"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
	"github.com/omeg4monx/router-gateway/internal/platform/logging"
)

func TestAppBootstrapsMigrationsRoutesAndShutdown(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		Host:       "127.0.0.1",
		Port:       0,
		DataDir:    t.TempDir(),
		LogLevel:   "error",
		APITimeout: 5 * time.Second,
	}
	build := system.BuildInfo{
		Version:   "test",
		Commit:    "abc1234",
		BuildTime: "2026-11-22T00:00:00Z",
	}

	app, err := buildApp(ctx, cfg, build)
	if err != nil {
		t.Fatalf("buildApp() error = %v", err)
	}
	t.Cleanup(func() {
		_ = app.Close(context.Background())
	})

	handler := app.Handler()

	// 1) GET /api/health -> 200 JSON with modules containing "system"
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/health status = %d, want %d (body: %s)", rec.Code, http.StatusOK, rec.Body.String())
	}
	var health struct {
		Status  string   `json:"status"`
		DB      string   `json:"db"`
		Modules []string `json:"modules"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &health); err != nil {
		t.Fatalf("decode /api/health JSON: %v", err)
	}
	if health.Status != "ok" || health.DB != "ok" {
		t.Fatalf("unexpected health payload: %+v", health)
	}
	if !slices.Contains(health.Modules, "system") {
		t.Fatalf("expected modules to contain \"system\", got %v", health.Modules)
	}

	// 2) GET /api/version -> 200 JSON with build metadata
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/version", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/version status = %d, want %d", rec.Code, http.StatusOK)
	}
	var ver struct {
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		BuildTime string `json:"buildTime"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &ver); err != nil {
		t.Fatalf("decode /api/version JSON: %v", err)
	}
	if ver.Version != "test" || ver.Commit != "abc1234" || ver.BuildTime != "2026-11-22T00:00:00Z" {
		t.Fatalf("unexpected version payload: %+v", ver)
	}

	// 3) GET /api/init -> 200 {initialized: false}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/init", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/init status = %d, want %d (body: %s)", rec.Code, http.StatusOK, rec.Body.String())
	}
	var initResp struct {
		Initialized  bool `json:"initialized"`
		RequireLogin bool `json:"requireLogin"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &initResp); err != nil {
		t.Fatalf("decode /api/init JSON: %v", err)
	}
	if initResp.Initialized {
		t.Fatalf("expected initialized=false, got true")
	}

	// 4) GET /api/nonexistent -> 404 JSON
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/nonexistent", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET /api/nonexistent status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("GET /api/nonexistent Content-Type = %q, want application/json", ct)
	}

	// 5) app.Close(ctx) -> nil (and idempotent on second call)
	if err := app.Close(ctx); err != nil {
		t.Fatalf("app.Close() error = %v", err)
	}
	if err := app.Close(ctx); err != nil {
		t.Fatalf("second app.Close() error = %v", err)
	}
}

func TestAppRunAndShutdownEndpointTriggersGracefulExit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := config.Config{
		Host:          "127.0.0.1",
		Port:          0,
		DataDir:       t.TempDir(),
		AdminPassword: "shutdown-test-pass-2026",
		LogLevel:      "error",
		APITimeout:    5 * time.Second,
	}

	app, err := buildApp(ctx, cfg, system.BuildInfo{Version: "test"})
	if err != nil {
		t.Fatalf("buildApp() error = %v", err)
	}

	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- app.Run(ctx)
	}()

	var addr string
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		candidate := app.Addr()
		if candidate != "" && !strings.HasSuffix(candidate, ":0") {
			addr = candidate
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if addr == "" {
		t.Fatalf("server failed to bind ephemeral address")
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 2 * time.Second,
	}

	healthResp, err := client.Get("http://" + addr + "/api/health")
	if err != nil {
		t.Fatalf("GET /api/health failed: %v", err)
	}
	_, _ = io.ReadAll(healthResp.Body)
	_ = healthResp.Body.Close()
	if healthResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/health status = %d, want 200", healthResp.StatusCode)
	}

	// Unauthenticated POST /api/shutdown must be rejected when require_login=true.
	unauthShutResp, err := client.Post("http://"+addr+"/api/shutdown", "application/json", nil)
	if err != nil {
		t.Fatalf("unauthenticated POST /api/shutdown failed: %v", err)
	}
	_, _ = io.ReadAll(unauthShutResp.Body)
	_ = unauthShutResp.Body.Close()
	if unauthShutResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated POST /api/shutdown status = %d, want 401", unauthShutResp.StatusCode)
	}

	loginResp, err := client.Post(
		"http://"+addr+"/api/auth/login",
		"application/json",
		strings.NewReader(`{"password":"shutdown-test-pass-2026"}`),
	)
	if err != nil {
		t.Fatalf("POST /api/auth/login failed: %v", err)
	}
	loginBody, _ := io.ReadAll(loginResp.Body)
	_ = loginResp.Body.Close()
	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("POST /api/auth/login status = %d, want 200 (body: %s)", loginResp.StatusCode, string(loginBody))
	}

	shutResp, err := client.Post("http://"+addr+"/api/shutdown", "application/json", nil)
	if err != nil {
		t.Fatalf("POST /api/shutdown failed: %v", err)
	}
	_, _ = io.ReadAll(shutResp.Body)
	_ = shutResp.Body.Close()
	if shutResp.StatusCode != http.StatusAccepted {
		t.Fatalf("POST /api/shutdown status = %d, want 202", shutResp.StatusCode)
	}

	select {
	case err := <-runErrCh:
		if err != nil {
			t.Fatalf("app.Run() returned error after shutdown: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for app.Run() to exit after POST /api/shutdown")
	}
}

func TestBuildAppCleansUpOnPartialFailure(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		Host:          "127.0.0.1",
		Port:          0,
		DataDir:       t.TempDir(),
		AdminPassword: "short", // triggers system.Module.Start validation failure after DB is opened
		LogLevel:      "error",
		APITimeout:    5 * time.Second,
	}

	app, err := buildApp(ctx, cfg, system.BuildInfo{Version: "test"})
	if err == nil {
		_ = app.Close(ctx)
		t.Fatalf("expected buildApp() to fail on invalid AdminPassword, got nil")
	}
	if app != nil {
		t.Fatalf("expected nil app on buildApp failure")
	}
}

type orderStubModule struct {
	name   string
	onStop func(context.Context) error
}

func (m *orderStubModule) Name() string                         { return m.name }
func (m *orderStubModule) Dependencies() []string               { return nil }
func (m *orderStubModule) Migrations() fs.FS                    { return nil }
func (m *orderStubModule) RegisterRoutes(lifecycle.RouteMounts) {}
func (m *orderStubModule) Start(context.Context) error          { return nil }
func (m *orderStubModule) Stop(ctx context.Context) error {
	if m.onStop != nil {
		return m.onStop(ctx)
	}
	return nil
}

func TestAppCloseTeardownOrderIsLIFO(t *testing.T) {
	ctx := context.Background()
	logger, _ := logging.New("error", io.Discard, 16)

	database, err := db.Open(filepath.Join(t.TempDir(), "lifo.sqlite"))
	if err != nil {
		t.Fatalf("db.Open() error = %v", err)
	}

	var (
		mu     sync.Mutex
		events []string
	)
	record := func(ev string) {
		mu.Lock()
		events = append(events, ev)
		mu.Unlock()
	}

	jobStarted := make(chan struct{})
	var startOnce sync.Once

	jobRunner := jobs.NewRunner(logger)
	jobRunner.Every("lifo-check", time.Millisecond, func(jobCtx context.Context) error {
		startOnce.Do(func() { close(jobStarted) })
		<-jobCtx.Done()
		// Yield briefly so an incorrect teardown order that stops modules before
		// jobs reliably records module_stopped first.
		time.Sleep(15 * time.Millisecond)
		record("job_stopped")
		return jobCtx.Err()
	})

	registry := lifecycle.NewRegistry(nil)
	if err := registry.Register(&orderStubModule{
		name: "stub",
		onStop: func(stopCtx context.Context) error {
			if err := database.Ping(stopCtx); err != nil {
				t.Errorf("database was closed before module Stop ran: %v", err)
			}
			record("module_stopped")
			return nil
		},
	}); err != nil {
		t.Fatalf("registry.Register() error = %v", err)
	}
	if err := registry.StartAll(ctx); err != nil {
		t.Fatalf("registry.StartAll() error = %v", err)
	}
	jobRunner.Start(context.Background())

	select {
	case <-jobStarted:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for background job to start")
	}

	app := &App{
		db:       database,
		jobs:     jobRunner,
		registry: registry,
	}

	closeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := app.Close(closeCtx); err != nil {
		t.Fatalf("app.Close() error = %v", err)
	}

	if err := database.Ping(ctx); err == nil {
		t.Fatalf("expected database to be closed after app.Close()")
	}

	mu.Lock()
	got := slices.Clone(events)
	mu.Unlock()

	want := []string{"job_stopped", "module_stopped"}
	if !slices.Equal(got, want) {
		t.Fatalf("teardown order = %v, want %v", got, want)
	}
}

func TestAppWave1AAuthSettingsKeysEndToEnd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cfg := config.Config{
		Host:          "127.0.0.1",
		Port:          0,
		DataDir:       t.TempDir(),
		AdminPassword: "test-admin-pass-2026",
		LogLevel:      "error",
		APITimeout:    5 * time.Second,
	}
	build := system.BuildInfo{
		Version:   "test",
		Commit:    "wave1a",
		BuildTime: "2026-11-23T00:00:00Z",
	}

	app, err := buildApp(ctx, cfg, build)
	if err != nil {
		t.Fatalf("buildApp() error = %v", err)
	}
	t.Cleanup(func() {
		_ = app.Close(context.Background())
	})

	ts := httptest.NewServer(app.Handler())
	defer ts.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 5 * time.Second,
	}

	// 1) GET /api/health includes system, auth, settings, keys in modules.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/api/health", nil)
	if err != nil {
		t.Fatalf("new GET /api/health request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET /api/health failed: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/health status = %d, want %d (body: %s)", resp.StatusCode, http.StatusOK, string(body))
	}
	var health struct {
		Status  string   `json:"status"`
		DB      string   `json:"db"`
		Modules []string `json:"modules"`
	}
	if err := json.Unmarshal(body, &health); err != nil {
		t.Fatalf("decode /api/health JSON: %v", err)
	}
	for _, wantMod := range []string{"system", "auth", "settings", "keys"} {
		if !slices.Contains(health.Modules, wantMod) {
			t.Fatalf("GET /api/health modules = %v, missing %q", health.Modules, wantMod)
		}
	}

	// 2) Unauthenticated GET /api/keys and GET /api/settings return 401 when require_login=true.
	for _, path := range []string{"/api/keys", "/api/settings"} {
		unauthReq, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+path, nil)
		if err != nil {
			t.Fatalf("new unauthenticated GET %s request: %v", path, err)
		}
		unauthResp, err := client.Do(unauthReq)
		if err != nil {
			t.Fatalf("unauthenticated GET %s failed: %v", path, err)
		}
		unauthBody, _ := io.ReadAll(unauthResp.Body)
		_ = unauthResp.Body.Close()
		if unauthResp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("unauthenticated GET %s status = %d, want %d (body: %s)", path, unauthResp.StatusCode, http.StatusUnauthorized, string(unauthBody))
		}
	}

	// 3) Authenticate via POST /api/auth/login -> capture om4_session cookie.
	loginReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ts.URL+"/api/auth/login",
		strings.NewReader(`{"password":"test-admin-pass-2026"}`),
	)
	if err != nil {
		t.Fatalf("new POST /api/auth/login request: %v", err)
	}
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, err := client.Do(loginReq)
	if err != nil {
		t.Fatalf("POST /api/auth/login failed: %v", err)
	}
	loginBody, _ := io.ReadAll(loginResp.Body)
	_ = loginResp.Body.Close()
	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("POST /api/auth/login status = %d, want %d (body: %s)", loginResp.StatusCode, http.StatusOK, string(loginBody))
	}
	var sessionCookie *http.Cookie
	for _, c := range loginResp.Cookies() {
		if c.Name == "om4_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil || sessionCookie.Value == "" {
		t.Fatalf("POST /api/auth/login missing non-empty om4_session cookie: %v", loginResp.Cookies())
	}

	// 4) With cookie: GET /api/settings returns 200.
	getSettingsReq, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/api/settings", nil)
	if err != nil {
		t.Fatalf("new GET /api/settings request: %v", err)
	}
	getSettingsResp, err := client.Do(getSettingsReq)
	if err != nil {
		t.Fatalf("GET /api/settings failed: %v", err)
	}
	getSettingsBody, _ := io.ReadAll(getSettingsResp.Body)
	_ = getSettingsResp.Body.Close()
	if getSettingsResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/settings status = %d, want %d (body: %s)", getSettingsResp.StatusCode, http.StatusOK, string(getSettingsBody))
	}

	// 5) With cookie: PATCH /api/settings updates ui.theme.
	patchSettingsReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		ts.URL+"/api/settings",
		strings.NewReader(`{"updates":{"ui.theme":"light"},"reason":"e2e integration test"}`),
	)
	if err != nil {
		t.Fatalf("new PATCH /api/settings request: %v", err)
	}
	patchSettingsReq.Header.Set("Content-Type", "application/json")
	patchSettingsResp, err := client.Do(patchSettingsReq)
	if err != nil {
		t.Fatalf("PATCH /api/settings failed: %v", err)
	}
	patchSettingsBody, _ := io.ReadAll(patchSettingsResp.Body)
	_ = patchSettingsResp.Body.Close()
	if patchSettingsResp.StatusCode != http.StatusOK {
		t.Fatalf("PATCH /api/settings status = %d, want %d (body: %s)", patchSettingsResp.StatusCode, http.StatusOK, string(patchSettingsBody))
	}
	var settingsPayload struct {
		Settings []struct {
			Key      string `json:"key"`
			Value    any    `json:"value"`
			Revision int64  `json:"revision"`
		} `json:"settings"`
	}
	if err := json.Unmarshal(patchSettingsBody, &settingsPayload); err != nil {
		t.Fatalf("decode PATCH /api/settings JSON: %v", err)
	}
	foundTheme := false
	for _, item := range settingsPayload.Settings {
		if item.Key == "ui.theme" {
			foundTheme = true
			if item.Value != "light" {
				t.Fatalf("ui.theme value = %v, want %q", item.Value, "light")
			}
			if item.Revision < 1 {
				t.Fatalf("ui.theme revision = %d, want >= 1", item.Revision)
			}
		}
	}
	if !foundTheme {
		t.Fatalf("PATCH /api/settings response missing ui.theme: %s", string(patchSettingsBody))
	}

	// 6) With cookie: POST /api/keys creates key.
	createKeyReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ts.URL+"/api/keys",
		strings.NewReader(`{"name":"wave1a-e2e-key"}`),
	)
	if err != nil {
		t.Fatalf("new POST /api/keys request: %v", err)
	}
	createKeyReq.Header.Set("Content-Type", "application/json")
	createKeyResp, err := client.Do(createKeyReq)
	if err != nil {
		t.Fatalf("POST /api/keys failed: %v", err)
	}
	createKeyBody, _ := io.ReadAll(createKeyResp.Body)
	_ = createKeyResp.Body.Close()
	if createKeyResp.StatusCode != http.StatusCreated {
		t.Fatalf("POST /api/keys status = %d, want %d (body: %s)", createKeyResp.StatusCode, http.StatusCreated, string(createKeyBody))
	}
	var createdKey struct {
		Key struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"key"`
		Secret string `json:"secret"`
	}
	if err := json.Unmarshal(createKeyBody, &createdKey); err != nil {
		t.Fatalf("decode POST /api/keys JSON: %v", err)
	}
	if createdKey.Key.ID == "" || createdKey.Key.Name != "wave1a-e2e-key" || !strings.HasPrefix(createdKey.Secret, "sk-om4-") {
		t.Fatalf("unexpected POST /api/keys response: %+v", createdKey)
	}

	// 7) POST /api/auth/logout clears session -> subsequent GET /api/keys returns 401.
	logoutReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.URL+"/api/auth/logout", nil)
	if err != nil {
		t.Fatalf("new POST /api/auth/logout request: %v", err)
	}
	logoutResp, err := client.Do(logoutReq)
	if err != nil {
		t.Fatalf("POST /api/auth/logout failed: %v", err)
	}
	logoutBody, _ := io.ReadAll(logoutResp.Body)
	_ = logoutResp.Body.Close()
	if logoutResp.StatusCode != http.StatusOK {
		t.Fatalf("POST /api/auth/logout status = %d, want %d (body: %s)", logoutResp.StatusCode, http.StatusOK, string(logoutBody))
	}

	postLogoutReq, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/api/keys", nil)
	if err != nil {
		t.Fatalf("new post-logout GET /api/keys request: %v", err)
	}
	postLogoutResp, err := client.Do(postLogoutReq)
	if err != nil {
		t.Fatalf("post-logout GET /api/keys failed: %v", err)
	}
	postLogoutBody, _ := io.ReadAll(postLogoutResp.Body)
	_ = postLogoutResp.Body.Close()
	if postLogoutResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("post-logout GET /api/keys status = %d, want %d (body: %s)", postLogoutResp.StatusCode, http.StatusUnauthorized, string(postLogoutBody))
	}
}

func TestAppWave1BGatewayEndToEnd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cfg := config.Config{
		Host:          "127.0.0.1",
		Port:          0,
		DataDir:       t.TempDir(),
		AdminPassword: "test-admin-pass-wave1b",
		LogLevel:      "error",
		APITimeout:    5 * time.Second,
	}
	build := system.BuildInfo{
		Version:   "test",
		Commit:    "wave1b",
		BuildTime: "2026-11-23T00:00:00Z",
	}

	app, err := buildApp(ctx, cfg, build)
	if err != nil {
		t.Fatalf("buildApp() error = %v", err)
	}
	t.Cleanup(func() {
		_ = app.Close(context.Background())
	})

	ts := httptest.NewServer(app.Handler())
	defer ts.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 5 * time.Second,
	}

	doJSON := func(method, url string, payload any, headers map[string]string) (int, []byte, http.Header) {
		t.Helper()
		var bodyReader io.Reader
		if payload != nil {
			raw, mErr := json.Marshal(payload)
			if mErr != nil {
				t.Fatalf("marshal %s %s payload: %v", method, url, mErr)
			}
			bodyReader = bytes.NewReader(raw)
		}
		req, rErr := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if rErr != nil {
			t.Fatalf("new request %s %s: %v", method, url, rErr)
		}
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, dErr := client.Do(req)
		if dErr != nil {
			t.Fatalf("%s %s failed: %v", method, url, dErr)
		}
		defer resp.Body.Close()
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("read %s %s body: %v", method, url, readErr)
		}
		return resp.StatusCode, respBody, resp.Header.Clone()
	}

	// 1) GET /api/health includes all 7 Wave 1A + 1B modules.
	status, healthBody, _ := doJSON(http.MethodGet, ts.URL+"/api/health", nil, nil)
	if status != http.StatusOK {
		t.Fatalf("GET /api/health status = %d, want 200 (body: %s)", status, string(healthBody))
	}
	var health struct {
		Status  string   `json:"status"`
		DB      string   `json:"db"`
		Modules []string `json:"modules"`
	}
	if err := json.Unmarshal(healthBody, &health); err != nil {
		t.Fatalf("decode /api/health JSON: %v", err)
	}
	for _, wantMod := range []string{"system", "auth", "settings", "keys", "providers", "combos", "gateway"} {
		if !slices.Contains(health.Modules, wantMod) {
			t.Fatalf("GET /api/health modules = %v, missing %q", health.Modules, wantMod)
		}
	}

	// 2) Authenticate admin via POST /api/auth/login.
	status, loginBody, _ := doJSON(http.MethodPost, ts.URL+"/api/auth/login", map[string]string{
		"password": "test-admin-pass-wave1b",
	}, nil)
	if status != http.StatusOK {
		t.Fatalf("POST /api/auth/login status = %d, want 200 (body: %s)", status, string(loginBody))
	}

	// 3) Start mock OpenAI upstream (returns 429 Retry-After: 60) and fallback mock Anthropic SSE upstream.
	var openaiHits atomic.Int32
	openaiMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		openaiHits.Add(1)
		if r.URL.Path != "/chat/completions" {
			t.Errorf("openai mock path = %q, want /chat/completions", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer oai-e2e-key" {
			t.Errorf("openai mock Authorization = %q, want Bearer oai-e2e-key", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = io.WriteString(w, `{"error":{"message":"rate limit exceeded","type":"rate_limit_error"}}`)
	}))
	defer openaiMock.Close()

	var anthropicHits atomic.Int32
	anthropicMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		anthropicHits.Add(1)
		if r.URL.Path != "/v1/messages" {
			t.Errorf("anthropic mock path = %q, want /v1/messages", r.URL.Path)
		}
		if got := r.Header.Get("x-api-key"); got != "ant-e2e-key" {
			t.Errorf("anthropic mock x-api-key = %q, want ant-e2e-key", got)
		}
		if got := r.Header.Get("anthropic-version"); got != "2023-06-01" {
			t.Errorf("anthropic mock anthropic-version = %q, want 2023-06-01", got)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		rc := http.NewResponseController(w)
		events := []string{
			"event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"id\":\"msg_e2e_1\",\"type\":\"message\",\"role\":\"assistant\",\"model\":\"claude-sonnet-4-20250514\",\"content\":[],\"usage\":{\"input_tokens\":12,\"output_tokens\":1}}}\n\n",
			"event: content_block_start\ndata: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"text\",\"text\":\"\"}}\n\n",
			"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"Failover from OpenAI 429 succeeded\"}}\n\n",
			"event: content_block_stop\ndata: {\"type\":\"content_block_stop\",\"index\":0}\n\n",
			"event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"},\"usage\":{\"output_tokens\":8}}\n\n",
			"event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n",
		}
		for _, ev := range events {
			_, _ = io.WriteString(w, ev)
			_ = rc.Flush()
		}
	}))
	defer anthropicMock.Close()

	// Create primary OpenAI connection.
	status, oaiConnBody, _ := doJSON(http.MethodPost, ts.URL+"/api/providers/connections", map[string]any{
		"provider_id": "openai",
		"name":        "primary-openai-mock",
		"base_url":    openaiMock.URL,
		"api_key":     "oai-e2e-key",
		"priority":    1,
		"enabled":     true,
		"models":      []string{"gpt-4o"},
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("create openai connection status = %d, want 201 (body: %s)", status, string(oaiConnBody))
	}

	// Create fallback Anthropic connection.
	status, antConnBody, _ := doJSON(http.MethodPost, ts.URL+"/api/providers/connections", map[string]any{
		"provider_id": "anthropic",
		"name":        "fallback-anthropic-mock",
		"base_url":    anthropicMock.URL,
		"api_key":     "ant-e2e-key",
		"priority":    2,
		"enabled":     true,
		"models":      []string{"claude-sonnet-4-20250514"},
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("create anthropic connection status = %d, want 201 (body: %s)", status, string(antConnBody))
	}

	// Create ordered combo spanning both connections.
	const comboName = "resilient-chat"
	status, comboBody, _ := doJSON(http.MethodPost, ts.URL+"/api/combos", map[string]any{
		"name":        comboName,
		"description": "E2E primary OpenAI with Anthropic fallback",
		"targets":     []string{"openai/gpt-4o", "anthropic/claude-sonnet-4-20250514"},
		"strategy":    "ordered",
		"enabled":     true,
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("create combo status = %d, want 201 (body: %s)", status, string(comboBody))
	}

	// Create gateway API key scoped strictly to the combo.
	status, keyBody, _ := doJSON(http.MethodPost, ts.URL+"/api/keys", map[string]any{
		"name":           "wave1b-scoped-key",
		"allowed_models": []string{comboName},
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("create scoped API key status = %d, want 201 (body: %s)", status, string(keyBody))
	}
	var createdKey struct {
		Secret string `json:"secret"`
	}
	if err := json.Unmarshal(keyBody, &createdKey); err != nil {
		t.Fatalf("decode POST /api/keys JSON: %v", err)
	}
	if !strings.HasPrefix(createdKey.Secret, "sk-om4-") {
		t.Fatalf("unexpected API key secret prefix: %q", createdKey.Secret)
	}

	authHeader := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", createdKey.Secret),
	}

	// 4) Send streaming POST /v1/chat/completions -> failover 429 -> Anthropic SSE -> OpenAI SSE + [DONE].
	status, chatBody, chatHeaders := doJSON(http.MethodPost, ts.URL+"/v1/chat/completions", map[string]any{
		"model":  comboName,
		"stream": true,
		"messages": []map[string]string{
			{"role": "user", "content": "Hello failover"},
		},
	}, authHeader)
	if status != http.StatusOK {
		t.Fatalf("POST /v1/chat/completions status = %d, want 200 (body: %s)", status, string(chatBody))
	}
	if ct := chatHeaders.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("POST /v1/chat/completions Content-Type = %q, want text/event-stream", ct)
	}
	sseText := string(chatBody)
	if !strings.Contains(sseText, "chat.completion.chunk") ||
		!strings.Contains(sseText, "Failover from OpenAI 429 succeeded") ||
		!strings.Contains(sseText, "data: [DONE]") {
		t.Fatalf("unexpected OpenAI SSE translated stream: %s", sseText)
	}
	if openaiHits.Load() != 1 || anthropicHits.Load() != 1 {
		t.Fatalf("upstream hit counts: openai=%d anthropic=%d, want 1 and 1", openaiHits.Load(), anthropicHits.Load())
	}

	// 5) GET /v1/models with scoped API key returns only the authorized combo.
	status, modelsBody, _ := doJSON(http.MethodGet, ts.URL+"/v1/models", nil, authHeader)
	if status != http.StatusOK {
		t.Fatalf("GET /v1/models status = %d, want 200 (body: %s)", status, string(modelsBody))
	}
	var modelsResp struct {
		Object string `json:"object"`
		Data   []struct {
			ID      string `json:"id"`
			OwnedBy string `json:"owned_by"`
		} `json:"data"`
	}
	if err := json.Unmarshal(modelsBody, &modelsResp); err != nil {
		t.Fatalf("decode GET /v1/models JSON: %v", err)
	}
	if len(modelsResp.Data) != 1 || modelsResp.Data[0].ID != comboName {
		t.Fatalf("GET /v1/models data = %+v, want single model %q", modelsResp.Data, comboName)
	}

	// 6) Unauthorized model returns 403 Forbidden.
	status, unauthModelBody, _ := doJSON(http.MethodPost, ts.URL+"/v1/chat/completions", map[string]any{
		"model":  "openai/gpt-4o",
		"stream": false,
		"messages": []map[string]string{
			{"role": "user", "content": "Should be forbidden"},
		},
	}, authHeader)
	if status != http.StatusForbidden {
		t.Fatalf("unauthorized model status = %d, want 403 (body: %s)", status, string(unauthModelBody))
	}
}

func TestAppWave1CUsageAndCLIToolsEndToEnd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cfg := config.Config{
		Host:          "127.0.0.1",
		Port:          0,
		DataDir:       t.TempDir(),
		AdminPassword: "test-admin-pass-wave1c",
		LogLevel:      "error",
		APITimeout:    5 * time.Second,
	}
	build := system.BuildInfo{
		Version:   "test",
		Commit:    "wave1c",
		BuildTime: "2026-11-23T00:00:00Z",
	}

	app, err := buildApp(ctx, cfg, build)
	if err != nil {
		t.Fatalf("buildApp() error = %v", err)
	}
	t.Cleanup(func() {
		_ = app.Close(context.Background())
	})

	ts := httptest.NewServer(app.Handler())
	defer ts.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 5 * time.Second,
	}

	doJSON := func(method, url string, payload any, headers map[string]string) (int, []byte) {
		t.Helper()
		var bodyReader io.Reader
		if payload != nil {
			raw, mErr := json.Marshal(payload)
			if mErr != nil {
				t.Fatalf("marshal %s %s payload: %v", method, url, mErr)
			}
			bodyReader = bytes.NewReader(raw)
		}
		req, rErr := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if rErr != nil {
			t.Fatalf("new request %s %s: %v", method, url, rErr)
		}
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, dErr := client.Do(req)
		if dErr != nil {
			t.Fatalf("%s %s failed: %v", method, url, dErr)
		}
		defer resp.Body.Close()
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("read %s %s body: %v", method, url, readErr)
		}
		return resp.StatusCode, respBody
	}

	// 1) GET /api/health lists all 9 Wave 1 modules.
	status, healthBody := doJSON(http.MethodGet, ts.URL+"/api/health", nil, nil)
	if status != http.StatusOK {
		t.Fatalf("GET /api/health status = %d, want 200 (body: %s)", status, string(healthBody))
	}
	var health struct {
		Status  string   `json:"status"`
		DB      string   `json:"db"`
		Modules []string `json:"modules"`
	}
	if err := json.Unmarshal(healthBody, &health); err != nil {
		t.Fatalf("decode /api/health JSON: %v", err)
	}
	wantModules := []string{"auth", "clitools", "combos", "gateway", "keys", "providers", "settings", "system", "usage"}
	for _, wantMod := range wantModules {
		if !slices.Contains(health.Modules, wantMod) {
			t.Fatalf("GET /api/health modules = %v, missing %q", health.Modules, wantMod)
		}
	}
	if len(health.Modules) != len(wantModules) {
		t.Fatalf("GET /api/health modules count = %d (%v), want %d (%v)", len(health.Modules), health.Modules, len(wantModules), wantModules)
	}

	// 2) Authenticate admin.
	status, loginBody := doJSON(http.MethodPost, ts.URL+"/api/auth/login", map[string]string{
		"password": "test-admin-pass-wave1c",
	}, nil)
	if status != http.StatusOK {
		t.Fatalf("POST /api/auth/login status = %d, want 200 (body: %s)", status, string(loginBody))
	}

	// 3) Mock OpenAI upstream returning usage: prompt_tokens=120, completion_tokens=45.
	openaiMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("openai mock path = %q, want /chat/completions", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"id":"chatcmpl-wave1c","object":"chat.completion","created":1732320000,"model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"wave1c telemetry ok"},"finish_reason":"stop"}],"usage":{"prompt_tokens":120,"completion_tokens":45,"total_tokens":165}}`)
	}))
	defer openaiMock.Close()

	status, connBody := doJSON(http.MethodPost, ts.URL+"/api/providers/connections", map[string]any{
		"provider_id": "openai",
		"name":        "wave1c-openai-mock",
		"base_url":    openaiMock.URL,
		"api_key":     "oai-wave1c-key",
		"priority":    1,
		"enabled":     true,
		"models":      []string{"gpt-4o"},
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("POST /api/providers/connections status = %d, want 201 (body: %s)", status, string(connBody))
	}

	status, keyBody := doJSON(http.MethodPost, ts.URL+"/api/keys", map[string]any{
		"name": "wave1c-e2e-key",
	}, nil)
	if status != http.StatusCreated {
		t.Fatalf("POST /api/keys status = %d, want 201 (body: %s)", status, string(keyBody))
	}
	var createdKey struct {
		Key struct {
			ID string `json:"id"`
		} `json:"key"`
		Secret string `json:"secret"`
	}
	if err := json.Unmarshal(keyBody, &createdKey); err != nil {
		t.Fatalf("decode POST /api/keys JSON: %v", err)
	}
	if createdKey.Key.ID == "" || !strings.HasPrefix(createdKey.Secret, "sk-om4-") {
		t.Fatalf("unexpected POST /api/keys response: %+v", createdKey)
	}

	status, chatBody := doJSON(http.MethodPost, ts.URL+"/v1/chat/completions", map[string]any{
		"model":  "openai/gpt-4o",
		"stream": false,
		"messages": []map[string]string{
			{"role": "user", "content": "Hello Wave 1C"},
		},
	}, map[string]string{
		"Authorization": "Bearer " + createdKey.Secret,
	})
	if status != http.StatusOK {
		t.Fatalf("POST /v1/chat/completions status = %d, want 200 (body: %s)", status, string(chatBody))
	}

	// 4) Poll GET /api/usage/logs until async batcher flushes, then verify log and summary.
	type usageLog struct {
		Model        string  `json:"model"`
		ProviderID   string  `json:"provider_id"`
		KeyID        string  `json:"key_id"`
		InputTokens  int64   `json:"input_tokens"`
		OutputTokens int64   `json:"output_tokens"`
		CostUSD      float64 `json:"cost_usd"`
		StatusCode   int     `json:"status_code"`
	}
	type usageLogsResp struct {
		Logs  []usageLog `json:"logs"`
		Total int64      `json:"total"`
	}

	var logsPayload usageLogsResp
	pollDeadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(pollDeadline) {
		lStatus, lBody := doJSON(http.MethodGet, ts.URL+"/api/usage/logs", nil, nil)
		if lStatus != http.StatusOK {
			t.Fatalf("GET /api/usage/logs status = %d, want 200 (body: %s)", lStatus, string(lBody))
		}
		if err := json.Unmarshal(lBody, &logsPayload); err != nil {
			t.Fatalf("decode GET /api/usage/logs JSON: %v", err)
		}
		if len(logsPayload.Logs) > 0 {
			break
		}
		time.Sleep(25 * time.Millisecond)
	}
	if len(logsPayload.Logs) == 0 {
		t.Fatalf("timed out waiting for async usage batcher to flush logs")
	}
	gotLog := logsPayload.Logs[0]
	if gotLog.Model != "openai/gpt-4o" ||
		gotLog.ProviderID != "openai" ||
		gotLog.KeyID != createdKey.Key.ID ||
		gotLog.InputTokens != 120 ||
		gotLog.OutputTokens != 45 ||
		gotLog.StatusCode != http.StatusOK ||
		gotLog.CostUSD <= 0 {
		t.Fatalf("unexpected usage log entry: %+v", gotLog)
	}

	status, summaryBody := doJSON(http.MethodGet, ts.URL+"/api/usage/summary?range=24h", nil, nil)
	if status != http.StatusOK {
		t.Fatalf("GET /api/usage/summary status = %d, want 200 (body: %s)", status, string(summaryBody))
	}
	var summary struct {
		Totals struct {
			Requests     int64   `json:"requests"`
			InputTokens  int64   `json:"input_tokens"`
			OutputTokens int64   `json:"output_tokens"`
			CostUSD      float64 `json:"cost_usd"`
		} `json:"totals"`
	}
	if err := json.Unmarshal(summaryBody, &summary); err != nil {
		t.Fatalf("decode GET /api/usage/summary JSON: %v", err)
	}
	if summary.Totals.Requests < 1 ||
		summary.Totals.InputTokens < 120 ||
		summary.Totals.OutputTokens < 45 ||
		summary.Totals.CostUSD <= 0 {
		t.Fatalf("unexpected usage summary totals: %+v", summary.Totals)
	}

	// 5) Verify GET /api/clitools and POST /api/clitools/claude-code/generate.
	status, toolsBody := doJSON(http.MethodGet, ts.URL+"/api/clitools", nil, nil)
	if status != http.StatusOK {
		t.Fatalf("GET /api/clitools status = %d, want 200 (body: %s)", status, string(toolsBody))
	}
	var toolsCatalog struct {
		LocalApplyEnabled bool `json:"local_apply_enabled"`
		Tools             []struct {
			Tool string `json:"tool"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(toolsBody, &toolsCatalog); err != nil {
		t.Fatalf("decode GET /api/clitools JSON: %v", err)
	}
	if !toolsCatalog.LocalApplyEnabled {
		t.Fatalf("expected local_apply_enabled=true for loopback host")
	}
	if len(toolsCatalog.Tools) != 6 {
		t.Fatalf("GET /api/clitools count = %d, want 6 (%+v)", len(toolsCatalog.Tools), toolsCatalog.Tools)
	}

	status, genBody := doJSON(http.MethodPost, ts.URL+"/api/clitools/claude-code/generate", map[string]any{
		"base_url": ts.URL,
		"api_key":  createdKey.Secret,
		"model":    "openai/gpt-4o",
	}, nil)
	if status != http.StatusOK {
		t.Fatalf("POST /api/clitools/claude-code/generate status = %d, want 200 (body: %s)", status, string(genBody))
	}
	var genResp struct {
		Tool    string `json:"tool"`
		Format  string `json:"format"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(genBody, &genResp); err != nil {
		t.Fatalf("decode POST /api/clitools/claude-code/generate JSON: %v", err)
	}
	if genResp.Tool != "claude-code" || genResp.Format != "json" {
		t.Fatalf("unexpected generate response metadata: %+v", genResp)
	}
	var claudeCfg struct {
		Env map[string]string `json:"env"`
	}
	if err := json.Unmarshal([]byte(genResp.Content), &claudeCfg); err != nil {
		t.Fatalf("generated claude-code content is not valid JSON: %v (content: %s)", err, genResp.Content)
	}
	if claudeCfg.Env["ANTHROPIC_BASE_URL"] != ts.URL {
		t.Fatalf("ANTHROPIC_BASE_URL = %q, want %q", claudeCfg.Env["ANTHROPIC_BASE_URL"], ts.URL)
	}
	if claudeCfg.Env["ANTHROPIC_AUTH_TOKEN"] != createdKey.Secret {
		t.Fatalf("ANTHROPIC_AUTH_TOKEN = %q, want %q", claudeCfg.Env["ANTHROPIC_AUTH_TOKEN"], createdKey.Secret)
	}
}

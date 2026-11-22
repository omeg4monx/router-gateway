package system_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/system"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubAdminAuth struct {
	requireLogin bool
	validToken   string
}

func (a *stubAdminAuth) AuthenticateAdmin(r *http.Request) (string, error) {
	if r.Header.Get("Authorization") == "Bearer "+a.validToken {
		return "admin", nil
	}
	return "", errors.New("unauthorized")
}

func (a *stubAdminAuth) RequireLogin(context.Context) bool {
	return a.requireLogin
}

type testHarness struct {
	db     *db.DB
	reg    *lifecycle.Registry
	server *httpserver.Server
	mod    *system.Module
}

func newTestHarness(t *testing.T, cfg config.Config, build system.BuildInfo, activeModules func() []string, shutdown func(), auth *stubAdminAuth) *testHarness {
	t.Helper()

	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.Port == 0 {
		cfg.Port = 20128
	}
	if cfg.APITimeout == 0 {
		cfg.APITimeout = 30 * time.Second
	}

	dbPath := filepath.Join(t.TempDir(), "router.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})

	reg := lifecycle.NewRegistry(nil)
	if activeModules == nil {
		activeModules = reg.ActiveNames
	}

	mod := system.New(system.Options{
		DB:            database,
		Config:        cfg,
		Build:         build,
		ActiveModules: activeModules,
		Shutdown:      shutdown,
	})

	if err := reg.Register(mod); err != nil {
		t.Fatalf("reg.Register: %v", err)
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

	return &testHarness{
		db:     database,
		reg:    reg,
		server: srv,
		mod:    mod,
	}
}

func TestSystemHealthAndVersion(t *testing.T) {
	build := system.BuildInfo{
		Version:   "0.1.0-dev",
		Commit:    "abc",
		BuildTime: "2026-11-22T00:00:00Z",
	}
	h := newTestHarness(t, config.Config{Host: "127.0.0.1"}, build, nil, nil, nil)
	if err := h.reg.StartAll(context.Background()); err != nil {
		t.Fatalf("StartAll: %v", err)
	}
	t.Cleanup(func() {
		_ = h.reg.StopAll(context.Background())
	})

	t.Run("GET /api/health healthy", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/health status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
		}

		var body struct {
			Status        string   `json:"status"`
			DB            string   `json:"db"`
			Modules       []string `json:"modules"`
			UptimeSeconds int64    `json:"uptimeSeconds"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshal health response: %v (raw=%s)", err, rec.Body.String())
		}
		if body.Status != "ok" {
			t.Errorf("status = %q, want ok", body.Status)
		}
		if body.DB != "ok" {
			t.Errorf("db = %q, want ok", body.DB)
		}
		if len(body.Modules) != 1 || body.Modules[0] != "system" {
			t.Errorf("modules = %v, want [system]", body.Modules)
		}
		if body.UptimeSeconds < 0 {
			t.Errorf("uptimeSeconds = %d, want >= 0", body.UptimeSeconds)
		}
	})

	t.Run("GET /api/version", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
		rec := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/version status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
		}

		var body struct {
			Version   string `json:"version"`
			Commit    string `json:"commit"`
			BuildTime string `json:"buildTime"`
			GoVersion string `json:"goVersion"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshal version response: %v (raw=%s)", err, rec.Body.String())
		}
		if body.Version != "0.1.0-dev" {
			t.Errorf("version = %q, want 0.1.0-dev", body.Version)
		}
		if body.Commit != "abc" {
			t.Errorf("commit = %q, want abc", body.Commit)
		}
		if body.BuildTime != "2026-11-22T00:00:00Z" {
			t.Errorf("buildTime = %q, want 2026-11-22T00:00:00Z", body.BuildTime)
		}
		if body.GoVersion != runtime.Version() {
			t.Errorf("goVersion = %q, want %q", body.GoVersion, runtime.Version())
		}
	})

	t.Run("GET /api/health returns empty JSON array when no modules active", func(t *testing.T) {
		hEmpty := newTestHarness(t, config.Config{Host: "127.0.0.1"}, build, func() []string { return nil }, nil, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()
		hEmpty.server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), `"modules":[]`) {
			t.Fatalf("expected non-null empty modules array, got body=%s", rec.Body.String())
		}
	})

	t.Run("GET /api/health returns 503 degraded when DB is closed", func(t *testing.T) {
		hBad := newTestHarness(t, config.Config{Host: "127.0.0.1"}, build, nil, nil, nil)
		if err := hBad.db.Close(); err != nil {
			t.Fatalf("db.Close: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()
		hBad.server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("GET /api/health on closed DB status = %d, want 503 (body=%s)", rec.Code, rec.Body.String())
		}

		var body struct {
			Status string `json:"status"`
			DB     string `json:"db"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshal degraded health response: %v", err)
		}
		if body.Status != "degraded" || body.DB != "error" {
			t.Errorf("got status=%q db=%q, want degraded/error", body.Status, body.DB)
		}
	})
}

func TestSystemInitLifecycle(t *testing.T) {
	t.Run("interactive init flow", func(t *testing.T) {
		h := newTestHarness(t, config.Config{Host: "127.0.0.1"}, system.BuildInfo{}, nil, nil, nil)
		if err := h.reg.StartAll(context.Background()); err != nil {
			t.Fatalf("StartAll: %v", err)
		}
		t.Cleanup(func() {
			_ = h.reg.StopAll(context.Background())
		})

		// Verify migration version 1 is recorded.
		versions, err := h.db.AppliedVersions(context.Background(), "system")
		if err != nil {
			t.Fatalf("AppliedVersions: %v", err)
		}
		if len(versions) != 1 || versions[0] != 1 {
			t.Fatalf("applied versions = %v, want [1]", versions)
		}

		// 1) GET /api/init -> {initialized:false, requireLogin:true}
		req := httptest.NewRequest(http.MethodGet, "/api/init", nil)
		rec := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/init status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
		}
		var initState struct {
			Initialized  bool `json:"initialized"`
			RequireLogin bool `json:"requireLogin"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &initState); err != nil {
			t.Fatalf("unmarshal GET /api/init: %v", err)
		}
		if initState.Initialized {
			t.Errorf("initialized = true, want false")
		}
		if !initState.RequireLogin {
			t.Errorf("requireLogin = false, want true")
		}

		// 2) POST /api/init with invalid inputs -> 400
		invalidBodies := []struct {
			name string
			body string
		}{
			{"short password", `{"password":"short"}`},
			{"too long password", `{"password":"` + strings.Repeat("a", 73) + `"}`},
			{"malformed json", `{"password":`},
			{"unknown field", `{"password":"validpassword123","extra":true}`},
			{"empty body", ``},
		}
		for _, tc := range invalidBodies {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/api/init", strings.NewReader(tc.body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				h.server.Handler().ServeHTTP(rec, req)
				if rec.Code != http.StatusBadRequest {
					t.Fatalf("POST /api/init (%s) status = %d, want 400 (body=%s)", tc.name, rec.Code, rec.Body.String())
				}
				var env httpserver.ErrorEnvelope
				if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
					t.Fatalf("unmarshal error envelope: %v", err)
				}
				if env.Error.Code == "" || env.Error.Message == "" {
					t.Fatalf("expected structured error envelope, got %+v", env)
				}
			})
		}

		// 3) POST /api/init with valid password -> 200 {initialized:true}
		const adminPass = "super-secret-pass-2026"
		req = httptest.NewRequest(http.MethodPost, "/api/init", strings.NewReader(`{"password":"`+adminPass+`"}`))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("POST /api/init status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
		}
		var initResp struct {
			Initialized bool `json:"initialized"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &initResp); err != nil {
			t.Fatalf("unmarshal POST /api/init: %v", err)
		}
		if !initResp.Initialized {
			t.Fatalf("POST /api/init initialized = false, want true")
		}

		// Confirm stored bcrypt hash in system_state verifies adminPass.
		var storedHash string
		if err := h.db.Reader().QueryRowContext(context.Background(), `SELECT value FROM system_state WHERE key = 'admin_password_hash'`).Scan(&storedHash); err != nil {
			t.Fatalf("query admin_password_hash: %v", err)
		}
		if !crypto.VerifyPassword(storedHash, adminPass) {
			t.Fatalf("stored bcrypt hash failed verification for adminPass")
		}

		// Confirm GET /api/init now reports initialized=true.
		req = httptest.NewRequest(http.MethodGet, "/api/init", nil)
		rec = httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/init after init status = %d, want 200", rec.Code)
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &initState); err != nil {
			t.Fatalf("unmarshal GET /api/init after init: %v", err)
		}
		if !initState.Initialized || !initState.RequireLogin {
			t.Fatalf("got initialized=%v requireLogin=%v, want true/true", initState.Initialized, initState.RequireLogin)
		}

		// 4) Second POST /api/init -> 409 Conflict
		req = httptest.NewRequest(http.MethodPost, "/api/init", strings.NewReader(`{"password":"another-valid-password"}`))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusConflict {
			t.Fatalf("second POST /api/init status = %d, want 409 (body=%s)", rec.Code, rec.Body.String())
		}
		var env httpserver.ErrorEnvelope
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("unmarshal 409 error envelope: %v", err)
		}
		if env.Error.Code == "" {
			t.Fatalf("expected non-empty error code in 409 response")
		}

		// Ensure original hash was not overwritten.
		var hashAfter string
		if err := h.db.Reader().QueryRowContext(context.Background(), `SELECT value FROM system_state WHERE key = 'admin_password_hash'`).Scan(&hashAfter); err != nil {
			t.Fatalf("query admin_password_hash after 409: %v", err)
		}
		if hashAfter != storedHash {
			t.Fatalf("admin_password_hash was mutated by second init attempt")
		}
	})

	t.Run("Config.AdminPassword auto-bootstrap in Start", func(t *testing.T) {
		const envPass = "env-bootstrapped-password"
		h := newTestHarness(t, config.Config{
			Host:          "127.0.0.1",
			AdminPassword: envPass,
		}, system.BuildInfo{}, nil, nil, nil)

		if err := h.reg.StartAll(context.Background()); err != nil {
			t.Fatalf("StartAll with AdminPassword: %v", err)
		}
		t.Cleanup(func() {
			_ = h.reg.StopAll(context.Background())
		})

		req := httptest.NewRequest(http.MethodGet, "/api/init", nil)
		rec := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/init status = %d, want 200", rec.Code)
		}
		var initState struct {
			Initialized  bool `json:"initialized"`
			RequireLogin bool `json:"requireLogin"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &initState); err != nil {
			t.Fatalf("unmarshal GET /api/init: %v", err)
		}
		if !initState.Initialized || !initState.RequireLogin {
			t.Fatalf("got initialized=%v requireLogin=%v, want true/true", initState.Initialized, initState.RequireLogin)
		}

		var storedHash string
		if err := h.db.Reader().QueryRowContext(context.Background(), `SELECT value FROM system_state WHERE key = 'admin_password_hash'`).Scan(&storedHash); err != nil {
			t.Fatalf("query admin_password_hash: %v", err)
		}
		if !crypto.VerifyPassword(storedHash, envPass) {
			t.Fatalf("stored hash does not match env password")
		}

		// Calling Start again directly must be idempotent and keep existing hash.
		if err := h.mod.Start(context.Background()); err != nil {
			t.Fatalf("idempotent Start: %v", err)
		}
		var hashAfter string
		if err := h.db.Reader().QueryRowContext(context.Background(), `SELECT value FROM system_state WHERE key = 'admin_password_hash'`).Scan(&hashAfter); err != nil {
			t.Fatalf("query admin_password_hash after second Start: %v", err)
		}
		if hashAfter != storedHash {
			t.Fatalf("second Start overwrote existing admin_password_hash")
		}
	})

	t.Run("Start fails fast on invalid Config.AdminPassword", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			pass string
		}{
			{"too short", "1234567"},
			{"too long", strings.Repeat("x", 73)},
		} {
			t.Run(tc.name, func(t *testing.T) {
				h := newTestHarness(t, config.Config{
					Host:          "127.0.0.1",
					AdminPassword: tc.pass,
				}, system.BuildInfo{}, nil, nil, nil)
				if err := h.reg.StartAll(context.Background()); err == nil {
					t.Fatalf("expected StartAll error for %s AdminPassword, got nil", tc.name)
				}
			})
		}
	})
}

func TestSystemInitConcurrentRace(t *testing.T) {
	h := newTestHarness(t, config.Config{Host: "127.0.0.1"}, system.BuildInfo{}, nil, nil, nil)
	if err := h.reg.StartAll(context.Background()); err != nil {
		t.Fatalf("StartAll: %v", err)
	}
	t.Cleanup(func() {
		_ = h.reg.StopAll(context.Background())
	})

	const workers = 10
	start := make(chan struct{})
	var wg sync.WaitGroup
	var okCount atomic.Int32
	var conflictCount atomic.Int32
	var otherCount atomic.Int32
	var winningPassMu sync.Mutex
	var winningPass string

	for i := 0; i < workers; i++ {
		pass := "concurrent-pass-" + string(rune('a'+i)) + "-2026"
		wg.Add(1)
		go func(password string) {
			defer wg.Done()
			<-start

			payload := bytes.NewBufferString(`{"password":"` + password + `"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/init", payload)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.server.Handler().ServeHTTP(rec, req)

			switch rec.Code {
			case http.StatusOK:
				okCount.Add(1)
				winningPassMu.Lock()
				winningPass = password
				winningPassMu.Unlock()
			case http.StatusConflict:
				conflictCount.Add(1)
			default:
				otherCount.Add(1)
				t.Errorf("unexpected status %d: %s", rec.Code, rec.Body.String())
			}
		}(pass)
	}

	close(start)
	wg.Wait()

	if okCount.Load() != 1 {
		t.Fatalf("200 OK count = %d, want exactly 1", okCount.Load())
	}
	if conflictCount.Load() != workers-1 {
		t.Fatalf("409 Conflict count = %d, want %d", conflictCount.Load(), workers-1)
	}
	if otherCount.Load() != 0 {
		t.Fatalf("unexpected status count = %d, want 0", otherCount.Load())
	}

	var storedHash string
	if err := h.db.Reader().QueryRowContext(context.Background(), `SELECT value FROM system_state WHERE key = 'admin_password_hash'`).Scan(&storedHash); err != nil {
		t.Fatalf("query admin_password_hash: %v", err)
	}
	if !crypto.VerifyPassword(storedHash, winningPass) {
		t.Fatalf("stored hash does not verify winning password %q", winningPass)
	}
}

func TestSystemShutdownRemoteBindGuard(t *testing.T) {
	t.Run("loopback allows shutdown and invokes callback once", func(t *testing.T) {
		var calls atomic.Int32
		done := make(chan struct{})
		h := newTestHarness(t, config.Config{
			Host: "127.0.0.1",
		}, system.BuildInfo{}, nil, func() {
			if calls.Add(1) == 1 {
				close(done)
			}
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/shutdown", nil)
		rec := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusAccepted {
			t.Fatalf("POST /api/shutdown status = %d, want 202 (body=%s)", rec.Code, rec.Body.String())
		}
		var resp struct {
			Status string `json:"status"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal shutdown response: %v", err)
		}
		if resp.Status != "shutting_down" {
			t.Errorf("status = %q, want shutting_down", resp.Status)
		}

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatalf("shutdown callback was not invoked")
		}

		// Calling shutdown a second time must not invoke callback again.
		req2 := httptest.NewRequest(http.MethodPost, "/api/shutdown", nil)
		rec2 := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec2, req2)
		if rec2.Code != http.StatusAccepted {
			t.Fatalf("second POST /api/shutdown status = %d, want 202", rec2.Code)
		}
		time.Sleep(50 * time.Millisecond)
		if calls.Load() != 1 {
			t.Fatalf("shutdown callback invocations = %d, want 1", calls.Load())
		}
	})

	t.Run("non-loopback without AllowInsecureRemote returns 403 and does not call callback", func(t *testing.T) {
		var called atomic.Bool
		auth := &stubAdminAuth{requireLogin: true, validToken: "admin-token"}
		h := newTestHarness(t, config.Config{
			Host:                "0.0.0.0",
			AllowInsecureRemote: false,
		}, system.BuildInfo{}, nil, func() {
			called.Store(true)
		}, auth)

		req := httptest.NewRequest(http.MethodPost, "/api/shutdown", nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		rec := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("POST /api/shutdown on 0.0.0.0 status = %d, want 403 (body=%s)", rec.Code, rec.Body.String())
		}
		var env httpserver.ErrorEnvelope
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("unmarshal 403 response: %v", err)
		}
		if env.Error.Code == "" {
			t.Errorf("expected non-empty error code in 403 response")
		}
		time.Sleep(50 * time.Millisecond)
		if called.Load() {
			t.Fatalf("shutdown callback was invoked despite 403 forbidden guard")
		}
	})

	t.Run("non-loopback with AllowInsecureRemote returns 202 and invokes callback", func(t *testing.T) {
		done := make(chan struct{})
		h := newTestHarness(t, config.Config{
			Host:                "0.0.0.0",
			AllowInsecureRemote: true,
		}, system.BuildInfo{}, nil, func() {
			close(done)
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/shutdown", nil)
		rec := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusAccepted {
			t.Fatalf("POST /api/shutdown with AllowInsecureRemote status = %d, want 202 (body=%s)", rec.Code, rec.Body.String())
		}
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatalf("shutdown callback was not invoked with AllowInsecureRemote=true")
		}
	})
}

func TestSystemAdminCredentialStore(t *testing.T) {
	h := newTestHarness(t, config.Config{Host: "127.0.0.1"}, system.BuildInfo{}, nil, nil, nil)
	ctx := context.Background()
	if err := h.reg.StartAll(ctx); err != nil {
		t.Fatalf("StartAll: %v", err)
	}
	t.Cleanup(func() {
		_ = h.reg.StopAll(context.Background())
	})

	var store contracts.AdminCredentialStore = h.mod

	// Before init: AdminPasswordHash and UpdateAdminPasswordHash return ErrNotInitialized; RequireLogin returns true, nil.
	if _, err := store.AdminPasswordHash(ctx); !errors.Is(err, contracts.ErrNotInitialized) {
		t.Fatalf("AdminPasswordHash before init err = %v, want contracts.ErrNotInitialized", err)
	}
	if err := store.UpdateAdminPasswordHash(ctx, "$2a$10$notinitializedplaceholderhash"); !errors.Is(err, contracts.ErrNotInitialized) {
		t.Fatalf("UpdateAdminPasswordHash before init err = %v, want contracts.ErrNotInitialized", err)
	}
	reqLogin, err := store.RequireLogin(ctx)
	if err != nil || !reqLogin {
		t.Fatalf("RequireLogin before init = (%v, %v), want (true, nil)", reqLogin, err)
	}

	// Initialize via POST /api/init.
	const initialPass = "initial-admin-password-2026"
	req := httptest.NewRequest(http.MethodPost, "/api/init", strings.NewReader(`{"password":"`+initialPass+`"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/init status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}

	hash, err := store.AdminPasswordHash(ctx)
	if err != nil {
		t.Fatalf("AdminPasswordHash after init err = %v", err)
	}
	if !crypto.VerifyPassword(hash, initialPass) {
		t.Fatalf("AdminPasswordHash returned hash that does not verify initial password")
	}

	// UpdateAdminPasswordHash updates stored hash atomically.
	const updatedPass = "updated-admin-password-2026"
	newHash, err := crypto.HashPassword(updatedPass)
	if err != nil {
		t.Fatalf("crypto.HashPassword: %v", err)
	}
	if err := store.UpdateAdminPasswordHash(ctx, newHash); err != nil {
		t.Fatalf("UpdateAdminPasswordHash err = %v", err)
	}
	gotHash, err := store.AdminPasswordHash(ctx)
	if err != nil {
		t.Fatalf("AdminPasswordHash after update err = %v", err)
	}
	if gotHash != newHash {
		t.Fatalf("AdminPasswordHash after update = %q, want %q", gotHash, newHash)
	}
	if !crypto.VerifyPassword(gotHash, updatedPass) {
		t.Fatalf("updated hash does not verify updated password")
	}

	// SetRequireLogin(false) updates system_state to "false" and RequireLogin returns false, nil.
	if err := store.SetRequireLogin(ctx, false); err != nil {
		t.Fatalf("SetRequireLogin(false) err = %v", err)
	}
	var rawRequireLogin string
	if err := h.db.Reader().QueryRowContext(ctx, `SELECT value FROM system_state WHERE key = 'require_login'`).Scan(&rawRequireLogin); err != nil {
		t.Fatalf("query require_login after SetRequireLogin(false): %v", err)
	}
	if rawRequireLogin != "false" {
		t.Fatalf("system_state require_login = %q, want %q", rawRequireLogin, "false")
	}
	reqLogin, err = store.RequireLogin(ctx)
	if err != nil || reqLogin {
		t.Fatalf("RequireLogin after SetRequireLogin(false) = (%v, %v), want (false, nil)", reqLogin, err)
	}

	// SetRequireLogin(true) restores true, nil.
	if err := store.SetRequireLogin(ctx, true); err != nil {
		t.Fatalf("SetRequireLogin(true) err = %v", err)
	}
	reqLogin, err = store.RequireLogin(ctx)
	if err != nil || !reqLogin {
		t.Fatalf("RequireLogin after SetRequireLogin(true) = (%v, %v), want (true, nil)", reqLogin, err)
	}
}

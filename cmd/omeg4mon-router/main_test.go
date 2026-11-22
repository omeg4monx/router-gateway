package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/modules/system"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
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
		Host:       "127.0.0.1",
		Port:       0,
		DataDir:    t.TempDir(),
		LogLevel:   "error",
		APITimeout: 5 * time.Second,
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

	client := &http.Client{Timeout: 2 * time.Second}

	healthResp, err := client.Get("http://" + addr + "/api/health")
	if err != nil {
		t.Fatalf("GET /api/health failed: %v", err)
	}
	_, _ = io.ReadAll(healthResp.Body)
	_ = healthResp.Body.Close()
	if healthResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/health status = %d, want 200", healthResp.StatusCode)
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

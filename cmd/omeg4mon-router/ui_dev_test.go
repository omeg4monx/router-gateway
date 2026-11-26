//go:build !embedui

package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/modules/system"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
)

func TestDevUIHandlerDirectsToViteDevServer(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		Host:       "127.0.0.1",
		Port:       0,
		DataDir:    t.TempDir(),
		LogLevel:   "error",
		APITimeout: 5 * time.Second,
	}

	app, err := buildApp(ctx, cfg, system.BuildInfo{Version: "dev"})
	if err != nil {
		t.Fatalf("buildApp() error = %v", err)
	}
	t.Cleanup(func() {
		_ = app.Close(context.Background())
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/dashboard/", nil)
	app.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /dashboard/ status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "http://127.0.0.1:20127") {
		t.Fatalf("expected dev UI notice to reference http://127.0.0.1:20127, got %q", rec.Body.String())
	}
}

package httpserver_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"testing/fstest"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

type stubAuthVerifier struct {
	requireLogin bool
	validToken   string
}

func (s *stubAuthVerifier) AuthenticateAdmin(r *http.Request) (string, error) {
	if r.Header.Get("Authorization") == "Bearer "+s.validToken {
		return "admin", nil
	}
	return "", errors.New("invalid credentials")
}

func (s *stubAuthVerifier) RequireLogin(ctx context.Context) bool {
	return s.requireLogin
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestServerIsolatesJSON404FromSPAFallback(t *testing.T) {
	uiFS := fstest.MapFS{
		"index.html":             {Data: []byte("<!DOCTYPE html><html><body>root-spa</body></html>")},
		"dashboard/index.html":   {Data: []byte("<!DOCTYPE html><html><body>dashboard-page</body></html>")},
		"about.html":             {Data: []byte("<!DOCTYPE html><html><body>about-page</body></html>")},
		"assets/index-abc123.js": {Data: []byte("console.log('vite-ok');")},
		"_next/static/app.js":    {Data: []byte("console.log('ok');")},
	}

	srv := httpserver.New(httpserver.Options{
		Config: config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: time.Second},
		Logger: testLogger(),
		UI:     uiFS,
	})
	handler := srv.Handler()

	t.Run("unknown API and gateway paths return JSON 404 not HTML", func(t *testing.T) {
		paths := []string{
			"/api/unknown",
			"/api",
			"/v1/unknown",
			"/v1beta/unknown",
			"/codex/unknown",
			"/responses",
			"/models/unknown",
		}
		for _, p := range paths {
			req := httptest.NewRequest(http.MethodGet, p, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("GET %s: status = %d, want 404 (body=%q)", p, rec.Code, rec.Body.String())
			}
			ct := rec.Header().Get("Content-Type")
			if !strings.Contains(ct, "application/json") {
				t.Fatalf("GET %s: Content-Type = %q, want application/json", p, ct)
			}

			var env httpserver.ErrorEnvelope
			if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
				t.Fatalf("GET %s: invalid JSON error envelope %q: %v", p, rec.Body.String(), err)
			}
			if env.Error.Code == "" || env.Error.Message == "" {
				t.Fatalf("GET %s: expected non-empty error code and message, got %+v", p, env.Error)
			}
			if env.Error.RequestID == "" {
				t.Fatalf("GET %s: expected non-empty request_id, got %+v", p, env.Error)
			}
		}
	})

	t.Run("static UI resolves directory index, html file, SPA fallback, and immutable assets", func(t *testing.T) {
		cases := []struct {
			name            string
			method          string
			path            string
			wantStatus      int
			wantContentType string
			wantBodySubstr  string
			wantCacheSubstr string
		}{
			{
				name:            "dashboard directory index without trailing slash",
				method:          http.MethodGet,
				path:            "/dashboard",
				wantStatus:      http.StatusOK,
				wantContentType: "text/html",
				wantBodySubstr:  "dashboard-page",
				wantCacheSubstr: "no-store",
			},
			{
				name:            "dashboard directory index with trailing slash",
				method:          http.MethodGet,
				path:            "/dashboard/",
				wantStatus:      http.StatusOK,
				wantContentType: "text/html",
				wantBodySubstr:  "dashboard-page",
				wantCacheSubstr: "no-store",
			},
			{
				name:            "html suffix resolution",
				method:          http.MethodGet,
				path:            "/about",
				wantStatus:      http.StatusOK,
				wantContentType: "text/html",
				wantBodySubstr:  "about-page",
				wantCacheSubstr: "no-store",
			},
			{
				name:            "client-side SPA route falls back to root index.html",
				method:          http.MethodGet,
				path:            "/settings/providers",
				wantStatus:      http.StatusOK,
				wantContentType: "text/html",
				wantBodySubstr:  "root-spa",
				wantCacheSubstr: "no-store",
			},
			{
				name:            "immutable vite static asset",
				method:          http.MethodGet,
				path:            "/assets/index-abc123.js",
				wantStatus:      http.StatusOK,
				wantBodySubstr:  "console.log('vite-ok');",
				wantCacheSubstr: "immutable",
			},
			{
				name:            "immutable next static asset",
				method:          http.MethodGet,
				path:            "/_next/static/app.js",
				wantStatus:      http.StatusOK,
				wantBodySubstr:  "console.log('ok');",
				wantCacheSubstr: "immutable",
			},
			{
				name:            "HEAD on html route returns 200 with empty body",
				method:          http.MethodHead,
				path:            "/dashboard",
				wantStatus:      http.StatusOK,
				wantContentType: "text/html",
				wantCacheSubstr: "no-store",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(tc.method, tc.path, nil)
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)

				if rec.Code != tc.wantStatus {
					t.Fatalf("%s %s: status = %d, want %d (body=%q)", tc.method, tc.path, rec.Code, tc.wantStatus, rec.Body.String())
				}
				if tc.wantContentType != "" && !strings.Contains(rec.Header().Get("Content-Type"), tc.wantContentType) {
					t.Fatalf("%s %s: Content-Type = %q, want substring %q", tc.method, tc.path, rec.Header().Get("Content-Type"), tc.wantContentType)
				}
				if tc.wantCacheSubstr != "" && !strings.Contains(rec.Header().Get("Cache-Control"), tc.wantCacheSubstr) {
					t.Fatalf("%s %s: Cache-Control = %q, want substring %q", tc.method, tc.path, rec.Header().Get("Cache-Control"), tc.wantCacheSubstr)
				}
				if tc.method == http.MethodHead {
					if rec.Body.Len() != 0 {
						t.Fatalf("HEAD %s: expected empty body, got %q", tc.path, rec.Body.String())
					}
				} else if tc.wantBodySubstr != "" && !strings.Contains(rec.Body.String(), tc.wantBodySubstr) {
					t.Fatalf("%s %s: body = %q, want substring %q", tc.method, tc.path, rec.Body.String(), tc.wantBodySubstr)
				}
			})
		}
	})

	t.Run("missing assets return 404 instead of SPA HTML", func(t *testing.T) {
		for _, p := range []string{"/assets/missing.js", "/assets/chunk", "/_next/static/missing.js", "/_next/data/build/page", "/assets/logo.png"} {
			req := httptest.NewRequest(http.MethodGet, p, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("GET %s: status = %d, want 404 (body=%q)", p, rec.Code, rec.Body.String())
			}
			if strings.Contains(rec.Body.String(), "root-spa") {
				t.Fatalf("GET %s: missing asset fell back to root index.html", p)
			}
		}
	})

	t.Run("static UI rejects non-GET/HEAD methods and path traversal", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/dashboard", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("POST /dashboard: status = %d, want 405", rec.Code)
		}
		if !strings.Contains(rec.Header().Get("Allow"), "GET") {
			t.Fatalf("POST /dashboard: Allow = %q, want GET, HEAD", rec.Header().Get("Allow"))
		}

		staticHandler := httpserver.NewStaticUIHandler(uiFS)
		for _, badPath := range []string{"/../../etc/passwd", "/dashboard/../../index.html", `/\_next/static/app.js`} {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.URL.Path = badPath
			rec := httptest.NewRecorder()
			staticHandler.ServeHTTP(rec, req)
			if rec.Code == http.StatusOK {
				t.Fatalf("GET %q: expected traversal rejection, got 200 (body=%q)", badPath, rec.Body.String())
			}
		}
	})
}

func TestGatewaySSEStreamsPastAPITimeout(t *testing.T) {
	srv := httpserver.New(httpserver.Options{
		Config: config.Config{
			Host:       "127.0.0.1",
			Port:       20128,
			APITimeout: 60 * time.Millisecond,
		},
		Logger: testLogger(),
	})

	mounts := srv.Mounts()

	mounts.AdminAPI.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(150 * time.Millisecond):
			httpserver.WriteJSON(w, http.StatusOK, map[string]string{"status": "done"})
		case <-r.Context().Done():
			httpserver.WriteError(w, r, http.StatusGatewayTimeout, "timeout", "request timed out")
		}
	})

	mounts.PublicAPI.Get("/slow-silent", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(150 * time.Millisecond):
			httpserver.WriteJSON(w, http.StatusOK, map[string]string{"status": "done"})
		case <-r.Context().Done():
			// Return without writing so middleware emits 504 JSON.
		}
	})

	mounts.Gateway.Post("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "no_flusher", "streaming unsupported")
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte("data: chunk-1\n\n"))
		flusher.Flush()

		select {
		case <-time.After(90 * time.Millisecond):
		case <-r.Context().Done():
			return
		}

		_, _ = w.Write([]byte("data: chunk-2\n\n"))
		flusher.Flush()
	})

	// Also verify compatibility aliases route cleanly on Gateway without double-prefixing.
	mounts.Gateway.Post("/v1beta/models/{model}", func(w http.ResponseWriter, r *http.Request) {
		httpserver.WriteJSON(w, http.StatusOK, map[string]string{"model": chi.URLParam(r, "model")})
	})
	mounts.Gateway.Post("/codex/responses", func(w http.ResponseWriter, r *http.Request) {
		httpserver.WriteJSON(w, http.StatusOK, map[string]string{"route": "codex"})
	})
	mounts.Gateway.Post("/responses", func(w http.ResponseWriter, r *http.Request) {
		httpserver.WriteJSON(w, http.StatusOK, map[string]string{"route": "responses"})
	})
	mounts.Gateway.Get("/models/{id}", func(w http.ResponseWriter, r *http.Request) {
		httpserver.WriteJSON(w, http.StatusOK, map[string]string{"id": chi.URLParam(r, "id")})
	})

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	t.Run("AdminAPI slow handler times out at APITimeout", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/slow")
		if err != nil {
			t.Fatalf("GET /api/slow: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusGatewayTimeout {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("GET /api/slow: status = %d, want 504 (body=%q)", resp.StatusCode, string(body))
		}
	})

	t.Run("PublicAPI silent timeout handler gets 504 JSON from middleware", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/slow-silent")
		if err != nil {
			t.Fatalf("GET /api/slow-silent: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusGatewayTimeout {
			t.Fatalf("GET /api/slow-silent: status = %d, want 504 (body=%q)", resp.StatusCode, string(body))
		}
		if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
			t.Fatalf("GET /api/slow-silent: Content-Type = %q, want application/json", resp.Header.Get("Content-Type"))
		}
		var env httpserver.ErrorEnvelope
		if err := json.Unmarshal(body, &env); err != nil {
			t.Fatalf("GET /api/slow-silent: invalid JSON %q: %v", string(body), err)
		}
		if env.Error.Code != "timeout" {
			t.Fatalf("GET /api/slow-silent: code = %q, want timeout", env.Error.Code)
		}
	})

	t.Run("Gateway SSE streams past APITimeout cleanly", func(t *testing.T) {
		resp, err := http.Post(ts.URL+"/v1/chat/completions", "application/json", strings.NewReader(`{"stream":true}`))
		if err != nil {
			t.Fatalf("POST /v1/chat/completions: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read SSE body: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("POST /v1/chat/completions: status = %d, want 200 (body=%q)", resp.StatusCode, string(body))
		}
		got := string(body)
		if !strings.Contains(got, "data: chunk-1") || !strings.Contains(got, "data: chunk-2") {
			t.Fatalf("SSE body = %q, want both chunk-1 and chunk-2", got)
		}
	})

	t.Run("Gateway compatibility aliases route cleanly", func(t *testing.T) {
		checks := []struct {
			method string
			path   string
			substr string
		}{
			{http.MethodPost, "/v1beta/models/gemini-2.5", `"model":"gemini-2.5"`},
			{http.MethodPost, "/codex/responses", `"route":"codex"`},
			{http.MethodPost, "/responses", `"route":"responses"`},
			{http.MethodGet, "/models/gpt-4o", `"id":"gpt-4o"`},
		}
		for _, c := range checks {
			req, err := http.NewRequest(c.method, ts.URL+c.path, nil)
			if err != nil {
				t.Fatalf("new request %s %s: %v", c.method, c.path, err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("%s %s: %v", c.method, c.path, err)
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), c.substr) {
				t.Fatalf("%s %s: status=%d body=%q, want 200 containing %q", c.method, c.path, resp.StatusCode, string(body), c.substr)
			}
		}
	})
}

func TestAdminAPIAuthEnforcement(t *testing.T) {
	t.Run("requireLogin=true blocks unauthenticated AdminAPI while allowing PublicAPI", func(t *testing.T) {
		srv := httpserver.New(httpserver.Options{
			Config: config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: time.Second},
			Logger: testLogger(),
			Auth:   &stubAuthVerifier{requireLogin: true, validToken: "secret"},
		})
		mounts := srv.Mounts()
		mounts.PublicAPI.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			httpserver.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})
		mounts.AdminAPI.Post("/system/shutdown", func(w http.ResponseWriter, r *http.Request) {
			httpserver.WriteJSON(w, http.StatusAccepted, map[string]string{
				"status":  "shutting_down",
				"subject": httpserver.AdminSubject(r.Context()),
			})
		})
		mounts.AdminAPI.Route("/providers", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				httpserver.WriteJSON(w, http.StatusOK, map[string]string{"ok": "true"})
			})
		})

		h := srv.Handler()

		// PublicAPI -> 200 without auth
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/health", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/health: status = %d, want 200", rec.Code)
		}

		// AdminAPI without auth -> 401 JSON
		for _, tc := range []struct {
			method string
			path   string
		}{
			{http.MethodPost, "/api/system/shutdown"},
			{http.MethodGet, "/api/providers/"},
		} {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("%s %s (no auth): status = %d, want 401", tc.method, tc.path, rec.Code)
			}
			if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
				t.Fatalf("%s %s (no auth): Content-Type = %q, want application/json", tc.method, tc.path, rec.Header().Get("Content-Type"))
			}
			var env httpserver.ErrorEnvelope
			if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
				t.Fatalf("%s %s: invalid JSON %q: %v", tc.method, tc.path, rec.Body.String(), err)
			}
			if env.Error.Code != "unauthorized" {
				t.Fatalf("%s %s: code = %q, want unauthorized", tc.method, tc.path, env.Error.Code)
			}
		}

		// AdminAPI with valid auth -> 202 and subject populated in context
		req := httptest.NewRequest(http.MethodPost, "/api/system/shutdown", nil)
		req.Header.Set("Authorization", "Bearer secret")
		rec = httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusAccepted {
			t.Fatalf("POST /api/system/shutdown (with auth): status = %d, want 202", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), `"subject":"admin"`) {
			t.Fatalf("POST /api/system/shutdown: expected subject=admin in body %q", rec.Body.String())
		}

		// Wrong method on API route -> 405 JSON
		rec = httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/api/health", nil))
		if rec.Code != http.StatusMethodNotAllowed || !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
			t.Fatalf("DELETE /api/health: status=%d ct=%q, want 405 JSON", rec.Code, rec.Header().Get("Content-Type"))
		}
	})

	t.Run("non-loopback host enforces AdminAPI auth even when RequireLogin is false", func(t *testing.T) {
		srv := httpserver.New(httpserver.Options{
			Config: config.Config{
				Host:                "0.0.0.0",
				Port:                20128,
				AllowInsecureRemote: false,
				APITimeout:          time.Second,
			},
			Logger: testLogger(),
			Auth:   &stubAuthVerifier{requireLogin: false, validToken: "secret"},
		})
		srv.Mounts().AdminAPI.Get("/secret", func(w http.ResponseWriter, r *http.Request) {
			httpserver.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		})

		h := srv.Handler()

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/secret", nil))
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("non-loopback unauthenticated AdminAPI: status = %d, want 401", rec.Code)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/secret", nil)
		req.Header.Set("Authorization", "Bearer secret")
		rec = httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("non-loopback authenticated AdminAPI: status = %d, want 200", rec.Code)
		}
	})

	t.Run("non-loopback with AllowInsecureRemote=true and RequireLogin=false permits unauthenticated AdminAPI", func(t *testing.T) {
		srv := httpserver.New(httpserver.Options{
			Config: config.Config{
				Host:                "0.0.0.0",
				Port:                20128,
				AllowInsecureRemote: true,
				APITimeout:          time.Second,
			},
			Logger: testLogger(),
			Auth:   &stubAuthVerifier{requireLogin: false, validToken: "secret"},
		})
		srv.Mounts().AdminAPI.Get("/secret", func(w http.ResponseWriter, r *http.Request) {
			httpserver.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		})

		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/secret", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})

	t.Run("non-loopback with nil Auth and AllowInsecureRemote=false fails closed with 401 JSON", func(t *testing.T) {
		srv := httpserver.New(httpserver.Options{
			Config: config.Config{
				Host:                "0.0.0.0",
				Port:                20128,
				AllowInsecureRemote: false,
				APITimeout:          time.Second,
			},
			Logger: testLogger(),
			Auth:   nil,
		})
		srv.Mounts().AdminAPI.Get("/secret", func(w http.ResponseWriter, r *http.Request) {
			httpserver.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		})

		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/secret", nil))
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("non-loopback nil Auth AdminAPI: status = %d, want 401", rec.Code)
		}
		if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
			t.Fatalf("non-loopback nil Auth AdminAPI: Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
		}
		var env httpserver.ErrorEnvelope
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("non-loopback nil Auth AdminAPI: invalid JSON %q: %v", rec.Body.String(), err)
		}
		if env.Error.Code != "unauthorized" {
			t.Fatalf("non-loopback nil Auth AdminAPI: code = %q, want unauthorized", env.Error.Code)
		}
	})

	t.Run("loopback with nil Auth permits AdminAPI", func(t *testing.T) {
		srv := httpserver.New(httpserver.Options{
			Config: config.Config{
				Host:                "127.0.0.1",
				Port:                20128,
				AllowInsecureRemote: false,
				APITimeout:          time.Second,
			},
			Logger: testLogger(),
			Auth:   nil,
		})
		srv.Mounts().AdminAPI.Get("/secret", func(w http.ResponseWriter, r *http.Request) {
			httpserver.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		})

		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/secret", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("loopback nil Auth AdminAPI: status = %d, want 200", rec.Code)
		}
	})

	t.Run("non-loopback with nil Auth and AllowInsecureRemote=true permits AdminAPI", func(t *testing.T) {
		srv := httpserver.New(httpserver.Options{
			Config: config.Config{
				Host:                "0.0.0.0",
				Port:                20128,
				AllowInsecureRemote: true,
				APITimeout:          time.Second,
			},
			Logger: testLogger(),
			Auth:   nil,
		})
		srv.Mounts().AdminAPI.Get("/secret", func(w http.ResponseWriter, r *http.Request) {
			httpserver.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		})

		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/secret", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("non-loopback AllowInsecureRemote nil Auth AdminAPI: status = %d, want 200", rec.Code)
		}
	})
}

func TestPanicRecoveryReturnsJSON500(t *testing.T) {
	srv := httpserver.New(httpserver.Options{
		Config: config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: time.Second},
		Logger: testLogger(),
	})
	srv.Mounts().PublicAPI.Get("/boom", func(w http.ResponseWriter, r *http.Request) {
		panic("api exploded")
	})
	srv.Mounts().Gateway.Post("/v1/boom", func(w http.ResponseWriter, r *http.Request) {
		panic("gateway exploded")
	})

	h := srv.Handler()
	for _, tc := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/boom"},
		{http.MethodPost, "/v1/boom"},
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("%s %s: status = %d, want 500", tc.method, tc.path, rec.Code)
		}
		if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
			t.Fatalf("%s %s: Content-Type = %q, want application/json", tc.method, tc.path, rec.Header().Get("Content-Type"))
		}
		var env httpserver.ErrorEnvelope
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("%s %s: invalid JSON %q: %v", tc.method, tc.path, rec.Body.String(), err)
		}
		if env.Error.Code != "internal_error" || env.Error.RequestID == "" {
			t.Fatalf("%s %s: unexpected envelope %+v", tc.method, tc.path, env.Error)
		}
	}
}

func TestDevUIHandlerFallback(t *testing.T) {
	devHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("dev-ui:" + r.URL.Path))
	})
	srv := httpserver.New(httpserver.Options{
		Config:       config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: time.Second},
		Logger:       testLogger(),
		DevUIHandler: devHandler,
	})
	h := srv.Handler()

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/dashboard", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "dev-ui:/dashboard" {
		t.Fatalf("GET /dashboard: status=%d body=%q", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/missing", nil))
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("GET /api/missing: status=%d ct=%q", rec.Code, rec.Header().Get("Content-Type"))
	}
}

func TestServerStartAndGracefulShutdown(t *testing.T) {
	srv := httpserver.New(httpserver.Options{
		Config: config.Config{
			Host:       "127.0.0.1",
			Port:       0,
			APITimeout: 5 * time.Second,
		},
		Logger: testLogger(),
	})

	entered := make(chan struct{})
	release := make(chan struct{})
	var completed atomic.Bool

	srv.Mounts().PublicAPI.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpserver.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	srv.Mounts().PublicAPI.Get("/inflight", func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-release
		completed.Store(true)
		httpserver.WriteJSON(w, http.StatusOK, map[string]string{"drained": "true"})
	})

	ctx := context.Background()
	startErrCh := make(chan error, 1)
	go func() {
		startErrCh <- srv.Start(ctx)
	}()

	// Wait for listener to bind an ephemeral port.
	var addr string
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		a := srv.Addr()
		if a != "" && !strings.HasSuffix(a, ":0") {
			addr = a
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if addr == "" {
		t.Fatalf("server did not bind ephemeral port in time")
	}

	client := &http.Client{
		Transport: &http.Transport{},
	}
	t.Cleanup(client.CloseIdleConnections)

	// Verify /api/health serves live HTTP.
	resp, err := client.Get("http://" + addr + "/api/health")
	if err != nil {
		t.Fatalf("GET /api/health: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/health: status = %d, want 200", resp.StatusCode)
	}

	// Start an in-flight request, then trigger graceful shutdown.
	inflightRespCh := make(chan *http.Response, 1)
	inflightErrCh := make(chan error, 1)
	go func() {
		r, err := client.Get("http://" + addr + "/api/inflight")
		if err != nil {
			inflightErrCh <- err
			return
		}
		inflightRespCh <- r
	}()

	<-entered

	shutdownErrCh := make(chan error, 1)
	go func() {
		shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		shutdownErrCh <- srv.Shutdown(shutCtx)
	}()

	// Let Shutdown begin draining, then release the in-flight request.
	time.Sleep(30 * time.Millisecond)
	close(release)

	select {
	case err := <-inflightErrCh:
		t.Fatalf("in-flight request failed during graceful shutdown: %v", err)
	case r := <-inflightRespCh:
		_, _ = io.Copy(io.Discard, r.Body)
		r.Body.Close()
		if r.StatusCode != http.StatusOK {
			t.Fatalf("in-flight status = %d, want 200", r.StatusCode)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for in-flight request")
	}

	if !completed.Load() {
		t.Fatal("expected in-flight handler to complete before shutdown")
	}

	if err := <-shutdownErrCh; err != nil {
		t.Fatalf("Shutdown returned error: %v", err)
	}
	if err := <-startErrCh; err != nil {
		t.Fatalf("Start returned error after Shutdown: %v", err)
	}
}

func TestIsHTTPSDetection(t *testing.T) {
	cases := []struct {
		name    string
		tls     bool
		headers map[string]string
		want    bool
	}{
		{name: "plain http", want: false},
		{name: "direct tls", tls: true, want: true},
		{name: "x-forwarded-proto https", headers: map[string]string{"X-Forwarded-Proto": "https"}, want: true},
		{name: "x-forwarded-proto list", headers: map[string]string{"X-Forwarded-Proto": "HTTPS, http"}, want: true},
		{name: "x-forwarded-proto http", headers: map[string]string{"X-Forwarded-Proto": "http"}, want: false},
		{name: "rfc7239 forwarded proto https", headers: map[string]string{"Forwarded": `for=192.0.2.60;proto=https;by=203.0.113.43`}, want: true},
		{name: "rfc7239 forwarded quoted https", headers: map[string]string{"Forwarded": `proto="https", for=192.0.2.60`}, want: true},
		{name: "rfc7239 forwarded proto http", headers: map[string]string{"Forwarded": `for=192.0.2.60;proto=http`}, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
			if tc.tls {
				req.TLS = &tls.ConnectionState{}
			}
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			if got := httpserver.IsHTTPS(req); got != tc.want {
				t.Fatalf("IsHTTPS() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	srv := httpserver.New(httpserver.Options{
		Config: config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: time.Second},
		Logger: testLogger(),
	})
	srv.Mounts().PublicAPI.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		httpserver.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Plain HTTP: baseline headers without HSTS.
	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("X-Content-Type-Options = %q, want nosniff", rec.Header().Get("X-Content-Type-Options"))
	}
	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatalf("X-Frame-Options = %q, want DENY", rec.Header().Get("X-Frame-Options"))
	}
	if rec.Header().Get("Referrer-Policy") != "strict-origin-when-cross-origin" {
		t.Fatalf("Referrer-Policy = %q, want strict-origin-when-cross-origin", rec.Header().Get("Referrer-Policy"))
	}
	if hsts := rec.Header().Get("Strict-Transport-Security"); hsts != "" {
		t.Fatalf("expected no HSTS on plain HTTP, got %q", hsts)
	}

	// HTTPS via X-Forwarded-Proto: includes HSTS.
	httpsReq := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	httpsReq.Header.Set("X-Forwarded-Proto", "https")
	httpsRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(httpsRec, httpsReq)
	if hsts := httpsRec.Header().Get("Strict-Transport-Security"); !strings.Contains(hsts, "max-age=") {
		t.Fatalf("expected HSTS header on HTTPS request, got %q", hsts)
	}
}

func TestAdminAuthEnforcedOverHTTPSOnLoopback(t *testing.T) {
	srv := httpserver.New(httpserver.Options{
		Config: config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: time.Second},
		Logger: testLogger(),
		Auth:   &stubAuthVerifier{requireLogin: false, validToken: "secret"},
	})
	srv.Mounts().AdminAPI.Get("/secret", func(w http.ResponseWriter, r *http.Request) {
		httpserver.WriteJSON(w, http.StatusOK, map[string]string{"ok": "true"})
	})

	// Local plain HTTP with requireLogin=false is allowed.
	httpReq := httptest.NewRequest(http.MethodGet, "/api/secret", nil)
	httpRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(httpRec, httpReq)
	if httpRec.Code != http.StatusOK {
		t.Fatalf("plain HTTP loopback status = %d, want 200", httpRec.Code)
	}

	// Reverse-proxied HTTPS to loopback must enforce login.
	httpsReq := httptest.NewRequest(http.MethodGet, "/api/secret", nil)
	httpsReq.Header.Set("X-Forwarded-Proto", "https")
	httpsRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(httpsRec, httpsReq)
	if httpsRec.Code != http.StatusUnauthorized {
		t.Fatalf("proxied HTTPS without token status = %d, want 401", httpsRec.Code)
	}

	// Valid token over proxied HTTPS succeeds.
	authedReq := httptest.NewRequest(http.MethodGet, "/api/secret", nil)
	authedReq.Header.Set("X-Forwarded-Proto", "https")
	authedReq.Header.Set("Authorization", "Bearer secret")
	authedRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(authedRec, authedReq)
	if authedRec.Code != http.StatusOK {
		t.Fatalf("proxied HTTPS with valid token status = %d, want 200", authedRec.Code)
	}
}

func generateSelfSignedCertFiles(t *testing.T) (certPath, keyPath string, certPool *x509.CertPool) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "127.0.0.1"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	dir := t.TempDir()
	certPath = filepath.Join(dir, "server.crt")
	keyPath = filepath.Join(dir, "server.key")
	if err := os.WriteFile(certPath, certPEM, 0o600); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, keyPEM, 0o600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(certPEM) {
		t.Fatal("failed to add cert to pool")
	}
	return certPath, keyPath, pool
}

func TestServerStartNativeTLS(t *testing.T) {
	certPath, keyPath, pool := generateSelfSignedCertFiles(t)

	srv := httpserver.New(httpserver.Options{
		Config: config.Config{
			Host:        "127.0.0.1",
			Port:        0,
			APITimeout:  time.Second,
			TLSCertFile: certPath,
			TLSKeyFile:  keyPath,
		},
		Logger: testLogger(),
	})
	srv.Mounts().PublicAPI.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpserver.WriteJSON(w, http.StatusOK, map[string]bool{"tls": r.TLS != nil})
	})

	startErrCh := make(chan error, 1)
	go func() {
		startErrCh <- srv.Start(context.Background())
	}()

	var addr string
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		a := srv.Addr()
		if a != "" && !strings.HasSuffix(a, ":0") {
			addr = a
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if addr == "" {
		t.Fatal("TLS server did not bind ephemeral port in time")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    pool,
				MinVersion: tls.VersionTLS12,
			},
		},
	}
	t.Cleanup(client.CloseIdleConnections)

	resp, err := client.Get("https://" + addr + "/api/health")
	if err != nil {
		t.Fatalf("HTTPS GET /api/health: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HTTPS GET /api/health status = %d, want 200", resp.StatusCode)
	}
	if hsts := resp.Header.Get("Strict-Transport-Security"); !strings.Contains(hsts, "max-age=") {
		t.Fatalf("expected HSTS header on native TLS response, got %q", hsts)
	}

	shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		t.Fatalf("Shutdown TLS server: %v", err)
	}
	if err := <-startErrCh; err != nil {
		t.Fatalf("Start TLS server returned error: %v", err)
	}
}

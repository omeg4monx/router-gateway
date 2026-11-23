package pxpipe

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

func setupTestPxpipe(t *testing.T) (*Module, http.Handler) {
	t.Helper()
	database, err := db.Open(filepath.Join(t.TempDir(), "pxpipe.sqlite"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	mod := New(Options{
		DB:  database,
		Bus: eventbus.New(),
	})

	ctx := context.Background()
	if err := database.MigrateModule(ctx, mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := mod.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() { _ = mod.Stop(context.Background()) })

	r := chi.NewRouter()
	r.Route("/api", func(api chi.Router) {
		mod.RegisterRoutes(lifecycle.RouteMounts{
			PublicAPI: api,
			AdminAPI:  api,
			Gateway:   r,
		})
	})

	return mod, r
}

func doPxpipeJSON(t *testing.T, handler http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		buf = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestPxpipeLifecycleStatusAndHealth(t *testing.T) {
	mod, handler := setupTestPxpipe(t)
	if mod.Name() != "pxpipe" {
		t.Fatalf("expected name pxpipe, got %q", mod.Name())
	}
	if mod.Priority() != 35 {
		t.Fatalf("expected priority 35, got %d", mod.Priority())
	}

	// Running by default (auto_start=true)
	rec := doPxpipeJSON(t, handler, http.MethodGet, "/api/pxpipe/status", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/pxpipe/status: %d %s", rec.Code, rec.Body.String())
	}
	var st Status
	if err := json.Unmarshal(rec.Body.Bytes(), &st); err != nil {
		t.Fatalf("decode status: %v", err)
	}
	if !st.Running || st.Mode != "in-process" {
		t.Fatalf("expected running in-process pxpipe, got %+v", st)
	}

	// Health check passes when running
	hRec := doPxpipeJSON(t, handler, http.MethodGet, "/api/pxpipe/health", nil)
	if hRec.Code != http.StatusOK {
		t.Fatalf("GET /api/pxpipe/health: %d %s", hRec.Code, hRec.Body.String())
	}
	var health HealthReport
	if err := json.Unmarshal(hRec.Body.Bytes(), &health); err != nil {
		t.Fatalf("decode health: %v", err)
	}
	if !health.Healthy || len(health.Checks) != 3 {
		t.Fatalf("expected healthy 3-check report, got %+v", health)
	}

	// Stop -> health reports unhealthy (pipeline_active=false)
	stopRec := doPxpipeJSON(t, handler, http.MethodPost, "/api/pxpipe/stop", nil)
	if stopRec.Code != http.StatusOK {
		t.Fatalf("POST /api/pxpipe/stop: %d", stopRec.Code)
	}
	hRec2 := doPxpipeJSON(t, handler, http.MethodGet, "/api/pxpipe/health", nil)
	var health2 HealthReport
	_ = json.Unmarshal(hRec2.Body.Bytes(), &health2)
	if health2.Healthy {
		t.Fatalf("expected unhealthy when stopped, got %+v", health2)
	}

	// Restart brings it back online
	rstRec := doPxpipeJSON(t, handler, http.MethodPost, "/api/pxpipe/restart", nil)
	if rstRec.Code != http.StatusOK {
		t.Fatalf("POST /api/pxpipe/restart: %d", rstRec.Code)
	}
	var restarted Status
	_ = json.Unmarshal(rstRec.Body.Bytes(), &restarted)
	if !restarted.Running {
		t.Fatalf("expected running after restart")
	}
}

func TestPxpipeInterceptorTransformsAndRuntimeGating(t *testing.T) {
	mod, handler := setupTestPxpipe(t)

	// Enable max_history_messages=2
	rec := doPxpipeJSON(t, handler, http.MethodPatch, "/api/pxpipe/config", map[string]any{
		"max_history_messages": 2,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH config: %d %s", rec.Code, rec.Body.String())
	}

	rawBody := []byte(`{
		"model": "gpt-4o",
		"user": null,
		"metadata": null,
		"messages": [
			{"role": "system", "content": "You   are   helpful."},
			{"role": "system", "content": "You are helpful."},
			{"role": "user", "content": "first   message"},
			{"role": "assistant", "content": "second   reply"},
			{"role": "user", "content": "latest    question"}
		]
	}`)

	creq := &contracts.CanonicalRequest{
		RequestID: "req-1",
		Protocol:  "openai",
		Model:     "gpt-4o",
		Body:      append([]byte(nil), rawBody...),
	}
	if err := mod.OnRequest(context.Background(), creq); err != nil {
		t.Fatalf("OnRequest: %v", err)
	}
	if len(creq.Body) >= len(rawBody) {
		t.Fatalf("expected transformed body to be smaller: before=%d after=%d", len(rawBody), len(creq.Body))
	}

	var parsed map[string]any
	if err := json.Unmarshal(creq.Body, &parsed); err != nil {
		t.Fatalf("transformed body invalid JSON: %v", err)
	}
	if _, hasNull := parsed["user"]; hasNull {
		t.Fatalf("expected null field 'user' to be stripped")
	}
	msgs, _ := parsed["messages"].([]any)
	// 1 deduplicated system message + 2 windowed non-system messages = 3
	if len(msgs) != 3 {
		t.Fatalf("expected 3 messages (1 deduped system + 2 windowed), got %d: %s", len(msgs), string(creq.Body))
	}

	// Verify stats & logs recorded
	statsRec := doPxpipeJSON(t, handler, http.MethodGet, "/api/pxpipe/stats", nil)
	var stats Stats
	_ = json.Unmarshal(statsRec.Body.Bytes(), &stats)
	if stats.RequestsTransformed != 1 || stats.BytesSaved <= 0 {
		t.Fatalf("expected 1 transformed request with positive bytes saved, got %+v", stats)
	}

	logsRec := doPxpipeJSON(t, handler, http.MethodGet, "/api/pxpipe/logs", nil)
	var logsResp LogListResponse
	_ = json.Unmarshal(logsRec.Body.Bytes(), &logsResp)
	if len(logsResp.Logs) != 1 || logsResp.Logs[0].RequestID != "req-1" {
		t.Fatalf("expected 1 transform log for req-1, got %+v", logsResp.Logs)
	}

	// Stop pipeline -> OnRequest passes through untouched
	_ = doPxpipeJSON(t, handler, http.MethodPost, "/api/pxpipe/stop", nil)
	creq2 := &contracts.CanonicalRequest{
		RequestID: "req-2",
		Model:     "gpt-4o",
		Body:      append([]byte(nil), rawBody...),
	}
	if err := mod.OnRequest(context.Background(), creq2); err != nil {
		t.Fatalf("OnRequest while stopped: %v", err)
	}
	if !bytes.Equal(creq2.Body, rawBody) {
		t.Fatalf("expected body untouched when pxpipe stopped")
	}

	// Reset stats & clear logs
	_ = doPxpipeJSON(t, handler, http.MethodPost, "/api/pxpipe/stats/reset", nil)
	_ = doPxpipeJSON(t, handler, http.MethodDelete, "/api/pxpipe/logs", nil)

	// Preview endpoint works even when pipeline is stopped
	prevRec := doPxpipeJSON(t, handler, http.MethodPost, "/api/pxpipe/preview", map[string]any{
		"payload": map[string]any{
			"model": "gpt-4o",
			"empty": nil,
			"messages": []map[string]any{
				{"role": "user", "content": "hello     world"},
			},
		},
	})
	if prevRec.Code != http.StatusOK {
		t.Fatalf("POST /api/pxpipe/preview: %d %s", prevRec.Code, prevRec.Body.String())
	}
	var prev PreviewResponse
	if err := json.Unmarshal(prevRec.Body.Bytes(), &prev); err != nil {
		t.Fatalf("decode preview: %v", err)
	}
	if prev.BytesSaved <= 0 || len(prev.StagesApplied) == 0 {
		t.Fatalf("expected preview savings and stages, got %+v", prev)
	}
}

func TestPxpipePreservesCodeFencesAvoidsOrphanToolsAndEscaping(t *testing.T) {
	mod, handler := setupTestPxpipe(t)
	_ = doPxpipeJSON(t, handler, http.MethodPatch, "/api/pxpipe/config", map[string]any{
		"max_history_messages": 2,
	})

	body := []byte(`{"model":"gpt-4o","messages":[{"role":"user","content":"q1"},{"role":"assistant","content":"call"},{"role":"tool","tool_call_id":"tc1","content":"result"},{"role":"user","content":"keep <tag>&val\n` + "```python\\n    def f():\\n        return 1\\n```" + `"}]}`)
	creq := &contracts.CanonicalRequest{
		RequestID: "req-fence",
		Model:     "gpt-4o",
		Body:      append([]byte(nil), body...),
	}
	if err := mod.OnRequest(context.Background(), creq); err != nil {
		t.Fatalf("OnRequest: %v", err)
	}
	out := string(creq.Body)
	if !bytes.Contains(creq.Body, []byte("    def f():")) {
		t.Fatalf("expected fenced code indentation preserved, got: %s", out)
	}
	if !bytes.Contains(creq.Body, []byte("<tag>&val")) {
		t.Fatalf("expected HTML characters not escaped, got: %s", out)
	}
	var parsed struct {
		Messages []struct {
			Role string `json:"role"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(creq.Body, &parsed); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(parsed.Messages) == 0 || parsed.Messages[0].Role == "tool" {
		t.Fatalf("expected leading orphan tool message to be skipped, got: %+v", parsed.Messages)
	}
}

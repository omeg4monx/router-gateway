package tokensaver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type testEnv struct {
	mod    *Module
	router chi.Router
}

func setupTokenSaverTest(t *testing.T) *testEnv {
	t.Helper()
	database, err := db.Open(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	mod := New(Options{
		DB:     database,
		Bus:    eventbus.New(),
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	ctx := context.Background()
	if err := database.MigrateModule(ctx, mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("MigrateModule: %v", err)
	}
	if err := mod.Start(ctx); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() { _ = mod.Stop(context.Background()) })

	r := chi.NewRouter()
	r.Route("/api", func(ar chi.Router) {
		mod.RegisterRoutes(lifecycle.RouteMounts{AdminAPI: ar})
	})
	return &testEnv{mod: mod, router: r}
}

func doJSON(t *testing.T, h http.Handler, method, path string, payload any) *httptest.ResponseRecorder {
	t.Helper()
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		body = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(t *testing.T) {
	env := setupTokenSaverTest(t)

	if env.mod.Priority() != 30 {
		t.Fatalf("expected priority 30, got %d", env.mod.Priority())
	}

	// Enable aggressive strategy with small tool output limit.
	rec := doJSON(t, env.router, http.MethodPatch, "/api/tokensaver/config", map[string]any{
		"enabled":               true,
		"strategy":              "aggressive",
		"max_history_messages":  4,
		"max_tool_output_chars": 60,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("patch config: status %d body %s", rec.Code, rec.Body.String())
	}

	codeBlock := "```python\ndef f():\n    x = 1\n\n\n    return x\n```"
	prose := "Hello     world!\n\n\n\nThis   has   extra    spaces.\n" + codeBlock
	longTool := strings.Repeat("tool-output-chunk-", 20)

	payload, _ := json.Marshal(map[string]any{
		"model": "gpt-4o",
		"messages": []any{
			map[string]any{"role": "system", "content": "You   are   helpful."},
			map[string]any{"role": "user", "content": "turn 1"},
			map[string]any{"role": "assistant", "content": "reply 1"},
			map[string]any{"role": "user", "content": "turn 2"},
			map[string]any{"role": "tool", "content": longTool},
			map[string]any{"role": "user", "content": prose},
		},
	})

	req := &contracts.CanonicalRequest{
		RequestID: "req-compress",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      payload,
		Metadata:  make(map[string]any),
	}
	if err := env.mod.OnRequest(context.Background(), req); err != nil {
		t.Fatalf("OnRequest: %v", err)
	}
	if len(req.Body) >= len(payload) {
		t.Fatalf("expected compressed body smaller than %d, got %d", len(payload), len(req.Body))
	}
	if !strings.Contains(string(req.Body), `def f():\n    x = 1\n\n\n    return x`) {
		t.Fatalf("fenced code block indentation/blank lines were corrupted: %s", req.Body)
	}
	if !strings.Contains(string(req.Body), "[...truncated ") {
		t.Fatalf("expected long tool output truncation marker, got: %s", req.Body)
	}

	// Bypass flag skips compression completely.
	bypassReq := &contracts.CanonicalRequest{
		RequestID: "req-bypass",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      append([]byte(nil), payload...),
		Metadata:  map[string]any{"bypass_tokensaver": true},
	}
	if err := env.mod.OnRequest(context.Background(), bypassReq); err != nil {
		t.Fatalf("OnRequest bypass: %v", err)
	}
	if !bytes.Equal(bypassReq.Body, payload) {
		t.Fatal("expected bypass_tokensaver to leave body untouched")
	}
}

func TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(t *testing.T) {
	env := setupTokenSaverTest(t)

	rec := doJSON(t, env.router, http.MethodPatch, "/api/tokensaver/config", map[string]any{
		"enabled":              true,
		"strategy":             "whitespace",
		"headroom_enabled":     true,
		"inject_cache_control": true,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("patch config: status %d", rec.Code)
	}

	makeAnthropicReq := func(reqID string) *contracts.CanonicalRequest {
		return &contracts.CanonicalRequest{
			RequestID: reqID,
			Protocol:  "anthropic",
			Model:     "claude-sonnet-4",
			Body:      []byte(`{"model":"claude-sonnet-4","system":[{"type":"text","text":"Large stable system instructions for coding."}],"messages":[{"role":"user","content":"Question 1"}]}`),
			Metadata:  make(map[string]any),
		}
	}

	// First request -> miss, then OnResponseHeaders records provider/connection affinity.
	r1 := makeAnthropicReq("req-1")
	if err := env.mod.OnRequest(context.Background(), r1); err != nil {
		t.Fatalf("OnRequest 1: %v", err)
	}
	if r1.Metadata["headroom_cache_hit"] != false {
		t.Fatalf("expected first request cache miss, got %v", r1.Metadata["headroom_cache_hit"])
	}
	if !strings.Contains(string(r1.Body), `"cache_control":{"type":"ephemeral"}`) {
		t.Fatalf("expected Anthropic ephemeral cache_control injection, got: %s", r1.Body)
	}

	m1 := &contracts.ResponseMeta{
		RequestID:    "req-1",
		ProviderID:   "anthropic",
		ConnectionID: "conn-1",
		Headers:      make(http.Header),
		Metadata:     r1.Metadata,
	}
	_ = env.mod.OnResponseHeaders(context.Background(), m1)
	if m1.Headers.Get("X-Omeg4mon-Headroom-Cache") != "miss" {
		t.Fatalf("expected X-Omeg4mon-Headroom-Cache=miss, got %q", m1.Headers.Get("X-Omeg4mon-Headroom-Cache"))
	}

	// Second request with same prefix -> hit with affinity metadata.
	r2 := makeAnthropicReq("req-2")
	if err := env.mod.OnRequest(context.Background(), r2); err != nil {
		t.Fatalf("OnRequest 2: %v", err)
	}
	if r2.Metadata["headroom_cache_hit"] != true {
		t.Fatalf("expected second request cache hit, got %v", r2.Metadata["headroom_cache_hit"])
	}
	if r2.Metadata["headroom_affinity_connection"] != "conn-1" {
		t.Fatalf("expected affinity connection conn-1, got %v", r2.Metadata["headroom_affinity_connection"])
	}
	m2 := &contracts.ResponseMeta{
		RequestID:    "req-2",
		ProviderID:   "anthropic",
		ConnectionID: "conn-1",
		Headers:      make(http.Header),
		Metadata:     r2.Metadata,
	}
	_ = env.mod.OnResponseHeaders(context.Background(), m2)
	if m2.Headers.Get("X-Omeg4mon-Headroom-Cache") != "hit" {
		t.Fatalf("expected X-Omeg4mon-Headroom-Cache=hit, got %q", m2.Headers.Get("X-Omeg4mon-Headroom-Cache"))
	}
}

func TestHeadroomLRUBoundedCapacityAndConcurrency(t *testing.T) {
	env := setupTokenSaverTest(t)

	rec := doJSON(t, env.router, http.MethodPatch, "/api/tokensaver/config", map[string]any{
		"enabled":           true,
		"headroom_enabled":  true,
		"max_cache_entries": 16,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("patch config: status %d", rec.Code)
	}

	var wg sync.WaitGroup
	for i := range 64 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			req := &contracts.CanonicalRequest{
				RequestID: fmt.Sprintf("req-%d", idx),
				Protocol:  "openai-chat",
				Model:     "gpt-4o",
				Body:      []byte(fmt.Sprintf(`{"model":"gpt-4o","messages":[{"role":"system","content":"sys-%d"},{"role":"user","content":"u"}]}`, idx)),
				Metadata:  make(map[string]any),
			}
			_ = env.mod.OnRequest(context.Background(), req)
			_ = env.mod.OnResponseHeaders(context.Background(), &contracts.ResponseMeta{
				RequestID:    req.RequestID,
				ProviderID:   "openai",
				ConnectionID: "c1",
				Headers:      make(http.Header),
				Metadata:     req.Metadata,
			})
		}(i)
	}
	wg.Wait()

	statsRec := doJSON(t, env.router, http.MethodGet, "/api/tokensaver/stats", nil)
	if statsRec.Code != http.StatusOK {
		t.Fatalf("get stats: status %d", statsRec.Code)
	}
	var stats StatsResponse
	if err := json.Unmarshal(statsRec.Body.Bytes(), &stats); err != nil {
		t.Fatalf("unmarshal stats: %v", err)
	}
	if stats.Headroom.ActiveEntries > 16 {
		t.Fatalf("expected active_entries <= 16, got %d", stats.Headroom.ActiveEntries)
	}

	// Clear cache and reset stats.
	if r := doJSON(t, env.router, http.MethodPost, "/api/tokensaver/cache/clear", nil); r.Code != http.StatusNoContent {
		t.Fatalf("cache clear: status %d", r.Code)
	}
	if r := doJSON(t, env.router, http.MethodPost, "/api/tokensaver/stats/reset", nil); r.Code != http.StatusNoContent {
		t.Fatalf("stats reset: status %d", r.Code)
	}
	// Dry-run preview endpoint.
	prevRec := doJSON(t, env.router, http.MethodPost, "/api/tokensaver/preview", map[string]any{
		"protocol": "openai-chat",
		"model":    "gpt-4o",
		"body": map[string]any{
			"model":    "gpt-4o",
			"messages": []any{map[string]any{"role": "user", "content": "Hello     world\n\n\n\nline 2"}},
		},
	})
	if prevRec.Code != http.StatusOK {
		t.Fatalf("preview: status %d body %s", prevRec.Code, prevRec.Body.String())
	}
}

func TestTokenSaverFinalReviewFindings(t *testing.T) {
	env := setupTokenSaverTest(t)

	rec := doJSON(t, env.router, http.MethodPatch, "/api/tokensaver/config", map[string]any{
		"enabled":               true,
		"strategy":              "aggressive",
		"max_history_messages":  4,
		"max_tool_output_chars": 32,
		"headroom_enabled":      true,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("patch config: status %d", rec.Code)
	}

	// Finding #6: Preserve newlines adjacent to fenced ``` code blocks.
	md := "Hello   world.\n```python\ndef f():\n    return 1\n```\nDone."
	gotMD := compressText(md)
	if !strings.Contains(gotMD, "Hello world.\n```python\n") || !strings.Contains(gotMD, "\n```\nDone.") {
		t.Fatalf("newlines around code fence were stripped: %q", gotMD)
	}

	// Finding #7a: Rune-safe UTF-8 truncation in truncateToolOutput.
	utf8Tool := "日本語の長いツール出力テスト文字列です"
	trunc := truncateToolOutput(utf8Tool, 8)
	if strings.ContainsRune(trunc, '\uFFFD') {
		t.Fatalf("truncateToolOutput corrupted UTF-8: %q", trunc)
	}

	// Finding #7b: OpenAI Responses "input" gets distinct Headroom keys and compresses text.
	rA := &contracts.CanonicalRequest{
		RequestID: "resp-a",
		Protocol:  "openai-responses",
		Model:     "gpt-4o",
		Body:      []byte(`{"model":"gpt-4o","input":"First     unique     prompt"}`),
		Metadata:  make(map[string]any),
	}
	rB := &contracts.CanonicalRequest{
		RequestID: "resp-b",
		Protocol:  "openai-responses",
		Model:     "gpt-4o",
		Body:      []byte(`{"model":"gpt-4o","input":"Second    distinct   prompt"}`),
		Metadata:  make(map[string]any),
	}
	_ = env.mod.OnRequest(context.Background(), rA)
	_ = env.mod.OnRequest(context.Background(), rB)
	if rA.Metadata["headroom_cache_key"] == rB.Metadata["headroom_cache_key"] {
		t.Fatalf("expected distinct Headroom keys for different Responses inputs, both got %v", rA.Metadata["headroom_cache_key"])
	}
	if !strings.Contains(string(rA.Body), "First unique prompt") {
		t.Fatalf("expected Responses input whitespace compressed, got: %s", rA.Body)
	}

	// Finding #8: Failed request (OnRequest without OnResponseHeaders) must NOT poison Headroom cache.
	fail1 := &contracts.CanonicalRequest{
		RequestID: "fail-1",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      []byte(`{"model":"gpt-4o","messages":[{"role":"user","content":"unconfirmed prefix"}]}`),
		Metadata:  make(map[string]any),
	}
	fail2 := &contracts.CanonicalRequest{
		RequestID: "fail-2",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      []byte(`{"model":"gpt-4o","messages":[{"role":"user","content":"unconfirmed prefix"}]}`),
		Metadata:  make(map[string]any),
	}
	_ = env.mod.OnRequest(context.Background(), fail1)
	// No OnResponseHeaders for fail1 (simulating upstream 502).
	_ = env.mod.OnRequest(context.Background(), fail2)
	if fail2.Metadata["headroom_cache_hit"] != false {
		t.Fatalf("expected second request after failed first request to still be a miss, got %v", fail2.Metadata["headroom_cache_hit"])
	}
}

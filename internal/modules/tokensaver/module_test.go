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
	rec := doJSON(t, env.router, http.MethodPatch, "/api/token-saver/config", map[string]any{
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

	rec := doJSON(t, env.router, http.MethodPatch, "/api/token-saver/config", map[string]any{
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

	rec := doJSON(t, env.router, http.MethodPatch, "/api/token-saver/config", map[string]any{
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

	statsRec := doJSON(t, env.router, http.MethodGet, "/api/token-saver/stats", nil)
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
	if r := doJSON(t, env.router, http.MethodDelete, "/api/token-saver/cache", nil); r.Code != http.StatusNoContent {
		t.Fatalf("cache clear: status %d", r.Code)
	}
	if r := doJSON(t, env.router, http.MethodDelete, "/api/token-saver/stats", nil); r.Code != http.StatusNoContent {
		t.Fatalf("stats reset: status %d", r.Code)
	}
	// Dry-run preview endpoint.
	prevRec := doJSON(t, env.router, http.MethodPost, "/api/token-saver/preview", map[string]any{
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

	rec := doJSON(t, env.router, http.MethodPatch, "/api/token-saver/config", map[string]any{
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

func TestInjectAnthropicCacheControlCoversToolsAndSystem(t *testing.T) {
	body := []byte(`{"model":"claude-sonnet-4","system":"You are a coding assistant.","tools":[{"name":"read_file","description":"Reads a file"},{"name":"run_cmd","description":"Runs a command"}],"messages":[{"role":"user","content":"hi"}]}`)
	updated, changed := injectAnthropicCacheControl(body)
	if !changed {
		t.Fatal("expected cache_control injection to mutate payload")
	}

	var root map[string]any
	if err := json.Unmarshal(updated, &root); err != nil {
		t.Fatalf("unmarshal updated: %v", err)
	}
	sysBlocks, ok := root["system"].([]any)
	if !ok || len(sysBlocks) != 1 {
		t.Fatalf("expected system converted to 1 block, got %#v", root["system"])
	}
	if _, ok := sysBlocks[0].(map[string]any)["cache_control"]; !ok {
		t.Fatal("expected cache_control on system block")
	}

	tools, ok := root["tools"].([]any)
	if !ok || len(tools) != 2 {
		t.Fatalf("expected 2 tools, got %#v", root["tools"])
	}
	if _, exists := tools[0].(map[string]any)["cache_control"]; exists {
		t.Fatal("expected cache_control only on last tool, not first tool")
	}
	lastToolCC, ok := tools[1].(map[string]any)["cache_control"].(map[string]any)
	if !ok || lastToolCC["type"] != "ephemeral" {
		t.Fatalf("expected ephemeral cache_control on last tool, got %#v", tools[1])
	}

	// Second call must be a no-op (idempotent).
	if _, changedAgain := injectAnthropicCacheControl(updated); changedAgain {
		t.Fatal("expected second injectAnthropicCacheControl call to be idempotent")
	}
}

func TestTokenSaverFeatureTogglesCavemanPonytailAndTrimHistory(t *testing.T) {
	env := setupTokenSaverTest(t)

	rec := doJSON(t, env.router, http.MethodPatch, "/api/token-saver/config", map[string]any{
		"enabled":              false,
		"caveman_enabled":      true,
		"caveman_level":        "ultra",
		"ponytail_enabled":     true,
		"ponytail_level":       "full",
		"trim_history_enabled": true,
		"max_history_messages": 4,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("patch config: status %d body %s", rec.Code, rec.Body.String())
	}

	payload, _ := json.Marshal(map[string]any{
		"model": "gpt-4o",
		"messages": []any{
			map[string]any{"role": "system", "content": "Base instructions."},
			map[string]any{"role": "user", "content": "msg 1"},
			map[string]any{"role": "assistant", "content": "reply 1"},
			map[string]any{"role": "user", "content": "msg 2"},
			map[string]any{"role": "assistant", "content": "reply 2"},
			map[string]any{"role": "user", "content": "msg 3"},
		},
	})
	req := &contracts.CanonicalRequest{
		RequestID: "req-toggles",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      payload,
		Metadata:  make(map[string]any),
	}
	if err := env.mod.OnRequest(context.Background(), req); err != nil {
		t.Fatalf("OnRequest: %v", err)
	}

	var root map[string]any
	if err := json.Unmarshal(req.Body, &root); err != nil {
		t.Fatalf("unmarshal mutated body: %v", err)
	}
	msgs, ok := root["messages"].([]any)
	if !ok {
		t.Fatalf("expected messages array, got %#v", root["messages"])
	}
	if len(msgs) != 4 {
		t.Fatalf("expected trim_history_enabled to keep 4 messages, got %d", len(msgs))
	}
	sysMsg, _ := msgs[0].(map[string]any)["content"].(string)
	for _, want := range []string{
		"Base instructions.",
		"[caveman:ultra]",
		"Auto-Clarity:",
		"Preserve the user's dominant language",
		"[ponytail:full]",
		"stop at the first rung that holds:",
		"Never simplify away:",
	} {
		if !strings.Contains(sysMsg, want) {
			t.Fatalf("expected %q in injected system prompt, got: %q", want, sysMsg)
		}
	}
}

func TestTokenSaverToggleIsolationAndAllOffPassThrough(t *testing.T) {
	env := setupTokenSaverTest(t)

	// 1. RTK ON + Trim History OFF: must truncate tool output & clean whitespace,
	//    but MUST NOT drop chat history messages.
	rec := doJSON(t, env.router, http.MethodPatch, "/api/token-saver/config", map[string]any{
		"enabled":               true,
		"strategy":              "aggressive",
		"trim_history_enabled":  false,
		"max_history_messages":  4,
		"max_tool_output_chars": 40,
		"caveman_enabled":       false,
		"ponytail_enabled":      false,
		"headroom_enabled":      false,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("patch RTK only: status %d body %s", rec.Code, rec.Body.String())
	}

	payload, _ := json.Marshal(map[string]any{
		"model": "gpt-4o",
		"messages": []any{
			map[string]any{"role": "system", "content": "You   are   concise."},
			map[string]any{"role": "user", "content": "m1"},
			map[string]any{"role": "assistant", "content": "a1"},
			map[string]any{"role": "user", "content": "m2"},
			map[string]any{"role": "tool", "content": strings.Repeat("log-line-", 20)},
			map[string]any{"role": "user", "content": "m3"},
		},
	})
	req := &contracts.CanonicalRequest{
		RequestID: "req-rtk-only",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      append([]byte(nil), payload...),
		Metadata:  make(map[string]any),
	}
	if err := env.mod.OnRequest(context.Background(), req); err != nil {
		t.Fatalf("OnRequest RTK only: %v", err)
	}
	var root map[string]any
	if err := json.Unmarshal(req.Body, &root); err != nil {
		t.Fatalf("unmarshal RTK only: %v", err)
	}
	msgs := root["messages"].([]any)
	if len(msgs) != 6 {
		t.Fatalf("RTK ON with trim_history_enabled=false must keep all 6 messages, got %d", len(msgs))
	}
	if !strings.Contains(string(req.Body), "[...truncated ") {
		t.Fatalf("expected RTK to truncate long tool output, got: %s", req.Body)
	}

	// 2. Gemini payload receives Caveman directive in systemInstruction.
	rec = doJSON(t, env.router, http.MethodPatch, "/api/token-saver/config", map[string]any{
		"enabled":         false,
		"caveman_enabled": true,
		"caveman_level":   "lite",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("patch caveman: status %d", rec.Code)
	}
	gemReq := &contracts.CanonicalRequest{
		RequestID: "req-gemini",
		Protocol:  "gemini",
		Model:     "gemini-2.5-pro",
		Body:      []byte(`{"contents":[{"role":"user","parts":[{"text":"explain tcp"}]}]}`),
		Metadata:  make(map[string]any),
	}
	if err := env.mod.OnRequest(context.Background(), gemReq); err != nil {
		t.Fatalf("OnRequest gemini: %v", err)
	}
	if !strings.Contains(string(gemReq.Body), "[caveman:lite]") {
		t.Fatalf("expected Gemini systemInstruction to contain [caveman:lite], got: %s", gemReq.Body)
	}

	// 2b. Anthropic payload WITHOUT system field must write to top-level "system",
	//     NEVER prepend role="system" into messages (Anthropic rejects that with 400).
	anthReq := &contracts.CanonicalRequest{
		RequestID: "req-anthropic-nosys",
		Protocol:  "anthropic",
		Model:     "claude-sonnet-4",
		Body:      []byte(`{"model":"claude-sonnet-4","messages":[{"role":"user","content":"hi"}]}`),
		Metadata:  make(map[string]any),
	}
	if err := env.mod.OnRequest(context.Background(), anthReq); err != nil {
		t.Fatalf("OnRequest anthropic no-system: %v", err)
	}
	var anthRoot map[string]any
	if err := json.Unmarshal(anthReq.Body, &anthRoot); err != nil {
		t.Fatalf("unmarshal anthropic: %v", err)
	}
	if sysStr, ok := anthRoot["system"].(string); !ok || !strings.Contains(sysStr, "[caveman:lite]") {
		t.Fatalf("expected Anthropic top-level system string with directive, got: %#v", anthRoot["system"])
	}
	if anthMsgs := anthRoot["messages"].([]any); len(anthMsgs) != 1 {
		t.Fatalf("expected Anthropic messages untouched (len=1), got %d: %#v", len(anthMsgs), anthMsgs)
	}
	// Re-running OnRequest on already-injected Anthropic payload must be idempotent (no fallthrough).
	beforeRepeat := append([]byte(nil), anthReq.Body...)
	if err := env.mod.OnRequest(context.Background(), anthReq); err != nil {
		t.Fatalf("OnRequest anthropic repeat: %v", err)
	}
	if !bytes.Equal(anthReq.Body, beforeRepeat) {
		t.Fatalf("expected idempotent repeat on Anthropic payload, got: %s", anthReq.Body)
	}

	// 2c. OpenAI Responses payload WITHOUT instructions must set top-level "instructions".
	respReq := &contracts.CanonicalRequest{
		RequestID: "req-responses-noinst",
		Protocol:  "openai-responses",
		Model:     "gpt-4o",
		Body:      []byte(`{"model":"gpt-4o","input":"explain quicksort"}`),
		Metadata:  make(map[string]any),
	}
	if err := env.mod.OnRequest(context.Background(), respReq); err != nil {
		t.Fatalf("OnRequest responses: %v", err)
	}
	if !strings.Contains(string(respReq.Body), `"instructions":"[caveman:lite]`) {
		t.Fatalf("expected OpenAI Responses instructions injected, got: %s", respReq.Body)
	}

	// 3. All toggles OFF -> payload passes through untouched byte-for-byte.
	rec = doJSON(t, env.router, http.MethodPatch, "/api/token-saver/config", map[string]any{
		"enabled":              false,
		"headroom_enabled":     false,
		"caveman_enabled":      false,
		"ponytail_enabled":     false,
		"trim_history_enabled": false,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("patch all off: status %d", rec.Code)
	}
	offReq := &contracts.CanonicalRequest{
		RequestID: "req-all-off",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      append([]byte(nil), payload...),
		Metadata:  make(map[string]any),
	}
	if err := env.mod.OnRequest(context.Background(), offReq); err != nil {
		t.Fatalf("OnRequest all off: %v", err)
	}
	if !bytes.Equal(offReq.Body, payload) {
		t.Fatal("expected all-off config to leave payload untouched byte-for-byte")
	}
}

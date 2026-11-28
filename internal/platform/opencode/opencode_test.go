package opencode

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestSessionStoreStablePerConversationAndEvictsOnForget(t *testing.T) {
	store := NewSessionStore()
	re := regexp.MustCompile(`^ses_[0-9a-f]{12}[0-9A-Za-z]{14}$`)

	turn1 := []byte(`{"model":"big-pickle","messages":[{"role":"system","content":"You are helpful."},{"role":"user","content":"Hello world"}]}`)
	turn2 := []byte(`{"model":"big-pickle","messages":[{"role":"system","content":"You are helpful."},{"role":"user","content":"Hello world"},{"role":"assistant","content":"Hi there!"},{"role":"user","content":"Follow up"}]}`)
	otherConv := []byte(`{"model":"big-pickle","messages":[{"role":"system","content":"You are helpful."},{"role":"user","content":"Different topic"}]}`)

	s1 := store.Resolve("noauth:opencode", "claude-code", nil, turn1)
	if !re.MatchString(s1) {
		t.Fatalf("session id %q does not match canonical ses_ format", s1)
	}
	s2 := store.Resolve("noauth:opencode", "claude-code", nil, turn2)
	if s1 != s2 {
		t.Fatalf("expected same session across turns of same conversation: %q != %q", s1, s2)
	}

	sOther := store.Resolve("noauth:opencode", "claude-code", nil, otherConv)
	if sOther == s1 {
		t.Fatalf("expected distinct session for different conversation, got same %q", sOther)
	}

	// Native valid x-opencode-session header is passed through untouched.
	native := "ses_123456789abcABCDEF12345678"
	hdrs := http.Header{"X-Opencode-Session": []string{native}}
	if got := store.Resolve("noauth:opencode", "opencode", hdrs, turn1); got != native {
		t.Fatalf("expected native session %q, got %q", native, got)
	}

	// ForgetSession drops the cached minted session so next Resolve mints a new one.
	if !store.Forget(s1) {
		t.Fatalf("expected Forget(%q) to return true", s1)
	}
	s3 := store.Resolve("noauth:opencode", "claude-code", nil, turn1)
	if s3 == s1 {
		t.Fatalf("expected fresh session after Forget, still got %q", s3)
	}
}

func TestTransformRequestCloaksGatedToolsAndRestoresResponses(t *testing.T) {
	reqBody := []byte(`{
		"model": "big-pickle",
		"stream": false,
		"reasoning_effort": "auto",
		"tool_choice": {"type": "function", "function": {"name": "ReportFindings"}},
		"messages": [
			{"role": "system", "content": "Be concise."},
			{"role": "user", "content": "Check repo"},
			{"role": "assistant", "tool_calls": [{"id": "c1", "type": "function", "function": {"name": "ReportFindings", "arguments": "{}"}}]}
		],
		"tools": [
			{"type": "function", "function": {"name": "ReportFindings", "description": "Report", "parameters": {"type": "object", "properties": {}}}},
			{"type": "function", "function": {"name": "Read", "description": "Read file", "parameters": {"type": "object", "properties": {"file_path": {"type": "string"}}}}},
			{"type": "function", "function": {"name": "read", "description": "Custom read", "parameters": {"type": "object", "properties": {"uri": {"type": "string"}}}}}
		]
	}`)

	outBytes, renameMap, err := TransformRequest("opencode", "big-pickle", reqBody, false, false)
	if err != nil {
		t.Fatalf("TransformRequest error: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(outBytes, &parsed); err != nil {
		t.Fatalf("unmarshal transformed body: %v", err)
	}
	if parsed["stream"] != true {
		t.Fatalf("expected stream=true forced on opencode free, got %v", parsed["stream"])
	}
	if _, hasEffort := parsed["reasoning_effort"]; hasEffort {
		t.Fatalf("expected reasoning_effort=auto dropped on chat completions")
	}

	// Check <env> appended to existing system message without losing caller content.
	msgs := parsed["messages"].([]any)
	sysContent := msgs[0].(map[string]any)["content"].(string)
	if !strings.Contains(sysContent, "Be concise.") || !strings.Contains(sysContent, "<env>") {
		t.Fatalf("system prompt missing caller text or <env>: %q", sysContent)
	}

	// Check gated tools aliased and all 5 canonical tools present.
	tools := parsed["tools"].([]any)
	toolNames := make(map[string]bool)
	for _, raw := range tools {
		fn := raw.(map[string]any)["function"].(map[string]any)
		toolNames[fn["name"].(string)] = true
	}
	for _, canonical := range []string{"edit", "glob", "grep", "read", "shell"} {
		if !toolNames[canonical] {
			t.Fatalf("missing canonical tool %q in %v", canonical, toolNames)
		}
	}
	if toolNames["ReportFindings"] || toolNames["Read"] {
		t.Fatalf("gated tool names should have been aliased: %v", toolNames)
	}
	if len(renameMap) < 3 {
		t.Fatalf("expected at least 3 renamed tools (ReportFindings, Read, read), got %v", renameMap)
	}

	// Verify tool_choice was retargeted to the alias of ReportFindings.
	tc := parsed["tool_choice"].(map[string]any)["function"].(map[string]any)["name"].(string)
	if renameMap[tc] != "ReportFindings" {
		t.Fatalf("tool_choice name %q does not map back to ReportFindings in %v", tc, renameMap)
	}

	// Verify SSE restoration restores aliased tool name back to ReportFindings.
	sseInput := "data: {\"id\":\"chatcmpl-1\",\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"c2\",\"type\":\"function\",\"function\":{\"name\":\"" + tc + "\",\"arguments\":\"{\\\"ok\\\":true}\"}}]}}]}\n\ndata: [DONE]\n\n"
	restoredReader := NewToolNameRestoringSSEReader(io.NopCloser(strings.NewReader(sseInput)), renameMap)
	restoredBytes, err := io.ReadAll(restoredReader)
	_ = restoredReader.Close()
	if err != nil {
		t.Fatalf("read restored SSE: %v", err)
	}
	if !strings.Contains(string(restoredBytes), `"name":"ReportFindings"`) {
		t.Fatalf("expected SSE restoration to restore ReportFindings, got: %s", string(restoredBytes))
	}

	// Verify SSE-to-JSON aggregation for non-stream clients.
	aggSSE := strings.Join([]string{
		`data: {"id":"chatcmpl-9","model":"big-pickle","choices":[{"index":0,"delta":{"role":"assistant","reasoning_content":"think ","content":"Hello "}}]}`,
		`data: {"id":"chatcmpl-9","model":"big-pickle","choices":[{"index":0,"delta":{"content":"world","tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"` + tc + `","arguments":"{\"a\":1}"}}]},"finish_reason":"tool_calls"}]}`,
		`data: {"id":"chatcmpl-9","model":"big-pickle","choices":[],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}`,
		`data: [DONE]`,
	}, "\n") + "\n"
	aggJSON, err := AggregateOpenAISSEToJSON([]byte(aggSSE), "big-pickle", renameMap)
	if err != nil {
		t.Fatalf("AggregateOpenAISSEToJSON error: %v", err)
	}
	var agg map[string]any
	if err := json.Unmarshal(aggJSON, &agg); err != nil {
		t.Fatalf("unmarshal aggregated JSON: %v", err)
	}
	choice := agg["choices"].([]any)[0].(map[string]any)
	msg := choice["message"].(map[string]any)
	if msg["content"] != "Hello world" {
		t.Fatalf("aggregated content = %v, want Hello world", msg["content"])
	}
	if msg["reasoning_content"] != "think " {
		t.Fatalf("aggregated reasoning_content = %v, want 'think '", msg["reasoning_content"])
	}
	calls := msg["tool_calls"].([]any)
	if calls[0].(map[string]any)["function"].(map[string]any)["name"] != "ReportFindings" {
		t.Fatalf("expected restored tool call name ReportFindings, got %v", calls[0])
	}
}

func TestTransformRequestEnforcesToolSurfaceCapAndResponsesSanitization(t *testing.T) {
	// Build a responses request with prior reasoning items and huge extra tools.
	var bigTools []map[string]any
	bigDesc := strings.Repeat("x", 6000)
	for i := 0; i < 10; i++ {
		bigTools = append(bigTools, map[string]any{
			"type":        "function",
			"name":        "extra_tool_" + string(rune('a'+i)),
			"description": bigDesc,
			"parameters":  map[string]any{"type": "object", "properties": map[string]any{}},
		})
	}
	payload := map[string]any{
		"model":            "muse-wonder-1.0",
		"max_tokens":       1024,
		"reasoning_effort": "max",
		"tool_choice":      "none",
		"input": []map[string]any{
			{"type": "reasoning", "encrypted_content": "secret-blob"},
			{"type": "message", "role": "user", "content": "hi"},
		},
		"tools": bigTools,
	}
	raw, _ := json.Marshal(payload)

	outBytes, _, err := TransformRequest("opencode", "muse-wonder-1.0", raw, true, false)
	if err != nil {
		t.Fatalf("TransformRequest responses error: %v", err)
	}
	var parsed map[string]any
	_ = json.Unmarshal(outBytes, &parsed)

	if parsed["store"] != false {
		t.Fatalf("expected store=false, got %v", parsed["store"])
	}
	if parsed["max_output_tokens"] != float64(1024) {
		t.Fatalf("expected max_output_tokens=1024, got %v", parsed["max_output_tokens"])
	}
	if _, hasMax := parsed["max_tokens"]; hasMax {
		t.Fatalf("expected max_tokens removed")
	}
	if parsed["tool_choice"] != "auto" {
		t.Fatalf("expected tool_choice=auto on responses after stripping none, got %v", parsed["tool_choice"])
	}
	inputs := parsed["input"].([]any)
	if len(inputs) != 1 || inputs[0].(map[string]any)["type"] != "message" {
		t.Fatalf("expected prior reasoning item stripped, got %v", inputs)
	}
	toolsBytes, _ := json.Marshal(parsed["tools"])
	if len(toolsBytes) > MaxToolSurfaceChars {
		t.Fatalf("tools exceeded WAF cap %d: got %d", MaxToolSurfaceChars, len(toolsBytes))
	}
}

func TestResolveUpstreamURLMultiEndpoint(t *testing.T) {
	tests := []struct {
		provider string
		baseURL  string
		model    string
		proto    string
		stream   bool
		want     string
	}{
		{"opencode", "", "big-pickle", "openai-chat", true, "https://opencode.ai/inference/openai/v1/chat/completions"},
		{"opencode", "", "muse-wonder-1.0", "openai-responses", true, "https://opencode.ai/inference/openai/v1/responses"},
		{"opencode", "", "union-alpha", "anthropic", true, "https://opencode.ai/inference/openai/v1/messages"},
		{"opencode-zen", "", "claude-sonnet-4-6", "anthropic", true, "https://opencode.ai/zen/v1/messages"},
		{"opencode-zen", "", "gpt-5.4", "openai-responses", true, "https://opencode.ai/zen/v1/responses"},
		{"opencode-go", "", "gemini-4-argon-contributor", "gemini", true, "https://opencode.ai/inference/go/v1/models/gemini-4-argon-contributor:streamGenerateContent?alt=sse"},
		{"opencode-go", "", "glm-5.3", "openai-chat", false, "https://opencode.ai/zen/go/v1/chat/completions"},
	}
	for _, tc := range tests {
		got := ResolveUpstreamURL(tc.provider, tc.baseURL, tc.model, tc.proto, tc.stream)
		if got != tc.want {
			t.Errorf("ResolveUpstreamURL(%q,%q,%q,%q,%v) = %q, want %q", tc.provider, tc.baseURL, tc.model, tc.proto, tc.stream, got, tc.want)
		}
	}
}

func TestClassifyRejectionAndParseResetAfter(t *testing.T) {
	if got := ClassifyRejection(http.StatusTooManyRequests, []byte(`{"error":{"message":"The inference service is overloaded. Please try again later."}}`)); got != RejectionOverload {
		t.Fatalf("expected RejectionOverload, got %v", got)
	}
	if got := ClassifyRejection(http.StatusTooManyRequests, []byte(`{"error":{"type":"FreeUsageLimitError","resets_in_seconds":45}}`)); got != RejectionSession {
		t.Fatalf("expected RejectionSession, got %v", got)
	}
	if got := ClassifyRejection(http.StatusForbidden, []byte(`{"error":{"type":"FreeTierError"}}`)); got != RejectionSession {
		t.Fatalf("expected RejectionSession for 403, got %v", got)
	}

	now := time.Unix(1700000000, 0)
	if d := ParseResetAfter(nil, []byte(`{"error":{"resets_at":1700000090}}`), now); d != 90*time.Second {
		t.Fatalf("expected 90s from resets_at, got %v", d)
	}
	if d := ParseResetAfter(nil, []byte(`{"error":{"resets_in_seconds":30}}`), now); d != 30*time.Second {
		t.Fatalf("expected 30s from resets_in_seconds, got %v", d)
	}
	if d := ParseResetAfter(http.Header{"Retry-After": []string{"15"}}, nil, now); d != 15*time.Second {
		t.Fatalf("expected 15s from Retry-After, got %v", d)
	}
}

func TestVersionResolverRefreshesFromUpdater(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"2.9.3"}`))
	}))
	defer srv.Close()

	vr := NewVersionResolver(srv.URL, srv.Client())
	if got := vr.UserAgent(); !strings.HasPrefix(got, "opencode/latest/") {
		t.Fatalf("unexpected cold-start UA: %q", got)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if vr.Version() == "2.9.3" {
			break
		}
		time.Sleep(15 * time.Millisecond)
	}
	if got := vr.UserAgent(); got != "opencode/latest/2.9.3/cli" {
		t.Fatalf("expected refreshed UA opencode/latest/2.9.3/cli, got %q", got)
	}
}

func TestPacerSpacesRequestsAndRespectsCancel(t *testing.T) {
	pacer := NewPacer(30 * time.Millisecond)
	start := time.Now()
	if err := pacer.Wait(context.Background()); err != nil {
		t.Fatalf("first Wait: %v", err)
	}
	if err := pacer.Wait(context.Background()); err != nil {
		t.Fatalf("second Wait: %v", err)
	}
	if elapsed := time.Since(start); elapsed < 25*time.Millisecond {
		t.Fatalf("expected pacer gap >= 25ms, got %v", elapsed)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := pacer.Wait(ctx); err == nil {
		t.Fatalf("expected error on canceled context")
	}
	_ = bytes.NewReader(nil)
}

func TestResolveTargetProtocolMatchesModelAndPassthrough(t *testing.T) {
	tests := []struct {
		provider string
		model    string
		inProto  string
		defProto string
		want     string
	}{
		{"opencode", "big-pickle", "anthropic", "openai-chat", "openai-chat"},
		{"opencode", "muse-wonder-1.0", "openai-chat", "openai-chat", "openai-responses"},
		{"opencode", "union-alpha", "openai-chat", "openai-chat", "anthropic"},
		{"opencode-zen", "claude-sonnet-4-6", "openai-chat", "openai-chat", "anthropic"},
		{"opencode-zen", "claude-sonnet-4-6", "anthropic", "openai-chat", "anthropic"},
		{"opencode-zen", "gpt-5.4", "openai-chat", "openai-chat", "openai-responses"},
		{"opencode-go", "gemini-4-argon-contributor", "openai-chat", "openai-chat", "gemini"},
		{"opencode-go", "deepseek-v4.1-flash", "anthropic", "openai-chat", "anthropic"},
		{"opencode-go", "glm-5.3", "anthropic", "openai-chat", "openai-chat"},
	}
	for _, tc := range tests {
		got := ResolveTargetProtocol(tc.provider, tc.model, tc.inProto, tc.defProto)
		if got != tc.want {
			t.Errorf("ResolveTargetProtocol(%q,%q,%q,%q) = %q, want %q", tc.provider, tc.model, tc.inProto, tc.defProto, got, tc.want)
		}
	}
}

func TestAggregateSSEToJSONSupportsResponsesAndAnthropic(t *testing.T) {
	renameMap := map[string]string{"run_command": "ReportFindings"}

	// 1. OpenAI Responses SSE -> extracts response.completed object and restores tool name.
	respSSE := strings.Join([]string{
		`event: response.created`,
		`data: {"type":"response.created","response":{"id":"resp_1","status":"in_progress"}}`,
		`event: response.completed`,
		`data: {"type":"response.completed","response":{"id":"resp_1","status":"completed","model":"muse-wonder-1.0","output":[{"type":"function_call","id":"fc_1","call_id":"call_1","name":"run_command","arguments":"{\"ok\":true}"}]}}`,
	}, "\n") + "\n"
	outResp, err := AggregateSSEToJSON("openai-responses", []byte(respSSE), "muse-wonder-1.0", renameMap)
	if err != nil {
		t.Fatalf("AggregateSSEToJSON(openai-responses) error: %v", err)
	}
	if !strings.Contains(string(outResp), `"name":"ReportFindings"`) || !strings.Contains(string(outResp), `"id":"resp_1"`) {
		t.Fatalf("unexpected responses aggregation: %s", string(outResp))
	}

	// 2. Anthropic Messages SSE -> builds message JSON and restores tool_use name.
	antSSE := strings.Join([]string{
		`event: message_start`,
		`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"union-alpha","content":[],"usage":{"input_tokens":11,"output_tokens":0}}}`,
		`event: content_block_start`,
		`data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
		`event: content_block_delta`,
		`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hi "}}`,
		`event: content_block_delta`,
		`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"there"}}`,
		`event: content_block_start`,
		`data: {"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_1","name":"run_command","input":{}}}`,
		`event: content_block_delta`,
		`data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"a\":1}"}}`,
		`event: message_delta`,
		`data: {"type":"message_delta","delta":{"stop_reason":"tool_use"},"usage":{"output_tokens":7}}`,
	}, "\n") + "\n"
	outAnt, err := AggregateSSEToJSON("anthropic", []byte(antSSE), "union-alpha", renameMap)
	if err != nil {
		t.Fatalf("AggregateSSEToJSON(anthropic) error: %v", err)
	}
	var antObj map[string]any
	if err := json.Unmarshal(outAnt, &antObj); err != nil {
		t.Fatalf("unmarshal anthropic agg: %v", err)
	}
	if antObj["stop_reason"] != "tool_use" {
		t.Fatalf("stop_reason = %v, want tool_use", antObj["stop_reason"])
	}
	blocks := antObj["content"].([]any)
	if len(blocks) != 2 {
		t.Fatalf("expected 2 content blocks, got %d: %s", len(blocks), string(outAnt))
	}
	if blocks[0].(map[string]any)["text"] != "Hi there" {
		t.Fatalf("text block = %v, want 'Hi there'", blocks[0])
	}
	toolBlock := blocks[1].(map[string]any)
	if toolBlock["name"] != "ReportFindings" {
		t.Fatalf("tool_use name = %v, want ReportFindings", toolBlock["name"])
	}
}

func TestResetClientEvictsPooledTLSClient(t *testing.T) {
	base := &http.Client{Transport: &http.Transport{}}
	c1 := ResolveClient(context.Background(), base, "https://opencode.ai/inference/openai/v1/chat/completions")
	c2 := ResolveClient(context.Background(), base, "https://opencode.ai/inference/openai/v1/chat/completions")
	if c1 != c2 {
		t.Fatalf("expected pooled client reuse before reset")
	}
	ResetClient(context.Background(), base, "https://opencode.ai/inference/openai/v1/chat/completions")
	c3 := ResolveClient(context.Background(), base, "https://opencode.ai/inference/openai/v1/chat/completions")
	if c3 == c1 {
		t.Fatalf("expected fresh pooled client after ResetClient")
	}
}

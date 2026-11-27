package plugins_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

func putPluginConfig(t *testing.T, handler http.Handler, key, body string) {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/plugins/"+key, strings.NewReader(body))
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT /api/plugins/%s = %d: %s", key, rec.Code, rec.Body.String())
	}
}

func TestPluginsProtocolAwareReasoningAndSpeedMode(t *testing.T) {
	mod, handler := setupPluginsTest(t)

	putPluginConfig(t, handler, "thinkDeeper", `{"enabled":true,"models":["claude-*","o3*"],"config":{"effort":"high","budget_tokens":16384}}`)
	putPluginConfig(t, handler, "speedMode", `{"enabled":true,"models":["gpt-4o-mini","claude-haiku*"],"config":{"max_tokens":512}}`)

	// Claude thinkDeeper uses native thinking.budget_tokens.
	claudeReq := &contracts.CanonicalRequest{
		RequestID: "req_td_claude",
		Protocol:  "anthropic-messages",
		Model:     "claude-3-7-sonnet",
		Body:      []byte(`{"model":"claude-3-7-sonnet","messages":[{"role":"user","content":"Prove it"}]}`),
		Metadata:  map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), claudeReq); err != nil {
		t.Fatalf("OnRequest claude thinkDeeper: %v", err)
	}
	var claudeOut map[string]any
	if err := json.Unmarshal(claudeReq.Body, &claudeOut); err != nil {
		t.Fatalf("unmarshal claudeOut: %v", err)
	}
	thinking, _ := claudeOut["thinking"].(map[string]any)
	if thinking["type"] != "enabled" || thinking["budget_tokens"] != float64(16384) {
		t.Fatalf("expected claude thinking enabled with budget 16384, got %v", claudeOut["thinking"])
	}

	// OpenAI thinkDeeper uses reasoning_effort + reasoning.effort.
	oaiReq := &contracts.CanonicalRequest{
		RequestID: "req_td_oai",
		Protocol:  "openai-chat",
		Model:     "o3-mini",
		Body:      []byte(`{"model":"o3-mini","messages":[{"role":"user","content":"Solve"}]}`),
		Metadata:  map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), oaiReq); err != nil {
		t.Fatalf("OnRequest openai thinkDeeper: %v", err)
	}
	var oaiOut map[string]any
	if err := json.Unmarshal(oaiReq.Body, &oaiOut); err != nil {
		t.Fatalf("unmarshal oaiOut: %v", err)
	}
	if oaiOut["reasoning_effort"] != "high" {
		t.Fatalf("expected reasoning_effort=high, got %v", oaiOut["reasoning_effort"])
	}

	// OpenAI speedMode disables reasoning_effort and applies optional max_tokens.
	speedReq := &contracts.CanonicalRequest{
		RequestID: "req_sm_oai",
		Protocol:  "openai-chat",
		Model:     "gpt-4o-mini",
		Body:      []byte(`{"model":"gpt-4o-mini","enable_thinking":true,"messages":[{"role":"user","content":"Hi"}]}`),
		Metadata:  map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), speedReq); err != nil {
		t.Fatalf("OnRequest openai speedMode: %v", err)
	}
	var speedOut map[string]any
	if err := json.Unmarshal(speedReq.Body, &speedOut); err != nil {
		t.Fatalf("unmarshal speedOut: %v", err)
	}
	if speedOut["reasoning_effort"] != "none" || speedOut["max_tokens"] != float64(512) {
		t.Fatalf("expected reasoning_effort=none and max_tokens=512, got %v", speedOut)
	}
	if _, exists := speedOut["enable_thinking"]; exists {
		t.Fatalf("expected enable_thinking removed in speedMode")
	}
}

func TestPluginsImageVisionOnlyInjectsWhenImagesPresent(t *testing.T) {
	mod, handler := setupPluginsTest(t)
	putPluginConfig(t, handler, "imageVision", `{"enabled":true,"models":[]}`)

	textOnly := &contracts.CanonicalRequest{
		RequestID: "req_iv_text",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      []byte(`{"model":"gpt-4o","messages":[{"role":"user","content":"Hello"}]}`),
		Metadata:  map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), textOnly); err != nil {
		t.Fatalf("OnRequest textOnly: %v", err)
	}
	if strings.Contains(string(textOnly.Body), "Image Vision is active") {
		t.Fatalf("expected no visual directive when no images present, got %s", string(textOnly.Body))
	}
	if active, _ := textOnly.Metadata["plugins.vision_active"].(bool); !active {
		t.Fatalf("expected plugins.vision_active=true in metadata")
	}

	withImage := &contracts.CanonicalRequest{
		RequestID: "req_iv_img",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      []byte(`{"model":"gpt-4o","messages":[{"role":"user","content":[{"type":"text","text":"What is this?"},{"type":"image_url","image_url":{"url":"https://example.com/a.png"}}]}]}`),
		Metadata:  map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), withImage); err != nil {
		t.Fatalf("OnRequest withImage: %v", err)
	}
	if !strings.Contains(string(withImage.Body), "Image Vision is active") {
		t.Fatalf("expected visual directive when image_url present, got %s", string(withImage.Body))
	}
}

func TestPluginsContextSqueezerTokenAwareEvictionAndElision(t *testing.T) {
	mod, handler := setupPluginsTest(t)
	putPluginConfig(t, handler, "contextSqueezer", `{"enabled":true,"models":[],"config":{"context_window":1024,"fill_ratio":0.75,"min_turns":2}}`)

	longText := strings.Repeat("A", 5000)
	payload := map[string]any{
		"model": "gpt-4o",
		"messages": []any{
			map[string]any{"role": "system", "content": "You are helpful."},
			map[string]any{"role": "user", "content": strings.Repeat("old question ", 120)},
			map[string]any{"role": "assistant", "content": strings.Repeat("old answer ", 120)},
			map[string]any{"role": "user", "content": longText},
			map[string]any{"role": "user", "content": "Latest question"},
		},
	}
	raw, _ := json.Marshal(payload)
	req := &contracts.CanonicalRequest{
		RequestID: "req_cs",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      raw,
		Metadata:  map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), req); err != nil {
		t.Fatalf("OnRequest contextSqueezer: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(req.Body, &out); err != nil {
		t.Fatalf("unmarshal squeezed: %v", err)
	}
	msgs, _ := out["messages"].([]any)
	if len(msgs) != 5 {
		t.Fatalf("expected 5 messages (system + recap + remaining 3 turns within budget), got %d: %s", len(msgs), string(req.Body))
	}
	recapMsg, _ := msgs[1].(map[string]any)
	if content, _ := recapMsg["content"].(string); !strings.Contains(content, "[Context Squeezer]") {
		t.Fatalf("expected recap system message, got %v", recapMsg)
	}
	shortenedMsg, _ := msgs[3].(map[string]any)
	if content, _ := shortenedMsg["content"].(string); !strings.Contains(content, "characters elided by Context Squeezer") {
		t.Fatalf("expected oversized message to be elided, got len=%d", len(content))
	}
}

func TestPluginsJSONGuardResponseRepairAndSchemaFiltering(t *testing.T) {
	mod, handler := setupPluginsTest(t)
	putPluginConfig(t, handler, "jsonGuard", `{"enabled":true,"models":["gpt-4o"]}`)

	req := &contracts.CanonicalRequest{
		RequestID: "req_jg",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body: []byte(`{
			"model":"gpt-4o",
			"messages":[{"role":"user","content":"Return JSON"}],
			"tools":[{
				"type":"function",
				"function":{
					"name":"read_file",
					"parameters":{
						"type":"object",
						"properties":{"path":{"type":"string"},"recursive":{"type":"boolean"}},
						"required":["path","recursive"]
					}
				}
			}]
		}`),
		Metadata: map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), req); err != nil {
		t.Fatalf("OnRequest jsonGuard: %v", err)
	}

	// Simulate messy model response: fenced JSON with Python literals + trailing comma,
	// and tool call with undeclared arg + missing required boolean.
	respBody, _ := json.Marshal(map[string]any{
		"choices": []any{
			map[string]any{
				"message": map[string]any{
					"role":    "assistant",
					"content": "Sure! Here is the result:\n```json\n{'ok': True, 'items': [1, 2,],}\n```",
					"tool_calls": []any{
						map[string]any{
							"id":   "call_1",
							"type": "function",
							"function": map[string]any{
								"name":      "read_file",
								"arguments": "{'path': '/tmp/a', 'hallucinated': 99,}",
							},
						},
					},
				},
			},
		},
	})

	meta := &contracts.ResponseMeta{
		RequestID: "req_jg",
		Model:     "gpt-4o",
		Metadata:  req.Metadata,
	}
	repairedBytes, err := mod.OnResponseBody(context.Background(), meta, respBody)
	if err != nil {
		t.Fatalf("OnResponseBody jsonGuard: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(repairedBytes, &out); err != nil {
		t.Fatalf("unmarshal repaired: %v", err)
	}
	choice := out["choices"].([]any)[0].(map[string]any)
	msg := choice["message"].(map[string]any)

	var parsedContent map[string]any
	if err := json.Unmarshal([]byte(msg["content"].(string)), &parsedContent); err != nil {
		t.Fatalf("expected content repaired to valid JSON, got %q: %v", msg["content"], err)
	}
	if parsedContent["ok"] != true {
		t.Fatalf("expected ok=true, got %v", parsedContent)
	}

	tc := msg["tool_calls"].([]any)[0].(map[string]any)
	fn := tc["function"].(map[string]any)
	var args map[string]any
	if err := json.Unmarshal([]byte(fn["arguments"].(string)), &args); err != nil {
		t.Fatalf("expected valid tool arguments JSON, got %q", fn["arguments"])
	}
	if _, hasHallucinated := args["hallucinated"]; hasHallucinated {
		t.Fatalf("expected undeclared arg dropped, got %v", args)
	}
	if args["path"] != "/tmp/a" || args["recursive"] != false {
		t.Fatalf("expected path=/tmp/a and required recursive=false filled, got %v", args)
	}
}

func TestPluginsOpenAIToolBridgeExtractsProseCallsAndDropsHallucinations(t *testing.T) {
	mod, handler := setupPluginsTest(t)
	putPluginConfig(t, handler, "openaiToolBridge", `{"enabled":true,"models":["kimi-*"]}`)

	req := &contracts.CanonicalRequest{
		RequestID: "req_tb",
		Protocol:  "openai-chat",
		Model:     "provider/kimi-web",
		Body: []byte(`{
			"model":"provider/kimi-web",
			"messages":[{"role":"user","content":"Check the file"}],
			"tools":[{
				"type":"function",
				"function":{
					"name":"read_file",
					"parameters":{"type":"object","properties":{"path":{"type":"string"}}}
				}
			}]
		}`),
		Metadata: map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), req); err != nil {
		t.Fatalf("OnRequest toolBridge: %v", err)
	}

	respBody := []byte(`{
		"choices":[{
			"message":{
				"role":"assistant",
				"content":"I will inspect the file now.\n<tool_call>{\"name\":\"read-file\",\"arguments\":{\"path\":\"README.md\",\"extra\":123}}</tool_call>\n<tool_call>{\"name\":\"delete_everything\",\"arguments\":{}}</tool_call>"
			},
			"finish_reason":"stop"
		}]
	}`)

	meta := &contracts.ResponseMeta{
		RequestID: "req_tb",
		Model:     "provider/kimi-web",
		Metadata:  req.Metadata,
	}
	bridgedBytes, err := mod.OnResponseBody(context.Background(), meta, respBody)
	if err != nil {
		t.Fatalf("OnResponseBody toolBridge: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(bridgedBytes, &out); err != nil {
		t.Fatalf("unmarshal bridged: %v", err)
	}
	choice := out["choices"].([]any)[0].(map[string]any)
	if choice["finish_reason"] != "tool_calls" {
		t.Fatalf("expected finish_reason=tool_calls, got %v", choice["finish_reason"])
	}
	msg := choice["message"].(map[string]any)
	if msg["content"] != "I will inspect the file now." {
		t.Fatalf("expected cleaned prose content, got %q", msg["content"])
	}
	tcs, _ := msg["tool_calls"].([]any)
	if len(tcs) != 1 {
		t.Fatalf("expected 1 accepted tool_call (hallucinated tool dropped), got %d: %v", len(tcs), tcs)
	}
	fn := tcs[0].(map[string]any)["function"].(map[string]any)
	if fn["name"] != "read_file" {
		t.Fatalf("expected normalized tool name read_file, got %v", fn["name"])
	}
	if fn["arguments"] != `{"path":"README.md"}` {
		t.Fatalf("expected undeclared arg filtered out, got %s", fn["arguments"])
	}
}

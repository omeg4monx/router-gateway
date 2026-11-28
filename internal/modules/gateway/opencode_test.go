package gateway_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/gateway"
	"github.com/omeg4monx/router-gateway/internal/platform/opencode"
)

func TestOpenCodeFreeEndToEndToolCloakingOverloadRetryAndSessionLifecycle(t *testing.T) {
	t.Parallel()

	var (
		mu          sync.Mutex
		seenSession []string
		calls       atomic.Int32
	)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if r.URL.Path != "/inference/openai/v1/chat/completions" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer public" {
			t.Errorf("Authorization = %q, want Bearer public", got)
		}
		sid := r.Header.Get("x-opencode-session")
		if !opencode.IsCanonicalSessionID(sid) || r.Header.Get("x-opencode-session-id") != sid {
			t.Errorf("invalid session headers: x-opencode-session=%q x-opencode-session-id=%q", sid, r.Header.Get("x-opencode-session-id"))
		}
		if !opencode.HasValidV2UserAgent(r.Header.Get("User-Agent")) {
			t.Errorf("invalid v2 User-Agent: %q", r.Header.Get("User-Agent"))
		}

		mu.Lock()
		seenSession = append(seenSession, sid)
		mu.Unlock()

		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if body["stream"] != true {
			t.Errorf("expected forced stream=true, got %v", body["stream"])
		}

		// Find aliased name of ReportFindings in tools.
		alias := ""
		if tools, ok := body["tools"].([]any); ok {
			for _, rawTool := range tools {
				fn, _ := rawTool.(map[string]any)["function"].(map[string]any)
				name, _ := fn["name"].(string)
				if name == "ReportFindings" {
					t.Errorf("gated tool ReportFindings leaked unaliased onto wire")
				}
				if strings.HasPrefix(name, "run_command") {
					alias = name
				}
			}
		}

		// Call 1: transient overload 429 -> gateway must retry automatically.
		if n == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"message":"The inference service is overloaded. Please try again later."}}`))
			return
		}

		// Call 4: 403 session rejection -> gateway must evict cached session.
		if n == 4 {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":{"type":"FreeTierError","message":"session rejected"}}`))
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(
			`data: {"id":"chatcmpl-oc-1","model":"big-pickle","choices":[{"index":0,"delta":{"role":"assistant","content":"Done!","tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"` + alias + `","arguments":"{\"ok\":true}"}}]}}]}` + "\n\n" +
				`data: {"id":"chatcmpl-oc-1","model":"big-pickle","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}],"usage":{"prompt_tokens":12,"completion_tokens":4,"total_tokens":16}}` + "\n\n" +
				"data: [DONE]\n\n",
		))
	}))
	defer upstream.Close()

	env := setupGatewayTest(t, nil, 0, true)
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ConnectionID:  opencode.VirtualConnectionID,
			ProviderID:    "opencode",
			ProviderType:  gateway.ProtocolOpenAIChat,
			UpstreamModel: "big-pickle",
			Endpoint:      upstream.URL + "/inference/openai/v1",
		}}, nil
	}

	reqPayload := `{
		"model":"opencode/big-pickle",
		"stream":false,
		"messages":[{"role":"user","content":"Audit this repo"}],
		"tools":[{"type":"function","function":{"name":"ReportFindings","description":"Report","parameters":{"type":"object","properties":{}}}}]
	}`

	// 1. Non-stream request: retries 429 overload on call 1, succeeds on call 2, aggregates SSE -> JSON, restores ReportFindings.
	req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewBufferString(reqPayload))
	req.Header.Set("Authorization", "Bearer sk-valid")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("non-stream request failed: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non-stream status = %d, want 200: %s", resp.StatusCode, string(body))
	}
	if !strings.Contains(string(body), `"name":"ReportFindings"`) {
		t.Fatalf("non-stream response did not restore ReportFindings: %s", string(body))
	}

	// 2. Streaming request (call 3): reuses same conversation session & restores ReportFindings in SSE.
	streamPayload := strings.Replace(reqPayload, `"stream":false`, `"stream":true`, 1)
	reqStream, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewBufferString(streamPayload))
	reqStream.Header.Set("Authorization", "Bearer sk-valid")
	reqStream.Header.Set("Content-Type", "application/json")
	respStream, err := http.DefaultClient.Do(reqStream)
	if err != nil {
		t.Fatalf("stream request failed: %v", err)
	}
	streamBody, _ := io.ReadAll(respStream.Body)
	_ = respStream.Body.Close()
	if respStream.StatusCode != http.StatusOK {
		t.Fatalf("stream status = %d, want 200: %s", respStream.StatusCode, string(streamBody))
	}
	if !strings.Contains(string(streamBody), `"name":"ReportFindings"`) {
		t.Fatalf("stream SSE did not restore ReportFindings: %s", string(streamBody))
	}

	// 3. Call 4 gets 403 -> evicts cached session; Call 5 mints a brand new session ID.
	for i := 0; i < 2; i++ {
		r, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewBufferString(reqPayload))
		r.Header.Set("Authorization", "Bearer sk-valid")
		r.Header.Set("Content-Type", "application/json")
		res, doErr := http.DefaultClient.Do(r)
		if doErr != nil {
			t.Fatalf("request %d failed: %v", i+4, doErr)
		}
		_, _ = io.ReadAll(res.Body)
		_ = res.Body.Close()
	}

	mu.Lock()
	defer mu.Unlock()
	if len(seenSession) != 5 {
		t.Fatalf("expected 5 upstream calls, got %d", len(seenSession))
	}
	if seenSession[0] != seenSession[1] || seenSession[1] != seenSession[2] || seenSession[2] != seenSession[3] {
		t.Fatalf("expected stable session across calls 1-4, got %v", seenSession)
	}
	if seenSession[4] == seenSession[3] {
		t.Fatalf("expected fresh session on call 5 after 403 eviction, still got %q", seenSession[4])
	}
}

func TestOpenCodeMultiEndpointProtocolTranslationAndResponsesAggregation(t *testing.T) {
	t.Parallel()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)

		switch r.URL.Path {
		case "/zen/v1/messages":
			if r.Header.Get("x-api-key") != "sk-zen" {
				t.Errorf("zen /messages x-api-key = %q, want sk-zen", r.Header.Get("x-api-key"))
			}
			if r.Header.Get("anthropic-version") == "" {
				t.Errorf("zen /messages missing anthropic-version header")
			}
			if _, hasMaxTokens := body["max_tokens"]; !hasMaxTokens {
				t.Errorf("zen /messages expected translated Anthropic payload with max_tokens, got %s", string(raw))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"msg_zen_1","type":"message","role":"assistant","model":"claude-sonnet-4-6","content":[{"type":"text","text":"Hello from Zen Claude"}],"stop_reason":"end_turn","usage":{"input_tokens":5,"output_tokens":4}}`))
		case "/inference/openai/v1/responses":
			if _, hasInput := body["input"]; !hasInput {
				t.Errorf("free /responses expected translated Responses payload with input, got %s", string(raw))
			}
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = w.Write([]byte(
				"event: response.completed\n" +
					`data: {"type":"response.completed","response":{"id":"resp_oc_1","status":"completed","model":"muse-wonder-1.0","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"Hello from Muse Wonder"}]}],"usage":{"input_tokens":8,"output_tokens":4}}}` + "\n\n",
			))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer upstream.Close()

	env := setupGatewayTest(t, nil, 0, true)
	env.providers.creds["conn-zen-1"] = contracts.UpstreamCredentials{
		ConnectionID: "conn-zen-1",
		ProviderID:   "opencode-zen",
		Protocol:     gateway.ProtocolOpenAIChat,
		BaseURL:      upstream.URL + "/zen/v1",
		AuthType:     "api_key",
		APIKey:       "sk-zen",
	}

	// 1. OpenAI Chat inbound -> opencode-zen/claude-sonnet-4-6 translates to Anthropic /zen/v1/messages.
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ConnectionID:  "conn-zen-1",
			ProviderID:    "opencode-zen",
			ProviderType:  gateway.ProtocolOpenAIChat,
			UpstreamModel: "claude-sonnet-4-6",
			Endpoint:      upstream.URL + "/zen/v1",
		}}, nil
	}
	reqZen, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewBufferString(`{
		"model":"opencode-zen/claude-sonnet-4-6",
		"messages":[{"role":"user","content":"say hi"}]
	}`))
	reqZen.Header.Set("Authorization", "Bearer sk-valid")
	reqZen.Header.Set("Content-Type", "application/json")
	respZen, err := http.DefaultClient.Do(reqZen)
	if err != nil {
		t.Fatalf("zen request failed: %v", err)
	}
	zenBody, _ := io.ReadAll(respZen.Body)
	_ = respZen.Body.Close()
	if respZen.StatusCode != http.StatusOK || !strings.Contains(string(zenBody), "Hello from Zen Claude") {
		t.Fatalf("zen response status=%d body=%s", respZen.StatusCode, string(zenBody))
	}

	// 2. OpenAI Chat non-stream inbound -> opencode/muse-wonder-1.0 translates to /responses & aggregates Responses SSE.
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ConnectionID:  opencode.VirtualConnectionID,
			ProviderID:    "opencode",
			ProviderType:  gateway.ProtocolOpenAIChat,
			UpstreamModel: "muse-wonder-1.0",
			Endpoint:      upstream.URL + "/inference/openai/v1",
		}}, nil
	}
	reqFreeResp, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewBufferString(`{
		"model":"opencode/muse-wonder-1.0",
		"stream":false,
		"messages":[{"role":"user","content":"say hi"}]
	}`))
	reqFreeResp.Header.Set("Authorization", "Bearer sk-valid")
	reqFreeResp.Header.Set("Content-Type", "application/json")
	respFree, err := http.DefaultClient.Do(reqFreeResp)
	if err != nil {
		t.Fatalf("free responses request failed: %v", err)
	}
	freeBody, _ := io.ReadAll(respFree.Body)
	_ = respFree.Body.Close()
	if respFree.StatusCode != http.StatusOK || !strings.Contains(string(freeBody), "Hello from Muse Wonder") {
		t.Fatalf("free responses status=%d body=%s", respFree.StatusCode, string(freeBody))
	}
}

func TestOpenCodeSystemOneLaneEndToEnd(t *testing.T) {
	t.Parallel()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/inference/v1/systemone" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer public" || !opencode.IsCanonicalSessionID(r.Header.Get("x-opencode-session")) {
			t.Errorf("invalid systemone headers: %v", r.Header)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		if body["model"] != "jev-1.13-free" {
			t.Errorf("systemone model = %v, want jev-1.13-free", body["model"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"answers":{"proceed":true},"usage":{"input_tokens":6,"output_tokens":2}}`))
	}))
	defer upstream.Close()

	env := setupGatewayTest(t, nil, 0, true)
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ConnectionID:  opencode.VirtualConnectionID,
			ProviderID:    "opencode",
			ProviderType:  gateway.ProtocolOpenAIChat,
			UpstreamModel: "jev-1.13-free",
			Endpoint:      upstream.URL + "/inference/openai/v1",
		}}, nil
	}

	// 1. Missing required field -> 400.
	badReq, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/systemone", bytes.NewBufferString(`{"model":"opencode/jev-1.13-free"}`))
	badReq.Header.Set("Authorization", "Bearer sk-valid")
	badReq.Header.Set("Content-Type", "application/json")
	badResp, err := http.DefaultClient.Do(badReq)
	if err != nil {
		t.Fatalf("bad systemone request failed: %v", err)
	}
	_ = badResp.Body.Close()
	if badResp.StatusCode != http.StatusBadRequest {
		t.Fatalf("missing state status = %d, want 400", badResp.StatusCode)
	}

	// 2. Valid System One decision payload -> 200 passthrough.
	req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/systemone", bytes.NewBufferString(`{
		"model":"opencode/jev-1.13-free",
		"state":{"step":"review"},
		"questions":{"proceed":"boolean"}
	}`))
	req.Header.Set("Authorization", "Bearer sk-valid")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("systemone request failed: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), `"proceed":true`) {
		t.Fatalf("systemone status=%d body=%s", resp.StatusCode, string(body))
	}
}

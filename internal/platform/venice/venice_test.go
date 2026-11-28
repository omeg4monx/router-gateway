package venice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestIsVeniceKeyless(t *testing.T) {
	cases := []struct {
		name  string
		creds contracts.UpstreamCredentials
		want  bool
	}{
		{
			name:  "venice with authType none",
			creds: contracts.UpstreamCredentials{ProviderID: "venice", AuthType: "none", BaseURL: DefaultInferenceURL},
			want:  true,
		},
		{
			name:  "venice alias with empty credentials",
			creds: contracts.UpstreamCredentials{ProviderID: "venice-ai", AuthType: "api_key"},
			want:  true,
		},
		{
			name:  "venice outerface URL even with placeholder key",
			creds: contracts.UpstreamCredentials{ProviderID: "venice", AuthType: "api_key", APIKey: "dummy", BaseURL: DefaultInferenceURL},
			want:  true,
		},
		{
			name:  "paid venice API with api.venice.ai and real key",
			creds: contracts.UpstreamCredentials{ProviderID: "venice", AuthType: "api_key", APIKey: "vnc_live_123", BaseURL: "https://api.venice.ai/api/v1"},
			want:  false,
		},
		{
			name:  "other provider",
			creds: contracts.UpstreamCredentials{ProviderID: "openai", AuthType: "none"},
			want:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsVeniceKeyless(tc.creds); got != tc.want {
				t.Fatalf("IsVeniceKeyless(%+v) = %v, want %v", tc.creds, got, tc.want)
			}
		})
	}
}

func TestBuildChatPayload(t *testing.T) {
	openAIBody := []byte(`{
		"model": "venice/venice-uncensored-1-2",
		"temperature": 0.2,
		"top_p": 0.8,
		"messages": [
			{"role": "system", "content": "Be concise."},
			{"role": "user", "content": "Check weather"},
			{"role": "assistant", "content": "", "tool_calls": [{"id": "call_1", "type": "function", "function": {"name": "get_weather", "arguments": "{\"city\":\"Tokyo\"}"}}]},
			{"role": "tool", "name": "get_weather", "content": "22C sunny"}
		],
		"tools": [
			{"type": "function", "function": {"name": "get_weather", "description": "Get weather", "parameters": {"type": "object"}}}
		]
	}`)

	payloadBytes, resolvedModel, err := BuildChatPayload("venice/venice-uncensored-1-2", openAIBody)
	if err != nil {
		t.Fatalf("BuildChatPayload err = %v", err)
	}
	if resolvedModel != "venice-uncensored-1-2" {
		t.Fatalf("resolvedModel = %q, want venice-uncensored-1-2", resolvedModel)
	}

	var got struct {
		RequestID    string              `json:"requestId"`
		ModelID      string              `json:"modelId"`
		Prompt       []map[string]string `json:"prompt"`
		SystemPrompt string              `json:"systemPrompt"`
		Temperature  float64             `json:"temperature"`
		TopP         float64             `json:"topP"`
		SimpleMode   bool                `json:"simpleMode"`
	}
	if err := json.Unmarshal(payloadBytes, &got); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if got.RequestID == "" {
		t.Fatal("requestId must not be empty")
	}
	if got.ModelID != "venice-uncensored-1-2" {
		t.Fatalf("modelId = %q, want venice-uncensored-1-2", got.ModelID)
	}
	if got.Temperature != 0.2 || got.TopP != 0.8 || got.SimpleMode {
		t.Fatalf("unexpected sampling params: temp=%v topP=%v simpleMode=%v", got.Temperature, got.TopP, got.SimpleMode)
	}
	if !strings.Contains(got.SystemPrompt, "[Tools Available]") || !strings.Contains(got.SystemPrompt, "get_weather") || !strings.Contains(got.SystemPrompt, "Be concise.") {
		t.Fatalf("systemPrompt missing tools or system content: %q", got.SystemPrompt)
	}
	if len(got.Prompt) != 3 {
		t.Fatalf("len(prompt) = %d, want 3: %+v", len(got.Prompt), got.Prompt)
	}
	if got.Prompt[1]["role"] != "assistant" || !strings.Contains(got.Prompt[1]["content"], "<tool_call>") {
		t.Fatalf("assistant tool_call message not formatted: %+v", got.Prompt[1])
	}
	if got.Prompt[2]["role"] != "user" || !strings.Contains(got.Prompt[2]["content"], "[Tool Result (get_weather)]") {
		t.Fatalf("tool result message not formatted: %+v", got.Prompt[2])
	}
}

func TestAttestationManagerBatchCacheExpiryAndSingleflight(t *testing.T) {
	var mintCalls atomic.Int32
	now := time.Unix(1_800_000_000, 0)
	var mu sync.Mutex
	currentNow := now

	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != "/api/app/user/client_attestation/mint" {
				t.Errorf("unexpected mint path: %s", r.URL.Path)
			}
			if got := r.Header.Get("User-Agent"); got != "venice-middleface/"+MiddlefaceVersion {
				t.Errorf("User-Agent = %q", got)
			}
			if got := r.Header.Get("x-venice-middleface-version"); got != MiddlefaceVersion {
				t.Errorf("x-venice-middleface-version = %q", got)
			}
			if got := r.Header.Get("Origin"); got != "https://venice.ai" {
				t.Errorf("Origin = %q", got)
			}
			callNum := mintCalls.Add(1)
			mu.Lock()
			exp := currentNow.Add(5 * time.Minute).Unix()
			mu.Unlock()
			body := map[string]any{
				"expiresAt": exp,
				"tokens": []map[string]string{
					{"token": "tok-a-" + string('0'+callNum)},
					{"token": "tok-b-" + string('0'+callNum)},
					{"token": "tok-c-" + string('0'+callNum)},
					{"token": "tok-d-" + string('0'+callNum)},
					{"token": "tok-e-" + string('0'+callNum)},
				},
			}
			raw, _ := json.Marshal(body)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(raw)),
			}, nil
		}),
	}

	mgr := NewAttestationManager(func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return currentNow
	})

	// 5 concurrent requests should trigger only 1 mint call and receive 5 distinct tokens.
	var wg sync.WaitGroup
	tokens := make([]string, 5)
	errs := make([]error, 5)
	for i := range 5 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			tokens[idx], errs[idx] = mgr.GetToken(context.Background(), client, DefaultInferenceURL)
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("GetToken[%d] err = %v", i, err)
		}
		if tokens[i] == "" {
			t.Fatalf("GetToken[%d] returned empty token", i)
		}
	}
	if got := mintCalls.Load(); got != 1 {
		t.Fatalf("mintCalls after 5 concurrent gets = %d, want 1", got)
	}

	// Next call triggers mint #2 (queue exhausted).
	tok, err := mgr.GetToken(context.Background(), client, DefaultInferenceURL)
	if err != nil || !strings.HasSuffix(tok, "-2") {
		t.Fatalf("after queue drain: tok=%q err=%v", tok, err)
	}
	if got := mintCalls.Load(); got != 2 {
		t.Fatalf("mintCalls = %d, want 2", got)
	}

	// Advance clock into 60s expiry grace window -> forces preemptive mint #3.
	mu.Lock()
	currentNow = currentNow.Add(4*time.Minute + 10*time.Second)
	mu.Unlock()
	tok, err = mgr.GetToken(context.Background(), client, DefaultInferenceURL)
	if err != nil || !strings.HasSuffix(tok, "-3") {
		t.Fatalf("after expiry grace: tok=%q err=%v", tok, err)
	}
	if got := mintCalls.Load(); got != 3 {
		t.Fatalf("mintCalls = %d, want 3", got)
	}

	// Evict forces mint #4 even when tokens remain in queue.
	mgr.Evict(DefaultInferenceURL)
	tok, err = mgr.GetToken(context.Background(), client, DefaultInferenceURL)
	if err != nil || !strings.HasSuffix(tok, "-4") {
		t.Fatalf("after Evict: tok=%q err=%v", tok, err)
	}
	if got := mintCalls.Load(); got != 4 {
		t.Fatalf("mintCalls = %d, want 4", got)
	}
}

func TestParseNDJSONToOpenAI(t *testing.T) {
	ndjson := strings.Join([]string{
		`{"kind":"meta","references":[],"servingModelId":"venice-uncensored-1-2"}`,
		`{"kind":"meta","completion_id":"chatcmpl-real-123"}`,
		`{"kind":"content","content":"VEN"}`,
		`{"kind":"content","content":"ICE_OK"}`,
	}, "\n")

	out, text, err := ParseNDJSONToOpenAI([]byte(ndjson), "venice/venice-uncensored-1-2")
	if err != nil {
		t.Fatalf("ParseNDJSONToOpenAI err = %v", err)
	}
	if text != "VENICE_OK" {
		t.Fatalf("text = %q, want VENICE_OK", text)
	}

	var parsed struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("unmarshal OpenAI response: %v", err)
	}
	if parsed.ID != "chatcmpl-real-123" || parsed.Object != "chat.completion" || parsed.Model != "venice/venice-uncensored-1-2" {
		t.Fatalf("unexpected envelope: %+v", parsed)
	}
	if len(parsed.Choices) != 1 || parsed.Choices[0].Message.Content != "VENICE_OK" || parsed.Choices[0].FinishReason != "stop" {
		t.Fatalf("unexpected choices: %+v", parsed.Choices)
	}
}

func TestNewNDJSONToOpenAISSEReader(t *testing.T) {
	ndjson := strings.Join([]string{
		`{"kind":"meta","completion_id":"chatcmpl-stream-1"}`,
		`{"kind":"content","content":"Hello "}`,
		`{"kind":"content","content":"Venice"}`,
	}, "\n") + "\n"

	rc := NewNDJSONToOpenAISSEReader(context.Background(), io.NopCloser(strings.NewReader(ndjson)), "venice/venice-uncensored-1-2")
	defer rc.Close()

	raw, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll err = %v", err)
	}
	sse := string(raw)
	if !strings.Contains(sse, `"content":"Hello "`) || !strings.Contains(sse, `"content":"Venice"`) {
		t.Fatalf("missing streamed content chunks in SSE:\n%s", sse)
	}
	if !strings.Contains(sse, `"finish_reason":"stop"`) || !strings.Contains(sse, "data: [DONE]") {
		t.Fatalf("missing stop chunk or [DONE] in SSE:\n%s", sse)
	}
}

func TestSelectHTTPClientPreservesProxyTransportAndAddsCookieJar(t *testing.T) {
	proxyTransport := &http.Transport{
		Proxy: func(*http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:8888")
		},
	}
	configured := &http.Client{
		Transport: proxyTransport,
		Timeout:   12 * time.Second,
	}

	got := SelectHTTPClient(configured)
	if got.Transport != proxyTransport {
		t.Fatal("SelectHTTPClient must preserve configured *http.Transport so provider proxy pools remain active")
	}
	if got.Jar == nil {
		t.Fatal("SelectHTTPClient must attach a CookieJar when configured.Jar is nil")
	}
	if got.Timeout != configured.Timeout {
		t.Fatalf("Timeout = %v, want %v", got.Timeout, configured.Timeout)
	}
}

func TestResolveSessionUsesChromeUTLSAndIsolatesPerProxy(t *testing.T) {
	proxyA, _ := url.Parse("http://127.0.0.1:8801")
	proxyB, _ := url.Parse("http://127.0.0.1:8802")

	var seq atomic.Int32
	configured := &http.Client{
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				if seq.Add(1)%2 == 1 {
					return proxyA, nil
				}
				return proxyB, nil
			},
		},
		Timeout: 15 * time.Second,
	}

	clientA, keyA := ResolveSession(context.Background(), configured, DefaultInferenceURL)
	clientB, keyB := ResolveSession(context.Background(), configured, DefaultInferenceURL)
	t.Cleanup(func() {
		EvictSession(keyA)
		EvictSession(keyB)
	})

	if keyA != proxyA.String() || keyB != proxyB.String() {
		t.Fatalf("proxy keys = (%q, %q), want (%q, %q)", keyA, keyB, proxyA.String(), proxyB.String())
	}
	if clientA.Jar == nil || clientB.Jar == nil || clientA.Jar == clientB.Jar {
		t.Fatal("expected distinct CookieJar per proxy session")
	}
	trA, ok := clientA.Transport.(*http.Transport)
	if !ok || trA.DialTLSContext == nil {
		t.Fatal("expected ResolveSession to return Chrome uTLS transport with DialTLSContext")
	}
	httpReq, _ := http.NewRequest(http.MethodPost, "http://outerface.venice.ai/api/inference/chat", http.NoBody)
	gotProxyA, err := trA.Proxy(httpReq)
	if err != nil || gotProxyA == nil || gotProxyA.String() != proxyA.String() {
		t.Fatalf("session A http proxy = %v (err=%v), want %s", gotProxyA, err, proxyA)
	}
	httpsReq, _ := http.NewRequest(http.MethodPost, DefaultInferenceURL, http.NoBody)
	if gotHTTPSProxy, httpsErr := trA.Proxy(httpsReq); httpsErr != nil || gotHTTPSProxy != nil {
		t.Fatalf("session A https proxy = %v (err=%v), want nil so DialTLSContext handles CONNECT with uTLS", gotHTTPSProxy, httpsErr)
	}

	// Attestation buckets must also isolate per proxyKey.
	var mints atomic.Int32
	stubClient := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			n := mints.Add(1)
			raw, _ := json.Marshal(map[string]any{
				"expiresAt": time.Now().Add(5 * time.Minute).Unix(),
				"tokens":    []string{"tok-proxy-" + string('0'+n)},
			})
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(raw)),
			}, nil
		}),
	}
	mgr := NewAttestationManager(nil)
	tokA, err := mgr.GetTokenForProxy(context.Background(), stubClient, DefaultInferenceURL, keyA)
	if err != nil || tokA != "tok-proxy-1" {
		t.Fatalf("tokA = %q, err = %v", tokA, err)
	}
	tokB, err := mgr.GetTokenForProxy(context.Background(), stubClient, DefaultInferenceURL, keyB)
	if err != nil || tokB != "tok-proxy-2" {
		t.Fatalf("tokB = %q, err = %v", tokB, err)
	}
	if mints.Load() != 2 {
		t.Fatalf("mints = %d, want 2 isolated mints across 2 proxies", mints.Load())
	}
}

func TestAttestationManagerRetriesTransientNetworkError(t *testing.T) {
	var attempts atomic.Int32
	client := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			if attempts.Add(1) == 1 {
				return nil, errors.New("connection reset by peer")
			}
			raw := []byte(`{"expiresAt":1999999999,"tokens":["recovered-tok"]}`)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(raw)),
			}, nil
		}),
	}
	mgr := NewAttestationManager(nil)
	tok, err := mgr.GetToken(context.Background(), client, DefaultInferenceURL)
	if err != nil {
		t.Fatalf("GetToken err = %v, want recovery on retry", err)
	}
	if tok != "recovered-tok" || attempts.Load() != 2 {
		t.Fatalf("tok=%q attempts=%d, want recovered-tok / 2", tok, attempts.Load())
	}
}

func TestParseNDJSONToOpenAIWithToolsAndStreamingToolBridge(t *testing.T) {
	tools := ExtractRequestTools([]byte(`{
		"tools": [
			{
				"type": "function",
				"function": {
					"name": "get_weather",
					"parameters": {
						"type": "object",
						"properties": {"city": {"type": "string"}}
					}
				}
			}
		]
	}`))
	if len(tools) != 1 {
		t.Fatalf("ExtractRequestTools len = %d, want 1", len(tools))
	}

	ndjson := strings.Join([]string{
		`{"kind":"meta","completion_id":"chatcmpl-tool-1"}`,
		`{"kind":"content","content":"Checking weather now.\n"}`,
		`{"kind":"content","content":"<tool_call>{\"name\":\"get-weather\",\"arguments\":{\"city\":\"Tokyo\",\"invented\":123}}</tool_call>"}`,
		`{"kind":"content","content":"<tool_call>{\"name\":\"hallucinated_tool\",\"arguments\":{}}</tool_call>"}`,
	}, "\n") + "\n"

	out, _, err := ParseNDJSONToOpenAIWithTools([]byte(ndjson), "venice/deepseek-v4-1-flash-uncensored", tools)
	if err != nil {
		t.Fatalf("ParseNDJSONToOpenAIWithTools err = %v", err)
	}
	var parsed struct {
		Choices []struct {
			FinishReason string `json:"finish_reason"`
			Message      struct {
				Content   string     `json:"content"`
				ToolCalls []ToolCall `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(parsed.Choices) != 1 || parsed.Choices[0].FinishReason != "tool_calls" {
		t.Fatalf("finish_reason = %q, want tool_calls", parsed.Choices[0].FinishReason)
	}
	if len(parsed.Choices[0].Message.ToolCalls) != 1 {
		t.Fatalf("tool_calls = %+v, want 1 declared call", parsed.Choices[0].Message.ToolCalls)
	}
	tc := parsed.Choices[0].Message.ToolCalls[0]
	if tc.Function.Name != "get_weather" || tc.Function.Arguments != `{"city":"Tokyo"}` {
		t.Fatalf("unexpected bridged tool call: %+v", tc)
	}
	if parsed.Choices[0].Message.Content != "Checking weather now." {
		t.Fatalf("cleaned content = %q", parsed.Choices[0].Message.Content)
	}

	// Streaming tool bridge across chunk boundaries.
	streamRC := NewNDJSONToOpenAISSEReaderWithTools(context.Background(), io.NopCloser(strings.NewReader(ndjson)), "venice/deepseek-v4-1-flash-uncensored", tools)
	defer streamRC.Close()
	sseRaw, err := io.ReadAll(streamRC)
	if err != nil {
		t.Fatalf("ReadAll stream: %v", err)
	}
	sse := string(sseRaw)
	if !strings.Contains(sse, `"tool_calls"`) || !strings.Contains(sse, `"get_weather"`) || !strings.Contains(sse, `"finish_reason":"tool_calls"`) {
		t.Fatalf("stream SSE missing bridged tool_calls:\n%s", sse)
	}

	// Ordinary prose containing '<' or '```python' must flush in real time before EOF.
	pr, pw := io.Pipe()
	proseRC := NewNDJSONToOpenAISSEReaderWithTools(context.Background(), pr, "venice/deepseek-v4-1-flash-uncensored", tools)
	defer proseRC.Close()

	go func() {
		_, _ = io.WriteString(pw, `{"kind":"content","content":"if a < b then "}`+"\n")
		_, _ = io.WriteString(pw, "```python\nprint(1)\n```"+"\n")
	}()

	buf := make([]byte, 4096)
	n, readErr := proseRC.Read(buf)
	if readErr != nil || !strings.Contains(string(buf[:n]), `"content":"if a "`) {
		t.Fatalf("expected prose before '<' to stream immediately before EOF: n=%d err=%v chunk=%q", n, readErr, string(buf[:n]))
	}
	n2, readErr2 := proseRC.Read(buf)
	if readErr2 != nil || !strings.Contains(string(buf[:n2]), `"content":"\u003c b then "`) {
		t.Fatalf("expected non-tool '< b then ' to flush immediately before EOF: n=%d err=%v chunk=%q", n2, readErr2, string(buf[:n2]))
	}
	_ = pw.Close()
}

func TestVeniceLiveEndToEnd(t *testing.T) {
	if os.Getenv("VENICE_LIVE_TEST") != "1" {
		t.Skip("set VENICE_LIVE_TEST=1 to run live Venice AI verification")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	client := SharedHTTPClient()
	mgr := NewAttestationManager(nil)

	// 1. Non-streaming call (mints 5-token batch).
	t0 := time.Now()
	tok1, err := mgr.GetToken(ctx, client, DefaultInferenceURL)
	if err != nil {
		t.Fatalf("GetToken #1 failed: %v", err)
	}
	t.Logf("mint #1 ok: token_len=%d latency=%v", len(tok1), time.Since(t0).Round(time.Millisecond))

	body1, _, err := BuildChatPayload("venice/venice-uncensored-1-2", []byte(`{"messages":[{"role":"user","content":"Reply with exactly: OMEG4MON_VENICE_OK"}],"temperature":0.1}`))
	if err != nil {
		t.Fatalf("BuildChatPayload #1: %v", err)
	}
	req1, err := http.NewRequestWithContext(ctx, http.MethodPost, DefaultInferenceURL, bytes.NewReader(body1))
	if err != nil {
		t.Fatalf("NewRequest #1: %v", err)
	}
	ApplyHeaders(req1, tok1, true)
	resp1, err := client.Do(req1)
	if err != nil {
		t.Fatalf("chat #1 failed: %v", err)
	}
	raw1, _ := io.ReadAll(resp1.Body)
	_ = resp1.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("chat #1 status=%d body=%s", resp1.StatusCode, string(raw1))
	}
	openAIJSON, text1, err := ParseNDJSONToOpenAI(raw1, "venice/venice-uncensored-1-2")
	if err != nil {
		t.Fatalf("ParseNDJSONToOpenAI: %v", err)
	}
	t.Logf("non-stream ok: status=%d text=%q json_bytes=%d", resp1.StatusCode, strings.TrimSpace(text1), len(openAIJSON))
	if strings.TrimSpace(text1) == "" {
		t.Fatal("non-stream returned empty text")
	}

	// 2. Streaming call (consumes cached token from batch #1 -> 0ms mint).
	t1 := time.Now()
	tok2, err := mgr.GetToken(ctx, client, DefaultInferenceURL)
	if err != nil {
		t.Fatalf("GetToken #2 failed: %v", err)
	}
	t.Logf("mint #2 cached: token_len=%d latency=%v", len(tok2), time.Since(t1))

	body2, _, err := BuildChatPayload("venice/venice-uncensored-1-2", []byte(`{"messages":[{"role":"user","content":"Count 1 2 3"}],"temperature":0.1}`))
	if err != nil {
		t.Fatalf("BuildChatPayload #2: %v", err)
	}
	req2, err := http.NewRequestWithContext(ctx, http.MethodPost, DefaultInferenceURL, bytes.NewReader(body2))
	if err != nil {
		t.Fatalf("NewRequest #2: %v", err)
	}
	ApplyHeaders(req2, tok2, true)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("chat #2 failed: %v", err)
	}
	sseReader := NewNDJSONToOpenAISSEReader(ctx, resp2.Body, "venice/venice-uncensored-1-2")
	sseBytes, err := io.ReadAll(sseReader)
	_ = sseReader.Close()
	if err != nil {
		t.Fatalf("stream read failed: %v", err)
	}
	sseStr := string(sseBytes)
	chunks := strings.Count(sseStr, "data: {")
	hasDone := strings.Contains(sseStr, "data: [DONE]")
	t.Logf("stream ok: status=%d sse_chunks=%d has_done=%v", resp2.StatusCode, chunks, hasDone)
	if resp2.StatusCode != http.StatusOK || chunks == 0 || !hasDone {
		t.Fatalf("unexpected stream result: status=%d chunks=%d hasDone=%v", resp2.StatusCode, chunks, hasDone)
	}
}

func TestAttestationManagerLeaderCancelDoesNotAbortWaiters(t *testing.T) {
	t.Parallel()

	mgr := NewAttestationManager(nil)
	leaderStarted := make(chan struct{})
	releaseMint := make(chan struct{})
	var calls atomic.Int32

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls.Add(1)
			close(leaderStarted)
			select {
			case <-releaseMint:
				body := `{"tokens":[{"token":"tok-a"},{"token":"tok-b"},{"token":"tok-c"}],"expiresAt":4102444800}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
		}),
	}

	leaderCtx, cancelLeader := context.WithCancel(context.Background())
	defer cancelLeader()

	var (
		wg        sync.WaitGroup
		leaderErr error
		waiterTok string
		waiterErr error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, leaderErr = mgr.GetToken(leaderCtx, client, "https://venice.example/api/inference/chat")
	}()

	<-leaderStarted

	wg.Add(1)
	go func() {
		defer wg.Done()
		waiterTok, waiterErr = mgr.GetToken(context.Background(), client, "https://venice.example/api/inference/chat")
	}()

	time.Sleep(25 * time.Millisecond)
	cancelLeader()
	time.Sleep(25 * time.Millisecond)
	close(releaseMint)
	wg.Wait()

	if !errors.Is(leaderErr, context.Canceled) {
		t.Fatalf("leaderErr = %v, want context.Canceled", leaderErr)
	}
	if waiterErr != nil {
		t.Fatalf("waiterErr = %v, want nil (leader cancel must not abort waiter)", waiterErr)
	}
	if waiterTok == "" {
		t.Fatal("waiterTok is empty, want minted token")
	}
	if calls.Load() != 1 {
		t.Fatalf("mint calls = %d, want 1", calls.Load())
	}
}

func TestAttestationManagerPanicDoesNotDeadlockWaiters(t *testing.T) {
	t.Parallel()

	mgr := NewAttestationManager(nil)
	leaderStarted := make(chan struct{})
	releasePanic := make(chan struct{})
	var calls atomic.Int32

	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			if calls.Add(1) == 1 {
				close(leaderStarted)
				<-releasePanic
				panic("mint transport boom")
			}
			body := `{"tokens":[{"token":"tok-recovered"}],"expiresAt":4102444800}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		}),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() { _ = recover() }()
		_, _ = mgr.GetToken(context.Background(), client, "https://venice-panic.example/api/inference/chat")
	}()

	<-leaderStarted

	waiterDone := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		_, err := mgr.GetToken(ctx, client, "https://venice-panic.example/api/inference/chat")
		waiterDone <- err
	}()

	time.Sleep(25 * time.Millisecond)
	close(releasePanic)

	select {
	case err := <-waiterDone:
		if errors.Is(err, context.DeadlineExceeded) {
			t.Fatal("waiter deadlocked after leader panic")
		}
		if err == nil {
			t.Fatal("expected error from panicked mint")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for waiter after leader panic")
	}
	wg.Wait()

	tok, err := mgr.GetToken(context.Background(), client, "https://venice-panic.example/api/inference/chat")
	if err != nil || tok != "tok-recovered" {
		t.Fatalf("subsequent GetToken = (%q, %v), want (tok-recovered, nil)", tok, err)
	}
}

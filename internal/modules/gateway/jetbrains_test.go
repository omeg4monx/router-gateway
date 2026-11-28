package gateway_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/gateway"
)

func TestJetBrainsGatewayTransformsHeadersAndTools(t *testing.T) {
	var calls atomic.Int32
	var seenHeaders http.Header
	var seenPayload map[string]any

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		seenHeaders = r.Header.Clone()
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &seenPayload)

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("data: {\"id\":\"chatcmpl-jb\",\"object\":\"chat.completion.chunk\",\"model\":\"gpt-6-luna\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"hello from luna\"},\"finish_reason\":null}]}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer upstream.Close()

	env := setupGatewayTest(t, nil, 0, false)
	env.providers.creds["conn-jb-1"] = contracts.UpstreamCredentials{
		ConnectionID: "conn-jb-1",
		ProviderID:   "jetbrains",
		Protocol:     gateway.ProtocolOpenAIChat,
		AuthType:     "oauth",
		BaseURL:      upstream.URL,
		AccessToken:  "jb-oauth-access",
		ExtraHeaders: map[string]string{
			"Grazie-Authenticate-JWT": "grazie-jwt-token-123",
		},
	}
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ConnectionID:  "conn-jb-1",
			ProviderID:    "jetbrains",
			UpstreamModel: "gpt-6-luna",
		}}, nil
	}

	reqBody := `{
		"model": "jetbrains/gpt-6-luna",
		"stream": true,
		"max_tokens": 120,
		"tools": [{"type": "function", "function": {"name": "ping", "parameters": {"type": "object"}}}],
		"messages": [{"role": "user", "content": "hi"}]
	}`
	req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-valid")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", resp.StatusCode, string(raw))
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1", calls.Load())
	}

	if got := seenHeaders.Get("Grazie-Authenticate-JWT"); got != "grazie-jwt-token-123" {
		t.Fatalf("Grazie-Authenticate-JWT = %q, want grazie-jwt-token-123", got)
	}
	if got := seenHeaders.Get("Authorization"); got != "" {
		t.Fatalf("Authorization = %q, want empty (JetBrains uses Grazie-Authenticate-JWT)", got)
	}
	if got := seenHeaders.Get("X-Forced-License-Free"); got != "2" {
		t.Fatalf("X-Forced-License-Free = %q, want 2", got)
	}
	if got := seenHeaders.Get("Grazie-Agent"); got != `{"name":"Pulse","version":"1"}` {
		t.Fatalf("Grazie-Agent = %q, want Pulse v1 JSON", got)
	}

	if _, hasMaxTokens := seenPayload["max_tokens"]; hasMaxTokens {
		t.Fatalf("expected max_tokens to be deleted")
	}
	if mct, ok := seenPayload["max_completion_tokens"].(float64); !ok || int(mct) != 4096 {
		t.Fatalf("max_completion_tokens = %v, want 4096", seenPayload["max_completion_tokens"])
	}
	if got := seenPayload["reasoning_effort"]; got != "none" {
		t.Fatalf("reasoning_effort = %v, want none when tools are present", got)
	}
	if so, ok := seenPayload["stream_options"].(map[string]any); !ok || so["include_usage"] != true {
		t.Fatalf("stream_options = %v, want include_usage=true", seenPayload["stream_options"])
	}
}

func TestJetBrainsGatewayRejectsOversizedInput(t *testing.T) {
	var calls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	env := setupGatewayTest(t, nil, 0, false)
	env.providers.creds["conn-jb-1"] = contracts.UpstreamCredentials{
		ConnectionID: "conn-jb-1",
		ProviderID:   "jetbrains",
		Protocol:     gateway.ProtocolOpenAIChat,
		AuthType:     "oauth",
		BaseURL:      upstream.URL,
		ExtraHeaders: map[string]string{
			"Grazie-Authenticate-JWT": "grazie-jwt-token-123",
		},
	}
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ConnectionID:  "conn-jb-1",
			ProviderID:    "jetbrains",
			UpstreamModel: "gpt-6-luna",
		}}, nil
	}

	oversized := strings.Repeat("a", 5_200_005)
	payload, _ := json.Marshal(map[string]any{
		"model":    "jetbrains/gpt-6-luna",
		"messages": []map[string]any{{"role": "user", "content": oversized}},
	})
	req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-valid")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", resp.StatusCode, string(raw))
	}
	if calls.Load() != 0 {
		t.Fatalf("calls = %d, want 0 (oversized input rejected locally)", calls.Load())
	}
	if !bytes.Contains(raw, []byte("context_length_exceeded")) {
		t.Fatalf("body = %s, want context_length_exceeded", string(raw))
	}
}

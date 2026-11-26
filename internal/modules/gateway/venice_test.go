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
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/gateway"
	"github.com/omeg4monx/router-gateway/internal/platform/venice"
)

func newMockVeniceServer(t *testing.T, mintCalls, chatCalls *atomic.Int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case venice.MintPath:
			mintCalls.Add(1)
			if r.Header.Get("Authorization") != "" {
				t.Errorf("mint unexpected Authorization header: %q", r.Header.Get("Authorization"))
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"expiresAt": time.Now().Add(5 * time.Minute).Unix(),
				"tokens": []map[string]string{
					{"token": "att-tok-1"},
					{"token": "att-tok-2"},
					{"token": "att-tok-3"},
				},
			})
		case "/api/inference/chat":
			chatCalls.Add(1)
			if r.Header.Get("Authorization") != "" {
				t.Errorf("chat unexpected Authorization header: %q", r.Header.Get("Authorization"))
			}
			if tok := r.Header.Get("x-venice-client-attestation"); tok == "" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"missing attestation"}`))
				return
			}
			raw, _ := io.ReadAll(r.Body)
			var payload struct {
				ModelID string              `json:"modelId"`
				Prompt  []map[string]string `json:"prompt"`
			}
			if err := json.Unmarshal(raw, &payload); err != nil || payload.ModelID == "" || len(payload.Prompt) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error":"invalid venice payload"}`))
				return
			}
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(
				`{"kind":"meta","completion_id":"chatcmpl-venice-e2e"}` + "\n" +
					`{"kind":"content","content":"Hello "}` + "\n" +
					`{"kind":"content","content":"from Venice"}` + "\n",
			))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"unexpected path ` + r.URL.Path + `"}`))
		}
	}))
}

func TestVeniceKeylessNonStreamingAndStreamingEndToEnd(t *testing.T) {
	t.Parallel()

	var mintCalls, chatCalls atomic.Int32
	upstream := newMockVeniceServer(t, &mintCalls, &chatCalls)
	defer upstream.Close()

	env := setupGatewayTest(t, nil, 0, true)
	// Intentionally leave env.providers.creds empty — virtual noauth:venice must resolve without DB creds.
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ConnectionID:  venice.VirtualConnectionID,
			ProviderID:    "venice",
			ProviderType:  gateway.ProtocolOpenAIChat,
			UpstreamModel: "venice-uncensored-1-2",
			Endpoint:      upstream.URL + "/api/inference/chat",
		}}, nil
	}

	// 1. Non-streaming OpenAI chat completion.
	req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewBufferString(`{
		"model":"venice/venice-uncensored-1-2",
		"messages":[{"role":"user","content":"say hi"}]
	}`))
	req.Header.Set("Authorization", "Bearer sk-valid")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("non-stream request failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non-stream status = %d, want 200: %s", resp.StatusCode, string(body))
	}
	if !strings.Contains(string(body), "Hello from Venice") {
		t.Fatalf("non-stream body missing Venice text: %s", string(body))
	}

	// 2. Streaming OpenAI chat completion (reuses cached attestation batch from mint #1).
	reqStream, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewBufferString(`{
		"model":"venice/venice-uncensored-1-2",
		"stream":true,
		"messages":[{"role":"user","content":"say hi"}]
	}`))
	reqStream.Header.Set("Authorization", "Bearer sk-valid")
	reqStream.Header.Set("Content-Type", "application/json")
	respStream, err := http.DefaultClient.Do(reqStream)
	if err != nil {
		t.Fatalf("stream request failed: %v", err)
	}
	defer respStream.Body.Close()
	streamBody, _ := io.ReadAll(respStream.Body)
	if respStream.StatusCode != http.StatusOK {
		t.Fatalf("stream status = %d, want 200: %s", respStream.StatusCode, string(streamBody))
	}
	if !strings.Contains(string(streamBody), `"content":"Hello "`) || !strings.Contains(string(streamBody), "data: [DONE]") {
		t.Fatalf("stream SSE body unexpected:\n%s", string(streamBody))
	}

	// 3. Anthropic /v1/messages inbound translated to Venice keyless upstream.
	reqAnt, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/messages", bytes.NewBufferString(`{
		"model":"venice/venice-uncensored-1-2",
		"max_tokens":64,
		"messages":[{"role":"user","content":"say hi"}]
	}`))
	reqAnt.Header.Set("x-api-key", "sk-valid")
	reqAnt.Header.Set("anthropic-version", "2023-06-01")
	reqAnt.Header.Set("Content-Type", "application/json")
	respAnt, err := http.DefaultClient.Do(reqAnt)
	if err != nil {
		t.Fatalf("anthropic request failed: %v", err)
	}
	defer respAnt.Body.Close()
	antBody, _ := io.ReadAll(respAnt.Body)
	if respAnt.StatusCode != http.StatusOK {
		t.Fatalf("anthropic status = %d, want 200: %s", respAnt.StatusCode, string(antBody))
	}
	if !strings.Contains(string(antBody), "Hello from Venice") {
		t.Fatalf("anthropic translated body missing Venice text: %s", string(antBody))
	}

	if got := mintCalls.Load(); got != 1 {
		t.Fatalf("mintCalls = %d, want 1 (batched tokens reused across 3 requests)", got)
	}
	if got := chatCalls.Load(); got != 3 {
		t.Fatalf("chatCalls = %d, want 3", got)
	}
}

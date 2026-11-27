package gateway_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/gateway"
)

func TestCodeBuddyGatewayHeadersAndRequestTransform(t *testing.T) {
	cases := []struct {
		providerID string
		wantDomain string
	}{
		{providerID: "codebuddy-cn", wantDomain: "copilot.tencent.com"},
		{providerID: "codebuddy-intl", wantDomain: "www.codebuddy.ai"},
	}

	for _, tc := range cases {
		t.Run(tc.providerID, func(t *testing.T) {
			var gotHeaders http.Header
			var gotBody map[string]any

			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotHeaders = r.Header.Clone()
				raw, _ := io.ReadAll(r.Body)
				_ = json.Unmarshal(raw, &gotBody)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"chatcmpl-cb","object":"chat.completion","model":"glm-5.3","choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`))
			}))
			defer upstream.Close()

			env := setupGatewayTest(t, nil, 0, false)
			env.providers.creds["conn-cb"] = contracts.UpstreamCredentials{
				ConnectionID: "conn-cb",
				ProviderID:   tc.providerID,
				Protocol:     gateway.ProtocolOpenAIChat,
				AuthType:     "oauth",
				BaseURL:      upstream.URL + "/v2",
				AccessToken:  "cb-tok-1",
			}
			env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
				return []contracts.ResolvedTarget{{
					ConnectionID:  "conn-cb",
					ProviderID:    tc.providerID,
					UpstreamModel: "glm-5.3",
				}}, nil
			}

			reqBody := `{"model":"glm-5.3","reasoning_effort":"high","messages":[{"role":"system","content":"You are Claude Code, Anthropic's official CLI for Claude. Help carefully."},{"role":"user","content":"hello"}]}`
			req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewBufferString(reqBody))
			req.Header.Set("Authorization", "Bearer sk-valid")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Conversation-ID", "conv-explicit-99")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b, _ := io.ReadAll(resp.Body)
				t.Fatalf("status = %d, want 200: %s", resp.StatusCode, string(b))
			}
			if gotHeaders.Get("X-CodeBuddy-Request") != "1" {
				t.Errorf("X-CodeBuddy-Request = %q, want 1", gotHeaders.Get("X-CodeBuddy-Request"))
			}
			if gotHeaders.Get("X-Domain") != tc.wantDomain {
				t.Errorf("X-Domain = %q, want %q", gotHeaders.Get("X-Domain"), tc.wantDomain)
			}
			if gotHeaders.Get("X-Conversation-ID") != "conv-explicit-99" {
				t.Errorf("X-Conversation-ID = %q, want conv-explicit-99", gotHeaders.Get("X-Conversation-ID"))
			}
			if gotHeaders.Get("X-Request-ID") == "" || gotHeaders.Get("X-Conversation-Request-ID") == "" {
				t.Errorf("missing CodeBuddy request/turn headers")
			}
			if gotBody["reasoning_summary"] != "auto" {
				t.Errorf("reasoning_summary = %v, want auto", gotBody["reasoning_summary"])
			}

			msgs, ok := gotBody["messages"].([]any)
			if !ok || len(msgs) != 2 {
				t.Fatalf("messages = %#v, want 2 entries", gotBody["messages"])
			}
			sysMsg, _ := msgs[0].(map[string]any)
			sysContent, _ := sysMsg["content"].(string)
			if strings.Contains(strings.ToLower(sysContent), "claude code") {
				t.Errorf("system prompt still contains Claude Code marker: %q", sysContent)
			}
			userMsg, _ := msgs[1].(map[string]any)
			if _, isArr := userMsg["content"].([]any); !isArr {
				t.Errorf("user content should be typed array block, got %T (%#v)", userMsg["content"], userMsg["content"])
			}
		})
	}
}

package gateway_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/gateway"
)

func TestFreePIGatewayInjectsHeadersClampsTokensAndRetries426(t *testing.T) {
	gateway.SetFreePIClientVersionForTest("0.6.28")
	t.Cleanup(func() { gateway.SetFreePIClientVersionForTest("") })

	var calls atomic.Int32
	var seenVersions []string
	var seenMaxTokens []float64

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		seenVersions = append(seenVersions, r.Header.Get("x-client-version"))
		if r.Header.Get("x-session-id") == "" || r.Header.Get("x-turn-id") == "" {
			t.Errorf("missing freepi session/turn headers on call %d", n)
		}
		raw, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(raw, &payload)
		if mt, ok := payload["max_tokens"].(float64); ok {
			seenMaxTokens = append(seenMaxTokens, mt)
		}

		if n == 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUpgradeRequired)
			_, _ = w.Write([]byte(`{"error":"upgrade_required","message":"your free-pi-cli 0.6.28 is below the required minimum (0.6.31); update with: npx free-pi-cli@latest"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"chatcmpl-fp","object":"chat.completion","model":"particle","choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`))
	}))
	defer upstream.Close()

	env := setupGatewayTest(t, nil, 0, false)
	env.providers.creds["conn-fp-1"] = contracts.UpstreamCredentials{
		ConnectionID: "conn-fp-1",
		ProviderID:   "freepi",
		Protocol:     gateway.ProtocolOpenAIChat,
		AuthType:     "oauth",
		BaseURL:      upstream.URL,
		AccessToken:  "fp-tok-1",
	}
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ConnectionID:  "conn-fp-1",
			ProviderID:    "freepi",
			UpstreamModel: "particle",
		}}, nil
	}

	reqBody := `{"model":"freepi/particle","max_tokens":30,"messages":[{"role":"user","content":"hi"}]}`
	req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-valid")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 200; body=%s", resp.StatusCode, string(b))
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d, want 2 (initial 426 + retry)", calls.Load())
	}
	if len(seenVersions) != 2 || seenVersions[0] != "0.6.28" || seenVersions[1] != "0.6.31" {
		t.Fatalf("seenVersions = %v, want [0.6.28 0.6.31]", seenVersions)
	}
	if len(seenMaxTokens) == 0 || int(seenMaxTokens[0]) != 256 {
		t.Fatalf("seenMaxTokens = %v, want clamped to 256", seenMaxTokens)
	}
}

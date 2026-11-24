package providers_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/modules/providers"
)

func makeIntegrationTestJWT(t *testing.T, payload map[string]any) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal jwt payload: %v", err)
	}
	return header + "." + base64.RawURLEncoding.EncodeToString(bodyBytes) + ".sig"
}

func TestCatalogExposesXRouterOAuthFlowMetadata(t *testing.T) {
	t.Parallel()

	h := newTestHarness(t, nil, nil)
	defer func() { _ = h.mod.Stop(context.Background()) }()

	rec := h.doRequest(t, http.MethodGet, "/api/providers/catalog", "", false)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/providers/catalog status = %d: %s", rec.Code, rec.Body.String())
	}

	var catalog struct {
		Adapters []struct {
			ID    string `json:"id"`
			OAuth *struct {
				FlowType     string `json:"flow_type"`
				FixedPort    int    `json:"fixed_port"`
				CallbackPath string `json:"callback_path"`
			} `json:"oauth"`
		} `json:"adapters"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &catalog); err != nil {
		t.Fatalf("decode catalog: %v", err)
	}

	byID := map[string]struct {
		FlowType     string
		FixedPort    int
		CallbackPath string
	}{}
	for _, a := range catalog.Adapters {
		if a.OAuth != nil {
			byID[a.ID] = struct {
				FlowType     string
				FixedPort    int
				CallbackPath string
			}{
				FlowType:     a.OAuth.FlowType,
				FixedPort:    a.OAuth.FixedPort,
				CallbackPath: a.OAuth.CallbackPath,
			}
		}
	}

	if got := byID["codex"]; got.FlowType != "authorization_code_pkce" || got.FixedPort != 1455 || got.CallbackPath != "/auth/callback" {
		t.Fatalf("codex catalog oauth = %+v", got)
	}
	if got := byID["github-copilot"]; got.FlowType != "device_code" {
		t.Fatalf("github-copilot catalog oauth = %+v", got)
	}
	if got := byID["aws-quick"]; got.FlowType != "loopback_proxy" {
		t.Fatalf("aws-quick catalog oauth = %+v", got)
	}
}

func TestCodexOAuthExchangeEnrichesNameFromIDTokenJWT(t *testing.T) {
	t.Parallel()

	idToken := makeIntegrationTestJWT(t, map[string]any{
		"email": "codex.user@example.com",
		"https://api.openai.com/auth": map[string]any{
			"chatgpt_account_id": "acct_codex_77",
			"chatgpt_plan_type":  "pro",
		},
	})

	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "codex_access_tok",
			"refresh_token": "codex_refresh_tok",
			"id_token":      idToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
		})
	}))
	defer tokenSrv.Close()

	h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
		"codex": {
			TokenURL: tokenSrv.URL,
		},
	})
	defer func() { _ = h.mod.Stop(context.Background()) }()

	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/codex/start", `{}`, false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("start status = %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		State        string `json:"state"`
		FlowType     string `json:"flow_type"`
		FixedPort    int    `json:"fixed_port"`
		CallbackPath string `json:"callback_path"`
	}
	if err := json.Unmarshal(startRec.Body.Bytes(), &startResp); err != nil {
		t.Fatalf("decode startResp: %v", err)
	}
	if startResp.FlowType != "authorization_code_pkce" || startResp.FixedPort != 1455 || startResp.CallbackPath != "/auth/callback" {
		t.Fatalf("unexpected startResp metadata: %+v", startResp)
	}

	exchBody, _ := json.Marshal(map[string]any{
		"code":  "codex_code_123",
		"state": startResp.State,
	})
	exchRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/codex/exchange", string(exchBody), false)
	if exchRec.Code != http.StatusCreated {
		t.Fatalf("exchange status = %d: %s", exchRec.Code, exchRec.Body.String())
	}
	var conn providers.ConnectionSummary
	if err := json.Unmarshal(exchRec.Body.Bytes(), &conn); err != nil {
		t.Fatalf("decode conn: %v", err)
	}
	if !strings.Contains(conn.Name, "codex.user@example.com") {
		t.Fatalf("expected connection name enriched with email, got %q", conn.Name)
	}
}

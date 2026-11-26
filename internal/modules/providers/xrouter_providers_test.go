package providers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/modules/providers"
)

func TestXRouterNativeOAuthAdaptersRegisteredInCatalog(t *testing.T) {
	t.Parallel()

	h := newTestHarness(t, nil, nil)
	defer func() { _ = h.mod.Stop(context.Background()) }()

	rec := h.doRequest(t, http.MethodGet, "/api/providers/catalog", "", false)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/providers/catalog status = %d: %s", rec.Code, rec.Body.String())
	}

	var catalog struct {
		Adapters []struct {
			ID       string `json:"id"`
			AuthType string `json:"auth_type"`
			OAuth    *struct {
				FlowType           string `json:"flow_type"`
				ClientIDConfigured bool   `json:"client_id_configured"`
			} `json:"oauth"`
		} `json:"adapters"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &catalog); err != nil {
		t.Fatalf("decode catalog: %v", err)
	}

	byID := make(map[string]struct {
		flowType   string
		configured bool
	}, len(catalog.Adapters))
	for _, a := range catalog.Adapters {
		if a.OAuth != nil {
			byID[a.ID] = struct {
				flowType   string
				configured bool
			}{a.OAuth.FlowType, a.OAuth.ClientIDConfigured}
		}
	}

	wantFlows := map[string]string{
		"kimi":                "device_code",
		"grok-cli":            "device_code",
		"nous":                "device_code",
		"minimax-code":        "device_code",
		"minimax-code-global": "device_code",
		"muse":                "device_code",
		"kilocode":            "device_code",
		"iflow":               "authorization_code",
		"cline":               "authorization_code",
		"clinepass":           "authorization_code",
		"gitlab":              "authorization_code_pkce",
	}

	for id, wantFlow := range wantFlows {
		got, ok := byID[id]
		if !ok {
			t.Errorf("missing built-in OAuth adapter %q in catalog", id)
			continue
		}
		if got.flowType != wantFlow {
			t.Errorf("adapter %q flow_type = %q, want %q", id, got.flowType, wantFlow)
		}
		if !got.configured {
			t.Errorf("adapter %q client_id_configured = false, want true", id)
		}
	}
}

func TestDeclarativeRFC8628DeviceCodeFlowForMiniMaxPKCE(t *testing.T) {
	t.Parallel()

	var sawChallenge bool
	var sawVerifier bool
	var pollCount int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = r.ParseForm()
		switch r.URL.Path {
		case "/oauth2/device/code":
			if r.FormValue("code_challenge") != "" && r.FormValue("code_challenge_method") == "S256" {
				sawChallenge = true
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"device_code":               "mm_dev_123",
				"user_code":                 "MM-9999",
				"verification_uri":          "https://account.minimax.cn/verify",
				"verification_uri_complete": "https://account.minimax.cn/verify?code=MM-9999",
				"expires_in":                300,
				"interval":                  1,
			})
		case "/oauth2/token":
			if r.FormValue("code_verifier") != "" {
				sawVerifier = true
			}
			if atomic.AddInt32(&pollCount, 1) == 1 {
				_ = json.NewEncoder(w).Encode(map[string]any{"error": "authorization_pending"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token":  "mm_access_tok",
				"refresh_token": "mm_refresh_tok",
				"token_type":    "Bearer",
				"expires_in":    3600,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
		"minimax-code": {
			AuthorizeURL: srv.URL + "/oauth2/device/code",
			TokenURL:     srv.URL + "/oauth2/token",
		},
	})
	defer func() { _ = h.mod.Stop(context.Background()) }()

	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/minimax-code/start", `{}`, false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("start status = %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		State    string `json:"state"`
		UserCode string `json:"user_code"`
		FlowType string `json:"flow_type"`
	}
	_ = json.Unmarshal(startRec.Body.Bytes(), &startResp)
	if startResp.FlowType != "device_code" || startResp.UserCode != "MM-9999" {
		t.Fatalf("unexpected startResp: %+v", startResp)
	}
	if !sawChallenge {
		t.Fatalf("expected PKCE S256 code_challenge on minimax-code device code request")
	}

	// Poll 1 -> pending
	p1 := h.doRequest(t, http.MethodGet, "/api/providers/oauth/minimax-code/status?state="+startResp.State, "", false)
	if p1.Code != http.StatusOK || !strings.Contains(p1.Body.String(), `"pending"`) {
		t.Fatalf("poll1 = %d %s, want pending", p1.Code, p1.Body.String())
	}

	// Poll 2 -> done
	p2 := h.doRequest(t, http.MethodGet, "/api/providers/oauth/minimax-code/status?state="+startResp.State, "", false)
	if p2.Code != http.StatusOK || !strings.Contains(p2.Body.String(), `"done"`) {
		t.Fatalf("poll2 = %d %s, want done", p2.Code, p2.Body.String())
	}
	if !sawVerifier {
		t.Fatalf("expected PKCE code_verifier on minimax-code token poll")
	}

	// Exchange -> creates connection
	exRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/minimax-code/exchange", `{"state":"`+startResp.State+`"}`, false)
	if exRec.Code != http.StatusCreated {
		t.Fatalf("exchange status = %d: %s", exRec.Code, exRec.Body.String())
	}
}

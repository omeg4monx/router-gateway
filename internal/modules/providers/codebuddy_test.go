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

func TestCodeBuddyOAuthDeviceFlowEndToEnd(t *testing.T) {
	cases := []struct {
		providerID   string
		wantPlatform string
		wantDomain   string
	}{
		{providerID: "codebuddy-cn", wantPlatform: "CLI", wantDomain: "copilot.tencent.com"},
		{providerID: "codebuddy-intl", wantPlatform: "ide", wantDomain: "www.codebuddy.ai"},
	}

	for _, tc := range cases {
		t.Run(tc.providerID, func(t *testing.T) {
			var pollCount atomic.Int32
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/v2/plugin/auth/state":
					if r.Method != http.MethodPost {
						t.Errorf("state method = %s, want POST", r.Method)
					}
					if got := r.URL.Query().Get("platform"); got != tc.wantPlatform {
						t.Errorf("platform = %q, want %q", got, tc.wantPlatform)
					}
					if got := r.Header.Get("X-Domain"); got != tc.wantDomain {
						t.Errorf("X-Domain = %q, want %q", got, tc.wantDomain)
					}
					_ = json.NewEncoder(w).Encode(map[string]any{
						"code": 0,
						"data": map[string]any{
							"state":   "cb-state-123",
							"authUrl": "https://" + tc.wantDomain + "/login?state=cb-state-123",
						},
					})
				case "/v2/plugin/auth/token":
					if r.Method != http.MethodGet {
						t.Errorf("token poll method = %s, want GET", r.Method)
					}
					if got := r.URL.Query().Get("state"); got != "cb-state-123" {
						t.Errorf("state query = %q, want cb-state-123", got)
					}
					if pollCount.Add(1) == 1 {
						_ = json.NewEncoder(w).Encode(map[string]any{
							"code": 11217,
							"msg":  "RetryFetchToken",
						})
						return
					}
					_ = json.NewEncoder(w).Encode(map[string]any{
						"code": 0,
						"data": map[string]any{
							"accessToken":  "cb-access-xyz",
							"refreshToken": "cb-refresh-xyz",
							"tokenType":    "Bearer",
							"expiresIn":    86400,
						},
					})
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer mockServer.Close()

			h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
				tc.providerID: {
					AuthorizeURL: mockServer.URL + "/v2/plugin/auth/state",
					TokenURL:     mockServer.URL + "/v2/plugin/auth/token",
				},
			})
			defer func() { _ = h.mod.Stop(context.Background()) }()

			startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/"+tc.providerID+"/start", `{"connection_name":"My `+tc.providerID+`"}`, false)
			if startRec.Code != http.StatusOK {
				t.Fatalf("start status = %d: %s", startRec.Code, startRec.Body.String())
			}
			var startResp struct {
				State           string `json:"state"`
				AuthorizeURL    string `json:"authorize_url"`
				VerificationURI string `json:"verification_uri"`
			}
			if err := json.Unmarshal(startRec.Body.Bytes(), &startResp); err != nil {
				t.Fatalf("unmarshal start: %v", err)
			}
			if startResp.State == "" || !strings.Contains(startResp.AuthorizeURL, "cb-state-123") {
				t.Fatalf("unexpected startResp: %+v", startResp)
			}

			// First poll -> pending (code 11217)
			p1 := h.doRequest(t, http.MethodGet, "/api/providers/oauth/"+tc.providerID+"/status?state="+startResp.State, "", false)
			if p1.Code != http.StatusOK || !strings.Contains(p1.Body.String(), `"pending"`) {
				t.Fatalf("first poll = %d %s, want pending", p1.Code, p1.Body.String())
			}

			// Second poll -> done (code 0)
			p2 := h.doRequest(t, http.MethodGet, "/api/providers/oauth/"+tc.providerID+"/status?state="+startResp.State, "", false)
			if p2.Code != http.StatusOK || !strings.Contains(p2.Body.String(), `"done"`) {
				t.Fatalf("second poll = %d %s, want done", p2.Code, p2.Body.String())
			}

			// Exchange -> creates connection
			exRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/"+tc.providerID+"/exchange", `{"state":"`+startResp.State+`","connection_name":"My `+tc.providerID+`"}`, false)
			if exRec.Code != http.StatusCreated {
				t.Fatalf("exchange status = %d: %s", exRec.Code, exRec.Body.String())
			}
			var conn providers.ConnectionSummary
			if err := json.Unmarshal(exRec.Body.Bytes(), &conn); err != nil {
				t.Fatalf("unmarshal conn: %v", err)
			}
			creds, err := h.mod.ResolveCredentials(context.Background(), conn.ID)
			if err != nil {
				t.Fatalf("GetCredentials: %v", err)
			}
			if creds.AccessToken != "cb-access-xyz" {
				t.Fatalf("access_token = %q, want cb-access-xyz", creds.AccessToken)
			}
		})
	}
}

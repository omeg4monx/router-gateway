package providers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/providers"
)

func TestExtendedCatalogAdaptersAndAntigravityHeaderGuard(t *testing.T) {
	t.Parallel()

	h := newTestHarness(t, nil, nil)

	rec := h.doRequest(t, http.MethodGet, "/api/providers/catalog", "", false)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/providers/catalog status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	var catalog struct {
		Adapters []struct {
			ID             string   `json:"id"`
			DisplayName    string   `json:"display_name"`
			Category       string   `json:"category"`
			Protocol       string   `json:"protocol"`
			AuthType       string   `json:"auth_type"`
			DefaultBaseURL string   `json:"default_base_url"`
			DefaultModels  []string `json:"default_models"`
			OAuth          *struct {
				AuthorizeURL       string   `json:"authorize_url"`
				TokenURL           string   `json:"token_url"`
				Scopes             []string `json:"scopes"`
				DefaultRedirectURI string   `json:"default_redirect_uri"`
			} `json:"oauth"`
		} `json:"adapters"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &catalog); err != nil {
		t.Fatalf("unmarshal catalog: %v", err)
	}

	byID := make(map[string]int, len(catalog.Adapters))
	for i, a := range catalog.Adapters {
		byID[a.ID] = i
	}

	type wantAdapter struct {
		protocol           string
		authType           string
		wantOAuth          bool
		wantRedirectNoPath bool
	}
	wants := map[string]wantAdapter{
		"aws-quick": {protocol: "openai-chat", authType: "oauth", wantOAuth: true, wantRedirectNoPath: true},
		"venice":    {protocol: "openai-chat", authType: "api_key", wantOAuth: false},
		"freepi":    {protocol: "openai-chat", authType: "oauth", wantOAuth: true},
	}

	for id, want := range wants {
		idx, ok := byID[id]
		if !ok {
			t.Fatalf("missing built-in adapter %q in catalog", id)
		}
		got := catalog.Adapters[idx]
		if got.Protocol != want.protocol {
			t.Errorf("adapter %q protocol = %q, want %q", id, got.Protocol, want.protocol)
		}
		if got.AuthType != want.authType {
			t.Errorf("adapter %q auth_type = %q, want %q", id, got.AuthType, want.authType)
		}
		if got.DefaultBaseURL == "" {
			t.Errorf("adapter %q default_base_url is empty", id)
		}
		if len(got.DefaultModels) == 0 {
			t.Errorf("adapter %q default_models is empty", id)
		}
		if want.wantOAuth {
			if got.OAuth == nil || got.OAuth.AuthorizeURL == "" || got.OAuth.TokenURL == "" {
				t.Errorf("adapter %q oauth config missing or incomplete: %+v", id, got.OAuth)
			}
			if want.wantRedirectNoPath && strings.HasSuffix(got.OAuth.DefaultRedirectURI, "/") {
				t.Errorf("adapter %q default_redirect_uri = %q, must be bare loopback origin without trailing slash", id, got.OAuth.DefaultRedirectURI)
			}
		} else if got.OAuth != nil {
			t.Errorf("adapter %q should not expose oauth config", id)
		}
	}

	// Verify Antigravity test request strips consumer/legacy x-goog-user-project
	// values while forwarding real provisioned GCP project IDs.
	cases := []struct {
		name        string
		projectID   string
		wantProject string
	}{
		{name: "strip_aicode_consumers", projectID: "aicode-consumers", wantProject: ""},
		{name: "strip_legacy_fake_pid", projectID: "useful-fuze-1234abcd", wantProject: ""},
		{name: "strip_cloudaicompanion_fake", projectID: "cloudaicompanion-proj-99aa", wantProject: ""},
		{name: "keep_real_gcp_project", projectID: "my-real-gcp-project-42", wantProject: "my-real-gcp-project-42"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var sawProject string
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sawProject = r.Header.Get("x-goog-user-project")
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"models":[]}`))
			}))
			defer upstream.Close()

			payload, _ := json.Marshal(map[string]any{
				"provider_id":  "antigravity",
				"name":         "AG " + tc.name,
				"base_url":     upstream.URL,
				"access_token": "tok-ag",
				"headers": map[string]string{
					"x-goog-user-project": tc.projectID,
				},
			})
			createRec := h.doRequest(t, http.MethodPost, "/api/providers/connections", string(payload), false)
			if createRec.Code != http.StatusCreated {
				t.Fatalf("create antigravity connection status = %d: %s", createRec.Code, createRec.Body.String())
			}
			var conn providers.ConnectionSummary
			if err := json.Unmarshal(createRec.Body.Bytes(), &conn); err != nil {
				t.Fatalf("unmarshal connection: %v", err)
			}

			testRec := h.doRequest(t, http.MethodPost, "/api/providers/connections/"+conn.ID+"/test", "", false)
			if testRec.Code != http.StatusOK {
				t.Fatalf("test antigravity connection status = %d: %s", testRec.Code, testRec.Body.String())
			}
			if sawProject != tc.wantProject {
				t.Errorf("x-goog-user-project = %q, want %q", sawProject, tc.wantProject)
			}
		})
	}
}

func TestOAuthRefreshCooldownGuardPreventsRetryStorm(t *testing.T) {
	t.Parallel()

	var tokenHits atomic.Int32
	var shouldSucceed atomic.Bool

	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		tokenHits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if !shouldSucceed.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"server_error"}`))
			return
		}
		_, _ = w.Write([]byte(`{"access_token":"fresh-tok-2","refresh_token":"rt-2","expires_in":3600}`))
	}))
	defer tokenSrv.Close()

	h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
		"claude": {
			TokenURL: tokenSrv.URL,
			ClientID: "test-client",
		},
	})

	createBody, _ := json.Marshal(map[string]any{
		"provider_id":   "claude",
		"name":          "Claude Cooldown Test",
		"access_token":  "expired-tok",
		"refresh_token": "rt-1",
	})
	createRec := h.doRequest(t, http.MethodPost, "/api/providers/connections", string(createBody), false)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create connection status = %d: %s", createRec.Code, createRec.Body.String())
	}
	var conn providers.ConnectionSummary
	if err := json.Unmarshal(createRec.Body.Bytes(), &conn); err != nil {
		t.Fatalf("unmarshal connection: %v", err)
	}

	if _, err := h.mod.RefreshOAuthToken(context.Background(), conn.ID); err == nil {
		t.Fatal("first refresh expected upstream 500 error, got nil")
	}
	if got := tokenHits.Load(); got != 1 {
		t.Fatalf("upstream token hits after first failure = %d, want 1", got)
	}

	// Second immediate call must be blocked by failure cooldown without hitting upstream.
	if _, err := h.mod.RefreshOAuthToken(context.Background(), conn.ID); err == nil {
		t.Fatal("second refresh during cooldown expected error, got nil")
	}
	if got := tokenHits.Load(); got != 1 {
		t.Fatalf("upstream token hits during cooldown = %d, want 1 (cooldown should short-circuit)", got)
	}

	// Updating credentials via PATCH clears the cooldown so refresh can proceed.
	shouldSucceed.Store(true)
	patchBody, _ := json.Marshal(map[string]any{
		"refresh_token": "rt-updated",
	})
	patchRec := h.doRequest(t, http.MethodPatch, "/api/providers/connections/"+conn.ID, string(patchBody), false)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("patch connection status = %d: %s", patchRec.Code, patchRec.Body.String())
	}

	creds, err := h.mod.RefreshOAuthToken(context.Background(), conn.ID)
	if err != nil {
		t.Fatalf("refresh after credential update failed: %v", err)
	}
	if creds.AccessToken != "fresh-tok-2" {
		t.Errorf("refreshed access_token = %q, want fresh-tok-2", creds.AccessToken)
	}
	if got := tokenHits.Load(); got != 2 {
		t.Errorf("upstream token hits after cooldown clear = %d, want 2", got)
	}

	// Connection without a refresh_token returns ErrOAuthRefreshUnsupported and arms cooldown.
	noRTBody, _ := json.Marshal(map[string]any{
		"provider_id":  "claude",
		"name":         "Claude No RT",
		"access_token": "only-access-tok",
	})
	noRTRec := h.doRequest(t, http.MethodPost, "/api/providers/connections", string(noRTBody), false)
	if noRTRec.Code != http.StatusCreated {
		t.Fatalf("create no-RT connection status = %d: %s", noRTRec.Code, noRTRec.Body.String())
	}
	var noRTConn providers.ConnectionSummary
	if err := json.Unmarshal(noRTRec.Body.Bytes(), &noRTConn); err != nil {
		t.Fatalf("unmarshal no-RT connection: %v", err)
	}

	for i := 0; i < 2; i++ {
		_, err := h.mod.RefreshOAuthToken(context.Background(), noRTConn.ID)
		if !errors.Is(err, contracts.ErrOAuthRefreshUnsupported) {
			t.Fatalf("call %d for no-RT connection err = %v, want ErrOAuthRefreshUnsupported", i+1, err)
		}
	}
}

package providers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/modules/providers"
)

func TestAllBuiltinOAuthAdaptersHaveClientIDConfiguredAndStartWithoutManualClientID(t *testing.T) {
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
				ClientIDConfigured bool `json:"client_id_configured"`
			} `json:"oauth"`
		} `json:"adapters"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &catalog); err != nil {
		t.Fatalf("decode catalog: %v", err)
	}

	for _, a := range catalog.Adapters {
		if a.AuthType != "oauth" {
			continue
		}
		if a.OAuth == nil || !a.OAuth.ClientIDConfigured {
			t.Errorf("oauth adapter %q has client_id_configured = false, want true", a.ID)
		}
	}

	// All 5 PKCE providers must start cleanly with an empty JSON body.
	for _, pid := range []string{"claude", "codex", "gemini-oauth", "antigravity", "qwen"} {
		startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/"+pid+"/start", `{}`, false)
		if startRec.Code != http.StatusOK {
			t.Errorf("POST /api/providers/oauth/%s/start status = %d: %s", pid, startRec.Code, startRec.Body.String())
			continue
		}
		var resp struct {
			State        string `json:"state"`
			AuthorizeURL string `json:"authorize_url"`
		}
		if err := json.Unmarshal(startRec.Body.Bytes(), &resp); err != nil {
			t.Errorf("decode %s startResp: %v", pid, err)
			continue
		}
		parsed, err := url.Parse(resp.AuthorizeURL)
		if err != nil || parsed.Query().Get("client_id") == "" {
			t.Errorf("provider %q start authorize_url missing client_id: %s", pid, resp.AuthorizeURL)
		}
	}
}

func TestGitHubCopilotDeviceCodeFlow(t *testing.T) {
	t.Parallel()

	var pollCount atomic.Int32
	mockGH := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/login/device/code":
			_, _ = w.Write([]byte(`{"device_code":"gh-dev-1","user_code":"ABCD-1234","verification_uri":"https://github.com/login/device","expires_in":900,"interval":5}`))
		case "/login/oauth/access_token":
			n := pollCount.Add(1)
			if n == 1 {
				_, _ = w.Write([]byte(`{"error":"authorization_pending"}`))
				return
			}
			_, _ = w.Write([]byte(`{"access_token":"ghu_test_token","token_type":"bearer","scope":"read:user"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockGH.Close()

	h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
		"github-copilot": {
			AuthorizeURL: mockGH.URL + "/login/device/code",
			TokenURL:     mockGH.URL + "/login/oauth/access_token",
		},
	})

	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/github-copilot/start", `{"connection_name":"My Copilot"}`, false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("start github-copilot status = %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		State           string `json:"state"`
		FlowType        string `json:"flow_type"`
		UserCode        string `json:"user_code"`
		VerificationURI string `json:"verification_uri"`
		Interval        int    `json:"interval"`
	}
	_ = json.Unmarshal(startRec.Body.Bytes(), &startResp)
	if startResp.FlowType != "device_code" || startResp.UserCode != "ABCD-1234" || startResp.VerificationURI != "https://github.com/login/device" {
		t.Fatalf("unexpected github-copilot startResp: %+v", startResp)
	}

	// Poll 1 -> pending
	s1 := h.doRequest(t, http.MethodGet, "/api/providers/oauth/github-copilot/status?state="+url.QueryEscape(startResp.State), "", false)
	if s1.Code != http.StatusOK || !containsJSONStatus(s1.Body.Bytes(), "pending") {
		t.Fatalf("poll 1 = %d %s, want pending", s1.Code, s1.Body.String())
	}

	// Poll 2 -> done
	s2 := h.doRequest(t, http.MethodGet, "/api/providers/oauth/github-copilot/status?state="+url.QueryEscape(startResp.State), "", false)
	if s2.Code != http.StatusOK || !containsJSONStatus(s2.Body.Bytes(), "done") {
		t.Fatalf("poll 2 = %d %s, want done", s2.Code, s2.Body.String())
	}

	// Exchange -> 201 Created
	exRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/github-copilot/exchange", `{"state":"`+startResp.State+`"}`, false)
	if exRec.Code != http.StatusCreated {
		t.Fatalf("exchange status = %d: %s", exRec.Code, exRec.Body.String())
	}
	var conn providers.ConnectionSummary
	_ = json.Unmarshal(exRec.Body.Bytes(), &conn)
	creds, err := h.mod.ResolveCredentials(context.Background(), conn.ID)
	if err != nil || creds.AccessToken != "ghu_test_token" {
		t.Fatalf("resolved creds = %+v, err = %v", creds, err)
	}
}

func TestKiroDeviceCodeFlowAndRefresh(t *testing.T) {
	t.Parallel()

	mockAWS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/client/register":
			_, _ = w.Write([]byte(`{"clientId":"kiro-cid","clientSecret":"kiro-csec"}`))
		case "/device_authorization":
			_, _ = w.Write([]byte(`{"deviceCode":"kiro-dev","userCode":"KIRO-9999","verificationUri":"https://device.sso.us-east-1.amazonaws.com/","verificationUriComplete":"https://device.sso.us-east-1.amazonaws.com/?user_code=KIRO-9999","expiresIn":600,"interval":5}`))
		case "/token":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["grantType"] == "refresh_token" {
				_, _ = w.Write([]byte(`{"accessToken":"kiro-at-2","refreshToken":"kiro-rt-2","expiresIn":3600}`))
				return
			}
			_, _ = w.Write([]byte(`{"accessToken":"kiro-at-1","refreshToken":"kiro-rt-1","expiresIn":3600}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockAWS.Close()

	h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
		"kiro": {
			TokenURL: mockAWS.URL + "/token",
		},
	})

	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/kiro/start", `{}`, false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("start kiro status = %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		State    string `json:"state"`
		FlowType string `json:"flow_type"`
		UserCode string `json:"user_code"`
	}
	_ = json.Unmarshal(startRec.Body.Bytes(), &startResp)
	if startResp.FlowType != "device_code" || startResp.UserCode != "KIRO-9999" {
		t.Fatalf("unexpected kiro startResp: %+v", startResp)
	}

	s := h.doRequest(t, http.MethodGet, "/api/providers/oauth/kiro/status?state="+url.QueryEscape(startResp.State), "", false)
	if s.Code != http.StatusOK || !containsJSONStatus(s.Body.Bytes(), "done") {
		t.Fatalf("kiro poll = %d %s, want done", s.Code, s.Body.String())
	}

	exRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/kiro/exchange", `{"state":"`+startResp.State+`"}`, false)
	if exRec.Code != http.StatusCreated {
		t.Fatalf("kiro exchange status = %d: %s", exRec.Code, exRec.Body.String())
	}
	var conn providers.ConnectionSummary
	_ = json.Unmarshal(exRec.Body.Bytes(), &conn)

	refreshed, err := h.mod.RefreshOAuthToken(context.Background(), conn.ID)
	if err != nil || refreshed.AccessToken != "kiro-at-2" {
		t.Fatalf("kiro refresh = %+v, err = %v", refreshed, err)
	}
}

func TestFreePIDeviceCodeFlow(t *testing.T) {
	t.Parallel()

	var consentCalled atomic.Bool
	mockFP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/auth/github/device":
			_, _ = w.Write([]byte(`{"session_id":"fp-sess-1","user_code":"FREE-4242","verification_uri":"https://github.com/login/device","expires_in":899,"interval":5}`))
		case "/auth/consent":
			consentCalled.Store(true)
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/auth/token":
			_, _ = w.Write([]byte(`{"token":"fp-tok-1","handle":"octocat"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockFP.Close()

	h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
		"freepi": {
			TokenURL: mockFP.URL + "/auth/token",
		},
	})

	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/freepi/start", `{}`, false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("start freepi status = %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		State    string `json:"state"`
		FlowType string `json:"flow_type"`
		UserCode string `json:"user_code"`
	}
	_ = json.Unmarshal(startRec.Body.Bytes(), &startResp)
	if startResp.FlowType != "device_code" || startResp.UserCode != "FREE-4242" || !consentCalled.Load() {
		t.Fatalf("unexpected freepi startResp: %+v, consent=%v", startResp, consentCalled.Load())
	}

	s := h.doRequest(t, http.MethodGet, "/api/providers/oauth/freepi/status?state="+url.QueryEscape(startResp.State), "", false)
	if s.Code != http.StatusOK || !containsJSONStatus(s.Body.Bytes(), "done") {
		t.Fatalf("freepi poll = %d %s, want done", s.Code, s.Body.String())
	}

	exRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/freepi/exchange", `{"state":"`+startResp.State+`"}`, false)
	if exRec.Code != http.StatusCreated {
		t.Fatalf("freepi exchange status = %d: %s", exRec.Code, exRec.Body.String())
	}
	var conn providers.ConnectionSummary
	_ = json.Unmarshal(exRec.Body.Bytes(), &conn)
	creds, err := h.mod.ResolveCredentials(context.Background(), conn.ID)
	if err != nil || creds.AccessToken != "fp-tok-1" {
		t.Fatalf("freepi creds = %+v, err = %v", creds, err)
	}
}

func containsJSONStatus(raw []byte, want string) bool {
	var payload struct {
		Status string `json:"status"`
	}
	return json.Unmarshal(raw, &payload) == nil && payload.Status == want
}

func TestFreePIConsentFailureFailsStart(t *testing.T) {
	t.Parallel()

	t.Run("4xx_fails_fast_without_retry", func(t *testing.T) {
		var consentHits atomic.Int32
		mockFP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/auth/github/device":
				_, _ = w.Write([]byte(`{"session_id":"fp-4xx","user_code":"FREE-403","verification_uri":"https://github.com/login/device"}`))
			case "/auth/consent":
				consentHits.Add(1)
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"code":"invalid_session","message":"bad session"}`))
			default:
				http.NotFound(w, r)
			}
		}))
		defer mockFP.Close()

		h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
			"freepi": {TokenURL: mockFP.URL + "/auth/token"},
		})

		startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/freepi/start", `{}`, false)
		if startRec.Code == http.StatusOK {
			t.Fatalf("expected start to fail when consent returns 400, got 200: %s", startRec.Body.String())
		}
		if got := consentHits.Load(); got != 1 {
			t.Fatalf("consent attempts on 4xx = %d, want 1", got)
		}
		if !strings.Contains(startRec.Body.String(), "consent") {
			t.Fatalf("expected error body to mention consent, got: %s", startRec.Body.String())
		}
	})

	t.Run("5xx_retries_once_and_succeeds", func(t *testing.T) {
		var consentHits atomic.Int32
		mockFP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/auth/github/device":
				_, _ = w.Write([]byte(`{"session_id":"fp-5xx-ok","user_code":"FREE-503","verification_uri":"https://github.com/login/device"}`))
			case "/auth/consent":
				if consentHits.Add(1) == 1 {
					w.WriteHeader(http.StatusServiceUnavailable)
					_, _ = w.Write([]byte(`{"code":"unavailable","message":"try again"}`))
					return
				}
				_, _ = w.Write([]byte(`{"ok":true}`))
			default:
				http.NotFound(w, r)
			}
		}))
		defer mockFP.Close()

		h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
			"freepi": {TokenURL: mockFP.URL + "/auth/token"},
		})

		startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/freepi/start", `{}`, false)
		if startRec.Code != http.StatusOK {
			t.Fatalf("expected start to succeed after 5xx retry, got %d: %s", startRec.Code, startRec.Body.String())
		}
		if got := consentHits.Load(); got != 2 {
			t.Fatalf("consent attempts = %d, want 2", got)
		}
	})
}

func TestFreePIPollSurfacesCodeMessageAndRetryAfter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		tokenStatus  int
		tokenBody    string
		wantStatus   string
		wantInterval int
		wantErrPart  string
	}{
		{
			name:         "pending_honors_retry_after",
			tokenStatus:  http.StatusOK,
			tokenBody:    `{"status":"pending","retry_after":12}`,
			wantStatus:   "pending",
			wantInterval: 12,
		},
		{
			name:         "pending_caps_retry_after_at_60s",
			tokenStatus:  http.StatusOK,
			tokenBody:    `{"status":"pending","retry_after":120}`,
			wantStatus:   "pending",
			wantInterval: 60,
		},
		{
			name:        "surfaces_403_consent_required",
			tokenStatus: http.StatusForbidden,
			tokenBody:   `{"code":"consent_required","message":"Consent record missing"}`,
			wantStatus:  "error",
			wantErrPart: "consent_required",
		},
		{
			name:        "surfaces_404_not_found",
			tokenStatus: http.StatusNotFound,
			tokenBody:   `{"code":"not_found","message":"unknown session"}`,
			wantStatus:  "error",
			wantErrPart: "not_found",
		},
		{
			name:        "surfaces_429_rate_limited",
			tokenStatus: http.StatusTooManyRequests,
			tokenBody:   `{"code":"rate_limited","message":"too many polls"}`,
			wantStatus:  "error",
			wantErrPart: "rate_limited",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockFP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/auth/github/device":
					_, _ = w.Write([]byte(`{"session_id":"fp-poll-1","user_code":"FREE-POLL","verification_uri":"https://github.com/login/device"}`))
				case "/auth/consent":
					_, _ = w.Write([]byte(`{"ok":true}`))
				case "/auth/token":
					w.WriteHeader(tc.tokenStatus)
					_, _ = w.Write([]byte(tc.tokenBody))
				default:
					http.NotFound(w, r)
				}
			}))
			defer mockFP.Close()

			h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
				"freepi": {TokenURL: mockFP.URL + "/auth/token"},
			})

			startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/freepi/start", `{}`, false)
			if startRec.Code != http.StatusOK {
				t.Fatalf("start status = %d: %s", startRec.Code, startRec.Body.String())
			}
			var startResp struct {
				State string `json:"state"`
			}
			_ = json.Unmarshal(startRec.Body.Bytes(), &startResp)

			pollRec := h.doRequest(t, http.MethodGet, "/api/providers/oauth/freepi/status?state="+url.QueryEscape(startResp.State), "", false)
			if pollRec.Code != http.StatusOK {
				t.Fatalf("poll HTTP status = %d: %s", pollRec.Code, pollRec.Body.String())
			}
			var pollResp struct {
				Status   string `json:"status"`
				Interval int    `json:"interval"`
				Error    string `json:"error"`
			}
			if err := json.Unmarshal(pollRec.Body.Bytes(), &pollResp); err != nil {
				t.Fatalf("unmarshal poll response: %v", err)
			}
			if pollResp.Status != tc.wantStatus {
				t.Fatalf("poll status = %q, want %q (body=%s)", pollResp.Status, tc.wantStatus, pollRec.Body.String())
			}
			if tc.wantInterval > 0 && pollResp.Interval != tc.wantInterval {
				t.Fatalf("poll interval = %d, want %d", pollResp.Interval, tc.wantInterval)
			}
			if tc.wantErrPart != "" && !strings.Contains(pollResp.Error, tc.wantErrPart) {
				t.Fatalf("poll error = %q, want substring %q", pollResp.Error, tc.wantErrPart)
			}
		})
	}
}

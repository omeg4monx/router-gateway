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

func TestAwsQuickTwoHopLoopbackOAuthAndRefresh(t *testing.T) {
	t.Parallel()

	var (
		registerCalls atomic.Int32
		tokenCalls    atomic.Int32
		refreshCalls  atomic.Int32
		userInfoCalls atomic.Int32
	)

	mockAWS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/client/register":
			registerCalls.Add(1)
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "bad json", http.StatusBadRequest)
				return
			}
			if body["clientName"] != "AmazonQuick-Desktop" || body["clientType"] != "public" {
				http.Error(w, "unexpected client register payload", http.StatusBadRequest)
				return
			}
			_, _ = w.Write([]byte(`{"clientId":"dyn-client-1","clientSecret":"dyn-secret-1"}`))
		case "/token":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "bad json", http.StatusBadRequest)
				return
			}
			switch body["grantType"] {
			case "authorization_code":
				tokenCalls.Add(1)
				if body["clientId"] != "dyn-client-1" || body["clientSecret"] != "dyn-secret-1" || body["code"] != "quick-auth-code" {
					http.Error(w, "invalid token exchange body", http.StatusUnauthorized)
					return
				}
				if body["codeVerifier"] == "" || body["redirectUri"] == "" {
					http.Error(w, "missing pkce/redirect", http.StatusBadRequest)
					return
				}
				_, _ = w.Write([]byte(`{"accessToken":"opaque-access-1","idToken":"quick-id-token-1","refreshToken":"quick-rt-1","expiresIn":3600}`))
			case "refresh_token":
				refreshCalls.Add(1)
				if body["clientId"] != "dyn-client-1" || body["clientSecret"] != "dyn-secret-1" || body["refreshToken"] != "quick-rt-1" {
					http.Error(w, "invalid refresh body", http.StatusUnauthorized)
					return
				}
				_, _ = w.Write([]byte(`{"accessToken":"opaque-access-2","idToken":"quick-id-token-2","refreshToken":"quick-rt-2","expiresIn":3600}`))
			default:
				http.Error(w, "unsupported grantType", http.StatusBadRequest)
			}
		case "/integration/quick-work/user-info":
			userInfoCalls.Add(1)
			if r.Header.Get("Authorization") != "Bearer opaque-access-1" {
				http.Error(w, "missing bearer", http.StatusUnauthorized)
				return
			}
			_, _ = w.Write([]byte(`{"email":"alice@example.com","displayName":"Alice Quick","userArn":"arn:aws:iam::123:user/alice"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockAWS.Close()

	h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
		"aws-quick": {
			TokenURL: mockAWS.URL + "/token",
			ExtraTokenParams: map[string]string{
				"user_info_url": mockAWS.URL + "/integration/quick-work/user-info",
			},
		},
	})
	defer func() { _ = h.mod.Stop(context.Background()) }()

	// 1. Start without client_id — must succeed and return quick.aws.com/sn discovery URL.
	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/aws-quick/start", `{}`, false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("POST /oauth/aws-quick/start status = %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		State        string `json:"state"`
		AuthorizeURL string `json:"authorize_url"`
		RedirectURI  string `json:"redirect_uri"`
	}
	if err := json.Unmarshal(startRec.Body.Bytes(), &startResp); err != nil {
		t.Fatalf("decode startResp: %v", err)
	}
	if startResp.State == "" || !strings.HasPrefix(startResp.AuthorizeURL, "https://quick.aws.com/sn?qwSurface=") {
		t.Fatalf("unexpected aws-quick startResp: %+v", startResp)
	}
	if !strings.HasPrefix(startResp.RedirectURI, "http://127.0.0.1:") {
		t.Fatalf("unexpected redirect_uri: %q", startResp.RedirectURI)
	}

	// Extract nonce from qwSurface=desktop:<nonce>:<port>:<version>.
	parsedStart, err := url.Parse(startResp.AuthorizeURL)
	if err != nil {
		t.Fatalf("parse authorize_url: %v", err)
	}
	surfaceParts := strings.Split(parsedStart.Query().Get("qwSurface"), ":")
	if len(surfaceParts) != 4 || surfaceParts[0] != "desktop" {
		t.Fatalf("unexpected qwSurface: %q", parsedStart.Query().Get("qwSurface"))
	}
	nonce := surfaceParts[1]

	// 2. Hop 1 loopback callback: /callback?nonce=...&issuerUrl=...&tenantUrl=...&idcRegion=us-east-1
	noRedirectClient := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	hop1URL := startResp.RedirectURI + "/callback?" + url.Values{
		"nonce":     {nonce},
		"issuerUrl": {"https://identitycenter.amazonaws.com/ssoins-test"},
		"tenantUrl": {"https://tenant.quick.aws.com"},
		"idcRegion": {"us-east-1"},
	}.Encode()
	hop1Resp, err := noRedirectClient.Get(hop1URL)
	if err != nil {
		t.Fatalf("hop-1 GET failed: %v", err)
	}
	_ = hop1Resp.Body.Close()
	if hop1Resp.StatusCode != http.StatusFound {
		t.Fatalf("hop-1 status = %d, want 302", hop1Resp.StatusCode)
	}
	loc := hop1Resp.Header.Get("Location")
	parsedLoc, err := url.Parse(loc)
	if err != nil {
		t.Fatalf("parse hop-1 Location: %v", err)
	}
	if parsedLoc.Query().Get("client_id") != "dyn-client-1" || parsedLoc.Query().Get("code_challenge_method") != "S256" || parsedLoc.Query().Get("state") != startResp.State {
		t.Fatalf("unexpected hop-1 authorize redirect: %s", loc)
	}

	// 3. Hop 2 loopback callback: /?code=quick-auth-code&state=<state>
	hop2URL := startResp.RedirectURI + "/?" + url.Values{
		"code":  {"quick-auth-code"},
		"state": {startResp.State},
	}.Encode()
	hop2Resp, err := noRedirectClient.Get(hop2URL)
	if err != nil {
		t.Fatalf("hop-2 GET failed: %v", err)
	}
	_ = hop2Resp.Body.Close()
	if hop2Resp.StatusCode != http.StatusOK {
		t.Fatalf("hop-2 status = %d, want 200", hop2Resp.StatusCode)
	}

	// 4. Poll status -> done.
	statusRec := h.doRequest(t, http.MethodGet, "/api/providers/oauth/aws-quick/status?state="+url.QueryEscape(startResp.State), "", false)
	if statusRec.Code != http.StatusOK {
		t.Fatalf("GET status = %d: %s", statusRec.Code, statusRec.Body.String())
	}
	var statusResp struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(statusRec.Body.Bytes(), &statusResp); err != nil {
		t.Fatalf("decode statusResp: %v", err)
	}
	if statusResp.Status != "done" {
		t.Fatalf("status = %q, want done", statusResp.Status)
	}

	// 5. Exchange settles connection in DB.
	exRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/aws-quick/exchange", `{"state":"`+startResp.State+`"}`, false)
	if exRec.Code != http.StatusCreated {
		t.Fatalf("POST exchange status = %d: %s", exRec.Code, exRec.Body.String())
	}
	var conn providers.ConnectionSummary
	if err := json.Unmarshal(exRec.Body.Bytes(), &conn); err != nil {
		t.Fatalf("decode conn: %v", err)
	}
	if conn.ProviderID != "aws-quick" || conn.BaseURL != "https://tenant.quick.aws.com" {
		t.Fatalf("unexpected conn: %+v", conn)
	}

	creds, err := h.mod.ResolveCredentials(context.Background(), conn.ID)
	if err != nil {
		t.Fatalf("ResolveCredentials: %v", err)
	}
	if creds.AccessToken != "quick-id-token-1" {
		t.Fatalf("AccessToken = %q, want quick-id-token-1", creds.AccessToken)
	}

	// 6. Refresh token uses camelCase JSON and returns rotated id_token.
	refreshed, err := h.mod.RefreshOAuthToken(context.Background(), conn.ID)
	if err != nil {
		t.Fatalf("RefreshOAuthToken: %v", err)
	}
	if refreshed.AccessToken != "quick-id-token-2" {
		t.Fatalf("refreshed AccessToken = %q, want quick-id-token-2", refreshed.AccessToken)
	}
	if registerCalls.Load() != 1 || tokenCalls.Load() != 1 || userInfoCalls.Load() != 1 || refreshCalls.Load() != 1 {
		t.Fatalf("call counts: register=%d token=%d userInfo=%d refresh=%d",
			registerCalls.Load(), tokenCalls.Load(), userInfoCalls.Load(), refreshCalls.Load())
	}
}

func TestAwsQuickTwoStepPasteBackFlow(t *testing.T) {
	t.Parallel()

	mockAWS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/client/register":
			_, _ = w.Write([]byte(`{"clientId":"paste-client","clientSecret":"paste-secret"}`))
		case "/token":
			_, _ = w.Write([]byte(`{"accessToken":"paste-access","idToken":"paste-id-tok","refreshToken":"paste-rt","expiresIn":3600}`))
		case "/integration/quick-work/user-info":
			_, _ = w.Write([]byte(`{"email":"remote@example.com","displayName":"Remote User"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockAWS.Close()

	h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
		"aws-quick": {
			TokenURL: mockAWS.URL + "/token",
			ExtraTokenParams: map[string]string{
				"user_info_url": mockAWS.URL + "/integration/quick-work/user-info",
			},
		},
	})
	defer func() { _ = h.mod.Stop(context.Background()) }()

	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/aws-quick/start", `{"connection_name":"Remote Quick"}`, false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("start status = %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		State        string `json:"state"`
		AuthorizeURL string `json:"authorize_url"`
		RedirectURI  string `json:"redirect_uri"`
	}
	_ = json.Unmarshal(startRec.Body.Bytes(), &startResp)
	parsedStart, _ := url.Parse(startResp.AuthorizeURL)
	nonce := strings.Split(parsedStart.Query().Get("qwSurface"), ":")[1]

	// Paste Hop-1 discovery URL.
	hop1Pasted := startResp.RedirectURI + "/callback?" + url.Values{
		"nonce":     {nonce},
		"issuerUrl": {"https://identitycenter.amazonaws.com/ssoins-remote"},
		"tenantUrl": {"https://remote.quick.aws.com"},
		"idcRegion": {"us-west-2"},
	}.Encode()
	paste1Body, _ := json.Marshal(map[string]string{"url": hop1Pasted})
	paste1Rec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/aws-quick/paste", string(paste1Body), false)
	if paste1Rec.Code != http.StatusOK {
		t.Fatalf("paste hop-1 status = %d: %s", paste1Rec.Code, paste1Rec.Body.String())
	}
	var paste1Resp struct {
		Status       string `json:"status"`
		State        string `json:"state"`
		AuthorizeURL string `json:"authorize_url"`
	}
	_ = json.Unmarshal(paste1Rec.Body.Bytes(), &paste1Resp)
	if paste1Resp.Status != "discovery_saved" || paste1Resp.State != startResp.State || !strings.Contains(paste1Resp.AuthorizeURL, "client_id=paste-client") {
		t.Fatalf("unexpected paste1Resp: %+v", paste1Resp)
	}

	// Paste Hop-2 code URL.
	hop2Pasted := startResp.RedirectURI + "/?" + url.Values{
		"code":  {"remote-code-123"},
		"state": {startResp.State},
	}.Encode()
	paste2Body, _ := json.Marshal(map[string]string{"url": hop2Pasted})
	paste2Rec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/aws-quick/paste", string(paste2Body), false)
	if paste2Rec.Code != http.StatusOK {
		t.Fatalf("paste hop-2 status = %d: %s", paste2Rec.Code, paste2Rec.Body.String())
	}
	var paste2Resp struct {
		Status string `json:"status"`
		State  string `json:"state"`
	}
	_ = json.Unmarshal(paste2Rec.Body.Bytes(), &paste2Resp)
	if paste2Resp.Status != "done" || paste2Resp.State != startResp.State {
		t.Fatalf("unexpected paste2Resp: %+v", paste2Resp)
	}

	// Exchange creates the connection.
	exRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/aws-quick/exchange", `{"state":"`+startResp.State+`"}`, false)
	if exRec.Code != http.StatusCreated {
		t.Fatalf("exchange status = %d: %s", exRec.Code, exRec.Body.String())
	}
	var conn providers.ConnectionSummary
	_ = json.Unmarshal(exRec.Body.Bytes(), &conn)
	if conn.Name != "Remote Quick" {
		t.Fatalf("conn.Name = %q, want Remote Quick", conn.Name)
	}
}

func TestAwsQuickEntitlementsProbeTestAndAutoRefresh(t *testing.T) {
	t.Parallel()

	var qbsCalls int32
	var refreshCalls int32
	var seenAccountID string
	var seenTraceHeader string

	var srvURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/client/register":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"clientId":     "quick-probe-cid",
				"clientSecret": "quick-probe-csec",
			})
		case "/token":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["grantType"] == "refresh_token" {
				atomic.AddInt32(&refreshCalls, 1)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"idToken":      "quick-id-token-refreshed",
					"accessToken":  "quick-access-refreshed",
					"refreshToken": "quick-refresh-rotated",
					"expiresIn":    3600,
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"idToken":      "quick-id-token-expired",
				"accessToken":  "quick-access-initial",
				"refreshToken": "quick-refresh-initial",
				"expiresIn":    3600,
			})
		case "/integration/quick-work/user-info":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"email":        "probe@quick.aws",
				"displayName":  "Probe User",
				"awsAccountId": "123456789012",
				"userArn":      "arn:aws:quicksight:us-east-1:123456789012:user/default/probe@quick.aws",
			})
		case "/integration/quick-work/qbs-event":
			call := atomic.AddInt32(&qbsCalls, 1)
			seenTraceHeader = r.Header.Get("x-amzn-qs-internal-trace-id")
			var payload struct {
				EventType string `json:"eventType"`
				EventBody struct {
					GetAvailableEntitlementUnits struct {
						AWSAccountID string `json:"awsAccountId"`
					} `json:"getAvailableEntitlementUnits"`
				} `json:"eventBody"`
			}
			_ = json.NewDecoder(r.Body).Decode(&payload)
			seenAccountID = payload.EventBody.GetAvailableEntitlementUnits.AWSAccountID
			if payload.EventType != "getAvailableEntitlementUnits" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if call == 1 {
				// First call with expired token -> 401 triggers refresh + retry.
				if r.Header.Get("Authorization") != "Bearer quick-id-token-expired" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if call == 2 {
				// Second call after refresh -> 200 OK.
				if r.Header.Get("Authorization") != "Bearer quick-id-token-refreshed" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"eventResponse": map[string]any{
						"getAvailableEntitlementUnitsResponse": map[string]any{
							"entitlements": []map[string]any{
								{"unitType": "AUTOMATION_AGENT_HOUR_UNIT", "numMinsProvisioned": 600, "numMinsAvailable": 480},
							},
						},
					},
				})
				return
			}
			// Third call: 429 throttling -> soft success (healthy account must not be marked broken).
			w.WriteHeader(http.StatusTooManyRequests)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	srvURL = srv.URL

	h := newTestHarness(t, nil, map[string]providers.OAuthConfig{
		"aws-quick": {
			TokenURL: srvURL + "/token",
			ExtraTokenParams: map[string]string{
				"user_info_url": srvURL + "/integration/quick-work/user-info",
			},
		},
	})
	defer func() { _ = h.mod.Stop(context.Background()) }()

	// Complete paste-back OAuth flow to create connection with tenant_url = srvURL.
	startRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/aws-quick/start", `{}`, false)
	if startRec.Code != http.StatusOK {
		t.Fatalf("start status = %d: %s", startRec.Code, startRec.Body.String())
	}
	var startResp struct {
		State        string `json:"state"`
		AuthorizeURL string `json:"authorize_url"`
		RedirectURI  string `json:"redirect_uri"`
	}
	_ = json.Unmarshal(startRec.Body.Bytes(), &startResp)
	parsedAuth, _ := url.Parse(startResp.AuthorizeURL)
	nonce := strings.Split(parsedAuth.Query().Get("qwSurface"), ":")[1]

	p1, _ := json.Marshal(map[string]string{
		"url": startResp.RedirectURI + "/callback?" + url.Values{
			"nonce":     {nonce},
			"issuerUrl": {"https://identitycenter.amazonaws.com/ssoins-probe"},
			"tenantUrl": {srvURL},
			"idcRegion": {"us-east-1"},
		}.Encode(),
	})
	if rec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/aws-quick/paste", string(p1), false); rec.Code != http.StatusOK {
		t.Fatalf("paste1 status = %d: %s", rec.Code, rec.Body.String())
	}
	p2, _ := json.Marshal(map[string]string{
		"url": startResp.RedirectURI + "/?" + url.Values{"code": {"c1"}, "state": {startResp.State}}.Encode(),
	})
	if rec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/aws-quick/paste", string(p2), false); rec.Code != http.StatusOK {
		t.Fatalf("paste2 status = %d: %s", rec.Code, rec.Body.String())
	}
	exRec := h.doRequest(t, http.MethodPost, "/api/providers/oauth/aws-quick/exchange", `{"state":"`+startResp.State+`"}`, false)
	if exRec.Code != http.StatusCreated {
		t.Fatalf("exchange status = %d: %s", exRec.Code, exRec.Body.String())
	}
	var conn providers.ConnectionSummary
	_ = json.Unmarshal(exRec.Body.Bytes(), &conn)

	// 1. Test connection -> 401 -> auto refresh -> retry -> 200 OK.
	testRec := h.doRequest(t, http.MethodPost, "/api/providers/connections/"+conn.ID+"/test", "", false)
	if testRec.Code != http.StatusOK {
		t.Fatalf("test status = %d: %s", testRec.Code, testRec.Body.String())
	}
	var testResp struct {
		OK         bool   `json:"ok"`
		StatusCode int    `json:"status_code"`
		Error      string `json:"error"`
	}
	_ = json.Unmarshal(testRec.Body.Bytes(), &testResp)
	if !testResp.OK || testResp.StatusCode != http.StatusOK {
		t.Fatalf("expected ok=true status=200 after auto-refresh, got %+v", testResp)
	}
	if atomic.LoadInt32(&refreshCalls) != 1 {
		t.Fatalf("expected 1 refresh call, got %d", atomic.LoadInt32(&refreshCalls))
	}
	if seenAccountID != "123456789012" {
		t.Fatalf("expected awsAccountId 123456789012 in qbs-event probe, got %q", seenAccountID)
	}
	if seenTraceHeader == "" {
		t.Fatalf("expected x-amzn-qs-internal-trace-id header to be set")
	}

	// 2. Test connection when Quick throttles (429) -> soft success (ok=true).
	testRec2 := h.doRequest(t, http.MethodPost, "/api/providers/connections/"+conn.ID+"/test", "", false)
	var testResp2 struct {
		OK         bool   `json:"ok"`
		StatusCode int    `json:"status_code"`
		Error      string `json:"error"`
	}
	_ = json.Unmarshal(testRec2.Body.Bytes(), &testResp2)
	if !testResp2.OK || testResp2.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected soft-success ok=true on 429 throttle, got %+v", testResp2)
	}
}

package providers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func makeTestJWT(t *testing.T, payload map[string]any) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal jwt payload: %v", err)
	}
	body := base64.RawURLEncoding.EncodeToString(bodyBytes)
	return header + "." + body + ".sig"
}

func TestOAuthProviderRegistryMetadata(t *testing.T) {
	m := New(Options{})

	expected := map[string]struct {
		flowType     string
		fixedPort    int
		callbackPath string
	}{
		"claude":         {flowType: oauthFlowAuthCodePKCE},
		"codex":          {flowType: oauthFlowAuthCodePKCE, fixedPort: 1455, callbackPath: "/auth/callback"},
		"gemini-oauth":   {flowType: oauthFlowAuthCodePKCE},
		"qwen":           {flowType: oauthFlowAuthCodePKCE},
		"antigravity":    {flowType: oauthFlowAuthCodePKCE},
		"github-copilot": {flowType: oauthFlowDeviceCode},
		"kiro":           {flowType: oauthFlowDeviceCode},
		"freepi":         {flowType: oauthFlowDeviceCode},
		"aws-quick":      {flowType: oauthFlowLoopbackProxy},
	}

	for id, want := range expected {
		strat := m.getOAuthStrategy(id)
		if strat == nil {
			t.Fatalf("expected oauth strategy for %q, got nil", id)
		}
		if strat.FlowType != want.flowType {
			t.Errorf("provider %q flowType = %q, want %q", id, strat.FlowType, want.flowType)
		}
		if strat.FixedPort != want.fixedPort {
			t.Errorf("provider %q fixedPort = %d, want %d", id, strat.FixedPort, want.fixedPort)
		}
		if want.callbackPath != "" && strat.CallbackPath != want.callbackPath {
			t.Errorf("provider %q callbackPath = %q, want %q", id, strat.CallbackPath, want.callbackPath)
		}
	}
}

func TestExtractCodexAccountInfo(t *testing.T) {
	idToken := makeTestJWT(t, map[string]any{
		"email": "operator@example.com",
		"https://api.openai.com/auth": map[string]any{
			"chatgpt_account_id": "acct_12345",
			"chatgpt_plan_type":  "plus",
		},
	})

	info := extractCodexAccountInfo(idToken)
	if info.Email != "operator@example.com" {
		t.Fatalf("email = %q, want operator@example.com", info.Email)
	}
	if info.ChatGPTAccountID != "acct_12345" {
		t.Fatalf("chatgpt_account_id = %q, want acct_12345", info.ChatGPTAccountID)
	}
	if info.ChatGPTPlanType != "plus" {
		t.Fatalf("chatgpt_plan_type = %q, want plus", info.ChatGPTPlanType)
	}
}

func TestOAuthLoopbackProxyMultiProvider(t *testing.T) {
	proxy := newOAuthLoopbackProxy()
	defer proxy.stopAll()

	port, err := proxy.start("test-provider", 0, "/custom/callback", "st_multi_1")
	if err != nil {
		t.Fatalf("proxy.start failed: %v", err)
	}
	if port <= 0 {
		t.Fatalf("expected ephemeral port > 0, got %d", port)
	}

	if st, _ := proxy.status("test-provider", "st_multi_1"); st != "waiting" {
		t.Fatalf("initial status = %q, want waiting", st)
	}

	callbackURL := (&url.URL{
		Scheme:   "http",
		Host:     "127.0.0.1:" + strconv.Itoa(port),
		Path:     "/custom/callback",
		RawQuery: "code=auth_code_99&state=st_multi_1",
	}).String()

	resp, err := http.Get(callbackURL)
	if err != nil {
		t.Fatalf("GET callback failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("callback HTTP status = %d, want 200", resp.StatusCode)
	}

	if st, errMsg := proxy.status("test-provider", "st_multi_1"); st != "done" || errMsg != "" {
		t.Fatalf("after callback status = (%q, %q), want (done, '')", st, errMsg)
	}

	code, ok := proxy.consumeCode("test-provider", "st_multi_1")
	if !ok || code != "auth_code_99" {
		t.Fatalf("consumeCode = (%q, %v), want (auth_code_99, true)", code, ok)
	}

	if _, ok := proxy.consumeCode("test-provider", "st_multi_1"); ok {
		t.Fatalf("expected second consumeCode to return false")
	}

	proxy.stop("test-provider")
}

func TestEnrichAntigravityFiltersBlockedProjectAndSetsHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/userinfo":
			_ = json.NewEncoder(w).Encode(map[string]any{"email": "ag@example.com"})
		case "/loadCodeAssist":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"cloudaicompanionProject": map[string]any{"id": "my-real-gcp-proj"},
				"allowedTiers": []map[string]any{
					{"id": "pro-tier", "isDefault": true},
				},
			})
		case "/loadBlocked":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"cloudaicompanionProject": "aicode-consumers",
			})
		}
	}))
	defer srv.Close()

	m := New(Options{HTTPClient: srv.Client()})

	res := enrichAntigravity(context.Background(), m, srv.URL+"/userinfo", srv.URL+"/loadCodeAssist", srv.URL+"/onboardUser", "tok")
	if res.DisplayLabel != "Antigravity (ag@example.com)" {
		t.Fatalf("DisplayLabel = %q, want Antigravity (ag@example.com)", res.DisplayLabel)
	}
	if res.ExtraHeaders["x-goog-user-project"] != "my-real-gcp-proj" {
		t.Fatalf("x-goog-user-project = %q, want my-real-gcp-proj", res.ExtraHeaders["x-goog-user-project"])
	}
	if res.Metadata["tier_id"] != "pro-tier" {
		t.Fatalf("tier_id = %v, want pro-tier", res.Metadata["tier_id"])
	}

	blocked := enrichAntigravity(context.Background(), m, srv.URL+"/userinfo", srv.URL+"/loadBlocked", srv.URL+"/onboardUser", "tok")
	if _, hasHeader := blocked.ExtraHeaders["x-goog-user-project"]; hasHeader {
		t.Fatalf("expected blocked project ID to be excluded from x-goog-user-project")
	}
}

func TestEnrichAntigravityAutoOnboardsNewAccountAndSendsIDEFingerprint(t *testing.T) {
	var sawUA string
	var sawReqSource string
	var sawClientMeta string
	var onboardCalled bool

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/userinfo":
			_ = json.NewEncoder(w).Encode(map[string]any{"email": "new.ag@example.com"})
		case "/loadCodeAssist":
			sawUA = r.Header.Get("User-Agent")
			sawReqSource = r.Header.Get("x-request-source")
			sawClientMeta = r.Header.Get("Client-Metadata")
			// New account has no cloudaicompanionProject yet.
			_ = json.NewEncoder(w).Encode(map[string]any{
				"allowedTiers": []map[string]any{
					{"id": "ultra-tier", "isDefault": true},
				},
			})
		case "/onboardUser":
			onboardCalled = true
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["tierId"] != "ultra-tier" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"done": true,
				"response": map[string]any{
					"cloudaicompanionProject": map[string]any{"id": "onboarded-gcp-proj-42"},
				},
			})
		}
	}))
	defer srv.Close()

	m := New(Options{HTTPClient: srv.Client()})
	res := enrichAntigravity(context.Background(), m, srv.URL+"/userinfo", srv.URL+"/loadCodeAssist", srv.URL+"/onboardUser", "tok")
	if !onboardCalled {
		t.Fatalf("expected onboardUser to be called when loadCodeAssist returns no projectId")
	}
	if sawUA != antigravityIDEUserAgent {
		t.Fatalf("User-Agent = %q, want %q", sawUA, antigravityIDEUserAgent)
	}
	if sawReqSource != "local" {
		t.Fatalf("x-request-source = %q, want local", sawReqSource)
	}
	if sawClientMeta == "" {
		t.Fatalf("expected Client-Metadata header to be set")
	}
	if res.ExtraHeaders["x-goog-user-project"] != "onboarded-gcp-proj-42" {
		t.Fatalf("x-goog-user-project = %q, want onboarded-gcp-proj-42", res.ExtraHeaders["x-goog-user-project"])
	}
	if res.Metadata["tier_id"] != "ultra-tier" {
		t.Fatalf("tier_id = %v, want ultra-tier", res.Metadata["tier_id"])
	}
}

func TestEnrichGitHubCopilotAndKiroProfileARN(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/copilot":
			_ = json.NewEncoder(w).Encode(map[string]any{"token": "tid=copilot_abc", "expires_at": 1893456000})
		case "/user":
			_ = json.NewEncoder(w).Encode(map[string]any{"login": "octocat", "email": "octo@github.com"})
		case "/ListAvailableProfiles":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"profiles": []map[string]any{{"arn": "arn:aws:codewhisperer:us-east-1:123456789012:profile/ABC"}},
			})
		}
	}))
	defer srv.Close()

	m := New(Options{HTTPClient: srv.Client()})

	gh := enrichGitHubCopilot(context.Background(), m, srv.URL+"/copilot", srv.URL+"/user", "gh_tok")
	if gh.DisplayLabel != "GitHub Copilot (@octocat)" {
		t.Fatalf("DisplayLabel = %q, want GitHub Copilot (@octocat)", gh.DisplayLabel)
	}
	if gh.Metadata["copilot_token"] != "tid=copilot_abc" {
		t.Fatalf("copilot_token = %v, want tid=copilot_abc", gh.Metadata["copilot_token"])
	}

	arn := m.fetchKiroProfileARN(context.Background(), srv.URL+"/ListAvailableProfiles", "kiro_tok")
	if arn != "arn:aws:codewhisperer:us-east-1:123456789012:profile/ABC" {
		t.Fatalf("profile_arn = %q", arn)
	}
}

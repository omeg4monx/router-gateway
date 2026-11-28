package providers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

const (
	oauthFlowAuthCodePKCE  = "authorization_code_pkce"
	oauthFlowAuthCode      = "authorization_code"
	oauthFlowDeviceCode    = "device_code"
	oauthFlowLoopbackProxy = "loopback_proxy"

	antigravityIDEVersion     = "2.34.0"
	antigravityIDEUserAgent   = "antigravity/hub/" + antigravityIDEVersion + " linux/x64"
	antigravityClientMetadata = `{"ideType":"ANTIGRAVITY","platform":"PLATFORM_UNSPECIFIED","pluginType":"GEMINI"}`

	googleUserInfoV1URL   = "https://www.googleapis.com/oauth2/v1/userinfo"
	googleUserInfoURL     = "https://www.googleapis.com/oauth2/v2/userinfo"
	googleLoadCodeAssist  = "https://cloudcode-pa.googleapis.com/v1internal:loadCodeAssist"
	googleOnboardUserURL  = "https://cloudcode-pa.googleapis.com/v1internal:onboardUser"
	githubCopilotTokenURL = "https://api.github.com/copilot_internal/v2/token" //nolint:gosec // G101: public GitHub Copilot token endpoint URL
	githubUserInfoURL     = "https://api.github.com/user"
	kiroListProfilesURL   = "https://codewhisperer.us-east-1.amazonaws.com/ListAvailableProfiles"
)

// oauthPostExchangeResult carries enriched account info and provider-specific
// headers/metadata discovered after token exchange (adopted from xrouter's
// postExchange + mapTokens lifecycle).
type oauthPostExchangeResult struct {
	Email        string
	DisplayLabel string
	ProjectID    string
	ExtraHeaders map[string]string
	Metadata     map[string]any
}

// oauthProviderStrategy encapsulates provider-specific OAuth behavior in a
// pluggable registry entry, mirroring xrouter/src/lib/oauth/providers/*.js.
type oauthProviderStrategy struct {
	ID               string
	FlowType         string
	FixedPort        int
	CallbackPath     string
	ClientIDOptional bool

	BuildAuthURL func(cfg OAuthConfig, redirectURI, state, challenge string) (string, error)
	ExchangeCode func(ctx context.Context, m *Module, cfg OAuthConfig, sess oauthSessionRecord, code, state, verifier, clientSecret string) (oauthTokenResponse, error)
	PostExchange func(ctx context.Context, m *Module, tokenURL string, tok oauthTokenResponse) oauthPostExchangeResult
	RefreshToken func(ctx context.Context, m *Module, tokenURL, clientID, clientSecret string, stored storedCredentials) (oauthTokenResponse, error)

	StartDevice func(ctx context.Context, m *Module, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error)
	PollDevice  func(ctx context.Context, m *Module, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string)
}

func defaultOAuthStrategies() map[string]*oauthProviderStrategy {
	return map[string]*oauthProviderStrategy{
		"claude": {
			ID:           "claude",
			FlowType:     oauthFlowAuthCodePKCE,
			CallbackPath: "/api/providers/oauth/callback",
			ExchangeCode: exchangeClaudeCode,
			RefreshToken: refreshClaudeToken,
		},
		providerIDCodex: {
			ID:           providerIDCodex,
			FlowType:     oauthFlowAuthCodePKCE,
			FixedPort:    1455,
			CallbackPath: "/auth/callback",
			PostExchange: postExchangeCodex,
		},
		"gemini-oauth": {
			ID:           "gemini-oauth",
			FlowType:     oauthFlowAuthCodePKCE,
			CallbackPath: "/api/providers/oauth/callback",
			PostExchange: postExchangeGeminiOAuth,
		},
		"qwen": {
			ID:           "qwen",
			FlowType:     oauthFlowAuthCodePKCE,
			CallbackPath: "/api/providers/oauth/callback",
		},
		"antigravity": {
			ID:           "antigravity",
			FlowType:     oauthFlowAuthCodePKCE,
			CallbackPath: "/api/providers/oauth/callback",
			PostExchange: postExchangeAntigravity,
		},
		providerIDGitHubCopilot: {
			ID:               providerIDGitHubCopilot,
			FlowType:         oauthFlowDeviceCode,
			ClientIDOptional: true,
			StartDevice: func(ctx context.Context, m *Module, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error) {
				return m.startGitHubCopilotDevice(ctx, state, cfg, connName, baseURL, models)
			},
			PollDevice: func(ctx context.Context, m *Module, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
				return m.pollGitHubCopilotDevice(ctx, sess)
			},
			PostExchange: postExchangeGitHubCopilot,
		},
		providerIDKiro: {
			ID:               providerIDKiro,
			FlowType:         oauthFlowDeviceCode,
			ClientIDOptional: true,
			StartDevice: func(ctx context.Context, m *Module, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error) {
				return m.startKiroDevice(ctx, state, cfg, connName, baseURL, models)
			},
			PollDevice: func(ctx context.Context, m *Module, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
				return m.pollKiroDevice(ctx, sess)
			},
			PostExchange: postExchangeKiro,
			RefreshToken: refreshAWSOIDCToken,
		},
		providerIDFreePI: {
			ID:               providerIDFreePI,
			FlowType:         oauthFlowDeviceCode,
			ClientIDOptional: true,
			StartDevice: func(ctx context.Context, m *Module, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error) {
				return m.startFreePIDevice(ctx, state, cfg, connName, baseURL, models)
			},
			PollDevice: func(ctx context.Context, m *Module, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
				return m.pollFreePIDevice(ctx, sess)
			},
		},
		providerIDAWSQuick: {
			ID:               providerIDAWSQuick,
			FlowType:         oauthFlowLoopbackProxy,
			ClientIDOptional: true,
			CallbackPath:     "/callback",
			RefreshToken:     refreshAWSOIDCToken,
		},
		providerIDCodeBuddyCN:   newCodeBuddyOAuthStrategy(providerIDCodeBuddyCN),
		providerIDCodeBuddyIntl: newCodeBuddyOAuthStrategy(providerIDCodeBuddyIntl),
		providerIDJetBrains:     newJetBrainsOAuthStrategy(),
		"kimi": newRFC8628Strategy(rfc8628DeviceSpec{ //nolint:gosec // G101: public OAuth device client ID and token endpoint URL
			ProviderID:              "kimi",
			DefaultDeviceCodeURL:    "https://auth.kimi.com/api/oauth/device_authorization",
			DefaultTokenURL:         "https://auth.kimi.com/api/oauth/token",
			DefaultClientID:         "17e5f671-d194-4dfb-9706-5516cb48c098",
			FallbackVerificationURI: "https://www.kimi.com/code/authorize_device",
			DisplayPrefix:           "Kimi Code",
		}, nil),
		"grok-cli": newRFC8628Strategy(rfc8628DeviceSpec{ //nolint:gosec // G101: public OAuth device client ID and token endpoint URL
			ProviderID:           "grok-cli",
			DefaultDeviceCodeURL: "https://auth.x.ai/oauth2/device/code",
			DefaultTokenURL:      "https://auth.x.ai/oauth2/token",
			DefaultClientID:      "b1a00492-073a-47ea-816f-4c329264a828",
			DefaultScope:         "openid profile email offline_access grok-cli:access api:access conversations:read conversations:write",
			DisplayPrefix:        "Grok CLI",
		}, nil),
		"nous": newRFC8628Strategy(rfc8628DeviceSpec{ //nolint:gosec // G101: public OAuth device client ID and token endpoint URL
			ProviderID:           "nous",
			DefaultDeviceCodeURL: "https://portal.nousresearch.com/api/oauth/device/code",
			DefaultTokenURL:      "https://portal.nousresearch.com/api/oauth/token",
			DefaultClientID:      "hermes-cli",
			DefaultScope:         "inference:invoke inference:mint_agent_key",
			DisplayPrefix:        "Nous",
		}, postExchangeNous),
		"minimax-code": newRFC8628Strategy(rfc8628DeviceSpec{ //nolint:gosec // G101: public OAuth device client ID and token endpoint URL
			ProviderID:           "minimax-code",
			DefaultDeviceCodeURL: "https://account.minimax.cn/oauth2/device/code",
			DefaultTokenURL:      "https://account.minimax.cn/oauth2/token",
			DefaultClientID:      "mcode-public",
			DefaultScope:         "agent.default",
			Audience:             "agent-backend",
			UsePKCE:              true,
			DisplayPrefix:        "MiniMax Coding",
		}, nil),
		"minimax-code-global": newRFC8628Strategy(rfc8628DeviceSpec{ //nolint:gosec // G101: public OAuth device client ID and token endpoint URL
			ProviderID:           "minimax-code-global",
			DefaultDeviceCodeURL: "https://account.minimax.io/oauth2/device/code",
			DefaultTokenURL:      "https://account.minimax.io/oauth2/token",
			DefaultClientID:      "mcode-public",
			DefaultScope:         "agent.default",
			Audience:             "agent-backend",
			UsePKCE:              true,
			DisplayPrefix:        "MiniMax Coding Global",
		}, nil),
		"muse": newRFC8628Strategy(rfc8628DeviceSpec{ //nolint:gosec // G101: public OAuth device client ID and token endpoint URL
			ProviderID:           "muse",
			DefaultDeviceCodeURL: "https://auth.meta.com/oidc/device/authorization/",
			DefaultTokenURL:      "https://auth.meta.com/oidc/device/token/",
			DefaultClientID:      "1031625952748946",
			DisplayPrefix:        "Meta Muse",
		}, nil),
		"kilocode": {
			ID:               "kilocode",
			FlowType:         oauthFlowDeviceCode,
			ClientIDOptional: true,
			StartDevice: func(ctx context.Context, m *Module, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error) {
				return m.startKiloCodeDevice(ctx, state, cfg, connName, baseURL, models)
			},
			PollDevice: func(ctx context.Context, m *Module, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
				return m.pollKiloCodeDevice(ctx, sess)
			},
		},
		"iflow": {
			ID:           "iflow",
			FlowType:     oauthFlowAuthCode,
			CallbackPath: "/api/providers/oauth/callback",
			PostExchange: postExchangeIFlow,
		},
		"cline": {
			ID:               "cline",
			FlowType:         oauthFlowAuthCode,
			ClientIDOptional: true,
			CallbackPath:     "/api/providers/oauth/callback",
		},
		"clinepass": {
			ID:               "clinepass",
			FlowType:         oauthFlowAuthCode,
			ClientIDOptional: true,
			CallbackPath:     "/api/providers/oauth/callback",
		},
		"gitlab": {
			ID:               "gitlab",
			FlowType:         oauthFlowAuthCodePKCE,
			ClientIDOptional: true,
			CallbackPath:     "/api/providers/oauth/callback",
		},
	}
}

func newRFC8628Strategy(
	spec rfc8628DeviceSpec,
	postExchange func(context.Context, *Module, string, oauthTokenResponse) oauthPostExchangeResult,
) *oauthProviderStrategy {
	return &oauthProviderStrategy{
		ID:               spec.ProviderID,
		FlowType:         oauthFlowDeviceCode,
		ClientIDOptional: true,
		StartDevice: func(ctx context.Context, m *Module, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error) {
			return m.startRFC8628Device(ctx, spec, state, cfg, connName, baseURL, models)
		},
		PollDevice: func(ctx context.Context, m *Module, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
			return m.pollRFC8628Device(ctx, spec, sess)
		},
		PostExchange: postExchange,
	}
}

func (m *Module) getOAuthStrategy(providerID string) *oauthProviderStrategy {
	m.mu.RLock()
	strategy, ok := m.oauthStrategies[providerID]
	m.mu.RUnlock()
	if ok && strategy != nil {
		return strategy
	}
	return &oauthProviderStrategy{
		ID:           providerID,
		FlowType:     oauthFlowAuthCodePKCE,
		CallbackPath: "/api/providers/oauth/callback",
	}
}

// --- JWT Helpers (adopted from xrouter/src/lib/oauth/providerHelpers.js) ---

func decodeJWTPayload(jwt string) map[string]any {
	parts := strings.Split(strings.TrimSpace(jwt), ".")
	if len(parts) != 3 || parts[1] == "" {
		return nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		// Retry with standard padded base64url if upstream includes '='.
		raw, err = base64.URLEncoding.DecodeString(parts[1])
		if err != nil {
			return nil
		}
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	return payload
}

type codexAccountInfo struct {
	Email            string
	ChatGPTAccountID string
	ChatGPTPlanType  string
}

func extractCodexAccountInfo(idToken string) codexAccountInfo {
	payload := decodeJWTPayload(idToken)
	if payload == nil {
		return codexAccountInfo{}
	}
	info := codexAccountInfo{
		Email: stringClaim(payload, "email"),
	}
	if authClaims, ok := payload["https://api.openai.com/auth"].(map[string]any); ok {
		info.ChatGPTAccountID = stringClaim(authClaims, "chatgpt_account_id")
		info.ChatGPTPlanType = stringClaim(authClaims, "chatgpt_plan_type")
	}
	if info.ChatGPTAccountID == "" {
		info.ChatGPTAccountID = stringClaim(payload, "account_id")
	}
	if info.ChatGPTPlanType == "" {
		info.ChatGPTPlanType = stringClaim(payload, "plan_type")
	}
	return info
}

func extractEmailFromJWT(token string) string {
	payload := decodeJWTPayload(token)
	if payload == nil {
		return ""
	}
	if v := stringClaim(payload, "email"); v != "" {
		return v
	}
	if v := stringClaim(payload, "preferred_username"); v != "" {
		return v
	}
	return stringClaim(payload, "sub")
}

func stringClaim(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(s)
}

// --- Provider-specific Exchange & Refresh Handlers ---

func exchangeClaudeCode(ctx context.Context, m *Module, cfg OAuthConfig, sess oauthSessionRecord, code, state, verifier, clientSecret string) (oauthTokenResponse, error) {
	if !isAnthropicTokenURL(cfg.TokenURL) {
		return defaultExchangePKCECode(ctx, m, cfg, sess, code, verifier, clientSecret)
	}
	authCode, codeState := code, state
	if idx := strings.IndexByte(authCode, '#'); idx >= 0 {
		codeState = authCode[idx+1:]
		authCode = authCode[:idx]
		if codeState == "" {
			codeState = state
		}
	}
	return m.postJSONTokenEndpoint(ctx, cfg.TokenURL, map[string]any{
		"grant_type":    "authorization_code",
		"code":          authCode,
		"state":         codeState,
		"client_id":     sess.ClientID,
		"redirect_uri":  sess.RedirectURI,
		"code_verifier": verifier,
	})
}

func defaultExchangePKCECode(ctx context.Context, m *Module, cfg OAuthConfig, sess oauthSessionRecord, code, verifier, clientSecret string) (oauthTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", sess.RedirectURI)
	form.Set("client_id", sess.ClientID)
	if verifier != "" {
		form.Set("code_verifier", verifier)
	}
	if clientSecret != "" {
		form.Set("client_secret", clientSecret)
	}
	for k, v := range cfg.ExtraTokenParams {
		if strings.TrimSpace(k) != "" && v != "" {
			form.Set(k, v)
		}
	}
	return m.postTokenEndpoint(ctx, cfg.TokenURL, form)
}

func refreshClaudeToken(ctx context.Context, m *Module, tokenURL, clientID, clientSecret string, stored storedCredentials) (oauthTokenResponse, error) {
	if !isAnthropicTokenURL(tokenURL) {
		return defaultRefreshFormToken(ctx, m, tokenURL, clientID, clientSecret, stored)
	}
	payload := map[string]any{
		"grant_type":    "refresh_token",
		"refresh_token": stored.RefreshToken,
	}
	if clientID != "" {
		payload["client_id"] = clientID
	}
	return m.postJSONTokenEndpoint(ctx, tokenURL, payload)
}

func refreshAWSOIDCToken(ctx context.Context, m *Module, tokenURL, clientID, clientSecret string, stored storedCredentials) (oauthTokenResponse, error) {
	awsTok, err := m.postAwsQuickTokenJSON(ctx, tokenURL, map[string]any{
		"clientId":     clientID,
		"clientSecret": clientSecret,
		"grantType":    "refresh_token",
		"refreshToken": stored.RefreshToken,
	})
	if err != nil {
		return oauthTokenResponse{}, err
	}
	accessTok := awsTok.IDToken
	if accessTok == "" {
		accessTok = awsTok.AccessToken
	}
	return oauthTokenResponse{
		AccessToken:  accessTok,
		RefreshToken: awsTok.RefreshToken,
		IDToken:      awsTok.IDToken,
		ExpiresIn:    awsTok.ExpiresIn,
	}, nil
}

func defaultRefreshFormToken(ctx context.Context, m *Module, tokenURL, clientID, clientSecret string, stored storedCredentials) (oauthTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", stored.RefreshToken)
	if clientID != "" {
		form.Set("client_id", clientID)
	}
	if clientSecret != "" {
		form.Set("client_secret", clientSecret)
	}
	return m.postTokenEndpoint(ctx, tokenURL, form)
}

// --- Provider-specific PostExchange Hooks (adopted from xrouter) ---

func postExchangeCodex(_ context.Context, _ *Module, _ string, tok oauthTokenResponse) oauthPostExchangeResult {
	info := extractCodexAccountInfo(tok.IDToken)
	email := info.Email
	if email == "" {
		email = extractEmailFromJWT(tok.AccessToken)
	}
	meta := map[string]any{}
	if email != "" {
		meta["email"] = email
	}
	if info.ChatGPTAccountID != "" {
		meta["chatgpt_account_id"] = info.ChatGPTAccountID
	}
	if info.ChatGPTPlanType != "" {
		meta["chatgpt_plan_type"] = info.ChatGPTPlanType
	}
	label := ""
	if email != "" {
		label = "Codex (" + email + ")"
	}
	return oauthPostExchangeResult{
		Email:        email,
		DisplayLabel: label,
		Metadata:     meta,
	}
}

func isHostMatch(rawURL, expectedSuffix string) bool {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	return host == expectedSuffix || strings.HasSuffix(host, "."+expectedSuffix)
}

func oauthPlatformEnum() int {
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return 2
		}
		return 1
	case "linux":
		if runtime.GOARCH == "arm64" {
			return 4
		}
		return 3
	case "windows":
		return 5
	default:
		return 0
	}
}

func (m *Module) fetchGoogleUserEmail(ctx context.Context, userInfoURL, accessToken string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL+"?alt=json", http.NoBody)
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-request-source", "local")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ""
	}
	var payload struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return ""
	}
	return strings.TrimSpace(payload.Email)
}

func extractCompanionProjectID(data map[string]any) string {
	if data == nil {
		return ""
	}
	switch proj := data["cloudaicompanionProject"].(type) {
	case string:
		if s := strings.TrimSpace(proj); s != "" {
			return s
		}
	case map[string]any:
		if s := stringClaim(proj, "id"); s != "" {
			return s
		}
	}
	if nested, ok := data["response"].(map[string]any); ok {
		return extractCompanionProjectID(nested)
	}
	return ""
}

func (m *Module) fetchGoogleCodeAssistInfo(ctx context.Context, endpoint, accessToken string, useAntigravityHeaders bool) (projectID, tierID string) {
	body := map[string]any{
		"metadata": map[string]any{
			"ideType":    9,
			"platform":   oauthPlatformEnum(),
			"pluginType": 2,
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", ""
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(raw)))
	if err != nil {
		return "", ""
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	if useAntigravityHeaders {
		req.Header.Set("User-Agent", antigravityIDEUserAgent)
		req.Header.Set("Client-Metadata", antigravityClientMetadata)
		req.Header.Set("x-request-source", "local")
	}
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", ""
	}
	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", ""
	}
	projectID = extractCompanionProjectID(data)
	if tiers, ok := data["allowedTiers"].([]any); ok {
		for _, item := range tiers {
			if tm, ok := item.(map[string]any); ok {
				if isDef, _ := tm["isDefault"].(bool); isDef {
					if id := stringClaim(tm, "id"); id != "" {
						tierID = id
						break
					}
				}
			}
		}
	}
	return projectID, tierID
}

func (m *Module) onboardAntigravityUser(ctx context.Context, onboardEndpoint, accessToken, tierID string) string {
	payload := map[string]any{
		"tierId": tierID,
		"metadata": map[string]any{
			"ideType":    "ANTIGRAVITY",
			"ideName":    "antigravity",
			"ideVersion": antigravityIDEVersion,
			"platform":   oauthPlatformEnum(),
			"pluginType": 2,
		},
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, onboardEndpoint, strings.NewReader(string(raw)))
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", antigravityIDEUserAgent)
	req.Header.Set("Client-Metadata", antigravityClientMetadata)
	req.Header.Set("x-request-source", "local")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ""
	}
	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ""
	}
	return extractCompanionProjectID(data)
}

func postExchangeAntigravity(ctx context.Context, m *Module, tokenURL string, tok oauthTokenResponse) oauthPostExchangeResult {
	if !isHostMatch(tokenURL, "googleapis.com") || tok.AccessToken == "" {
		return oauthPostExchangeResult{}
	}
	return enrichAntigravity(ctx, m, googleUserInfoV1URL, googleLoadCodeAssist, googleOnboardUserURL, tok.AccessToken)
}

func enrichAntigravity(ctx context.Context, m *Module, userInfoURL, loadEndpoint, onboardEndpoint, accessToken string) oauthPostExchangeResult {
	email := m.fetchGoogleUserEmail(ctx, userInfoURL, accessToken)
	projectID, tierID := m.fetchGoogleCodeAssistInfo(ctx, loadEndpoint, accessToken, true)
	if tierID == "" {
		tierID = "free-tier"
	}

	// Auto-onboard new Antigravity accounts (xrouter antigravity.js parity):
	// provisions cloudaicompanionProject and records the resolved tier.
	if projectID == "" && onboardEndpoint != "" {
		if provisioned := m.onboardAntigravityUser(ctx, onboardEndpoint, accessToken, tierID); provisioned != "" {
			projectID = provisioned
		} else {
			projectID, _ = m.fetchGoogleCodeAssistInfo(ctx, loadEndpoint, accessToken, true)
		}
	}

	res := oauthPostExchangeResult{
		Email:    email,
		Metadata: map[string]any{"tier_id": tierID},
	}
	if email != "" {
		res.DisplayLabel = "Antigravity (" + email + ")"
		res.Metadata["email"] = email
	}
	if projectID != "" && !isBlockedAntigravityProjectID(projectID) {
		res.ProjectID = projectID
		res.Metadata["project_id"] = projectID
		res.ExtraHeaders = map[string]string{
			"x-goog-user-project": projectID,
		}
	}
	return res
}

func postExchangeGeminiOAuth(ctx context.Context, m *Module, tokenURL string, tok oauthTokenResponse) oauthPostExchangeResult {
	if !isHostMatch(tokenURL, "googleapis.com") || tok.AccessToken == "" {
		return oauthPostExchangeResult{}
	}
	return enrichGeminiOAuth(ctx, m, googleUserInfoURL, googleLoadCodeAssist, tok.AccessToken)
}

func enrichGeminiOAuth(ctx context.Context, m *Module, userInfoURL, loadEndpoint, accessToken string) oauthPostExchangeResult {
	email := m.fetchGoogleUserEmail(ctx, userInfoURL, accessToken)
	projectID, _ := m.fetchGoogleCodeAssistInfo(ctx, loadEndpoint, accessToken, false)

	res := oauthPostExchangeResult{
		Email:     email,
		ProjectID: projectID,
		Metadata:  map[string]any{},
	}
	if email != "" {
		res.DisplayLabel = "Google Gemini (" + email + ")"
		res.Metadata["email"] = email
	}
	if projectID != "" {
		res.Metadata["project_id"] = projectID
	}
	return res
}

func postExchangeGitHubCopilot(ctx context.Context, m *Module, tokenURL string, tok oauthTokenResponse) oauthPostExchangeResult {
	if !isHostMatch(tokenURL, "github.com") || tok.AccessToken == "" {
		return oauthPostExchangeResult{}
	}
	return enrichGitHubCopilot(ctx, m, githubCopilotTokenURL, githubUserInfoURL, tok.AccessToken)
}

func enrichGitHubCopilot(ctx context.Context, m *Module, copilotURL, userURL, accessToken string) oauthPostExchangeResult {
	res := oauthPostExchangeResult{Metadata: map[string]any{}}

	fetchGH := func(endpoint string, dst any) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
		if err != nil {
			return
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		req.Header.Set("User-Agent", "GithubCopilot/1.155.0")
		resp, err := m.httpClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			_ = json.NewDecoder(resp.Body).Decode(dst)
		}
	}

	var copilot struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}
	fetchGH(copilotURL, &copilot)
	if copilot.Token != "" {
		res.Metadata["copilot_token"] = copilot.Token
		res.Metadata["copilot_token_expires_at"] = copilot.ExpiresAt
	}

	var user struct {
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	fetchGH(userURL, &user)
	if user.Login != "" {
		res.Metadata["github_login"] = user.Login
		res.DisplayLabel = "GitHub Copilot (@" + user.Login + ")"
	}
	if user.Email != "" {
		res.Email = user.Email
		res.Metadata["email"] = user.Email
	}
	return res
}

func postExchangeKiro(ctx context.Context, m *Module, tokenURL string, tok oauthTokenResponse) oauthPostExchangeResult {
	res := oauthPostExchangeResult{Metadata: map[string]any{}}
	if email := extractEmailFromJWT(tok.AccessToken); email != "" {
		res.Email = email
		res.Metadata["email"] = email
		res.DisplayLabel = "Kiro (" + email + ")"
	}
	if isHostMatch(tokenURL, "amazonaws.com") && tok.AccessToken != "" {
		if arn := m.fetchKiroProfileARN(ctx, kiroListProfilesURL, tok.AccessToken); arn != "" {
			res.Metadata["profile_arn"] = arn
		}
	}
	return res
}

func (m *Module) fetchKiroProfileARN(ctx context.Context, endpoint, accessToken string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(`{"maxResults":10}`))
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ""
	}
	var payload struct {
		Profiles []struct {
			ARN string `json:"arn"`
		} `json:"profiles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return ""
	}
	for _, p := range payload.Profiles {
		if arn := strings.TrimSpace(p.ARN); arn != "" {
			return arn
		}
	}
	return ""
}

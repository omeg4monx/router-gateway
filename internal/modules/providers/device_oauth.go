package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

const (
	providerIDGitHubCopilot = "github-copilot"
	providerIDKiro          = "kiro"
	providerIDFreePI        = "freepi"

	githubDeviceCodeURL     = "https://github.com/login/device/code"
	githubOAuthTokenURL     = "https://github.com/login/oauth/access_token"
	kiroOIDCBaseURL         = "https://oidc.us-east-1.amazonaws.com"
	kiroStartURL            = "https://view.awsapps.com/start"
	kiroIssuerURL           = "https://identitycenter.amazonaws.com/ssoins-722374e8c3c8e6c6"
	freepiDeviceCodeURL     = "https://api.freepi.ai/auth/github/device"
	freepiConsentURL        = "https://api.freepi.ai/auth/consent"
	freepiTokenURL          = "https://api.freepi.ai/auth/token"
	freepiConsentVersion    = "4"
	defaultDevicePollSec    = 5
	defaultDeviceSessionTTL = 15 * time.Minute
)

func isDeviceCodeProvider(providerID string) bool {
	switch providerID {
	case providerIDGitHubCopilot, providerIDKiro, providerIDFreePI:
		return true
	default:
		return false
	}
}

type deviceOAuthResult struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	Scope        string
	ExpiresIn    int64
	ClientID     string
	ClientSecret string
	TokenURL     string
	DisplayLabel string
}

type deviceOAuthSession struct {
	State           string
	ProviderID      string
	DeviceCode      string
	UserCode        string
	VerificationURI string
	AuthorizeURL    string
	Interval        int
	ClientID        string
	ClientSecret    string
	TokenURL        string
	CodeVerifier    string
	ConnectionName  string
	BaseURL         string
	Models          []string
	Status          string // "pending" | "done" | "error"
	Result          *deviceOAuthResult
	ErrMsg          string
	ExpiresAt       time.Time
}

type deviceOAuthRuntime struct {
	mu       sync.Mutex
	sessions map[string]*deviceOAuthSession
}

func newDeviceOAuthRuntime() *deviceOAuthRuntime {
	return &deviceOAuthRuntime{
		sessions: make(map[string]*deviceOAuthSession),
	}
}

func (rt *deviceOAuthRuntime) sweepLocked(now time.Time) {
	for st, sess := range rt.sessions {
		if !sess.ExpiresAt.IsZero() && !now.Before(sess.ExpiresAt) {
			delete(rt.sessions, st)
		}
	}
}

func (m *Module) postFormJSON(ctx context.Context, targetURL string, form url.Values) (int, []byte, error) {
	validURL, err := validateHTTPURL(targetURL)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, validURL, strings.NewReader(form.Encode()))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxOAuthResponseBytes))
	return resp.StatusCode, body, err
}

func (m *Module) postRawJSON(ctx context.Context, targetURL string, payload any) (int, []byte, error) {
	validURL, err := validateHTTPURL(targetURL)
	if err != nil {
		return 0, nil, err
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, validURL, bytes.NewReader(raw))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxOAuthResponseBytes))
	return resp.StatusCode, body, err
}

func originOfURL(rawURL string) string {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Host
}

func (m *Module) handleDeviceOAuthStart(w http.ResponseWriter, r *http.Request, providerID string, adapter ProviderAdapter, req oauthStartRequest) {
	connName := strings.TrimSpace(req.ConnectionName)
	if connName != "" {
		validated, err := validateName(connName)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		connName = validated
	}

	baseURL := strings.TrimSpace(req.BaseURL)
	if baseURL == "" {
		baseURL = adapter.DefaultBaseURL()
	} else {
		validated, err := validateHTTPURL(baseURL)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "base_url: "+err.Error())
			return
		}
		baseURL = validated
	}

	var models []string
	if req.Models != nil {
		validated, err := validateStringSlice("models", req.Models)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		models = validated
	} else {
		models = adapter.DefaultModels()
	}

	state, err := generateOAuthState()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate OAuth state")
		return
	}

	cfg, _ := m.resolveOAuthConfig(providerID, req.ClientID, req.ClientSecret)

	strat := m.getOAuthStrategy(providerID)
	if strat.StartDevice == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_not_supported", "provider does not support device code OAuth")
		return
	}
	sess, err := strat.StartDevice(r.Context(), m, state, cfg, connName, baseURL, models)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadGateway, "oauth_start_failed", truncateError(err.Error()))
		return
	}

	m.deviceOAuth.mu.Lock()
	m.deviceOAuth.sweepLocked(time.Now().UTC())
	m.deviceOAuth.sessions[state] = sess
	m.deviceOAuth.mu.Unlock()

	httpserver.WriteJSON(w, http.StatusOK, oauthStartResponse{
		State:           sess.State,
		AuthorizeURL:    sess.AuthorizeURL,
		RedirectURI:     sess.VerificationURI,
		ExpiresAt:       sess.ExpiresAt.Format(rfc3339NanoFixed),
		FlowType:        "device_code",
		UserCode:        sess.UserCode,
		VerificationURI: sess.VerificationURI,
		Interval:        sess.Interval,
	})
}

func (m *Module) startGitHubCopilotDevice(ctx context.Context, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error) {
	deviceURL := githubDeviceCodeURL
	if strings.Contains(cfg.AuthorizeURL, "/device/code") {
		deviceURL = cfg.AuthorizeURL
	} else if orig := originOfURL(cfg.TokenURL); orig != "" && !strings.Contains(orig, "github.com") {
		deviceURL = orig + "/login/device/code"
	}
	tokenURL := cfg.TokenURL
	if tokenURL == "" {
		tokenURL = githubOAuthTokenURL
	}

	scope := "read:user"
	if len(cfg.Scopes) > 0 {
		scope = strings.Join(cfg.Scopes, " ")
	}
	form := url.Values{}
	form.Set("client_id", cfg.ClientID)
	form.Set("scope", scope)

	status, body, err := m.postFormJSON(ctx, deviceURL, form)
	if err != nil {
		return nil, fmt.Errorf("github device code request failed: %w", err)
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("github device code returned HTTP %d", status)
	}

	var raw struct {
		DeviceCode      string `json:"device_code"`
		UserCode        string `json:"user_code"`
		VerificationURI string `json:"verification_uri"`
		ExpiresIn       int    `json:"expires_in"`
		Interval        int    `json:"interval"`
		Error           string `json:"error"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode github device code response: %w", err)
	}
	if raw.Error != "" || raw.DeviceCode == "" {
		return nil, fmt.Errorf("github device code error: %s", raw.Error)
	}
	interval := raw.Interval
	if interval <= 0 {
		interval = defaultDevicePollSec
	}
	ttl := defaultDeviceSessionTTL
	if raw.ExpiresIn > 0 {
		ttl = time.Duration(raw.ExpiresIn) * time.Second
	}
	return &deviceOAuthSession{
		State:           state,
		ProviderID:      providerIDGitHubCopilot,
		DeviceCode:      raw.DeviceCode,
		UserCode:        raw.UserCode,
		VerificationURI: raw.VerificationURI,
		AuthorizeURL:    raw.VerificationURI,
		Interval:        interval,
		ClientID:        cfg.ClientID,
		TokenURL:        tokenURL,
		ConnectionName:  connName,
		BaseURL:         baseURL,
		Models:          models,
		Status:          "pending",
		ExpiresAt:       time.Now().UTC().Add(ttl),
	}, nil
}

func (m *Module) startKiroDevice(ctx context.Context, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error) {
	oidcBase := kiroOIDCBaseURL
	tokenURL := kiroOIDCBaseURL + "/token"
	m.mu.RLock()
	override, hasOverride := m.oauthOverrides[providerIDKiro]
	m.mu.RUnlock()
	if hasOverride && strings.TrimSpace(override.TokenURL) != "" {
		if orig := originOfURL(override.TokenURL); orig != "" {
			oidcBase = orig
			tokenURL = strings.TrimSpace(override.TokenURL)
		}
	}

	regStatus, regBody, err := m.postRawJSON(ctx, oidcBase+"/client/register", map[string]any{
		"clientName": "kiro-oauth-client",
		"clientType": "public",
		"scopes":     []string{"codewhisperer:completions", "codewhisperer:analysis", "codewhisperer:conversations"},
		"grantTypes": []string{"urn:ietf:params:oauth:grant-type:device_code", "refresh_token"},
		"issuerUrl":  kiroIssuerURL,
	})
	if err != nil || regStatus < 200 || regStatus >= 300 {
		return nil, fmt.Errorf("kiro client registration failed (status %d): %v", regStatus, err)
	}
	var clientInfo struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}
	if err := json.Unmarshal(regBody, &clientInfo); err != nil || clientInfo.ClientID == "" {
		return nil, errors.New("kiro client registration returned empty clientId")
	}

	devStatus, devBody, err := m.postRawJSON(ctx, oidcBase+"/device_authorization", map[string]any{
		"clientId":     clientInfo.ClientID,
		"clientSecret": clientInfo.ClientSecret,
		"startUrl":     kiroStartURL,
	})
	if err != nil || devStatus < 200 || devStatus >= 300 {
		return nil, fmt.Errorf("kiro device authorization failed (status %d): %v", devStatus, err)
	}
	var devData struct {
		DeviceCode              string `json:"deviceCode"`
		UserCode                string `json:"userCode"`
		VerificationURI         string `json:"verificationUri"`
		VerificationURIComplete string `json:"verificationUriComplete"`
		ExpiresIn               int    `json:"expiresIn"`
		Interval                int    `json:"interval"`
	}
	if err := json.Unmarshal(devBody, &devData); err != nil || devData.DeviceCode == "" {
		return nil, errors.New("kiro device authorization returned empty deviceCode")
	}

	interval := devData.Interval
	if interval <= 0 {
		interval = defaultDevicePollSec
	}
	ttl := defaultDeviceSessionTTL
	if devData.ExpiresIn > 0 {
		ttl = time.Duration(devData.ExpiresIn) * time.Second
	}
	authURL := devData.VerificationURIComplete
	if authURL == "" {
		authURL = devData.VerificationURI
	}
	return &deviceOAuthSession{
		State:           state,
		ProviderID:      providerIDKiro,
		DeviceCode:      devData.DeviceCode,
		UserCode:        devData.UserCode,
		VerificationURI: devData.VerificationURI,
		AuthorizeURL:    authURL,
		Interval:        interval,
		ClientID:        clientInfo.ClientID,
		ClientSecret:    clientInfo.ClientSecret,
		TokenURL:        tokenURL,
		ConnectionName:  connName,
		BaseURL:         baseURL,
		Models:          models,
		Status:          "pending",
		ExpiresAt:       time.Now().UTC().Add(ttl),
	}, nil
}

func (m *Module) startFreePIDevice(ctx context.Context, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error) {
	devURL := freepiDeviceCodeURL
	consentURL := freepiConsentURL
	tokenURL := freepiTokenURL
	m.mu.RLock()
	override, hasOverride := m.oauthOverrides[providerIDFreePI]
	m.mu.RUnlock()
	if hasOverride && strings.TrimSpace(override.TokenURL) != "" {
		if orig := originOfURL(override.TokenURL); orig != "" {
			devURL = orig + "/auth/github/device"
			consentURL = orig + "/auth/consent"
			tokenURL = strings.TrimSpace(override.TokenURL)
		}
	}

	status, body, err := m.postRawJSON(ctx, devURL, map[string]any{})
	if err != nil || status < 200 || status >= 300 {
		return nil, fmt.Errorf("freepi device code request failed (status %d): %v", status, err)
	}
	var raw struct {
		SessionID               string `json:"session_id"`
		UserCode                string `json:"user_code"`
		VerificationURI         string `json:"verification_uri"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		Interval                int    `json:"interval"`
		ExpiresIn               int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &raw); err != nil || raw.SessionID == "" {
		return nil, errors.New("freepi device code returned empty session_id")
	}

	// Best-effort automatic consent submission (matches xrouter freepi.js).
	_, _, _ = m.postRawJSON(ctx, consentURL, map[string]any{
		"session_id":      raw.SessionID,
		"consent_version": freepiConsentVersion,
	})

	interval := raw.Interval
	if interval <= 0 {
		interval = defaultDevicePollSec
	}
	ttl := defaultDeviceSessionTTL
	if raw.ExpiresIn > 0 {
		ttl = time.Duration(raw.ExpiresIn) * time.Second
	}
	authURL := raw.VerificationURIComplete
	if authURL == "" {
		authURL = raw.VerificationURI
	}
	return &deviceOAuthSession{
		State:           state,
		ProviderID:      providerIDFreePI,
		DeviceCode:      raw.SessionID,
		UserCode:        raw.UserCode,
		VerificationURI: raw.VerificationURI,
		AuthorizeURL:    authURL,
		Interval:        interval,
		TokenURL:        tokenURL,
		ConnectionName:  connName,
		BaseURL:         baseURL,
		Models:          models,
		Status:          "pending",
		ExpiresAt:       time.Now().UTC().Add(ttl),
	}, nil
}

func (m *Module) handleDeviceOAuthStatus(w http.ResponseWriter, r *http.Request, providerID string) {
	state := strings.TrimSpace(r.URL.Query().Get("state"))
	if state == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "state query parameter is required")
		return
	}

	m.deviceOAuth.mu.Lock()
	m.deviceOAuth.sweepLocked(time.Now().UTC())
	sess, ok := m.deviceOAuth.sessions[state]
	if !ok || sess.ProviderID != providerID {
		m.deviceOAuth.mu.Unlock()
		httpserver.WriteJSON(w, http.StatusOK, oauthStatusResponse{Status: "unknown"})
		return
	}
	if sess.Status == "done" {
		m.deviceOAuth.mu.Unlock()
		httpserver.WriteJSON(w, http.StatusOK, oauthStatusResponse{Status: "done"})
		return
	}
	if sess.Status == "error" {
		errMsg := sess.ErrMsg
		delete(m.deviceOAuth.sessions, state)
		m.deviceOAuth.mu.Unlock()
		httpserver.WriteJSON(w, http.StatusOK, oauthStatusResponse{Status: "error", Error: errMsg})
		return
	}
	snapshot := *sess
	m.deviceOAuth.mu.Unlock()

	var (
		res       *deviceOAuthResult
		pending   bool
		nextInt   int
		pollError string
	)
	if strat := m.getOAuthStrategy(providerID); strat.PollDevice != nil {
		res, pending, nextInt, pollError = strat.PollDevice(r.Context(), m, snapshot)
	}

	if pending {
		if nextInt > 0 {
			m.deviceOAuth.mu.Lock()
			if cur, exists := m.deviceOAuth.sessions[state]; exists {
				cur.Interval = nextInt
			}
			m.deviceOAuth.mu.Unlock()
		}
		httpserver.WriteJSON(w, http.StatusOK, oauthStatusResponse{
			Status:   "pending",
			Interval: nextInt,
		})
		return
	}

	if pollError != "" || res == nil {
		if pollError == "" {
			pollError = "device authorization failed"
		}
		m.deviceOAuth.mu.Lock()
		delete(m.deviceOAuth.sessions, state)
		m.deviceOAuth.mu.Unlock()
		httpserver.WriteJSON(w, http.StatusOK, oauthStatusResponse{
			Status: "error",
			Error:  pollError,
		})
		return
	}

	m.deviceOAuth.mu.Lock()
	if cur, exists := m.deviceOAuth.sessions[state]; exists {
		cur.Status = "done"
		cur.Result = res
	}
	m.deviceOAuth.mu.Unlock()

	httpserver.WriteJSON(w, http.StatusOK, oauthStatusResponse{Status: "done"})
}

func (m *Module) pollGitHubCopilotDevice(ctx context.Context, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
	form := url.Values{}
	form.Set("client_id", sess.ClientID)
	form.Set("device_code", sess.DeviceCode)
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	_, body, err := m.postFormJSON(ctx, sess.TokenURL, form)
	if err != nil {
		return nil, true, 0, ""
	}
	var raw struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		TokenType        string `json:"token_type"`
		Scope            string `json:"scope"`
		ExpiresIn        int64  `json:"expires_in"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	_ = json.Unmarshal(body, &raw)
	if raw.AccessToken != "" {
		return &deviceOAuthResult{
			AccessToken:  raw.AccessToken,
			RefreshToken: raw.RefreshToken,
			TokenType:    raw.TokenType,
			Scope:        raw.Scope,
			ExpiresIn:    raw.ExpiresIn,
			ClientID:     sess.ClientID,
			TokenURL:     sess.TokenURL,
		}, false, 0, ""
	}
	switch raw.Error {
	case "", "authorization_pending":
		return nil, true, 0, ""
	case "slow_down":
		return nil, true, sess.Interval + 5, ""
	default:
		msg := raw.ErrorDescription
		if msg == "" {
			msg = raw.Error
		}
		return nil, false, 0, msg
	}
}

func (m *Module) pollKiroDevice(ctx context.Context, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
	_, body, err := m.postRawJSON(ctx, sess.TokenURL, map[string]any{
		"clientId":     sess.ClientID,
		"clientSecret": sess.ClientSecret,
		"deviceCode":   sess.DeviceCode,
		"grantType":    "urn:ietf:params:oauth:grant-type:device_code",
	})
	if err != nil {
		return nil, true, 0, ""
	}
	var raw struct {
		AccessToken      string `json:"accessToken"`
		RefreshToken     string `json:"refreshToken"`
		ExpiresIn        int64  `json:"expiresIn"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		Message          string `json:"message"`
	}
	_ = json.Unmarshal(body, &raw)
	if raw.AccessToken != "" {
		return &deviceOAuthResult{
			AccessToken:  raw.AccessToken,
			RefreshToken: raw.RefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    raw.ExpiresIn,
			ClientID:     sess.ClientID,
			ClientSecret: sess.ClientSecret,
			TokenURL:     sess.TokenURL,
		}, false, 0, ""
	}
	switch raw.Error {
	case "", "authorization_pending":
		return nil, true, 0, ""
	case "slow_down":
		return nil, true, sess.Interval + 5, ""
	default:
		msg := raw.ErrorDescription
		if msg == "" {
			msg = raw.Message
		}
		if msg == "" {
			msg = raw.Error
		}
		return nil, false, 0, msg
	}
}

func (m *Module) pollFreePIDevice(ctx context.Context, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
	_, body, err := m.postRawJSON(ctx, sess.TokenURL, map[string]any{
		"session_id": sess.DeviceCode,
		"timezone":   "UTC",
		"locale":     "en-US",
	})
	if err != nil {
		return nil, true, 0, ""
	}
	var raw struct {
		Status           string `json:"status"`
		Token            string `json:"token"`
		AccessToken      string `json:"access_token"`
		AccessTokenCamel string `json:"accessToken"`
		Handle           string `json:"handle"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	_ = json.Unmarshal(body, &raw)
	tok := raw.Token
	if tok == "" {
		tok = raw.AccessToken
	}
	if tok == "" {
		tok = raw.AccessTokenCamel
	}
	if tok != "" {
		label := ""
		if raw.Handle != "" {
			label = "FreePI (@" + raw.Handle + ")"
		}
		return &deviceOAuthResult{
			AccessToken:  tok,
			TokenType:    "Bearer",
			TokenURL:     sess.TokenURL,
			DisplayLabel: label,
		}, false, 0, ""
	}
	if raw.Status == "pending" || raw.Error == "authorization_pending" || raw.Error == "" {
		return nil, true, 0, ""
	}
	msg := raw.ErrorDescription
	if msg == "" {
		msg = raw.Error
	}
	return nil, false, 0, msg
}

func (m *Module) handleDeviceOAuthExchange(w http.ResponseWriter, r *http.Request, providerID string, adapter ProviderAdapter, req oauthExchangeRequest) {
	state := strings.TrimSpace(req.State)
	if state == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "oauth state is required")
		return
	}

	var validatedConnName string
	if strings.TrimSpace(req.ConnectionName) != "" {
		var err error
		validatedConnName, err = validateName(req.ConnectionName)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}

	priority := defaultPriorityValue
	if req.Priority != nil {
		if *req.Priority < 0 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "priority must be non-negative")
			return
		}
		priority = *req.Priority
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	headers, err := validateHeaders(req.Headers)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	m.deviceOAuth.mu.Lock()
	m.deviceOAuth.sweepLocked(time.Now().UTC())
	sess, ok := m.deviceOAuth.sessions[state]
	if !ok || sess.ProviderID != providerID || sess.Status != "done" || sess.Result == nil {
		m.deviceOAuth.mu.Unlock()
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_oauth_state", "device oauth session not completed yet")
		return
	}
	res := sess.Result
	sessConnName := sess.ConnectionName
	sessBaseURL := sess.BaseURL
	sessModels := cloneStringSlice(sess.Models)
	delete(m.deviceOAuth.sessions, state)
	m.deviceOAuth.mu.Unlock()

	var postRes oauthPostExchangeResult
	if strat := m.getOAuthStrategy(providerID); strat.PostExchange != nil {
		postRes = strat.PostExchange(r.Context(), m, res.TokenURL, oauthTokenResponse{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			TokenType:    res.TokenType,
			Scope:        res.Scope,
			ExpiresIn:    res.ExpiresIn,
		})
	}
	for k, v := range postRes.ExtraHeaders {
		if _, exists := headers[k]; !exists && strings.TrimSpace(v) != "" {
			headers[k] = v
		}
	}

	stored := storedCredentials{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		TokenType:    res.TokenType,
		Scope:        res.Scope,
		ClientID:     res.ClientID,
		ClientSecret: res.ClientSecret,
		TokenURL:     res.TokenURL,
		Metadata:     postRes.Metadata,
	}
	sealed, err := m.sealCredentials(stored)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to encrypt OAuth tokens")
		return
	}

	var tokenExpiresAt *string
	if res.ExpiresIn > 0 {
		exp := time.Now().UTC().Add(time.Duration(res.ExpiresIn) * time.Second).Format(rfc3339NanoFixed)
		tokenExpiresAt = &exp
	}

	connName := validatedConnName
	if connName == "" {
		connName = sessConnName
	}
	if connName == "" {
		connName = res.DisplayLabel
	}
	if connName == "" {
		connName = postRes.DisplayLabel
	}
	if connName == "" {
		connName = adapter.DisplayName()
	}

	connID, err := generateID(defaultConnectionPrefix)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate connection identifier")
		return
	}

	summary, err := m.repo.CreateConnection(r.Context(), createConnectionParams{
		ID:             connID,
		ProviderID:     providerID,
		Name:           connName,
		Protocol:       adapter.Protocol(),
		BaseURL:        sessBaseURL,
		AuthType:       authTypeOAuth,
		CredentialsEnc: sealed,
		Priority:       priority,
		Enabled:        enabled,
		Models:         sessModels,
		Headers:        headers,
		TokenExpiresAt: tokenExpiresAt,
	})
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to persist OAuth connection")
		return
	}

	m.publishChanged("created", "connection", summary.ID, summary.ProviderID, actorFromContext(r.Context()), summary.UpdatedAt)
	httpserver.WriteJSON(w, http.StatusCreated, summary)
}

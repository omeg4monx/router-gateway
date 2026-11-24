package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// rfc8628DeviceSpec declaratively configures a standard RFC 8628 Device Code
// OAuth flow (with optional PKCE S256 for providers like MiniMax Coding).
type rfc8628DeviceSpec struct {
	ProviderID              string
	DefaultDeviceCodeURL    string
	DefaultTokenURL         string
	DefaultClientID         string
	DefaultScope            string
	Audience                string
	UsePKCE                 bool
	ExtraDeviceParams       map[string]string
	FallbackVerificationURI string
	DisplayPrefix           string
}

func (m *Module) makeRFC8628StartFn(spec rfc8628DeviceSpec) func(context.Context, *Module, string, OAuthConfig, string, string, []string) (*deviceOAuthSession, error) {
	return func(ctx context.Context, mod *Module, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error) {
		return mod.startRFC8628Device(ctx, spec, state, cfg, connName, baseURL, models)
	}
}

func (m *Module) makeRFC8628PollFn(spec rfc8628DeviceSpec) func(context.Context, *Module, deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
	return func(ctx context.Context, mod *Module, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
		return mod.pollRFC8628Device(ctx, spec, sess)
	}
}

func (m *Module) startRFC8628Device(
	ctx context.Context,
	spec rfc8628DeviceSpec,
	state string,
	cfg OAuthConfig,
	connName, baseURL string,
	models []string,
) (*deviceOAuthSession, error) {
	deviceURL := strings.TrimSpace(cfg.AuthorizeURL)
	if deviceURL == "" {
		deviceURL = spec.DefaultDeviceCodeURL
	}
	tokenURL := strings.TrimSpace(cfg.TokenURL)
	if tokenURL == "" {
		tokenURL = spec.DefaultTokenURL
	}
	clientID := strings.TrimSpace(cfg.ClientID)
	if clientID == "" {
		clientID = spec.DefaultClientID
	}
	scope := spec.DefaultScope
	if len(cfg.Scopes) > 0 {
		scope = strings.Join(cfg.Scopes, " ")
	}

	form := url.Values{}
	if clientID != "" {
		form.Set("client_id", clientID)
	}
	if scope != "" {
		form.Set("scope", scope)
	}
	if spec.Audience != "" {
		form.Set("audience", spec.Audience)
	}
	for k, v := range spec.ExtraDeviceParams {
		form.Set(k, v)
	}

	var verifier string
	if spec.UsePKCE {
		var challenge string
		var err error
		verifier, challenge, err = generatePKCE()
		if err != nil {
			return nil, fmt.Errorf("%s pkce generation failed: %w", spec.ProviderID, err)
		}
		form.Set("code_challenge", challenge)
		form.Set("code_challenge_method", "S256")
	}

	status, body, err := m.postFormJSON(ctx, deviceURL, form)
	if err != nil {
		return nil, fmt.Errorf("%s device code request failed: %w", spec.ProviderID, err)
	}
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("%s device code returned HTTP %d", spec.ProviderID, status)
	}

	var raw struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURI         string `json:"verification_uri"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		ExpiresIn               int    `json:"expires_in"`
		Interval                int    `json:"interval"`
		Error                   string `json:"error"`
		ErrorDescription        string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decode %s device code response: %w", spec.ProviderID, err)
	}
	if raw.Error != "" || raw.DeviceCode == "" {
		msg := raw.ErrorDescription
		if msg == "" {
			msg = raw.Error
		}
		return nil, fmt.Errorf("%s device code error: %s", spec.ProviderID, msg)
	}

	interval := raw.Interval
	if interval <= 0 {
		interval = defaultDevicePollSec
	}
	ttl := defaultDeviceSessionTTL
	if raw.ExpiresIn > 0 {
		ttl = time.Duration(raw.ExpiresIn) * time.Second
	}
	verifyURI := strings.TrimSpace(raw.VerificationURI)
	if verifyURI == "" {
		verifyURI = spec.FallbackVerificationURI
	}
	authURL := strings.TrimSpace(raw.VerificationURIComplete)
	if authURL == "" {
		authURL = verifyURI
	}

	return &deviceOAuthSession{
		State:           state,
		ProviderID:      spec.ProviderID,
		DeviceCode:      raw.DeviceCode,
		UserCode:        raw.UserCode,
		VerificationURI: verifyURI,
		AuthorizeURL:    authURL,
		Interval:        interval,
		ClientID:        clientID,
		TokenURL:        tokenURL,
		CodeVerifier:    verifier,
		ConnectionName:  connName,
		BaseURL:         baseURL,
		Models:          models,
		Status:          "pending",
		ExpiresAt:       time.Now().UTC().Add(ttl),
	}, nil
}

func (m *Module) pollRFC8628Device(
	ctx context.Context,
	spec rfc8628DeviceSpec,
	sess deviceOAuthSession,
) (*deviceOAuthResult, bool, int, string) {
	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	form.Set("device_code", sess.DeviceCode)
	if sess.ClientID != "" {
		form.Set("client_id", sess.ClientID)
	}
	if spec.UsePKCE && sess.CodeVerifier != "" {
		form.Set("code_verifier", sess.CodeVerifier)
	}

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

	if strings.TrimSpace(raw.AccessToken) != "" {
		label := ""
		if email := extractEmailFromJWT(raw.AccessToken); email != "" && spec.DisplayPrefix != "" {
			label = spec.DisplayPrefix + " (" + email + ")"
		}
		return &deviceOAuthResult{
			AccessToken:  raw.AccessToken,
			RefreshToken: raw.RefreshToken,
			TokenType:    raw.TokenType,
			Scope:        raw.Scope,
			ExpiresIn:    raw.ExpiresIn,
			ClientID:     sess.ClientID,
			TokenURL:     sess.TokenURL,
			DisplayLabel: label,
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

// KiloCode custom REST device-auth flow (adopted from xrouter/src/lib/oauth/providers/kilocode.js).
func (m *Module) startKiloCodeDevice(
	ctx context.Context,
	state string,
	cfg OAuthConfig,
	connName, baseURL string,
	models []string,
) (*deviceOAuthSession, error) {
	initURL := strings.TrimSpace(cfg.AuthorizeURL)
	if initURL == "" {
		initURL = "https://api.kilo.ai/api/device-auth/codes"
	}
	pollBase := strings.TrimRight(strings.TrimSpace(cfg.TokenURL), "/")
	if pollBase == "" {
		pollBase = "https://api.kilo.ai/api/device-auth/codes"
	}
	status, body, err := m.postRawJSON(ctx, initURL, map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("kilocode device auth failed: %w", err)
	}
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("kilocode device auth returned HTTP %d", status)
	}
	var raw struct {
		Code            string `json:"code"`
		VerificationURL string `json:"verificationUrl"`
		ExpiresIn       int    `json:"expiresIn"`
	}
	if err := json.Unmarshal(body, &raw); err != nil || raw.Code == "" {
		return nil, fmt.Errorf("kilocode device auth returned invalid response")
	}
	ttl := defaultDeviceSessionTTL
	if raw.ExpiresIn > 0 {
		ttl = time.Duration(raw.ExpiresIn) * time.Second
	}
	return &deviceOAuthSession{
		State:           state,
		ProviderID:      "kilocode",
		DeviceCode:      raw.Code,
		UserCode:        raw.Code,
		VerificationURI: raw.VerificationURL,
		AuthorizeURL:    raw.VerificationURL,
		Interval:        3,
		TokenURL:        pollBase + "/" + url.PathEscape(raw.Code),
		ConnectionName:  connName,
		BaseURL:         baseURL,
		Models:          models,
		Status:          "pending",
		ExpiresAt:       time.Now().UTC().Add(ttl),
	}, nil
}

func (m *Module) pollKiloCodeDevice(ctx context.Context, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sess.TokenURL, http.NoBody)
	if err != nil {
		return nil, true, 0, ""
	}
	req.Header.Set("Accept", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, true, 0, ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxOAuthResponseBytes))
	if resp.StatusCode == http.StatusAccepted {
		return nil, true, 0, ""
	}
	var raw struct {
		Status    string `json:"status"`
		Token     string `json:"token"`
		UserEmail string `json:"userEmail"`
		Error     string `json:"error"`
	}
	_ = json.Unmarshal(body, &raw)
	if raw.Token != "" {
		label := ""
		if raw.UserEmail != "" {
			label = "KiloCode (" + raw.UserEmail + ")"
		}
		return &deviceOAuthResult{
			AccessToken:  raw.Token,
			TokenType:    "Bearer",
			TokenURL:     sess.TokenURL,
			DisplayLabel: label,
		}, false, 0, ""
	}
	if raw.Status == "pending" || resp.StatusCode == http.StatusNotFound {
		return nil, true, 0, ""
	}
	msg := raw.Error
	if msg == "" {
		msg = fmt.Sprintf("kilocode poll returned HTTP %d", resp.StatusCode)
	}
	return nil, false, 0, msg
}

func postExchangeNous(ctx context.Context, m *Module, tokenURL string, tok oauthTokenResponse) oauthPostExchangeResult {
	if !isHostMatch(tokenURL, "nousresearch.com") || tok.AccessToken == "" {
		return oauthPostExchangeResult{}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://portal.nousresearch.com/api/oauth/agent-key", http.NoBody)
	if err != nil {
		return oauthPostExchangeResult{}
	}
	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	req.Header.Set("Accept", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return oauthPostExchangeResult{}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return oauthPostExchangeResult{}
	}
	var payload struct {
		AgentKey string `json:"agent_key"`
		APIKey   string `json:"api_key"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return oauthPostExchangeResult{}
	}
	meta := map[string]any{}
	if key := strings.TrimSpace(payload.AgentKey); key != "" {
		meta["agent_key"] = key
	} else if key := strings.TrimSpace(payload.APIKey); key != "" {
		meta["agent_key"] = key
	}
	label := ""
	if email := strings.TrimSpace(payload.Email); email != "" {
		meta["email"] = email
		label = "Nous (" + email + ")"
	}
	return oauthPostExchangeResult{
		Email:        strings.TrimSpace(payload.Email),
		DisplayLabel: label,
		Metadata:     meta,
	}
}

func postExchangeIFlow(ctx context.Context, m *Module, tokenURL string, tok oauthTokenResponse) oauthPostExchangeResult {
	if !isHostMatch(tokenURL, "iflow.cn") || tok.AccessToken == "" {
		return oauthPostExchangeResult{}
	}
	endpoint := "https://iflow.cn/api/oauth/getUserInfo?accessToken=" + url.QueryEscape(tok.AccessToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return oauthPostExchangeResult{}
	}
	req.Header.Set("Accept", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return oauthPostExchangeResult{}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return oauthPostExchangeResult{}
	}
	var payload struct {
		Data struct {
			APIKey   string `json:"apiKey"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return oauthPostExchangeResult{}
	}
	meta := map[string]any{}
	if payload.Data.APIKey != "" {
		meta["api_key"] = payload.Data.APIKey
	}
	who := strings.TrimSpace(payload.Data.Email)
	if who == "" {
		who = strings.TrimSpace(payload.Data.Username)
	}
	label := ""
	if who != "" {
		meta["email"] = who
		label = "iFlow (" + who + ")"
	}
	return oauthPostExchangeResult{
		Email:        who,
		DisplayLabel: label,
		Metadata:     meta,
	}
}

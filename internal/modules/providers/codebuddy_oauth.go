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
	"time"
)

const (
	providerIDCodeBuddyCN   = "codebuddy-cn"
	providerIDCodeBuddyIntl = "codebuddy-intl"

	codebuddyCNOrigin           = "https://copilot.tencent.com"
	codebuddyCNDefaultBaseURL   = codebuddyCNOrigin + "/v2"
	codebuddyCNStateURL         = codebuddyCNOrigin + "/v2/plugin/auth/state"
	codebuddyCNTokenURL         = codebuddyCNOrigin + "/v2/plugin/auth/token"
	codebuddyCNRefreshURL       = codebuddyCNOrigin + "/v2/plugin/auth/token/refresh"
	codebuddyCNHost             = "copilot.tencent.com"
	codebuddyCNAuthUserAgent    = "CLI/2.63.2 CodeBuddy/2.63.2"
	codebuddyCNPlatform         = "CLI"
	codebuddyIntlOrigin         = "https://www.codebuddy.ai"
	codebuddyIntlDefaultBaseURL = codebuddyIntlOrigin + "/v2"
	codebuddyIntlStateURL       = codebuddyIntlOrigin + "/v2/plugin/auth/state"
	codebuddyIntlTokenURL       = codebuddyIntlOrigin + "/v2/plugin/auth/token"
	codebuddyIntlRefreshURL     = codebuddyIntlOrigin + "/v2/plugin/auth/token/refresh"
	codebuddyIntlHost           = "www.codebuddy.ai"
	codebuddyIntlAuthUserAgent  = "IDE/2.63.2 CodeBuddy/2.63.2"
	codebuddyIntlPlatform       = "ide"

	codebuddyQuotaPath     = "/v2/billing/meter/get-user-resource"
	codebuddyCLIUserAgent  = "CLI/2.108.1 CodeBuddy/2.108.1"
	codebuddyPendingCode   = 11217
	codebuddyOAuthTimeout  = 15 * time.Second
	codebuddyDefaultExpiry = 86400
)

type codebuddyAuthSpec struct {
	providerID string
	domain     string
	stateURL   string
	tokenURL   string
	refreshURL string
	userAgent  string
	platform   string
	label      string
}

func codebuddySpecFor(providerID string) codebuddyAuthSpec {
	if providerID == providerIDCodeBuddyIntl {
		return codebuddyAuthSpec{
			providerID: providerIDCodeBuddyIntl,
			domain:     codebuddyIntlHost,
			stateURL:   codebuddyIntlStateURL,
			tokenURL:   codebuddyIntlTokenURL,
			refreshURL: codebuddyIntlRefreshURL,
			userAgent:  codebuddyIntlAuthUserAgent,
			platform:   codebuddyIntlPlatform,
			label:      "CodeBuddy Intl",
		}
	}
	return codebuddyAuthSpec{
		providerID: providerIDCodeBuddyCN,
		domain:     codebuddyCNHost,
		stateURL:   codebuddyCNStateURL,
		tokenURL:   codebuddyCNTokenURL,
		refreshURL: codebuddyCNRefreshURL,
		userAgent:  codebuddyCNAuthUserAgent,
		platform:   codebuddyCNPlatform,
		label:      "CodeBuddy CN",
	}
}

func newCodeBuddyOAuthStrategy(providerID string) *oauthProviderStrategy {
	return &oauthProviderStrategy{
		ID:               providerID,
		FlowType:         oauthFlowDeviceCode,
		ClientIDOptional: true,
		StartDevice: func(ctx context.Context, m *Module, state string, cfg OAuthConfig, connName, baseURL string, models []string) (*deviceOAuthSession, error) {
			return m.startCodeBuddyDevice(ctx, providerID, state, cfg, connName, baseURL, models)
		},
		PollDevice: func(ctx context.Context, m *Module, sess deviceOAuthSession) (*deviceOAuthResult, bool, int, string) {
			return m.pollCodeBuddyDevice(ctx, providerID, sess)
		},
		RefreshToken: func(ctx context.Context, m *Module, tokenURL, _ string, _ string, stored storedCredentials) (oauthTokenResponse, error) {
			return refreshCodeBuddyToken(ctx, m, providerID, tokenURL, stored)
		},
	}
}

func (m *Module) startCodeBuddyDevice(
	ctx context.Context,
	providerID, state string,
	_ OAuthConfig,
	connName, baseURL string,
	models []string,
) (*deviceOAuthSession, error) {
	spec := codebuddySpecFor(providerID)
	stateURL := spec.stateURL
	tokenURL := spec.tokenURL

	m.mu.RLock()
	override, hasOverride := m.oauthOverrides[providerID]
	m.mu.RUnlock()
	if hasOverride {
		if strings.TrimSpace(override.AuthorizeURL) != "" {
			stateURL = strings.TrimSpace(override.AuthorizeURL)
		}
		if strings.TrimSpace(override.TokenURL) != "" {
			tokenURL = strings.TrimSpace(override.TokenURL)
		}
	}

	reqURL := stateURL
	if strings.Contains(reqURL, "?") {
		reqURL += "&platform=" + url.QueryEscape(spec.platform)
	} else {
		reqURL += "?platform=" + url.QueryEscape(spec.platform)
	}

	validURL, err := validateHTTPURL(reqURL)
	if err != nil {
		return nil, fmt.Errorf("%s invalid state URL: %w", spec.label, err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, codebuddyOAuthTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, validURL, bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, fmt.Errorf("%s build state request: %w", spec.label, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", spec.userAgent)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-Domain", spec.domain)
	req.Header.Set("X-No-Authorization", "true")
	req.Header.Set("X-No-User-Id", "true")
	req.Header.Set("X-Product", "SaaS")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s state request failed: %w", spec.label, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("%s read state response: %w", spec.label, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s state request failed (status %d)", spec.label, resp.StatusCode)
	}

	var raw struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			State   string `json:"state"`
			AuthURL string `json:"authUrl"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("%s decode state response: %w", spec.label, err)
	}
	if raw.Code != 0 || strings.TrimSpace(raw.Data.State) == "" || strings.TrimSpace(raw.Data.AuthURL) == "" {
		msg := strings.TrimSpace(raw.Msg)
		if msg == "" {
			msg = "missing state/authUrl"
		}
		return nil, fmt.Errorf("%s state error: %s", spec.label, msg)
	}

	return &deviceOAuthSession{
		State:           state,
		ProviderID:      providerID,
		DeviceCode:      strings.TrimSpace(raw.Data.State),
		UserCode:        "",
		VerificationURI: strings.TrimSpace(raw.Data.AuthURL),
		AuthorizeURL:    strings.TrimSpace(raw.Data.AuthURL),
		Interval:        defaultDevicePollSec,
		TokenURL:        tokenURL,
		ConnectionName:  connName,
		BaseURL:         baseURL,
		Models:          models,
		Status:          "pending",
		ExpiresAt:       time.Now().UTC().Add(defaultDeviceSessionTTL),
	}, nil
}

func (m *Module) pollCodeBuddyDevice(
	ctx context.Context,
	providerID string,
	sess deviceOAuthSession,
) (result *deviceOAuthResult, pending bool, retryAfter int, errMsg string) {
	spec := codebuddySpecFor(providerID)
	tokenURL := strings.TrimSpace(sess.TokenURL)
	if tokenURL == "" {
		tokenURL = spec.tokenURL
	}
	pollURL := tokenURL
	if strings.Contains(pollURL, "?") {
		pollURL += "&state=" + url.QueryEscape(sess.DeviceCode)
	} else {
		pollURL += "?state=" + url.QueryEscape(sess.DeviceCode)
	}

	validURL, err := validateHTTPURL(pollURL)
	if err != nil {
		return nil, false, 0, fmt.Sprintf("%s invalid token poll URL: %v", spec.label, err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, codebuddyOAuthTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, validURL, http.NoBody)
	if err != nil {
		return nil, false, 0, err.Error()
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", spec.userAgent)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-Domain", spec.domain)
	req.Header.Set("X-No-Authorization", "true")
	req.Header.Set("X-No-User-Id", "true")
	req.Header.Set("X-No-Enterprise-Id", "true")
	req.Header.Set("X-No-Department-Info", "true")
	req.Header.Set("X-Product", "SaaS")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, true, 0, ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, true, 0, ""
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, false, 0, fmt.Sprintf("%s token poll failed (status %d)", spec.label, resp.StatusCode)
	}

	var raw struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
			TokenType    string `json:"tokenType"`
			ExpiresIn    int64  `json:"expiresIn"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, true, 0, ""
	}

	if raw.Code == 0 && strings.TrimSpace(raw.Data.AccessToken) != "" {
		tokenType := strings.TrimSpace(raw.Data.TokenType)
		if tokenType == "" {
			tokenType = "Bearer"
		}
		expiresIn := raw.Data.ExpiresIn
		if expiresIn <= 0 {
			expiresIn = codebuddyDefaultExpiry
		}
		return &deviceOAuthResult{
			AccessToken:  strings.TrimSpace(raw.Data.AccessToken),
			RefreshToken: strings.TrimSpace(raw.Data.RefreshToken),
			TokenType:    tokenType,
			ExpiresIn:    expiresIn,
			TokenURL:     spec.refreshURL,
			DisplayLabel: spec.label,
		}, false, 0, ""
	}

	if raw.Code == codebuddyPendingCode {
		return nil, true, 0, ""
	}

	msg := strings.TrimSpace(raw.Msg)
	if msg == "" {
		msg = fmt.Sprintf("%s authorization failed (code %d)", spec.label, raw.Code)
	}
	return nil, false, 0, msg
}

func refreshCodeBuddyToken(
	ctx context.Context,
	m *Module,
	providerID, tokenURL string,
	stored storedCredentials,
) (oauthTokenResponse, error) {
	if strings.TrimSpace(stored.RefreshToken) == "" {
		return oauthTokenResponse{}, errors.New("missing CodeBuddy refresh_token")
	}

	spec := codebuddySpecFor(providerID)
	refreshURL := spec.refreshURL
	trimmed := strings.TrimSpace(tokenURL)
	if trimmed != "" {
		if strings.HasSuffix(trimmed, "/refresh") {
			refreshURL = trimmed
		} else if orig := originOfURL(trimmed); orig != "" {
			refreshURL = orig + "/v2/plugin/auth/token/refresh"
		}
	}

	validURL, err := validateHTTPURL(refreshURL)
	if err != nil {
		return oauthTokenResponse{}, fmt.Errorf("%s invalid refresh URL: %w", spec.label, err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, codebuddyOAuthTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, validURL, bytes.NewReader([]byte("{}")))
	if err != nil {
		return oauthTokenResponse{}, fmt.Errorf("%s build refresh request: %w", spec.label, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", spec.userAgent)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-Domain", spec.domain)
	req.Header.Set("X-Refresh-Token", stored.RefreshToken)
	req.Header.Set("X-Auth-Refresh-Source", "plugin")
	req.Header.Set("X-Product", "SaaS")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return oauthTokenResponse{}, fmt.Errorf("%s token refresh request failed: %w", spec.label, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return oauthTokenResponse{}, fmt.Errorf("%s read refresh response: %w", spec.label, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return oauthTokenResponse{}, fmt.Errorf("%s token refresh failed (status %d)", spec.label, resp.StatusCode)
	}

	var raw struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
			TokenType    string `json:"tokenType"`
			ExpiresIn    int64  `json:"expiresIn"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return oauthTokenResponse{}, fmt.Errorf("%s decode refresh response: %w", spec.label, err)
	}
	if raw.Code != 0 || strings.TrimSpace(raw.Data.AccessToken) == "" {
		msg := strings.TrimSpace(raw.Msg)
		if msg == "" {
			msg = fmt.Sprintf("code %d", raw.Code)
		}
		return oauthTokenResponse{}, fmt.Errorf("%s token refresh rejected: %s", spec.label, msg)
	}

	nextRefresh := strings.TrimSpace(raw.Data.RefreshToken)
	if nextRefresh == "" {
		nextRefresh = stored.RefreshToken
	}
	expiresIn := raw.Data.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = codebuddyDefaultExpiry
	}
	tokenType := strings.TrimSpace(raw.Data.TokenType)
	if tokenType == "" {
		tokenType = "Bearer"
	}

	return oauthTokenResponse{
		AccessToken:  strings.TrimSpace(raw.Data.AccessToken),
		RefreshToken: nextRefresh,
		TokenType:    tokenType,
		ExpiresIn:    expiresIn,
	}, nil
}

func applyCodeBuddyProbeHeaders(req *http.Request, providerID string) {
	spec := codebuddySpecFor(providerID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", codebuddyCLIUserAgent)
	req.Header.Set("X-Product", "SaaS")
	req.Header.Set("X-IDE-Type", "CLI")
	req.Header.Set("X-IDE-Name", "CLI")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("x-codebuddy-request", "1")
	req.Header.Set("X-Domain", spec.domain)
}

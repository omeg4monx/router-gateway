package providers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

const (
	providerIDAWSQuick         = "aws-quick"
	awsQuickDiscoveryBaseURL   = "https://quick.aws.com/sn"
	awsQuickSocialDataPlaneURL = "https://cell001.dp.appintegrations.us-east-1.prod.plato.ai.aws.dev"
	awsQuickScopeProd          = "quicksight:read_write"
	awsQuickAppVersion         = "1.2.4660"
	awsQuickClientName         = "AmazonQuick-Desktop"
	awsQuickClientType         = "public"
	awsQuickSessionTTL         = 10 * time.Minute
	awsQuickDiscoveryMaxAge    = 15 * time.Minute
)

var awsRegionPattern = regexp.MustCompile(`^[a-z]{2}-[a-z]+-\d{1,2}$`)

type awsQuickDiscovery struct {
	IssuerURL     string
	TenantURL     string
	Account       string
	IDCInstanceID string
	IDCRegion     string
}

type awsQuickUserInfo struct {
	Email        string
	UserName     string
	DisplayName  string
	UserARN      string
	AWSAccountID string
	AccountName  string
}

type awsQuickTokens struct {
	AccessToken  string
	RefreshToken string
	IDToken      string
	ExpiresIn    int64
}

type awsQuickResult struct {
	Tokens       awsQuickTokens
	ClientID     string
	ClientSecret string
	Discovery    awsQuickDiscovery
	UserInfo     awsQuickUserInfo
	TokenURL     string
}

type awsQuickSession struct {
	State            string
	Nonce            string
	Verifier         string
	Challenge        string
	Port             int
	RedirectURI      string
	ConnectionName   string
	BaseURL          string
	Models           []string
	Discovery        *awsQuickDiscovery
	DiscoverySavedAt time.Time
	AuthorizeURL     string
	ClientID         string
	ClientSecret     string
	Status           string // "pending" | "discovery_saved" | "done" | "error"
	Result           *awsQuickResult
	ErrMsg           string
	CreatedAt        time.Time
	ExpiresAt        time.Time
}

type awsQuickRuntime struct {
	mu        sync.Mutex
	listener  net.Listener
	server    *http.Server
	port      int
	idleTimer *time.Timer
	sessions  map[string]*awsQuickSession
}

func newAwsQuickRuntime() *awsQuickRuntime {
	return &awsQuickRuntime{
		sessions: make(map[string]*awsQuickSession),
	}
}

func (rt *awsQuickRuntime) sweepLocked(now time.Time) {
	for st, sess := range rt.sessions {
		if !sess.ExpiresAt.IsZero() && !now.Before(sess.ExpiresAt) {
			delete(rt.sessions, st)
		}
	}
}

func (rt *awsQuickRuntime) findByNonceLocked(nonce string) *awsQuickSession {
	if nonce == "" {
		return nil
	}
	for _, sess := range rt.sessions {
		if sess.Nonce == nonce {
			return sess
		}
	}
	return nil
}

func generateAwsQuickNonce() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func validateAwsRegion(region string) (string, error) {
	trimmed := strings.TrimSpace(region)
	if trimmed == "" {
		return "us-east-1", nil
	}
	if !awsRegionPattern.MatchString(trimmed) {
		return "", fmt.Errorf("invalid AWS region %q", trimmed)
	}
	return trimmed, nil
}

func parseAwsQuickDiscovery(q url.Values) awsQuickDiscovery {
	first := func(keys ...string) string {
		for _, k := range keys {
			if v := strings.TrimSpace(q.Get(k)); v != "" {
				return v
			}
		}
		return ""
	}
	return awsQuickDiscovery{
		IssuerURL:     first("issuerUrl", "issuer_url", "issuer"),
		TenantURL:     first("tenantUrl", "tenant_url"),
		Account:       first("account"),
		IDCInstanceID: first("idcInstanceId", "idc_instance_id"),
		IDCRegion:     first("idcRegion", "idc_region"),
	}
}

func parseAwsQuickPastedQuery(pasted string) (url.Values, bool) {
	s := strings.TrimSpace(pasted)
	if s == "" {
		return nil, false
	}
	if idx := strings.IndexByte(s, '#'); idx >= 0 {
		s = s[:idx]
	}
	var rawQuery string
	if qi := strings.IndexByte(s, '?'); qi >= 0 {
		rawQuery = s[qi+1:]
	} else if strings.Contains(s, "=") && !strings.ContainsAny(s, " \t\r\n/") {
		rawQuery = s
	} else {
		return nil, false
	}
	vals, err := url.ParseQuery(rawQuery)
	if err != nil || len(vals) == 0 {
		return nil, false
	}
	return vals, true
}

func buildAwsQuickSigninURL(port int, nonce string) string {
	surface := url.QueryEscape(fmt.Sprintf("desktop:%s:%d:%s", nonce, port, awsQuickAppVersion))
	return fmt.Sprintf("%s?qwSurface=%s", awsQuickDiscoveryBaseURL, surface)
}

func (m *Module) awsQuickEndpoints(idcRegion string) (registerURL, authorizeBase, tokenURL, userInfoOverride string) {
	m.mu.RLock()
	override, hasOverride := m.oauthOverrides[providerIDAWSQuick]
	m.mu.RUnlock()

	if hasOverride && strings.TrimSpace(override.TokenURL) != "" {
		tokURL := strings.TrimSpace(override.TokenURL)
		if parsed, err := url.Parse(tokURL); err == nil && parsed.Scheme != "" && parsed.Host != "" {
			origin := parsed.Scheme + "://" + parsed.Host
			registerURL = origin + "/client/register"
			authorizeBase = origin + "/authorize"
			tokenURL = tokURL
			if override.ExtraTokenParams != nil {
				userInfoOverride = strings.TrimSpace(override.ExtraTokenParams["user_info_url"])
			}
			return registerURL, authorizeBase, tokenURL, userInfoOverride
		}
	}

	base := fmt.Sprintf("https://oidc.%s.amazonaws.com", idcRegion)
	return base + "/client/register", base + "/authorize", base + "/token", ""
}

func (m *Module) postAwsQuickJSON(ctx context.Context, targetURL string, payload any, extraHeaders map[string]string) ([]byte, error) {
	validURL, err := validateHTTPURL(targetURL)
	if err != nil {
		return nil, fmt.Errorf("providers: invalid aws-quick endpoint URL: %w", err)
	}
	rawBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("providers: marshal aws-quick request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, validURL, bytes.NewReader(rawBody))
	if err != nil {
		return nil, fmt.Errorf("providers: build aws-quick request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("providers: aws-quick request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxOAuthResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("providers: read aws-quick response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return nil, fmt.Errorf("providers: aws-quick endpoint returned HTTP %d: %s", resp.StatusCode, snippet)
	}
	return body, nil
}

func (m *Module) registerAwsQuickOIDCClient(ctx context.Context, discovery awsQuickDiscovery, redirectURI string) (clientID, clientSecret, idcRegion string, err error) {
	idcRegion, err = validateAwsRegion(discovery.IDCRegion)
	if err != nil {
		return "", "", "", err
	}
	registerURL, _, _, _ := m.awsQuickEndpoints(idcRegion)
	bareRedirect := strings.TrimRight(strings.TrimSpace(redirectURI), "/")

	body, err := m.postAwsQuickJSON(ctx, registerURL, map[string]any{
		"clientName":   awsQuickClientName,
		"clientType":   awsQuickClientType,
		"scopes":       []string{awsQuickScopeProd},
		"grantTypes":   []string{"authorization_code", "refresh_token"},
		"redirectUris": []string{bareRedirect},
		"issuerUrl":    discovery.IssuerURL,
	}, nil)
	if err != nil {
		return "", "", "", err
	}

	var parsed struct {
		ClientID       string `json:"clientId"`
		ClientIDSnake  string `json:"client_id"`
		ClientSecret   string `json:"clientSecret"`
		ClientSecSnake string `json:"client_secret"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", "", "", fmt.Errorf("providers: decode aws-quick client register response: %w", err)
	}
	clientID = strings.TrimSpace(parsed.ClientID)
	if clientID == "" {
		clientID = strings.TrimSpace(parsed.ClientIDSnake)
	}
	clientSecret = strings.TrimSpace(parsed.ClientSecret)
	if clientSecret == "" {
		clientSecret = strings.TrimSpace(parsed.ClientSecSnake)
	}
	if clientID == "" {
		return "", "", "", errors.New("providers: aws-quick client registration returned empty clientId")
	}
	return clientID, clientSecret, idcRegion, nil
}

func (m *Module) buildAwsQuickAuthorizeURL(clientID, redirectURI, state, challenge, idcRegion string) string {
	_, authorizeBase, _, _ := m.awsQuickEndpoints(idcRegion)
	parsed, err := url.Parse(authorizeBase)
	if err != nil {
		return authorizeBase
	}
	q := parsed.Query()
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("redirect_uri", strings.TrimRight(redirectURI, "/"))
	q.Set("scope", awsQuickScopeProd)
	q.Set("state", state)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", "S256")
	parsed.RawQuery = q.Encode()
	return parsed.String()
}

func (m *Module) postAwsQuickTokenJSON(ctx context.Context, tokenURL string, payload map[string]any) (awsQuickTokens, error) {
	body, err := m.postAwsQuickJSON(ctx, tokenURL, payload, nil)
	if err != nil {
		return awsQuickTokens{}, err
	}
	var raw struct {
		AccessToken       string `json:"accessToken"`
		AccessTokenSnake  string `json:"access_token"`
		RefreshToken      string `json:"refreshToken"`
		RefreshTokenSnake string `json:"refresh_token"`
		IDToken           string `json:"idToken"`
		IDTokenSnake      string `json:"id_token"`
		ExpiresIn         int64  `json:"expiresIn"`
		ExpiresInSnake    int64  `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return awsQuickTokens{}, fmt.Errorf("providers: decode aws-quick token response: %w", err)
	}
	tok := awsQuickTokens{
		AccessToken:  strings.TrimSpace(raw.AccessToken),
		RefreshToken: strings.TrimSpace(raw.RefreshToken),
		IDToken:      strings.TrimSpace(raw.IDToken),
		ExpiresIn:    raw.ExpiresIn,
	}
	if tok.AccessToken == "" {
		tok.AccessToken = strings.TrimSpace(raw.AccessTokenSnake)
	}
	if tok.RefreshToken == "" {
		tok.RefreshToken = strings.TrimSpace(raw.RefreshTokenSnake)
	}
	if tok.IDToken == "" {
		tok.IDToken = strings.TrimSpace(raw.IDTokenSnake)
	}
	if tok.ExpiresIn == 0 {
		tok.ExpiresIn = raw.ExpiresInSnake
	}
	if tok.IDToken == "" && tok.AccessToken == "" {
		return awsQuickTokens{}, errors.New("providers: aws-quick token response missing idToken and accessToken")
	}
	return tok, nil
}

func (m *Module) fetchAwsQuickUserInfo(ctx context.Context, accessToken, tenantURL, idcRegion string) awsQuickUserInfo {
	if strings.TrimSpace(accessToken) == "" {
		return awsQuickUserInfo{}
	}
	_, _, _, userInfoOverride := m.awsQuickEndpoints(idcRegion)
	targetURL := userInfoOverride
	if targetURL == "" {
		base := strings.TrimRight(strings.TrimSpace(tenantURL), "/")
		if !strings.Contains(base, "dp.appintegrations") {
			base = awsQuickSocialDataPlaneURL
		}
		targetURL = base + "/integration/quick-work/user-info"
	}

	body, err := m.postAwsQuickJSON(ctx, targetURL, map[string]any{}, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	if err != nil {
		return awsQuickUserInfo{}
	}

	var envelope map[string]any
	if err := json.Unmarshal(body, &envelope); err != nil {
		return awsQuickUserInfo{}
	}
	if nested, ok := envelope["userInfo"].(map[string]any); ok {
		envelope = nested
	}
	pick := func(keys ...string) string {
		for _, k := range keys {
			if v, ok := envelope[k].(string); ok && strings.TrimSpace(v) != "" {
				return strings.TrimSpace(v)
			}
		}
		return ""
	}
	return awsQuickUserInfo{
		Email:        pick("email", "emailAddress"),
		UserName:     pick("userName", "user_name", "username"),
		DisplayName:  pick("displayName", "display_name"),
		UserARN:      pick("userArn", "user_arn"),
		AWSAccountID: pick("awsAccountId", "aws_account_id", "accountId"),
		AccountName:  pick("accountName", "account_name"),
	}
}

func (m *Module) startAwsQuickProxy() (port int, baseURL string, err error) {
	rt := m.awsQuick
	rt.mu.Lock()
	defer rt.mu.Unlock()

	if rt.server != nil && rt.port > 0 {
		if rt.idleTimer != nil {
			rt.idleTimer.Reset(awsQuickSessionTTL)
		}
		return rt.port, fmt.Sprintf("http://127.0.0.1:%d", rt.port), nil
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, "", fmt.Errorf("providers: start aws-quick loopback listener: %w", err)
	}
	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok || tcpAddr.Port <= 0 {
		_ = ln.Close()
		return 0, "", errors.New("providers: resolve aws-quick loopback port")
	}

	srv := &http.Server{
		Handler:           http.HandlerFunc(m.handleAwsQuickLoopback),
		ReadHeaderTimeout: 10 * time.Second,
	}
	rt.listener = ln
	rt.server = srv
	rt.port = tcpAddr.Port
	rt.idleTimer = time.AfterFunc(awsQuickSessionTTL, func() {
		m.stopAwsQuickProxy("")
	})

	go func() {
		_ = srv.Serve(ln)
	}()

	return rt.port, fmt.Sprintf("http://127.0.0.1:%d", rt.port), nil
}

func (m *Module) stopAwsQuickProxy(keepState string) {
	rt := m.awsQuick
	if rt == nil {
		return
	}
	rt.mu.Lock()
	defer rt.mu.Unlock()

	for st, sess := range rt.sessions {
		if st != keepState && (sess.Status == "pending" || sess.Status == "discovery_saved") {
			if keepState != "" {
				delete(rt.sessions, keepState)
			}
			if rt.idleTimer != nil {
				rt.idleTimer.Reset(awsQuickSessionTTL)
			}
			return
		}
	}

	if rt.idleTimer != nil {
		rt.idleTimer.Stop()
		rt.idleTimer = nil
	}
	if rt.server != nil {
		srv := rt.server
		rt.server = nil
		rt.listener = nil
		rt.port = 0
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = srv.Shutdown(ctx)
		}()
	}

	now := time.Now().UTC()
	for st, sess := range rt.sessions {
		if st == keepState || sess.Status == "done" {
			sess.ExpiresAt = now.Add(awsQuickSessionTTL)
		} else {
			delete(rt.sessions, st)
		}
	}
}

func (m *Module) shutdownAwsQuickProxy() {
	rt := m.awsQuick
	if rt == nil {
		return
	}
	rt.mu.Lock()
	if rt.idleTimer != nil {
		rt.idleTimer.Stop()
		rt.idleTimer = nil
	}
	srv := rt.server
	rt.server = nil
	rt.listener = nil
	rt.port = 0
	rt.mu.Unlock()

	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}
}

func renderAwsQuickResultHTML(ok bool, message string) string {
	title := "Amazon Quick Sign-In"
	status := "Authentication Failed"
	color := "#ef4444"
	if ok {
		status = "Connected to Amazon Quick"
		color = "#10b981"
	}
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head><meta charset="utf-8"><title>%s</title>
<style>body{font-family:system-ui,-apple-system,sans-serif;background:#0f172a;color:#e2e8f0;display:flex;align-items:center;justify-content:center;min-height:100vh;margin:0}
.card{max-width:460px;padding:28px;border-radius:16px;background:#1e293b;border:1px solid #334155;text-align:center}
h1{font-size:18px;margin:0 0 10px;color:%s}p{font-size:14px;color:#cbd5e1;margin:0;line-height:1.5}</style></head>
<body><div class="card"><h1>%s</h1><p>%s</p></div></body></html>`,
		html.EscapeString(title), color, html.EscapeString(status), html.EscapeString(message))
}

func (m *Module) completeAwsQuickHop1(ctx context.Context, sess *awsQuickSession, discovery awsQuickDiscovery) (string, error) {
	if discovery.IssuerURL == "" {
		return "", errors.New("discovery callback did not provide issuerUrl")
	}
	clientID, clientSecret, idcRegion, err := m.registerAwsQuickOIDCClient(ctx, discovery, sess.RedirectURI)
	if err != nil {
		return "", err
	}
	discovery.IDCRegion = idcRegion
	authorizeURL := m.buildAwsQuickAuthorizeURL(clientID, sess.RedirectURI, sess.State, sess.Challenge, idcRegion)

	m.awsQuick.mu.Lock()
	sess.Discovery = &discovery
	sess.DiscoverySavedAt = time.Now().UTC()
	sess.ClientID = clientID
	sess.ClientSecret = clientSecret
	sess.AuthorizeURL = authorizeURL
	sess.Status = "discovery_saved"
	sess.ErrMsg = ""
	m.awsQuick.mu.Unlock()

	return authorizeURL, nil
}

func (m *Module) completeAwsQuickHop2(ctx context.Context, sess *awsQuickSession, q url.Values) (*awsQuickResult, error) {
	m.awsQuick.mu.Lock()
	if sess.Status == "done" && sess.Result != nil {
		res := sess.Result
		m.awsQuick.mu.Unlock()
		return res, nil
	}
	expectedState := sess.State
	clientID := sess.ClientID
	clientSecret := sess.ClientSecret
	verifier := sess.Verifier
	redirectURI := sess.RedirectURI
	var discovery awsQuickDiscovery
	if sess.Discovery != nil {
		discovery = *sess.Discovery
	}
	m.awsQuick.mu.Unlock()

	if q.Get("state") != expectedState {
		return nil, errors.New("oauth state mismatch")
	}
	if errCode := strings.TrimSpace(q.Get("error")); errCode != "" {
		desc := strings.TrimSpace(q.Get("error_description"))
		return nil, fmt.Errorf("sign-in failed: %s %s", errCode, desc)
	}
	code := strings.TrimSpace(q.Get("code"))
	if code == "" {
		return nil, errors.New("oauth callback missing authorization code")
	}
	idcRegion, err := validateAwsRegion(discovery.IDCRegion)
	if err != nil {
		return nil, err
	}
	_, _, tokenURL, _ := m.awsQuickEndpoints(idcRegion)
	tokens, err := m.postAwsQuickTokenJSON(ctx, tokenURL, map[string]any{
		"clientId":     clientID,
		"clientSecret": clientSecret,
		"grantType":    "authorization_code",
		"code":         code,
		"redirectUri":  strings.TrimRight(redirectURI, "/"),
		"codeVerifier": verifier,
	})
	if err != nil {
		return nil, err
	}
	userInfo := m.fetchAwsQuickUserInfo(ctx, tokens.AccessToken, discovery.TenantURL, idcRegion)

	res := &awsQuickResult{
		Tokens:       tokens,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Discovery:    discovery,
		UserInfo:     userInfo,
		TokenURL:     tokenURL,
	}

	m.awsQuick.mu.Lock()
	sess.Status = "done"
	sess.Result = res
	sess.ErrMsg = ""
	m.awsQuick.mu.Unlock()

	return res, nil
}

func (m *Module) handleAwsQuickLoopback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	q := r.URL.Query()

	// Hop 1: discovery callback (/callback?nonce=...&issuerUrl=...&tenantUrl=...&idcRegion=...)
	if r.URL.Path == "/callback" {
		m.awsQuick.mu.Lock()
		m.awsQuick.sweepLocked(time.Now().UTC())
		sess := m.awsQuick.findByNonceLocked(strings.TrimSpace(q.Get("nonce")))
		m.awsQuick.mu.Unlock()

		if sess == nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, renderAwsQuickResultHTML(false, "No active Amazon Quick login matches this link. Restart the login flow."))
			return
		}
		if errCode := strings.TrimSpace(q.Get("error")); errCode != "" {
			msg := strings.TrimSpace("Login failed: " + errCode + " " + strings.TrimSpace(q.Get("error_description")))
			m.awsQuick.mu.Lock()
			sess.Status = "error"
			sess.ErrMsg = msg
			m.awsQuick.mu.Unlock()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, renderAwsQuickResultHTML(false, msg))
			m.stopAwsQuickProxy(sess.State)
			return
		}

		discovery := parseAwsQuickDiscovery(q)
		authorizeURL, err := m.completeAwsQuickHop1(r.Context(), sess, discovery)
		if err != nil {
			m.awsQuick.mu.Lock()
			sess.Status = "error"
			sess.ErrMsg = err.Error()
			m.awsQuick.mu.Unlock()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, renderAwsQuickResultHTML(false, err.Error()))
			m.stopAwsQuickProxy(sess.State)
			return
		}

		w.Header().Set("Cache-Control", "no-store")
		http.Redirect(w, r, authorizeURL, http.StatusFound)
		return
	}

	// Hop 2: OIDC callback (/?code=...&state=...)
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		if !q.Has("code") && !q.Has("state") && !q.Has("error") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, renderAwsQuickResultHTML(false, "Waiting for Amazon Quick sign-in — complete the login in the quick.aws.com tab."))
			return
		}

		m.awsQuick.mu.Lock()
		m.awsQuick.sweepLocked(time.Now().UTC())
		sess := m.awsQuick.sessions[strings.TrimSpace(q.Get("state"))]
		m.awsQuick.mu.Unlock()

		if sess == nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, renderAwsQuickResultHTML(false, "No active Amazon Quick login session."))
			return
		}

		res, err := m.completeAwsQuickHop2(r.Context(), sess, q)
		if err != nil {
			m.awsQuick.mu.Lock()
			if sess.Status != "done" {
				sess.Status = "error"
				sess.ErrMsg = err.Error()
			}
			m.awsQuick.mu.Unlock()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, renderAwsQuickResultHTML(false, err.Error()))
			m.stopAwsQuickProxy(sess.State)
			return
		}

		who := res.UserInfo.Email
		if who == "" {
			who = res.UserInfo.DisplayName
		}
		if who == "" {
			who = res.UserInfo.UserName
		}
		msg := "You can close this window."
		if who != "" {
			msg = fmt.Sprintf("Signed in as %s. You can close this window.", who)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, renderAwsQuickResultHTML(true, msg))
		m.stopAwsQuickProxy("")
		return
	}

	http.NotFound(w, r)
}

func (m *Module) handleAwsQuickOAuthStart(w http.ResponseWriter, r *http.Request, adapter ProviderAdapter, req oauthStartRequest) {
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
	if baseURL != "" {
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

	port, redirectURI, err := m.startAwsQuickProxy()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to start Amazon Quick loopback callback server")
		return
	}

	verifier, challenge, err := generatePKCE()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate PKCE verifier")
		return
	}
	state, err := generateOAuthState()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate OAuth state")
		return
	}
	nonce, err := generateAwsQuickNonce()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate OAuth nonce")
		return
	}

	now := time.Now().UTC()
	expiresAt := now.Add(awsQuickSessionTTL)
	signinURL := buildAwsQuickSigninURL(port, nonce)

	m.awsQuick.mu.Lock()
	m.awsQuick.sweepLocked(now)
	m.awsQuick.sessions[state] = &awsQuickSession{
		State:          state,
		Nonce:          nonce,
		Verifier:       verifier,
		Challenge:      challenge,
		Port:           port,
		RedirectURI:    redirectURI,
		ConnectionName: connName,
		BaseURL:        baseURL,
		Models:         models,
		AuthorizeURL:   signinURL,
		Status:         "pending",
		CreatedAt:      now,
		ExpiresAt:      expiresAt,
	}
	m.awsQuick.mu.Unlock()

	httpserver.WriteJSON(w, http.StatusOK, oauthStartResponse{
		State:        state,
		AuthorizeURL: signinURL,
		RedirectURI:  redirectURI,
		ExpiresAt:    expiresAt.Format(rfc3339NanoFixed),
	})
}

type oauthStatusResponse struct {
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
	AuthorizeURL string `json:"authorize_url,omitempty"`
	Interval     int    `json:"interval,omitempty"`
}

func (m *Module) handleOAuthStatus(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "provider"))
	strategy := m.getOAuthStrategy(providerID)
	if strategy.FixedPort > 0 {
		m.handleLoopbackOAuthStatus(w, r, providerID)
		return
	}
	if strategy.FlowType == oauthFlowDeviceCode {
		m.handleDeviceOAuthStatus(w, r, providerID)
		return
	}
	if strategy.FlowType != oauthFlowLoopbackProxy {
		httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_not_supported", "oauth status polling is only supported for loopback and device-code providers")
		return
	}
	state := strings.TrimSpace(r.URL.Query().Get("state"))
	if state == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "state query parameter is required")
		return
	}

	m.awsQuick.mu.Lock()
	m.awsQuick.sweepLocked(time.Now().UTC())
	sess, ok := m.awsQuick.sessions[state]
	if !ok {
		m.awsQuick.mu.Unlock()
		httpserver.WriteJSON(w, http.StatusOK, oauthStatusResponse{Status: "unknown"})
		return
	}
	resp := oauthStatusResponse{
		Status:       sess.Status,
		Error:        sess.ErrMsg,
		AuthorizeURL: sess.AuthorizeURL,
	}
	if sess.Status == "error" {
		delete(m.awsQuick.sessions, state)
		m.awsQuick.mu.Unlock()
		m.stopAwsQuickProxy("")
		httpserver.WriteJSON(w, http.StatusOK, resp)
		return
	}
	m.awsQuick.mu.Unlock()

	httpserver.WriteJSON(w, http.StatusOK, resp)
}

type oauthPasteRequest struct {
	URL string `json:"url"`
}

type oauthPasteResponse struct {
	Status       string `json:"status"`
	State        string `json:"state"`
	AuthorizeURL string `json:"authorize_url,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
}

func (m *Module) handleOAuthPaste(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "provider"))
	if providerID != providerIDAWSQuick {
		httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_not_supported", "oauth paste-back is only supported for aws-quick")
		return
	}

	var req oauthPasteRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	q, ok := parseAwsQuickPastedQuery(req.URL)
	if !ok {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request",
			"Not a completed-login address. Paste the full URL from the browser address bar after Amazon Quick sign-in finishes.")
		return
	}

	m.awsQuick.mu.Lock()
	m.awsQuick.sweepLocked(time.Now().UTC())
	m.awsQuick.mu.Unlock()

	// Hop 1: discovery URL (/callback?nonce=...&issuerUrl=...)
	if q.Has("nonce") {
		m.awsQuick.mu.Lock()
		sess := m.awsQuick.findByNonceLocked(strings.TrimSpace(q.Get("nonce")))
		m.awsQuick.mu.Unlock()
		if sess == nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_oauth_state",
				"Login session not found — restart the Amazon Quick sign-in flow.")
			return
		}
		if errCode := strings.TrimSpace(q.Get("error")); errCode != "" {
			httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_error",
				strings.TrimSpace("Login failed: "+errCode+" "+strings.TrimSpace(q.Get("error_description"))))
			return
		}
		discovery := parseAwsQuickDiscovery(q)
		authorizeURL, err := m.completeAwsQuickHop1(r.Context(), sess, discovery)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_error", err.Error())
			return
		}
		httpserver.WriteJSON(w, http.StatusOK, oauthPasteResponse{
			Status:       "discovery_saved",
			State:        sess.State,
			AuthorizeURL: authorizeURL,
			RedirectURI:  sess.RedirectURI,
		})
		return
	}

	// Hop 2: authorization code URL (/?code=...&state=...)
	if q.Has("code") {
		m.awsQuick.mu.Lock()
		sess := m.awsQuick.sessions[strings.TrimSpace(q.Get("state"))]
		m.awsQuick.mu.Unlock()
		if sess == nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_oauth_state",
				"Login session not found — restart the Amazon Quick sign-in flow.")
			return
		}
		if _, err := m.completeAwsQuickHop2(r.Context(), sess, q); err != nil {
			httpserver.WriteError(w, r, http.StatusBadGateway, "oauth_exchange_failed", err.Error())
			return
		}
		httpserver.WriteJSON(w, http.StatusOK, oauthPasteResponse{
			Status: "done",
			State:  sess.State,
		})
		return
	}

	if errCode := strings.TrimSpace(q.Get("error")); errCode != "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_error",
			strings.TrimSpace("Login failed: "+errCode+" "+strings.TrimSpace(q.Get("error_description"))))
		return
	}

	httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request",
		"Not a completed-login address. Paste the full URL containing nonce=... (Step 1) or code=...&state=... (Step 2).")
}

type oauthStopResponse struct {
	Stopped bool `json:"stopped"`
}

func (m *Module) handleOAuthStop(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "provider"))
	strategy := m.getOAuthStrategy(providerID)
	if strategy.FixedPort > 0 {
		m.oauthProxy.stop(providerID)
		httpserver.WriteJSON(w, http.StatusOK, oauthStopResponse{Stopped: true})
		return
	}
	if strategy.FlowType != oauthFlowLoopbackProxy {
		httpserver.WriteError(w, r, http.StatusBadRequest, "oauth_not_supported", "oauth stop is only supported for loopback OAuth providers")
		return
	}
	m.stopAwsQuickProxy("")
	httpserver.WriteJSON(w, http.StatusOK, oauthStopResponse{Stopped: true})
}

func (m *Module) handleAwsQuickOAuthExchange(w http.ResponseWriter, r *http.Request, adapter ProviderAdapter, req oauthExchangeRequest) {
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

	m.awsQuick.mu.Lock()
	m.awsQuick.sweepLocked(time.Now().UTC())
	sess, ok := m.awsQuick.sessions[state]
	if !ok || sess.Status != "done" || sess.Result == nil {
		errMsg := "oauth session not completed; finish sign-in and try again"
		if ok && sess.ErrMsg != "" {
			errMsg = sess.ErrMsg
		}
		m.awsQuick.mu.Unlock()
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_oauth_state", errMsg)
		return
	}
	res := sess.Result
	sessConnName := sess.ConnectionName
	sessBaseURL := sess.BaseURL
	sessModels := cloneStringSlice(sess.Models)
	delete(m.awsQuick.sessions, state)
	m.awsQuick.mu.Unlock()

	m.stopAwsQuickProxy("")

	accessToken := res.Tokens.IDToken
	if accessToken == "" {
		accessToken = res.Tokens.AccessToken
	}
	tenantURL := strings.TrimRight(strings.TrimSpace(res.Discovery.TenantURL), "/")
	wsHost := ""
	if parsedTenant, parseErr := url.Parse(tenantURL); parseErr == nil {
		wsHost = parsedTenant.Host
	}
	region := strings.TrimSpace(res.Discovery.IDCRegion)
	if region == "" {
		region = "us-east-1"
	}
	meta := map[string]any{
		"tenant_url":     tenantURL,
		"ws_host":        wsHost,
		"actor_arn":      res.UserInfo.UserARN,
		"aws_account_id": res.UserInfo.AWSAccountID,
		"account_name":   res.UserInfo.AccountName,
		"email":          res.UserInfo.Email,
		"region":         region,
		"issuer_url":     res.Discovery.IssuerURL,
		"app_version":    awsQuickAppVersion,
	}
	stored := storedCredentials{
		AccessToken:  accessToken,
		RefreshToken: res.Tokens.RefreshToken,
		IDToken:      res.Tokens.IDToken,
		TokenType:    "Bearer",
		Scope:        awsQuickScopeProd,
		ClientID:     res.ClientID,
		ClientSecret: res.ClientSecret,
		TokenURL:     res.TokenURL,
		Metadata:     meta,
	}
	sealed, err := m.sealCredentials(stored)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to encrypt OAuth tokens")
		return
	}

	var tokenExpiresAt *string
	if res.Tokens.ExpiresIn > 0 {
		exp := time.Now().UTC().Add(time.Duration(res.Tokens.ExpiresIn) * time.Second).Format(rfc3339NanoFixed)
		tokenExpiresAt = &exp
	}

	connName := validatedConnName
	if connName == "" {
		connName = sessConnName
	}
	if connName == "" {
		switch {
		case res.UserInfo.Email != "":
			connName = "Amazon Quick (" + res.UserInfo.Email + ")"
		case res.UserInfo.DisplayName != "":
			connName = "Amazon Quick (" + res.UserInfo.DisplayName + ")"
		default:
			connName = adapter.DisplayName()
		}
	}

	baseURL := sessBaseURL
	if baseURL == "" {
		baseURL = strings.TrimSpace(res.Discovery.TenantURL)
	}
	if baseURL == "" {
		baseURL = adapter.DefaultBaseURL()
	}

	connID, err := generateID(defaultConnectionPrefix)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate connection identifier")
		return
	}

	summary, err := m.repo.CreateConnection(r.Context(), createConnectionParams{
		ID:             connID,
		ProviderID:     providerIDAWSQuick,
		Name:           connName,
		Protocol:       adapter.Protocol(),
		BaseURL:        baseURL,
		AuthType:       authTypeOAuth,
		CredentialsEnc: sealed,
		Priority:       priority,
		Enabled:        enabled,
		Models:         sessModels,
		Headers:        headers,
		TokenExpiresAt: tokenExpiresAt,
	})
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to persist Amazon Quick connection")
		return
	}

	m.publishChanged("created", "connection", summary.ID, summary.ProviderID, actorFromContext(r.Context()), summary.UpdatedAt)
	httpserver.WriteJSON(w, http.StatusCreated, summary)
}

func resolveAwsQuickTenantBaseURL(baseURL string, meta map[string]any) string {
	if tenant := strings.TrimRight(stringClaim(meta, "tenant_url"), "/"); tenant != "" {
		return tenant
	}
	if wsHost := strings.TrimRight(stringClaim(meta, "ws_host"), "/"); wsHost != "" {
		switch {
		case strings.HasPrefix(wsHost, "wss://"):
			return "https://" + strings.TrimPrefix(wsHost, "wss://")
		case strings.HasPrefix(wsHost, "ws://"):
			return "http://" + strings.TrimPrefix(wsHost, "ws://")
		case strings.HasPrefix(wsHost, "http://"), strings.HasPrefix(wsHost, "https://"):
			return wsHost
		default:
			return "https://" + wsHost
		}
	}
	return strings.TrimRight(strings.TrimSpace(baseURL), "/")
}

func generateTraceID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("trace-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
}

// testAwsQuickConnection probes Amazon Quick's entitlements qbs-event endpoint
// (adopted from xrouter commits 105b0ff & 9d62aae): validates id_token + tenant
// + account without spending agent-hours, auto-refreshes on 401/403 (persisting
// the rotated refresh token), and treats 408/429/5xx/timeouts as soft success.
func (m *Module) testAwsQuickConnection(ctx context.Context, connectionID string, creds contracts.UpstreamCredentials) (result testConnectionResponse, healthStatus string) {
	var meta map[string]any
	if rec, err := m.repo.GetConnectionRecord(ctx, connectionID); err == nil {
		if stored, unsealErr := m.unsealCredentials(rec.CredentialsEnc); unsealErr == nil {
			meta = stored.Metadata
		}
	}

	tenantBase := resolveAwsQuickTenantBaseURL(creds.BaseURL, meta)
	if tenantBase == "" {
		return testConnectionResponse{OK: false, Error: "Missing Quick tenant/account metadata"}, "error"
	}
	accountID := stringClaim(meta, "aws_account_id")
	probeURL := tenantBase + "/integration/quick-work/qbs-event"

	doProbe := func(token string) (int, error) {
		payload := map[string]any{
			"eventType": "getAvailableEntitlementUnits",
			"eventBody": map[string]any{
				"getAvailableEntitlementUnits": map[string]any{
					"awsAccountId": accountID,
				},
			},
		}
		rawBody, err := json.Marshal(payload)
		if err != nil {
			return 0, err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, probeURL, bytes.NewReader(rawBody))
		if err != nil {
			return 0, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("x-amzn-qs-internal-trace-id", generateTraceID())

		resp, err := m.httpClient.Do(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxTestErrorBodyBytes))
		return resp.StatusCode, nil
	}

	start := time.Now()
	statusCode, probeErr := doProbe(creds.AccessToken)
	if (probeErr != nil || statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden) && connectionID != "" {
		if refreshed, refreshErr := m.RefreshOAuthToken(ctx, connectionID); refreshErr == nil && refreshed.AccessToken != "" {
			statusCode, probeErr = doProbe(refreshed.AccessToken)
		}
	}
	latencyMS := time.Since(start).Milliseconds()

	if probeErr != nil {
		// Upstream timeout/network blip says nothing about credentials: soft success.
		return testConnectionResponse{
			OK:        true,
			LatencyMS: latencyMS,
			Error:     "Quick API probe failed or timed out — quota unknown",
		}, "ok"
	}
	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		return testConnectionResponse{OK: true, StatusCode: statusCode, LatencyMS: latencyMS}, "ok"
	}
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		return testConnectionResponse{
			OK:         false,
			StatusCode: statusCode,
			LatencyMS:  latencyMS,
			Error:      "Token invalid or revoked",
		}, "error"
	}
	if statusCode == http.StatusRequestTimeout || statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError {
		return testConnectionResponse{
			OK:         true,
			StatusCode: statusCode,
			LatencyMS:  latencyMS,
			Error:      fmt.Sprintf("Quick API returned %d — quota unknown", statusCode),
		}, "ok"
	}
	return testConnectionResponse{
		OK:         false,
		StatusCode: statusCode,
		LatencyMS:  latencyMS,
		Error:      fmt.Sprintf("Quick API returned %d", statusCode),
	}, "error"
}

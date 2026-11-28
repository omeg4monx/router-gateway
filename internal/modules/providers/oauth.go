package providers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const (
	oauthSessionTTL             = 10 * time.Minute
	oauthVerifierBytes          = 32 // 43 base64url chars (RFC 7636 §4.1)
	oauthStateBytes             = 24 // 32 base64url chars
	maxOAuthResponseBytes       = 1 << 20
	oauthRefreshTimeout         = 30 * time.Second
	oauthNoRefreshTokenCooldown = 10 * time.Minute
	oauthRefreshFailureCooldown = 30 * time.Second
	maxRefreshCooldownEntries   = 1024
)

type refreshCooldownEntry struct {
	until  time.Time
	reason string
}

type refreshCall struct {
	done  chan struct{}
	creds contracts.UpstreamCredentials
	err   error
}

type oauthTokenResponse struct {
	AccessToken      string         `json:"access_token"`
	RefreshToken     string         `json:"refresh_token"`
	IDToken          string         `json:"id_token"`
	TokenType        string         `json:"token_type"`
	ExpiresIn        int64          `json:"expires_in"`
	Scope            string         `json:"scope"`
	Error            string         `json:"error"`
	ErrorDescription string         `json:"error_description"`
	Metadata         map[string]any `json:"-"`
}

func generatePKCE() (verifier, challenge string, err error) {
	buf := make([]byte, oauthVerifierBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("providers: generate pkce verifier: %w", err)
	}
	verifier = base64.RawURLEncoding.EncodeToString(buf)
	sum := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(sum[:])
	return verifier, challenge, nil
}

func generateOAuthState() (string, error) {
	buf := make([]byte, oauthStateBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("providers: generate oauth state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func buildAuthorizeURL(cfg OAuthConfig, redirectURI, state, challenge string) (string, error) {
	base, err := url.Parse(cfg.AuthorizeURL)
	if err != nil {
		return "", fmt.Errorf("providers: parse authorize_url: %w", err)
	}
	q := base.Query()
	q.Set("response_type", "code")
	q.Set("client_id", cfg.ClientID)
	q.Set("redirect_uri", redirectURI)
	if len(cfg.Scopes) > 0 {
		q.Set("scope", strings.Join(cfg.Scopes, " "))
	}
	q.Set("state", state)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", "S256")
	for k, v := range cfg.ExtraAuthParams {
		if strings.TrimSpace(k) != "" && v != "" {
			q.Set(k, v)
		}
	}
	base.RawQuery = q.Encode()
	return base.String(), nil
}

func parseOAuthCallbackInput(callbackURL, directCode, directState string) (code, state, oauthErr, oauthErrDesc string, err error) {
	code = strings.TrimSpace(directCode)
	state = strings.TrimSpace(directState)

	if trimmedURL := strings.TrimSpace(callbackURL); trimmedURL != "" {
		parsed, parseErr := url.Parse(trimmedURL)
		if parseErr != nil {
			return "", "", "", "", errors.New("callback_url must be a valid URL")
		}
		q := parsed.Query()
		if qErr := strings.TrimSpace(q.Get("error")); qErr != "" {
			return "", "", qErr, strings.TrimSpace(q.Get("error_description")), nil
		}
		urlCode := strings.TrimSpace(q.Get("code"))
		urlState := strings.TrimSpace(q.Get("state"))
		if urlCode != "" {
			if code != "" && code != urlCode {
				return "", "", "", "", errors.New("code conflicts with callback_url")
			}
			code = urlCode
		}
		if urlState != "" {
			if state != "" && state != urlState {
				return "", "", "", "", errors.New("state conflicts with callback_url")
			}
			state = urlState
		}
	}

	if code == "" {
		return "", "", "", "", errors.New("authorization code is required (provide code+state or callback_url)")
	}
	if state == "" {
		return "", "", "", "", errors.New("oauth state is required (provide code+state or callback_url)")
	}
	return code, state, "", "", nil
}

func (m *Module) resolveOAuthConfig(providerID, reqClientID, reqClientSecret string) (OAuthConfig, error) {
	adapter, ok := m.GetAdapter(providerID)
	if !ok {
		return OAuthConfig{}, fmt.Errorf("providers: unknown provider %q", providerID)
	}
	base := adapter.OAuthConfig()
	if base == nil {
		return OAuthConfig{}, contracts.ErrOAuthRefreshUnsupported
	}
	cfg := *cloneOAuthConfig(base)

	// Merge environment variables (never hardcoded secrets).
	prefix := oauthEnvPrefix(providerID)
	if envID := strings.TrimSpace(os.Getenv(prefix + "_CLIENT_ID")); envID != "" {
		cfg.ClientID = envID
	}
	if envSecret := strings.TrimSpace(os.Getenv(prefix + "_SECRET")); envSecret != "" {
		cfg.ClientSecret = envSecret
	} else if envClientSecret := strings.TrimSpace(os.Getenv(prefix + "_CLIENT_SECRET")); envClientSecret != "" {
		cfg.ClientSecret = envClientSecret
	}

	// Merge Options.OAuthOverrides (used for custom deployments and tests).
	m.mu.RLock()
	override, hasOverride := m.oauthOverrides[providerID]
	m.mu.RUnlock()
	if hasOverride {
		if strings.TrimSpace(override.AuthorizeURL) != "" {
			cfg.AuthorizeURL = strings.TrimSpace(override.AuthorizeURL)
		}
		if strings.TrimSpace(override.TokenURL) != "" {
			cfg.TokenURL = strings.TrimSpace(override.TokenURL)
		}
		if len(override.Scopes) > 0 {
			cfg.Scopes = cloneStringSlice(override.Scopes)
		}
		if strings.TrimSpace(override.ClientID) != "" {
			cfg.ClientID = strings.TrimSpace(override.ClientID)
		}
		if strings.TrimSpace(override.ClientSecret) != "" {
			cfg.ClientSecret = strings.TrimSpace(override.ClientSecret)
		}
		if strings.TrimSpace(override.DefaultRedirectURI) != "" {
			cfg.DefaultRedirectURI = strings.TrimSpace(override.DefaultRedirectURI)
		}
		if len(override.ExtraAuthParams) > 0 {
			maps.Copy(cfg.ExtraAuthParams, override.ExtraAuthParams)
		}
		if len(override.ExtraTokenParams) > 0 {
			maps.Copy(cfg.ExtraTokenParams, override.ExtraTokenParams)
		}
	}

	// Apply per-request overrides last.
	if trimmedID := strings.TrimSpace(reqClientID); trimmedID != "" {
		cfg.ClientID = trimmedID
	}
	if trimmedSecret := strings.TrimSpace(reqClientSecret); trimmedSecret != "" {
		cfg.ClientSecret = trimmedSecret
	}

	if cfg.DefaultRedirectURI == "" {
		cfg.DefaultRedirectURI = defaultOAuthRedirectURI
	}
	return cfg, nil
}

func isAnthropicTokenURL(rawURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	return host == "api.anthropic.com" || host == "console.anthropic.com"
}

func (m *Module) postJSONTokenEndpoint(ctx context.Context, tokenURL string, payload map[string]any) (oauthTokenResponse, error) {
	status, body, err := m.postRawJSON(ctx, tokenURL, payload)
	if err != nil {
		return oauthTokenResponse{}, fmt.Errorf("providers: token endpoint request failed: %w", err)
	}
	var tok oauthTokenResponse
	if len(body) > 0 {
		_ = json.Unmarshal(body, &tok)
	}
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		if tok.Error != "" {
			return oauthTokenResponse{}, fmt.Errorf("providers: token endpoint returned %d (%s)", status, tok.Error)
		}
		return oauthTokenResponse{}, fmt.Errorf("providers: token endpoint returned status %d", status)
	}
	if strings.TrimSpace(tok.AccessToken) == "" {
		return oauthTokenResponse{}, errors.New("providers: token endpoint response missing access_token")
	}
	return tok, nil
}

func (m *Module) postTokenEndpoint(ctx context.Context, tokenURL string, form url.Values) (oauthTokenResponse, error) {
	validURL, err := validateHTTPURL(tokenURL)
	if err != nil {
		return oauthTokenResponse{}, fmt.Errorf("providers: invalid token endpoint URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, validURL, strings.NewReader(form.Encode()))
	if err != nil {
		return oauthTokenResponse{}, fmt.Errorf("providers: create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return oauthTokenResponse{}, fmt.Errorf("providers: token endpoint request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxOAuthResponseBytes))
	if err != nil {
		return oauthTokenResponse{}, fmt.Errorf("providers: read token endpoint response: %w", err)
	}

	var tok oauthTokenResponse
	if len(body) > 0 {
		_ = json.Unmarshal(body, &tok)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if tok.Error != "" {
			return oauthTokenResponse{}, fmt.Errorf("providers: token endpoint returned %d (%s)", resp.StatusCode, tok.Error)
		}
		return oauthTokenResponse{}, fmt.Errorf("providers: token endpoint returned status %d", resp.StatusCode)
	}

	if strings.TrimSpace(tok.AccessToken) == "" {
		return oauthTokenResponse{}, errors.New("providers: token endpoint response missing access_token")
	}
	return tok, nil
}

// RefreshOAuthToken refreshes the OAuth access token for connectionID,
// coalescing concurrent refresh calls for the same connectionID via an
// in-flight single-flight map so only one upstream token-endpoint POST fires,
// and applying per-connection cooldowns after failed or unsupported refreshes.
func (m *Module) RefreshOAuthToken(ctx context.Context, connectionID string) (contracts.UpstreamCredentials, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return contracts.UpstreamCredentials{}, err
	}

	now := time.Now().UTC()
	m.refreshMu.Lock()
	if existing, ok := m.refreshCalls[connectionID]; ok {
		m.refreshMu.Unlock()
		select {
		case <-existing.done:
			return cloneUpstreamCredentials(existing.creds), existing.err
		case <-ctx.Done():
			return contracts.UpstreamCredentials{}, ctx.Err()
		}
	}
	if entry, cooling := m.getRefreshCooldownLocked(connectionID, now); cooling {
		m.refreshMu.Unlock()
		if entry.reason == "no_refresh_token" {
			return contracts.UpstreamCredentials{}, contracts.ErrOAuthRefreshUnsupported
		}
		return contracts.UpstreamCredentials{}, fmt.Errorf("providers: oauth refresh cooling down for %s", connectionID)
	}

	call := &refreshCall{done: make(chan struct{})}
	m.refreshCalls[connectionID] = call
	m.refreshMu.Unlock()

	go func() {
		refreshCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), oauthRefreshTimeout)
		defer cancel()

		creds, err := m.doRefreshOAuthToken(refreshCtx, connectionID)
		call.creds = creds
		call.err = err

		doneNow := time.Now().UTC()
		m.refreshMu.Lock()
		delete(m.refreshCalls, connectionID)
		switch {
		case err == nil:
			delete(m.refreshCooldowns, connectionID)
		case errors.Is(err, contracts.ErrOAuthRefreshUnsupported):
			m.markRefreshCooldownLocked(connectionID, "no_refresh_token", doneNow)
		default:
			m.markRefreshCooldownLocked(connectionID, "refresh_failed", doneNow)
		}
		close(call.done)
		m.refreshMu.Unlock()
	}()

	select {
	case <-call.done:
		return cloneUpstreamCredentials(call.creds), call.err
	case <-ctx.Done():
		return contracts.UpstreamCredentials{}, ctx.Err()
	}
}

func (m *Module) getRefreshCooldownLocked(connectionID string, now time.Time) (refreshCooldownEntry, bool) {
	entry, ok := m.refreshCooldowns[connectionID]
	if !ok {
		return refreshCooldownEntry{}, false
	}
	if !now.Before(entry.until) {
		delete(m.refreshCooldowns, connectionID)
		return refreshCooldownEntry{}, false
	}
	return entry, true
}

func (m *Module) markRefreshCooldownLocked(connectionID, reason string, now time.Time) {
	ttl := oauthRefreshFailureCooldown
	if reason == "no_refresh_token" {
		ttl = oauthNoRefreshTokenCooldown
	}
	if len(m.refreshCooldowns) >= maxRefreshCooldownEntries {
		for k, v := range m.refreshCooldowns {
			if !now.Before(v.until) {
				delete(m.refreshCooldowns, k)
			}
		}
		if len(m.refreshCooldowns) >= maxRefreshCooldownEntries {
			var oldestKey string
			var oldestUntil time.Time
			for k, v := range m.refreshCooldowns {
				if oldestKey == "" || v.until.Before(oldestUntil) {
					oldestKey = k
					oldestUntil = v.until
				}
			}
			if oldestKey != "" {
				delete(m.refreshCooldowns, oldestKey)
			}
		}
	}
	m.refreshCooldowns[connectionID] = refreshCooldownEntry{
		until:  now.Add(ttl),
		reason: reason,
	}
}

func (m *Module) clearRefreshCooldown(connectionID string) {
	m.refreshMu.Lock()
	delete(m.refreshCooldowns, connectionID)
	m.refreshMu.Unlock()
}

func (m *Module) doRefreshOAuthToken(ctx context.Context, connectionID string) (contracts.UpstreamCredentials, error) {
	rec, err := m.repo.GetConnectionRecord(ctx, connectionID)
	if err != nil {
		return contracts.UpstreamCredentials{}, err
	}
	if rec.Summary.AuthType != authTypeOAuth {
		return contracts.UpstreamCredentials{}, contracts.ErrOAuthRefreshUnsupported
	}

	stored, err := m.unsealCredentials(rec.CredentialsEnc)
	if err != nil {
		return contracts.UpstreamCredentials{}, err
	}
	if strings.TrimSpace(stored.RefreshToken) == "" {
		return contracts.UpstreamCredentials{}, contracts.ErrOAuthRefreshUnsupported
	}

	cfg, cfgErr := m.resolveOAuthConfig(rec.Summary.ProviderID, "", "")
	tokenURL := strings.TrimSpace(stored.TokenURL)
	clientID := strings.TrimSpace(stored.ClientID)
	clientSecret := strings.TrimSpace(stored.ClientSecret)

	// Allow explicit OAuthOverrides.TokenURL to redirect refresh requests (e.g. in tests or custom proxies).
	m.mu.RLock()
	override, hasOverride := m.oauthOverrides[rec.Summary.ProviderID]
	m.mu.RUnlock()
	if hasOverride && strings.TrimSpace(override.TokenURL) != "" {
		tokenURL = strings.TrimSpace(override.TokenURL)
	}
	if cfgErr == nil {
		if tokenURL == "" {
			tokenURL = cfg.TokenURL
		}
		if clientID == "" {
			clientID = cfg.ClientID
		}
		if clientSecret == "" {
			clientSecret = cfg.ClientSecret
		}
	}
	if tokenURL == "" {
		return contracts.UpstreamCredentials{}, contracts.ErrOAuthRefreshUnsupported
	}

	strategy := m.getOAuthStrategy(rec.Summary.ProviderID)
	var tok oauthTokenResponse
	var postErr error
	if strategy.RefreshToken != nil {
		tok, postErr = strategy.RefreshToken(ctx, m, tokenURL, clientID, clientSecret, stored)
	} else {
		tok, postErr = defaultRefreshFormToken(ctx, m, tokenURL, clientID, clientSecret, stored)
	}
	if postErr != nil {
		return contracts.UpstreamCredentials{}, postErr
	}

	stored.AccessToken = tok.AccessToken
	if strings.TrimSpace(tok.RefreshToken) != "" {
		stored.RefreshToken = strings.TrimSpace(tok.RefreshToken)
	}
	if strings.TrimSpace(tok.IDToken) != "" {
		stored.IDToken = strings.TrimSpace(tok.IDToken)
	}
	if strings.TrimSpace(tok.TokenType) != "" {
		stored.TokenType = strings.TrimSpace(tok.TokenType)
	}
	if strings.TrimSpace(tok.Scope) != "" {
		stored.Scope = strings.TrimSpace(tok.Scope)
	}
	stored.ClientID = clientID
	stored.ClientSecret = clientSecret
	stored.TokenURL = tokenURL
	if len(tok.Metadata) > 0 {
		if stored.Metadata == nil {
			stored.Metadata = make(map[string]any, len(tok.Metadata))
		}
		maps.Copy(stored.Metadata, tok.Metadata)
	}

	sealed, err := m.sealCredentials(stored)
	if err != nil {
		return contracts.UpstreamCredentials{}, err
	}

	var tokenExpiresAt *string
	if tok.ExpiresIn > 0 {
		exp := time.Now().UTC().Add(time.Duration(tok.ExpiresIn) * time.Second).Format(rfc3339NanoFixed)
		tokenExpiresAt = &exp
	}

	updatedAt, err := m.repo.UpdateOAuthCredentials(ctx, connectionID, sealed, tokenExpiresAt)
	if err != nil {
		return contracts.UpstreamCredentials{}, err
	}

	m.publishChanged("token_refreshed", "connection", connectionID, rec.Summary.ProviderID, "system", updatedAt)

	extraHeaders := cloneStringMap(rec.Summary.Headers)
	applyStoredProviderHeaders(rec.Summary.ProviderID, stored, extraHeaders)

	return contracts.UpstreamCredentials{
		ConnectionID: rec.Summary.ID,
		ProviderID:   rec.Summary.ProviderID,
		Protocol:     rec.Summary.Protocol,
		BaseURL:      rec.Summary.BaseURL,
		AuthType:     rec.Summary.AuthType,
		APIKey:       stored.APIKey,
		AccessToken:  stored.AccessToken,
		ExtraHeaders: extraHeaders,
	}, nil
}

func applyStoredProviderHeaders(providerID string, stored storedCredentials, headers map[string]string) {
	if providerID == providerIDJetBrains && stored.Metadata != nil {
		if jwt, ok := stored.Metadata["grazie_jwt"].(string); ok && strings.TrimSpace(jwt) != "" {
			headers["Grazie-Authenticate-JWT"] = strings.TrimSpace(jwt)
		}
	}
}

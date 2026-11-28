package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	providerIDJetBrains         = "jetbrains"
	jetbrainsDefaultBaseURL     = "https://api.jetbrains.ai/user/v5/llm/openai/v1"
	jetbrainsDefaultModel       = "gpt-6-luna"
	jetbrainsOAuthClientID      = "junie-cli"
	jetbrainsOAuthAuthorizeURL  = "https://junie.jetbrains.com/cli-auth"
	jetbrainsOAuthTokenURL      = "https://oauth.account.jetbrains.com/oauth2/token" //nolint:gosec // G101: public JetBrains OAuth2 token endpoint URL
	jetbrainsOAuthFixedPort     = 62345
	jetbrainsOAuthRedirectURI   = "http://localhost:62345"
	jetbrainsGrazieAgentHeader  = `{"name":"Pulse","version":"1"}`
	jetbrainsGrazieProvisionTTL = 15 * time.Second
)

// jetbrainsAIBaseURL is overridable in unit tests.
var jetbrainsAIBaseURL = "https://api.jetbrains.ai"

func newJetBrainsOAuthStrategy() *oauthProviderStrategy {
	return &oauthProviderStrategy{
		ID:           providerIDJetBrains,
		FlowType:     oauthFlowAuthCodePKCE,
		FixedPort:    jetbrainsOAuthFixedPort,
		CallbackPath: "/",
		BuildAuthURL: buildJetBrainsAuthURL,
		ExchangeCode: exchangeJetBrainsCode,
		PostExchange: postExchangeJetBrains,
		RefreshToken: refreshJetBrainsToken,
	}
}

func normalizeJetBrainsRedirectURI(redirectURI string) string {
	trimmed := strings.TrimSpace(redirectURI)
	if strings.Contains(trimmed, "localhost:62345") {
		return trimmed
	}
	return jetbrainsOAuthRedirectURI
}

func buildJetBrainsAuthURL(cfg OAuthConfig, redirectURI, state, challenge string) (string, error) {
	authBase := strings.TrimSpace(cfg.AuthorizeURL)
	if authBase == "" {
		authBase = jetbrainsOAuthAuthorizeURL
	}
	parsed, err := url.Parse(authBase)
	if err != nil {
		return "", fmt.Errorf("providers: invalid jetbrains authorize url: %w", err)
	}
	clientID := strings.TrimSpace(cfg.ClientID)
	if clientID == "" {
		clientID = jetbrainsOAuthClientID
	}
	scope := strings.TrimSpace(strings.Join(cfg.Scopes, " "))
	if scope == "" {
		scope = "offline_access openid jb-authn-service"
	}

	q := parsed.Query()
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("redirect_uri", normalizeJetBrainsRedirectURI(redirectURI))
	q.Set("scope", scope)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", "S256")
	q.Set("state", state)
	parsed.RawQuery = q.Encode()
	return parsed.String(), nil
}

func cleanJetBrainsAuthCode(raw string) string {
	clean := strings.TrimSpace(raw)
	if clean == "" {
		return ""
	}
	if strings.Contains(clean, "://") || strings.Contains(clean, "code=") || strings.HasPrefix(clean, "/") {
		candidate := clean
		if !strings.HasPrefix(candidate, "http://") && !strings.HasPrefix(candidate, "https://") {
			if !strings.HasPrefix(candidate, "/") {
				candidate = "/" + candidate
			}
			candidate = "http://localhost" + candidate
		}
		if parsed, err := url.Parse(candidate); err == nil {
			if extracted := strings.TrimSpace(parsed.Query().Get("code")); extracted != "" {
				return extracted
			}
		}
	}
	return clean
}

func exchangeJetBrainsCode(
	ctx context.Context,
	m *Module,
	cfg OAuthConfig,
	sess oauthSessionRecord,
	code, _, verifier, _ string,
) (oauthTokenResponse, error) {
	tokenURL := strings.TrimSpace(cfg.TokenURL)
	if tokenURL == "" {
		tokenURL = jetbrainsOAuthTokenURL
	}
	clientID := strings.TrimSpace(cfg.ClientID)
	if clientID == "" {
		clientID = jetbrainsOAuthClientID
	}
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {clientID},
		"code":          {cleanJetBrainsAuthCode(code)},
		"redirect_uri":  {normalizeJetBrainsRedirectURI(sess.RedirectURI)},
		"code_verifier": {verifier},
	}
	reqCtx, cancel := context.WithTimeout(ctx, jetbrainsGrazieProvisionTTL)
	defer cancel()
	return m.postTokenEndpoint(reqCtx, tokenURL, form)
}

func provisionJetBrainsGrazieJWT(ctx context.Context, client *http.Client, idToken string) (licenseID, grazieJWT string) {
	trimmedID := strings.TrimSpace(idToken)
	if trimmedID == "" || client == nil {
		return "", ""
	}

	reqCtx, cancel := context.WithTimeout(ctx, jetbrainsGrazieProvisionTTL)
	defer cancel()

	base := strings.TrimRight(jetbrainsAIBaseURL, "/")
	obtainURL := base + "/auth/jetbrains-jwt/license/obtain/grazie-lite"
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, obtainURL, bytes.NewReader([]byte("{}")))
	if err != nil {
		return "", ""
	}
	req.Header.Set("Authorization", "Bearer "+trimmedID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", ""
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxOAuthResponseBytes))
	_ = resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", ""
	}

	var liteResp struct {
		Token   string `json:"token"`
		License struct {
			LicenseID string `json:"licenseId"`
		} `json:"license"`
	}
	if err := json.Unmarshal(body, &liteResp); err != nil {
		return "", ""
	}
	licenseID = strings.TrimSpace(liteResp.License.LicenseID)
	grazieJWT = strings.TrimSpace(liteResp.Token)

	if grazieJWT == "" && licenseID != "" {
		provideURL := base + "/auth/jetbrains-jwt/provide-access/license/v2"
		payload, _ := json.Marshal(map[string]string{"licenseId": licenseID})
		pReq, pErr := http.NewRequestWithContext(reqCtx, http.MethodPost, provideURL, bytes.NewReader(payload))
		if pErr != nil {
			return licenseID, ""
		}
		pReq.Header.Set("Authorization", "Bearer "+trimmedID)
		pReq.Header.Set("Content-Type", "application/json")
		pReq.Header.Set("Accept", "application/json")

		pResp, pDoErr := client.Do(pReq)
		if pDoErr != nil {
			return licenseID, ""
		}
		pBody, _ := io.ReadAll(io.LimitReader(pResp.Body, maxOAuthResponseBytes))
		_ = pResp.Body.Close()
		if pResp.StatusCode >= http.StatusOK && pResp.StatusCode < http.StatusMultipleChoices {
			var accessResp struct {
				Token string `json:"token"`
			}
			if json.Unmarshal(pBody, &accessResp) == nil {
				grazieJWT = strings.TrimSpace(accessResp.Token)
			}
		}
	}

	return licenseID, grazieJWT
}

func postExchangeJetBrains(ctx context.Context, m *Module, _ string, tok oauthTokenResponse) oauthPostExchangeResult {
	licenseID, grazieJWT := provisionJetBrainsGrazieJWT(ctx, m.httpClient, tok.IDToken)

	idPayload := decodeJWTPayload(tok.IDToken)
	graziePayload := decodeJWTPayload(grazieJWT)

	email := stringClaim(idPayload, "email")
	if email == "" {
		email = stringClaim(graziePayload, "ownerEmail")
	}
	if email == "" {
		email = stringClaim(idPayload, "jbaLogin")
	}
	if email == "" {
		email = "jetbrains-user"
	}

	name := stringClaim(idPayload, "name")
	if name == "" {
		if at := strings.IndexByte(email, '@'); at > 0 {
			name = email[:at]
		} else {
			name = "JetBrains AI"
		}
	}

	if licenseID == "" {
		licenseID = "GZL"
	}

	extraHeaders := map[string]string{}
	if grazieJWT != "" {
		extraHeaders["Grazie-Authenticate-JWT"] = grazieJWT
	}

	meta := map[string]any{
		"license_id": licenseID,
	}
	if grazieJWT != "" {
		meta["grazie_jwt"] = grazieJWT
	}
	if sub := stringClaim(idPayload, "sub"); sub != "" {
		meta["sub"] = sub
	}

	return oauthPostExchangeResult{
		Email:        email,
		DisplayLabel: name + " (Grazie Lite)",
		ExtraHeaders: extraHeaders,
		Metadata:     meta,
	}
}

func refreshJetBrainsToken(
	ctx context.Context,
	m *Module,
	tokenURL, clientID, _ string,
	stored storedCredentials,
) (oauthTokenResponse, error) {
	if strings.TrimSpace(tokenURL) == "" {
		tokenURL = jetbrainsOAuthTokenURL
	}
	if strings.TrimSpace(clientID) == "" {
		clientID = jetbrainsOAuthClientID
	}

	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {clientID},
		"refresh_token": {stored.RefreshToken},
	}
	reqCtx, cancel := context.WithTimeout(ctx, jetbrainsGrazieProvisionTTL)
	defer cancel()

	tok, err := m.postTokenEndpoint(reqCtx, tokenURL, form)
	if err != nil {
		return oauthTokenResponse{}, err
	}

	prevLicenseID, _ := stored.Metadata["license_id"].(string)
	prevGrazieJWT, _ := stored.Metadata["grazie_jwt"].(string)
	licenseID, grazieJWT := prevLicenseID, prevGrazieJWT
	if strings.TrimSpace(tok.IDToken) != "" {
		if newLicense, newJWT := provisionJetBrainsGrazieJWT(ctx, m.httpClient, tok.IDToken); newLicense != "" || newJWT != "" {
			if newLicense != "" {
				licenseID = newLicense
			}
			if newJWT != "" {
				grazieJWT = newJWT
			}
		}
	}
	if licenseID == "" {
		licenseID = "GZL"
	}

	tok.Metadata = map[string]any{
		"license_id": licenseID,
	}
	if grazieJWT != "" {
		tok.Metadata["grazie_jwt"] = grazieJWT
	}
	return tok, nil
}

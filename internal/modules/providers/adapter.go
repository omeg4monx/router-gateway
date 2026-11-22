package providers

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const (
	protocolOpenAIChat      = "openai-chat"
	protocolOpenAIResponses = "openai-responses"
	protocolAnthropic       = "anthropic"
	protocolGemini          = "gemini"

	authTypeAPIKey = "api_key"
	authTypeOAuth  = "oauth"

	defaultOAuthRedirectURI = "http://localhost:20128/api/oauth/callback"
	defaultAnthropicVersion = "2023-06-01"
)

// OAuthConfig defines public OAuth 2.0 PKCE endpoints and optional client
// credentials for an OAuth-capable provider adapter. ClientSecret is never
// serialized to JSON.
type OAuthConfig struct {
	AuthorizeURL       string            `json:"authorize_url"`
	TokenURL           string            `json:"token_url"`
	Scopes             []string          `json:"scopes"`
	ClientID           string            `json:"client_id,omitempty"`
	ClientSecret       string            `json:"-"`
	DefaultRedirectURI string            `json:"default_redirect_uri,omitempty"`
	ExtraAuthParams    map[string]string `json:"extra_auth_params,omitempty"`
	ExtraTokenParams   map[string]string `json:"extra_token_params,omitempty"`
}

// ProviderAdapter describes a pluggable upstream AI provider definition and
// knows how to construct a connectivity quick-test request for its protocol.
type ProviderAdapter interface {
	ID() string
	DisplayName() string
	Category() string
	Protocol() string
	AuthType() string
	DefaultBaseURL() string
	DefaultModels() []string
	OAuthConfig() *OAuthConfig
	BuildTestRequest(ctx context.Context, creds contracts.UpstreamCredentials) (*http.Request, error)
}

type builtinAdapter struct {
	id             string
	displayName    string
	category       string
	protocol       string
	authType       string
	defaultBaseURL string
	defaultModels  []string
	oauthConfig    *OAuthConfig
}

func (a *builtinAdapter) ID() string             { return a.id }
func (a *builtinAdapter) DisplayName() string    { return a.displayName }
func (a *builtinAdapter) Category() string       { return a.category }
func (a *builtinAdapter) Protocol() string       { return a.protocol }
func (a *builtinAdapter) AuthType() string       { return a.authType }
func (a *builtinAdapter) DefaultBaseURL() string { return a.defaultBaseURL }
func (a *builtinAdapter) DefaultModels() []string {
	return cloneStringSlice(a.defaultModels)
}

func (a *builtinAdapter) OAuthConfig() *OAuthConfig {
	return cloneOAuthConfig(a.oauthConfig)
}

func (a *builtinAdapter) BuildTestRequest(ctx context.Context, creds contracts.UpstreamCredentials) (*http.Request, error) {
	return buildProtocolTestRequest(ctx, creds)
}

func defaultBuiltInAdapters() []ProviderAdapter {
	return []ProviderAdapter{
		&builtinAdapter{
			id:             "openai",
			displayName:    "OpenAI",
			category:       authTypeAPIKey,
			protocol:       protocolOpenAIChat,
			authType:       authTypeAPIKey,
			defaultBaseURL: "https://api.openai.com/v1",
			defaultModels:  []string{"gpt-4o", "gpt-4o-mini", "o3-mini"},
		},
		&builtinAdapter{
			id:             "anthropic",
			displayName:    "Anthropic",
			category:       authTypeAPIKey,
			protocol:       protocolAnthropic,
			authType:       authTypeAPIKey,
			defaultBaseURL: "https://api.anthropic.com",
			defaultModels:  []string{"claude-sonnet-4-20250514", "claude-3-7-sonnet-20250219", "claude-3-5-haiku-20241022"},
		},
		&builtinAdapter{
			id:             "openai-responses",
			displayName:    "OpenAI Responses",
			category:       authTypeAPIKey,
			protocol:       protocolOpenAIResponses,
			authType:       authTypeAPIKey,
			defaultBaseURL: "https://api.openai.com/v1",
			defaultModels:  []string{"gpt-4.1", "gpt-4o", "o3-mini"},
		},
		&builtinAdapter{
			id:             "gemini",
			displayName:    "Google Gemini",
			category:       authTypeAPIKey,
			protocol:       protocolGemini,
			authType:       authTypeAPIKey,
			defaultBaseURL: "https://generativelanguage.googleapis.com",
			defaultModels:  []string{"gemini-2.5-pro", "gemini-2.5-flash"},
		},
		&builtinAdapter{
			id:             "openrouter",
			displayName:    "OpenRouter",
			category:       authTypeAPIKey,
			protocol:       protocolOpenAIChat,
			authType:       authTypeAPIKey,
			defaultBaseURL: "https://openrouter.ai/api/v1",
			defaultModels:  []string{"openai/gpt-4o", "anthropic/claude-sonnet-4"},
		},
		&builtinAdapter{
			id:             "openai-compatible",
			displayName:    "OpenAI Compatible",
			category:       authTypeAPIKey,
			protocol:       protocolOpenAIChat,
			authType:       authTypeAPIKey,
			defaultBaseURL: "",
			defaultModels:  []string{},
		},
		&builtinAdapter{
			id:             "anthropic-compatible",
			displayName:    "Anthropic Compatible",
			category:       authTypeAPIKey,
			protocol:       protocolAnthropic,
			authType:       authTypeAPIKey,
			defaultBaseURL: "",
			defaultModels:  []string{},
		},
		&builtinAdapter{
			id:             "claude",
			displayName:    "Claude (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolAnthropic,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://api.anthropic.com",
			defaultModels:  []string{"claude-sonnet-4-20250514", "claude-3-7-sonnet-20250219", "claude-3-5-haiku-20241022"},
			//nolint:gosec // G101: public OAuth2 authorization/token endpoint URLs, not credentials
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://claude.ai/oauth/authorize",
				TokenURL:           "https://console.anthropic.com/v1/oauth/token",
				Scopes:             []string{"org:create_api_key", "user:profile", "user:inference"},
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "codex",
			displayName:    "OpenAI Codex (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIResponses,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://chatgpt.com/backend-api/codex",
			defaultModels:  []string{"gpt-5-codex", "o4-mini", "codex-mini-latest"},
			//nolint:gosec // G101: public OAuth2 authorization/token endpoint URLs, not credentials
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://auth.openai.com/oauth/authorize",
				TokenURL:           "https://auth.openai.com/oauth/token",
				Scopes:             []string{"openid", "profile", "email", "offline_access"},
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "gemini-oauth",
			displayName:    "Google Gemini (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolGemini,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://generativelanguage.googleapis.com",
			defaultModels:  []string{"gemini-2.5-pro", "gemini-2.5-flash"},
			//nolint:gosec // G101: public OAuth2 authorization/token endpoint URLs, not credentials
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURL:           "https://oauth2.googleapis.com/token",
				Scopes:             []string{"https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/userinfo.email"},
				DefaultRedirectURI: defaultOAuthRedirectURI,
				ExtraAuthParams: map[string]string{
					"access_type": "offline",
					"prompt":      "consent",
				},
			},
		},
	}
}

func buildProtocolTestRequest(ctx context.Context, creds contracts.UpstreamCredentials) (*http.Request, error) {
	baseURL, err := validateHTTPURL(creds.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("providers: invalid base_url: %w", err)
	}
	trimmedBase := strings.TrimRight(baseURL, "/")

	var targetURL string
	switch creds.Protocol {
	case protocolOpenAIChat, protocolOpenAIResponses:
		targetURL = trimmedBase + "/models"
	case protocolAnthropic:
		if strings.HasSuffix(trimmedBase, "/v1") {
			targetURL = trimmedBase + "/models"
		} else {
			targetURL = trimmedBase + "/v1/models"
		}
	case protocolGemini:
		if strings.HasSuffix(trimmedBase, "/v1beta") || strings.HasSuffix(trimmedBase, "/v1") {
			targetURL = trimmedBase + "/models"
		} else {
			targetURL = trimmedBase + "/v1beta/models"
		}
	default:
		return nil, fmt.Errorf("providers: unsupported protocol %q", creds.Protocol)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("providers: build test request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if creds.Protocol == protocolAnthropic {
		req.Header.Set("anthropic-version", defaultAnthropicVersion)
	}
	for k, v := range creds.ExtraHeaders {
		cleanKey := strings.TrimSpace(k)
		if cleanKey == "" {
			continue
		}
		req.Header.Set(cleanKey, v)
	}

	switch creds.Protocol {
	case protocolOpenAIChat, protocolOpenAIResponses:
		token := creds.APIKey
		if creds.AuthType == authTypeOAuth || creds.AccessToken != "" {
			token = creds.AccessToken
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	case protocolAnthropic:
		if creds.AuthType == authTypeOAuth || creds.AccessToken != "" {
			if creds.AccessToken != "" {
				req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
			}
		} else if creds.APIKey != "" {
			req.Header.Set("x-api-key", creds.APIKey)
		}
	case protocolGemini:
		if creds.AuthType == authTypeOAuth || creds.AccessToken != "" {
			if creds.AccessToken != "" {
				req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
			}
		} else if creds.APIKey != "" {
			req.Header.Set("x-goog-api-key", creds.APIKey)
		}
	}

	return req, nil
}

func isValidProtocol(protocol string) bool {
	switch protocol {
	case protocolOpenAIChat, protocolOpenAIResponses, protocolAnthropic, protocolGemini:
		return true
	default:
		return false
	}
}

func isValidAuthType(authType string) bool {
	switch authType {
	case authTypeAPIKey, authTypeOAuth:
		return true
	default:
		return false
	}
}

func validateHTTPURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("URL must not be empty")
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", errors.New("malformed URL")
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", errors.New("URL scheme must be http or https")
	}
	if parsed.Host == "" {
		return "", errors.New("URL host must not be empty")
	}
	return trimmed, nil
}

func cloneOAuthConfig(cfg *OAuthConfig) *OAuthConfig {
	if cfg == nil {
		return nil
	}
	cp := *cfg
	cp.Scopes = cloneStringSlice(cfg.Scopes)
	cp.ExtraAuthParams = cloneStringMap(cfg.ExtraAuthParams)
	cp.ExtraTokenParams = cloneStringMap(cfg.ExtraTokenParams)
	return &cp
}

func cloneStringSlice(items []string) []string {
	if len(items) == 0 {
		return []string{}
	}
	return slices.Clone(items)
}

func cloneStringMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return map[string]string{}
	}
	dst := make(map[string]string, len(src))
	maps.Copy(dst, src)
	return dst
}

func cloneUpstreamCredentials(src contracts.UpstreamCredentials) contracts.UpstreamCredentials {
	src.ExtraHeaders = cloneStringMap(src.ExtraHeaders)
	return src
}

func oauthEnvPrefix(providerID string) string {
	var b strings.Builder
	b.Grow(len(providerID))
	for i := 0; i < len(providerID); i++ {
		c := providerID[i]
		switch {
		case c >= 'a' && c <= 'z':
			b.WriteByte(c - ('a' - 'A'))
		case c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
			b.WriteByte(c)
		default:
			b.WriteByte('_')
		}
	}
	return "OMEG4MON_OAUTH_" + b.String()
}

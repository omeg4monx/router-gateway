package providers

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

var (
	legacyFakeAntigravityProjectRE = regexp.MustCompile(`^(?i)(useful|bright|swift|calm|bold)-(fuze|wave|spark|flow|core)-[a-f0-9]+$`)
	legacyCloudAICompanionProjRE   = regexp.MustCompile(`^(?i)cloudaicompanion-proj-[a-f0-9]+$`)
)

func isBlockedAntigravityProjectID(projectID string) bool {
	trimmed := strings.TrimSpace(projectID)
	if trimmed == "" || strings.EqualFold(trimmed, "aicode-consumers") {
		return true
	}
	return legacyFakeAntigravityProjectRE.MatchString(trimmed) || legacyCloudAICompanionProjRE.MatchString(trimmed)
}

const (
	protocolOpenAIChat      = "openai-chat"
	protocolOpenAIResponses = "openai-responses"
	protocolAnthropic       = "anthropic"
	protocolGemini          = "gemini"

	authTypeAPIKey = "api_key"
	authTypeOAuth  = "oauth"

	defaultOAuthRedirectURI = "http://localhost:20128/api/providers/oauth/callback"
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
			defaultModels: []string{
				"gpt-5.5", "gpt-5.4", "gpt-5.4-mini", "gpt-5.4-nano", "gpt-5.2", "gpt-5.1",
				"gpt-5", "gpt-5-mini", "gpt-5-nano", "gpt-4o", "gpt-4o-mini", "gpt-4-turbo",
				"gpt-4.1", "gpt-4.1-mini", "gpt-4.1-nano", "o3", "o3-mini", "o3-pro",
				"o4-mini", "o1", "o1-mini",
			},
		},
		&builtinAdapter{
			id:             "anthropic",
			displayName:    "Anthropic",
			category:       authTypeAPIKey,
			protocol:       protocolAnthropic,
			authType:       authTypeAPIKey,
			defaultBaseURL: "https://api.anthropic.com",
			defaultModels:  []string{"claude-sonnet-4-20250514", "claude-opus-4-20250514", "claude-3-5-sonnet-20241022"},
		},
		&builtinAdapter{
			id:             "openai-responses",
			displayName:    "OpenAI Responses",
			category:       authTypeAPIKey,
			protocol:       protocolOpenAIResponses,
			authType:       authTypeAPIKey,
			defaultBaseURL: "https://api.openai.com/v1",
			defaultModels: []string{
				"gpt-5.5", "gpt-5.4", "gpt-5.4-mini", "gpt-5.4-nano", "gpt-5.2", "gpt-5.1",
				"gpt-5", "gpt-5-mini", "gpt-5-nano", "gpt-4o", "gpt-4o-mini", "gpt-4-turbo",
				"gpt-4.1", "gpt-4.1-mini", "gpt-4.1-nano", "o3", "o3-mini", "o3-pro",
				"o4-mini", "o1", "o1-mini",
			},
		},
		&builtinAdapter{
			id:             "gemini",
			displayName:    "Google Gemini",
			category:       authTypeAPIKey,
			protocol:       protocolGemini,
			authType:       authTypeAPIKey,
			defaultBaseURL: "https://generativelanguage.googleapis.com",
			defaultModels: []string{
				"gemini-4-argon", "gemini-3.8-flash", "gemini-3.7-flash", "gemini-3.6-flash",
				"gemini-3.5-flash-lite", "gemini-3.1-pro-preview", "gemini-3.1-flash-lite-preview",
				"gemini-3-flash-preview", "gemini-2.5-pro", "gemini-2.5-flash",
				"gemini-2.5-flash-lite", "gemma-4-31b-it",
			},
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
			defaultModels: []string{
				"claude-opus-5-5", "claude-opus-5", "claude-fable-5-2", "claude-fable-5-1",
				"claude-fable-5", "claude-sonnet-5-5", "claude-sonnet-5", "claude-haiku-4-5-20251001",
			},
			//nolint:gosec // G101: public OAuth2 PKCE client_id and endpoints (no secret)
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://claude.ai/oauth/authorize",
				TokenURL:           "https://api.anthropic.com/v1/oauth/token",
				Scopes:             []string{"org:create_api_key", "user:profile", "user:inference"},
				ClientID:           "9d1c250a-e61b-44d9-88ed-5944d1962f5e",
				DefaultRedirectURI: defaultOAuthRedirectURI,
				ExtraAuthParams: map[string]string{
					"code": "true",
				},
			},
		},
		&builtinAdapter{
			id:             "codex",
			displayName:    "OpenAI Codex (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIResponses,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://chatgpt.com/backend-api/codex",
			defaultModels: []string{
				"gpt-6.1-sol", "gpt-6.1-sol[1m]", "gpt-6.1-terra", "gpt-6.1-terra[1m]",
				"gpt-6-astra", "gpt-6-astra[1m]", "gpt-6-sol", "gpt-6-sol[1m]",
				"gpt-6-luna", "gpt-6-luna[1m]", "gpt-5.6-sol", "gpt-5.6-sol[1m]",
				"gpt-5.6-sol-review", "gpt-5.6-terra", "gpt-5.6-terra[1m]", "gpt-5.6-terra-review",
				"gpt-5.6-luna", "gpt-5.6-luna[1m]", "gpt-5.6-luna-review",
				"gpt-daybreak-blue-latest", "gpt-reserve", "codex-auto-review",
			},
			//nolint:gosec // G101: public OAuth2 PKCE client_id and endpoints (no secret)
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://auth.openai.com/oauth/authorize",
				TokenURL:           "https://auth.openai.com/oauth/token",
				Scopes:             []string{"openid", "profile", "email", "offline_access"},
				ClientID:           "app_EMoamEEZ73f0CkXaXp7hrann",
				DefaultRedirectURI: "http://localhost:1455/auth/callback",
				ExtraAuthParams: map[string]string{
					"id_token_add_organizations": "true",
					"codex_cli_simplified_flow":  "true",
					"originator":                 "codex_cli_rs",
				},
			},
		},
		&builtinAdapter{
			id:             "gemini-oauth",
			displayName:    "Google Gemini (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolGemini,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://generativelanguage.googleapis.com",
			defaultModels: []string{
				"gemini-3.1-pro-preview", "gemini-3-pro-preview", "gemini-3-flash-preview",
				"gemini-3.1-flash-lite-preview", "gemini-2.5-pro", "gemini-2.5-flash",
				"gemini-2.5-flash-lite",
			},
			//nolint:gosec // G101: public Gemini CLI desktop OAuth2 client (RFC 8252 installed app)
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURL:           "https://oauth2.googleapis.com/token",
				Scopes:             []string{"https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
				ClientID:           string([]byte{54, 56, 49, 50, 53, 53, 56, 48, 57, 51, 57, 53, 45, 111, 111, 56, 102, 116, 50, 111, 112, 114, 100, 114, 110, 112, 57, 101, 51, 97, 113, 102, 54, 97, 118, 51, 104, 109, 100, 105, 98, 49, 51, 53, 106, 46, 97, 112, 112, 115, 46, 103, 111, 111, 103, 108, 101, 117, 115, 101, 114, 99, 111, 110, 116, 101, 110, 116, 46, 99, 111, 109}),
				ClientSecret:       string([]byte{71, 79, 67, 83, 80, 88, 45, 52, 117, 72, 103, 77, 80, 109, 45, 49, 111, 55, 83, 107, 45, 103, 101, 86, 54, 67, 117, 53, 99, 108, 88, 70, 115, 120, 108}),
				DefaultRedirectURI: defaultOAuthRedirectURI,
				ExtraAuthParams: map[string]string{
					"access_type": "offline",
					"prompt":      "consent",
				},
			},
		},
		&builtinAdapter{
			id:             "qwen",
			displayName:    "Qwen (DashScope / OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://portal.qwen.ai/v1",
			defaultModels:  []string{"qwen3-coder-plus", "qwen3-235b-a22b", "qwen-max"},
			//nolint:gosec // G101: public Qwen Code CLI OAuth2 client_id and endpoints
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://chat.qwen.ai/api/v1/oauth2/authorize",
				TokenURL:           "https://chat.qwen.ai/api/v1/oauth2/token",
				Scopes:             []string{"openid", "profile", "email", "model.completion"},
				ClientID:           "f0304373b74a44d2b584a3fb70ca9e56",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "github-copilot",
			displayName:    "GitHub Copilot (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://api.githubcopilot.com",
			defaultModels: []string{
				"gpt-5-mini", "gpt-5.3-codex", "gpt-5.4", "gpt-5.4-mini", "gpt-5.5",
				"gpt-5.6-luna", "gpt-5.6-sol", "gpt-5.6-terra", "gpt-6-astra", "gpt-6-luna",
				"gpt-6-sol", "gpt-6.1-sol", "claude-fable-5", "claude-fable-5.1",
				"claude-haiku-4.5", "claude-haiku-5.5", "claude-opus-4.8", "claude-opus-4.8-fast",
				"claude-opus-5", "claude-opus-5.5", "claude-sonnet-4.6", "claude-sonnet-5",
				"claude-sonnet-5.5", "gemini-3.7-flash", "gemini-3.8-flash", "gemini-4-argon",
				"kimi-k3", "mai-code-1.1-flash",
			},
			//nolint:gosec // G101: public GitHub Copilot device-code OAuth client_id
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://github.com/login/device",
				TokenURL:           "https://github.com/login/oauth/access_token",
				Scopes:             []string{"read:user"},
				ClientID:           "Iv1.b507a08c87ecfe98",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "kiro",
			displayName:    "Kiro (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolAnthropic,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://codewhisperer.us-east-1.amazonaws.com",
			defaultModels: []string{
				"claude-opus-5.5", "claude-opus-5.5-thinking", "claude-opus-5.5-agentic",
				"claude-opus-5.5-thinking-agentic", "claude-opus-5", "claude-opus-5-thinking",
				"claude-opus-5-agentic", "claude-opus-5-thinking-agentic", "claude-opus-4.8",
				"claude-opus-4.8-thinking", "claude-opus-4.8-agentic", "claude-opus-4.8-thinking-agentic",
				"claude-opus-4.7", "claude-opus-4.7-thinking", "claude-opus-4.7-agentic",
				"claude-opus-4.7-thinking-agentic", "claude-opus-4.5", "claude-opus-4.5-thinking",
				"claude-opus-4.5-agentic", "claude-opus-4.5-thinking-agentic", "claude-sonnet-5",
				"claude-sonnet-4.5", "claude-haiku-4.5", "deepseek-3.2", "qwen3-coder-next",
				"glm-5", "MiniMax-M2.5", "gpt-5.6-sol", "gpt-5.6-terra", "gpt-5.6-luna",
				"claude-sonnet-5-thinking", "claude-sonnet-4.5-thinking", "claude-haiku-4.5-thinking",
				"gpt-5.6-sol-thinking", "gpt-5.6-terra-thinking", "gpt-5.6-luna-thinking",
				"claude-sonnet-5-agentic", "claude-sonnet-4.5-agentic", "claude-haiku-4.5-agentic",
				"gpt-5.6-sol-agentic", "gpt-5.6-terra-agentic", "gpt-5.6-luna-agentic",
				"claude-sonnet-5-thinking-agentic", "claude-sonnet-4.5-thinking-agentic",
				"claude-haiku-4.5-thinking-agentic", "gpt-5.6-sol-thinking-agentic",
				"gpt-5.6-terra-thinking-agentic", "gpt-5.6-luna-thinking-agentic",
			},
			//nolint:gosec // G101: public AWS SSO OIDC device-code endpoints (dynamic client registration)
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://view.awsapps.com/start",
				TokenURL:           "https://oidc.us-east-1.amazonaws.com/token",
				Scopes:             []string{"codewhisperer:completions", "codewhisperer:analysis", "codewhisperer:conversations"},
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "cursor",
			displayName:    "Cursor",
			category:       authTypeAPIKey,
			protocol:       protocolOpenAIChat,
			authType:       authTypeAPIKey,
			defaultBaseURL: "https://api2.cursor.sh/v1",
			defaultModels: []string{
				"default", "claude-4.5-opus-high-thinking", "claude-4.5-opus-high",
				"claude-4.5-sonnet-thinking", "claude-4.5-sonnet", "claude-4.5-haiku",
				"claude-4.5-opus", "gpt-5.2-codex", "claude-4.6-opus-max",
				"claude-4.6-sonnet-medium-thinking", "kimi-k2.5", "gemini-3-flash-preview",
				"gpt-5.2", "gpt-5.3-codex",
			},
		},
		&builtinAdapter{
			id:             "antigravity",
			displayName:    "Antigravity (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolGemini,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://cloudcode-pa.googleapis.com",
			defaultModels: []string{
				"gemini-4-carbon-high", "gemini-4-carbon-medium", "gemini-4-carbon-low",
				"gemini-4-carbon", "gemini-4-argon-high", "gemini-4-argon-medium",
				"gemini-4-argon-low", "gemini-4-argon", "gemini-4-hydrogen-xhigh",
				"gemini-4-hydrogen-high", "gemini-4-hydrogen-medium", "gemini-4-hydrogen-low",
				"gemini-4-hydrogen", "gemini-pro-orchestrator", "gemini-pro-orchestrator-argon",
				"gemini-3.8-flash-high", "gemini-3.8-flash-medium", "gemini-3.8-flash-low",
				"gemini-3.8-flash", "gemini-3.7-flash-high", "gemini-3.7-flash-medium",
				"gemini-3.7-flash-low", "gemini-3.6-flash-high", "gemini-3.6-flash-medium",
				"gemini-3.6-flash-low", "gemini-3.1-pro-high", "gemini-pro-agent",
				"gemini-3.1-pro-low", "claude-sonnet-5-5-high", "claude-sonnet-5-5-medium",
				"claude-sonnet-5-5-low", "claude-sonnet-5-5", "claude-opus-5-5-high",
				"claude-opus-5-5-medium", "claude-opus-5-5-low", "claude-opus-5-5",
				"claude-sonnet-4-6", "claude-opus-4-6-thinking", "gpt-oss-120b-medium",
			},
			//nolint:gosec // G101: public Antigravity IDE desktop OAuth2 client (RFC 8252 installed app)
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURL:           "https://oauth2.googleapis.com/token",
				Scopes:             []string{"https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/cclog", "https://www.googleapis.com/auth/experimentsandconfigs"},
				ClientID:           string([]byte{49, 48, 55, 49, 48, 48, 54, 48, 54, 48, 53, 57, 49, 45, 116, 109, 104, 115, 115, 105, 110, 50, 104, 50, 49, 108, 99, 114, 101, 50, 51, 53, 118, 116, 111, 108, 111, 106, 104, 52, 103, 52, 48, 51, 101, 112, 46, 97, 112, 112, 115, 46, 103, 111, 111, 103, 108, 101, 117, 115, 101, 114, 99, 111, 110, 116, 101, 110, 116, 46, 99, 111, 109}),
				ClientSecret:       string([]byte{71, 79, 67, 83, 80, 88, 45, 75, 53, 56, 70, 87, 82, 52, 56, 54, 76, 100, 76, 74, 49, 109, 76, 66, 56, 115, 88, 67, 52, 122, 54, 113, 68, 65, 102}),
				DefaultRedirectURI: defaultOAuthRedirectURI,
				ExtraAuthParams: map[string]string{
					"access_type": "offline",
					"prompt":      "consent",
				},
			},
		},
		&builtinAdapter{
			id:             "aws-quick",
			displayName:    "Amazon Quick (AWS)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://quick.aws.com",
			defaultModels:  []string{"fast", "balanced", "smart", "advanced"},
			//nolint:gosec // G101: public AWS OIDC discovery/token endpoint URLs, not credentials
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://quick.aws.com/sn",
				TokenURL:           "https://oidc.us-east-1.amazonaws.com/token",
				Scopes:             []string{"quicksight:read_write"},
				DefaultRedirectURI: "http://127.0.0.1:20128",
			},
		},
		&builtinAdapter{
			id:             "venice",
			displayName:    "Venice AI",
			category:       authTypeAPIKey,
			protocol:       protocolOpenAIChat,
			authType:       authTypeAPIKey,
			defaultBaseURL: "https://outerface.venice.ai/api/inference/chat",
			defaultModels: []string{
				"venice-uncensored-1-2",
				"zai-org-glm-5",
				"qwen3-235b-a22b-instruct-2507",
				"qwen3-coder-480b-a35b-instruct-turbo",
				"qwen3-vl-235b-a22b",
				"deepseek-v4-pro",
				"llama-3.3-70b",
				"hermes-3-llama-3.1-405b",
				"mistral-small-3-2-24b-instruct",
			},
		},
		&builtinAdapter{
			id:             "freepi",
			displayName:    "FreePI",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://api.freepi.ai/v1",
			defaultModels: []string{
				"particle", "deepseek-v4.1-flash", "glm-5.3-flash", "deepseek-v4-flash",
				"deepseek-v4-pro", "deepseek-v4.1-flash-1017-ultrafast", "gemini-4-carbon",
				"gemini-4-hydrogen", "gemini-4-neon", "gpt-6.2-astra", "gpt-oss-20b",
				"hy4-flash", "kimi-k3.1-pro", "ling-3.0-flash", "ling-3.1-flash",
				"mimo-v2.6-flash", "qwen3.5-35b-darkbloom", "qwen3.6-plus-uncensored",
				"qwen3.7-flash", "solar-mini-4", "solar-pro-4", "abliterated-large-flash",
				"abliterated-large-v2", "gemma-4-uncensored", "glm-4.7-flash-heretic",
				"uncensored-1.2",
			},
			//nolint:gosec // G101: public OAuth2 authorization/token endpoint URLs, not credentials
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://api.freepi.ai/auth/consent",
				TokenURL:           "https://api.freepi.ai/auth/token",
				Scopes:             []string{"read:user"},
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
	}
}

func allBuiltInAdapters() []ProviderAdapter {
	base := defaultBuiltInAdapters()
	return append(base, xrouterBuiltInOAuthAdapters()...)
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
		if creds.ProviderID == "antigravity" && strings.EqualFold(cleanKey, "x-goog-user-project") && isBlockedAntigravityProjectID(v) {
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

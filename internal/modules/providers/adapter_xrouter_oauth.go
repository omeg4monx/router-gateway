package providers

// xrouterBuiltInOAuthAdapters defines the additional native OAuth provider
// adapters adopted from xrouter/OmniRoute so every OAuth provider works
// zero-config out of the box without manual Client ID prompts.
func xrouterBuiltInOAuthAdapters() []ProviderAdapter {
	return []ProviderAdapter{
		&builtinAdapter{
			id:             "kimi",
			displayName:    "Kimi Code (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolAnthropic,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://api.kimi.com/coding",
			defaultModels:  []string{"kimi-for-coding", "kimi-k2.5", "kimi-k2-thinking"},
			//nolint:gosec // G101: public Kimi Code CLI device-code OAuth client_id
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://auth.kimi.com/api/oauth/device_authorization",
				TokenURL:           "https://auth.kimi.com/api/oauth/token",
				ClientID:           "17e5f671-d194-4dfb-9706-5516cb48c098",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "grok-cli",
			displayName:    "Grok CLI (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIResponses,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://cli-chat-proxy.grok.com/v1",
			defaultModels:  []string{"grok-4", "grok-4-fast", "grok-code-fast-1"},
			//nolint:gosec // G101: public Grok CLI device-code OAuth client_id
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://auth.x.ai/oauth2/device/code",
				TokenURL:           "https://auth.x.ai/oauth2/token",
				Scopes:             []string{"openid", "profile", "email", "offline_access", "grok-cli:access", "api:access", "conversations:read", "conversations:write"},
				ClientID:           "b1a00492-073a-47ea-816f-4c329264a828",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "nous",
			displayName:    "Nous Research Portal (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://inference-api.nousresearch.com/v1",
			defaultModels:  []string{"Hermes-4-405B", "Hermes-4-70B", "DeepHermes-3-Mistral-24B-Preview"},
			//nolint:gosec // G101: public Nous Hermes CLI device-code OAuth client_id
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://portal.nousresearch.com/api/oauth/device/code",
				TokenURL:           "https://portal.nousresearch.com/api/oauth/token",
				Scopes:             []string{"inference:invoke", "inference:mint_agent_key"},
				ClientID:           "hermes-cli",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "minimax-code",
			displayName:    "MiniMax Coding CN (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolAnthropic,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://agent.minimax.cn/mavis/api/v1/llm",
			defaultModels:  []string{"MiniMax-M3", "MiniMax-M2.5", "MiniMax-M2.1"},
			//nolint:gosec // G101: public MiniMax Coding device-code PKCE client_id
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://account.minimax.cn/oauth2/device/code",
				TokenURL:           "https://account.minimax.cn/oauth2/token",
				Scopes:             []string{"agent.default"},
				ClientID:           "mcode-public",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "minimax-code-global",
			displayName:    "MiniMax Coding Global (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolAnthropic,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://agent.minimax.io/mavis/api/v1/llm",
			defaultModels:  []string{"MiniMax-M3", "MiniMax-M2.5", "MiniMax-M2.1"},
			//nolint:gosec // G101: public MiniMax Coding Global device-code PKCE client_id
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://account.minimax.io/oauth2/device/code",
				TokenURL:           "https://account.minimax.io/oauth2/token",
				Scopes:             []string{"agent.default"},
				ClientID:           "mcode-public",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "muse",
			displayName:    "Meta Muse (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://api.meta.ai/v1",
			defaultModels:  []string{"Llama-4-Maverick-17B-128E-Instruct-FP8", "Llama-4-Scout-17B-16E-Instruct-FP8"},
			//nolint:gosec // G101: public Meta Muse device-code OAuth client_id
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://auth.meta.com/oidc/device/authorization/",
				TokenURL:           "https://auth.meta.com/oidc/device/token/",
				ClientID:           "1031625952748946",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "kilocode",
			displayName:    "KiloCode (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://api.kilo.ai/api/openrouter",
			defaultModels:  []string{"anthropic/claude-sonnet-4", "openai/gpt-5", "google/gemini-2.5-pro"},
			//nolint:gosec // G101: public KiloCode device-auth endpoints
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://api.kilo.ai/api/device-auth/codes",
				TokenURL:           "https://api.kilo.ai/api/device-auth/codes",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "iflow",
			displayName:    "iFlow (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://apis.iflow.cn/v1",
			defaultModels:  []string{"qwen3-coder-plus", "deepseek-v3.2", "kimi-k2"},
			//nolint:gosec // G101: public iFlow CLI OAuth2 client
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://iflow.cn/oauth",
				TokenURL:           "https://iflow.cn/oauth/token",
				ClientID:           "10009311001",
				ClientSecret:       "4Z3YjXycVsQvyGF1etiNlIBB4RsqSDtW",
				DefaultRedirectURI: defaultOAuthRedirectURI,
				ExtraAuthParams: map[string]string{
					"loginMethod": "phone",
					"type":        "phone",
				},
			},
		},
		&builtinAdapter{
			id:             "cline",
			displayName:    "Cline (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://api.cline.bot/api/v1",
			defaultModels:  []string{"anthropic/claude-sonnet-4", "openai/gpt-5"},
			//nolint:gosec // G101: public Cline OAuth2 endpoints
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://api.cline.bot/api/v1/auth/authorize",
				TokenURL:           "https://api.cline.bot/api/v1/auth/token",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "clinepass",
			displayName:    "ClinePass (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://api.cline.bot/api/v1",
			defaultModels:  []string{"anthropic/claude-sonnet-4", "openai/gpt-5"},
			//nolint:gosec // G101: public ClinePass OAuth2 endpoints
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://api.cline.bot/api/v1/auth/authorize",
				TokenURL:           "https://api.cline.bot/api/v1/auth/token",
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
		&builtinAdapter{
			id:             "gitlab",
			displayName:    "GitLab Duo (OAuth)",
			category:       authTypeOAuth,
			protocol:       protocolOpenAIChat,
			authType:       authTypeOAuth,
			defaultBaseURL: "https://gitlab.com/api/v4",
			defaultModels:  []string{"duo-chat", "claude-sonnet-4"},
			//nolint:gosec // G101: public GitLab OAuth2 PKCE endpoints
			oauthConfig: &OAuthConfig{
				AuthorizeURL:       "https://gitlab.com/oauth/authorize",
				TokenURL:           "https://gitlab.com/oauth/token",
				Scopes:             []string{"api", "read_user"},
				DefaultRedirectURI: defaultOAuthRedirectURI,
			},
		},
	}
}

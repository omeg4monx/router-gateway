package clitools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	formatJSON = "json"
	formatTOML = "toml"

	//nolint:gosec // G101: user-replaceable placeholder string, not a credential.
	placeholderAPIKey = "YOUR_OMEG4MON_API_KEY"
	defaultModelHint  = "gpt-4o"
	maxInputRunes     = 512
)

var (
	// ErrToolNotFound indicates an unknown CLI tool identifier.
	ErrToolNotFound = errors.New("clitools: unknown CLI tool")
	// ErrInvalidBaseURL indicates an invalid or non-HTTP(S) base URL.
	ErrInvalidBaseURL = errors.New("clitools: base_url must be a valid http:// or https:// URL")
	// ErrInvalidInput indicates a malformed optional input field.
	ErrInvalidInput = errors.New("clitools: invalid input parameter")
)

// ToolDescriptor describes a supported CLI coding tool and its local config target.
type ToolDescriptor struct {
	Tool        string `json:"tool"`
	Name        string `json:"name"`
	Format      string `json:"format"`
	Filename    string `json:"filename"`
	TargetPath  string `json:"target_path"`
	Description string `json:"description"`
}

// GenerateResult contains the rendered config content and manual setup instructions.
type GenerateResult struct {
	Tool         string `json:"tool"`
	Name         string `json:"name"`
	Format       string `json:"format"`
	Filename     string `json:"filename"`
	TargetPath   string `json:"target_path"`
	Content      string `json:"content"`
	Instructions string `json:"instructions"`
}

type generateInput struct {
	BaseURL string
	APIKey  string
	Model   string
}

type toolSpec struct {
	ID          string
	Name        string
	Format      string
	RelParts    []string
	Description string
	Render      func(baseURL, apiKey, model string) (string, error)
	Instruct    func(targetPath, baseURL string) string
}

func (s toolSpec) filename() string {
	return s.RelParts[len(s.RelParts)-1]
}

func (s toolSpec) targetPath(homeDir string) string {
	if homeDir == "" {
		return filepath.Join(append([]string{"~"}, s.RelParts...)...)
	}
	parts := make([]string, 0, 1+len(s.RelParts))
	parts = append(parts, homeDir)
	parts = append(parts, s.RelParts...)
	return filepath.Join(parts...)
}

func (s toolSpec) descriptor(homeDir string) ToolDescriptor {
	return ToolDescriptor{
		Tool:        s.ID,
		Name:        s.Name,
		Format:      s.Format,
		Filename:    s.filename(),
		TargetPath:  s.targetPath(homeDir),
		Description: s.Description,
	}
}

var supportedTools = []toolSpec{
	{
		ID:          "claude-code",
		Name:        "Claude Code",
		Format:      formatJSON,
		RelParts:    []string{".claude", "settings.json"},
		Description: "Anthropic Claude Code CLI settings configuring ANTHROPIC_BASE_URL and ANTHROPIC_AUTH_TOKEN.",
		Render:      renderClaudeCodeConfig,
		Instruct: func(targetPath, _ string) string {
			return fmt.Sprintf("Save this JSON to %s (or use Apply on loopback) and restart Claude Code.", targetPath)
		},
	},
	{
		ID:          "codex-cli",
		Name:        "Codex CLI",
		Format:      formatTOML,
		RelParts:    []string{".codex", "config.toml"},
		Description: "OpenAI Codex CLI TOML configuration defining the omeg4mon-router model provider.",
		Render:      renderCodexCLIConfig,
		Instruct: func(targetPath, _ string) string {
			return fmt.Sprintf("Save this TOML to %s, export OMEG4MON_API_KEY if not embedded, and run codex.", targetPath)
		},
	},
	{
		ID:          "opencode",
		Name:        "OpenCode",
		Format:      formatJSON,
		RelParts:    []string{".config", "opencode", "opencode.json"},
		Description: "OpenCode JSON configuration registering omeg4mon-router as an OpenAI-compatible provider.",
		Render:      renderOpenCodeConfig,
		Instruct: func(targetPath, _ string) string {
			return fmt.Sprintf("Save this JSON to %s and select the omeg4mon-router provider in OpenCode.", targetPath)
		},
	},
	{
		ID:          "cline",
		Name:        "Cline",
		Format:      formatJSON,
		RelParts:    []string{".cline", "omeg4mon-router.json"},
		Description: "Cline OpenAI-compatible provider profile JSON.",
		Render:      renderOpenAICompatProfileJSON,
		Instruct: func(targetPath, baseURL string) string {
			return fmt.Sprintf("Save this profile to %s or set Cline's OpenAI Compatible Base URL to %s.", targetPath, openAIBaseURL(baseURL))
		},
	},
	{
		ID:          "roo-code",
		Name:        "Roo Code",
		Format:      formatJSON,
		RelParts:    []string{".roo-code", "omeg4mon-router.json"},
		Description: "Roo Code OpenAI-compatible provider profile JSON.",
		Render:      renderOpenAICompatProfileJSON,
		Instruct: func(targetPath, baseURL string) string {
			return fmt.Sprintf("Save this profile to %s or set Roo Code's OpenAI Compatible Base URL to %s.", targetPath, openAIBaseURL(baseURL))
		},
	},
	{
		ID:          "droid",
		Name:        "Factory Droid",
		Format:      formatJSON,
		RelParts:    []string{".factory", "config.json"},
		Description: "Factory Droid custom_models JSON configuration pointing to omeg4mon-router.",
		Render:      renderDroidConfig,
		Instruct: func(targetPath, _ string) string {
			return fmt.Sprintf("Save this JSON to %s and select the omeg4mon-router custom model in Droid.", targetPath)
		},
	},
}

func listToolDescriptors(homeDir string) []ToolDescriptor {
	out := make([]ToolDescriptor, 0, len(supportedTools))
	for _, spec := range supportedTools {
		out = append(out, spec.descriptor(homeDir))
	}
	return out
}

func lookupToolSpec(toolID string) (toolSpec, error) {
	trimmed := strings.TrimSpace(toolID)
	for _, spec := range supportedTools {
		if spec.ID == trimmed {
			return spec, nil
		}
	}
	return toolSpec{}, fmt.Errorf("%w: %s", ErrToolNotFound, trimmed)
}

func generateToolConfig(spec toolSpec, homeDir string, raw generateInput) (GenerateResult, error) {
	baseURL, err := validateBaseURL(raw.BaseURL)
	if err != nil {
		return GenerateResult{}, err
	}
	apiKey, err := sanitizeOptionalField(raw.APIKey, maxInputRunes, "api_key")
	if err != nil {
		return GenerateResult{}, err
	}
	model, err := sanitizeOptionalField(raw.Model, maxInputRunes, "model")
	if err != nil {
		return GenerateResult{}, err
	}

	content, err := spec.Render(baseURL, apiKey, model)
	if err != nil {
		return GenerateResult{}, err
	}
	targetPath := spec.targetPath(homeDir)

	return GenerateResult{
		Tool:         spec.ID,
		Name:         spec.Name,
		Format:       spec.Format,
		Filename:     spec.filename(),
		TargetPath:   targetPath,
		Content:      content,
		Instructions: spec.Instruct(targetPath, baseURL),
	}, nil
}

func validateBaseURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || hasControlChar(trimmed) {
		return "", ErrInvalidBaseURL
	}
	if utf8.RuneCountInString(trimmed) > maxInputRunes {
		return "", ErrInvalidBaseURL
	}
	u, err := url.Parse(trimmed)
	if err != nil || u == nil {
		return "", ErrInvalidBaseURL
	}
	if !strings.EqualFold(u.Scheme, "http") && !strings.EqualFold(u.Scheme, "https") {
		return "", ErrInvalidBaseURL
	}
	if u.Host == "" || u.Hostname() == "" {
		return "", ErrInvalidBaseURL
	}
	return strings.TrimRight(trimmed, "/"), nil
}

func sanitizeOptionalField(raw string, maxRunes int, fieldName string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if hasControlChar(trimmed) {
		return "", fmt.Errorf("%w: %s must not contain control characters", ErrInvalidInput, fieldName)
	}
	if utf8.RuneCountInString(trimmed) > maxRunes {
		return "", fmt.Errorf("%w: %s must not exceed %d characters", ErrInvalidInput, fieldName, maxRunes)
	}
	return trimmed, nil
}

func hasControlChar(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}

func openAIBaseURL(baseURL string) string {
	trimmed := strings.TrimRight(baseURL, "/")
	if strings.HasSuffix(trimmed, "/v1") {
		return trimmed
	}
	return trimmed + "/v1"
}

func apiKeyOrPlaceholder(apiKey string) string {
	if apiKey != "" {
		return apiKey
	}
	return placeholderAPIKey
}

func marshalPrettyJSON(v any) (string, error) {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(raw) + "\n", nil
}

func renderClaudeCodeConfig(baseURL, apiKey, model string) (string, error) {
	env := map[string]string{
		"ANTHROPIC_BASE_URL":   baseURL,
		"ANTHROPIC_AUTH_TOKEN": apiKeyOrPlaceholder(apiKey),
	}
	if model != "" {
		env["ANTHROPIC_MODEL"] = model
	}
	return marshalPrettyJSON(map[string]any{
		"env": env,
	})
}

func renderCodexCLIConfig(baseURL, apiKey, model string) (string, error) {
	var b strings.Builder
	b.WriteString("model_provider = ")
	b.WriteString(strconv.Quote("omeg4mon-router"))
	b.WriteByte('\n')
	if model != "" {
		b.WriteString("model = ")
		b.WriteString(strconv.Quote(model))
		b.WriteByte('\n')
	}
	b.WriteString("\n[model_providers.omeg4mon-router]\n")
	b.WriteString("name = ")
	b.WriteString(strconv.Quote("omeg4mon-router"))
	b.WriteByte('\n')
	b.WriteString("base_url = ")
	b.WriteString(strconv.Quote(openAIBaseURL(baseURL)))
	b.WriteByte('\n')
	b.WriteString("env_key = ")
	b.WriteString(strconv.Quote("OMEG4MON_API_KEY"))
	b.WriteByte('\n')
	if apiKey != "" {
		b.WriteString("experimental_bearer_token = ")
		b.WriteString(strconv.Quote(apiKey))
		b.WriteByte('\n')
	}
	b.WriteString("wire_api = ")
	b.WriteString(strconv.Quote("responses"))
	b.WriteByte('\n')
	return b.String(), nil
}

func renderOpenCodeConfig(baseURL, apiKey, model string) (string, error) {
	models := map[string]any{}
	if model != "" {
		models[model] = map[string]any{
			"name": model,
		}
	}
	payload := map[string]any{
		"$schema": "https://opencode.ai/config.json",
		"provider": map[string]any{
			"omeg4mon-router": map[string]any{
				"npm":  "@ai-sdk/openai-compatible",
				"name": "omeg4mon-router",
				"options": map[string]any{
					"baseURL": openAIBaseURL(baseURL),
					"apiKey":  apiKeyOrPlaceholder(apiKey),
				},
				"models": models,
			},
		},
	}
	if model != "" {
		payload["model"] = "omeg4mon-router/" + model
	}
	return marshalPrettyJSON(payload)
}

func renderOpenAICompatProfileJSON(baseURL, apiKey, model string) (string, error) {
	payload := map[string]any{
		"apiProvider":   "openai",
		"openAiBaseUrl": openAIBaseURL(baseURL),
		"openAiApiKey":  apiKeyOrPlaceholder(apiKey),
	}
	if model != "" {
		payload["openAiModelId"] = model
	}
	return marshalPrettyJSON(payload)
}

func renderDroidConfig(baseURL, apiKey, model string) (string, error) {
	resolvedModel := model
	if resolvedModel == "" {
		resolvedModel = defaultModelHint
	}
	return marshalPrettyJSON(map[string]any{
		"custom_models": []map[string]any{
			{
				"model_display_name": "omeg4mon-router (" + resolvedModel + ")",
				"model":              resolvedModel,
				"base_url":           openAIBaseURL(baseURL),
				"api_key":            apiKeyOrPlaceholder(apiKey),
				"provider":           "generic-chat-completion-api",
			},
		},
	})
}

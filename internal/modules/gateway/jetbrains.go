package gateway

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const (
	providerIDJetBrains          = "jetbrains"
	jetbrainsDefaultModel        = "gpt-6-luna"
	jetbrainsMaxInputChars       = 5_200_000 // ~922K tokens hard ceiling
	jetbrainsMinCompletionTokens = 500
	jetbrainsDefaultMaxTokens    = 4096
	jetbrainsGrazieAgentHeader   = `{"name":"Pulse","version":"1"}`
	jetbrainsContextLimitMessage = "Input tokens exceed the configured limit of 922000 tokens. Please reduce the length of the messages."
)

// exceedsJetBrainsInputLimit reports whether the combined message text in body
// exceeds JetBrains' 922K-token (~5.2M character) ceiling.
func exceedsJetBrainsInputLimit(body []byte) bool {
	if len(body) <= jetbrainsMaxInputChars {
		return false
	}
	var payload struct {
		Messages []struct {
			Content any `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return false
	}
	totalChars := 0
	for _, msg := range payload.Messages {
		switch c := msg.Content.(type) {
		case string:
			totalChars += len(c)
		case []any:
			for _, part := range c {
				if pmap, ok := part.(map[string]any); ok {
					if ptype, _ := pmap["type"].(string); ptype == "text" {
						if txt, ok := pmap["text"].(string); ok {
							totalChars += len(txt)
						}
					}
				}
			}
		}
		if totalChars > jetbrainsMaxInputChars {
			return true
		}
	}
	return false
}

// transformJetBrainsRequest normalizes an OpenAI-chat payload for JetBrains AI
// (Grazie): enforces model gpt-6-luna, converts max_tokens to
// max_completion_tokens (minimum 500, default 4096), injects
// stream_options.include_usage on streaming requests, strips empty tools/functions
// arrays, and sets reasoning_effort="none" when tool calls are present.
func transformJetBrainsRequest(body []byte, stream bool) []byte {
	if len(body) == 0 {
		return body
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil || payload == nil {
		return body
	}

	payload["model"] = jetbrainsDefaultModel

	if rawMax, hasMax := payload["max_tokens"]; hasMax {
		payload["max_completion_tokens"] = rawMax
		delete(payload, "max_tokens")
	}

	mct := 0
	if rawVal, ok := payload["max_completion_tokens"].(float64); ok {
		mct = int(rawVal)
	}
	if mct < jetbrainsMinCompletionTokens {
		payload["max_completion_tokens"] = jetbrainsDefaultMaxTokens
	}

	isStream := stream
	if b, ok := payload["stream"].(bool); ok && b {
		isStream = true
	}
	if isStream {
		payload["stream"] = true
		streamOpts, ok := payload["stream_options"].(map[string]any)
		if !ok || streamOpts == nil {
			streamOpts = map[string]any{}
		}
		streamOpts["include_usage"] = true
		payload["stream_options"] = streamOpts
	}

	hasTools := false
	if tools, ok := payload["tools"].([]any); ok {
		if len(tools) == 0 {
			delete(payload, "tools")
		} else {
			hasTools = true
		}
	}
	if funcs, ok := payload["functions"].([]any); ok {
		if len(funcs) == 0 {
			delete(payload, "functions")
		} else {
			hasTools = true
		}
	}
	if hasTools {
		payload["reasoning_effort"] = "none"
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return body
	}
	return encoded
}

// applyJetBrainsHeaders configures the Grazie inference headers and authenticates
// via Grazie-Authenticate-JWT instead of an Authorization Bearer header.
func applyJetBrainsHeaders(req *http.Request, creds contracts.UpstreamCredentials, stream bool) {
	req.Header.Del("Authorization")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Grazie-Agent", jetbrainsGrazieAgentHeader)
	req.Header.Set("X-Forced-License-Free", "2")
	req.Header.Set("X-Accept-EAP-License", "false")
	req.Header.Set("X-Accept-Release-License", "false")

	jwt := strings.TrimSpace(creds.ExtraHeaders["Grazie-Authenticate-JWT"])
	if jwt == "" {
		jwt = strings.TrimSpace(creds.AccessToken)
	}
	if jwt == "" {
		jwt = strings.TrimSpace(creds.APIKey)
	}
	if jwt != "" {
		req.Header.Set("Grazie-Authenticate-JWT", jwt)
	}
	if stream {
		req.Header.Set("Accept", "text/event-stream")
	}
}

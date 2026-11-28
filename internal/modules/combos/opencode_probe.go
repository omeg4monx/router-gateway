package combos

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/opencode"
	"github.com/omeg4monx/router-gateway/internal/platform/venice"
)

func selectProbeHTTPClient(configured *http.Client, creds contracts.UpstreamCredentials) *http.Client {
	if venice.IsVeniceKeyless(creds) {
		return venice.SelectHTTPClient(configured)
	}
	if isOpenCodeFreeCreds(creds) {
		return opencode.ResolveClient(context.Background(), configured, creds.BaseURL)
	}
	if configured != nil {
		return configured
	}
	return &http.Client{Timeout: defaultProbeTimeout}
}

func isOpenCodeFreeCreds(creds contracts.UpstreamCredentials) bool {
	return opencode.IsFreeCreds(creds)
}

func marshalProbeJSON(payload any) ([]byte, error) {
	return opencode.MarshalNoEscapeHTML(payload)
}

func applyOpenCodeFreePayload(payload map[string]any, userPrompt string) {
	payload["stream"] = true
	payload["messages"] = []map[string]string{
		{"role": "system", "content": opencode.EnvFingerprintPrompt()},
		{"role": "user", "content": userPrompt},
	}
	payload["tools"] = opencode.CanonicalTools()
}

func applyOpenCodeFreeHeaders(req *http.Request) {
	sessionID := opencode.GenerateSessionID(time.Now())
	opencode.ApplyFreeHeaders(req, sessionID, nil, true, req.URL.String())
}

func extractSSETokenUsage(protocol string, body []byte) (inputTokens, outputTokens, cachedTokens int64) {
	for _, line := range bytes.Split(body, []byte("\n")) {
		trimmed := bytes.TrimSpace(line)
		if !bytes.HasPrefix(trimmed, []byte("data:")) {
			continue
		}
		data := bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data:")))
		if len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) || data[0] != '{' {
			continue
		}
		inTok, outTok, cachedTok := extractProbeTokenUsage(protocol, data)
		if inTok > 0 || outTok > 0 || cachedTok > 0 {
			inputTokens, outputTokens, cachedTokens = inTok, outTok, cachedTok
		}
	}
	return inputTokens, outputTokens, cachedTokens
}

func parseSSEIdentityPayload(body []byte) (reportedModel, idPrefix, fingerprint, selfReport string, hasReasoning bool) {
	var textBuf strings.Builder
	for _, line := range bytes.Split(body, []byte("\n")) {
		trimmed := bytes.TrimSpace(line)
		if !bytes.HasPrefix(trimmed, []byte("data:")) {
			continue
		}
		data := bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data:")))
		if len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) || data[0] != '{' {
			continue
		}
		var chunk map[string]any
		if err := json.Unmarshal(data, &chunk); err != nil {
			continue
		}
		if reportedModel == "" {
			if m, ok := chunk["model"].(string); ok {
				reportedModel = strings.TrimSpace(m)
			}
		}
		if fingerprint == "" {
			if fp, ok := chunk["system_fingerprint"].(string); ok {
				fingerprint = strings.TrimSpace(fp)
			}
		}
		if idPrefix == "" {
			if id, ok := chunk["id"].(string); ok {
				idPrefix = extractIDPrefix(strings.TrimSpace(id))
			}
		}
		if choices, ok := chunk["choices"].([]any); ok && len(choices) > 0 {
			if first, ok := choices[0].(map[string]any); ok {
				if delta, ok := first["delta"].(map[string]any); ok {
					if rc, ok := delta["reasoning_content"].(string); ok && strings.TrimSpace(rc) != "" {
						hasReasoning = true
					}
					if c, ok := delta["content"].(string); ok && c != "" {
						textBuf.WriteString(c)
					}
				}
			}
		}
	}
	return reportedModel, idPrefix, fingerprint, strings.TrimSpace(textBuf.String()), hasReasoning
}

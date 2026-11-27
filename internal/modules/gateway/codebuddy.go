package gateway

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

const (
	providerIDCodeBuddyCN   = "codebuddy-cn"
	providerIDCodeBuddyIntl = "codebuddy-intl"

	codebuddySystemPrompt = "You are CodeBuddy Code."
	codebuddyCLIVersion   = "2.108.1"
	codebuddyCLIUserAgent = "CLI/" + codebuddyCLIVersion + " CodeBuddy/" + codebuddyCLIVersion
)

var (
	codebuddyChannelMarkers = []string{
		"You are Claude Code, Anthropic's official CLI for Claude.",
		"x-anthropic-billing-header",
		"Main branch (you will usually use this for PRs):",
	}
	codebuddyChannelMarkerRE = regexp.MustCompile(
		`(?i)you are claude code|claude.?code.+official.+cli|anthropic.+official.+cli|anxthxropic.+official.+cli|you are (?:cursor|windsurf|cline|aider|continue|copilot|cody)|cc_entrypoint\s*=\s*(?:cli|vscode|jetbrains|gui)|claude.?code.+issues|give feedback.+claude.?code|OhMyOpenCode`,
	)
)

func isCodeBuddyProvider(providerID string) bool {
	return providerID == providerIDCodeBuddyCN || providerID == providerIDCodeBuddyIntl
}

func codebuddyDomainFor(providerID string) string {
	if providerID == providerIDCodeBuddyIntl {
		return "www.codebuddy.ai"
	}
	return "copilot.tencent.com"
}

func sanitizeCodeBuddyChannelMarkers(text string) string {
	if text == "" {
		return text
	}
	out := text
	for _, marker := range codebuddyChannelMarkers {
		if strings.Contains(out, marker) {
			out = strings.ReplaceAll(out, marker, "[external-client-identity]")
		}
	}
	return codebuddyChannelMarkerRE.ReplaceAllString(out, "[external-client-identity]")
}

func applyCodeBuddyHeaders(upReq, clientReq *http.Request, providerID string) {
	reqID := newRandomHexID(16)
	convID := ""
	convReqID := ""
	convMsgID := ""
	if clientReq != nil {
		convID = strings.TrimSpace(clientReq.Header.Get("X-Conversation-ID"))
		convReqID = strings.TrimSpace(clientReq.Header.Get("X-Conversation-Request-ID"))
		convMsgID = strings.TrimSpace(clientReq.Header.Get("X-Conversation-Message-ID"))
	}
	if convID == "" {
		convID = newRandomHexID(16)
	}
	if convReqID == "" {
		convReqID = newRandomHexID(16)
	}
	if convMsgID == "" {
		convMsgID = reqID
	}

	upReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	upReq.Header.Set("User-Agent", codebuddyCLIUserAgent)
	upReq.Header.Set("X-Product", "SaaS")
	upReq.Header.Set("X-App", "cli")
	upReq.Header.Set("X-IDE-Type", "CLI")
	upReq.Header.Set("X-IDE-Name", "CLI")
	upReq.Header.Set("X-IDE-Version", codebuddyCLIVersion)
	upReq.Header.Set("X-Private-Data", "false")
	upReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	upReq.Header.Set("X-CodeBuddy-Request", "1")
	upReq.Header.Set("X-Domain", codebuddyDomainFor(providerID))
	upReq.Header.Set("X-Request-ID", reqID)
	upReq.Header.Set("X-Conversation-ID", convID)
	upReq.Header.Set("X-Conversation-Request-ID", convReqID)
	upReq.Header.Set("X-Conversation-Message-ID", convMsgID)
	upReq.Header.Set("X-Agent-Intent", "craft")
	upReq.Header.Set("X-Agent-Purpose", "conversation")
}

// transformCodeBuddyRequest normalizes an OpenAI-chat payload for CodeBuddy CN/Intl:
//   - strips "none"/"off" reasoning_effort or sets reasoning_summary="auto"
//   - consolidates system/developer prompts with sanitized channel markers
//   - converts bare string user content to typed text blocks
func transformCodeBuddyRequest(body []byte) []byte {
	if len(body) == 0 {
		return body
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil || payload == nil {
		return body
	}

	if effort, ok := payload["reasoning_effort"].(string); ok {
		trimmed := strings.ToLower(strings.TrimSpace(effort))
		if trimmed == "none" || trimmed == "off" {
			delete(payload, "reasoning_effort")
		} else if trimmed != "" {
			payload["reasoning_summary"] = "auto"
		}
	}

	rawMsgs, ok := payload["messages"].([]any)
	if !ok {
		return body
	}

	systemParts := []string{codebuddySystemPrompt}
	otherMsgs := make([]any, 0, len(rawMsgs))

	for _, item := range rawMsgs {
		msg, ok := item.(map[string]any)
		if !ok {
			continue
		}
		role, _ := msg["role"].(string)
		switch role {
		case "system", "developer":
			rawText := extractCodeBuddyText(msg["content"])
			cleaned := strings.TrimSpace(sanitizeCodeBuddyChannelMarkers(rawText))
			if cleaned != "" && cleaned != codebuddySystemPrompt {
				systemParts = append(systemParts, cleaned)
			}
		case "user":
			copied := cloneShallowMap(msg)
			if text, isStr := copied["content"].(string); isStr {
				copied["content"] = []map[string]any{{"type": "text", "text": text}}
			}
			otherMsgs = append(otherMsgs, copied)
		case "assistant":
			copied := cloneShallowMap(msg)
			switch c := copied["content"].(type) {
			case string:
				copied["content"] = sanitizeCodeBuddyChannelMarkers(c)
			case []any:
				nextBlocks := make([]any, len(c))
				for i, b := range c {
					if block, isMap := b.(map[string]any); isMap {
						if txt, isStr := block["text"].(string); isStr {
							bc := cloneShallowMap(block)
							bc["text"] = sanitizeCodeBuddyChannelMarkers(txt)
							nextBlocks[i] = bc
							continue
						}
					}
					nextBlocks[i] = b
				}
				copied["content"] = nextBlocks
			}
			otherMsgs = append(otherMsgs, copied)
		default:
			otherMsgs = append(otherMsgs, msg)
		}
	}

	finalMsgs := make([]any, 0, 1+len(otherMsgs))
	finalMsgs = append(finalMsgs, map[string]any{
		"role":    "system",
		"content": strings.Join(systemParts, "\n\n"),
	})
	finalMsgs = append(finalMsgs, otherMsgs...)
	payload["messages"] = finalMsgs

	encoded, err := json.Marshal(payload)
	if err != nil {
		return body
	}
	return encoded
}

func extractCodeBuddyText(content any) string {
	switch c := content.(type) {
	case string:
		return c
	case []any:
		parts := make([]string, 0, len(c))
		for _, b := range c {
			switch bv := b.(type) {
			case string:
				parts = append(parts, bv)
			case map[string]any:
				if txt, ok := bv["text"].(string); ok {
					parts = append(parts, txt)
				}
			}
		}
		return strings.Join(parts, "\n")
	default:
		return ""
	}
}

func cloneShallowMap(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

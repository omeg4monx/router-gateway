package plugins

import "strings"

const systemDirectiveSep = "\n\n"

func hasSystemSegment(haystack, directive string) bool {
	if haystack == "" || directive == "" {
		return false
	}
	if haystack == directive {
		return true
	}
	if strings.HasPrefix(haystack, directive+systemDirectiveSep) ||
		strings.HasSuffix(haystack, systemDirectiveSep+directive) ||
		strings.Contains(haystack, systemDirectiveSep+directive+systemDirectiveSep) {
		return true
	}
	return false
}

func appendSystemSegment(existing, directive string) (string, bool) {
	trimmed := strings.TrimSpace(existing)
	if trimmed == "" {
		return directive, true
	}
	if hasSystemSegment(existing, directive) {
		return existing, false
	}
	return existing + systemDirectiveSep + directive, true
}

func isGeminiFormat(protocol string, payload map[string]any) bool {
	p := strings.ToLower(strings.TrimSpace(protocol))
	if strings.Contains(p, "gemini") || strings.Contains(p, "vertex") {
		return true
	}
	if _, hasContents := payload["contents"].([]any); hasContents {
		return true
	}
	if reqObj, ok := payload["request"].(map[string]any); ok {
		if _, hasContents := reqObj["contents"].([]any); hasContents {
			return true
		}
	}
	return false
}

// injectSystemDirectives appends directive into the appropriate system prompt slot
// for the target protocol/wire shape without breaking prompt caching or duplicating
// existing directives. Returns true if payload was modified.
func injectSystemDirectives(payload map[string]any, protocol, model, directive string) bool {
	if payload == nil || strings.TrimSpace(directive) == "" {
		return false
	}
	if isClaudeFormat(protocol, model, payload) {
		return injectClaudeSystem(payload, directive)
	}
	if isGeminiFormat(protocol, payload) {
		return injectGeminiSystem(payload, directive)
	}
	if _, hasInstr := payload["instructions"].(string); hasInstr {
		return injectInstructionsSystem(payload, directive)
	}
	if _, hasMsgs := payload["messages"].([]any); hasMsgs {
		return injectChatMessagesSystem(payload, directive)
	}
	if _, hasInput := payload["input"]; hasInput {
		return injectInstructionsSystem(payload, directive)
	}
	return false
}

func injectClaudeSystem(payload map[string]any, directive string) bool {
	switch sys := payload["system"].(type) {
	case string:
		next, changed := appendSystemSegment(sys, directive)
		if !changed {
			return false
		}
		payload["system"] = next
		return true
	case []any:
		lastCacheIdx := -1
		for i, item := range sys {
			block, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if txt, ok := block["text"].(string); ok && hasSystemSegment(txt, directive) {
				return false
			}
			if block["cache_control"] != nil {
				lastCacheIdx = i
			}
		}
		newBlock := map[string]any{"type": "text", "text": directive}
		if lastCacheIdx >= 0 {
			updated := make([]any, 0, len(sys)+1)
			updated = append(updated, sys[:lastCacheIdx]...)
			updated = append(updated, newBlock)
			updated = append(updated, sys[lastCacheIdx:]...)
			payload["system"] = updated
		} else {
			payload["system"] = append(sys, newBlock)
		}
		return true
	default:
		payload["system"] = directive
		return true
	}
}

func injectGeminiSystem(payload map[string]any, directive string) bool {
	target := payload
	if wrapped, ok := payload["request"].(map[string]any); ok {
		target = wrapped
	}
	key := "systemInstruction"
	if _, hasSnake := target["system_instruction"]; hasSnake {
		key = "system_instruction"
	}
	if sysObj, ok := target[key].(map[string]any); ok {
		if parts, ok := sysObj["parts"].([]any); ok {
			for _, rawPart := range parts {
				part, ok := rawPart.(map[string]any)
				if !ok {
					continue
				}
				if txt, ok := part["text"].(string); ok && hasSystemSegment(txt, directive) {
					return false
				}
			}
			sysObj["parts"] = append(parts, map[string]any{"text": directive})
			return true
		}
	}
	target[key] = map[string]any{
		"parts": []any{
			map[string]any{"text": directive},
		},
	}
	return true
}

func injectInstructionsSystem(payload map[string]any, directive string) bool {
	existing, _ := payload["instructions"].(string)
	next, changed := appendSystemSegment(existing, directive)
	if !changed {
		return false
	}
	payload["instructions"] = next
	return true
}

func injectChatMessagesSystem(payload map[string]any, directive string) bool {
	msgs, ok := payload["messages"].([]any)
	if !ok {
		return false
	}

	firstSysIdx := -1
	for i, item := range msgs {
		msg, ok := item.(map[string]any)
		if !ok {
			continue
		}
		role, _ := msg["role"].(string)
		if role != "system" && role != "developer" {
			continue
		}
		if firstSysIdx == -1 {
			firstSysIdx = i
		}
		switch content := msg["content"].(type) {
		case string:
			if hasSystemSegment(content, directive) {
				return false
			}
		case []any:
			for _, rawPart := range content {
				part, ok := rawPart.(map[string]any)
				if !ok {
					continue
				}
				if txt, ok := part["text"].(string); ok && hasSystemSegment(txt, directive) {
					return false
				}
			}
		}
	}

	if firstSysIdx >= 0 {
		msg, ok := msgs[firstSysIdx].(map[string]any)
		if !ok {
			return false
		}
		switch content := msg["content"].(type) {
		case string:
			next, changed := appendSystemSegment(content, directive)
			if !changed {
				return false
			}
			msg["content"] = next
			return true
		case []any:
			msg["content"] = append(content, map[string]any{
				"type": "text",
				"text": directive,
			})
			return true
		default:
			msg["content"] = directive
			return true
		}
	}

	updated := make([]any, 0, len(msgs)+1)
	updated = append(updated, map[string]any{
		"role":    "system",
		"content": directive,
	})
	updated = append(updated, msgs...)
	payload["messages"] = updated
	return true
}

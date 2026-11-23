package modeleditor

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
)

const (
	ActionRewriteModel     = "rewrite_model"
	ActionBypassOverrides  = "bypass_overrides"
	ActionBypassTokenSaver = "bypass_tokensaver"
	ActionBypassHeadroom   = "bypass_headroom"
)

// PreviewResult summarizes a dry-run execution of bypass + override rules.
type PreviewResult struct {
	Model              string         `json:"model"`
	Body               map[string]any `json:"body"`
	AppliedOverrides   []string       `json:"applied_overrides"`
	AppliedBypassRules []string       `json:"applied_bypass_rules"`
	Metadata           map[string]any `json:"metadata"`
}

func matchPattern(pattern, value string) bool {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" || pattern == "*" {
		return true
	}
	if pattern == value {
		return true
	}
	matched, err := path.Match(pattern, value)
	return err == nil && matched
}

func bypassRuleMatches(rule BypassRule, model, keyID, protocol string, headers http.Header) bool {
	if !rule.Enabled {
		return false
	}
	if !matchPattern(rule.ModelPattern, model) {
		return false
	}
	if rule.KeyID != "" && rule.KeyID != keyID {
		return false
	}
	if rule.Protocol != "" && rule.Protocol != protocol {
		return false
	}
	if rule.HeaderName != "" {
		val := ""
		if headers != nil {
			val = headers.Get(rule.HeaderName)
		}
		if rule.HeaderValue == "" {
			if val == "" {
				return false
			}
		} else if !strings.EqualFold(strings.TrimSpace(val), strings.TrimSpace(rule.HeaderValue)) {
			return false
		}
	}
	return true
}

func applyPipeline(
	protocol string,
	model string,
	keyID string,
	headers http.Header,
	body []byte,
	meta map[string]any,
	bypassRules []BypassRule,
	overrides []ModelOverride,
) (outModel string, outBody []byte, appliedBypass, appliedOverrides []string) {
	outModel = model
	outBody = body
	appliedBypass = make([]string, 0)
	appliedOverrides = make([]string, 0)

	if meta == nil {
		meta = make(map[string]any)
	}

	var root map[string]any
	hasJSON := len(body) > 0 && json.Unmarshal(body, &root) == nil && root != nil
	mutated := false
	skipOverrides := false

	for _, br := range bypassRules {
		if !bypassRuleMatches(br, outModel, keyID, protocol, headers) {
			continue
		}
		appliedBypass = append(appliedBypass, br.ID)
		switch br.Action {
		case ActionRewriteModel:
			if br.TargetModel != "" {
				outModel = br.TargetModel
				if hasJSON && protocol != "gemini" {
					root["model"] = outModel
					mutated = true
				}
			}
		case ActionBypassOverrides:
			skipOverrides = true
			meta["bypass_overrides"] = true
		case ActionBypassTokenSaver:
			meta["bypass_tokensaver"] = true
		case ActionBypassHeadroom:
			meta["bypass_headroom"] = true
		}
	}

	if len(appliedBypass) > 0 {
		meta["model_bypass_applied"] = strings.Join(appliedBypass, ",")
	}

	if !skipOverrides {
		for _, ov := range overrides {
			if !ov.Enabled || !matchPattern(ov.ModelPattern, outModel) {
				continue
			}
			appliedOverrides = append(appliedOverrides, ov.ID)
			if ov.TargetModel != "" {
				outModel = ov.TargetModel
				if hasJSON && protocol != "gemini" {
					root["model"] = outModel
					mutated = true
				}
			}
			if hasJSON {
				if applyOverrideToJSON(protocol, root, ov) {
					mutated = true
				}
			}
		}
	}

	if len(appliedOverrides) > 0 {
		meta["model_override_applied"] = strings.Join(appliedOverrides, ",")
	}

	if mutated {
		if encoded, err := json.Marshal(root); err == nil {
			outBody = encoded
		}
	}
	return outModel, outBody, appliedBypass, appliedOverrides
}

func applyOverrideToJSON(protocol string, root map[string]any, ov ModelOverride) bool {
	mutated := false

	for _, k := range ov.StripParams {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		if _, exists := root[k]; exists {
			delete(root, k)
			mutated = true
		}
		if protocol == "gemini" {
			if gc, ok := root["generationConfig"].(map[string]any); ok {
				if _, exists := gc[k]; exists {
					delete(gc, k)
					mutated = true
				}
			}
		}
	}

	if protocol == "gemini" {
		if ov.Temperature != nil || ov.MaxTokens != nil || ov.TopP != nil {
			gc, ok := root["generationConfig"].(map[string]any)
			if !ok || gc == nil {
				gc = make(map[string]any)
				root["generationConfig"] = gc
			}
			if ov.Temperature != nil {
				gc["temperature"] = *ov.Temperature
				mutated = true
			}
			if ov.MaxTokens != nil {
				gc["maxOutputTokens"] = *ov.MaxTokens
				mutated = true
			}
			if ov.TopP != nil {
				gc["topP"] = *ov.TopP
				mutated = true
			}
		}
	} else {
		if ov.Temperature != nil {
			root["temperature"] = *ov.Temperature
			mutated = true
		}
		if ov.MaxTokens != nil {
			if protocol == "openai-responses" {
				root["max_output_tokens"] = *ov.MaxTokens
			} else {
				root["max_tokens"] = *ov.MaxTokens
			}
			mutated = true
		}
		if ov.TopP != nil {
			root["top_p"] = *ov.TopP
			mutated = true
		}
	}

	for k, v := range ov.ExtraParams {
		root[k] = v
		mutated = true
	}

	if ov.SystemPromptPrefix != "" || ov.SystemPromptSuffix != "" {
		if wrapSystemPrompt(protocol, root, ov.SystemPromptPrefix, ov.SystemPromptSuffix) {
			mutated = true
		}
	}

	return mutated
}

func wrapText(base, prefix, suffix string) string {
	out := base
	if prefix != "" {
		if out == "" {
			out = prefix
		} else {
			out = prefix + "\n" + out
		}
	}
	if suffix != "" {
		if out == "" {
			out = suffix
		} else {
			out = out + "\n" + suffix
		}
	}
	return out
}

func wrapSystemPrompt(protocol string, root map[string]any, prefix, suffix string) bool {
	switch protocol {
	case "anthropic":
		switch cur := root["system"].(type) {
		case string:
			root["system"] = wrapText(cur, prefix, suffix)
			return true
		case []any:
			blocks := cur
			if prefix != "" {
				blocks = append([]any{map[string]any{"type": "text", "text": prefix}}, blocks...)
			}
			if suffix != "" {
				blocks = append(blocks, map[string]any{"type": "text", "text": suffix})
			}
			root["system"] = blocks
			return true
		default:
			root["system"] = wrapText("", prefix, suffix)
			return true
		}
	case "openai-responses":
		cur, _ := root["instructions"].(string)
		root["instructions"] = wrapText(cur, prefix, suffix)
		return true
	case "gemini":
		si, ok := root["systemInstruction"].(map[string]any)
		if !ok || si == nil {
			root["systemInstruction"] = map[string]any{
				"parts": []any{map[string]any{"text": wrapText("", prefix, suffix)}},
			}
			return true
		}
		parts, ok := si["parts"].([]any)
		if !ok || len(parts) == 0 {
			si["parts"] = []any{map[string]any{"text": wrapText("", prefix, suffix)}}
			return true
		}
		if first, ok := parts[0].(map[string]any); ok {
			cur, _ := first["text"].(string)
			first["text"] = wrapText(cur, prefix, suffix)
			return true
		}
		return false
	default: // openai-chat
		msgs, ok := root["messages"].([]any)
		if !ok {
			return false
		}
		for i, raw := range msgs {
			msg, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			if role, _ := msg["role"].(string); role == "system" {
				if content, ok := msg["content"].(string); ok {
					msg["content"] = wrapText(content, prefix, suffix)
					msgs[i] = msg
					return true
				}
			}
		}
		sysMsg := map[string]any{
			"role":    "system",
			"content": wrapText("", prefix, suffix),
		}
		root["messages"] = append([]any{sysMsg}, msgs...)
		return true
	}
}

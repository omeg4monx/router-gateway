package modeleditor

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
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

func validateModelPattern(pattern string) error {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" || len(pattern) > 128 {
		return errors.New("model_pattern must be 1..128 characters")
	}
	inBracket := false
	for i := range len(pattern) {
		switch pattern[i] {
		case '[':
			if inBracket {
				return errors.New("invalid model_pattern: nested '['")
			}
			inBracket = true
		case ']':
			if !inBracket {
				return errors.New("invalid model_pattern: unmatched ']'")
			}
			inBracket = false
		}
	}
	if inBracket {
		return errors.New("invalid model_pattern: unclosed '['")
	}
	if _, err := globToRegexp(pattern); err != nil {
		return errors.New("invalid model_pattern syntax")
	}
	return nil
}

func globToRegexp(pattern string) (*regexp.Regexp, error) {
	var b strings.Builder
	b.WriteByte('^')
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			b.WriteString(".*")
		case '?':
			b.WriteByte('.')
		case '[':
			end := strings.IndexByte(pattern[i+1:], ']')
			if end < 0 {
				return nil, errors.New("unclosed bracket")
			}
			b.WriteString(pattern[i : i+end+2])
			i += end + 1
		default:
			b.WriteString(regexp.QuoteMeta(pattern[i : i+1]))
		}
	}
	b.WriteByte('$')
	return regexp.Compile(b.String())
}

func compileGlobPattern(pattern string) *regexp.Regexp {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" || pattern == "*" || !strings.ContainsAny(pattern, "*?[") {
		return nil
	}
	re, err := globToRegexp(pattern)
	if err != nil {
		return nil
	}
	return re
}

func matchCompiledPattern(pattern string, compiled *regexp.Regexp, value string) bool {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" || pattern == "*" {
		return true
	}
	if pattern == value {
		return true
	}
	if compiled != nil {
		return compiled.MatchString(value)
	}
	if !strings.ContainsAny(pattern, "*?[") {
		return false
	}
	re, err := globToRegexp(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(value)
}

func bypassRuleMatches(rule BypassRule, model, keyID, protocol string, headers http.Header) bool {
	if !rule.Enabled {
		return false
	}
	if !matchCompiledPattern(rule.ModelPattern, rule.compiledPattern, model) {
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
			if !ov.Enabled || !matchCompiledPattern(ov.ModelPattern, ov.compiledPattern, outModel) {
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

	promptMode, _ := ov.ExtraParams[internalPromptModeKey].(string)

	for k, v := range ov.ExtraParams {
		if k == internalPromptModeKey {
			continue
		}
		if existingMap, ok := root[k].(map[string]any); ok {
			if patchMap, ok2 := v.(map[string]any); ok2 {
				for pk, pv := range patchMap {
					existingMap[pk] = cloneJSONValue(pv)
				}
				root[k] = existingMap
				mutated = true
				continue
			}
		}
		root[k] = cloneJSONValue(v)
		mutated = true
	}

	if ov.SystemPromptPrefix != "" || ov.SystemPromptSuffix != "" {
		if wrapSystemPrompt(protocol, root, ov.SystemPromptPrefix, ov.SystemPromptSuffix) {
			mutated = true
		}
		if strings.EqualFold(strings.TrimSpace(promptMode), "dual") {
			anchor := ov.SystemPromptSuffix
			if anchor == "" {
				anchor = ov.SystemPromptPrefix
			}
			if injectLatestUserAnchor(protocol, root, anchor) {
				mutated = true
			}
		}
	}

	return mutated
}

const internalPromptModeKey = "_prompt_mode"

func cloneJSONValue(v any) any {
	switch typed := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for k, val := range typed {
			out[k] = cloneJSONValue(val)
		}
		return out
	case []any:
		out := make([]any, len(typed))
		for i, val := range typed {
			out[i] = cloneJSONValue(val)
		}
		return out
	default:
		return v
	}
}

func hasPromptSegment(haystack, needle string) bool {
	if haystack == "" || needle == "" {
		return false
	}
	if haystack == needle {
		return true
	}
	for _, seg := range strings.Split(haystack, "\n\n") {
		if seg == needle {
			return true
		}
	}
	return false
}

func hasBlockText(blocks []any, needle string) bool {
	for _, raw := range blocks {
		if m, ok := raw.(map[string]any); ok {
			if txt, ok2 := m["text"].(string); ok2 && hasPromptSegment(txt, needle) {
				return true
			}
		}
	}
	return false
}

func wrapText(base, prefix, suffix string) string {
	out := base
	if prefix != "" && !hasPromptSegment(out, prefix) {
		if out == "" {
			out = prefix
		} else {
			out = prefix + "\n\n" + out
		}
	}
	if suffix != "" && !hasPromptSegment(out, suffix) {
		if out == "" {
			out = suffix
		} else {
			out = out + "\n\n" + suffix
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
			if prefix != "" && !hasBlockText(blocks, prefix) {
				blocks = append([]any{map[string]any{"type": "text", "text": prefix}}, blocks...)
			}
			if suffix != "" && !hasBlockText(blocks, suffix) {
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
		siKey := "systemInstruction"
		if _, hasSnake := root["system_instruction"]; hasSnake {
			siKey = "system_instruction"
		}
		si, ok := root[siKey].(map[string]any)
		if !ok || si == nil {
			root[siKey] = map[string]any{
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
			role, _ := msg["role"].(string)
			if role != "system" && role != "developer" {
				continue
			}
			switch content := msg["content"].(type) {
			case string:
				msg["content"] = wrapText(content, prefix, suffix)
				msgs[i] = msg
				return true
			case []any:
				blocks := content
				if prefix != "" && !hasBlockText(blocks, prefix) {
					blocks = append([]any{map[string]any{"type": "text", "text": prefix}}, blocks...)
				}
				if suffix != "" && !hasBlockText(blocks, suffix) {
					blocks = append(blocks, map[string]any{"type": "text", "text": suffix})
				}
				msg["content"] = blocks
				msgs[i] = msg
				return true
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

// injectLatestUserAnchor prepends anchor to the most recent user turn (idempotently)
// to counter multi-turn attention drift while preserving upstream system prefix caching.
func injectLatestUserAnchor(protocol string, root map[string]any, anchor string) bool {
	if anchor == "" {
		return false
	}
	switch protocol {
	case "gemini":
		contents, ok := root["contents"].([]any)
		if !ok {
			return false
		}
		for i := len(contents) - 1; i >= 0; i-- {
			item, ok := contents[i].(map[string]any)
			if !ok {
				continue
			}
			role, _ := item["role"].(string)
			if role != "" && role != "user" {
				continue
			}
			parts, ok := item["parts"].([]any)
			if !ok || len(parts) == 0 {
				continue
			}
			for j, rawPart := range parts {
				part, ok := rawPart.(map[string]any)
				if !ok {
					continue
				}
				if txt, ok := part["text"].(string); ok {
					next := wrapText(txt, anchor, "")
					if next == txt {
						return false
					}
					part["text"] = next
					parts[j] = part
					item["parts"] = parts
					contents[i] = item
					root["contents"] = contents
					return true
				}
			}
		}
		return false
	case "openai-responses":
		inputs, ok := root["input"].([]any)
		if !ok {
			return false
		}
		for i := len(inputs) - 1; i >= 0; i-- {
			item, ok := inputs[i].(map[string]any)
			if !ok {
				continue
			}
			if role, _ := item["role"].(string); role != "user" {
				continue
			}
			if injectMessageAnchor(item, anchor) {
				inputs[i] = item
				root["input"] = inputs
				return true
			}
			return false
		}
		return false
	default: // openai-chat & anthropic
		msgs, ok := root["messages"].([]any)
		if !ok {
			return false
		}
		for i := len(msgs) - 1; i >= 0; i-- {
			msg, ok := msgs[i].(map[string]any)
			if !ok {
				continue
			}
			if role, _ := msg["role"].(string); role != "user" {
				continue
			}
			if injectMessageAnchor(msg, anchor) {
				msgs[i] = msg
				root["messages"] = msgs
				return true
			}
			return false
		}
		return false
	}
}

func injectMessageAnchor(msg map[string]any, anchor string) bool {
	switch content := msg["content"].(type) {
	case string:
		next := wrapText(content, anchor, "")
		if next == content {
			return false
		}
		msg["content"] = next
		return true
	case []any:
		if hasBlockText(content, anchor) {
			return false
		}
		msg["content"] = append([]any{map[string]any{"type": "text", "text": anchor}}, content...)
		return true
	default:
		return false
	}
}

package plugins

import (
	"encoding/json"
	"regexp"
	"strings"
)

var (
	jsonWantPromptRe = regexp.MustCompile(`(?i)\b(?:respond|reply|answer|return|output|produce)\b[^.!?]{0,40}\bjson\b|\bjson\s+only\b|\bonly\s+(?:output|return)\s+json\b`)
	fencedJSONRe     = regexp.MustCompile("(?i)```(?:json|jsonc|json5)?[ \\t]*\\r?\\n([\\s\\S]*?)```")
	trailingCommaRe  = regexp.MustCompile(`,\s*([}\]])`)
	pyTrueRe         = regexp.MustCompile(`\bTrue\b`)
	pyFalseRe        = regexp.MustCompile(`\bFalse\b`)
	pyNoneRe         = regexp.MustCompile(`\bNone\b`)
	pyNaNRe          = regexp.MustCompile(`\bNaN\b`)
	pyNegInfRe       = regexp.MustCompile(`-Infinity\b`)
	pyPosInfRe       = regexp.MustCompile(`\bInfinity\b`)
	singleQuoteStrRe = regexp.MustCompile(`'([^'\\]*(?:\\.[^'\\]*)*)'`)
)

func promptWantsJSON(payload map[string]any) bool {
	var b strings.Builder
	if sys, ok := payload["system"].(string); ok {
		b.WriteString(sys)
		b.WriteByte(' ')
	}
	if msgs, ok := payload["messages"].([]any); ok {
		count := 0
		for _, item := range msgs {
			m, ok := item.(map[string]any)
			if !ok {
				continue
			}
			role, _ := m["role"].(string)
			if role != "system" && role != "user" {
				continue
			}
			if content, ok := m["content"].(string); ok {
				b.WriteString(content)
				b.WriteByte(' ')
			}
			count++
			if count >= 4 {
				break
			}
		}
	}
	haystack := b.String()
	if len(haystack) > 4000 {
		haystack = haystack[:4000]
	}
	return jsonWantPromptRe.MatchString(haystack)
}

func tryParseJSON(text string) (any, bool) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return nil, false
	}
	var out any
	if err := json.Unmarshal([]byte(trimmed), &out); err != nil {
		return nil, false
	}
	return out, true
}

func scanBalancedJSON(text string) (string, bool) {
	start := strings.IndexAny(text, "{[")
	if start == -1 {
		return "", false
	}
	stack := make([]byte, 0, 8)
	inString := false
	escaped := false

	for i := start; i < len(text); i++ {
		ch := text[i]
		if inString {
			switch {
			case escaped:
				escaped = false
			case ch == '\\':
				escaped = true
			case ch == '"':
				inString = false
			}
			continue
		}
		if ch == '"' {
			inString = true
			continue
		}
		if ch == '{' {
			stack = append(stack, '}')
			continue
		}
		if ch == '[' {
			stack = append(stack, ']')
			continue
		}
		if ch == '}' || ch == ']' {
			if len(stack) == 0 || stack[len(stack)-1] != ch {
				return "", false
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				return text[start : i+1], true
			}
		}
	}
	return "", false
}

func extractFencedJSON(text string) (string, bool) {
	m := fencedJSONRe.FindStringSubmatch(text)
	if len(m) < 2 {
		return "", false
	}
	return strings.TrimSpace(m[1]), true
}

func mapOutsideStrings(text string, fn func(string) string) string {
	var out strings.Builder
	var seg strings.Builder
	inString := false
	escaped := false

	flush := func() {
		if seg.Len() > 0 {
			out.WriteString(fn(seg.String()))
			seg.Reset()
		}
	}

	for i := 0; i < len(text); i++ {
		ch := text[i]
		if inString {
			out.WriteByte(ch)
			switch {
			case escaped:
				escaped = false
			case ch == '\\':
				escaped = true
			case ch == '"':
				inString = false
			}
			continue
		}
		if ch == '"' {
			flush()
			out.WriteByte(ch)
			inString = true
			continue
		}
		seg.WriteByte(ch)
	}
	flush()
	return out.String()
}

func repairJSONText(text string) string {
	out := mapOutsideStrings(text, func(seg string) string {
		s := pyTrueRe.ReplaceAllString(seg, "true")
		s = pyFalseRe.ReplaceAllString(s, "false")
		s = pyNoneRe.ReplaceAllString(s, "null")
		s = pyNaNRe.ReplaceAllString(s, "null")
		s = pyNegInfRe.ReplaceAllString(s, "null")
		s = pyPosInfRe.ReplaceAllString(s, "null")
		return s
	})

	if !strings.Contains(out, `"`) && strings.Contains(out, `'`) {
		out = singleQuoteStrRe.ReplaceAllString(out, `"$1"`)
	}
	return trailingCommaRe.ReplaceAllString(out, "$1")
}

func closeTruncatedJSON(text string) string {
	if strings.TrimSpace(text) == "" {
		return text
	}
	openerPos := make([]int, 0, 8)
	inString := false
	escaped := false
	lastComma := -1
	lastOpener := -1

	for i := 0; i < len(text); i++ {
		ch := text[i]
		if inString {
			switch {
			case escaped:
				escaped = false
			case ch == '\\':
				escaped = true
			case ch == '"':
				inString = false
			}
			continue
		}
		if ch == '"' {
			inString = true
			continue
		}
		if ch == '{' || ch == '[' {
			openerPos = append(openerPos, i)
			lastOpener = i
			lastComma = -1
			continue
		}
		if ch == '}' || ch == ']' {
			if len(openerPos) == 0 {
				return text
			}
			openerPos = openerPos[:len(openerPos)-1]
			lastComma = -1
			continue
		}
		if ch == ',' {
			lastComma = i
		}
	}

	if len(openerPos) == 0 {
		return text
	}

	var closers strings.Builder
	for i := len(openerPos) - 1; i >= 0; i-- {
		if text[openerPos[i]] == '{' {
			closers.WriteByte('}')
		} else {
			closers.WriteByte(']')
		}
	}
	closerStr := closers.String()

	candidates := make([]string, 0, 3)
	if !inString {
		candidates = append(candidates, text+closerStr)
	}
	if lastComma >= 0 {
		candidates = append(candidates, text[:lastComma]+closerStr)
	}
	if lastOpener >= 0 {
		candidates = append(candidates, text[:lastOpener+1]+closerStr)
	}

	for _, candidate := range candidates {
		if _, ok := tryParseJSON(candidate); ok {
			return candidate
		}
	}
	if len(candidates) > 0 {
		return candidates[0]
	}
	return text
}

type coercedJSON struct {
	Value     any
	Repaired  bool
	Truncated bool
}

func coerceJSONOutput(text string) (coercedJSON, bool) {
	if val, ok := tryParseJSON(text); ok {
		return coercedJSON{Value: val, Repaired: false}, true
	}
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return coercedJSON{}, false
	}

	if fence, ok := extractFencedJSON(trimmed); ok {
		if val, parsed := tryParseJSON(fence); parsed {
			return coercedJSON{Value: val, Repaired: true}, true
		}
		if val, parsed := tryParseJSON(repairJSONText(fence)); parsed {
			return coercedJSON{Value: val, Repaired: true}, true
		}
	}

	if balanced, ok := scanBalancedJSON(trimmed); ok {
		if val, parsed := tryParseJSON(balanced); parsed {
			return coercedJSON{Value: val, Repaired: true}, true
		}
		if val, parsed := tryParseJSON(repairJSONText(balanced)); parsed {
			return coercedJSON{Value: val, Repaired: true}, true
		}
	}

	repaired := repairJSONText(trimmed)
	if val, ok := tryParseJSON(repaired); ok {
		return coercedJSON{Value: val, Repaired: true}, true
	}

	closed := closeTruncatedJSON(repaired)
	if val, ok := tryParseJSON(closed); ok {
		return coercedJSON{Value: val, Repaired: true, Truncated: true}, true
	}

	return coercedJSON{}, false
}

func wantsJSONAnswer(text string) bool {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return false
	}
	if fencedJSONRe.MatchString(trimmed) {
		return true
	}
	opener := strings.IndexAny(trimmed, "{[")
	if opener == -1 {
		return false
	}
	_, balanced := scanBalancedJSON(trimmed[opener:])
	return balanced
}

func indexToolSchemas(tools []any) map[string]map[string]any {
	if len(tools) == 0 {
		return nil
	}
	out := make(map[string]map[string]any, len(tools))
	for _, raw := range tools {
		toolMap, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		fn := toolMap
		if inner, ok := toolMap["function"].(map[string]any); ok {
			fn = inner
		}
		name, _ := fn["name"].(string)
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if params, ok := fn["parameters"].(map[string]any); ok {
			out[name] = params
		} else if params, ok := fn["input_schema"].(map[string]any); ok {
			out[name] = params
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func filterToolArguments(args, schema map[string]any) (filtered map[string]any, changed bool) {
	if args == nil || schema == nil {
		return args, false
	}
	props, _ := schema["properties"].(map[string]any)
	result := make(map[string]any, len(args))
	for k, v := range args {
		if props != nil {
			if _, allowed := props[k]; !allowed {
				changed = true
				continue
			}
		}
		result[k] = v
	}

	if reqList, ok := schema["required"].([]any); ok {
		for _, item := range reqList {
			key, ok := item.(string)
			if !ok || key == "" {
				continue
			}
			if _, exists := result[key]; exists {
				continue
			}
			prop, _ := props[key].(map[string]any)
			if prop == nil {
				continue
			}
			if defVal, hasDef := prop["default"]; hasDef {
				result[key] = defVal
				changed = true
				continue
			}
			propType, _ := prop["type"].(string)
			switch propType {
			case "boolean":
				result[key] = false
				changed = true
			case "number", "integer":
				result[key] = 0
				changed = true
			case "array":
				result[key] = []any{}
				changed = true
			}
		}
	}
	return result, changed
}

func applyJSONGuardResponse(payload map[string]any, tools []any) bool {
	if payload == nil {
		return false
	}
	schemas := indexToolSchemas(tools)
	changed := false

	if choices, ok := payload["choices"].([]any); ok {
		for _, item := range choices {
			choice, ok := item.(map[string]any)
			if !ok {
				continue
			}
			msg, ok := choice["message"].(map[string]any)
			if !ok {
				continue
			}
			if content, ok := msg["content"].(string); ok && wantsJSONAnswer(content) {
				if coerced, ok := coerceJSONOutput(content); ok && coerced.Repaired {
					if encoded, err := json.Marshal(coerced.Value); err == nil {
						msg["content"] = string(encoded)
						changed = true
					}
				}
			}
			if toolCalls, ok := msg["tool_calls"].([]any); ok {
				for _, tcItem := range toolCalls {
					tc, ok := tcItem.(map[string]any)
					if !ok {
						continue
					}
					fn, ok := tc["function"].(map[string]any)
					if !ok {
						continue
					}
					argsStr, ok := fn["arguments"].(string)
					if !ok {
						continue
					}
					coerced, parsed := coerceJSONOutput(argsStr)
					if !parsed {
						continue
					}
					argsMap, ok := coerced.Value.(map[string]any)
					if !ok {
						continue
					}
					dirty := coerced.Repaired
					if schemas != nil {
						fnName, _ := fn["name"].(string)
						if cleaned, schemaChanged := filterToolArguments(argsMap, schemas[fnName]); schemaChanged {
							argsMap = cleaned
							dirty = true
						}
					}
					if dirty {
						if encoded, err := json.Marshal(argsMap); err == nil {
							fn["arguments"] = string(encoded)
							changed = true
						}
					}
				}
			}
		}
	}

	// Anthropic Messages response content blocks.
	if blocks, ok := payload["content"].([]any); ok {
		for _, item := range blocks {
			block, ok := item.(map[string]any)
			if !ok {
				continue
			}
			blockType, _ := block["type"].(string)
			switch blockType {
			case "tool_use":
				if inputStr, ok := block["input"].(string); ok {
					if coerced, parsed := coerceJSONOutput(inputStr); parsed {
						if argsMap, isObj := coerced.Value.(map[string]any); isObj {
							block["input"] = argsMap
							changed = true
						}
					}
				}
			case "text":
				if text, ok := block["text"].(string); ok && wantsJSONAnswer(text) {
					if coerced, parsed := coerceJSONOutput(text); parsed && coerced.Repaired {
						if encoded, err := json.Marshal(coerced.Value); err == nil {
							block["text"] = string(encoded)
							changed = true
						}
					}
				}
			}
		}
	}

	return changed
}

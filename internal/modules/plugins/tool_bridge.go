package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
)

var (
	toolCallTagRe   = regexp.MustCompile(`(?is)<tool_call>\s*(\{[\s\S]*?\})\s*</tool_call>`)
	xmlToolUseRe    = regexp.MustCompile(`(?is)<tool_use\s+name\s*=\s*["']([^"']+)["']\s*>([\s\S]*?)</tool_use>`)
	openCodeToolRe  = regexp.MustCompile(`(?m)^\s*(?:assistant\s+)?to\s*=\s*([A-Za-z_][\w.-]*)\s*(\{[\s\S]*?\})\s*$`)
	fencedToolRe    = regexp.MustCompile("(?is)```(?:json|tool_call|tool)?\\s*(\\{[^{}]*?\"(?:tool_)?(?:name|call)\"[\\s\\S]*?\\})\\s*```")
	mistralToolRe   = regexp.MustCompile(`(?s)\[TOOL_CALLS\]\s*(\[[\s\S]*?\])`)
	invokeXMLRe     = regexp.MustCompile(`(?is)<(?:antml:)?invoke\s+name\s*=\s*["']([^"']+)["']\s*>([\s\S]*?)</(?:antml:)?invoke>`)
	invokeParamRe   = regexp.MustCompile(`(?is)<parameter\s+name\s*=\s*["']([^"']+)["']\s*>([\s\S]*?)</parameter>`)
	functionTagRe   = regexp.MustCompile(`(?is)<function\s*=\s*["']?([A-Za-z_][\w.-]*)["']?\s*>([\s\S]*?)</function>`)
	bareToolObjRe   = regexp.MustCompile(`(?m)^\s*(\{\s*"function"\s*:\s*"[^"]+"[\s\S]*?\})\s*$`)
	toolIDCounter   atomic.Uint64
	toolNormSplitRe = regexp.MustCompile(`[-_]+`)
)

type extractedToolCall struct {
	Name      string
	Arguments map[string]any
}

func nextToolCallID(seed string) string {
	n := toolIDCounter.Add(1)
	if seed == "" {
		seed = "otb"
	}
	return fmt.Sprintf("call_%s_%d", seed, n)
}

func parseLenientJSONObject(text string) map[string]any {
	if val, ok := tryParseJSON(text); ok {
		if m, isObj := val.(map[string]any); isObj {
			return m
		}
	}
	if val, ok := tryParseJSON(repairJSONText(text)); ok {
		if m, isObj := val.(map[string]any); isObj {
			return m
		}
	}
	return nil
}

func parseLenientJSONArray(text string) []any {
	if val, ok := tryParseJSON(text); ok {
		if arr, isArr := val.([]any); isArr {
			return arr
		}
	}
	if val, ok := tryParseJSON(repairJSONText(text)); ok {
		if arr, isArr := val.([]any); isArr {
			return arr
		}
	}
	return nil
}

func toolNameFromMap(m map[string]any) string {
	if m == nil {
		return ""
	}
	for _, key := range []string{"name", "tool", "tool_name", "function"} {
		if s, ok := m[key].(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
		if sub, ok := m[key].(map[string]any); ok {
			if s, ok := sub["name"].(string); ok && strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s)
			}
		}
	}
	if fc, ok := m["function_call"].(map[string]any); ok {
		if s, ok := fc["name"].(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func normalizeArgsValue(raw any) map[string]any {
	if raw == nil {
		return map[string]any{}
	}
	if m, ok := raw.(map[string]any); ok {
		return m
	}
	if s, ok := raw.(string); ok {
		if parsed := parseLenientJSONObject(s); parsed != nil {
			return parsed
		}
	}
	return map[string]any{}
}

func toolArgsFromMap(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	for _, key := range []string{"arguments", "parameters", "input", "args"} {
		if v, exists := m[key]; exists && v != nil {
			return normalizeArgsValue(v)
		}
	}
	if fn, ok := m["function"].(map[string]any); ok {
		if v, exists := fn["arguments"]; exists && v != nil {
			return normalizeArgsValue(v)
		}
	}
	reserved := map[string]struct{}{
		"name": {}, "tool": {}, "tool_name": {}, "type": {}, "id": {}, "call": {}, "function": {},
	}
	rest := make(map[string]any)
	for k, v := range m {
		if _, skip := reserved[k]; !skip {
			rest[k] = v
		}
	}
	return rest
}

func extractToolCallsFromProse(text string) (calls []extractedToolCall, cleaned string) {
	if strings.TrimSpace(text) == "" {
		return nil, text
	}
	cleaned = text

	harvest := func(fullMatch, nameHint, rawArgs string) {
		payload := parseLenientJSONObject(rawArgs)
		name := toolNameFromMap(payload)
		if name == "" {
			name = strings.TrimSpace(nameHint)
		}
		if name == "" {
			return
		}
		calls = append(calls, extractedToolCall{
			Name:      name,
			Arguments: toolArgsFromMap(payload),
		})
		cleaned = strings.Replace(cleaned, fullMatch, "", 1)
	}

	for _, m := range toolCallTagRe.FindAllStringSubmatch(text, -1) {
		harvest(m[0], "", m[1])
	}
	if len(calls) == 0 {
		for _, m := range xmlToolUseRe.FindAllStringSubmatch(text, -1) {
			harvest(m[0], m[1], m[2])
		}
	}
	if len(calls) == 0 {
		for _, m := range openCodeToolRe.FindAllStringSubmatch(text, -1) {
			harvest(m[0], m[1], m[2])
		}
	}
	if len(calls) == 0 {
		for _, m := range fencedToolRe.FindAllStringSubmatch(text, -1) {
			harvest(m[0], "", m[1])
		}
	}
	if len(calls) == 0 {
		for _, m := range mistralToolRe.FindAllStringSubmatch(text, -1) {
			items := parseLenientJSONArray(m[1])
			took := 0
			for _, item := range items {
				obj, ok := item.(map[string]any)
				if !ok {
					continue
				}
				name := toolNameFromMap(obj)
				if name == "" {
					continue
				}
				calls = append(calls, extractedToolCall{
					Name:      name,
					Arguments: toolArgsFromMap(obj),
				})
				took++
			}
			if took > 0 {
				cleaned = strings.Replace(cleaned, m[0], "", 1)
			}
		}
	}
	if len(calls) == 0 {
		for _, m := range invokeXMLRe.FindAllStringSubmatch(text, -1) {
			name := strings.TrimSpace(m[1])
			if name == "" {
				continue
			}
			args := make(map[string]any)
			for _, p := range invokeParamRe.FindAllStringSubmatch(m[2], -1) {
				args[strings.TrimSpace(p[1])] = strings.TrimSpace(p[2])
			}
			if len(args) > 0 {
				calls = append(calls, extractedToolCall{Name: name, Arguments: args})
				cleaned = strings.Replace(cleaned, m[0], "", 1)
			}
		}
	}
	if len(calls) == 0 {
		for _, m := range functionTagRe.FindAllStringSubmatch(text, -1) {
			harvest(m[0], m[1], m[2])
		}
	}
	if len(calls) == 0 {
		for _, m := range bareToolObjRe.FindAllStringSubmatch(text, -1) {
			harvest(m[0], "", m[1])
		}
	}

	return calls, strings.TrimSpace(cleaned)
}

func declaredArgNames(schema map[string]any, depth int, out map[string]struct{}) {
	if schema == nil || depth > 6 {
		return
	}
	if props, ok := schema["properties"].(map[string]any); ok {
		for k := range props {
			out[k] = struct{}{}
		}
	}
	for _, key := range []string{"items", "anyOf", "oneOf", "allOf"} {
		switch child := schema[key].(type) {
		case map[string]any:
			declaredArgNames(child, depth+1, out)
		case []any:
			for _, item := range child {
				if m, ok := item.(map[string]any); ok {
					declaredArgNames(m, depth+1, out)
				}
			}
		}
	}
}

func normalizeToolKey(name string) string {
	return toolNormSplitRe.ReplaceAllString(strings.ToLower(strings.TrimSpace(name)), "_")
}

func applyOpenAIToolBridgeResponse(payload map[string]any, tools []any) bool {
	if payload == nil || len(tools) == 0 {
		return false
	}
	catalog := indexToolSchemas(tools)
	if len(catalog) == 0 {
		return false
	}

	resolveName := func(raw string) string {
		if _, ok := catalog[raw]; ok {
			return raw
		}
		target := normalizeToolKey(raw)
		for declared := range catalog {
			if normalizeToolKey(declared) == target {
				return declared
			}
		}
		return ""
	}

	choices, ok := payload["choices"].([]any)
	if !ok {
		return false
	}

	changed := false
	for _, item := range choices {
		choice, ok := item.(map[string]any)
		if !ok {
			continue
		}
		msg, ok := choice["message"].(map[string]any)
		if !ok {
			continue
		}
		if existing, ok := msg["tool_calls"].([]any); ok && len(existing) > 0 {
			continue
		}
		content, ok := msg["content"].(string)
		if !ok || strings.TrimSpace(content) == "" {
			continue
		}

		extracted, cleaned := extractToolCallsFromProse(content)
		if len(extracted) == 0 {
			continue
		}

		toolCalls := make([]any, 0, len(extracted))
		for _, call := range extracted {
			resolved := resolveName(call.Name)
			if resolved == "" {
				continue
			}
			args := call.Arguments
			declared := make(map[string]struct{})
			declaredArgNames(catalog[resolved], 0, declared)
			if len(declared) > 0 {
				filtered := make(map[string]any, len(args))
				for k, v := range args {
					if _, allowed := declared[k]; allowed {
						filtered[k] = v
					}
				}
				if len(filtered) > 0 {
					args = filtered
				}
			}
			argsJSON, err := json.Marshal(args)
			if err != nil {
				argsJSON = []byte("{}")
			}
			toolCalls = append(toolCalls, map[string]any{
				"id":   nextToolCallID(resolved),
				"type": "function",
				"function": map[string]any{
					"name":      resolved,
					"arguments": string(argsJSON),
				},
			})
		}

		if len(toolCalls) == 0 {
			continue
		}
		msg["content"] = cleaned
		msg["tool_calls"] = toolCalls
		choice["finish_reason"] = "tool_calls"
		changed = true
	}

	return changed
}

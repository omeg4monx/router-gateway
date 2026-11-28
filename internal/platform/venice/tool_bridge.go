package venice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
)

// ToolCall represents an OpenAI-compatible function tool call extracted from
// Venice model prose when the caller declared tools in the request.
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction carries the invoked tool's canonical name and JSON arguments string.
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

var (
	toolTagRE     = regexp.MustCompile(`(?is)<tool_call>\s*(\{[\s\S]*?\})\s*</tool_call>`)
	xmlToolUseRE  = regexp.MustCompile(`(?is)<tool_use\s+name\s*=\s*["']([^"']+)["']\s*>([\s\S]*?)</tool_use>`)
	fencedToolRE  = regexp.MustCompile("(?is)```(?:json|tool_call|tool)?\\s*(\\{[\\s\\S]*?\\})\\s*```")
	invokeXMLRE   = regexp.MustCompile(`(?is)<invoke\s+name\s*=\s*["']([^"']+)["']\s*>([\s\S]*?)</invoke>`)
	invokeParamRE = regexp.MustCompile(`(?is)<parameter\s+name\s*=\s*["']([^"']+)["']\s*>([\s\S]*?)</parameter>`)
	toolCallSeq   atomic.Uint64
)

// ExtractRequestTools returns the declared tool definitions (from "tools" or
// legacy "functions") carried by an OpenAI Chat Completions request body.
func ExtractRequestTools(openAIBody []byte) []json.RawMessage {
	trimmed := bytes.TrimSpace(openAIBody)
	if len(trimmed) == 0 {
		return nil
	}
	var req struct {
		Tools     []json.RawMessage `json:"tools"`
		Functions []json.RawMessage `json:"functions"`
	}
	if json.Unmarshal(trimmed, &req) != nil {
		return nil
	}
	if len(req.Tools) > 0 {
		return req.Tools
	}
	if len(req.Functions) > 0 {
		return req.Functions
	}
	return nil
}

type declaredTool struct {
	name          string
	declaredProps map[string]struct{}
}

func buildToolCatalog(tools []json.RawMessage) map[string]declaredTool {
	catalog := make(map[string]declaredTool, len(tools))
	for _, rawTool := range tools {
		var wrapper struct {
			Function *struct {
				Name       string          `json:"name"`
				Parameters json.RawMessage `json:"parameters"`
			} `json:"function"`
			Name        string          `json:"name"`
			Parameters  json.RawMessage `json:"parameters"`
			InputSchema json.RawMessage `json:"input_schema"`
		}
		if json.Unmarshal(rawTool, &wrapper) != nil {
			continue
		}
		name := strings.TrimSpace(wrapper.Name)
		schemaRaw := wrapper.Parameters
		if len(bytes.TrimSpace(schemaRaw)) == 0 {
			schemaRaw = wrapper.InputSchema
		}
		if wrapper.Function != nil {
			name = strings.TrimSpace(wrapper.Function.Name)
			schemaRaw = wrapper.Function.Parameters
		}
		if name == "" {
			continue
		}
		catalog[name] = declaredTool{
			name:          name,
			declaredProps: collectSchemaProperties(schemaRaw),
		}
	}
	return catalog
}

func collectSchemaProperties(schemaRaw json.RawMessage) map[string]struct{} {
	trimmed := bytes.TrimSpace(schemaRaw)
	if len(trimmed) == 0 {
		return nil
	}
	var root map[string]any
	if json.Unmarshal(trimmed, &root) != nil {
		return nil
	}
	props := make(map[string]struct{})
	var visit func(node any, depth int)
	visit = func(node any, depth int) {
		if depth > 6 {
			return
		}
		m, ok := node.(map[string]any)
		if !ok {
			return
		}
		if p, ok := m["properties"].(map[string]any); ok {
			for k := range p {
				props[k] = struct{}{}
			}
		}
		for _, key := range []string{"items", "anyOf", "oneOf", "allOf"} {
			child, exists := m[key]
			if !exists {
				continue
			}
			if arr, ok := child.([]any); ok {
				for _, item := range arr {
					visit(item, depth+1)
				}
			} else {
				visit(child, depth+1)
			}
		}
	}
	visit(root, 0)
	if len(props) == 0 {
		return nil
	}
	return props
}

func normalizeToolName(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, "-", "_")
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}
	return s
}

func resolveDeclaredTool(catalog map[string]declaredTool, rawName string) (declaredTool, bool) {
	trimmed := strings.TrimSpace(rawName)
	if trimmed == "" {
		return declaredTool{}, false
	}
	if dt, ok := catalog[trimmed]; ok {
		return dt, true
	}
	normTarget := normalizeToolName(trimmed)
	for _, dt := range catalog {
		if normalizeToolName(dt.name) == normTarget {
			return dt, true
		}
	}
	return declaredTool{}, false
}

type parsedCandidate struct {
	name string
	args map[string]any
}

func parseToolObject(rawJSON string) map[string]any {
	trimmed := strings.TrimSpace(rawJSON)
	if trimmed == "" {
		return nil
	}
	var obj map[string]any
	if json.Unmarshal([]byte(trimmed), &obj) == nil && obj != nil {
		return obj
	}
	return nil
}

func extractNameFromObject(obj map[string]any) string {
	if obj == nil {
		return ""
	}
	for _, key := range []string{"name", "tool", "tool_name", "function"} {
		if v, ok := obj[key]; ok {
			switch typed := v.(type) {
			case string:
				if s := strings.TrimSpace(typed); s != "" {
					return s
				}
			case map[string]any:
				if s, ok := typed["name"].(string); ok && strings.TrimSpace(s) != "" {
					return strings.TrimSpace(s)
				}
			}
		}
	}
	if fc, ok := obj["function_call"].(map[string]any); ok {
		if s, ok := fc["name"].(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func extractArgsFromObject(obj map[string]any) map[string]any {
	if obj == nil {
		return map[string]any{}
	}
	for _, key := range []string{"arguments", "parameters", "input", "args"} {
		if v, ok := obj[key]; ok && v != nil {
			if m := normalizeArgsMap(v); m != nil {
				return m
			}
		}
	}
	if fn, ok := obj["function"].(map[string]any); ok {
		if v, ok := fn["arguments"]; ok && v != nil {
			if m := normalizeArgsMap(v); m != nil {
				return m
			}
		}
	}
	reserved := map[string]struct{}{
		"name": {}, "tool": {}, "tool_name": {}, "type": {}, "id": {}, "call": {}, "function": {},
	}
	rest := make(map[string]any)
	for k, v := range obj {
		if _, isRes := reserved[k]; !isRes {
			rest[k] = v
		}
	}
	return rest
}

func normalizeArgsMap(v any) map[string]any {
	switch typed := v.(type) {
	case map[string]any:
		return typed
	case string:
		if parsed := parseToolObject(typed); parsed != nil {
			return parsed
		}
	}
	return nil
}

func extractCandidatesFromText(text string) (candidates []parsedCandidate, cleaned string) {
	cleaned = text

	harvest := func(fullMatch, nameHint, argsRaw string) {
		obj := parseToolObject(argsRaw)
		name := extractNameFromObject(obj)
		if name == "" {
			name = strings.TrimSpace(nameHint)
		}
		if name == "" {
			return
		}
		candidates = append(candidates, parsedCandidate{
			name: name,
			args: extractArgsFromObject(obj),
		})
		cleaned = strings.Replace(cleaned, fullMatch, "", 1)
	}

	for _, m := range toolTagRE.FindAllStringSubmatch(text, -1) {
		harvest(m[0], "", m[1])
	}
	if len(candidates) == 0 {
		for _, m := range xmlToolUseRE.FindAllStringSubmatch(text, -1) {
			harvest(m[0], m[1], m[2])
		}
	}
	if len(candidates) == 0 {
		for _, m := range fencedToolRE.FindAllStringSubmatch(text, -1) {
			obj := parseToolObject(m[1])
			if extractNameFromObject(obj) != "" {
				harvest(m[0], "", m[1])
			}
		}
	}
	if len(candidates) == 0 {
		for _, m := range invokeXMLRE.FindAllStringSubmatch(text, -1) {
			name := strings.TrimSpace(m[1])
			if name == "" {
				continue
			}
			args := make(map[string]any)
			for _, pm := range invokeParamRE.FindAllStringSubmatch(m[2], -1) {
				args[strings.TrimSpace(pm[1])] = strings.TrimSpace(pm[2])
			}
			candidates = append(candidates, parsedCandidate{name: name, args: args})
			cleaned = strings.Replace(cleaned, m[0], "", 1)
		}
	}

	return candidates, strings.TrimSpace(cleaned)
}

// ApplyToolBridge extracts structured OpenAI tool_calls from Venice model prose
// when the caller declared tools. Undeclared tool names are dropped so model
// hallucinations never reach the client.
func ApplyToolBridge(content string, tools []json.RawMessage) (cleaned string, calls []ToolCall, applied bool) {
	if len(tools) == 0 || strings.TrimSpace(content) == "" {
		return content, nil, false
	}
	catalog := buildToolCatalog(tools)
	if len(catalog) == 0 {
		return content, nil, false
	}

	candidates, stripped := extractCandidatesFromText(content)
	if len(candidates) == 0 {
		return content, nil, false
	}

	calls = make([]ToolCall, 0, len(candidates))
	for _, cand := range candidates {
		dt, ok := resolveDeclaredTool(catalog, cand.name)
		if !ok {
			continue
		}
		args := cand.args
		if args == nil {
			args = map[string]any{}
		}
		if len(dt.declaredProps) > 0 && len(args) > 0 {
			filtered := make(map[string]any, len(args))
			for k, v := range args {
				if _, allowed := dt.declaredProps[k]; allowed {
					filtered[k] = v
				}
			}
			if len(filtered) > 0 {
				args = filtered
			}
		}
		argsBytes, err := json.Marshal(args)
		if err != nil {
			argsBytes = []byte("{}")
		}
		seq := toolCallSeq.Add(1)
		calls = append(calls, ToolCall{
			ID:   fmt.Sprintf("call_vnc_%s", strconv.FormatUint(seq, 36)),
			Type: "function",
			Function: ToolCallFunction{
				Name:      dt.name,
				Arguments: string(argsBytes),
			},
		})
	}
	if len(calls) == 0 {
		return content, nil, false
	}
	return stripped, calls, true
}

func mayContainToolMarker(s string) bool {
	lower := strings.ToLower(s)
	return strings.Contains(lower, "<tool_call") ||
		strings.Contains(lower, "<tool_use") ||
		strings.Contains(lower, "<invoke") ||
		strings.Contains(lower, "```")
}

// couldBeToolMarkerPrefix reports whether buf (starting at '<' or '`') is either
// a complete/in-progress tool marker or a prefix that could still become one.
// Returns false as soon as buf is proven to be ordinary prose/code (e.g. "a < b"
// or "```python"), allowing streaming to flush immediately.
func couldBeToolMarkerPrefix(buf string) bool {
	if buf == "" {
		return false
	}
	lower := strings.ToLower(buf)
	prefixes := []string{
		"<tool_call",
		"<tool_use",
		"<invoke",
		"```json",
		"```tool_call",
		"```tool",
		"```\n",
		"```\r",
		"``` ",
		"```{",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(lower, p) || strings.HasPrefix(p, lower) {
			return true
		}
	}
	return false
}

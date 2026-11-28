package opencode

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
)

const (
	// MaxToolSurfaceChars enforces OpenCode free-tier Cloudflare WAF's 38K char tool payload limit.
	MaxToolSurfaceChars = 38000
	maxToolNameLen      = 128
	maxSSELineBytes     = 2 << 20
)

var (
	canonicalToolOrder = []string{"edit", "glob", "grep", "read", "shell"}
	gatedToolNameRe    = regexp.MustCompile(`(?i)terminal|reportfindings|exec_command`)
)

type gatedNameSchemaRule struct {
	re    *regexp.Regexp
	props []string
	alias string
}

var gatedNameSchemaRules = []gatedNameSchemaRule{
	{re: regexp.MustCompile(`(?i)^edit$`), props: []string{"old_string", "new_string"}, alias: "file_edit"},
	{re: regexp.MustCompile(`^Read$`), props: []string{"file_path"}, alias: "file_read"},
}

type canonicalToolDef struct {
	description string
	parameters  map[string]any
	required    []string
	propNames   []string
}

func canonicalToolDefinitions() map[string]canonicalToolDef {
	return map[string]canonicalToolDef{
		"edit": {
			description: "Edit the contents of a file by finding and replacing exact text. When editing text from Read output, preserve the exact indentation (tabs or spaces) and omit the line-number prefix, such as `1: `. Never include the prefix in oldString or newString. Write line breaks as \\n; the file's existing CRLF or LF line endings are matched and preserved. To change line endings themselves, make oldString and newString differ only in \\r; that edit applies to the exact bytes. If there is no exact match, the edit retries while ignoring typographic quote, dash, and space variants, then trailing whitespace on each line. The edit fails if oldString is not found. By default, oldString must identify a UNIQUE location. Multiple matches FAIL unless replaceAll is true. Add more surrounding context to disambiguate, or set replaceAll to true to replace every occurrence. Use replaceAll when the change should apply to every occurrence, such as renaming a variable.",
			parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"path":       map[string]any{"type": "string", "description": "File to edit"},
					"oldString":  map[string]any{"type": "string", "description": "Exact text to find and replace"},
					"newString":  map[string]any{"type": "string", "description": "Text to replace oldString with (must differ from oldString)"},
					"replaceAll": map[string]any{"type": "boolean", "description": "Whether to replace every occurrence of oldString. When false, oldString must match exactly once. Defaults to false."},
				},
				"required":             []string{"path", "oldString", "newString"},
				"additionalProperties": false,
			},
			required:  []string{"path", "oldString", "newString"},
			propNames: []string{"path", "oldString", "newString", "replaceAll"},
		},
		"glob": {
			description: "Find files matching a glob pattern relative to a directory. Output is a list of file paths. Always use forward slashes (/) as path separators even on Windows. The pattern supports standard glob features: * matches characters within a directory, ** matches across directories, ? matches a single character, and [abc] matches a character set. Prefer specific patterns over broad ones.",
			parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"pattern": map[string]any{"type": "string", "description": "Glob pattern to match files against"},
					"path":    map[string]any{"type": "string", "description": "Directory to search. Defaults to the current working directory."},
					"hidden":  map[string]any{"type": "boolean", "description": "Include hidden files and directories (default: false)."},
					"limit":   map[string]any{"type": "integer", "exclusiveMinimum": 0, "description": "Maximum number of matching files to return (default: 100)"},
				},
				"required":             []string{"pattern"},
				"additionalProperties": false,
			},
			required:  []string{"pattern"},
			propNames: []string{"pattern", "path", "hidden", "limit"},
		},
		"grep": {
			description: "Search for a regular expression pattern within files. Lines are returned with 1-based line numbers. By default, searches recursively from the working directory. Output is sorted by file path and line number.",
			parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"pattern":       map[string]any{"type": "string", "minLength": 1, "description": "Regular expression or literal text to match in file contents."},
					"path":          map[string]any{"type": "string", "description": "File or directory to search. Defaults to the current working directory."},
					"include":       map[string]any{"type": "string", "description": "Glob pattern to filter files (for example, \"*.js\" or \"*.{ts,tsx}\")"},
					"literal":       map[string]any{"type": "boolean", "description": "Treat pattern as exact text instead of a regular expression (default: false)."},
					"caseSensitive": map[string]any{"type": "boolean", "description": "Use case-sensitive matching (default: true)."},
					"limit":         map[string]any{"type": "integer", "exclusiveMinimum": 0, "description": "Maximum number of matching lines to return (default: 100)"},
				},
				"required":             []string{"pattern"},
				"additionalProperties": false,
			},
			required:  []string{"pattern"},
			propNames: []string{"pattern", "path", "include", "literal", "caseSensitive", "limit"},
		},
		"read": {
			description: "Read the contents of a file or directory. Supports text files, images, and PDFs. Images and PDFs are presented directly to the model. Each text line is prefixed by its 1-based line number as <line>: <content>. The prefix is for reference and is not part of the file content. Directory entries are returned one per line. Use offset and limit to read large files or directories in sections. Prefer one larger read over many small slices, and use grep to find specific content in large files.",
			parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"path":   map[string]any{"type": "string", "description": "File or directory to read"},
					"offset": map[string]any{"type": "integer", "minimum": 0, "description": "The line or directory entry to start reading from (1-based)"},
					"limit":  map[string]any{"type": "integer", "minimum": 0, "description": "The maximum number of lines or directory entries to read (defaults to and capped at 2000)"},
				},
				"required":             []string{"path"},
				"additionalProperties": false,
			},
			required:  []string{"path"},
			propNames: []string{"path", "offset", "limit"},
		},
		"shell": {
			description: "Execute a shell command and return its output. Commands run on Linux using bash. Quote file paths containing spaces or special characters. Prefer dedicated tools over shell commands when possible. When output is large, the full result is saved to a file and a truncated preview is returned. Rely on automatic truncation unless filtering the output is more useful. Commands accept an optional timeout, background commands have no timeout by default. Background commands return immediately, and you will be notified when they complete.",
			parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"command":    map[string]any{"type": "string", "description": "Shell command string to execute"},
					"workdir":    map[string]any{"type": "string", "description": "Working directory to execute the command in. Defaults to the current working directory. When possible, avoid changing directories in the command and set the working directory here instead."},
					"timeout":    map[string]any{"type": "integer", "minimum": 0, "description": "Timeout in milliseconds. Set to 0 to disable the timeout. Defaults to 120000 for foreground commands. Background commands have no timeout by default."},
					"background": map[string]any{"type": "boolean", "description": "Run the command in the background and return immediately (useful for dev servers and long-running builds). You do not need to use '&' at the end of the command when using this parameter. You will be notified when it completes. DO NOT poll for completion."},
				},
				"required":             []string{"command"},
				"additionalProperties": false,
			},
			required:  []string{"command"},
			propNames: []string{"command", "workdir", "timeout", "background"},
		},
	}
}

// CanonicalTools returns the 5 verbatim OpenCode v2 tool definitions in OpenAI Chat shape.
func CanonicalTools() []map[string]any {
	defs := canonicalToolDefinitions()
	out := make([]map[string]any, 0, len(canonicalToolOrder))
	for _, name := range canonicalToolOrder {
		d := defs[name]
		out = append(out, map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        name,
				"description": d.description,
				"parameters":  d.parameters,
			},
		})
	}
	return out
}

// EnvFingerprintPrompt builds the CLI <env> block required by OpenCode's free-tier gate.
func EnvFingerprintPrompt() string {
	cwd, err := os.Getwd()
	if err != nil || strings.TrimSpace(cwd) == "" {
		cwd = "/tmp"
	}
	tmpDir := filepath.Join(os.TempDir(), "opencode")
	return strings.Join([]string{
		"Here is some useful information about the environment you are running in:",
		"<env>",
		"  Working directory: " + cwd,
		"  Workspace root folder: " + cwd,
		"  Is directory a git repo: no",
		"  Platform: " + runtime.GOOS,
		"  Prefer " + tmpDir + " over generic system temporary directories such as /tmp; it is pre-created and approved for external access.",
		"</env>",
	}, "\n")
}

// MarshalNoEscapeHTML encodes v as JSON without escaping `<` or `>` (preserving `<env>`).
func MarshalNoEscapeHTML(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimSpace(buf.Bytes()), nil
}

// TransformRequest applies OpenCode request normalization, Responses sanitization,
// non-destructive <env> injection, and tool cloaking. Returned renameMap maps
// `sentName -> originalCallerName` for response-side restoration.
func TransformRequest(providerID, model string, rawBody []byte, isResponses, isClaude bool) (outBody []byte, renameMap map[string]string, err error) {
	var body map[string]any
	if unmarshalErr := json.Unmarshal(rawBody, &body); unmarshalErr != nil {
		return rawBody, nil, unmarshalErr
	}

	if model != "" {
		if existing, _ := body["model"].(string); strings.TrimSpace(existing) == "" {
			body["model"] = model
		}
	}

	isFree := IsFreeProviderID(providerID)
	if isFree {
		// Upstream rejects non-streaming requests on free models with 403 FreeTierError.
		body["stream"] = true
	}

	if isResponses {
		normalizeResponsesPayload(body)
	} else {
		dropAdaptiveReasoningEffort(body)
	}

	if isFree {
		renameMap = applyFingerprintTools(body, isResponses, isClaude)
	}

	outBody, err = MarshalNoEscapeHTML(body)
	if err != nil {
		return rawBody, nil, err
	}
	return outBody, renameMap, nil
}

func normalizeResponsesPayload(body map[string]any) {
	if _, hasMaxOut := body["max_output_tokens"]; !hasMaxOut {
		if v, ok := body["max_completion_tokens"]; ok && v != nil {
			body["max_output_tokens"] = v
		} else if v, ok := body["max_tokens"]; ok && v != nil {
			body["max_output_tokens"] = v
		}
	}
	delete(body, "max_tokens")
	delete(body, "max_completion_tokens")
	body["store"] = false

	// Normalize reasoning_effort -> reasoning:{effort, summary:"auto"}.
	if effortRaw, ok := body["reasoning_effort"].(string); ok {
		effort := strings.ToLower(strings.TrimSpace(effortRaw))
		delete(body, "reasoning_effort")
		if effort != "" && effort != "auto" {
			if effort == "max" || effort == "ultra" {
				effort = "xhigh"
			}
			reasoning, _ := body["reasoning"].(map[string]any)
			if reasoning == nil {
				reasoning = make(map[string]any)
			}
			reasoning["effort"] = effort
			if _, hasSummary := reasoning["summary"]; !hasSummary {
				reasoning["summary"] = "auto"
			}
			body["reasoning"] = reasoning
		}
	}

	// Sanitize input items: strip prior-turn reasoning items & encrypted_content
	// which fail validation across pooled upstream accounts.
	if inputs, ok := body["input"].([]any); ok {
		filtered := make([]any, 0, len(inputs))
		for _, raw := range inputs {
			item, ok := raw.(map[string]any)
			if !ok {
				filtered = append(filtered, raw)
				continue
			}
			itemType, _ := item["type"].(string)
			if itemType == "reasoning" {
				continue
			}
			delete(item, "encrypted_content")
			delete(item, "reasoning_encrypted_content")
			if itemType == "function_call" {
				name, _ := item["name"].(string)
				name = strings.TrimSpace(name)
				if name == "" {
					continue
				}
				if len(name) > maxToolNameLen {
					name = name[:maxToolNameLen]
				}
				item["name"] = name
			}
			filtered = append(filtered, item)
		}
		body["input"] = filtered
	}
}

func dropAdaptiveReasoningEffort(body map[string]any) {
	effort, ok := body["reasoning_effort"].(string)
	if !ok {
		return
	}
	switch strings.ToLower(strings.TrimSpace(effort)) {
	case "auto":
		delete(body, "reasoning_effort")
	case "xhigh", "max":
		body["reasoning_effort"] = "high"
	}
}

func applyFingerprintTools(body map[string]any, flat, isClaude bool) map[string]string {
	defs := canonicalToolDefinitions()
	rawTools, _ := body["tools"].([]any)
	cloaked, renameMap := concealToolNames(rawTools, defs)
	cloaked = appendMissingCanonicalTools(cloaked, defs, flat)
	cloaked = enforceToolSurfaceCap(cloaked, defs)
	body["tools"] = cloaked

	if len(renameMap) > 0 {
		retargetToolChoice(body, renameMap)
		aliasHistoryToolCalls(body, renameMap)
	}

	applyEnvFingerprint(body, flat, isClaude)

	// OpenCode rejects tool_choice: "none" with 403 FreeTierError.
	if tcStr, ok := body["tool_choice"].(string); ok && strings.EqualFold(strings.TrimSpace(tcStr), "none") {
		delete(body, "tool_choice")
	} else if tcMap, ok := body["tool_choice"].(map[string]any); ok {
		if t, _ := tcMap["type"].(string); strings.EqualFold(strings.TrimSpace(t), "none") {
			delete(body, "tool_choice")
		}
	}
	if flat {
		if _, hasChoice := body["tool_choice"]; !hasChoice {
			body["tool_choice"] = "auto"
		}
	}

	if _, hasLegacyCache := body["prompt_cache_retention"]; hasLegacyCache {
		delete(body, "prompt_cache_retention")
		if _, hasOpts := body["prompt_cache_options"]; !hasOpts {
			body["prompt_cache_options"] = map[string]any{"mode": "implicit"}
		}
	}

	return renameMap
}

func toolNameOf(tool map[string]any) string {
	if name, ok := tool["name"].(string); ok && strings.TrimSpace(name) != "" {
		return strings.TrimSpace(name)
	}
	if fn, ok := tool["function"].(map[string]any); ok {
		if name, ok := fn["name"].(string); ok {
			return strings.TrimSpace(name)
		}
	}
	return ""
}

func declaredToolProps(tool map[string]any) []string {
	var params map[string]any
	if fn, ok := tool["function"].(map[string]any); ok {
		params, _ = fn["parameters"].(map[string]any)
	} else {
		params, _ = tool["parameters"].(map[string]any)
	}
	props, _ := params["properties"].(map[string]any)
	if len(props) == 0 {
		return nil
	}
	out := make([]string, 0, len(props))
	for k := range props {
		out = append(out, k)
	}
	return out
}

func withToolName(tool map[string]any, newName string) map[string]any {
	cloned := make(map[string]any, len(tool))
	for k, v := range tool {
		cloned[k] = v
	}
	if fn, ok := tool["function"].(map[string]any); ok {
		fnClone := make(map[string]any, len(fn))
		for k, v := range fn {
			fnClone[k] = v
		}
		fnClone["name"] = newName
		cloned["function"] = fnClone
	} else {
		cloned["name"] = newName
	}
	return cloned
}

func canNormalizeInPlace(tool map[string]any, def canonicalToolDef) bool {
	declared := declaredToolProps(tool)
	if len(declared) == 0 {
		return false
	}
	for _, prop := range declared {
		if !slices.Contains(def.propNames, prop) {
			return false
		}
	}
	for _, req := range def.required {
		if !slices.Contains(declared, req) {
			return false
		}
	}
	return true
}

func withCanonicalDefinition(tool map[string]any, def canonicalToolDef) map[string]any {
	cloned := make(map[string]any, len(tool))
	for k, v := range tool {
		cloned[k] = v
	}
	if fn, ok := tool["function"].(map[string]any); ok {
		fnClone := make(map[string]any, len(fn))
		for k, v := range fn {
			fnClone[k] = v
		}
		fnClone["description"] = def.description
		fnClone["parameters"] = def.parameters
		cloned["function"] = fnClone
	} else {
		cloned["description"] = def.description
		cloned["parameters"] = def.parameters
	}
	return cloned
}

func nextAvailableAlias(base string, taken map[string]bool) string {
	alias := base
	for i := 2; taken[alias]; i++ {
		alias = fmt.Sprintf("%s%d", base, i)
	}
	taken[alias] = true
	return alias
}

func concealToolNames(rawTools []any, defs map[string]canonicalToolDef) (cloaked []any, renameMap map[string]string) {
	renameMap = make(map[string]string)
	if len(rawTools) == 0 {
		return nil, renameMap
	}

	taken := make(map[string]bool, len(rawTools)+len(canonicalToolOrder))
	for _, raw := range rawTools {
		if m, ok := raw.(map[string]any); ok {
			if n := toolNameOf(m); n != "" {
				taken[n] = true
			}
		}
	}
	for _, name := range canonicalToolOrder {
		taken[name] = true
	}

	seenCanonical := make(map[string]bool)
	out := make([]any, 0, len(rawTools))
	for _, raw := range rawTools {
		tool, ok := raw.(map[string]any)
		if !ok {
			out = append(out, raw)
			continue
		}
		name := toolNameOf(tool)
		if name == "" {
			continue
		}

		// Gate 1: Gated tool names (terminal, ReportFindings, exec_command).
		if gatedToolNameRe.MatchString(name) {
			alias := nextAvailableAlias("run_command", taken)
			renameMap[alias] = name
			out = append(out, withToolName(tool, alias))
			continue
		}

		// Gate 2: Gated name+schema pairs (edit+old_string/new_string, Read+file_path).
		props := declaredToolProps(tool)
		matchedSchemaGate := false
		for _, rule := range gatedNameSchemaRules {
			if !rule.re.MatchString(name) {
				continue
			}
			for _, p := range rule.props {
				if !slices.Contains(props, p) {
					continue
				}
				alias := nextAvailableAlias(rule.alias, taken)
				renameMap[alias] = name
				out = append(out, withToolName(tool, alias))
				matchedSchemaGate = true
				break
			}
			if matchedSchemaGate {
				break
			}
		}
		if matchedSchemaGate {
			continue
		}

		// Gate 3: Collision with canonical fingerprint tool names.
		def, isCanonical := defs[name]
		if !isCanonical {
			out = append(out, tool)
			continue
		}
		if seenCanonical[name] {
			continue
		}
		seenCanonical[name] = true

		if canNormalizeInPlace(tool, def) {
			out = append(out, withCanonicalDefinition(tool, def))
			continue
		}
		alias := nextAvailableAlias(name+"_client", taken)
		renameMap[alias] = name
		out = append(out, withToolName(tool, alias))
	}
	return out, renameMap
}

func appendMissingCanonicalTools(tools []any, defs map[string]canonicalToolDef, flat bool) []any {
	present := make(map[string]bool, len(tools))
	for _, raw := range tools {
		if m, ok := raw.(map[string]any); ok {
			if n := toolNameOf(m); n != "" {
				present[n] = true
			}
		}
	}
	for _, name := range canonicalToolOrder {
		if present[name] {
			continue
		}
		d := defs[name]
		if flat {
			tools = append(tools, map[string]any{
				"type":        "function",
				"name":        name,
				"description": d.description,
				"parameters":  d.parameters,
				"strict":      false,
			})
		} else {
			tools = append(tools, map[string]any{
				"type": "function",
				"function": map[string]any{
					"name":        name,
					"description": d.description,
					"parameters":  d.parameters,
				},
			})
		}
	}
	return tools
}

func enforceToolSurfaceCap(tools []any, defs map[string]canonicalToolDef) []any {
	encoded, err := json.Marshal(tools)
	if err != nil || len(encoded) <= MaxToolSurfaceChars {
		return tools
	}
	var canonical, extras []any
	for _, raw := range tools {
		m, ok := raw.(map[string]any)
		if ok {
			if _, isCan := defs[toolNameOf(m)]; isCan {
				canonical = append(canonical, raw)
				continue
			}
		}
		extras = append(extras, raw)
	}
	canBytes, _ := json.Marshal(canonical)
	budget := MaxToolSurfaceChars - len(canBytes)
	kept := make([]any, 0, len(canonical)+len(extras))
	kept = append(kept, canonical...)
	for _, ex := range extras {
		b, err := json.Marshal(ex)
		if err != nil {
			continue
		}
		if budget-(len(b)+1) > 0 {
			kept = append(kept, ex)
			budget -= len(b) + 1
		}
	}
	return kept
}

func invertRenameMap(renameMap map[string]string) map[string]string {
	inv := make(map[string]string, len(renameMap))
	for alias, orig := range renameMap {
		inv[orig] = alias
	}
	return inv
}

func retargetToolChoice(body map[string]any, renameMap map[string]string) {
	choice, ok := body["tool_choice"].(map[string]any)
	if !ok {
		return
	}
	byOrig := invertRenameMap(renameMap)
	if name, ok := choice["name"].(string); ok {
		if alias, found := byOrig[strings.TrimSpace(name)]; found {
			choice["name"] = alias
		}
		return
	}
	if fn, ok := choice["function"].(map[string]any); ok {
		if name, ok := fn["name"].(string); ok {
			if alias, found := byOrig[strings.TrimSpace(name)]; found {
				fn["name"] = alias
			}
		}
	}
}

func aliasHistoryToolCalls(body map[string]any, renameMap map[string]string) {
	byOrig := invertRenameMap(renameMap)
	if inputs, ok := body["input"].([]any); ok {
		for _, raw := range inputs {
			item, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			if t, _ := item["type"].(string); t == "function_call" {
				if name, _ := item["name"].(string); name != "" {
					if alias, found := byOrig[name]; found {
						item["name"] = alias
					}
				}
			}
		}
	}
	if msgs, ok := body["messages"].([]any); ok {
		for _, raw := range msgs {
			msg, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			calls, ok := msg["tool_calls"].([]any)
			if !ok {
				continue
			}
			for _, rawCall := range calls {
				call, ok := rawCall.(map[string]any)
				if !ok {
					continue
				}
				if fn, ok := call["function"].(map[string]any); ok {
					if name, _ := fn["name"].(string); name != "" {
						if alias, found := byOrig[name]; found {
							fn["name"] = alias
						}
					}
				}
			}
		}
	}
}

func containsEnvTag(v any) bool {
	switch s := v.(type) {
	case string:
		return strings.Contains(s, "<env>")
	case []any:
		for _, block := range s {
			if m, ok := block.(map[string]any); ok {
				if t, _ := m["text"].(string); strings.Contains(t, "<env>") {
					return true
				}
			}
		}
	}
	return false
}

func appendSystemText(existing any, envBlock string) any {
	switch s := existing.(type) {
	case string:
		if strings.TrimSpace(s) == "" {
			return envBlock
		}
		return s + "\n\n" + envBlock
	case []any:
		out := slices.Clone(s)
		for i := len(out) - 1; i >= 0; i-- {
			if m, ok := out[i].(map[string]any); ok {
				if t, ok := m["text"].(string); ok {
					cloned := make(map[string]any, len(m))
					for k, v := range m {
						cloned[k] = v
					}
					if strings.TrimSpace(t) == "" {
						cloned["text"] = envBlock
					} else {
						cloned["text"] = t + "\n\n" + envBlock
					}
					out[i] = cloned
					return out
				}
			}
		}
		return append(out, map[string]any{"type": "text", "text": envBlock})
	default:
		return envBlock
	}
}

func applyEnvFingerprint(body map[string]any, flat, isClaude bool) {
	envBlock := EnvFingerprintPrompt()
	if flat {
		if inst, ok := body["instructions"]; ok && inst != nil {
			if _, isStr := inst.(string); !isStr {
				return
			}
			if containsEnvTag(inst) {
				return
			}
			body["instructions"] = appendSystemText(inst, envBlock)
			return
		}
		body["instructions"] = envBlock
		return
	}

	if isClaude || body["system"] != nil {
		if containsEnvTag(body["system"]) {
			return
		}
		body["system"] = appendSystemText(body["system"], envBlock)
		return
	}

	msgs, ok := body["messages"].([]any)
	if !ok {
		return
	}
	for i, raw := range msgs {
		msg, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		if role, _ := msg["role"].(string); role == "system" {
			if containsEnvTag(msg["content"]) {
				return
			}
			msg["content"] = appendSystemText(msg["content"], envBlock)
			msgs[i] = msg
			body["messages"] = msgs
			return
		}
	}
	sysMsg := map[string]any{"role": "system", "content": envBlock}
	body["messages"] = append([]any{sysMsg}, msgs...)
}

// RestoreToolNamesJSON restores any aliased tool names in a JSON response or SSE
// event payload (supports OpenAI Chat, OpenAI Responses, and Claude shapes).
func RestoreToolNamesJSON(raw []byte, renameMap map[string]string) []byte {
	if len(renameMap) == 0 || len(raw) == 0 {
		return raw
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return raw
	}
	if !restoreToolNamesInMap(payload, renameMap) {
		return raw
	}
	out, err := MarshalNoEscapeHTML(payload)
	if err != nil {
		return raw
	}
	return out
}

func restoreToolNamesInMap(payload map[string]any, renameMap map[string]string) bool {
	changed := false

	// Claude streaming content_block_start.
	if t, _ := payload["type"].(string); t == "content_block_start" {
		if cb, ok := payload["content_block"].(map[string]any); ok {
			if ct, _ := cb["type"].(string); ct == "tool_use" {
				if name, _ := cb["name"].(string); name != "" {
					if orig, ok := renameMap[name]; ok {
						cb["name"] = orig
						changed = true
					}
				}
			}
		}
	}

	// Claude non-streaming message content blocks.
	if blocks, ok := payload["content"].([]any); ok {
		for _, raw := range blocks {
			if b, ok := raw.(map[string]any); ok {
				if bt, _ := b["type"].(string); bt == "tool_use" {
					if name, _ := b["name"].(string); name != "" {
						if orig, ok := renameMap[name]; ok {
							b["name"] = orig
							changed = true
						}
					}
				}
			}
		}
	}

	// OpenAI Chat Completions (delta & message).
	if choices, ok := payload["choices"].([]any); ok {
		for _, rawChoice := range choices {
			choice, ok := rawChoice.(map[string]any)
			if !ok {
				continue
			}
			for _, holder := range []string{"delta", "message"} {
				hMap, ok := choice[holder].(map[string]any)
				if !ok {
					continue
				}
				calls, ok := hMap["tool_calls"].([]any)
				if !ok {
					continue
				}
				for _, rawCall := range calls {
					call, ok := rawCall.(map[string]any)
					if !ok {
						continue
					}
					if fn, ok := call["function"].(map[string]any); ok {
						if name, _ := fn["name"].(string); name != "" {
							if orig, ok := renameMap[name]; ok {
								fn["name"] = orig
								changed = true
							}
						}
					}
				}
			}
		}
	}

	// OpenAI Responses final JSON output array.
	if outputs, ok := payload["output"].([]any); ok {
		for _, raw := range outputs {
			if item, ok := raw.(map[string]any); ok {
				if it, _ := item["type"].(string); it == "function_call" {
					if name, _ := item["name"].(string); name != "" {
						if orig, ok := renameMap[name]; ok {
							item["name"] = orig
							changed = true
						}
					}
				}
			}
		}
	}

	// OpenAI Responses SSE event item.
	if item, ok := payload["item"].(map[string]any); ok {
		if it, _ := item["type"].(string); it == "function_call" {
			if name, _ := item["name"].(string); name != "" {
				if orig, ok := renameMap[name]; ok {
					item["name"] = orig
					changed = true
				}
			}
		}
	}

	return changed
}

type sseRestoringReader struct {
	src       io.ReadCloser
	reader    *bufio.Reader
	renameMap map[string]string
	buf       bytes.Buffer
	eof       bool
}

// NewToolNameRestoringSSEReader wraps src and restores any aliased tool names on
// SSE `data: {...}` lines. When renameMap is empty, src is returned untouched.
func NewToolNameRestoringSSEReader(src io.ReadCloser, renameMap map[string]string) io.ReadCloser {
	if len(renameMap) == 0 || src == nil {
		return src
	}
	return &sseRestoringReader{
		src:       src,
		reader:    bufio.NewReaderSize(src, 16<<10),
		renameMap: renameMap,
	}
}

func (r *sseRestoringReader) Read(p []byte) (int, error) {
	for r.buf.Len() == 0 {
		if r.eof {
			return 0, io.EOF
		}
		line, err := r.reader.ReadBytes('\n')
		if len(line) > 0 {
			trimmed := bytes.TrimSpace(line)
			if bytes.HasPrefix(trimmed, []byte("data:")) {
				data := bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data:")))
				if len(data) > 0 && data[0] == '{' {
					restored := RestoreToolNamesJSON(data, r.renameMap)
					r.buf.WriteString("data: ")
					r.buf.Write(restored)
					if bytes.HasSuffix(line, []byte("\n")) {
						r.buf.WriteByte('\n')
					}
				} else {
					r.buf.Write(line)
				}
			} else {
				r.buf.Write(line)
			}
		}
		if err == io.EOF {
			r.eof = true
		} else if err != nil {
			if r.buf.Len() > 0 {
				break
			}
			return 0, err
		}
	}
	return r.buf.Read(p)
}

func (r *sseRestoringReader) Close() error {
	return r.src.Close()
}

type aggregatedToolCall struct {
	id   string
	name string
	args strings.Builder
}

// AggregateOpenAISSEToJSON converts a forced-stream OpenAI Chat Completions SSE
// response from OpenCode Free into a standard non-streaming Chat Completion JSON
// body, restoring any aliased tool names.
func AggregateOpenAISSEToJSON(sseBody []byte, fallbackModel string, renameMap map[string]string) ([]byte, error) {
	var (
		respID        string
		model         = fallbackModel
		created       int64
		finishReason  string
		contentBuf    strings.Builder
		reasoningBuf  strings.Builder
		toolCalls     = make(map[int]*aggregatedToolCall)
		toolCallOrder []int
		usage         map[string]any
		sawEvent      bool
	)

	scanner := bufio.NewScanner(bytes.NewReader(sseBody))
	scanner.Buffer(make([]byte, 0, 64<<10), maxSSELineBytes)
	for scanner.Scan() {
		trimmed := bytes.TrimSpace(scanner.Bytes())
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
		sawEvent = true
		if respID == "" {
			if id, ok := chunk["id"].(string); ok {
				respID = id
			}
		}
		if m, ok := chunk["model"].(string); ok && strings.TrimSpace(m) != "" {
			model = m
		}
		if c, ok := chunk["created"].(float64); ok && created == 0 {
			created = int64(c)
		}
		if u, ok := chunk["usage"].(map[string]any); ok && len(u) > 0 {
			usage = u
		}
		choices, _ := chunk["choices"].([]any)
		if len(choices) == 0 {
			continue
		}
		first, ok := choices[0].(map[string]any)
		if !ok {
			continue
		}
		if fr, ok := first["finish_reason"].(string); ok && fr != "" {
			finishReason = fr
		}
		delta, ok := first["delta"].(map[string]any)
		if !ok {
			continue
		}
		if c, ok := delta["content"].(string); ok && c != "" {
			contentBuf.WriteString(c)
		}
		if rc, ok := delta["reasoning_content"].(string); ok && rc != "" {
			reasoningBuf.WriteString(rc)
		}
		if tcs, ok := delta["tool_calls"].([]any); ok {
			for _, rawTC := range tcs {
				tc, ok := rawTC.(map[string]any)
				if !ok {
					continue
				}
				idx := 0
				if f, ok := tc["index"].(float64); ok {
					idx = int(f)
				}
				entry, exists := toolCalls[idx]
				if !exists {
					entry = &aggregatedToolCall{}
					toolCalls[idx] = entry
					toolCallOrder = append(toolCallOrder, idx)
				}
				if id, ok := tc["id"].(string); ok && id != "" {
					entry.id = id
				}
				if fn, ok := tc["function"].(map[string]any); ok {
					if name, ok := fn["name"].(string); ok && name != "" {
						entry.name += name
					}
					if args, ok := fn["arguments"].(string); ok && args != "" {
						entry.args.WriteString(args)
					}
				}
			}
		}
	}

	if !sawEvent {
		// If upstream unexpectedly returned plain JSON, restore tool names and return as-is.
		trimmed := bytes.TrimSpace(sseBody)
		if len(trimmed) > 0 && trimmed[0] == '{' && json.Valid(trimmed) {
			return RestoreToolNamesJSON(trimmed, renameMap), nil
		}
		return nil, fmt.Errorf("opencode: empty or invalid SSE stream")
	}

	if respID == "" {
		respID = "chatcmpl-oc-" + RandomHex(8)
	}
	if finishReason == "" {
		if len(toolCallOrder) > 0 {
			finishReason = "tool_calls"
		} else {
			finishReason = "stop"
		}
	}

	msg := map[string]any{
		"role":    "assistant",
		"content": contentBuf.String(),
	}
	if reasoningBuf.Len() > 0 {
		msg["reasoning_content"] = reasoningBuf.String()
	}
	if len(toolCallOrder) > 0 {
		slices.Sort(toolCallOrder)
		calls := make([]map[string]any, 0, len(toolCallOrder))
		for _, idx := range toolCallOrder {
			entry := toolCalls[idx]
			name := entry.name
			if orig, ok := renameMap[name]; ok {
				name = orig
			}
			calls = append(calls, map[string]any{
				"id":   entry.id,
				"type": "function",
				"function": map[string]any{
					"name":      name,
					"arguments": entry.args.String(),
				},
			})
		}
		msg["tool_calls"] = calls
	}

	resp := map[string]any{
		"id":      respID,
		"object":  "chat.completion",
		"created": created,
		"model":   model,
		"choices": []map[string]any{
			{
				"index":         0,
				"message":       msg,
				"finish_reason": finishReason,
			},
		},
	}
	if usage != nil {
		resp["usage"] = usage
	}
	return MarshalNoEscapeHTML(resp)
}

// AggregateSSEToJSON aggregates a forced-stream OpenCode Free SSE response into
// the non-streaming JSON shape matching protocol (`openai-chat`, `openai-responses`,
// or `anthropic`), restoring any aliased tool names.
func AggregateSSEToJSON(protocol string, sseBody []byte, fallbackModel string, renameMap map[string]string) ([]byte, error) {
	switch strings.ToLower(strings.TrimSpace(protocol)) {
	case "openai-responses":
		return aggregateResponsesSSEToJSON(sseBody, renameMap)
	case "anthropic":
		return aggregateAnthropicSSEToJSON(sseBody, fallbackModel, renameMap)
	default:
		return AggregateOpenAISSEToJSON(sseBody, fallbackModel, renameMap)
	}
}

func aggregateResponsesSSEToJSON(sseBody []byte, renameMap map[string]string) ([]byte, error) {
	var completed map[string]any
	scanner := bufio.NewScanner(bytes.NewReader(sseBody))
	scanner.Buffer(make([]byte, 0, 64<<10), maxSSELineBytes)
	for scanner.Scan() {
		trimmed := bytes.TrimSpace(scanner.Bytes())
		if !bytes.HasPrefix(trimmed, []byte("data:")) {
			continue
		}
		data := bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data:")))
		if len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) || data[0] != '{' {
			continue
		}
		var evt map[string]any
		if err := json.Unmarshal(data, &evt); err != nil {
			continue
		}
		if t, _ := evt["type"].(string); t == "response.completed" || t == "response.done" {
			if respObj, ok := evt["response"].(map[string]any); ok {
				completed = respObj
			}
		}
	}
	if completed != nil {
		restoreToolNamesInMap(completed, renameMap)
		return MarshalNoEscapeHTML(completed)
	}
	trimmed := bytes.TrimSpace(sseBody)
	if len(trimmed) > 0 && trimmed[0] == '{' && json.Valid(trimmed) {
		return RestoreToolNamesJSON(trimmed, renameMap), nil
	}
	return nil, fmt.Errorf("opencode: responses SSE missing response.completed event")
}

type aggregatedAnthropicBlock struct {
	blockType string
	id        string
	name      string
	text      strings.Builder
	thinking  strings.Builder
	inputJSON strings.Builder
}

func aggregateAnthropicSSEToJSON(sseBody []byte, fallbackModel string, renameMap map[string]string) ([]byte, error) {
	var (
		msgID      string
		model      = fallbackModel
		stopReason string
		usage      = map[string]any{}
		blocks     = make(map[int]*aggregatedAnthropicBlock)
		blockOrder []int
		sawEvent   bool
	)

	scanner := bufio.NewScanner(bytes.NewReader(sseBody))
	scanner.Buffer(make([]byte, 0, 64<<10), maxSSELineBytes)
	for scanner.Scan() {
		trimmed := bytes.TrimSpace(scanner.Bytes())
		if !bytes.HasPrefix(trimmed, []byte("data:")) {
			continue
		}
		data := bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data:")))
		if len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) || data[0] != '{' {
			continue
		}
		var evt map[string]any
		if err := json.Unmarshal(data, &evt); err != nil {
			continue
		}
		sawEvent = true
		evtType, _ := evt["type"].(string)
		switch evtType {
		case "message_start":
			if msg, ok := evt["message"].(map[string]any); ok {
				if id, ok := msg["id"].(string); ok && id != "" {
					msgID = id
				}
				if m, ok := msg["model"].(string); ok && m != "" {
					model = m
				}
				if u, ok := msg["usage"].(map[string]any); ok {
					for k, v := range u {
						usage[k] = v
					}
				}
			}
		case "content_block_start":
			idx := 0
			if f, ok := evt["index"].(float64); ok {
				idx = int(f)
			}
			cb, _ := evt["content_block"].(map[string]any)
			bt, _ := cb["type"].(string)
			entry := &aggregatedAnthropicBlock{blockType: bt}
			if id, ok := cb["id"].(string); ok {
				entry.id = id
			}
			if name, ok := cb["name"].(string); ok {
				entry.name = name
			}
			if t, ok := cb["text"].(string); ok && t != "" {
				entry.text.WriteString(t)
			}
			if th, ok := cb["thinking"].(string); ok && th != "" {
				entry.thinking.WriteString(th)
			}
			blocks[idx] = entry
			if !slices.Contains(blockOrder, idx) {
				blockOrder = append(blockOrder, idx)
			}
		case "content_block_delta":
			idx := 0
			if f, ok := evt["index"].(float64); ok {
				idx = int(f)
			}
			entry, exists := blocks[idx]
			if !exists {
				entry = &aggregatedAnthropicBlock{blockType: "text"}
				blocks[idx] = entry
				blockOrder = append(blockOrder, idx)
			}
			delta, _ := evt["delta"].(map[string]any)
			dt, _ := delta["type"].(string)
			switch dt {
			case "text_delta":
				if t, ok := delta["text"].(string); ok {
					entry.text.WriteString(t)
				}
			case "thinking_delta":
				if th, ok := delta["thinking"].(string); ok {
					entry.thinking.WriteString(th)
				}
			case "input_json_delta":
				if pj, ok := delta["partial_json"].(string); ok {
					entry.inputJSON.WriteString(pj)
				}
			}
		case "message_delta":
			if delta, ok := evt["delta"].(map[string]any); ok {
				if sr, ok := delta["stop_reason"].(string); ok && sr != "" {
					stopReason = sr
				}
			}
			if u, ok := evt["usage"].(map[string]any); ok {
				for k, v := range u {
					usage[k] = v
				}
			}
		}
	}

	if !sawEvent {
		trimmed := bytes.TrimSpace(sseBody)
		if len(trimmed) > 0 && trimmed[0] == '{' && json.Valid(trimmed) {
			return RestoreToolNamesJSON(trimmed, renameMap), nil
		}
		return nil, fmt.Errorf("opencode: empty or invalid anthropic SSE stream")
	}

	if msgID == "" {
		msgID = "msg_oc_" + RandomHex(8)
	}
	slices.Sort(blockOrder)
	content := make([]map[string]any, 0, len(blockOrder))
	for _, idx := range blockOrder {
		b := blocks[idx]
		switch b.blockType {
		case "tool_use":
			name := b.name
			if orig, ok := renameMap[name]; ok {
				name = orig
			}
			var input any = map[string]any{}
			if raw := strings.TrimSpace(b.inputJSON.String()); raw != "" {
				var parsed any
				if json.Unmarshal([]byte(raw), &parsed) == nil {
					input = parsed
				}
			}
			content = append(content, map[string]any{
				"type":  "tool_use",
				"id":    b.id,
				"name":  name,
				"input": input,
			})
		case "thinking":
			content = append(content, map[string]any{
				"type":     "thinking",
				"thinking": b.thinking.String(),
			})
		default:
			content = append(content, map[string]any{
				"type": "text",
				"text": b.text.String(),
			})
		}
	}
	if stopReason == "" {
		stopReason = "end_turn"
	}

	resp := map[string]any{
		"id":          msgID,
		"type":        "message",
		"role":        "assistant",
		"model":       model,
		"content":     content,
		"stop_reason": stopReason,
		"usage":       usage,
	}
	return MarshalNoEscapeHTML(resp)
}

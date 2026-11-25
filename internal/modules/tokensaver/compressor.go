package tokensaver

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	StrategyOff        = "off"
	StrategyWhitespace = "whitespace"
	StrategyDedupe     = "dedupe"
	StrategyAggressive = "aggressive"

	LevelLite  = "lite"
	LevelFull  = "full"
	LevelUltra = "ultra"
)

func isValidLevel(v string) bool {
	return v == LevelLite || v == LevelFull || v == LevelUltra
}

var multiSpaceRe = regexp.MustCompile(`[ \t]{2,}`)

// estimateTokens gives a fast zero-table approximation (~4 chars per token).
func estimateTokens(b []byte) int64 {
	if len(b) == 0 {
		return 0
	}
	return int64((len(b) + 3) / 4)
}

// compressText normalizes redundant whitespace outside line-anchored fenced
// code blocks while preserving newlines adjacent to fences.
func compressText(input string) string {
	lines := strings.Split(input, "\n")
	out := make([]string, 0, len(lines))
	inFence := false
	fenceLen := 0
	blankRun := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if ticks := leadingBackticks(trimmed); ticks >= 3 {
			if !inFence {
				inFence = true
				fenceLen = ticks
				blankRun = 0
				out = append(out, line)
				continue
			}
			if ticks >= fenceLen && strings.Trim(trimmed, "`") == "" {
				inFence = false
				fenceLen = 0
				blankRun = 0
				out = append(out, line)
				continue
			}
		}
		if inFence {
			out = append(out, line)
			continue
		}
		cleaned := strings.TrimRight(line, " \t")
		cleaned = multiSpaceRe.ReplaceAllString(cleaned, " ")
		if cleaned == "" {
			blankRun++
			if blankRun > 1 {
				continue
			}
		} else {
			blankRun = 0
		}
		out = append(out, cleaned)
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}

func leadingBackticks(s string) int {
	n := 0
	for n < len(s) && s[n] == '`' {
		n++
	}
	return n
}

func truncateToolOutput(s string, maxChars int) string {
	runes := []rune(s)
	if maxChars <= 0 || len(runes) <= maxChars {
		return s
	}
	omitted := len(runes) - maxChars
	head := maxChars / 2
	tail := maxChars - head
	return fmt.Sprintf("%s\n[...truncated %d chars...]\n%s", string(runes[:head]), omitted, string(runes[len(runes)-tail:]))
}

// compressPayload applies the configured Go-native token compression strategy
// to a JSON inference payload without corrupting fenced code blocks.
func compressPayload(body []byte, cfg Config) (out []byte, bytesSaved, tokensSaved int64, mutated bool) {
	if cfg.Strategy == StrategyOff || len(body) == 0 {
		return body, 0, 0, false
	}
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil || root == nil {
		return body, 0, 0, false
	}

	if msgs, ok := root["messages"].([]any); ok {
		root["messages"] = compressMessages(msgs, cfg)
	}
	switch sys := root["system"].(type) {
	case string:
		root["system"] = compressText(sys)
	case []any:
		root["system"] = compressContentBlocks(sys, false, cfg)
	}
	if inst, ok := root["instructions"].(string); ok {
		root["instructions"] = compressText(inst)
	}
	switch inp := root["input"].(type) {
	case string:
		root["input"] = compressText(inp)
	case []any:
		root["input"] = compressMessages(inp, cfg)
	}
	for _, siKey := range []string{"systemInstruction", "system_instruction"} {
		if si, ok := root[siKey].(map[string]any); ok {
			if parts, ok := si["parts"].([]any); ok {
				for j, pr := range parts {
					if pm, ok := pr.(map[string]any); ok {
						if txt, ok := pm["text"].(string); ok {
							pm["text"] = compressText(txt)
							parts[j] = pm
						}
					}
				}
				si["parts"] = parts
			}
			root[siKey] = si
		}
	}
	if contents, ok := root["contents"].([]any); ok {
		for i, raw := range contents {
			item, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			if parts, ok := item["parts"].([]any); ok {
				for j, pr := range parts {
					pm, ok := pr.(map[string]any)
					if !ok {
						continue
					}
					if txt, ok := pm["text"].(string); ok {
						pm["text"] = compressText(txt)
						parts[j] = pm
					}
				}
				item["parts"] = parts
			}
			contents[i] = item
		}
		root["contents"] = contents
	}

	encoded, err := json.Marshal(root)
	if err != nil || len(encoded) >= len(body) {
		return body, 0, 0, false
	}
	savedBytes := int64(len(body) - len(encoded))
	savedTokens := estimateTokens(body) - estimateTokens(encoded)
	if savedTokens < 0 {
		savedTokens = 0
	}
	return encoded, savedBytes, savedTokens, true
}

func compressMessages(msgs []any, cfg Config) []any {
	out := make([]any, 0, len(msgs))
	var prevSig string

	for _, raw := range msgs {
		msg, ok := raw.(map[string]any)
		if !ok {
			out = append(out, raw)
			continue
		}
		role, _ := msg["role"].(string)
		isTool := role == "tool"

		switch content := msg["content"].(type) {
		case string:
			if cfg.Strategy == StrategyAggressive && isTool {
				msg["content"] = compressText(filterToolOutput(content, cfg.MaxToolOutputChars))
			} else {
				msg["content"] = compressText(content)
			}
		case []any:
			msg["content"] = compressContentBlocks(content, isTool, cfg)
		}

		if cfg.Strategy == StrategyDedupe || cfg.Strategy == StrategyAggressive {
			sigBytes, _ := json.Marshal(msg)
			sig := string(sigBytes)
			if sig == prevSig {
				continue
			}
			prevSig = sig
		}
		out = append(out, msg)
	}

	if cfg.TrimHistoryEnabled && cfg.MaxHistoryMessages > 2 && len(out) > cfg.MaxHistoryMessages {
		out = trimChatHistoryMessages(out, cfg.MaxHistoryMessages)
	}

	return out
}

func trimChatHistoryMessages(msgs []any, maxHistory int) []any {
	if maxHistory <= 2 || len(msgs) <= maxHistory {
		return msgs
	}
	var leading []any
	restIdx := 0
	for restIdx < len(msgs) {
		m, ok := msgs[restIdx].(map[string]any)
		if !ok {
			break
		}
		if role, _ := m["role"].(string); role == "system" || role == "developer" {
			leading = append(leading, m)
			restIdx++
		} else {
			break
		}
	}
	if restIdx >= len(msgs) {
		return msgs
	}
	keepTail := maxHistory - len(leading)
	if keepTail < 1 {
		keepTail = 1
	}
	start := len(msgs) - keepTail
	if start < restIdx {
		start = restIdx
	}
	// Avoid starting the retained tail on an orphaned tool result.
	for start < len(msgs)-1 {
		m, ok := msgs[start].(map[string]any)
		if !ok {
			break
		}
		if role, _ := m["role"].(string); role == "tool" {
			start++
			continue
		}
		break
	}
	tail := msgs[start:]
	combined := make([]any, 0, len(leading)+len(tail))
	combined = append(combined, leading...)
	combined = append(combined, tail...)
	return combined
}

func trimPayloadHistory(body []byte, maxHistory int) (out []byte, bytesSaved, tokensSaved int64, mutated bool) {
	if maxHistory <= 2 || len(body) == 0 {
		return body, 0, 0, false
	}
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil || root == nil {
		return body, 0, 0, false
	}
	changed := false
	if msgs, ok := root["messages"].([]any); ok && len(msgs) > maxHistory {
		root["messages"] = trimChatHistoryMessages(msgs, maxHistory)
		changed = true
	}
	if inp, ok := root["input"].([]any); ok && len(inp) > maxHistory {
		root["input"] = trimChatHistoryMessages(inp, maxHistory)
		changed = true
	}
	if !changed {
		return body, 0, 0, false
	}
	encoded, err := json.Marshal(root)
	if err != nil || len(encoded) >= len(body) {
		return body, 0, 0, false
	}
	savedBytes := int64(len(body) - len(encoded))
	savedTokens := estimateTokens(body) - estimateTokens(encoded)
	if savedTokens < 0 {
		savedTokens = 0
	}
	return encoded, savedBytes, savedTokens, true
}

const (
	cavemanSharedRules = "Code blocks, identifiers, file paths, commands, errors, URLs: keep exact verbatim. " +
		"Auto-Clarity: switch to clear full sentences for security warnings, irreversible actions, or ambiguous multi-step sequences, then resume terse style. " +
		"Preserve the user's dominant language (e.g. Indonesian -> reply Indonesian, English -> reply English). " +
		"No invented abbreviations (standard tech acronyms like API/DB/HTTP/JSON OK). " +
		"No self-reference (never announce the style), no decorative emoji, no filler/pleasantries/tool narration. " +
		"Active every response; never drift back to verbose."

	ponytailSharedRules = "You are a lazy senior developer: lazy means efficient, not careless; the best code is code never written. " +
		"Before writing code, stop at the first rung that holds: 1) Does it need to exist? (YAGNI) 2) Stdlib does it? Use it. 3) Native platform covers it? Use it. 4) Existing dependency solves it? Never add a new one for a few lines. 5) One-liner? Ship it. 6) Only then: minimum working code. " +
		"No unrequested abstractions (no single-impl interface, no factory for one type, no config for constants). Fewest files, shortest working diff, deletion over addition, boring over clever. Mark deliberate simplifications with a `ponytail:` comment naming the ceiling and upgrade path. " +
		"Code first, then at most three short lines: what was skipped and when to add it. " +
		"Never simplify away: trust-boundary input validation, error handling preventing data loss, security, accessibility, or explicitly requested behavior. Non-trivial logic keeps one runnable check/test. " +
		"Active every response; no drift back to over-building."
)

func buildOutputDirective(cfg Config) string {
	var parts []string
	if cfg.CavemanEnabled {
		switch cfg.CavemanLevel {
		case LevelLite:
			parts = append(parts, "[caveman:lite] Respond tersely. Keep grammar and full sentences; drop filler, hedging, and pleasantries. Pattern: state thing, action, reason, next step. "+cavemanSharedRules)
		case LevelUltra:
			parts = append(parts, "[caveman:ultra] Respond ultra-terse and telegraphic. Maximum compression: strip articles, filler, and conjunctions; one word when enough. Pattern: [thing] [action] [reason]. [next step]. "+cavemanSharedRules)
		default:
			parts = append(parts, "[caveman:full] Respond terse. Technical substance stays exact, fluff dies. Drop articles (a/an/the), filler, pleasantries, hedging; fragments OK; short synonyms. Pattern: [thing] [action] [reason]. [next step]. "+cavemanSharedRules)
		}
	}
	if cfg.PonytailEnabled {
		switch cfg.PonytailLevel {
		case LevelLite:
			parts = append(parts, "[ponytail:lite] Build what was asked, and name the lazier alternative in one line so the user can pick. "+ponytailSharedRules)
		case LevelUltra:
			parts = append(parts, "[ponytail:ultra] YAGNI extremist: deletion before addition; ship the shortest working solution and challenge unnecessary scope in the same response. "+ponytailSharedRules)
		default:
			parts = append(parts, "[ponytail:full] Enforce the 6-rung ladder strictly: stdlib and native first, shortest diff, shortest explanation. "+ponytailSharedRules)
		}
	}
	return strings.Join(parts, "\n")
}

// injectOutputDirectives appends Caveman / Ponytail system directives idempotently
// using the wire-valid system field for each protocol.
func injectOutputDirectives(protocol string, body []byte, cfg Config) ([]byte, bool) {
	directive := buildOutputDirective(cfg)
	if directive == "" || len(body) == 0 {
		return body, false
	}
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil || root == nil {
		return body, false
	}

	marker := "[caveman:"
	if !cfg.CavemanEnabled {
		marker = "[ponytail:"
	}

	mutated := false
	switch {
	case protocol == "anthropic" || root["system"] != nil:
		switch sys := root["system"].(type) {
		case nil:
			root["system"] = directive
			mutated = true
		case string:
			if strings.Contains(sys, marker) {
				return body, false
			}
			root["system"] = strings.TrimSpace(sys + "\n\n" + directive)
			mutated = true
		case []any:
			for _, item := range sys {
				if bm, ok := item.(map[string]any); ok {
					if txt, _ := bm["text"].(string); strings.Contains(txt, marker) {
						return body, false
					}
				}
			}
			root["system"] = append(sys, map[string]any{"type": "text", "text": directive})
			mutated = true
		}

	case protocol == "openai-responses" || root["instructions"] != nil || root["input"] != nil:
		inst, _ := root["instructions"].(string)
		if strings.Contains(inst, marker) {
			return body, false
		}
		if strings.TrimSpace(inst) == "" {
			root["instructions"] = directive
		} else {
			root["instructions"] = strings.TrimSpace(inst + "\n\n" + directive)
		}
		mutated = true

	case protocol == "gemini" || root["contents"] != nil:
		siKey := "systemInstruction"
		if _, ok := root["system_instruction"]; ok {
			siKey = "system_instruction"
		}
		si, _ := root[siKey].(map[string]any)
		if si == nil {
			si = map[string]any{"parts": []any{}}
		}
		parts, _ := si["parts"].([]any)
		for _, pr := range parts {
			if pm, ok := pr.(map[string]any); ok {
				if txt, _ := pm["text"].(string); strings.Contains(txt, marker) {
					return body, false
				}
			}
		}
		si["parts"] = append(parts, map[string]any{"text": directive})
		root[siKey] = si
		mutated = true

	default:
		msgs, ok := root["messages"].([]any)
		if !ok {
			return body, false
		}
		for _, raw := range msgs {
			if m, ok := raw.(map[string]any); ok {
				if c, _ := m["content"].(string); strings.Contains(c, marker) {
					return body, false
				}
			}
		}
		if len(msgs) > 0 {
			if first, ok := msgs[0].(map[string]any); ok {
				role, _ := first["role"].(string)
				if role == "system" || role == "developer" {
					if content, ok := first["content"].(string); ok {
						first["content"] = strings.TrimSpace(content + "\n\n" + directive)
						msgs[0] = first
						root["messages"] = msgs
						mutated = true
					}
				}
			}
		}
		if !mutated {
			sysMsg := map[string]any{"role": "system", "content": directive}
			root["messages"] = append([]any{sysMsg}, msgs...)
			mutated = true
		}
	}

	if !mutated {
		return body, false
	}
	encoded, err := json.Marshal(root)
	if err != nil {
		return body, false
	}
	return encoded, true
}

func compressContentBlocks(blocks []any, isTool bool, cfg Config) []any {
	for i, raw := range blocks {
		bm, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		btype, _ := bm["type"].(string)
		if txt, ok := bm["text"].(string); ok {
			if cfg.Strategy == StrategyAggressive && (isTool || btype == "tool_result") {
				bm["text"] = compressText(filterToolOutput(txt, cfg.MaxToolOutputChars))
			} else {
				bm["text"] = compressText(txt)
			}
		}
		if contentStr, ok := bm["content"].(string); ok && btype == "tool_result" {
			if cfg.Strategy == StrategyAggressive {
				bm["content"] = compressText(filterToolOutput(contentStr, cfg.MaxToolOutputChars))
			} else {
				bm["content"] = compressText(contentStr)
			}
		}
		blocks[i] = bm
	}
	return blocks
}

// injectAnthropicCacheControl adds ephemeral cache_control to the last system
// content block and the last tool definition when eligible and not already present.
func injectAnthropicCacheControl(body []byte) ([]byte, bool) {
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil || root == nil {
		return body, false
	}

	mutated := false
	switch sys := root["system"].(type) {
	case string:
		if strings.TrimSpace(sys) != "" {
			root["system"] = []any{
				map[string]any{
					"type":          "text",
					"text":          sys,
					"cache_control": map[string]any{"type": "ephemeral"},
				},
			}
			mutated = true
		}
	case []any:
		if len(sys) > 0 {
			if last, ok := sys[len(sys)-1].(map[string]any); ok {
				if _, exists := last["cache_control"]; !exists {
					last["cache_control"] = map[string]any{"type": "ephemeral"}
					sys[len(sys)-1] = last
					root["system"] = sys
					mutated = true
				}
			}
		}
	}

	if tools, ok := root["tools"].([]any); ok && len(tools) > 0 {
		if lastTool, ok := tools[len(tools)-1].(map[string]any); ok {
			if _, exists := lastTool["cache_control"]; !exists {
				lastTool["cache_control"] = map[string]any{"type": "ephemeral"}
				tools[len(tools)-1] = lastTool
				root["tools"] = tools
				mutated = true
			}
		}
	}

	if !mutated {
		return body, false
	}
	encoded, err := json.Marshal(root)
	if err != nil {
		return body, false
	}
	return encoded, true
}

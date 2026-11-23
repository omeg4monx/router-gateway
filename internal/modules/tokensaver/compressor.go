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
)

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
			c := compressText(content)
			if cfg.Strategy == StrategyAggressive && isTool {
				c = truncateToolOutput(c, cfg.MaxToolOutputChars)
			}
			msg["content"] = c
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

	if cfg.Strategy == StrategyAggressive && cfg.MaxHistoryMessages > 2 && len(out) > cfg.MaxHistoryMessages {
		var leading []any
		restIdx := 0
		for restIdx < len(out) {
			m, ok := out[restIdx].(map[string]any)
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
		if restIdx < len(out) {
			keepTail := cfg.MaxHistoryMessages - len(leading)
			if keepTail < 1 {
				keepTail = 1
			}
			start := len(out) - keepTail
			if start < restIdx {
				start = restIdx
			}
			// Avoid starting the retained tail on an orphaned tool result.
			for start < len(out)-1 {
				m, ok := out[start].(map[string]any)
				if !ok {
					break
				}
				if role, _ := m["role"].(string); role == "tool" {
					start++
					continue
				}
				break
			}
			tail := out[start:]
			combined := make([]any, 0, len(leading)+len(tail))
			combined = append(combined, leading...)
			combined = append(combined, tail...)
			out = combined
		}
	}

	return out
}

func compressContentBlocks(blocks []any, isTool bool, cfg Config) []any {
	for i, raw := range blocks {
		bm, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		btype, _ := bm["type"].(string)
		if txt, ok := bm["text"].(string); ok {
			c := compressText(txt)
			if cfg.Strategy == StrategyAggressive && (isTool || btype == "tool_result") {
				c = truncateToolOutput(c, cfg.MaxToolOutputChars)
			}
			bm["text"] = c
		}
		if contentStr, ok := bm["content"].(string); ok && btype == "tool_result" {
			c := compressText(contentStr)
			if cfg.Strategy == StrategyAggressive {
				c = truncateToolOutput(c, cfg.MaxToolOutputChars)
			}
			bm["content"] = c
		}
		blocks[i] = bm
	}
	return blocks
}

// injectAnthropicCacheControl adds ephemeral cache_control to the last system
// content block when eligible and not already present.
func injectAnthropicCacheControl(body []byte) ([]byte, bool) {
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil || root == nil {
		return body, false
	}
	switch sys := root["system"].(type) {
	case string:
		if strings.TrimSpace(sys) == "" {
			return body, false
		}
		root["system"] = []any{
			map[string]any{
				"type":          "text",
				"text":          sys,
				"cache_control": map[string]any{"type": "ephemeral"},
			},
		}
	case []any:
		if len(sys) == 0 {
			return body, false
		}
		last, ok := sys[len(sys)-1].(map[string]any)
		if !ok {
			return body, false
		}
		if _, exists := last["cache_control"]; exists {
			return body, false
		}
		last["cache_control"] = map[string]any{"type": "ephemeral"}
		sys[len(sys)-1] = last
		root["system"] = sys
	default:
		return body, false
	}
	encoded, err := json.Marshal(root)
	if err != nil {
		return body, false
	}
	return encoded, true
}

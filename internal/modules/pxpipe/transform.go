package pxpipe

import (
	"encoding/json"
	"strings"
)

type transformResult struct {
	out           []byte
	body          map[string]any
	stagesApplied []string
	changed       bool
}

func runTransformPipeline(raw []byte, cfg Config) (transformResult, error) {
	if len(raw) == 0 {
		return transformResult{out: raw}, nil
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		// Non-JSON bodies pass through untouched without error.
		return transformResult{out: raw}, nil
	}

	stages := make([]string, 0, 4)

	if cfg.StripNullFields {
		if stripNulls(payload) {
			stages = append(stages, "strip_null_fields")
		}
	}

	if msgs, ok := payload["messages"].([]any); ok && len(msgs) > 0 {
		if cfg.CollapseWhitespace {
			if collapseMessageWhitespace(msgs) {
				stages = append(stages, "collapse_whitespace")
			}
		}
		if cfg.DedupeSystemPrompts {
			deduped, changed := dedupeSystemMessages(msgs)
			if changed {
				msgs = deduped
				payload["messages"] = msgs
				stages = append(stages, "dedupe_system_prompts")
			}
		}
		if cfg.MaxHistoryMessages > 0 {
			windowed, changed := windowHistoryMessages(msgs, cfg.MaxHistoryMessages)
			if changed {
				payload["messages"] = windowed
				stages = append(stages, "max_history_window")
			}
		}
	}

	if sys, ok := payload["system"].(string); ok && cfg.CollapseWhitespace {
		collapsed := collapseRuns(sys)
		if collapsed != sys {
			payload["system"] = collapsed
			if !containsString(stages, "collapse_whitespace") {
				stages = append(stages, "collapse_whitespace")
			}
		}
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return transformResult{}, err
	}

	changed := len(stages) > 0 || len(encoded) < len(raw)
	if len(encoded) < len(raw) && !containsString(stages, "json_minify") {
		stages = append(stages, "json_minify")
	}

	return transformResult{
		out:           encoded,
		body:          payload,
		stagesApplied: stages,
		changed:       changed,
	}, nil
}

func stripNulls(m map[string]any) bool {
	removed := false
	for k, v := range m {
		if v == nil {
			delete(m, k)
			removed = true
			continue
		}
		if nested, ok := v.(map[string]any); ok {
			if stripNulls(nested) {
				removed = true
			}
		}
	}
	return removed
}

func collapseMessageWhitespace(msgs []any) bool {
	changed := false
	for _, item := range msgs {
		msg, ok := item.(map[string]any)
		if !ok {
			continue
		}
		switch content := msg["content"].(type) {
		case string:
			collapsed := collapseRuns(content)
			if collapsed != content {
				msg["content"] = collapsed
				changed = true
			}
		case []any:
			for _, part := range content {
				partMap, ok := part.(map[string]any)
				if !ok {
					continue
				}
				if text, ok := partMap["text"].(string); ok {
					collapsed := collapseRuns(text)
					if collapsed != text {
						partMap["text"] = collapsed
						changed = true
					}
				}
			}
		}
	}
	return changed
}

func collapseRuns(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.Join(strings.Fields(line), " ")
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func dedupeSystemMessages(msgs []any) ([]any, bool) {
	seen := make(map[string]struct{})
	out := make([]any, 0, len(msgs))
	changed := false
	for _, item := range msgs {
		msg, ok := item.(map[string]any)
		if !ok {
			out = append(out, item)
			continue
		}
		role, _ := msg["role"].(string)
		if strings.EqualFold(role, "system") {
			if content, ok := msg["content"].(string); ok {
				norm := strings.TrimSpace(content)
				if _, exists := seen[norm]; exists {
					changed = true
					continue
				}
				seen[norm] = struct{}{}
			}
		}
		out = append(out, item)
	}
	return out, changed
}

func windowHistoryMessages(msgs []any, maxNonSystem int) ([]any, bool) {
	systemMsgs := make([]any, 0, 2)
	otherMsgs := make([]any, 0, len(msgs))
	for _, item := range msgs {
		msg, ok := item.(map[string]any)
		if ok {
			role, _ := msg["role"].(string)
			if strings.EqualFold(role, "system") {
				systemMsgs = append(systemMsgs, item)
				continue
			}
		}
		otherMsgs = append(otherMsgs, item)
	}
	if len(otherMsgs) <= maxNonSystem {
		return msgs, false
	}
	trimmed := otherMsgs[len(otherMsgs)-maxNonSystem:]
	out := make([]any, 0, len(systemMsgs)+len(trimmed))
	out = append(out, systemMsgs...)
	out = append(out, trimmed...)
	return out, true
}

func containsString(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

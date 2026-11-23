package plugins

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const interceptorPriority = 60

var (
	_ contracts.GatewayInterceptor = (*Module)(nil)

	multiNewlineRe = regexp.MustCompile(`\n{3,}`)
	multiSpaceRe   = regexp.MustCompile(`[ \t]{2,}`)
)

// Priority returns the deterministic execution order in the gateway interceptor chain.
func (m *Module) Priority() int {
	return interceptorPriority
}

// OnRequest applies enabled model-scoped request-enhancement plugins to JSON gateway payloads.
func (m *Module) OnRequest(ctx context.Context, req *contracts.CanonicalRequest) error {
	if req == nil || len(req.Body) == 0 {
		return nil
	}
	states, err := m.snapshotStates(ctx)
	if err != nil {
		return nil
	}
	mutated, applied, changed := applyPluginsToBody(req.Model, req.Body, states)
	if !changed {
		return nil
	}
	req.Body = mutated
	if req.Metadata == nil {
		req.Metadata = make(map[string]any)
	}
	req.Metadata["plugins.applied"] = applied
	return nil
}

// OnResponseHeaders is a no-op for the plugins interceptor.
func (m *Module) OnResponseHeaders(context.Context, *contracts.ResponseMeta) error {
	return nil
}

// OnStreamChunk is a no-op for the plugins interceptor.
func (m *Module) OnStreamChunk(context.Context, *contracts.StreamEvent) error {
	return nil
}

func applyPluginsToBody(model string, rawBody []byte, states []PluginState) (mutated []byte, applied []string, changed bool) {
	var payload map[string]any
	if err := json.Unmarshal(rawBody, &payload); err != nil || payload == nil {
		return rawBody, nil, false
	}

	effectiveModel := strings.TrimSpace(model)
	if effectiveModel == "" {
		if mVal, ok := payload["model"].(string); ok {
			effectiveModel = strings.TrimSpace(mVal)
		}
	}

	applied = make([]string, 0)
	directives := make([]string, 0)

	for _, st := range states {
		if !st.Enabled || !modelMatches(effectiveModel, st.Models) {
			continue
		}
		switch st.Key {
		case "imageVision":
			directives = append(directives, "When images are provided, inspect concrete visual details first before summarizing.")
			applied = append(applied, st.Key)
		case "thinkDeeper":
			directives = append(directives, "Reason carefully step-by-step and verify edge cases before giving the final answer.")
			if _, exists := payload["reasoning_effort"]; !exists {
				if effort, ok := st.Config["effort"].(string); ok && strings.TrimSpace(effort) != "" {
					payload["reasoning_effort"] = strings.TrimSpace(effort)
				}
			}
			applied = append(applied, st.Key)
		case "speedMode":
			directives = append(directives, "Respond concisely and directly without unnecessary elaboration.")
			if _, exists := payload["max_tokens"]; !exists {
				if mt := intFromAny(st.Config["max_tokens"]); mt > 0 {
					payload["max_tokens"] = mt
				}
			}
			applied = append(applied, st.Key)
		case "jsonGuard":
			if _, exists := payload["response_format"]; !exists {
				payload["response_format"] = map[string]any{"type": "json_object"}
			}
			directives = append(directives, "Output strictly valid JSON with no surrounding markdown fences or commentary.")
			applied = append(applied, st.Key)
		case "contextSqueezer":
			minChars := intFromAny(st.Config["min_chars"])
			if minChars <= 0 {
				minChars = 256
			}
			squeezePayloadMessages(payload, minChars)
			applied = append(applied, st.Key)
		case "openaiToolBridge":
			if tools, ok := payload["tools"].([]any); ok && len(tools) > 0 {
				if _, exists := payload["tool_choice"]; !exists {
					payload["tool_choice"] = "auto"
				}
				if _, exists := payload["parallel_tool_calls"]; !exists {
					if ptc, ok := st.Config["parallel_tool_calls"].(bool); ok {
						payload["parallel_tool_calls"] = ptc
					}
				}
			}
			applied = append(applied, st.Key)
		case "antiSlop":
			lvl, _ := st.Config["level"].(string)
			directives = append(directives, antiSlopDirective(lvl))
			applied = append(applied, st.Key)
		}
	}

	if len(applied) == 0 {
		return rawBody, nil, false
	}

	if len(directives) > 0 {
		injectSystemDirectives(payload, strings.Join(directives, " "))
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return rawBody, nil, false
	}
	return encoded, applied, true
}

func modelMatches(model string, allowlist []string) bool {
	if len(allowlist) == 0 {
		return true
	}
	normModel := strings.ToLower(strings.TrimSpace(model))
	for _, entry := range allowlist {
		normEntry := strings.ToLower(strings.TrimSpace(entry))
		if normEntry == "" {
			continue
		}
		if normEntry == "*" || normEntry == normModel {
			return true
		}
		if strings.HasSuffix(normEntry, "*") {
			prefix := strings.TrimSuffix(normEntry, "*")
			if strings.HasPrefix(normModel, prefix) {
				return true
			}
		}
	}
	return false
}

func antiSlopDirective(level string) string {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "light":
		return "Avoid filler openings and keep answers direct."
	case "standard":
		return "Be direct and concrete. Omit filler preambles, hedging, and repetitive disclaimers."
	default:
		return "Zero filler, zero hedging, zero conversational preamble or closing pleasantries. Deliver the exact answer directly."
	}
}

func injectSystemDirectives(payload map[string]any, directive string) {
	if msgs, ok := payload["messages"].([]any); ok {
		sysMsg := map[string]any{
			"role":    "system",
			"content": directive,
		}
		updated := make([]any, 0, len(msgs)+1)
		updated = append(updated, sysMsg)
		updated = append(updated, msgs...)
		payload["messages"] = updated
		return
	}
	if sysStr, ok := payload["system"].(string); ok {
		if strings.TrimSpace(sysStr) == "" {
			payload["system"] = directive
		} else {
			payload["system"] = directive + "\n\n" + sysStr
		}
		return
	}
	if _, hasInput := payload["input"]; hasInput {
		if instr, ok := payload["instructions"].(string); ok && strings.TrimSpace(instr) != "" {
			payload["instructions"] = directive + "\n\n" + instr
		} else {
			payload["instructions"] = directive
		}
	}
}

func squeezePayloadMessages(payload map[string]any, minChars int) {
	msgs, ok := payload["messages"].([]any)
	if !ok {
		return
	}
	for _, item := range msgs {
		msgMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		content, ok := msgMap["content"].(string)
		if !ok || len(content) < minChars {
			continue
		}
		compacted := multiNewlineRe.ReplaceAllString(content, "\n\n")
		compacted = multiSpaceRe.ReplaceAllString(compacted, " ")
		msgMap["content"] = strings.TrimSpace(compacted)
	}
}

func intFromAny(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	default:
		return 0
	}
}

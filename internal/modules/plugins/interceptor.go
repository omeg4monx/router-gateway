package plugins

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const interceptorPriority = 60

var (
	_ contracts.GatewayInterceptor      = (*Module)(nil)
	_ contracts.ResponseBodyInterceptor = (*Module)(nil)
)

const (
	metaKeyApplied     = "plugins.applied"
	metaKeyVision      = "plugins.vision_active"
	metaKeyJSONGuard   = "plugins.json_guard"
	metaKeyToolBridge  = "plugins.tool_bridge"
	metaKeyToolCatalog = "plugins.tools"
)

// Priority returns the deterministic execution order in the gateway interceptor chain.
func (m *Module) Priority() int {
	return interceptorPriority
}

type requestApplyOutcome struct {
	mutated      []byte
	applied      []string
	changed      bool
	visionActive bool
	jsonGuard    bool
	toolBridge   bool
	toolCatalog  []any
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
	outcome := applyPluginsToBody(req.Model, req.Protocol, req.Body, states)
	if len(outcome.applied) == 0 {
		return nil
	}
	if outcome.changed {
		req.Body = outcome.mutated
	}
	if req.Metadata == nil {
		req.Metadata = make(map[string]any)
	}
	req.Metadata[metaKeyApplied] = outcome.applied
	if outcome.visionActive {
		req.Metadata[metaKeyVision] = true
	}
	if outcome.jsonGuard {
		req.Metadata[metaKeyJSONGuard] = true
	}
	if outcome.toolBridge {
		req.Metadata[metaKeyToolBridge] = true
	}
	if len(outcome.toolCatalog) > 0 {
		req.Metadata[metaKeyToolCatalog] = outcome.toolCatalog
	}
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

// OnResponseBody applies response-side repairs (OpenAI Tool Bridge and JSON Guard)
// to non-streaming gateway responses when those plugins were activated on the request.
func (m *Module) OnResponseBody(_ context.Context, meta *contracts.ResponseMeta, body []byte) ([]byte, error) {
	if meta == nil || len(body) == 0 || meta.Metadata == nil {
		return body, nil
	}
	jsonGuardActive, _ := meta.Metadata[metaKeyJSONGuard].(bool)
	toolBridgeActive, _ := meta.Metadata[metaKeyToolBridge].(bool)
	if !jsonGuardActive && !toolBridgeActive {
		return body, nil
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil || payload == nil {
		return body, nil
	}

	tools, _ := meta.Metadata[metaKeyToolCatalog].([]any)
	changed := false

	if toolBridgeActive {
		if applyOpenAIToolBridgeResponse(payload, tools) {
			changed = true
		}
	}
	if jsonGuardActive {
		if applyJSONGuardResponse(payload, tools) {
			changed = true
		}
	}

	if !changed {
		return body, nil
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return body, nil
	}
	return encoded, nil
}

func isClaudeFormat(protocol, model string, payload map[string]any) bool {
	p := strings.ToLower(strings.TrimSpace(protocol))
	if strings.Contains(p, "anthropic") || strings.Contains(p, "claude") {
		return true
	}
	if p == "" {
		if _, hasSystem := payload["system"]; hasSystem {
			return true
		}
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(model)), "claude") {
			return true
		}
	}
	return false
}

func payloadHasImages(payload map[string]any) bool {
	msgs, ok := payload["messages"].([]any)
	if !ok {
		return false
	}
	for _, item := range msgs {
		msg, ok := item.(map[string]any)
		if !ok {
			continue
		}
		blocks, ok := msg["content"].([]any)
		if !ok {
			continue
		}
		for _, rawBlock := range blocks {
			block, ok := rawBlock.(map[string]any)
			if !ok {
				continue
			}
			bType, _ := block["type"].(string)
			if bType == "image_url" || bType == "image" || bType == "input_image" {
				return true
			}
			if bType == "tool_result" {
				if nested, ok := block["content"].([]any); ok {
					for _, sub := range nested {
						if subMap, ok := sub.(map[string]any); ok {
							st, _ := subMap["type"].(string)
							if st == "image_url" || st == "image" {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func extractRequestTools(payload map[string]any) []any {
	var out []any
	if tools, ok := payload["tools"].([]any); ok && len(tools) > 0 {
		out = append(out, tools...)
	}
	if fns, ok := payload["functions"].([]any); ok && len(fns) > 0 {
		out = append(out, fns...)
	}
	return out
}

func applyPluginsToBody(model, protocol string, rawBody []byte, states []PluginState) requestApplyOutcome {
	var payload map[string]any
	if err := json.Unmarshal(rawBody, &payload); err != nil || payload == nil {
		return requestApplyOutcome{mutated: rawBody}
	}

	effectiveModel := strings.TrimSpace(model)
	if effectiveModel == "" {
		if mVal, ok := payload["model"].(string); ok {
			effectiveModel = strings.TrimSpace(mVal)
		}
	}
	claudeFmt := isClaudeFormat(protocol, effectiveModel, payload)

	outcome := requestApplyOutcome{
		mutated: rawBody,
		applied: make([]string, 0),
	}
	directives := make([]string, 0)
	mutatedBody := false

	for _, st := range states {
		if !st.Enabled || !modelMatches(effectiveModel, st.Models) {
			continue
		}
		switch st.Key {
		case "imageVision":
			outcome.visionActive = true
			if payloadHasImages(payload) {
				directives = append(directives, "Image Vision is active: inspect concrete visual details in all attached images before summarizing.")
				mutatedBody = true
			}
			outcome.applied = append(outcome.applied, st.Key)

		case "thinkDeeper":
			if claudeFmt {
				budget := intFromAny(st.Config["budget_tokens"])
				if budget <= 0 {
					budget = 10240
				}
				if existing, ok := payload["thinking"].(map[string]any); ok {
					if existingBudget := intFromAny(existing["budget_tokens"]); existingBudget > 0 {
						budget = existingBudget
					}
				}
				payload["thinking"] = map[string]any{
					"type":          "enabled",
					"budget_tokens": budget,
				}
				mutatedBody = true
			} else {
				effort, _ := st.Config["effort"].(string)
				effort = strings.TrimSpace(effort)
				if effort == "" {
					effort = "high"
				}
				if _, exists := payload["reasoning_effort"]; !exists {
					payload["reasoning_effort"] = effort
					mutatedBody = true
				}
				if _, exists := payload["reasoning"]; !exists {
					payload["reasoning"] = map[string]any{"effort": effort}
					mutatedBody = true
				}
			}
			outcome.applied = append(outcome.applied, st.Key)

		case "speedMode":
			if claudeFmt {
				payload["thinking"] = map[string]any{"type": "disabled"}
				delete(payload, "reasoning_effort")
				mutatedBody = true
			} else {
				payload["reasoning_effort"] = "none"
				if thinking, ok := payload["thinking"].(map[string]any); ok {
					delete(thinking, "budget_tokens")
				}
				mutatedBody = true
			}
			delete(payload, "enable_thinking")
			delete(payload, "thinking_budget")
			if _, exists := payload["max_tokens"]; !exists {
				if mt := intFromAny(st.Config["max_tokens"]); mt > 0 {
					payload["max_tokens"] = mt
					mutatedBody = true
				}
			}
			outcome.applied = append(outcome.applied, st.Key)

		case "jsonGuard":
			outcome.jsonGuard = true
			forceJSON, _ := st.Config["force"].(bool)
			_, hasRespFmt := payload["response_format"]
			if !hasRespFmt && (forceJSON || promptWantsJSON(payload)) {
				if claudeFmt {
					directives = append(directives, "JSON Guard is active: answer with one valid JSON value and nothing else — no prose, no markdown fence.")
				} else {
					payload["response_format"] = map[string]any{"type": "json_object"}
					directives = append(directives, "Output strictly valid JSON with no surrounding markdown fences or commentary.")
				}
				mutatedBody = true
			}
			outcome.applied = append(outcome.applied, st.Key)

		case "contextSqueezer":
			if rawMsgs, ok := payload["messages"].([]any); ok {
				ctxWin := intFromAny(st.Config["context_window"])
				if ctxWin <= 0 {
					ctxWin = defaultContextWindow
				}
				fillRatio := floatFromAny(st.Config["fill_ratio"])
				if fillRatio <= 0 || fillRatio > 1 {
					fillRatio = defaultFillRatio
				}
				minTurns := intFromAny(st.Config["min_turns"])
				if minTurns <= 0 {
					minTurns = defaultMinTurns
				}
				if squeezed, didSqueeze, _ := squeezeConversation(rawMsgs, squeezeOptions{
					ContextWindow: ctxWin,
					FillRatio:     fillRatio,
					MinTurns:      minTurns,
				}); didSqueeze {
					payload["messages"] = squeezed
					mutatedBody = true
				}
			}
			outcome.applied = append(outcome.applied, st.Key)

		case "openaiToolBridge":
			outcome.toolBridge = true
			if tools, ok := payload["tools"].([]any); ok && len(tools) > 0 {
				if _, exists := payload["tool_choice"]; !exists {
					payload["tool_choice"] = "auto"
					mutatedBody = true
				}
				if _, exists := payload["parallel_tool_calls"]; !exists {
					if ptc, ok := st.Config["parallel_tool_calls"].(bool); ok {
						payload["parallel_tool_calls"] = ptc
						mutatedBody = true
					}
				}
			}
			outcome.applied = append(outcome.applied, st.Key)

		case "antiSlop":
			lvl, _ := st.Config["level"].(string)
			directives = append(directives, antiSlopDirective(lvl))
			mutatedBody = true
			outcome.applied = append(outcome.applied, st.Key)
		}
	}

	if outcome.jsonGuard || outcome.toolBridge {
		outcome.toolCatalog = extractRequestTools(payload)
	}

	if len(outcome.applied) == 0 {
		return outcome
	}

	if len(directives) > 0 {
		injectSystemDirectives(payload, strings.Join(directives, "\n\n"))
		mutatedBody = true
	}

	if !mutatedBody {
		return outcome
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return outcome
	}
	outcome.mutated = encoded
	outcome.changed = true
	return outcome
}

func modelMatches(model string, allowlist []string) bool {
	if len(allowlist) == 0 {
		return true
	}
	normModel := strings.ToLower(strings.TrimSpace(model))
	bareModel := normModel
	if idx := strings.LastIndex(normModel, "/"); idx >= 0 {
		bareModel = normModel[idx+1:]
	}
	for _, entry := range allowlist {
		normEntry := strings.ToLower(strings.TrimSpace(entry))
		if normEntry == "" {
			continue
		}
		if normEntry == "*" || normEntry == normModel || normEntry == bareModel {
			return true
		}
		if strings.HasSuffix(normEntry, "*") {
			prefix := strings.TrimSuffix(normEntry, "*")
			if strings.HasPrefix(normModel, prefix) || strings.HasPrefix(bareModel, prefix) {
				return true
			}
		}
		if idx := strings.LastIndex(normEntry, "/"); idx >= 0 {
			bareEntry := normEntry[idx+1:]
			if bareEntry == bareModel {
				return true
			}
		}
	}
	return false
}

func injectSystemDirectives(payload map[string]any, directive string) {
	if sysStr, ok := payload["system"].(string); ok {
		if strings.TrimSpace(sysStr) == "" {
			payload["system"] = directive
		} else {
			payload["system"] = directive + "\n\n" + sysStr
		}
		return
	}
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
	if _, hasInput := payload["input"]; hasInput {
		if instr, ok := payload["instructions"].(string); ok && strings.TrimSpace(instr) != "" {
			payload["instructions"] = directive + "\n\n" + instr
		} else {
			payload["instructions"] = directive
		}
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

func floatFromAny(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		return 0
	}
}

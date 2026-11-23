package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

type memoryResponseWriter struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func newMemoryResponseWriter() *memoryResponseWriter {
	return &memoryResponseWriter{
		header: make(http.Header),
		status: http.StatusOK,
	}
}

func (w *memoryResponseWriter) Header() http.Header {
	return w.header
}

func (w *memoryResponseWriter) WriteHeader(code int) {
	w.status = code
}

func (w *memoryResponseWriter) Write(p []byte) (int, error) {
	return w.body.Write(p)
}

// CompleteChat executes an internal admin-scoped non-streaming OpenAI chat
// completion through the full gateway pipeline (interceptors, resolution,
// credentials/OAuth refresh, fallback, translation, and usage accounting).
func (m *Module) CompleteChat(ctx context.Context, req contracts.ChatCompletionRequest) (contracts.ChatCompletionResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	model := strings.TrimSpace(req.Model)
	if model == "" {
		return contracts.ChatCompletionResult{}, fmt.Errorf("%w: model is required", ErrInvalidPayload)
	}
	if len(req.Messages) == 0 {
		return contracts.ChatCompletionResult{}, fmt.Errorf("%w: at least one message is required", ErrInvalidPayload)
	}

	// Pre-check model resolution so sentinel errors (e.g. ErrModelNotFound,
	// ErrAllAccountsCoolingDown) propagate cleanly to callers.
	if _, err := m.resolver.Resolve(ctx, model, nil); err != nil {
		return contracts.ChatCompletionResult{}, err
	}

	payload := map[string]any{
		"model":    model,
		"messages": req.Messages,
		"stream":   false,
	}
	if req.Temperature != nil {
		payload["temperature"] = *req.Temperature
	}
	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return contracts.ChatCompletionResult{}, fmt.Errorf("marshal chat completion request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return contracts.ChatCompletionResult{}, fmt.Errorf("build internal chat request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	start := time.Now()
	rec := newMemoryResponseWriter()
	m.handleInferenceWithKey(rec, httpReq, contracts.KeyInfo{
		ID:   "admin-internal",
		Name: "admin-internal",
	}, ProtocolOpenAIChat, "", false, false)
	latencyMs := time.Since(start).Milliseconds()

	if rec.status < 200 || rec.status >= 300 {
		msg := extractGatewayErrorMessage(rec.body.Bytes())
		if msg == "" {
			msg = fmt.Sprintf("gateway returned HTTP %d", rec.status)
		}
		return contracts.ChatCompletionResult{Model: model, LatencyMs: latencyMs}, errors.New(msg)
	}

	var parsed struct {
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				Content any `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int64 `json:"prompt_tokens"`
			CompletionTokens int64 `json:"completion_tokens"`
			TotalTokens      int64 `json:"total_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(rec.body.Bytes(), &parsed); err != nil {
		return contracts.ChatCompletionResult{Model: model, LatencyMs: latencyMs}, fmt.Errorf("decode gateway chat response: %w", err)
	}

	content := ""
	if len(parsed.Choices) > 0 {
		content = stringifyMessageContent(parsed.Choices[0].Message.Content)
	}
	resolvedModel := strings.TrimSpace(parsed.Model)
	if resolvedModel == "" {
		resolvedModel = model
	}

	total := parsed.Usage.TotalTokens
	if total == 0 && (parsed.Usage.PromptTokens > 0 || parsed.Usage.CompletionTokens > 0) {
		total = parsed.Usage.PromptTokens + parsed.Usage.CompletionTokens
	}

	return contracts.ChatCompletionResult{
		Model:            resolvedModel,
		Content:          content,
		PromptTokens:     parsed.Usage.PromptTokens,
		CompletionTokens: parsed.Usage.CompletionTokens,
		TotalTokens:      total,
		LatencyMs:        latencyMs,
	}, nil
}

func extractGatewayErrorMessage(raw []byte) string {
	var envelope struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &envelope); err == nil && strings.TrimSpace(envelope.Error.Message) != "" {
		return strings.TrimSpace(envelope.Error.Message)
	}
	return strings.TrimSpace(string(raw))
}

func stringifyMessageContent(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []any:
		var parts []string
		for _, item := range val {
			if m, ok := item.(map[string]any); ok {
				if txt, ok := m["text"].(string); ok && txt != "" {
					parts = append(parts, txt)
				}
			}
		}
		return strings.Join(parts, "")
	default:
		return ""
	}
}

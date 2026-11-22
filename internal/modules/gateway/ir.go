// Package gateway implements the Canonical IR, direct fast-path model
// rewriting, multi-protocol request/response/SSE stream translators, and
// deterministic token estimation for the router gateway.
package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Supported gateway protocol identifiers matching providers and combos verbatim.
const (
	ProtocolOpenAIChat      = "openai-chat"
	ProtocolOpenAIResponses = "openai-responses"
	ProtocolAnthropic       = "anthropic"
	ProtocolGemini          = "gemini"
)

// Canonical content block types.
const (
	BlockTypeText       = "text"
	BlockTypeImage      = "image"
	BlockTypeThinking   = "thinking"
	BlockTypeToolUse    = "tool_use"
	BlockTypeToolResult = "tool_result"
)

// Canonical stop reasons.
const (
	StopReasonStop          = "stop"
	StopReasonLength        = "length"
	StopReasonToolCalls     = "tool_calls"
	StopReasonContentFilter = "content_filter"
)

// Canonical stream chunk lifecycle event types.
const (
	ChunkMessageStart = "message_start"
	ChunkBlockStart   = "block_start"
	ChunkBlockDelta   = "block_delta"
	ChunkBlockStop    = "block_stop"
	ChunkMessageDelta = "message_delta"
	ChunkMessageStop  = "message_stop"
)

var (
	// ErrUnsupportedProtocol indicates an unknown or unsupported wire protocol.
	ErrUnsupportedProtocol = errors.New("gateway: unsupported protocol")
	// ErrInvalidPayload indicates a malformed or invalid protocol payload.
	ErrInvalidPayload = errors.New("gateway: invalid payload")
)

// CacheControl represents prompt caching hints (e.g. Anthropic cache_control).
type CacheControl struct {
	Type string `json:"type"`
	TTL  string `json:"ttl,omitempty"`
}

// SystemBlock represents a single system instruction segment with optional cache control.
type SystemBlock struct {
	Text         string        `json:"text"`
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}

// ImageSource represents an inline base64 or URL-referenced image payload.
type ImageSource struct {
	URL       string `json:"url,omitempty"`
	MediaType string `json:"media_type,omitempty"`
	Data      string `json:"data,omitempty"`
	Detail    string `json:"detail,omitempty"`
}

// ContentBlock represents a typed multimodal, reasoning, or tool content block.
type ContentBlock struct {
	Type         string          `json:"type"`
	Text         string          `json:"text,omitempty"`
	Image        *ImageSource    `json:"image,omitempty"`
	Thinking     string          `json:"thinking,omitempty"`
	Signature    string          `json:"signature,omitempty"`
	ToolUseID    string          `json:"tool_use_id,omitempty"`
	ToolName     string          `json:"tool_name,omitempty"`
	ToolInput    json.RawMessage `json:"tool_input,omitempty"`
	ToolContent  []ContentBlock  `json:"tool_content,omitempty"`
	IsError      bool            `json:"is_error,omitempty"`
	CacheControl *CacheControl   `json:"cache_control,omitempty"`
}

// CanonicalMessage represents a single conversation turn in the Canonical IR.
type CanonicalMessage struct {
	Role    string         `json:"role"` // user, assistant, tool
	Name    string         `json:"name,omitempty"`
	Content []ContentBlock `json:"content"`
}

// ToolDefinition defines a callable function tool in the Canonical IR.
type ToolDefinition struct {
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Parameters   json.RawMessage `json:"parameters,omitempty"`
	Strict       *bool           `json:"strict,omitempty"`
	CacheControl *CacheControl   `json:"cache_control,omitempty"`
}

// ToolChoice constrains how the upstream model selects tools.
type ToolChoice struct {
	Type                   string `json:"type"` // auto, none, required, tool
	Name                   string `json:"name,omitempty"`
	DisableParallelToolUse *bool  `json:"disable_parallel_tool_use,omitempty"`
}

// ResponseFormat describes structured output constraints.
type ResponseFormat struct {
	Type   string          `json:"type"` // text, json_object, json_schema
	Name   string          `json:"name,omitempty"`
	Schema json.RawMessage `json:"schema,omitempty"`
	Strict *bool           `json:"strict,omitempty"`
}

// CanonicalChatRequest is the protocol-agnostic request intermediate representation.
type CanonicalChatRequest struct {
	Model              string             `json:"model"`
	System             []SystemBlock      `json:"system,omitempty"`
	Messages           []CanonicalMessage `json:"messages"`
	Tools              []ToolDefinition   `json:"tools,omitempty"`
	ToolChoice         *ToolChoice        `json:"tool_choice,omitempty"`
	ResponseFormat     *ResponseFormat    `json:"response_format,omitempty"`
	MaxTokens          *int               `json:"max_tokens,omitempty"`
	Temperature        *float64           `json:"temperature,omitempty"`
	TopP               *float64           `json:"top_p,omitempty"`
	StopSequences      []string           `json:"stop_sequences,omitempty"`
	Stream             bool               `json:"stream,omitempty"`
	ThinkingBudget     *int               `json:"thinking_budget,omitempty"`
	ReasoningEffort    string             `json:"reasoning_effort,omitempty"`
	PreviousResponseID string             `json:"previous_response_id,omitempty"`
	Metadata           map[string]any     `json:"metadata,omitempty"`
}

// CanonicalUsage captures normalized token usage across all upstream protocols.
type CanonicalUsage struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	CachedTokens int64 `json:"cached_tokens,omitempty"`
}

// CanonicalChatResponse is the protocol-agnostic non-streaming response IR.
type CanonicalChatResponse struct {
	ID         string         `json:"id"`
	Model      string         `json:"model"`
	Content    []ContentBlock `json:"content"`
	StopReason string         `json:"stop_reason"`
	Usage      CanonicalUsage `json:"usage"`
}

// CanonicalStreamChunk represents a normalized incremental streaming event.
type CanonicalStreamChunk struct {
	ID                 string          `json:"id,omitempty"`
	Model              string          `json:"model,omitempty"`
	EventType          string          `json:"event_type,omitempty"`
	BlockIndex         int             `json:"block_index,omitempty"`
	BlockType          string          `json:"block_type,omitempty"`
	TextDelta          string          `json:"text_delta,omitempty"`
	ThinkingDelta      string          `json:"thinking_delta,omitempty"`
	SignatureDelta     string          `json:"signature_delta,omitempty"`
	ToolUseID          string          `json:"tool_use_id,omitempty"`
	ToolName           string          `json:"tool_name,omitempty"`
	ToolInputJSONDelta string          `json:"tool_input_json_delta,omitempty"`
	StopReason         string          `json:"stop_reason,omitempty"`
	Usage              *CanonicalUsage `json:"usage,omitempty"`
	Done               bool            `json:"done,omitempty"`
}

// StreamTranslator converts upstream SSE events (outProto) into ready-to-write
// inbound SSE frames (inProto) while tracking normalized token usage.
type StreamTranslator interface {
	TranslateEvent(eventType string, data []byte) (frames [][]byte, usage *CanonicalUsage, done bool, err error)
	Finalize() (frames [][]byte, usage *CanonicalUsage)
}

type streamDecoder interface {
	Decode(eventType string, data []byte) ([]CanonicalStreamChunk, error)
}

type streamEncoder interface {
	Encode(chunk CanonicalStreamChunk) ([][]byte, error)
	Finalize(usage CanonicalUsage) [][]byte
}

// IsSupportedProtocol reports whether proto is one of the four supported protocols.
func IsSupportedProtocol(proto string) bool {
	switch proto {
	case ProtocolOpenAIChat, ProtocolOpenAIResponses, ProtocolAnthropic, ProtocolGemini:
		return true
	default:
		return false
	}
}

// CanFastPath returns true when inbound and outbound protocols match and are supported,
// enabling zero-loss top-level model rewriting without IR translation.
func CanFastPath(inProto, outProto string) bool {
	return IsSupportedProtocol(inProto) && inProto == outProto
}

// RewriteModelField rewrites only the top-level "model" field of a JSON request body
// while preserving all other provider-specific fields (such as Anthropic cache_control,
// thinking, betas, or OpenAI reasoning/response_format) byte-for-byte.
// For Gemini, where the model lives in the URL path, the validated body passes through unchanged.
func RewriteModelField(body []byte, proto, upstreamModel string) ([]byte, error) {
	if !IsSupportedProtocol(proto) {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProtocol, proto)
	}
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("%w: empty body", ErrInvalidPayload)
	}

	if proto == ProtocolGemini {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(trimmed, &obj); err != nil || obj == nil {
			return nil, fmt.Errorf("%w: expected JSON object", ErrInvalidPayload)
		}
		return bytes.Clone(body), nil
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &raw); err != nil || raw == nil {
		return nil, fmt.Errorf("%w: expected JSON object", ErrInvalidPayload)
	}

	if upstreamModel != "" {
		encodedModel, err := json.Marshal(upstreamModel)
		if err != nil {
			return nil, fmt.Errorf("gateway: marshal model: %w", err)
		}
		raw["model"] = encodedModel
	}

	out, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("gateway: marshal rewritten body: %w", err)
	}
	return out, nil
}

// DecodeRequest decodes a protocol-specific request body into the Canonical IR.
// If modelOverride is non-empty, it replaces the decoded Model field.
func DecodeRequest(proto string, body []byte, modelOverride string) (*CanonicalChatRequest, error) {
	if len(bytes.TrimSpace(body)) == 0 {
		return nil, fmt.Errorf("%w: empty request body", ErrInvalidPayload)
	}

	var (
		req *CanonicalChatRequest
		err error
	)
	switch proto {
	case ProtocolOpenAIChat:
		req, err = decodeOpenAIChatRequest(body)
	case ProtocolOpenAIResponses:
		req, err = decodeOpenAIResponsesRequest(body)
	case ProtocolAnthropic:
		req, err = decodeAnthropicRequest(body)
	case ProtocolGemini:
		req, err = decodeGeminiRequest(body)
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProtocol, proto)
	}
	if err != nil {
		return nil, err
	}
	if modelOverride != "" {
		req.Model = modelOverride
	}
	return req, nil
}

// EncodeRequest encodes a Canonical IR request into the target protocol wire format.
func EncodeRequest(proto string, req *CanonicalChatRequest) ([]byte, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: nil canonical request", ErrInvalidPayload)
	}
	switch proto {
	case ProtocolOpenAIChat:
		return encodeOpenAIChatRequest(req)
	case ProtocolOpenAIResponses:
		return encodeOpenAIResponsesRequest(req)
	case ProtocolAnthropic:
		return encodeAnthropicRequest(req)
	case ProtocolGemini:
		return encodeGeminiRequest(req)
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProtocol, proto)
	}
}

// DecodeResponse decodes a protocol-specific non-streaming response into the Canonical IR.
func DecodeResponse(proto string, body []byte) (*CanonicalChatResponse, error) {
	if len(bytes.TrimSpace(body)) == 0 {
		return nil, fmt.Errorf("%w: empty response body", ErrInvalidPayload)
	}
	switch proto {
	case ProtocolOpenAIChat:
		return decodeOpenAIChatResponse(body)
	case ProtocolOpenAIResponses:
		return decodeOpenAIResponsesResponse(body)
	case ProtocolAnthropic:
		return decodeAnthropicResponse(body)
	case ProtocolGemini:
		return decodeGeminiResponse(body)
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProtocol, proto)
	}
}

// EncodeResponse encodes a Canonical IR response into the target protocol wire format.
func EncodeResponse(proto string, resp *CanonicalChatResponse) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("%w: nil canonical response", ErrInvalidPayload)
	}
	switch proto {
	case ProtocolOpenAIChat:
		return encodeOpenAIChatResponse(resp)
	case ProtocolOpenAIResponses:
		return encodeOpenAIResponsesResponse(resp)
	case ProtocolAnthropic:
		return encodeAnthropicResponse(resp)
	case ProtocolGemini:
		return encodeGeminiResponse(resp)
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProtocol, proto)
	}
}

// NewStreamTranslator constructs a stateful SSE translator that decodes events from
// outProto and emits ready-to-write SSE frames in inProto.
func NewStreamTranslator(inProto, outProto, model string) (StreamTranslator, error) {
	if !IsSupportedProtocol(inProto) {
		return nil, fmt.Errorf("%w: inbound %q", ErrUnsupportedProtocol, inProto)
	}
	if !IsSupportedProtocol(outProto) {
		return nil, fmt.Errorf("%w: outbound %q", ErrUnsupportedProtocol, outProto)
	}

	var dec streamDecoder
	switch outProto {
	case ProtocolOpenAIChat:
		dec = newOpenAIStreamDecoder(model)
	case ProtocolOpenAIResponses:
		dec = newResponsesStreamDecoder(model)
	case ProtocolAnthropic:
		dec = newAnthropicStreamDecoder(model)
	case ProtocolGemini:
		dec = newGeminiStreamDecoder(model)
	}

	var enc streamEncoder
	switch inProto {
	case ProtocolOpenAIChat:
		enc = newOpenAIStreamEncoder(model)
	case ProtocolOpenAIResponses:
		enc = newResponsesStreamEncoder(model)
	case ProtocolAnthropic:
		enc = newAnthropicStreamEncoder(model)
	case ProtocolGemini:
		enc = newGeminiStreamEncoder(model)
	}

	return &composedStreamTranslator{
		decoder: dec,
		encoder: enc,
	}, nil
}

type composedStreamTranslator struct {
	decoder  streamDecoder
	encoder  streamEncoder
	usage    CanonicalUsage
	hasUsage bool
	done     bool
}

func (t *composedStreamTranslator) TranslateEvent(eventType string, data []byte) (frames [][]byte, usage *CanonicalUsage, done bool, err error) {
	chunks, err := t.decoder.Decode(eventType, data)
	if err != nil {
		return nil, nil, false, err
	}

	for _, chunk := range chunks {
		if chunk.Usage != nil {
			t.mergeUsage(*chunk.Usage)
		}
		if t.hasUsage && (chunk.EventType == ChunkMessageDelta || chunk.EventType == ChunkMessageStop) {
			merged := t.usage
			chunk.Usage = &merged
		}
		var encoded [][]byte
		encoded, err = t.encoder.Encode(chunk)
		if err != nil {
			return nil, nil, false, err
		}
		frames = append(frames, encoded...)
		if chunk.Done {
			t.done = true
		}
	}

	if t.hasUsage {
		u := t.usage
		usage = &u
	}
	return frames, usage, t.done, nil
}

func (t *composedStreamTranslator) Finalize() ([][]byte, *CanonicalUsage) {
	var usageCopy *CanonicalUsage
	if t.hasUsage {
		u := t.usage
		usageCopy = &u
	}
	if t.done {
		return nil, usageCopy
	}
	t.done = true
	return t.encoder.Finalize(t.usage), usageCopy
}

func (t *composedStreamTranslator) mergeUsage(u CanonicalUsage) {
	if u.InputTokens > 0 {
		t.usage.InputTokens = u.InputTokens
		t.hasUsage = true
	}
	if u.OutputTokens > 0 {
		t.usage.OutputTokens = u.OutputTokens
		t.hasUsage = true
	}
	if u.CachedTokens > 0 {
		t.usage.CachedTokens = u.CachedTokens
		t.hasUsage = true
	}
}

func formatSSEFrame(eventType string, payload []byte) []byte {
	var b bytes.Buffer
	if eventType != "" {
		b.WriteString("event: ")
		b.WriteString(eventType)
		b.WriteByte('\n')
	}
	b.WriteString("data: ")
	b.Write(payload)
	b.WriteString("\n\n")
	return b.Bytes()
}

func marshalSSEFrame(eventType string, v any) ([]byte, error) {
	payload, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("gateway: marshal SSE frame: %w", err)
	}
	return formatSSEFrame(eventType, payload), nil
}

func parseDataURI(raw string) (mediaType, data string, ok bool) {
	if !strings.HasPrefix(raw, "data:") {
		return "", "", false
	}
	rest := strings.TrimPrefix(raw, "data:")
	comma := strings.IndexByte(rest, ',')
	if comma < 0 {
		return "", "", false
	}
	meta := rest[:comma]
	payload := rest[comma+1:]
	if !strings.HasSuffix(meta, ";base64") {
		return "", "", false
	}
	media := strings.TrimSuffix(meta, ";base64")
	if media == "" {
		media = "application/octet-stream"
	}
	return media, payload, true
}

func formatDataURI(mediaType, data string) string {
	if mediaType == "" {
		mediaType = "image/png"
	}
	return "data:" + mediaType + ";base64," + data
}

func parseJSONObjectOrWrap(raw string) json.RawMessage {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return json.RawMessage(`{}`)
	}
	if json.Valid([]byte(trimmed)) {
		var buf bytes.Buffer
		if err := json.Compact(&buf, []byte(trimmed)); err == nil {
			return json.RawMessage(buf.Bytes())
		}
		return json.RawMessage(trimmed)
	}
	wrapped, err := json.Marshal(map[string]string{"raw": raw})
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(wrapped)
}

func ensureJSONObject(raw json.RawMessage) json.RawMessage {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return json.RawMessage(`{}`)
	}
	if trimmed[0] == '{' && json.Valid(trimmed) {
		return trimmed
	}
	if json.Valid(trimmed) {
		var s string
		if err := json.Unmarshal(trimmed, &s); err == nil {
			inner := bytes.TrimSpace([]byte(s))
			if len(inner) > 0 && inner[0] == '{' && json.Valid(inner) {
				return json.RawMessage(inner)
			}
		}
		wrapped, err := json.Marshal(map[string]json.RawMessage{"value": trimmed})
		if err == nil {
			return json.RawMessage(wrapped)
		}
	}
	return json.RawMessage(`{}`)
}

func rawJSONToString(raw json.RawMessage) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return "{}"
	}
	if json.Valid(trimmed) {
		var buf bytes.Buffer
		if err := json.Compact(&buf, trimmed); err == nil {
			return buf.String()
		}
	}
	return string(trimmed)
}

func joinSystemText(blocks []SystemBlock) string {
	parts := make([]string, 0, len(blocks))
	for _, b := range blocks {
		if b.Text != "" {
			parts = append(parts, b.Text)
		}
	}
	return strings.Join(parts, "\n\n")
}

func toolResultText(b ContentBlock) string {
	if b.Text != "" {
		return b.Text
	}
	if len(b.ToolContent) == 0 {
		return ""
	}
	parts := make([]string, 0, len(b.ToolContent))
	for _, child := range b.ToolContent {
		if child.Text != "" {
			parts = append(parts, child.Text)
		}
	}
	return strings.Join(parts, "\n")
}

func thinkingBudgetToEffort(budget int) string {
	switch {
	case budget <= 0:
		return ""
	case budget <= 2048:
		return "low"
	case budget <= 8192:
		return "medium"
	default:
		return "high"
	}
}

func effortToThinkingBudget(effort string) int {
	switch strings.ToLower(strings.TrimSpace(effort)) {
	case "low":
		return 2048
	case "medium":
		return 8192
	case "high", "xhigh":
		return 16384
	default:
		return 0
	}
}

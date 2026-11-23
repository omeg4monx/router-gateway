package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

const defaultAnthropicMaxTokens = 4096

type anthropicRequest struct {
	Model         string             `json:"model"`
	MaxTokens     int                `json:"max_tokens"`
	System        json.RawMessage    `json:"system,omitempty"`
	Messages      []anthropicMessage `json:"messages"`
	Tools         []anthropicTool    `json:"tools,omitempty"`
	ToolChoice    *anthropicChoice   `json:"tool_choice,omitempty"`
	Temperature   *float64           `json:"temperature,omitempty"`
	TopP          *float64           `json:"top_p,omitempty"`
	StopSequences []string           `json:"stop_sequences,omitempty"`
	Stream        bool               `json:"stream,omitempty"`
	Thinking      *anthropicThinking `json:"thinking,omitempty"`
	Metadata      map[string]any     `json:"metadata,omitempty"`
}

type anthropicThinking struct {
	Type         string `json:"type"`
	BudgetTokens int    `json:"budget_tokens,omitempty"`
}

type anthropicChoice struct {
	Type                   string `json:"type"`
	Name                   string `json:"name,omitempty"`
	DisableParallelToolUse *bool  `json:"disable_parallel_tool_use,omitempty"`
}

type anthropicTool struct {
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	InputSchema  json.RawMessage `json:"input_schema"`
	CacheControl *CacheControl   `json:"cache_control,omitempty"`
}

type anthropicMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type anthropicContentBlock struct {
	Type         string                `json:"type"`
	Text         string                `json:"text,omitempty"`
	Thinking     string                `json:"thinking,omitempty"`
	Signature    string                `json:"signature,omitempty"`
	ID           string                `json:"id,omitempty"`
	Name         string                `json:"name,omitempty"`
	Input        json.RawMessage       `json:"input,omitempty"`
	ToolUseID    string                `json:"tool_use_id,omitempty"`
	Content      json.RawMessage       `json:"content,omitempty"`
	IsError      bool                  `json:"is_error,omitempty"`
	Source       *anthropicImageSource `json:"source,omitempty"`
	CacheControl *CacheControl         `json:"cache_control,omitempty"`
}

type anthropicImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type,omitempty"`
	Data      string `json:"data,omitempty"`
	URL       string `json:"url,omitempty"`
}

type anthropicResponse struct {
	ID           string                  `json:"id"`
	Type         string                  `json:"type"`
	Role         string                  `json:"role"`
	Model        string                  `json:"model"`
	Content      []anthropicContentBlock `json:"content"`
	StopReason   string                  `json:"stop_reason"`
	StopSequence *string                 `json:"stop_sequence,omitempty"`
	Usage        anthropicUsage          `json:"usage"`
}

type anthropicUsage struct {
	InputTokens              int64 `json:"input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens,omitempty"`
}

func decodeAnthropicRequest(body []byte) (*CanonicalChatRequest, error) {
	var raw anthropicRequest
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("%w: decode anthropic request: %w", ErrInvalidPayload, err)
	}

	req := &CanonicalChatRequest{
		Model:         raw.Model,
		Temperature:   raw.Temperature,
		TopP:          raw.TopP,
		StopSequences: raw.StopSequences,
		Stream:        raw.Stream,
		Metadata:      raw.Metadata,
	}
	if raw.MaxTokens > 0 {
		mt := raw.MaxTokens
		req.MaxTokens = &mt
	}
	if raw.Thinking != nil && raw.Thinking.Type == "enabled" && raw.Thinking.BudgetTokens > 0 {
		budget := raw.Thinking.BudgetTokens
		req.ThinkingBudget = &budget
		req.ReasoningEffort = thinkingBudgetToEffort(budget)
	}

	if len(bytes.TrimSpace(raw.System)) > 0 && !bytes.Equal(bytes.TrimSpace(raw.System), []byte("null")) {
		trimmed := bytes.TrimSpace(raw.System)
		if trimmed[0] == '"' {
			var sysStr string
			if err := json.Unmarshal(trimmed, &sysStr); err != nil {
				return nil, fmt.Errorf("%w: invalid anthropic system string: %w", ErrInvalidPayload, err)
			}
			if sysStr != "" {
				req.System = append(req.System, SystemBlock{Text: sysStr})
			}
		} else {
			var sysBlocks []anthropicContentBlock
			if err := json.Unmarshal(trimmed, &sysBlocks); err != nil {
				return nil, fmt.Errorf("%w: invalid anthropic system blocks: %w", ErrInvalidPayload, err)
			}
			for _, sb := range sysBlocks {
				if sb.Text != "" {
					req.System = append(req.System, SystemBlock{
						Text:         sb.Text,
						CacheControl: sb.CacheControl,
					})
				}
			}
		}
	}

	for _, t := range raw.Tools {
		req.Tools = append(req.Tools, ToolDefinition{
			Name:         t.Name,
			Description:  t.Description,
			Parameters:   t.InputSchema,
			CacheControl: t.CacheControl,
		})
	}

	if raw.ToolChoice != nil {
		tc := &ToolChoice{
			Name:                   raw.ToolChoice.Name,
			DisableParallelToolUse: raw.ToolChoice.DisableParallelToolUse,
		}
		switch raw.ToolChoice.Type {
		case "any":
			tc.Type = "required"
		default:
			tc.Type = raw.ToolChoice.Type
		}
		req.ToolChoice = tc
	}

	for _, m := range raw.Messages {
		blocks, err := decodeAnthropicContent(m.Content)
		if err != nil {
			return nil, err
		}
		req.Messages = append(req.Messages, CanonicalMessage{
			Role:    m.Role,
			Content: blocks,
		})
	}

	return req, nil
}

func decodeAnthropicContent(raw json.RawMessage) ([]ContentBlock, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}
	if trimmed[0] == '"' {
		var text string
		if err := json.Unmarshal(trimmed, &text); err != nil {
			return nil, fmt.Errorf("%w: invalid anthropic string content: %w", ErrInvalidPayload, err)
		}
		return []ContentBlock{{Type: BlockTypeText, Text: text}}, nil
	}

	var blocks []anthropicContentBlock
	if err := json.Unmarshal(trimmed, &blocks); err != nil {
		return nil, fmt.Errorf("%w: invalid anthropic content array: %w", ErrInvalidPayload, err)
	}
	out := make([]ContentBlock, 0, len(blocks))
	for _, b := range blocks {
		out = append(out, decodeAnthropicBlock(b))
	}
	return out, nil
}

func decodeAnthropicBlock(b anthropicContentBlock) ContentBlock {
	switch b.Type {
	case "text":
		return ContentBlock{
			Type:         BlockTypeText,
			Text:         b.Text,
			CacheControl: b.CacheControl,
		}
	case "thinking":
		return ContentBlock{
			Type:         BlockTypeThinking,
			Thinking:     b.Thinking,
			Signature:    b.Signature,
			CacheControl: b.CacheControl,
		}
	case "image":
		img := &ImageSource{}
		if b.Source != nil {
			img.MediaType = b.Source.MediaType
			img.Data = b.Source.Data
			img.URL = b.Source.URL
		}
		return ContentBlock{
			Type:         BlockTypeImage,
			Image:        img,
			CacheControl: b.CacheControl,
		}
	case "tool_use":
		return ContentBlock{
			Type:         BlockTypeToolUse,
			ToolUseID:    b.ID,
			ToolName:     b.Name,
			ToolInput:    ensureJSONObject(b.Input),
			CacheControl: b.CacheControl,
		}
	case "tool_result":
		cb := ContentBlock{
			Type:         BlockTypeToolResult,
			ToolUseID:    b.ToolUseID,
			IsError:      b.IsError,
			CacheControl: b.CacheControl,
		}
		trimmed := bytes.TrimSpace(b.Content)
		if len(trimmed) > 0 && trimmed[0] == '"' {
			var s string
			if err := json.Unmarshal(trimmed, &s); err == nil {
				cb.Text = s
			}
		} else if len(trimmed) > 0 && trimmed[0] == '[' {
			if children, err := decodeAnthropicContent(trimmed); err == nil {
				cb.ToolContent = children
				var textParts []string
				for _, ch := range children {
					if ch.Type == BlockTypeText {
						textParts = append(textParts, ch.Text)
					}
				}
				cb.Text = strings.Join(textParts, "\n")
			}
		}
		return cb
	default:
		return ContentBlock{
			Type: BlockTypeText,
			Text: b.Text,
		}
	}
}

func encodeAnthropicRequest(req *CanonicalChatRequest) ([]byte, error) {
	maxTokens := defaultAnthropicMaxTokens
	if req.MaxTokens != nil && *req.MaxTokens > 0 {
		maxTokens = *req.MaxTokens
	}

	var thinking *anthropicThinking
	budget := 0
	if req.ThinkingBudget != nil && *req.ThinkingBudget > 0 {
		budget = *req.ThinkingBudget
	} else if req.ReasoningEffort != "" {
		budget = effortToThinkingBudget(req.ReasoningEffort)
	}
	if budget > 0 {
		thinking = &anthropicThinking{
			Type:         "enabled",
			BudgetTokens: budget,
		}
		if maxTokens <= budget {
			maxTokens = budget + 1024
		}
	}

	out := anthropicRequest{
		Model:         req.Model,
		MaxTokens:     maxTokens,
		Temperature:   req.Temperature,
		TopP:          req.TopP,
		StopSequences: req.StopSequences,
		Stream:        req.Stream,
		Thinking:      thinking,
		Metadata:      req.Metadata,
	}

	if len(req.System) > 0 {
		hasCache := false
		for _, sb := range req.System {
			if sb.CacheControl != nil {
				hasCache = true
				break
			}
		}
		if len(req.System) == 1 && !hasCache {
			rawSys, err := json.Marshal(req.System[0].Text)
			if err != nil {
				return nil, fmt.Errorf("gateway: marshal anthropic system: %w", err)
			}
			out.System = rawSys
		} else {
			sysBlocks := make([]anthropicContentBlock, 0, len(req.System))
			for _, sb := range req.System {
				sysBlocks = append(sysBlocks, anthropicContentBlock{
					Type:         "text",
					Text:         sb.Text,
					CacheControl: sb.CacheControl,
				})
			}
			rawSys, err := json.Marshal(sysBlocks)
			if err != nil {
				return nil, fmt.Errorf("gateway: marshal anthropic system blocks: %w", err)
			}
			out.System = rawSys
		}
	}

	for _, t := range req.Tools {
		out.Tools = append(out.Tools, anthropicTool{
			Name:         t.Name,
			Description:  t.Description,
			InputSchema:  ensureJSONObject(t.Parameters),
			CacheControl: t.CacheControl,
		})
	}

	if req.ToolChoice != nil {
		tc := &anthropicChoice{
			Name:                   req.ToolChoice.Name,
			DisableParallelToolUse: req.ToolChoice.DisableParallelToolUse,
		}
		switch req.ToolChoice.Type {
		case "required":
			tc.Type = "any"
		default:
			tc.Type = req.ToolChoice.Type
		}
		out.ToolChoice = tc
	}

	// Merge consecutive tool/user messages into valid alternating user/assistant messages.
	type pendingMsg struct {
		role   string
		blocks []anthropicContentBlock
	}
	var merged []pendingMsg
	for _, m := range req.Messages {
		role := m.Role
		if role == "tool" || role == "" {
			role = "user"
		}
		encodedBlocks := encodeAnthropicBlocks(m.Content)
		if len(encodedBlocks) == 0 {
			continue
		}
		if len(merged) > 0 && merged[len(merged)-1].role == role {
			merged[len(merged)-1].blocks = append(merged[len(merged)-1].blocks, encodedBlocks...)
		} else {
			merged = append(merged, pendingMsg{role: role, blocks: encodedBlocks})
		}
	}

	for _, m := range merged {
		rawBlocks, err := json.Marshal(m.blocks)
		if err != nil {
			return nil, fmt.Errorf("gateway: marshal anthropic message blocks: %w", err)
		}
		out.Messages = append(out.Messages, anthropicMessage{
			Role:    m.role,
			Content: rawBlocks,
		})
	}

	return json.Marshal(out)
}

func encodeAnthropicBlocks(blocks []ContentBlock) []anthropicContentBlock {
	out := make([]anthropicContentBlock, 0, len(blocks))
	for _, b := range blocks {
		switch b.Type {
		case BlockTypeText:
			out = append(out, anthropicContentBlock{
				Type:         "text",
				Text:         b.Text,
				CacheControl: b.CacheControl,
			})
		case BlockTypeThinking:
			out = append(out, anthropicContentBlock{
				Type:         "thinking",
				Thinking:     b.Thinking,
				Signature:    b.Signature,
				CacheControl: b.CacheControl,
			})
		case BlockTypeImage:
			if b.Image == nil {
				continue
			}
			var src *anthropicImageSource
			if b.Image.Data != "" {
				mt := b.Image.MediaType
				if mt == "" {
					mt = "image/png"
				}
				src = &anthropicImageSource{
					Type:      "base64",
					MediaType: mt,
					Data:      b.Image.Data,
				}
			} else if b.Image.URL != "" {
				if mt, data, ok := parseDataURI(b.Image.URL); ok {
					src = &anthropicImageSource{
						Type:      "base64",
						MediaType: mt,
						Data:      data,
					}
				} else {
					src = &anthropicImageSource{
						Type: "url",
						URL:  b.Image.URL,
					}
				}
			}
			if src != nil {
				out = append(out, anthropicContentBlock{
					Type:         "image",
					Source:       src,
					CacheControl: b.CacheControl,
				})
			}
		case BlockTypeToolUse:
			out = append(out, anthropicContentBlock{
				Type:         "tool_use",
				ID:           b.ToolUseID,
				Name:         b.ToolName,
				Input:        ensureJSONObject(b.ToolInput),
				CacheControl: b.CacheControl,
			})
		case BlockTypeToolResult:
			var rawContent json.RawMessage
			if len(b.ToolContent) > 0 {
				if encoded, err := json.Marshal(encodeAnthropicBlocks(b.ToolContent)); err == nil {
					rawContent = encoded
				}
			} else {
				if encoded, err := json.Marshal(b.Text); err == nil {
					rawContent = encoded
				}
			}
			out = append(out, anthropicContentBlock{
				Type:         "tool_result",
				ToolUseID:    b.ToolUseID,
				Content:      rawContent,
				IsError:      b.IsError,
				CacheControl: b.CacheControl,
			})
		}
	}
	return out
}

func decodeAnthropicResponse(body []byte) (*CanonicalChatResponse, error) {
	var raw anthropicResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("%w: decode anthropic response: %w", ErrInvalidPayload, err)
	}

	resp := &CanonicalChatResponse{
		ID:         raw.ID,
		Model:      raw.Model,
		StopReason: mapAnthropicStopReasonToCanonical(raw.StopReason),
		Usage: CanonicalUsage{
			InputTokens:  raw.Usage.InputTokens,
			OutputTokens: raw.Usage.OutputTokens,
			CachedTokens: raw.Usage.CacheReadInputTokens,
		},
	}
	for _, b := range raw.Content {
		resp.Content = append(resp.Content, decodeAnthropicBlock(b))
	}
	return resp, nil
}

func encodeAnthropicResponse(resp *CanonicalChatResponse) ([]byte, error) {
	out := anthropicResponse{
		ID:         resp.ID,
		Type:       "message",
		Role:       "assistant",
		Model:      resp.Model,
		Content:    encodeAnthropicBlocks(resp.Content),
		StopReason: mapCanonicalStopReasonToAnthropic(resp.StopReason),
		Usage: anthropicUsage{
			InputTokens:          resp.Usage.InputTokens,
			OutputTokens:         resp.Usage.OutputTokens,
			CacheReadInputTokens: resp.Usage.CachedTokens,
		},
	}
	return json.Marshal(out)
}

func mapAnthropicStopReasonToCanonical(reason string) string {
	switch reason {
	case "max_tokens":
		return StopReasonLength
	case "tool_use":
		return StopReasonToolCalls
	case "refusal":
		return StopReasonContentFilter
	default:
		return StopReasonStop
	}
}

func mapCanonicalStopReasonToAnthropic(reason string) string {
	switch reason {
	case StopReasonLength:
		return "max_tokens"
	case StopReasonToolCalls:
		return "tool_use"
	case StopReasonContentFilter:
		return "refusal"
	default:
		return "end_turn"
	}
}

// --- Anthropic Streaming Decoder & Encoder ---

type anthropicStreamDecoder struct {
	model      string
	blockTypes map[int]string
}

func newAnthropicStreamDecoder(model string) *anthropicStreamDecoder {
	return &anthropicStreamDecoder{
		model:      model,
		blockTypes: make(map[int]string),
	}
}

func (d *anthropicStreamDecoder) Decode(eventType string, data []byte) ([]CanonicalStreamChunk, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, nil
	}

	var envelope struct {
		Type    string `json:"type"`
		Index   int    `json:"index"`
		Message *struct {
			ID    string         `json:"id"`
			Model string         `json:"model"`
			Usage anthropicUsage `json:"usage"`
		} `json:"message,omitempty"`
		ContentBlock *struct {
			Type string `json:"type"`
			ID   string `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"content_block,omitempty"`
		Delta *struct {
			Type        string `json:"type,omitempty"`
			Text        string `json:"text,omitempty"`
			Thinking    string `json:"thinking,omitempty"`
			Signature   string `json:"signature,omitempty"`
			PartialJSON string `json:"partial_json,omitempty"`
			StopReason  string `json:"stop_reason,omitempty"`
		} `json:"delta,omitempty"`
		Usage *anthropicUsage `json:"usage,omitempty"`
	}
	if err := json.Unmarshal(trimmed, &envelope); err != nil {
		return nil, fmt.Errorf("%w: decode anthropic SSE event: %w", ErrInvalidPayload, err)
	}

	ev := envelope.Type
	if ev == "" {
		ev = eventType
	}

	switch ev {
	case "ping":
		return nil, nil

	case "message_start":
		chunk := CanonicalStreamChunk{
			Model:     d.model,
			EventType: ChunkMessageStart,
		}
		if envelope.Message != nil {
			chunk.ID = envelope.Message.ID
			if envelope.Message.Model != "" {
				chunk.Model = envelope.Message.Model
			}
			chunk.Usage = &CanonicalUsage{
				InputTokens:  envelope.Message.Usage.InputTokens,
				OutputTokens: envelope.Message.Usage.OutputTokens,
				CachedTokens: envelope.Message.Usage.CacheReadInputTokens,
			}
		}
		return []CanonicalStreamChunk{chunk}, nil

	case "content_block_start":
		if envelope.ContentBlock == nil {
			return nil, nil
		}
		bt := envelope.ContentBlock.Type
		d.blockTypes[envelope.Index] = bt
		return []CanonicalStreamChunk{{
			EventType:  ChunkBlockStart,
			BlockIndex: envelope.Index,
			BlockType:  bt,
			ToolUseID:  envelope.ContentBlock.ID,
			ToolName:   envelope.ContentBlock.Name,
		}}, nil

	case "content_block_delta":
		if envelope.Delta == nil {
			return nil, nil
		}
		bt := d.blockTypes[envelope.Index]
		chunk := CanonicalStreamChunk{
			EventType:  ChunkBlockDelta,
			BlockIndex: envelope.Index,
			BlockType:  bt,
		}
		switch envelope.Delta.Type {
		case "text_delta":
			chunk.BlockType = BlockTypeText
			chunk.TextDelta = envelope.Delta.Text
		case "thinking_delta":
			chunk.BlockType = BlockTypeThinking
			chunk.ThinkingDelta = envelope.Delta.Thinking
		case "signature_delta":
			chunk.BlockType = BlockTypeThinking
			chunk.SignatureDelta = envelope.Delta.Signature
		case "input_json_delta":
			chunk.BlockType = BlockTypeToolUse
			chunk.ToolInputJSONDelta = envelope.Delta.PartialJSON
		}
		return []CanonicalStreamChunk{chunk}, nil

	case "content_block_stop":
		bt := d.blockTypes[envelope.Index]
		return []CanonicalStreamChunk{{
			EventType:  ChunkBlockStop,
			BlockIndex: envelope.Index,
			BlockType:  bt,
		}}, nil

	case "message_delta":
		chunk := CanonicalStreamChunk{
			EventType: ChunkMessageDelta,
		}
		if envelope.Delta != nil && envelope.Delta.StopReason != "" {
			chunk.StopReason = mapAnthropicStopReasonToCanonical(envelope.Delta.StopReason)
		}
		if envelope.Usage != nil {
			chunk.Usage = &CanonicalUsage{
				InputTokens:  envelope.Usage.InputTokens,
				OutputTokens: envelope.Usage.OutputTokens,
				CachedTokens: envelope.Usage.CacheReadInputTokens,
			}
		}
		return []CanonicalStreamChunk{chunk}, nil

	case "message_stop":
		return []CanonicalStreamChunk{{
			EventType: ChunkMessageStop,
			Done:      true,
		}}, nil

	default:
		return nil, nil
	}
}

type anthropicStreamEncoder struct {
	id          string
	model       string
	started     bool
	openBlocks  map[int]string
	stopReason  string
	usage       CanonicalUsage
	sawToolUse  bool
	sentStopped bool
}

func newAnthropicStreamEncoder(model string) *anthropicStreamEncoder {
	return &anthropicStreamEncoder{
		id:         "msg_gateway",
		model:      model,
		openBlocks: make(map[int]string),
		stopReason: StopReasonStop,
	}
}

func (e *anthropicStreamEncoder) ensureStarted(frames *[][]byte) error {
	if e.started {
		return nil
	}
	e.started = true
	frame, err := marshalSSEFrame("message_start", map[string]any{
		"type": "message_start",
		"message": map[string]any{
			"id":            e.id,
			"type":          "message",
			"role":          "assistant",
			"model":         e.model,
			"content":       []any{},
			"stop_reason":   nil,
			"stop_sequence": nil,
			"usage": anthropicUsage{
				InputTokens:          e.usage.InputTokens,
				OutputTokens:         0,
				CacheReadInputTokens: e.usage.CachedTokens,
			},
		},
	})
	if err != nil {
		return err
	}
	*frames = append(*frames, frame)
	return nil
}

func (e *anthropicStreamEncoder) Encode(chunk CanonicalStreamChunk) ([][]byte, error) {
	if chunk.ID != "" {
		e.id = chunk.ID
	}
	if chunk.Model != "" {
		e.model = chunk.Model
	}
	if chunk.Usage != nil {
		if chunk.Usage.InputTokens > 0 {
			e.usage.InputTokens = chunk.Usage.InputTokens
		}
		if chunk.Usage.OutputTokens > 0 {
			e.usage.OutputTokens = chunk.Usage.OutputTokens
		}
		if chunk.Usage.CachedTokens > 0 {
			e.usage.CachedTokens = chunk.Usage.CachedTokens
		}
	}

	var frames [][]byte
	switch chunk.EventType {
	case ChunkMessageStart:
		if err := e.ensureStarted(&frames); err != nil {
			return nil, err
		}
		return frames, nil

	case ChunkBlockStart:
		if err := e.ensureStarted(&frames); err != nil {
			return nil, err
		}
		var cb map[string]any
		switch chunk.BlockType {
		case BlockTypeThinking:
			cb = map[string]any{"type": "thinking", "thinking": ""}
		case BlockTypeToolUse:
			e.sawToolUse = true
			cb = map[string]any{
				"type":  "tool_use",
				"id":    chunk.ToolUseID,
				"name":  chunk.ToolName,
				"input": map[string]any{},
			}
		default:
			cb = map[string]any{"type": "text", "text": ""}
		}
		e.openBlocks[chunk.BlockIndex] = chunk.BlockType
		frame, err := marshalSSEFrame("content_block_start", map[string]any{
			"type":          "content_block_start",
			"index":         chunk.BlockIndex,
			"content_block": cb,
		})
		if err != nil {
			return nil, err
		}
		frames = append(frames, frame)
		return frames, nil

	case ChunkBlockDelta:
		if err := e.ensureStarted(&frames); err != nil {
			return nil, err
		}
		if _, open := e.openBlocks[chunk.BlockIndex]; !open {
			startFrames, err := e.Encode(CanonicalStreamChunk{
				EventType:  ChunkBlockStart,
				BlockIndex: chunk.BlockIndex,
				BlockType:  chunk.BlockType,
				ToolUseID:  chunk.ToolUseID,
				ToolName:   chunk.ToolName,
			})
			if err != nil {
				return nil, err
			}
			frames = append(frames, startFrames...)
		}
		var delta map[string]any
		switch {
		case chunk.ThinkingDelta != "":
			delta = map[string]any{"type": "thinking_delta", "thinking": chunk.ThinkingDelta}
		case chunk.SignatureDelta != "":
			delta = map[string]any{"type": "signature_delta", "signature": chunk.SignatureDelta}
		case chunk.TextDelta != "":
			delta = map[string]any{"type": "text_delta", "text": chunk.TextDelta}
		case chunk.ToolInputJSONDelta != "":
			delta = map[string]any{"type": "input_json_delta", "partial_json": chunk.ToolInputJSONDelta}
		default:
			return frames, nil
		}
		frame, err := marshalSSEFrame("content_block_delta", map[string]any{
			"type":  "content_block_delta",
			"index": chunk.BlockIndex,
			"delta": delta,
		})
		if err != nil {
			return nil, err
		}
		frames = append(frames, frame)
		return frames, nil

	case ChunkBlockStop:
		if _, open := e.openBlocks[chunk.BlockIndex]; !open {
			return nil, nil
		}
		delete(e.openBlocks, chunk.BlockIndex)
		frame, err := marshalSSEFrame("content_block_stop", map[string]any{
			"type":  "content_block_stop",
			"index": chunk.BlockIndex,
		})
		if err != nil {
			return nil, err
		}
		return [][]byte{frame}, nil

	case ChunkMessageDelta:
		if chunk.StopReason != "" {
			e.stopReason = chunk.StopReason
		}
		return nil, nil

	case ChunkMessageStop:
		return e.emitStop(e.usage)
	}
	return nil, nil
}

func (e *anthropicStreamEncoder) emitStop(usage CanonicalUsage) ([][]byte, error) {
	if e.sentStopped {
		return nil, nil
	}
	e.sentStopped = true
	var frames [][]byte
	if err := e.ensureStarted(&frames); err != nil {
		return nil, err
	}
	openIndices := make([]int, 0, len(e.openBlocks))
	for idx := range e.openBlocks {
		openIndices = append(openIndices, idx)
	}
	slices.Sort(openIndices)
	for _, idx := range openIndices {
		if frame, err := marshalSSEFrame("content_block_stop", map[string]any{
			"type":  "content_block_stop",
			"index": idx,
		}); err == nil {
			frames = append(frames, frame)
		}
		delete(e.openBlocks, idx)
	}
	stopReason := e.stopReason
	if stopReason == StopReasonStop && e.sawToolUse {
		stopReason = StopReasonToolCalls
	}
	deltaFrame, err := marshalSSEFrame("message_delta", map[string]any{
		"type": "message_delta",
		"delta": map[string]any{
			"stop_reason":   mapCanonicalStopReasonToAnthropic(stopReason),
			"stop_sequence": nil,
		},
		"usage": anthropicUsage{
			InputTokens:          usage.InputTokens,
			OutputTokens:         usage.OutputTokens,
			CacheReadInputTokens: usage.CachedTokens,
		},
	})
	if err != nil {
		return nil, err
	}
	stopFrame, err := marshalSSEFrame("message_stop", map[string]any{
		"type": "message_stop",
	})
	if err != nil {
		return nil, err
	}
	frames = append(frames, deltaFrame, stopFrame)
	return frames, nil
}

func (e *anthropicStreamEncoder) Finalize(usage CanonicalUsage) [][]byte {
	frames, _ := e.emitStop(usage)
	return frames
}

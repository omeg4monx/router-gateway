package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type openAIChatRequest struct {
	Model               string                   `json:"model"`
	Messages            []openAIChatMessage      `json:"messages"`
	Tools               []openAIChatTool         `json:"tools,omitempty"`
	ToolChoice          json.RawMessage          `json:"tool_choice,omitempty"`
	ParallelToolCalls   *bool                    `json:"parallel_tool_calls,omitempty"`
	ResponseFormat      *openAIResponseFormat    `json:"response_format,omitempty"`
	MaxTokens           *int                     `json:"max_tokens,omitempty"`
	MaxCompletionTokens *int                     `json:"max_completion_tokens,omitempty"`
	Temperature         *float64                 `json:"temperature,omitempty"`
	TopP                *float64                 `json:"top_p,omitempty"`
	Stop                json.RawMessage          `json:"stop,omitempty"`
	Stream              bool                     `json:"stream,omitempty"`
	ReasoningEffort     string                   `json:"reasoning_effort,omitempty"`
	Metadata            map[string]any           `json:"metadata,omitempty"`
	StreamOptions       *openAIStreamOptionsHint `json:"stream_options,omitempty"`
}

type openAIStreamOptionsHint struct {
	IncludeUsage bool `json:"include_usage,omitempty"`
}

type openAIChatMessage struct {
	Role             string               `json:"role"`
	Name             string               `json:"name,omitempty"`
	Content          json.RawMessage      `json:"content,omitempty"`
	ReasoningContent string               `json:"reasoning_content,omitempty"`
	ToolCalls        []openAIChatToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string               `json:"tool_call_id,omitempty"`
}

type openAIChatContentPart struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	ImageURL *openAIImageURL `json:"image_url,omitempty"`
}

type openAIImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

type openAIChatToolCall struct {
	Index    *int                   `json:"index,omitempty"`
	ID       string                 `json:"id,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Function openAIChatFunctionCall `json:"function"`
}

type openAIChatFunctionCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments"`
}

type openAIChatTool struct {
	Type     string                 `json:"type"`
	Function openAIChatFunctionSpec `json:"function"`
}

type openAIChatFunctionSpec struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
	Strict      *bool           `json:"strict,omitempty"`
}

type openAIResponseFormat struct {
	Type       string                 `json:"type"`
	JSONSchema *openAIJSONSchemaField `json:"json_schema,omitempty"`
}

type openAIJSONSchemaField struct {
	Name   string          `json:"name"`
	Schema json.RawMessage `json:"schema,omitempty"`
	Strict *bool           `json:"strict,omitempty"`
}

type openAIChatResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created,omitempty"`
	Model   string                 `json:"model"`
	Choices []openAIChatChoice     `json:"choices"`
	Usage   *openAIChatUsage       `json:"usage,omitempty"`
	Error   map[string]interface{} `json:"error,omitempty"`
}

type openAIChatChoice struct {
	Index        int                  `json:"index"`
	Message      openAIChatOutMessage `json:"message"`
	FinishReason string               `json:"finish_reason"`
}

type openAIChatOutMessage struct {
	Role             string               `json:"role"`
	Content          *string              `json:"content"`
	ReasoningContent string               `json:"reasoning_content,omitempty"`
	ToolCalls        []openAIChatToolCall `json:"tool_calls,omitempty"`
}

type openAIChatUsage struct {
	PromptTokens        int64                    `json:"prompt_tokens"`
	CompletionTokens    int64                    `json:"completion_tokens"`
	TotalTokens         int64                    `json:"total_tokens"`
	PromptTokensDetails *openAIPromptTokenDetail `json:"prompt_tokens_details,omitempty"`
}

type openAIPromptTokenDetail struct {
	CachedTokens int64 `json:"cached_tokens,omitempty"`
}

func decodeOpenAIChatRequest(body []byte) (*CanonicalChatRequest, error) {
	var raw openAIChatRequest
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("%w: decode openai-chat request: %w", ErrInvalidPayload, err)
	}

	req := &CanonicalChatRequest{
		Model:           raw.Model,
		Temperature:     raw.Temperature,
		TopP:            raw.TopP,
		Stream:          raw.Stream,
		ReasoningEffort: raw.ReasoningEffort,
		Metadata:        raw.Metadata,
	}
	if raw.MaxCompletionTokens != nil {
		req.MaxTokens = raw.MaxCompletionTokens
	} else if raw.MaxTokens != nil {
		req.MaxTokens = raw.MaxTokens
	}

	if len(bytes.TrimSpace(raw.Stop)) > 0 && !bytes.Equal(bytes.TrimSpace(raw.Stop), []byte("null")) {
		var single string
		if err := json.Unmarshal(raw.Stop, &single); err == nil {
			if single != "" {
				req.StopSequences = []string{single}
			}
		} else {
			var multi []string
			if err := json.Unmarshal(raw.Stop, &multi); err != nil {
				return nil, fmt.Errorf("%w: invalid openai stop field: %w", ErrInvalidPayload, err)
			}
			req.StopSequences = multi
		}
	}

	if raw.ResponseFormat != nil {
		rf := &ResponseFormat{Type: raw.ResponseFormat.Type}
		if raw.ResponseFormat.JSONSchema != nil {
			rf.Name = raw.ResponseFormat.JSONSchema.Name
			rf.Schema = raw.ResponseFormat.JSONSchema.Schema
			rf.Strict = raw.ResponseFormat.JSONSchema.Strict
		}
		req.ResponseFormat = rf
	}

	for _, t := range raw.Tools {
		if t.Function.Name == "" {
			continue
		}
		req.Tools = append(req.Tools, ToolDefinition{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			Parameters:  t.Function.Parameters,
			Strict:      t.Function.Strict,
		})
	}

	if len(bytes.TrimSpace(raw.ToolChoice)) > 0 && !bytes.Equal(bytes.TrimSpace(raw.ToolChoice), []byte("null")) {
		var tcStr string
		if err := json.Unmarshal(raw.ToolChoice, &tcStr); err == nil {
			req.ToolChoice = &ToolChoice{Type: tcStr}
		} else {
			var tcObj struct {
				Type     string `json:"type"`
				Function struct {
					Name string `json:"name"`
				} `json:"function"`
			}
			if err := json.Unmarshal(raw.ToolChoice, &tcObj); err != nil {
				return nil, fmt.Errorf("%w: invalid openai tool_choice: %w", ErrInvalidPayload, err)
			}
			if tcObj.Function.Name != "" {
				req.ToolChoice = &ToolChoice{Type: "tool", Name: tcObj.Function.Name}
			} else if tcObj.Type != "" {
				req.ToolChoice = &ToolChoice{Type: tcObj.Type}
			}
		}
	}
	if raw.ParallelToolCalls != nil {
		disable := !*raw.ParallelToolCalls
		if req.ToolChoice == nil {
			req.ToolChoice = &ToolChoice{Type: "auto"}
		}
		req.ToolChoice.DisableParallelToolUse = &disable
	}

	for _, m := range raw.Messages {
		blocks, err := decodeOpenAIContent(m.Content)
		if err != nil {
			return nil, err
		}

		switch m.Role {
		case "system", "developer":
			for _, b := range blocks {
				if b.Type == BlockTypeText && b.Text != "" {
					req.System = append(req.System, SystemBlock{Text: b.Text})
				}
			}
		case "tool":
			var textParts []string
			for _, b := range blocks {
				if b.Type == BlockTypeText {
					textParts = append(textParts, b.Text)
				}
			}
			req.Messages = append(req.Messages, CanonicalMessage{
				Role: "tool",
				Name: m.Name,
				Content: []ContentBlock{{
					Type:      BlockTypeToolResult,
					ToolUseID: m.ToolCallID,
					ToolName:  m.Name,
					Text:      strings.Join(textParts, "\n"),
				}},
			})
		case "assistant":
			var content []ContentBlock
			if m.ReasoningContent != "" {
				content = append(content, ContentBlock{
					Type:     BlockTypeThinking,
					Thinking: m.ReasoningContent,
				})
			}
			content = append(content, blocks...)
			for _, tc := range m.ToolCalls {
				content = append(content, ContentBlock{
					Type:      BlockTypeToolUse,
					ToolUseID: tc.ID,
					ToolName:  tc.Function.Name,
					ToolInput: parseJSONObjectOrWrap(tc.Function.Arguments),
				})
			}
			req.Messages = append(req.Messages, CanonicalMessage{
				Role:    "assistant",
				Name:    m.Name,
				Content: content,
			})
		default:
			req.Messages = append(req.Messages, CanonicalMessage{
				Role:    m.Role,
				Name:    m.Name,
				Content: blocks,
			})
		}
	}

	return req, nil
}

func decodeOpenAIContent(raw json.RawMessage) ([]ContentBlock, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}
	if trimmed[0] == '"' {
		var text string
		if err := json.Unmarshal(trimmed, &text); err != nil {
			return nil, fmt.Errorf("%w: invalid openai string content: %w", ErrInvalidPayload, err)
		}
		return []ContentBlock{{Type: BlockTypeText, Text: text}}, nil
	}

	var parts []openAIChatContentPart
	if err := json.Unmarshal(trimmed, &parts); err != nil {
		return nil, fmt.Errorf("%w: invalid openai content array: %w", ErrInvalidPayload, err)
	}
	blocks := make([]ContentBlock, 0, len(parts))
	for _, p := range parts {
		switch p.Type {
		case "text", "input_text", "output_text":
			blocks = append(blocks, ContentBlock{Type: BlockTypeText, Text: p.Text})
		case "image_url":
			if p.ImageURL == nil {
				continue
			}
			img := &ImageSource{
				URL:    p.ImageURL.URL,
				Detail: p.ImageURL.Detail,
			}
			if media, data, ok := parseDataURI(p.ImageURL.URL); ok {
				img.MediaType = media
				img.Data = data
			}
			blocks = append(blocks, ContentBlock{Type: BlockTypeImage, Image: img})
		}
	}
	return blocks, nil
}

func encodeOpenAIChatRequest(req *CanonicalChatRequest) ([]byte, error) {
	out := openAIChatRequest{
		Model:           req.Model,
		MaxTokens:       req.MaxTokens,
		Temperature:     req.Temperature,
		TopP:            req.TopP,
		Stream:          req.Stream,
		ReasoningEffort: req.ReasoningEffort,
		Metadata:        req.Metadata,
	}
	if out.ReasoningEffort == "" && req.ThinkingBudget != nil && *req.ThinkingBudget > 0 {
		out.ReasoningEffort = thinkingBudgetToEffort(*req.ThinkingBudget)
	}
	if len(req.StopSequences) > 0 {
		stopJSON, err := json.Marshal(req.StopSequences)
		if err != nil {
			return nil, fmt.Errorf("gateway: marshal openai stop: %w", err)
		}
		out.Stop = stopJSON
	}
	if req.ResponseFormat != nil {
		rf := &openAIResponseFormat{Type: req.ResponseFormat.Type}
		if req.ResponseFormat.Type == "json_schema" {
			rf.JSONSchema = &openAIJSONSchemaField{
				Name:   req.ResponseFormat.Name,
				Schema: req.ResponseFormat.Schema,
				Strict: req.ResponseFormat.Strict,
			}
		}
		out.ResponseFormat = rf
	}
	for _, t := range req.Tools {
		out.Tools = append(out.Tools, openAIChatTool{
			Type: "function",
			Function: openAIChatFunctionSpec{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  ensureJSONObject(t.Parameters),
				Strict:      t.Strict,
			},
		})
	}
	if req.ToolChoice != nil {
		switch req.ToolChoice.Type {
		case "auto", "none", "required":
			rawTC, err := json.Marshal(req.ToolChoice.Type)
			if err != nil {
				return nil, fmt.Errorf("gateway: marshal openai tool_choice: %w", err)
			}
			out.ToolChoice = rawTC
		case "tool":
			rawTC, err := json.Marshal(map[string]any{
				"type":     "function",
				"function": map[string]string{"name": req.ToolChoice.Name},
			})
			if err != nil {
				return nil, fmt.Errorf("gateway: marshal openai tool_choice function: %w", err)
			}
			out.ToolChoice = rawTC
		}
		if req.ToolChoice.DisableParallelToolUse != nil {
			parallel := !*req.ToolChoice.DisableParallelToolUse
			out.ParallelToolCalls = &parallel
		}
	}

	if sysText := joinSystemText(req.System); sysText != "" {
		rawSys, err := json.Marshal(sysText)
		if err != nil {
			return nil, fmt.Errorf("gateway: marshal openai system: %w", err)
		}
		out.Messages = append(out.Messages, openAIChatMessage{
			Role:    "system",
			Content: rawSys,
		})
	}

	for _, msg := range req.Messages {
		encodedMsgs, err := encodeOpenAIMessage(msg)
		if err != nil {
			return nil, err
		}
		out.Messages = append(out.Messages, encodedMsgs...)
	}

	return json.Marshal(out)
}

func encodeOpenAIMessage(msg CanonicalMessage) ([]openAIChatMessage, error) {
	if msg.Role == "assistant" {
		var thinkingParts []string
		var textParts []string
		var toolCalls []openAIChatToolCall
		for _, b := range msg.Content {
			switch b.Type {
			case BlockTypeThinking:
				if b.Thinking != "" {
					thinkingParts = append(thinkingParts, b.Thinking)
				}
			case BlockTypeText:
				textParts = append(textParts, b.Text)
			case BlockTypeToolUse:
				toolCalls = append(toolCalls, openAIChatToolCall{
					ID:   b.ToolUseID,
					Type: "function",
					Function: openAIChatFunctionCall{
						Name:      b.ToolName,
						Arguments: rawJSONToString(b.ToolInput),
					},
				})
			}
		}
		var rawContent json.RawMessage
		switch {
		case len(textParts) > 0:
			b, err := json.Marshal(strings.Join(textParts, ""))
			if err != nil {
				return nil, fmt.Errorf("gateway: marshal assistant text: %w", err)
			}
			rawContent = b
		case len(toolCalls) > 0:
			rawContent = json.RawMessage("null")
		default:
			rawContent = json.RawMessage(`""`)
		}
		return []openAIChatMessage{{
			Role:             "assistant",
			Name:             msg.Name,
			Content:          rawContent,
			ReasoningContent: strings.Join(thinkingParts, "\n"),
			ToolCalls:        toolCalls,
		}}, nil
	}

	// For user or tool messages, emit each tool_result as a dedicated role:"tool" message
	// and keep any remaining user text/image blocks in a user message.
	var result []openAIChatMessage
	var userBlocks []ContentBlock
	for _, b := range msg.Content {
		if b.Type == BlockTypeToolResult {
			rawText, err := json.Marshal(toolResultText(b))
			if err != nil {
				return nil, fmt.Errorf("gateway: marshal tool result: %w", err)
			}
			result = append(result, openAIChatMessage{
				Role:       "tool",
				Name:       b.ToolName,
				ToolCallID: b.ToolUseID,
				Content:    rawText,
			})
			continue
		}
		userBlocks = append(userBlocks, b)
	}

	if len(userBlocks) > 0 {
		rawContent, err := encodeOpenAIUserContent(userBlocks)
		if err != nil {
			return nil, err
		}
		role := msg.Role
		if role == "tool" || role == "" {
			role = "user"
		}
		result = append(result, openAIChatMessage{
			Role:    role,
			Name:    msg.Name,
			Content: rawContent,
		})
	}
	return result, nil
}

func encodeOpenAIUserContent(blocks []ContentBlock) (json.RawMessage, error) {
	hasImage := false
	for _, b := range blocks {
		if b.Type == BlockTypeImage {
			hasImage = true
			break
		}
	}
	if !hasImage {
		var sb strings.Builder
		for _, b := range blocks {
			if b.Type == BlockTypeText {
				sb.WriteString(b.Text)
			}
		}
		return json.Marshal(sb.String())
	}

	parts := make([]openAIChatContentPart, 0, len(blocks))
	for _, b := range blocks {
		switch b.Type {
		case BlockTypeText:
			parts = append(parts, openAIChatContentPart{
				Type: "text",
				Text: b.Text,
			})
		case BlockTypeImage:
			if b.Image == nil {
				continue
			}
			imgURL := b.Image.URL
			if imgURL == "" && b.Image.Data != "" {
				imgURL = formatDataURI(b.Image.MediaType, b.Image.Data)
			}
			if imgURL != "" {
				parts = append(parts, openAIChatContentPart{
					Type: "image_url",
					ImageURL: &openAIImageURL{
						URL:    imgURL,
						Detail: b.Image.Detail,
					},
				})
			}
		}
	}
	return json.Marshal(parts)
}

func decodeOpenAIChatResponse(body []byte) (*CanonicalChatResponse, error) {
	var raw openAIChatResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("%w: decode openai-chat response: %w", ErrInvalidPayload, err)
	}

	resp := &CanonicalChatResponse{
		ID:         raw.ID,
		Model:      raw.Model,
		StopReason: StopReasonStop,
	}
	if len(raw.Choices) > 0 {
		ch := raw.Choices[0]
		resp.StopReason = mapOpenAIFinishReasonToCanonical(ch.FinishReason)
		if ch.Message.ReasoningContent != "" {
			resp.Content = append(resp.Content, ContentBlock{
				Type:     BlockTypeThinking,
				Thinking: ch.Message.ReasoningContent,
			})
		}
		if ch.Message.Content != nil && *ch.Message.Content != "" {
			resp.Content = append(resp.Content, ContentBlock{
				Type: BlockTypeText,
				Text: *ch.Message.Content,
			})
		}
		for _, tc := range ch.Message.ToolCalls {
			resp.Content = append(resp.Content, ContentBlock{
				Type:      BlockTypeToolUse,
				ToolUseID: tc.ID,
				ToolName:  tc.Function.Name,
				ToolInput: parseJSONObjectOrWrap(tc.Function.Arguments),
			})
		}
	}
	if raw.Usage != nil {
		resp.Usage.InputTokens = raw.Usage.PromptTokens
		resp.Usage.OutputTokens = raw.Usage.CompletionTokens
		if raw.Usage.PromptTokensDetails != nil {
			resp.Usage.CachedTokens = raw.Usage.PromptTokensDetails.CachedTokens
		}
	}
	return resp, nil
}

func encodeOpenAIChatResponse(resp *CanonicalChatResponse) ([]byte, error) {
	var thinkingParts []string
	var textParts []string
	var toolCalls []openAIChatToolCall
	for _, b := range resp.Content {
		switch b.Type {
		case BlockTypeThinking:
			if b.Thinking != "" {
				thinkingParts = append(thinkingParts, b.Thinking)
			}
		case BlockTypeText:
			textParts = append(textParts, b.Text)
		case BlockTypeToolUse:
			toolCalls = append(toolCalls, openAIChatToolCall{
				ID:   b.ToolUseID,
				Type: "function",
				Function: openAIChatFunctionCall{
					Name:      b.ToolName,
					Arguments: rawJSONToString(b.ToolInput),
				},
			})
		}
	}

	var contentPtr *string
	if len(textParts) > 0 || len(toolCalls) == 0 {
		joined := strings.Join(textParts, "")
		contentPtr = &joined
	}

	usage := &openAIChatUsage{
		PromptTokens:     resp.Usage.InputTokens,
		CompletionTokens: resp.Usage.OutputTokens,
		TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
	}
	if resp.Usage.CachedTokens > 0 {
		usage.PromptTokensDetails = &openAIPromptTokenDetail{
			CachedTokens: resp.Usage.CachedTokens,
		}
	}

	out := openAIChatResponse{
		ID:     resp.ID,
		Object: "chat.completion",
		Model:  resp.Model,
		Choices: []openAIChatChoice{{
			Index: 0,
			Message: openAIChatOutMessage{
				Role:             "assistant",
				Content:          contentPtr,
				ReasoningContent: strings.Join(thinkingParts, "\n"),
				ToolCalls:        toolCalls,
			},
			FinishReason: mapCanonicalStopReasonToOpenAI(resp.StopReason),
		}},
		Usage: usage,
	}
	return json.Marshal(out)
}

func mapOpenAIFinishReasonToCanonical(reason string) string {
	switch reason {
	case "length":
		return StopReasonLength
	case "tool_calls", "function_call":
		return StopReasonToolCalls
	case "content_filter":
		return StopReasonContentFilter
	default:
		return StopReasonStop
	}
}

func mapCanonicalStopReasonToOpenAI(reason string) string {
	switch reason {
	case StopReasonLength:
		return "length"
	case StopReasonToolCalls:
		return "tool_calls"
	case StopReasonContentFilter:
		return "content_filter"
	default:
		return "stop"
	}
}

// --- OpenAI Chat Streaming Decoder & Encoder ---

type openAIStreamChunkPayload struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role             string               `json:"role,omitempty"`
			Content          string               `json:"content,omitempty"`
			ReasoningContent string               `json:"reasoning_content,omitempty"`
			ToolCalls        []openAIChatToolCall `json:"tool_calls,omitempty"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason,omitempty"`
	} `json:"choices"`
	Usage *openAIChatUsage `json:"usage,omitempty"`
}

type openAIStreamDecoder struct {
	model          string
	started        bool
	activeBlock    string
	activeBlockIdx int
	nextBlockIdx   int
	toolBlocks     map[int]int // openAI tool_call index -> canonical block index
	openToolBlocks map[int]bool
}

func newOpenAIStreamDecoder(model string) *openAIStreamDecoder {
	return &openAIStreamDecoder{
		model:          model,
		toolBlocks:     make(map[int]int),
		openToolBlocks: make(map[int]bool),
	}
}

func (d *openAIStreamDecoder) closeActiveBlock(chunks *[]CanonicalStreamChunk) {
	if d.activeBlock == "" {
		return
	}
	*chunks = append(*chunks, CanonicalStreamChunk{
		EventType:  ChunkBlockStop,
		BlockIndex: d.activeBlockIdx,
		BlockType:  d.activeBlock,
	})
	d.activeBlock = ""
}

func (d *openAIStreamDecoder) closeAllOpenBlocks(chunks *[]CanonicalStreamChunk) {
	d.closeActiveBlock(chunks)
	if len(d.openToolBlocks) == 0 {
		return
	}
	indices := make([]int, 0, len(d.openToolBlocks))
	for idx := range d.openToolBlocks {
		indices = append(indices, idx)
	}
	slices.Sort(indices)
	for _, idx := range indices {
		*chunks = append(*chunks, CanonicalStreamChunk{
			EventType:  ChunkBlockStop,
			BlockIndex: idx,
			BlockType:  BlockTypeToolUse,
		})
		delete(d.openToolBlocks, idx)
	}
}

func (d *openAIStreamDecoder) Decode(_ string, data []byte) ([]CanonicalStreamChunk, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, nil
	}
	if bytes.Equal(trimmed, []byte("[DONE]")) {
		var chunks []CanonicalStreamChunk
		d.closeAllOpenBlocks(&chunks)
		chunks = append(chunks, CanonicalStreamChunk{
			EventType: ChunkMessageStop,
			Done:      true,
		})
		return chunks, nil
	}

	var payload openAIStreamChunkPayload
	if err := json.Unmarshal(trimmed, &payload); err != nil {
		return nil, fmt.Errorf("%w: decode openai stream chunk: %w", ErrInvalidPayload, err)
	}

	var chunks []CanonicalStreamChunk
	model := payload.Model
	if model == "" {
		model = d.model
	}

	var usage *CanonicalUsage
	if payload.Usage != nil {
		u := CanonicalUsage{
			InputTokens:  payload.Usage.PromptTokens,
			OutputTokens: payload.Usage.CompletionTokens,
		}
		if payload.Usage.PromptTokensDetails != nil {
			u.CachedTokens = payload.Usage.PromptTokensDetails.CachedTokens
		}
		usage = &u
	}

	if !d.started {
		d.started = true
		chunks = append(chunks, CanonicalStreamChunk{
			ID:        payload.ID,
			Model:     model,
			EventType: ChunkMessageStart,
			Usage:     usage,
		})
	}

	if len(payload.Choices) == 0 {
		if usage != nil {
			chunks = append(chunks, CanonicalStreamChunk{
				EventType: ChunkMessageDelta,
				Usage:     usage,
			})
		}
		return chunks, nil
	}

	ch := payload.Choices[0]
	if ch.Delta.ReasoningContent != "" {
		if d.activeBlock != BlockTypeThinking || len(d.openToolBlocks) > 0 {
			d.closeAllOpenBlocks(&chunks)
			d.activeBlock = BlockTypeThinking
			d.activeBlockIdx = d.nextBlockIdx
			d.nextBlockIdx++
			chunks = append(chunks, CanonicalStreamChunk{
				EventType:  ChunkBlockStart,
				BlockIndex: d.activeBlockIdx,
				BlockType:  BlockTypeThinking,
			})
		}
		chunks = append(chunks, CanonicalStreamChunk{
			EventType:     ChunkBlockDelta,
			BlockIndex:    d.activeBlockIdx,
			BlockType:     BlockTypeThinking,
			ThinkingDelta: ch.Delta.ReasoningContent,
		})
	}

	if ch.Delta.Content != "" {
		if d.activeBlock != BlockTypeText || len(d.openToolBlocks) > 0 {
			d.closeAllOpenBlocks(&chunks)
			d.activeBlock = BlockTypeText
			d.activeBlockIdx = d.nextBlockIdx
			d.nextBlockIdx++
			chunks = append(chunks, CanonicalStreamChunk{
				EventType:  ChunkBlockStart,
				BlockIndex: d.activeBlockIdx,
				BlockType:  BlockTypeText,
			})
		}
		chunks = append(chunks, CanonicalStreamChunk{
			EventType:  ChunkBlockDelta,
			BlockIndex: d.activeBlockIdx,
			BlockType:  BlockTypeText,
			TextDelta:  ch.Delta.Content,
		})
	}

	for _, tc := range ch.Delta.ToolCalls {
		d.closeActiveBlock(&chunks)
		tcIdx := 0
		if tc.Index != nil {
			tcIdx = *tc.Index
		}
		blockIdx, seen := d.toolBlocks[tcIdx]
		if !seen {
			blockIdx = d.nextBlockIdx
			d.nextBlockIdx++
			d.toolBlocks[tcIdx] = blockIdx
			d.openToolBlocks[blockIdx] = true
			chunks = append(chunks, CanonicalStreamChunk{
				EventType:  ChunkBlockStart,
				BlockIndex: blockIdx,
				BlockType:  BlockTypeToolUse,
				ToolUseID:  tc.ID,
				ToolName:   tc.Function.Name,
			})
		}
		if tc.Function.Arguments != "" {
			chunks = append(chunks, CanonicalStreamChunk{
				EventType:          ChunkBlockDelta,
				BlockIndex:         blockIdx,
				BlockType:          BlockTypeToolUse,
				ToolInputJSONDelta: tc.Function.Arguments,
			})
		}
	}

	if ch.FinishReason != "" {
		d.closeAllOpenBlocks(&chunks)
		chunks = append(chunks, CanonicalStreamChunk{
			EventType:  ChunkMessageDelta,
			StopReason: mapOpenAIFinishReasonToCanonical(ch.FinishReason),
			Usage:      usage,
		})
	} else if usage != nil && len(chunks) == 0 {
		chunks = append(chunks, CanonicalStreamChunk{
			EventType: ChunkMessageDelta,
			Usage:     usage,
		})
	}

	return chunks, nil
}

type openAIStreamEncoder struct {
	id            string
	model         string
	sentRole      bool
	sentFinish    bool
	sentUsage     bool
	stopReason    string
	toolIndexByID map[int]int
	nextToolIndex int
}

func newOpenAIStreamEncoder(model string) *openAIStreamEncoder {
	return &openAIStreamEncoder{
		id:            "chatcmpl-gateway",
		model:         model,
		stopReason:    StopReasonStop,
		toolIndexByID: make(map[int]int),
	}
}

func (e *openAIStreamEncoder) makeChunk(delta map[string]any, finishReason *string, usage *CanonicalUsage) ([]byte, error) {
	if !e.sentRole && delta != nil {
		delta["role"] = "assistant"
		e.sentRole = true
	}
	choice := map[string]any{
		"index":         0,
		"delta":         delta,
		"finish_reason": finishReason,
	}
	payload := map[string]any{
		"id":      e.id,
		"object":  "chat.completion.chunk",
		"model":   e.model,
		"choices": []any{choice},
	}
	if usage != nil && (usage.InputTokens > 0 || usage.OutputTokens > 0 || usage.CachedTokens > 0) {
		u := &openAIChatUsage{
			PromptTokens:     usage.InputTokens,
			CompletionTokens: usage.OutputTokens,
			TotalTokens:      usage.InputTokens + usage.OutputTokens,
		}
		if usage.CachedTokens > 0 {
			u.PromptTokensDetails = &openAIPromptTokenDetail{CachedTokens: usage.CachedTokens}
		}
		payload["usage"] = u
	}
	return marshalSSEFrame("", payload)
}

func (e *openAIStreamEncoder) Encode(chunk CanonicalStreamChunk) ([][]byte, error) {
	if chunk.ID != "" {
		e.id = chunk.ID
	}
	if chunk.Model != "" {
		e.model = chunk.Model
	}

	switch chunk.EventType {
	case ChunkMessageStart:
		if !e.sentRole {
			frame, err := e.makeChunk(map[string]any{}, nil, nil)
			if err != nil {
				return nil, err
			}
			return [][]byte{frame}, nil
		}
		return nil, nil

	case ChunkBlockStart:
		if chunk.BlockType == BlockTypeToolUse {
			toolIdx := e.nextToolIndex
			e.nextToolIndex++
			e.toolIndexByID[chunk.BlockIndex] = toolIdx
			frame, err := e.makeChunk(map[string]any{
				"tool_calls": []map[string]any{{
					"index": toolIdx,
					"id":    chunk.ToolUseID,
					"type":  "function",
					"function": map[string]any{
						"name":      chunk.ToolName,
						"arguments": "",
					},
				}},
			}, nil, nil)
			if err != nil {
				return nil, err
			}
			return [][]byte{frame}, nil
		}
		return nil, nil

	case ChunkBlockDelta:
		switch chunk.BlockType {
		case BlockTypeThinking:
			if chunk.ThinkingDelta == "" {
				return nil, nil
			}
			frame, err := e.makeChunk(map[string]any{
				"reasoning_content": chunk.ThinkingDelta,
			}, nil, nil)
			if err != nil {
				return nil, err
			}
			return [][]byte{frame}, nil
		case BlockTypeText:
			if chunk.TextDelta == "" {
				return nil, nil
			}
			frame, err := e.makeChunk(map[string]any{
				"content": chunk.TextDelta,
			}, nil, nil)
			if err != nil {
				return nil, err
			}
			return [][]byte{frame}, nil
		case BlockTypeToolUse:
			if chunk.ToolInputJSONDelta == "" {
				return nil, nil
			}
			toolIdx, ok := e.toolIndexByID[chunk.BlockIndex]
			if !ok {
				toolIdx = e.nextToolIndex
				e.nextToolIndex++
				e.toolIndexByID[chunk.BlockIndex] = toolIdx
			}
			frame, err := e.makeChunk(map[string]any{
				"tool_calls": []map[string]any{{
					"index": toolIdx,
					"function": map[string]any{
						"arguments": chunk.ToolInputJSONDelta,
					},
				}},
			}, nil, nil)
			if err != nil {
				return nil, err
			}
			return [][]byte{frame}, nil
		}
		return nil, nil

	case ChunkMessageDelta:
		if chunk.StopReason != "" {
			e.stopReason = chunk.StopReason
			fin := mapCanonicalStopReasonToOpenAI(chunk.StopReason)
			e.sentFinish = true
			if chunk.Usage != nil && (chunk.Usage.InputTokens > 0 || chunk.Usage.OutputTokens > 0 || chunk.Usage.CachedTokens > 0) {
				e.sentUsage = true
			}
			frame, err := e.makeChunk(map[string]any{}, &fin, chunk.Usage)
			if err != nil {
				return nil, err
			}
			return [][]byte{frame}, nil
		}
		if e.sentFinish && !e.sentUsage && chunk.Usage != nil && (chunk.Usage.InputTokens > 0 || chunk.Usage.OutputTokens > 0 || chunk.Usage.CachedTokens > 0) {
			e.sentUsage = true
			u := &openAIChatUsage{
				PromptTokens:     chunk.Usage.InputTokens,
				CompletionTokens: chunk.Usage.OutputTokens,
				TotalTokens:      chunk.Usage.InputTokens + chunk.Usage.OutputTokens,
			}
			if chunk.Usage.CachedTokens > 0 {
				u.PromptTokensDetails = &openAIPromptTokenDetail{CachedTokens: chunk.Usage.CachedTokens}
			}
			frame, err := marshalSSEFrame("", map[string]any{
				"id":      e.id,
				"object":  "chat.completion.chunk",
				"model":   e.model,
				"choices": []any{},
				"usage":   u,
			})
			if err != nil {
				return nil, err
			}
			return [][]byte{frame}, nil
		}
		return nil, nil

	case ChunkMessageStop:
		var frames [][]byte
		if !e.sentFinish {
			fin := mapCanonicalStopReasonToOpenAI(e.stopReason)
			e.sentFinish = true
			frame, err := e.makeChunk(map[string]any{}, &fin, chunk.Usage)
			if err != nil {
				return nil, err
			}
			frames = append(frames, frame)
		}
		frames = append(frames, formatSSEFrame("", []byte("[DONE]")))
		return frames, nil
	}

	return nil, nil
}

func (e *openAIStreamEncoder) Finalize(usage CanonicalUsage) [][]byte {
	var frames [][]byte
	if !e.sentFinish {
		fin := mapCanonicalStopReasonToOpenAI(e.stopReason)
		e.sentFinish = true
		if frame, err := e.makeChunk(map[string]any{}, &fin, &usage); err == nil {
			frames = append(frames, frame)
		}
	}
	frames = append(frames, formatSSEFrame("", []byte("[DONE]")))
	return frames
}

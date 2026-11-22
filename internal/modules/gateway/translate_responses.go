package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type responsesRequest struct {
	Model              string            `json:"model"`
	Instructions       string            `json:"instructions,omitempty"`
	Input              json.RawMessage   `json:"input"`
	Tools              []responsesTool   `json:"tools,omitempty"`
	ToolChoice         json.RawMessage   `json:"tool_choice,omitempty"`
	ParallelToolCalls  *bool             `json:"parallel_tool_calls,omitempty"`
	Text               *responsesTextCfg `json:"text,omitempty"`
	MaxOutputTokens    *int              `json:"max_output_tokens,omitempty"`
	Temperature        *float64          `json:"temperature,omitempty"`
	TopP               *float64          `json:"top_p,omitempty"`
	Stream             bool              `json:"stream,omitempty"`
	Reasoning          *responsesReason  `json:"reasoning,omitempty"`
	PreviousResponseID string            `json:"previous_response_id,omitempty"`
	Metadata           map[string]any    `json:"metadata,omitempty"`
}

type responsesReason struct {
	Effort  string `json:"effort,omitempty"`
	Summary string `json:"summary,omitempty"`
}

type responsesTextCfg struct {
	Format *responsesFormat `json:"format,omitempty"`
}

type responsesFormat struct {
	Type   string          `json:"type"`
	Name   string          `json:"name,omitempty"`
	Schema json.RawMessage `json:"schema,omitempty"`
	Strict *bool           `json:"strict,omitempty"`
}

type responsesTool struct {
	Type        string                  `json:"type"`
	Name        string                  `json:"name,omitempty"`
	Description string                  `json:"description,omitempty"`
	Parameters  json.RawMessage         `json:"parameters,omitempty"`
	Strict      *bool                   `json:"strict,omitempty"`
	Function    *openAIChatFunctionSpec `json:"function,omitempty"`
}

type responsesItem struct {
	Type             string                 `json:"type,omitempty"`
	ID               string                 `json:"id,omitempty"`
	Role             string                 `json:"role,omitempty"`
	Status           string                 `json:"status,omitempty"`
	Content          json.RawMessage        `json:"content,omitempty"`
	Summary          []responsesSummaryText `json:"summary,omitempty"`
	EncryptedContent string                 `json:"encrypted_content,omitempty"`
	CallID           string                 `json:"call_id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	Arguments        string                 `json:"arguments,omitempty"`
	Output           json.RawMessage        `json:"output,omitempty"`
}

type responsesSummaryText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type responsesContentPart struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	ImageURL json.RawMessage `json:"image_url,omitempty"`
	Detail   string          `json:"detail,omitempty"`
}

type responsesResponse struct {
	ID                string                      `json:"id"`
	Object            string                      `json:"object"`
	Status            string                      `json:"status"`
	Model             string                      `json:"model"`
	Output            []responsesItem             `json:"output"`
	IncompleteDetails *responsesIncompleteDetails `json:"incomplete_details,omitempty"`
	Usage             *responsesUsage             `json:"usage,omitempty"`
}

type responsesIncompleteDetails struct {
	Reason string `json:"reason"`
}

type responsesUsage struct {
	InputTokens        int64                        `json:"input_tokens"`
	OutputTokens       int64                        `json:"output_tokens"`
	TotalTokens        int64                        `json:"total_tokens"`
	InputTokensDetails *responsesInputTokensDetails `json:"input_tokens_details,omitempty"`
}

type responsesInputTokensDetails struct {
	CachedTokens int64 `json:"cached_tokens,omitempty"`
}

func decodeOpenAIResponsesRequest(body []byte) (*CanonicalChatRequest, error) {
	var raw responsesRequest
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("%w: decode openai-responses request: %w", ErrInvalidPayload, err)
	}

	req := &CanonicalChatRequest{
		Model:              raw.Model,
		MaxTokens:          raw.MaxOutputTokens,
		Temperature:        raw.Temperature,
		TopP:               raw.TopP,
		Stream:             raw.Stream,
		PreviousResponseID: raw.PreviousResponseID,
		Metadata:           raw.Metadata,
	}
	if raw.Instructions != "" {
		req.System = append(req.System, SystemBlock{Text: raw.Instructions})
	}
	if raw.Reasoning != nil && raw.Reasoning.Effort != "" {
		req.ReasoningEffort = raw.Reasoning.Effort
	}
	if raw.Text != nil && raw.Text.Format != nil {
		req.ResponseFormat = &ResponseFormat{
			Type:   raw.Text.Format.Type,
			Name:   raw.Text.Format.Name,
			Schema: raw.Text.Format.Schema,
			Strict: raw.Text.Format.Strict,
		}
	}

	for _, t := range raw.Tools {
		switch {
		case t.Name != "":
			req.Tools = append(req.Tools, ToolDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
				Strict:      t.Strict,
			})
		case t.Function != nil && t.Function.Name != "":
			req.Tools = append(req.Tools, ToolDefinition{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
				Strict:      t.Function.Strict,
			})
		}
	}

	if len(bytes.TrimSpace(raw.ToolChoice)) > 0 && !bytes.Equal(bytes.TrimSpace(raw.ToolChoice), []byte("null")) {
		var tcStr string
		if err := json.Unmarshal(raw.ToolChoice, &tcStr); err == nil {
			req.ToolChoice = &ToolChoice{Type: tcStr}
		} else {
			var tcObj struct {
				Type string `json:"type"`
				Name string `json:"name"`
			}
			if err := json.Unmarshal(raw.ToolChoice, &tcObj); err != nil {
				return nil, fmt.Errorf("%w: invalid responses tool_choice: %w", ErrInvalidPayload, err)
			}
			if tcObj.Name != "" {
				req.ToolChoice = &ToolChoice{Type: "tool", Name: tcObj.Name}
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

	if err := decodeResponsesInput(raw.Input, req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeResponsesInput(raw json.RawMessage, req *CanonicalChatRequest) error {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil
	}
	if trimmed[0] == '"' {
		var text string
		if err := json.Unmarshal(trimmed, &text); err != nil {
			return fmt.Errorf("%w: invalid responses input string: %w", ErrInvalidPayload, err)
		}
		req.Messages = append(req.Messages, CanonicalMessage{
			Role:    "user",
			Content: []ContentBlock{{Type: BlockTypeText, Text: text}},
		})
		return nil
	}

	var items []responsesItem
	if err := json.Unmarshal(trimmed, &items); err != nil {
		return fmt.Errorf("%w: invalid responses input array: %w", ErrInvalidPayload, err)
	}

	appendAssistantBlock := func(block ContentBlock) {
		if len(req.Messages) > 0 && req.Messages[len(req.Messages)-1].Role == "assistant" {
			req.Messages[len(req.Messages)-1].Content = append(req.Messages[len(req.Messages)-1].Content, block)
		} else {
			req.Messages = append(req.Messages, CanonicalMessage{
				Role:    "assistant",
				Content: []ContentBlock{block},
			})
		}
	}

	for _, item := range items {
		itemType := item.Type
		if itemType == "" && item.Role != "" {
			itemType = "message"
		}
		switch itemType {
		case "message":
			blocks, err := decodeResponsesContent(item.Content)
			if err != nil {
				return err
			}
			switch item.Role {
			case "system", "developer":
				for _, b := range blocks {
					if b.Type == BlockTypeText && b.Text != "" {
						req.System = append(req.System, SystemBlock{Text: b.Text})
					}
				}
			case "assistant":
				for _, b := range blocks {
					appendAssistantBlock(b)
				}
			default:
				role := item.Role
				if role == "" {
					role = "user"
				}
				req.Messages = append(req.Messages, CanonicalMessage{
					Role:    role,
					Content: blocks,
				})
			}
		case "reasoning":
			var summaryParts []string
			for _, s := range item.Summary {
				if s.Text != "" {
					summaryParts = append(summaryParts, s.Text)
				}
			}
			appendAssistantBlock(ContentBlock{
				Type:      BlockTypeThinking,
				Thinking:  strings.Join(summaryParts, "\n"),
				Signature: item.EncryptedContent,
			})
		case "function_call":
			callID := item.CallID
			if callID == "" {
				callID = item.ID
			}
			appendAssistantBlock(ContentBlock{
				Type:      BlockTypeToolUse,
				ToolUseID: callID,
				ToolName:  item.Name,
				ToolInput: parseJSONObjectOrWrap(item.Arguments),
			})
		case "function_call_output":
			outText := decodeResponsesToolOutput(item.Output)
			req.Messages = append(req.Messages, CanonicalMessage{
				Role: "tool",
				Content: []ContentBlock{{
					Type:      BlockTypeToolResult,
					ToolUseID: item.CallID,
					Text:      outText,
				}},
			})
		}
	}
	return nil
}

func decodeResponsesContent(raw json.RawMessage) ([]ContentBlock, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}
	if trimmed[0] == '"' {
		var text string
		if err := json.Unmarshal(trimmed, &text); err != nil {
			return nil, fmt.Errorf("%w: invalid responses message text: %w", ErrInvalidPayload, err)
		}
		return []ContentBlock{{Type: BlockTypeText, Text: text}}, nil
	}
	var parts []responsesContentPart
	if err := json.Unmarshal(trimmed, &parts); err != nil {
		return nil, fmt.Errorf("%w: invalid responses content parts: %w", ErrInvalidPayload, err)
	}
	blocks := make([]ContentBlock, 0, len(parts))
	for _, p := range parts {
		switch p.Type {
		case "input_text", "output_text", "text":
			blocks = append(blocks, ContentBlock{Type: BlockTypeText, Text: p.Text})
		case "input_image":
			var imgURL string
			detail := p.Detail
			trimmedURL := bytes.TrimSpace(p.ImageURL)
			if len(trimmedURL) > 0 && trimmedURL[0] == '"' {
				_ = json.Unmarshal(trimmedURL, &imgURL)
			} else if len(trimmedURL) > 0 {
				var obj openAIImageURL
				if err := json.Unmarshal(trimmedURL, &obj); err == nil {
					imgURL = obj.URL
					if obj.Detail != "" {
						detail = obj.Detail
					}
				}
			}
			if imgURL != "" {
				img := &ImageSource{URL: imgURL, Detail: detail}
				if mt, data, ok := parseDataURI(imgURL); ok {
					img.MediaType = mt
					img.Data = data
				}
				blocks = append(blocks, ContentBlock{Type: BlockTypeImage, Image: img})
			}
		}
	}
	return blocks, nil
}

func decodeResponsesToolOutput(raw json.RawMessage) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return ""
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(trimmed, &s); err == nil {
			return s
		}
	}
	return string(trimmed)
}

func encodeOpenAIResponsesRequest(req *CanonicalChatRequest) ([]byte, error) {
	out := responsesRequest{
		Model:              req.Model,
		Instructions:       joinSystemText(req.System),
		MaxOutputTokens:    req.MaxTokens,
		Temperature:        req.Temperature,
		TopP:               req.TopP,
		Stream:             req.Stream,
		PreviousResponseID: req.PreviousResponseID,
		Metadata:           req.Metadata,
	}

	effort := req.ReasoningEffort
	if effort == "" && req.ThinkingBudget != nil && *req.ThinkingBudget > 0 {
		effort = thinkingBudgetToEffort(*req.ThinkingBudget)
	}
	if effort != "" {
		out.Reasoning = &responsesReason{
			Effort:  effort,
			Summary: "auto",
		}
	}

	if req.ResponseFormat != nil {
		out.Text = &responsesTextCfg{
			Format: &responsesFormat{
				Type:   req.ResponseFormat.Type,
				Name:   req.ResponseFormat.Name,
				Schema: req.ResponseFormat.Schema,
				Strict: req.ResponseFormat.Strict,
			},
		}
	}

	for _, t := range req.Tools {
		out.Tools = append(out.Tools, responsesTool{
			Type:        "function",
			Name:        t.Name,
			Description: t.Description,
			Parameters:  ensureJSONObject(t.Parameters),
			Strict:      t.Strict,
		})
	}

	if req.ToolChoice != nil {
		switch req.ToolChoice.Type {
		case "auto", "none", "required":
			rawTC, err := json.Marshal(req.ToolChoice.Type)
			if err != nil {
				return nil, fmt.Errorf("gateway: marshal responses tool_choice: %w", err)
			}
			out.ToolChoice = rawTC
		case "tool":
			rawTC, err := json.Marshal(map[string]string{
				"type": "function",
				"name": req.ToolChoice.Name,
			})
			if err != nil {
				return nil, fmt.Errorf("gateway: marshal responses tool_choice function: %w", err)
			}
			out.ToolChoice = rawTC
		}
		if req.ToolChoice.DisableParallelToolUse != nil {
			parallel := !*req.ToolChoice.DisableParallelToolUse
			out.ParallelToolCalls = &parallel
		}
	}

	items := make([]responsesItem, 0, len(req.Messages))
	for _, msg := range req.Messages {
		encoded, err := encodeResponsesMessage(msg)
		if err != nil {
			return nil, err
		}
		items = append(items, encoded...)
	}

	rawInput, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("gateway: marshal responses input: %w", err)
	}
	out.Input = rawInput

	return json.Marshal(out)
}

func encodeResponsesMessage(msg CanonicalMessage) ([]responsesItem, error) {
	var items []responsesItem

	if msg.Role == "assistant" {
		var textParts []responsesContentPart
		flushText := func() error {
			if len(textParts) == 0 {
				return nil
			}
			rawContent, err := json.Marshal(textParts)
			if err != nil {
				return fmt.Errorf("gateway: marshal responses assistant content: %w", err)
			}
			items = append(items, responsesItem{
				Type:    "message",
				Role:    "assistant",
				Content: rawContent,
			})
			textParts = nil
			return nil
		}

		for _, b := range msg.Content {
			switch b.Type {
			case BlockTypeThinking:
				if err := flushText(); err != nil {
					return nil, err
				}
				var summary []responsesSummaryText
				if b.Thinking != "" {
					summary = []responsesSummaryText{{Type: "summary_text", Text: b.Thinking}}
				}
				items = append(items, responsesItem{
					Type:             "reasoning",
					Summary:          summary,
					EncryptedContent: b.Signature,
				})
			case BlockTypeText:
				textParts = append(textParts, responsesContentPart{
					Type: "output_text",
					Text: b.Text,
				})
			case BlockTypeToolUse:
				if err := flushText(); err != nil {
					return nil, err
				}
				items = append(items, responsesItem{
					Type:      "function_call",
					CallID:    b.ToolUseID,
					Name:      b.ToolName,
					Arguments: rawJSONToString(b.ToolInput),
				})
			}
		}
		if err := flushText(); err != nil {
			return nil, err
		}
		return items, nil
	}

	var userParts []responsesContentPart
	flushUser := func() error {
		if len(userParts) == 0 {
			return nil
		}
		rawContent, err := json.Marshal(userParts)
		if err != nil {
			return fmt.Errorf("gateway: marshal responses user content: %w", err)
		}
		items = append(items, responsesItem{
			Type:    "message",
			Role:    "user",
			Content: rawContent,
		})
		userParts = nil
		return nil
	}

	for _, b := range msg.Content {
		switch b.Type {
		case BlockTypeToolResult:
			if err := flushUser(); err != nil {
				return nil, err
			}
			rawOut, err := json.Marshal(toolResultText(b))
			if err != nil {
				return nil, fmt.Errorf("gateway: marshal responses tool output: %w", err)
			}
			items = append(items, responsesItem{
				Type:   "function_call_output",
				CallID: b.ToolUseID,
				Output: rawOut,
			})
		case BlockTypeText:
			userParts = append(userParts, responsesContentPart{
				Type: "input_text",
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
				rawURL, err := json.Marshal(imgURL)
				if err != nil {
					return nil, fmt.Errorf("gateway: marshal responses image_url: %w", err)
				}
				userParts = append(userParts, responsesContentPart{
					Type:     "input_image",
					ImageURL: rawURL,
					Detail:   b.Image.Detail,
				})
			}
		}
	}
	if err := flushUser(); err != nil {
		return nil, err
	}
	return items, nil
}

func decodeOpenAIResponsesResponse(body []byte) (*CanonicalChatResponse, error) {
	var raw responsesResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("%w: decode openai-responses response: %w", ErrInvalidPayload, err)
	}

	resp := &CanonicalChatResponse{
		ID:         raw.ID,
		Model:      raw.Model,
		StopReason: StopReasonStop,
	}

	hasToolUse := false
	for _, item := range raw.Output {
		switch item.Type {
		case "reasoning":
			var parts []string
			for _, s := range item.Summary {
				if s.Text != "" {
					parts = append(parts, s.Text)
				}
			}
			resp.Content = append(resp.Content, ContentBlock{
				Type:      BlockTypeThinking,
				Thinking:  strings.Join(parts, "\n"),
				Signature: item.EncryptedContent,
			})
		case "message":
			blocks, err := decodeResponsesContent(item.Content)
			if err != nil {
				return nil, err
			}
			resp.Content = append(resp.Content, blocks...)
		case "function_call":
			hasToolUse = true
			callID := item.CallID
			if callID == "" {
				callID = item.ID
			}
			resp.Content = append(resp.Content, ContentBlock{
				Type:      BlockTypeToolUse,
				ToolUseID: callID,
				ToolName:  item.Name,
				ToolInput: parseJSONObjectOrWrap(item.Arguments),
			})
		}
	}

	resp.StopReason = mapResponsesStatusToCanonical(raw.Status, raw.IncompleteDetails, hasToolUse)
	if raw.Usage != nil {
		resp.Usage.InputTokens = raw.Usage.InputTokens
		resp.Usage.OutputTokens = raw.Usage.OutputTokens
		if raw.Usage.InputTokensDetails != nil {
			resp.Usage.CachedTokens = raw.Usage.InputTokensDetails.CachedTokens
		}
	}
	return resp, nil
}

func encodeOpenAIResponsesResponse(resp *CanonicalChatResponse) ([]byte, error) {
	status := "completed"
	var incomplete *responsesIncompleteDetails
	switch resp.StopReason {
	case StopReasonLength:
		status = "incomplete"
		incomplete = &responsesIncompleteDetails{Reason: "max_output_tokens"}
	case StopReasonContentFilter:
		status = "incomplete"
		incomplete = &responsesIncompleteDetails{Reason: "content_filter"}
	}

	items := make([]responsesItem, 0, len(resp.Content))
	var textParts []responsesContentPart
	flushText := func() error {
		if len(textParts) == 0 {
			return nil
		}
		rawContent, err := json.Marshal(textParts)
		if err != nil {
			return fmt.Errorf("gateway: marshal responses output text: %w", err)
		}
		items = append(items, responsesItem{
			Type:    "message",
			ID:      fmt.Sprintf("msg_%d", len(items)),
			Role:    "assistant",
			Status:  "completed",
			Content: rawContent,
		})
		textParts = nil
		return nil
	}

	for _, b := range resp.Content {
		switch b.Type {
		case BlockTypeThinking:
			if err := flushText(); err != nil {
				return nil, err
			}
			var summary []responsesSummaryText
			if b.Thinking != "" {
				summary = []responsesSummaryText{{Type: "summary_text", Text: b.Thinking}}
			}
			items = append(items, responsesItem{
				Type:             "reasoning",
				ID:               fmt.Sprintf("rs_%d", len(items)),
				Summary:          summary,
				EncryptedContent: b.Signature,
			})
		case BlockTypeText:
			textParts = append(textParts, responsesContentPart{
				Type: "output_text",
				Text: b.Text,
			})
		case BlockTypeToolUse:
			if err := flushText(); err != nil {
				return nil, err
			}
			items = append(items, responsesItem{
				Type:      "function_call",
				ID:        fmt.Sprintf("fc_%d", len(items)),
				CallID:    b.ToolUseID,
				Name:      b.ToolName,
				Arguments: rawJSONToString(b.ToolInput),
				Status:    "completed",
			})
		}
	}
	if err := flushText(); err != nil {
		return nil, err
	}

	usage := &responsesUsage{
		InputTokens:  resp.Usage.InputTokens,
		OutputTokens: resp.Usage.OutputTokens,
		TotalTokens:  resp.Usage.InputTokens + resp.Usage.OutputTokens,
	}
	if resp.Usage.CachedTokens > 0 {
		usage.InputTokensDetails = &responsesInputTokensDetails{
			CachedTokens: resp.Usage.CachedTokens,
		}
	}

	out := responsesResponse{
		ID:                resp.ID,
		Object:            "response",
		Status:            status,
		Model:             resp.Model,
		Output:            items,
		IncompleteDetails: incomplete,
		Usage:             usage,
	}
	return json.Marshal(out)
}

func mapResponsesStatusToCanonical(status string, incomplete *responsesIncompleteDetails, hasToolUse bool) string {
	if status == "incomplete" && incomplete != nil {
		switch incomplete.Reason {
		case "max_output_tokens":
			return StopReasonLength
		case "content_filter":
			return StopReasonContentFilter
		}
	}
	if hasToolUse {
		return StopReasonToolCalls
	}
	return StopReasonStop
}

// --- OpenAI Responses Streaming Decoder & Encoder ---

type responsesStreamDecoder struct {
	model      string
	blockTypes map[int]string
	hasToolUse bool
}

func newResponsesStreamDecoder(model string) *responsesStreamDecoder {
	return &responsesStreamDecoder{
		model:      model,
		blockTypes: make(map[int]string),
	}
}

func (d *responsesStreamDecoder) Decode(eventType string, data []byte) ([]CanonicalStreamChunk, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, nil
	}

	var ev struct {
		Type        string             `json:"type"`
		OutputIndex int                `json:"output_index"`
		Delta       string             `json:"delta"`
		Item        *responsesItem     `json:"item,omitempty"`
		Response    *responsesResponse `json:"response,omitempty"`
	}
	if err := json.Unmarshal(trimmed, &ev); err != nil {
		return nil, fmt.Errorf("%w: decode responses SSE event: %w", ErrInvalidPayload, err)
	}

	kind := ev.Type
	if kind == "" {
		kind = eventType
	}

	switch kind {
	case "response.created":
		chunk := CanonicalStreamChunk{
			Model:     d.model,
			EventType: ChunkMessageStart,
		}
		if ev.Response != nil {
			chunk.ID = ev.Response.ID
			if ev.Response.Model != "" {
				chunk.Model = ev.Response.Model
			}
		}
		return []CanonicalStreamChunk{chunk}, nil

	case "response.output_item.added":
		if ev.Item == nil {
			return nil, nil
		}
		switch ev.Item.Type {
		case "reasoning":
			d.blockTypes[ev.OutputIndex] = BlockTypeThinking
			return []CanonicalStreamChunk{{
				EventType:  ChunkBlockStart,
				BlockIndex: ev.OutputIndex,
				BlockType:  BlockTypeThinking,
			}}, nil
		case "message":
			d.blockTypes[ev.OutputIndex] = BlockTypeText
			return []CanonicalStreamChunk{{
				EventType:  ChunkBlockStart,
				BlockIndex: ev.OutputIndex,
				BlockType:  BlockTypeText,
			}}, nil
		case "function_call":
			d.hasToolUse = true
			d.blockTypes[ev.OutputIndex] = BlockTypeToolUse
			callID := ev.Item.CallID
			if callID == "" {
				callID = ev.Item.ID
			}
			return []CanonicalStreamChunk{{
				EventType:  ChunkBlockStart,
				BlockIndex: ev.OutputIndex,
				BlockType:  BlockTypeToolUse,
				ToolUseID:  callID,
				ToolName:   ev.Item.Name,
			}}, nil
		}
		return nil, nil

	case "response.reasoning_summary_text.delta", "response.reasoning_text.delta":
		return []CanonicalStreamChunk{{
			EventType:     ChunkBlockDelta,
			BlockIndex:    ev.OutputIndex,
			BlockType:     BlockTypeThinking,
			ThinkingDelta: ev.Delta,
		}}, nil

	case "response.output_text.delta":
		return []CanonicalStreamChunk{{
			EventType:  ChunkBlockDelta,
			BlockIndex: ev.OutputIndex,
			BlockType:  BlockTypeText,
			TextDelta:  ev.Delta,
		}}, nil

	case "response.function_call_arguments.delta":
		return []CanonicalStreamChunk{{
			EventType:          ChunkBlockDelta,
			BlockIndex:         ev.OutputIndex,
			BlockType:          BlockTypeToolUse,
			ToolInputJSONDelta: ev.Delta,
		}}, nil

	case "response.output_item.done":
		bt, ok := d.blockTypes[ev.OutputIndex]
		if !ok {
			return nil, nil
		}
		var chunks []CanonicalStreamChunk
		if ev.Item != nil && ev.Item.Type == "reasoning" && ev.Item.EncryptedContent != "" {
			chunks = append(chunks, CanonicalStreamChunk{
				EventType:      ChunkBlockDelta,
				BlockIndex:     ev.OutputIndex,
				BlockType:      BlockTypeThinking,
				SignatureDelta: ev.Item.EncryptedContent,
			})
		}
		chunks = append(chunks, CanonicalStreamChunk{
			EventType:  ChunkBlockStop,
			BlockIndex: ev.OutputIndex,
			BlockType:  bt,
		})
		return chunks, nil

	case "response.completed", "response.incomplete":
		var usage *CanonicalUsage
		stopReason := StopReasonStop
		if d.hasToolUse {
			stopReason = StopReasonToolCalls
		}
		if ev.Response != nil {
			stopReason = mapResponsesStatusToCanonical(ev.Response.Status, ev.Response.IncompleteDetails, d.hasToolUse)
			if ev.Response.Usage != nil {
				u := CanonicalUsage{
					InputTokens:  ev.Response.Usage.InputTokens,
					OutputTokens: ev.Response.Usage.OutputTokens,
				}
				if ev.Response.Usage.InputTokensDetails != nil {
					u.CachedTokens = ev.Response.Usage.InputTokensDetails.CachedTokens
				}
				usage = &u
			}
		}
		return []CanonicalStreamChunk{
			{
				EventType:  ChunkMessageDelta,
				StopReason: stopReason,
				Usage:      usage,
			},
			{
				EventType: ChunkMessageStop,
				Usage:     usage,
				Done:      true,
			},
		}, nil

	default:
		return nil, nil
	}
}

type responsesStreamBlockState struct {
	blockType string
	itemID    string
	callID    string
	toolName  string
	buf       strings.Builder
	signature string
}

type responsesStreamEncoder struct {
	id         string
	model      string
	started    bool
	sentDone   bool
	stopReason string
	blocks     map[int]*responsesStreamBlockState
	blockOrder []int
}

func newResponsesStreamEncoder(model string) *responsesStreamEncoder {
	return &responsesStreamEncoder{
		id:         "resp_gateway",
		model:      model,
		stopReason: StopReasonStop,
		blocks:     make(map[int]*responsesStreamBlockState),
	}
}

func (e *responsesStreamEncoder) ensureStarted(frames *[][]byte) error {
	if e.started {
		return nil
	}
	e.started = true
	frame, err := marshalSSEFrame("response.created", map[string]any{
		"type": "response.created",
		"response": map[string]any{
			"id":     e.id,
			"object": "response",
			"status": "in_progress",
			"model":  e.model,
			"output": []any{},
		},
	})
	if err != nil {
		return err
	}
	*frames = append(*frames, frame)
	return nil
}

func (e *responsesStreamEncoder) Encode(chunk CanonicalStreamChunk) ([][]byte, error) {
	if chunk.ID != "" {
		e.id = chunk.ID
	}
	if chunk.Model != "" {
		e.model = chunk.Model
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
		st := &responsesStreamBlockState{
			blockType: chunk.BlockType,
			callID:    chunk.ToolUseID,
			toolName:  chunk.ToolName,
		}
		var item map[string]any
		switch chunk.BlockType {
		case BlockTypeThinking:
			st.itemID = fmt.Sprintf("rs_%d", chunk.BlockIndex)
			item = map[string]any{
				"id":      st.itemID,
				"type":    "reasoning",
				"summary": []any{},
			}
		case BlockTypeToolUse:
			st.itemID = fmt.Sprintf("fc_%d", chunk.BlockIndex)
			item = map[string]any{
				"id":        st.itemID,
				"type":      "function_call",
				"call_id":   st.callID,
				"name":      st.toolName,
				"arguments": "",
			}
		default:
			st.itemID = fmt.Sprintf("msg_%d", chunk.BlockIndex)
			item = map[string]any{
				"id":      st.itemID,
				"type":    "message",
				"role":    "assistant",
				"status":  "in_progress",
				"content": []any{},
			}
		}
		e.blocks[chunk.BlockIndex] = st
		e.blockOrder = append(e.blockOrder, chunk.BlockIndex)

		addedFrame, err := marshalSSEFrame("response.output_item.added", map[string]any{
			"type":         "response.output_item.added",
			"output_index": chunk.BlockIndex,
			"item":         item,
		})
		if err != nil {
			return nil, err
		}
		frames = append(frames, addedFrame)
		if chunk.BlockType == BlockTypeText {
			partFrame, err := marshalSSEFrame("response.content_part.added", map[string]any{
				"type":          "response.content_part.added",
				"item_id":       st.itemID,
				"output_index":  chunk.BlockIndex,
				"content_index": 0,
				"part":          map[string]any{"type": "output_text", "text": ""},
			})
			if err != nil {
				return nil, err
			}
			frames = append(frames, partFrame)
		}
		return frames, nil

	case ChunkBlockDelta:
		if err := e.ensureStarted(&frames); err != nil {
			return nil, err
		}
		st, ok := e.blocks[chunk.BlockIndex]
		if !ok {
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
			st = e.blocks[chunk.BlockIndex]
		}
		switch {
		case chunk.ThinkingDelta != "":
			st.buf.WriteString(chunk.ThinkingDelta)
			frame, err := marshalSSEFrame("response.reasoning_summary_text.delta", map[string]any{
				"type":         "response.reasoning_summary_text.delta",
				"item_id":      st.itemID,
				"output_index": chunk.BlockIndex,
				"delta":        chunk.ThinkingDelta,
			})
			if err != nil {
				return nil, err
			}
			frames = append(frames, frame)
		case chunk.SignatureDelta != "":
			st.signature += chunk.SignatureDelta
		case chunk.TextDelta != "":
			st.buf.WriteString(chunk.TextDelta)
			frame, err := marshalSSEFrame("response.output_text.delta", map[string]any{
				"type":          "response.output_text.delta",
				"item_id":       st.itemID,
				"output_index":  chunk.BlockIndex,
				"content_index": 0,
				"delta":         chunk.TextDelta,
			})
			if err != nil {
				return nil, err
			}
			frames = append(frames, frame)
		case chunk.ToolInputJSONDelta != "":
			st.buf.WriteString(chunk.ToolInputJSONDelta)
			frame, err := marshalSSEFrame("response.function_call_arguments.delta", map[string]any{
				"type":         "response.function_call_arguments.delta",
				"item_id":      st.itemID,
				"output_index": chunk.BlockIndex,
				"delta":        chunk.ToolInputJSONDelta,
			})
			if err != nil {
				return nil, err
			}
			frames = append(frames, frame)
		}
		return frames, nil

	case ChunkBlockStop:
		st, ok := e.blocks[chunk.BlockIndex]
		if !ok {
			return nil, nil
		}
		doneFrame, err := marshalSSEFrame("response.output_item.done", map[string]any{
			"type":         "response.output_item.done",
			"output_index": chunk.BlockIndex,
			"item":         e.buildOutputItem(st),
		})
		if err != nil {
			return nil, err
		}
		return [][]byte{doneFrame}, nil

	case ChunkMessageDelta:
		if chunk.StopReason != "" {
			e.stopReason = chunk.StopReason
		}
		return nil, nil

	case ChunkMessageStop:
		var u CanonicalUsage
		if chunk.Usage != nil {
			u = *chunk.Usage
		}
		return e.emitCompleted(u)
	}
	return nil, nil
}

func (e *responsesStreamEncoder) buildOutputItem(st *responsesStreamBlockState) responsesItem {
	switch st.blockType {
	case BlockTypeThinking:
		var summary []responsesSummaryText
		if st.buf.Len() > 0 {
			summary = []responsesSummaryText{{Type: "summary_text", Text: st.buf.String()}}
		}
		return responsesItem{
			Type:             "reasoning",
			ID:               st.itemID,
			Summary:          summary,
			EncryptedContent: st.signature,
		}
	case BlockTypeToolUse:
		return responsesItem{
			Type:      "function_call",
			ID:        st.itemID,
			CallID:    st.callID,
			Name:      st.toolName,
			Arguments: st.buf.String(),
			Status:    "completed",
		}
	default:
		rawContent, _ := json.Marshal([]responsesContentPart{{
			Type: "output_text",
			Text: st.buf.String(),
		}})
		return responsesItem{
			Type:    "message",
			ID:      st.itemID,
			Role:    "assistant",
			Status:  "completed",
			Content: rawContent,
		}
	}
}

func (e *responsesStreamEncoder) emitCompleted(usage CanonicalUsage) ([][]byte, error) {
	if e.sentDone {
		return nil, nil
	}
	e.sentDone = true

	var frames [][]byte
	if err := e.ensureStarted(&frames); err != nil {
		return nil, err
	}

	output := make([]responsesItem, 0, len(e.blockOrder))
	for _, idx := range e.blockOrder {
		if st, ok := e.blocks[idx]; ok {
			output = append(output, e.buildOutputItem(st))
		}
	}

	status := "completed"
	eventName := "response.completed"
	var incomplete *responsesIncompleteDetails
	switch e.stopReason {
	case StopReasonLength:
		status = "incomplete"
		eventName = "response.incomplete"
		incomplete = &responsesIncompleteDetails{Reason: "max_output_tokens"}
	case StopReasonContentFilter:
		status = "incomplete"
		eventName = "response.incomplete"
		incomplete = &responsesIncompleteDetails{Reason: "content_filter"}
	}

	u := &responsesUsage{
		InputTokens:  usage.InputTokens,
		OutputTokens: usage.OutputTokens,
		TotalTokens:  usage.InputTokens + usage.OutputTokens,
	}
	if usage.CachedTokens > 0 {
		u.InputTokensDetails = &responsesInputTokensDetails{CachedTokens: usage.CachedTokens}
	}

	completedFrame, err := marshalSSEFrame(eventName, map[string]any{
		"type": eventName,
		"response": responsesResponse{
			ID:                e.id,
			Object:            "response",
			Status:            status,
			Model:             e.model,
			Output:            output,
			IncompleteDetails: incomplete,
			Usage:             u,
		},
	})
	if err != nil {
		return nil, err
	}
	frames = append(frames, completedFrame)
	return frames, nil
}

func (e *responsesStreamEncoder) Finalize(usage CanonicalUsage) [][]byte {
	frames, _ := e.emitCompleted(usage)
	return frames
}

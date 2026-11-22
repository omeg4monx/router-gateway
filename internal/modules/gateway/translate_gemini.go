package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type geminiRequest struct {
	SystemInstruction      *geminiContent    `json:"systemInstruction,omitempty"`
	SystemInstructionSnake *geminiContent    `json:"system_instruction,omitempty"`
	Contents               []geminiContent   `json:"contents"`
	Tools                  []geminiToolGroup `json:"tools,omitempty"`
	ToolConfig             *geminiToolConfig `json:"toolConfig,omitempty"`
	GenerationConfig       *geminiGenConfig  `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text                  string                  `json:"text,omitempty"`
	Thought               bool                    `json:"thought,omitempty"`
	InlineData            *geminiInlineData       `json:"inlineData,omitempty"`
	InlineDataSnake       *geminiInlineData       `json:"inline_data,omitempty"`
	FileData              *geminiFileData         `json:"fileData,omitempty"`
	FileDataSnake         *geminiFileData         `json:"file_data,omitempty"`
	FunctionCall          *geminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionCallSnake     *geminiFunctionCall     `json:"function_call,omitempty"`
	FunctionResponse      *geminiFunctionResponse `json:"functionResponse,omitempty"`
	FunctionResponseSnake *geminiFunctionResponse `json:"function_response,omitempty"`
}

type geminiInlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

type geminiFileData struct {
	MimeType string `json:"mimeType,omitempty"`
	FileURI  string `json:"fileUri"`
}

type geminiFunctionCall struct {
	ID   string          `json:"id,omitempty"`
	Name string          `json:"name"`
	Args json.RawMessage `json:"args,omitempty"`
}

type geminiFunctionResponse struct {
	ID       string          `json:"id,omitempty"`
	Name     string          `json:"name"`
	Response json.RawMessage `json:"response"`
}

type geminiToolGroup struct {
	FunctionDeclarations      []geminiFunctionDecl `json:"functionDeclarations,omitempty"`
	FunctionDeclarationsSnake []geminiFunctionDecl `json:"function_declarations,omitempty"`
}

type geminiFunctionDecl struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

type geminiToolConfig struct {
	FunctionCallingConfig *geminiFunctionCallingConfig `json:"functionCallingConfig,omitempty"`
}

type geminiFunctionCallingConfig struct {
	Mode                 string   `json:"mode,omitempty"`
	AllowedFunctionNames []string `json:"allowedFunctionNames,omitempty"`
}

type geminiGenConfig struct {
	MaxOutputTokens  *int               `json:"maxOutputTokens,omitempty"`
	Temperature      *float64           `json:"temperature,omitempty"`
	TopP             *float64           `json:"topP,omitempty"`
	StopSequences    []string           `json:"stopSequences,omitempty"`
	ResponseMimeType string             `json:"responseMimeType,omitempty"`
	ResponseSchema   json.RawMessage    `json:"responseSchema,omitempty"`
	ThinkingConfig   *geminiThinkingCfg `json:"thinkingConfig,omitempty"`
}

type geminiThinkingCfg struct {
	ThinkingBudget  *int `json:"thinkingBudget,omitempty"`
	IncludeThoughts bool `json:"includeThoughts,omitempty"`
}

type geminiResponse struct {
	ResponseID    string               `json:"responseId,omitempty"`
	ModelVersion  string               `json:"modelVersion,omitempty"`
	Candidates    []geminiCandidate    `json:"candidates"`
	UsageMetadata *geminiUsageMetadata `json:"usageMetadata,omitempty"`
}

type geminiCandidate struct {
	Index        int           `json:"index"`
	Content      geminiContent `json:"content"`
	FinishReason string        `json:"finishReason,omitempty"`
}

type geminiUsageMetadata struct {
	PromptTokenCount        int64 `json:"promptTokenCount,omitempty"`
	CandidatesTokenCount    int64 `json:"candidatesTokenCount,omitempty"`
	TotalTokenCount         int64 `json:"totalTokenCount,omitempty"`
	CachedContentTokenCount int64 `json:"cachedContentTokenCount,omitempty"`
	ThoughtsTokenCount      int64 `json:"thoughtsTokenCount,omitempty"`
}

func decodeGeminiRequest(body []byte) (*CanonicalChatRequest, error) {
	var raw geminiRequest
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("%w: decode gemini request: %w", ErrInvalidPayload, err)
	}

	req := &CanonicalChatRequest{}
	sys := raw.SystemInstruction
	if sys == nil {
		sys = raw.SystemInstructionSnake
	}
	if sys != nil {
		for _, p := range sys.Parts {
			if p.Text != "" {
				req.System = append(req.System, SystemBlock{Text: p.Text})
			}
		}
	}

	if raw.GenerationConfig != nil {
		gc := raw.GenerationConfig
		req.MaxTokens = gc.MaxOutputTokens
		req.Temperature = gc.Temperature
		req.TopP = gc.TopP
		req.StopSequences = gc.StopSequences
		if gc.ResponseMimeType == "application/json" {
			if len(bytes.TrimSpace(gc.ResponseSchema)) > 0 {
				req.ResponseFormat = &ResponseFormat{
					Type:   "json_schema",
					Schema: gc.ResponseSchema,
				}
			} else {
				req.ResponseFormat = &ResponseFormat{Type: "json_object"}
			}
		}
		if gc.ThinkingConfig != nil && gc.ThinkingConfig.ThinkingBudget != nil && *gc.ThinkingConfig.ThinkingBudget > 0 {
			budget := *gc.ThinkingConfig.ThinkingBudget
			req.ThinkingBudget = &budget
			req.ReasoningEffort = thinkingBudgetToEffort(budget)
		}
	}

	for _, tg := range raw.Tools {
		decls := tg.FunctionDeclarations
		if len(decls) == 0 {
			decls = tg.FunctionDeclarationsSnake
		}
		for _, fd := range decls {
			req.Tools = append(req.Tools, ToolDefinition{
				Name:        fd.Name,
				Description: fd.Description,
				Parameters:  fd.Parameters,
			})
		}
	}

	if raw.ToolConfig != nil && raw.ToolConfig.FunctionCallingConfig != nil {
		cfg := raw.ToolConfig.FunctionCallingConfig
		switch strings.ToUpper(cfg.Mode) {
		case "NONE":
			req.ToolChoice = &ToolChoice{Type: "none"}
		case "AUTO":
			req.ToolChoice = &ToolChoice{Type: "auto"}
		case "ANY":
			if len(cfg.AllowedFunctionNames) == 1 {
				req.ToolChoice = &ToolChoice{Type: "tool", Name: cfg.AllowedFunctionNames[0]}
			} else {
				req.ToolChoice = &ToolChoice{Type: "required"}
			}
		}
	}

	callCounter := 0
	lastCallIDByName := make(map[string]string)

	for _, c := range raw.Contents {
		role := "user"
		if c.Role == "model" || c.Role == "assistant" {
			role = "assistant"
		}
		var blocks []ContentBlock
		hasToolResult := false
		for _, p := range c.Parts {
			inline := p.InlineData
			if inline == nil {
				inline = p.InlineDataSnake
			}
			fileData := p.FileData
			if fileData == nil {
				fileData = p.FileDataSnake
			}
			fc := p.FunctionCall
			if fc == nil {
				fc = p.FunctionCallSnake
			}
			fr := p.FunctionResponse
			if fr == nil {
				fr = p.FunctionResponseSnake
			}

			switch {
			case p.Thought && p.Text != "":
				blocks = append(blocks, ContentBlock{
					Type:     BlockTypeThinking,
					Thinking: p.Text,
				})
			case p.Text != "":
				blocks = append(blocks, ContentBlock{
					Type: BlockTypeText,
					Text: p.Text,
				})
			case inline != nil:
				blocks = append(blocks, ContentBlock{
					Type: BlockTypeImage,
					Image: &ImageSource{
						MediaType: inline.MimeType,
						Data:      inline.Data,
						URL:       formatDataURI(inline.MimeType, inline.Data),
					},
				})
			case fileData != nil:
				blocks = append(blocks, ContentBlock{
					Type: BlockTypeImage,
					Image: &ImageSource{
						MediaType: fileData.MimeType,
						URL:       fileData.FileURI,
					},
				})
			case fc != nil:
				callCounter++
				callID := fc.ID
				if callID == "" {
					callID = fmt.Sprintf("call_gem_%d", callCounter)
				}
				lastCallIDByName[fc.Name] = callID
				blocks = append(blocks, ContentBlock{
					Type:      BlockTypeToolUse,
					ToolUseID: callID,
					ToolName:  fc.Name,
					ToolInput: ensureJSONObject(fc.Args),
				})
			case fr != nil:
				hasToolResult = true
				callID := fr.ID
				if callID == "" {
					callID = lastCallIDByName[fr.Name]
				}
				blocks = append(blocks, ContentBlock{
					Type:      BlockTypeToolResult,
					ToolUseID: callID,
					ToolName:  fr.Name,
					Text:      extractGeminiFunctionResponseText(fr.Response),
				})
			}
		}
		if hasToolResult && role == "user" {
			role = "tool"
		}
		if len(blocks) > 0 {
			req.Messages = append(req.Messages, CanonicalMessage{
				Role:    role,
				Content: blocks,
			})
		}
	}

	return req, nil
}

func extractGeminiFunctionResponseText(raw json.RawMessage) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return ""
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &obj); err == nil && len(obj) == 1 {
		for _, key := range []string{"result", "output", "content", "error"} {
			if val, ok := obj[key]; ok {
				var s string
				if err := json.Unmarshal(val, &s); err == nil {
					return s
				}
			}
		}
	}
	return string(trimmed)
}

func encodeGeminiRequest(req *CanonicalChatRequest) ([]byte, error) {
	out := geminiRequest{}

	if sysText := joinSystemText(req.System); sysText != "" {
		out.SystemInstruction = &geminiContent{
			Parts: []geminiPart{{Text: sysText}},
		}
	}

	var gc geminiGenConfig
	hasGenCfg := false
	if req.MaxTokens != nil {
		gc.MaxOutputTokens = req.MaxTokens
		hasGenCfg = true
	}
	if req.Temperature != nil {
		gc.Temperature = req.Temperature
		hasGenCfg = true
	}
	if req.TopP != nil {
		gc.TopP = req.TopP
		hasGenCfg = true
	}
	if len(req.StopSequences) > 0 {
		gc.StopSequences = req.StopSequences
		hasGenCfg = true
	}
	if req.ResponseFormat != nil {
		switch req.ResponseFormat.Type {
		case "json_object":
			gc.ResponseMimeType = "application/json"
			hasGenCfg = true
		case "json_schema":
			gc.ResponseMimeType = "application/json"
			gc.ResponseSchema = req.ResponseFormat.Schema
			hasGenCfg = true
		}
	}
	budget := 0
	if req.ThinkingBudget != nil && *req.ThinkingBudget > 0 {
		budget = *req.ThinkingBudget
	} else if req.ReasoningEffort != "" {
		budget = effortToThinkingBudget(req.ReasoningEffort)
	}
	if budget > 0 {
		b := budget
		gc.ThinkingConfig = &geminiThinkingCfg{
			ThinkingBudget:  &b,
			IncludeThoughts: true,
		}
		hasGenCfg = true
	}
	if hasGenCfg {
		out.GenerationConfig = &gc
	}

	if len(req.Tools) > 0 {
		decls := make([]geminiFunctionDecl, 0, len(req.Tools))
		for _, t := range req.Tools {
			decls = append(decls, geminiFunctionDecl{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  ensureJSONObject(t.Parameters),
			})
		}
		out.Tools = []geminiToolGroup{{FunctionDeclarations: decls}}
	}

	if req.ToolChoice != nil {
		cfg := &geminiFunctionCallingConfig{}
		switch req.ToolChoice.Type {
		case "none":
			cfg.Mode = "NONE"
		case "required":
			cfg.Mode = "ANY"
		case "tool":
			cfg.Mode = "ANY"
			if req.ToolChoice.Name != "" {
				cfg.AllowedFunctionNames = []string{req.ToolChoice.Name}
			}
		default:
			cfg.Mode = "AUTO"
		}
		out.ToolConfig = &geminiToolConfig{FunctionCallingConfig: cfg}
	}

	// Build tool_use_id -> tool_name lookup so tool_result messages without Name still populate functionResponse.name.
	toolNameByID := make(map[string]string)
	for _, m := range req.Messages {
		for _, b := range m.Content {
			if b.Type == BlockTypeToolUse && b.ToolUseID != "" && b.ToolName != "" {
				toolNameByID[b.ToolUseID] = b.ToolName
			}
		}
	}

	for _, m := range req.Messages {
		role := "user"
		if m.Role == "assistant" {
			role = "model"
		}
		parts := make([]geminiPart, 0, len(m.Content))
		for _, b := range m.Content {
			switch b.Type {
			case BlockTypeThinking:
				if b.Thinking != "" {
					parts = append(parts, geminiPart{Text: b.Thinking, Thought: true})
				}
			case BlockTypeText:
				parts = append(parts, geminiPart{Text: b.Text})
			case BlockTypeImage:
				if b.Image == nil {
					continue
				}
				if b.Image.Data != "" {
					mt := b.Image.MediaType
					if mt == "" {
						mt = "image/png"
					}
					parts = append(parts, geminiPart{
						InlineData: &geminiInlineData{
							MimeType: mt,
							Data:     b.Image.Data,
						},
					})
				} else if b.Image.URL != "" {
					if mt, data, ok := parseDataURI(b.Image.URL); ok {
						parts = append(parts, geminiPart{
							InlineData: &geminiInlineData{
								MimeType: mt,
								Data:     data,
							},
						})
					} else {
						parts = append(parts, geminiPart{
							FileData: &geminiFileData{
								MimeType: b.Image.MediaType,
								FileURI:  b.Image.URL,
							},
						})
					}
				}
			case BlockTypeToolUse:
				parts = append(parts, geminiPart{
					FunctionCall: &geminiFunctionCall{
						ID:   b.ToolUseID,
						Name: b.ToolName,
						Args: ensureJSONObject(b.ToolInput),
					},
				})
			case BlockTypeToolResult:
				name := b.ToolName
				if name == "" {
					name = toolNameByID[b.ToolUseID]
				}
				if name == "" {
					name = m.Name
				}
				respText := toolResultText(b)
				var respObj json.RawMessage
				trimmed := strings.TrimSpace(respText)
				if trimmed != "" && trimmed[0] == '{' && json.Valid([]byte(trimmed)) {
					respObj = json.RawMessage(trimmed)
				} else {
					key := "result"
					if b.IsError {
						key = "error"
					}
					encoded, err := json.Marshal(map[string]string{key: respText})
					if err != nil {
						return nil, fmt.Errorf("gateway: marshal gemini functionResponse: %w", err)
					}
					respObj = json.RawMessage(encoded)
				}
				parts = append(parts, geminiPart{
					FunctionResponse: &geminiFunctionResponse{
						ID:       b.ToolUseID,
						Name:     name,
						Response: respObj,
					},
				})
			}
		}
		if len(parts) > 0 {
			out.Contents = append(out.Contents, geminiContent{
				Role:  role,
				Parts: parts,
			})
		}
	}

	return json.Marshal(out)
}

func decodeGeminiResponse(body []byte) (*CanonicalChatResponse, error) {
	var raw geminiResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("%w: decode gemini response: %w", ErrInvalidPayload, err)
	}

	resp := &CanonicalChatResponse{
		ID:         raw.ResponseID,
		Model:      raw.ModelVersion,
		StopReason: StopReasonStop,
	}

	hasToolUse := false
	callCounter := 0
	if len(raw.Candidates) > 0 {
		cand := raw.Candidates[0]
		for _, p := range cand.Content.Parts {
			fc := p.FunctionCall
			if fc == nil {
				fc = p.FunctionCallSnake
			}
			switch {
			case p.Thought && p.Text != "":
				resp.Content = append(resp.Content, ContentBlock{
					Type:     BlockTypeThinking,
					Thinking: p.Text,
				})
			case p.Text != "":
				resp.Content = append(resp.Content, ContentBlock{
					Type: BlockTypeText,
					Text: p.Text,
				})
			case fc != nil:
				hasToolUse = true
				callCounter++
				callID := fc.ID
				if callID == "" {
					callID = fmt.Sprintf("call_gem_%d", callCounter)
				}
				resp.Content = append(resp.Content, ContentBlock{
					Type:      BlockTypeToolUse,
					ToolUseID: callID,
					ToolName:  fc.Name,
					ToolInput: ensureJSONObject(fc.Args),
				})
			}
		}
		resp.StopReason = mapGeminiFinishReasonToCanonical(cand.FinishReason, hasToolUse)
	}

	if raw.UsageMetadata != nil {
		resp.Usage = geminiUsageToCanonical(raw.UsageMetadata)
	}
	return resp, nil
}

func encodeGeminiResponse(resp *CanonicalChatResponse) ([]byte, error) {
	parts := make([]geminiPart, 0, len(resp.Content))
	for _, b := range resp.Content {
		switch b.Type {
		case BlockTypeThinking:
			if b.Thinking != "" {
				parts = append(parts, geminiPart{Text: b.Thinking, Thought: true})
			}
		case BlockTypeText:
			parts = append(parts, geminiPart{Text: b.Text})
		case BlockTypeToolUse:
			parts = append(parts, geminiPart{
				FunctionCall: &geminiFunctionCall{
					ID:   b.ToolUseID,
					Name: b.ToolName,
					Args: ensureJSONObject(b.ToolInput),
				},
			})
		}
	}

	out := geminiResponse{
		ResponseID:   resp.ID,
		ModelVersion: resp.Model,
		Candidates: []geminiCandidate{{
			Index: 0,
			Content: geminiContent{
				Role:  "model",
				Parts: parts,
			},
			FinishReason: mapCanonicalStopReasonToGemini(resp.StopReason),
		}},
		UsageMetadata: &geminiUsageMetadata{
			PromptTokenCount:        resp.Usage.InputTokens,
			CandidatesTokenCount:    resp.Usage.OutputTokens,
			TotalTokenCount:         resp.Usage.InputTokens + resp.Usage.OutputTokens,
			CachedContentTokenCount: resp.Usage.CachedTokens,
		},
	}
	return json.Marshal(out)
}

func geminiUsageToCanonical(u *geminiUsageMetadata) CanonicalUsage {
	outTokens := u.CandidatesTokenCount
	if u.ThoughtsTokenCount > 0 && u.PromptTokenCount+u.CandidatesTokenCount+u.ThoughtsTokenCount == u.TotalTokenCount {
		outTokens += u.ThoughtsTokenCount
	}
	return CanonicalUsage{
		InputTokens:  u.PromptTokenCount,
		OutputTokens: outTokens,
		CachedTokens: u.CachedContentTokenCount,
	}
}

func mapGeminiFinishReasonToCanonical(reason string, hasToolUse bool) string {
	switch strings.ToUpper(reason) {
	case "MAX_TOKENS":
		return StopReasonLength
	case "SAFETY", "RECITATION", "BLOCKLIST", "PROHIBITED_CONTENT", "SPII":
		return StopReasonContentFilter
	default:
		if hasToolUse {
			return StopReasonToolCalls
		}
		return StopReasonStop
	}
}

func mapCanonicalStopReasonToGemini(reason string) string {
	switch reason {
	case StopReasonLength:
		return "MAX_TOKENS"
	case StopReasonContentFilter:
		return "SAFETY"
	default:
		return "STOP"
	}
}

// --- Gemini Streaming Decoder & Encoder ---

type geminiStreamDecoder struct {
	model          string
	started        bool
	activeBlock    string
	activeBlockIdx int
	nextBlockIdx   int
	callCounter    int
	hasToolUse     bool
}

func newGeminiStreamDecoder(model string) *geminiStreamDecoder {
	return &geminiStreamDecoder{model: model}
}

func (d *geminiStreamDecoder) closeActiveBlock(chunks *[]CanonicalStreamChunk) {
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

func (d *geminiStreamDecoder) Decode(_ string, data []byte) ([]CanonicalStreamChunk, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, nil
	}

	var payload geminiResponse
	if err := json.Unmarshal(trimmed, &payload); err != nil {
		return nil, fmt.Errorf("%w: decode gemini stream chunk: %w", ErrInvalidPayload, err)
	}

	var chunks []CanonicalStreamChunk
	model := payload.ModelVersion
	if model == "" {
		model = d.model
	}

	var usage *CanonicalUsage
	if payload.UsageMetadata != nil {
		u := geminiUsageToCanonical(payload.UsageMetadata)
		usage = &u
	}

	if !d.started {
		d.started = true
		chunks = append(chunks, CanonicalStreamChunk{
			ID:        payload.ResponseID,
			Model:     model,
			EventType: ChunkMessageStart,
			Usage:     usage,
		})
	}

	if len(payload.Candidates) == 0 {
		if usage != nil {
			chunks = append(chunks, CanonicalStreamChunk{
				EventType: ChunkMessageDelta,
				Usage:     usage,
			})
		}
		return chunks, nil
	}

	cand := payload.Candidates[0]
	for _, p := range cand.Content.Parts {
		fc := p.FunctionCall
		if fc == nil {
			fc = p.FunctionCallSnake
		}
		switch {
		case p.Thought && p.Text != "":
			if d.activeBlock != BlockTypeThinking {
				d.closeActiveBlock(&chunks)
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
				ThinkingDelta: p.Text,
			})
		case p.Text != "":
			if d.activeBlock != BlockTypeText {
				d.closeActiveBlock(&chunks)
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
				TextDelta:  p.Text,
			})
		case fc != nil:
			d.closeActiveBlock(&chunks)
			d.hasToolUse = true
			d.callCounter++
			callID := fc.ID
			if callID == "" {
				callID = fmt.Sprintf("call_gem_%d", d.callCounter)
			}
			idx := d.nextBlockIdx
			d.nextBlockIdx++
			chunks = append(chunks,
				CanonicalStreamChunk{
					EventType:  ChunkBlockStart,
					BlockIndex: idx,
					BlockType:  BlockTypeToolUse,
					ToolUseID:  callID,
					ToolName:   fc.Name,
				},
				CanonicalStreamChunk{
					EventType:          ChunkBlockDelta,
					BlockIndex:         idx,
					BlockType:          BlockTypeToolUse,
					ToolInputJSONDelta: rawJSONToString(ensureJSONObject(fc.Args)),
				},
				CanonicalStreamChunk{
					EventType:  ChunkBlockStop,
					BlockIndex: idx,
					BlockType:  BlockTypeToolUse,
				},
			)
		}
	}

	if cand.FinishReason != "" {
		d.closeActiveBlock(&chunks)
		chunks = append(chunks,
			CanonicalStreamChunk{
				EventType:  ChunkMessageDelta,
				StopReason: mapGeminiFinishReasonToCanonical(cand.FinishReason, d.hasToolUse),
				Usage:      usage,
			},
			CanonicalStreamChunk{
				EventType: ChunkMessageStop,
				Usage:     usage,
				Done:      true,
			},
		)
	} else if usage != nil && len(chunks) == 0 {
		chunks = append(chunks, CanonicalStreamChunk{
			EventType: ChunkMessageDelta,
			Usage:     usage,
		})
	}

	return chunks, nil
}

type geminiPendingTool struct {
	id   string
	name string
	args strings.Builder
}

type geminiStreamEncoder struct {
	model        string
	stopReason   string
	sentFinal    bool
	pendingTools map[int]*geminiPendingTool
}

func newGeminiStreamEncoder(model string) *geminiStreamEncoder {
	return &geminiStreamEncoder{
		model:        model,
		stopReason:   StopReasonStop,
		pendingTools: make(map[int]*geminiPendingTool),
	}
}

func (e *geminiStreamEncoder) Encode(chunk CanonicalStreamChunk) ([][]byte, error) {
	if chunk.Model != "" {
		e.model = chunk.Model
	}

	switch chunk.EventType {
	case ChunkBlockStart:
		if chunk.BlockType == BlockTypeToolUse {
			e.pendingTools[chunk.BlockIndex] = &geminiPendingTool{
				id:   chunk.ToolUseID,
				name: chunk.ToolName,
			}
		}
		return nil, nil

	case ChunkBlockDelta:
		switch chunk.BlockType {
		case BlockTypeThinking:
			if chunk.ThinkingDelta == "" {
				return nil, nil
			}
			frame, err := marshalSSEFrame("", geminiResponse{
				ModelVersion: e.model,
				Candidates: []geminiCandidate{{
					Index:   0,
					Content: geminiContent{Role: "model", Parts: []geminiPart{{Text: chunk.ThinkingDelta, Thought: true}}},
				}},
			})
			if err != nil {
				return nil, err
			}
			return [][]byte{frame}, nil
		case BlockTypeText:
			if chunk.TextDelta == "" {
				return nil, nil
			}
			frame, err := marshalSSEFrame("", geminiResponse{
				ModelVersion: e.model,
				Candidates: []geminiCandidate{{
					Index:   0,
					Content: geminiContent{Role: "model", Parts: []geminiPart{{Text: chunk.TextDelta}}},
				}},
			})
			if err != nil {
				return nil, err
			}
			return [][]byte{frame}, nil
		case BlockTypeToolUse:
			pt, ok := e.pendingTools[chunk.BlockIndex]
			if !ok {
				pt = &geminiPendingTool{id: chunk.ToolUseID, name: chunk.ToolName}
				e.pendingTools[chunk.BlockIndex] = pt
			}
			pt.args.WriteString(chunk.ToolInputJSONDelta)
			return nil, nil
		}
		return nil, nil

	case ChunkBlockStop:
		pt, ok := e.pendingTools[chunk.BlockIndex]
		if !ok {
			return nil, nil
		}
		delete(e.pendingTools, chunk.BlockIndex)
		frame, err := marshalSSEFrame("", geminiResponse{
			ModelVersion: e.model,
			Candidates: []geminiCandidate{{
				Index: 0,
				Content: geminiContent{
					Role: "model",
					Parts: []geminiPart{{
						FunctionCall: &geminiFunctionCall{
							ID:   pt.id,
							Name: pt.name,
							Args: parseJSONObjectOrWrap(pt.args.String()),
						},
					}},
				},
			}},
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
		var u CanonicalUsage
		if chunk.Usage != nil {
			u = *chunk.Usage
		}
		return e.emitFinal(u)
	}
	return nil, nil
}

func (e *geminiStreamEncoder) emitFinal(usage CanonicalUsage) ([][]byte, error) {
	if e.sentFinal {
		return nil, nil
	}
	e.sentFinal = true

	var parts []geminiPart
	pendingIndices := make([]int, 0, len(e.pendingTools))
	for idx := range e.pendingTools {
		pendingIndices = append(pendingIndices, idx)
	}
	slices.Sort(pendingIndices)
	for _, idx := range pendingIndices {
		pt := e.pendingTools[idx]
		parts = append(parts, geminiPart{
			FunctionCall: &geminiFunctionCall{
				ID:   pt.id,
				Name: pt.name,
				Args: parseJSONObjectOrWrap(pt.args.String()),
			},
		})
		delete(e.pendingTools, idx)
	}

	frame, err := marshalSSEFrame("", geminiResponse{
		ModelVersion: e.model,
		Candidates: []geminiCandidate{{
			Index: 0,
			Content: geminiContent{
				Role:  "model",
				Parts: parts,
			},
			FinishReason: mapCanonicalStopReasonToGemini(e.stopReason),
		}},
		UsageMetadata: &geminiUsageMetadata{
			PromptTokenCount:        usage.InputTokens,
			CandidatesTokenCount:    usage.OutputTokens,
			TotalTokenCount:         usage.InputTokens + usage.OutputTokens,
			CachedContentTokenCount: usage.CachedTokens,
		},
	})
	if err != nil {
		return nil, err
	}
	return [][]byte{frame}, nil
}

func (e *geminiStreamEncoder) Finalize(usage CanonicalUsage) [][]byte {
	frames, _ := e.emitFinal(usage)
	return frames
}

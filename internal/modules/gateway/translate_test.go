package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestProtocolConstantsAndFastPath(t *testing.T) {
	t.Parallel()

	if ProtocolOpenAIChat != "openai-chat" {
		t.Fatalf("ProtocolOpenAIChat = %q, want %q", ProtocolOpenAIChat, "openai-chat")
	}
	if ProtocolOpenAIResponses != "openai-responses" {
		t.Fatalf("ProtocolOpenAIResponses = %q, want %q", ProtocolOpenAIResponses, "openai-responses")
	}
	if ProtocolAnthropic != "anthropic" {
		t.Fatalf("ProtocolAnthropic = %q, want %q", ProtocolAnthropic, "anthropic")
	}
	if ProtocolGemini != "gemini" {
		t.Fatalf("ProtocolGemini = %q, want %q", ProtocolGemini, "gemini")
	}

	protos := []string{
		ProtocolOpenAIChat,
		ProtocolOpenAIResponses,
		ProtocolAnthropic,
		ProtocolGemini,
	}
	for _, in := range protos {
		for _, out := range protos {
			got := CanFastPath(in, out)
			want := in == out
			if got != want {
				t.Fatalf("CanFastPath(%q, %q) = %v, want %v", in, out, got, want)
			}
		}
	}

	if CanFastPath("unknown", "unknown") {
		t.Fatal("CanFastPath should return false for unsupported protocol")
	}
}

func TestRewriteModelFieldPreservesProviderSpecificFields(t *testing.T) {
	t.Parallel()

	t.Run("anthropic preserves cache_control thinking and betas byte-for-byte", func(t *testing.T) {
		raw := []byte(`{"model":"claude-alias","max_tokens":8192,"betas":["prompt-caching-2024-07-31","interleaved-thinking-2025-05-14"],"thinking":{"type":"enabled","budget_tokens":4096},"system":[{"type":"text","text":"Be concise.","cache_control":{"type":"ephemeral","ttl":"1h"}}],"tools":[{"name":"lookup","input_schema":{"type":"object"},"cache_control":{"type":"ephemeral"}}],"messages":[{"role":"user","content":[{"type":"text","text":"Hi","cache_control":{"type":"ephemeral"}}]}]}`)

		rewritten, err := RewriteModelField(raw, ProtocolAnthropic, "claude-sonnet-4-20250514")
		if err != nil {
			t.Fatalf("RewriteModelField error: %v", err)
		}

		var origMap, newMap map[string]json.RawMessage
		if err := json.Unmarshal(raw, &origMap); err != nil {
			t.Fatalf("unmarshal raw: %v", err)
		}
		if err := json.Unmarshal(rewritten, &newMap); err != nil {
			t.Fatalf("unmarshal rewritten: %v", err)
		}

		var gotModel string
		if err := json.Unmarshal(newMap["model"], &gotModel); err != nil {
			t.Fatalf("unmarshal model: %v", err)
		}
		if gotModel != "claude-sonnet-4-20250514" {
			t.Fatalf("model = %q, want %q", gotModel, "claude-sonnet-4-20250514")
		}

		for _, key := range []string{"max_tokens", "betas", "thinking", "system", "tools", "messages"} {
			if !bytes.Equal(origMap[key], newMap[key]) {
				t.Fatalf("field %q mutated:\ngot:  %s\nwant: %s", key, newMap[key], origMap[key])
			}
		}
	})

	t.Run("openai chat preserves reasoning_effort and unknown fields", func(t *testing.T) {
		raw := []byte(`{"model":"gpt-alias","reasoning_effort":"high","store":false,"response_format":{"type":"json_schema","json_schema":{"name":"out","strict":true,"schema":{"type":"object"}}},"messages":[{"role":"user","content":"hello"}]}`)

		rewritten, err := RewriteModelField(raw, ProtocolOpenAIChat, "o3-mini")
		if err != nil {
			t.Fatalf("RewriteModelField error: %v", err)
		}

		var origMap, newMap map[string]json.RawMessage
		if err := json.Unmarshal(raw, &origMap); err != nil {
			t.Fatalf("unmarshal raw: %v", err)
		}
		if err := json.Unmarshal(rewritten, &newMap); err != nil {
			t.Fatalf("unmarshal rewritten: %v", err)
		}
		if string(newMap["model"]) != `"o3-mini"` {
			t.Fatalf("model = %s, want %q", newMap["model"], "o3-mini")
		}
		for _, key := range []string{"reasoning_effort", "store", "response_format", "messages"} {
			if !bytes.Equal(origMap[key], newMap[key]) {
				t.Fatalf("field %q mutated:\ngot:  %s\nwant: %s", key, newMap[key], origMap[key])
			}
		}
	})

	t.Run("gemini passes body through byte-for-byte unchanged", func(t *testing.T) {
		raw := []byte(`{ "contents": [{"role":"user","parts":[{"text":"hi"}]}], "generationConfig": {"thinkingConfig": {"thinkingBudget": 2048}} }`)

		rewritten, err := RewriteModelField(raw, ProtocolGemini, "gemini-2.5-pro")
		if err != nil {
			t.Fatalf("RewriteModelField error: %v", err)
		}
		if !bytes.Equal(raw, rewritten) {
			t.Fatalf("gemini body mutated:\ngot:  %s\nwant: %s", rewritten, raw)
		}
	})

	t.Run("rejects invalid protocol and payloads", func(t *testing.T) {
		if _, err := RewriteModelField([]byte(`{"model":"x"}`), "bogus", "y"); !errors.Is(err, ErrUnsupportedProtocol) {
			t.Fatalf("err = %v, want ErrUnsupportedProtocol", err)
		}
		for _, bad := range [][]byte{[]byte(""), []byte("null"), []byte("[1,2]"), []byte("{bad")} {
			if _, err := RewriteModelField(bad, ProtocolOpenAIChat, "gpt-4o"); !errors.Is(err, ErrInvalidPayload) {
				t.Fatalf("payload %q: err = %v, want ErrInvalidPayload", string(bad), err)
			}
		}
	})
}

func TestOpenAIChatToAnthropicRoundTrip(t *testing.T) {
	t.Parallel()

	openAIReq := []byte(`{
		"model": "claude-sonnet-4-20250514",
		"max_tokens": 2048,
		"temperature": 0.3,
		"top_p": 0.9,
		"stop": ["DONE"],
		"stream": true,
		"tools": [
			{
				"type": "function",
				"function": {
					"name": "get_weather",
					"description": "Fetch weather by city",
					"parameters": {"type":"object","properties":{"city":{"type":"string"}},"required":["city"]}
				}
			},
			{
				"type": "function",
				"function": {
					"name": "get_time",
					"description": "Fetch local time",
					"parameters": {"type":"object","properties":{"tz":{"type":"string"}},"required":["tz"]}
				}
			}
		],
		"tool_choice": "auto",
		"messages": [
			{"role": "system", "content": "You are a helpful router assistant."},
			{
				"role": "user",
				"content": [
					{"type": "text", "text": "Check weather and time for this image."},
					{"type": "image_url", "image_url": {"url": "data:image/png;base64,iVBORw0KGgo="}}
				]
			},
			{
				"role": "assistant",
				"content": null,
				"tool_calls": [
					{"id": "call_1", "type": "function", "function": {"name": "get_weather", "arguments": "{\"city\":\"Tokyo\"}"}},
					{"id": "call_2", "type": "function", "function": {"name": "get_time", "arguments": "{\"tz\":\"Asia/Tokyo\"}"}}
				]
			},
			{"role": "tool", "tool_call_id": "call_1", "content": "{\"temp_c\":22}"},
			{"role": "tool", "tool_call_id": "call_2", "content": "14:30 JST"}
		]
	}`)

	ir, err := DecodeRequest(ProtocolOpenAIChat, openAIReq, "")
	if err != nil {
		t.Fatalf("DecodeRequest(openai-chat) error: %v", err)
	}
	if len(ir.System) != 1 || ir.System[0].Text != "You are a helpful router assistant." {
		t.Fatalf("unexpected system blocks: %+v", ir.System)
	}
	if len(ir.Tools) != 2 || ir.Tools[0].Name != "get_weather" || ir.Tools[1].Name != "get_time" {
		t.Fatalf("unexpected tools: %+v", ir.Tools)
	}

	anthBytes, err := EncodeRequest(ProtocolAnthropic, ir)
	if err != nil {
		t.Fatalf("EncodeRequest(anthropic) error: %v", err)
	}

	var anthPayload struct {
		Model         string   `json:"model"`
		MaxTokens     int      `json:"max_tokens"`
		StopSequences []string `json:"stop_sequences"`
		Stream        bool     `json:"stream"`
		System        any      `json:"system"`
		Tools         []struct {
			Name        string          `json:"name"`
			InputSchema json.RawMessage `json:"input_schema"`
		} `json:"tools"`
		Messages []struct {
			Role    string `json:"role"`
			Content []struct {
				Type      string          `json:"type"`
				Text      string          `json:"text"`
				ID        string          `json:"id"`
				Name      string          `json:"name"`
				Input     json.RawMessage `json:"input"`
				ToolUseID string          `json:"tool_use_id"`
				Content   any             `json:"content"`
				Source    *struct {
					Type      string `json:"type"`
					MediaType string `json:"media_type"`
					Data      string `json:"data"`
				} `json:"source"`
			} `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(anthBytes, &anthPayload); err != nil {
		t.Fatalf("unmarshal anthropic payload: %v", err)
	}
	if anthPayload.MaxTokens != 2048 || !anthPayload.Stream || len(anthPayload.StopSequences) != 1 {
		t.Fatalf("unexpected anthropic top-level fields: %+v", anthPayload)
	}
	if len(anthPayload.Tools) != 2 {
		t.Fatalf("expected 2 anthropic tools, got %d", len(anthPayload.Tools))
	}
	// Consecutive OpenAI tool messages must merge into a single Anthropic user message!
	if len(anthPayload.Messages) != 3 {
		t.Fatalf("expected 3 anthropic messages (user, assistant, merged user tool_results), got %d: %s", len(anthPayload.Messages), string(anthBytes))
	}
	if anthPayload.Messages[0].Content[1].Type != "image" ||
		anthPayload.Messages[0].Content[1].Source == nil ||
		anthPayload.Messages[0].Content[1].Source.MediaType != "image/png" ||
		anthPayload.Messages[0].Content[1].Source.Data != "iVBORw0KGgo=" {
		t.Fatalf("unexpected image block in anthropic request: %+v", anthPayload.Messages[0].Content[1])
	}
	if len(anthPayload.Messages[1].Content) != 2 ||
		anthPayload.Messages[1].Content[0].Type != "tool_use" ||
		anthPayload.Messages[1].Content[1].Type != "tool_use" {
		t.Fatalf("expected 2 parallel tool_use blocks, got %+v", anthPayload.Messages[1].Content)
	}
	if anthPayload.Messages[2].Role != "user" || len(anthPayload.Messages[2].Content) != 2 ||
		anthPayload.Messages[2].Content[0].Type != "tool_result" ||
		anthPayload.Messages[2].Content[0].ToolUseID != "call_1" ||
		anthPayload.Messages[2].Content[1].ToolUseID != "call_2" {
		t.Fatalf("expected merged user tool_result blocks, got %+v", anthPayload.Messages[2])
	}

	// Response round-trip: Anthropic -> OpenAI Chat
	anthResp := []byte(`{
		"id": "msg_123",
		"type": "message",
		"role": "assistant",
		"model": "claude-sonnet-4-20250514",
		"content": [
			{"type": "thinking", "thinking": "Let me inspect both tools.", "signature": "sig_abc"},
			{"type": "text", "text": "Calling tools now."},
			{"type": "tool_use", "id": "toolu_1", "name": "get_weather", "input": {"city": "Tokyo"}}
		],
		"stop_reason": "tool_use",
		"usage": {
			"input_tokens": 120,
			"output_tokens": 45,
			"cache_read_input_tokens": 30
		}
	}`)
	canonResp, err := DecodeResponse(ProtocolAnthropic, anthResp)
	if err != nil {
		t.Fatalf("DecodeResponse(anthropic) error: %v", err)
	}
	if canonResp.StopReason != StopReasonToolCalls || canonResp.Usage.CachedTokens != 30 {
		t.Fatalf("unexpected canonical response: %+v", canonResp)
	}

	openAIRespBytes, err := EncodeResponse(ProtocolOpenAIChat, canonResp)
	if err != nil {
		t.Fatalf("EncodeResponse(openai-chat) error: %v", err)
	}
	var openAIResp struct {
		ID      string `json:"id"`
		Model   string `json:"model"`
		Choices []struct {
			FinishReason string `json:"finish_reason"`
			Message      struct {
				Role             string `json:"role"`
				Content          string `json:"content"`
				ReasoningContent string `json:"reasoning_content"`
				ToolCalls        []struct {
					ID       string `json:"id"`
					Type     string `json:"type"`
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens        int64 `json:"prompt_tokens"`
			CompletionTokens    int64 `json:"completion_tokens"`
			TotalTokens         int64 `json:"total_tokens"`
			PromptTokensDetails *struct {
				CachedTokens int64 `json:"cached_tokens"`
			} `json:"prompt_tokens_details"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(openAIRespBytes, &openAIResp); err != nil {
		t.Fatalf("unmarshal openai response: %v", err)
	}
	if len(openAIResp.Choices) != 1 || openAIResp.Choices[0].FinishReason != "tool_calls" {
		t.Fatalf("unexpected openai choices: %+v", openAIResp.Choices)
	}
	msg := openAIResp.Choices[0].Message
	if msg.ReasoningContent != "Let me inspect both tools." || msg.Content != "Calling tools now." {
		t.Fatalf("unexpected openai message content/reasoning: %+v", msg)
	}
	if len(msg.ToolCalls) != 1 || msg.ToolCalls[0].ID != "toolu_1" || msg.ToolCalls[0].Function.Name != "get_weather" {
		t.Fatalf("unexpected openai tool_calls: %+v", msg.ToolCalls)
	}
	if openAIResp.Usage.PromptTokens != 120 || openAIResp.Usage.CompletionTokens != 45 ||
		openAIResp.Usage.PromptTokensDetails == nil || openAIResp.Usage.PromptTokensDetails.CachedTokens != 30 {
		t.Fatalf("unexpected openai usage: %+v", openAIResp.Usage)
	}
}

func TestAnthropicToOpenAIChatRoundTrip(t *testing.T) {
	t.Parallel()

	anthReq := []byte(`{
		"model": "gpt-4o",
		"max_tokens": 1024,
		"thinking": {"type": "enabled", "budget_tokens": 8192},
		"system": [
			{"type": "text", "text": "System part 1", "cache_control": {"type": "ephemeral"}},
			{"type": "text", "text": "System part 2"}
		],
		"messages": [
			{"role": "user", "content": "Run calc"},
			{
				"role": "assistant",
				"content": [
					{"type": "thinking", "thinking": "Need calculator", "signature": "sig1"},
					{"type": "tool_use", "id": "toolu_9", "name": "calc", "input": {"expr": "2+2"}}
				]
			},
			{
				"role": "user",
				"content": [
					{"type": "tool_result", "tool_use_id": "toolu_9", "content": "4"},
					{"type": "text", "text": "Now explain why."}
				]
			}
		]
	}`)

	ir, err := DecodeRequest(ProtocolAnthropic, anthReq, "gpt-4o-override")
	if err != nil {
		t.Fatalf("DecodeRequest(anthropic) error: %v", err)
	}
	if ir.Model != "gpt-4o-override" {
		t.Fatalf("model = %q, want gpt-4o-override", ir.Model)
	}
	if ir.ThinkingBudget == nil || *ir.ThinkingBudget != 8192 {
		t.Fatalf("unexpected thinking budget: %v", ir.ThinkingBudget)
	}

	openAIBytes, err := EncodeRequest(ProtocolOpenAIChat, ir)
	if err != nil {
		t.Fatalf("EncodeRequest(openai-chat) error: %v", err)
	}

	var payload struct {
		Model           string `json:"model"`
		ReasoningEffort string `json:"reasoning_effort"`
		Messages        []struct {
			Role             string `json:"role"`
			Content          any    `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
			ToolCallID       string `json:"tool_call_id"`
			ToolCalls        []struct {
				ID       string `json:"id"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(openAIBytes, &payload); err != nil {
		t.Fatalf("unmarshal openai payload: %v", err)
	}
	if payload.ReasoningEffort != "medium" {
		t.Fatalf("reasoning_effort = %q, want medium", payload.ReasoningEffort)
	}
	// Expected messages: system, user("Run calc"), assistant(reasoning+tool_call), tool(toolu_9), user("Now explain why.")
	if len(payload.Messages) != 5 {
		t.Fatalf("expected 5 messages, got %d: %s", len(payload.Messages), string(openAIBytes))
	}
	if payload.Messages[0].Role != "system" {
		t.Fatalf("msg[0] role = %q, want system", payload.Messages[0].Role)
	}
	if payload.Messages[2].ReasoningContent != "Need calculator" || len(payload.Messages[2].ToolCalls) != 1 {
		t.Fatalf("unexpected assistant message: %+v", payload.Messages[2])
	}
	if payload.Messages[3].Role != "tool" || payload.Messages[3].ToolCallID != "toolu_9" || payload.Messages[3].Content != "4" {
		t.Fatalf("unexpected tool message: %+v", payload.Messages[3])
	}
	if payload.Messages[4].Role != "user" || payload.Messages[4].Content != "Now explain why." {
		t.Fatalf("unexpected trailing user message: %+v", payload.Messages[4])
	}
}

func TestOpenAIResponsesAndAnthropicDirectPair(t *testing.T) {
	t.Parallel()

	// Direction 1: OpenAI Responses -> Anthropic Messages
	responsesReq := []byte(`{
		"model": "claude-sonnet-4-20250514",
		"instructions": "Think step by step.",
		"max_output_tokens": 4096,
		"previous_response_id": "resp_prev_1",
		"reasoning": {"effort": "high", "summary": "auto"},
		"tools": [
			{
				"type": "function",
				"name": "search_docs",
				"description": "Search internal docs",
				"parameters": {"type":"object","properties":{"q":{"type":"string"}}}
			}
		],
		"tool_choice": {"type": "function", "name": "search_docs"},
		"input": [
			{"type": "message", "role": "user", "content": [{"type": "input_text", "text": "Find SSE spec"}]},
			{"type": "reasoning", "summary": [{"type": "summary_text", "text": "Searching docs first"}], "encrypted_content": "enc_sig_1"},
			{"type": "function_call", "call_id": "fc_1", "name": "search_docs", "arguments": "{\"q\":\"SSE\"}"},
			{"type": "function_call_output", "call_id": "fc_1", "output": "Found 3 matches"}
		]
	}`)

	irFromResp, err := DecodeRequest(ProtocolOpenAIResponses, responsesReq, "")
	if err != nil {
		t.Fatalf("DecodeRequest(openai-responses) error: %v", err)
	}
	if irFromResp.PreviousResponseID != "resp_prev_1" || irFromResp.ReasoningEffort != "high" {
		t.Fatalf("unexpected responses IR metadata: %+v", irFromResp)
	}

	anthBytes, err := EncodeRequest(ProtocolAnthropic, irFromResp)
	if err != nil {
		t.Fatalf("EncodeRequest(anthropic) error: %v", err)
	}
	var anthOut struct {
		Thinking *struct {
			Type         string `json:"type"`
			BudgetTokens int    `json:"budget_tokens"`
		} `json:"thinking"`
		ToolChoice *struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"tool_choice"`
		Messages []struct {
			Role    string `json:"role"`
			Content []struct {
				Type      string `json:"type"`
				Thinking  string `json:"thinking"`
				Signature string `json:"signature"`
				ID        string `json:"id"`
				Name      string `json:"name"`
				ToolUseID string `json:"tool_use_id"`
			} `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(anthBytes, &anthOut); err != nil {
		t.Fatalf("unmarshal anthropic from responses: %v", err)
	}
	if anthOut.Thinking == nil || anthOut.Thinking.Type != "enabled" || anthOut.Thinking.BudgetTokens < 8192 {
		t.Fatalf("expected thinking enabled from high effort, got %+v", anthOut.Thinking)
	}
	if anthOut.ToolChoice == nil || anthOut.ToolChoice.Type != "tool" || anthOut.ToolChoice.Name != "search_docs" {
		t.Fatalf("unexpected anthropic tool_choice: %+v", anthOut.ToolChoice)
	}
	if len(anthOut.Messages) != 3 {
		t.Fatalf("expected 3 messages (user, assistant[thinking+tool_use], user[tool_result]), got %d: %s", len(anthOut.Messages), string(anthBytes))
	}
	if anthOut.Messages[1].Content[0].Type != "thinking" || anthOut.Messages[1].Content[0].Signature != "enc_sig_1" ||
		anthOut.Messages[1].Content[1].Type != "tool_use" || anthOut.Messages[1].Content[1].ID != "fc_1" {
		t.Fatalf("unexpected assistant blocks: %+v", anthOut.Messages[1].Content)
	}
	if anthOut.Messages[2].Content[0].Type != "tool_result" || anthOut.Messages[2].Content[0].ToolUseID != "fc_1" {
		t.Fatalf("unexpected tool_result block: %+v", anthOut.Messages[2].Content)
	}

	// Direction 2: Anthropic Messages -> OpenAI Responses
	respReqBytes, err := EncodeRequest(ProtocolOpenAIResponses, irFromResp)
	if err != nil {
		t.Fatalf("EncodeRequest(openai-responses) error: %v", err)
	}
	var respOut struct {
		Instructions string `json:"instructions"`
		Reasoning    *struct {
			Effort  string `json:"effort"`
			Summary string `json:"summary"`
		} `json:"reasoning"`
		Input []struct {
			Type   string `json:"type"`
			CallID string `json:"call_id"`
			Name   string `json:"name"`
			Output string `json:"output"`
		} `json:"input"`
	}
	if err := json.Unmarshal(respReqBytes, &respOut); err != nil {
		t.Fatalf("unmarshal responses request: %v", err)
	}
	if respOut.Instructions != "Think step by step." {
		t.Fatalf("instructions = %q, want %q", respOut.Instructions, "Think step by step.")
	}
	if respOut.Reasoning == nil || respOut.Reasoning.Effort != "high" {
		t.Fatalf("unexpected reasoning config: %+v", respOut.Reasoning)
	}
	if len(respOut.Input) != 4 || respOut.Input[1].Type != "reasoning" || respOut.Input[2].Type != "function_call" || respOut.Input[3].Type != "function_call_output" {
		t.Fatalf("unexpected responses input items: %+v", respOut.Input)
	}

	// Response round-trip: OpenAI Responses -> Anthropic
	responsesResp := []byte(`{
		"id": "resp_999",
		"object": "response",
		"status": "completed",
		"model": "gpt-4.1",
		"output": [
			{"type": "reasoning", "id": "rs_1", "summary": [{"type": "summary_text", "text": "Reasoning trace"}], "encrypted_content": "sig_r1"},
			{"type": "message", "id": "msg_1", "role": "assistant", "status": "completed", "content": [{"type": "output_text", "text": "Here is the answer."}]},
			{"type": "function_call", "id": "fc_item_1", "call_id": "call_99", "name": "search_docs", "arguments": "{\"q\":\"router\"}"}
		],
		"usage": {
			"input_tokens": 200,
			"output_tokens": 80,
			"total_tokens": 280,
			"input_tokens_details": {"cached_tokens": 50}
		}
	}`)
	canonResp, err := DecodeResponse(ProtocolOpenAIResponses, responsesResp)
	if err != nil {
		t.Fatalf("DecodeResponse(openai-responses) error: %v", err)
	}
	if canonResp.StopReason != StopReasonToolCalls || canonResp.Usage.CachedTokens != 50 {
		t.Fatalf("unexpected canonical responses response: %+v", canonResp)
	}
	anthRespBytes, err := EncodeResponse(ProtocolAnthropic, canonResp)
	if err != nil {
		t.Fatalf("EncodeResponse(anthropic) error: %v", err)
	}
	var anthRespOut struct {
		ID         string `json:"id"`
		StopReason string `json:"stop_reason"`
		Content    []struct {
			Type      string `json:"type"`
			Thinking  string `json:"thinking"`
			Signature string `json:"signature"`
			Text      string `json:"text"`
			ID        string `json:"id"`
			Name      string `json:"name"`
		} `json:"content"`
		Usage struct {
			InputTokens          int64 `json:"input_tokens"`
			OutputTokens         int64 `json:"output_tokens"`
			CacheReadInputTokens int64 `json:"cache_read_input_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(anthRespBytes, &anthRespOut); err != nil {
		t.Fatalf("unmarshal anthropic response: %v", err)
	}
	if anthRespOut.StopReason != "tool_use" || len(anthRespOut.Content) != 3 ||
		anthRespOut.Content[0].Type != "thinking" || anthRespOut.Content[0].Thinking != "Reasoning trace" ||
		anthRespOut.Content[1].Type != "text" || anthRespOut.Content[1].Text != "Here is the answer." ||
		anthRespOut.Content[2].Type != "tool_use" || anthRespOut.Content[2].ID != "call_99" ||
		anthRespOut.Usage.CacheReadInputTokens != 50 {
		t.Fatalf("unexpected anthropic response from responses: %+v", anthRespOut)
	}
}

func TestGeminiAndOpenAIChatRoundTrip(t *testing.T) {
	t.Parallel()

	openAIReq := []byte(`{
		"model": "gemini-2.5-pro",
		"max_tokens": 1024,
		"temperature": 0.4,
		"response_format": {"type": "json_object"},
		"tools": [
			{
				"type": "function",
				"function": {
					"name": "lookup_price",
					"description": "Lookup price",
					"parameters": {"type":"object","properties":{"sku":{"type":"string"}}}
				}
			}
		],
		"messages": [
			{"role": "system", "content": "Return JSON."},
			{
				"role": "user",
				"content": [
					{"type": "text", "text": "Price this item"},
					{"type": "image_url", "image_url": {"url": "data:image/jpeg;base64,/9j/4AAQ"}}
				]
			},
			{
				"role": "assistant",
				"tool_calls": [
					{"id": "call_sku_1", "type": "function", "function": {"name": "lookup_price", "arguments": "{\"sku\":\"A1\"}"}}
				]
			},
			{"role": "tool", "tool_call_id": "call_sku_1", "content": "19.99"}
		]
	}`)

	ir, err := DecodeRequest(ProtocolOpenAIChat, openAIReq, "")
	if err != nil {
		t.Fatalf("DecodeRequest(openai-chat) error: %v", err)
	}
	gemBytes, err := EncodeRequest(ProtocolGemini, ir)
	if err != nil {
		t.Fatalf("EncodeRequest(gemini) error: %v", err)
	}

	var gemReq struct {
		SystemInstruction *struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"systemInstruction"`
		GenerationConfig *struct {
			MaxOutputTokens  int    `json:"maxOutputTokens"`
			ResponseMimeType string `json:"responseMimeType"`
		} `json:"generationConfig"`
		Tools []struct {
			FunctionDeclarations []struct {
				Name string `json:"name"`
			} `json:"functionDeclarations"`
		} `json:"tools"`
		Contents []struct {
			Role  string `json:"role"`
			Parts []struct {
				Text       string `json:"text"`
				InlineData *struct {
					MimeType string `json:"mimeType"`
					Data     string `json:"data"`
				} `json:"inlineData"`
				FunctionCall *struct {
					Name string          `json:"name"`
					Args json.RawMessage `json:"args"`
				} `json:"functionCall"`
				FunctionResponse *struct {
					Name     string          `json:"name"`
					Response json.RawMessage `json:"response"`
				} `json:"functionResponse"`
			} `json:"parts"`
		} `json:"contents"`
	}
	if err := json.Unmarshal(gemBytes, &gemReq); err != nil {
		t.Fatalf("unmarshal gemini req: %v", err)
	}
	if gemReq.SystemInstruction == nil || gemReq.SystemInstruction.Parts[0].Text != "Return JSON." {
		t.Fatalf("unexpected gemini systemInstruction: %+v", gemReq.SystemInstruction)
	}
	if gemReq.GenerationConfig == nil || gemReq.GenerationConfig.ResponseMimeType != "application/json" {
		t.Fatalf("unexpected gemini generationConfig: %+v", gemReq.GenerationConfig)
	}
	if len(gemReq.Contents) != 3 ||
		gemReq.Contents[0].Parts[1].InlineData == nil ||
		gemReq.Contents[0].Parts[1].InlineData.MimeType != "image/jpeg" ||
		gemReq.Contents[1].Role != "model" ||
		gemReq.Contents[1].Parts[0].FunctionCall == nil ||
		gemReq.Contents[2].Parts[0].FunctionResponse == nil ||
		gemReq.Contents[2].Parts[0].FunctionResponse.Name != "lookup_price" {
		t.Fatalf("unexpected gemini contents: %s", string(gemBytes))
	}

	// Gemini response -> OpenAI Chat response
	gemResp := []byte(`{
		"responseId": "gem_resp_1",
		"modelVersion": "gemini-2.5-pro",
		"candidates": [
			{
				"index": 0,
				"content": {
					"role": "model",
					"parts": [
						{"text": "Thinking about price", "thought": true},
						{"text": "Checking inventory"},
						{"functionCall": {"id": "fc_g1", "name": "lookup_price", "args": {"sku": "A1"}}}
					]
				},
				"finishReason": "STOP"
			}
		],
		"usageMetadata": {
			"promptTokenCount": 90,
			"candidatesTokenCount": 25,
			"cachedContentTokenCount": 15,
			"totalTokenCount": 115
		}
	}`)
	canonGemResp, err := DecodeResponse(ProtocolGemini, gemResp)
	if err != nil {
		t.Fatalf("DecodeResponse(gemini) error: %v", err)
	}
	if canonGemResp.StopReason != StopReasonToolCalls || canonGemResp.Usage.CachedTokens != 15 {
		t.Fatalf("unexpected canonical gemini response: %+v", canonGemResp)
	}

	// Also decode Gemini request back to OpenAI Chat to verify bidirectional request translation.
	irFromGem, err := DecodeRequest(ProtocolGemini, gemBytes, "gpt-4o")
	if err != nil {
		t.Fatalf("DecodeRequest(gemini) error: %v", err)
	}
	openAIRoundTrip, err := EncodeRequest(ProtocolOpenAIChat, irFromGem)
	if err != nil {
		t.Fatalf("EncodeRequest(openai-chat) from gemini error: %v", err)
	}
	if !strings.Contains(string(openAIRoundTrip), `"lookup_price"`) || !strings.Contains(string(openAIRoundTrip), `"data:image/jpeg;base64,/9j/4AAQ"`) {
		t.Fatalf("round-tripped openai request missing expected fields: %s", string(openAIRoundTrip))
	}
}

func TestStreamTranslators(t *testing.T) {
	t.Parallel()

	t.Run("anthropic SSE to openai-chat SSE", func(t *testing.T) {
		tr, err := NewStreamTranslator(ProtocolOpenAIChat, ProtocolAnthropic, "claude-sonnet-4-20250514")
		if err != nil {
			t.Fatalf("NewStreamTranslator error: %v", err)
		}

		events := []struct {
			ev   string
			data string
		}{
			{"message_start", `{"type":"message_start","message":{"id":"msg_s1","type":"message","role":"assistant","model":"claude-sonnet-4-20250514","content":[],"usage":{"input_tokens":50,"output_tokens":0,"cache_read_input_tokens":12}}}`},
			{"content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"thinking","thinking":""}}`},
			{"content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"thinking_delta","thinking":"Plan: "}}`},
			{"content_block_stop", `{"type":"content_block_stop","index":0}`},
			{"content_block_start", `{"type":"content_block_start","index":1,"content_block":{"type":"text","text":""}}`},
			{"content_block_delta", `{"type":"content_block_delta","index":1,"delta":{"type":"text_delta","text":"Hello!"}}`},
			{"content_block_stop", `{"type":"content_block_stop","index":1}`},
			{"content_block_start", `{"type":"content_block_start","index":2,"content_block":{"type":"tool_use","id":"toolu_s1","name":"get_weather","input":{}}}`},
			{"content_block_delta", `{"type":"content_block_delta","index":2,"delta":{"type":"input_json_delta","partial_json":"{\"city\":\"NYC\"}"}}`},
			{"content_block_stop", `{"type":"content_block_stop","index":2}`},
			{"message_delta", `{"type":"message_delta","delta":{"stop_reason":"tool_use"},"usage":{"output_tokens":20}}`},
			{"message_stop", `{"type":"message_stop"}`},
		}

		var allFrames strings.Builder
		var finalUsage *CanonicalUsage
		var sawDone bool
		for _, item := range events {
			frames, usage, done, err := tr.TranslateEvent(item.ev, []byte(item.data))
			if err != nil {
				t.Fatalf("TranslateEvent(%s) error: %v", item.ev, err)
			}
			for _, f := range frames {
				allFrames.Write(f)
			}
			if usage != nil {
				finalUsage = usage
			}
			if done {
				sawDone = true
			}
		}

		out := allFrames.String()
		if !sawDone {
			t.Fatal("expected done=true on message_stop")
		}
		if finalUsage == nil || finalUsage.InputTokens != 50 || finalUsage.OutputTokens != 20 || finalUsage.CachedTokens != 12 {
			t.Fatalf("unexpected usage: %+v", finalUsage)
		}
		for _, want := range []string{`"reasoning_content":"Plan: "`, `"content":"Hello!"`, `"toolu_s1"`, `"get_weather"`, `"finish_reason":"tool_calls"`, "data: [DONE]\n\n"} {
			if !strings.Contains(out, want) {
				t.Fatalf("translated openai SSE missing %s:\n%s", want, out)
			}
		}
	})

	t.Run("openai-chat SSE to anthropic SSE", func(t *testing.T) {
		tr, err := NewStreamTranslator(ProtocolAnthropic, ProtocolOpenAIChat, "gpt-4o")
		if err != nil {
			t.Fatalf("NewStreamTranslator error: %v", err)
		}

		chunks := []string{
			`{"id":"chatcmpl-1","object":"chat.completion.chunk","model":"gpt-4o","choices":[{"index":0,"delta":{"role":"assistant","content":"Hi "},"finish_reason":null}]}`,
			`{"id":"chatcmpl-1","object":"chat.completion.chunk","model":"gpt-4o","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_a","type":"function","function":{"name":"lookup","arguments":""}}]},"finish_reason":null}]}`,
			`{"id":"chatcmpl-1","object":"chat.completion.chunk","model":"gpt-4o","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"k\":1}"}}]},"finish_reason":null}]}`,
			`{"id":"chatcmpl-1","object":"chat.completion.chunk","model":"gpt-4o","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}],"usage":{"prompt_tokens":33,"completion_tokens":14,"total_tokens":47}}`,
			`[DONE]`,
		}

		var allFrames strings.Builder
		var finalUsage *CanonicalUsage
		var sawDone bool
		for _, c := range chunks {
			frames, usage, done, err := tr.TranslateEvent("", []byte(c))
			if err != nil {
				t.Fatalf("TranslateEvent error: %v", err)
			}
			for _, f := range frames {
				allFrames.Write(f)
			}
			if usage != nil {
				finalUsage = usage
			}
			if done {
				sawDone = true
			}
		}
		out := allFrames.String()
		if !sawDone {
			t.Fatal("expected done=true on [DONE]")
		}
		if finalUsage == nil || finalUsage.InputTokens != 33 || finalUsage.OutputTokens != 14 {
			t.Fatalf("unexpected usage: %+v", finalUsage)
		}
		for _, want := range []string{"event: message_start\n", "event: content_block_start\n", `"text_delta"`, `"input_json_delta"`, `"stop_reason":"tool_use"`, "event: message_stop\n"} {
			if !strings.Contains(out, want) {
				t.Fatalf("translated anthropic SSE missing %s:\n%s", want, out)
			}
		}
	})

	t.Run("openai-responses SSE to anthropic SSE direct pair", func(t *testing.T) {
		tr, err := NewStreamTranslator(ProtocolAnthropic, ProtocolOpenAIResponses, "gpt-5-codex")
		if err != nil {
			t.Fatalf("NewStreamTranslator error: %v", err)
		}

		events := []struct {
			ev   string
			data string
		}{
			{"response.created", `{"type":"response.created","response":{"id":"resp_s1","model":"gpt-5-codex","status":"in_progress"}}`},
			{"response.output_item.added", `{"type":"response.output_item.added","output_index":0,"item":{"id":"rs_1","type":"reasoning","summary":[]}}`},
			{"response.reasoning_summary_text.delta", `{"type":"response.reasoning_summary_text.delta","output_index":0,"delta":"Step 1"}`},
			{"response.output_item.done", `{"type":"response.output_item.done","output_index":0,"item":{"id":"rs_1","type":"reasoning"}}`},
			{"response.output_item.added", `{"type":"response.output_item.added","output_index":1,"item":{"id":"msg_1","type":"message","role":"assistant","content":[]}}`},
			{"response.output_text.delta", `{"type":"response.output_text.delta","output_index":1,"delta":"Answer"}`},
			{"response.output_item.done", `{"type":"response.output_item.done","output_index":1,"item":{"id":"msg_1","type":"message"}}`},
			{"response.output_item.added", `{"type":"response.output_item.added","output_index":2,"item":{"id":"fc_1","type":"function_call","call_id":"call_r1","name":"run_cmd","arguments":""}}`},
			{"response.function_call_arguments.delta", `{"type":"response.function_call_arguments.delta","output_index":2,"item_id":"fc_1","delta":"{\"cmd\":\"ls\"}"}`},
			{"response.output_item.done", `{"type":"response.output_item.done","output_index":2,"item":{"id":"fc_1","type":"function_call","call_id":"call_r1","name":"run_cmd"}}`},
			{"response.completed", `{"type":"response.completed","response":{"id":"resp_s1","status":"completed","usage":{"input_tokens":100,"output_tokens":40,"total_tokens":140,"input_tokens_details":{"cached_tokens":25}}}}`},
		}

		var allFrames strings.Builder
		var finalUsage *CanonicalUsage
		var sawDone bool
		for _, item := range events {
			frames, usage, done, err := tr.TranslateEvent(item.ev, []byte(item.data))
			if err != nil {
				t.Fatalf("TranslateEvent(%s) error: %v", item.ev, err)
			}
			for _, f := range frames {
				allFrames.Write(f)
			}
			if usage != nil {
				finalUsage = usage
			}
			if done {
				sawDone = true
			}
		}
		out := allFrames.String()
		if !sawDone || finalUsage == nil || finalUsage.InputTokens != 100 || finalUsage.OutputTokens != 40 || finalUsage.CachedTokens != 25 {
			t.Fatalf("done=%v usage=%+v", sawDone, finalUsage)
		}
		for _, want := range []string{"event: message_start\n", `"thinking_delta"`, `"Step 1"`, `"text_delta"`, `"Answer"`, `"input_json_delta"`, `"call_r1"`, `"stop_reason":"tool_use"`, "event: message_stop\n"} {
			if !strings.Contains(out, want) {
				t.Fatalf("responses->anthropic SSE missing %s:\n%s", want, out)
			}
		}
	})

	t.Run("anthropic SSE to openai-responses SSE direct pair", func(t *testing.T) {
		tr, err := NewStreamTranslator(ProtocolOpenAIResponses, ProtocolAnthropic, "claude-sonnet-4-20250514")
		if err != nil {
			t.Fatalf("NewStreamTranslator error: %v", err)
		}

		events := []struct {
			ev   string
			data string
		}{
			{"message_start", `{"type":"message_start","message":{"id":"msg_ar1","type":"message","role":"assistant","model":"claude-sonnet-4-20250514","content":[],"usage":{"input_tokens":60,"output_tokens":0,"cache_read_input_tokens":10}}}`},
			{"content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"thinking","thinking":""}}`},
			{"content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"thinking_delta","thinking":"Deep thought"}}`},
			{"content_block_stop", `{"type":"content_block_stop","index":0}`},
			{"content_block_start", `{"type":"content_block_start","index":1,"content_block":{"type":"text","text":""}}`},
			{"content_block_delta", `{"type":"content_block_delta","index":1,"delta":{"type":"text_delta","text":"Done!"}}`},
			{"content_block_stop", `{"type":"content_block_stop","index":1}`},
			{"message_delta", `{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":18}}`},
			{"message_stop", `{"type":"message_stop"}`},
		}

		var allFrames strings.Builder
		var finalUsage *CanonicalUsage
		for _, item := range events {
			frames, usage, _, err := tr.TranslateEvent(item.ev, []byte(item.data))
			if err != nil {
				t.Fatalf("TranslateEvent(%s) error: %v", item.ev, err)
			}
			for _, f := range frames {
				allFrames.Write(f)
			}
			if usage != nil {
				finalUsage = usage
			}
		}
		out := allFrames.String()
		if finalUsage == nil || finalUsage.InputTokens != 60 || finalUsage.OutputTokens != 18 || finalUsage.CachedTokens != 10 {
			t.Fatalf("unexpected usage: %+v", finalUsage)
		}
		for _, want := range []string{"event: response.created\n", "event: response.reasoning_summary_text.delta\n", `"Deep thought"`, "event: response.output_text.delta\n", `"Done!"`, "event: response.completed\n", `"cached_tokens":10`} {
			if !strings.Contains(out, want) {
				t.Fatalf("anthropic->responses SSE missing %s:\n%s", want, out)
			}
		}
	})

	t.Run("gemini SSE to openai-chat SSE and Finalize", func(t *testing.T) {
		tr, err := NewStreamTranslator(ProtocolOpenAIChat, ProtocolGemini, "gemini-2.5-flash")
		if err != nil {
			t.Fatalf("NewStreamTranslator error: %v", err)
		}

		frames1, _, done1, err := tr.TranslateEvent("", []byte(`{"candidates":[{"index":0,"content":{"role":"model","parts":[{"text":"Hello "}]}}]}`))
		if err != nil || done1 || len(frames1) == 0 {
			t.Fatalf("chunk 1: frames=%d done=%v err=%v", len(frames1), done1, err)
		}
		frames2, usage2, done2, err := tr.TranslateEvent("", []byte(`{"candidates":[{"index":0,"content":{"role":"model","parts":[{"text":"world"}]},"finishReason":"STOP"}],"usageMetadata":{"promptTokenCount":15,"candidatesTokenCount":5,"cachedContentTokenCount":3,"totalTokenCount":20}}`))
		if err != nil || !done2 {
			t.Fatalf("chunk 2: done=%v err=%v", done2, err)
		}
		if usage2 == nil || usage2.InputTokens != 15 || usage2.OutputTokens != 5 || usage2.CachedTokens != 3 {
			t.Fatalf("unexpected gemini stream usage: %+v", usage2)
		}
		var combined strings.Builder
		for _, f := range append(frames1, frames2...) {
			combined.Write(f)
		}
		if !strings.Contains(combined.String(), `"Hello "`) || !strings.Contains(combined.String(), `"world"`) || !strings.Contains(combined.String(), "data: [DONE]\n\n") {
			t.Fatalf("unexpected gemini->openai SSE:\n%s", combined.String())
		}
		// Finalize after completion must be idempotent and emit no duplicate [DONE].
		extraFrames, _ := tr.Finalize()
		if len(extraFrames) != 0 {
			t.Fatalf("expected 0 frames from idempotent Finalize, got %d", len(extraFrames))
		}
	})
}

func TestEstimateTokens(t *testing.T) {
	t.Parallel()

	if got := EstimateTokens(nil); got != 0 {
		t.Fatalf("EstimateTokens(nil) = %d, want 0", got)
	}
	if got := EstimateTokens(&CanonicalChatRequest{}); got != 0 {
		t.Fatalf("EstimateTokens(empty) = %d, want 0", got)
	}

	small := &CanonicalChatRequest{
		System: []SystemBlock{{Text: "Be concise."}},
		Messages: []CanonicalMessage{
			{Role: "user", Content: []ContentBlock{{Type: BlockTypeText, Text: "Hello"}}},
		},
	}
	large := &CanonicalChatRequest{
		System: []SystemBlock{{Text: "Be concise and thorough."}},
		Messages: []CanonicalMessage{
			{Role: "user", Content: []ContentBlock{
				{Type: BlockTypeText, Text: "Hello, please inspect the attached architecture diagram and call the tool."},
				{Type: BlockTypeImage, Image: &ImageSource{MediaType: "image/png", Data: strings.Repeat("A", 8000)}},
			}},
		},
		Tools: []ToolDefinition{
			{
				Name:        "analyze_diagram",
				Description: "Inspect an architecture diagram and return structured findings",
				Parameters:  json.RawMessage(`{"type":"object","properties":{"focus":{"type":"string"}}}`),
			},
		},
	}

	smallCount := EstimateTokens(small)
	if smallCount <= 0 {
		t.Fatalf("EstimateTokens(small) = %d, want > 0", smallCount)
	}
	if second := EstimateTokens(small); second != smallCount {
		t.Fatalf("EstimateTokens not deterministic: %d vs %d", smallCount, second)
	}

	largeCount := EstimateTokens(large)
	if largeCount <= smallCount {
		t.Fatalf("EstimateTokens(large)=%d should be > small=%d", largeCount, smallCount)
	}
	// Image base64 data should not be counted raw as text tokens (8000/4 = 2000).
	if largeCount > 600 {
		t.Fatalf("EstimateTokens(large)=%d appears to count raw base64 bytes instead of fixed image heuristic", largeCount)
	}
}

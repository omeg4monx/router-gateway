package gateway

import "strings"

const (
	tokensPerSystemBlock = 4
	tokensPerMessage     = 4
	tokensPerToolDef     = 8
	tokensPerToolUse     = 8
	tokensPerToolResult  = 4
	tokensPerImageLow    = 85
	tokensPerImageHigh   = 256
)

// EstimateTokens computes a deterministic, zero-table heuristic token count for a
// CanonicalChatRequest (used for /v1/messages/count_tokens fallback and telemetry).
// Text payloads are estimated at ceil(utf8Bytes / 4) plus structural framing overhead,
// while inline images use a fixed detail-based token estimate instead of raw base64 size.
func EstimateTokens(req *CanonicalChatRequest) int64 {
	if req == nil {
		return 0
	}

	var totalBytes int64
	var structuralTokens int64
	hasContent := false

	for _, sb := range req.System {
		hasContent = true
		structuralTokens += tokensPerSystemBlock
		totalBytes += int64(len(sb.Text))
	}

	for _, msg := range req.Messages {
		hasContent = true
		structuralTokens += tokensPerMessage
		totalBytes += int64(len(msg.Role)) + int64(len(msg.Name))
		bytesAdd, structAdd := estimateBlocks(msg.Content)
		totalBytes += bytesAdd
		structuralTokens += structAdd
	}

	for _, tool := range req.Tools {
		hasContent = true
		structuralTokens += tokensPerToolDef
		totalBytes += int64(len(tool.Name)) + int64(len(tool.Description)) + int64(len(tool.Parameters))
	}

	if req.ResponseFormat != nil && (req.ResponseFormat.Name != "" || len(req.ResponseFormat.Schema) > 0) {
		hasContent = true
		totalBytes += int64(len(req.ResponseFormat.Name)) + int64(len(req.ResponseFormat.Schema))
	}

	if !hasContent {
		return 0
	}

	textTokens := (totalBytes + 3) / 4
	total := textTokens + structuralTokens
	if total < 1 {
		return 1
	}
	return total
}

func estimateBlocks(blocks []ContentBlock) (textBytes, structuralTokens int64) {
	for _, b := range blocks {
		switch b.Type {
		case BlockTypeText:
			textBytes += int64(len(b.Text))
		case BlockTypeThinking:
			textBytes += int64(len(b.Thinking))
		case BlockTypeImage:
			if b.Image != nil && strings.EqualFold(b.Image.Detail, "low") {
				structuralTokens += tokensPerImageLow
			} else {
				structuralTokens += tokensPerImageHigh
			}
		case BlockTypeToolUse:
			structuralTokens += tokensPerToolUse
			textBytes += int64(len(b.ToolUseID)) + int64(len(b.ToolName)) + int64(len(b.ToolInput))
		case BlockTypeToolResult:
			structuralTokens += tokensPerToolResult
			textBytes += int64(len(b.ToolUseID)) + int64(len(b.ToolName)) + int64(len(b.Text))
			childBytes, childStruct := estimateBlocks(b.ToolContent)
			textBytes += childBytes
			structuralTokens += childStruct
		default:
			textBytes += int64(len(b.Text))
		}
	}
	return textBytes, structuralTokens
}

package plugins

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	defaultContextWindow = 128000
	defaultFillRatio     = 0.75
	defaultMinTurns      = 4
	longContentChars     = 4000
	minKeepChars         = 400
)

type squeezeOptions struct {
	ContextWindow int
	FillRatio     float64
	MinTurns      int
}

type squeezeStats struct {
	TokensBefore      int
	TokensAfter       int
	DroppedTurns      int
	ShortenedMessages int
	Reason            string
}

func estimateValueTokens(value any) int {
	if value == nil {
		return 0
	}
	if text, ok := value.(string); ok {
		return (len(text) + 3) / 4
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return 0
	}
	return (len(raw) + 3) / 4
}

func estimateMessageTokens(msg map[string]any) int {
	if msg == nil {
		return 0
	}
	tokens := 4 // role + framing overhead
	if content, ok := msg["content"]; ok {
		tokens += estimateValueTokens(content)
	}
	for _, field := range []string{"reasoning_content", "thinking", "tool_calls", "tool_call_id"} {
		if v, ok := msg[field]; ok && v != nil {
			tokens += estimateValueTokens(v)
		}
	}
	return tokens
}

func estimateConversationTokens(msgs []map[string]any) int {
	total := 0
	for _, msg := range msgs {
		total += estimateMessageTokens(msg)
	}
	return total
}

func buildSqueezeRecap(dropped []map[string]any) string {
	var userCount, assistantCount, toolCount, otherCount int
	for _, m := range dropped {
		role, _ := m["role"].(string)
		switch {
		case role == "user":
			userCount++
		case role == "assistant":
			assistantCount++
		case role == "tool" || m["tool_call_id"] != nil:
			toolCount++
		default:
			otherCount++
		}
	}
	parts := make([]string, 0, 4)
	if userCount > 0 {
		suffix := ""
		if userCount > 1 {
			suffix = "s"
		}
		parts = append(parts, fmt.Sprintf("%d user message%s", userCount, suffix))
	}
	if assistantCount > 0 {
		noun := "reply"
		if assistantCount > 1 {
			noun = "replies"
		}
		parts = append(parts, fmt.Sprintf("%d assistant %s", assistantCount, noun))
	}
	if toolCount > 0 {
		suffix := ""
		if toolCount > 1 {
			suffix = "s"
		}
		parts = append(parts, fmt.Sprintf("%d tool result%s", toolCount, suffix))
	}
	if otherCount > 0 {
		suffix := ""
		if otherCount > 1 {
			suffix = "s"
		}
		parts = append(parts, fmt.Sprintf("%d other turn%s", otherCount, suffix))
	}
	what := "earlier turns"
	if len(parts) > 0 {
		what = strings.Join(parts, ", ")
	}
	turnSuffix := ""
	if len(dropped) > 1 {
		turnSuffix = "s"
	}
	return fmt.Sprintf(
		"[Context Squeezer] %d earlier turn%s were dropped to fit this model's context window: %s. "+
			"The turns below are verbatim and current. If the dropped history matters, ask the user to restate it.",
		len(dropped), turnSuffix, what,
	)
}

func shortenMessageContent(msg map[string]any, charsToCut int) bool {
	text, ok := msg["content"].(string)
	if !ok {
		return false
	}
	runes := []rune(text)
	maxCut := len(runes) - minKeepChars
	if maxCut <= 0 {
		return false
	}
	cut := charsToCut
	if cut > maxCut {
		cut = maxCut
	}
	if cut <= 0 {
		return false
	}
	keepChars := len(runes) - cut
	headChars := int(float64(keepChars) * 0.7)
	if headChars < minKeepChars {
		headChars = minKeepChars
	}
	if headChars > keepChars {
		headChars = keepChars
	}
	tailChars := keepChars - headChars
	head := string(runes[:headChars])
	tail := ""
	if tailChars > 0 {
		tail = string(runes[len(runes)-tailChars:])
	}
	elided := len(runes) - headChars - tailChars
	msg["content"] = fmt.Sprintf("%s\n[... %d characters elided by Context Squeezer ...]\n%s", head, elided, tail)
	return true
}

func squeezeConversation(rawMsgs []any, opts squeezeOptions) (messages []any, changed bool, stats squeezeStats) {
	if opts.ContextWindow <= 0 || len(rawMsgs) == 0 {
		return rawMsgs, false, squeezeStats{Reason: "no-window"}
	}
	fillRatio := opts.FillRatio
	if fillRatio <= 0 || fillRatio > 1 {
		fillRatio = defaultFillRatio
	}
	minTurns := opts.MinTurns
	if minTurns <= 0 {
		minTurns = defaultMinTurns
	}

	working := make([]map[string]any, 0, len(rawMsgs))
	for _, item := range rawMsgs {
		m, ok := item.(map[string]any)
		if !ok {
			return rawMsgs, false, squeezeStats{Reason: "non-object-message"}
		}
		working = append(working, cloneMap(m))
	}

	budget := int(float64(opts.ContextWindow) * fillRatio)
	tokensBefore := estimateConversationTokens(working)
	if tokensBefore <= budget {
		return rawMsgs, false, squeezeStats{
			TokensBefore: tokensBefore,
			TokensAfter:  tokensBefore,
			Reason:       "fits",
		}
	}

	firstTurn := 0
	for firstTurn < len(working) {
		role, _ := working[firstTurn]["role"].(string)
		if role != "system" {
			break
		}
		firstTurn++
	}

	// Pass 1: shorten oversized string messages so eviction sees realistic sizes.
	shortened := 0
	perMsgTokenCap := budget / 4
	if perMsgTokenCap < 100 {
		perMsgTokenCap = 100
	}
	for i := firstTurn; i < len(working); i++ {
		content, ok := working[i]["content"].(string)
		if !ok || len(content) < longContentChars {
			continue
		}
		cost := estimateMessageTokens(working[i])
		if cost > perMsgTokenCap && shortenMessageContent(working[i], (cost-perMsgTokenCap)*4) {
			shortened++
		}
	}

	// Pass 2: evict oldest conversational turns until within budget, respecting minTurns.
	tokens := estimateConversationTokens(working)
	cut := firstTurn
	keepFloor := len(working) - minTurns
	if keepFloor < firstTurn {
		keepFloor = firstTurn
	}
	for tokens > budget && cut < keepFloor {
		tokens -= estimateMessageTokens(working[cut])
		cut++
	}

	if cut == firstTurn {
		if shortened == 0 {
			return rawMsgs, false, squeezeStats{
				TokensBefore: tokensBefore,
				TokensAfter:  tokensBefore,
				Reason:       "floor",
			}
		}
		out := make([]any, len(working))
		for i, m := range working {
			out[i] = m
		}
		return out, true, squeezeStats{
			TokensBefore:      tokensBefore,
			TokensAfter:       estimateConversationTokens(working),
			ShortenedMessages: shortened,
			Reason:            "floor",
		}
	}

	dropped := working[firstTurn:cut]
	tail := working[cut:]
	recap := map[string]any{
		"role":    "system",
		"content": buildSqueezeRecap(dropped),
	}
	result := make([]map[string]any, 0, firstTurn+1+len(tail))
	result = append(result, working[:firstTurn]...)
	result = append(result, recap)
	result = append(result, tail...)

	out := make([]any, len(result))
	for i, m := range result {
		out[i] = m
	}
	return out, true, squeezeStats{
		TokensBefore:      tokensBefore,
		TokensAfter:       estimateConversationTokens(result),
		DroppedTurns:      len(dropped),
		ShortenedMessages: shortened,
		Reason:            "trimmed",
	}
}

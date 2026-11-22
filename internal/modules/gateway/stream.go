package gateway

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

type streamResult struct {
	usage        CanonicalUsage
	ttft         time.Duration
	responseID   string
	assistantMsg *CanonicalMessage
}

func (m *Module) streamUpstreamSSE(
	ctx context.Context,
	w http.ResponseWriter,
	body io.Reader,
	reqID string,
	inProto string,
	outProto string,
	requestedModel string,
	startTime time.Time,
) (streamResult, error) {
	var res streamResult

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	rc := http.NewResponseController(w)
	_ = rc.Flush()

	fastPath := CanFastPath(inProto, outProto)
	var translator StreamTranslator
	if !fastPath {
		var err error
		translator, err = NewStreamTranslator(inProto, outProto, requestedModel)
		if err != nil {
			return res, err
		}
	}

	captureResponsesChain := inProto == ProtocolOpenAIResponses && outProto != ProtocolOpenAIResponses
	reader := bufio.NewReaderSize(body, 32*1024)

	var (
		eventType string
		dataLines [][]byte
		done      bool
	)

	flushFrames := func(frames [][]byte, evDone bool) error {
		for _, frame := range frames {
			if len(frame) == 0 {
				continue
			}
			if res.ttft == 0 {
				res.ttft = time.Since(startTime)
			}
			outFrame, err := m.applyStreamInterceptors(ctx, reqID, frame, evDone)
			if err != nil {
				return err
			}
			if captureResponsesChain {
				captureResponsesFrame(outFrame, &res)
			}
			if _, err := w.Write(outFrame); err != nil { //nolint:gosec // G705: SSE text/event-stream relays JSON frames, not HTML
				return err
			}
			if err := rc.Flush(); err != nil {
				return err
			}
		}
		return nil
	}

	dispatchEvent := func(evName string, rawData []byte) error {
		if len(rawData) == 0 {
			return nil
		}
		if fastPath {
			if res.ttft == 0 {
				res.ttft = time.Since(startTime)
			}
			mergeFastPathSSEUsage(outProto, evName, rawData, &res.usage)

			evDone := bytes.Equal(bytes.TrimSpace(rawData), []byte("[DONE]")) || evName == "message_stop"
			if len(m.interceptors) == 0 {
				if _, err := w.Write(formatSSEFrame(evName, rawData)); err != nil {
					return err
				}
				return rc.Flush()
			}

			ev := &contracts.StreamEvent{
				RequestID: reqID,
				Event:     evName,
				Data:      bytes.Clone(rawData),
				Done:      evDone,
				Metadata:  make(map[string]any),
			}
			for _, ic := range m.interceptors {
				if err := ic.OnStreamChunk(ctx, ev); err != nil {
					return err
				}
			}
			if _, err := w.Write(formatSSEFrame(ev.Event, ev.Data)); err != nil {
				return err
			}
			return rc.Flush()
		}

		frames, usage, trDone, err := translator.TranslateEvent(evName, rawData)
		if err != nil {
			return err
		}
		if usage != nil {
			res.usage = *usage
		}
		if err := flushFrames(frames, trDone); err != nil {
			return err
		}
		if trDone {
			done = true
			finalFrames, finalUsage := translator.Finalize()
			if finalUsage != nil {
				res.usage = *finalUsage
			}
			if err := flushFrames(finalFrames, true); err != nil {
				return err
			}
		}
		return nil
	}

	for {
		if ctx.Err() != nil {
			return res, ctx.Err()
		}
		line, err := reader.ReadString('\n')
		if line != "" {
			trimmed := strings.TrimRight(line, "\r\n")
			switch {
			case trimmed == "":
				if len(dataLines) > 0 {
					joined := bytes.Join(dataLines, []byte("\n"))
					if dispErr := dispatchEvent(eventType, joined); dispErr != nil {
						return res, dispErr
					}
					eventType = ""
					dataLines = nil
					if done {
						return res, nil
					}
				} else {
					eventType = ""
				}
			case strings.HasPrefix(trimmed, ":"):
				// SSE comment/heartbeat; ignore.
			case strings.HasPrefix(trimmed, "event:"):
				eventType = strings.TrimSpace(strings.TrimPrefix(trimmed, "event:"))
			case strings.HasPrefix(trimmed, "data:"):
				payload := strings.TrimPrefix(trimmed, "data:")
				payload = strings.TrimPrefix(payload, " ")
				dataLines = append(dataLines, []byte(payload))
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(dataLines) > 0 {
					joined := bytes.Join(dataLines, []byte("\n"))
					if dispErr := dispatchEvent(eventType, joined); dispErr != nil {
						return res, dispErr
					}
				}
				if !fastPath && !done && translator != nil {
					finalFrames, finalUsage := translator.Finalize()
					if finalUsage != nil {
						res.usage = *finalUsage
					}
					if flushErr := flushFrames(finalFrames, true); flushErr != nil {
						return res, flushErr
					}
				}
				return res, nil
			}
			return res, err
		}
	}
}

func (m *Module) applyStreamInterceptors(ctx context.Context, reqID string, frame []byte, done bool) ([]byte, error) {
	if len(m.interceptors) == 0 {
		return frame, nil
	}
	evName, evData := parseSSEFrame(frame)
	ev := &contracts.StreamEvent{
		RequestID: reqID,
		Event:     evName,
		Data:      bytes.Clone(evData),
		Done:      done || bytes.Equal(bytes.TrimSpace(evData), []byte("[DONE]")),
		Metadata:  make(map[string]any),
	}
	for _, ic := range m.interceptors {
		if err := ic.OnStreamChunk(ctx, ev); err != nil {
			return nil, err
		}
	}
	if ev.Event == evName && bytes.Equal(ev.Data, evData) {
		return frame, nil
	}
	return formatSSEFrame(ev.Event, ev.Data), nil
}

func parseSSEFrame(frame []byte) (evName string, data []byte) {
	var dataParts [][]byte
	for line := range bytes.SplitSeq(frame, []byte("\n")) {
		trimmed := bytes.TrimSuffix(line, []byte("\r"))
		switch {
		case bytes.HasPrefix(trimmed, []byte("event:")):
			evName = strings.TrimSpace(string(bytes.TrimPrefix(trimmed, []byte("event:"))))
		case bytes.HasPrefix(trimmed, []byte("data:")):
			p := bytes.TrimPrefix(trimmed, []byte("data:"))
			p = bytes.TrimPrefix(p, []byte(" "))
			dataParts = append(dataParts, p)
		}
	}
	return evName, bytes.Join(dataParts, []byte("\n"))
}

func mergeFastPathSSEUsage(proto, eventType string, data []byte, dst *CanonicalUsage) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("[DONE]")) {
		return
	}
	switch proto {
	case ProtocolOpenAIChat:
		var chunk struct {
			Usage *openAIChatUsage `json:"usage"`
		}
		if json.Unmarshal(trimmed, &chunk) == nil && chunk.Usage != nil {
			if chunk.Usage.PromptTokens > 0 {
				dst.InputTokens = chunk.Usage.PromptTokens
			}
			if chunk.Usage.CompletionTokens > 0 {
				dst.OutputTokens = chunk.Usage.CompletionTokens
			}
			if chunk.Usage.PromptTokensDetails != nil && chunk.Usage.PromptTokensDetails.CachedTokens > 0 {
				dst.CachedTokens = chunk.Usage.PromptTokensDetails.CachedTokens
			}
		}

	case ProtocolAnthropic:
		switch eventType {
		case "message_start":
			var ev struct {
				Message struct {
					Usage *anthropicUsage `json:"usage"`
				} `json:"message"`
			}
			if json.Unmarshal(trimmed, &ev) == nil && ev.Message.Usage != nil {
				if ev.Message.Usage.InputTokens > 0 {
					dst.InputTokens = ev.Message.Usage.InputTokens
				}
				if ev.Message.Usage.OutputTokens > 0 {
					dst.OutputTokens = ev.Message.Usage.OutputTokens
				}
				if ev.Message.Usage.CacheReadInputTokens > 0 {
					dst.CachedTokens = ev.Message.Usage.CacheReadInputTokens
				}
			}
		case "message_delta":
			var ev struct {
				Usage *anthropicUsage `json:"usage"`
			}
			if json.Unmarshal(trimmed, &ev) == nil && ev.Usage != nil {
				if ev.Usage.InputTokens > 0 {
					dst.InputTokens = ev.Usage.InputTokens
				}
				if ev.Usage.OutputTokens > 0 {
					dst.OutputTokens = ev.Usage.OutputTokens
				}
				if ev.Usage.CacheReadInputTokens > 0 {
					dst.CachedTokens = ev.Usage.CacheReadInputTokens
				}
			}
		}

	case ProtocolOpenAIResponses:
		var ev struct {
			Response *struct {
				Usage *responsesUsage `json:"usage"`
			} `json:"response"`
		}
		if json.Unmarshal(trimmed, &ev) == nil && ev.Response != nil && ev.Response.Usage != nil {
			u := ev.Response.Usage
			if u.InputTokens > 0 {
				dst.InputTokens = u.InputTokens
			}
			if u.OutputTokens > 0 {
				dst.OutputTokens = u.OutputTokens
			}
			if u.InputTokensDetails != nil && u.InputTokensDetails.CachedTokens > 0 {
				dst.CachedTokens = u.InputTokensDetails.CachedTokens
			}
		}

	case ProtocolGemini:
		var chunk geminiResponse
		if json.Unmarshal(trimmed, &chunk) == nil && chunk.UsageMetadata != nil {
			u := geminiUsageToCanonical(chunk.UsageMetadata)
			if u.InputTokens > 0 {
				dst.InputTokens = u.InputTokens
			}
			if u.OutputTokens > 0 {
				dst.OutputTokens = u.OutputTokens
			}
			if u.CachedTokens > 0 {
				dst.CachedTokens = u.CachedTokens
			}
		}
	}
}

func captureResponsesFrame(frame []byte, res *streamResult) {
	evName, evData := parseSSEFrame(frame)
	if evName != "response.completed" && evName != "response.incomplete" {
		return
	}
	var wrapper struct {
		Response json.RawMessage `json:"response"`
	}
	if err := json.Unmarshal(evData, &wrapper); err != nil || len(wrapper.Response) == 0 {
		return
	}
	canon, err := decodeOpenAIResponsesResponse(wrapper.Response)
	if err != nil {
		return
	}
	res.responseID = canon.ID
	if len(canon.Content) > 0 {
		res.assistantMsg = &CanonicalMessage{
			Role:    "assistant",
			Content: canon.Content,
		}
	}
}

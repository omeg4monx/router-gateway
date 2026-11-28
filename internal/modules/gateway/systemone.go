package gateway

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
	"github.com/omeg4monx/router-gateway/internal/platform/opencode"
)

const defaultSystemOneModel = "jev-1.13-free"

// handleOpenCodeSystemOne proxies OpenCode System One (Jev) decision requests
// natively to /inference/v1/systemone over the BoringSSL/Bun uTLS client.
func (m *Module) handleOpenCodeSystemOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startTime := time.Now()
	reqID := ensureRequestID(r)

	keyInfo, ok := m.authenticateRequest(w, r, ProtocolOpenAIChat)
	if !ok {
		return
	}

	rawBody, err := io.ReadAll(io.LimitReader(r.Body, maxInboundBodyBytes+1))
	if err != nil {
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusBadRequest, "invalid_request_error", "read_error", "failed to read request body")
		return
	}
	if len(rawBody) > maxInboundBodyBytes {
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusRequestEntityTooLarge, "invalid_request_error", "body_too_large", "request body exceeds maximum size")
		return
	}

	var payload map[string]any
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusBadRequest, "invalid_request_error", "invalid_json", "request body must be valid JSON")
		return
	}
	if payload["state"] == nil {
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusBadRequest, "invalid_request_error", "missing_state", "missing required field: state")
		return
	}
	if q, ok := payload["questions"].(map[string]any); !ok || q == nil {
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusBadRequest, "invalid_request_error", "missing_questions", "missing required field: questions")
		return
	}

	requestedModel, _ := payload["model"].(string)
	requestedModel = strings.TrimSpace(requestedModel)
	if requestedModel == "" {
		requestedModel = "opencode/" + defaultSystemOneModel
	}
	upstreamModel := requestedModel
	for _, prefix := range []string{"opencode/", "oc/"} {
		if strings.HasPrefix(strings.ToLower(upstreamModel), prefix) {
			upstreamModel = strings.TrimSpace(upstreamModel[len(prefix):])
			break
		}
	}
	if upstreamModel == "" {
		upstreamModel = defaultSystemOneModel
	}
	payload["model"] = upstreamModel

	targetURL := opencode.DefaultSystemOneURL
	connID := opencode.VirtualConnectionID
	providerID := "opencode"
	if m.resolver != nil {
		if targets, resErr := m.resolver.Resolve(ctx, requestedModel, nil); resErr == nil && len(targets) > 0 {
			if targets[0].Endpoint != "" {
				targetURL = opencode.ResolveSystemOneURL(targets[0].Endpoint)
			}
			if targets[0].ConnectionID != "" {
				connID = targets[0].ConnectionID
			}
			if targets[0].ProviderID != "" {
				providerID = targets[0].ProviderID
			}
		}
	}

	outBody, err := opencode.MarshalNoEscapeHTML(payload)
	if err != nil {
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusInternalServerError, "server_error", "encode_error", err.Error())
		return
	}

	upCtx := httpclient.WithProviderID(ctx, providerID)
	upReq, err := http.NewRequestWithContext(upCtx, http.MethodPost, targetURL, bytes.NewReader(outBody))
	if err != nil {
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusBadGateway, "server_error", "request_error", err.Error())
		return
	}
	sessionID := opencode.GenerateSessionID(time.Now())
	if m.opencodeSessions != nil {
		sessionID = m.opencodeSessions.Resolve(connID, "", r.Header, rawBody)
	}
	opencode.ApplyFreeHeaders(upReq, sessionID, r.Header, false, targetURL)
	upReq.Header.Set("Accept", "application/json")

	reqClient := opencode.ResolveClient(upCtx, m.httpClient, targetURL)
	resp, doErr := reqClient.Do(upReq)
	if doErr != nil {
		if ctx.Err() != nil {
			return
		}
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusBadGateway, "server_error", "upstream_error", doErr.Error())
		return
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes))
	if readErr != nil {
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusBadGateway, "server_error", "upstream_read_error", readErr.Error())
		return
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		opencode.ResetClient(upCtx, m.httpClient, targetURL)
	}

	var inputTok, outputTok int64
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var parsed struct {
			Usage struct {
				InputTokens  int64 `json:"input_tokens"`
				OutputTokens int64 `json:"output_tokens"`
			} `json:"usage"`
		}
		if json.Unmarshal(respBody, &parsed) == nil {
			inputTok = parsed.Usage.InputTokens
			outputTok = parsed.Usage.OutputTokens
		}
	}

	m.emitUsage(ctx, contracts.UsageRecord{
		RequestID:    reqID,
		KeyID:        keyInfo.ID,
		Model:        requestedModel,
		ProviderID:   providerID,
		ConnectionID: connID,
		Protocol:     "systemone",
		InputTokens:  inputTok,
		OutputTokens: outputTok,
		Latency:      time.Since(startTime),
		TTFT:         time.Since(startTime),
		StatusCode:   resp.StatusCode,
		Stream:       false,
		Timestamp:    time.Now().UTC(),
	})

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := extractUpstreamErrorMessage(respBody, resp.StatusCode)
		writeGatewayError(w, r, ProtocolOpenAIChat, resp.StatusCode, "upstream_error", "systemone_error", errMsg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(respBody)
}

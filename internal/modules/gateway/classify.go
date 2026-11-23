package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const defaultAnthropicVersion = "2023-06-01"

type errorDisposition int

const (
	dispositionAccountError errorDisposition = iota
	dispositionClientError
)

func buildUpstreamURL(baseURL, proto, upstreamModel string, stream, isCountTokens, isCompact bool) (string, error) {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmed == "" {
		return "", fmt.Errorf("gateway: empty upstream base URL for protocol %q", proto)
	}
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Host == "" {
		return "", fmt.Errorf("gateway: invalid upstream base URL %q", baseURL)
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
	default:
		return "", fmt.Errorf("gateway: invalid upstream base URL scheme %q", baseURL)
	}

	switch proto {
	case ProtocolOpenAIChat:
		if strings.HasSuffix(trimmed, "/chat/completions") {
			return trimmed, nil
		}
		return trimmed + "/chat/completions", nil

	case ProtocolOpenAIResponses:
		if isCompact {
			switch {
			case strings.HasSuffix(trimmed, "/responses/compact"):
				return trimmed, nil
			case strings.HasSuffix(trimmed, "/responses"):
				return trimmed + "/compact", nil
			default:
				return trimmed + "/responses/compact", nil
			}
		}
		if strings.HasSuffix(trimmed, "/responses") {
			return trimmed, nil
		}
		return trimmed + "/responses", nil

	case ProtocolAnthropic:
		if isCountTokens {
			switch {
			case strings.HasSuffix(trimmed, "/messages/count_tokens"):
				return trimmed, nil
			case strings.HasSuffix(trimmed, "/v1/messages"), strings.HasSuffix(trimmed, "/messages"):
				return trimmed + "/count_tokens", nil
			case strings.HasSuffix(trimmed, "/v1"):
				return trimmed + "/messages/count_tokens", nil
			default:
				return trimmed + "/v1/messages/count_tokens", nil
			}
		}
		switch {
		case strings.HasSuffix(trimmed, "/v1/messages"), strings.HasSuffix(trimmed, "/messages"):
			return trimmed, nil
		case strings.HasSuffix(trimmed, "/v1"):
			return trimmed + "/messages", nil
		default:
			return trimmed + "/v1/messages", nil
		}

	case ProtocolGemini:
		action := "generateContent"
		if stream {
			action = "streamGenerateContent?alt=sse"
		}
		escapedModel := url.PathEscape(upstreamModel)
		if strings.HasSuffix(trimmed, "/v1beta") || strings.HasSuffix(trimmed, "/v1") {
			return trimmed + "/models/" + escapedModel + ":" + action, nil
		}
		return trimmed + "/v1beta/models/" + escapedModel + ":" + action, nil

	default:
		return "", fmt.Errorf("%w: %q", ErrUnsupportedProtocol, proto)
	}
}

func applyUpstreamHeaders(req *http.Request, creds contracts.UpstreamCredentials, stream bool) {
	for k, v := range creds.ExtraHeaders {
		cleanKey := strings.TrimSpace(k)
		if cleanKey == "" || isProtectedRequestHeader(cleanKey) {
			continue
		}
		req.Header.Set(cleanKey, v)
	}

	req.Header.Set("Content-Type", "application/json")
	if stream {
		req.Header.Set("Accept", "text/event-stream")
	} else {
		req.Header.Set("Accept", "application/json")
	}

	switch creds.Protocol {
	case ProtocolAnthropic:
		if req.Header.Get("anthropic-version") == "" {
			req.Header.Set("anthropic-version", defaultAnthropicVersion)
		}
		if creds.AuthType == "oauth" {
			req.Header.Del("x-api-key")
			if creds.AccessToken != "" {
				req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
			}
		} else {
			req.Header.Del("Authorization")
			if creds.APIKey != "" {
				req.Header.Set("x-api-key", creds.APIKey)
			}
		}

	case ProtocolOpenAIChat, ProtocolOpenAIResponses:
		token := creds.APIKey
		if creds.AuthType == "oauth" || token == "" {
			token = creds.AccessToken
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

	case ProtocolGemini:
		if creds.AuthType == "oauth" {
			req.Header.Del("x-goog-api-key")
			if creds.AccessToken != "" {
				req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
			}
		} else {
			req.Header.Del("Authorization")
			if creds.APIKey != "" {
				req.Header.Set("x-goog-api-key", creds.APIKey)
			}
		}
	}
}

func isProtectedRequestHeader(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "host", "content-type", "content-length", "transfer-encoding", "connection":
		return true
	default:
		return false
	}
}

func classifyStatus(statusCode int, body []byte) errorDisposition {
	switch statusCode {
	case http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusRequestEntityTooLarge,
		http.StatusUnsupportedMediaType,
		http.StatusUnprocessableEntity:
		return dispositionClientError
	case http.StatusUnauthorized,
		http.StatusPaymentRequired,
		http.StatusForbidden,
		http.StatusRequestTimeout,
		http.StatusTooManyRequests:
		return dispositionAccountError
	default:
		if statusCode >= 500 {
			return dispositionAccountError
		}
		if statusCode >= 400 && isClientPayloadErrorBody(body) {
			return dispositionClientError
		}
		return dispositionAccountError
	}
}

func isClientPayloadErrorBody(body []byte) bool {
	lower := strings.ToLower(string(body))
	for _, needle := range []string{
		"context_length_exceeded",
		"max_tokens",
		"invalid_request_error",
		"invalid argument",
		"malformed",
	} {
		if strings.Contains(lower, needle) {
			return true
		}
	}
	return false
}

func parseRetryAfter(headerVal string, now time.Time) time.Duration {
	trimmed := strings.TrimSpace(headerVal)
	if trimmed == "" {
		return 0
	}
	if secs, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
		if secs <= 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if ts, err := http.ParseTime(trimmed); err == nil {
		d := ts.Sub(now)
		if d > 0 {
			return d
		}
	}
	return 0
}

func extractUpstreamErrorMessage(body []byte, fallbackStatus int) string {
	var envelope struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil {
		if strings.TrimSpace(envelope.Error.Message) != "" {
			return envelope.Error.Message
		}
		if strings.TrimSpace(envelope.Message) != "" {
			return envelope.Message
		}
	}
	if text := strings.TrimSpace(string(body)); text != "" && len(text) <= 512 {
		return text
	}
	return http.StatusText(fallbackStatus)
}

package gateway

import (
	"errors"
	"net/http"
	"strings"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

func extractAPIKey(r *http.Request) string {
	if r == nil {
		return ""
	}
	if auth := strings.TrimSpace(r.Header.Get("Authorization")); len(auth) > 7 && strings.EqualFold(auth[:7], "Bearer ") {
		if tok := strings.TrimSpace(auth[7:]); tok != "" {
			return tok
		}
	}
	if key := strings.TrimSpace(r.Header.Get("x-api-key")); key != "" {
		return key
	}
	if key := strings.TrimSpace(r.Header.Get("x-goog-api-key")); key != "" {
		return key
	}
	if r.URL != nil {
		if key := strings.TrimSpace(r.URL.Query().Get("key")); key != "" {
			return key
		}
	}
	return ""
}

func (m *Module) authenticateRequest(w http.ResponseWriter, r *http.Request, proto string) (contracts.KeyInfo, bool) {
	rawKey := extractAPIKey(r)
	if rawKey == "" {
		writeGatewayError(w, r, proto, http.StatusUnauthorized, "authentication_error", "missing_api_key", "missing API key")
		return contracts.KeyInfo{}, false
	}

	info, err := m.keys.ValidateKey(r.Context(), rawKey)
	if err != nil {
		code := "invalid_api_key"
		msg := "invalid API key"
		errMsg := strings.ToLower(err.Error())
		switch {
		case strings.Contains(errMsg, "disabled"):
			code = "key_disabled"
			msg = "API key is disabled"
		case strings.Contains(errMsg, "expired"):
			code = "key_expired"
			msg = "API key has expired"
		}
		writeGatewayError(w, r, proto, http.StatusUnauthorized, "authentication_error", code, msg)
		return contracts.KeyInfo{}, false
	}
	if info.Disabled {
		writeGatewayError(w, r, proto, http.StatusUnauthorized, "authentication_error", "key_disabled", "API key is disabled")
		return contracts.KeyInfo{}, false
	}
	return info, true
}

func writeGatewayError(w http.ResponseWriter, _ *http.Request, proto string, status int, errType, code, msg string) {
	if msg == "" {
		msg = http.StatusText(status)
	}
	switch proto {
	case ProtocolAnthropic:
		httpserver.WriteJSON(w, status, map[string]any{
			"type": "error",
			"error": map[string]any{
				"type":    anthropicErrorType(status, errType),
				"message": msg,
			},
		})
	case ProtocolGemini:
		httpserver.WriteJSON(w, status, map[string]any{
			"error": map[string]any{
				"code":    status,
				"message": msg,
				"status":  geminiErrorStatus(status),
			},
		})
	default:
		httpserver.WriteJSON(w, status, map[string]any{
			"error": map[string]any{
				"message": msg,
				"type":    openAIErrorType(status, errType),
				"code":    code,
			},
		})
	}
}

func writeResolutionError(w http.ResponseWriter, r *http.Request, proto string, err error) {
	switch {
	case errors.Is(err, contracts.ErrModelForbidden):
		writeGatewayError(w, r, proto, http.StatusForbidden, "permission_error", "model_forbidden", "model is not permitted for this API key")
	case errors.Is(err, contracts.ErrModelNotFound):
		writeGatewayError(w, r, proto, http.StatusNotFound, "invalid_request_error", "model_not_found", "requested model was not found")
	case errors.Is(err, contracts.ErrAllAccountsCoolingDown):
		writeGatewayError(w, r, proto, http.StatusServiceUnavailable, "server_error", "all_accounts_cooling_down", "all upstream accounts for this model are temporarily cooling down")
	default:
		writeGatewayError(w, r, proto, http.StatusInternalServerError, "server_error", "resolution_failed", "failed to resolve upstream model")
	}
}

func openAIErrorType(status int, fallback string) string {
	switch status {
	case http.StatusUnauthorized:
		return "authentication_error"
	case http.StatusForbidden:
		return "permission_error"
	case http.StatusTooManyRequests:
		return "rate_limit_error"
	default:
		if status >= 500 {
			return "server_error"
		}
		if fallback != "" {
			return fallback
		}
		return "invalid_request_error"
	}
}

func anthropicErrorType(status int, fallback string) string {
	switch status {
	case http.StatusUnauthorized:
		return "authentication_error"
	case http.StatusForbidden:
		return "permission_error"
	case http.StatusNotFound:
		return "not_found_error"
	case http.StatusTooManyRequests:
		return "rate_limit_error"
	case http.StatusServiceUnavailable, http.StatusBadGateway, http.StatusGatewayTimeout:
		return "overloaded_error"
	default:
		if status >= 500 {
			return "api_error"
		}
		if fallback != "" {
			return fallback
		}
		return "invalid_request_error"
	}
}

func geminiErrorStatus(status int) string {
	switch status {
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return "INVALID_ARGUMENT"
	case http.StatusUnauthorized:
		return "UNAUTHENTICATED"
	case http.StatusForbidden:
		return "PERMISSION_DENIED"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusTooManyRequests:
		return "RESOURCE_EXHAUSTED"
	case http.StatusServiceUnavailable, http.StatusBadGateway:
		return "UNAVAILABLE"
	case http.StatusGatewayTimeout:
		return "DEADLINE_EXCEEDED"
	default:
		return "INTERNAL"
	}
}

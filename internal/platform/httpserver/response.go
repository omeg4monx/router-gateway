// Package httpserver provides the chi-based HTTP server with split API and
// streaming gateway middleware stacks, structured JSON error responses, and
// SPA static asset serving.
package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// ErrorDetail carries structured machine-readable error metadata in JSON error
// responses.
type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// ErrorEnvelope is the canonical JSON envelope returned for API and gateway
// errors.
type ErrorEnvelope struct {
	Error ErrorDetail `json:"error"`
}

// WriteJSON marshals v as JSON and writes it with status and an explicit
// application/json Content-Type header.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"code":"internal_error","message":"failed to encode response"}}`))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// WriteError writes a structured JSON ErrorEnvelope response with status,
// machine-readable code, human-readable msg, and the request's RequestID when
// present in context.
func WriteError(w http.ResponseWriter, r *http.Request, status int, code, msg string) {
	var reqID string
	if r != nil {
		reqID = middleware.GetReqID(r.Context())
	}
	WriteJSON(w, status, ErrorEnvelope{
		Error: ErrorDetail{
			Code:      code,
			Message:   msg,
			RequestID: reqID,
		},
	})
}

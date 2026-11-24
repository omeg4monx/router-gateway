package providers

import (
	"net/http"
	"strings"

	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

const (
	providerIDCodex   = "codex"
	codexCallbackPort = 1455
	codexCallbackPath = "/auth/callback"
)

// codexOAuthRuntime is a thin compatibility wrapper over oauthLoopbackProxy
// bound to Codex's fixed port 1455 /auth/callback.
type codexOAuthRuntime struct {
	proxy *oauthLoopbackProxy
}

func newCodexOAuthRuntime() *codexOAuthRuntime {
	return &codexOAuthRuntime{
		proxy: newOAuthLoopbackProxy(),
	}
}

func (rt *codexOAuthRuntime) registerSession(state string) {
	rt.proxy.registerSession(providerIDCodex, codexCallbackPath, state)
}

func (rt *codexOAuthRuntime) start(state string) {
	// Port 1455 busy — gracefully fall back to manual callback paste.
	_, _ = rt.proxy.start(providerIDCodex, codexCallbackPort, codexCallbackPath, state)
}

func (rt *codexOAuthRuntime) stop() {
	rt.proxy.stop(providerIDCodex)
}

func (rt *codexOAuthRuntime) status(state string) (string, string) {
	return rt.proxy.status(providerIDCodex, state)
}

func (rt *codexOAuthRuntime) consumeCode(state string) (string, bool) {
	return rt.proxy.consumeCode(providerIDCodex, state)
}

func (rt *codexOAuthRuntime) handleCallback(w http.ResponseWriter, r *http.Request) {
	rt.proxy.handleCallback(providerIDCodex, w, r)
}

func (m *Module) handleLoopbackOAuthStatus(w http.ResponseWriter, r *http.Request, providerID string) {
	state := strings.TrimSpace(r.URL.Query().Get("state"))
	if state == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "state query parameter is required")
		return
	}
	status, errMsg := m.oauthProxy.status(providerID, state)
	httpserver.WriteJSON(w, http.StatusOK, oauthStatusResponse{
		Status: status,
		Error:  errMsg,
	})
}

package providers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCodexLoopbackRuntimeCaptureAndStatus(t *testing.T) {
	rt := newCodexOAuthRuntime()
	defer rt.stop()

	rt.registerSession("codex-state-1")

	if status, errMsg := rt.status("codex-state-1"); status != "waiting" || errMsg != "" {
		t.Fatalf("expected waiting, got status=%q err=%q", status, errMsg)
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/callback?code=codex-code-xyz&state=codex-state-1", http.NoBody)
	rec := httptest.NewRecorder()
	rt.handleCallback(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"oauth_callback"`) {
		t.Fatalf("expected callback HTML to contain oauth_callback bridge")
	}

	if status, errMsg := rt.status("codex-state-1"); status != "done" || errMsg != "" {
		t.Fatalf("expected done, got status=%q err=%q", status, errMsg)
	}

	code, ok := rt.consumeCode("codex-state-1")
	if !ok || code != "codex-code-xyz" {
		t.Fatalf("expected captured code codex-code-xyz, got ok=%v code=%q", ok, code)
	}
}

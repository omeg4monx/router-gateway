package translator_test

import (
	"bufio"
	"context"
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/modules/gateway"
	"github.com/omeg4monx/router-gateway/internal/modules/translator"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
	"github.com/omeg4monx/router-gateway/internal/platform/logging"
)

var _ lifecycle.Module = (*translator.Module)(nil)

type stubDep struct{ name string }

func (s stubDep) Name() string                         { return s.name }
func (s stubDep) Dependencies() []string               { return nil }
func (s stubDep) Migrations() fs.FS                    { return nil }
func (s stubDep) RegisterRoutes(lifecycle.RouteMounts) {}
func (s stubDep) Start(context.Context) error          { return nil }
func (s stubDep) Stop(context.Context) error           { return nil }

type harness struct {
	server *httpserver.Server
	logger *slog.Logger
	ring   *logging.RingBuffer
}

func newHarness(t *testing.T) *harness {
	t.Helper()

	logger, ring := logging.New("debug", nil, 100)
	gw := gateway.New(gateway.Options{Logger: logger})
	mod := translator.New(translator.Options{
		Translator: gw,
		LogRing:    ring,
	})

	cfg := config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: 15 * time.Second}
	srv := httpserver.New(httpserver.Options{Config: cfg})

	reg := lifecycle.NewRegistry(nil)
	for _, m := range []lifecycle.Module{stubDep{name: "system"}, stubDep{name: "gateway"}, mod} {
		if err := reg.Register(m); err != nil {
			t.Fatalf("register %s: %v", m.Name(), err)
		}
	}
	if err := reg.MountAll(srv.Mounts()); err != nil {
		t.Fatalf("MountAll: %v", err)
	}
	if err := reg.StartAll(context.Background()); err != nil {
		t.Fatalf("StartAll: %v", err)
	}
	t.Cleanup(func() { _ = reg.StopAll(context.Background()) })

	return &harness{server: srv, logger: logger, ring: ring}
}

func (h *harness) doReq(method, path, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, http.NoBody)
	}
	rec := httptest.NewRecorder()
	h.server.Handler().ServeHTTP(rec, req)
	return rec
}

func TestTranslatorTranslateRequestAndResponseAllPairs(t *testing.T) {
	t.Parallel()

	h := newHarness(t)

	// GET /api/translator/protocols
	protoRec := h.doReq(http.MethodGet, "/api/translator/protocols", "")
	if protoRec.Code != http.StatusOK {
		t.Fatalf("GET /api/translator/protocols = %d: %s", protoRec.Code, protoRec.Body.String())
	}
	var protoResp struct {
		Protocols []string `json:"protocols"`
		Modes     []string `json:"modes"`
	}
	if err := json.Unmarshal(protoRec.Body.Bytes(), &protoResp); err != nil {
		t.Fatalf("unmarshal protocols: %v", err)
	}
	if len(protoResp.Protocols) != 4 || len(protoResp.Modes) != 2 {
		t.Fatalf("unexpected protocols response: %+v", protoResp)
	}

	// OpenAI Chat -> Anthropic request translation
	reqCases := []struct {
		name    string
		body    string
		wantKey string
	}{
		{
			name:    "openai-chat to anthropic request",
			body:    `{"mode":"request","from_protocol":"openai-chat","to_protocol":"anthropic","payload":{"model":"claude-sonnet-4","messages":[{"role":"system","content":"Be terse."},{"role":"user","content":"Hi"}]}}`,
			wantKey: "messages",
		},
		{
			name:    "openai-chat to gemini request",
			body:    `{"mode":"request","from_protocol":"openai-chat","to_protocol":"gemini","payload":{"model":"gemini-2.5-pro","messages":[{"role":"user","content":"Hello Gemini"}]}}`,
			wantKey: "contents",
		},
		{
			name:    "anthropic to openai-responses request",
			body:    `{"mode":"request","from_protocol":"anthropic","to_protocol":"openai-responses","payload":{"model":"gpt-4.1","max_tokens":128,"messages":[{"role":"user","content":"Ping"}]}}`,
			wantKey: "input",
		},
		{
			name:    "anthropic to openai-chat response",
			body:    `{"mode":"response","from_protocol":"anthropic","to_protocol":"openai-chat","payload":{"id":"msg_1","type":"message","role":"assistant","model":"claude-sonnet-4","content":[{"type":"text","text":"Pong"}],"stop_reason":"end_turn","usage":{"input_tokens":5,"output_tokens":3}}}`,
			wantKey: "choices",
		},
	}

	for _, tc := range reqCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := h.doReq(http.MethodPost, "/api/translator/translate", tc.body)
			if rec.Code != http.StatusOK {
				t.Fatalf("POST /api/translator/translate = %d: %s", rec.Code, rec.Body.String())
			}
			var out struct {
				Mode       string          `json:"mode"`
				FastPath   bool            `json:"fast_path"`
				DurationUS int64           `json:"duration_us"`
				Output     json.RawMessage `json:"output"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
				t.Fatalf("unmarshal translate resp: %v", err)
			}
			var obj map[string]any
			if err := json.Unmarshal(out.Output, &obj); err != nil {
				t.Fatalf("output is not JSON object: %v (%s)", err, string(out.Output))
			}
			if _, ok := obj[tc.wantKey]; !ok {
				t.Fatalf("translated output missing key %q: %s", tc.wantKey, string(out.Output))
			}
		})
	}

	// Invalid protocol, invalid mode, and non-object payload return 400.
	for _, bad := range []string{
		`{"mode":"request","from_protocol":"unknown","to_protocol":"anthropic","payload":{}}`,
		`{"mode":"invalid","from_protocol":"openai-chat","to_protocol":"anthropic","payload":{}}`,
		`{"mode":"request","from_protocol":"openai-chat","to_protocol":"anthropic","payload":[1,2,3]}`,
	} {
		rec := h.doReq(http.MethodPost, "/api/translator/translate", bad)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for %s, got %d: %s", bad, rec.Code, rec.Body.String())
		}
	}
}

func TestTranslatorLogsSnapshotAndLiveSSEStream(t *testing.T) {
	t.Parallel()

	h := newHarness(t)
	h.logger.Info("boot-info", "component", "test")
	h.logger.Warn("high-latency", "ms", 420)

	// Filtered snapshot: level=WARN returns only WARN+ entries.
	snapRec := h.doReq(http.MethodGet, "/api/translator/logs?level=WARN&limit=50", "")
	if snapRec.Code != http.StatusOK {
		t.Fatalf("GET /api/translator/logs = %d: %s", snapRec.Code, snapRec.Body.String())
	}
	var snap struct {
		Entries []logging.LogEntry `json:"entries"`
		Total   int                `json:"total"`
	}
	if err := json.Unmarshal(snapRec.Body.Bytes(), &snap); err != nil {
		t.Fatalf("unmarshal logs: %v", err)
	}
	if snap.Total != 1 || len(snap.Entries) != 1 || snap.Entries[0].Message != "high-latency" {
		t.Fatalf("unexpected filtered logs: %+v", snap)
	}

	// Live SSE stream via real httptest.Server so ResponseController.Flush streams chunks.
	ts := httptest.NewServer(h.server.Handler())
	t.Cleanup(ts.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/api/translator/logs/stream?tail=10", http.NoBody)
	if err != nil {
		t.Fatalf("new sse req: %v", err)
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("connect sse: %v", err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want text/event-stream", ct)
	}

	scanner := bufio.NewScanner(resp.Body)
	sawReplay := false
	sawLive := false
	liveEmitted := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			if strings.Contains(line, "high-latency") {
				sawReplay = true
				if !liveEmitted {
					liveEmitted = true
					h.logger.Error("live-stream-event", "req", "r1")
				}
			}
			if strings.Contains(line, "live-stream-event") {
				sawLive = true
				cancel()
				break
			}
		}
	}

	if !sawReplay || !sawLive {
		t.Fatalf("SSE stream missing events: sawReplay=%v sawLive=%v", sawReplay, sawLive)
	}
}

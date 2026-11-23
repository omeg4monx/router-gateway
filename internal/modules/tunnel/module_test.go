package tunnel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type fakeHandle struct {
	pid    int
	out    *io.PipeReader
	pw     *io.PipeWriter
	done   chan error
	killed atomic.Bool
	once   sync.Once
}

func (f *fakeHandle) PID() int          { return f.pid }
func (f *fakeHandle) Output() io.Reader { return f.out }
func (f *fakeHandle) Wait() error       { return <-f.done }
func (f *fakeHandle) Kill() error {
	f.killed.Store(true)
	f.once.Do(func() {
		_ = f.pw.Close()
		f.done <- nil
		close(f.done)
	})
	return nil
}

type fakeRunner struct {
	mu           sync.Mutex
	binaries     map[string]string
	lastHandle   *fakeHandle
	runCalls     [][]string
	runOutputs   map[string]string
	startLogLine string
}

func (r *fakeRunner) LookPath(name string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.binaries[name]; ok && p != "" {
		return p, nil
	}
	return "", errors.New("not found")
}

func (r *fakeRunner) Start(_ context.Context, name string, args ...string) (ProcessHandle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runCalls = append(r.runCalls, append([]string{name}, args...))
	pr, pw := io.Pipe()
	h := &fakeHandle{
		pid:  4242,
		out:  pr,
		pw:   pw,
		done: make(chan error, 1),
	}
	r.lastHandle = h
	line := r.startLogLine
	go func() {
		if line != "" {
			_, _ = io.WriteString(pw, line+"\n")
		}
	}()
	return h, nil
}

func (r *fakeRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	call := append([]string{name}, args...)
	r.runCalls = append(r.runCalls, call)
	key := strings.Join(call, " ")
	for prefix, out := range r.runOutputs {
		if strings.HasPrefix(key, prefix) {
			return []byte(out), nil
		}
	}
	return []byte("ok\n"), nil
}

func setupTunnelTest(t *testing.T, runner ProcessRunner) (*Module, *chi.Mux) {
	t.Helper()
	database, err := db.Open(filepath.Join(t.TempDir(), "tunnel.sqlite"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	mod := New(Options{
		DB:     database,
		Config: config.Config{Host: "127.0.0.1", Port: 20128},
		Bus:    eventbus.New(),
		Runner: runner,
	})
	if err := database.MigrateModule(context.Background(), mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() { _ = mod.Stop(context.Background()) })

	r := chi.NewRouter()
	mod.RegisterRoutes(lifecycle.RouteMounts{AdminAPI: r})
	return mod, r
}

func doTunnelReq(t *testing.T, h http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var r io.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		r = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, r)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestTunnelConfigValidationAndPersistence(t *testing.T) {
	runner := &fakeRunner{binaries: map[string]string{"cloudflared": "/usr/local/bin/cloudflared"}}
	_, r := setupTunnelTest(t, runner)

	// Default status.
	rec := doTunnelReq(t, r, http.MethodGet, "/tunnel/status", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 status, got %d: %s", rec.Code, rec.Body.String())
	}
	var st Status
	if err := json.Unmarshal(rec.Body.Bytes(), &st); err != nil {
		t.Fatalf("decode status: %v", err)
	}
	if st.Status != StateStopped || st.Provider != ProviderCloudflare || st.LocalURL != "http://127.0.0.1:20128" {
		t.Fatalf("unexpected default status: %+v", st)
	}
	if !st.Binaries.Cloudflared.Available || st.Binaries.Tailscale.Available {
		t.Fatalf("unexpected binary detection: %+v", st.Binaries)
	}

	// Reject non-loopback local_url.
	for _, bad := range []string{"http://10.0.0.5:8080", "http://169.254.169.254", "https://example.com"} {
		badRec := doTunnelReq(t, r, http.MethodPatch, "/tunnel/config", map[string]any{
			"local_url": bad,
		})
		if badRec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for non-loopback %q, got %d: %s", bad, badRec.Code, badRec.Body.String())
		}
	}

	// Persist valid config.
	patchRec := doTunnelReq(t, r, http.MethodPatch, "/tunnel/config", map[string]any{
		"provider":       "tailscale",
		"local_url":      "http://localhost:20128",
		"tailscale_mode": "serve",
		"auto_start":     true,
	})
	if patchRec.Code != http.StatusOK {
		t.Fatalf("expected 200 patch, got %d: %s", patchRec.Code, patchRec.Body.String())
	}
	var updated Status
	_ = json.Unmarshal(patchRec.Body.Bytes(), &updated)
	if updated.Provider != ProviderTailscale || updated.TailscaleMode != "serve" || !updated.AutoStart {
		t.Fatalf("unexpected patched status: %+v", updated)
	}
}

func TestTunnelCloudflareLifecycleAndURLCapture(t *testing.T) {
	runner := &fakeRunner{
		binaries:     map[string]string{"cloudflared": "/usr/local/bin/cloudflared"},
		startLogLine: "2026-11-23T10:00:00Z INF |  https://ember-test-1234.trycloudflare.com  |",
	}
	_, r := setupTunnelTest(t, runner)

	startRec := doTunnelReq(t, r, http.MethodPost, "/tunnel/start", map[string]any{
		"provider": "cloudflare",
	})
	if startRec.Code != http.StatusOK {
		t.Fatalf("expected 200 start, got %d: %s", startRec.Code, startRec.Body.String())
	}

	// Poll briefly for async URL capture from child pipe.
	deadline := time.Now().Add(2 * time.Second)
	var st Status
	for time.Now().Before(deadline) {
		rec := doTunnelReq(t, r, http.MethodGet, "/tunnel/status", nil)
		_ = json.Unmarshal(rec.Body.Bytes(), &st)
		if st.Status == StateRunning && st.PublicURL != "" {
			break
		}
		time.Sleep(15 * time.Millisecond)
	}
	if st.Status != StateRunning || st.PublicURL != "https://ember-test-1234.trycloudflare.com" {
		t.Fatalf("expected running with captured trycloudflare URL, got %+v", st)
	}

	// Double start returns 409.
	dupRec := doTunnelReq(t, r, http.MethodPost, "/tunnel/start", nil)
	if dupRec.Code != http.StatusConflict {
		t.Fatalf("expected 409 double start, got %d: %s", dupRec.Code, dupRec.Body.String())
	}

	// Stop cleans up process.
	stopRec := doTunnelReq(t, r, http.MethodPost, "/tunnel/stop", nil)
	if stopRec.Code != http.StatusOK {
		t.Fatalf("expected 200 stop, got %d: %s", stopRec.Code, stopRec.Body.String())
	}
	var stopped Status
	_ = json.Unmarshal(stopRec.Body.Bytes(), &stopped)
	if stopped.Status != StateStopped || stopped.PublicURL != "" {
		t.Fatalf("expected stopped with empty public_url, got %+v", stopped)
	}
	if !runner.lastHandle.killed.Load() {
		t.Fatalf("expected child process handle to be killed on stop")
	}

	// Events recorded.
	evRec := doTunnelReq(t, r, http.MethodGet, "/tunnel/events", nil)
	if evRec.Code != http.StatusOK {
		t.Fatalf("expected 200 events, got %d", evRec.Code)
	}
	if !strings.Contains(evRec.Body.String(), "ember-test-1234.trycloudflare.com") {
		t.Fatalf("expected captured URL in tunnel events, got %s", evRec.Body.String())
	}
}

func TestTunnelTailscaleLifecycleAndMissingBinary(t *testing.T) {
	// Missing binary returns 422.
	emptyRunner := &fakeRunner{binaries: map[string]string{}}
	_, rMissing := setupTunnelTest(t, emptyRunner)
	missRec := doTunnelReq(t, rMissing, http.MethodPost, "/tunnel/start", map[string]any{
		"provider": "cloudflare",
	})
	if missRec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 when binary missing, got %d: %s", missRec.Code, missRec.Body.String())
	}

	// Tailscale start + stop.
	tsRunner := &fakeRunner{
		binaries: map[string]string{"tailscale": "/usr/bin/tailscale"},
		runOutputs: map[string]string{
			"tailscale funnel --bg": "Available on the internet:\nhttps://router.tailnet.ts.net/\n",
		},
	}
	_, rTS := setupTunnelTest(t, tsRunner)
	startRec := doTunnelReq(t, rTS, http.MethodPost, "/tunnel/start", map[string]any{
		"provider":       "tailscale",
		"tailscale_mode": "funnel",
	})
	if startRec.Code != http.StatusOK {
		t.Fatalf("expected 200 tailscale start, got %d: %s", startRec.Code, startRec.Body.String())
	}
	var st Status
	_ = json.Unmarshal(startRec.Body.Bytes(), &st)
	if st.Status != StateRunning || st.PublicURL != "https://router.tailnet.ts.net/" {
		t.Fatalf("expected tailscale running with URL, got %+v", st)
	}

	stopRec := doTunnelReq(t, rTS, http.MethodPost, "/tunnel/stop", nil)
	if stopRec.Code != http.StatusOK {
		t.Fatalf("expected 200 tailscale stop, got %d: %s", stopRec.Code, stopRec.Body.String())
	}
}

func TestTunnelStopTerminatesRunningProcessOnShutdown(t *testing.T) {
	runner := &fakeRunner{
		binaries:     map[string]string{"cloudflared": "/usr/local/bin/cloudflared"},
		startLogLine: "INF https://shutdown-test.trycloudflare.com",
	}
	mod, r := setupTunnelTest(t, runner)
	rec := doTunnelReq(t, r, http.MethodPost, "/tunnel/start", map[string]any{"provider": "cloudflare"})
	if rec.Code != http.StatusOK {
		t.Fatalf("start failed: %d", rec.Code)
	}
	if err := mod.Stop(context.Background()); err != nil {
		t.Fatalf("mod.Stop: %v", err)
	}
	if runner.lastHandle == nil || !runner.lastHandle.killed.Load() {
		t.Fatalf("expected mod.Stop to kill active tunnel process")
	}
}

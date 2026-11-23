package mitm

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

func setupTestMitm(t *testing.T) (*Module, http.Handler, string) {
	t.Helper()
	dir := t.TempDir()
	database, err := db.Open(filepath.Join(dir, "mitm.sqlite"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	cfg := config.Config{
		Host:    "127.0.0.1",
		Port:    20128,
		DataDir: dir,
	}
	mod := New(Options{
		DB:     database,
		Config: cfg,
		Bus:    eventbus.New(),
	})

	ctx := context.Background()
	if err := database.MigrateModule(ctx, mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	if err := mod.Start(ctx); err != nil {
		t.Fatalf("start module: %v", err)
	}
	t.Cleanup(func() { _ = mod.Stop(context.Background()) })

	r := chi.NewRouter()
	r.Route("/api", func(api chi.Router) {
		mod.RegisterRoutes(lifecycle.RouteMounts{
			PublicAPI: api,
			AdminAPI:  api,
			Gateway:   r,
		})
	})

	return mod, r, dir
}

func doMitmJSON(t *testing.T, handler http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		buf = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestMitmModuleMetadataAndStatus(t *testing.T) {
	mod, handler, _ := setupTestMitm(t)
	if mod.Name() != "mitm" {
		t.Fatalf("expected name mitm, got %q", mod.Name())
	}

	rec := doMitmJSON(t, handler, http.MethodGet, "/api/mitm/status", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/mitm/status: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var status Status
	if err := json.Unmarshal(rec.Body.Bytes(), &status); err != nil {
		t.Fatalf("decode status: %v", err)
	}
	if status.Status != StateStopped {
		t.Fatalf("expected initial status stopped, got %q", status.Status)
	}
	if status.CA.Exists {
		t.Fatalf("expected CA.Exists=false before generation")
	}
	if len(status.Tools) != 4 {
		t.Fatalf("expected 4 IDE tools in catalog, got %d", len(status.Tools))
	}
	if status.InstallGuide.HostsSnippet == "" || status.InstallGuide.CACommand == "" {
		t.Fatalf("expected non-empty install guide, got %+v", status.InstallGuide)
	}
}

func TestMitmConfigValidationLoopbackOnly(t *testing.T) {
	_, handler, _ := setupTestMitm(t)

	// Reject non-loopback listen_addr
	rec := doMitmJSON(t, handler, http.MethodPatch, "/api/mitm/config", map[string]any{
		"listen_addr": "0.0.0.0:8443",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-loopback listen_addr, got %d: %s", rec.Code, rec.Body.String())
	}

	// Reject non-loopback target_gateway_url
	rec = doMitmJSON(t, handler, http.MethodPatch, "/api/mitm/config", map[string]any{
		"target_gateway_url": "http://10.0.0.5:20128",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-loopback target_gateway_url, got %d: %s", rec.Code, rec.Body.String())
	}

	// Reject unknown tool
	rec = doMitmJSON(t, handler, http.MethodPatch, "/api/mitm/config", map[string]any{
		"enabled_tools": []string{"unknown-ide"},
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for unknown tool, got %d: %s", rec.Code, rec.Body.String())
	}

	// Reject out-of-range or non-numeric port
	for _, bad := range []string{"127.0.0.1:70000", "127.0.0.1:abc"} {
		rec = doMitmJSON(t, handler, http.MethodPatch, "/api/mitm/config", map[string]any{
			"listen_addr": bad,
		})
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for invalid port %q, got %d: %s", bad, rec.Code, rec.Body.String())
		}
	}

	// Verify bracketed IPv6 defaultGatewayURL does not double-bracket
	m6 := New(Options{Config: config.Config{Host: "[::1]", Port: 20128, DataDir: t.TempDir()}})
	if got := m6.defaultGatewayURL(); got != "http://[::1]:20128" {
		t.Fatalf("expected http://[::1]:20128, got %q", got)
	}

	// Accept valid loopback config
	rec = doMitmJSON(t, handler, http.MethodPatch, "/api/mitm/config", map[string]any{
		"listen_addr":        "127.0.0.1:0",
		"target_gateway_url": "http://127.0.0.1:20128",
		"enabled_tools":      []string{"copilot", "cursor"},
		"auto_start":         false,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid loopback config, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestMitmCAGenerateDownloadAndKeyPermissions(t *testing.T) {
	_, handler, dir := setupTestMitm(t)

	// Download before generation returns 404
	rec := doMitmJSON(t, handler, http.MethodGet, "/api/mitm/ca/download", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 before CA exists, got %d", rec.Code)
	}

	// Generate Root CA
	rec = doMitmJSON(t, handler, http.MethodPost, "/api/mitm/ca/generate", map[string]any{
		"force": false,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/mitm/ca/generate: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var firstCA CAInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &firstCA); err != nil {
		t.Fatalf("decode CAInfo: %v", err)
	}
	if !firstCA.Exists || firstCA.Fingerprint == "" {
		t.Fatalf("expected generated CA metadata, got %+v", firstCA)
	}

	// Check private key file permissions on disk (0600)
	keyPath := filepath.Join(dir, "mitm", "rootCA.key")
	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("stat rootCA.key: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected rootCA.key mode 0600, got %o", info.Mode().Perm())
	}

	// Download CA PEM and verify it never leaks private key
	rec = doMitmJSON(t, handler, http.MethodGet, "/api/mitm/ca/download", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/mitm/ca/download: expected 200, got %d", rec.Code)
	}
	pemText := rec.Body.String()
	if !strings.Contains(pemText, "BEGIN CERTIFICATE") {
		t.Fatalf("expected PEM certificate, got: %s", pemText)
	}
	if strings.Contains(pemText, "PRIVATE KEY") {
		t.Fatalf("CA download must NEVER contain private key material")
	}

	// Rotate with force=true changes fingerprint
	rec = doMitmJSON(t, handler, http.MethodPost, "/api/mitm/ca/generate", map[string]any{
		"force": true,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("force rotate CA: expected 200, got %d", rec.Code)
	}
	var rotatedCA CAInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &rotatedCA); err != nil {
		t.Fatalf("decode rotated CA: %v", err)
	}
	if rotatedCA.Fingerprint == firstCA.Fingerprint {
		t.Fatalf("expected rotated CA to have a new fingerprint")
	}
}

func TestMitmAliasesCRUD(t *testing.T) {
	_, handler, _ := setupTestMitm(t)

	rec := doMitmJSON(t, handler, http.MethodPost, "/api/mitm/aliases", map[string]any{
		"tool":         "copilot",
		"source_model": "gpt-4o",
		"target_model": "claude-sonnet-4",
		"enabled":      true,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/mitm/aliases: expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var created Alias
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode alias: %v", err)
	}

	// Duplicate returns 409
	dup := doMitmJSON(t, handler, http.MethodPost, "/api/mitm/aliases", map[string]any{
		"tool":         "copilot",
		"source_model": "gpt-4o",
		"target_model": "other-model",
		"enabled":      true,
	})
	if dup.Code != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate alias, got %d: %s", dup.Code, dup.Body.String())
	}

	// List
	listRec := doMitmJSON(t, handler, http.MethodGet, "/api/mitm/aliases", nil)
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /api/mitm/aliases: expected 200, got %d", listRec.Code)
	}
	var listResp AliasListResponse
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("decode aliases: %v", err)
	}
	if len(listResp.Aliases) != 1 || listResp.Aliases[0].ID != created.ID {
		t.Fatalf("expected 1 alias matching %q, got %+v", created.ID, listResp.Aliases)
	}

	// Delete
	delRec := doMitmJSON(t, handler, http.MethodDelete, "/api/mitm/aliases/"+created.ID, nil)
	if delRec.Code != http.StatusOK {
		t.Fatalf("DELETE /api/mitm/aliases: expected 200, got %d", delRec.Code)
	}
}

func TestMitmLiveTLSProxyInterceptAndAliasRewrite(t *testing.T) {
	_, handler, _ := setupTestMitm(t)

	// Mock local gateway target
	var receivedModel string
	var receivedPath string
	mockGateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		var payload map[string]any
		_ = json.NewDecoder(r.Body).Decode(&payload)
		if m, ok := payload["model"].(string); ok {
			receivedModel = m
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"chatcmpl-mitm","object":"chat.completion"}`))
	}))
	t.Cleanup(mockGateway.Close)

	// Configure MITM on ephemeral loopback port targeting mockGateway
	rec := doMitmJSON(t, handler, http.MethodPatch, "/api/mitm/config", map[string]any{
		"listen_addr":        "127.0.0.1:0",
		"target_gateway_url": mockGateway.URL,
		"enabled_tools":      []string{"copilot"},
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("patch config: %d %s", rec.Code, rec.Body.String())
	}

	// Create alias for copilot: gpt-4o -> claude-sonnet-4
	rec = doMitmJSON(t, handler, http.MethodPost, "/api/mitm/aliases", map[string]any{
		"tool":         "copilot",
		"source_model": "gpt-4o",
		"target_model": "claude-sonnet-4",
		"enabled":      true,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create alias: %d %s", rec.Code, rec.Body.String())
	}

	// Start MITM proxy (auto-generates CA if missing)
	rec = doMitmJSON(t, handler, http.MethodPost, "/api/mitm/start", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("start mitm: %d %s", rec.Code, rec.Body.String())
	}
	var st Status
	if err := json.Unmarshal(rec.Body.Bytes(), &st); err != nil {
		t.Fatalf("decode start status: %v", err)
	}
	if st.Status != StateRunning || st.BoundAddr == "" {
		t.Fatalf("expected running MITM with bound_addr, got %+v", st)
	}

	// Download Root CA PEM and build TLS client trusting only that CA
	caRec := doMitmJSON(t, handler, http.MethodGet, "/api/mitm/ca/download", nil)
	if caRec.Code != http.StatusOK {
		t.Fatalf("download CA: %d", caRec.Code)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caRec.Body.Bytes()) {
		t.Fatalf("failed to parse downloaded CA PEM")
	}

	tlsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    pool,
				ServerName: "api.individual.githubcopilot.com",
				MinVersion: tls.VersionTLS12,
			},
		},
	}
	t.Cleanup(tlsClient.CloseIdleConnections)

	req, err := http.NewRequest(http.MethodPost, "https://"+st.BoundAddr+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}`))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Host = "api.individual.githubcopilot.com"
	req.Header.Set("Content-Type", "application/json")

	resp, err := tlsClient.Do(req)
	if err != nil {
		t.Fatalf("TLS request through MITM proxy failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 from MITM proxy, got %d: %s", resp.StatusCode, string(body))
	}
	if receivedPath != "/v1/chat/completions" {
		t.Fatalf("expected forwarded path /v1/chat/completions, got %q", receivedPath)
	}
	if receivedModel != "claude-sonnet-4" {
		t.Fatalf("expected rewritten model claude-sonnet-4, got %q", receivedModel)
	}

	// Verify event logged
	evRec := doMitmJSON(t, handler, http.MethodGet, "/api/mitm/events", nil)
	if evRec.Code != http.StatusOK {
		t.Fatalf("GET /api/mitm/events: %d", evRec.Code)
	}
	var evResp EventListResponse
	if err := json.Unmarshal(evRec.Body.Bytes(), &evResp); err != nil {
		t.Fatalf("decode events: %v", err)
	}
	if len(evResp.Events) != 1 {
		t.Fatalf("expected 1 intercepted event, got %d", len(evResp.Events))
	}
	if evResp.Events[0].Tool != "copilot" || evResp.Events[0].SourceModel != "gpt-4o" || evResp.Events[0].TargetModel != "claude-sonnet-4" {
		t.Fatalf("unexpected event record: %+v", evResp.Events[0])
	}

	// Clear events
	clrRec := doMitmJSON(t, handler, http.MethodDelete, "/api/mitm/events", nil)
	if clrRec.Code != http.StatusOK {
		t.Fatalf("DELETE /api/mitm/events: %d", clrRec.Code)
	}

	// Stop MITM proxy
	stopRec := doMitmJSON(t, handler, http.MethodPost, "/api/mitm/stop", nil)
	if stopRec.Code != http.StatusOK {
		t.Fatalf("POST /api/mitm/stop: %d", stopRec.Code)
	}
	var stopped Status
	if err := json.Unmarshal(stopRec.Body.Bytes(), &stopped); err != nil {
		t.Fatalf("decode stop status: %v", err)
	}
	if stopped.Status != StateStopped || stopped.BoundAddr != "" {
		t.Fatalf("expected stopped MITM, got %+v", stopped)
	}
}

func TestMitmProxyStreamsSSEAndBlocksRedirectsAndStripsHopByHop(t *testing.T) {
	_, handler, _ := setupTestMitm(t)

	chunk1Received := make(chan struct{})
	var sawCustomHop bool
	mockGateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/redirect":
			http.Redirect(w, r, "http://169.254.169.254/latest/meta-data", http.StatusFound)
		case "/hop":
			if r.Header.Get("X-Custom-Hop") != "" {
				sawCustomHop = true
			}
			w.WriteHeader(http.StatusOK)
		case "/sse":
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			rc := http.NewResponseController(w)
			_, _ = io.WriteString(w, "data: chunk1\n\n")
			_ = rc.Flush()
			<-chunk1Received
			_, _ = io.WriteString(w, "data: chunk2\n\n")
			_ = rc.Flush()
		}
	}))
	t.Cleanup(mockGateway.Close)

	rec := doMitmJSON(t, handler, http.MethodPatch, "/api/mitm/config", map[string]any{
		"listen_addr":        "127.0.0.1:0",
		"target_gateway_url": mockGateway.URL,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("patch config: %d", rec.Code)
	}
	startRec := doMitmJSON(t, handler, http.MethodPost, "/api/mitm/start", nil)
	if startRec.Code != http.StatusOK {
		t.Fatalf("start mitm: %d", startRec.Code)
	}
	var st Status
	_ = json.Unmarshal(startRec.Body.Bytes(), &st)

	caRec := doMitmJSON(t, handler, http.MethodGet, "/api/mitm/ca/download", nil)
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caRec.Body.Bytes())
	tlsClient := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    pool,
				ServerName: "api.individual.githubcopilot.com",
				MinVersion: tls.VersionTLS12,
			},
		},
	}
	t.Cleanup(tlsClient.CloseIdleConnections)

	// 1. Redirect is NOT followed by MITM proxy (SSRF prevention)
	redReq, _ := http.NewRequest(http.MethodGet, "https://"+st.BoundAddr+"/redirect", nil)
	redResp, err := tlsClient.Do(redReq)
	if err != nil {
		t.Fatalf("redirect request failed: %v", err)
	}
	_ = redResp.Body.Close()
	if redResp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 returned verbatim without following redirect, got %d", redResp.StatusCode)
	}

	// 2. Hop-by-hop headers stripped
	hopReq, _ := http.NewRequest(http.MethodGet, "https://"+st.BoundAddr+"/hop", nil)
	hopReq.Header.Set("Connection", "X-Custom-Hop")
	hopReq.Header.Set("X-Custom-Hop", "should-be-stripped")
	hopResp, err := tlsClient.Do(hopReq)
	if err != nil {
		t.Fatalf("hop request failed: %v", err)
	}
	_ = hopResp.Body.Close()
	if sawCustomHop {
		t.Fatalf("expected Connection-nominated hop-by-hop header X-Custom-Hop to be stripped")
	}

	// 3. SSE chunks flush immediately before upstream finishes
	sseReq, _ := http.NewRequest(http.MethodGet, "https://"+st.BoundAddr+"/sse", nil)
	sseResp, err := tlsClient.Do(sseReq)
	if err != nil {
		t.Fatalf("sse request failed: %v", err)
	}
	defer func() { _ = sseResp.Body.Close() }()

	buf := make([]byte, 64)
	n, err := sseResp.Body.Read(buf)
	if err != nil || !strings.Contains(string(buf[:n]), "chunk1") {
		close(chunk1Received)
		t.Fatalf("expected immediate chunk1 before stream close, got n=%d err=%v data=%q", n, err, string(buf[:n]))
	}
	close(chunk1Received)
}

func TestMitmConcurrentStartAndChunkedCARotation(t *testing.T) {
	_, handler, _ := setupTestMitm(t)

	// Generate initial CA
	rec1 := doMitmJSON(t, handler, http.MethodPost, "/api/mitm/ca/generate", map[string]any{"force": false})
	var ca1 CAInfo
	_ = json.Unmarshal(rec1.Body.Bytes(), &ca1)

	// Chunked request (ContentLength == -1) with {"force":true} must rotate CA
	req := httptest.NewRequest(http.MethodPost, "/api/mitm/ca/generate", strings.NewReader(`{"force":true}`))
	req.ContentLength = -1
	req.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req)
	if rec2.Code != http.StatusOK {
		t.Fatalf("chunked rotate CA: expected 200, got %d: %s", rec2.Code, rec2.Body.String())
	}
	var ca2 CAInfo
	_ = json.Unmarshal(rec2.Body.Bytes(), &ca2)
	if ca2.Fingerprint == ca1.Fingerprint {
		t.Fatalf("expected chunked force=true to rotate CA fingerprint")
	}

	// Concurrent start: exactly 1 succeeds (200), others fail with 409
	_ = doMitmJSON(t, handler, http.MethodPatch, "/api/mitm/config", map[string]any{"listen_addr": "127.0.0.1:0"})
	codes := make(chan int, 5)
	for range 5 {
		go func() {
			r := doMitmJSON(t, handler, http.MethodPost, "/api/mitm/start", nil)
			codes <- r.Code
		}()
	}
	okCount := 0
	conflictCount := 0
	for range 5 {
		switch c := <-codes; c {
		case http.StatusOK:
			okCount++
		case http.StatusConflict:
			conflictCount++
		}
	}
	if okCount != 1 || conflictCount != 4 {
		t.Fatalf("expected 1x200 and 4x409 on concurrent start, got ok=%d conflict=%d", okCount, conflictCount)
	}
}

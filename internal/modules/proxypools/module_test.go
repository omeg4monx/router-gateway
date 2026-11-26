package proxypools_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/proxypools"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
	"github.com/omeg4monx/router-gateway/internal/platform/venice"
)

var (
	_ lifecycle.Module              = (*proxypools.Module)(nil)
	_ contracts.ProxyDialerResolver = (*proxypools.Module)(nil)
)

type stubDepModule struct {
	name string
}

func (s stubDepModule) Name() string                         { return s.name }
func (s stubDepModule) Dependencies() []string               { return nil }
func (s stubDepModule) Migrations() fs.FS                    { return nil }
func (s stubDepModule) RegisterRoutes(lifecycle.RouteMounts) {}
func (s stubDepModule) Start(context.Context) error          { return nil }
func (s stubDepModule) Stop(context.Context) error           { return nil }

type testHarness struct {
	db     *db.DB
	mod    *proxypools.Module
	server *httpserver.Server
}

func newHarness(t *testing.T) *testHarness {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "proxypools.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	masterKey := make([]byte, crypto.MasterKeySize)
	if _, err := rand.Read(masterKey); err != nil {
		t.Fatalf("rand.Read: %v", err)
	}
	vault, err := crypto.NewVault(masterKey)
	if err != nil {
		t.Fatalf("crypto.NewVault: %v", err)
	}

	bus := eventbus.New()
	mod := proxypools.New(proxypools.Options{
		DB:           database,
		Vault:        vault,
		Bus:          bus,
		ProbeTimeout: 2 * time.Second,
	})

	cfg := config.Config{Host: "127.0.0.1", Port: 20128, APITimeout: 15 * time.Second}
	srv := httpserver.New(httpserver.Options{Config: cfg})

	reg := lifecycle.NewRegistry(nil)
	for _, m := range []lifecycle.Module{stubDepModule{name: "system"}, stubDepModule{name: "providers"}, mod} {
		if err := reg.Register(m); err != nil {
			t.Fatalf("register %s: %v", m.Name(), err)
		}
	}
	if err := reg.MigrateAll(context.Background(), database); err != nil {
		t.Fatalf("MigrateAll: %v", err)
	}
	if err := reg.MountAll(srv.Mounts()); err != nil {
		t.Fatalf("MountAll: %v", err)
	}
	if err := reg.StartAll(context.Background()); err != nil {
		t.Fatalf("StartAll: %v", err)
	}
	t.Cleanup(func() { _ = reg.StopAll(context.Background()) })

	return &testHarness{db: database, mod: mod, server: srv}
}

func (h *testHarness) doReq(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
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

func TestProxyPoolsCRUDCredentialRedactionAndValidation(t *testing.T) {
	t.Parallel()

	h := newHarness(t)

	// Invalid schemes and missing hosts rejected with 400.
	for _, bad := range []string{
		`{"name":"Bad FTP","url":"ftp://proxy.local:21"}`,
		`{"name":"Bad File","url":"file:///etc/passwd"}`,
		`{"name":"No Host","url":"http://:8080"}`,
		`{"name":"","url":"http://127.0.0.1:8080"}`,
	} {
		rec := h.doReq(t, http.MethodPost, "/api/proxy-pools", bad)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("POST /api/proxy-pools (%s) = %d, want 400: %s", bad, rec.Code, rec.Body.String())
		}
	}

	// Create HTTP proxy with embedded credentials.
	rec := h.doReq(t, http.MethodPost, "/api/proxy-pools", `{
		"name": "Corp HTTP",
		"url": "http://alice:supersecretpass@proxy.corp.local:8080"
	}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create HTTP pool = %d, want 201: %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "supersecretpass") || strings.Contains(rec.Body.String(), "alice") {
		t.Fatalf("response leaked proxy credentials: %s", rec.Body.String())
	}

	var created proxypools.PoolSummary
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal pool: %v", err)
	}
	if !created.HasCredentials || created.Scheme != "http" || created.DisplayURL != "http://***@proxy.corp.local:8080" {
		t.Fatalf("unexpected created pool: %+v", created)
	}

	// Create SOCKS5 pool without credentials.
	recSocks := h.doReq(t, http.MethodPost, "/api/proxy-pools", `{
		"name": "Edge SOCKS5",
		"url": "socks5h://127.0.0.1:1080"
	}`)
	if recSocks.Code != http.StatusCreated {
		t.Fatalf("create SOCKS5 pool = %d, want 201: %s", recSocks.Code, recSocks.Body.String())
	}

	// Bind provider and verify ProxyURLForProvider returns decrypted URL with credentials.
	bindRec := h.doReq(t, http.MethodPut, "/api/proxy-pool-bindings/openai", `{"pool_id":"`+created.ID+`"}`)
	if bindRec.Code != http.StatusOK {
		t.Fatalf("bind provider = %d, want 200: %s", bindRec.Code, bindRec.Body.String())
	}

	gotURL, err := h.mod.ProxyURLForProvider(context.Background(), "openai")
	if err != nil {
		t.Fatalf("ProxyURLForProvider: %v", err)
	}
	if gotURL == nil || gotURL.String() != "http://alice:supersecretpass@proxy.corp.local:8080" {
		t.Fatalf("ProxyURLForProvider(openai) = %v, want full decrypted URL", gotURL)
	}

	// Unbound provider returns nil.
	unboundURL, err := h.mod.ProxyURLForProvider(context.Background(), "anthropic")
	if err != nil || unboundURL != nil {
		t.Fatalf("ProxyURLForProvider(anthropic) = (%v, %v), want (nil, nil)", unboundURL, err)
	}

	// Disable pool -> ProxyURLForProvider returns nil.
	patchRec := h.doReq(t, http.MethodPatch, "/api/proxy-pools/"+created.ID, `{"enabled":false}`)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("disable pool = %d: %s", patchRec.Code, patchRec.Body.String())
	}
	disabledURL, err := h.mod.ProxyURLForProvider(context.Background(), "openai")
	if err != nil || disabledURL != nil {
		t.Fatalf("ProxyURLForProvider(disabled) = (%v, %v), want (nil, nil)", disabledURL, err)
	}
}

func TestProxyPoolsLatencyTestAndAutoScanner(t *testing.T) {
	t.Parallel()

	h := newHarness(t)

	// Start a local TCP listener representing a reachable proxy port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	// Obtain a closed port for an unreachable proxy candidate.
	closedLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen closed: %v", err)
	}
	closedAddr := closedLn.Addr().String()
	_ = closedLn.Close()

	// Create pool pointing to reachable listener and run /test.
	createRec := h.doReq(t, http.MethodPost, "/api/proxy-pools", `{"name":"Local Proxy","url":"http://`+ln.Addr().String()+`"}`)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create pool = %d: %s", createRec.Code, createRec.Body.String())
	}
	var pool proxypools.PoolSummary
	if err := json.Unmarshal(createRec.Body.Bytes(), &pool); err != nil {
		t.Fatalf("unmarshal pool: %v", err)
	}

	testRec := h.doReq(t, http.MethodPost, "/api/proxy-pools/"+pool.ID+"/test", `{}`)
	if testRec.Code != http.StatusOK {
		t.Fatalf("test pool = %d: %s", testRec.Code, testRec.Body.String())
	}
	var testRes proxypools.PoolTestResponse
	if err := json.Unmarshal(testRec.Body.Bytes(), &testRes); err != nil {
		t.Fatalf("unmarshal test response: %v", err)
	}
	if !testRes.Healthy || testRes.Status != "healthy" || testRes.LatencyMS < 0 {
		t.Fatalf("expected healthy test result, got %+v", testRes)
	}

	// Run auto-scanner with one healthy and one unreachable candidate, with auto_add=true.
	scanBody, _ := json.Marshal(map[string]any{
		"candidates":  []string{"http://" + ln.Addr().String(), "http://" + closedAddr},
		"auto_add":    true,
		"name_prefix": "Scanned",
		"timeout_ms":  1000,
	})
	scanRec := h.doReq(t, http.MethodPost, "/api/proxy-pools/scan", string(scanBody))
	if scanRec.Code != http.StatusOK {
		t.Fatalf("scan = %d: %s", scanRec.Code, scanRec.Body.String())
	}
	var scanRes proxypools.ScanResponse
	if err := json.Unmarshal(scanRec.Body.Bytes(), &scanRes); err != nil {
		t.Fatalf("unmarshal scan response: %v", err)
	}
	if len(scanRes.Results) != 2 {
		t.Fatalf("scan results len = %d, want 2", len(scanRes.Results))
	}
	if !scanRes.Results[0].Healthy || scanRes.Results[0].AddedPoolID == "" {
		t.Fatalf("expected first candidate healthy & auto-added, got %+v", scanRes.Results[0])
	}
	if scanRes.Results[1].Healthy || scanRes.Results[1].Error == "" {
		t.Fatalf("expected second candidate unhealthy with error, got %+v", scanRes.Results[1])
	}
}

func TestProxyPoolsBindingLifecycleAndCacheInvalidation(t *testing.T) {
	t.Parallel()

	h := newHarness(t)

	rec := h.doReq(t, http.MethodPost, "/api/proxy-pools", `{"name":"Primary","url":"http://127.0.0.1:8888"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create pool = %d: %s", rec.Code, rec.Body.String())
	}
	var pool proxypools.PoolSummary
	if err := json.Unmarshal(rec.Body.Bytes(), &pool); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if r := h.doReq(t, http.MethodPut, "/api/proxy-pool-bindings/qwen", `{"pool_id":"`+pool.ID+`"}`); r.Code != http.StatusOK {
		t.Fatalf("bind qwen = %d: %s", r.Code, r.Body.String())
	}

	// Update URL -> cache invalidates immediately.
	if r := h.doReq(t, http.MethodPatch, "/api/proxy-pools/"+pool.ID, `{"url":"http://127.0.0.1:9999"}`); r.Code != http.StatusOK {
		t.Fatalf("patch pool = %d: %s", r.Code, r.Body.String())
	}
	got, err := h.mod.ProxyURLForProvider(context.Background(), "qwen")
	if err != nil || got == nil || got.String() != "http://127.0.0.1:9999" {
		t.Fatalf("ProxyURLForProvider after patch = (%v, %v), want http://127.0.0.1:9999", got, err)
	}

	// Unbind -> returns nil.
	if r := h.doReq(t, http.MethodDelete, "/api/proxy-pool-bindings/qwen", ""); r.Code != http.StatusNoContent {
		t.Fatalf("unbind qwen = %d: %s", r.Code, r.Body.String())
	}
	got, err = h.mod.ProxyURLForProvider(context.Background(), "qwen")
	if err != nil || got != nil {
		t.Fatalf("ProxyURLForProvider after unbind = (%v, %v), want (nil, nil)", got, err)
	}

	// Rebind then delete pool -> binding cascades cleanly.
	if r := h.doReq(t, http.MethodPut, "/api/proxy-pool-bindings/qwen", `{"pool_id":"`+pool.ID+`"}`); r.Code != http.StatusOK {
		t.Fatalf("rebind qwen = %d: %s", r.Code, r.Body.String())
	}
	if r := h.doReq(t, http.MethodDelete, "/api/proxy-pools/"+pool.ID, ""); r.Code != http.StatusNoContent {
		t.Fatalf("delete pool = %d: %s", r.Code, r.Body.String())
	}
	got, err = h.mod.ProxyURLForProvider(context.Background(), "qwen")
	if err != nil || got != nil {
		t.Fatalf("ProxyURLForProvider after delete = (%v, %v), want (nil, nil)", got, err)
	}
}

func startAuthHTTPProxy(t *testing.T, expectedUserPass, tag string, hits *atomic.Int32) *httptest.Server {
	t.Helper()
	wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(expectedUserPass))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Proxy-Authorization") != wantAuth {
			w.WriteHeader(http.StatusProxyAuthRequired)
			_, _ = w.Write([]byte("proxy auth required"))
			return
		}
		hits.Add(1)
		outReq, err := http.NewRequestWithContext(r.Context(), r.Method, r.URL.String(), r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		outReq.Header = r.Header.Clone()
		outReq.Header.Del("Proxy-Authorization")
		outReq.Header.Set("X-Proxy-Tag", tag)
		resp, err := http.DefaultTransport.RoundTrip(outReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		for k, vals := range resp.Header {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}
		w.Header().Set("X-Via-Proxy", tag)
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestProxyPoolsEndToEndHTTPForwardingRotationAndAuthProbe(t *testing.T) {
	t.Parallel()

	h := newHarness(t)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"proxy":"` + r.Header.Get("X-Proxy-Tag") + `"}`))
	}))
	t.Cleanup(upstream.Close)

	var hitsA, hitsB atomic.Int32
	proxyA := startAuthHTTPProxy(t, "userA:passA", "proxy-a", &hitsA)
	proxyB := startAuthHTTPProxy(t, "userB:passB", "proxy-b", &hitsB)

	hostA := strings.TrimPrefix(proxyA.URL, "http://")
	hostB := strings.TrimPrefix(proxyB.URL, "http://")

	// 1. Create two authenticated proxy pools + one with wrong credentials to verify 407 probe rejection.
	createPool := func(name, rawURL string) proxypools.PoolSummary {
		t.Helper()
		rec := h.doReq(t, http.MethodPost, "/api/proxy-pools", `{"name":"`+name+`","url":"`+rawURL+`"}`)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create %s = %d: %s", name, rec.Code, rec.Body.String())
		}
		var out proxypools.PoolSummary
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatalf("unmarshal %s: %v", name, err)
		}
		return out
	}

	poolA := createPool("Webshare A", "http://userA:passA@"+hostA)
	_ = createPool("Webshare B", "http://userB:passB@"+hostB)
	badAuthPool := createPool("Bad Auth", "http://wrong:creds@"+hostA)

	// Probe with target_url must succeed for valid auth and fail on 407 for bad auth.
	okTest := h.doReq(t, http.MethodPost, "/api/proxy-pools/"+poolA.ID+"/test", `{"target_url":"`+upstream.URL+`"}`)
	var okRes proxypools.PoolTestResponse
	if err := json.Unmarshal(okTest.Body.Bytes(), &okRes); err != nil || !okRes.Healthy {
		t.Fatalf("valid auth probe = (%d, %+v, err=%v), want healthy", okTest.Code, okRes, err)
	}
	badTest := h.doReq(t, http.MethodPost, "/api/proxy-pools/"+badAuthPool.ID+"/test", `{"target_url":"`+upstream.URL+`"}`)
	var badRes proxypools.PoolTestResponse
	if err := json.Unmarshal(badTest.Body.Bytes(), &badRes); err != nil || badRes.Healthy {
		t.Fatalf("407 bad auth probe = (%d, %+v, err=%v), want unhealthy", badTest.Code, badRes, err)
	}
	// Remove badAuthPool so rotation pool only contains valid pools A & B.
	_ = h.doReq(t, http.MethodDelete, "/api/proxy-pools/"+badAuthPool.ID, "")

	// Build real router httpClient + venice.SelectHTTPClient (proves Venice uses proxy pools too!).
	routerClient := httpclient.NewClient(httpclient.Options{
		ProxyFunc: httpclient.NewProviderProxyFunc(h.mod),
	})
	veniceClient := venice.SelectHTTPClient(routerClient)

	doThrough := func(providerID string) string {
		t.Helper()
		ctx := httpclient.WithProviderID(context.Background(), providerID)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, upstream.URL, http.NoBody)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		resp, err := veniceClient.Do(req)
		if err != nil {
			t.Fatalf("Do(%s): %v", providerID, err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Do(%s) status=%d body=%s", providerID, resp.StatusCode, string(b))
		}
		return string(b)
	}

	// 2. Bind venice to poolA -> traverses proxyA; unbound openai dials directly.
	if rec := h.doReq(t, http.MethodPut, "/api/proxy-pool-bindings/venice", `{"pool_id":"`+poolA.ID+`"}`); rec.Code != http.StatusOK {
		t.Fatalf("bind venice to poolA = %d: %s", rec.Code, rec.Body.String())
	}
	if body := doThrough("venice"); body != `{"proxy":"proxy-a"}` {
		t.Fatalf("venice via poolA body = %s, want proxy-a", body)
	}
	if body := doThrough("openai"); body != `{"proxy":""}` {
		t.Fatalf("unbound openai body = %s, want direct empty proxy tag", body)
	}

	// 3. Bind venice to __rotate__ -> round-robin across enabled pools A and B.
	if rec := h.doReq(t, http.MethodPut, "/api/proxy-pool-bindings/venice", `{"pool_id":"`+proxypools.RotatePoolID+`"}`); rec.Code != http.StatusOK {
		t.Fatalf("bind venice to __rotate__ = %d: %s", rec.Code, rec.Body.String())
	}
	seenA, seenB := 0, 0
	for range 4 {
		switch doThrough("venice") {
		case `{"proxy":"proxy-a"}`:
			seenA++
		case `{"proxy":"proxy-b"}`:
			seenB++
		}
	}
	if seenA != 2 || seenB != 2 {
		t.Fatalf("rotation distribution = (A:%d, B:%d), want (2, 2)", seenA, seenB)
	}

	// 4. Disable poolA -> __rotate__ automatically routes 100% to poolB.
	if rec := h.doReq(t, http.MethodPatch, "/api/proxy-pools/"+poolA.ID, `{"enabled":false}`); rec.Code != http.StatusOK {
		t.Fatalf("disable poolA = %d: %s", rec.Code, rec.Body.String())
	}
	for range 2 {
		if body := doThrough("venice"); body != `{"proxy":"proxy-b"}` {
			t.Fatalf("rotation after disabling poolA = %s, want proxy-b", body)
		}
	}
}

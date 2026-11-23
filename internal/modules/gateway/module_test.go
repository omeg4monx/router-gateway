package gateway_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/gateway"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubModule struct {
	name string
	deps []string
}

func (s stubModule) Name() string                         { return s.name }
func (s stubModule) Dependencies() []string               { return s.deps }
func (s stubModule) Migrations() fs.FS                    { return nil }
func (s stubModule) RegisterRoutes(lifecycle.RouteMounts) {}
func (s stubModule) Start(context.Context) error          { return nil }
func (s stubModule) Stop(context.Context) error           { return nil }

var _ lifecycle.Module = (*gateway.Module)(nil)

type stubKeyValidator struct {
	keys map[string]contracts.KeyInfo
	errs map[string]error
}

func (s *stubKeyValidator) ValidateKey(_ context.Context, rawKey string) (contracts.KeyInfo, error) {
	if err, ok := s.errs[rawKey]; ok {
		return contracts.KeyInfo{}, err
	}
	if info, ok := s.keys[rawKey]; ok {
		return info, nil
	}
	return contracts.KeyInfo{}, errors.New("keys: invalid API key")
}

type stubResolver struct {
	fn func(ctx context.Context, model string, allowed []string) ([]contracts.ResolvedTarget, error)
}

func (s *stubResolver) Resolve(ctx context.Context, model string, allowed []string) ([]contracts.ResolvedTarget, error) {
	return s.fn(ctx, model, allowed)
}

type stubModelLister struct {
	models []contracts.ModelDescriptor
}

func (s *stubModelLister) ListModels(_ context.Context, allowed []string) ([]contracts.ModelDescriptor, error) {
	if len(allowed) == 0 {
		return append([]contracts.ModelDescriptor(nil), s.models...), nil
	}
	allowSet := make(map[string]bool, len(allowed))
	for _, a := range allowed {
		allowSet[a] = true
	}
	var out []contracts.ModelDescriptor
	for _, m := range s.models {
		if allowSet["*"] || allowSet[m.ID] {
			out = append(out, m)
		}
	}
	return out, nil
}

type stubProviderStore struct {
	mu           sync.Mutex
	creds        map[string]contracts.UpstreamCredentials
	refreshCalls map[string]int
	refreshFn    func(ctx context.Context, connectionID string) (contracts.UpstreamCredentials, error)
}

func (s *stubProviderStore) ListEnabledConnections(context.Context) ([]contracts.ProviderConnection, error) {
	return nil, nil
}

func (s *stubProviderStore) ResolveCredentials(_ context.Context, connectionID string) (contracts.UpstreamCredentials, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.creds[connectionID]
	if !ok {
		return contracts.UpstreamCredentials{}, contracts.ErrConnectionNotFound
	}
	return c, nil
}

func (s *stubProviderStore) RefreshOAuthToken(ctx context.Context, connectionID string) (contracts.UpstreamCredentials, error) {
	s.mu.Lock()
	if s.refreshCalls == nil {
		s.refreshCalls = make(map[string]int)
	}
	s.refreshCalls[connectionID]++
	fn := s.refreshFn
	s.mu.Unlock()
	if fn == nil {
		return contracts.UpstreamCredentials{}, contracts.ErrOAuthRefreshUnsupported
	}
	c, err := fn(ctx, connectionID)
	if err == nil {
		s.mu.Lock()
		s.creds[connectionID] = c
		s.mu.Unlock()
	}
	return c, err
}

func (s *stubProviderStore) RefreshCount(connectionID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.refreshCalls[connectionID]
}

type failureRecord struct {
	ConnectionID string
	StatusCode   int
	RetryAfter   time.Duration
}

type stubHealthTracker struct {
	mu        sync.Mutex
	failures  []failureRecord
	successes []string
	cooldowns map[string]bool
}

func (s *stubHealthTracker) IsInCooldown(connectionID string, _ time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cooldowns[connectionID]
}

func (s *stubHealthTracker) RecordFailure(connectionID string, statusCode int, retryAfter time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failures = append(s.failures, failureRecord{
		ConnectionID: connectionID,
		StatusCode:   statusCode,
		RetryAfter:   retryAfter,
	})
}

func (s *stubHealthTracker) RecordSuccess(connectionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.successes = append(s.successes, connectionID)
}

func (s *stubHealthTracker) Failures() []failureRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]failureRecord(nil), s.failures...)
}

func (s *stubHealthTracker) Successes() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.successes...)
}

type stubUsageRecorder struct {
	mu      sync.Mutex
	records []contracts.UsageRecord
}

func (s *stubUsageRecorder) Record(_ context.Context, rec contracts.UsageRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append(s.records, rec)
}

func (s *stubUsageRecorder) Records() []contracts.UsageRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]contracts.UsageRecord(nil), s.records...)
}

type recordingInterceptor struct {
	name     string
	priority int
	logMu    *sync.Mutex
	events   *[]string

	onRequestFn         func(ctx context.Context, req *contracts.CanonicalRequest) error
	onResponseHeadersFn func(ctx context.Context, meta *contracts.ResponseMeta) error
	onStreamChunkFn     func(ctx context.Context, ev *contracts.StreamEvent) error
}

func (r *recordingInterceptor) Name() string  { return r.name }
func (r *recordingInterceptor) Priority() int { return r.priority }

func (r *recordingInterceptor) record(stage string) {
	if r.logMu != nil && r.events != nil {
		r.logMu.Lock()
		*r.events = append(*r.events, fmt.Sprintf("%s:%d:%s", r.name, r.priority, stage))
		r.logMu.Unlock()
	}
}

func (r *recordingInterceptor) OnRequest(ctx context.Context, req *contracts.CanonicalRequest) error {
	r.record("OnRequest")
	if r.onRequestFn != nil {
		return r.onRequestFn(ctx, req)
	}
	return nil
}

func (r *recordingInterceptor) OnResponseHeaders(ctx context.Context, meta *contracts.ResponseMeta) error {
	r.record("OnResponseHeaders")
	if r.onResponseHeadersFn != nil {
		return r.onResponseHeadersFn(ctx, meta)
	}
	return nil
}

func (r *recordingInterceptor) OnStreamChunk(ctx context.Context, ev *contracts.StreamEvent) error {
	r.record("OnStreamChunk")
	if r.onStreamChunkFn != nil {
		return r.onStreamChunkFn(ctx, ev)
	}
	return nil
}

type testEnv struct {
	t         *testing.T
	db        *db.DB
	server    *httpserver.Server
	baseURL   string
	keys      *stubKeyValidator
	providers *stubProviderStore
	health    *stubHealthTracker
	usage     *stubUsageRecorder
	models    *stubModelLister
	resolver  *stubResolver
	module    *gateway.Module
}

func setupGatewayTest(t *testing.T, interceptors []contracts.GatewayInterceptor, chainTTL time.Duration, nilUsage bool) *testEnv {
	t.Helper()

	database, err := db.Open(filepath.Join(t.TempDir(), "gateway.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	keys := &stubKeyValidator{
		keys: map[string]contracts.KeyInfo{
			"sk-valid": {
				ID:            "key-1",
				Name:          "primary",
				AllowedModels: []string{"*"},
			},
			"sk-restricted": {
				ID:            "key-2",
				Name:          "restricted",
				AllowedModels: []string{"gpt-4o"},
			},
			"sk-disabled": {
				ID:       "key-3",
				Name:     "disabled-struct",
				Disabled: true,
			},
		},
		errs: map[string]error{
			"sk-err-disabled": errors.New("keys: API key is disabled"),
			"sk-err-expired":  errors.New("keys: API key has expired"),
		},
	}

	providers := &stubProviderStore{
		creds: make(map[string]contracts.UpstreamCredentials),
	}
	health := &stubHealthTracker{
		cooldowns: make(map[string]bool),
	}
	usage := &stubUsageRecorder{}
	models := &stubModelLister{}
	resolver := &stubResolver{
		fn: func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return nil, contracts.ErrModelNotFound
		},
	}

	var usageRec contracts.UsageRecorder
	if !nilUsage {
		usageRec = usage
	}

	mod := gateway.New(gateway.Options{
		DB:           database,
		HTTPClient:   &http.Client{},
		Keys:         keys,
		Resolver:     resolver,
		Models:       models,
		Providers:    providers,
		Health:       health,
		Interceptors: interceptors,
		Usage:        usageRec,
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		ChainTTL:     chainTTL,
	})

	reg := lifecycle.NewRegistry(nil)
	for _, dep := range []lifecycle.Module{
		stubModule{name: "system"},
		stubModule{name: "keys", deps: []string{"system"}},
		stubModule{name: "providers", deps: []string{"system"}},
		stubModule{name: "combos", deps: []string{"system", "providers"}},
		mod,
	} {
		if err := reg.Register(dep); err != nil {
			t.Fatalf("reg.Register(%s): %v", dep.Name(), err)
		}
	}

	ctx := context.Background()
	if err := reg.MigrateAll(ctx, database); err != nil {
		t.Fatalf("reg.MigrateAll: %v", err)
	}

	srv := httpserver.New(httpserver.Options{
		Config: config.Config{Host: "127.0.0.1", Port: 0},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	if err := reg.MountAll(srv.Mounts()); err != nil {
		t.Fatalf("reg.MountAll: %v", err)
	}
	if err := reg.StartAll(ctx); err != nil {
		t.Fatalf("reg.StartAll: %v", err)
	}
	t.Cleanup(func() {
		_ = reg.StopAll(context.Background())
	})

	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)

	return &testEnv{
		t:         t,
		db:        database,
		server:    srv,
		baseURL:   ts.URL,
		keys:      keys,
		providers: providers,
		health:    health,
		usage:     usage,
		models:    models,
		resolver:  resolver,
		module:    mod,
	}
}

func TestGatewayAuthExtractionAndProtocolErrors(t *testing.T) {
	env := setupGatewayTest(t, nil, time.Hour, true)

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/chat/completions"):
			_, _ = w.Write([]byte(`{"id":"chat-1","object":"chat.completion","model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"ok-bearer"},"finish_reason":"stop"}],"usage":{"prompt_tokens":3,"completion_tokens":2}}`))
		case strings.HasSuffix(r.URL.Path, "/v1/messages"):
			_, _ = w.Write([]byte(`{"id":"msg-1","type":"message","role":"assistant","model":"claude-sonnet-4","content":[{"type":"text","text":"ok-xapikey"}],"stop_reason":"end_turn","usage":{"input_tokens":3,"output_tokens":2}}`))
		case strings.Contains(r.URL.Path, ":generateContent"):
			_, _ = w.Write([]byte(`{"responseId":"gem-1","candidates":[{"index":0,"content":{"role":"model","parts":[{"text":"ok-query"}]},"finishReason":"STOP"}],"usageMetadata":{"promptTokenCount":3,"candidatesTokenCount":2}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer up.Close()

	env.providers.creds["c-oai"] = contracts.UpstreamCredentials{
		ConnectionID: "c-oai", ProviderID: "openai", Protocol: "openai-chat",
		BaseURL: up.URL, AuthType: "api_key", APIKey: "up-oai",
	}
	env.providers.creds["c-ant"] = contracts.UpstreamCredentials{
		ConnectionID: "c-ant", ProviderID: "anthropic", Protocol: "anthropic",
		BaseURL: up.URL, AuthType: "api_key", APIKey: "up-ant",
	}
	env.providers.creds["c-gem"] = contracts.UpstreamCredentials{
		ConnectionID: "c-gem", ProviderID: "gemini", Protocol: "gemini",
		BaseURL: up.URL, AuthType: "api_key", APIKey: "up-gem",
	}

	env.resolver.fn = func(_ context.Context, model string, allowed []string) ([]contracts.ResolvedTarget, error) {
		for _, a := range allowed {
			if a == "*" || a == model {
				switch model {
				case "gpt-4o":
					return []contracts.ResolvedTarget{{ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c-oai", Model: model, UpstreamModel: "gpt-4o", Endpoint: up.URL}}, nil
				case "claude-sonnet-4":
					return []contracts.ResolvedTarget{{ProviderID: "anthropic", ProviderType: "anthropic", ConnectionID: "c-ant", Model: model, UpstreamModel: "claude-sonnet-4", Endpoint: up.URL}}, nil
				case "gemini-2.5-flash":
					return []contracts.ResolvedTarget{{ProviderID: "gemini", ProviderType: "gemini", ConnectionID: "c-gem", Model: model, UpstreamModel: "gemini-2.5-flash", Endpoint: up.URL}}, nil
				default:
					return nil, contracts.ErrModelNotFound
				}
			}
		}
		return nil, contracts.ErrModelForbidden
	}

	t.Run("missing and invalid keys return protocol-appropriate 401", func(t *testing.T) {
		// OpenAI missing auth
		resp, err := http.Post(env.baseURL+"/v1/chat/completions", "application/json", strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}`))
		if err != nil {
			t.Fatalf("POST: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d: %s", resp.StatusCode, body)
		}
		var oaiErr struct {
			Error struct {
				Type string `json:"type"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &oaiErr); err != nil || oaiErr.Error.Type != "authentication_error" {
			t.Fatalf("expected OpenAI authentication_error JSON, got %s", body)
		}

		// Anthropic disabled key
		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/messages", strings.NewReader(`{"model":"claude-sonnet-4","max_tokens":16,"messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "sk-disabled")
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST anthropic: %v", err)
		}
		body, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d: %s", resp.StatusCode, body)
		}
		var antErr struct {
			Type  string `json:"type"`
			Error struct {
				Type string `json:"type"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &antErr); err != nil || antErr.Type != "error" || antErr.Error.Type != "authentication_error" {
			t.Fatalf("expected Anthropic error JSON, got %s", body)
		}

		// Gemini expired key via ?key=
		resp, err = http.Post(env.baseURL+"/v1beta/models/gemini-2.5-flash:generateContent?key=sk-err-expired", "application/json", strings.NewReader(`{"contents":[{"role":"user","parts":[{"text":"hi"}]}]}`))
		if err != nil {
			t.Fatalf("POST gemini: %v", err)
		}
		body, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d: %s", resp.StatusCode, body)
		}
		var gemErr struct {
			Error struct {
				Code   int    `json:"code"`
				Status string `json:"status"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &gemErr); err != nil || gemErr.Error.Code != 401 || gemErr.Error.Status != "UNAUTHENTICATED" {
			t.Fatalf("expected Gemini UNAUTHENTICATED JSON, got %s", body)
		}
	})

	t.Run("valid keys via Bearer, x-api-key, and ?key= succeed", func(t *testing.T) {
		// Bearer
		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST bearer: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "ok-bearer") {
			t.Fatalf("unexpected bearer response %d: %s", resp.StatusCode, body)
		}

		// x-api-key
		req, _ = http.NewRequest(http.MethodPost, env.baseURL+"/v1/messages", strings.NewReader(`{"model":"claude-sonnet-4","max_tokens":16,"messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "sk-valid")
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST x-api-key: %v", err)
		}
		body, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "ok-xapikey") {
			t.Fatalf("unexpected x-api-key response %d: %s", resp.StatusCode, body)
		}

		// ?key=
		resp, err = http.Post(env.baseURL+"/v1beta/models/gemini-2.5-flash:generateContent?key=sk-valid", "application/json", strings.NewReader(`{"contents":[{"role":"user","parts":[{"text":"hi"}]}]}`))
		if err != nil {
			t.Fatalf("POST ?key=: %v", err)
		}
		body, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "ok-query") {
			t.Fatalf("unexpected ?key= response %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("forbidden model returns 403", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/messages", strings.NewReader(`{"model":"claude-sonnet-4","max_tokens":16,"messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "sk-restricted")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST forbidden: %v", err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", resp.StatusCode)
		}
	})
}

func TestGatewayDirectFastPathIncrementalStreaming(t *testing.T) {
	env := setupGatewayTest(t, nil, time.Hour, false)

	chunk1Sent := make(chan struct{})
	allowChunk2 := make(chan struct{})
	var upstreamReqBody []byte
	var upstreamAuth string

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamAuth = r.Header.Get("Authorization")
		upstreamReqBody, _ = io.ReadAll(r.Body)

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		rc := http.NewResponseController(w)

		_, _ = fmt.Fprint(w, "data: {\"id\":\"chatcmpl-1\",\"object\":\"chat.completion.chunk\",\"model\":\"gpt-4o-2024\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"hello \"},\"finish_reason\":null}]}\n\n")
		_ = rc.Flush()
		close(chunk1Sent)

		<-allowChunk2
		_, _ = fmt.Fprint(w, "data: {\"id\":\"chatcmpl-1\",\"object\":\"chat.completion.chunk\",\"model\":\"gpt-4o-2024\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"world\"},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":11,\"completion_tokens\":5}}\n\n")
		_, _ = fmt.Fprint(w, "data: [DONE]\n\n")
		_ = rc.Flush()
	}))
	defer up.Close()

	env.providers.creds["conn-1"] = contracts.UpstreamCredentials{
		ConnectionID: "conn-1",
		ProviderID:   "openai",
		Protocol:     "openai-chat",
		BaseURL:      up.URL,
		AuthType:     "api_key",
		APIKey:       "upstream-secret",
		ExtraHeaders: map[string]string{
			"X-Provider-Tag": "fastpath",
			"Host":           "evil.example.com",
		},
	}
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ProviderID:    "openai",
			ProviderType:  "openai-chat",
			ConnectionID:  "conn-1",
			Model:         "openai/gpt-4o-2024",
			UpstreamModel: "gpt-4o-2024",
			Endpoint:      up.URL,
		}}, nil
	}

	reqBody := `{"model":"openai/gpt-4o-2024","stream":true,"reasoning_effort":"high","messages":[{"role":"user","content":"hi"}]}`
	req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-valid")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("expected text/event-stream, got %q", ct)
	}

	<-chunk1Sent
	reader := bufio.NewReader(resp.Body)

	// Read first event BEFORE releasing chunk 2 — proves zero-buffering incremental flush!
	var firstLines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("read first chunk before upstream close: %v", err)
		}
		firstLines = append(firstLines, line)
		if line == "\n" {
			break
		}
	}
	firstEvent := strings.Join(firstLines, "")
	if !strings.Contains(firstEvent, "hello ") {
		t.Fatalf("expected first event to contain 'hello ', got %q", firstEvent)
	}

	// Now allow upstream to send chunk 2 and finish.
	close(allowChunk2)
	rest, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll rest: %v", err)
	}
	if !strings.Contains(string(rest), "world") || !strings.Contains(string(rest), "[DONE]") {
		t.Fatalf("expected second chunk and [DONE], got %q", string(rest))
	}

	if upstreamAuth != "Bearer upstream-secret" {
		t.Fatalf("unexpected upstream auth: %q", upstreamAuth)
	}
	if !strings.Contains(string(upstreamReqBody), `"model":"gpt-4o-2024"`) || !strings.Contains(string(upstreamReqBody), `"reasoning_effort":"high"`) {
		t.Fatalf("expected rewritten model and preserved reasoning_effort, got %s", upstreamReqBody)
	}

	records := env.usage.Records()
	if len(records) != 1 || records[0].InputTokens != 11 || records[0].OutputTokens != 5 || !records[0].Stream {
		t.Fatalf("unexpected usage records: %+v", records)
	}
}

func TestGatewayCrossProtocolStreaming(t *testing.T) {
	env := setupGatewayTest(t, nil, time.Hour, false)

	t.Run("OpenAI Chat to Anthropic SSE upstream", func(t *testing.T) {
		var gotHeaders http.Header
		up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotHeaders = r.Header.Clone()
			if r.URL.Path != "/v1/messages" {
				t.Errorf("expected /v1/messages, got %s", r.URL.Path)
			}
			w.Header().Set("Content-Type", "text/event-stream")
			rc := http.NewResponseController(w)
			events := []string{
				"event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"id\":\"msg_1\",\"type\":\"message\",\"role\":\"assistant\",\"model\":\"claude-sonnet-4\",\"content\":[],\"usage\":{\"input_tokens\":14,\"output_tokens\":1}}}\n\n",
				"event: content_block_start\ndata: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"text\",\"text\":\"\"}}\n\n",
				"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"Hello from Claude\"}}\n\n",
				"event: content_block_stop\ndata: {\"type\":\"content_block_stop\",\"index\":0}\n\n",
				"event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"},\"usage\":{\"output_tokens\":9}}\n\n",
				"event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n",
			}
			for _, ev := range events {
				_, _ = io.WriteString(w, ev)
				_ = rc.Flush()
			}
		}))
		defer up.Close()

		env.providers.creds["c-ant"] = contracts.UpstreamCredentials{
			ConnectionID: "c-ant", ProviderID: "anthropic", Protocol: "anthropic",
			BaseURL: up.URL, AuthType: "api_key", APIKey: "ant-key-123",
		}
		env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{{
				ProviderID: "anthropic", ProviderType: "anthropic", ConnectionID: "c-ant",
				Model: "claude-sonnet-4", UpstreamModel: "claude-sonnet-4", Endpoint: up.URL,
			}}, nil
		}

		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", strings.NewReader(`{"model":"claude-sonnet-4","stream":true,"messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
		}
		if gotHeaders.Get("x-api-key") != "ant-key-123" || gotHeaders.Get("anthropic-version") != "2023-06-01" {
			t.Fatalf("unexpected anthropic headers: %+v", gotHeaders)
		}
		if !strings.Contains(string(body), "Hello from Claude") || !strings.Contains(string(body), "data: [DONE]") {
			t.Fatalf("expected translated OpenAI SSE stream, got: %s", body)
		}
	})

	t.Run("Anthropic Messages to OpenAI SSE upstream", func(t *testing.T) {
		up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/chat/completions" {
				t.Errorf("expected /chat/completions, got %s", r.URL.Path)
			}
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-99\",\"object\":\"chat.completion.chunk\",\"model\":\"gpt-4o\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"Hi from OpenAI\"},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":8,\"completion_tokens\":4}}\n\n")
			_, _ = io.WriteString(w, "data: [DONE]\n\n")
		}))
		defer up.Close()

		env.providers.creds["c-oai"] = contracts.UpstreamCredentials{
			ConnectionID: "c-oai", ProviderID: "openai", Protocol: "openai-chat",
			BaseURL: up.URL, AuthType: "api_key", APIKey: "oai-key-123",
		}
		env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{{
				ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c-oai",
				Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up.URL,
			}}, nil
		}

		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/messages", strings.NewReader(`{"model":"gpt-4o","max_tokens":100,"stream":true,"messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
		}
		s := string(body)
		if !strings.Contains(s, "event: message_start") || !strings.Contains(s, "Hi from OpenAI") || !strings.Contains(s, "event: message_stop") {
			t.Fatalf("expected Anthropic SSE stream, got: %s", s)
		}
	})

	t.Run("Gemini streamGenerateContent to OpenAI SSE upstream", func(t *testing.T) {
		up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-gem\",\"object\":\"chat.completion.chunk\",\"model\":\"gpt-4o\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"Gemini translated\"},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":6,\"completion_tokens\":3}}\n\n")
			_, _ = io.WriteString(w, "data: [DONE]\n\n")
		}))
		defer up.Close()

		env.providers.creds["c-oai"] = contracts.UpstreamCredentials{
			ConnectionID: "c-oai", ProviderID: "openai", Protocol: "openai-chat",
			BaseURL: up.URL, AuthType: "api_key", APIKey: "oai-key-123",
		}
		env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{{
				ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c-oai",
				Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up.URL,
			}}, nil
		}

		resp, err := http.Post(env.baseURL+"/v1beta/models/gpt-4o:streamGenerateContent?alt=sse&key=sk-valid", "application/json", strings.NewReader(`{"contents":[{"role":"user","parts":[{"text":"hi"}]}]}`))
		if err != nil {
			t.Fatalf("POST: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
		}
		if !strings.Contains(string(body), "Gemini translated") || !strings.Contains(string(body), `"candidates"`) {
			t.Fatalf("expected Gemini SSE response, got: %s", body)
		}
	})
}

func TestGatewayErrorClassificationAndOAuthSingleRetry(t *testing.T) {
	t.Run("upstream 400 returns immediately without RecordFailure or fallback", func(t *testing.T) {
		env := setupGatewayTest(t, nil, time.Hour, true)

		var cand1Calls, cand2Calls int
		up1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cand1Calls++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"message":"context_length_exceeded: too many tokens","type":"invalid_request_error"}}`))
		}))
		defer up1.Close()

		up2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cand2Calls++
			w.WriteHeader(http.StatusOK)
		}))
		defer up2.Close()

		env.providers.creds["c1"] = contracts.UpstreamCredentials{
			ConnectionID: "c1", ProviderID: "openai", Protocol: "openai-chat", BaseURL: up1.URL, AuthType: "api_key", APIKey: "k1",
		}
		env.providers.creds["c2"] = contracts.UpstreamCredentials{
			ConnectionID: "c2", ProviderID: "openai", Protocol: "openai-chat", BaseURL: up2.URL, AuthType: "api_key", APIKey: "k2",
		}
		env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{
				{ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c1", Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up1.URL},
				{ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c2", Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up2.URL},
			}, nil
		}

		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d: %s", resp.StatusCode, body)
		}
		if !strings.Contains(string(body), "context_length_exceeded") {
			t.Fatalf("expected upstream message preserved, got: %s", body)
		}
		if cand1Calls != 1 || cand2Calls != 0 {
			t.Fatalf("expected cand1=1, cand2=0, got %d/%d", cand1Calls, cand2Calls)
		}
		if len(env.health.Failures()) != 0 {
			t.Fatalf("expected zero RecordFailure calls on 400, got: %+v", env.health.Failures())
		}
	})

	t.Run("upstream 429 with Retry-After records cooldown and fails over to candidate 2", func(t *testing.T) {
		env := setupGatewayTest(t, nil, time.Hour, true)

		up1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "45")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"message":"rate limited"}}`))
		}))
		defer up1.Close()

		up2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"chat-fallback","object":"chat.completion","model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"recovered"},"finish_reason":"stop"}]}`))
		}))
		defer up2.Close()

		env.providers.creds["c1"] = contracts.UpstreamCredentials{
			ConnectionID: "c1", ProviderID: "openai", Protocol: "openai-chat", BaseURL: up1.URL, AuthType: "api_key", APIKey: "k1",
		}
		env.providers.creds["c2"] = contracts.UpstreamCredentials{
			ConnectionID: "c2", ProviderID: "openai", Protocol: "openai-chat", BaseURL: up2.URL, AuthType: "api_key", APIKey: "k2",
		}
		env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{
				{ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c1", Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up1.URL},
				{ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c2", Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up2.URL},
			}, nil
		}

		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "recovered") {
			t.Fatalf("expected 200 recovered, got %d: %s", resp.StatusCode, body)
		}
		failures := env.health.Failures()
		if len(failures) != 1 || failures[0].ConnectionID != "c1" || failures[0].StatusCode != 429 || failures[0].RetryAfter != 45*time.Second {
			t.Fatalf("unexpected failures: %+v", failures)
		}
		successes := env.health.Successes()
		if len(successes) != 1 || successes[0] != "c2" {
			t.Fatalf("expected RecordSuccess(c2), got: %+v", successes)
		}
	})

	t.Run("upstream 401 on OAuth triggers single RefreshOAuthToken retry and succeeds", func(t *testing.T) {
		env := setupGatewayTest(t, nil, time.Hour, true)

		var tokensSeen []string
		up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			tokensSeen = append(tokensSeen, auth)
			if auth == "Bearer stale-token" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":{"message":"expired token"}}`))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"msg-oauth","type":"message","role":"assistant","model":"claude-sonnet-4","content":[{"type":"text","text":"oauth-ok"}],"stop_reason":"end_turn","usage":{"input_tokens":2,"output_tokens":2}}`))
		}))
		defer up.Close()

		env.providers.creds["c-oauth"] = contracts.UpstreamCredentials{
			ConnectionID: "c-oauth", ProviderID: "claude", Protocol: "anthropic",
			BaseURL: up.URL, AuthType: "oauth", AccessToken: "stale-token",
		}
		env.providers.refreshFn = func(_ context.Context, id string) (contracts.UpstreamCredentials, error) {
			return contracts.UpstreamCredentials{
				ConnectionID: id, ProviderID: "claude", Protocol: "anthropic",
				BaseURL: up.URL, AuthType: "oauth", AccessToken: "fresh-token",
			}, nil
		}
		env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{{
				ProviderID: "claude", ProviderType: "anthropic", ConnectionID: "c-oauth",
				Model: "claude-sonnet-4", UpstreamModel: "claude-sonnet-4", Endpoint: up.URL,
			}}, nil
		}

		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/messages", strings.NewReader(`{"model":"claude-sonnet-4","max_tokens":16,"messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "oauth-ok") {
			t.Fatalf("expected 200 oauth-ok, got %d: %s", resp.StatusCode, body)
		}
		if env.providers.RefreshCount("c-oauth") != 1 {
			t.Fatalf("expected 1 refresh call, got %d", env.providers.RefreshCount("c-oauth"))
		}
		if len(tokensSeen) != 2 || tokensSeen[0] != "Bearer stale-token" || tokensSeen[1] != "Bearer fresh-token" {
			t.Fatalf("unexpected tokens sequence: %+v", tokensSeen)
		}
		if len(env.health.Failures()) != 0 || len(env.health.Successes()) != 1 {
			t.Fatalf("expected 0 failures and 1 success, got failures=%+v successes=%+v", env.health.Failures(), env.health.Successes())
		}
	})

	t.Run("upstream 401 on OAuth fails at most once per candidate then fails over", func(t *testing.T) {
		env := setupGatewayTest(t, nil, time.Hour, true)

		var c1Calls int
		up1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c1Calls++
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":{"message":"still unauthorized"}}`))
		}))
		defer up1.Close()

		up2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"chat-c2","object":"chat.completion","model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"cand2-ok"},"finish_reason":"stop"}]}`))
		}))
		defer up2.Close()

		env.providers.creds["c1"] = contracts.UpstreamCredentials{
			ConnectionID: "c1", ProviderID: "codex", Protocol: "openai-chat",
			BaseURL: up1.URL, AuthType: "oauth", AccessToken: "t1",
		}
		env.providers.creds["c2"] = contracts.UpstreamCredentials{
			ConnectionID: "c2", ProviderID: "openai", Protocol: "openai-chat",
			BaseURL: up2.URL, AuthType: "api_key", APIKey: "k2",
		}
		env.providers.refreshFn = func(_ context.Context, id string) (contracts.UpstreamCredentials, error) {
			return contracts.UpstreamCredentials{
				ConnectionID: id, ProviderID: "codex", Protocol: "openai-chat",
				BaseURL: up1.URL, AuthType: "oauth", AccessToken: "t2",
			}, nil
		}
		env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{
				{ProviderID: "codex", ProviderType: "openai-chat", ConnectionID: "c1", Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up1.URL},
				{ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c2", Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up2.URL},
			}, nil
		}

		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "cand2-ok") {
			t.Fatalf("expected 200 cand2-ok, got %d: %s", resp.StatusCode, body)
		}
		if c1Calls != 2 || env.providers.RefreshCount("c1") != 1 {
			t.Fatalf("expected c1Calls=2 and refreshCount=1, got %d / %d", c1Calls, env.providers.RefreshCount("c1"))
		}
		if failures := env.health.Failures(); len(failures) != 1 || failures[0].ConnectionID != "c1" || failures[0].StatusCode != 401 {
			t.Fatalf("expected failure recorded for c1 401, got %+v", failures)
		}
	})

	t.Run("canceled client context during OAuth refresh does not call RecordFailure", func(t *testing.T) {
		env := setupGatewayTest(t, nil, time.Hour, true)

		up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":{"message":"expired token"}}`))
		}))
		defer up.Close()

		refreshEntered := make(chan struct{})
		env.providers.creds["c-oauth"] = contracts.UpstreamCredentials{
			ConnectionID: "c-oauth", ProviderID: "claude", Protocol: "anthropic",
			BaseURL: up.URL, AuthType: "oauth", AccessToken: "stale-token",
		}
		env.providers.refreshFn = func(ctx context.Context, _ string) (contracts.UpstreamCredentials, error) {
			close(refreshEntered)
			<-ctx.Done()
			return contracts.UpstreamCredentials{}, ctx.Err()
		}
		env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{{
				ProviderID: "claude", ProviderType: "anthropic", ConnectionID: "c-oauth",
				Model: "claude-sonnet-4", UpstreamModel: "claude-sonnet-4", Endpoint: up.URL,
			}}, nil
		}

		reqCtx, cancel := context.WithCancel(context.Background())
		req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(`{"model":"claude-sonnet-4","max_tokens":16,"messages":[{"role":"user","content":"hi"}]}`)).WithContext(reqCtx)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "sk-valid")
		rec := httptest.NewRecorder()

		handlerDone := make(chan struct{})
		go func() {
			defer close(handlerDone)
			env.server.Handler().ServeHTTP(rec, req)
		}()

		<-refreshEntered
		cancel()
		<-handlerDone

		if env.providers.RefreshCount("c-oauth") != 1 {
			t.Fatalf("expected 1 refresh call, got %d", env.providers.RefreshCount("c-oauth"))
		}
		if failures := env.health.Failures(); len(failures) != 0 {
			t.Fatalf("expected zero RecordFailure calls on canceled client context, got: %+v", failures)
		}
	})
}

func TestGatewayStreamFinalMultiFrameDoneFlag(t *testing.T) {
	type seenChunk struct {
		event string
		done  bool
	}
	var (
		mu     sync.Mutex
		chunks []seenChunk
	)
	ic := &recordingInterceptor{
		name:     "done-checker",
		priority: 1,
		onStreamChunkFn: func(_ context.Context, ev *contracts.StreamEvent) error {
			mu.Lock()
			chunks = append(chunks, seenChunk{event: ev.Event, done: ev.Done})
			mu.Unlock()
			return nil
		},
	}

	env := setupGatewayTest(t, []contracts.GatewayInterceptor{ic}, time.Hour, true)

	// Upstream is OpenAI chat; inbound is Anthropic so a terminal finish_reason chunk
	// emits multiple frames (content_block_stop, message_delta, message_stop) in one batch.
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"id\":\"chat-1\",\"object\":\"chat.completion.chunk\",\"model\":\"gpt-4o\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"hi\"},\"finish_reason\":null}]}\n\n")
		_, _ = io.WriteString(w, "data: {\"id\":\"chat-1\",\"object\":\"chat.completion.chunk\",\"model\":\"gpt-4o\",\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}))
	defer up.Close()

	env.providers.creds["c1"] = contracts.UpstreamCredentials{
		ConnectionID: "c1", ProviderID: "openai", Protocol: "openai-chat",
		BaseURL: up.URL, AuthType: "api_key", APIKey: "k1",
	}
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c1",
			Model: "claude-sonnet-4", UpstreamModel: "gpt-4o", Endpoint: up.URL,
		}}, nil
	}

	req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/messages", strings.NewReader(`{"model":"claude-sonnet-4","max_tokens":16,"stream":true,"messages":[{"role":"user","content":"hi"}]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "sk-valid")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	_, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(chunks) < 3 {
		t.Fatalf("expected at least 3 stream chunks, got %+v", chunks)
	}
	var doneCount int
	for i, ch := range chunks {
		if ch.done {
			doneCount++
			if i != len(chunks)-1 {
				t.Fatalf("chunk[%d] (%s) unexpectedly had Done=true before final chunk %+v", i, ch.event, chunks)
			}
		}
	}
	if doneCount != 1 || !chunks[len(chunks)-1].done {
		t.Fatalf("expected exactly the final chunk to have Done=true, got %+v", chunks)
	}
}

func TestGatewayInterceptorsPriorityOrdering(t *testing.T) {
	var mu sync.Mutex
	var events []string

	// Register out of priority order (20 before 10) to verify sorting!
	i20 := &recordingInterceptor{
		name:     "second",
		priority: 20,
		logMu:    &mu,
		events:   &events,
		onResponseHeadersFn: func(_ context.Context, meta *contracts.ResponseMeta) error {
			meta.Headers.Set("X-Order", meta.Headers.Get("X-Order")+",20")
			return nil
		},
		onStreamChunkFn: func(_ context.Context, ev *contracts.StreamEvent) error {
			if bytes.Contains(ev.Data, []byte(`"msg"`)) {
				ev.Data = bytes.ReplaceAll(ev.Data, []byte(`"msg"`), []byte(`"msg-20"`))
			}
			return nil
		},
	}
	i10 := &recordingInterceptor{
		name:     "first",
		priority: 10,
		logMu:    &mu,
		events:   &events,
		onResponseHeadersFn: func(_ context.Context, meta *contracts.ResponseMeta) error {
			meta.Headers.Set("X-Order", "10")
			return nil
		},
		onStreamChunkFn: func(_ context.Context, ev *contracts.StreamEvent) error {
			if bytes.Contains(ev.Data, []byte(`"hi"`)) {
				ev.Data = bytes.ReplaceAll(ev.Data, []byte(`"hi"`), []byte(`"msg"`))
			}
			return nil
		},
	}

	env := setupGatewayTest(t, []contracts.GatewayInterceptor{i20, i10}, time.Hour, true)

	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"id\":\"c1\",\"object\":\"chat.completion.chunk\",\"model\":\"gpt-4o\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"hi\"},\"finish_reason\":\"stop\"}]}\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}))
	defer up.Close()

	env.providers.creds["c1"] = contracts.UpstreamCredentials{
		ConnectionID: "c1", ProviderID: "openai", Protocol: "openai-chat",
		BaseURL: up.URL, AuthType: "api_key", APIKey: "k1",
	}
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c1",
			Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up.URL,
		}}, nil
	}

	req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4o","stream":true,"messages":[{"role":"user","content":"hi"}]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-valid")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	if resp.Header.Get("X-Order") != "10,20" {
		t.Fatalf("expected X-Order=10,20, got %q", resp.Header.Get("X-Order"))
	}
	if !strings.Contains(string(body), `"msg-20"`) {
		t.Fatalf("expected chained stream mutation 'msg-20', got: %s", body)
	}

	mu.Lock()
	defer mu.Unlock()
	expectedPrefix := []string{
		"first:10:OnRequest",
		"second:20:OnRequest",
		"first:10:OnResponseHeaders",
		"second:20:OnResponseHeaders",
		"first:10:OnStreamChunk",
		"second:20:OnStreamChunk",
	}
	if len(events) < len(expectedPrefix) {
		t.Fatalf("expected at least %d events, got %+v", len(expectedPrefix), events)
	}
	for i, want := range expectedPrefix {
		if events[i] != want {
			t.Fatalf("event[%d]: want %q, got %q (all=%+v)", i, want, events[i], events)
		}
	}
}

func TestGatewayModelsCountTokensResponseChainsAndAliases(t *testing.T) {
	env := setupGatewayTest(t, nil, time.Hour, true)

	env.models.models = []contracts.ModelDescriptor{
		{ID: "gpt-4o", OwnedBy: "openai", Kind: "model"},
		{ID: "claude-sonnet-4", OwnedBy: "anthropic", Kind: "model"},
	}

	var seenRequests []string
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		seenRequests = append(seenRequests, r.URL.Path+":"+string(raw))
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/messages":
			n := len(seenRequests)
			_, _ = fmt.Fprintf(w, `{"id":"msg_%d","type":"message","role":"assistant","model":"claude-sonnet-4","content":[{"type":"text","text":"answer-%d"}],"stop_reason":"end_turn","usage":{"input_tokens":5,"output_tokens":3}}`, n, n)
		case "/responses":
			_, _ = w.Write([]byte(`{"id":"resp_native","object":"response","status":"completed","model":"gpt-4o","output":[{"type":"message","id":"m1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"native-ok"}]}]}`))
		case "/responses/compact":
			_, _ = w.Write([]byte(`{"id":"cmp_1","object":"response.compaction","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"compacted"}]}]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer up.Close()

	env.providers.creds["c-ant"] = contracts.UpstreamCredentials{
		ConnectionID: "c-ant", ProviderID: "anthropic", Protocol: "anthropic", BaseURL: up.URL, AuthType: "api_key", APIKey: "k-ant",
	}
	env.providers.creds["c-resp"] = contracts.UpstreamCredentials{
		ConnectionID: "c-resp", ProviderID: "openai-responses", Protocol: "openai-responses", BaseURL: up.URL, AuthType: "api_key", APIKey: "k-resp",
	}
	env.resolver.fn = func(_ context.Context, model string, allowed []string) ([]contracts.ResolvedTarget, error) {
		switch model {
		case "claude-sonnet-4":
			return []contracts.ResolvedTarget{{ProviderID: "anthropic", ProviderType: "anthropic", ConnectionID: "c-ant", Model: model, UpstreamModel: "claude-sonnet-4", Endpoint: up.URL}}, nil
		case "gpt-4o":
			return []contracts.ResolvedTarget{{ProviderID: "openai-responses", ProviderType: "openai-responses", ConnectionID: "c-resp", Model: model, UpstreamModel: "gpt-4o", Endpoint: up.URL}}, nil
		default:
			return nil, contracts.ErrModelNotFound
		}
	}

	t.Run("/v1/models and /models compatibility alias filter by key", func(t *testing.T) {
		for _, path := range []string{"/v1/models", "/models", "/v1beta/models"} {
			req, _ := http.NewRequest(http.MethodGet, env.baseURL+path, http.NoBody)
			req.Header.Set("Authorization", "Bearer sk-restricted")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("GET %s: %v", path, err)
			}
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("GET %s: status %d: %s", path, resp.StatusCode, body)
			}
			if !strings.Contains(string(body), "gpt-4o") || strings.Contains(string(body), "claude-sonnet-4") {
				t.Fatalf("GET %s: unexpected filtered models: %s", path, body)
			}
		}

		// GET /v1/models/{id} and /models/{id}
		for _, path := range []string{"/v1/models/gpt-4o", "/models/gpt-4o"} {
			req, _ := http.NewRequest(http.MethodGet, env.baseURL+path, http.NoBody)
			req.Header.Set("Authorization", "Bearer sk-restricted")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("GET %s: %v", path, err)
			}
			_ = resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("GET %s: expected 200, got %d", path, resp.StatusCode)
			}
		}
	})

	t.Run("/v1/messages/count_tokens heuristic fallback", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/messages/count_tokens", strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"hello world"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST count_tokens: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", resp.StatusCode, body)
		}
		var out struct {
			InputTokens int64 `json:"input_tokens"`
		}
		if err := json.Unmarshal(body, &out); err != nil || out.InputTokens <= 0 {
			t.Fatalf("expected positive input_tokens, got %s", body)
		}
	})

	t.Run("/v1/responses previous_response_id chain persistence scoped by key_id", func(t *testing.T) {
		// Turn 1 on Anthropic upstream via /v1/responses
		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/responses", strings.NewReader(`{"model":"claude-sonnet-4","input":"Remember code 42"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("turn 1: %v", err)
		}
		body1, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("turn 1 status %d: %s", resp.StatusCode, body1)
		}
		var r1 struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(body1, &r1); err != nil || r1.ID == "" {
			t.Fatalf("turn 1 missing id: %s", body1)
		}

		// Turn 2 with previous_response_id using same key -> must prepend prior user+assistant messages
		req2, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/responses", strings.NewReader(fmt.Sprintf(`{"model":"claude-sonnet-4","previous_response_id":%q,"input":"What was the code?"}`, r1.ID)))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer sk-valid")
		resp2, err := http.DefaultClient.Do(req2)
		if err != nil {
			t.Fatalf("turn 2: %v", err)
		}
		body2, _ := io.ReadAll(resp2.Body)
		_ = resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			t.Fatalf("turn 2 status %d: %s", resp2.StatusCode, body2)
		}

		lastUpstream := seenRequests[len(seenRequests)-1]
		if !strings.Contains(lastUpstream, "Remember code 42") || !strings.Contains(lastUpstream, "answer-1") || !strings.Contains(lastUpstream, "What was the code?") {
			t.Fatalf("expected chained conversation history in upstream request, got: %s", lastUpstream)
		}

		// Turn 3 with different key (key-2) attempting to read key-1's chain -> must NOT load key-1's history!
		env.keys.keys["sk-other"] = contracts.KeyInfo{ID: "key-other", AllowedModels: []string{"*"}}
		req3, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/responses", strings.NewReader(fmt.Sprintf(`{"model":"claude-sonnet-4","previous_response_id":%q,"input":"Leak?"}`, r1.ID)))
		req3.Header.Set("Content-Type", "application/json")
		req3.Header.Set("Authorization", "Bearer sk-other")
		resp3, err := http.DefaultClient.Do(req3)
		if err != nil {
			t.Fatalf("turn 3: %v", err)
		}
		_ = resp3.Body.Close()

		isolatedUpstream := seenRequests[len(seenRequests)-1]
		if strings.Contains(isolatedUpstream, "Remember code 42") {
			t.Fatalf("key-other was able to read key-1 response chain! got: %s", isolatedUpstream)
		}
	})

	t.Run("compatibility aliases /responses, /codex/responses, and /v1/responses/compact", func(t *testing.T) {
		for _, path := range []string{"/responses", "/codex/responses"} {
			req, _ := http.NewRequest(http.MethodPost, env.baseURL+path, strings.NewReader(`{"model":"gpt-4o","input":"ping"}`))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer sk-valid")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("POST %s: %v", path, err)
			}
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "native-ok") {
				t.Fatalf("POST %s: unexpected %d: %s", path, resp.StatusCode, body)
			}
		}

		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/responses/compact", strings.NewReader(`{"model":"gpt-4o","input":"long conversation"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST compact: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "compacted") {
			t.Fatalf("POST compact: unexpected %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("expired response_chains entry is ignored and lazily pruned", func(t *testing.T) {
		past := time.Now().UTC().Add(-2 * time.Hour).Format("2006-01-02T15:04:05.000000000Z")
		_, err := env.db.Writer().Exec(
			`INSERT OR REPLACE INTO response_chains (response_id, key_id, model, messages_json, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
			"resp_expired", "key-1", "claude-sonnet-4", `[{"role":"user","content":[{"type":"text","text":"SECRET_EXPIRED"}]}]`, past, past,
		)
		if err != nil {
			t.Fatalf("insert expired chain: %v", err)
		}

		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/responses", strings.NewReader(`{"model":"claude-sonnet-4","previous_response_id":"resp_expired","input":"fresh turn"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST expired chain: %v", err)
		}
		_ = resp.Body.Close()

		lastUpstream := seenRequests[len(seenRequests)-1]
		if strings.Contains(lastUpstream, "SECRET_EXPIRED") {
			t.Fatalf("expired chain was loaded into upstream request: %s", lastUpstream)
		}

		var count int
		if err := env.db.Reader().QueryRow(`SELECT COUNT(*) FROM response_chains WHERE response_id = 'resp_expired'`).Scan(&count); err != nil {
			t.Fatalf("query expired count: %v", err)
		}
		if count != 0 {
			t.Fatalf("expected expired row to be lazily pruned, got count=%d", count)
		}
	})
}

func TestGatewayLargeSSELineAndAPIKey401Failover(t *testing.T) {
	t.Run("handles >64KB SSE data line without scanner truncation", func(t *testing.T) {
		env := setupGatewayTest(t, nil, time.Hour, true)
		largeText := strings.Repeat("A", 96*1024)

		up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = fmt.Fprintf(w, "data: {\"id\":\"c-big\",\"object\":\"chat.completion.chunk\",\"model\":\"gpt-4o\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":%q},\"finish_reason\":\"stop\"}]}\n\n", largeText)
			_, _ = io.WriteString(w, "data: [DONE]\n\n")
		}))
		defer up.Close()

		env.providers.creds["c1"] = contracts.UpstreamCredentials{
			ConnectionID: "c1", ProviderID: "openai", Protocol: "openai-chat", BaseURL: up.URL, AuthType: "api_key", APIKey: "k1",
		}
		env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{{
				ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c1",
				Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up.URL,
			}}, nil
		}

		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4o","stream":true,"messages":[{"role":"user","content":"big"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), largeText) {
			t.Fatalf("expected 96KB SSE payload intact, got status=%d len=%d", resp.StatusCode, len(body))
		}
	})

	t.Run("401 on api_key skips OAuth refresh and fails over to candidate 2", func(t *testing.T) {
		env := setupGatewayTest(t, nil, time.Hour, true)

		up1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":{"message":"bad api key"}}`))
		}))
		defer up1.Close()

		up2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"chat-2","object":"chat.completion","model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"fallback-2"},"finish_reason":"stop"}]}`))
		}))
		defer up2.Close()

		env.providers.creds["c1"] = contracts.UpstreamCredentials{
			ConnectionID: "c1", ProviderID: "openai", Protocol: "openai-chat", BaseURL: up1.URL, AuthType: "api_key", APIKey: "bad",
		}
		env.providers.creds["c2"] = contracts.UpstreamCredentials{
			ConnectionID: "c2", ProviderID: "openai", Protocol: "openai-chat", BaseURL: up2.URL, AuthType: "api_key", APIKey: "good",
		}
		env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{
				{ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c1", Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up1.URL},
				{ProviderID: "openai", ProviderType: "openai-chat", ConnectionID: "c2", Model: "gpt-4o", UpstreamModel: "gpt-4o", Endpoint: up2.URL},
			}, nil
		}

		req, _ := http.NewRequest(http.MethodPost, env.baseURL+"/v1/chat/completions", strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sk-valid")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "fallback-2") {
			t.Fatalf("expected 200 fallback-2, got %d: %s", resp.StatusCode, body)
		}
		if env.providers.RefreshCount("c1") != 0 {
			t.Fatalf("expected 0 OAuth refresh calls for api_key, got %d", env.providers.RefreshCount("c1"))
		}
		if failures := env.health.Failures(); len(failures) != 1 || failures[0].ConnectionID != "c1" || failures[0].StatusCode != 401 {
			t.Fatalf("expected failure recorded for c1 401, got %+v", failures)
		}
	})
}

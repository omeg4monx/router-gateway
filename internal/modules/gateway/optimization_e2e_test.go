package gateway_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/gateway"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

func TestGatewayCoalescingAndDeterministicCache(t *testing.T) {
	t.Parallel()

	var upstreamHits atomic.Int32
	releaseUpstream := make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		upstreamHits.Add(1)
		<-releaseUpstream
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"chatcmpl-cache","object":"chat.completion","model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"cached-ok"},"finish_reason":"stop"}],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}`))
	}))
	defer upstream.Close()

	env := setupGatewayTest(t, nil, 0, false)
	env.providers.creds["conn-cache"] = contracts.UpstreamCredentials{
		ConnectionID: "conn-cache",
		ProviderID:   "openai",
		Protocol:     gateway.ProtocolOpenAIChat,
		BaseURL:      upstream.URL,
		APIKey:       "sk-upstream",
	}
	env.resolver.fn = func(_ context.Context, _ string, _ []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ProviderID:    "openai",
			ProviderType:  gateway.ProtocolOpenAIChat,
			ConnectionID:  "conn-cache",
			Model:         "gpt-4o",
			UpstreamModel: "gpt-4o",
			Endpoint:      upstream.URL,
		}}, nil
	}

	payload := []byte(`{"model":"gpt-4o","temperature":0,"messages":[{"role":"user","content":"deterministic"}]}`)
	const concurrency = 5
	var wg sync.WaitGroup
	cacheHeaders := make([]string, concurrency)

	for i := range concurrency {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewReader(payload))
			if err != nil {
				t.Errorf("new req[%d]: %v", idx, err)
				return
			}
			req.Header.Set("Authorization", "Bearer sk-valid")
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("do req[%d]: %v", idx, err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				t.Errorf("status[%d] = %d, body = %s", idx, resp.StatusCode, string(body))
				return
			}
			cacheHeaders[idx] = resp.Header.Get("X-Router-Cache")
		}(i)
	}

	time.Sleep(50 * time.Millisecond)
	close(releaseUpstream)
	wg.Wait()

	if upstreamHits.Load() != 1 {
		t.Fatalf("upstreamHits = %d, want 1 (coalesced concurrent requests)", upstreamHits.Load())
	}

	// Follow-up identical request must be served from TTL cache without hitting upstream.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("new cached req: %v", err)
	}
	req.Header.Set("Authorization", "Bearer sk-valid")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do cached req: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("cached status = %d, want 200", resp.StatusCode)
	}
	if got := resp.Header.Get("X-Router-Cache"); got != "HIT" {
		t.Fatalf("X-Router-Cache = %q, want HIT", got)
	}
	if upstreamHits.Load() != 1 {
		t.Fatalf("upstreamHits after HIT = %d, want 1", upstreamHits.Load())
	}
}

func TestGatewaySingleCandidateTransientRetry(t *testing.T) {
	t.Parallel()

	var attempts atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := attempts.Add(1)
		if n == 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"message":"rate limit momentary"}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"chatcmpl-retry","object":"chat.completion","model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"recovered"},"finish_reason":"stop"}]}`))
	}))
	defer upstream.Close()

	env := setupGatewayTest(t, nil, 0, false)
	env.providers.creds["conn-solo"] = contracts.UpstreamCredentials{
		ConnectionID: "conn-solo",
		ProviderID:   "openai",
		Protocol:     gateway.ProtocolOpenAIChat,
		BaseURL:      upstream.URL,
		APIKey:       "sk-solo",
	}
	env.resolver.fn = func(_ context.Context, _ string, _ []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{{
			ProviderID:    "openai",
			ProviderType:  gateway.ProtocolOpenAIChat,
			ConnectionID:  "conn-solo",
			Model:         "gpt-4o",
			UpstreamModel: "gpt-4o",
			Endpoint:      upstream.URL,
		}}, nil
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, env.baseURL+"/v1/chat/completions", bytes.NewReader([]byte(`{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}`)))
	if err != nil {
		t.Fatalf("new req: %v", err)
	}
	req.Header.Set("Authorization", "Bearer sk-valid")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do req: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want 200, body = %s", resp.StatusCode, string(body))
	}
	if attempts.Load() != 2 {
		t.Fatalf("attempts = %d, want 2", attempts.Load())
	}
	if failures := env.health.Failures(); len(failures) != 0 {
		t.Fatalf("expected 0 health failures after transient recovery, got %+v", failures)
	}
}

func TestGatewayGlobalConcurrencyLoadShedding(t *testing.T) {
	t.Parallel()

	started := make(chan struct{})
	release := make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-started:
		default:
			close(started)
		}
		select {
		case <-release:
		case <-r.Context().Done():
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"chatcmpl-1","object":"chat.completion","model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`))
	}))
	defer upstream.Close()

	database, err := db.Open(filepath.Join(t.TempDir(), "shed.sqlite"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	srv := httpserver.New(httpserver.Options{
		Config: config.Config{Host: "127.0.0.1", Port: 0},
	})
	keys := &stubKeyValidator{
		keys: map[string]contracts.KeyInfo{
			"sk-valid": {ID: "key_1", Name: "Primary", AllowedModels: []string{"*"}, Scopes: []string{"gateway"}},
		},
	}
	providers := &stubProviderStore{
		creds: map[string]contracts.UpstreamCredentials{
			"conn-1": {
				ConnectionID: "conn-1",
				ProviderID:   "openai",
				Protocol:     gateway.ProtocolOpenAIChat,
				BaseURL:      upstream.URL,
				APIKey:       "sk-1",
			},
		},
	}
	resolver := &stubResolver{
		fn: func(_ context.Context, _ string, _ []string) ([]contracts.ResolvedTarget, error) {
			return []contracts.ResolvedTarget{{
				ProviderID:    "openai",
				ProviderType:  gateway.ProtocolOpenAIChat,
				ConnectionID:  "conn-1",
				Model:         "gpt-4o",
				UpstreamModel: "gpt-4o",
				Endpoint:      upstream.URL,
			}}, nil
		},
	}
	mod := gateway.New(gateway.Options{
		DB:                database,
		Keys:              keys,
		Resolver:          resolver,
		Models:            &stubModelLister{},
		Providers:         providers,
		Health:            &stubHealthTracker{},
		Usage:             &stubUsageRecorder{},
		MaxGlobalInFlight: 1,
		MaxConnInFlight:   1,
		QueueTimeout:      40 * time.Millisecond,
	})
	reg := lifecycle.NewRegistry(nil)
	for _, m := range []lifecycle.Module{
		stubModule{name: "system"},
		stubModule{name: "keys"},
		stubModule{name: "providers"},
		stubModule{name: "combos"},
		mod,
	} {
		if err := reg.Register(m); err != nil {
			t.Fatalf("register %s: %v", m.Name(), err)
		}
	}
	if err := reg.MigrateAll(context.Background(), database); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	reg.MountAll(srv.Mounts())
	if err := reg.StartAll(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() { _ = reg.StopAll(context.Background()) })

	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)

	// Start first request that occupies the sole global slot.
	firstDone := make(chan int, 1)
	go func() {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, ts.URL+"/v1/chat/completions", bytes.NewReader([]byte(`{"model":"gpt-4o","messages":[{"role":"user","content":"slow"}]}`)))
		req.Header.Set("Authorization", "Bearer sk-valid")
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			firstDone <- 0
			return
		}
		_ = resp.Body.Close()
		firstDone <- resp.StatusCode
	}()

	select {
	case <-started:
	case status := <-firstDone:
		t.Fatalf("first request returned early with status %d before reaching upstream", status)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for first request to reach upstream")
	}

	// Second concurrent request must be shed with 503 + Retry-After.
	req2, err := http.NewRequestWithContext(context.Background(), http.MethodPost, ts.URL+"/v1/chat/completions", bytes.NewReader([]byte(`{"model":"gpt-4o","messages":[{"role":"user","content":"shed"}]}`)))
	if err != nil {
		t.Fatalf("new req2: %v", err)
	}
	req2.Header.Set("Authorization", "Bearer sk-valid")
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("do req2: %v", err)
	}
	_ = resp2.Body.Close()

	close(release)
	if status1 := <-firstDone; status1 != http.StatusOK {
		t.Fatalf("first request status = %d, want 200", status1)
	}
	if resp2.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("shed request status = %d, want 503", resp2.StatusCode)
	}
	if resp2.Header.Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header on 503 load-shed response")
	}
}

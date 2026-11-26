package gateway

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

func TestTrafficShaperLimitsPerConnectionConcurrencyAndQueues(t *testing.T) {
	t.Parallel()

	shaper := newTrafficShaper(10, 2, 2*time.Second, nil)

	var (
		inFlight    atomic.Int32
		maxInFlight atomic.Int32
		completed   atomic.Int32
		wg          sync.WaitGroup
	)

	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			release, err := shaper.acquireConnection(context.Background(), "conn-single")
			if err != nil {
				t.Errorf("acquireConnection failed: %v", err)
				return
			}
			cur := inFlight.Add(1)
			for {
				prev := maxInFlight.Load()
				if cur <= prev || maxInFlight.CompareAndSwap(prev, cur) {
					break
				}
			}
			time.Sleep(20 * time.Millisecond)
			inFlight.Add(-1)
			release()
			completed.Add(1)
		}()
	}
	wg.Wait()

	if completed.Load() != 8 {
		t.Fatalf("completed = %d, want 8", completed.Load())
	}
	if maxInFlight.Load() > 2 {
		t.Fatalf("maxInFlight = %d, want <= 2", maxInFlight.Load())
	}
}

func TestTrafficShaperShedsLoadWhenGlobalFull(t *testing.T) {
	t.Parallel()

	shaper := newTrafficShaper(1, 2, 30*time.Millisecond, nil)

	release, err := shaper.acquireGlobal(context.Background())
	if err != nil {
		t.Fatalf("first acquireGlobal: %v", err)
	}
	defer release()

	_, err = shaper.acquireGlobal(context.Background())
	if !errors.Is(err, errGlobalOverloaded) {
		t.Fatalf("second acquireGlobal err = %v, want %v", err, errGlobalOverloaded)
	}
}

func TestTrafficShaperAdaptiveHeaderPacing(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 11, 25, 12, 0, 0, 0, time.UTC)
	var mu sync.Mutex
	clock := func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return now
	}

	shaper := newTrafficShaper(10, 4, 2*time.Second, clock)
	headers := http.Header{}
	headers.Set("x-ratelimit-remaining-requests", "0")
	headers.Set("x-ratelimit-reset-requests", "60ms")

	shaper.observeRateLimitHeaders("conn-paced", headers)

	start := time.Now()
	// Advance wall wait via real time while checking wait duration.
	shaperReal := newTrafficShaper(10, 4, 2*time.Second, time.Now)
	shaperReal.observeRateLimitHeaders("conn-paced", headers)
	release, err := shaperReal.acquireConnection(context.Background(), "conn-paced")
	if err != nil {
		t.Fatalf("acquireConnection after pace: %v", err)
	}
	release()
	if elapsed := time.Since(start); elapsed < 40*time.Millisecond {
		t.Fatalf("expected adaptive pacing to wait ~60ms, waited only %v", elapsed)
	}
}

func TestResponseCoalescerDeduplicatesAndCachesPerKey(t *testing.T) {
	t.Parallel()

	coalescer := newResponseCoalescer(16, 30*time.Second, nil)

	if !coalescer.isCacheableRequest(ProtocolOpenAIChat, false, false, []byte(`{"model":"gpt-4o","temperature":0,"messages":[{"role":"user","content":"hi"}]}`)) {
		t.Fatal("expected temperature=0 non-streaming request to be cacheable")
	}
	if coalescer.isCacheableRequest(ProtocolOpenAIChat, true, false, []byte(`{"model":"gpt-4o","temperature":0}`)) {
		t.Fatal("expected streaming request not to be cacheable")
	}
	if coalescer.isCacheableRequest(ProtocolOpenAIChat, false, false, []byte(`{"model":"gpt-4o","temperature":0.7}`)) {
		t.Fatal("expected temperature>0 request not to be cacheable")
	}
	if coalescer.isCacheableRequest(ProtocolOpenAIChat, false, false, []byte(`{"model":"gpt-4o","messages":[]}`)) {
		t.Fatal("expected omitted temperature request not to be cacheable")
	}

	body := []byte(`{"model":"gpt-4o","temperature":0,"messages":[{"role":"user","content":"ping"}]}`)
	keyA := coalescer.cacheKey("key-a", ProtocolOpenAIChat, "gpt-4o", body)
	keyB := coalescer.cacheKey("key-b", ProtocolOpenAIChat, "gpt-4o", body)
	if keyA == keyB {
		t.Fatal("cache keys must be isolated per API key")
	}

	var upstreamCalls atomic.Int32
	gate := make(chan struct{})
	var wg sync.WaitGroup
	statuses := make([]string, 5)

	for i := range 5 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			resp, status, err := coalescer.do(context.Background(), keyA, func(_ context.Context) (cachedGatewayResponse, error) {
				upstreamCalls.Add(1)
				<-gate
				return cachedGatewayResponse{
					statusCode:   http.StatusOK,
					body:         []byte(`{"id":"ok"}`),
					providerID:   "openai",
					connectionID: "conn-1",
				}, nil
			})
			if err != nil {
				t.Errorf("coalescer.do[%d]: %v", idx, err)
				return
			}
			if string(resp.body) != `{"id":"ok"}` {
				t.Errorf("unexpected body[%d]: %s", idx, string(resp.body))
			}
			statuses[idx] = status
		}(i)
	}

	time.Sleep(30 * time.Millisecond)
	close(gate)
	wg.Wait()

	if upstreamCalls.Load() != 1 {
		t.Fatalf("upstreamCalls = %d, want 1 (singleflight coalescing)", upstreamCalls.Load())
	}

	// Subsequent call must be a cache HIT without calling upstream.
	resp, status, err := coalescer.do(context.Background(), keyA, func(_ context.Context) (cachedGatewayResponse, error) {
		upstreamCalls.Add(1)
		return cachedGatewayResponse{}, nil
	})
	if err != nil {
		t.Fatalf("cached do: %v", err)
	}
	if status != cacheStatusHit {
		t.Fatalf("status = %q, want %q", status, cacheStatusHit)
	}
	if upstreamCalls.Load() != 1 {
		t.Fatalf("upstreamCalls after HIT = %d, want 1", upstreamCalls.Load())
	}
	if string(resp.body) != `{"id":"ok"}` {
		t.Fatalf("cached body = %q", string(resp.body))
	}
}

func TestPrioritizeHeadroomAffinityPromotesMatchingCandidateWithoutLosingFallbacks(t *testing.T) {
	t.Parallel()

	base := []contracts.ResolvedTarget{
		{ProviderID: "anthropic", ConnectionID: "conn-1"},
		{ProviderID: "anthropic", ConnectionID: "conn-2"},
		{ProviderID: "openai", ConnectionID: "conn-3"},
	}

	got := prioritizeHeadroomAffinity(base, map[string]any{
		"headroom_affinity_provider":   "anthropic",
		"headroom_affinity_connection": "conn-2",
	})
	if len(got) != 3 || got[0].ConnectionID != "conn-2" || got[1].ConnectionID != "conn-1" || got[2].ConnectionID != "conn-3" {
		t.Fatalf("unexpected candidate order after affinity promotion: %+v", got)
	}
	// Original slice must not be mutated.
	if base[0].ConnectionID != "conn-1" {
		t.Fatalf("expected input candidates slice to remain unmodified, got %+v", base)
	}

	// Mismatched provider/connection leaves order unchanged.
	unchanged := prioritizeHeadroomAffinity(base, map[string]any{
		"headroom_affinity_provider":   "openai",
		"headroom_affinity_connection": "conn-2",
	})
	if unchanged[0].ConnectionID != "conn-1" {
		t.Fatalf("expected mismatched provider+connection not to promote, got %+v", unchanged)
	}
}

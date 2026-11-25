package keys

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

func TestKeyLimiterShardedConcurrencyAndRPM(t *testing.T) {
	t.Parallel()

	limiter := newKeyLimiter()
	limiter.setLimits("key-hot", keyLimitSnapshot{
		name:         "hot-key",
		rateLimitRPM: 25,
	})
	limiter.setLimits("key-other", keyLimitSnapshot{
		name:         "other-key",
		rateLimitRPM: 25,
	})

	// ensureLimits must not overwrite already-cached snapshot.
	limiter.ensureLimits("key-hot", keyLimitSnapshot{
		name:         "stale-name",
		rateLimitRPM: 999,
	})
	if snap, ok := limiter.getLimits("key-hot"); !ok || snap.rateLimitRPM != 25 {
		t.Fatalf("expected ensureLimits to preserve existing snapshot, got %+v (ok=%v)", snap, ok)
	}

	mod := &Module{limiter: limiter}
	const workers = 60
	var (
		wg        sync.WaitGroup
		allowed   atomic.Int32
		rejected  atomic.Int32
		otherPass atomic.Int32
	)
	for range workers {
		wg.Add(2)
		go func() {
			defer wg.Done()
			err := mod.OnRequest(context.Background(), &contracts.CanonicalRequest{KeyID: "key-hot"})
			if err == nil {
				allowed.Add(1)
			} else {
				rejected.Add(1)
			}
		}()
		go func() {
			defer wg.Done()
			if err := mod.OnRequest(context.Background(), &contracts.CanonicalRequest{KeyID: "key-other"}); err == nil {
				otherPass.Add(1)
			}
		}()
	}
	wg.Wait()

	if allowed.Load() != 25 {
		t.Fatalf("allowed for key-hot = %d, want 25", allowed.Load())
	}
	if rejected.Load() != workers-25 {
		t.Fatalf("rejected for key-hot = %d, want %d", rejected.Load(), workers-25)
	}
	if otherPass.Load() != 25 {
		t.Fatalf("allowed for key-other = %d, want 25", otherPass.Load())
	}

	// Verify last-used throttle under concurrent touches.
	now := time.Now().UTC()
	var touched atomic.Int32
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.shouldTouchLastUsed("key-hot", now) {
				touched.Add(1)
			}
		}()
	}
	wg.Wait()
	if touched.Load() != 1 {
		t.Fatalf("shouldTouchLastUsed returned true %d times, want 1", touched.Load())
	}

	// Verify purgeKey & clear clean up state cleanly.
	limiter.purgeKey("key-hot")
	if _, ok := limiter.getLimits("key-hot"); ok {
		t.Fatal("expected key-hot limits to be purged")
	}
	limiter.clear()
	if _, ok := limiter.getLimits("key-other"); ok {
		t.Fatal("expected key-other limits to be cleared")
	}
}

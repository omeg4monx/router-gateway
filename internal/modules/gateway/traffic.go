package gateway

import (
	"context"
	"errors"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultMaxGlobalInFlight = 256
	defaultMaxConnInFlight   = 16
	defaultQueueTimeout      = 15 * time.Second
	maxTrackedConnections    = 2048
	maxAdaptivePaceWait      = 10 * time.Second

	singleCandidateRetryMax   = 2
	singleCandidateRetryBase  = 250 * time.Millisecond
	singleCandidateRetryMaxMs = 3000
)

var (
	errGlobalOverloaded = errors.New("gateway overloaded: global concurrency limit reached")
	errConnectionBusy   = errors.New("upstream connection busy: concurrency queue timed out")
)

// trafficShaper bounds global and per-connection in-flight concurrency, queues
// bursty traffic cleanly with context cancellation, and applies proactive
// header-aware pacing before an upstream connection hits 429.
type trafficShaper struct {
	globalSem    chan struct{}
	connLimit    int
	queueTimeout time.Duration
	now          func() time.Time

	mu        sync.Mutex
	connSems  map[string]chan struct{}
	paceUntil map[string]time.Time
}

func newTrafficShaper(globalMax, connMax int, queueTimeout time.Duration, clock func() time.Time) *trafficShaper {
	if globalMax <= 0 {
		globalMax = defaultMaxGlobalInFlight
	}
	if connMax <= 0 {
		connMax = defaultMaxConnInFlight
	}
	if queueTimeout <= 0 {
		queueTimeout = defaultQueueTimeout
	}
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	return &trafficShaper{
		globalSem:    make(chan struct{}, globalMax),
		connLimit:    connMax,
		queueTimeout: queueTimeout,
		now:          clock,
		connSems:     make(map[string]chan struct{}),
		paceUntil:    make(map[string]time.Time),
	}
}

func (s *trafficShaper) acquireGlobal(ctx context.Context) (func(), error) {
	if err := acquireSem(ctx, s.globalSem, s.queueTimeout, errGlobalOverloaded); err != nil {
		return nil, err
	}
	var once sync.Once
	return func() {
		once.Do(func() { <-s.globalSem })
	}, nil
}

func (s *trafficShaper) acquireConnection(ctx context.Context, connectionID string) (func(), error) {
	id := strings.TrimSpace(connectionID)
	if id == "" {
		return func() {}, nil
	}

	if err := s.waitForPace(ctx, id); err != nil {
		return nil, err
	}

	sem := s.connSemaphore(id)
	if err := acquireSem(ctx, sem, s.queueTimeout, errConnectionBusy); err != nil {
		return nil, err
	}
	var once sync.Once
	return func() {
		once.Do(func() { <-sem })
	}, nil
}

func (s *trafficShaper) connSemaphore(connectionID string) chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sem, ok := s.connSems[connectionID]; ok {
		return sem
	}
	if len(s.connSems) >= maxTrackedConnections {
		for k, ch := range s.connSems {
			if len(ch) == 0 {
				delete(s.connSems, k)
				delete(s.paceUntil, k)
			}
			if len(s.connSems) < maxTrackedConnections {
				break
			}
		}
	}
	sem := make(chan struct{}, s.connLimit)
	s.connSems[connectionID] = sem
	return sem
}

func (s *trafficShaper) waitForPace(ctx context.Context, connectionID string) error {
	s.mu.Lock()
	until, ok := s.paceUntil[connectionID]
	s.mu.Unlock()
	if !ok {
		return nil
	}
	wait := until.Sub(s.now())
	if wait <= 0 {
		s.mu.Lock()
		if cur, exists := s.paceUntil[connectionID]; exists && !cur.After(s.now()) {
			delete(s.paceUntil, connectionID)
		}
		s.mu.Unlock()
		return nil
	}
	if wait > s.queueTimeout {
		wait = s.queueTimeout
	}
	return sleepContext(ctx, wait)
}

// observeRateLimitHeaders inspects standard OpenAI/Anthropic rate-limit headers
// and schedules proactive pacing when remaining quota reaches zero.
func (s *trafficShaper) observeRateLimitHeaders(connectionID string, headers http.Header) {
	id := strings.TrimSpace(connectionID)
	if id == "" || len(headers) == 0 {
		return
	}
	now := s.now()
	var resetWait time.Duration

	if rem, ok := parseHeaderInt(headers, "x-ratelimit-remaining-requests", "anthropic-ratelimit-requests-remaining"); ok && rem <= 0 {
		if d := parseResetDuration(headers, now, "x-ratelimit-reset-requests", "anthropic-ratelimit-requests-reset"); d > resetWait {
			resetWait = d
		}
	}
	if rem, ok := parseHeaderInt(headers, "x-ratelimit-remaining-tokens", "anthropic-ratelimit-tokens-remaining"); ok && rem <= 0 {
		if d := parseResetDuration(headers, now, "x-ratelimit-reset-tokens", "anthropic-ratelimit-tokens-reset"); d > resetWait {
			resetWait = d
		}
	}
	if resetWait <= 0 {
		return
	}
	if resetWait > maxAdaptivePaceWait {
		resetWait = maxAdaptivePaceWait
	}

	target := now.Add(resetWait)
	s.mu.Lock()
	if len(s.paceUntil) >= maxTrackedConnections {
		for k, ts := range s.paceUntil {
			if !ts.After(now) {
				delete(s.paceUntil, k)
			}
		}
	}
	if existing, ok := s.paceUntil[id]; !ok || target.After(existing) {
		s.paceUntil[id] = target
	}
	s.mu.Unlock()
}

func parseHeaderInt(headers http.Header, keys ...string) (int64, bool) {
	for _, key := range keys {
		raw := strings.TrimSpace(headers.Get(key))
		if raw == "" {
			continue
		}
		v, err := strconv.ParseInt(raw, 10, 64)
		if err == nil {
			return v, true
		}
	}
	return 0, false
}

func parseResetDuration(headers http.Header, now time.Time, keys ...string) time.Duration {
	for _, key := range keys {
		raw := strings.TrimSpace(headers.Get(key))
		if raw == "" {
			continue
		}
		if d, err := time.ParseDuration(raw); err == nil && d > 0 {
			return d
		}
		if secs, err := strconv.ParseFloat(raw, 64); err == nil && secs > 0 {
			return time.Duration(secs * float64(time.Second))
		}
		if ts, err := time.Parse(time.RFC3339Nano, raw); err == nil {
			if d := ts.Sub(now); d > 0 {
				return d
			}
		}
	}
	return 0
}

func acquireSem(ctx context.Context, sem chan struct{}, timeout time.Duration, timeoutErr error) error {
	select {
	case sem <- struct{}{}:
		return nil
	default:
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return timeoutErr
	}
}

func sleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func isTransientSingleCandidateStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// computeJitteredRetryDelay computes bounded exponential backoff with full
// jitter so concurrent requests on a single key do not stampede upstream.
func computeJitteredRetryDelay(attempt int, retryAfter time.Duration) time.Duration {
	maxWait := time.Duration(singleCandidateRetryMaxMs) * time.Millisecond
	if retryAfter > maxWait {
		return 0
	}
	base := singleCandidateRetryBase << attempt
	if retryAfter > base {
		base = retryAfter
	}
	if base > maxWait {
		base = maxWait
	}
	if base <= 0 {
		return 0
	}
	// Full jitter in [base/2, base] to spread concurrent retries smoothly.
	half := base / 2
	if half <= 0 {
		return base
	}
	jitter := time.Duration(rand.Int64N(int64(half) + 1)) //nolint:gosec // G404: non-cryptographic PRNG is intentional for retry jitter
	return half + jitter
}

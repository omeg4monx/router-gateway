package gateway

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	defaultResponseCacheMaxEntries = 512
	defaultResponseCacheTTL        = 60 * time.Second
	defaultCoalesceExecutionLimit  = 2 * time.Minute
	maxCachedResponseBytes         = 256 << 10 // 256 KiB per cached entry

	cacheStatusHit       = "HIT"
	cacheStatusMiss      = "MISS"
	cacheStatusCoalesced = "COALESCED"
)

type cachedGatewayResponse struct {
	statusCode   int
	body         []byte
	usage        CanonicalUsage
	providerID   string
	connectionID string
	expiresAt    time.Time
}

type inflightCall struct {
	done chan struct{}
	resp cachedGatewayResponse
	err  error
}

// responseCoalescer combines singleflight in-flight deduplication with a
// bounded per-key TTL cache for strictly deterministic (temperature == 0,
// non-streaming, stateless) inference requests.
type responseCoalescer struct {
	maxEntries int
	ttl        time.Duration
	now        func() time.Time

	mu       sync.Mutex
	entries  map[string]cachedGatewayResponse
	inflight map[string]*inflightCall
}

func newResponseCoalescer(maxEntries int, ttl time.Duration, clock func() time.Time) *responseCoalescer {
	if maxEntries <= 0 {
		maxEntries = defaultResponseCacheMaxEntries
	}
	if ttl <= 0 {
		ttl = defaultResponseCacheTTL
	}
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	return &responseCoalescer{
		maxEntries: maxEntries,
		ttl:        ttl,
		now:        clock,
		entries:    make(map[string]cachedGatewayResponse),
		inflight:   make(map[string]*inflightCall),
	}
}

// isCacheableRequest returns true only for non-streaming, non-compact,
// stateless requests where temperature is explicitly set to 0.
func (c *responseCoalescer) isCacheableRequest(_ string, stream, isCompact bool, body []byte) bool {
	if c == nil || stream || isCompact || len(body) == 0 {
		return false
	}
	var probe struct {
		Temperature        *float64 `json:"temperature"`
		PreviousResponseID string   `json:"previous_response_id"`
		GenerationConfig   *struct {
			Temperature *float64 `json:"temperature"`
		} `json:"generationConfig"`
	}
	if err := json.Unmarshal(body, &probe); err != nil {
		return false
	}
	if probe.PreviousResponseID != "" {
		return false
	}
	if probe.Temperature != nil && *probe.Temperature == 0 {
		return true
	}
	if probe.GenerationConfig != nil && probe.GenerationConfig.Temperature != nil && *probe.GenerationConfig.Temperature == 0 {
		return true
	}
	return false
}

func (c *responseCoalescer) cacheKey(keyID, inProto, model string, body []byte) string {
	h := sha256.New()
	_, _ = h.Write([]byte(keyID))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(inProto))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(model))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write(bytes.TrimSpace(body))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *responseCoalescer) get(key string) (cachedGatewayResponse, bool) {
	now := c.now()
	c.mu.Lock()
	entry, ok := c.entries[key]
	if ok && !now.Before(entry.expiresAt) {
		delete(c.entries, key)
		ok = false
	}
	c.mu.Unlock()
	if !ok {
		return cachedGatewayResponse{}, false
	}
	entry.body = bytes.Clone(entry.body)
	return entry, true
}

func (c *responseCoalescer) set(key string, resp cachedGatewayResponse) {
	if resp.statusCode < http.StatusOK || resp.statusCode >= http.StatusMultipleChoices {
		return
	}
	if len(resp.body) == 0 || len(resp.body) > maxCachedResponseBytes {
		return
	}
	now := c.now()
	stored := cachedGatewayResponse{
		statusCode:   resp.statusCode,
		body:         bytes.Clone(resp.body),
		usage:        resp.usage,
		providerID:   resp.providerID,
		connectionID: resp.connectionID,
		expiresAt:    now.Add(c.ttl),
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.entries) >= c.maxEntries {
		// Evict expired entries first, then oldest entry if still full.
		var oldestKey string
		var oldestExp time.Time
		for k, v := range c.entries {
			if !now.Before(v.expiresAt) {
				delete(c.entries, k)
				continue
			}
			if oldestKey == "" || v.expiresAt.Before(oldestExp) {
				oldestKey = k
				oldestExp = v.expiresAt
			}
		}
		if len(c.entries) >= c.maxEntries && oldestKey != "" {
			delete(c.entries, oldestKey)
		}
	}
	c.entries[key] = stored
}

func (c *responseCoalescer) do(
	ctx context.Context,
	key string,
	fn func(context.Context) (cachedGatewayResponse, error),
) (cachedGatewayResponse, string, error) {
	if cached, ok := c.get(key); ok {
		return cached, cacheStatusHit, nil
	}

	c.mu.Lock()
	if entry, ok := c.entries[key]; ok && c.now().Before(entry.expiresAt) {
		c.mu.Unlock()
		entry.body = bytes.Clone(entry.body)
		return entry, cacheStatusHit, nil
	}
	if call, ok := c.inflight[key]; ok {
		c.mu.Unlock()
		select {
		case <-call.done:
			if call.err != nil {
				return cachedGatewayResponse{}, cacheStatusCoalesced, call.err
			}
			res := call.resp
			res.body = bytes.Clone(res.body)
			return res, cacheStatusCoalesced, nil
		case <-ctx.Done():
			return cachedGatewayResponse{}, cacheStatusCoalesced, ctx.Err()
		}
	}

	call := &inflightCall{done: make(chan struct{})}
	c.inflight[key] = call
	c.mu.Unlock()

	go func() {
		callCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), defaultCoalesceExecutionLimit)
		defer cancel()

		var (
			resp cachedGatewayResponse
			err  error
		)
		defer func() {
			if rec := recover(); rec != nil {
				err = fmt.Errorf("gateway: coalesced execution panicked: %v", rec)
			}
			call.resp = resp
			call.err = err
			if err == nil {
				c.set(key, resp)
			}
			c.mu.Lock()
			delete(c.inflight, key)
			close(call.done)
			c.mu.Unlock()
		}()

		resp, err = fn(callCtx)
	}()

	select {
	case <-call.done:
		if call.err != nil {
			return cachedGatewayResponse{}, cacheStatusMiss, call.err
		}
		res := call.resp
		res.body = bytes.Clone(res.body)
		return res, cacheStatusMiss, nil
	case <-ctx.Done():
		return cachedGatewayResponse{}, cacheStatusMiss, ctx.Err()
	}
}

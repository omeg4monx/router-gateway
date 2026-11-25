package keys

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
)

const (
	interceptorPriority = 1
	usageSubBufSize     = 256
	lastUsedThrottle    = 5 * time.Second
	limiterShardCount   = 32
)

type statusError struct {
	code int
	msg  string
}

func (e statusError) Error() string   { return e.msg }
func (e statusError) StatusCode() int { return e.code }

type keyLimitSnapshot struct {
	name         string
	rateLimitRPM int
	quotaUSD     float64
}

type keyLimiterShard struct {
	mu          sync.RWMutex
	limits      map[string]keyLimitSnapshot
	rpmWindows  map[string][]time.Time
	spent       map[string]float64
	spentSeeded map[string]bool
	lastTouched map[string]time.Time
}

type keyLimiter struct {
	shards [limiterShardCount]keyLimiterShard
}

func newKeyLimiter() *keyLimiter {
	l := &keyLimiter{}
	for i := range l.shards {
		l.shards[i] = keyLimiterShard{
			limits:      make(map[string]keyLimitSnapshot),
			rpmWindows:  make(map[string][]time.Time),
			spent:       make(map[string]float64),
			spentSeeded: make(map[string]bool),
			lastTouched: make(map[string]time.Time),
		}
	}
	return l
}

func (l *keyLimiter) shardFor(id string) *keyLimiterShard {
	const (
		fnvOffset32 uint32 = 2166136261
		fnvPrime32  uint32 = 16777619
	)
	h := fnvOffset32
	for i := 0; i < len(id); i++ {
		h ^= uint32(id[i])
		h *= fnvPrime32
	}
	return &l.shards[h&(limiterShardCount-1)]
}

func (l *keyLimiter) getLimits(id string) (keyLimitSnapshot, bool) {
	sh := l.shardFor(id)
	sh.mu.RLock()
	snap, ok := sh.limits[id]
	sh.mu.RUnlock()
	return snap, ok
}

func (l *keyLimiter) setLimits(id string, snap keyLimitSnapshot) {
	sh := l.shardFor(id)
	sh.mu.Lock()
	sh.limits[id] = snap
	sh.mu.Unlock()
}

// ensureLimits populates the cached limit snapshot only when absent, avoiding
// write-lock contention on the ValidateKey cache-hit hot path.
func (l *keyLimiter) ensureLimits(id string, snap keyLimitSnapshot) {
	sh := l.shardFor(id)
	sh.mu.RLock()
	_, ok := sh.limits[id]
	sh.mu.RUnlock()
	if ok {
		return
	}
	sh.mu.Lock()
	if _, exists := sh.limits[id]; !exists {
		sh.limits[id] = snap
	}
	sh.mu.Unlock()
}

func (l *keyLimiter) invalidateLimits(id string) {
	sh := l.shardFor(id)
	sh.mu.Lock()
	delete(sh.limits, id)
	sh.mu.Unlock()
}

func (l *keyLimiter) purgeKey(id string) {
	sh := l.shardFor(id)
	sh.mu.Lock()
	delete(sh.limits, id)
	delete(sh.rpmWindows, id)
	delete(sh.spent, id)
	delete(sh.spentSeeded, id)
	delete(sh.lastTouched, id)
	sh.mu.Unlock()
}

func (l *keyLimiter) clear() {
	for i := range l.shards {
		sh := &l.shards[i]
		sh.mu.Lock()
		sh.limits = make(map[string]keyLimitSnapshot)
		sh.rpmWindows = make(map[string][]time.Time)
		sh.spent = make(map[string]float64)
		sh.spentSeeded = make(map[string]bool)
		sh.lastTouched = make(map[string]time.Time)
		sh.mu.Unlock()
	}
}

func (l *keyLimiter) shouldTouchLastUsed(id string, now time.Time) bool {
	sh := l.shardFor(id)
	sh.mu.RLock()
	if prev, ok := sh.lastTouched[id]; ok && now.Sub(prev) < lastUsedThrottle {
		sh.mu.RUnlock()
		return false
	}
	sh.mu.RUnlock()

	sh.mu.Lock()
	defer sh.mu.Unlock()
	if prev, ok := sh.lastTouched[id]; ok && now.Sub(prev) < lastUsedThrottle {
		return false
	}
	sh.lastTouched[id] = now
	return true
}

func (l *keyLimiter) addSpend(id string, costUSD float64) {
	if id == "" || costUSD <= 0 {
		return
	}
	sh := l.shardFor(id)
	sh.mu.Lock()
	sh.spent[id] += costUSD
	sh.mu.Unlock()
}

// Priority runs per-key RPM and USD quota enforcement before team/content interceptors.
func (m *Module) Priority() int {
	return interceptorPriority
}

// OnRequest enforces per-key sliding-window RPM and cumulative USD quota limits.
func (m *Module) OnRequest(ctx context.Context, req *contracts.CanonicalRequest) error {
	if req == nil {
		return nil
	}
	keyID := strings.TrimSpace(req.KeyID)
	if keyID == "" || strings.HasPrefix(keyID, "admin-") {
		return nil
	}

	snap, ok := m.resolveKeyLimits(ctx, keyID)
	if !ok {
		return nil
	}

	// Fast path: keys without RPM or USD quota limits need no mutex acquisition.
	if snap.quotaUSD <= 0 && snap.rateLimitRPM <= 0 {
		return nil
	}

	// Seed historical spend from usage_logs once per key when quota_usd is configured.
	if snap.quotaUSD > 0 {
		m.ensureSpendSeeded(ctx, keyID)
	}

	now := time.Now()
	sh := m.limiter.shardFor(keyID)
	sh.mu.Lock()
	defer sh.mu.Unlock()

	if snap.quotaUSD > 0 {
		spent := sh.spent[keyID]
		if spent >= snap.quotaUSD {
			return statusError{
				code: http.StatusTooManyRequests,
				msg:  fmt.Sprintf("API key %q USD quota exhausted ($%.4f / $%.4f)", snap.name, spent, snap.quotaUSD),
			}
		}
	}

	if snap.rateLimitRPM > 0 {
		cutoff := now.Add(-time.Minute)
		window := sh.rpmWindows[keyID]
		pruned := window[:0]
		for _, ts := range window {
			if ts.After(cutoff) {
				pruned = append(pruned, ts)
			}
		}
		if len(pruned) >= snap.rateLimitRPM {
			sh.rpmWindows[keyID] = pruned
			return statusError{
				code: http.StatusTooManyRequests,
				msg:  fmt.Sprintf("API key %q rate limit exceeded (%d RPM)", snap.name, snap.rateLimitRPM),
			}
		}
		sh.rpmWindows[keyID] = append(pruned, now)
	}

	return nil
}

// OnResponseHeaders is a no-op for the keys module.
func (m *Module) OnResponseHeaders(_ context.Context, _ *contracts.ResponseMeta) error {
	return nil
}

// OnStreamChunk is a no-op for the keys module.
func (m *Module) OnStreamChunk(_ context.Context, _ *contracts.StreamEvent) error {
	return nil
}

func (m *Module) resolveKeyLimits(ctx context.Context, keyID string) (keyLimitSnapshot, bool) {
	if snap, ok := m.limiter.getLimits(keyID); ok {
		return snap, true
	}

	if m.repo == nil {
		return keyLimitSnapshot{}, false
	}
	summary, err := m.repo.GetByID(ctx, keyID)
	if err != nil {
		return keyLimitSnapshot{}, false
	}
	snap := keyLimitSnapshot{
		name:         summary.Name,
		rateLimitRPM: summary.RateLimitRPM,
		quotaUSD:     summary.QuotaUSD,
	}
	m.limiter.setLimits(keyID, snap)
	return snap, true
}

func (m *Module) ensureSpendSeeded(ctx context.Context, keyID string) {
	sh := m.limiter.shardFor(keyID)
	sh.mu.RLock()
	seeded := sh.spentSeeded[keyID]
	sh.mu.RUnlock()
	if seeded || m.repo == nil {
		return
	}

	historical, ok := m.repo.SumUsageCost(ctx, keyID)
	sh.mu.Lock()
	if !sh.spentSeeded[keyID] {
		sh.spentSeeded[keyID] = true
		if ok && historical > 0 {
			sh.spent[keyID] += historical
		}
	}
	sh.mu.Unlock()
}

func (m *Module) recordLastUsed(ctx context.Context, keyID string) {
	now := time.Now().UTC()
	if !m.limiter.shouldTouchLastUsed(keyID, now) {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	_ = m.repo.TouchLastUsed(ctx, keyID, now.Format(rfc3339NanoFixed))
}

func (m *Module) consumeUsage(ctx context.Context, ch <-chan eventbus.Event, unsub func()) {
	defer m.wg.Done()
	defer unsub()

	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			keyID, costUSD, valid := extractUsageCost(evt.Payload)
			if !valid {
				continue
			}
			m.limiter.addSpend(keyID, costUSD)
		}
	}
}

func extractUsageCost(payload any) (keyID string, costUSD float64, ok bool) {
	if payload == nil {
		return "", 0, false
	}
	if mp, isMap := payload.(map[string]any); isMap {
		id, _ := mp["key_id"].(string)
		id = strings.TrimSpace(id)
		if id == "" || strings.HasPrefix(id, "admin-") {
			return "", 0, false
		}
		switch v := mp["cost_usd"].(type) {
		case float64:
			return id, v, v > 0
		case int:
			return id, float64(v), v > 0
		case int64:
			return id, float64(v), v > 0
		default:
			return "", 0, false
		}
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return "", 0, false
	}
	var decoded struct {
		KeyID   string  `json:"key_id"`
		CostUSD float64 `json:"cost_usd"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return "", 0, false
	}
	id := strings.TrimSpace(decoded.KeyID)
	if id == "" || strings.HasPrefix(id, "admin-") || decoded.CostUSD <= 0 {
		return "", 0, false
	}
	return id, decoded.CostUSD, true
}

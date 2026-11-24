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

type keyLimiter struct {
	mu          sync.Mutex
	limits      map[string]keyLimitSnapshot
	rpmWindows  map[string][]time.Time
	spent       map[string]float64
	spentSeeded map[string]bool
	lastTouched map[string]time.Time
}

func newKeyLimiter() *keyLimiter {
	return &keyLimiter{
		limits:      make(map[string]keyLimitSnapshot),
		rpmWindows:  make(map[string][]time.Time),
		spent:       make(map[string]float64),
		spentSeeded: make(map[string]bool),
		lastTouched: make(map[string]time.Time),
	}
}

func (l *keyLimiter) setLimits(id string, snap keyLimitSnapshot) {
	l.mu.Lock()
	l.limits[id] = snap
	l.mu.Unlock()
}

func (l *keyLimiter) invalidateLimits(id string) {
	l.mu.Lock()
	delete(l.limits, id)
	l.mu.Unlock()
}

func (l *keyLimiter) purgeKey(id string) {
	l.mu.Lock()
	delete(l.limits, id)
	delete(l.rpmWindows, id)
	delete(l.spent, id)
	delete(l.spentSeeded, id)
	delete(l.lastTouched, id)
	l.mu.Unlock()
}

func (l *keyLimiter) clear() {
	l.mu.Lock()
	l.limits = make(map[string]keyLimitSnapshot)
	l.rpmWindows = make(map[string][]time.Time)
	l.spent = make(map[string]float64)
	l.spentSeeded = make(map[string]bool)
	l.lastTouched = make(map[string]time.Time)
	l.mu.Unlock()
}

func (l *keyLimiter) shouldTouchLastUsed(id string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if prev, ok := l.lastTouched[id]; ok && now.Sub(prev) < lastUsedThrottle {
		return false
	}
	l.lastTouched[id] = now
	return true
}

func (l *keyLimiter) addSpend(id string, costUSD float64) {
	if id == "" || costUSD <= 0 {
		return
	}
	l.mu.Lock()
	l.spent[id] += costUSD
	l.mu.Unlock()
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

	// Seed historical spend from usage_logs once per key when quota_usd is configured.
	if snap.quotaUSD > 0 {
		m.ensureSpendSeeded(ctx, keyID)
	}

	now := time.Now()
	m.limiter.mu.Lock()
	defer m.limiter.mu.Unlock()

	if snap.quotaUSD > 0 {
		spent := m.limiter.spent[keyID]
		if spent >= snap.quotaUSD {
			return statusError{
				code: http.StatusTooManyRequests,
				msg:  fmt.Sprintf("API key %q USD quota exhausted ($%.4f / $%.4f)", snap.name, spent, snap.quotaUSD),
			}
		}
	}

	if snap.rateLimitRPM > 0 {
		cutoff := now.Add(-time.Minute)
		window := m.limiter.rpmWindows[keyID]
		pruned := window[:0]
		for _, ts := range window {
			if ts.After(cutoff) {
				pruned = append(pruned, ts)
			}
		}
		if len(pruned) >= snap.rateLimitRPM {
			m.limiter.rpmWindows[keyID] = pruned
			return statusError{
				code: http.StatusTooManyRequests,
				msg:  fmt.Sprintf("API key %q rate limit exceeded (%d RPM)", snap.name, snap.rateLimitRPM),
			}
		}
		m.limiter.rpmWindows[keyID] = append(pruned, now)
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
	m.limiter.mu.Lock()
	snap, ok := m.limiter.limits[keyID]
	m.limiter.mu.Unlock()
	if ok {
		return snap, true
	}

	summary, err := m.repo.GetByID(ctx, keyID)
	if err != nil {
		return keyLimitSnapshot{}, false
	}
	snap = keyLimitSnapshot{
		name:         summary.Name,
		rateLimitRPM: summary.RateLimitRPM,
		quotaUSD:     summary.QuotaUSD,
	}
	m.limiter.setLimits(keyID, snap)
	return snap, true
}

func (m *Module) ensureSpendSeeded(ctx context.Context, keyID string) {
	m.limiter.mu.Lock()
	seeded := m.limiter.spentSeeded[keyID]
	m.limiter.mu.Unlock()
	if seeded {
		return
	}

	historical, ok := m.repo.SumUsageCost(ctx, keyID)
	m.limiter.mu.Lock()
	if !m.limiter.spentSeeded[keyID] {
		m.limiter.spentSeeded[keyID] = true
		if ok && historical > 0 {
			m.limiter.spent[keyID] += historical
		}
	}
	m.limiter.mu.Unlock()
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

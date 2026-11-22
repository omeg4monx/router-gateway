package combos

import (
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	defaultBackoffSchedule = []time.Duration{
		15 * time.Second,
		30 * time.Second,
		60 * time.Second,
		120 * time.Second,
		300 * time.Second,
	}

	authBackoffSchedule = []time.Duration{
		60 * time.Second,
		120 * time.Second,
		300 * time.Second,
	}
)

// CooldownEntry describes an active account cooldown window exposed by the
// admin observability endpoint.
type CooldownEntry struct {
	ConnectionID string `json:"connection_id"`
	Failures     int    `json:"failures"`
	Until        string `json:"until"`
	LastStatus   int    `json:"last_status"`
	UpdatedAt    string `json:"updated_at"`
}

type cooldownState struct {
	failures   int
	until      time.Time
	lastStatus int
	updatedAt  time.Time
}

type cooldownTracker struct {
	mu      sync.RWMutex
	now     func() time.Time
	entries map[string]cooldownState
}

func newCooldownTracker(clock func() time.Time) *cooldownTracker {
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	return &cooldownTracker{
		now:     clock,
		entries: make(map[string]cooldownState),
	}
}

func (t *cooldownTracker) currentTime(now time.Time) time.Time {
	if now.IsZero() {
		return t.now().UTC()
	}
	return now.UTC()
}

func (t *cooldownTracker) IsInCooldown(connectionID string, now time.Time) bool {
	id := strings.TrimSpace(connectionID)
	if id == "" {
		return false
	}
	ts := t.currentTime(now)

	t.mu.RLock()
	st, ok := t.entries[id]
	t.mu.RUnlock()

	return ok && ts.Before(st.until)
}

func (t *cooldownTracker) RecordFailure(connectionID string, statusCode int, retryAfter time.Duration) {
	id := strings.TrimSpace(connectionID)
	if id == "" {
		return
	}
	ts := t.now().UTC()

	t.mu.Lock()
	st := t.entries[id]
	st.failures++
	backoff := computeBackoff(st.failures, statusCode, retryAfter)
	if nextUntil := ts.Add(backoff); nextUntil.After(st.until) {
		st.until = nextUntil
	}
	st.lastStatus = statusCode
	st.updatedAt = ts
	t.entries[id] = st
	t.mu.Unlock()
}

func (t *cooldownTracker) RecordSuccess(connectionID string) {
	id := strings.TrimSpace(connectionID)
	if id == "" {
		return
	}
	t.mu.Lock()
	delete(t.entries, id)
	t.mu.Unlock()
}

func (t *cooldownTracker) ClearCooldown(connectionID string) bool {
	id := strings.TrimSpace(connectionID)
	if id == "" {
		return false
	}
	t.mu.Lock()
	_, existed := t.entries[id]
	delete(t.entries, id)
	t.mu.Unlock()
	return existed
}

func (t *cooldownTracker) ListCooldowns(now time.Time) []CooldownEntry {
	ts := t.currentTime(now)

	t.mu.RLock()
	out := make([]CooldownEntry, 0, len(t.entries))
	for id, st := range t.entries {
		if !ts.Before(st.until) {
			continue
		}
		out = append(out, CooldownEntry{
			ConnectionID: id,
			Failures:     st.failures,
			Until:        st.until.UTC().Format(rfc3339NanoFixed),
			LastStatus:   st.lastStatus,
			UpdatedAt:    st.updatedAt.UTC().Format(rfc3339NanoFixed),
		})
	}
	t.mu.RUnlock()

	sort.Slice(out, func(i, j int) bool {
		return out[i].ConnectionID < out[j].ConnectionID
	})
	return out
}

func computeBackoff(failures, statusCode int, retryAfter time.Duration) time.Duration {
	if failures < 1 {
		failures = 1
	}

	schedule := defaultBackoffSchedule
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		schedule = authBackoffSchedule
	}

	idx := failures - 1
	if idx >= len(schedule) {
		idx = len(schedule) - 1
	}
	base := schedule[idx]
	if retryAfter > base {
		return retryAfter
	}
	return base
}

package auth

import (
	"math"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	defaultLoginMaxFailures = 10
	defaultLoginWindow      = time.Minute
	maxRateLimitEntries     = 1024
)

type loginAttemptWindow struct {
	failures  int
	windowEnd time.Time
}

type loginRateLimiter struct {
	mu          sync.Mutex
	maxFailures int
	window      time.Duration
	entries     map[string]loginAttemptWindow
}

func newLoginRateLimiter(maxFailures int, window time.Duration) *loginRateLimiter {
	if maxFailures <= 0 {
		maxFailures = defaultLoginMaxFailures
	}
	if window <= 0 {
		window = defaultLoginWindow
	}
	return &loginRateLimiter{
		maxFailures: maxFailures,
		window:      window,
		entries:     make(map[string]loginAttemptWindow),
	}
}

// check returns (allowed, retryAfterSec) for clientKey at now.
func (l *loginRateLimiter) check(clientKey string, now time.Time) (allowed bool, retryAfterSec int) {
	if l == nil || clientKey == "" {
		return true, 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.entries[clientKey]
	if !ok {
		return true, 0
	}
	if !now.Before(entry.windowEnd) {
		delete(l.entries, clientKey)
		return true, 0
	}
	if entry.failures >= l.maxFailures {
		retrySec := int(math.Ceil(entry.windowEnd.Sub(now).Seconds()))
		if retrySec < 1 {
			retrySec = 1
		}
		return false, retrySec
	}
	return true, 0
}

func (l *loginRateLimiter) recordFailure(clientKey string, now time.Time) {
	if l == nil || clientKey == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.entries[clientKey]
	if !ok || !now.Before(entry.windowEnd) {
		if len(l.entries) >= maxRateLimitEntries {
			l.evictLocked(now)
		}
		l.entries[clientKey] = loginAttemptWindow{
			failures:  1,
			windowEnd: now.Add(l.window),
		}
		return
	}
	entry.failures++
	l.entries[clientKey] = entry
}

func (l *loginRateLimiter) reset(clientKey string) {
	if l == nil || clientKey == "" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, clientKey)
}

func (l *loginRateLimiter) evictLocked(now time.Time) {
	for k, v := range l.entries {
		if !now.Before(v.windowEnd) {
			delete(l.entries, k)
		}
	}
	if len(l.entries) < maxRateLimitEntries {
		return
	}
	var oldestKey string
	var oldestEnd time.Time
	first := true
	for k, v := range l.entries {
		if first || v.windowEnd.Before(oldestEnd) {
			oldestKey = k
			oldestEnd = v.windowEnd
			first = false
		}
	}
	if !first {
		delete(l.entries, oldestKey)
	}
}

func clientRateLimitKey(r *http.Request) string {
	if r == nil {
		return "unknown"
	}
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		first, _, _ := strings.Cut(xff, ",")
		if ip := strings.TrimSpace(first); ip != "" {
			return ip
		}
	}
	if xri := strings.TrimSpace(r.Header.Get("X-Real-IP")); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	if addr := strings.TrimSpace(r.RemoteAddr); addr != "" {
		return addr
	}
	return "unknown"
}

package opencode

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	base62Chars       = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	sessionPrefixMask = (uint64(1) << 48) - 1
	maxSessionEntries = 1000
	sessionEntryTTL   = 30 * time.Minute
	maxScopeTextLen   = 4000

	// OverloadRetryAttempts mirrors OpenCode CLI's free-lane load-shedding retry budget (~20s).
	OverloadRetryAttempts = 8
	// OverloadBaseBackoff is the linear backoff base for transient 429 overload responses.
	OverloadBaseBackoff = 800 * time.Millisecond
	// OverloadMaxBackoff caps each individual backoff sleep.
	OverloadMaxBackoff = 4 * time.Second
	// DefaultStartGap enforces minimum spacing between free-lane request starts to avoid burst 429s.
	DefaultStartGap = 400 * time.Millisecond
)

var (
	canonicalSessionRe = regexp.MustCompile(`^ses_[0-9a-f]{12}[0-9A-Za-z]{14}$`)
	overloadErrorRe    = regexp.MustCompile(`(?i)overloaded`)
)

// RejectionKind classifies upstream 403/429 responses for OpenCode Free.
type RejectionKind int

const (
	// RejectionNone indicates a response that does not require session eviction or overload retry.
	RejectionNone RejectionKind = iota
	// RejectionOverload indicates transient free-lane load shedding worth retrying with backoff.
	RejectionOverload
	// RejectionSession indicates a spent or flagged session (403 or non-overload 429).
	RejectionSession
)

var (
	sessionGenMu      sync.Mutex
	sessionGenLastMS  int64
	sessionGenCounter uint64
)

// GenerateSessionID mints a canonical OpenCode v2 session ID (`ses_<12hex><14base62>`)
// with a real monotonic counter (0..15 per ms) so upstream timestamp validation passes.
func GenerateSessionID(now time.Time) string {
	ms := max(now.UnixMilli(), 0)

	sessionGenMu.Lock()
	if ms != sessionGenLastMS {
		sessionGenLastMS = ms
		sessionGenCounter = 0
	}
	sessionGenCounter = (sessionGenCounter + 1) & 0x0F
	counterVal := sessionGenCounter
	sessionGenMu.Unlock()

	prefix := (^(uint64(ms)*0x1000 + counterVal)) & sessionPrefixMask
	buf := make([]byte, 14)
	_, _ = rand.Read(buf)
	var suffix strings.Builder
	suffix.Grow(14)
	for _, b := range buf {
		suffix.WriteByte(base62Chars[int(b)%len(base62Chars)])
	}
	return fmt.Sprintf("ses_%012x%s", prefix, suffix.String())
}

// IsCanonicalSessionID reports whether s matches OpenCode's canonical ses_* format.
func IsCanonicalSessionID(s string) bool {
	return canonicalSessionRe.MatchString(strings.TrimSpace(s))
}

type sessionEntry struct {
	sessionID string
	lastUsed  time.Time
}

// SessionStore maintains bounded, TTL-evicted canonical sessions per
// (connectionID, clientIdentity, conversationScope) so multi-turn chats reuse
// one session and preserve upstream prefix caching without burning quota.
type SessionStore struct {
	mu      sync.Mutex
	entries map[string]*sessionEntry
	ttl     time.Duration
	maxSize int
}

// NewSessionStore constructs a bounded SessionStore with default TTL and capacity.
func NewSessionStore() *SessionStore {
	return &SessionStore{
		entries: make(map[string]*sessionEntry),
		ttl:     sessionEntryTTL,
		maxSize: maxSessionEntries,
	}
}

// Resolve returns a valid canonical OpenCode session ID for the request.
// Native valid x-opencode-session headers from genuine OpenCode CLIs pass through verbatim.
func (s *SessionStore) Resolve(connectionID, clientIdentity string, incoming http.Header, reqBody []byte) string {
	if incoming != nil {
		if native := strings.TrimSpace(incoming.Get("x-opencode-session")); IsCanonicalSessionID(native) {
			return native
		}
	}

	client := strings.TrimSpace(clientIdentity)
	if client == "" && incoming != nil {
		if ua := strings.TrimSpace(incoming.Get("User-Agent")); ua != "" {
			sum := sha256.Sum256([]byte(ua))
			client = hex.EncodeToString(sum[:8])
		}
	}
	if client == "" {
		client = "generic"
	}

	scope := deriveConversationScope(reqBody)
	key := strings.TrimSpace(connectionID) + "\x00" + client + "\x00" + scope

	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, ok := s.entries[key]; ok {
		if now.Sub(entry.lastUsed) <= s.ttl {
			entry.lastUsed = now
			return entry.sessionID
		}
		delete(s.entries, key)
	}

	if len(s.entries) >= s.maxSize {
		s.evictLocked(now)
	}

	sid := GenerateSessionID(now)
	s.entries[key] = &sessionEntry{
		sessionID: sid,
		lastUsed:  now,
	}
	return sid
}

// Forget removes any cached mapping pointing to sessionID (matched by value so
// concurrent replacements survive). Returns true if at least one entry was dropped.
func (s *SessionStore) Forget(sessionID string) bool {
	sid := strings.TrimSpace(sessionID)
	if sid == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	removed := false
	for k, entry := range s.entries {
		if entry.sessionID == sid {
			delete(s.entries, k)
			removed = true
		}
	}
	return removed
}

func (s *SessionStore) evictLocked(now time.Time) {
	// First sweep expired entries.
	for k, entry := range s.entries {
		if now.Sub(entry.lastUsed) > s.ttl {
			delete(s.entries, k)
		}
	}
	if len(s.entries) < s.maxSize {
		return
	}
	// Evict least-recently-used entry.
	var victimKey string
	var oldest time.Time
	for k, entry := range s.entries {
		if victimKey == "" || entry.lastUsed.Before(oldest) {
			victimKey = k
			oldest = entry.lastUsed
		}
	}
	if victimKey != "" {
		delete(s.entries, victimKey)
	}
}

// deriveConversationScope hashes only the conversation prefix (system prompt +
// first user turn) so later turns in the same thread keep the same scope and
// preserve upstream prefix cache.
func deriveConversationScope(reqBody []byte) string {
	if len(reqBody) == 0 {
		return ""
	}
	var payload map[string]any
	if err := json.Unmarshal(reqBody, &payload); err != nil {
		return ""
	}

	var sysBuf strings.Builder
	switch sys := payload["system"].(type) {
	case string:
		sysBuf.WriteString(sys)
		sysBuf.WriteByte('\n')
	case []any:
		for _, block := range sys {
			sysBuf.WriteString(extractTextPart(block))
		}
		sysBuf.WriteByte('\n')
	}
	if inst, ok := payload["instructions"].(string); ok && inst != "" {
		sysBuf.WriteString(inst)
		sysBuf.WriteByte('\n')
	}

	var arr []any
	if msgs, ok := payload["messages"].([]any); ok {
		arr = msgs
	} else if input, ok := payload["input"].([]any); ok {
		arr = input
	}

	var firstUser string
	for _, raw := range arr {
		msg, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		role, _ := msg["role"].(string)
		switch role {
		case "system", "developer":
			sysBuf.WriteString(extractContentText(msg["content"]))
			sysBuf.WriteByte('\n')
		case "user":
			if firstUser == "" {
				firstUser = extractContentText(msg["content"])
			}
		}
		if firstUser != "" {
			break
		}
	}

	sysText := truncateString(sysBuf.String(), maxScopeTextLen)
	userText := truncateString(firstUser, maxScopeTextLen)
	if strings.TrimSpace(sysText) == "" && strings.TrimSpace(userText) == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(sysText + "\x01" + userText))
	return hex.EncodeToString(sum[:12])
}

func extractContentText(content any) string {
	switch v := content.(type) {
	case string:
		return v
	case []any:
		var b strings.Builder
		for _, part := range v {
			b.WriteString(extractTextPart(part))
		}
		return b.String()
	default:
		return ""
	}
}

func extractTextPart(part any) string {
	switch v := part.(type) {
	case string:
		return v
	case map[string]any:
		if t, ok := v["text"].(string); ok {
			return t
		}
	}
	return ""
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// Pacer enforces a minimum interval between free-lane request starts to avoid
// rapid burst 429s from Cloudflare while remaining strictly ctx-cancellable.
type Pacer struct {
	mu       sync.Mutex
	minGap   time.Duration
	nextSlot time.Time
}

// NewPacer constructs a Pacer with the given minimum inter-request gap.
func NewPacer(minGap time.Duration) *Pacer {
	if minGap <= 0 {
		minGap = DefaultStartGap
	}
	return &Pacer{minGap: minGap}
}

// Wait blocks until the caller's paced slot arrives or ctx is canceled.
func (p *Pacer) Wait(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p.mu.Lock()
	now := time.Now()
	slot := p.nextSlot
	if slot.Before(now) {
		slot = now
	}
	p.nextSlot = slot.Add(p.minGap)
	waitDur := slot.Sub(now)
	p.mu.Unlock()

	if waitDur <= 0 {
		return nil
	}
	timer := time.NewTimer(waitDur)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// ClassifyRejection categorizes OpenCode 403/429 responses into overload retry vs session eviction.
func ClassifyRejection(statusCode int, body []byte) RejectionKind {
	switch statusCode {
	case http.StatusForbidden:
		return RejectionSession
	case http.StatusTooManyRequests:
		if overloadErrorRe.Match(body) {
			return RejectionOverload
		}
		return RejectionSession
	default:
		return RejectionNone
	}
}

// OverloadBackoff returns the backoff duration for retry attempt (1-based).
func OverloadBackoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	d := time.Duration(attempt) * OverloadBaseBackoff
	if d > OverloadMaxBackoff {
		d = OverloadMaxBackoff
	}
	return d
}

// ParseResetAfter extracts a cooldown duration from OpenCode FreeUsageLimitError JSON
// (`error.resets_at`, `error.resets_in_seconds`) or the HTTP `Retry-After` header.
func ParseResetAfter(headers http.Header, body []byte, now time.Time) time.Duration {
	if len(body) > 0 {
		var parsed struct {
			Error struct {
				ResetsAt        float64 `json:"resets_at"`
				ResetsInSeconds float64 `json:"resets_in_seconds"`
			} `json:"error"`
		}
		if json.Unmarshal(body, &parsed) == nil {
			if parsed.Error.ResetsAt > 0 {
				target := time.Unix(int64(parsed.Error.ResetsAt), 0)
				if diff := target.Sub(now); diff > 0 {
					return diff
				}
			}
			if parsed.Error.ResetsInSeconds > 0 {
				return time.Duration(parsed.Error.ResetsInSeconds * float64(time.Second))
			}
		}
	}
	if headers != nil {
		if raw := strings.TrimSpace(headers.Get("Retry-After")); raw != "" {
			if sec, err := strconv.Atoi(raw); err == nil && sec > 0 {
				return time.Duration(sec) * time.Second
			}
		}
	}
	return 0
}

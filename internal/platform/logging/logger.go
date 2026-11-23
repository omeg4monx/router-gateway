// Package logging provides structured JSON logging via log/slog paired with an
// in-memory circular ring buffer for recent log inspection.
package logging

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"
)

// LogEntry represents a captured structured log record stored in a RingBuffer.
type LogEntry struct {
	Time    time.Time      `json:"time"`
	Level   string         `json:"level"`
	Message string         `json:"message"`
	Attrs   map[string]any `json:"attrs,omitempty"`
}

const subscriberBufferCap = 64

// RingBuffer stores up to a fixed capacity of the most recent LogEntry values
// in a mutex-protected circular buffer and supports non-blocking live log
// subscriptions.
type RingBuffer struct {
	mu      sync.Mutex
	entries []LogEntry
	next    int
	count   int
	cap     int

	nextSubID uint64
	subs      map[uint64]chan LogEntry
}

// Subscribe registers a live log listener returning a read-only channel of
// newly recorded LogEntry values and an idempotent unsubscribe function.
// Slow subscribers drop excess entries rather than blocking loggers.
func (r *RingBuffer) Subscribe() (entries <-chan LogEntry, unsubscribe func()) {
	ch := make(chan LogEntry, subscriberBufferCap)
	if r == nil {
		close(ch)
		return ch, func() {}
	}

	r.mu.Lock()
	if r.subs == nil {
		r.subs = make(map[uint64]chan LogEntry)
	}
	r.nextSubID++
	id := r.nextSubID
	r.subs[id] = ch
	r.mu.Unlock()

	var once sync.Once
	unsub := func() {
		once.Do(func() {
			r.mu.Lock()
			if existing, ok := r.subs[id]; ok {
				delete(r.subs, id)
				close(existing)
			}
			r.mu.Unlock()
		})
	}
	return ch, unsub
}

// Recent returns up to limit most recent log entries in chronological order.
// If limit <= 0, it returns an empty slice.
func (r *RingBuffer) Recent(limit int) []LogEntry {
	if r == nil || limit <= 0 {
		return []LogEntry{}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == 0 {
		return []LogEntry{}
	}

	n := min(limit, r.count)
	out := make([]LogEntry, n)
	start := (r.next - n + r.cap) % r.cap
	for i := range n {
		entry := r.entries[(start+i)%r.cap]
		entry.Attrs = cloneAttrMap(entry.Attrs)
		out[i] = entry
	}
	return out
}

func cloneAttrMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		if nested, ok := v.(map[string]any); ok {
			dst[k] = cloneAttrMap(nested)
		} else {
			dst[k] = v
		}
	}
	return dst
}

func (r *RingBuffer) add(entry LogEntry) {
	if r == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cap > 0 {
		r.entries[r.next] = entry
		r.next = (r.next + 1) % r.cap
		if r.count < r.cap {
			r.count++
		}
	}

	for _, ch := range r.subs {
		cloned := entry
		cloned.Attrs = cloneAttrMap(entry.Attrs)
		select {
		case ch <- cloned:
		default:
			// Drop when subscriber cannot keep up so logging never blocks.
		}
	}
}

// New constructs a *slog.Logger writing JSON records to w (or io.Discard if w is
// nil) at the parsed log level, while teeing every enabled record into the
// returned RingBuffer of capacity ringCap.
func New(level string, w io.Writer, ringCap int) (*slog.Logger, *RingBuffer) {
	if w == nil {
		w = io.Discard
	}
	if ringCap < 0 {
		ringCap = 0
	}

	ring := &RingBuffer{
		entries: make([]LogEntry, ringCap),
		cap:     ringCap,
	}

	jsonHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: parseLevel(level),
	})

	return slog.New(&teeHandler{
		jsonHandler: jsonHandler,
		ring:        ring,
	}), ring
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type boundAttrs struct {
	groups []string
	attrs  []slog.Attr
}

type teeHandler struct {
	jsonHandler slog.Handler
	ring        *RingBuffer
	groups      []string
	bound       []boundAttrs
}

func (h *teeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.jsonHandler.Enabled(ctx, level)
}

func (h *teeHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := make(map[string]any)
	for _, b := range h.bound {
		for _, a := range b.attrs {
			addAttrAt(attrs, b.groups, a)
		}
	}
	r.Attrs(func(a slog.Attr) bool {
		addAttrAt(attrs, h.groups, a)
		return true
	})

	h.ring.add(LogEntry{
		Time:    r.Time,
		Level:   r.Level.String(),
		Message: r.Message,
		Attrs:   attrs,
	})

	return h.jsonHandler.Handle(ctx, r)
}

func (h *teeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	nextBound := make([]boundAttrs, len(h.bound)+1)
	copy(nextBound, h.bound)
	nextBound[len(h.bound)] = boundAttrs{
		groups: slices.Clone(h.groups),
		attrs:  slices.Clone(attrs),
	}

	return &teeHandler{
		jsonHandler: h.jsonHandler.WithAttrs(attrs),
		ring:        h.ring,
		groups:      h.groups,
		bound:       nextBound,
	}
}

func (h *teeHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	nextGroups := make([]string, len(h.groups)+1)
	copy(nextGroups, h.groups)
	nextGroups[len(h.groups)] = name

	return &teeHandler{
		jsonHandler: h.jsonHandler.WithGroup(name),
		ring:        h.ring,
		groups:      nextGroups,
		bound:       h.bound,
	}
}

func addAttrAt(dst map[string]any, path []string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}

	if a.Value.Kind() == slog.KindGroup {
		sub := a.Value.Group()
		if len(sub) == 0 {
			return
		}
		nextPath := path
		if a.Key != "" {
			nextPath = append(slices.Clone(path), a.Key)
		}
		for _, sa := range sub {
			addAttrAt(dst, nextPath, sa)
		}
		return
	}

	target := ensurePath(dst, path)
	target[a.Key] = attrValue(a.Value)
}

func ensurePath(root map[string]any, path []string) map[string]any {
	cur := root
	for _, seg := range path {
		if existing, ok := cur[seg].(map[string]any); ok {
			cur = existing
			continue
		}
		next := make(map[string]any)
		cur[seg] = next
		cur = next
	}
	return cur
}

func attrValue(v slog.Value) any {
	switch v.Kind() {
	case slog.KindString:
		return v.String()
	case slog.KindInt64:
		return v.Int64()
	case slog.KindUint64:
		return v.Uint64()
	case slog.KindFloat64:
		return v.Float64()
	case slog.KindBool:
		return v.Bool()
	case slog.KindDuration:
		return v.Duration()
	case slog.KindTime:
		return v.Time()
	default:
		a := v.Any()
		if _, jm := a.(json.Marshaler); !jm {
			if err, ok := a.(error); ok {
				return err.Error()
			}
		}
		return a
	}
}

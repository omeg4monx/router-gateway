package logging_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/logging"
)

func TestRingBufferCapsAndPreservesOrder(t *testing.T) {
	var buf bytes.Buffer
	logger, ring := logging.New("info", &buf, 3)

	messages := []string{"first", "second", "third", "fourth", "fifth"}
	for i, msg := range messages {
		logger.Info(msg, "seq", i+1)
	}

	entries := ring.Recent(10)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	want := []string{"third", "fourth", "fifth"}
	for i, w := range want {
		if entries[i].Message != w {
			t.Errorf("entry[%d].Message = %q, want %q", i, entries[i].Message, w)
		}
		if entries[i].Level != "INFO" {
			t.Errorf("entry[%d].Level = %q, want %q", i, entries[i].Level, "INFO")
		}
		if entries[i].Time.IsZero() {
			t.Errorf("entry[%d].Time is zero", i)
		}
		wantSeq := int64(i + 3)
		if gotSeq, ok := entries[i].Attrs["seq"]; !ok || gotSeq != wantSeq {
			t.Errorf("entry[%d].Attrs[seq] = %v (%T), want %d", i, gotSeq, gotSeq, wantSeq)
		}
	}

	// Verify Recent(2) returns the last 2 entries in chronological order.
	lastTwo := ring.Recent(2)
	if len(lastTwo) != 2 {
		t.Fatalf("expected 2 entries from Recent(2), got %d", len(lastTwo))
	}
	if lastTwo[0].Message != "fourth" || lastTwo[1].Message != "fifth" {
		t.Errorf("Recent(2) = [%q, %q], want [fourth, fifth]", lastTwo[0].Message, lastTwo[1].Message)
	}

	// Verify JSON lines were also written to w.
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 JSON log lines, got %d", len(lines))
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(lines[4]), &parsed); err != nil {
		t.Fatalf("unmarshal last JSON line: %v", err)
	}
	if parsed["msg"] != "fifth" {
		t.Errorf("JSON msg = %v, want fifth", parsed["msg"])
	}
}

func TestLevelFilteringAndChildLoggerAttrsAndGroups(t *testing.T) {
	var buf bytes.Buffer
	logger, ring := logging.New("warn", &buf, 10)

	// Below threshold -> must not be recorded in ring or writer.
	logger.Debug("ignored debug")
	logger.Info("ignored info")
	if got := ring.Recent(10); len(got) != 0 {
		t.Fatalf("expected 0 entries below warn level, got %d", len(got))
	}

	child := logger.With("component", "gateway").WithGroup("http").With("method", "POST")
	child.Warn("upstream slow", "status", 504, slog.Group("timing", slog.Int("ms", 1200)))

	// Sibling child without group must not inherit "http" group.
	sibling := logger.With("component", "auth")
	sibling.Error("token invalid")

	entries := ring.Recent(10)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	first := entries[0]
	if first.Level != "WARN" || first.Message != "upstream slow" {
		t.Fatalf("unexpected first entry: %+v", first)
	}
	if first.Attrs["component"] != "gateway" {
		t.Errorf("Attrs[component] = %v, want gateway", first.Attrs["component"])
	}
	httpGroup, ok := first.Attrs["http"].(map[string]any)
	if !ok {
		t.Fatalf("Attrs[http] = %T (%v), want map[string]any", first.Attrs["http"], first.Attrs["http"])
	}
	if httpGroup["method"] != "POST" {
		t.Errorf("http.method = %v, want POST", httpGroup["method"])
	}
	if httpGroup["status"] != int64(504) {
		t.Errorf("http.status = %v (%T), want 504", httpGroup["status"], httpGroup["status"])
	}
	timingGroup, ok := httpGroup["timing"].(map[string]any)
	if !ok || timingGroup["ms"] != int64(1200) {
		t.Errorf("http.timing = %v, want map with ms=1200", httpGroup["timing"])
	}

	second := entries[1]
	if second.Attrs["component"] != "auth" {
		t.Errorf("sibling Attrs[component] = %v, want auth", second.Attrs["component"])
	}
	if _, exists := second.Attrs["http"]; exists {
		t.Errorf("sibling leaked parent/other group 'http': %v", second.Attrs)
	}
}

func TestRingBufferConcurrentWritesAndReads(t *testing.T) {
	logger, ring := logging.New("debug", nil, 16)

	var wg sync.WaitGroup
	for g := range 8 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			child := logger.With("worker", id)
			for i := range 50 {
				child.Debug("tick", "iter", i)
				_ = ring.Recent(5)
			}
		}(g)
	}
	wg.Wait()

	entries := ring.Recent(100)
	if len(entries) != 16 {
		t.Fatalf("expected capped 16 entries, got %d", len(entries))
	}
	if empty := ring.Recent(0); len(empty) != 0 {
		t.Errorf("Recent(0) = %d entries, want 0", len(empty))
	}
}

type jsonMarshalerError struct {
	Code int `json:"code"`
}

func (e jsonMarshalerError) Error() string {
	return fmt.Sprintf("code=%d", e.Code)
}

func (e jsonMarshalerError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code int `json:"code"`
	}{Code: e.Code})
}

func TestErrorAttrsSerializeAsMessageInRingBufferAndJSON(t *testing.T) {
	var buf bytes.Buffer
	logger, ring := logging.New("info", &buf, 10)

	wrapped := fmt.Errorf("upstream failed: %w", errors.New("connection refused"))
	logger.Error("request failed",
		"err", errors.New("boom"),
		"wrapped", wrapped,
		"structured", jsonMarshalerError{Code: 502},
	)

	entries := ring.Recent(10)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].Attrs["err"] != "boom" {
		t.Errorf("Attrs[err] = %#v, want %q", entries[0].Attrs["err"], "boom")
	}
	if entries[0].Attrs["wrapped"] != "upstream failed: connection refused" {
		t.Errorf("Attrs[wrapped] = %#v, want %q", entries[0].Attrs["wrapped"], "upstream failed: connection refused")
	}

	raw, err := json.Marshal(entries[0])
	if err != nil {
		t.Fatalf("json.Marshal(LogEntry): %v", err)
	}
	var decoded struct {
		Attrs map[string]any `json:"attrs"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("json.Unmarshal(LogEntry): %v", err)
	}
	if decoded.Attrs["err"] != "boom" {
		t.Errorf("JSON attrs.err = %#v, want %q", decoded.Attrs["err"], "boom")
	}
	if decoded.Attrs["wrapped"] != "upstream failed: connection refused" {
		t.Errorf("JSON attrs.wrapped = %#v, want %q", decoded.Attrs["wrapped"], "upstream failed: connection refused")
	}
	structured, ok := decoded.Attrs["structured"].(map[string]any)
	if !ok || structured["code"] != float64(502) {
		t.Errorf("JSON attrs.structured = %#v, want {\"code\":502}", decoded.Attrs["structured"])
	}
}

func TestRingBufferRecentDeepClonesAttrs(t *testing.T) {
	logger, ring := logging.New("info", nil, 5)
	logger.With("top", "original").WithGroup("http").Info("ok", "method", "GET")

	first := ring.Recent(10)
	if len(first) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(first))
	}

	// Mutate both top-level and nested group maps on the returned slice.
	first[0].Attrs["top"] = "mutated"
	first[0].Attrs["injected"] = "bad"
	httpMap := first[0].Attrs["http"].(map[string]any)
	httpMap["method"] = "DELETE"

	second := ring.Recent(10)
	if second[0].Attrs["top"] != "original" {
		t.Errorf("buffer top attr mutated: got %v, want original", second[0].Attrs["top"])
	}
	if _, exists := second[0].Attrs["injected"]; exists {
		t.Errorf("buffer leaked injected key into stored attrs")
	}
	gotHTTP := second[0].Attrs["http"].(map[string]any)
	if gotHTTP["method"] != "GET" {
		t.Errorf("buffer nested http.method mutated: got %v, want GET", gotHTTP["method"])
	}

	// Concurrent Recent() + caller-side map mutation must remain race-free.
	var wg sync.WaitGroup
	for g := range 4 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for range 50 {
				snap := ring.Recent(5)
				if len(snap) > 0 {
					snap[0].Attrs["top"] = id
					if hm, ok := snap[0].Attrs["http"].(map[string]any); ok {
						hm["method"] = "PATCH"
					}
				}
			}
		}(g)
	}
	wg.Wait()
}

func TestRingBufferSubscribeNonBlockingAndUnsubscribe(t *testing.T) {
	logger, ring := logging.New("debug", nil, 32)

	ch, unsub := ring.Subscribe()
	defer unsub()

	// Emit far more entries than subscriber channel buffer capacity without reading.
	// Logger must never block even when subscriber is stalled.
	done := make(chan struct{})
	go func() {
		for i := range 200 {
			logger.Info("burst", "seq", i)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("logger.Info blocked on slow RingBuffer subscriber")
	}

	// Drain at least one entry and verify deep-cloned attrs.
	select {
	case entry := <-ch:
		if entry.Message != "burst" {
			t.Fatalf("unexpected entry message %q", entry.Message)
		}
	default:
		t.Fatal("expected buffered entries on subscriber channel")
	}

	// Unsubscribe must be idempotent and close the channel cleanly.
	unsub()
	unsub()
	for range ch {
	}

	// Logging after unsubscribe must not panic.
	logger.Warn("after-unsub")
}

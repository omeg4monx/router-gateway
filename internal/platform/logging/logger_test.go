package logging_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"testing"

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

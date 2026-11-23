package system

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParseProcHelpers(t *testing.T) {
	cpuinfo := "processor\t: 0\nvendor_id\t: GenuineIntel\nmodel name\t: Intel(R) Xeon(R) Gold 6248 CPU @ 2.50GHz\n"
	if got := parseCPUModel(cpuinfo); got != "Intel(R) Xeon(R) Gold 6248 CPU @ 2.50GHz" {
		t.Fatalf("parseCPUModel = %q, want Intel Xeon", got)
	}

	meminfo := "MemTotal:       16384000 kB\nMemFree:         2048000 kB\nMemAvailable:    6144000 kB\n"
	total, used := parseMemInfo(meminfo)
	if total != 16384000*1024 {
		t.Fatalf("mem total = %d, want %d", total, 16384000*1024)
	}
	if used != (16384000-6144000)*1024 {
		t.Fatalf("mem used = %d, want %d", used, (16384000-6144000)*1024)
	}

	load := parseLoadAvg("0.42 0.35 0.28 1/312 12345\n")
	if load[0] != 0.42 || load[1] != 0.35 || load[2] != 0.28 {
		t.Fatalf("parseLoadAvg = %v, want [0.42 0.35 0.28]", load)
	}

	if got := parseHostUptime("123456.78 456789.12\n"); got != 123456 {
		t.Fatalf("parseHostUptime = %d, want 123456", got)
	}

	osRelease := "NAME=\"Ubuntu\"\nPRETTY_NAME=\"Ubuntu 24.04.2 LTS\"\n"
	if got := parseOSPrettyName(osRelease); got != "Ubuntu 24.04.2 LTS" {
		t.Fatalf("parseOSPrettyName = %q, want Ubuntu 24.04.2 LTS", got)
	}

	s1 := parseProcStatCPU("cpu  100 0 50 850 0 0 0 0 0 0\n")
	s2 := parseProcStatCPU("cpu  150 0 100 950 0 0 0 0 0 0\n")
	if pct := computeCPUPercent(s1, s2); pct != 50.0 {
		t.Fatalf("computeCPUPercent = %v, want 50.0", pct)
	}
	if pct := computeCPUPercent(s2, s1); pct != 0 {
		t.Fatalf("underflow guard computeCPUPercent = %v, want 0", pct)
	}
}

func TestTelemetryCacheAvoidsRedundantWork(t *testing.T) {
	var c telemetryCollector
	now := time.Now().UTC()
	_, _, t1 := c.snapshot(now)
	_, _, t2 := c.snapshot(now.Add(500 * time.Millisecond))
	if !t1.Equal(t2) {
		t.Fatalf("expected cached sample within TTL: t1=%v t2=%v", t1, t2)
	}
	_, _, t3 := c.snapshot(now.Add(telemetryCacheTTL + time.Millisecond))
	if t3.Equal(t1) {
		t.Fatalf("expected fresh sample after TTL expired")
	}
}

func TestHandleTelemetryReturnsHostSnapshot(t *testing.T) {
	m := New(Options{})
	req := httptest.NewRequest(http.MethodGet, "/api/system/telemetry", nil)
	rec := httptest.NewRecorder()

	m.handleTelemetry(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var resp systemTelemetryResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.CPUCores < 1 {
		t.Fatalf("cpuCores = %d, want >= 1", resp.CPUCores)
	}
	if resp.GoVersion == "" {
		t.Fatal("expected non-empty goVersion")
	}
	if resp.OS == "" {
		t.Fatal("expected non-empty os")
	}
	if resp.DBStatus != "error" && resp.DBStatus != "ok" {
		t.Fatalf("unexpected dbStatus = %q", resp.DBStatus)
	}
}

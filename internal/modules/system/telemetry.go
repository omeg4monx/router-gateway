package system

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

const (
	telemetryCacheTTL  = 1500 * time.Millisecond
	maxProcFileBytes   = 128 << 10 // 128 KiB hard cap against runaway reads
)

type systemTelemetryResponse struct {
	Hostname          string     `json:"hostname"`
	OS                string     `json:"os"`
	Arch              string     `json:"arch"`
	GoVersion         string     `json:"goVersion"`
	CPUModel          string     `json:"cpuModel"`
	CPUCores          int        `json:"cpuCores"`
	CPUUsagePercent   float64    `json:"cpuUsagePercent"`
	LoadAvg           [3]float64 `json:"loadAvg"`
	MemTotalBytes     uint64     `json:"memTotalBytes"`
	MemUsedBytes      uint64     `json:"memUsedBytes"`
	DiskTotalBytes    uint64     `json:"diskTotalBytes"`
	DiskUsedBytes     uint64     `json:"diskUsedBytes"`
	HostUptimeSeconds int64      `json:"hostUptimeSeconds"`
	AppUptimeSeconds  int64      `json:"appUptimeSeconds"`
	DBStatus          string     `json:"dbStatus"`
	SampledAt         string     `json:"sampledAt"`
}

type staticHostInfo struct {
	hostname string
	osLabel  string
	cpuModel string
	cpuCores int
}

type cpuTickSample struct {
	valid bool
	total uint64
	idle  uint64
}

type hostTelemetrySnapshot struct {
	cpuUsagePercent   float64
	loadAvg           [3]float64
	memTotalBytes     uint64
	memUsedBytes      uint64
	diskTotalBytes    uint64
	diskUsedBytes     uint64
	hostUptimeSeconds int64
}

type telemetryCollector struct {
	staticOnce sync.Once
	staticInfo staticHostInfo

	mu          sync.Mutex
	lastSampled time.Time
	cached      hostTelemetrySnapshot
	prevCPU     cpuTickSample
}

func (c *telemetryCollector) getStatic() staticHostInfo {
	c.staticOnce.Do(func() {
		hostname, _ := os.Hostname()
		hostname = strings.TrimSpace(hostname)
		if hostname == "" {
			hostname = "localhost"
		}

		prettyOS, cpuModel := collectStaticHostInfo()
		osLabel := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
		if prettyOS != "" {
			osLabel = fmt.Sprintf("%s · %s/%s", prettyOS, runtime.GOOS, runtime.GOARCH)
		}
		if cpuModel == "" {
			cpuModel = fmt.Sprintf("%s (%s)", runtime.GOOS, runtime.GOARCH)
		}

		c.staticInfo = staticHostInfo{
			hostname: hostname,
			osLabel:  osLabel,
			cpuModel: cpuModel,
			cpuCores: runtime.NumCPU(),
		}
	})
	return c.staticInfo
}

func (c *telemetryCollector) snapshot(now time.Time) (staticHostInfo, hostTelemetrySnapshot, time.Time) {
	info := c.getStatic()

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.lastSampled.IsZero() && now.Sub(c.lastSampled) < telemetryCacheTTL {
		return info, c.cached, c.lastSampled
	}

	snap, nextCPU := collectDynamicHostTelemetry(c.prevCPU, info.cpuCores)
	snap.cpuUsagePercent = clampPercent(snap.cpuUsagePercent)
	if snap.memUsedBytes > snap.memTotalBytes {
		snap.memUsedBytes = snap.memTotalBytes
	}
	if snap.diskUsedBytes > snap.diskTotalBytes {
		snap.diskUsedBytes = snap.diskTotalBytes
	}
	if snap.hostUptimeSeconds < 0 {
		snap.hostUptimeSeconds = 0
	}

	c.prevCPU = nextCPU
	c.cached = snap
	c.lastSampled = now
	return info, snap, now
}

func (m *Module) handleTelemetry(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	info, snap, sampledAt := m.telemetry.snapshot(now)

	dbStatus := "ok"
	if m.db == nil || m.db.Ping(r.Context()) != nil {
		dbStatus = "error"
	}

	w.Header().Set("Cache-Control", "no-store")
	httpserver.WriteJSON(w, http.StatusOK, systemTelemetryResponse{
		Hostname:          info.hostname,
		OS:                info.osLabel,
		Arch:              runtime.GOARCH,
		GoVersion:         runtime.Version(),
		CPUModel:          info.cpuModel,
		CPUCores:          info.cpuCores,
		CPUUsagePercent:   snap.cpuUsagePercent,
		LoadAvg:           snap.loadAvg,
		MemTotalBytes:     snap.memTotalBytes,
		MemUsedBytes:      snap.memUsedBytes,
		DiskTotalBytes:    snap.diskTotalBytes,
		DiskUsedBytes:     snap.diskUsedBytes,
		HostUptimeSeconds: snap.hostUptimeSeconds,
		AppUptimeSeconds:  m.uptimeSeconds(),
		DBStatus:          dbStatus,
		SampledAt:         sampledAt.Format(time.RFC3339),
	})
}

func clampPercent(v float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return math.Round(v*10) / 10
}

func computeCPUPercent(prev, curr cpuTickSample) float64 {
	if !prev.valid || !curr.valid || curr.total <= prev.total {
		return 0
	}
	deltaTotal := curr.total - prev.total
	var deltaIdle uint64
	if curr.idle >= prev.idle {
		deltaIdle = curr.idle - prev.idle
	}
	if deltaIdle >= deltaTotal {
		return 0
	}
	return clampPercent(float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100)
}

func parseProcStatCPU(raw string) cpuTickSample {
	for _, line := range strings.Split(raw, "\n") {
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			return cpuTickSample{}
		}
		var total, idle uint64
		for i := 1; i < len(fields); i++ {
			v, err := strconv.ParseUint(fields[i], 10, 64)
			if err != nil {
				continue
			}
			total += v
			// field 4 = idle, field 5 = iowait
			if i == 4 || i == 5 {
				idle += v
			}
		}
		return cpuTickSample{valid: total > 0, total: total, idle: idle}
	}
	return cpuTickSample{}
}

func parseCPUModel(raw string) string {
	for _, line := range strings.Split(raw, "\n") {
		idx := strings.IndexByte(line, ':')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		if key == "model name" || key == "Hardware" || key == "Processor" {
			if val := strings.TrimSpace(line[idx+1:]); val != "" {
				return val
			}
		}
	}
	return ""
}

func parseMemInfo(raw string) (totalBytes uint64, usedBytes uint64) {
	var memTotalKB, memAvailableKB, memFreeKB uint64
	var hasAvailable bool
	for _, line := range strings.Split(raw, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			memTotalKB = val
		case "MemAvailable:":
			memAvailableKB = val
			hasAvailable = true
		case "MemFree:":
			memFreeKB = val
		}
	}
	totalBytes = memTotalKB * 1024
	if hasAvailable && memTotalKB >= memAvailableKB {
		usedBytes = (memTotalKB - memAvailableKB) * 1024
	} else if memTotalKB >= memFreeKB {
		usedBytes = (memTotalKB - memFreeKB) * 1024
	}
	return totalBytes, usedBytes
}

func parseLoadAvg(raw string) [3]float64 {
	var out [3]float64
	fields := strings.Fields(raw)
	for i := 0; i < 3 && i < len(fields); i++ {
		if v, err := strconv.ParseFloat(fields[i], 64); err == nil && !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0 {
			out[i] = v
		}
	}
	return out
}

func parseHostUptime(raw string) int64 {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return 0
	}
	if v, err := strconv.ParseFloat(fields[0], 64); err == nil && !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0 {
		return int64(v)
	}
	return 0
}

func parseOSPrettyName(raw string) string {
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "PRETTY_NAME=") {
			val := strings.TrimPrefix(trimmed, "PRETTY_NAME=")
			return strings.Trim(val, `"'`)
		}
	}
	return ""
}

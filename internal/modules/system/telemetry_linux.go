//go:build linux

package system

import (
	"io"
	"os"
	"syscall"
)

func readBoundedFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer func() { _ = f.Close() }()
	data, err := io.ReadAll(io.LimitReader(f, maxProcFileBytes))
	if err != nil {
		return ""
	}
	return string(data)
}

func collectStaticHostInfo() (osPrettyName, cpuModel string) {
	osPrettyName = parseOSPrettyName(readBoundedFile("/etc/os-release"))
	cpuModel = parseCPUModel(readBoundedFile("/proc/cpuinfo"))
	return osPrettyName, cpuModel
}

func collectDynamicHostTelemetry(prevCPU cpuTickSample, cpuCores int) (hostTelemetrySnapshot, cpuTickSample) {
	var snap hostTelemetrySnapshot

	currCPU := parseProcStatCPU(readBoundedFile("/proc/stat"))
	if prevCPU.valid && currCPU.valid {
		snap.cpuUsagePercent = computeCPUPercent(prevCPU, currCPU)
	}

	snap.memTotalBytes, snap.memUsedBytes = parseMemInfo(readBoundedFile("/proc/meminfo"))
	snap.loadAvg = parseLoadAvg(readBoundedFile("/proc/loadavg"))
	snap.hostUptimeSeconds = parseHostUptime(readBoundedFile("/proc/uptime"))

	// On the very first sample before a delta exists, approximate CPU % from 1m loadavg.
	if !prevCPU.valid && cpuCores > 0 && snap.loadAvg[0] > 0 {
		snap.cpuUsagePercent = clampPercent((snap.loadAvg[0] / float64(cpuCores)) * 100)
	}

	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err == nil && stat.Bsize > 0 {
		bsize := uint64(stat.Bsize)
		total := stat.Blocks * bsize
		avail := stat.Bavail * bsize
		snap.diskTotalBytes = total
		if total >= avail {
			snap.diskUsedBytes = total - avail
		}
	}

	return snap, currCPU
}

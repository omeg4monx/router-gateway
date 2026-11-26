//go:build darwin

package system

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func collectStaticHostInfo() (osPrettyName string, cpuModel string) {
	if ver := runSysctl("kern.osproductversion"); ver != "" {
		osPrettyName = "macOS " + ver
	} else {
		osPrettyName = "macOS"
	}
	cpuModel = runSysctl("machdep.cpu.brand_string")
	return osPrettyName, cpuModel
}

func collectDynamicHostTelemetry(prevCPU cpuTickSample, cpuCores int) (hostTelemetrySnapshot, cpuTickSample) {
	var snap hostTelemetrySnapshot

	if memStr := runSysctl("hw.memsize"); memStr != "" {
		if total, err := strconv.ParseUint(memStr, 10, 64); err == nil {
			snap.memTotalBytes = total
			snap.memUsedBytes = estimateDarwinUsedMem(total)
		}
	}

	if loadStr := runSysctl("vm.loadavg"); loadStr != "" {
		cleaned := strings.Trim(loadStr, "{} ")
		snap.loadAvg = parseLoadAvg(cleaned)
		if cpuCores > 0 && snap.loadAvg[0] > 0 {
			snap.cpuUsagePercent = clampPercent((snap.loadAvg[0] / float64(cpuCores)) * 100)
		}
	}

	if bootStr := runSysctl("kern.boottime"); bootStr != "" {
		if idx := strings.Index(bootStr, "sec = "); idx >= 0 {
			rest := bootStr[idx+len("sec = "):]
			if comma := strings.IndexByte(rest, ','); comma > 0 {
				if sec, err := strconv.ParseInt(strings.TrimSpace(rest[:comma]), 10, 64); err == nil {
					if now := time.Now().Unix(); now > sec {
						snap.hostUptimeSeconds = now - sec
					}
				}
			}
		}
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

	return snap, prevCPU
}

func runSysctl(key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "sysctl", "-n", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func estimateDarwinUsedMem(totalBytes uint64) uint64 {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "vm_stat").Output()
	if err != nil {
		return 0
	}
	pageSize := uint64(4096)
	var active, wired, compressed uint64
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "page size of ") {
			fields := strings.Fields(line)
			for i, f := range fields {
				if f == "of" && i+1 < len(fields) {
					if ps, parseErr := strconv.ParseUint(fields[i+1], 10, 64); parseErr == nil && ps > 0 {
						pageSize = ps
					}
				}
			}
			continue
		}
		idx := strings.IndexByte(line, ':')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		valStr := strings.TrimRight(strings.TrimSpace(line[idx+1:]), ".")
		val, parseErr := strconv.ParseUint(valStr, 10, 64)
		if parseErr != nil {
			continue
		}
		switch key {
		case "Pages active":
			active = val
		case "Pages wired down":
			wired = val
		case "Pages occupied by compressor":
			compressed = val
		}
	}
	used := (active + wired + compressed) * pageSize
	if used > totalBytes {
		return totalBytes
	}
	return used
}

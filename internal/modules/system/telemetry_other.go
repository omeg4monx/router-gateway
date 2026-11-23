//go:build !linux && !darwin && !windows

package system

func collectStaticHostInfo() (osPrettyName string, cpuModel string) {
	return "", ""
}

func collectDynamicHostTelemetry(prevCPU cpuTickSample, _ int) (hostTelemetrySnapshot, cpuTickSample) {
	return hostTelemetrySnapshot{}, prevCPU
}

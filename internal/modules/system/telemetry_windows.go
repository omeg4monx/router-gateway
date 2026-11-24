//go:build windows

package system

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

type memoryStatusEx struct {
	length               uint32
	memoryLoad           uint32
	totalPhys            uint64
	availPhys            uint64
	totalPageFile        uint64
	availPageFile        uint64
	totalVirtual         uint64
	availVirtual         uint64
	availExtendedVirtual uint64
}

type winFileTime struct {
	lowDateTime  uint32
	highDateTime uint32
}

func (ft winFileTime) uint64() uint64 {
	return (uint64(ft.highDateTime) << 32) | uint64(ft.lowDateTime)
}

var (
	winKernel32              = syscall.NewLazyDLL("kernel32.dll")
	procGlobalMemoryStatusEx = winKernel32.NewProc("GlobalMemoryStatusEx")
	procGetTickCount64       = winKernel32.NewProc("GetTickCount64")
	procGetDiskFreeSpaceExW  = winKernel32.NewProc("GetDiskFreeSpaceExW")
	procGetSystemTimes       = winKernel32.NewProc("GetSystemTimes")
)

func collectStaticHostInfo() (osPrettyName string, cpuModel string) {
	osPrettyName = readWindowsRegistryString(
		syscall.HKEY_LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		"ProductName",
	)
	if osPrettyName == "" {
		osPrettyName = "Windows"
	}

	cpuModel = readWindowsRegistryString(
		syscall.HKEY_LOCAL_MACHINE,
		`HARDWARE\DESCRIPTION\System\CentralProcessor\0`,
		"ProcessorNameString",
	)
	if cpuModel == "" {
		cpuModel = strings.TrimSpace(os.Getenv("PROCESSOR_IDENTIFIER"))
	}
	return osPrettyName, cpuModel
}

func collectDynamicHostTelemetry(prevCPU cpuTickSample, _ int) (hostTelemetrySnapshot, cpuTickSample) {
	var snap hostTelemetrySnapshot
	var currCPU cpuTickSample

	if procGetSystemTimes.Find() == nil {
		var idle, kernel, user winFileTime
		r1, _, _ := procGetSystemTimes.Call(
			uintptr(unsafe.Pointer(&idle)),
			uintptr(unsafe.Pointer(&kernel)),
			uintptr(unsafe.Pointer(&user)),
		)
		if r1 != 0 {
			// On Windows, kernel time includes idle time.
			total := kernel.uint64() + user.uint64()
			currCPU = cpuTickSample{valid: total > 0, total: total, idle: idle.uint64()}
			if prevCPU.valid {
				snap.cpuUsagePercent = computeCPUPercent(prevCPU, currCPU)
			}
		}
	}

	if procGlobalMemoryStatusEx.Find() == nil {
		var mem memoryStatusEx
		mem.length = uint32(unsafe.Sizeof(mem))
		if r1, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&mem))); r1 != 0 {
			snap.memTotalBytes = mem.totalPhys
			if mem.totalPhys >= mem.availPhys {
				snap.memUsedBytes = mem.totalPhys - mem.availPhys
			}
		}
	}

	if procGetTickCount64.Find() == nil {
		if ms, _, _ := procGetTickCount64.Call(); ms > 0 {
			snap.hostUptimeSeconds = int64(ms / 1000)
		}
	}

	if procGetDiskFreeSpaceExW.Find() == nil {
		root := os.Getenv("SystemDrive")
		if root == "" {
			root = "C:"
		}
		root = filepath.Clean(root) + `\`
		if rootPtr, err := syscall.UTF16PtrFromString(root); err == nil {
			var freeAvail, totalBytes, totalFree uint64
			r1, _, _ := procGetDiskFreeSpaceExW.Call(
				uintptr(unsafe.Pointer(rootPtr)),
				uintptr(unsafe.Pointer(&freeAvail)),
				uintptr(unsafe.Pointer(&totalBytes)),
				uintptr(unsafe.Pointer(&totalFree)),
			)
			if r1 != 0 {
				snap.diskTotalBytes = totalBytes
				if totalBytes >= freeAvail {
					snap.diskUsedBytes = totalBytes - freeAvail
				}
			}
		}
	}

	return snap, currCPU
}

func readWindowsRegistryString(root syscall.Handle, subKey, valueName string) string {
	subKeyPtr, err := syscall.UTF16PtrFromString(subKey)
	if err != nil {
		return ""
	}
	var handle syscall.Handle
	if err := syscall.RegOpenKeyEx(root, subKeyPtr, 0, syscall.KEY_READ, &handle); err != nil {
		return ""
	}
	defer func() { _ = syscall.RegCloseKey(handle) }()

	valPtr, err := syscall.UTF16PtrFromString(valueName)
	if err != nil {
		return ""
	}
	var valType uint32
	var bufLen uint32
	if err := syscall.RegQueryValueEx(handle, valPtr, nil, &valType, nil, &bufLen); err != nil || bufLen == 0 {
		return ""
	}
	buf := make([]uint16, bufLen/2+1)
	if err := syscall.RegQueryValueEx(handle, valPtr, nil, &valType, (*byte)(unsafe.Pointer(&buf[0])), &bufLen); err != nil {
		return ""
	}
	return strings.TrimSpace(syscall.UTF16ToString(buf))
}

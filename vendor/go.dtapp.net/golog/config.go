package golog

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

type systemResult struct {
	SystemHostname      string  // 主机名
	SystemOs            string  // 系统类型
	SystemVersion       string  // 系统版本
	SystemKernel        string  // 系统内核
	SystemKernelVersion string  // 系统内核版本
	SystemUpTime        uint64  // 系统运行时间
	SystemBootTime      uint64  // 系统开机时间
	CpuCores            int     // CPU核数
	CpuModelName        string  // CPU型号名称
	CpuMhz              float64 // CPU兆赫
}

func getSystem() (result systemResult) {

	hInfo, _ := host.Info()

	result.SystemHostname = hInfo.Hostname
	result.SystemOs = hInfo.OS
	result.SystemVersion = hInfo.PlatformVersion
	result.SystemKernel = hInfo.KernelArch
	result.SystemKernelVersion = hInfo.KernelVersion
	result.SystemUpTime = hInfo.Uptime
	if hInfo.BootTime != 0 {
		result.SystemBootTime = hInfo.BootTime
	}

	hCpu, _ := cpu.Times(true)

	result.CpuCores = len(hCpu)

	cInfo, _ := cpu.Info()

	if len(cInfo) > 0 {
		result.CpuModelName = cInfo[0].ModelName
		result.CpuMhz = cInfo[0].Mhz
	}

	return result
}

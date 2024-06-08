package golog

import (
	"context"
	"go.dtapp.net/gorequest"
	"runtime"
)

func (gg *GinGorm) setConfig(ctx context.Context) {

	gg.config.GoVersion = runtime.Version()
	gg.config.SdkVersion = Version

	info := getSystem()
	gg.config.system.SystemVersion = info.SystemVersion
	gg.config.system.SystemOs = info.SystemOs
	gg.config.system.SystemArch = info.SystemKernel
	gg.config.system.SystemInsideIP = gorequest.GetInsideIp(ctx)
	gg.config.system.SystemCpuModel = info.CpuModelName
	gg.config.system.SystemCpuCores = info.CpuCores
	gg.config.system.SystemCpuMhz = info.CpuMhz

}

package golog

import (
	"context"
	"go.dtapp.net/gorequest"
	"runtime"
)

func (hg *HertzGorm) setConfig(ctx context.Context) {

	hg.config.GoVersion = runtime.Version()
	hg.config.SdkVersion = Version

	info := getSystem()
	hg.config.system.SystemVersion = info.SystemVersion
	hg.config.system.SystemOs = info.SystemOs
	hg.config.system.SystemArch = info.SystemKernel
	hg.config.system.SystemInsideIP = gorequest.GetInsideIp(ctx)
	hg.config.system.SystemCpuModel = info.CpuModelName
	hg.config.system.SystemCpuCores = info.CpuCores
	hg.config.system.SystemCpuMhz = info.CpuMhz

}

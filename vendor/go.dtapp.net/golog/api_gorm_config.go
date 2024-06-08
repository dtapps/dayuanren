package golog

import (
	"context"
	"go.dtapp.net/gorequest"
	"runtime"
)

func (ag *ApiGorm) setConfig(ctx context.Context) {

	ag.config.GoVersion = runtime.Version()
	ag.config.SdkVersion = Version

	info := getSystem()
	ag.config.system.SystemVersion = info.SystemVersion
	ag.config.system.SystemOs = info.SystemOs
	ag.config.system.SystemArch = info.SystemKernel
	ag.config.system.SystemInsideIP = gorequest.GetInsideIp(ctx)
	ag.config.system.SystemCpuModel = info.CpuModelName
	ag.config.system.SystemCpuCores = info.CpuCores
	ag.config.system.SystemCpuMhz = info.CpuMhz

}

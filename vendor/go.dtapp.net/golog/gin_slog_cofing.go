package golog

// ConfigSLogClientFun 日志配置
func (gl *GinSLog) ConfigSLogClientFun(sLogFun SLogFun) {
	sLog := sLogFun()
	if sLog != nil {
		gl.slog.client = sLog
		gl.slog.status = true
	}
}

// ConfigSLogResultClientFun 日志配置然后返回
func (gl *GinSLog) ConfigSLogResultClientFun(sLogFun SLogFun) *GinSLog {
	sLog := sLogFun()
	if sLog != nil {
		gl.slog.client = sLog
		gl.slog.status = true
	}
	return gl
}

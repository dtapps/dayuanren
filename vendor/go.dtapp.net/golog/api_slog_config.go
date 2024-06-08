package golog

// ConfigSLogClientFun 日志配置
func (al *ApiSLog) ConfigSLogClientFun(sLogFun SLogFun) {
	sLog := sLogFun()
	if sLog != nil {
		al.slog.client = sLog
		al.slog.status = true
	}
}

// ConfigSLogResultClientFun 日志配置然后返回
func (al *ApiSLog) ConfigSLogResultClientFun(sLogFun SLogFun) *ApiSLog {
	sLog := sLogFun()
	if sLog != nil {
		al.slog.client = sLog
		al.slog.status = true
	}
	return al
}

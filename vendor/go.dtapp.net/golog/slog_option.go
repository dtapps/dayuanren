package golog

import "gopkg.in/natefinch/lumberjack.v2"

type SLogOption func(*SLog)

// WithSLogLumberjack Lumberjack配置
// Filename 日志文件的位置
// MaxSize 文件最大尺寸（以MB为单位）
// MaxAge 留旧文件的最大天数
// MaxBackups 保留的最大旧文件数量
// Compress 是否压缩/归档旧文件
// LocalTime 使用本地时间创建时间戳
func WithSLogLumberjack(config lumberjack.Logger) SLogOption {
	return func(sl *SLog) {
		sl.option.lumberjackConfig = config
		sl.option.lumberjackConfigStatus = true
	}
}

// WithSLogShowLine 显示代码行
func WithSLogShowLine() SLogOption {
	return func(sl *SLog) {
		sl.option.showLine = true
	}
}

// WithSLogSetDefault 设置为默认的实例
func WithSLogSetDefault() SLogOption {
	return func(sl *SLog) {
		sl.option.setDefault = true
	}
}

// WithSLogSetDefaultCtx 设置默认上下文
func WithSLogSetDefaultCtx() SLogOption {
	return func(sl *SLog) {
		sl.option.setDefaultCtx = true
	}
}

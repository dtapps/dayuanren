package golog

import (
	"context"
)

// ApiSLog 接口日志
type ApiSLog struct {
	slog struct {
		status bool  // 状态
		client *SLog // 日志服务
	}
}

// ApiSLogFun 接口日志驱动
type ApiSLogFun func() *ApiSLog

func NewApiSlog(ctx context.Context) *ApiSLog {
	sl := &ApiSLog{}
	return sl
}

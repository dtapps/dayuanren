package gorequest

import (
	"context"
	"fmt"
)

var (
	xRequestID = "X-Request-ID"
	tNil       = "%!s(<nil>)"
)

func GetRequestIDContext(ctx context.Context) string {
	return customGetIDContext(ctx, xRequestID)
}

// customGetIDContext 通过自定义上下文获取跟踪编号
func customGetIDContext(ctx context.Context, key string) string {
	traceId := fmt.Sprintf("%s", ctx.Value(key))
	if traceId == tNil {
		return ""
	}
	if len(traceId) <= 0 {
		return ""
	}
	return traceId
}

package gorequest

import (
	"context"
	"fmt"
)

var (
	tID1 = "trace_id"
	tID2 = "X-Request-ID"
	tNil = "%!s(<nil>)"
)

func getTraceIDContext(ctx context.Context) string {
	return customGetIDContext(ctx, tID1)
}

func getRequestIDContext(ctx context.Context) string {
	return customGetIDContext(ctx, tID2)
}

func getIDContext(ctx context.Context) string {
	id := getTraceIDContext(ctx)
	if id == "" {
		id = getRequestIDContext(ctx)
	}
	return id
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

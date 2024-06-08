package golog

import (
	"context"
	"go.dtapp.net/gotrace_id"
	"log/slog"
)

type ContextHandler struct {
	slog.Handler
}

// Handle 添加上下文属性到 Record 中，然后调用底层的 handler
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	traceIdKey := gotrace_id.CustomGetTraceIdContext(ctx, gotrace_id.TraceIdKey)
	if traceIdKey != "" {
		r.AddAttrs(slog.String(gotrace_id.TraceIdKey, traceIdKey))
	} else {
		traceIdRequestKey := gotrace_id.CustomGetTraceIdContext(ctx, gotrace_id.TraceIdRequestKey)
		if traceIdRequestKey != "" {
			r.AddAttrs(slog.String(gotrace_id.TraceIdRequestKey, traceIdRequestKey))
		} else {
			traceIDRequestKey := gotrace_id.CustomGetTraceIdContext(ctx, gotrace_id.TraceIDRequestKey)
			if traceIDRequestKey != "" {
				r.AddAttrs(slog.String(gotrace_id.TraceIDRequestKey, traceIDRequestKey))
			}
		}
	}
	return h.Handler.Handle(ctx, r)
}

package gorequest

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SetTrace 设置OpenTelemetry链路追踪
func (c *App) SetTrace(trace bool) {
	c.trace = trace
}

// TraceStartSpan 开始OpenTelemetry链路追踪状态
func (c *App) TraceStartSpan(ctx context.Context, spanName string) context.Context {
	if c.trace {
		// 获取 Tracer
		tr := otel.Tracer("go.dtapp.net/gorequest", trace.WithInstrumentationVersion(Version))
		// 创建一个新的 span
		ctx, c.span = tr.Start(ctx, "gorequest."+spanName, trace.WithSpanKind(trace.SpanKindClient))
	}
	return ctx
}

// TraceEndSpan 结束OpenTelemetry链路追踪状态
func (c *App) TraceEndSpan() {
	if c.trace && c.span != nil {
		c.span.End()
	}
}

// TraceSetAttributes 设置OpenTelemetry链路追踪属性
func (c *App) TraceSetAttributes(kv ...attribute.KeyValue) {
	if c.trace && c.span != nil {
		c.span.SetAttributes(kv...)
	}
}

// TraceSetStatus 设置OpenTelemetry链路追踪状态
func (c *App) TraceSetStatus(code codes.Code, description string) {
	if c.trace && c.span != nil {
		c.span.SetStatus(code, description)
	}
}

// TraceRecordError 记录OpenTelemetry链路追踪错误
func (c *App) TraceRecordError(err error, options ...trace.EventOption) {
	if c.trace && c.span != nil {
		c.span.RecordError(err, options...)
	}
}

// TraceGetTraceID 获取OpenTelemetry链路追踪TraceID
func (c *App) TraceGetTraceID() (traceID string) {
	if c.trace && c.span != nil {
		traceID = c.span.SpanContext().TraceID().String()
	}
	return traceID
}

// TraceGetSpanID 获取OpenTelemetry链路追踪SpanID
func (c *App) TraceGetSpanID() (spanID string) {
	if c.trace && c.span != nil {
		spanID = c.span.SpanContext().SpanID().String()
	}
	return spanID
}

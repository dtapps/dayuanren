package dayuanren

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// SetTrace 设置OpenTelemetry链路追踪
func (c *Client) SetTrace(trace bool) {
	c.trace = trace
	c.httpClient.SetTrace(trace)
}

func (c *Client) TraceStartSpan(ctx context.Context, spanName string) context.Context {
	if c.trace {
		tr := otel.Tracer("go.dtapp.net/dayuanren", trace.WithInstrumentationVersion(Version))
		ctx, c.span = tr.Start(ctx, spanName)
	}
	return ctx
}

func (c *Client) TraceEndSpan() {
	if c.trace {
		c.span.End()
	}
}

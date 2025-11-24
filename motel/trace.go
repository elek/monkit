package motel

import (
	"context"

	"github.com/spacemonkeygo/monkit/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type TraceProvider struct {
	noop.TracerProvider
	registry *monkit.Registry
}

func NewTraceProvider(registry *monkit.Registry) *TraceProvider {
	return &TraceProvider{
		registry: registry,
	}
}

func (t *TraceProvider) Tracer(name string, options ...trace.TracerOption) trace.Tracer {
	return &MonkitTracer{
		scope:    t.registry.ScopeNamed(name),
		provider: t,
	}
}

var _ trace.TracerProvider = &TraceProvider{}

type MonkitTracer struct {
	noop.Tracer
	scope    *monkit.Scope
	provider *TraceProvider
}

func (m *MonkitTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	f := m.scope.FuncNamed(spanName)
	task := f.Task(&ctx)
	return ctx, &MonkitSpan{
		span:     monkit.SpanFromCtx(ctx),
		task:     task,
		provider: m.provider,
	}
}

var _ trace.Tracer = &MonkitTracer{}

type MonkitSpan struct {
	noop.Span
	span     *monkit.Span
	task     func(*error)
	provider *TraceProvider
}

func (m *MonkitSpan) End(options ...trace.SpanEndOption) {
	m.task(nil)
}

func (m *MonkitSpan) AddEvent(name string, options ...trace.EventOption) {
	//TODO implement me
	panic("implement me")
}

func (m *MonkitSpan) AddLink(link trace.Link) {
	//TODO implement me
	panic("implement me")
}

func (m *MonkitSpan) IsRecording() bool {
	//TODO implement me
	panic("implement me")
}

func (m *MonkitSpan) RecordError(err error, options ...trace.EventOption) {
	//TODO implement me
	panic("implement me")
}

func (m *MonkitSpan) SpanContext() trace.SpanContext {
	//TODO implement me
	panic("implement me")
}

func (m *MonkitSpan) SetStatus(code codes.Code, description string) {
	//TODO implement me
	panic("implement me")
}

func (m *MonkitSpan) SetName(name string) {
	//TODO implement me
	panic("implement me")
}

func (m *MonkitSpan) SetAttributes(kv ...attribute.KeyValue) {
	//TODO implement me
	panic("implement me")
}

func (m *MonkitSpan) TracerProvider() trace.TracerProvider {
	return m.provider
}

var _ trace.Span = &MonkitSpan{}

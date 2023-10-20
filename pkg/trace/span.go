package trace

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	End(options ...trace.SpanEndOption)
	AddEvent(name string, options ...trace.EventOption)
	IsRecording() bool
	RecordError(err error, options ...trace.EventOption)
	SpanContext() trace.SpanContext
	SetStatus(code codes.Code, description string)
	SetName(name string)
	SetAttributes(kv ...attribute.KeyValue)
	TracerProvider() trace.TracerProvider
}

type noopSpan struct {
}

func (noopSpan) SpanContext() trace.SpanContext { return trace.SpanContext{} }

func (noopSpan) IsRecording() bool { return false }

func (noopSpan) SetStatus(codes.Code, string) {}

func (noopSpan) SetError(bool) {}

func (noopSpan) SetAttributes(...attribute.KeyValue) {}

func (noopSpan) End(...trace.SpanEndOption) {}

func (noopSpan) RecordError(error, ...trace.EventOption) {}

func (noopSpan) AddEvent(string, ...trace.EventOption) {}

func (noopSpan) SetName(string) {}

func (noopSpan) TracerProvider() trace.TracerProvider { return noopTracerProvider{} }

type noopTracerProvider struct{}

// Tracer returns noop implementation of Tracer.
func (p noopTracerProvider) Tracer(string, ...trace.TracerOption) trace.Tracer {
	return noopTracer{}
}

type noopTracer struct{}

func (t noopTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	span := trace.SpanFromContext(ctx)
	if _, ok := span.(nonRecordingSpan); !ok {
		// span is likely already a noopSpan, but let's be sure
		span = noopSpan{}
	}
	return trace.ContextWithSpan(ctx, span), span
}

type nonRecordingSpan struct {
	noopSpan

	sc trace.SpanContext
}

// SpanContext returns the wrapped SpanContext.
func (s nonRecordingSpan) SpanContext() trace.SpanContext { return s.sc }

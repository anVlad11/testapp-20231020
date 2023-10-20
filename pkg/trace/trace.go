package trace

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultServiceName = "trace-engine"
)

var (
	isEnable = false
)

func SetTracerProvider(enable bool, url string, serviceName string) error {
	isEnable = enable
	if !enable {
		return nil
	}

	if serviceName == "" {
		serviceName = defaultServiceName
	}

	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tp)

	return nil
}

func CreateSpan(ctx context.Context, operationName string) (Span, context.Context) {
	if !isEnable {
		return noopSpan{}, ctx
	}

	if ctx == nil {
		return noopSpan{}, ctx
	}

	tr := otel.Tracer(operationName)
	ctx, span := tr.Start(ctx, operationName)
	span.SetAttributes(attribute.String("name", operationName))

	// Get caller function name, file and line
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	callerDetails := fmt.Sprintf("%s - %s#%d", frame.Function, frame.File, frame.Line)
	span.SetAttributes(attribute.String("caller", callerDetails))

	return span, ctx
}

func InjectHTTPCarrier(ctx context.Context, req *http.Request) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
}

func SetTraceStatus(span trace.Span, err error) {
	if span != nil && err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
}

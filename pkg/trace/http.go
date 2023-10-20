package trace

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerKey = "otel-tracer-echo"
)

type middlewareConfig struct {
	TracerProvider oteltrace.TracerProvider
	Skipper        middleware.Skipper
}

type Option func(config *middlewareConfig)

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return func(cfg *middlewareConfig) {
		if provider != nil {
			cfg.TracerProvider = provider
		}
	}
}

// WithSkipper specifies a skipper for allowing requests to skip generating spans.
func WithSkipper(skipper middleware.Skipper) Option {
	return func(cfg *middlewareConfig) {
		cfg.Skipper = skipper
	}
}

// MakeTracingMiddleware returns echo middleware which will trace incoming requests.
func MakeTracingMiddleware(service string, tracerName string, version string, options ...Option) echo.MiddlewareFunc {
	cfg := middlewareConfig{}
	for _, option := range options {
		option(&cfg)
	}

	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}

	tracer := cfg.TracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(version),
	)

	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}

			c.Set(tracerKey, tracer)
			request := c.Request()
			savedCtx := request.Context()
			defer func() {
				request = request.WithContext(savedCtx)
				c.SetRequest(request)
			}()
			ctx := otel.GetTextMapPropagator().Extract(savedCtx, propagation.HeaderCarrier(request.Header))
			opts := []oteltrace.SpanStartOption{
				oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
				oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
				oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(service, c.Path(), request)...),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}
			spanName := c.Path()
			if spanName == "" {
				spanName = fmt.Sprintf("HTTP %s route not found", request.Method)
			}

			ctx, span := tracer.Start(ctx, spanName, opts...)
			defer span.End()

			// pass the span through the request context
			c.SetRequest(request.WithContext(ctx))

			// serve the request to the next middleware
			err := next(c)
			if err != nil {
				span.SetAttributes(attribute.String("echo.error", err.Error()))
				// invokes the registered HTTP error handler
				c.Error(err)
			}

			attrs := semconv.HTTPAttributesFromHTTPStatusCode(c.Response().Status)
			spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(c.Response().Status, oteltrace.SpanKindServer)
			span.SetAttributes(attrs...)
			span.SetStatus(spanStatus, spanMessage)

			return nil
		}
	}
}

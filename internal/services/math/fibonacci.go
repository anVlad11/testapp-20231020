package math

import (
	"context"
	"github.com/anVlad11/testapp-20231020/pkg/trace"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *Service) Fibonacci(ctx context.Context, n int64) (decimal.Decimal, error) {
	var err error

	var span trace.Span
	span, ctx = trace.CreateSpan(ctx, "services.math.Fibonacci")
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	if n <= 0 {
		return decimal.Zero, nil
	}

	if n == 1 {
		return decimal.NewFromInt(1), nil
	}

	if result, exists := s.fibonacciCache.Load(n); exists {
		span.SetAttributes(attribute.Bool("cache_hit", true))
		return result, nil
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))

	f1 := decimal.NewFromInt(1)
	f2 := decimal.Zero
	for i := int64(0); i < n; i++ {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				s.logger.Info("context timeout")
			case context.Canceled:
				s.logger.Info("context cancelled")
			}
			span.SetAttributes(attribute.Bool("cancel", true))

			return decimal.Zero, ctx.Err()
		default:

		}

		f1, f2 = f2, f1.Add(f2)
		s.fibonacciCache.Store(i, f1)
	}

	s.fibonacciCache.Store(n, f2)

	return f2, nil
}

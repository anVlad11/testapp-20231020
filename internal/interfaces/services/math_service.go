package services

import (
	"context"
	"github.com/shopspring/decimal"
)

type MathService interface {
	Fibonacci(ctx context.Context, n int64) (decimal.Decimal, error)
}

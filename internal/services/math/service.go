package math

import (
	"github.com/anVlad11/testapp-20231020/pkg/log"
	"github.com/anVlad11/testapp-20231020/pkg/xsync"
	"github.com/shopspring/decimal"
)

type Service struct {
	logger         log.Logger
	fibonacciCache *xsync.Map[int64, decimal.Decimal]
}

func NewService(logger log.Logger) *Service {
	return &Service{
		logger:         logger,
		fibonacciCache: xsync.NewMap[int64, decimal.Decimal](),
	}
}

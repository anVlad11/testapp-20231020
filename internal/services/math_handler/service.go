package math_handler

import (
	"github.com/anVlad11/testapp-20231020/internal/interfaces/services"
	"github.com/anVlad11/testapp-20231020/pkg/log"
)

type Service struct {
	logger      log.Logger
	mathService services.MathService
}

func NewService(
	logger log.Logger,
	mathService services.MathService,
) *Service {
	return &Service{
		logger:      logger,
		mathService: mathService,
	}
}

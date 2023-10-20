package math_handler

import (
	"github.com/anVlad11/testapp-20231020/pkg/errors"
	"github.com/anVlad11/testapp-20231020/pkg/log"
	contract "github.com/anVlad11/testapp-20231020/pkg/testapp"
	"github.com/anVlad11/testapp-20231020/pkg/trace"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel/codes"
	"net/http"
)

func (s *Service) GetFibonacciNumber(c echo.Context) error {
	var err error

	span, ctx := trace.CreateSpan(
		c.Request().Context(),
		"services.math_handler.GetFibonacciNumber",
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	response := contract.V1FibonacciPostResponseBody{}

	var request *contract.V1FibonacciPostRequestBody

	err = c.Bind(&request)
	if err != nil {
		s.logger.Error(
			"could not bind request",
			log.Error(err),
		)

		response.Errors = append(response.Errors, errors.RequestBodyIsMalformed.Error())

		return c.JSON(http.StatusBadRequest, response)
	}

	var number decimal.Decimal
	number, err = s.mathService.Fibonacci(ctx, request.Position)
	if err != nil {
		s.logger.Error(
			"could not calculate fibonacci number",
			log.Int64("position", request.Position),
			log.Error(err),
		)

		response.Errors = append(response.Errors, errors.InternalError.Error())

		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Status = true
	response.Data = contract.V1FibonacciPostResponseBodyData{Result: number.String()}

	return c.JSON(http.StatusOK, response)
}

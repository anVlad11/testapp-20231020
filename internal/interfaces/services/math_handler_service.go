package services

import "github.com/labstack/echo/v4"

type MathHandlerService interface {
	GetFibonacciNumber(c echo.Context) error
}

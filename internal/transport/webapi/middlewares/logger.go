package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/vitaliy-ukiru/bank-service/pkg/logging"
)

func WrapRequestContextWithLogger(log logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestLogger := log
			reqId := c.Response().Header().Get(echo.HeaderXRequestID)
			if reqId != "" {
				requestLogger = requestLogger.With(logging.String("request_id", reqId))
			}
			request := c.Request()
			ctx := logging.Context(request.Context(), requestLogger)

			c.SetRequest(request.WithContext(ctx))
			return next(c)
		}
	}
}

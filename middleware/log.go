package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RequestMetadataMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			requestID := req.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
				req.Header.Set("X-Request-ID", requestID) // Optionally set it in the request header
			}
			contextualLogger := logger.With(zap.String("request_id", requestID), zap.String("function", c.Path()))
			c.Set("logger", contextualLogger)
			c.Set("request_id", requestID)
			return next(c)
		}
	}
}

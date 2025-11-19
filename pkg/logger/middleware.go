package logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		Log.Info("request",
			zap.String("method", c.Request().Method),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().Status),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.RealIP()),
		)

		return err
	}
}

package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func Log(logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		status := c.Response().StatusCode()
		fields := []zap.Field{
			zap.String("request_id", GetRequestID(c)),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.IP()),
		}

		switch {
		case err != nil && status >= 500:
			fields = append(fields, zap.String("error", err.Error()))
			logger.Error("request", fields...)
		case err != nil:
			fields = append(fields, zap.String("error", err.Error()))
			logger.Warn("request", fields...)
		default:
			logger.Info("request", fields...)
		}

		return err
	}
}

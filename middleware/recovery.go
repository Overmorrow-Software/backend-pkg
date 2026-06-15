package middleware

import (
	"github.com/Overmorrow-Software/backend-pkg/apierror"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic", zap.Any("recover", r), zap.Stack("stack"))
				err = apierror.Internal("internal server error")
			}
		}()
		return c.Next()
	}
}

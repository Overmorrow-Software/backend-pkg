package middleware

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type RequestIDKey string

const requestIDKey RequestIDKey = "request_id"

const HeaderRequestID = "X-Request-ID"

type ctxKey int

const reqIDCtxKey ctxKey = iota

func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Get(HeaderRequestID)
		if id == "" {
			id = uuid.New().String()
		}
		c.Locals(requestIDKey, id)
		c.Set(HeaderRequestID, id)
		c.SetContext(context.WithValue(c.Context(), reqIDCtxKey, id))
		return c.Next()
	}
}

func GetRequestID(c fiber.Ctx) string {
	id, ok := c.Locals(requestIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

func GetRequestIDCtx(ctx context.Context) string {
	id, ok := ctx.Value(reqIDCtxKey).(string)
	if !ok {
		return ""
	}
	return id
}

package req

import (
	"strconv"
	"time"

	"backend-pkg/pkg/apierror"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func Parse[T any](c fiber.Ctx) (T, error) {
	var r T
	if err := c.Bind().Body(&r); err != nil {
		return r, apierror.BadRequest("invalid body")
	}
	if err := validate.Struct(&r); err != nil {
		return r, apierror.BadRequest("invalid req content", err)
	}
	return r, nil
}

func ParseID(c fiber.Ctx, param string) (uint64, error) {
	id, err := strconv.ParseUint(c.Params(param), 10, 64)
	if err != nil || id == 0 {
		return 0, apierror.BadRequest("invalid " + param)
	}
	return id, nil
}

func ParseOptionalIDQuery(c fiber.Ctx, key string) (*uint64, error) {
	raw := c.Query(key)
	if raw == "" {
		return nil, nil
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return nil, apierror.BadRequest("invalid " + key)
	}
	return &id, nil
}

func ParseTimeQuery(c fiber.Ctx, key string) (time.Time, error) {
	raw := c.Query(key)
	if raw == "" {
		return time.Time{}, apierror.BadRequest("missing query param: " + key)
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, apierror.BadRequest("invalid time format for " + key + ", expected RFC3339")
	}
	return t, nil
}

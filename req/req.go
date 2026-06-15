package req

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Overmorrow-Software/backend-pkg/apierror"
	"github.com/Overmorrow-Software/backend-pkg/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Validator struct {
	v *validator.Validate
}

type ValidatorOption func(*validator.Validate)

func WithValidation(tag string, fn validator.Func) ValidatorOption {
	return func(v *validator.Validate) {
		_ = v.RegisterValidation(tag, fn)
	}
}

func NewValidator(opts ...ValidatorOption) *Validator {
	v := validator.New()
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "" || name == "-" {
			return f.Name
		}
		return name
	})
	for _, opt := range opts {
		opt(v)
	}
	return &Validator{v: v}
}

var std = NewValidator()

func Parse[T any](c fiber.Ctx) (T, error) {
	return ParseWith[T](c, std)
}

func ParseWith[T any](c fiber.Ctx, v *Validator) (T, error) {
	var r T
	if err := c.Bind().Body(&r); err != nil {
		return r, apierror.BadRequest("invalid body")
	}
	if err := v.v.Struct(&r); err != nil {
		errs, _ := err.(validator.ValidationErrors)
		return r, apierror.Validation(toFieldErrors(errs))
	}
	return r, nil
}

func toFieldErrors(errs validator.ValidationErrors) []apierror.FieldError {
	out := make([]apierror.FieldError, 0, len(errs))
	for _, fe := range errs {
		out = append(out, apierror.FieldError{
			Field:   fe.Field(),
			Message: fieldMessage(fe),
		})
	}
	return out
}

func fieldMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required"
	case "email":
		return "must be a valid email"
	case "min":
		return "min length is " + fe.Param()
	case "max":
		return "max length is " + fe.Param()
	case "gte":
		return "must be >= " + fe.Param()
	case "lte":
		return "must be <= " + fe.Param()
	case "gt":
		return "must be > " + fe.Param()
	case "lt":
		return "must be < " + fe.Param()
	case "len":
		return "length must be " + fe.Param()
	case "oneof":
		return "must be one of: " + fe.Param()
	case "uuid", "uuid4":
		return "must be a valid UUID"
	case "url", "uri":
		return "must be a valid URL"
	default:
		return "invalid value"
	}
}

func ParseID(c fiber.Ctx, param string) (uint64, error) {
	id, err := strconv.ParseUint(c.Params(param), 10, 64)
	if err != nil || id == 0 {
		return 0, apierror.BadRequest("invalid " + param)
	}
	return id, nil
}

func ParseUUID(c fiber.Ctx, param string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params(param))
	if err != nil {
		return uuid.Nil, apierror.BadRequest("invalid " + param)
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

type queryOptions struct {
	PageNum   uint64 `query:"page_num"`
	PageSize  uint64 `query:"page_size"`
	OrderBy   string `query:"order_by"`
	OrderType string `query:"order_type"`
}

func ParseOptions(c fiber.Ctx) (*repository.Options, error) {
	var q queryOptions
	if err := c.Bind().Query(&q); err != nil {
		return nil, apierror.BadRequest("invalid query params")
	}
	return &repository.Options{
		Pagination: repository.Pagination{
			PageNum:  q.PageNum,
			PageSize: q.PageSize,
		},
		Order: repository.Order{
			OrderBy:   q.OrderBy,
			OrderType: q.OrderType,
		},
	}, nil
}

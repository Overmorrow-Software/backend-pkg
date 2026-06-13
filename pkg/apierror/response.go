package apierror

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

type Response[T any] struct {
	Data  T            `json:"data"`
	Error *ErrorPublic `json:"error"`
}

func OK[T any](c fiber.Ctx, data T) error {
	return c.JSON(Response[T]{Data: data})
}

func Empty(c fiber.Ctx) error {
	return c.JSON(Response[any]{})
}

type Paginated[T any] struct {
	Items    []T    `json:"items"`
	Total    int    `json:"total"`
	PageNum  uint64 `json:"page_num"`
	PageSize uint64 `json:"page_size"`
}

func Page[T any](c fiber.Ctx, items []T, total int, pageNum, pageSize uint64) error {
	if items == nil {
		items = make([]T, 0)
	}
	return c.JSON(Response[Paginated[T]]{
		Data: Paginated[T]{
			Items:    items,
			Total:    total,
			PageNum:  pageNum,
			PageSize: pageSize,
		},
	})
}

func ErrorHandler(c fiber.Ctx, err error) error {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		pub := apiErr.Public()
		return c.Status(apiErr.Status).JSON(Response[any]{Error: &pub})
	}

	pub := ErrorPublic{
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
		Status:  fiber.StatusInternalServerError,
	}
	return c.Status(fiber.StatusInternalServerError).JSON(Response[any]{Error: &pub})
}

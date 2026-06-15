package response

import (
	"errors"

	"github.com/Overmorrow-Software/backend-pkg/apierror"
	"github.com/Overmorrow-Software/backend-pkg/repository"
	"github.com/gofiber/fiber/v3"
)

var defaultErr = apierror.Error{
	Code:    "INTERNAL_ERROR",
	Message: "internal server error",
	Status:  fiber.StatusInternalServerError,
}

type Response[T any] struct {
	Data  T               `json:"data"`
	Error *apierror.Error `json:"error"`
}

type Paginated[T any] struct {
	Items    []T    `json:"items"`
	Total    int    `json:"total"`
	PageNum  uint64 `json:"page_num"`
	PageSize uint64 `json:"page_size"`
}

func OK[T any](c fiber.Ctx, data T) error {
	return c.JSON(Response[T]{Data: data})
}

func Empty(c fiber.Ctx) error {
	return c.JSON(Response[any]{})
}

func Page[T any](c fiber.Ctx, result repository.PageResult[T], pagination repository.Pagination) error {
	items := result.Items
	if items == nil {
		items = make([]T, 0)
	}
	return c.JSON(Response[Paginated[T]]{
		Data: Paginated[T]{
			Items:    items,
			Total:    result.Total,
			PageNum:  pagination.PageNum,
			PageSize: pagination.PageSize,
		},
	})
}

func ErrorHandler(c fiber.Ctx, err error) error {
	var apiErr *apierror.Error
	if errors.As(err, &apiErr) {
		return c.Status(apiErr.Status).JSON(Response[any]{Error: apiErr})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(Response[any]{Error: &defaultErr})
}

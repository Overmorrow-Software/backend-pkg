package apierror

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/goccy/go-json"
)

type Error struct {
	BaseError error
	Code      string
	Message   string
	path      string
	Status    int
}

func newError(status int, message string, err ...error) *Error {
	var baseError error
	if len(err) > 0 {
		baseError = err[0]
	}

	var e *Error
	if errors.As(baseError, &e) {
		return e.wrap(message, status)
	}

	var path string
	if _, file, line, ok := runtime.Caller(2); ok {
		path = fmt.Sprintf("%s:%d", file, line)
	}

	return &Error{
		Code:      defaultCode(status),
		BaseError: baseError,
		Message:   message,
		Status:    status,
		path:      path,
	}
}

func defaultCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusPaymentRequired:
		return "PAYMENT_REQUIRED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "ALREADY_EXISTS"
	case http.StatusGone:
		return "GONE"
	case http.StatusTooManyRequests:
		return "TOO_MANY_REQUESTS"
	case StatusLoginTimeout:
		return "LOGIN_TIMEOUT"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	default:
		return "UNKNOWN"
	}
}

type ErrorPublic struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *Error) Public() ErrorPublic {
	return ErrorPublic{
		Code:    e.Code,
		Message: e.Message,
		Status:  e.Status,
	}
}

func (e *Error) WithCode(code string) *Error {
	e.Code = code
	return e
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	b, err := json.Marshal(e.Public())
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (e *Error) wrap(message string, httpStatus int) *Error {
	if e == nil {
		return newError(httpStatus, message)
	}
	return &Error{
		Code:      defaultCode(httpStatus),
		BaseError: e.BaseError,
		Message:   fmt.Sprintf("%s: %s", message, e.Message),
		Status:    httpStatus,
		path:      e.path,
	}
}

func (e *Error) Path() string {
	return e.path
}

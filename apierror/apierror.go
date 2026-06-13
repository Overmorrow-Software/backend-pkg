package apierror

import (
	"net/http"

	"github.com/goccy/go-json"
)

type Error struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Status  int          `json:"status"`
	Fields  []FieldError `json:"fields,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func newError(status int, message string) *Error {
	return &Error{
		Code:    defaultCode(status),
		Message: message,
		Status:  status,
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

func (e *Error) WithCode(code string) *Error {
	e.Code = code
	return e
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	b, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

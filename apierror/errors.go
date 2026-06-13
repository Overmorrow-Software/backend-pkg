package apierror

import "net/http"

const StatusLoginTimeout = 440

func Validation(fields []FieldError) *Error {
	return &Error{
		Code:    "VALIDATION_ERROR",
		Message: "validation failed",
		Status:  http.StatusBadRequest,
		Fields:  fields,
	}
}

func BadRequest(message string) *Error      { return newError(http.StatusBadRequest, message) }
func Unauthorized(message string) *Error    { return newError(http.StatusUnauthorized, message) }
func PaymentRequired(message string) *Error { return newError(http.StatusPaymentRequired, message) }
func Forbidden(message string) *Error       { return newError(http.StatusForbidden, message) }
func NotFound(message string) *Error        { return newError(http.StatusNotFound, message) }
func Gone(message string) *Error            { return newError(http.StatusGone, message) }
func AlreadyExist(message string) *Error    { return newError(http.StatusConflict, message) }
func TooManyRequests(message string) *Error { return newError(http.StatusTooManyRequests, message) }
func LoginTimeout(message string) *Error    { return newError(StatusLoginTimeout, message) }
func Internal(message string) *Error        { return newError(http.StatusInternalServerError, message) }

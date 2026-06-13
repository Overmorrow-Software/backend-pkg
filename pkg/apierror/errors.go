package apierror

import "net/http"

const (
	StatusLoginTimeout = 440
)

func BadRequest(message string, err ...error) *Error {
	return newError(http.StatusBadRequest, message, err...)
}

func Unauthorized(message string, err ...error) *Error {
	return newError(http.StatusUnauthorized, message, err...)
}

func PaymentRequired(message string, err ...error) *Error {
	return newError(http.StatusPaymentRequired, message, err...)
}

func Forbidden(message string, err ...error) *Error {
	return newError(http.StatusForbidden, message, err...)
}

func NotFound(message string, err ...error) *Error {
	return newError(http.StatusNotFound, message, err...)
}

func Gone(message string, err ...error) *Error {
	return newError(http.StatusGone, message, err...)
}

func AlreadyExist(message string, err ...error) *Error {
	return newError(http.StatusConflict, message, err...)
}

func TooManyRequests(message string, err ...error) *Error {
	return newError(http.StatusTooManyRequests, message, err...)
}

func LoginTimeout(message string, err ...error) *Error {
	return newError(StatusLoginTimeout, message, err...)
}

func Internal(message string, err ...error) *Error {
	return newError(http.StatusInternalServerError, message, err...)
}

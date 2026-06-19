package apierror

import "net/http"

const StatusLoginTimeout = 440

func Validation(reqID string, fields []FieldError) *Error {
	return &Error{
		Code:    "VALIDATION_ERROR",
		Message: "validation failed",
		Status:  http.StatusBadRequest,
		ReqID:   reqID,
		Fields:  fields,
	}
}

func BadRequest(message, reqID string) *Error      { return newError(http.StatusBadRequest, message, reqID) }
func Unauthorized(message, reqID string) *Error    { return newError(http.StatusUnauthorized, message, reqID) }
func PaymentRequired(message, reqID string) *Error { return newError(http.StatusPaymentRequired, message, reqID) }
func Forbidden(message, reqID string) *Error       { return newError(http.StatusForbidden, message, reqID) }
func NotFound(message, reqID string) *Error        { return newError(http.StatusNotFound, message, reqID) }
func Gone(message, reqID string) *Error            { return newError(http.StatusGone, message, reqID) }
func AlreadyExist(message, reqID string) *Error    { return newError(http.StatusConflict, message, reqID) }
func TooManyRequests(message, reqID string) *Error { return newError(http.StatusTooManyRequests, message, reqID) }
func LoginTimeout(message, reqID string) *Error    { return newError(StatusLoginTimeout, message, reqID) }
func Internal(message, reqID string) *Error        { return newError(http.StatusInternalServerError, message, reqID) }

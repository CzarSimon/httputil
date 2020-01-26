package httputil

import (
	"fmt"
	"net/http"

	"github.com/CzarSimon/httputil/id"
)

// Error error containing status code and error.
type Error struct {
	ID        string `json:"id,omitempty"`
	RequestID string `json:"requestId,omitempty"`
	Status    int    `json:"status,omitempty"`
	Message   string `json:"message,omitempty"`
	Err       error  `json:"-"`
}

// Error retruns a string representation of an Error and
// makes the type compliant with the go error interface.
func (err *Error) Error() string {
	return fmt.Sprintf("%d - %s", err.Status, err.Err)
}

// Unwrap returns the enclosed error.
func (err *Error) Unwrap() error {
	return err.Err
}

// BadRequestError creates a 400 - Bad Request error.
func BadRequestError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusBadRequest, err)
}

// UnauthorizedError creates a 401 - Unauthorized error.
func UnauthorizedError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusUnauthorized, err)
}

// ForbiddenError creates a 403 - Forbidden error.
func ForbiddenError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusForbidden, err)
}

// NotFoundError creates a 404 - Not Found error.
func NotFoundError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusNotFound, err)
}

// MethodNotAllowedError creates a 405 - Method Not Allowed error.
func MethodNotAllowedError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusMethodNotAllowed, err)
}

// ConflictError creates a 409 - Conflict error.
func ConflictError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusConflict, err)
}

// PreconditionRequiredError creates a 428 - Precondition Required error.
func PreconditionRequiredError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusPreconditionRequired, err)
}

// TooManyRequestsError creates a 429 - Too Many Requests error.
func TooManyRequestsError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusTooManyRequests, err)
}

// InternalServerError creates a 500 - Internal Server Error.
func InternalServerError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusInternalServerError, err)
}

// NotImplementedError creates a 501 - Not Implemented error.
func NotImplementedError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusNotImplemented, err)
}

// BadGatewayError creates a 502 - Bad Gateway error.
func BadGatewayError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusBadGateway, err)
}

// ServiceUnavailableError creates a 503 - Service Unavailable error.
func ServiceUnavailableError(requestID string, err error) *Error {
	return errorFromStatus(requestID, http.StatusServiceUnavailable, err)
}

func errorFromStatus(requestID string, status int, err error) *Error {
	return NewError(requestID, http.StatusText(status), status, err)
}

// NewError creates a new Error based on a supplied status code
// attempts to derive the error message.
func NewError(requestID, message string, status int, err error) *Error {
	return &Error{
		ID:        id.New(),
		RequestID: requestID,
		Status:    status,
		Message:   message,
		Err:       err,
	}
}

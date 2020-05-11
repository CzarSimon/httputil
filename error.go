package httputil

import (
	"fmt"
	"net/http"

	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracelog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
)

var errLog = logger.GetDefaultLogger("httputil/error-log")

// Error error containing status code and error.
type Error struct {
	ID      string `json:"id,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Err     error  `json:"-"`
}

// Error retruns a string representation of an Error and
// makes the type compliant with the go error interface.
func (err *Error) Error() string {
	return fmt.Sprintf("Error(id=%s, message=%s, status=%d, err=%v)", err.ID, err.Message, err.Status, err.Err)
}

// Unwrap returns the enclosed error.
func (err *Error) Unwrap() error {
	return err.Err
}

// HandleErrors wrapper function to deal with encountered errors
// during request handling.
func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := getFirstError(c)
		if err == nil {
			return
		}

		logError(c, err)
		c.AbortWithStatusJSON(err.Status, err)
	}
}

// getFirstError returns the first error in the gin.Context, nil if not present.
func getFirstError(c *gin.Context) *Error {
	allErrors := c.Errors
	if len(allErrors) == 0 {
		return nil
	}
	err := allErrors[0].Err

	var httpError *Error
	switch err.(type) {
	case *Error:
		httpError = err.(*Error)
		break
	default:
		httpError = InternalServerError(err)
		break
	}

	return httpError
}

func logError(c *gin.Context, err *Error) {
	span := opentracing.SpanFromContext(c.Request.Context())
	if span != nil {
		span.LogFields(tracelog.Error(err))
		ext.HTTPStatusCode.Set(span, uint16(err.Status))
	}

	if err.Status < 500 {
		errLog.Info(err.Message,
			zap.Int("status", err.Status),
			zap.String("errorId", err.ID),
			zap.Error(err.Err))
		return
	}
	errLog.Error(err.Message,
		zap.Int("status", err.Status),
		zap.String("errorId", err.ID),
		zap.Error(err.Err))
}

// BadRequestError creates a 400 - Bad Request error.
func BadRequestError(err error) *Error {
	return errorFromStatus(http.StatusBadRequest, err)
}

// UnauthorizedError creates a 401 - Unauthorized error.
func UnauthorizedError(err error) *Error {
	return errorFromStatus(http.StatusUnauthorized, err)
}

// ForbiddenError creates a 403 - Forbidden error.
func ForbiddenError(err error) *Error {
	return errorFromStatus(http.StatusForbidden, err)
}

// NotFoundError creates a 404 - Not Found error.
func NotFoundError(err error) *Error {
	return errorFromStatus(http.StatusNotFound, err)
}

// MethodNotAllowedError creates a 405 - Method Not Allowed error.
func MethodNotAllowedError(err error) *Error {
	return errorFromStatus(http.StatusMethodNotAllowed, err)
}

// ConflictError creates a 409 - Conflict error.
func ConflictError(err error) *Error {
	return errorFromStatus(http.StatusConflict, err)
}

// UnsupportedMediaTypeError creates a 415 - Unsupported Media Type error.
func UnsupportedMediaTypeError(err error) *Error {
	return errorFromStatus(http.StatusUnsupportedMediaType, err)
}

// PreconditionRequiredError creates a 428 - Precondition Required error.
func PreconditionRequiredError(err error) *Error {
	return errorFromStatus(http.StatusPreconditionRequired, err)
}

// TooManyRequestsError creates a 429 - Too Many Requests error.
func TooManyRequestsError(err error) *Error {
	return errorFromStatus(http.StatusTooManyRequests, err)
}

// InternalServerError creates a 500 - Internal Server Error.
func InternalServerError(err error) *Error {
	return errorFromStatus(http.StatusInternalServerError, err)
}

// NotImplementedError creates a 501 - Not Implemented error.
func NotImplementedError(err error) *Error {
	return errorFromStatus(http.StatusNotImplemented, err)
}

// BadGatewayError creates a 502 - Bad Gateway error.
func BadGatewayError(err error) *Error {
	return errorFromStatus(http.StatusBadGateway, err)
}

// ServiceUnavailableError creates a 503 - Service Unavailable error.
func ServiceUnavailableError(err error) *Error {
	return errorFromStatus(http.StatusServiceUnavailable, err)
}

func errorFromStatus(status int, err error) *Error {
	return NewError(http.StatusText(status), status, err)
}

// NewError creates a new Error based on a supplied status code
// attempts to derive the error message.
func NewError(message string, status int, err error) *Error {
	return &Error{
		ID:      id.New(),
		Status:  status,
		Message: message,
		Err:     err,
	}
}

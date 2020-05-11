package httputil_test

import (
	"errors"
	"testing"

	"github.com/CzarSimon/httputil"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	assert := assert.New(t)

	baseErr := errors.New("base error")

	err := httputil.BadRequestError(baseErr)
	assert.Equal(400, err.Status)
	assert.Equal("Bad Request", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.UnauthorizedError(baseErr)
	assert.Equal(401, err.Status)
	assert.Equal("Unauthorized", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.ForbiddenError(baseErr)
	assert.Equal(403, err.Status)
	assert.Equal("Forbidden", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.NotFoundError(baseErr)
	assert.Equal(404, err.Status)
	assert.Equal("Not Found", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.MethodNotAllowedError(baseErr)
	assert.Equal(405, err.Status)
	assert.Equal("Method Not Allowed", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.ConflictError(baseErr)
	assert.Equal(409, err.Status)
	assert.Equal("Conflict", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.UnsupportedMediaTypeError(baseErr)
	assert.Equal(415, err.Status)
	assert.Equal("Unsupported Media Type", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.PreconditionRequiredError(baseErr)
	assert.Equal(428, err.Status)
	assert.Equal("Precondition Required", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.TooManyRequestsError(baseErr)
	assert.Equal(429, err.Status)
	assert.Equal("Too Many Requests", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.InternalServerError(baseErr)
	assert.Equal(500, err.Status)
	assert.Equal("Internal Server Error", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.NotImplementedError(baseErr)
	assert.Equal(501, err.Status)
	assert.Equal("Not Implemented", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.BadGatewayError(baseErr)
	assert.Equal(502, err.Status)
	assert.Equal("Bad Gateway", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.ServiceUnavailableError(baseErr)
	assert.Equal(503, err.Status)
	assert.Equal("Service Unavailable", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))
}

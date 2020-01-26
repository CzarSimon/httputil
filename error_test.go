package httputil_test

import (
	"errors"
	"testing"

	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/id"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	assert := assert.New(t)

	requestID := id.New()
	baseErr := errors.New("base error")

	err := httputil.BadRequestError(requestID, baseErr)
	assert.Equal(400, err.Status)
	assert.Equal("Bad Request", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.UnauthorizedError(requestID, baseErr)
	assert.Equal(401, err.Status)
	assert.Equal("Unauthorized", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.ForbiddenError(requestID, baseErr)
	assert.Equal(403, err.Status)
	assert.Equal("Forbidden", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.NotFoundError(requestID, baseErr)
	assert.Equal(404, err.Status)
	assert.Equal("Not Found", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.MethodNotAllowedError(requestID, baseErr)
	assert.Equal(405, err.Status)
	assert.Equal("Method Not Allowed", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.ConflictError(requestID, baseErr)
	assert.Equal(409, err.Status)
	assert.Equal("Conflict", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.PreconditionRequiredError(requestID, baseErr)
	assert.Equal(428, err.Status)
	assert.Equal("Precondition Required", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.TooManyRequestsError(requestID, baseErr)
	assert.Equal(429, err.Status)
	assert.Equal("Too Many Requests", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.InternalServerError(requestID, baseErr)
	assert.Equal(500, err.Status)
	assert.Equal("Internal Server Error", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.NotImplementedError(requestID, baseErr)
	assert.Equal(501, err.Status)
	assert.Equal("Not Implemented", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.BadGatewayError(requestID, baseErr)
	assert.Equal(502, err.Status)
	assert.Equal("Bad Gateway", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))

	err = httputil.ServiceUnavailableError(requestID, baseErr)
	assert.Equal(503, err.Status)
	assert.Equal("Service Unavailable", err.Message)
	assert.Error(err)
	assert.True(errors.Is(err, baseErr))
}

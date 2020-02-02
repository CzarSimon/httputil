package httputil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/CzarSimon/httputil/jwt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracelog "github.com/opentracing/opentracing-go/log"
)

// RBAC adds role based access controll checks extracting roles from jwt.
type RBAC struct {
	Verifier jwt.Verifier
}

// Secure checks if a request was made with a jwt containing a specified list of roles.
func (r *RBAC) Secure(roles ...string) gin.HandlerFunc {
	validRoles := make([]string, 0)
	for _, role := range roles {
		if role != "" {
			validRoles = append(validRoles, role)
		}
	}

	return func(c *gin.Context) {
		user, err := extractUserFromRequest(c, r.Verifier)
		if err != nil {
			c.AbortWithStatusJSON(err.Status, err)
			return
		}

		span := opentracing.SpanFromContext(c.Request.Context())
		if span != nil {
			span.SetBaggageItem("user-id", user.ID)
		}

		for _, role := range validRoles {
			if user.HasRole(role) {
				c.Next()
				return
			}
		}

		msg := fmt.Sprintf("%s %s access denied for %s", c.Request.Method, c.Request.URL.Path, user)
		err = ForbiddenError(errors.New(msg))
		if span != nil {
			span.LogFields(tracelog.Error(err))
			ext.HTTPStatusCode.Set(span, uint16(err.Status))
		}

		c.AbortWithStatusJSON(err.Status, err)
	}
}

func extractUserFromRequest(c *gin.Context, verifier jwt.Verifier) (jwt.User, *Error) {
	token, err := exctractToken(c)
	if err != nil {
		return jwt.User{}, err
	}

	user, jwtErr := verifier.Verify(token)
	if jwtErr != nil {
		return jwt.User{}, UnauthorizedError(err)
	}

	return user, nil
}

func exctractToken(c *gin.Context) (string, *Error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		err := errors.New("no authorization header provided")
		return "", UnauthorizedError(err)
	}

	token := strings.Replace(header, "Bearer ", "", 1)
	return token, nil
}

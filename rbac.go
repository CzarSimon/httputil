package httputil

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/CzarSimon/httputil/jwt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

const (
	userKey = "X-JWT-User"
)

// GetPrincipal returns the authenticated user if exists.
func GetPrincipal(c *gin.Context) (jwt.User, bool) {
	val, ok := c.Get(userKey)
	if !ok {
		return jwt.User{}, false
	}

	user, ok := val.(jwt.User)
	return user, ok
}

// RBAC adds role based access controll checks extracting roles from jwt.
type RBAC struct {
	Verifier jwt.Verifier
}

// NewRBAC creates a new RBAC struct with sane defaults.
func NewRBAC(creds jwt.Credentials) RBAC {
	return RBAC{
		Verifier: jwt.NewVerifier(creds, time.Minute),
	}
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
			logError(c, err)
			c.AbortWithStatusJSON(err.Status, err)
			return
		}

		span := opentracing.SpanFromContext(c.Request.Context())
		if span != nil {
			span.SetBaggageItem("user-id", user.ID)
			span.SetBaggageItem("user-roles", strings.Join(user.Roles, ";"))
		}
		c.Set(userKey, user)

		for _, role := range validRoles {
			if user.HasRole(role) {
				c.Next()
				return
			}
		}

		msg := fmt.Sprintf("%s %s access denied for %s", c.Request.Method, c.Request.URL.Path, user)
		err = ForbiddenError(errors.New(msg))
		logError(c, err)
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

package httputil_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/stretchr/testify/assert"
)

func TestRBAC(t *testing.T) {
	assert := assert.New(t)
	r := httputil.NewRouter("httputil-test", func() error {
		return nil
	})
	rbac := httputil.RBAC{
		Verifier: jwt.NewVerifier(getTestJWTCredentials(), time.Minute),
	}
	secured := r.Group("", rbac.Secure(jwt.AnonymousRole, jwt.AdminRole))
	secured.GET("/test", httputil.SendOK)

	req := createTestRequest("/health", http.MethodGet, "", nil)
	res := performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/test", http.MethodGet, jwt.AnonymousRole, nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/test", http.MethodGet, jwt.AdminRole, nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/test", http.MethodGet, "", nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusUnauthorized, res.Code)

	req = createTestRequest("/test", http.MethodGet, "OTHER_ROLE", nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusForbidden, res.Code)
}

func TestRBAC_WithConstructor(t *testing.T) {
	assert := assert.New(t)
	r := httputil.NewRouter("httputil-test", func() error {
		return nil
	})
	rbac := httputil.NewRBAC(getTestJWTCredentials())
	secured := r.Group("", rbac.Secure(jwt.AnonymousRole, jwt.AdminRole))
	secured.GET("/test", httputil.SendOK)

	req := createTestRequest("/health", http.MethodGet, "", nil)
	res := performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/test", http.MethodGet, jwt.AnonymousRole, nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/test", http.MethodGet, jwt.AdminRole, nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/test", http.MethodGet, "", nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusUnauthorized, res.Code)

	req = createTestRequest("/test", http.MethodGet, "OTHER_ROLE", nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusForbidden, res.Code)
}

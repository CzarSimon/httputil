package httputil_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/client/rpc"
	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHealth(t *testing.T) {
	assert := assert.New(t)
	r := httputil.NewRouter("httputil-test", func() error {
		return nil
	})

	req := createTestRequest("/health", http.MethodGet, "", nil)
	res := performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	r = httputil.NewRouter("httputil-test-fail", func() error {
		return httputil.ServiceUnavailableError(nil)
	})

	req = createTestRequest("/health", http.MethodGet, "", nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusServiceUnavailable, res.Code)
}

func TestAllowContentType(t *testing.T) {
	assert := assert.New(t)
	r := httputil.NewRouter("httputil-test", func() error {
		return nil
	})
	r.Use(httputil.AllowContentType("application/json"))
	r.GET("/test", httputil.SendOK)

	req := createTestRequest("/test", http.MethodGet, "", nil)
	req.Header.Add("Content-Type", "application/json")
	res := performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/test", http.MethodGet, "", nil)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	res = performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/test", http.MethodGet, "", nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusUnsupportedMediaType, res.Code)

	req = createTestRequest("/test", http.MethodGet, "", nil)
	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	res = performTestRequest(r, req)
	assert.Equal(http.StatusUnsupportedMediaType, res.Code)

	req = createTestRequest("/test", http.MethodGet, "", nil)
	req.Header.Add("Content-Type", "application/xml")
	res = performTestRequest(r, req)
	assert.Equal(http.StatusUnsupportedMediaType, res.Code)
}

func TestAllowJSON(t *testing.T) {
	assert := assert.New(t)
	r := httputil.NewRouter("httputil-test", func() error {
		return nil
	})
	r.Use(httputil.AllowJSON())
	r.GET("/test", httputil.SendOK)

	req := createTestRequest("/test", http.MethodGet, "", nil)
	req.Header.Add("Content-Type", "application/json")
	res := performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/test", http.MethodGet, "", nil)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	res = performTestRequest(r, req)
	assert.Equal(http.StatusOK, res.Code)

	req = createTestRequest("/test", http.MethodGet, "", nil)
	res = performTestRequest(r, req)
	assert.Equal(http.StatusUnsupportedMediaType, res.Code)

	req = createTestRequest("/test", http.MethodGet, "", nil)
	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	res = performTestRequest(r, req)
	assert.Equal(http.StatusUnsupportedMediaType, res.Code)

	req = createTestRequest("/test", http.MethodGet, "", nil)
	req.Header.Add("Content-Type", "application/xml")
	res = performTestRequest(r, req)
	assert.Equal(http.StatusUnsupportedMediaType, res.Code)
}

// ---- Test utils ----

func performTestRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func createTestRequest(route, method, role string, body interface{}) *http.Request {
	client := rpc.NewClient(time.Second)
	req, err := client.CreateRequest(method, route, body)
	if err != nil {
		log.Fatal("Failed to create request", zap.Error(err))
	}

	if role == "" {
		return req
	}

	issuer := jwt.NewIssuer(getTestJWTCredentials())
	token, err := issuer.Issue(jwt.User{
		ID:    id.New(),
		Roles: []string{role},
	}, time.Hour)
	if err != nil {
		log.Fatal(err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+token)
	return req
}

func getTestJWTCredentials() jwt.Credentials {
	return jwt.Credentials{
		Issuer: "httputil_test",
		Secret: "very-secret-secret",
	}
}

package auth

import (
	"log"
	"net/http"

	"github.com/CzarSimon/httputil"
)

// Test Function for testing whether a request is authorized.
type Test func(r *http.Request) bool

// Wrapper Struct for encapsulating a http.Handler after an authentication test.
type Wrapper struct {
	auth Test
}

// NewWrapper Creates a new wrapper with the supplied authentication test.
func NewWrapper(auth Test) Wrapper {
	return Wrapper{
		auth: auth,
	}
}

// Wrap Converts a supplied handler to an AuthHandler
func (wr Wrapper) Wrap(h http.Handler) http.Handler {
	return AuthHandler{
		handler: h,
		auth:    wr.auth,
	}
}

// AuthHandler http.Handler compliant struct with an authentication test.
type AuthHandler struct {
	handler http.Handler
	auth    Test
}

// ServeHTTP Tests authentication before executing the nested handlers ServeHTTP.
func (auth AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	TestAuth(auth.auth, auth.handler.ServeHTTP, w, r)
}

// WrapFunc Wraps a authentication test around a hander function.
func (wr Wrapper) WrapFunc(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		TestAuth(wr.auth, f, w, r)
	}
}

// TestAuth Executes authentication test and handles result.
func TestAuth(test Test, f http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	if !test(r) {
		LogStatus("Auth failed", r)
		httputil.SendErr(w, httputil.NotAuthorized)
		return
	}
	LogStatus("Auth success", r)
	f(w, r)
}

// LogStatus Logs outcome of authorization challange.
func LogStatus(msg string, r *http.Request) {
	log.Printf("%s from: %s, %s\n", r.URL.Path, r.RemoteAddr, msg)
}

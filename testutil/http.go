package testutil

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

// PerformRequest perform a test request against a given handler, returing a recorder of the response.
func PerformRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// CreateRequest creates a http request for testing.
func CreateRequest(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalf("failed to create request %s", err.Error())
	}

	return req
}

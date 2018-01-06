package httputil

import (
	"net/http"
)

var HealthCheck = NewHandler(healthCheckFunc)

// HandlerFunc Function for dealing with a http request
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// Handler http.Handler compliant handler with build in error handling
type Handler struct {
	Handle HandlerFunc
}

// NewHandler Creates a new Handler struct
func NewHandler(fn HandlerFunc) Handler {
	return Handler{
		Handle: fn,
	}
}

// ServeHTTP Invocaion method for the Handle function.
// Makes Handler compliant with http.Handler
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.Handle(w, r)
	if err == nil {
		return
	}
	switch e := err.(type) {
	case Error:
		SendErr(w, e)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// healthCheckFunc Sends an ok status to the requestor to confim health.
func healthCheckFunc(w http.ResponseWriter, r *http.Request) error {
	SendOK(w)
	return nil
}

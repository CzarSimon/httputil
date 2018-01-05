package httputil

import (
	"fmt"
	"net/http"
)

// HttpError Error containing status code and error
type HttpError struct {
	Status int   `json:"status"`
	Err    error `json:"error"`
}

// Error Retruns a string representation of an HttpError and
// makes the type compliant with the go error interface
func (err HttpError) Error() string {
	return fmt.Sprintf("%d - %s", err.Status, err.Err)
}

// NewError Creates a new HttpError based on a supplied status code
// attempts to derive the error message
func NewError(status int) HttpError {
	return HttpError{
		Status: status,
		Err:    fmt.Errorf(http.StatusText(status)),
	}
}

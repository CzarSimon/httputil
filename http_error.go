package httputil

import (
	"fmt"
	"net/http"
)

// Common http errors
var (
	BadRequest          = NewError(http.StatusBadRequest)
	InternalServerError = NewError(http.StatusInternalServerError)
	MethodNotAllowed    = NewError(http.StatusMethodNotAllowed)
	NotAuthorized       = NewError(http.StatusUnauthorized)
)

// Error error containing status code and error.
type Error struct {
	Status int   `json:"status"`
	Err    error `json:"error"`
}

// Error retruns a string representation of an Error and
// makes the type compliant with the go error interface.
func (err Error) Error() string {
	return fmt.Sprintf("%d - %s", err.Status, err.Err)
}

// NewError creates a new Error based on a supplied status code
// attempts to derive the error message.
func NewError(status int) Error {
	return Error{
		Status: status,
		Err:    fmt.Errorf(http.StatusText(status)),
	}
}

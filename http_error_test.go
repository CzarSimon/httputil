package httputil

import (
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(http.StatusOK)
	if err.Status != 200 {
		t.Errorf("Wrong err.Status Expected=200 Got=%d", err.Status)
	}
	if err.Error() != "200 - OK" {
		t.Errorf("Wrong err.Error() Expected=[200 - OK] Got=[%s]", err.Error())
	}
}

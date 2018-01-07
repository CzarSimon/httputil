package query

import (
	"fmt"
	"net/http"
)

// ParseValue parses a query value from request.
func ParseValue(r *http.Request, key string) (string, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return value, fmt.Errorf("No value found for key: %s", key)
	}
	return value, nil
}

// ParseValues parses query values from a request
func ParseValues(r *http.Request, key string) ([]string, error) {
	values := r.URL.Query()[key]
	if len(values) < 1 {
		return nil, fmt.Errorf("No values found for key: %s", key)
	}
	return values, nil
}

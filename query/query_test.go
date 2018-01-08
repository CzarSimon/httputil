package query

import (
	"net/http/httptest"
	"testing"
)

func TestParseValue(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost/test?key=value", nil)
	val, err := ParseValue(r, "key")
	expectedValue := "value"
	if err != nil {
		t.Errorf("%s Expected=%s", err, expectedValue)
	}
	if val != expectedValue {
		t.Errorf("ParseValue returned wrong value. Expected=%s Got=%s", expectedValue, val)
	}
	r = httptest.NewRequest("GET", "http://localhost/test?notKey=value", nil)
	val, err = ParseValue(r, "key")
	if err == nil {
		t.Error("No error recived from ParseValue when expected")
	}
	if val != "" {
		t.Errorf("ParseValue returned wrong value. Exprected=\"\" Got=%s", val)
	}
}

func TestParseValues(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost/test?key=val1&key=val2", nil)
	values, err := ParseValues(r, "key")
	expectedValues := []string{"val1", "val2"}
	if len(values) != len(expectedValues) {
		t.Fatalf("Unexptected number of parsed values. Expected=2 Got=%d", len(values))
	}
	for i, expectedValue := range expectedValues {
		if values[i] != expectedValue {
			t.Errorf("%d - ParseValues returned wrong value. Expected=%s Got=%s",
				expectedValue, values[i])
		}
	}
	r = httptest.NewRequest("GET", "http://localhost/test?notKey=val1&notKey=val2", nil)
	values, err = ParseValues(r, "key")
	if err == nil {
		t.Error("No error recived from ParseValues when expected")
	}
	if len(values) != 0 {
		t.Errorf("ParseValue returned values. Exprected=[] Got=%v", values)
	}
}

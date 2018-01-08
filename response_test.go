package httputil

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CzarSimon/util"
)

func TestSendOK(t *testing.T) {
	w := httptest.NewRecorder()
	SendOK(w)
	if w.Code != http.StatusOK {
		t.Errorf("SendOK: Wrong status code. Expected=%d Got=%s", http.StatusOK, w.Code)
	}
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("SendOK: Unexpected error. Got=[%s]", err)
	}
	if string(body) != "OK" {
		t.Errorf("SendOK: Wrong body. Expeceted=[OK] Got=[%s]", string(body))
	}
}

func TestPing(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/ping", nil)
	Ping(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Ping: Wrong status code. Expected=%d Got=%s", http.StatusOK, w.Code)
	}
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Ping: Unexpected error reading body. Got=[%s]", err)
	}
	if string(body) != "OK" {
		t.Errorf("Ping: Wrong body. Expeceted=[OK] Got=[%s]", string(body))
	}
}

type testType struct {
	Value string `json:"value"`
	Num   int    `json:"num"`
}

func TestSendJSONHappyPath(t *testing.T) {
	w := httptest.NewRecorder()
	expectedJSON := testType{Value: "val", Num: 10}
	err := SendJSON(w, expectedJSON)
	if err != nil {
		t.Errorf("SendJSON: Unexpected error. Got=[%s]", err)
	}
	if w.Code != http.StatusOK {
		t.Errorf("SendJSON: Wrong status code. Expected=%d Got=%s", http.StatusOK, w.Code)
	}
	resp := w.Result()
	contentType := resp.Header.Get("Content-Type")
	if contentType != JSON {
		t.Errorf("SendJSON: Wrong Content-Type: Expected=%s Got=%s", JSON, contentType)
	}
	var res testType
	err = util.DecodeJSON(resp.Body, &res)
	if err != nil {
		t.Errorf("SendJSON: Unexpected error reading body. Got=[%s]", err)
	}
	if res.Num != expectedJSON.Num {
		t.Errorf("SendJSON: Wrong res.Num Expected=%d Got=%d", expectedJSON.Num, res.Num)
	}
	if res.Value != expectedJSON.Value {
		t.Errorf("SendJSON: Wrong res.Value Expected=%s Got=%s", expectedJSON.Value, res.Value)
	}
}

func TestSendJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	nonJSON := make(chan int)
	err := SendJSON(w, nonJSON)
	if err == nil {
		t.Errorf("SendJSON: Expected error got nil")
	}
	if err != InternalServerError {
		t.Errorf("SendJSON: Wrong error. Expected=[%s] Got=[%s]",
			InternalServerError.Error(), err.Error())
	}
}

func TestSendErr(t *testing.T) {
	w := httptest.NewRecorder()
	SendErr(w, BadRequest)
	if w.Code != http.StatusBadRequest {
		t.Errorf("SendErr: Wrong status code. Expected=%d Got=%d",
			http.StatusBadRequest, w.Code)
	}
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("SendErr: Unexpected error reading body. Got=[%s]", err)
	}
	if string(body) != BadRequest.Error()+"\n" {
		t.Errorf("SendErr: Wrong body. Expeceted=[%s] Got=[%s]",
			BadRequest.Error(), string(body))
	}
}

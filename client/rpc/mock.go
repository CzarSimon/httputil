package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/CzarSimon/httputil"
)

// MockResponse mocked rpc response.
type MockResponse struct {
	Body interface{}
	Err  error
}

// MockResponses is a MockResponse map
type MockResponses map[string]MockResponse

// MockClient mock implementation of a client.
type MockClient struct {
	Client
	Responses MockResponses
}

// Do perform a mked request.
func (c *MockClient) Do(req *http.Request) (*http.Response, error) {
	time.Sleep(10 * time.Millisecond)
	key := fmt.Sprintf("%s:%s", req.Method, req.URL)
	mockRes, ok := c.Responses[key]
	if !ok {
		err := fmt.Errorf("could not find uri %s", key)
		return nil, httputil.NotFoundError(err)
	}

	if mockRes.Err != nil {
		return nil, mockRes.Err
	}

	var body io.ReadCloser
	headers := http.Header{}
	if mockRes.Body != nil {
		bytesBody, err := json.Marshal(mockRes.Body)
		if err != nil {
			return nil, err
		}
		body = ioutil.NopCloser(bytes.NewBuffer(bytesBody))
		headers.Add("Content-Type", "application/json")
	} else {
		body = http.NoBody
	}

	status := http.StatusOK
	res := &http.Response{
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Status:     fmt.Sprintf("%d - %s", status, http.StatusText(status)),
		StatusCode: status,
		Body:       body,
		Header:     headers,
	}

	return res, nil
}

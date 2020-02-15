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
	body interface{}
	err  error
}

// MockClient mock implementation of a client.
type MockClient struct {
	Client
	Responses map[string]MockResponse
}

// Do perform a mocked request.
func (c *MockClient) Do(req *http.Request) (*http.Response, error) {
	time.Sleep(10 * time.Millisecond)
	key := fmt.Sprintf("%s:%s", req.Method, req.URL)
	mockRes, ok := c.Responses[key]
	if !ok {
		err := fmt.Errorf("could not find uri %s", key)
		return nil, httputil.NotFoundError(err)
	}

	if mockRes.err != nil {
		return nil, mockRes.err
	}

	var body io.ReadCloser
	headers := http.Header{}
	if mockRes.body != nil {
		bytesBody, err := json.Marshal(mockRes.body)
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

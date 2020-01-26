package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/id"
)

// Common errors
var (
	ErrUnsupportedContentType = errors.New("unsupported content type")
)

const (
	contentTypeJSON = "application/json"
	contentTypeText = "text/plain"

	headerContentType = "Content-Type"
)

// HasStatus check if an error is a HTTPError with the specified status.
func HasStatus(err error, status int) bool {
	var httpErr *httputil.Error
	ok := errors.As(err, &httpErr)
	return ok && httpErr.Status == status
}

// Client rpc client interface for creating and executing requrests.
type Client interface {
	CreateRequest(method, url string, body interface{}) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
}

// NewClient creates a new rpc client using the default implementation.
func NewClient(timeout time.Duration) Client {
	return &httpClient{
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

type httpClient struct {
	http *http.Client
}

func (c *httpClient) CreateRequest(method, url string, body interface{}) (*http.Request, error) {
	r, err := createBody(body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body\n%w", err)
	}

	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	res, err := c.http.Do(req)
	err = wrapRequestError(req, res, err)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func wrapRequestError(req *http.Request, res *http.Response, err error) error {
	httpErr := &httputil.Error{
		ID:      id.New(),
		Status:  http.StatusOK,
		Message: "remote request failed",
		Err:     err,
	}

	if err != nil {
		httpErr.Status = http.StatusServiceUnavailable
		httpErr.Err = err
	}

	if res != nil {
		httpErr.Status = res.StatusCode
		httpErr.Message = fmt.Sprintf("request failed, status: %s", res.Status)
	}

	if httpErr.Status >= 300 {
		return httpErr
	}

	return err
}

// DecodeJSON decodes a json response body into a value reciever.
func DecodeJSON(res *http.Response, v interface{}) error {
	contentType := res.Header.Get(headerContentType)
	if !strings.HasPrefix(contentType, contentTypeJSON) && contentType == "" {
		return fmt.Errorf("%w: %s", ErrUnsupportedContentType, contentType)
	}

	err := json.NewDecoder(res.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("failed to parse response body\n%w", err)
	}

	return nil
}

// DecodeText decodes a text/plain response body and returns it as a string.
func DecodeText(res *http.Response) (string, error) {
	contentType := res.Header.Get(headerContentType)
	if !strings.HasPrefix(contentType, contentTypeText) && contentType == "" {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedContentType, contentType)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func createBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	var bodyReader io.Reader
	if body != nil {
		bytesBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(bytesBody)
	}

	return bodyReader, nil
}

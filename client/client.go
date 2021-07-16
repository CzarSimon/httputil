package client

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/httputil/client/rpc"
	"github.com/CzarSimon/httputil/jwt"
	"github.com/CzarSimon/httputil/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var log = logger.GetDefaultLogger("httputil/client")

// Prometheus metrics.
var (
	rpcsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_requests_total",
			Help: "The total number of remote procedure calls",
		},
		[]string{"endpoint", "method", "status"},
	)
	rpcLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rpc_request_latency_ms",
			Help:    "Remote procedure call duration in milliseconds",
			Buckets: []float64{1, 5, 10, 20, 50, 100, 200, 500, 1000, 2000, 5000},
		},
		[]string{"endpoint", "method", "status"},
	)
)

// Client rest client.
type Client struct {
	Issuer    jwt.Issuer
	BaseURL   string
	Role      string
	UserAgent string
	RPCClient rpc.Client
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, v interface{}) error {
	return c.request(ctx, http.MethodGet, path, nil, v)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, body, v interface{}) error {
	return c.request(ctx, http.MethodPut, path, body, v)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body, v interface{}) error {
	return c.request(ctx, http.MethodPost, path, body, v)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string, v interface{}) error {
	return c.request(ctx, http.MethodDelete, path, nil, v)
}

func (c *Client) request(ctx context.Context, method, path string, body, v interface{}) error {
	req, err := c.RPCClient.CreateRequest(method, c.BaseURL+path, body)
	if err != nil {
		return fmt.Errorf("failed to create request\n%w", err)
	}

	c.addToken(req)
	injectSpan(ctx, req)

	timer := createTimer()
	res, err := c.RPCClient.Do(req)
	if err != nil {
		c.recordMetricsOnError(timer, path, method, res)
		return err
	}

	defer c.recordMetrics(timer, path, method, res.StatusCode)
	if v == nil {
		return nil
	}

	return rpc.DecodeJSON(res, v)
}

func (c *Client) recordMetricsOnError(timer calcDuration, path, method string, res *http.Response) {
	statusCode := http.StatusServiceUnavailable
	if res != nil {
		statusCode = res.StatusCode
	}

	c.recordMetrics(timer, path, method, statusCode)
}

func (c *Client) recordMetrics(timer calcDuration, path, method string, statusCode int) {
	latency := timer()
	endpoint := stripQueryAndUUIDs(c.BaseURL + path)
	status := strconv.Itoa(statusCode)

	rpcsTotal.WithLabelValues(endpoint, method, status).Inc()
	rpcLatency.WithLabelValues(endpoint, method, status).Observe(latency)
}

func (c *Client) addToken(req *http.Request) {
	token, err := c.Issuer.Issue(jwt.User{
		ID:    c.UserAgent,
		Roles: []string{c.Role},
	}, 24*time.Hour)
	if err != nil {
		log.Warn("failed to create auth token", zap.Error(err))
	}

	req.Header.Add("Authorization", "Bearer "+token)
}

func injectSpan(ctx context.Context, req *http.Request) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		reqID := span.BaggageItem(httputil.RequestIDHeader)
		if reqID != "" {
			req.Header.Set(httputil.RequestIDHeader, reqID)
		}

		opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
	}
}

type calcDuration func() float64

func createTimer() calcDuration {
	start := time.Now()

	// Returns latency in milliseconds.
	return func() float64 {
		end := time.Now()
		return toMilliseconds(end.Sub(start))
	}
}

func toMilliseconds(d time.Duration) float64 {
	return float64(d) / 1e6
}

var uuidRegexp = regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)

func stripQueryAndUUIDs(url string) string {
	withoutQuery := strings.Split(url, "?")[0]
	return uuidRegexp.ReplaceAllString(withoutQuery, ":id")
}

func stripQueryParameters(url string) string {
	split := strings.Split(url, "?")
	domainAndPath := split[0]
	if len(split) != 2 {
		return domainAndPath
	}

	query := strings.Split(split[1], "&")
	strippedQuery := make([]string, 0, len(query))
	for _, part := range query {
		paramAndValue := strings.Split(part, "=")
		param := paramAndValue[0]
		stripped := ":value"
		if len(paramAndValue) == 2 && len(strings.Split(paramAndValue[1], ",")) > 1 {
			stripped = ":values"
		}
		strippedQuery = append(strippedQuery, fmt.Sprintf("%s=%s", param, stripped))
	}

	return fmt.Sprintf("%s?%s", domainAndPath, strings.Join(strippedQuery, "&"))
}

package httputil

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

//  Common tracing headers
const (
	RequestIDHeader = "X-Request-ID"
	ClientIDHeader  = "X-Client-ID"
	SessionIDHeader = "X-Session-ID"
)

var requestLog = logger.GetDefaultLogger("httputil/request-log")

// Prometheus metrics.
var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "The total number served requests",
		},
		[]string{"endpoint", "method", "status"},
	)
	requestsLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_latency_ms",
			Help:    "Request latency in milliseconds",
			Buckets: []float64{1, 5, 10, 20, 50, 100, 200, 500, 1000, 2000, 5000},
		},
		[]string{"endpoint", "method", "status"},
	)
)

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Metrics records metrics about a request.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == metricsPath {
			c.Next()
			return
		}
		stop := createTimer()
		endpoint := c.FullPath()
		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		latency := stop()
		requestsTotal.WithLabelValues(endpoint, method, status).Inc()
		requestsLatency.WithLabelValues(endpoint, method, status).Observe(latency)
	}
}

// RequestID ensures that a request contains a request id. If not sets a UUID as there request id.
func RequestID(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(key)
		if reqID == "" {
			reqID = id.New()
			c.Request.Header.Set(key, reqID)
		}

		c.Header(key, reqID)
		c.Next()
	}
}

// Trace captures open tracing span and attaches it to the request context.
func Trace(app string, headers ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		wireContext, _ := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)

		spanName := fmt.Sprintf("%s - %s %s", app, c.Request.Method, c.Request.URL.Path)
		span := opentracing.StartSpan(spanName, ext.RPCServerOption(wireContext))
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())

		for _, h := range headers {
			setBaggageIfMissing(span, h, c.GetHeader(h))
		}

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))
		c.Next()

		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		span.Finish()
	}
}

func setBaggageIfMissing(span opentracing.Span, key, val string) {
	baggage := span.BaggageItem(key)
	if baggage != "" || val == "" {
		return
	}

	span.SetBaggageItem(key, val)
}

// Logger request logging middleware.
// Accepts a list of paths to skip logging incomming requests and non 500 outgoing requests to.
func Logger(skip ...string) gin.HandlerFunc {
	skipPaths := make(map[string]bool)
	for _, path := range skip {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		stop := createTimer()
		path := c.Request.URL.Path
		reqID := c.GetHeader(RequestIDHeader)
		_, skippablePath := skipPaths[path]

		if !skippablePath {
			requestLog.Info(
				fmt.Sprintf("Incomming request: %s %s", c.Request.Method, path),
				zap.String("requestId", reqID),
			)
		}

		c.Next()

		latency := stop()
		status := c.Writer.Status()

		if skippablePath && status < http.StatusInternalServerError {
			return
		}

		logFn := requestLog.Info
		if status >= http.StatusInternalServerError {
			logFn = requestLog.Error
		}

		logFn(
			fmt.Sprintf("Outgoing request: %s %s", c.Request.Method, path),
			zap.String("requestId", reqID),
			zap.Int("status", status),
			zap.Float64("latency", latency),
		)
	}
}

type calcDuration func() float64

func createTimer() calcDuration {
	start := time.Now()

	// Returns latency in milliseconds.
	return func() float64 {
		end := time.Now()
		return float64(end.Sub(start)) / 1e6
	}
}

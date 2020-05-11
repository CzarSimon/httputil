package httputil

import (
	"fmt"
	"strconv"
	"time"

	"github.com/CzarSimon/httputil/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	metricsPath = "/metrics"
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
			Name: "http_request_latency_ms",
			Help: "Request latency in milliseconds",
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

// Trace captures open tracing span and attaches it to the request context.
func Trace(app string) gin.HandlerFunc {
	return func(c *gin.Context) {
		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			errLog.Debug("failed to extract wireContext from request", zap.Error(err))
		}

		span := opentracing.StartSpan(app, ext.RPCServerOption(wireContext))
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))
		c.Next()

		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		span.Finish()
	}
}

// Logger request logging middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		stop := createTimer()
		path := c.Request.URL.Path
		requestLog.Info(fmt.Sprintf("Incomming request: %s %s", c.Request.Method, path))

		c.Next()

		latency := stop()
		requestLog.Info(fmt.Sprintf("Outgoing request: %s %s", c.Request.Method, path),
			zap.Int("status", c.Writer.Status()),
			zap.Float64("latency", latency))
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

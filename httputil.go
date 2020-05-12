package httputil

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthFunc health check function signature.
type HealthFunc func() error

// NewRouter creates a default router.
func NewRouter(appName string, healthCheck HealthFunc) *gin.Engine {
	return NewCustomRouter(
		healthCheck,
		gin.Recovery(),
		Trace(appName),
		Metrics(),
		Logger(),
		HandleErrors(),
	)
}

// NewCustomRouter creates a new router with a custom list of base middlewares.
func NewCustomRouter(healthCheck HealthFunc, middlewares ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(middlewares...)

	r.GET("/health", checkHealth(healthCheck))
	r.GET(metricsPath, prometheusHandler())
	return r
}

// SendOK sends an ok status and message to the client.
func SendOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// ParseQueryValue parses a query value from request.
func ParseQueryValue(c *gin.Context, key string) (string, error) {
	value, ok := c.GetQuery(key)
	if !ok {
		err := fmt.Errorf("No value found for param: %s", key)
		return "", BadRequestError(err)
	}
	return value, nil
}

// ParseQueryValues parses query values from a request.
func ParseQueryValues(c *gin.Context, key string) ([]string, error) {
	values, ok := c.GetQueryArray(key)
	if !ok {
		err := fmt.Errorf("No value found for param: %s", key)
		return nil, BadRequestError(err)
	}
	return values, nil
}

// AllowContentType whitelists a given list of content types.
func AllowContentType(types ...string) gin.HandlerFunc {
	allowedTypes := make(map[string]bool)
	for _, t := range types {
		allowedTypes[t] = true
	}

	return func(c *gin.Context) {
		ct := c.ContentType()
		if _, ok := allowedTypes[ct]; ok {
			c.Next()
			return
		}

		err := UnsupportedMediaTypeError(fmt.Errorf("unsupported content-type: %s", ct))
		logError(c, err)
		c.AbortWithStatusJSON(err.Status, err)
	}
}

// AllowJSON only allowes request with content type json.
func AllowJSON() gin.HandlerFunc {
	return AllowContentType(gin.MIMEJSON)
}

func checkHealth(check HealthFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := check()
		if err == nil {
			SendOK(c)
			return
		}

		c.Error(err)
	}
}

package environ

import (
	"os"

	"github.com/CzarSimon/httputil/logger"
	"go.uber.org/zap"
)

var log = logger.MustGetLogger("httputil/environ", zap.InfoLevel)

// Get gets environment variable with default value.
func Get(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// MustGet gets environment variable and panics if not found
func MustGet(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Panic("Failed to get environment variable", zap.String("name", name))
	}

	return value
}

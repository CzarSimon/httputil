package logger

import (
	"log"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const logLevelEnvironmentVariable = "LOG_LEVEL"

var logEventsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "log_events_total",
		Help: "Total log events with name and level",
	},
	[]string{"name", "level"},
)

// GetLogger creates a named logger for internal application logs.
func GetLogger(name string, level zapcore.Level) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey: "message",
		LevelKey:   "level",
		TimeKey:    "time",
		CallerKey:  "caller",

		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	hooks := zap.Hooks(metricsHook)
	logger, err := cfg.Build(hooks)
	if err != nil {
		return logger, err
	}

	logger = logger.With(zap.String("logger", name))
	return logger.Named(name), nil
}

// MustGetLogger creates a names logger and panics on failure.
func MustGetLogger(name string, level zapcore.Level) *zap.Logger {
	logger, err := GetLogger(name, level)
	if err != nil {
		log.Fatalln("Failed to get zap.Logger "+name, err)
	}

	return logger
}

// GetDefaultLogger gets a named logger using the default implementation.
func GetDefaultLogger(name string) *zap.Logger {
	logLevel := getDefaultLogLevel()
	return MustGetLogger(name, logLevel)
}

func metricsHook(entry zapcore.Entry) error {
	logEventsTotal.WithLabelValues(entry.LoggerName, entry.Level.String()).Inc()
	return nil
}

func getDefaultLogLevel() zapcore.Level {
	level := getLogLevelFromEnv()

	switch level {
	case zap.DebugLevel.String():
		return zap.DebugLevel
	case zap.InfoLevel.String():
		return zap.InfoLevel
	case zap.WarnLevel.String():
		return zap.WarnLevel
	case zap.ErrorLevel.String():
		return zap.ErrorLevel
	default:
		return zap.DebugLevel
	}
}

func getLogLevelFromEnv() string {
	value := os.Getenv(logLevelEnvironmentVariable)
	if value == "" {
		return zap.DebugLevel.String()
	}

	return strings.ToLower(value)
}

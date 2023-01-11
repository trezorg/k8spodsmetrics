package logger

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/exp/slog"
)

var log *slog.Logger
var onceLog sync.Once

var defaultLogger = "INFO"
var serviceName = "k8spodmetrics"

func levelToSlogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARNING", "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		panic(fmt.Sprintf("Unknown level: %s", level))
	}
}

func initLogger(logLevel string) {
	opts := slog.HandlerOptions{
		Level: levelToSlogLevel(logLevel),
	}
	jsonHandler := opts.NewJSONHandler(os.Stdout).WithAttrs([]slog.Attr{slog.String("service", serviceName)})
	log = slog.New(jsonHandler)
}

// InitLogger initializes logger instance
func InitLogger(logLevel string, logPrettyPrint bool) {
	onceLog.Do(func() {
		initLogger(logLevel)
	})
}

// InitDefaultLogger initializes logger instance
func InitDefaultLogger() {
	InitLogger(defaultLogger, true)
}

// Info logs a message at level Info
func Info(message string, args ...any) {
	log.Info(message, args...)
}

// Error logs a message at level Info
func Error(message string, err error, args ...any) {
	log.Error(message, err, args...)
}

// Warn logs a message at level Info
func Warn(message string, args ...any) {
	log.Warn(message, args...)
}
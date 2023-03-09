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

func LevelToSlogLevel(level string) (slog.Level, error) {
	switch level {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARNING", "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelError, fmt.Errorf("Unknown level: %s", level)
	}
}

func initLogger(logLevel string) {
	level, err := LevelToSlogLevel(logLevel)
	if err != nil {
		panic(err)
	}
	opts := slog.HandlerOptions{
		Level: level,
	}
	jsonHandler := opts.NewJSONHandler(os.Stdout).WithAttrs([]slog.Attr{slog.String("service", serviceName)})
	log = slog.New(jsonHandler)
}

// InitLogger initializes logger instance
func InitLogger(logLevel string) {
	onceLog.Do(func() {
		initLogger(logLevel)
	})
}

// InitDefaultLogger initializes logger instance
func InitDefaultLogger() {
	InitLogger(defaultLogger)
}

// Info logs a message at level Info
func Info(message string, args ...any) {
	log.Info(message, args...)
}

// Error logs a message at level Error
func Error(message string, err error, args ...any) {
	log.Error(message, err, args...)
}

// Warn logs a message at level Warn
func Warn(message string, args ...any) {
	log.Warn(message, args...)
}

func Debug(message string, args ...any) {
	log.Debug(message, args...)
}

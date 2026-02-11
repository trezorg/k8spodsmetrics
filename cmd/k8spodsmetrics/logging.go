package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

const serviceName = "k8spodmetrics"

func levelFromString(levelStr string) (slog.Level, error) {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO", "":
		return slog.LevelInfo, nil
	case "WARN", "WARNING":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelError, fmt.Errorf("unknown log level: %s", levelStr)
	}
}

func initLogger(levelStr string) error {
	level, err := levelFromString(levelStr)
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler).With(slog.String("service", serviceName)))
	return err
}

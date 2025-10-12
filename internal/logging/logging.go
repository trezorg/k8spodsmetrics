package logging

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

const (
	serviceName = "k8spodmetrics"
)

func levelFromString(levelStr string) (level slog.Level, err error) {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO", "":
		level = slog.LevelInfo
	case "WARN", "WARNING":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		err = fmt.Errorf("unknown log level: %s", levelStr)
	}
	return level, err
}

func Init(levelStr string) error {
	level, err := levelFromString(levelStr)
	if err != nil {
		return err
	}
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler).With(slog.String("service", serviceName)))
	return nil
}

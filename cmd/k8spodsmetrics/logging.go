package main

import (
	"log/slog"
	"os"
	"strings"
)

func levelFromString(s string) slog.Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}

func initLogger(levelStr string) {
	lvl := levelFromString(levelStr)
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: lvl})
	slog.SetDefault(slog.New(handler).With(slog.String("service", "k8spodmetrics")))
}

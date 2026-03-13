package main

import (
	"log/slog"
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/adapters/stdin"
	"github.com/urfave/cli/v2"
)

var version = "0.0.1"

func run(args []string) error {
	app := stdin.NewApp(version)
	// Initialize logging once using global flags via Before hook
	app.Before = func(c *cli.Context) error {
		return initLogger(c.String("loglevel"))
	}
	return app.Run(args)
}

func main() {
	if err := initLogger(""); err != nil {
		slog.Error("command failed", "error", err)
		os.Exit(1)
	}

	if err := run(os.Args); err != nil {
		slog.Error("command failed", "error", err)
		os.Exit(1)
	}
}

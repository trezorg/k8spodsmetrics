package main

import (
	"log"
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/adapters/stdin"
	"github.com/urfave/cli/v2"
)

var version = "0.0.1"

func main() {
	app := stdin.NewApp(version)
	// Initialize logging once using global flags via Before hook
	app.Before = func(c *cli.Context) error {
		return initLogger(c.String("loglevel"))
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

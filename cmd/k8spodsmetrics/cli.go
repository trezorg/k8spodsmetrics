package main

import (
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/urfave/cli/v2"
)

func processArgs() error {

	config := metricsresources.Config{}

	app := cli.NewApp()
	app.Version = version
	app.Authors = []*cli.Author{{
		Name:  "Igor Nemilentsev",
		Email: "trezorg@gmail.com",
	}}
	app.Usage = "K8S pod metrics"
	app.AllowExtFlags = true
	app.EnableBashCompletion = true
	app.Description = "The application shows pods metrics"
	app.Action = func(c *cli.Context) error {
		return processK8sMetrics(config)
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "kubeconfig",
			Aliases:     []string{"k"},
			Value:       "",
			Usage:       "K8S config",
			Destination: &config.KubeConfig,
		},
		&cli.StringFlag{
			Name:        "context",
			Aliases:     []string{"c"},
			Value:       "",
			Usage:       "K8S config context",
			Destination: &config.KubeContext,
		},
		&cli.StringFlag{
			Name:        "namespace",
			Aliases:     []string{"n"},
			Value:       "",
			Usage:       "K8S namespace",
			Destination: &config.Namespace,
		},
		&cli.StringFlag{
			Name:        "label",
			Aliases:     []string{"l"},
			Value:       "",
			Usage:       "K8S pod label",
			Destination: &config.Label,
		},
		&cli.BoolFlag{
			Name:        "alerts",
			Aliases:     []string{"a"},
			Value:       false,
			Usage:       "Show only metrics with alert",
			Destination: &config.OnlyAlert,
		},
		&cli.BoolFlag{
			Name:        "watch",
			Aliases:     []string{"w"},
			Value:       false,
			Usage:       "Watch for metrics for some period",
			Destination: &config.WatchMetrics,
		},
		&cli.UintFlag{
			Name:        "watch-period",
			Aliases:     []string{"p"},
			Value:       5,
			Usage:       "Watch period",
			Destination: &config.WatchPeriod,
		},
		&cli.StringFlag{
			Name:        "loglevel",
			Aliases:     []string{"level"},
			Value:       "INFO",
			Usage:       "Log level",
			Destination: &config.LogLevel,
		},
		&cli.UintFlag{
			Name:        "kloglevel",
			Aliases:     []string{"klevel"},
			Value:       3,
			Usage:       "k8s client log level",
			Destination: &config.KLogLevel,
		},
	}
	if err := app.Run(os.Args); err != nil {
		return err
	}
	if _, err := logger.LevelToSlogLevel(config.LogLevel); err != nil {
		return err
	}
	return nil
}

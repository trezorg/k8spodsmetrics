package stdin

import (
	"fmt"

	"github.com/trezorg/k8spodsmetrics/internal/alert"
	"github.com/trezorg/k8spodsmetrics/internal/output"
	"github.com/urfave/cli/v2"
)

const (
	defaultKlogLevel          = 3
	defaultWatchPeriodSeconds = 5
)

func commonFlags(config *commonConfig) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Value:       "",
			Usage:       "Config file path (YAML format)",
			Destination: &config.ConfigFile,
		},
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
			Name:    "loglevel",
			Aliases: []string{"level"},
			Value:   "INFO",
			Usage:   "Log level",
		},
		&cli.UintFlag{
			Name:        "kloglevel",
			Aliases:     []string{"klevel"},
			Value:       defaultKlogLevel,
			Usage:       "k8s client log level",
			Destination: &config.KLogLevel,
		},
		&cli.StringFlag{
			Name:        "alerts",
			Aliases:     []string{"a"},
			Value:       string(alert.None),
			Usage:       fmt.Sprintf("Alert format. [%s]", alert.StringListDefault()),
			Destination: &config.Alert,
			Action: func(_ *cli.Context, value string) error {
				if err := alert.Valid(alert.Alert(value)); err != nil {
					return err
				}
				config.Alert = value
				return nil
			},
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
			Value:       defaultWatchPeriodSeconds,
			Usage:       "Watch period",
			Destination: &config.WatchPeriod,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Value:       string(output.Table),
			Usage:       fmt.Sprintf("Output format. [%s]", output.StringListDefault()),
			Destination: &config.Output,
			Action: func(_ *cli.Context, value string) error {
				if err := output.Valid(output.Output(value)); err != nil {
					return err
				}
				config.Output = value
				return nil
			},
		},
	}
}

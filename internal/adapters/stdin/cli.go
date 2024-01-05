package stdin

import (
	"fmt"
	"os"

	metricsjson "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/json/metricsresources"
	metricsscreen "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/screen/metricsresources"
	metricsstring "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/string/metricsresources"
	metricstable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/metricsresources"
	metricsyaml "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/yaml/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/alert"

	nodesjson "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/json/noderesources"
	nodesscreen "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/screen/noderesources"
	nodesstring "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/string/noderesources"
	nodestable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/noderesources"
	nodesyaml "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/yaml/noderesources"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/output"
	metricssorting "github.com/trezorg/k8spodsmetrics/internal/sorting/metricsresources"
	nodesorting "github.com/trezorg/k8spodsmetrics/internal/sorting/noderesources"
	"github.com/urfave/cli/v2"
)

type config struct {
	KubeConfig   string
	KubeContext  string
	LogLevel     string
	Output       string
	Alert        string
	KLogLevel    uint
	WatchPeriod  uint
	WatchMetrics bool
}

type podConfig struct {
	Namespace string
	Label     string
	Sorting   string
	config
	Reverse bool
}

type summaryConfig struct {
	Name    string
	Label   string
	Sorting string
	config
	Reverse bool
}

func metricsResourcesConfig(c podConfig) metricsresources.Config {
	return metricsresources.Config{
		KubeConfig:   c.KubeConfig,
		KubeContext:  c.KubeContext,
		Namespace:    c.Namespace,
		Label:        c.Label,
		LogLevel:     c.LogLevel,
		Output:       c.Output,
		Sorting:      c.Sorting,
		Reverse:      c.Reverse,
		KLogLevel:    c.KLogLevel,
		Alert:        c.Alert,
		WatchMetrics: c.WatchMetrics,
		WatchPeriod:  c.WatchPeriod,
	}
}

func nodeResourcesConfig(c summaryConfig) noderesources.Config {
	return noderesources.Config{
		KubeConfig:   c.KubeConfig,
		KubeContext:  c.KubeContext,
		LogLevel:     c.LogLevel,
		Label:        c.Label,
		Name:         c.Name,
		Output:       c.Output,
		Sorting:      c.Sorting,
		Reverse:      c.Reverse,
		KLogLevel:    c.KLogLevel,
		Alert:        c.Alert,
		WatchMetrics: c.WatchMetrics,
		WatchPeriod:  c.WatchPeriod,
	}
}

type SummaryProcessor interface {
	Process(noderesources.SuccessProcessor) error
}

type PodsProcessor interface {
	Process(metricsresources.SuccessProcessor) error
}

type SummaryWatcher interface {
	ProcessWatch(noderesources.SuccessProcessor, noderesources.ErrorProcessor) error
}

type PodsWatcher interface {
	ProcessWatch(metricsresources.SuccessProcessor, metricsresources.ErrorProcessor) error
}

func summaryOutputProcessor(out output.Output) noderesources.SuccessProcessor {
	switch out {
	case output.Table:
		return nodestable.Table(nodestable.Print)
	case output.Json:
		return nodesjson.Json(nodesjson.Print)
	case output.Yaml:
		return nodesyaml.Yaml(nodesyaml.Print)
	case output.String:
		return nodesstring.String(nodesstring.Print)
	default:
		return nodestable.Table(nodestable.Print)
	}
}

func podsOutputProcessor(out output.Output) metricsresources.SuccessProcessor {
	switch out {
	case output.Table:
		return metricstable.Table(metricstable.Print)
	case output.Json:
		return metricsjson.Json(metricsjson.Print)
	case output.Yaml:
		return metricsyaml.Yaml(metricsyaml.Print)
	case output.String:
		return metricsstring.String(metricsstring.Print)
	default:
		return metricstable.Table(metricstable.Print)
	}
}

func summary(processor SummaryProcessor, successProcessor noderesources.SuccessProcessor) error {
	return processor.Process(successProcessor)
}

func summaryWatch(processor SummaryWatcher, successProcessor noderesources.SuccessProcessor, errorProcessor noderesources.ErrorProcessor) error {
	return processor.ProcessWatch(nodesscreen.NewScreenSuccessWriter(successProcessor), nodesscreen.NewScreenErrorWriter(errorProcessor))
}

func pods(processor PodsProcessor, successProcessor metricsresources.SuccessProcessor) error {
	return processor.Process(successProcessor)
}

func podsWatch(processor PodsWatcher, successProcessor metricsresources.SuccessProcessor, errorProcessor metricsresources.ErrorProcessor) error {
	return processor.ProcessWatch(metricsscreen.NewScreenSuccessWriter(successProcessor), metricsscreen.NewScreenErrorWriter(errorProcessor))
}

func Start(version string) error {
	config := config{}

	app := cli.NewApp()
	app.Version = version
	app.Authors = []*cli.Author{{
		Name:  "Igor Nemilentsev",
		Email: "trezorg@gmail.com",
	}}
	app.Usage = "K8S pod and node metrics"
	app.AllowExtFlags = true
	app.EnableBashCompletion = true
	app.Description = "The application shows pod and node metrics"
	app.Action = func(c *cli.Context) error {
		return nil
	}
	app.Commands = []*cli.Command{
		{
			Name:    "summary",
			Aliases: []string{"s"},
			Action: func(c *cli.Context) error {
				summaryConfig := summaryConfig{config: config}
				summaryConfig.Name = c.String("name")
				summaryConfig.Label = c.String("label")
				summaryConfig.Sorting = c.String("sorting")
				summaryConfig.Reverse = c.Bool("reverse")
				config := nodeResourcesConfig(summaryConfig)
				outputProcessor := summaryOutputProcessor(output.Output(summaryConfig.Output))
				errorProcessor := outputProcessor.(noderesources.ErrorProcessor)
				if summaryConfig.WatchMetrics {
					return summaryWatch(config, outputProcessor, errorProcessor)
				}
				return summary(config, outputProcessor)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "label",
					Aliases: []string{"l"},
					Value:   "",
					Usage:   "K8S node label",
				},
				&cli.StringFlag{
					Name:    "name",
					Aliases: []string{"n"},
					Value:   "",
					Usage:   "K8S node name",
				},
				&cli.StringFlag{
					Name:    "sorting",
					Aliases: []string{"s"},
					Value:   "name",
					Usage:   fmt.Sprintf("Sorting. [%s]", nodesorting.StringListDefault()),
					Action: func(_ *cli.Context, value string) error {
						if err := nodesorting.Valid(nodesorting.Sorting(value)); err != nil {
							return err
						}
						return nil
					},
				},
				&cli.BoolFlag{
					Name:    "reverse",
					Aliases: []string{"r"},
					Value:   false,
					Usage:   "Reverse sort",
				},
			},
		},
		{
			Name:    "pods",
			Aliases: []string{"p"},
			Action: func(c *cli.Context) error {
				podConfig := podConfig{config: config}
				podConfig.Namespace = c.String("namespace")
				podConfig.Label = c.String("label")
				podConfig.Sorting = c.String("sorting")
				podConfig.Reverse = c.Bool("reverse")
				config := metricsResourcesConfig(podConfig)
				outputProcessor := podsOutputProcessor(output.Output(podConfig.Output))
				errorProcessor := outputProcessor.(metricsresources.ErrorProcessor)
				if podConfig.WatchMetrics {
					return podsWatch(config, outputProcessor, errorProcessor)
				}
				return pods(config, outputProcessor)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "namespace",
					Aliases: []string{"n"},
					Value:   "",
					Usage:   "K8S namespace",
				},
				&cli.StringFlag{
					Name:    "label",
					Aliases: []string{"l"},
					Value:   "",
					Usage:   "K8S pod label",
				},
				&cli.StringFlag{
					Name:    "sorting",
					Aliases: []string{"s"},
					Value:   "namespace",
					Usage:   fmt.Sprintf("Sorting. [%s]", metricssorting.StringListDefault()),
					Action: func(_ *cli.Context, value string) error {
						if err := metricssorting.Valid(metricssorting.Sorting(value)); err != nil {
							return err
						}
						return nil
					},
				},
				&cli.BoolFlag{
					Name:    "reverse",
					Aliases: []string{"r"},
					Value:   false,
					Usage:   "Reverse sort",
				},
			},
		},
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
			Value:       5,
			Usage:       "Watch period",
			Destination: &config.WatchPeriod,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Value:       string(output.String),
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
	if err := app.Run(os.Args); err != nil {
		return err
	}
	return nil
}

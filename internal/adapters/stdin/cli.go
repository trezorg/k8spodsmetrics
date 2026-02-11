package stdin

import (
	metricsjson "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/json/metricsresources"
	metricsscreen "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/screen/metricsresources"
	metricsstring "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/string/metricsresources"
	metricstable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/metricsresources"
	metricsyaml "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/yaml/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"

	nodesjson "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/json/noderesources"
	nodesscreen "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/screen/noderesources"
	nodesstring "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/string/noderesources"
	nodestable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/noderesources"
	nodesyaml "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/yaml/noderesources"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/output"
	"github.com/urfave/cli/v2"
)

type commonConfig struct {
	KubeConfig   string
	KubeContext  string
	Output       string
	Alert        string
	KLogLevel    uint
	WatchPeriod  uint
	WatchMetrics bool
}

type podConfig struct {
	Namespace     string
	Label         string
	FieldSelector string
	Nodes         []string
	Sorting       string
	Resources     []string
	commonConfig
	Reverse bool
}

type summaryConfig struct {
	Name      string
	Label     string
	Sorting   string
	Resources []string
	commonConfig
	Reverse bool
}

type SummaryProcessor interface {
	Process(noderesources.SuccessProcessor) error
}

type SummaryOutputProcessor interface {
	noderesources.SuccessProcessor
	noderesources.ErrorProcessor
}

type PodsProcessor interface {
	Process(metricsresources.SuccessProcessor) error
}

type PodsOutputProcessor interface {
	metricsresources.SuccessProcessor
	metricsresources.ErrorProcessor
}

type SummaryWatcher interface {
	ProcessWatch(noderesources.SuccessProcessor, noderesources.ErrorProcessor) error
}

type PodsWatcher interface {
	ProcessWatch(metricsresources.SuccessProcessor, metricsresources.ErrorProcessor) error
}

func summaryOutputProcessor(out output.Output, res resources.Resources) SummaryOutputProcessor {
	switch out {
	case output.Table:
		return nodestable.ToTable(res)
	case output.JSON:
		return nodesjson.JSON(nodesjson.Print)
	case output.Yaml:
		return nodesyaml.Yaml(nodesyaml.Print)
	case output.String:
		return nodesstring.String(nodesstring.Print)
	}
	return nodestable.ToTable(res)
}

func podsOutputProcessor(out output.Output, res resources.Resources) PodsOutputProcessor {
	switch out {
	case output.Table:
		return metricstable.ToTable(res)
	case output.JSON:
		return metricsjson.JSON(metricsjson.Print)
	case output.Yaml:
		return metricsyaml.Yaml(metricsyaml.Print)
	case output.String:
		return metricsstring.String(metricsstring.Print)
	}
	return metricstable.ToTable(res)
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

func NewApp(version string) *cli.App {
	config := commonConfig{}

	app := cli.NewApp()
	app.Version = version
	app.DefaultCommand = "summary"
	app.Authors = []*cli.Author{{
		Name:  "Igor Nemilentsev",
		Email: "trezorg@gmail.com",
	}}
	app.Usage = "K8S pod and node metrics"
	app.AllowExtFlags = true
	app.EnableBashCompletion = true
	app.Description = "The application shows pod and node metrics"
	app.Action = func(_ *cli.Context) error {
		return nil
	}
	app.Commands = []*cli.Command{
		{
			Name:    "summary",
			Aliases: []string{"s"},
			Action: func(c *cli.Context) error {
				summaryActionConfig := summaryConfig{commonConfig: config}
				summaryActionConfig.Name = c.String("name")
				summaryActionConfig.Label = c.String("label")
				summaryActionConfig.Sorting = c.String("sorting")
				summaryActionConfig.Reverse = c.Bool("reverse")
				cmdResources := c.StringSlice("resource")
				outputResources := resources.FromStrings(cmdResources...)
				if err := resources.Valid(outputResources...); err != nil {
					return err
				}
				summaryActionConfig.Resources = resources.ToStrings(outputResources...)
				summaryConfig := nodeResourcesConfig(summaryActionConfig)
				outputProcessor := summaryOutputProcessor(output.Output(summaryActionConfig.Output), outputResources)
				if summaryActionConfig.WatchMetrics {
					return summaryWatch(summaryConfig, outputProcessor, outputProcessor)
				}
				return summary(summaryConfig, outputProcessor)
			},
			Flags: summaryFlags(),
		},
		{
			Name:    "pods",
			Aliases: []string{"p"},
			Action: func(c *cli.Context) error {
				podActionConfig := podConfig{commonConfig: config}
				podActionConfig.Namespace = c.String("namespace")
				podActionConfig.Label = c.String("label")
				podActionConfig.FieldSelector = c.String("field-selector")
				podActionConfig.Sorting = c.String("sorting")
				podActionConfig.Reverse = c.Bool("reverse")
				podActionConfig.Nodes = c.StringSlice("node")
				cmdResources := c.StringSlice("resource")
				outputResources := resources.FromStrings(cmdResources...)
				if err := resources.Valid(outputResources...); err != nil {
					return err
				}
				podActionConfig.Resources = resources.ToStrings(outputResources...)
				podConfig := metricsResourcesConfig(podActionConfig)
				outputProcessor := podsOutputProcessor(output.Output(podActionConfig.Output), outputResources)
				if podActionConfig.WatchMetrics {
					return podsWatch(podConfig, outputProcessor, outputProcessor)
				}
				return pods(podConfig, outputProcessor)
			},
			Flags: podsFlags(),
		},
	}
	app.Flags = commonFlags(&config)
	return app
}

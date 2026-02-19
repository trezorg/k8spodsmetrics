package stdin

import (
	metricsjson "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/json/metricsresources"
	metricsscreen "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/screen/metricsresources"
	metricsstring "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/string/metricsresources"
	metricstable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/metricsresources"
	metricsyaml "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/yaml/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/config"
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
	ConfigFile   string
	KubeConfig   string
	KubeContext  string
	Output       string
	Alert        string
	WatchPeriod  uint
	WatchMetrics bool
	fileConfig   *config.Config
}

type podConfig struct {
	Namespaces    []string
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
	case output.Text:
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
	case output.Text:
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

func loadFileConfig(configFile string) (*config.Config, error) {
	if configFile == "" {
		return nil, nil
	}
	return config.Load(configFile)
}

// applyCommonConfig merges file config with CLI common config values.
// CLI values take precedence over file config for string and numeric types.
func applyCommonConfig(cfg *commonConfig, fileConfig *config.Config) config.Common {
	merged := config.Common{
		KubeConfig:   cfg.KubeConfig,
		KubeContext:  cfg.KubeContext,
		Output:       cfg.Output,
		Alert:        cfg.Alert,
		WatchPeriod:  cfg.WatchPeriod,
		WatchMetrics: cfg.WatchMetrics,
	}
	if fileConfig != nil {
		fileConfig.MergeCommon(&merged)
	}
	return merged
}

// applyPodsConfig merges file config with CLI pods command config values.
// CLI values take precedence over file config for string and slice types.
func applyPodsConfig(podCfg *podConfig, fileConfig *config.Config) config.Pods {
	merged := config.Pods{
		Namespaces:    podCfg.Namespaces,
		Label:         podCfg.Label,
		FieldSelector: podCfg.FieldSelector,
		Nodes:         podCfg.Nodes,
		Sorting:       podCfg.Sorting,
		Reverse:       podCfg.Reverse,
		Resources:     podCfg.Resources,
	}
	if fileConfig != nil {
		fileConfig.MergePods(&merged)
	}
	return merged
}

// applySummaryConfig merges file config with CLI summary command config values.
// CLI values take precedence over file config for string and slice types.
func applySummaryConfig(summaryCfg *summaryConfig, fileConfig *config.Config) config.Summary {
	merged := config.Summary{
		Name:      summaryCfg.Name,
		Label:     summaryCfg.Label,
		Sorting:   summaryCfg.Sorting,
		Reverse:   summaryCfg.Reverse,
		Resources: summaryCfg.Resources,
	}
	if fileConfig != nil {
		fileConfig.MergeSummary(&merged)
	}
	return merged
}

func NewApp(version string) *cli.App {
	cfg := commonConfig{}

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
			Before: func(_ *cli.Context) error {
				var err error
				cfg.fileConfig, err = loadFileConfig(cfg.ConfigFile)
				return err
			},
			Action: func(c *cli.Context) error {
				summaryActionConfig := summaryConfig{commonConfig: cfg}
				summaryActionConfig.Name = c.String("name")
				summaryActionConfig.Label = c.String("label")
				summaryActionConfig.Sorting = c.String("sorting")
				summaryActionConfig.Reverse = c.Bool("reverse")
				cmdResources := c.StringSlice("resource")

				mergedSummary := applySummaryConfig(&summaryActionConfig, cfg.fileConfig)
				summaryActionConfig.Name = mergedSummary.Name
				summaryActionConfig.Label = mergedSummary.Label
				summaryActionConfig.Sorting = mergedSummary.Sorting
				summaryActionConfig.Reverse = mergedSummary.Reverse
				summaryActionConfig.Resources = mergedSummary.Resources

				mergedCommon := applyCommonConfig(&cfg, cfg.fileConfig)
				summaryActionConfig.KubeConfig = mergedCommon.KubeConfig
				summaryActionConfig.KubeContext = mergedCommon.KubeContext
				summaryActionConfig.Output = mergedCommon.Output
				summaryActionConfig.Alert = mergedCommon.Alert
				summaryActionConfig.WatchPeriod = mergedCommon.WatchPeriod
				summaryActionConfig.WatchMetrics = mergedCommon.WatchMetrics

				if len(cmdResources) == 0 && len(mergedSummary.Resources) > 0 {
					cmdResources = mergedSummary.Resources
				}

				outputResources := resources.FromStrings(cmdResources...)
				if err := resources.Valid(outputResources...); err != nil {
					return err
				}
				summaryActionConfig.Resources = resources.ToStrings(outputResources...)
				summaryCfg := nodeResourcesConfig(summaryActionConfig)
				outputProcessor := summaryOutputProcessor(output.Output(summaryActionConfig.Output), outputResources)
				if summaryActionConfig.WatchMetrics {
					return summaryWatch(summaryCfg, outputProcessor, outputProcessor)
				}
				return summary(summaryCfg, outputProcessor)
			},
			Flags: summaryFlags(),
		},
		{
			Name:    "pods",
			Aliases: []string{"p"},
			Before: func(_ *cli.Context) error {
				var err error
				cfg.fileConfig, err = loadFileConfig(cfg.ConfigFile)
				return err
			},
			Action: func(c *cli.Context) error {
				podActionConfig := podConfig{commonConfig: cfg}
				podActionConfig.Namespaces = c.StringSlice("namespace")
				podActionConfig.Label = c.String("label")
				podActionConfig.FieldSelector = c.String("field-selector")
				podActionConfig.Sorting = c.String("sorting")
				podActionConfig.Reverse = c.Bool("reverse")
				podActionConfig.Nodes = c.StringSlice("node")
				cmdResources := c.StringSlice("resource")

				mergedPods := applyPodsConfig(&podActionConfig, cfg.fileConfig)
				podActionConfig.Namespaces = mergedPods.Namespaces
				podActionConfig.Label = mergedPods.Label
				podActionConfig.FieldSelector = mergedPods.FieldSelector
				podActionConfig.Nodes = mergedPods.Nodes
				podActionConfig.Sorting = mergedPods.Sorting
				podActionConfig.Reverse = mergedPods.Reverse
				podActionConfig.Resources = mergedPods.Resources

				mergedCommon := applyCommonConfig(&cfg, cfg.fileConfig)
				podActionConfig.KubeConfig = mergedCommon.KubeConfig
				podActionConfig.KubeContext = mergedCommon.KubeContext
				podActionConfig.Output = mergedCommon.Output
				podActionConfig.Alert = mergedCommon.Alert
				podActionConfig.WatchPeriod = mergedCommon.WatchPeriod
				podActionConfig.WatchMetrics = mergedCommon.WatchMetrics

				if len(cmdResources) == 0 && len(mergedPods.Resources) > 0 {
					cmdResources = mergedPods.Resources
				}

				outputResources := resources.FromStrings(cmdResources...)
				if err := resources.Valid(outputResources...); err != nil {
					return err
				}
				podActionConfig.Resources = resources.ToStrings(outputResources...)
				podCfg := metricsResourcesConfig(podActionConfig)
				outputProcessor := podsOutputProcessor(output.Output(podActionConfig.Output), outputResources)
				if podActionConfig.WatchMetrics {
					return podsWatch(podCfg, outputProcessor, outputProcessor)
				}
				return pods(podCfg, outputProcessor)
			},
			Flags: podsFlags(),
		},
	}
	app.Flags = commonFlags(&cfg)
	return app
}

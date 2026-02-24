package stdin

import (
	"io"

	metricsjson "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/json/metricsresources"
	metricsscreen "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/screen/metricsresources"
	metricstable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/metricsresources"
	metricstext "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/text/metricsresources"
	metricsyaml "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/yaml/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/columns"
	"github.com/trezorg/k8spodsmetrics/internal/config"
	"github.com/trezorg/k8spodsmetrics/internal/resources"

	nodesjson "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/json/noderesources"
	nodesscreen "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/screen/noderesources"
	nodestable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/noderesources"
	nodestext "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/text/noderesources"
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
	Columns      []string
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

func summaryOutputProcessor(out output.Output, res resources.Resources, cols []columns.Column) SummaryOutputProcessor {
	switch out {
	case output.Table:
		return nodestable.ToTable(res, cols)
	case output.JSON:
		return nodesjson.JSON(nodesjson.Print)
	case output.Yaml:
		return nodesyaml.Yaml(nodesyaml.Print)
	case output.Text:
		return nodestext.Text(nodestext.Print)
	}
	return nodestable.ToTable(res, cols)
}

func summaryWatchRenderer(out output.Output, res resources.Resources, cols []columns.Column) func(io.Writer, noderesources.NodeResourceList) {
	switch out {
	case output.Table:
		return nodestable.ToWriter(res, cols)
	case output.JSON:
		return nodesjson.PrintTo
	case output.Yaml:
		return nodesyaml.PrintTo
	case output.Text:
		return nodestext.PrintTo
	}
	return nodestable.ToWriter(res, cols)
}

func podsOutputProcessor(out output.Output, res resources.Resources, cols []columns.Column) PodsOutputProcessor {
	switch out {
	case output.Table:
		return metricstable.ToTable(res, cols)
	case output.JSON:
		return metricsjson.JSON(metricsjson.Print)
	case output.Yaml:
		return metricsyaml.Yaml(metricsyaml.Print)
	case output.Text:
		return metricstext.Text(metricstext.Print)
	}
	return metricstable.ToTable(res, cols)
}

func podsWatchRenderer(out output.Output, res resources.Resources, cols []columns.Column) func(io.Writer, metricsresources.PodMetricsResourceList) {
	switch out {
	case output.Table:
		return metricstable.ToWriter(res, cols)
	case output.JSON:
		return metricsjson.PrintTo
	case output.Yaml:
		return metricsyaml.PrintTo
	case output.Text:
		return metricstext.PrintTo
	}
	return metricstable.ToWriter(res, cols)
}

func parseColumnsForOutput(
	out output.Output,
	values []string,
	parse func([]string) []columns.Column,
	validate func([]columns.Column) error,
) ([]columns.Column, error) {
	if out != output.Table {
		return nil, nil
	}

	cols := parse(values)
	if err := validate(cols); err != nil {
		return nil, err
	}

	return cols, nil
}

func summary(processor SummaryProcessor, successProcessor noderesources.SuccessProcessor) error {
	return processor.Process(successProcessor)
}

func summaryWatch(
	processor SummaryWatcher,
	successRenderer func(io.Writer, noderesources.NodeResourceList),
	errorProcessor noderesources.ErrorProcessor,
) error {
	return processor.ProcessWatch(nodesscreen.NewScreenSuccessWriter(successRenderer), nodesscreen.NewScreenErrorWriter(errorProcessor))
}

func pods(processor PodsProcessor, successProcessor metricsresources.SuccessProcessor) error {
	return processor.Process(successProcessor)
}

func podsWatch(
	processor PodsWatcher,
	successRenderer func(io.Writer, metricsresources.PodMetricsResourceList),
	errorProcessor metricsresources.ErrorProcessor,
) error {
	return processor.ProcessWatch(metricsscreen.NewScreenSuccessWriter(successRenderer), metricsscreen.NewScreenErrorWriter(errorProcessor))
}

func loadFileConfig(configFile string) (*config.Config, error) {
	if configFile == "" {
		return nil, nil
	}
	return config.Load(configFile)
}

// applyCommonConfig merges file config with CLI common config values.
// CLI values take precedence over file config for string and numeric types.
func applyCommonConfig(cfg *commonConfig, fileConfig *config.Config, watchMetricsSet bool) config.Common {
	merged := config.Common{
		KubeConfig:   cfg.KubeConfig,
		KubeContext:  cfg.KubeContext,
		Output:       cfg.Output,
		Alert:        cfg.Alert,
		WatchPeriod:  cfg.WatchPeriod,
		WatchMetrics: cfg.WatchMetrics,
		Columns:      cfg.Columns,
	}
	if fileConfig != nil {
		fileConfig.MergeCommon(&merged)
	}
	if watchMetricsSet {
		merged.WatchMetrics = cfg.WatchMetrics
	}
	return merged
}

// applyPodsConfig merges file config with CLI pods command config values.
// CLI values take precedence over file config for string and slice types.
func applyPodsConfig(podCfg *podConfig, fileConfig *config.Config, reverseSet bool) config.Pods {
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
	if reverseSet {
		merged.Reverse = podCfg.Reverse
	}
	return merged
}

// applySummaryConfig merges file config with CLI summary command config values.
// CLI values take precedence over file config for string and slice types.
func applySummaryConfig(summaryCfg *summaryConfig, fileConfig *config.Config, reverseSet bool) config.Summary {
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
	if reverseSet {
		merged.Reverse = summaryCfg.Reverse
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
			Before:  loadConfigBefore(&cfg),
			Action: func(c *cli.Context) error {
				return runSummaryAction(c, cfg)
			},
			Flags: summaryFlags(),
		},
		{
			Name:    "pods",
			Aliases: []string{"p"},
			Before:  loadConfigBefore(&cfg),
			Action: func(c *cli.Context) error {
				return runPodsAction(c, cfg)
			},
			Flags: podsFlags(),
		},
	}
	app.Flags = commonFlags(&cfg)
	return app
}

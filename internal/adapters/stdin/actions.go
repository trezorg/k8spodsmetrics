package stdin

import (
	metricstable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/metricsresources"
	nodestable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/output"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
	"github.com/urfave/cli/v2"
)

type actionFlags struct {
	reverseSet bool
	watchSet   bool
	timeoutSet bool
	resources  []string
}

func loadConfigBefore(cfg *commonConfig) func(*cli.Context) error {
	return func(_ *cli.Context) error {
		var err error
		cfg.fileConfig, err = loadFileConfig(cfg.ConfigFile)
		return err
	}
}

func parseActionFlags(c *cli.Context) actionFlags {
	return actionFlags{
		reverseSet: c.IsSet("reverse"),
		watchSet:   c.IsSet("watch"),
		timeoutSet: c.IsSet("timeout"),
		resources:  c.StringSlice("resources"),
	}
}

func applyMergedCommonConfig(cfg *commonConfig, flags actionFlags) {
	mergedCommon := applyCommonConfig(cfg, cfg.fileConfig, flags.watchSet, flags.timeoutSet)
	cfg.KubeConfig = mergedCommon.KubeConfig
	cfg.KubeContext = mergedCommon.KubeContext
	cfg.Output = mergedCommon.Output
	cfg.Alert = mergedCommon.Alert
	cfg.WatchPeriod = mergedCommon.WatchPeriod
	cfg.WatchMetrics = mergedCommon.WatchMetrics
	cfg.Columns = mergedCommon.Columns
	cfg.Timeout = mergedCommon.Timeout
}

func mergedResources(cliResources []string, configResources []string) []string {
	if len(cliResources) == 0 && len(configResources) > 0 {
		return configResources
	}
	return cliResources
}

func mergeSummaryActionConfig(cfg *summaryConfig, flags actionFlags) {
	mergedSummary := applySummaryConfig(cfg, cfg.fileConfig, flags.reverseSet)
	cfg.Name = mergedSummary.Name
	cfg.Label = mergedSummary.Label
	cfg.Sorting = mergedSummary.Sorting
	cfg.Reverse = mergedSummary.Reverse
	applyMergedCommonConfig(&cfg.commonConfig, flags)
	cfg.Resources = mergedResources(flags.resources, mergedSummary.Resources)
}

func mergePodsActionConfig(cfg *podConfig, flags actionFlags) {
	mergedPods := applyPodsConfig(cfg, cfg.fileConfig, flags.reverseSet)
	cfg.Namespaces = mergedPods.Namespaces
	cfg.Label = mergedPods.Label
	cfg.FieldSelector = mergedPods.FieldSelector
	cfg.Nodes = mergedPods.Nodes
	cfg.Sorting = mergedPods.Sorting
	cfg.Reverse = mergedPods.Reverse
	applyMergedCommonConfig(&cfg.commonConfig, flags)
	cfg.Resources = mergedResources(flags.resources, mergedPods.Resources)
}

func runSummaryAction(c *cli.Context, cfg commonConfig) error {
	summaryActionConfig := summaryConfig{commonConfig: cfg}
	summaryActionConfig.Name = c.String("name")
	summaryActionConfig.Label = c.String("label")
	summaryActionConfig.Sorting = c.String("sorting")
	summaryActionConfig.Reverse = c.Bool("reverse")
	summaryActionConfig.Columns = c.StringSlice("columns")
	mergeSummaryActionConfig(&summaryActionConfig, parseActionFlags(c))

	if err := summaryActionConfig.Validate(); err != nil {
		return err
	}

	outputResources := resources.FromStrings(summaryActionConfig.Resources...)
	nodeCols, err := parseColumnsForOutput(
		output.Output(summaryActionConfig.Output),
		summaryActionConfig.Columns,
		nodestable.ParseColumns,
		nodestable.ValidateColumns,
	)
	if err != nil {
		return err
	}

	summaryCfg := nodeResourcesConfig(summaryActionConfig)
	outputProcessor := summaryOutputProcessor(output.Output(summaryActionConfig.Output), outputResources, nodeCols)
	if summaryActionConfig.WatchMetrics {
		return summaryWatch(&summaryCfg, summaryWatchRenderer(output.Output(summaryActionConfig.Output), outputResources, nodeCols), outputProcessor)
	}

	return summary(&summaryCfg, outputProcessor)
}

func runPodsAction(c *cli.Context, cfg commonConfig) error {
	podActionConfig := podConfig{commonConfig: cfg}
	podActionConfig.Namespaces = c.StringSlice("namespace")
	podActionConfig.Label = c.String("label")
	podActionConfig.FieldSelector = c.String("field-selector")
	podActionConfig.Sorting = c.String("sorting")
	podActionConfig.Reverse = c.Bool("reverse")
	podActionConfig.Nodes = c.StringSlice("node")
	podActionConfig.Columns = c.StringSlice("columns")
	mergePodsActionConfig(&podActionConfig, parseActionFlags(c))

	if err := podActionConfig.Validate(); err != nil {
		return err
	}

	outputResources := resources.FromStrings(podActionConfig.Resources...)
	podCols, err := parseColumnsForOutput(
		output.Output(podActionConfig.Output),
		podActionConfig.Columns,
		metricstable.ParseColumns,
		metricstable.ValidateColumns,
	)
	if err != nil {
		return err
	}

	podCfg := metricsResourcesConfig(podActionConfig)
	outputProcessor := podsOutputProcessor(output.Output(podActionConfig.Output), outputResources, podCols)
	if podActionConfig.WatchMetrics {
		return podsWatch(&podCfg, podsWatchRenderer(output.Output(podActionConfig.Output), outputResources, podCols), outputProcessor)
	}

	return pods(&podCfg, outputProcessor)
}

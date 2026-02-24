package stdin

import (
	metricstable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/metricsresources"
	nodestable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/output"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
	"github.com/urfave/cli/v2"
)

func loadConfigBefore(cfg *commonConfig) func(*cli.Context) error {
	return func(_ *cli.Context) error {
		var err error
		cfg.fileConfig, err = loadFileConfig(cfg.ConfigFile)
		return err
	}
}

func runSummaryAction(c *cli.Context, cfg commonConfig) error {
	summaryActionConfig := summaryConfig{commonConfig: cfg}
	summaryActionConfig.Name = c.String("name")
	summaryActionConfig.Label = c.String("label")
	summaryActionConfig.Sorting = c.String("sorting")
	summaryActionConfig.Reverse = c.Bool("reverse")
	summaryActionConfig.Columns = c.StringSlice("columns")
	cmdResources := c.StringSlice("resources")

	reverseSet := c.IsSet("reverse")
	watchSet := c.IsSet("watch")

	mergedSummary := applySummaryConfig(&summaryActionConfig, cfg.fileConfig, reverseSet)
	summaryActionConfig.Name = mergedSummary.Name
	summaryActionConfig.Label = mergedSummary.Label
	summaryActionConfig.Sorting = mergedSummary.Sorting
	summaryActionConfig.Reverse = mergedSummary.Reverse
	summaryActionConfig.Resources = mergedSummary.Resources

	mergedCommon := applyCommonConfig(&summaryActionConfig.commonConfig, cfg.fileConfig, watchSet)
	summaryActionConfig.KubeConfig = mergedCommon.KubeConfig
	summaryActionConfig.KubeContext = mergedCommon.KubeContext
	summaryActionConfig.Output = mergedCommon.Output
	summaryActionConfig.Alert = mergedCommon.Alert
	summaryActionConfig.WatchPeriod = mergedCommon.WatchPeriod
	summaryActionConfig.WatchMetrics = mergedCommon.WatchMetrics
	summaryActionConfig.Columns = mergedCommon.Columns

	if len(cmdResources) == 0 && len(mergedSummary.Resources) > 0 {
		cmdResources = mergedSummary.Resources
	}
	summaryActionConfig.Resources = cmdResources

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
	cmdResources := c.StringSlice("resources")

	reverseSet := c.IsSet("reverse")
	watchSet := c.IsSet("watch")

	mergedPods := applyPodsConfig(&podActionConfig, cfg.fileConfig, reverseSet)
	podActionConfig.Namespaces = mergedPods.Namespaces
	podActionConfig.Label = mergedPods.Label
	podActionConfig.FieldSelector = mergedPods.FieldSelector
	podActionConfig.Nodes = mergedPods.Nodes
	podActionConfig.Sorting = mergedPods.Sorting
	podActionConfig.Reverse = mergedPods.Reverse
	podActionConfig.Resources = mergedPods.Resources

	mergedCommon := applyCommonConfig(&podActionConfig.commonConfig, cfg.fileConfig, watchSet)
	podActionConfig.KubeConfig = mergedCommon.KubeConfig
	podActionConfig.KubeContext = mergedCommon.KubeContext
	podActionConfig.Output = mergedCommon.Output
	podActionConfig.Alert = mergedCommon.Alert
	podActionConfig.WatchPeriod = mergedCommon.WatchPeriod
	podActionConfig.WatchMetrics = mergedCommon.WatchMetrics
	podActionConfig.Columns = mergedCommon.Columns

	if len(cmdResources) == 0 && len(mergedPods.Resources) > 0 {
		cmdResources = mergedPods.Resources
	}
	podActionConfig.Resources = cmdResources

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

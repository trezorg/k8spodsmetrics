package stdin

import (
	metricstable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/metricsresources"
	nodestable "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/table/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/output"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
	"github.com/urfave/cli/v2"
)

type actionFlags struct {
	reverseSet     bool
	watchSet       bool
	watchPeriodSet bool
	timeoutSet     bool
	outputSet      bool
	alertSet       bool
	columnsSet     bool
	sortingSet     bool
	resourcesSet   bool
	resources      []string
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
		reverseSet:     c.IsSet("reverse"),
		watchSet:       c.IsSet("watch"),
		watchPeriodSet: c.IsSet("watch-period"),
		timeoutSet:     c.IsSet("timeout"),
		outputSet:      c.IsSet("output"),
		alertSet:       c.IsSet("alert"),
		columnsSet:     c.IsSet("columns"),
		sortingSet:     c.IsSet("sorting"),
		resourcesSet:   c.IsSet("resources"),
		resources:      c.StringSlice("resources"),
	}
}

func resolveCommonConfig(cfg commonConfig, flags actionFlags) commonConfig {
	mergeCandidate := cfg
	if !flags.outputSet {
		mergeCandidate.Output = ""
	}
	if !flags.alertSet {
		mergeCandidate.Alert = ""
	}
	if !flags.watchPeriodSet {
		mergeCandidate.WatchPeriod = 0
	}
	if !flags.columnsSet {
		mergeCandidate.Columns = nil
	}

	mergedCommon := applyCommonConfig(&mergeCandidate, cfg.fileConfig, flags.watchSet, flags.timeoutSet)
	if mergedCommon.Output == "" {
		mergedCommon.Output = string(output.Table)
	}
	if mergedCommon.Alert == "" {
		mergedCommon.Alert = "none"
	}
	if mergedCommon.WatchPeriod == 0 {
		mergedCommon.WatchPeriod = defaultWatchPeriodSeconds
	}

	return commonConfig{
		ConfigFile:   cfg.ConfigFile,
		KubeConfig:   mergedCommon.KubeConfig,
		KubeContext:  mergedCommon.KubeContext,
		Output:       mergedCommon.Output,
		Alert:        mergedCommon.Alert,
		WatchPeriod:  mergedCommon.WatchPeriod,
		WatchMetrics: mergedCommon.WatchMetrics,
		Columns:      mergedCommon.Columns,
		Timeout:      mergedCommon.Timeout,
		fileConfig:   cfg.fileConfig,
	}
}

func mergedResources(cliResources []string, configResources []string) []string {
	if len(cliResources) == 0 && len(configResources) > 0 {
		return configResources
	}
	if len(cliResources) == 0 {
		return []string{"all"}
	}
	return cliResources
}

func resolveSummaryActionConfig(c *cli.Context, cfg commonConfig) summaryConfig {
	flags := parseActionFlags(c)
	resolved := summaryConfig{
		Name:      c.String("name"),
		Label:     c.String("label"),
		Sorting:   c.String("sorting"),
		Reverse:   c.Bool("reverse"),
		Resources: flags.resources,
	}
	resolved.commonConfig = resolveCommonConfig(cfg, flags)
	if flags.columnsSet {
		resolved.Columns = c.StringSlice("columns")
	}
	if !flags.sortingSet {
		resolved.Sorting = ""
	}
	if !flags.resourcesSet {
		resolved.Resources = nil
	}

	mergedSummary := applySummaryConfig(&resolved, resolved.fileConfig, flags.reverseSet)
	resolved.Name = mergedSummary.Name
	resolved.Label = mergedSummary.Label
	resolved.Sorting = mergedSummary.Sorting
	resolved.Reverse = mergedSummary.Reverse
	if resolved.Sorting == "" {
		resolved.Sorting = "name"
	}
	resourcesFromCLI := []string(nil)
	if flags.resourcesSet {
		resourcesFromCLI = flags.resources
	}
	resolved.Resources = mergedResources(resourcesFromCLI, mergedSummary.Resources)

	return resolved
}

func resolvePodsActionConfig(c *cli.Context, cfg commonConfig) podConfig {
	flags := parseActionFlags(c)
	resolved := podConfig{
		Namespaces:    c.StringSlice("namespace"),
		Label:         c.String("label"),
		FieldSelector: c.String("field-selector"),
		Sorting:       c.String("sorting"),
		Reverse:       c.Bool("reverse"),
		Nodes:         c.StringSlice("node"),
		Resources:     flags.resources,
	}
	resolved.commonConfig = resolveCommonConfig(cfg, flags)
	if flags.columnsSet {
		resolved.Columns = c.StringSlice("columns")
	}
	if !flags.sortingSet {
		resolved.Sorting = ""
	}
	if !flags.resourcesSet {
		resolved.Resources = nil
	}

	mergedPods := applyPodsConfig(&resolved, resolved.fileConfig, flags.reverseSet)
	resolved.Namespaces = mergedPods.Namespaces
	resolved.Label = mergedPods.Label
	resolved.FieldSelector = mergedPods.FieldSelector
	resolved.Nodes = mergedPods.Nodes
	resolved.Sorting = mergedPods.Sorting
	resolved.Reverse = mergedPods.Reverse
	if resolved.Sorting == "" {
		resolved.Sorting = "namespace"
	}
	resourcesFromCLI := []string(nil)
	if flags.resourcesSet {
		resourcesFromCLI = flags.resources
	}
	resolved.Resources = mergedResources(resourcesFromCLI, mergedPods.Resources)

	return resolved
}

func runSummaryAction(c *cli.Context, cfg commonConfig) error {
	summaryActionConfig := resolveSummaryActionConfig(c, cfg)

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
	podActionConfig := resolvePodsActionConfig(c, cfg)

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

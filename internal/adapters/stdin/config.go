package stdin

import (
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

func metricsResourcesConfig(c podConfig) metricsresources.Config {
	return metricsresources.Config{
		KubeConfig:    c.KubeConfig,
		KubeContext:   c.KubeContext,
		Namespace:     c.Namespace,
		Label:         c.Label,
		FieldSelector: c.FieldSelector,
		Nodes:         c.Nodes,
		Output:        c.Output,
		Sorting:       c.Sorting,
		Reverse:       c.Reverse,
		Resources:     c.Resources,
		KLogLevel:     c.KLogLevel,
		Alert:         c.Alert,
		WatchMetrics:  c.WatchMetrics,
		WatchPeriod:   c.WatchPeriod,
	}
}

func nodeResourcesConfig(c summaryConfig) noderesources.Config {
	return noderesources.Config{
		KubeConfig:   c.KubeConfig,
		KubeContext:  c.KubeContext,
		Label:        c.Label,
		Name:         c.Name,
		Output:       c.Output,
		Sorting:      c.Sorting,
		Reverse:      c.Reverse,
		Resources:    c.Resources,
		KLogLevel:    c.KLogLevel,
		Alert:        c.Alert,
		WatchMetrics: c.WatchMetrics,
		WatchPeriod:  c.WatchPeriod,
	}
}

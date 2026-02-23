package stdin

import (
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

func metricsResourcesConfig(c podConfig) metricsresources.Config {
	return metricsresources.Config{
		KubeConfig:    c.KubeConfig,
		KubeContext:   c.KubeContext,
		Namespaces:    c.Namespaces,
		Label:         c.Label,
		FieldSelector: c.FieldSelector,
		Nodes:         c.Nodes,
		Sorting:       c.Sorting,
		Reverse:       c.Reverse,
		Alert:         c.Alert,
		WatchPeriod:   c.WatchPeriod,
	}
}

func nodeResourcesConfig(c summaryConfig) noderesources.Config {
	return noderesources.Config{
		KubeConfig:  c.KubeConfig,
		KubeContext: c.KubeContext,
		Label:       c.Label,
		Name:        c.Name,
		Sorting:     c.Sorting,
		Reverse:     c.Reverse,
		Alert:       c.Alert,
		WatchPeriod: c.WatchPeriod,
	}
}

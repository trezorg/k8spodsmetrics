package stdin

import (
	"errors"

	"github.com/trezorg/k8spodsmetrics/internal/alert"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/output"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
	metricssorting "github.com/trezorg/k8spodsmetrics/internal/sorting/metricsresources"
	nodesorting "github.com/trezorg/k8spodsmetrics/internal/sorting/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/tableview"
)

func (c *commonConfig) Validate() error {
	if c.WatchMetrics && c.WatchPeriod == 0 {
		return errors.New("watch period must be greater than 0")
	}
	if err := output.Valid(output.Output(c.Output)); err != nil {
		return err
	}
	view := tableview.View(c.TableView)
	if view == "" {
		view = tableview.Compact
	}
	if err := tableview.Valid(view); err != nil {
		return err
	}
	return alert.Valid(alert.Alert(c.Alert))
}

func (c *podConfig) Validate() error {
	if err := c.commonConfig.Validate(); err != nil {
		return err
	}
	if err := metricssorting.Valid(metricssorting.Sorting(c.Sorting)); err != nil {
		return err
	}
	outputResources := resources.FromStrings(c.Resources...)
	if err := resources.Valid(outputResources...); err != nil {
		return err
	}
	c.Resources = resources.ToStrings(outputResources...)
	return nil
}

func (c *summaryConfig) Validate() error {
	if err := c.commonConfig.Validate(); err != nil {
		return err
	}
	if err := nodesorting.Valid(nodesorting.Sorting(c.Sorting)); err != nil {
		return err
	}
	outputResources := resources.FromStrings(c.Resources...)
	if err := resources.Valid(outputResources...); err != nil {
		return err
	}
	c.Resources = resources.ToStrings(outputResources...)
	return nil
}

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
		Timeout:       c.Timeout,
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
		Timeout:     c.Timeout,
	}
}

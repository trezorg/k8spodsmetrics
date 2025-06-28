package metricsresources

import (
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

type String func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	_, _ = os.Stdout.WriteString(list.String() + "\n")
}

func (j String) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (String) Error(err error) {
	logger.Error("", err)
}

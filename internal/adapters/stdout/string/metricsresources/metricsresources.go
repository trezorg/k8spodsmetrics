package metricsresources

import (
	"fmt"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

type String func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	fmt.Println(list)
}

func (j String) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (j String) Error(err error) {
	logger.Error("", err)
}

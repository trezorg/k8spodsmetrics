package metricsresources

import (
	"fmt"
	"gopkg.in/yaml.v3"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

type Yaml func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	data, err := yaml.Marshal(list)
	if err != nil {
		logger.Error("", err)
		return
	}
	fmt.Println(string(data))
}

func (j Yaml) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (j Yaml) Error(err error) {
	logger.Error("", err)
}

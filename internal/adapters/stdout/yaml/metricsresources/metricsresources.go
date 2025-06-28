package metricsresources

import (
	"os"

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
	_, _ = os.Stdout.WriteString(string(data) + "\n")
}

func (j Yaml) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (Yaml) Error(err error) {
	logger.Error("", err)
}

package metricsresources

import (
	"encoding/json"
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

type JSON func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(list); err != nil {
		logger.Error("", err)
	}
}

func (j JSON) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (JSON) Error(err error) {
	logger.Error("", err)
}

package metricsresources

import (
	"encoding/json"
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

type Json func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(list); err != nil {
		logger.Error("", err)
	}
}

func (j Json) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (j Json) Error(err error) {
	logger.Error("", err)
}

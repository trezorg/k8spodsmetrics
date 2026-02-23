package metricsresources

import (
	stdoutcommon "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/common"
	"gopkg.in/yaml.v3"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"log/slog"
)

type Yaml func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	data, err := yaml.Marshal(list)
	if err != nil {
		slog.Error("failed to marshal metrics resources to yaml", "error", err)
		return
	}
	stdoutcommon.WriteStringLine(string(data))
}

func (j Yaml) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (Yaml) Error(err error) {
	slog.Error("yaml metrics resources output failed", "error", err)
}

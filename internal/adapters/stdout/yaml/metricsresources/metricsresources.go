package metricsresources

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"log/slog"
)

type Yaml func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	data, err := yaml.Marshal(list)
	if err != nil {
		slog.Error("", slog.Any("error", err))
		return
	}
	_, _ = os.Stdout.WriteString(string(data) + "\n")
}

func (j Yaml) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (Yaml) Error(err error) {
	slog.Error("", slog.Any("error", err))
}

package metricsresources

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"log/slog"
)

type Yaml func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	PrintTo(os.Stdout, list)
}

func PrintTo(w io.Writer, list metricsresources.PodMetricsResourceList) {
	data, err := yaml.Marshal(list)
	if err != nil {
		slog.Error("failed to marshal metrics resources to yaml", "error", err)
		return
	}
	_, _ = w.Write(data)
	_, _ = w.Write([]byte("\n"))
}

func (Yaml) SuccessTo(w io.Writer, list metricsresources.PodMetricsResourceList) {
	PrintTo(w, list)
}

func (j Yaml) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (Yaml) Error(err error) {
	slog.Error("yaml metrics resources output failed", "error", err)
}

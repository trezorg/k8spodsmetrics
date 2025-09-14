package metricsresources

import (
	"encoding/json"
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"log/slog"
)

type JSON func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(list); err != nil {
		slog.Error("", slog.Any("error", err))
	}
}

func (j JSON) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (JSON) Error(err error) {
	slog.Error("", slog.Any("error", err))
}

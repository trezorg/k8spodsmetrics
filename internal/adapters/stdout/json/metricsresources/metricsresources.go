package metricsresources

import (
	"encoding/json"
	"io"
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"log/slog"
)

type JSON func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	PrintTo(os.Stdout, list)
}

func PrintTo(w io.Writer, list metricsresources.PodMetricsResourceList) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(list); err != nil {
		slog.Error("failed to encode metrics resources as json", "error", err)
	}
}

func (JSON) SuccessTo(w io.Writer, list metricsresources.PodMetricsResourceList) {
	PrintTo(w, list)
}

func (j JSON) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (JSON) Error(err error) {
	slog.Error("json metrics resources output failed", "error", err)
}

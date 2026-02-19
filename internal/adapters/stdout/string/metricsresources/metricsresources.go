package metricsresources

import (
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"log/slog"
)

type String func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	_, _ = os.Stdout.WriteString(list.String() + "\n")
}

func (j String) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (String) Error(err error) {
	slog.Error("string metrics resources output failed", "error", err)
}

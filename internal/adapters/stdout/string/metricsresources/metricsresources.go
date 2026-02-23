package metricsresources

import (
	stdoutcommon "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/common"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"log/slog"
)

type String func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	stdoutcommon.WriteStringLine(list.String())
}

func (j String) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (String) Error(err error) {
	slog.Error("string metrics resources output failed", "error", err)
}

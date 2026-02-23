package metricsresources

import (
	stdoutcommon "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/common"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

type ScreenSuccessWriter func(list metricsresources.PodMetricsResourceList)
type ScreenErrorWriter func(err error)

func NewScreenSuccessWriter(writer metricsresources.SuccessProcessor) ScreenSuccessWriter {
	return ScreenSuccessWriter(stdoutcommon.WrapScreenSuccess(writer.Success))
}

func NewScreenErrorWriter(writer metricsresources.ErrorProcessor) ScreenErrorWriter {
	return ScreenErrorWriter(stdoutcommon.WrapScreenError(writer.Error))
}

func (s ScreenSuccessWriter) Success(list metricsresources.PodMetricsResourceList) {
	s(list)
}

func (s ScreenErrorWriter) Error(err error) {
	s(err)
}

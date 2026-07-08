package metricsresources

import (
	"io"

	"github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/screenutil"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

type ScreenSuccessWriter func(list metricsresources.PodMetricsResourceList)
type ScreenErrorWriter func(err error)

func NewScreenSuccessWriter(writer func(io.Writer, metricsresources.PodMetricsResourceList)) ScreenSuccessWriter {
	return ScreenSuccessWriter(screenutil.WrapScreenSuccess(writer))
}

func NewScreenErrorWriter(writer metricsresources.ErrorProcessor) ScreenErrorWriter {
	return ScreenErrorWriter(screenutil.WrapScreenError(writer.Error))
}

func (s ScreenSuccessWriter) Success(list metricsresources.PodMetricsResourceList) {
	s(list)
}

func (s ScreenErrorWriter) Error(err error) {
	s(err)
}

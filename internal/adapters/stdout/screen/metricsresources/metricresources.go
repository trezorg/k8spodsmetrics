package metricsresources

import (
	screen "github.com/aditya43/clear-shell-screen-golang"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

type ScreenSuccessWriter func(list metricsresources.PodMetricsResourceList)
type ScreenErrorWriter func(err error)

func NewScreenSuccessWriter(writer metricsresources.SuccessProcessor) ScreenSuccessWriter {
	return func(list metricsresources.PodMetricsResourceList) {
		screen.Clear()
		screen.MoveTopLeft()
		writer.Success(list)
		screen.MoveTopLeft()
	}
}

func NewScreenErrorWriter(writer metricsresources.ErrorProcessor) ScreenErrorWriter {
	return func(err error) {
		screen.Clear()
		screen.MoveTopLeft()
		writer.Error(err)
		screen.MoveTopLeft()
	}
}

func (s ScreenSuccessWriter) Success(list metricsresources.PodMetricsResourceList) {
	s(list)
}

func (s ScreenErrorWriter) Error(err error) {
	s(err)
}

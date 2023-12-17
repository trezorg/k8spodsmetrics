package noderesources

import (
	screen "github.com/aditya43/clear-shell-screen-golang"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type ScreenSuccessWriter func(list noderesources.NodeResourceList)
type ScreenErrorWriter func(err error)

func NewScreenSuccessWriter(writer noderesources.SuccessProcessor) ScreenSuccessWriter {
	return func(list noderesources.NodeResourceList) {
		screen.Clear()
		screen.MoveTopLeft()
		writer.Success(list)
		screen.MoveTopLeft()
	}
}

func NewScreenErrorWriter(writer noderesources.ErrorProcessor) ScreenErrorWriter {
	return func(err error) {
		screen.Clear()
		screen.MoveTopLeft()
		writer.Error(err)
		screen.MoveTopLeft()
	}
}

func (s ScreenSuccessWriter) Success(list noderesources.NodeResourceList) {
	s(list)
}

func (s ScreenErrorWriter) Error(err error) {
	s(err)
}

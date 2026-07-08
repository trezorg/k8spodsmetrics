package noderesources

import (
	"io"

	"github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/screenutil"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type ScreenSuccessWriter func(list noderesources.NodeResourceList)
type ScreenErrorWriter func(err error)

func NewScreenSuccessWriter(writer func(io.Writer, noderesources.NodeResourceList)) ScreenSuccessWriter {
	return ScreenSuccessWriter(screenutil.WrapScreenSuccess(writer))
}

func NewScreenErrorWriter(writer noderesources.ErrorProcessor) ScreenErrorWriter {
	return ScreenErrorWriter(screenutil.WrapScreenError(writer.Error))
}

func (s ScreenSuccessWriter) Success(list noderesources.NodeResourceList) {
	s(list)
}

func (s ScreenErrorWriter) Error(err error) {
	s(err)
}

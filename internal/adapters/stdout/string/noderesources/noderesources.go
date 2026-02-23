package noderesources

import (
	stdoutcommon "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/common"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"log/slog"
)

type String func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	stdoutcommon.WriteStringLine(list.String())
}

func (j String) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (String) Error(err error) {
	slog.Error("string node resources output failed", "error", err)
}

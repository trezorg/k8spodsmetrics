package noderesources

import (
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type String func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	_, _ = os.Stdout.WriteString(list.String() + "\n")
}

func (j String) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (String) Error(err error) {
	logger.Error("", err)
}

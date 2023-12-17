package noderesources

import (
	"fmt"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type String func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	fmt.Println(list)
}

func (j String) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (j String) Error(err error) {
	logger.Error("", err)
}

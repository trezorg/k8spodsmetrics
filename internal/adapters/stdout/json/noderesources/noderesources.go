package noderesources

import (
	"encoding/json"
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type Json func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	envelop := noderesources.NodeResourceListEnvelop{Items: list}
	if err := enc.Encode(envelop); err != nil {
		logger.Error("", err)
	}
}

func (j Json) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (j Json) Error(err error) {
	logger.Error("", err)
}

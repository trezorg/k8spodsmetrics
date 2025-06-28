package noderesources

import (
	"encoding/json"
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type JSON func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	envelop := noderesources.NodeResourceListEnvelop{Items: list}
	if err := enc.Encode(envelop); err != nil {
		logger.Error("", err)
	}
}

func (j JSON) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (JSON) Error(err error) {
	logger.Error("", err)
}

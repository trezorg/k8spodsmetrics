package noderesources

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type Yaml func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	enc := yaml.NewEncoder(os.Stdout)
	envelop := noderesources.NodeResourceListEnvelop{Items: list}
	if err := enc.Encode(envelop); err != nil {
		logger.Error("", err)
	}
}

func (j Yaml) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (j Yaml) Error(err error) {
	logger.Error("", err)
}

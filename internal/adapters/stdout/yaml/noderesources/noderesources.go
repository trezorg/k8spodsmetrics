package noderesources

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"log/slog"
)

type Yaml func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	enc := yaml.NewEncoder(os.Stdout)
	envelop := noderesources.NodeResourceListEnvelop{Items: list}
	if err := enc.Encode(envelop); err != nil {
		slog.Error("failed to encode node resources as yaml", "error", err)
	}
}

func (j Yaml) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (Yaml) Error(err error) {
	slog.Error("yaml node resources output failed", "error", err)
}

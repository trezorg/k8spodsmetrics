package noderesources

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"log/slog"
)

type Yaml func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	PrintTo(os.Stdout, list)
}

func PrintTo(w io.Writer, list noderesources.NodeResourceList) {
	enc := yaml.NewEncoder(w)
	defer func() { _ = enc.Close() }()
	envelope := noderesources.NodeResourceListEnvelope{Items: list}
	if err := enc.Encode(envelope); err != nil {
		slog.Error("failed to encode node resources as yaml", "error", err)
	}
}

func (Yaml) SuccessTo(w io.Writer, list noderesources.NodeResourceList) {
	PrintTo(w, list)
}

func (j Yaml) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (Yaml) Error(err error) {
	slog.Error("yaml node resources output failed", "error", err)
}

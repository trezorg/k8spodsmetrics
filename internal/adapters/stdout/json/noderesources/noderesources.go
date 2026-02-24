package noderesources

import (
	"encoding/json"
	"io"
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"log/slog"
)

type JSON func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	PrintTo(os.Stdout, list)
}

func PrintTo(w io.Writer, list noderesources.NodeResourceList) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	envelope := noderesources.NodeResourceListEnvelope{Items: list}
	if err := enc.Encode(envelope); err != nil {
		slog.Error("failed to encode node resources as json", "error", err)
	}
}

func (JSON) SuccessTo(w io.Writer, list noderesources.NodeResourceList) {
	PrintTo(w, list)
}

func (j JSON) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (JSON) Error(err error) {
	slog.Error("json node resources output failed", "error", err)
}

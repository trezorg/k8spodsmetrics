package noderesources

import (
	"encoding/json"
	"os"

	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"log/slog"
)

type JSON func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	envelop := noderesources.NodeResourceListEnvelop{Items: list}
	if err := enc.Encode(envelop); err != nil {
		slog.Error("", slog.Any("error", err))
	}
}

func (j JSON) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (JSON) Error(err error) {
	slog.Error("", slog.Any("error", err))
}

package noderesources

import (
	"bytes"
	"fmt"
	"log/slog"

	stdoutcommon "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/common"
	formatnoderesources "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/format/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type Text func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	var buffer bytes.Buffer
	for _, node := range list {
		formatter := formatnoderesources.New(node)
		_, _ = fmt.Fprintf(&buffer, "Name: %s\n", node.Name)
		_, _ = fmt.Fprintf(&buffer, "Memory: %s\n", formatter.MemoryTemplate())
		_, _ = fmt.Fprintf(&buffer, "CPU: %s\n", formatter.CPUTemplate())
	}
	stdoutcommon.WriteStringLine(buffer.String())
}

func (j Text) Success(list noderesources.NodeResourceList) {
	j(list)
}

func (Text) Error(err error) {
	slog.Error("text node resources output failed", "error", err)
}

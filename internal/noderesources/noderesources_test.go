package noderesources

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
)

func nodeResourceList(name string) NodeResourceList {
	return []NodeResource{
		{
			Name:              name,
			CPU:               1024,
			Memory:            1024,
			PodsCPURequest:    512,
			PodsMemoryRequest: 512,
			PodsCPULimit:      512,
			PodsMemoryLimit:   512,
		},
	}
}

func TestStringify(t *testing.T) {
	logger.InitDefaultLogger()
	nodeResourceList := nodeResourceList("foo")
	text := nodeResourceList.String()
	require.Greater(t, len(text), 0)
	require.NotContains(t, text, "/", text)
}

package noderesources

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		config := Config{}
		require.Empty(t, config.KubeConfig)
		require.Empty(t, config.Name)
	})

	t.Run("with values", func(t *testing.T) {
		config := Config{
			KubeConfig:  "/path/to/config",
			KubeContext: "test-context",
			Name:        "node1",
			Label:       "node-role.kubernetes.io/worker",
			Sorting:     "name",
			Reverse:     true,
			Alert:       "memory",
			WatchPeriod: 10,
		}
		require.Equal(t, "/path/to/config", config.KubeConfig)
		require.Equal(t, "test-context", config.KubeContext)
		require.Equal(t, "node1", config.Name)
	})
}

func TestWatchResponse(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		resp := WatchResponse{Error: errors.New("test error")}
		require.Error(t, resp.Error)
		require.Empty(t, resp.Data)
	})

	t.Run("with data", func(t *testing.T) {
		data := NodeResourceList{{Name: "node1"}}
		resp := WatchResponse{Data: data}
		require.NoError(t, resp.Error)
		require.Len(t, resp.Data, 1)
	})
}

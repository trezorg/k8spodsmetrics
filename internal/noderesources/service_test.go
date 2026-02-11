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
		require.False(t, config.WatchMetrics)
	})

	t.Run("with values", func(t *testing.T) {
		config := Config{
			KubeConfig:   "/path/to/config",
			KubeContext:  "test-context",
			Name:         "node1",
			Label:        "node-role.kubernetes.io/worker",
			Output:       "json",
			Sorting:      "name",
			Reverse:      true,
			Resources:    []string{"cpu", "memory"},
			KLogLevel:    5,
			Alert:        "memory",
			WatchPeriod:  10,
			WatchMetrics: true,
		}
		require.Equal(t, "/path/to/config", config.KubeConfig)
		require.Equal(t, "test-context", config.KubeContext)
		require.Equal(t, "node1", config.Name)
		require.True(t, config.WatchMetrics)
	})
}

func TestWatchResponse(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		resp := WatchResponse{error: errors.New("test error")}
		require.Error(t, resp.error)
		require.Empty(t, resp.data)
	})

	t.Run("with data", func(t *testing.T) {
		data := NodeResourceList{{Name: "node1"}}
		resp := WatchResponse{data: data}
		require.NoError(t, resp.error)
		require.Len(t, resp.data, 1)
	})
}

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

func TestConfigValidate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := Config{Sorting: "name", Alert: "none"}
		require.NoError(t, cfg.Validate())
	})

	t.Run("invalid sorting", func(t *testing.T) {
		cfg := Config{Sorting: "invalid", Alert: "none"}
		require.ErrorContains(t, cfg.Validate(), "sorting should be one of")
	})

	t.Run("invalid alert", func(t *testing.T) {
		cfg := Config{Sorting: "name", Alert: "invalid"}
		require.ErrorContains(t, cfg.Validate(), "alert should be one of")
	})
}

func TestConfigValidateWatch(t *testing.T) {
	t.Run("zero watch period", func(t *testing.T) {
		cfg := Config{Sorting: "name", Alert: "none", WatchPeriod: 0}
		require.ErrorContains(t, cfg.ValidateWatch(), "watch period must be greater than 0")
	})
}

type noopSuccessProcessor struct{}

func (noopSuccessProcessor) Success(NodeResourceList) {}

type noopErrorProcessor struct{}

func (noopErrorProcessor) Error(error) {}

func TestProcessValidationError(t *testing.T) {
	cfg := Config{KubeConfig: "dummy", Sorting: "invalid", Alert: "none"}

	err := cfg.Process(noopSuccessProcessor{})
	require.ErrorContains(t, err, "sorting should be one of")
}

func TestProcessWatchValidationError(t *testing.T) {
	cfg := Config{KubeConfig: "dummy", Sorting: "name", Alert: "none", WatchPeriod: 0}

	err := cfg.ProcessWatch(noopSuccessProcessor{}, noopErrorProcessor{})
	require.ErrorContains(t, err, "watch period must be greater than 0")
}

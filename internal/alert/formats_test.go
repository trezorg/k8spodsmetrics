package alert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	t.Run("valid alerts", func(t *testing.T) {
		validAlerts := []Alert{
			Any, Memory, MemoryRequest, MemoryLimit,
			CPU, CPURequest, CPULimit,
			Storage, StorageEphemeral, None,
		}
		for _, alert := range validAlerts {
			err := Valid(alert)
			require.NoError(t, err)
		}
	})

	t.Run("invalid alert", func(t *testing.T) {
		err := Valid(Alert("invalid"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "alert should be one of")
	})
}

func TestStringList(t *testing.T) {
	t.Run("default separator", func(t *testing.T) {
		list := StringListDefault()
		require.NotEmpty(t, list)
		require.Contains(t, list, "|")
	})

	t.Run("custom separator", func(t *testing.T) {
		list := StringList(",")
		require.NotEmpty(t, list)
		require.Contains(t, list, ",")
	})

	t.Run("all alerts included", func(t *testing.T) {
		list := StringListDefault()
		expectedAlerts := []string{"any", "memory", "memory_request", "memory_limit", "cpu", "cpu_request", "cpu_limit", "storage", "storage_ephemeral", "none"}
		for _, alert := range expectedAlerts {
			require.Contains(t, list, alert)
		}
	})
}

func TestAlertString(t *testing.T) {
	t.Run("any alert", func(t *testing.T) {
		require.Equal(t, "any", string(Any))
	})

	t.Run("memory alert", func(t *testing.T) {
		require.Equal(t, "memory", string(Memory))
	})

	t.Run("cpu alert", func(t *testing.T) {
		require.Equal(t, "cpu", string(CPU))
	})

	t.Run("none alert", func(t *testing.T) {
		require.Equal(t, "none", string(None))
	})
}

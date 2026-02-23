package noderesources

import (
	"testing"

	escapes "github.com/snugfox/ansi-escapes"
	"github.com/stretchr/testify/require"
	servicenoderesources "github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

func TestFormatterMemoryTemplate(t *testing.T) {
	resource := servicenoderesources.NodeResource{
		Memory:        1024,
		UsedMemory:    256,
		MemoryRequest: 2048,
		MemoryLimit:   4096,
	}

	formatted := New(resource).MemoryTemplate()
	require.Contains(t, formatted, "Node=1KiB/256B")
	require.Contains(t, formatted, "Requests=")
	require.Contains(t, formatted, "Limits=")
	require.Contains(t, formatted, escapes.TextColorYellow)
	require.Contains(t, formatted, escapes.TextColorRed)
}

func TestFormatterStorageUsedString(t *testing.T) {
	t.Run("non alerted", func(t *testing.T) {
		resource := servicenoderesources.NodeResource{Storage: 100, UsedStorage: 80}
		formatted := New(resource).StorageUsedString()
		require.Equal(t, "80B", formatted)
	})

	t.Run("alerted", func(t *testing.T) {
		resource := servicenoderesources.NodeResource{Storage: 100, UsedStorage: 96}
		formatted := New(resource).StorageUsedString()
		require.Contains(t, formatted, escapes.TextColorRed)
		require.Contains(t, formatted, escapes.ColorReset)
	})
}

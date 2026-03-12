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

func TestFormatterCompactCapacityStrings(t *testing.T) {
	resource := servicenoderesources.NodeResource{
		AllocatableCPU:              3900,
		UsedCPU:                     1200,
		FreeCPU:                     2700,
		AllocatableMemory:           8 * 1024,
		UsedMemory:                  3 * 1024,
		FreeMemory:                  5 * 1024,
		AllocatableStorage:          10 * 1024,
		UsedStorage:                 4 * 1024,
		FreeStorage:                 6 * 1024,
		AllocatableStorageEphemeral: 20 * 1024,
		UsedStorageEphemeral:        7 * 1024,
		FreeStorageEphemeral:        13 * 1024,
	}

	formatter := New(resource)
	require.Equal(t, "3900/1200/2700", formatter.CPUCapacityCompactString())
	require.Equal(t, "8KiB/3KiB/5KiB", formatter.MemoryCapacityCompactString())
	require.Equal(t, "10KiB/4KiB/6KiB", formatter.StorageCapacityCompactString())
	require.Equal(t, "20KiB/7KiB/13KiB", formatter.StorageEphemeralCapacityCompactString())
}

func TestFormatterCompactDemandStrings(t *testing.T) {
	resource := servicenoderesources.NodeResource{
		CPU:           8000,
		CPURequest:    2200,
		CPULimit:      6000,
		Memory:        32 * 1024,
		MemoryRequest: 8 * 1024,
		MemoryLimit:   16 * 1024,
	}

	formatter := New(resource)
	require.Equal(t, "2200/6000", formatter.CPUDemandCompactString())
	require.Equal(t, "8KiB/16KiB", formatter.MemoryDemandCompactString())
}

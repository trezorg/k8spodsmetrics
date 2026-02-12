package noderesources

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
)

func TestHeaderFooter(t *testing.T) {
	t.Run("with CPU only", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU}
		result := headerFooter(outputResources, "Test")
		require.Len(t, result, 8)
		require.Equal(t, "Test", result[0])
		require.Equal(t, "CPU", result[1])
	})

	t.Run("with Memory only", func(t *testing.T) {
		outputResources := resources.Resources{resources.Memory}
		result := headerFooter(outputResources, "Test")
		require.Len(t, result, 8)
		require.Equal(t, "Test", result[0])
		require.Equal(t, "Memory", result[1])
	})

	t.Run("with all resources", func(t *testing.T) {
		outputResources := resources.Resources{resources.All}
		result := headerFooter(outputResources, "Test")
		require.Len(t, result, 23)
		require.Equal(t, "Test", result[0])
	})
}

func TestSecondaryHeader(t *testing.T) {
	t.Run("with CPU only", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU}
		result := secondaryHeader(outputResources)
		require.Len(t, result, 8)
		require.Equal(t, "Total", result[1])
		require.Equal(t, "Allocatable", result[2])
		require.Equal(t, "Used", result[3])
	})

	t.Run("with Memory only", func(t *testing.T) {
		outputResources := resources.Resources{resources.Memory}
		result := secondaryHeader(outputResources)
		require.Len(t, result, 8)
		require.Equal(t, "Total", result[1])
		require.Equal(t, "Allocatable", result[2])
		require.Equal(t, "Used", result[3])
	})
}

func TestRow(t *testing.T) {
	t.Run("with CPU only", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU}
		resource := noderesources.NodeResource{
			Name:           "node-1",
			CPU:            4000,
			AllocatableCPU: 3900,
			UsedCPU:        100,
		}
		result := row(resource, outputResources)
		require.Len(t, result, 8)
		require.Equal(t, "node-1", result[0])
	})

	t.Run("with Memory only", func(t *testing.T) {
		outputResources := resources.Resources{resources.Memory}
		resource := noderesources.NodeResource{
			Name:              "node-1",
			Memory:            8 * 1024 * 1024 * 1024,
			AllocatableMemory: 7 * 1024 * 1024 * 1024,
			UsedMemory:        1024 * 1024 * 1024,
		}
		result := row(resource, outputResources)
		require.Len(t, result, 8)
		require.Equal(t, "node-1", result[0])
	})
}

func TestToTable(t *testing.T) {
	t.Run("creates table function", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU, resources.Memory}
		tableFunc := ToTable(outputResources)
		require.NotNil(t, tableFunc)

		list := noderesources.NodeResourceList{}
		require.NotPanics(t, func() {
			tableFunc(list)
		})
	})
}

func TestTable_Success(t *testing.T) {
	t.Run("calls table function", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU, resources.Memory}
		tableFunc := ToTable(outputResources)
		require.NotNil(t, tableFunc)

		list := noderesources.NodeResourceList{}
		require.NotPanics(t, func() {
			tableFunc.Success(list)
		})
	})
}

func TestTable_Error(t *testing.T) {
	t.Run("logs error without panicking", func(t *testing.T) {
		tableFunc := ToTable(resources.Resources{})
		require.NotPanics(t, func() {
			tableFunc.Error(nil)
		})
	})
}

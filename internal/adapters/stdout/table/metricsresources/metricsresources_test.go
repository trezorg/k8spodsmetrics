package metricsresources

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trezorg/k8spodsmetrics/internal/columns"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func TestHeaderFooter(t *testing.T) {
	t.Run("with CPU only", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU}
		result := headerFooter(outputResources, "Test")
		require.Len(t, result, 6)
		require.Equal(t, "Test", result[0])
		require.Equal(t, "CPU", result[3])
		require.Equal(t, "CPU", result[4])
		require.Equal(t, "CPU", result[5])
	})

	t.Run("with Memory only", func(t *testing.T) {
		outputResources := resources.Resources{resources.Memory}
		result := headerFooter(outputResources, "Test")
		require.Len(t, result, 6)
		require.Equal(t, "Test", result[0])
		require.Equal(t, "Memory", result[3])
		require.Equal(t, "Memory", result[4])
		require.Equal(t, "Memory", result[5])
	})

	t.Run("with all resources", func(t *testing.T) {
		outputResources := resources.Resources{resources.All}
		result := headerFooter(outputResources, "Test")
		require.Len(t, result, 15)
		require.Equal(t, "Test", result[0])
	})
}

func TestSecondaryHeader(t *testing.T) {
	t.Run("with CPU only", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU}
		result := secondaryHeader(outputResources)
		require.Len(t, result, 6)
		require.Equal(t, "Request", result[3])
		require.Equal(t, "Limit", result[4])
		require.Equal(t, "Used", result[5])
	})

	t.Run("with Memory only", func(t *testing.T) {
		outputResources := resources.Resources{resources.Memory}
		result := secondaryHeader(outputResources)
		require.Len(t, result, 6)
		require.Equal(t, "Request", result[3])
		require.Equal(t, "Limit", result[4])
		require.Equal(t, "Used", result[5])
	})
}

func TestContainerRow(t *testing.T) {
	t.Run("with CPU only", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU}
		container := metricsresources.ContainerMetricsResource{
			Name: "container-1",
			Requests: metricsresources.MetricsResource{
				CPURequest: 1000,
			},
			Limits: metricsresources.MetricsResource{
				CPURequest: 2000,
				CPUUsed:    1500,
			},
		}
		result := containerRow(container, outputResources)
		require.Len(t, result, 6)
		require.Equal(t, "container-1", result[0])
	})

	t.Run("with Memory only", func(t *testing.T) {
		outputResources := resources.Resources{resources.Memory}
		container := metricsresources.ContainerMetricsResource{
			Name: "container-1",
			Requests: metricsresources.MetricsResource{
				MemoryRequest: 1024 * 1024 * 100,
			},
			Limits: metricsresources.MetricsResource{
				MemoryRequest: 1024 * 1024 * 200,
				MemoryUsed:    1024 * 1024 * 150,
			},
		}
		result := containerRow(container, outputResources)
		require.Len(t, result, 6)
		require.Equal(t, "container-1", result[0])
	})
}

func TestRow(t *testing.T) {
	t.Run("with CPU only", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU}
		resource := metricsresources.PodMetricsResource{
			PodResource: pods.PodResource{
				NamespaceName: pods.NamespaceName{
					Name:      "test-pod",
					Namespace: "default",
				},
				NodeName: "node-1",
				Containers: []pods.ContainerResource{
					{
						Name: "container-1",
						Requests: pods.Resource{
							CPU: 1000,
						},
						Limits: pods.Resource{
							CPU: 2000,
						},
					},
				},
			},
			PodMetric: podmetrics.PodMetric{
				Name:      "test-pod",
				Namespace: "default",
				Containers: []podmetrics.ContainerMetric{
					{
						Name: "container-1",
						Metric: podmetrics.Metric{
							CPU: 1500,
						},
					},
				},
			},
		}
		result := row(resource, outputResources)
		require.Len(t, result, 6)
		require.Equal(t, "test-pod", result[0])
		require.Equal(t, "default", result[1])
		require.Equal(t, "node-1", result[2])
	})
}

func TestToTable(t *testing.T) {
	t.Run("creates table function with empty columns", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU, resources.Memory}
		tableFunc := ToTable(outputResources, nil)
		require.NotNil(t, tableFunc)

		list := metricsresources.PodMetricsResourceList{}
		require.NotPanics(t, func() {
			tableFunc(list)
		})
	})

	t.Run("creates table function with filtered columns", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU, resources.Memory}
		tableFunc := ToTable(outputResources, []columns.Column{columns.Used})
		require.NotNil(t, tableFunc)

		list := metricsresources.PodMetricsResourceList{}
		require.NotPanics(t, func() {
			tableFunc(list)
		})
	})
}

func TestTable_Success(t *testing.T) {
	t.Run("calls table function", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU, resources.Memory}
		tableFunc := ToTable(outputResources, nil)
		require.NotNil(t, tableFunc)

		list := metricsresources.PodMetricsResourceList{}
		require.NotPanics(t, func() {
			tableFunc.Success(list)
		})
	})
}

func TestTable_Error(t *testing.T) {
	t.Run("logs error without panicking", func(t *testing.T) {
		tableFunc := ToTable(resources.Resources{}, nil)
		require.NotPanics(t, func() {
			tableFunc.Error(nil)
		})
	})
}

func TestColumnSet(t *testing.T) {
	t.Run("empty columns shows all", func(t *testing.T) {
		cs := newColumnSet(nil)
		require.True(t, cs.Request)
		require.True(t, cs.Limit)
		require.True(t, cs.Used)
	})

	t.Run("selected columns only", func(t *testing.T) {
		cs := newColumnSet([]columns.Column{columns.Used})
		require.False(t, cs.Request)
		require.False(t, cs.Limit)
		require.True(t, cs.Used)
	})
}

func TestParseColumns(t *testing.T) {
	t.Run("removes duplicates", func(t *testing.T) {
		result := ParseColumns([]string{"used", "request", "used"})
		require.Len(t, result, 2)
		require.Equal(t, columns.Used, result[0])
		require.Equal(t, columns.Request, result[1])
	})
}

func TestValidateColumns(t *testing.T) {
	t.Run("valid columns for pods", func(t *testing.T) {
		err := ValidateColumns([]columns.Column{columns.Request, columns.Used})
		require.NoError(t, err)
	})

	t.Run("invalid column for pods - total not allowed", func(t *testing.T) {
		err := ValidateColumns([]columns.Column{columns.Total})
		require.Error(t, err)
	})

	t.Run("invalid column for pods - available not allowed", func(t *testing.T) {
		err := ValidateColumns([]columns.Column{columns.Available})
		require.Error(t, err)
	})
}

func TestColumnSetTotalRowStorageColumns(t *testing.T) {
	t.Run("storage request and limit totals use request fields", func(t *testing.T) {
		cs := ColumnSet{Request: true, Limit: true}
		outputResources := resources.Resources{resources.Storage}
		total := metricsresources.ContainerMetricsResource{
			Requests: metricsresources.MetricsResource{
				StorageRequest:          1024,
				StorageEphemeralRequest: 4096,
				StorageUsed:             123,
				StorageEphemeralUsed:    456,
			},
			Limits: metricsresources.MetricsResource{
				StorageRequest:          2048,
				StorageEphemeralRequest: 8192,
				StorageUsed:             789,
				StorageEphemeralUsed:    999,
			},
		}

		row := cs.totalRow(outputResources, total)
		require.Len(t, row, 7)
		require.Equal(t, "1KiB", row[3])
		require.Equal(t, "2KiB", row[4])
		require.Equal(t, "4KiB", row[5])
		require.Equal(t, "8KiB", row[6])
	})

	t.Run("storage used totals use used fields", func(t *testing.T) {
		cs := ColumnSet{Used: true}
		outputResources := resources.Resources{resources.Storage}
		total := metricsresources.ContainerMetricsResource{
			Requests: metricsresources.MetricsResource{
				StorageUsed:          3072,
				StorageEphemeralUsed: 5120,
			},
		}

		row := cs.totalRow(outputResources, total)
		require.Len(t, row, 5)
		require.Equal(t, "3KiB", row[3])
		require.Equal(t, "5KiB", row[4])
	})
}

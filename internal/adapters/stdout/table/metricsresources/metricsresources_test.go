package metricsresources

import (
	"testing"

	"github.com/stretchr/testify/require"

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
	t.Run("creates table function", func(t *testing.T) {
		outputResources := resources.Resources{resources.CPU, resources.Memory}
		tableFunc := ToTable(outputResources)
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
		tableFunc := ToTable(outputResources)
		require.NotNil(t, tableFunc)

		list := metricsresources.PodMetricsResourceList{}
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

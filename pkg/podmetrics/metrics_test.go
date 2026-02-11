package podmetrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetricFilter(t *testing.T) {
	t.Run("default filter", func(t *testing.T) {
		filter := MetricFilter{}
		require.Empty(t, filter.Namespace)
		require.Empty(t, filter.LabelSelector)
		require.Empty(t, filter.FieldSelector)
	})

	t.Run("with namespace", func(t *testing.T) {
		filter := MetricFilter{Namespace: "test-ns"}
		require.Equal(t, "test-ns", filter.Namespace)
	})

	t.Run("with label selector", func(t *testing.T) {
		filter := MetricFilter{LabelSelector: "app=test"}
		require.Equal(t, "app=test", filter.LabelSelector)
	})

	t.Run("with field selector", func(t *testing.T) {
		filter := MetricFilter{FieldSelector: "spec.nodeName=node1"}
		require.Equal(t, "spec.nodeName=node1", filter.FieldSelector)
	})
}

func TestMetric(t *testing.T) {
	t.Run("default metric", func(t *testing.T) {
		metric := Metric{}
		require.Equal(t, int64(0), metric.CPU)
		require.Equal(t, int64(0), metric.Memory)
	})

	t.Run("with values", func(t *testing.T) {
		metric := Metric{
			CPU:              1000,
			Memory:           1024 * 1024,
			Storage:          2048 * 1024,
			StorageEphemeral: 512 * 1024,
		}
		require.Equal(t, int64(1000), metric.CPU)
		require.Equal(t, int64(1024*1024), metric.Memory)
		require.Equal(t, int64(2048*1024), metric.Storage)
		require.Equal(t, int64(512*1024), metric.StorageEphemeral)
	})
}

func TestContainerMetric(t *testing.T) {
	t.Run("default metric", func(t *testing.T) {
		metric := ContainerMetric{}
		require.Empty(t, metric.Name)
		require.Equal(t, int64(0), metric.CPU)
	})

	t.Run("with values", func(t *testing.T) {
		metric := ContainerMetric{
			Name:   "container1",
			Metric: Metric{CPU: 1000, Memory: 1024 * 1024},
		}
		require.Equal(t, "container1", metric.Name)
		require.Equal(t, int64(1000), metric.CPU)
		require.Equal(t, int64(1024*1024), metric.Memory)
	})
}

func TestPodMetric(t *testing.T) {
	t.Run("default metric", func(t *testing.T) {
		metric := PodMetric{}
		require.Empty(t, metric.Namespace)
		require.Empty(t, metric.Name)
		require.Empty(t, metric.Containers)
	})

	t.Run("with values", func(t *testing.T) {
		metric := PodMetric{
			Namespace: "test-ns",
			Name:      "test-pod",
			Containers: []ContainerMetric{
				{Name: "c1", Metric: Metric{CPU: 100, Memory: 1024}},
			},
		}
		require.Equal(t, "test-ns", metric.Namespace)
		require.Equal(t, "test-pod", metric.Name)
		require.Len(t, metric.Containers, 1)
	})
}

func TestPodMetricList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		list := PodMetricList{}
		require.Empty(t, list)
	})

	t.Run("with metrics", func(t *testing.T) {
		list := PodMetricList{
			{Namespace: "ns1", Name: "pod1"},
			{Namespace: "ns2", Name: "pod2"},
		}
		require.Len(t, list, 2)
	})
}

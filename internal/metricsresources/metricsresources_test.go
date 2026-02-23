package metricsresources

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func posResourceList(name, namespace, container string) pods.PodResourceList {
	return []pods.PodResource{
		{
			NamespaceName: pods.NamespaceName{
				Name:      name,
				Namespace: namespace,
			},
			Containers: []pods.ContainerResource{
				{
					Name: container,
					Limits: pods.Resource{
						CPU:    1,
						Memory: 1024,
					},
					Requests: pods.Resource{
						CPU:    1,
						Memory: 1024,
					},
				},
			},
		},
	}
}

func posMetricsList(name, namespace, container string) podmetrics.PodMetricList {
	return []podmetrics.PodMetric{
		{
			Name:      name,
			Namespace: namespace,
			Containers: []podmetrics.ContainerMetric{
				{
					Name: container,
					Metric: podmetrics.Metric{
						CPU:    2000,
						Memory: 2048,
					},
				},
			},
		},
	}
}

func TestMergeSameNamespaceAndName(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})))
	podResourceList := posResourceList("foo", "bar", "foo-container")
	podMetricList := posMetricsList("foo", "bar", "foo-container")
	podResources := merge(podResourceList, podMetricList)
	require.Len(t, podResources, 1)
	for _, pod := range podResources {
		require.Len(t, pod.PodResource.Containers, 1)
		containers := pod.ContainersMetrics()
		require.Len(t, containers, 1)
		require.Equal(t, int64(2000), containers[0].Requests.CPUUsed)
		require.Equal(t, int64(2048), containers[0].Requests.MemoryUsed)
	}
}

func TestMergeDifferentNamespaceAndName(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})))
	podResourceList := posResourceList("foo", "bar", "foo-container")
	podMetricList := posMetricsList("foo1", "bar", "foo-container")
	podResources := merge(podResourceList, podMetricList)
	require.Len(t, podResources, 1)
	for _, pod := range podResources {
		containers := pod.ContainersMetrics()
		require.Len(t, containers, 1)
		require.Equal(t, unset, containers[0].Requests.CPUUsed)
		require.Equal(t, unset, containers[0].Requests.MemoryUsed)
	}
}

func TestContainersMetricsWithoutMetrics(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})))
	podResourceList := posResourceList("foo", "bar", "foo-container")
	podMetricList := []podmetrics.PodMetric{}
	podResources := merge(podResourceList, podMetricList)
	require.Len(t, podResources, 1)
	containers := podResources[0].ContainersMetrics()
	require.Len(t, containers, 1)
	require.Equal(t, unset, containers[0].Requests.CPUUsed)
	require.Equal(t, unset, containers[0].Requests.MemoryUsed)
}

func TestNewPodRepository(t *testing.T) {
	repo := NewPodRepository()
	require.NotNil(t, repo)
}

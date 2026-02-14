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
	pods := merge(podResourceList, podMetricList)
	require.Len(t, pods, 1)
	for _, pod := range pods {
		require.Len(t, pod.PodResource.Containers, 1)
		require.Contains(t, pod.String(), "/", pod.String())
	}
	require.Contains(t, pods.String(), "/", pods.String())
}

func TestMergeDifferentNamespaceAndName(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})))
	podResourceList := posResourceList("foo", "bar", "foo-container")
	podMetricList := posMetricsList("foo1", "bar", "foo-container")
	pods := merge(podResourceList, podMetricList)
	require.Len(t, pods, 1)
	for _, pod := range pods {
		require.NotContains(t, pod.String(), "/", pod.String())
	}
	require.NotContains(t, pods.String(), "/", pods.String())
}

func TestStringify(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})))
	podResourceList := posResourceList("foo", "bar", "foo-container")
	podMetricList := []podmetrics.PodMetric{}
	pods := merge(podResourceList, podMetricList)
	require.Len(t, pods, 1)
	text := pods[0].String()
	require.Greater(t, len(text), 0)
	require.NotContains(t, text, "/", text)
	require.NotContains(t, pods.String(), "/", pods.String())
}

func TestConfigBuilder(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		builder := NewConfigBuilder()
		config := builder.Build()
		require.Equal(t, uint(5), config.WatchPeriod)
	})

	t.Run("with all options", func(t *testing.T) {
		config := NewConfigBuilder().
			WithKubeConfig("/custom/kubeconfig").
			WithKubeContext("custom-context").
			WithNamespaces([]string{"production"}).
			WithLabel("app=web").
			WithFieldSelector("status.phase=Running").
			WithNodes([]string{"node1", "node2"}).
			WithOutput("json").
			WithSorting("cpu").
			WithResources([]string{"cpu", "memory"}).
			WithAlert("memory").
			WithWatchPeriod(10).
			WithReverse(true).
			WithWatchMetrics(true).
			Build()

		require.Equal(t, "/custom/kubeconfig", config.KubeConfig)
		require.Equal(t, "custom-context", config.KubeContext)
		require.Equal(t, []string{"production"}, config.Namespaces)
		require.Equal(t, "app=web", config.Label)
		require.Equal(t, "status.phase=Running", config.FieldSelector)
		require.Equal(t, []string{"node1", "node2"}, config.Nodes)
		require.Equal(t, "json", config.Output)
		require.Equal(t, "cpu", config.Sorting)
		require.Equal(t, []string{"cpu", "memory"}, config.Resources)
		require.Equal(t, "memory", config.Alert)
		require.Equal(t, uint(10), config.WatchPeriod)
		require.True(t, config.Reverse)
		require.True(t, config.WatchMetrics)
	})

	t.Run("chained calls", func(t *testing.T) {
		builder := NewConfigBuilder()
		require.NotNil(t, builder)

		result := builder.WithNamespaces([]string{"ns1"}).WithOutput("yaml")
		require.Equal(t, builder, result)
	})
}

func TestNewPodRepository(t *testing.T) {
	repo := NewPodRepository()
	require.NotNil(t, repo)
}

package metricsresources

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func TestConfig(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		config := Config{}
		require.Empty(t, config.KubeConfig)
		require.Empty(t, config.Namespaces)
		require.False(t, config.WatchMetrics)
	})

	t.Run("with values", func(t *testing.T) {
		config := Config{
			KubeConfig:    "/path/to/config",
			KubeContext:   "test-context",
			Namespaces:    []string{"test-ns"},
			Label:         "app=test",
			FieldSelector: "spec.nodeName=node1",
			Nodes:         []string{"node1", "node2"},
			Output:        "json",
			Sorting:       "name",
			Resources:     []string{"cpu", "memory"},
			Alert:         "memory",
			WatchPeriod:   10,
			Reverse:       true,
			WatchMetrics:  true,
		}
		require.Equal(t, "/path/to/config", config.KubeConfig)
		require.Equal(t, "test-context", config.KubeContext)
		require.Equal(t, []string{"test-ns"}, config.Namespaces)
		require.True(t, config.WatchMetrics)
	})
}

func TestWatchResponse(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		resp := WatchResponse{error: errors.New("test error")}
		require.Error(t, resp.error)
		require.Empty(t, resp.data)
	})

	t.Run("with data", func(t *testing.T) {
		data := PodMetricsResourceList{{PodResource: pods.PodResource{NamespaceName: pods.NamespaceName{Name: "test"}}}}
		resp := WatchResponse{data: data}
		require.NoError(t, resp.error)
		require.Len(t, resp.data, 1)
	})
}

func TestMerge(t *testing.T) {
	t.Run("empty lists", func(t *testing.T) {
		result := merge(pods.PodResourceList{}, podmetrics.PodMetricList{})
		require.Empty(t, result)
	})

	t.Run("pod resource without metrics", func(t *testing.T) {
		podResources := pods.PodResourceList{
			{NamespaceName: pods.NamespaceName{Namespace: "ns1", Name: "pod1"}},
		}
		result := merge(podResources, podmetrics.PodMetricList{})
		require.Len(t, result, 1)
		require.Equal(t, "pod1", result[0].NamespaceName.Name)
		require.Equal(t, "ns1", result[0].NamespaceName.Namespace)
	})

	t.Run("matching pod and metrics", func(t *testing.T) {
		podResources := pods.PodResourceList{
			{NamespaceName: pods.NamespaceName{Namespace: "ns1", Name: "pod1"}},
		}
		podMetrics := podmetrics.PodMetricList{
			{Namespace: "ns1", Name: "pod1"},
		}
		result := merge(podResources, podMetrics)
		require.Len(t, result, 1)
		require.Equal(t, "pod1", result[0].NamespaceName.Name)
		require.Equal(t, "ns1", result[0].NamespaceName.Namespace)
	})

	t.Run("mismatched namespace/name", func(t *testing.T) {
		podResources := pods.PodResourceList{
			{NamespaceName: pods.NamespaceName{Namespace: "ns1", Name: "pod1"}},
		}
		podMetrics := podmetrics.PodMetricList{
			{Namespace: "ns2", Name: "pod2"},
		}
		result := merge(podResources, podMetrics)
		require.Len(t, result, 1)

		output := result.toOutput()
		require.Len(t, output.Items, 1)
		require.Equal(t, "pod1", output.Items[0].Name)
		require.Equal(t, "ns1", output.Items[0].Namespace)
	})

	t.Run("mismatched namespace/name does not emit warn log", func(t *testing.T) {
		oldLogger := slog.Default()
		defer slog.SetDefault(oldLogger)

		var logBuffer bytes.Buffer
		handler := slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{Level: slog.LevelDebug})
		slog.SetDefault(slog.New(handler))

		podResources := pods.PodResourceList{
			{NamespaceName: pods.NamespaceName{Namespace: "ns1", Name: "pod1"}},
		}
		podMetrics := podmetrics.PodMetricList{
			{Namespace: "ns2", Name: "pod2"},
		}

		result := merge(podResources, podMetrics)
		require.Len(t, result, 1)

		logs := logBuffer.String()
		require.NotContains(t, logs, `"level":"WARN"`)
		require.Contains(t, logs, "Skipped unmatched pod metrics")
	})
}

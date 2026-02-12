package noderesources

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trezorg/k8spodsmetrics/pkg/nodemetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/nodes"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func nodeResourceList(name string) NodeResourceList {
	return []NodeResource{
		{
			Name:          name,
			CPU:           1024,
			Memory:        1024,
			CPURequest:    512,
			MemoryRequest: 512,
			CPULimit:      512,
			MemoryLimit:   512,
		},
	}
}

func TestStringify(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})))
	nodeResourceList := nodeResourceList("foo")
	text := nodeResourceList.String()
	require.Greater(t, len(text), 0)
	require.Contains(t, text, "/", text)
}

func TestConfigBuilder(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		builder := NewConfigBuilder()
		config := builder.Build()
		require.Equal(t, uint(5), config.WatchPeriod)
		require.Equal(t, uint(3), config.KLogLevel)
	})

	t.Run("with all options", func(t *testing.T) {
		config := NewConfigBuilder().
			WithKubeConfig("/custom/kubeconfig").
			WithKubeContext("custom-context").
			WithName("worker-1").
			WithLabel("node-role=worker").
			WithOutput("json").
			WithSorting("cpu").
			WithResources([]string{"cpu", "memory"}).
			WithAlert("cpu").
			WithKLogLevel(5).
			WithWatchPeriod(10).
			WithReverse(true).
			WithWatchMetrics(true).
			Build()

		require.Equal(t, "/custom/kubeconfig", config.KubeConfig)
		require.Equal(t, "custom-context", config.KubeContext)
		require.Equal(t, "worker-1", config.Name)
		require.Equal(t, "node-role=worker", config.Label)
		require.Equal(t, "json", config.Output)
		require.Equal(t, "cpu", config.Sorting)
		require.Equal(t, []string{"cpu", "memory"}, config.Resources)
		require.Equal(t, "cpu", config.Alert)
		require.Equal(t, uint(5), config.KLogLevel)
		require.Equal(t, uint(10), config.WatchPeriod)
		require.True(t, config.Reverse)
		require.True(t, config.WatchMetrics)
	})
}

func TestMergeNodeResources(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})))

	t.Run("empty lists", func(t *testing.T) {
		result := merge(pods.PodResourceList{}, nodes.NodeList{}, nodemetrics.List{})
		require.Empty(t, result)
	})

	t.Run("only nodes", func(t *testing.T) {
		nodeList := nodes.NodeList{{Name: "node1", CPU: 4000, Memory: 16 * 1024 * 1024 * 1024}}
		result := merge(pods.PodResourceList{}, nodeList, nodemetrics.List{})
		require.Len(t, result, 1)
		require.Equal(t, "node1", result[0].Name)
		require.Equal(t, int64(4000), result[0].CPU)
	})

	t.Run("node with pods", func(t *testing.T) {
		nodeList := nodes.NodeList{{Name: "node1", AllocatableCPU: 4000, AllocatableMemory: 16 * 1024 * 1024 * 1024}}
		podList := pods.PodResourceList{
			{
				NodeName: "node1",
				Containers: []pods.ContainerResource{
					{Name: "c1", Requests: pods.Resource{CPU: 500, Memory: 512 * 1024 * 1024}},
				},
			},
		}
		result := merge(podList, nodeList, nodemetrics.List{})
		require.Len(t, result, 1)
		require.Equal(t, int64(500), result[0].CPURequest)
		require.Equal(t, int64(512*1024*1024), result[0].MemoryRequest)
		require.Equal(t, int64(3500), result[0].AvailableCPU)
	})

	t.Run("node with metrics", func(t *testing.T) {
		nodeList := nodes.NodeList{{Name: "node1", AllocatableCPU: 4000, AllocatableMemory: 16 * 1024 * 1024 * 1024}}
		metricsList := nodemetrics.List{{Name: "node1", CPU: 2000, Memory: 8 * 1024 * 1024 * 1024}}
		result := merge(pods.PodResourceList{}, nodeList, metricsList)
		require.Len(t, result, 1)
		require.Equal(t, int64(2000), result[0].UsedCPU)
		require.Equal(t, int64(8*1024*1024*1024), result[0].UsedMemory)
		require.Equal(t, int64(2000), result[0].FreeCPU)
	})
}

func TestNewNodeRepository(t *testing.T) {
	repo := NewNodeRepository()
	require.NotNil(t, repo)
}

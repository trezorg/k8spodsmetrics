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

func TestNodeResource_IsStorageAlerted(t *testing.T) {
	t.Run("zero storage is not alerted", func(t *testing.T) {
		resource := NodeResource{Storage: 0, UsedStorage: 100}
		require.False(t, resource.IsStorageAlerted())
	})

	t.Run("high storage usage is alerted", func(t *testing.T) {
		resource := NodeResource{Storage: 100, UsedStorage: 96}
		require.True(t, resource.IsStorageAlerted())
	})

	t.Run("usage at threshold is not alerted", func(t *testing.T) {
		resource := NodeResource{Storage: 100, UsedStorage: 95}
		require.False(t, resource.IsStorageAlerted())
	})
}

func TestNodeResource_IsStorageEphemeralAlerted(t *testing.T) {
	t.Run("zero ephemeral storage is not alerted", func(t *testing.T) {
		resource := NodeResource{StorageEphemeral: 0, UsedStorageEphemeral: 100}
		require.False(t, resource.IsStorageEphemeralAlerted())
	})

	t.Run("high ephemeral storage usage is alerted", func(t *testing.T) {
		resource := NodeResource{StorageEphemeral: 100, UsedStorageEphemeral: 96}
		require.True(t, resource.IsStorageEphemeralAlerted())
	})

	t.Run("ephemeral usage at threshold is not alerted", func(t *testing.T) {
		resource := NodeResource{StorageEphemeral: 100, UsedStorageEphemeral: 95}
		require.False(t, resource.IsStorageEphemeralAlerted())
	})
}

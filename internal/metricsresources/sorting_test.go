package metricsresources

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trezorg/k8spodsmetrics/internal/sorting/metricsresources"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func testPodMetricsResource(name, namespace string, containers []pods.ContainerResource, metrics []podmetrics.ContainerMetric) PodMetricsResource {
	return PodMetricsResource{
		PodResource: pods.PodResource{
			NamespaceName: pods.NamespaceName{
				Name:      name,
				Namespace: namespace,
			},
			Containers: containers,
		},
		PodMetric: podmetrics.PodMetric{
			Name:       name,
			Namespace:  namespace,
			Containers: metrics,
		},
	}
}

func testPodMetricsResourceWithNode(name, namespace, nodeName string, containers []pods.ContainerResource, metrics []podmetrics.ContainerMetric) PodMetricsResource {
	r := testPodMetricsResource(name, namespace, containers, metrics)
	r.NodeName = nodeName
	return r
}

func testContainerResource(name string, cpuReq, cpuLimit, memReq, memLimit int64) pods.ContainerResource {
	return pods.ContainerResource{
		Name: name,
		Requests: pods.Resource{
			CPU:    cpuReq,
			Memory: memReq,
		},
		Limits: pods.Resource{
			CPU:    cpuLimit,
			Memory: memLimit,
		},
	}
}

func testContainerMetric(name string, cpu, memory int64) podmetrics.ContainerMetric {
	return podmetrics.ContainerMetric{
		Name: name,
		Metric: podmetrics.Metric{
			CPU:    cpu,
			Memory: memory,
		},
	}
}

func TestCPURequest(t *testing.T) {
	tests := []struct {
		name       string
		containers []pods.ContainerResource
		expected   int64
	}{
		{"empty", []pods.ContainerResource{}, 0},
		{"single container", []pods.ContainerResource{testContainerResource("c1", 100, 200, 1024, 2048)}, 100},
		{"multiple containers", []pods.ContainerResource{
			testContainerResource("c1", 100, 200, 1024, 2048),
			testContainerResource("c2", 50, 100, 512, 1024),
		}, 150},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := cpuRequest(tc.containers)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestCPULimit(t *testing.T) {
	tests := []struct {
		name       string
		containers []pods.ContainerResource
		expected   int64
	}{
		{"empty", []pods.ContainerResource{}, 0},
		{"single container", []pods.ContainerResource{testContainerResource("c1", 100, 200, 1024, 2048)}, 200},
		{"multiple containers", []pods.ContainerResource{
			testContainerResource("c1", 100, 200, 1024, 2048),
			testContainerResource("c2", 50, 100, 512, 1024),
		}, 300},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := cpuLimit(tc.containers)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestCPUUsed(t *testing.T) {
	tests := []struct {
		name       string
		containers []podmetrics.ContainerMetric
		expected   int64
	}{
		{"empty", []podmetrics.ContainerMetric{}, 0},
		{"single container", []podmetrics.ContainerMetric{testContainerMetric("c1", 150, 1024)}, 150},
		{"multiple containers", []podmetrics.ContainerMetric{
			testContainerMetric("c1", 150, 1024),
			testContainerMetric("c2", 75, 512),
		}, 225},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := cpuUsed(tc.containers)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestMemoryRequest(t *testing.T) {
	tests := []struct {
		name       string
		containers []pods.ContainerResource
		expected   int64
	}{
		{"empty", []pods.ContainerResource{}, 0},
		{"single container", []pods.ContainerResource{testContainerResource("c1", 100, 200, 1024, 2048)}, 1024},
		{"multiple containers", []pods.ContainerResource{
			testContainerResource("c1", 100, 200, 1024, 2048),
			testContainerResource("c2", 50, 100, 512, 1024),
		}, 1536},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := memoryRequest(tc.containers)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestMemoryLimit(t *testing.T) {
	tests := []struct {
		name       string
		containers []pods.ContainerResource
		expected   int64
	}{
		{"empty", []pods.ContainerResource{}, 0},
		{"single container", []pods.ContainerResource{testContainerResource("c1", 100, 200, 1024, 2048)}, 2048},
		{"multiple containers", []pods.ContainerResource{
			testContainerResource("c1", 100, 200, 1024, 2048),
			testContainerResource("c2", 50, 100, 512, 1024),
		}, 3072},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := memoryLimit(tc.containers)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestMemoryUsed(t *testing.T) {
	tests := []struct {
		name       string
		containers []podmetrics.ContainerMetric
		expected   int64
	}{
		{"empty", []podmetrics.ContainerMetric{}, 0},
		{"single container", []podmetrics.ContainerMetric{testContainerMetric("c1", 100, 1024)}, 1024},
		{"multiple containers", []podmetrics.ContainerMetric{
			testContainerMetric("c1", 100, 1024),
			testContainerMetric("c2", 50, 512),
		}, 1536},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := memoryUsed(tc.containers)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestStorageUsed(t *testing.T) {
	containers := []podmetrics.ContainerMetric{
		{Name: "c1", Metric: podmetrics.Metric{Storage: 1000}},
		{Name: "c2", Metric: podmetrics.Metric{Storage: 2000}},
	}
	require.Equal(t, int64(3000), storageUsed(containers))
	require.Equal(t, int64(0), storageUsed([]podmetrics.ContainerMetric{}))
}

func TestStorageEphemeralUsed(t *testing.T) {
	containers := []podmetrics.ContainerMetric{
		{Name: "c1", Metric: podmetrics.Metric{StorageEphemeral: 500}},
		{Name: "c2", Metric: podmetrics.Metric{StorageEphemeral: 700}},
	}
	require.Equal(t, int64(1200), storageEphemeralUsed(containers))
	require.Equal(t, int64(0), storageEphemeralUsed([]podmetrics.ContainerMetric{}))
}

func TestReverse(t *testing.T) {
	t.Run("reverses less function", func(t *testing.T) {
		less := func(i, j int) bool { return i < j }
		reversed := reverse(less)
		require.True(t, less(1, 2), "1 < 2 should be true")
		require.False(t, reversed(1, 2), "reversed(1, 2) should be false for i < j")
		require.True(t, reversed(2, 1), "reversed(2, 1) should be true for j < i")
	})
}

func TestSortByName(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-c", "ns1", nil, nil),
		testPodMetricsResource("pod-a", "ns1", nil, nil),
		testPodMetricsResource("pod-b", "ns1", nil, nil),
	}

	list.sortByName(false)
	require.Equal(t, "pod-a", list[0].Name)
	require.Equal(t, "pod-b", list[1].Name)
	require.Equal(t, "pod-c", list[2].Name)

	list.sortByName(true)
	require.Equal(t, "pod-c", list[0].Name)
	require.Equal(t, "pod-b", list[1].Name)
	require.Equal(t, "pod-a", list[2].Name)
}

func TestSortByNode(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResourceWithNode("pod-1", "ns1", "node-c", nil, nil),
		testPodMetricsResourceWithNode("pod-2", "ns1", "node-a", nil, nil),
		testPodMetricsResourceWithNode("pod-3", "ns1", "node-b", nil, nil),
	}

	list.sortByNode(false)
	require.Equal(t, "node-a", list[0].NodeName)
	require.Equal(t, "node-b", list[1].NodeName)
	require.Equal(t, "node-c", list[2].NodeName)

	list.sortByNode(true)
	require.Equal(t, "node-c", list[0].NodeName)
	require.Equal(t, "node-b", list[1].NodeName)
	require.Equal(t, "node-a", list[2].NodeName)
}

func TestSortByNamespace(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-b", "ns-c", nil, nil),
		testPodMetricsResource("pod-c", "ns-a", nil, nil),
		testPodMetricsResource("pod-a", "ns-b", nil, nil),
		testPodMetricsResource("pod-d", "ns-a", nil, nil),
	}

	list.sortByNamespace(false)
	require.Equal(t, "ns-a", list[0].Namespace)
	require.Equal(t, "pod-c", list[0].Name)
	require.Equal(t, "ns-a", list[1].Namespace)
	require.Equal(t, "pod-d", list[1].Name)
	require.Equal(t, "ns-b", list[2].Namespace)
	require.Equal(t, "ns-c", list[3].Namespace)

	list.sortByNamespace(true)
	require.Equal(t, "ns-c", list[0].Namespace)
	require.Equal(t, "ns-b", list[1].Namespace)
	require.Equal(t, "ns-a", list[2].Namespace)
}

func TestSortByRequestCPU(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-1", "ns1", []pods.ContainerResource{testContainerResource("c1", 300, 0, 0, 0)}, nil),
		testPodMetricsResource("pod-2", "ns1", []pods.ContainerResource{testContainerResource("c1", 100, 0, 0, 0)}, nil),
		testPodMetricsResource("pod-3", "ns1", []pods.ContainerResource{testContainerResource("c1", 200, 0, 0, 0)}, nil),
	}

	list.sortByRequestCPU(false)
	require.Equal(t, "pod-2", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-1", list[2].Name)

	list.sortByRequestCPU(true)
	require.Equal(t, "pod-1", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-2", list[2].Name)
}

func TestSortByLimitCPU(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-1", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 500, 0, 0)}, nil),
		testPodMetricsResource("pod-2", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 200, 0, 0)}, nil),
		testPodMetricsResource("pod-3", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 300, 0, 0)}, nil),
	}

	list.sortByLimitCPU(false)
	require.Equal(t, "pod-2", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-1", list[2].Name)

	list.sortByLimitCPU(true)
	require.Equal(t, "pod-1", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-2", list[2].Name)
}

func TestSortByUsedCPU(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-1", "ns1", nil, []podmetrics.ContainerMetric{testContainerMetric("c1", 400, 0)}),
		testPodMetricsResource("pod-2", "ns1", nil, []podmetrics.ContainerMetric{testContainerMetric("c1", 100, 0)}),
		testPodMetricsResource("pod-3", "ns1", nil, []podmetrics.ContainerMetric{testContainerMetric("c1", 250, 0)}),
	}

	list.sortByUsedCPU(false)
	require.Equal(t, "pod-2", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-1", list[2].Name)

	list.sortByUsedCPU(true)
	require.Equal(t, "pod-1", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-2", list[2].Name)
}

func TestSortByRequestMemory(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-1", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 0, 4096, 0)}, nil),
		testPodMetricsResource("pod-2", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 0, 1024, 0)}, nil),
		testPodMetricsResource("pod-3", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 0, 2048, 0)}, nil),
	}

	list.sortByRequestMemory(false)
	require.Equal(t, "pod-2", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-1", list[2].Name)

	list.sortByRequestMemory(true)
	require.Equal(t, "pod-1", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-2", list[2].Name)
}

func TestSortByLimitMemory(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-1", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 0, 0, 8192)}, nil),
		testPodMetricsResource("pod-2", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 0, 0, 2048)}, nil),
		testPodMetricsResource("pod-3", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 0, 0, 4096)}, nil),
	}

	list.sortByLimitMemory(false)
	require.Equal(t, "pod-2", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-1", list[2].Name)

	list.sortByLimitMemory(true)
	require.Equal(t, "pod-1", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-2", list[2].Name)
}

func TestSortByUsedMemory(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-1", "ns1", nil, []podmetrics.ContainerMetric{testContainerMetric("c1", 0, 5000)}),
		testPodMetricsResource("pod-2", "ns1", nil, []podmetrics.ContainerMetric{testContainerMetric("c1", 0, 1000)}),
		testPodMetricsResource("pod-3", "ns1", nil, []podmetrics.ContainerMetric{testContainerMetric("c1", 0, 3000)}),
	}

	list.sortByUsedMemory(false)
	require.Equal(t, "pod-2", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-1", list[2].Name)

	list.sortByUsedMemory(true)
	require.Equal(t, "pod-1", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-2", list[2].Name)
}

func TestSortByUsedStorage(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-1", "ns1", nil, []podmetrics.ContainerMetric{{Name: "c1", Metric: podmetrics.Metric{Storage: 3000}}}),
		testPodMetricsResource("pod-2", "ns1", nil, []podmetrics.ContainerMetric{{Name: "c1", Metric: podmetrics.Metric{Storage: 1000}}}),
		testPodMetricsResource("pod-3", "ns1", nil, []podmetrics.ContainerMetric{{Name: "c1", Metric: podmetrics.Metric{Storage: 2000}}}),
	}

	list.sortByUsedStorage(false)
	require.Equal(t, "pod-2", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-1", list[2].Name)

	list.sortByUsedStorage(true)
	require.Equal(t, "pod-1", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-2", list[2].Name)
}

func TestSortByUsedStorageEphemeral(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-1", "ns1", nil, []podmetrics.ContainerMetric{{Name: "c1", Metric: podmetrics.Metric{StorageEphemeral: 4000}}}),
		testPodMetricsResource("pod-2", "ns1", nil, []podmetrics.ContainerMetric{{Name: "c1", Metric: podmetrics.Metric{StorageEphemeral: 1000}}}),
		testPodMetricsResource("pod-3", "ns1", nil, []podmetrics.ContainerMetric{{Name: "c1", Metric: podmetrics.Metric{StorageEphemeral: 2000}}}),
	}

	list.sortByUsedStorageEphemeral(false)
	require.Equal(t, "pod-2", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-1", list[2].Name)

	list.sortByUsedStorageEphemeral(true)
	require.Equal(t, "pod-1", list[0].Name)
	require.Equal(t, "pod-3", list[1].Name)
	require.Equal(t, "pod-2", list[2].Name)
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		by       string
		reverse  bool
		list     PodMetricsResourceList
		expected []string
	}{
		{
			name:    "sort by name ascending",
			by:      string(metricsresources.Name),
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns1", nil, nil),
				testPodMetricsResource("pod-a", "ns1", nil, nil),
			},
			expected: []string{"pod-a", "pod-b"},
		},
		{
			name:    "sort by name descending",
			by:      string(metricsresources.Name),
			reverse: true,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-a", "ns1", nil, nil),
				testPodMetricsResource("pod-b", "ns1", nil, nil),
			},
			expected: []string{"pod-b", "pod-a"},
		},
		{
			name:    "sort by request_cpu ascending",
			by:      string(metricsresources.RequestCPU),
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns1", []pods.ContainerResource{testContainerResource("c1", 200, 0, 0, 0)}, nil),
				testPodMetricsResource("pod-a", "ns1", []pods.ContainerResource{testContainerResource("c1", 100, 0, 0, 0)}, nil),
			},
			expected: []string{"pod-a", "pod-b"},
		},
		{
			name:    "sort by request_cpu descending",
			by:      string(metricsresources.RequestCPU),
			reverse: true,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-a", "ns1", []pods.ContainerResource{testContainerResource("c1", 100, 0, 0, 0)}, nil),
				testPodMetricsResource("pod-b", "ns1", []pods.ContainerResource{testContainerResource("c1", 200, 0, 0, 0)}, nil),
			},
			expected: []string{"pod-b", "pod-a"},
		},
		{
			name:    "sort by namespace ascending",
			by:      string(metricsresources.Namespace),
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns-b", nil, nil),
				testPodMetricsResource("pod-a", "ns-a", nil, nil),
			},
			expected: []string{"pod-a", "pod-b"},
		},
		{
			name:    "sort by limit_cpu ascending",
			by:      string(metricsresources.LimitCPU),
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 200, 0, 0)}, nil),
				testPodMetricsResource("pod-a", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 100, 0, 0)}, nil),
			},
			expected: []string{"pod-a", "pod-b"},
		},
		{
			name:    "sort by used_cpu ascending",
			by:      string(metricsresources.UsedCPU),
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns1", nil, []podmetrics.ContainerMetric{testContainerMetric("c1", 200, 0)}),
				testPodMetricsResource("pod-a", "ns1", nil, []podmetrics.ContainerMetric{testContainerMetric("c1", 100, 0)}),
			},
			expected: []string{"pod-a", "pod-b"},
		},
		{
			name:    "sort by request_memory ascending",
			by:      string(metricsresources.RequestMemory),
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 0, 2048, 0)}, nil),
				testPodMetricsResource("pod-a", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 0, 1024, 0)}, nil),
			},
			expected: []string{"pod-a", "pod-b"},
		},
		{
			name:    "sort by limit_memory ascending",
			by:      string(metricsresources.LimitMemory),
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 0, 0, 4096)}, nil),
				testPodMetricsResource("pod-a", "ns1", []pods.ContainerResource{testContainerResource("c1", 0, 0, 0, 2048)}, nil),
			},
			expected: []string{"pod-a", "pod-b"},
		},
		{
			name:    "sort by used_memory ascending",
			by:      string(metricsresources.UsedMemory),
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns1", nil, []podmetrics.ContainerMetric{testContainerMetric("c1", 0, 2048)}),
				testPodMetricsResource("pod-a", "ns1", nil, []podmetrics.ContainerMetric{testContainerMetric("c1", 0, 1024)}),
			},
			expected: []string{"pod-a", "pod-b"},
		},
		{
			name:    "sort by used_storage ascending",
			by:      string(metricsresources.UsedStorage),
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns1", nil, []podmetrics.ContainerMetric{{Name: "c1", Metric: podmetrics.Metric{Storage: 2000}}}),
				testPodMetricsResource("pod-a", "ns1", nil, []podmetrics.ContainerMetric{{Name: "c1", Metric: podmetrics.Metric{Storage: 1000}}}),
			},
			expected: []string{"pod-a", "pod-b"},
		},
		{
			name:    "sort by used_storage_ephemeral ascending",
			by:      string(metricsresources.UsedStorageEphemeral),
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns1", nil, []podmetrics.ContainerMetric{{Name: "c1", Metric: podmetrics.Metric{StorageEphemeral: 2000}}}),
				testPodMetricsResource("pod-a", "ns1", nil, []podmetrics.ContainerMetric{{Name: "c1", Metric: podmetrics.Metric{StorageEphemeral: 1000}}}),
			},
			expected: []string{"pod-a", "pod-b"},
		},
		{
			name:    "unknown sort keeps order",
			by:      "unknown",
			reverse: false,
			list: PodMetricsResourceList{
				testPodMetricsResource("pod-b", "ns1", nil, nil),
				testPodMetricsResource("pod-a", "ns1", nil, nil),
			},
			expected: []string{"pod-b", "pod-a"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.list.sort(tc.by, tc.reverse)
			for i, name := range tc.expected {
				require.Equal(t, name, tc.list[i].Name, "position %d", i)
			}
		})
	}
}

func TestSortEmptyList(t *testing.T) {
	list := PodMetricsResourceList{}
	require.NotPanics(t, func() {
		list.sort(string(metricsresources.Name), false)
	})
	require.Empty(t, list)
}

func TestSortSingleElement(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-1", "ns1", nil, nil),
	}
	list.sort(string(metricsresources.Name), false)
	require.Len(t, list, 1)
	require.Equal(t, "pod-1", list[0].Name)
}

func TestSortWithEqualValues(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-b", "ns1", []pods.ContainerResource{testContainerResource("c1", 100, 0, 0, 0)}, nil),
		testPodMetricsResource("pod-a", "ns1", []pods.ContainerResource{testContainerResource("c1", 100, 0, 0, 0)}, nil),
		testPodMetricsResource("pod-c", "ns1", []pods.ContainerResource{testContainerResource("c1", 100, 0, 0, 0)}, nil),
	}

	list.sortByRequestCPU(false)
	require.Equal(t, "pod-b", list[0].Name)
	require.Equal(t, "pod-a", list[1].Name)
	require.Equal(t, "pod-c", list[2].Name)
}

func TestSortByNameUsesPodResourceIdentity(t *testing.T) {
	list := PodMetricsResourceList{
		{
			PodResource: pods.PodResource{
				NamespaceName: pods.NamespaceName{Namespace: "ns1", Name: "pod-b"},
			},
		},
		{
			PodResource: pods.PodResource{
				NamespaceName: pods.NamespaceName{Namespace: "ns1", Name: "pod-a"},
			},
		},
	}

	list.sortByName(false)
	require.Equal(t, "pod-a", list[0].PodResource.Name)
	require.Equal(t, "pod-b", list[1].PodResource.Name)
}

func TestSortByNamespaceUsesPodResourceIdentity(t *testing.T) {
	list := PodMetricsResourceList{
		{
			PodResource: pods.PodResource{
				NamespaceName: pods.NamespaceName{Namespace: "ns-b", Name: "pod-b"},
			},
		},
		{
			PodResource: pods.PodResource{
				NamespaceName: pods.NamespaceName{Namespace: "ns-a", Name: "pod-a"},
			},
		},
	}

	list.sortByNamespace(false)
	require.Equal(t, "ns-a", list[0].PodResource.Namespace)
	require.Equal(t, "ns-b", list[1].PodResource.Namespace)
}

func TestSortWithMultipleContainers(t *testing.T) {
	list := PodMetricsResourceList{
		testPodMetricsResource("pod-1", "ns1", []pods.ContainerResource{
			testContainerResource("c1", 100, 0, 0, 0),
			testContainerResource("c2", 200, 0, 0, 0),
		}, nil),
		testPodMetricsResource("pod-2", "ns1", []pods.ContainerResource{
			testContainerResource("c1", 50, 0, 0, 0),
		}, nil),
	}

	list.sortByRequestCPU(false)
	require.Equal(t, "pod-2", list[0].Name)
	require.Equal(t, "pod-1", list[1].Name)
}

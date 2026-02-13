package metricsresources

import (
	"strconv"
	"testing"

	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func makeBenchmarkPodMetricsResourceList(size, containers int) PodMetricsResourceList {
	list := make(PodMetricsResourceList, size)
	for i := range size {
		resourceContainers := make([]pods.ContainerResource, containers)
		metricContainers := make([]podmetrics.ContainerMetric, containers)
		for j := range containers {
			base := int64((i + 1) * (j + 3))
			resourceContainers[j] = pods.ContainerResource{
				Name: "c-" + strconv.Itoa(j),
				Requests: pods.Resource{
					CPU:              base * 10,
					Memory:           base * 1000,
					Storage:          base * 2000,
					StorageEphemeral: base * 500,
				},
				Limits: pods.Resource{
					CPU:              base * 15,
					Memory:           base * 1500,
					Storage:          base * 2500,
					StorageEphemeral: base * 700,
				},
			}
			metricContainers[j] = podmetrics.ContainerMetric{
				Name: "c-" + strconv.Itoa(j),
				Metric: podmetrics.Metric{
					CPU:              base * 8,
					Memory:           base * 900,
					Storage:          base * 1200,
					StorageEphemeral: base * 300,
				},
			}
		}
		list[i] = PodMetricsResource{
			PodResource: pods.PodResource{
				NamespaceName: pods.NamespaceName{
					Namespace: "ns-" + strconv.Itoa(i%20),
					Name:      "pod-" + strconv.Itoa(i),
				},
				NodeName:   "node-" + strconv.Itoa(i%100),
				Containers: resourceContainers,
			},
			PodMetric: podmetrics.PodMetric{
				Namespace:  "ns-" + strconv.Itoa(i%20),
				Name:       "pod-" + strconv.Itoa(i),
				Containers: metricContainers,
			},
		}
	}
	return list
}

func benchmarkSortByUsedCPU(b *testing.B, size, containers int) {
	base := makeBenchmarkPodMetricsResourceList(size, containers)
	work := make(PodMetricsResourceList, len(base))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(work, base)
		work.sortByUsedCPU(false)
	}
}

func benchmarkSortByRequestMemory(b *testing.B, size, containers int) {
	base := makeBenchmarkPodMetricsResourceList(size, containers)
	work := make(PodMetricsResourceList, len(base))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(work, base)
		work.sortByRequestMemory(false)
	}
}

func BenchmarkSortByUsedCPU1k(b *testing.B) {
	benchmarkSortByUsedCPU(b, 1000, 5)
}

func BenchmarkSortByUsedCPU10k(b *testing.B) {
	benchmarkSortByUsedCPU(b, 10000, 5)
}

func BenchmarkSortByRequestMemory1k(b *testing.B) {
	benchmarkSortByRequestMemory(b, 1000, 5)
}

func BenchmarkSortByRequestMemory10k(b *testing.B) {
	benchmarkSortByRequestMemory(b, 10000, 5)
}

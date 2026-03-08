package metricsresources

import (
	"cmp"
	"slices"

	"github.com/trezorg/k8spodsmetrics/internal/sorting/metricsresources"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func reverse(less func(i, j int) bool) func(i, j int) bool {
	return func(i, j int) bool {
		return less(j, i)
	}
}

func direction(reversed bool, result int) int {
	if reversed {
		return -result
	}
	return result
}

func cpuRequest(containers []pods.ContainerResource) int64 {
	var result int64
	for _, c := range containers {
		result += c.Requests.CPU
	}
	return result
}

func cpuLimit(containers []pods.ContainerResource) int64 {
	var result int64
	for _, c := range containers {
		result += c.Limits.CPU
	}
	return result
}

func cpuUsed(containers []podmetrics.ContainerMetric) int64 {
	var result int64
	for _, c := range containers {
		result += c.CPU
	}
	return result
}

func memoryRequest(containers []pods.ContainerResource) int64 {
	var result int64
	for _, c := range containers {
		result += c.Requests.Memory
	}
	return result
}

func memoryLimit(containers []pods.ContainerResource) int64 {
	var result int64
	for _, c := range containers {
		result += c.Limits.Memory
	}
	return result
}

func memoryUsed(containers []podmetrics.ContainerMetric) int64 {
	var result int64
	for _, c := range containers {
		result += c.Memory
	}
	return result
}

func storageUsed(containers []podmetrics.ContainerMetric) int64 {
	var result int64
	for _, c := range containers {
		result += c.Storage
	}
	return result
}

func storageEphemeralUsed(containers []podmetrics.ContainerMetric) int64 {
	var result int64
	for _, c := range containers {
		result += c.StorageEphemeral
	}
	return result
}

func (r PodMetricsResourceList) sortByNamespace(reversed bool) {
	slices.SortStableFunc(r, func(a, b PodMetricsResource) int {
		return direction(reversed, cmp.Or(
			cmp.Compare(a.PodResource.Namespace, b.PodResource.Namespace),
			cmp.Compare(a.PodResource.Name, b.PodResource.Name),
		))
	})
}

func (r PodMetricsResourceList) sortByName(reversed bool) {
	slices.SortStableFunc(r, func(a, b PodMetricsResource) int {
		return direction(reversed, cmp.Compare(a.PodResource.Name, b.PodResource.Name))
	})
}

func (r PodMetricsResourceList) sortByNode(reversed bool) {
	slices.SortStableFunc(r, func(a, b PodMetricsResource) int {
		return direction(reversed, cmp.Compare(a.NodeName, b.NodeName))
	})
}

func (r PodMetricsResourceList) sortPodResource(reversed bool, f func([]pods.ContainerResource) int64) {
	type sortItem struct {
		resource PodMetricsResource
		key      int64
	}
	items := make([]sortItem, len(r))
	for i := range r {
		items[i] = sortItem{
			resource: r[i],
			key:      f(r[i].PodResource.Containers),
		}
	}
	slices.SortStableFunc(items, func(a, b sortItem) int {
		return direction(reversed, cmp.Compare(a.key, b.key))
	})
	for i := range items {
		r[i] = items[i].resource
	}
}

func (r PodMetricsResourceList) sortPodMetric(reversed bool, f func([]podmetrics.ContainerMetric) int64) {
	type sortItem struct {
		resource PodMetricsResource
		key      int64
	}
	items := make([]sortItem, len(r))
	for i := range r {
		items[i] = sortItem{
			resource: r[i],
			key:      f(r[i].PodMetric.Containers),
		}
	}
	slices.SortStableFunc(items, func(a, b sortItem) int {
		return direction(reversed, cmp.Compare(a.key, b.key))
	})
	for i := range items {
		r[i] = items[i].resource
	}
}

func (r PodMetricsResourceList) sortByRequestCPU(reversed bool) {
	r.sortPodResource(reversed, cpuRequest)
}

func (r PodMetricsResourceList) sortByLimitCPU(reversed bool) {
	r.sortPodResource(reversed, cpuLimit)
}

func (r PodMetricsResourceList) sortByUsedCPU(reversed bool) {
	r.sortPodMetric(reversed, cpuUsed)
}

func (r PodMetricsResourceList) sortByRequestMemory(reversed bool) {
	r.sortPodResource(reversed, memoryRequest)
}

func (r PodMetricsResourceList) sortByLimitMemory(reversed bool) {
	r.sortPodResource(reversed, memoryLimit)
}

func (r PodMetricsResourceList) sortByUsedMemory(reversed bool) {
	r.sortPodMetric(reversed, memoryUsed)
}

func (r PodMetricsResourceList) sortByUsedStorage(reversed bool) {
	r.sortPodMetric(reversed, storageUsed)
}

func (r PodMetricsResourceList) sortByUsedStorageEphemeral(reversed bool) {
	r.sortPodMetric(reversed, storageEphemeralUsed)
}

func (r PodMetricsResourceList) sort(by string, reverse bool) {
	switch metricsresources.Sorting(by) {
	case metricsresources.Name:
		r.sortByName(reverse)
	case metricsresources.Namespace:
		r.sortByNamespace(reverse)
	case metricsresources.Node:
		r.sortByNode(reverse)
	case metricsresources.LimitCPU:
		r.sortByLimitCPU(reverse)
	case metricsresources.RequestCPU:
		r.sortByRequestCPU(reverse)
	case metricsresources.UsedCPU:
		r.sortByUsedCPU(reverse)
	case metricsresources.LimitMemory:
		r.sortByLimitMemory(reverse)
	case metricsresources.RequestMemory:
		r.sortByRequestMemory(reverse)
	case metricsresources.UsedMemory:
		r.sortByUsedMemory(reverse)
	case metricsresources.UsedStorage:
		r.sortByUsedStorage(reverse)
	case metricsresources.UsedStorageEphemeral:
		r.sortByUsedStorageEphemeral(reverse)
	default:
		// keep current order on unknown sorting
		return
	}
}

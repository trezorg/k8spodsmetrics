package metricsresources

import (
	"sort"

	"github.com/trezorg/k8spodsmetrics/internal/sorting/metricsresources"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func reverse(less func(i, j int) bool) func(i, j int) bool {
	return func(i, j int) bool {
		return !less(i, j)
	}
}

func cpuRequest(containers []pods.ContainerResource) int64 {
	result := int64(0)
	for i := 0; i < len(containers); i++ {
		result += containers[i].Requests.CPU
	}
	return result
}

func cpuLimit(containers []pods.ContainerResource) int64 {
	result := int64(0)
	for i := 0; i < len(containers); i++ {
		result += containers[i].Limits.CPU
	}
	return result
}

func cpuUsed(containers []podmetrics.ContainerMetric) int64 {
	result := int64(0)
	for i := 0; i < len(containers); i++ {
		result += containers[i].CPU
	}
	return result
}

func memoryRequest(containers []pods.ContainerResource) int64 {
	result := int64(0)
	for i := 0; i < len(containers); i++ {
		result += containers[i].Requests.Memory
	}
	return result
}

func memoryLimit(containers []pods.ContainerResource) int64 {
	result := int64(0)
	for i := range containers {
		result += containers[i].Limits.Memory
	}
	return result
}

func memoryUsed(containers []podmetrics.ContainerMetric) int64 {
	result := int64(0)
	for i := range containers {
		result += containers[i].Memory
	}
	return result
}

func (r PodMetricsResourceList) sortByNamespace(reversed bool) {
	less := func(i, j int) bool {
		if r[i].Namespace < r[j].Namespace {
			return true
		}
		if r[i].Namespace > r[j].Namespace {
			return false
		}
		return r[i].Name < r[j].Name
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(r, less)
}

func (r PodMetricsResourceList) sortByName(reversed bool) {
	less := func(i, j int) bool {
		return r[i].Name < r[j].Name
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(r, less)
}

func (r PodMetricsResourceList) sortPodResource(reversed bool, f func([]pods.ContainerResource) int64) {
	less := func(i, j int) bool {
		iCon := r[i].PodResource.Containers
		jCon := r[j].PodResource.Containers
		return f(iCon) < f(jCon)
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(r, less)
}

func (r PodMetricsResourceList) sortPodMetric(reversed bool, f func([]podmetrics.ContainerMetric) int64) {
	less := func(i, j int) bool {
		iCon := r[i].PodMetric.Containers
		jCon := r[j].PodMetric.Containers
		return f(iCon) < f(jCon)
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(r, less)
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

func (r PodMetricsResourceList) sort(by string, reverse bool) {
	switch metricsresources.Sorting(by) {
	case metricsresources.Name:
		r.sortByName(reverse)
	case metricsresources.Namespace:
		r.sortByNamespace(reverse)
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
	}
}

package metricsresources

import (
	"slices"

	alerts "github.com/trezorg/k8spodsmetrics/internal/alert"
)

func (r PodMetricsResourceList) filterBy(predicate containerMetricsPredicate) PodMetricsResourceList {
	var result PodMetricsResourceList
	for _, pod := range r {
		if predicate(pod.ContainersMetrics()) {
			result = append(result, pod)
		}
	}
	return result
}

func (r PodMetricsResourceList) filterByPodResource(predicate podResourceMetricsPredicate) PodMetricsResourceList {
	var result PodMetricsResourceList
	for _, pod := range r {
		if predicate(pod) {
			result = append(result, pod)
		}
	}
	return result
}

func (r PodMetricsResourceList) filterNodes(nodes []string) PodMetricsResourceList {
	if len(nodes) == 0 {
		return r
	}
	return r.filterByPodResource(func(r PodMetricsResource) bool {
		return slices.Contains(nodes, r.NodeName)
	})
}

func (r PodMetricsResourceList) filterByAlert(alert alerts.Alert) PodMetricsResourceList {
	switch alert {
	case alerts.Any:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsAlerted() })
	case alerts.Memory:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsMemoryAlerted() })
	case alerts.MemoryRequest:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsMemoryRequestAlerted() })
	case alerts.MemoryLimit:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsMemoryLimitAlerted() })
	case alerts.CPU:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsCPUAlerted() })
	case alerts.CPURequest:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsCPURequestAlerted() })
	case alerts.CPULimit:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsCPULimitAlerted() })
	case alerts.Storage, alerts.StorageEphemeral, alerts.None:
		return r
	}
	return r
}

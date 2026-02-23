package metricsresources

func (r PodMetricsResource) ContainersMetrics() ContainerMetricsResources {
	containerMetricsResources := make(ContainerMetricsResources, 0, len(r.PodMetric.Containers))
	for i, container := range r.PodResource.Containers {
		containerMetricsResource := ContainerMetricsResource{
			Name: container.Name,
		}
		var cpuMetric, memoryMetric, storageMetric, storageEphemeralMetric int64
		if i < len(r.PodMetric.Containers) {
			metrics := r.PodMetric.Containers[i]
			cpuMetric = metrics.CPU
			memoryMetric = metrics.Memory
			storageMetric = metrics.Storage
			storageEphemeralMetric = metrics.StorageEphemeral
		} else {
			cpuMetric = unset
			memoryMetric = unset
			storageMetric = unset
			storageEphemeralMetric = unset
		}
		containerMetricsResource.Requests = MetricsResource{
			CPURequest:              container.Requests.CPU,
			MemoryRequest:           container.Requests.Memory,
			CPUUsed:                 cpuMetric,
			MemoryUsed:              memoryMetric,
			StorageRequest:          container.Requests.Storage,
			StorageEphemeralRequest: container.Requests.StorageEphemeral,
			StorageUsed:             storageMetric,
			StorageEphemeralUsed:    storageEphemeralMetric,
		}
		containerMetricsResource.Limits = MetricsResource{
			CPURequest:              container.Limits.CPU,
			MemoryRequest:           container.Limits.Memory,
			CPUUsed:                 cpuMetric,
			MemoryUsed:              memoryMetric,
			StorageRequest:          container.Limits.Storage,
			StorageEphemeralRequest: container.Limits.StorageEphemeral,
			StorageUsed:             storageMetric,
			StorageEphemeralUsed:    storageEphemeralMetric,
		}
		containerMetricsResources = append(containerMetricsResources, containerMetricsResource)
	}
	return containerMetricsResources
}

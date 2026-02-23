package metricsresources

import (
	"fmt"

	escapes "github.com/snugfox/ansi-escapes"
	"github.com/trezorg/k8spodsmetrics/internal/humanize"
)

func (c ContainerMetricsResource) MemoryUsed() string {
	if c.Limits.MemoryAlert() {
		return c.Limits.MemoryUsedString(escapes.TextColorRed)
	}
	return c.Requests.MemoryUsedString(escapes.TextColorYellow)
}

func (c ContainerMetricsResource) CPUUsed() string {
	if c.Limits.CPUAlert() {
		return c.Limits.CPUUsedString(escapes.TextColorRed)
	}
	return c.Requests.CPUUsedString(escapes.TextColorRed)
}

func (c ContainerMetricsResource) StorageUsed() string {
	return c.Requests.StorageString()
}

func (c ContainerMetricsResource) StorageEphemeralUsed() string {
	return c.Requests.StorageEphemeralString()
}

func (m MetricsResource) CPU(alertColor string) string {
	if m.CPUUsed == unset {
		return fmt.Sprintf("%d", m.CPURequest)
	}
	cpuStartColor := ""
	cpuEndColor := ""
	if m.CPUAlert() {
		cpuStartColor = alertColor
		cpuEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%d/%s%d%s",
		m.CPURequest,
		cpuStartColor,
		m.CPUUsed,
		cpuEndColor,
	)
}

func (m MetricsResource) CPURequestString() string {
	return fmt.Sprintf("%d", m.CPURequest)
}

func (m MetricsResource) CPUUsedString(alertColor string) string {
	if m.CPUUsed == unset {
		return ""
	}
	cpuStartColor := ""
	cpuEndColor := ""
	if m.CPUAlert() {
		cpuStartColor = alertColor
		cpuEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuStartColor,
		m.CPUUsed,
		cpuEndColor,
	)
}

func (m MetricsResource) Memory(alertColor string) string {
	if m.MemoryUsed == unset {
		return humanize.Bytes(m.MemoryRequest)
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if m.MemoryAlert() {
		memoryStartColor = alertColor
		memoryEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s/%s%s%s",
		humanize.Bytes(m.MemoryRequest),
		memoryStartColor,
		humanize.Bytes(m.MemoryUsed),
		memoryEndColor,
	)
}

func (m MetricsResource) MemoryRequestString() string {
	return humanize.Bytes(m.MemoryRequest)
}

func (m MetricsResource) MemoryUsedString(alertColor string) string {
	if m.MemoryUsed == unset {
		return ""
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if m.MemoryAlert() {
		memoryStartColor = alertColor
		memoryEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryStartColor,
		humanize.Bytes(m.MemoryUsed),
		memoryEndColor,
	)
}

func (m MetricsResource) StringWithColor(alertColor string) string {
	cpuStartColor := ""
	cpuEndColor := ""
	if m.CPUAlert() {
		cpuStartColor = alertColor
		cpuEndColor = escapes.ColorReset
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if m.MemoryAlert() {
		memoryStartColor = alertColor
		memoryEndColor = escapes.ColorReset
	}
	if m.MemoryUsed == unset && m.CPUUsed == unset {
		return fmt.Sprintf(
			"CPU=%d, Memory=%s",
			m.CPURequest,
			humanize.Bytes(m.MemoryRequest),
		)
	}
	return fmt.Sprintf(
		"CPU=%d/%s%d%s, Memory=%s/%s%s%s",
		m.CPURequest,
		cpuStartColor,
		m.CPUUsed,
		cpuEndColor,
		humanize.Bytes(m.MemoryRequest),
		memoryStartColor,
		humanize.Bytes(m.MemoryUsed),
		memoryEndColor,
	)
}

func (m MetricsResource) String() string {
	return m.StringWithColor(escapes.TextColorRed)
}

func (m MetricsResource) StorageString() string {
	return humanize.Bytes(m.StorageUsed)
}

func (m MetricsResource) StorageEphemeralString() string {
	return humanize.Bytes(m.StorageEphemeralUsed)
}

func (m MetricsResource) StorageRequestString() string {
	return humanize.Bytes(m.StorageRequest)
}

func (m MetricsResource) StorageEphemeralRequestString() string {
	return humanize.Bytes(m.StorageEphemeralRequest)
}

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

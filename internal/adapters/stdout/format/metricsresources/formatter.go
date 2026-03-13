package metricsresources

import (
	"fmt"
	"strings"

	escapes "github.com/snugfox/ansi-escapes"
	"github.com/trezorg/k8spodsmetrics/internal/humanize"
	servicemetricsresources "github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

const unset = int64(-1)

type MetricsFormatter struct {
	resource servicemetricsresources.MetricsResource
}

func NewMetrics(resource servicemetricsresources.MetricsResource) MetricsFormatter {
	return MetricsFormatter{resource: resource}
}

func (f MetricsFormatter) CPU(alertColor string) string {
	if f.resource.CPUUsed == unset {
		return fmt.Sprintf("%d", f.resource.CPURequest)
	}
	cpuStartColor := ""
	cpuEndColor := ""
	if f.resource.CPUAlert() {
		cpuStartColor = alertColor
		cpuEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%d/%s%d%s",
		f.resource.CPURequest,
		cpuStartColor,
		f.resource.CPUUsed,
		cpuEndColor,
	)
}

func (f MetricsFormatter) CPURequestString() string {
	return fmt.Sprintf("%d", f.resource.CPURequest)
}

func (f MetricsFormatter) CPUUsedString(alertColor string) string {
	if f.resource.CPUUsed == unset {
		return ""
	}
	cpuStartColor := ""
	cpuEndColor := ""
	if f.resource.CPUAlert() {
		cpuStartColor = alertColor
		cpuEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuStartColor,
		f.resource.CPUUsed,
		cpuEndColor,
	)
}

func (f MetricsFormatter) Memory(alertColor string) string {
	if f.resource.MemoryUsed == unset {
		return humanize.Bytes(f.resource.MemoryRequest)
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if f.resource.MemoryAlert() {
		memoryStartColor = alertColor
		memoryEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s/%s%s%s",
		humanize.Bytes(f.resource.MemoryRequest),
		memoryStartColor,
		humanize.Bytes(f.resource.MemoryUsed),
		memoryEndColor,
	)
}

func (f MetricsFormatter) MemoryRequestString() string {
	return humanize.Bytes(f.resource.MemoryRequest)
}

func (f MetricsFormatter) MemoryUsedString(alertColor string) string {
	if f.resource.MemoryUsed == unset {
		return ""
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if f.resource.MemoryAlert() {
		memoryStartColor = alertColor
		memoryEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryStartColor,
		humanize.Bytes(f.resource.MemoryUsed),
		memoryEndColor,
	)
}

func (f MetricsFormatter) StringWithColor(alertColor string) string {
	cpuStartColor := ""
	cpuEndColor := ""
	if f.resource.CPUAlert() {
		cpuStartColor = alertColor
		cpuEndColor = escapes.ColorReset
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if f.resource.MemoryAlert() {
		memoryStartColor = alertColor
		memoryEndColor = escapes.ColorReset
	}
	if f.resource.MemoryUsed == unset && f.resource.CPUUsed == unset {
		return fmt.Sprintf(
			"CPU=%d, Memory=%s",
			f.resource.CPURequest,
			humanize.Bytes(f.resource.MemoryRequest),
		)
	}
	return fmt.Sprintf(
		"CPU=%d/%s%d%s, Memory=%s/%s%s%s",
		f.resource.CPURequest,
		cpuStartColor,
		f.resource.CPUUsed,
		cpuEndColor,
		humanize.Bytes(f.resource.MemoryRequest),
		memoryStartColor,
		humanize.Bytes(f.resource.MemoryUsed),
		memoryEndColor,
	)
}

func (f MetricsFormatter) String() string {
	return f.StringWithColor(escapes.TextColorRed)
}

func (f MetricsFormatter) StorageString() string {
	return humanize.Bytes(f.resource.StorageUsed)
}

func (f MetricsFormatter) StorageEphemeralString() string {
	return humanize.Bytes(f.resource.StorageEphemeralUsed)
}

func (f MetricsFormatter) StorageRequestString() string {
	return humanize.Bytes(f.resource.StorageRequest)
}

func (f MetricsFormatter) StorageEphemeralRequestString() string {
	return humanize.Bytes(f.resource.StorageEphemeralRequest)
}

type ContainerFormatter struct {
	resource servicemetricsresources.ContainerMetricsResource
}

func NewContainer(resource servicemetricsresources.ContainerMetricsResource) ContainerFormatter {
	return ContainerFormatter{resource: resource}
}

func (f ContainerFormatter) Name() string {
	return f.resource.Name
}

func (f ContainerFormatter) Requests() MetricsFormatter {
	return NewMetrics(f.resource.Requests)
}

func (f ContainerFormatter) Limits() MetricsFormatter {
	return NewMetrics(f.resource.Limits)
}

func (f ContainerFormatter) MemoryUsed() string {
	if f.resource.Limits.MemoryAlert() {
		return NewMetrics(f.resource.Limits).MemoryUsedString(escapes.TextColorRed)
	}
	return NewMetrics(f.resource.Requests).MemoryUsedString(escapes.TextColorYellow)
}

func (f ContainerFormatter) CPUUsed() string {
	if f.resource.Limits.CPUAlert() {
		return NewMetrics(f.resource.Limits).CPUUsedString(escapes.TextColorRed)
	}
	return NewMetrics(f.resource.Requests).CPUUsedString(escapes.TextColorRed)
}

func (f ContainerFormatter) StorageUsed() string {
	return NewMetrics(f.resource.Requests).StorageString()
}

func (f ContainerFormatter) StorageEphemeralUsed() string {
	return NewMetrics(f.resource.Requests).StorageEphemeralString()
}

func (f ContainerFormatter) CPUCompactString() string {
	return compactTriple(
		f.Requests().CPURequestString(),
		fallbackEmpty(f.CPUUsed()),
		f.Limits().CPURequestString(),
	)
}

func (f ContainerFormatter) MemoryCompactString() string {
	return compactTriple(
		f.Requests().MemoryRequestString(),
		fallbackEmpty(f.MemoryUsed()),
		f.Limits().MemoryRequestString(),
	)
}

func (f ContainerFormatter) StorageCompactString() string {
	return compactTriple(
		f.Requests().StorageRequestString(),
		fallbackEmpty(f.storageUsedCompactValue()),
		f.Limits().StorageRequestString(),
	)
}

func (f ContainerFormatter) StorageEphemeralCompactString() string {
	return compactTriple(
		f.Requests().StorageEphemeralRequestString(),
		fallbackEmpty(f.storageEphemeralUsedCompactValue()),
		f.Limits().StorageEphemeralRequestString(),
	)
}

func (f ContainerFormatter) storageUsedCompactValue() string {
	if f.resource.Requests.StorageUsed == unset {
		return ""
	}
	return f.StorageUsed()
}

func (f ContainerFormatter) storageEphemeralUsedCompactValue() string {
	if f.resource.Requests.StorageEphemeralUsed == unset {
		return ""
	}
	return f.StorageEphemeralUsed()
}

func compactTriple(first, second, third string) string {
	return strings.Join([]string{fallbackEmpty(first), fallbackEmpty(second), fallbackEmpty(third)}, "/")
}

func fallbackEmpty(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

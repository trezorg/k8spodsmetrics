package metricsresources

import (
	"testing"

	escapes "github.com/snugfox/ansi-escapes"
	"github.com/stretchr/testify/require"
	servicemetricsresources "github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

func TestMetricsFormatterStringWithColor(t *testing.T) {
	resource := servicemetricsresources.MetricsResource{
		CPURequest:    100,
		CPUUsed:       150,
		MemoryRequest: 1024,
		MemoryUsed:    2048,
	}

	formatted := NewMetrics(resource).StringWithColor(escapes.TextColorRed)
	require.Contains(t, formatted, "CPU=100/")
	require.Contains(t, formatted, "Memory=1KiB/")
	require.Contains(t, formatted, escapes.TextColorRed)
	require.Contains(t, formatted, escapes.ColorReset)
}

func TestContainerFormatterUsedStrings(t *testing.T) {
	container := servicemetricsresources.ContainerMetricsResource{
		Requests: servicemetricsresources.MetricsResource{
			CPURequest:    100,
			CPUUsed:       120,
			MemoryRequest: 1024,
			MemoryUsed:    1500,
		},
		Limits: servicemetricsresources.MetricsResource{
			CPURequest:    200,
			CPUUsed:       120,
			MemoryRequest: 2048,
			MemoryUsed:    1500,
		},
	}

	formatter := NewContainer(container)
	require.Contains(t, formatter.CPUUsed(), "120")
	require.Contains(t, formatter.MemoryUsed(), "1.5KiB")
}

func TestContainerFormatterCompactStrings(t *testing.T) {
	container := servicemetricsresources.ContainerMetricsResource{
		Requests: servicemetricsresources.MetricsResource{
			CPURequest:              100,
			CPUUsed:                 120,
			MemoryRequest:           1024,
			MemoryUsed:              1500,
			StorageRequest:          2048,
			StorageUsed:             512,
			StorageEphemeralRequest: 4096,
			StorageEphemeralUsed:    1024,
		},
		Limits: servicemetricsresources.MetricsResource{
			CPURequest:              200,
			CPUUsed:                 120,
			MemoryRequest:           2048,
			MemoryUsed:              1500,
			StorageRequest:          4096,
			StorageUsed:             512,
			StorageEphemeralRequest: 8192,
			StorageEphemeralUsed:    1024,
		},
	}

	formatter := NewContainer(container)
	require.Contains(t, formatter.CPUCompactString(), "100/")
	require.Contains(t, formatter.CPUCompactString(), "/200")
	require.Contains(t, formatter.MemoryCompactString(), "1KiB/")
	require.Contains(t, formatter.MemoryCompactString(), "/2KiB")
	require.Equal(t, "2KiB/512B/4KiB", formatter.StorageCompactString())
	require.Equal(t, "4KiB/1KiB/8KiB", formatter.StorageEphemeralCompactString())
}

func TestContainerFormatterCompactStringsWithUnsetUsed(t *testing.T) {
	container := servicemetricsresources.ContainerMetricsResource{
		Requests: servicemetricsresources.MetricsResource{
			CPURequest:              100,
			CPUUsed:                 unset,
			MemoryRequest:           1024,
			MemoryUsed:              unset,
			StorageRequest:          2048,
			StorageUsed:             unset,
			StorageEphemeralRequest: 4096,
			StorageEphemeralUsed:    unset,
		},
		Limits: servicemetricsresources.MetricsResource{
			CPURequest:              200,
			CPUUsed:                 unset,
			MemoryRequest:           2048,
			MemoryUsed:              unset,
			StorageRequest:          4096,
			StorageUsed:             unset,
			StorageEphemeralRequest: 8192,
			StorageEphemeralUsed:    unset,
		},
	}

	formatter := NewContainer(container)
	require.Equal(t, "100/-/200", formatter.CPUCompactString())
	require.Equal(t, "1KiB/-/2KiB", formatter.MemoryCompactString())
	require.Equal(t, "2KiB/-/4KiB", formatter.StorageCompactString())
	require.Equal(t, "4KiB/-/8KiB", formatter.StorageEphemeralCompactString())
}

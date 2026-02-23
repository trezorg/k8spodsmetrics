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

package metricsresources

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	servicemetricsresources "github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func TestCompactHeaderRow(t *testing.T) {
	row := compactHeaderRow(resources.Resources{resources.All})
	require.Equal(t, []any{"NAMESPACE", "POD", "NODE", "CPU(req/used/lim)", "MEM(req/used/lim)", "STO(req/used/lim)", "EPH(req/used/lim)"}, []any(row))
}

func TestAggregatePodContainers(t *testing.T) {
	resource := testCompactPodResource()

	aggregated := aggregatePodContainers(resource)
	require.Equal(t, int64(300), aggregated.Requests.CPURequest)
	require.Equal(t, int64(700), aggregated.Limits.CPURequest)
	require.Equal(t, int64(270), aggregated.Requests.CPUUsed)
	require.Equal(t, int64(3840), aggregated.Requests.MemoryRequest)
	require.Equal(t, int64(6144), aggregated.Limits.MemoryRequest)
	require.Equal(t, int64(3072), aggregated.Requests.MemoryUsed)
}

func TestCompactPodRow(t *testing.T) {
	resource := testCompactPodResource()
	aggregated := aggregatePodContainers(resource)

	row := compactPodRow(resource, aggregated, resources.Resources{resources.CPU, resources.Memory})
	require.Equal(t, "default", row[0])
	require.Equal(t, "api-server", row[1])
	require.Equal(t, "node-a", row[2])
	require.Contains(t, row[3], "300/")
	require.Contains(t, row[3], "/700")
	require.Contains(t, row[4], "3.8KiB/")
	require.Contains(t, row[4], "/6KiB")
}

func TestPrintCompactToOmitsContainerRows(t *testing.T) {
	var buf bytes.Buffer
	PrintCompactTo(&buf, servicemetricsresources.PodMetricsResourceList{testCompactPodResource()}, resources.Resources{resources.CPU})

	output := buf.String()
	require.Contains(t, output, "CPU(REQ/USED/LIM)")
	require.Contains(t, output, "api-server")
	require.NotContains(t, output, "frontend")
	require.NotContains(t, output, "worker")
	require.NotContains(t, output, "TOTAL")
}

func TestPrintCompactToIncludesTotalFooter(t *testing.T) {
	var buf bytes.Buffer
	list := servicemetricsresources.PodMetricsResourceList{testCompactPodResource(), testSecondCompactPodResource()}
	PrintCompactTo(&buf, list, resources.Resources{resources.CPU})

	output := buf.String()
	require.Contains(t, output, "TOTAL")
	require.Contains(t, output, "350/310/800")
}

func TestPrintCompactToRespectsResources(t *testing.T) {
	var buf bytes.Buffer
	PrintCompactTo(&buf, servicemetricsresources.PodMetricsResourceList{testCompactPodResource()}, resources.Resources{resources.Memory})

	output := buf.String()
	require.NotContains(t, output, "CPU(REQ/USED/LIM)")
	require.Contains(t, output, "MEM(REQ/USED/LIM)")
	require.NotContains(t, output, "STO(REQ/USED/LIM)")
	require.NotContains(t, output, "EPH(REQ/USED/LIM)")
}

func TestPrintCompactToDoesNotTruncateIdentityColumns(t *testing.T) {
	resource := testCompactPodResource()
	resource.PodResource.Namespace = "very-long-namespace-for-compact-output"
	resource.PodResource.NamespaceName.Namespace = resource.PodResource.Namespace
	resource.PodResource.Name = "very-long-pod-name-that-should-stay-complete"
	resource.PodResource.NamespaceName.Name = resource.PodResource.Name
	resource.PodResource.NodeName = "very-long-node-name-that-should-stay-complete"

	var buf bytes.Buffer
	PrintCompactTo(&buf, servicemetricsresources.PodMetricsResourceList{resource}, resources.Resources{resources.CPU})

	output := buf.String()
	require.Contains(t, output, "very-long-namespace-for-compact-output")
	require.Contains(t, output, "very-long-pod-name-that-should-stay-complete")
	require.Contains(t, output, "very-long-node-name-that-should-stay-complete")
	require.NotContains(t, output, "...")
}

func testCompactPodResource() servicemetricsresources.PodMetricsResource {
	return servicemetricsresources.PodMetricsResource{
		PodResource: pods.PodResource{
			NamespaceName: pods.NamespaceName{Name: "api-server", Namespace: "default"},
			NodeName:      "node-a",
			Containers: []pods.ContainerResource{{
				Name:     "frontend",
				Requests: pods.Resource{CPU: 100, Memory: 1024, Storage: 2048, StorageEphemeral: 4096},
				Limits:   pods.Resource{CPU: 200, Memory: 2048, Storage: 4096, StorageEphemeral: 8192},
			}, {
				Name:     "worker",
				Requests: pods.Resource{CPU: 200, Memory: 2816, Storage: 1024, StorageEphemeral: 2048},
				Limits:   pods.Resource{CPU: 500, Memory: 4096, Storage: 2048, StorageEphemeral: 4096},
			}},
		},
		PodMetric: podmetrics.PodMetric{
			Name:      "api-server",
			Namespace: "default",
			Containers: []podmetrics.ContainerMetric{{
				Name:   "frontend",
				Metric: podmetrics.Metric{CPU: 120, Memory: 1536, Storage: 512, StorageEphemeral: 1024},
			}, {
				Name:   "worker",
				Metric: podmetrics.Metric{CPU: 150, Memory: 1536, Storage: 256, StorageEphemeral: 512},
			}},
		},
	}
}

func testSecondCompactPodResource() servicemetricsresources.PodMetricsResource {
	return servicemetricsresources.PodMetricsResource{
		PodResource: pods.PodResource{
			NamespaceName: pods.NamespaceName{Name: "dns", Namespace: "kube-system"},
			NodeName:      "node-b",
			Containers: []pods.ContainerResource{{
				Name:     "dns",
				Requests: pods.Resource{CPU: 50, Memory: 512},
				Limits:   pods.Resource{CPU: 100, Memory: 1024},
			}},
		},
		PodMetric: podmetrics.PodMetric{
			Name:      "dns",
			Namespace: "kube-system",
			Containers: []podmetrics.ContainerMetric{{
				Name:   "dns",
				Metric: podmetrics.Metric{CPU: 40, Memory: 256},
			}},
		},
	}
}

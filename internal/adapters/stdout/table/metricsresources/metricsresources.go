package metricsresources

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
	"log/slog"
)

type Table func(list metricsresources.PodMetricsResourceList)

func ToTable(
	outputResources resources.Resources,
) Table {
	return Table(func(list metricsresources.PodMetricsResourceList) {
		Print(list, outputResources)
	})
}

func headerFooter(outputResources resources.Resources, firstColumn string) table.Row {
	result := table.Row{firstColumn, "Namespace", "Node Name"}
	if outputResources.IsCPU() {
		result = append(
			result,
			"CPU",
			"CPU",
			"CPU",
		)
	}
	if outputResources.IsMemory() {
		result = append(
			result,
			"Memory",
			"Memory",
			"Memory",
		)
	}
	if outputResources.IsStorage() {
		result = append(
			result,
			"Storage",
			"Storage",
			"Storage",
			"Storage Ephemeral",
			"Storage Ephemeral",
			"Storage Ephemeral",
		)
	}
	return result
}

func secondaryHeader(outputResources resources.Resources) table.Row {
	result := table.Row{"", "", ""}
	if outputResources.IsCPU() {
		result = append(
			result,
			"Request",
			"Limit",
			"Used",
		)
	}
	if outputResources.IsMemory() {
		result = append(
			result,
			"Request",
			"Limit",
			"Used",
		)
	}
	if outputResources.IsStorage() {
		result = append(
			result,
			"Request",
			"Limit",
			"Used",
			"Request",
			"Limit",
			"Used",
		)
	}
	return result
}

func row(resource metricsresources.PodMetricsResource, outputResources resources.Resources) table.Row {
	result := table.Row{resource.Name, resource.Namespace, resource.NodeName}
	containers := resource.ContainersMetrics()
	if len(containers) == 0 {
		return result
	}

	// Use first container for pod-level metrics
	container := containers[0]

	if outputResources.IsCPU() {
		result = append(
			result,
			container.Requests.CPURequestString(),
			container.Limits.CPURequestString(),
			container.CPUUsed(),
		)
	}
	if outputResources.IsMemory() {
		result = append(
			result,
			container.Requests.MemoryRequestString(),
			container.Limits.MemoryRequestString(),
			container.MemoryUsed(),
		)
	}
	if outputResources.IsStorage() {
		result = append(
			result,
			container.Requests.StorageRequestString(),
			container.Limits.StorageRequestString(),
			container.StorageUsed(),
			container.Requests.StorageEphemeralRequestString(),
			container.Limits.StorageEphemeralRequestString(),
			container.StorageEphemeralUsed(),
		)
	}
	return result
}

func containerRow(container metricsresources.ContainerMetricsResource, outputResources resources.Resources) table.Row {
	result := table.Row{"", "", container.Name}

	if outputResources.IsCPU() {
		result = append(
			result,
			container.Requests.CPURequestString(),
			container.Limits.CPURequestString(),
			container.CPUUsed(),
		)
	}
	if outputResources.IsMemory() {
		result = append(
			result,
			container.Requests.MemoryRequestString(),
			container.Limits.MemoryRequestString(),
			container.MemoryUsed(),
		)
	}
	if outputResources.IsStorage() {
		result = append(
			result,
			container.Requests.StorageRequestString(),
			container.Limits.StorageRequestString(),
			container.StorageUsed(),
			container.Requests.StorageEphemeralRequestString(),
			container.Limits.StorageEphemeralRequestString(),
			container.StorageEphemeralUsed(),
		)
	}
	return result
}

func Print(
	list metricsresources.PodMetricsResourceList,
	outputResources resources.Resources,
) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendHeader(headerFooter(outputResources, "Pod Name"), rowConfigAutoMerge)
	t.AppendHeader(secondaryHeader(outputResources))

	total := metricsresources.ContainerMetricsResource{}

	for _, resource := range list {
		containers := resource.ContainersMetrics()
		if len(containers) == 0 {
			continue
		}

		// Add pod row with first container
		t.AppendRow(row(resource, outputResources))

		// Add additional container rows
		for _, container := range containers[1:] {
			t.AppendRow(containerRow(container, outputResources))
		}

		t.AppendSeparator()

		// Calculate totals
		for _, container := range containers {
			total.Requests.CPURequest += container.Requests.CPURequest
			total.Limits.CPURequest += container.Limits.CPURequest
			total.Requests.CPUUsed += container.Requests.CPUUsed
			total.Requests.MemoryRequest += container.Requests.MemoryRequest
			total.Limits.MemoryRequest += container.Limits.MemoryRequest
			total.Requests.MemoryUsed += container.Requests.MemoryUsed
			total.Requests.StorageUsed += container.Requests.StorageUsed
			total.Limits.StorageUsed += container.Limits.StorageUsed
			total.Requests.StorageEphemeralUsed += container.Requests.StorageEphemeralUsed
			total.Limits.StorageEphemeralUsed += container.Limits.StorageEphemeralUsed
			total.Requests.StorageRequest += container.Requests.StorageRequest
			total.Limits.StorageRequest += container.Limits.StorageRequest
			total.Requests.StorageEphemeralRequest += container.Requests.StorageEphemeralRequest
			total.Limits.StorageEphemeralRequest += container.Limits.StorageEphemeralRequest
		}
	}

	// Add footer with totals
	t.AppendFooter(headerFooter(outputResources, "Total"), rowConfigAutoMerge)
	totalRow := table.Row{"Total", "", ""}
	if outputResources.IsCPU() {
		totalRow = append(
			totalRow,
			total.Requests.CPURequestString(),
			total.Limits.CPURequestString(),
			total.Requests.CPUUsedString(""),
		)
	}
	if outputResources.IsMemory() {
		totalRow = append(
			totalRow,
			total.Requests.MemoryRequestString(),
			total.Limits.MemoryRequestString(),
			total.Requests.MemoryUsedString(""),
		)
	}
	if outputResources.IsStorage() {
		totalRow = append(
			totalRow,
			total.Requests.StorageString(),
			total.Limits.StorageString(),
			total.Requests.StorageString(),
			total.Requests.StorageEphemeralString(),
			total.Limits.StorageEphemeralString(),
			total.Requests.StorageEphemeralString(),
		)
	}
	t.AppendFooter(totalRow)
	t.Render()
}

func (s Table) Success(list metricsresources.PodMetricsResourceList) {
	s(list)
}

func (Table) Error(err error) {
	slog.Error("", slog.Any("error", err))
}

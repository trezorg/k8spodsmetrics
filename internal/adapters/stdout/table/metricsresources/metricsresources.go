package metricsresources

import (
	"os"

	"log/slog"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/trezorg/k8spodsmetrics/internal/columns"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
)

type Table func(list metricsresources.PodMetricsResourceList)

type ColumnSet struct {
	Request bool
	Limit   bool
	Used    bool
}

func newColumnSet(cols []columns.Column) ColumnSet {
	if len(cols) == 0 {
		return ColumnSet{Request: true, Limit: true, Used: true}
	}
	cs := ColumnSet{}
	for _, col := range cols {
		//nolint:exhaustive // Total, Allocatable, Available, Free are node-only columns
		switch col {
		case columns.Request:
			cs.Request = true
		case columns.Limit:
			cs.Limit = true
		case columns.Used:
			cs.Used = true
		default:
			// Node-only columns (Total, Allocatable, Available, Free) ignored for pods
		}
	}
	return cs
}

func (cs ColumnSet) appendResourceHeaderRow(result table.Row, label string) table.Row {
	if cs.Request {
		result = append(result, label)
	}
	if cs.Limit {
		result = append(result, label)
	}
	if cs.Used {
		result = append(result, label)
	}
	return result
}

func (cs ColumnSet) appendStorageHeaderRow(result table.Row) table.Row {
	if cs.Request {
		result = append(result, "Storage")
	}
	if cs.Limit {
		result = append(result, "Storage")
	}
	if cs.Used {
		result = append(result, "Storage")
	}
	if cs.Request {
		result = append(result, "Storage Ephemeral")
	}
	if cs.Limit {
		result = append(result, "Storage Ephemeral")
	}
	if cs.Used {
		result = append(result, "Storage Ephemeral")
	}
	return result
}

func (cs ColumnSet) headerFooterRow(outputResources resources.Resources, firstColumn string) table.Row {
	result := table.Row{firstColumn, "", ""}
	if outputResources.IsCPU() {
		result = cs.appendResourceHeaderRow(result, "CPU")
	}
	if outputResources.IsMemory() {
		result = cs.appendResourceHeaderRow(result, "Memory")
	}
	if outputResources.IsStorage() {
		result = cs.appendStorageHeaderRow(result)
	}
	return result
}

func (cs ColumnSet) appendSecondaryHeaderRow(result table.Row) table.Row {
	if cs.Request {
		result = append(result, "Request")
	}
	if cs.Limit {
		result = append(result, "Limit")
	}
	if cs.Used {
		result = append(result, "Used")
	}
	return result
}

func (cs ColumnSet) appendStorageSecondaryHeaderRow(result table.Row) table.Row {
	if cs.Request {
		result = append(result, "Request")
	}
	if cs.Limit {
		result = append(result, "Limit")
	}
	if cs.Used {
		result = append(result, "Used")
	}
	if cs.Request {
		result = append(result, "Request")
	}
	if cs.Limit {
		result = append(result, "Limit")
	}
	if cs.Used {
		result = append(result, "Used")
	}
	return result
}

func (cs ColumnSet) secondaryHeaderRow(outputResources resources.Resources) table.Row {
	result := table.Row{"", "", ""}
	if outputResources.IsCPU() {
		result = cs.appendSecondaryHeaderRow(result)
	}
	if outputResources.IsMemory() {
		result = cs.appendSecondaryHeaderRow(result)
	}
	if outputResources.IsStorage() {
		result = cs.appendStorageSecondaryHeaderRow(result)
	}
	return result
}

func (cs ColumnSet) appendCPUColumns(result table.Row, container metricsresources.ContainerMetricsResource) table.Row {
	if cs.Request {
		result = append(result, container.Requests.CPURequestString())
	}
	if cs.Limit {
		result = append(result, container.Limits.CPURequestString())
	}
	if cs.Used {
		result = append(result, container.CPUUsed())
	}
	return result
}

func (cs ColumnSet) appendMemoryColumns(result table.Row, container metricsresources.ContainerMetricsResource) table.Row {
	if cs.Request {
		result = append(result, container.Requests.MemoryRequestString())
	}
	if cs.Limit {
		result = append(result, container.Limits.MemoryRequestString())
	}
	if cs.Used {
		result = append(result, container.MemoryUsed())
	}
	return result
}

func (cs ColumnSet) appendStorageColumns(result table.Row, container metricsresources.ContainerMetricsResource) table.Row {
	if cs.Request {
		result = append(result, container.Requests.StorageRequestString())
	}
	if cs.Limit {
		result = append(result, container.Limits.StorageRequestString())
	}
	if cs.Used {
		result = append(result, container.StorageUsed())
	}
	if cs.Request {
		result = append(result, container.Requests.StorageEphemeralRequestString())
	}
	if cs.Limit {
		result = append(result, container.Limits.StorageEphemeralRequestString())
	}
	if cs.Used {
		result = append(result, container.StorageEphemeralUsed())
	}
	return result
}

func (cs ColumnSet) dataRow(resource metricsresources.PodMetricsResource, outputResources resources.Resources) table.Row {
	result := table.Row{resource.Name, resource.Namespace, resource.NodeName}
	containers := resource.ContainersMetrics()
	if len(containers) == 0 {
		return result
	}
	container := containers[0]

	if outputResources.IsCPU() {
		for _, cn := range containers[1:] {
			container.Requests.CPURequest += cn.Requests.CPURequest
			container.Limits.CPURequest += cn.Limits.CPURequest
			container.Requests.CPUUsed += cn.Requests.CPUUsed
			container.Limits.CPUUsed += cn.Limits.CPUUsed
		}
		result = cs.appendCPUColumns(result, container)
	}
	if outputResources.IsMemory() {
		for _, cn := range containers[1:] {
			container.Requests.MemoryRequest += cn.Requests.MemoryRequest
			container.Limits.MemoryRequest += cn.Limits.MemoryRequest
			container.Requests.MemoryUsed += cn.Requests.MemoryUsed
			container.Limits.MemoryUsed += cn.Limits.MemoryUsed
		}
		result = cs.appendMemoryColumns(result, container)
	}
	if outputResources.IsStorage() {
		for _, cn := range containers[1:] {
			container.Requests.StorageRequest += cn.Requests.StorageRequest
			container.Limits.StorageRequest += cn.Limits.StorageRequest
			container.Requests.StorageUsed += cn.Requests.StorageUsed
			container.Limits.StorageUsed += cn.Limits.StorageUsed
			container.Requests.StorageEphemeralRequest += cn.Requests.StorageEphemeralRequest
			container.Limits.StorageEphemeralRequest += cn.Limits.StorageEphemeralRequest
			container.Requests.StorageEphemeralUsed += cn.Requests.StorageEphemeralUsed
			container.Limits.StorageEphemeralUsed += cn.Limits.StorageEphemeralUsed
		}
		result = cs.appendStorageColumns(result, container)
	}
	return result
}

func (cs ColumnSet) containerRow(container metricsresources.ContainerMetricsResource, outputResources resources.Resources) table.Row {
	result := table.Row{container.Name, "", ""}

	if outputResources.IsCPU() {
		result = cs.appendCPUColumns(result, container)
	}
	if outputResources.IsMemory() {
		result = cs.appendMemoryColumns(result, container)
	}
	if outputResources.IsStorage() {
		result = cs.appendStorageColumns(result, container)
	}
	return result
}

func (cs ColumnSet) totalRow(outputResources resources.Resources, total metricsresources.ContainerMetricsResource) table.Row {
	totalRow := table.Row{"", "", ""}
	if outputResources.IsCPU() {
		if cs.Request {
			totalRow = append(totalRow, total.Requests.CPURequestString())
		}
		if cs.Limit {
			totalRow = append(totalRow, total.Limits.CPURequestString())
		}
		if cs.Used {
			totalRow = append(totalRow, total.Requests.CPUUsedString(""))
		}
	}
	if outputResources.IsMemory() {
		if cs.Request {
			totalRow = append(totalRow, total.Requests.MemoryRequestString())
		}
		if cs.Limit {
			totalRow = append(totalRow, total.Limits.MemoryRequestString())
		}
		if cs.Used {
			totalRow = append(totalRow, total.Requests.MemoryUsedString(""))
		}
	}
	if outputResources.IsStorage() {
		if cs.Request {
			totalRow = append(totalRow, total.Requests.StorageRequestString())
		}
		if cs.Limit {
			totalRow = append(totalRow, total.Limits.StorageRequestString())
		}
		if cs.Used {
			totalRow = append(totalRow, total.Requests.StorageString())
		}
		if cs.Request {
			totalRow = append(totalRow, total.Requests.StorageEphemeralRequestString())
		}
		if cs.Limit {
			totalRow = append(totalRow, total.Limits.StorageEphemeralRequestString())
		}
		if cs.Used {
			totalRow = append(totalRow, total.Requests.StorageEphemeralString())
		}
	}
	return totalRow
}

func ToTable(
	outputResources resources.Resources,
	cols []columns.Column,
) Table {
	cs := newColumnSet(cols)
	return Table(func(list metricsresources.PodMetricsResourceList) {
		Print(list, outputResources, cs)
	})
}

func headerFooter(outputResources resources.Resources, firstColumn string) table.Row {
	result := table.Row{firstColumn, "", ""}
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
	container := containers[0]

	if outputResources.IsCPU() {
		for _, cn := range containers[1:] {
			container.Requests.CPURequest += cn.Requests.CPURequest
			container.Limits.CPURequest += cn.Limits.CPURequest
			container.Requests.CPUUsed += cn.Requests.CPUUsed
			container.Limits.CPUUsed += cn.Limits.CPUUsed
		}
		result = append(
			result,
			container.Requests.CPURequestString(),
			container.Limits.CPURequestString(),
			container.CPUUsed(),
		)
	}
	if outputResources.IsMemory() {
		for _, cn := range containers[1:] {
			container.Requests.MemoryRequest += cn.Requests.MemoryRequest
			container.Limits.MemoryRequest += cn.Limits.MemoryRequest
			container.Requests.MemoryUsed += cn.Requests.MemoryUsed
			container.Limits.MemoryUsed += cn.Limits.MemoryUsed
		}
		result = append(
			result,
			container.Requests.MemoryRequestString(),
			container.Limits.MemoryRequestString(),
			container.MemoryUsed(),
		)
	}
	if outputResources.IsStorage() {
		for _, cn := range containers[1:] {
			container.Requests.StorageRequest += cn.Requests.StorageRequest
			container.Limits.StorageRequest += cn.Limits.StorageRequest
			container.Requests.StorageUsed += cn.Requests.StorageUsed
			container.Limits.StorageUsed += cn.Limits.StorageUsed
		}
		for _, cn := range containers[1:] {
			container.Requests.StorageEphemeralRequest += cn.Requests.StorageEphemeralRequest
			container.Limits.StorageEphemeralRequest += cn.Limits.StorageEphemeralRequest
			container.Requests.StorageEphemeralUsed += cn.Requests.StorageEphemeralUsed
			container.Limits.StorageEphemeralUsed += cn.Limits.StorageEphemeralUsed
		}
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
	result := table.Row{container.Name, "", ""}

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
	cs ColumnSet,
) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendHeader(cs.headerFooterRow(outputResources, "Pod Name / Container Names"), rowConfigAutoMerge)
	t.AppendHeader(cs.secondaryHeaderRow(outputResources))

	total := metricsresources.ContainerMetricsResource{}

	for _, resource := range list {
		containers := resource.ContainersMetrics()
		if len(containers) == 0 {
			continue
		}

		// Add pod row with first container
		t.AppendRow(cs.dataRow(resource, outputResources))

		// Add additional container rows
		for _, container := range containers {
			t.AppendRow(cs.containerRow(container, outputResources))
		}

		t.AppendSeparator()

		// Calculate totals
		for _, container := range containers {
			total.Requests.CPURequest += container.Requests.CPURequest
			total.Requests.MemoryRequest += container.Requests.MemoryRequest
			total.Requests.StorageRequest += container.Requests.StorageRequest
			total.Requests.StorageEphemeralRequest += container.Requests.StorageEphemeralRequest

			total.Limits.CPURequest += container.Limits.CPURequest
			total.Limits.MemoryRequest += container.Limits.MemoryRequest
			total.Limits.StorageRequest += container.Limits.StorageRequest
			total.Limits.StorageEphemeralRequest += container.Limits.StorageEphemeralRequest

			total.Requests.CPUUsed += container.Requests.CPUUsed
			total.Requests.MemoryUsed += container.Requests.MemoryUsed
			total.Requests.StorageUsed += container.Requests.StorageUsed
			total.Requests.StorageEphemeralUsed += container.Requests.StorageEphemeralUsed
		}
	}

	// Add footer with totals
	t.AppendFooter(cs.headerFooterRow(outputResources, "Total"), rowConfigAutoMerge)
	t.AppendFooter(cs.secondaryHeaderRow(outputResources))
	t.AppendFooter(cs.totalRow(outputResources, total))
	t.Render()
}

func (s Table) Success(list metricsresources.PodMetricsResourceList) {
	s(list)
}

func (Table) Error(err error) {
	slog.Error("table metrics resources output failed", "error", err)
}

// ValidateColumns validates that all columns are valid for pod output.
func ValidateColumns(cols []columns.Column) error {
	return columns.ValidForPods(cols)
}

// ParseColumns converts string slice to Column slice, removing duplicates.
func ParseColumns(cols []string) []columns.Column {
	result := make([]columns.Column, 0, len(cols))
	seen := make(map[columns.Column]bool)
	for _, c := range cols {
		col := columns.Column(c)
		if !seen[col] {
			seen[col] = true
			result = append(result, col)
		}
	}
	return result
}

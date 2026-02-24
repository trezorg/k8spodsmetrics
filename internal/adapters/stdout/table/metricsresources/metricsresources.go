package metricsresources

import (
	"io"
	"os"

	"log/slog"

	"github.com/jedib0t/go-pretty/v6/table"
	formatmetricsresources "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/format/metricsresources"
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

func (cs ColumnSet) headerFooterRow(outputResources resources.Resources, columnNames ...string) table.Row {
	maxColumsNames := 3
	if len(columnNames) > maxColumsNames {
		columnNames = columnNames[:3]
	}
	result := table.Row{}
	for _, column := range columnNames {
		result = append(result, column)
	}
	for range maxColumsNames - len(columnNames) {
		result = append(result, "")
	}
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
	containerFormatter := formatmetricsresources.NewContainer(container)
	if cs.Request {
		result = append(result, containerFormatter.Requests().CPURequestString())
	}
	if cs.Limit {
		result = append(result, containerFormatter.Limits().CPURequestString())
	}
	if cs.Used {
		result = append(result, containerFormatter.CPUUsed())
	}
	return result
}

func (cs ColumnSet) appendMemoryColumns(result table.Row, container metricsresources.ContainerMetricsResource) table.Row {
	containerFormatter := formatmetricsresources.NewContainer(container)
	if cs.Request {
		result = append(result, containerFormatter.Requests().MemoryRequestString())
	}
	if cs.Limit {
		result = append(result, containerFormatter.Limits().MemoryRequestString())
	}
	if cs.Used {
		result = append(result, containerFormatter.MemoryUsed())
	}
	return result
}

func (cs ColumnSet) appendStorageColumns(result table.Row, container metricsresources.ContainerMetricsResource) table.Row {
	containerFormatter := formatmetricsresources.NewContainer(container)
	if cs.Request {
		result = append(result, containerFormatter.Requests().StorageRequestString())
	}
	if cs.Limit {
		result = append(result, containerFormatter.Limits().StorageRequestString())
	}
	if cs.Used {
		result = append(result, containerFormatter.StorageUsed())
	}
	if cs.Request {
		result = append(result, containerFormatter.Requests().StorageEphemeralRequestString())
	}
	if cs.Limit {
		result = append(result, containerFormatter.Limits().StorageEphemeralRequestString())
	}
	if cs.Used {
		result = append(result, containerFormatter.StorageEphemeralUsed())
	}
	return result
}

func (cs ColumnSet) dataRow(resource metricsresources.PodMetricsResource, outputResources resources.Resources) table.Row {
	result := table.Row{resource.PodResource.Name, resource.PodResource.Namespace, resource.NodeName}
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
	result := table.Row{"└─ " + container.Name, "", ""}

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
	totalRow := table.Row{"Total", "Total", "Total"}
	if outputResources.IsCPU() {
		if cs.Request {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Requests).CPURequestString())
		}
		if cs.Limit {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Limits).CPURequestString())
		}
		if cs.Used {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Requests).CPUUsedString(""))
		}
	}
	if outputResources.IsMemory() {
		if cs.Request {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Requests).MemoryRequestString())
		}
		if cs.Limit {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Limits).MemoryRequestString())
		}
		if cs.Used {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Requests).MemoryUsedString(""))
		}
	}
	if outputResources.IsStorage() {
		if cs.Request {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Requests).StorageRequestString())
		}
		if cs.Limit {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Limits).StorageRequestString())
		}
		if cs.Used {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Requests).StorageString())
		}
		if cs.Request {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Requests).StorageEphemeralRequestString())
		}
		if cs.Limit {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Limits).StorageEphemeralRequestString())
		}
		if cs.Used {
			totalRow = append(totalRow, formatmetricsresources.NewMetrics(total.Requests).StorageEphemeralString())
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
		PrintTo(os.Stdout, list, outputResources, cs)
	})
}

func ToWriter(
	outputResources resources.Resources,
	cols []columns.Column,
) func(io.Writer, metricsresources.PodMetricsResourceList) {
	cs := newColumnSet(cols)
	return func(w io.Writer, list metricsresources.PodMetricsResourceList) {
		PrintTo(w, list, outputResources, cs)
	}
}

func Print(
	list metricsresources.PodMetricsResourceList,
	outputResources resources.Resources,
	cs ColumnSet,
) {
	PrintTo(os.Stdout, list, outputResources, cs)
}

func PrintTo(
	w io.Writer,
	list metricsresources.PodMetricsResourceList,
	outputResources resources.Resources,
	cs ColumnSet,
) {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendHeader(cs.headerFooterRow(outputResources, "Pod/Container", "Namespace", "Node"), rowConfigAutoMerge)
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
	t.AppendFooter(cs.totalRow(outputResources, total), rowConfigAutoMerge)
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

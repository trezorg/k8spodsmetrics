package noderesources

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	formatnoderesources "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/format/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/columns"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
	"log/slog"
)

type Table func(
	list noderesources.NodeResourceList,
)

type ColumnSet struct {
	Total       bool
	Allocatable bool
	Used        bool
	Request     bool
	Limit       bool
	Available   bool
	Free        bool
}

func newColumnSet(cols []columns.Column) ColumnSet {
	if len(cols) == 0 {
		return ColumnSet{Total: true, Allocatable: true, Used: true, Request: true, Limit: true, Available: true, Free: true}
	}
	cs := ColumnSet{}
	for _, col := range cols {
		switch col {
		case columns.Total:
			cs.Total = true
		case columns.Allocatable:
			cs.Allocatable = true
		case columns.Used:
			cs.Used = true
		case columns.Request:
			cs.Request = true
		case columns.Limit:
			cs.Limit = true
		case columns.Available:
			cs.Available = true
		case columns.Free:
			cs.Free = true
		default:
			// Ignore invalid columns (validation happens elsewhere)
		}
	}
	return cs
}

func (cs ColumnSet) appendResourceHeader(result table.Row, label string) table.Row {
	if cs.Total {
		result = append(result, label)
	}
	if cs.Allocatable {
		result = append(result, label)
	}
	if cs.Used {
		result = append(result, label)
	}
	if cs.Request {
		result = append(result, label)
	}
	if cs.Limit {
		result = append(result, label)
	}
	if cs.Available {
		result = append(result, label)
	}
	if cs.Free {
		result = append(result, label)
	}
	return result
}

func (cs ColumnSet) appendStorageHeader(result table.Row) table.Row {
	if cs.Total {
		result = append(result, "Storage")
	}
	if cs.Allocatable {
		result = append(result, "Storage")
	}
	if cs.Used {
		result = append(result, "Storage")
	}
	if cs.Free {
		result = append(result, "Storage")
	}
	if cs.Total {
		result = append(result, "Storage Ephemeral")
	}
	if cs.Allocatable {
		result = append(result, "Storage Ephemeral")
	}
	if cs.Used {
		result = append(result, "Storage Ephemeral")
	}
	if cs.Free {
		result = append(result, "Storage Ephemeral")
	}
	return result
}

func (cs ColumnSet) headerFooterRow(outputResources resources.Resources, firstColumn string) table.Row {
	result := table.Row{firstColumn}
	if outputResources.IsCPU() {
		result = cs.appendResourceHeader(result, "CPU")
	}
	if outputResources.IsMemory() {
		result = cs.appendResourceHeader(result, "Memory")
	}
	if outputResources.IsStorage() {
		result = cs.appendStorageHeader(result)
	}
	return result
}

func (cs ColumnSet) appendSecondaryHeader(result table.Row) table.Row {
	if cs.Total {
		result = append(result, "Total")
	}
	if cs.Allocatable {
		result = append(result, "Allocatable")
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
	if cs.Available {
		result = append(result, "Available")
	}
	if cs.Free {
		result = append(result, "Free")
	}
	return result
}

func (cs ColumnSet) appendStorageSecondaryHeader(result table.Row) table.Row {
	if cs.Total {
		result = append(result, "Total")
	}
	if cs.Allocatable {
		result = append(result, "Allocatable")
	}
	if cs.Used {
		result = append(result, "Used")
	}
	if cs.Free {
		result = append(result, "Free")
	}
	if cs.Total {
		result = append(result, "Total")
	}
	if cs.Allocatable {
		result = append(result, "Allocatable")
	}
	if cs.Used {
		result = append(result, "Used")
	}
	if cs.Free {
		result = append(result, "Free")
	}
	return result
}

func (cs ColumnSet) secondaryHeaderRow(outputResources resources.Resources) table.Row {
	result := table.Row{""}
	if outputResources.IsCPU() {
		result = cs.appendSecondaryHeader(result)
	}
	if outputResources.IsMemory() {
		result = cs.appendSecondaryHeader(result)
	}
	if outputResources.IsStorage() {
		result = cs.appendStorageSecondaryHeader(result)
	}
	return result
}

func (cs ColumnSet) appendCPUColumns(result table.Row, resource noderesources.NodeResource) table.Row {
	formatter := formatnoderesources.New(resource)
	if cs.Total {
		result = append(result, resource.CPU)
	}
	if cs.Allocatable {
		result = append(result, resource.AllocatableCPU)
	}
	if cs.Used {
		result = append(result, resource.UsedCPU)
	}
	if cs.Request {
		result = append(result, formatter.CPURequestString())
	}
	if cs.Limit {
		result = append(result, formatter.CPULimitString())
	}
	if cs.Available {
		result = append(result, formatter.CPUAvailableString())
	}
	if cs.Free {
		result = append(result, formatter.CPUFreeString())
	}
	return result
}

func (cs ColumnSet) appendMemoryColumns(result table.Row, resource noderesources.NodeResource) table.Row {
	formatter := formatnoderesources.New(resource)
	if cs.Total {
		result = append(result, formatter.MemoryNodeString())
	}
	if cs.Allocatable {
		result = append(result, formatter.MemoryNodeAllocatableString())
	}
	if cs.Used {
		result = append(result, formatter.MemoryNodeUsedString())
	}
	if cs.Request {
		result = append(result, formatter.MemoryRequestString())
	}
	if cs.Limit {
		result = append(result, formatter.MemoryLimitString())
	}
	if cs.Available {
		result = append(result, formatter.MemoryAvailableString())
	}
	if cs.Free {
		result = append(result, formatter.MemoryFreeString())
	}
	return result
}

func (cs ColumnSet) appendStorageColumns(result table.Row, resource noderesources.NodeResource) table.Row {
	formatter := formatnoderesources.New(resource)
	if cs.Total {
		result = append(result, formatter.StorageString())
	}
	if cs.Allocatable {
		result = append(result, formatter.StorageAllocatableString())
	}
	if cs.Used {
		result = append(result, formatter.StorageUsedString())
	}
	if cs.Free {
		result = append(result, formatter.StorageFreeString())
	}
	if cs.Total {
		result = append(result, formatter.StorageEphemeralString())
	}
	if cs.Allocatable {
		result = append(result, formatter.StorageAllocatableEphemeralString())
	}
	if cs.Used {
		result = append(result, formatter.StorageUsedEphemeralString())
	}
	if cs.Free {
		result = append(result, formatter.StorageFreeEphemeralString())
	}
	return result
}

func (cs ColumnSet) dataRow(resource noderesources.NodeResource, outputResources resources.Resources) table.Row {
	result := table.Row{resource.Name}
	if outputResources.IsCPU() {
		result = cs.appendCPUColumns(result, resource)
	}
	if outputResources.IsMemory() {
		result = cs.appendMemoryColumns(result, resource)
	}
	if outputResources.IsStorage() {
		result = cs.appendStorageColumns(result, resource)
	}
	return result
}

func ToTable(
	outputResources resources.Resources,
	cols []columns.Column,
) Table {
	cs := newColumnSet(cols)
	return Table(func(list noderesources.NodeResourceList) {
		Print(list, outputResources, cs)
	})
}

func Print(
	list noderesources.NodeResourceList,
	outputResources resources.Resources,
	cs ColumnSet,
) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendHeader(cs.headerFooterRow(outputResources, "Name"), rowConfigAutoMerge)
	t.AppendHeader(cs.secondaryHeaderRow(outputResources))
	total := noderesources.NodeResource{}
	for _, resource := range list {
		t.AppendRow(cs.dataRow(resource, outputResources))
		t.AppendSeparator()
		if cs.Total {
			total.CPU += resource.CPU
			total.Memory += resource.Memory
			total.Storage += resource.Storage
			total.StorageEphemeral += resource.StorageEphemeral
		}
		if cs.Allocatable {
			total.AllocatableCPU += resource.AllocatableCPU
			total.AllocatableMemory += resource.AllocatableMemory
			total.AllocatableStorage += resource.AllocatableStorage
			total.AllocatableStorageEphemeral += resource.AllocatableStorageEphemeral
		}
		if cs.Used {
			total.UsedCPU += resource.UsedCPU
			total.UsedMemory += resource.UsedMemory
			total.UsedStorage += resource.UsedStorage
			total.UsedStorageEphemeral += resource.UsedStorageEphemeral
		}
		if cs.Request {
			total.CPURequest += resource.CPURequest
			total.MemoryRequest += resource.MemoryRequest
		}
		if cs.Limit {
			total.CPULimit += resource.CPULimit
			total.MemoryLimit += resource.MemoryLimit
		}
		if cs.Available {
			total.AvailableCPU += resource.AvailableCPU
			total.AvailableMemory += resource.AvailableMemory
		}
		if cs.Free {
			total.FreeCPU += resource.FreeCPU
			total.FreeMemory += resource.FreeMemory
			total.FreeStorage += resource.FreeStorage
			total.FreeStorageEphemeral += resource.FreeStorageEphemeral
		}
	}
	t.AppendFooter(cs.headerFooterRow(outputResources, "Total"), rowConfigAutoMerge)
	t.AppendFooter(cs.dataRow(total, outputResources))
	t.Render()
}

func (s Table) Success(list noderesources.NodeResourceList) {
	s(list)
}

func (Table) Error(err error) {
	slog.Error("table node resources output failed", "error", err)
}

// ValidateColumns validates that all columns are valid for node output.
func ValidateColumns(cols []columns.Column) error {
	return columns.ValidForNodes(cols)
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

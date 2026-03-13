package noderesources

import (
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	formatnoderesources "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/format/noderesources"
	servicenoderesources "github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
)

const (
	compactNameColumn   = 1
	compactFirstMetric  = 2
	compactSecondMetric = 3
	compactThirdMetric  = 4
	compactFourthMetric = 5
	compactFifthMetric  = 6
	maxCompactColumns   = 7
)

func ToCompactTable(outputResources resources.Resources) Table {
	return Table(func(list servicenoderesources.NodeResourceList) {
		PrintCompactTo(os.Stdout, list, outputResources)
	})
}

func ToCompactWriter(outputResources resources.Resources) func(io.Writer, servicenoderesources.NodeResourceList) {
	return func(w io.Writer, list servicenoderesources.NodeResourceList) {
		PrintCompactTo(w, list, outputResources)
	}
}

func PrintCompactTo(w io.Writer, list servicenoderesources.NodeResourceList, outputResources resources.Resources) {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	configureCompactTable(t)
	t.AppendHeader(compactHeaderRow(outputResources))

	var total servicenoderesources.NodeResource
	rendered := 0
	for _, resource := range list {
		t.AppendRow(compactNodeRow(resource, outputResources))
		accumulateTotal(&total, resource)
		rendered++
	}

	if rendered > 1 {
		t.AppendFooter(compactTotalRow(total, outputResources))
	}

	t.Render()
}

func compactHeaderRow(outputResources resources.Resources) table.Row {
	row := table.Row{"NAME"}
	if outputResources.IsCPU() {
		row = append(row, "CPU(alloc/used/free)", "CPU(req/lim)")
	}
	if outputResources.IsMemory() {
		row = append(row, "MEM(alloc/used/free)", "MEM(req/lim)")
	}
	if outputResources.IsStorage() {
		row = append(row, "STO(alloc/used/free)", "EPH(alloc/used/free)")
	}
	return row
}

func compactNodeRow(resource servicenoderesources.NodeResource, outputResources resources.Resources) table.Row {
	formatter := formatnoderesources.New(resource)
	row := table.Row{resource.Name}
	if outputResources.IsCPU() {
		row = append(row, formatter.CPUCapacityCompactString(), formatter.CPUDemandCompactString())
	}
	if outputResources.IsMemory() {
		row = append(row, formatter.MemoryCapacityCompactString(), formatter.MemoryDemandCompactString())
	}
	if outputResources.IsStorage() {
		row = append(row, formatter.StorageCapacityCompactString(), formatter.StorageEphemeralCapacityCompactString())
	}
	return row
}

func compactTotalRow(total servicenoderesources.NodeResource, outputResources resources.Resources) table.Row {
	formatter := formatnoderesources.New(total)
	row := table.Row{"TOTAL"}
	if outputResources.IsCPU() {
		row = append(row, formatter.CPUCapacityCompactString(), formatter.CPUDemandCompactString())
	}
	if outputResources.IsMemory() {
		row = append(row, formatter.MemoryCapacityCompactString(), formatter.MemoryDemandCompactString())
	}
	if outputResources.IsStorage() {
		row = append(row, formatter.StorageCapacityCompactString(), formatter.StorageEphemeralCapacityCompactString())
	}
	return row
}

func accumulateTotal(total *servicenoderesources.NodeResource, resource servicenoderesources.NodeResource) {
	total.CPU += resource.CPU
	total.Memory += resource.Memory
	total.UsedCPU += resource.UsedCPU
	total.UsedMemory += resource.UsedMemory
	total.AllocatableCPU += resource.AllocatableCPU
	total.AllocatableMemory += resource.AllocatableMemory
	total.CPURequest += resource.CPURequest
	total.MemoryRequest += resource.MemoryRequest
	total.CPULimit += resource.CPULimit
	total.MemoryLimit += resource.MemoryLimit
	total.AvailableCPU += resource.AvailableCPU
	total.AvailableMemory += resource.AvailableMemory
	total.FreeCPU += resource.FreeCPU
	total.FreeMemory += resource.FreeMemory
	total.Storage += resource.Storage
	total.AllocatableStorage += resource.AllocatableStorage
	total.UsedStorage += resource.UsedStorage
	total.FreeStorage += resource.FreeStorage
	total.StorageEphemeral += resource.StorageEphemeral
	total.AllocatableStorageEphemeral += resource.AllocatableStorageEphemeral
	total.UsedStorageEphemeral += resource.UsedStorageEphemeral
	total.FreeStorageEphemeral += resource.FreeStorageEphemeral
}

func configureCompactTable(t table.Writer) {
	applyTableStyle(t)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: compactNameColumn, Align: text.AlignLeft},
		{Number: compactFirstMetric, Align: text.AlignRight},
		{Number: compactSecondMetric, Align: text.AlignRight},
		{Number: compactThirdMetric, Align: text.AlignRight},
		{Number: compactFourthMetric, Align: text.AlignRight},
		{Number: compactFifthMetric, Align: text.AlignRight},
		{Number: maxCompactColumns, Align: text.AlignRight},
	})
}

func applyTableStyle(t table.Writer) {
	t.SetStyle(table.StyleLight)
}

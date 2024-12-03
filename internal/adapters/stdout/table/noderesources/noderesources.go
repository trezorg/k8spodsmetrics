package noderesources

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type Table func(list noderesources.NodeResourceList)

func Print(list noderesources.NodeResourceList) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendHeader(table.Row{
		"Name",
		"CPU",
		"CPU",
		"CPU",
		"CPU",
		"CPU",
		"CPU",
		"CPU",
		"Memory",
		"Memory",
		"Memory",
		"Memory",
		"Memory",
		"Memory",
		"Memory",
	}, rowConfigAutoMerge)
	t.AppendHeader(table.Row{
		"",
		"Total",
		"Allocatable",
		"Used",
		"Request",
		"Limit",
		"Available",
		"Free",
		"Total",
		"Allocatable",
		"Used",
		"Request",
		"Limit",
		"Available",
		"Free",
	}, rowConfigAutoMerge)
	total := noderesources.NodeResource{}
	for _, resource := range list {
		t.AppendRow(table.Row{
			resource.Name,
			resource.CPU,
			resource.AllocatableCPU,
			resource.UsedCPU,
			resource.CPURequestString(),
			resource.CPULimitString(),
			resource.CPUAvailableString(),
			resource.CPUFreeString(),
			resource.MemoryNodeString(),
			resource.MemoryNodeAlocatableString(),
			resource.MemoryNodeUsedString(),
			resource.MemoryRequestString(),
			resource.MemoryLimitString(),
			resource.MemoryAvailableString(),
			resource.MemoryFreeString(),
		})
		t.AppendSeparator()
		total.CPU += resource.CPU
		total.AllocatableCPU += resource.AllocatableCPU
		total.UsedCPU += resource.UsedCPU
		total.CPURequest += resource.CPURequest
		total.CPULimit += resource.CPULimit
		total.AvailableCPU += resource.AvailableCPU
		total.FreeCPU += resource.FreeCPU
		total.Memory += resource.Memory
		total.AllocatableMemory += resource.AllocatableMemory
		total.UsedMemory += resource.UsedMemory
		total.MemoryLimit += resource.MemoryLimit
		total.MemoryRequest += resource.MemoryRequest
		total.AvailableMemory += resource.AvailableMemory
		total.FreeMemory += resource.FreeMemory
	}
	t.AppendFooter(table.Row{
		"Total",
		"CPU",
		"CPU",
		"CPU",
		"CPU",
		"CPU",
		"CPU",
		"CPU",
		"Memory",
		"Memory",
		"Memory",
		"Memory",
		"Memory",
		"Memory",
		"Memory",
	}, rowConfigAutoMerge)
	t.AppendFooter(table.Row{
		"",
		total.CPU,
		total.AllocatableCPU,
		total.UsedCPU,
		total.CPURequestString(),
		total.CPULimitString(),
		total.CPUAvailableString(),
		total.CPUFreeString(),
		total.MemoryNodeString(),
		total.MemoryNodeAlocatableString(),
		total.MemoryNodeUsedString(),
		total.MemoryRequestString(),
		total.MemoryLimitString(),
		total.MemoryAvailableString(),
		total.MemoryFreeString(),
	}, rowConfigAutoMerge)
	t.Render()
}

func (s Table) Success(list noderesources.NodeResourceList) {
	s(list)
}

func (s Table) Error(err error) {
	logger.Error("", err)
}

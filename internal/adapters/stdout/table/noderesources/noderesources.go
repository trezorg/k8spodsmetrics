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
	t.AppendHeader(table.Row{"Name", "CPU", "CPU", "CPU", "CPU", "CPU", "Memory", "Memory", "Memory", "Memory", "Memory"}, rowConfigAutoMerge)
	t.AppendHeader(table.Row{"", "Total", "Allocatable", "Used", "Request", "Limit", "Total", "Allocatable", "Used", "Request", "Limit"}, rowConfigAutoMerge)
	for _, resource := range list {
		t.AppendRow(table.Row{
			resource.Name,
			resource.CPU,
			resource.AllocatableCPU,
			resource.UsedCPU,
			resource.CPURequestString(),
			resource.CPULimitString(),
			resource.MemoryNodeString(),
			resource.MemoryNodeAlocatableString(),
			resource.MemoryNodeUsedString(),
			resource.MemoryRequestString(),
			resource.MemoryLimitString(),
		})
		t.AppendSeparator()
	}
	t.Render()
}

func (s Table) Success(list noderesources.NodeResourceList) {
	s(list)
}

func (s Table) Error(err error) {
	logger.Error("", err)
}

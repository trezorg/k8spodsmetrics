package noderesources

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
)

type Table func(
	list noderesources.NodeResourceList,
)

func ToTable(
	outputResources resources.Resources,
) Table {
	return Table(func(list noderesources.NodeResourceList) {
		Print(list, outputResources)
	})
}

func headerFooter(outputResources resources.Resources, firstColumn string) table.Row {
	result := table.Row{firstColumn}
	if outputResources.IsCPU() {
		result = append(
			result,
			"CPU",
			"CPU",
			"CPU",
			"CPU",
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
			"Memory",
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
			"Storage",
			"Storage Ephemeral",
			"Storage Ephemeral",
			"Storage Ephemeral",
			"Storage Ephemeral",
		)
	}
	return result
}

func secondaryHeader(outputResources resources.Resources) table.Row {
	result := table.Row{""}
	if outputResources.IsCPU() {
		result = append(
			result,
			"Total",
			"Allocatable",
			"Used",
			"Request",
			"Limit",
			"Available",
			"Free",
		)
	}
	if outputResources.IsMemory() {
		result = append(
			result,
			"Total",
			"Allocatable",
			"Used",
			"Request",
			"Limit",
			"Available",
			"Free",
		)
	}
	if outputResources.IsStorage() {
		result = append(
			result,
			"Total",
			"Allocatable",
			"Used",
			"Free",
			"Total",
			"Allocatable",
			"Used",
			"Free",
		)
	}
	return result
}

func row(resource noderesources.NodeResource, outputResources resources.Resources) table.Row {
	result := table.Row{resource.Name}
	if outputResources.IsCPU() {
		result = append(
			result,
			resource.CPU,
			resource.AllocatableCPU,
			resource.UsedCPU,
			resource.CPURequestString(),
			resource.CPULimitString(),
			resource.CPUAvailableString(),
			resource.CPUFreeString(),
		)
	}
	if outputResources.IsMemory() {
		result = append(
			result,
			resource.MemoryNodeString(),
			resource.MemoryNodeAlocatableString(),
			resource.MemoryNodeUsedString(),
			resource.MemoryRequestString(),
			resource.MemoryLimitString(),
			resource.MemoryAvailableString(),
			resource.MemoryFreeString(),
		)
	}
	if outputResources.IsStorage() {
		result = append(
			result,
			resource.StorageString(),
			resource.StorageAllocatableString(),
			resource.StorageUsedString(),
			resource.StorageFreeString(),
			resource.StorageEphemeralString(),
			resource.StorageAllocatableEphemeralString(),
			resource.StorageUsedEphemeralString(),
			resource.StorageFreeEphemeralString(),
		)
	}
	return result
}

func Print(
	list noderesources.NodeResourceList,
	outputResources resources.Resources,
) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendHeader(headerFooter(outputResources, "Name"), rowConfigAutoMerge)
	t.AppendHeader(secondaryHeader(outputResources))
	total := noderesources.NodeResource{}
	for _, resource := range list {
		t.AppendRow(row(resource, outputResources))
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
		total.Storage += resource.Storage
		total.UsedStorage += resource.UsedStorage
		total.AllocatableStorage += resource.AllocatableStorage
		total.FreeStorage += resource.FreeStorage
		total.StorageEphemeral += resource.StorageEphemeral
		total.UsedStorageEphemeral += resource.UsedStorageEphemeral
		total.AllocatableStorageEphemeral += resource.AllocatableStorageEphemeral
		total.FreeStorageEphemeral += resource.FreeStorageEphemeral
	}
	t.AppendFooter(headerFooter(outputResources, "Total"), rowConfigAutoMerge)
	t.AppendFooter(row(total, outputResources))
	t.Render()
}

func (s Table) Success(list noderesources.NodeResourceList) {
	s(list)
}

func (Table) Error(err error) {
	logger.Error("", err)
}

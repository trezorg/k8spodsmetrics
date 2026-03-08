package noderesources

import (
	"cmp"
	"slices"

	"github.com/trezorg/k8spodsmetrics/internal/sorting/noderesources"
)

func direction(reversed bool, result int) int {
	if reversed {
		return -result
	}
	return result
}

func sortBy[T cmp.Ordered](n NodeResourceList, reversed bool, field func(NodeResource) T) {
	slices.SortFunc(n, func(a, b NodeResource) int {
		return direction(reversed, cmp.Compare(field(a), field(b)))
	})
}

func (n NodeResourceList) sortByName(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) string { return resource.Name })
}

func (n NodeResourceList) sortRequestCPU(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.CPURequest })
}

func (n NodeResourceList) sortLimitCPU(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.CPULimit })
}

func (n NodeResourceList) sortUsedCPU(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.UsedCPU })
}

func (n NodeResourceList) sortAvailableCPU(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.AvailableCPU })
}

func (n NodeResourceList) sortFreeCPU(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.FreeCPU })
}

func (n NodeResourceList) sortCPU(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.CPU })
}

func (n NodeResourceList) sortRequestMemory(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.MemoryRequest })
}

func (n NodeResourceList) sortLimitMemory(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.MemoryLimit })
}

func (n NodeResourceList) sortUsedMemory(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.UsedMemory })
}

func (n NodeResourceList) sortAvailableMemory(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.AvailableMemory })
}

func (n NodeResourceList) sortFreeMemory(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.FreeMemory })
}

func (n NodeResourceList) sortMemory(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.Memory })
}

func (n NodeResourceList) sortStorage(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.Storage })
}

func (n NodeResourceList) sortStorageEphemeral(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.StorageEphemeral })
}

func (n NodeResourceList) sortAvailableStorage(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.AllocatableStorage })
}

func (n NodeResourceList) sortAvailableStorageEphemeral(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.AllocatableStorageEphemeral })
}

func (n NodeResourceList) sortUsedStorageEphemeral(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.UsedStorageEphemeral })
}

func (n NodeResourceList) sortUsedStorage(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.UsedStorage })
}

func (n NodeResourceList) sortFreeStorage(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.FreeStorage })
}

func (n NodeResourceList) sortFreeStorageEphemeral(reversed bool) {
	sortBy(n, reversed, func(resource NodeResource) int64 { return resource.FreeStorageEphemeral })
}

func (n NodeResourceList) sort(by string, reversed bool) { //nolint:revive // it is ok
	switch noderesources.Sorting(by) {
	case noderesources.Name:
		n.sortByName(reversed)
	case noderesources.LimitCPU:
		n.sortLimitCPU(reversed)
	case noderesources.RequestCPU:
		n.sortRequestCPU(reversed)
	case noderesources.UsedCPU:
		n.sortUsedCPU(reversed)
	case noderesources.TotalCPU:
		n.sortCPU(reversed)
	case noderesources.AvailableCPU:
		n.sortAvailableCPU(reversed)
	case noderesources.FreeCPU:
		n.sortFreeCPU(reversed)
	case noderesources.LimitMemory:
		n.sortLimitMemory(reversed)
	case noderesources.RequestMemory:
		n.sortRequestMemory(reversed)
	case noderesources.UsedMemory:
		n.sortUsedMemory(reversed)
	case noderesources.TotalMemory:
		n.sortMemory(reversed)
	case noderesources.AvailableMemory:
		n.sortAvailableMemory(reversed)
	case noderesources.FreeMemory:
		n.sortFreeMemory(reversed)
	case noderesources.Storage:
		n.sortStorage(reversed)
	case noderesources.AllocatableStorage:
		n.sortAvailableStorage(reversed)
	case noderesources.UsedStorage:
		n.sortUsedStorage(reversed)
	case noderesources.StorageEphemeral:
		n.sortStorageEphemeral(reversed)
	case noderesources.AllocatableStorageEphemeral:
		n.sortAvailableStorageEphemeral(reversed)
	case noderesources.UsedStorageEphemeral:
		n.sortUsedStorageEphemeral(reversed)
	case noderesources.FreeStorage:
		n.sortFreeStorage(reversed)
	case noderesources.FreeStorageEphemeral:
		n.sortFreeStorageEphemeral(reversed)
	default:
		// keep current order on unknown sorting
		return
	}
}

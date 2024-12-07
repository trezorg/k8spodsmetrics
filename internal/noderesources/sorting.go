package noderesources

import (
	"sort"

	"github.com/trezorg/k8spodsmetrics/internal/sorting/noderesources"
)

func reverse(less func(i, j int) bool) func(i, j int) bool {
	return func(i, j int) bool {
		return !less(i, j)
	}
}

func (n NodeResourceList) sortByName(reversed bool) {
	less := func(i, j int) bool {
		return n[i].Name < n[j].Name
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortRequestCPU(reversed bool) {
	less := func(i, j int) bool {
		return n[i].CPURequest < n[j].CPURequest
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortLimitCPU(reversed bool) {
	less := func(i, j int) bool {
		return n[i].CPULimit < n[j].CPULimit
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortUsedCPU(reversed bool) {
	less := func(i, j int) bool {
		return n[i].UsedCPU < n[j].UsedCPU
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortAvailableCPU(reversed bool) {
	less := func(i, j int) bool {
		return n[i].AvailableCPU < n[j].AvailableCPU
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortFreeCPU(reversed bool) {
	less := func(i, j int) bool {
		return n[i].FreeCPU < n[j].FreeCPU
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortCPU(reversed bool) {
	less := func(i, j int) bool {
		return n[i].CPU < n[j].CPU
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortRequestMemory(reversed bool) {
	less := func(i, j int) bool {
		return n[i].MemoryRequest < n[j].MemoryRequest
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortLimitMemory(reversed bool) {
	less := func(i, j int) bool {
		return n[i].MemoryLimit < n[j].MemoryLimit
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortUsedMemory(reversed bool) {
	less := func(i, j int) bool {
		return n[i].UsedMemory < n[j].UsedMemory
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortAvailableMemory(reversed bool) {
	less := func(i, j int) bool {
		return n[i].AvailableMemory < n[j].AvailableMemory
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortFreeMemory(reversed bool) {
	less := func(i, j int) bool {
		return n[i].FreeMemory < n[j].FreeMemory
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sortMemory(reversed bool) {
	less := func(i, j int) bool {
		return n[i].Memory < n[j].Memory
	}
	if reversed {
		less = reverse(less)
	}
	sort.Slice(n, less)
}

func (n NodeResourceList) sort(by string, reversed bool) {
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
	}
}

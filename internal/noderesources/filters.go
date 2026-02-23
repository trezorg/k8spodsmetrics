package noderesources

import alerts "github.com/trezorg/k8spodsmetrics/internal/alert"

func (n NodeResourceList) filterBy(predicate nodePredicate) NodeResourceList {
	var result NodeResourceList
	for _, node := range n {
		if predicate(node) {
			result = append(result, node)
		}
	}
	return result
}

func (n NodeResourceList) filterByAlert(alert alerts.Alert) NodeResourceList {
	switch alert {
	case alerts.Any:
		return n.filterBy(func(n NodeResource) bool { return n.IsAlerted() })
	case alerts.Memory:
		return n.filterBy(func(n NodeResource) bool { return n.IsMemoryAlerted() })
	case alerts.MemoryRequest:
		return n.filterBy(func(n NodeResource) bool { return n.IsMemoryRequestAlerted() })
	case alerts.MemoryLimit:
		return n.filterBy(func(n NodeResource) bool { return n.IsMemoryLimitAlerted() })
	case alerts.CPU:
		return n.filterBy(func(n NodeResource) bool { return n.IsCPUAlerted() })
	case alerts.CPURequest:
		return n.filterBy(func(n NodeResource) bool { return n.IsCPURequestAlerted() })
	case alerts.CPULimit:
		return n.filterBy(func(n NodeResource) bool { return n.IsCPULimitAlerted() })
	case alerts.Storage:
		return n.filterBy(func(n NodeResource) bool { return n.IsStorageAlerted() })
	case alerts.StorageEphemeral:
		return n.filterBy(func(n NodeResource) bool { return n.IsStorageEphemeralAlerted() })
	case alerts.None:
		return n
	}
	return n
}

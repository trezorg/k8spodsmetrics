package noderesources

import (
	"bytes"
	"fmt"
	"html/template"

	escapes "github.com/snugfox/ansi-escapes"
	alerts "github.com/trezorg/k8spodsmetrics/internal/alert"
	"github.com/trezorg/k8spodsmetrics/internal/humanize"
	"github.com/trezorg/k8spodsmetrics/pkg/nodemetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/nodes"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	"log/slog"
)

const (
	storageUsedPercentAlert      = 95
	storageEphemeralPercentAlert = 95
)

type (
	NodeResource struct {
		Name                        string `json:"name" yaml:"name"`
		CPU                         int64  `json:"cpu" yaml:"cpu"`
		Memory                      int64  `json:"memory" yaml:"memory"`
		UsedCPU                     int64  `json:"used_cpu" yaml:"used_cpu"`
		UsedMemory                  int64  `json:"used_memory" yaml:"used_memory"`
		AllocatableCPU              int64  `json:"allocatable_cpu" yaml:"allocatable_cpu"`
		AllocatableMemory           int64  `json:"allocatable_memory" yaml:"allocatable_memory"`
		CPURequest                  int64  `json:"cpu_request" yaml:"cpu_request"`
		MemoryRequest               int64  `json:"memory_request" yaml:"memory_request"`
		CPULimit                    int64  `json:"cpu_limit" yaml:"cpu_limit"`
		MemoryLimit                 int64  `json:"memory_limit" yaml:"memory_limit"`
		AvailableCPU                int64  `json:"available_cpu" yaml:"available_cpu"`
		AvailableMemory             int64  `json:"available_memory" yaml:"available_memory"`
		FreeCPU                     int64  `json:"free_cpu" yaml:"free_cpu"`
		FreeMemory                  int64  `json:"free_memory" yaml:"free_memory"`
		Storage                     int64  `json:"storage" yaml:"storage"`
		AllocatableStorage          int64  `json:"allocatable_storage" yaml:"allocatable_storage"`
		UsedStorage                 int64  `json:"used_storage" yaml:"used_storage"`
		FreeStorage                 int64  `json:"free_storage" yaml:"free_storage"`
		StorageEphemeral            int64  `json:"storage_ephemeral" yaml:"storage_ephemeral"`
		AllocatableStorageEphemeral int64  `json:"allocatable_storage_ephemeral" yaml:"allocatable_storage_ephemeral"`
		UsedStorageEphemeral        int64  `json:"used_storage_ephemeral" yaml:"used_storage_ephemeral"`
		FreeStorageEphemeral        int64  `json:"free_storage_ephemeral" yaml:"free_storage_ephemeral"`
	}
	NodeResourceList        []NodeResource
	NodeResourceListEnvelop struct {
		Items NodeResourceList `json:"items,omitempty" yaml:"items,omitempty"`
	}
	nodePredicate func(n NodeResource) bool
)

var (
	nodeTemplate = template.Must(template.New("nodePod").Parse(`Name: {{.Name}}
Memory: {{.MemoryTemplate}}
CPU: {{.CPUTemplate}}
`))
	nodesTemplate = template.Must(template.New("nodesPod").Parse(`{{ range $index, $node := . -}}
{{ $node }}
{{ end -}}`))
)

func (n NodeResource) MemoryTemplate() string {
	memoryRequestStartColor := ""
	memoryRequestEndColor := ""
	memoryLimitStartColor := ""
	memoryLimitEndColor := ""
	if n.Memory <= n.MemoryRequest {
		memoryRequestStartColor = escapes.TextColorYellow
		memoryRequestEndColor = escapes.ColorReset
	}
	if n.Memory <= n.MemoryLimit {
		memoryLimitStartColor = escapes.TextColorRed
		memoryLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"Node=%s/%s, Requests=%s%s%s, Limits=%s%s%s",
		humanize.Bytes(n.Memory),
		humanize.Bytes(n.UsedMemory),
		memoryRequestStartColor,
		humanize.Bytes(n.MemoryRequest),
		memoryRequestEndColor,
		memoryLimitStartColor,
		humanize.Bytes(n.MemoryLimit),
		memoryLimitEndColor,
	)
}

func (n NodeResource) MemoryRequestString() string {
	memoryRequestStartColor := ""
	memoryRequestEndColor := ""
	if n.Memory <= n.MemoryRequest {
		memoryRequestStartColor = escapes.TextColorYellow
		memoryRequestEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryRequestStartColor,
		humanize.Bytes(n.MemoryRequest),
		memoryRequestEndColor,
	)
}

func (n NodeResource) MemoryLimitString() string {
	memoryLimitStartColor := ""
	memoryLimitEndColor := ""
	if n.Memory <= n.MemoryLimit {
		memoryLimitStartColor = escapes.TextColorRed
		memoryLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryLimitStartColor,
		humanize.Bytes(n.MemoryLimit),
		memoryLimitEndColor,
	)
}

func (n NodeResource) MemoryAvailableString() string {
	memoryAvailableStartColor := ""
	memoryAvailableEndColor := ""
	if n.AvailableMemory == 0 {
		memoryAvailableStartColor = escapes.TextColorRed
		memoryAvailableEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryAvailableStartColor,
		humanize.Bytes(n.AvailableMemory),
		memoryAvailableEndColor,
	)
}

func (n NodeResource) MemoryFreeString() string {
	memoryFreeStartColor := ""
	memoryFreeEndColor := ""
	if n.FreeMemory == 0 {
		memoryFreeStartColor = escapes.TextColorRed
		memoryFreeEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryFreeStartColor,
		humanize.Bytes(n.FreeMemory),
		memoryFreeEndColor,
	)
}

func (n NodeResourceList) filterBy(predicate nodePredicate) NodeResourceList {
	var result NodeResourceList
	for _, node := range n {
		if predicate(node) {
			result = append(result, node)
		}
	}
	return result
}

func (n NodeResource) MemoryNodeString() string {
	return humanize.Bytes(n.Memory)
}

func (n NodeResource) MemoryNodeUsedString() string {
	return humanize.Bytes(n.UsedMemory)
}

func (n NodeResource) MemoryNodeAlocatableString() string {
	return humanize.Bytes(n.AllocatableMemory)
}

func (n NodeResource) IsAlerted() bool {
	return n.IsCPUAlerted() || n.IsMemoryAlerted()
}

func (n NodeResource) IsMemoryAlerted() bool {
	return n.Memory <= n.MemoryLimit || n.Memory <= n.MemoryRequest
}

func (n NodeResource) IsMemoryRequestAlerted() bool {
	return n.Memory <= n.MemoryRequest
}

func (n NodeResource) IsMemoryLimitAlerted() bool {
	return n.Memory <= n.MemoryLimit
}

func (n NodeResource) IsCPUAlerted() bool {
	return n.CPU <= n.CPULimit || n.CPU <= n.CPURequest
}

func (n NodeResource) IsCPURequestAlerted() bool {
	return n.CPU <= n.CPURequest
}

func (n NodeResource) IsCPULimitAlerted() bool {
	return n.CPU <= n.CPULimit
}

func (n NodeResource) IsStorageAlerted() bool {
	return (float64(n.UsedStorage)/float64(n.Storage))*100 > storageUsedPercentAlert
}

func (n NodeResource) IsStorageEphemeralAlerted() bool {
	return (float64(n.UsedStorageEphemeral)/float64(n.StorageEphemeral))*100 > storageEphemeralPercentAlert
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

func (n NodeResource) CPUTemplate() string {
	cpuRequestStartColor := ""
	cpuRequestEndColor := ""
	cpuLimitStartColor := ""
	cpuLimitEndColor := ""
	if n.CPU <= n.CPURequest {
		cpuRequestStartColor = escapes.TextColorYellow
		cpuRequestEndColor = escapes.ColorReset
	}
	if n.CPU <= n.CPULimit {
		cpuLimitStartColor = escapes.TextColorRed
		cpuLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"Node=%d/%d, Requests=%s%d%s, Limits=%s%d%s",
		n.CPU,
		n.UsedCPU,
		cpuRequestStartColor,
		n.CPURequest,
		cpuRequestEndColor,
		cpuLimitStartColor,
		n.CPULimit,
		cpuLimitEndColor,
	)
}

func (n NodeResource) CPURequestString() string {
	cpuRequestStartColor := ""
	cpuRequestEndColor := ""
	if n.CPU <= n.CPURequest {
		cpuRequestStartColor = escapes.TextColorYellow
		cpuRequestEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuRequestStartColor,
		n.CPURequest,
		cpuRequestEndColor,
	)
}

func (n NodeResource) CPULimitString() string {
	cpuLimitStartColor := ""
	cpuLimitEndColor := ""
	if n.CPU <= n.CPULimit {
		cpuLimitStartColor = escapes.TextColorRed
		cpuLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuLimitStartColor,
		n.CPULimit,
		cpuLimitEndColor,
	)
}

func (n NodeResource) CPUAvailableString() string {
	cpuAvailableStartColor := ""
	cpuAvailableEndColor := ""
	if n.AvailableCPU == 0 {
		cpuAvailableStartColor = escapes.TextColorRed
		cpuAvailableEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuAvailableStartColor,
		n.AvailableCPU,
		cpuAvailableEndColor,
	)
}

func (n NodeResource) CPUFreeString() string {
	cpuFreeStartColor := ""
	cpuFreeEndColor := ""
	if n.FreeCPU == 0 {
		cpuFreeStartColor = escapes.TextColorRed
		cpuFreeEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuFreeStartColor,
		n.FreeCPU,
		cpuFreeEndColor,
	)
}

func (n NodeResource) StorageString() string {
	return humanize.Bytes(n.Storage)
}

func (n NodeResource) StorageAllocatableString() string {
	return humanize.Bytes(n.AllocatableStorage)
}

func (n NodeResource) StorageFreeString() string {
	return humanize.Bytes(n.FreeStorage)
}

func (n NodeResource) StorageUsedString() string {
	usedStorageStartColor := ""
	usedStorageEndColor := ""
	if n.IsStorageAlerted() {
		usedStorageStartColor = escapes.TextColorRed
		usedStorageEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		usedStorageStartColor,
		humanize.Bytes(n.UsedStorage),
		usedStorageEndColor,
	)
}

func (n NodeResource) StorageEphemeralString() string {
	return humanize.Bytes(n.StorageEphemeral)
}

func (n NodeResource) StorageAllocatableEphemeralString() string {
	return humanize.Bytes(n.AllocatableStorageEphemeral)
}

func (n NodeResource) StorageFreeEphemeralString() string {
	return humanize.Bytes(n.FreeStorageEphemeral)
}

func (n NodeResource) StorageUsedEphemeralString() string {
	usedStorageStartColor := ""
	usedStorageEndColor := ""
	if n.IsStorageEphemeralAlerted() {
		usedStorageStartColor = escapes.TextColorRed
		usedStorageEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		usedStorageStartColor,
		humanize.Bytes(n.UsedStorageEphemeral),
		usedStorageEndColor,
	)
}

func (n NodeResource) String() string {
	var buffer bytes.Buffer
	if err := nodeTemplate.Execute(&buffer, n); err != nil {
		panic(err)
	}
	return buffer.String()
}

func (n NodeResourceList) String() string {
	var buffer bytes.Buffer
	if err := nodesTemplate.Execute(&buffer, n); err != nil {
		panic(err)
	}
	return buffer.String()
}

func merge(podResourceList pods.PodResourceList, nodeList nodes.NodeList, nodeMetricList nodemetrics.List) NodeResourceList {
	nodesMap := make(map[string]*NodeResource)
	for _, node := range nodeList {
		nodesMap[node.Name] = &NodeResource{
			Name:                        node.Name,
			CPU:                         node.CPU,
			Memory:                      node.Memory,
			AllocatableCPU:              node.AllocatableCPU,
			AllocatableMemory:           node.AllocatableMemory,
			Storage:                     node.Storage,
			AllocatableStorage:          node.AllocatableStorage,
			UsedStorage:                 node.UsedStorage,
			StorageEphemeral:            node.StorageEphemeral,
			AllocatableStorageEphemeral: node.AllocatableStorageEphemeral,
			UsedStorageEphemeral:        node.UsedStorageEphemeral,
		}
	}
	for _, pod := range podResourceList {
		nodeResource, ok := nodesMap[pod.NodeName]
		if !ok {
			slog.Debug("Cannot find node", slog.String("node", pod.NodeName))
			continue
		}
		for _, container := range pod.Containers {
			nodeResource.CPULimit += container.Limits.CPU
			nodeResource.CPURequest += container.Requests.CPU
			nodeResource.MemoryLimit += container.Limits.Memory
			nodeResource.MemoryRequest += container.Requests.Memory
		}
		nodeResource.AvailableCPU = nodeResource.AllocatableCPU - nodeResource.CPURequest
		nodeResource.AvailableMemory = nodeResource.AllocatableMemory - nodeResource.MemoryRequest
	}
	for _, metric := range nodeMetricList {
		nodeResource, ok := nodesMap[metric.Name]
		if !ok {
			slog.Warn("Cannot find node", slog.String("node", metric.Name))
			continue
		}
		nodeResource.UsedCPU = metric.CPU
		nodeResource.UsedMemory = metric.Memory
		nodeResource.FreeCPU = nodeResource.AllocatableCPU - metric.CPU
		nodeResource.FreeMemory = nodeResource.AllocatableMemory - metric.Memory
		nodeResource.FreeStorage = nodeResource.AllocatableStorage - metric.Storage
		nodeResource.UsedStorage = metric.Storage
		nodeResource.FreeStorageEphemeral = nodeResource.AllocatableStorageEphemeral - metric.StorageEphemeral
		nodeResource.UsedStorageEphemeral = metric.StorageEphemeral
	}
	nodeResourceList := make(NodeResourceList, 0, len(nodesMap))
	for _, node := range nodesMap {
		nodeResourceList = append(nodeResourceList, *node)
	}
	return nodeResourceList
}

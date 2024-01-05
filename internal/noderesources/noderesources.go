package noderesources

import (
	"bytes"
	"fmt"
	"html/template"

	escapes "github.com/snugfox/ansi-escapes"
	alerts "github.com/trezorg/k8spodsmetrics/internal/alert"
	"github.com/trezorg/k8spodsmetrics/internal/humanize"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/pkg/nodemetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/nodes"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	"golang.org/x/exp/slog"
)

type (
	NodeResource struct {
		Name              string `json:"name,omitempty" yaml:"name,omitempty"`
		CPU               int64  `json:"cpu,omitempty" yaml:"cpu,omitempty"`
		Memory            int64  `json:"memory,omitempty" yaml:"memory,omitempty"`
		UsedCPU           int64  `json:"used_cpu,omitempty" yaml:"used_cpu,omitempty"`
		UsedMemory        int64  `json:"used_memory,omitempty" yaml:"used_memory,omitempty"`
		AllocatableCPU    int64  `json:"allocatable_cpu,omitempty" yaml:"allocatable_cpu,omitempty"`
		AllocatableMemory int64  `json:"allocatable_memory,omitempty" yaml:"allocatable_memory,omitempty"`
		CPURequest        int64  `json:"cpu_request,omitempty" yaml:"cpu_request,omitempty"`
		MemoryRequest     int64  `json:"memory_request,omitempty" yaml:"memory_request,omitempty"`
		CPULimit          int64  `json:"cpu_limit,omitempty" yaml:"cpu_limit,omitempty"`
		MemoryLimit       int64  `json:"memory_limit,omitempty" yaml:"memory_limit,omitempty"`
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
	default:
		return n
	}
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

func merge(podResourceList pods.PodResourceList, nodeList nodes.NodeList, nodeMetricList nodemetrics.NodeMetricsList) NodeResourceList {
	nodesMap := make(map[string]*NodeResource)
	for _, node := range nodeList {
		nodesMap[node.Name] = &NodeResource{
			Name:              node.Name,
			CPU:               node.CPU,
			Memory:            node.Memory,
			AllocatableCPU:    node.AllocatableCPU,
			AllocatableMemory: node.AllocatableMemory,
		}
	}
	for _, pod := range podResourceList {
		nodeResource, ok := nodesMap[pod.NodeName]
		if !ok {
			logger.Debug("Cannot find node", slog.String("node", pod.NodeName))
			continue
		}
		for _, container := range pod.Containers {
			nodeResource.CPULimit += container.Limits.CPU
			nodeResource.CPURequest += container.Requests.CPU
			nodeResource.MemoryLimit += container.Limits.Memory
			nodeResource.MemoryRequest += container.Requests.Memory
		}
	}
	for _, metric := range nodeMetricList {
		nodeResource, ok := nodesMap[metric.Name]
		if !ok {
			logger.Warn("Cannot find node", slog.String("node", metric.Name))
			continue
		}
		nodeResource.UsedCPU = metric.CPU
		nodeResource.UsedMemory = metric.Memory
	}
	nodeResourceList := make(NodeResourceList, 0, len(nodesMap))
	for _, node := range nodesMap {
		nodeResourceList = append(nodeResourceList, *node)
	}
	return nodeResourceList
}

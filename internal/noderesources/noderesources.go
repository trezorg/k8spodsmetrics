package noderesources

import (
	"bytes"
	"fmt"
	"html/template"

	escapes "github.com/snugfox/ansi-escapes"
	"github.com/trezorg/k8spodsmetrics/internal/humanize"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/pkg/nodemetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/nodes"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	"golang.org/x/exp/slog"
)

type (
	NodeResource struct {
		CPU               int64
		Memory            int64
		UsedCPU           int64
		UsedMemory        int64
		AllocatableCPU    int64
		AllocatableMemory int64
		PodsCPURequest    int64
		PodsMemoryRequest int64
		PodsCPULimit      int64
		PodsMemoryLimit   int64
		Name              string
	}
	NodeResourceList []NodeResource
)

var (
	nodeTemplate = template.Must(template.New("nodePod").Parse(`Name: {{.Name}}
Memory:
{{.MemoryTemplate}}
CPU:
{{.CPUTemplate}}
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
	if n.Memory <= n.PodsMemoryRequest {
		memoryRequestStartColor = escapes.TextColorYellow
		memoryRequestEndColor = escapes.ColorReset
	}
	if n.Memory <= n.PodsMemoryLimit {
		memoryLimitStartColor = escapes.TextColorRed
		memoryLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"Node=%s/%s, Requests=%s%s%s, Limits=%s%s%s",
		humanize.Bytes(n.Memory),
		humanize.Bytes(n.UsedMemory),
		memoryRequestStartColor,
		humanize.Bytes(n.PodsMemoryRequest),
		memoryRequestEndColor,
		memoryLimitStartColor,
		humanize.Bytes(n.PodsMemoryLimit),
		memoryLimitEndColor,
	)
}

func (n NodeResource) IsAlerted() bool {
	return n.CPU <= n.PodsCPULimit || n.CPU <= n.PodsCPURequest || n.Memory <= n.PodsMemoryLimit || n.Memory <= n.PodsMemoryRequest
}

func (n NodeResourceList) filterAlerts() NodeResourceList {
	var result NodeResourceList
	for _, node := range n {
		if node.IsAlerted() {
			result = append(result, node)
		}
	}
	return result
}

func (n NodeResource) CPUTemplate() string {
	cpuRequestStartColor := ""
	cpuRequestEndColor := ""
	cpuLimitStartColor := ""
	cpuLimitEndColor := ""
	if n.CPU <= n.PodsCPURequest {
		cpuRequestStartColor = escapes.TextColorYellow
		cpuRequestEndColor = escapes.ColorReset
	}
	if n.CPU <= n.PodsCPULimit {
		cpuLimitStartColor = escapes.TextColorRed
		cpuLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"Node=%d/%d, Requests=%s%d%s, Limits=%s%d%s",
		n.CPU,
		n.UsedCPU,
		cpuRequestStartColor,
		n.PodsCPURequest,
		cpuRequestEndColor,
		cpuLimitStartColor,
		n.PodsCPULimit,
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
			logger.Warn("Cannot find node", slog.String("node", pod.NodeName))
			continue
		}
		for _, container := range pod.Containers {
			nodeResource.PodsCPULimit += container.Limits.CPU
			nodeResource.PodsCPURequest += container.Requests.CPU
			nodeResource.PodsMemoryLimit += container.Limits.Memory
			nodeResource.PodsMemoryRequest += container.Requests.Memory
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

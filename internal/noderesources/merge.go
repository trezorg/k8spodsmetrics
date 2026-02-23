package noderesources

import (
	"log/slog"

	"github.com/trezorg/k8spodsmetrics/pkg/nodemetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/nodes"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

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

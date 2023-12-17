package nodes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type NodeFilter struct {
	LabelSelector string
	FieldSelector string
}

type Node struct {
	CPU               int64
	Memory            int64
	AllocatableCPU    int64
	AllocatableMemory int64
	UsedCPU           int64
	UsedMemory        int64
	Name              string
}

type NodeList []Node

func Nodes(ctx context.Context, corev1 corev1.CoreV1Interface, filter NodeFilter, name string) (NodeList, error) {
	var result NodeList
	var nodes *v1.NodeList
	var err error
	client := corev1.Nodes()
	if name == "" {
		nodes, err = client.List(ctx, metav1.ListOptions{
			LabelSelector: filter.LabelSelector,
			FieldSelector: filter.FieldSelector,
		})
		if err != nil {
			return result, err
		}
	} else {
		node, err := client.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		items := []v1.Node{*node}
		allNodes := v1.NodeList{Items: items}
		nodes = &allNodes
	}
	for _, node := range nodes.Items {
		memory, ok := node.Status.Capacity.Memory().AsInt64()
		if !ok {
			memory = int64(node.Status.Capacity.Memory().AsApproximateFloat64())
		}
		allocatableMemory, ok := node.Status.Allocatable.Memory().AsInt64()
		if !ok {
			allocatableMemory = int64(node.Status.Allocatable.Memory().AsApproximateFloat64())
		}
		nodeResource := Node{
			Name:              node.Name,
			CPU:               node.Status.Capacity.Cpu().MilliValue(),
			AllocatableCPU:    node.Status.Allocatable.Cpu().MilliValue(),
			Memory:            memory,
			AllocatableMemory: allocatableMemory,
		}
		nodeResource.UsedCPU = nodeResource.CPU - nodeResource.AllocatableCPU
		nodeResource.UsedMemory = nodeResource.Memory - nodeResource.AllocatableMemory
		result = append(result, nodeResource)
	}
	return result, err
}

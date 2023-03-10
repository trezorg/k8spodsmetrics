package nodes

import (
	"context"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"golang.org/x/exp/slog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type NodeFilter struct {
	LabelSelector string
	FieldSelector string
}

type Node struct {
	CPU    int64
	Memory int64
	Name   string
}

type NodeList []Node

func Nodes(ctx context.Context, corev1 corev1.CoreV1Interface, filter NodeFilter) (NodeList, error) {
	var result NodeList
	nodes, err := corev1.Nodes().List(ctx, metav1.ListOptions{
		LabelSelector: filter.LabelSelector,
		FieldSelector: filter.FieldSelector,
	})
	if err != nil {
		return result, err
	}
	for _, node := range nodes.Items {
		memory, ok := node.Status.Capacity.Memory().AsInt64()
		if !ok {
			logger.Warn("Cannot get node status memory", slog.String("node", node.Name))
		}
		nodeResource := Node{
			Name:   node.Name,
			CPU:    node.Status.Capacity.Cpu().MilliValue(),
			Memory: memory,
		}
		result = append(result, nodeResource)
	}
	return result, err
}

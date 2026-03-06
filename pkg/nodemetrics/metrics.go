package nodemetrics

import (
	"context"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type NodeMetric struct {
	Name             string
	CPU              int64
	Memory           int64
	Storage          int64
	StorageEphemeral int64
}

type List []NodeMetric

type MetricsFilter struct {
	LabelSelector string
	FieldSelector string
}

// Metrics get pod metrics for MetricFilter
func Metrics(ctx context.Context, api metricsv1beta1.MetricsV1beta1Interface, filter MetricsFilter, nodeName string) (List, error) {
	var nodeMetrics *v1beta1.NodeMetricsList
	var err error
	if nodeName == "" {
		nodeMetrics, err = listNodeMetrics(ctx, api, metav1.ListOptions{
			LabelSelector: filter.LabelSelector,
			FieldSelector: filter.FieldSelector,
		})
	} else {
		var nodeMetric *v1beta1.NodeMetrics
		nodeMetric, err = api.NodeMetricses().Get(ctx, nodeName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		allNodeMetrics := v1beta1.NodeMetricsList{Items: []v1beta1.NodeMetrics{*nodeMetric}}
		nodeMetrics = &allNodeMetrics
	}
	if err != nil {
		return nil, err
	}

	result := make(List, 0, len(nodeMetrics.Items))
	for _, nodeMetric := range nodeMetrics.Items {
		resourceList := nodeMetric.Usage
		metric := NodeMetric{Name: nodeMetric.Name}
		for name, quantity := range resourceList {
			switch name { //nolint:exhaustive // it is ok
			case v1.ResourceMemory:
				if memory, ok := quantity.AsInt64(); ok {
					metric.Memory = memory
				}
			case v1.ResourceCPU:
				metric.CPU = quantity.MilliValue()
			case v1.ResourceStorage:
				if storage, ok := quantity.AsInt64(); ok {
					metric.Storage = storage
				}
			case v1.ResourceEphemeralStorage:
				if storage, ok := quantity.AsInt64(); ok {
					metric.StorageEphemeral = storage
				}
			default:
				// Ignore other resource types
			}
		}
		result = append(result, metric)
	}
	return result, nil
}

func listNodeMetrics(
	ctx context.Context,
	api metricsv1beta1.MetricsV1beta1Interface,
	opts metav1.ListOptions,
) (*v1beta1.NodeMetricsList, error) {
	result := &v1beta1.NodeMetricsList{}
	for {
		nodeMetrics, err := api.NodeMetricses().List(ctx, opts)
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, nodeMetrics.Items...)
		if nodeMetrics.Continue == "" {
			return result, nil
		}
		opts.Continue = nodeMetrics.Continue
	}
}

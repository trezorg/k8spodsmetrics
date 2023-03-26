package nodemetrics

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type NodeMetric struct {
	Name   string
	CPU    int64
	Memory int64
}

type NodeMetricsList []NodeMetric

type MetricsFilter struct {
	LabelSelector string
	FieldSelector string
}

// Metrics get pod metrics for MetricFilter
func Metrics(ctx context.Context, api metricsv1beta1.MetricsV1beta1Interface, filter MetricsFilter) (NodeMetricsList, error) {
	var result NodeMetricsList
	nodeMetrics, err := api.NodeMetricses().List(ctx, metav1.ListOptions{
		LabelSelector: filter.LabelSelector,
		FieldSelector: filter.FieldSelector,
	})
	if err != nil {
		return result, err
	}
	for _, nodeMetric := range nodeMetrics.Items {
		resourceList := nodeMetric.Usage
		metric := NodeMetric{Name: nodeMetric.Name}
		for name, quantity := range resourceList {
			if name == "memory" {
				if memory, ok := quantity.AsInt64(); ok {
					metric.Memory = memory
				}
			}
			if name == "cpu" {
				metric.CPU = int64(quantity.ToDec().AsApproximateFloat64() * 1000)
			}
		}
		result = append(result, metric)
	}
	return result, nil
}

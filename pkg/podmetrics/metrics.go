package podmetrics

import (
	"context"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type Metric struct {
	CPU    int64
	Memory int64
}
type ContainerMetric struct {
	Name string
	Metric
}
type PodMetric struct {
	Namespace  string
	Name       string
	Containers []ContainerMetric
}

type PodMetricList []PodMetric

type MetricFilter struct {
	Namespace     string
	LabelSelector string
	FieldSelector string
}

// Metrics get pod metrics for MetricFilter
func Metrics(ctx context.Context, api metricsv1beta1.MetricsV1beta1Interface, filter MetricFilter) (PodMetricList, error) {
	var result PodMetricList
	podMetrics, err := api.PodMetricses(filter.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: filter.LabelSelector,
		FieldSelector: filter.FieldSelector,
	})
	if err != nil {
		return result, err
	}
	for _, podMetric := range podMetrics.Items {
		metric := PodMetric{Name: podMetric.Name, Namespace: podMetric.Namespace}
		for _, container := range podMetric.Containers {
			containerMetric := ContainerMetric{
				Name: container.Name,
			}
			cpu, ok := container.Usage.Cpu().AsInt64()
			if ok {
				containerMetric.CPU = cpu
			}
			memory, ok := container.Usage.Memory().AsInt64()
			if ok {
				containerMetric.Memory = memory
			}
			metric.Containers = append(metric.Containers, containerMetric)
		}
		sort.Slice(metric.Containers, func(i, j int) bool {
			return metric.Containers[i].Name < metric.Containers[j].Name
		})
	}
	return result, nil
}

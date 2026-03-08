package podmetrics

import (
	"cmp"
	"context"
	"errors"
	"slices"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type Metric struct {
	CPU              int64
	Memory           int64
	Storage          int64
	StorageEphemeral int64
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
	Namespaces    []string
	LabelSelector string
	FieldSelector string
}

// Metrics get pod metrics for MetricFilter
func Metrics(ctx context.Context, api metricsv1beta1.MetricsV1beta1Interface, filter MetricFilter) (PodMetricList, error) {
	filter.Namespaces = slices.DeleteFunc(filter.Namespaces, func(n string) bool { return n == "" })

	// If no namespaces specified, query all namespaces (empty string)
	if len(filter.Namespaces) == 0 {
		return listMetrics(ctx, api, filter, "")
	}

	// Single namespace
	if len(filter.Namespaces) == 1 {
		return listMetrics(ctx, api, filter, filter.Namespaces[0])
	}

	// Multiple namespaces: query each in parallel
	var wg sync.WaitGroup

	rErrors := make([]error, len(filter.Namespaces))
	metrics := make([]PodMetricList, len(filter.Namespaces))

	for idx, ns := range filter.Namespaces {
		wg.Go(func() {
			metrics[idx], rErrors[idx] = listMetrics(ctx, api, filter, ns)
		})
	}

	wg.Wait()

	if err := errors.Join(rErrors...); err != nil {
		return nil, err
	}

	return slices.Concat(metrics...), nil
}

func listMetrics(ctx context.Context, api metricsv1beta1.MetricsV1beta1Interface, filter MetricFilter, namespace string) (PodMetricList, error) {
	opts := metav1.ListOptions{
		LabelSelector: filter.LabelSelector,
		FieldSelector: filter.FieldSelector,
	}

	result := PodMetricList{}
	for {
		podMetrics, err := api.PodMetricses(namespace).List(ctx, opts)
		if err != nil {
			return nil, err
		}

		for _, podMetric := range podMetrics.Items {
			metric := PodMetric{Name: podMetric.Name, Namespace: podMetric.Namespace}
			for _, container := range podMetric.Containers {
				containerMetric := ContainerMetric{
					Name: container.Name,
				}
				containerMetric.CPU = container.Usage.Cpu().MilliValue()
				memory, ok := container.Usage.Memory().AsInt64()
				if ok {
					containerMetric.Memory = memory
				}
				storage, ok := container.Usage.Storage().AsInt64()
				if ok {
					containerMetric.Storage = storage
				}
				storage, ok = container.Usage.StorageEphemeral().AsInt64()
				if ok {
					containerMetric.StorageEphemeral = storage
				}
				metric.Containers = append(metric.Containers, containerMetric)
			}
			slices.SortFunc(metric.Containers, func(a, b ContainerMetric) int {
				return cmp.Compare(a.Name, b.Name)
			})
			result = append(result, metric)
		}
		if podMetrics.Continue == "" {
			return result, nil
		}
		opts.Continue = podMetrics.Continue
	}
}

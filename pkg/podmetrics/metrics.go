package podmetrics

import (
	"context"
	"errors"
	"slices"
	"sort"
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

	var errs error
	rErrors := make([]error, len(filter.Namespaces))
	metrics := make([]PodMetricList, len(filter.Namespaces))

	for idx, ns := range filter.Namespaces {
		wg.Go(func() {
			metrics[idx], rErrors[idx] = listMetrics(ctx, api, filter, ns)
		})
	}

	wg.Wait()

	for _, err := range rErrors {
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}

	resultLen := 0
	for _, m := range metrics {
		resultLen += len(m)
	}
	result := make(PodMetricList, 0, resultLen)
	for _, m := range metrics {
		result = append(result, m...)
	}
	return result, nil
}

func listMetrics(ctx context.Context, api metricsv1beta1.MetricsV1beta1Interface, filter MetricFilter, namespace string) (PodMetricList, error) {
	podMetrics, err := api.PodMetricses(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: filter.LabelSelector,
		FieldSelector: filter.FieldSelector,
	})
	if err != nil {
		return nil, err
	}

	result := make(PodMetricList, 0, len(podMetrics.Items))
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
		sort.Slice(metric.Containers, func(i, j int) bool {
			return metric.Containers[i].Name < metric.Containers[j].Name
		})
		result = append(result, metric)
	}
	return result, nil
}

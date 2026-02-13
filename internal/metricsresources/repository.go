package metricsresources

import (
	"context"
	"errors"
	"slices"
	"sync"

	"log/slog"

	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type PodRepository interface {
	FetchPods(
		ctx context.Context,
		podsClient corev1.CoreV1Interface,
		filter pods.PodFilter,
		nodeNames ...string,
	) (pods.PodResourceList, error)
	FetchMetrics(
		ctx context.Context,
		metricsClient metricsv1beta1.MetricsV1beta1Interface,
		filter podmetrics.MetricFilter,
	) (podmetrics.PodMetricList, error)
}

type podRepository struct{}

func NewPodRepository() PodRepository {
	return &podRepository{}
}

func (podRepository) FetchPods(
	ctx context.Context,
	podsClient corev1.CoreV1Interface,
	filter pods.PodFilter,
	nodeNames ...string,
) (pods.PodResourceList, error) {
	return pods.Pods(ctx, podsClient, filter, nodeNames...)
}

func (podRepository) FetchMetrics(
	ctx context.Context,
	metricsClient metricsv1beta1.MetricsV1beta1Interface,
	filter podmetrics.MetricFilter,
) (podmetrics.PodMetricList, error) {
	return podmetrics.Metrics(ctx, metricsClient, filter)
}

type FetchConfig struct {
	Namespaces    []string
	Label         string
	FieldSelector string
	Nodes         []string
}

func FetchPodMetrics(
	ctx context.Context,
	repo PodRepository,
	metricsClient metricsv1beta1.MetricsV1beta1Interface,
	podsClient corev1.CoreV1Interface,
	config FetchConfig,
) (PodMetricsResourceList, error) {
	slog.Debug("Getting metrics...")

	// Filter out empty namespace strings
	config.Namespaces = slices.DeleteFunc(config.Namespaces, func(n string) bool { return n == "" })

	// If no namespaces specified or single namespace, use existing logic
	if len(config.Namespaces) <= 1 {
		var ns string
		if len(config.Namespaces) == 1 {
			ns = config.Namespaces[0]
		}
		return fetchPodMetricsForNamespace(ctx, repo, metricsClient, podsClient, config, ns)
	}

	// Multiple namespaces: query each in parallel
	var wg sync.WaitGroup
	results := make([]PodMetricsResourceList, len(config.Namespaces))
	rErrors := make([]error, len(config.Namespaces))

	for idx, ns := range config.Namespaces {
		wg.Go(func() {
			results[idx], rErrors[idx] = fetchPodMetricsForNamespace(ctx, repo, metricsClient, podsClient, config, ns)
		})
	}

	wg.Wait()

	var rErr error
	for _, err := range rErrors {
		if err != nil {
			rErr = errors.Join(rErr, err)
		}
	}

	if rErr != nil {
		return nil, rErr
	}

	// Merge results from all namespaces
	resultLen := 0
	for _, r := range results {
		resultLen += len(r)
	}
	merged := make(PodMetricsResourceList, 0, resultLen)
	for _, r := range results {
		merged = append(merged, r...)
	}
	return merged, nil
}

func fetchPodMetricsForNamespace(
	ctx context.Context,
	repo PodRepository,
	metricsClient metricsv1beta1.MetricsV1beta1Interface,
	podsClient corev1.CoreV1Interface,
	config FetchConfig,
	namespace string,
) (PodMetricsResourceList, error) {
	var podMetricsResourceList PodMetricsResourceList

	cErrors := make([]error, 2)
	var podsList pods.PodResourceList
	var metricsList podmetrics.PodMetricList
	wg := sync.WaitGroup{}

	wg.Go(func() {
		metricsList, cErrors[0] = repo.FetchMetrics(ctx, metricsClient, podmetrics.MetricFilter{
			Namespaces:    []string{namespace},
			LabelSelector: config.Label,
			FieldSelector: config.FieldSelector,
		})
	})

	wg.Go(func() {
		podsList, cErrors[1] = repo.FetchPods(ctx, podsClient, pods.PodFilter{
			Namespaces:    []string{namespace},
			LabelSelector: config.Label,
			FieldSelector: config.FieldSelector,
		}, config.Nodes...)
	})

	wg.Wait()

	var rErr error

	for _, err := range cErrors {
		if err != nil {
			rErr = errors.Join(rErr, err)
		}
	}

	if rErr != nil {
		return podMetricsResourceList, rErr
	}

	podMetricsResourceList = merge(podsList, metricsList)
	return podMetricsResourceList, nil
}

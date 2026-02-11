package metricsresources

import (
	"context"
	"errors"
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
	Namespace     string
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
	var podMetricsResourceList PodMetricsResourceList

	cErrors := make([]error, 2)
	var podsList pods.PodResourceList
	var metricsList podmetrics.PodMetricList
	wg := sync.WaitGroup{}

	wg.Go(func() {
		metricsList, cErrors[0] = repo.FetchMetrics(ctx, metricsClient, podmetrics.MetricFilter{
			Namespace:     config.Namespace,
			LabelSelector: config.Label,
			FieldSelector: config.FieldSelector,
		})
	})

	wg.Go(func() {
		podsList, cErrors[1] = repo.FetchPods(ctx, podsClient, pods.PodFilter{
			Namespace:     config.Namespace,
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

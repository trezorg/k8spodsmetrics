package metricsresources

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type stubPodRepository struct {
	fetchPods    func(corev1.CoreV1Interface, pods.PodFilter, ...string) (pods.PodResourceList, error)
	fetchMetrics func(metricsv1beta1.MetricsV1beta1Interface, podmetrics.MetricFilter) (podmetrics.PodMetricList, error)
}

func (s stubPodRepository) FetchPods(
	_ context.Context,
	podsClient corev1.CoreV1Interface,
	filter pods.PodFilter,
	nodeNames ...string,
) (pods.PodResourceList, error) {
	if s.fetchPods != nil {
		return s.fetchPods(podsClient, filter, nodeNames...)
	}
	return nil, nil
}

func (s stubPodRepository) FetchMetrics(
	_ context.Context,
	metricsClient metricsv1beta1.MetricsV1beta1Interface,
	filter podmetrics.MetricFilter,
) (podmetrics.PodMetricList, error) {
	if s.fetchMetrics != nil {
		return s.fetchMetrics(metricsClient, filter)
	}
	return nil, nil
}

func TestFetchPodMetricsForNamespaceWrapsBranchErrors(t *testing.T) {
	t.Run("metrics branch", func(t *testing.T) {
		rootErr := errors.New("metrics api down")
		repo := stubPodRepository{
			fetchMetrics: func(metricsv1beta1.MetricsV1beta1Interface, podmetrics.MetricFilter) (podmetrics.PodMetricList, error) {
				return nil, rootErr
			},
		}

		_, err := fetchPodMetricsForNamespace(t.Context(), repo, nil, nil, FetchConfig{}, "team-a")
		require.Error(t, err)
		require.ErrorContains(t, err, `fetch pod usage metrics for namespace "team-a"`)
		require.ErrorIs(t, err, rootErr)
	})

	t.Run("pods branch", func(t *testing.T) {
		rootErr := errors.New("pods api down")
		repo := stubPodRepository{
			fetchPods: func(corev1.CoreV1Interface, pods.PodFilter, ...string) (pods.PodResourceList, error) {
				return nil, rootErr
			},
		}

		_, err := fetchPodMetricsForNamespace(t.Context(), repo, nil, nil, FetchConfig{}, "team-b")
		require.Error(t, err)
		require.ErrorContains(t, err, `fetch pod resources for namespace "team-b"`)
		require.ErrorIs(t, err, rootErr)
	})
}

func TestFetchPodMetricsWrapsAllNamespaceErrors(t *testing.T) {
	rootErr := errors.New("metrics api down")
	repo := stubPodRepository{
		fetchMetrics: func(metricsv1beta1.MetricsV1beta1Interface, podmetrics.MetricFilter) (podmetrics.PodMetricList, error) {
			return nil, rootErr
		},
	}

	_, err := FetchPodMetrics(t.Context(), repo, nil, nil, FetchConfig{})
	require.Error(t, err)
	require.ErrorContains(t, err, "across all namespaces")
	require.NotContains(t, err.Error(), `namespace ""`)
	require.ErrorIs(t, err, rootErr)
}

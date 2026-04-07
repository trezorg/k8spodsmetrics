package noderesources

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trezorg/k8spodsmetrics/pkg/nodemetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/nodes"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type stubNodeRepository struct {
	fetchNodes   func(corev1.CoreV1Interface, nodes.NodeFilter, string) (nodes.NodeList, error)
	fetchPods    func(corev1.CoreV1Interface, pods.PodFilter, string) (pods.PodResourceList, error)
	fetchMetrics func(metricsv1beta1.MetricsV1beta1Interface, nodemetrics.MetricsFilter, string) (nodemetrics.List, error)
}

func (s stubNodeRepository) FetchNodes(
	_ context.Context,
	coreClient corev1.CoreV1Interface,
	filter nodes.NodeFilter,
	name string,
) (nodes.NodeList, error) {
	if s.fetchNodes != nil {
		return s.fetchNodes(coreClient, filter, name)
	}
	return nil, nil
}

func (s stubNodeRepository) FetchPods(
	_ context.Context,
	coreClient corev1.CoreV1Interface,
	filter pods.PodFilter,
	name string,
) (pods.PodResourceList, error) {
	if s.fetchPods != nil {
		return s.fetchPods(coreClient, filter, name)
	}
	return nil, nil
}

func (s stubNodeRepository) FetchMetrics(
	_ context.Context,
	metricsClient metricsv1beta1.MetricsV1beta1Interface,
	filter nodemetrics.MetricsFilter,
	name string,
) (nodemetrics.List, error) {
	if s.fetchMetrics != nil {
		return s.fetchMetrics(metricsClient, filter, name)
	}
	return nil, nil
}

func TestFetchNodeMetricsWrapsBranchErrors(t *testing.T) {
	t.Run("nodes branch", func(t *testing.T) {
		rootErr := errors.New("nodes api down")
		repo := stubNodeRepository{
			fetchNodes: func(corev1.CoreV1Interface, nodes.NodeFilter, string) (nodes.NodeList, error) {
				return nil, rootErr
			},
		}

		_, err := FetchNodeMetrics(t.Context(), repo, nil, nil, FetchConfig{Name: "worker-1"})
		require.Error(t, err)
		require.ErrorContains(t, err, `fetch nodes for node "worker-1"`)
		require.ErrorIs(t, err, rootErr)
	})

	t.Run("pods branch", func(t *testing.T) {
		rootErr := errors.New("pods api down")
		repo := stubNodeRepository{
			fetchPods: func(corev1.CoreV1Interface, pods.PodFilter, string) (pods.PodResourceList, error) {
				return nil, rootErr
			},
		}

		_, err := FetchNodeMetrics(t.Context(), repo, nil, nil, FetchConfig{Name: "worker-2"})
		require.Error(t, err)
		require.ErrorContains(t, err, `fetch pods for node "worker-2"`)
		require.ErrorIs(t, err, rootErr)
	})

	t.Run("metrics branch", func(t *testing.T) {
		rootErr := errors.New("metrics api down")
		repo := stubNodeRepository{
			fetchMetrics: func(metricsv1beta1.MetricsV1beta1Interface, nodemetrics.MetricsFilter, string) (nodemetrics.List, error) {
				return nil, rootErr
			},
		}

		_, err := FetchNodeMetrics(t.Context(), repo, nil, nil, FetchConfig{Name: "worker-3"})
		require.Error(t, err)
		require.ErrorContains(t, err, `fetch node metrics for node "worker-3"`)
		require.ErrorIs(t, err, rootErr)
	})
}

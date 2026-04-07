package noderesources

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"log/slog"

	"github.com/trezorg/k8spodsmetrics/pkg/nodemetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/nodes"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type NodeRepository interface {
	FetchNodes(
		ctx context.Context,
		coreClient corev1.CoreV1Interface,
		filter nodes.NodeFilter,
		name string,
	) (nodes.NodeList, error)
	FetchPods(
		ctx context.Context,
		coreClient corev1.CoreV1Interface,
		filter pods.PodFilter,
		name string,
	) (pods.PodResourceList, error)
	FetchMetrics(
		ctx context.Context,
		metricsClient metricsv1beta1.MetricsV1beta1Interface,
		filter nodemetrics.MetricsFilter,
		name string,
	) (nodemetrics.List, error)
}

type nodeRepository struct{}

func NewNodeRepository() NodeRepository {
	return &nodeRepository{}
}

func (nodeRepository) FetchNodes(
	ctx context.Context,
	coreClient corev1.CoreV1Interface,
	filter nodes.NodeFilter,
	name string,
) (nodes.NodeList, error) {
	return nodes.Nodes(ctx, coreClient, filter, name)
}

func (nodeRepository) FetchPods(
	ctx context.Context,
	coreClient corev1.CoreV1Interface,
	filter pods.PodFilter,
	name string,
) (pods.PodResourceList, error) {
	return pods.Pods(ctx, coreClient, filter, name)
}

func (nodeRepository) FetchMetrics(
	ctx context.Context,
	metricsClient metricsv1beta1.MetricsV1beta1Interface,
	filter nodemetrics.MetricsFilter,
	name string,
) (nodemetrics.List, error) {
	return nodemetrics.Metrics(ctx, metricsClient, filter, name)
}

type FetchConfig struct {
	Label string
	Name  string
}

func FetchNodeMetrics(
	ctx context.Context,
	repo NodeRepository,
	coreClient corev1.CoreV1Interface,
	metricsClient metricsv1beta1.MetricsV1beta1Interface,
	config FetchConfig,
) (NodeResourceList, error) {
	slog.Debug("Getting nodes info...")
	var nodeResources NodeResourceList
	numberOfRequests := 3
	var podsList pods.PodResourceList
	var nodesList nodes.NodeList
	var nodeMetricsList nodemetrics.List
	cErrors := make([]error, numberOfRequests)

	wg := sync.WaitGroup{}

	wg.Go(func() {
		nodesList, cErrors[0] = repo.FetchNodes(ctx, coreClient, nodes.NodeFilter{LabelSelector: config.Label}, config.Name)
		if cErrors[0] != nil {
			cErrors[0] = wrapNodeFetchError("fetch nodes", config.Name, cErrors[0])
		}
	})

	wg.Go(func() {
		podsList, cErrors[1] = repo.FetchPods(ctx, coreClient, pods.PodFilter{}, config.Name)
		if cErrors[1] != nil {
			cErrors[1] = wrapNodeFetchError("fetch pods", config.Name, cErrors[1])
		}
	})

	wg.Go(func() {
		nodeMetricsList, cErrors[2] = repo.FetchMetrics(ctx, metricsClient, nodemetrics.MetricsFilter{LabelSelector: config.Label}, config.Name)
		if cErrors[2] != nil {
			cErrors[2] = wrapNodeFetchError("fetch node metrics", config.Name, cErrors[2])
		}
	})

	wg.Wait()

	if err := errors.Join(cErrors...); err != nil {
		return nodeResources, err
	}

	nodeResources = merge(podsList, nodesList, nodeMetricsList)
	return nodeResources, nil
}

func wrapNodeFetchError(action string, name string, err error) error {
	if err == nil {
		return nil
	}
	if name == "" {
		return fmt.Errorf("%s: %w", action, err)
	}
	return fmt.Errorf("%s for node %q: %w", action, name, err)
}

package noderesources

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/pkg/client"
	"github.com/trezorg/k8spodsmetrics/pkg/nodemetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/nodes"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type Config struct {
	KubeConfig  string
	KubeContext string
	LogLevel    string
	Label       string
	KLogLevel   uint
	OnlyAlert   bool
}

func (config Config) request(ctx context.Context, client corev1.CoreV1Interface, metricsClient metricsv1beta1.MetricsV1beta1Interface) (NodeResourceList, error) {
	logger.Debug("Getting nodes info...")
	var nodeResources NodeResourceList
	errors := make([]error, 3)
	var podsList pods.PodResourceList
	var nodesList nodes.NodeList
	var nodeMetricsList nodemetrics.NodeMetricsList
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		nodesList, errors[0] = nodes.Nodes(ctx, client, nodes.NodeFilter{LabelSelector: config.Label})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		podsList, errors[1] = pods.Pods(ctx, client, pods.PodFilter{})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		nodeMetricsList, errors[2] = nodemetrics.Metrics(ctx, metricsClient, nodemetrics.MetricsFilter{LabelSelector: config.Label})
	}()

	wg.Wait()

	var mErrs *multierror.Error

	for _, err := range errors {
		if err != nil {
			if err := multierror.Append(mErrs, err); err != nil {
				logger.Error("MultiError append error", err)
			}
		}
	}

	if err := mErrs.ErrorOrNil(); err != nil {
		return nodeResources, err
	}
	nodeResources = merge(podsList, nodesList, nodeMetricsList)
	if config.OnlyAlert {
		nodeResources = nodeResources.filterAlerts()
	}
	return nodeResources, nil
}

func (config Config) Request(ctx context.Context) (NodeResourceList, error) {
	var err error
	logger.Debug("Preparing client...")
	podsClient, err := client.CoreV1Client(config.KubeConfig, config.KubeContext)
	if err != nil {
		return nil, err
	}
	metricsClient, err := client.MetricsClient(config.KubeConfig, config.KubeContext)
	if err != nil {
		return nil, err
	}
	return config.request(ctx, podsClient, metricsClient)
}

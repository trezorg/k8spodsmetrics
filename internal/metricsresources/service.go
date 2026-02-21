package metricsresources

import (
	"context"
	"fmt"

	"github.com/trezorg/k8spodsmetrics/internal/alert"
	"github.com/trezorg/k8spodsmetrics/internal/serviceorchestration"
	"github.com/trezorg/k8spodsmetrics/pkg/client"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type Config struct {
	KubeConfig    string
	KubeContext   string
	Namespaces    []string
	Label         string
	FieldSelector string
	Nodes         []string
	Output        string
	Sorting       string
	Resources     []string
	Alert         string
	WatchPeriod   uint
	Reverse       bool
	WatchMetrics  bool
}

type WatchResponse struct {
	error error
	data  PodMetricsResourceList
}

func (c Config) apiRequest(
	ctx context.Context,
	repo PodRepository,
	metricsClient metricsv1beta1.MetricsV1beta1Interface,
	podsClient corev1.CoreV1Interface,
) (PodMetricsResourceList, error) {
	fetchConfig := FetchConfig{
		Namespaces:    c.Namespaces,
		Label:         c.Label,
		FieldSelector: c.FieldSelector,
		Nodes:         c.Nodes,
	}
	podMetricsResourceList, err := FetchPodMetrics(ctx, repo, metricsClient, podsClient, fetchConfig)
	if err != nil {
		return nil, err
	}
	podMetricsResourceList = podMetricsResourceList.filterByAlert(alert.Alert(c.Alert))
	podMetricsResourceList = podMetricsResourceList.filterNodes(c.Nodes)
	podMetricsResourceList.sort(c.Sorting, c.Reverse)
	return podMetricsResourceList, nil
}

func (c Config) Request(ctx context.Context) (PodMetricsResourceList, error) {
	repo := NewPodRepository()
	request := func(
		requestContext context.Context,
		metricsClient metricsv1beta1.MetricsV1beta1Interface,
		coreClient corev1.CoreV1Interface,
	) (PodMetricsResourceList, error) {
		return c.apiRequest(requestContext, repo, metricsClient, coreClient)
	}

	return serviceorchestration.RequestWithClients(
		ctx,
		c.KubeConfig,
		c.KubeContext,
		client.Clients,
		request,
	)
}

func (c Config) Watch(ctx context.Context) chan WatchResponse {
	ch := make(chan WatchResponse, 1)
	repo := NewPodRepository()
	request := func(
		requestContext context.Context,
		metricsClient metricsv1beta1.MetricsV1beta1Interface,
		coreClient corev1.CoreV1Interface,
	) (PodMetricsResourceList, error) {
		return c.apiRequest(requestContext, repo, metricsClient, coreClient)
	}

	responses := serviceorchestration.WatchWithClients(
		ctx,
		c.KubeConfig,
		c.KubeContext,
		c.WatchPeriod,
		client.Clients,
		request,
	)

	go func() {
		defer close(ch)
		for response := range responses {
			if response.Error != nil {
				ch <- WatchResponse{error: response.Error}
				continue
			}
			ch <- WatchResponse{data: response.Data}
		}
	}()

	return ch
}

func (c *Config) prepare() error {
	if c.KubeConfig == "" {
		var err error
		c.KubeConfig, err = client.FindKubeConfig()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Config) Process(successProcessor SuccessProcessor) error {
	return serviceorchestration.RunWithPreparedContext(c.prepare, func(ctx context.Context) error {
		resources, err := c.Request(ctx)
		if err != nil {
			return fmt.Errorf("cannot get k8s resources: %w", err)
		}
		successProcessor.Success(resources)
		return nil
	})
}

func (c Config) ProcessWatch(successProcessor SuccessProcessor, errorProcessor ErrorProcessor) error {
	return serviceorchestration.RunWithPreparedContext(c.prepare, func(ctx context.Context) error {
		for resources := range c.Watch(ctx) {
			if resources.error != nil {
				errorProcessor.Error(resources.error)
			} else {
				successProcessor.Success(resources.data)
			}
		}
		return nil
	})
}

type SuccessProcessor interface {
	Success(PodMetricsResourceList)
}

type ErrorProcessor interface {
	Error(error)
}

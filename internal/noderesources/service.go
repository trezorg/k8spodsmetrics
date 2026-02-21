package noderesources

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
	KubeConfig   string
	KubeContext  string
	Label        string
	Name         string
	Output       string
	Sorting      string
	Alert        string
	Resources    []string
	WatchPeriod  uint
	Reverse      bool
	WatchMetrics bool
}

type WatchResponse struct {
	error error
	data  NodeResourceList
}

func (c Config) apiRequest(
	ctx context.Context,
	repo NodeRepository,
	coreClient corev1.CoreV1Interface,
	metricsClient metricsv1beta1.MetricsV1beta1Interface,
) (NodeResourceList, error) {
	fetchConfig := FetchConfig{
		Label: c.Label,
		Name:  c.Name,
	}
	nodeResources, err := FetchNodeMetrics(ctx, repo, coreClient, metricsClient, fetchConfig)
	if err != nil {
		return nil, err
	}
	nodeResources = nodeResources.filterByAlert(alert.Alert(c.Alert))
	nodeResources.sort(c.Sorting, c.Reverse)
	return nodeResources, nil
}

func (c Config) Request(ctx context.Context) (NodeResourceList, error) {
	repo := NewNodeRepository()
	request := func(
		requestContext context.Context,
		metricsClient metricsv1beta1.MetricsV1beta1Interface,
		coreClient corev1.CoreV1Interface,
	) (NodeResourceList, error) {
		return c.apiRequest(requestContext, repo, coreClient, metricsClient)
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
	repo := NewNodeRepository()
	request := func(
		requestContext context.Context,
		metricsClient metricsv1beta1.MetricsV1beta1Interface,
		coreClient corev1.CoreV1Interface,
	) (NodeResourceList, error) {
		return c.apiRequest(requestContext, repo, coreClient, metricsClient)
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
	Success(NodeResourceList)
}

type ErrorProcessor interface {
	Error(error)
}

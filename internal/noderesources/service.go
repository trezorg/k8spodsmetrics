package noderesources

import (
	"context"

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

type WatchResponse = serviceorchestration.WatchResponse[NodeResourceList]

func (c Config) apiRequest(
	ctx context.Context,
	repo NodeRepository,
	metricsClient metricsv1beta1.MetricsV1beta1Interface,
	coreClient corev1.CoreV1Interface,
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

func (c *Config) Request(ctx context.Context) (NodeResourceList, error) {
	return serviceorchestration.RequestWithRepo(
		ctx,
		c.KubeConfig,
		c.KubeContext,
		client.Clients,
		NewNodeRepository,
		c.apiRequest,
	)
}

func (c *Config) Watch(ctx context.Context) chan WatchResponse {
	return serviceorchestration.WatchWithRepo(
		ctx,
		c.KubeConfig,
		c.KubeContext,
		c.WatchPeriod,
		client.Clients,
		NewNodeRepository,
		c.apiRequest,
	)
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

func (c *Config) Process(successProcessor SuccessProcessor) error {
	return serviceorchestration.ProcessRequest(c.prepare, c.Request, successProcessor.Success)
}

func (c *Config) ProcessWatch(successProcessor SuccessProcessor, errorProcessor ErrorProcessor) error {
	return serviceorchestration.ProcessWatch(c.prepare, c.Watch, successProcessor.Success, errorProcessor.Error)
}

type SuccessProcessor interface {
	Success(NodeResourceList)
}

type ErrorProcessor interface {
	Error(error)
}

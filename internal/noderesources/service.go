package noderesources

import (
	"context"
	"errors"

	"github.com/trezorg/k8spodsmetrics/internal/alert"
	"github.com/trezorg/k8spodsmetrics/internal/serviceorchestration"
	sorting "github.com/trezorg/k8spodsmetrics/internal/sorting/noderesources"
	"github.com/trezorg/k8spodsmetrics/pkg/client"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type Config struct {
	KubeConfig  string
	KubeContext string
	Label       string
	Name        string
	Sorting     string
	Alert       string
	WatchPeriod uint
	Timeout     uint
	Reverse     bool
}

type WatchResponse = serviceorchestration.WatchResponse[NodeResourceList]

func (c Config) Validate() error {
	if err := alert.Valid(alert.Alert(c.Alert)); err != nil {
		return err
	}
	return sorting.Valid(sorting.Sorting(c.Sorting))
}

func (c Config) ValidateWatch() error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.WatchPeriod == 0 {
		return errors.New("watch period must be greater than 0")
	}
	return nil
}

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
		c.Timeout,
		client.Clients,
		NewNodeRepository,
		c.apiRequest,
	)
}

func (c *Config) Watch(ctx context.Context) <-chan WatchResponse {
	return serviceorchestration.WatchWithRepo(
		ctx,
		c.KubeConfig,
		c.KubeContext,
		c.WatchPeriod,
		c.Timeout,
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

func (c *Config) prepareRequest() error {
	if err := c.Validate(); err != nil {
		return err
	}
	return c.prepare()
}

func (c *Config) prepareWatch() error {
	if err := c.ValidateWatch(); err != nil {
		return err
	}
	return c.prepare()
}

func (c *Config) Process(successProcessor SuccessProcessor) error {
	return serviceorchestration.ProcessRequest(c.prepareRequest, c.Request, successProcessor.Success)
}

func (c *Config) ProcessWatch(successProcessor SuccessProcessor, errorProcessor ErrorProcessor) error {
	return serviceorchestration.ProcessWatch(c.prepareWatch, c.Watch, successProcessor.Success, errorProcessor.Error)
}

type SuccessProcessor interface {
	Success(NodeResourceList)
}

type ErrorProcessor interface {
	Error(error)
}

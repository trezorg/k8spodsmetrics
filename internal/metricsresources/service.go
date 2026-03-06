package metricsresources

import (
	"context"
	"errors"

	"github.com/trezorg/k8spodsmetrics/internal/alert"
	"github.com/trezorg/k8spodsmetrics/internal/serviceorchestration"
	sorting "github.com/trezorg/k8spodsmetrics/internal/sorting/metricsresources"
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
	Sorting       string
	Alert         string
	WatchPeriod   uint
	Timeout       uint
	Reverse       bool
}

type WatchResponse = serviceorchestration.WatchResponse[PodMetricsResourceList]

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

func (c *Config) Request(ctx context.Context) (PodMetricsResourceList, error) {
	return serviceorchestration.RequestWithRepo(
		ctx,
		c.KubeConfig,
		c.KubeContext,
		c.Timeout,
		client.Clients,
		NewPodRepository,
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
		NewPodRepository,
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
	Success(PodMetricsResourceList)
}

type ErrorProcessor interface {
	Error(error)
}

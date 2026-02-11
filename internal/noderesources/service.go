package noderesources

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"log/slog"

	"github.com/trezorg/k8spodsmetrics/internal/alert"
	"github.com/trezorg/k8spodsmetrics/pkg/client"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/klog/v2"
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
	KLogLevel    uint
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
	var err error
	slog.Debug("Preparing client...")
	metricsClient, podsClient, err := client.Clients(c.KubeConfig, c.KubeContext)
	if err != nil {
		return nil, err
	}
	repo := NewNodeRepository()
	return c.apiRequest(ctx, repo, podsClient, metricsClient)
}

func (c Config) Watch(ctx context.Context) chan WatchResponse {
	ch := make(chan WatchResponse, 1)
	slog.Debug("Preparing client...")

	go func() {
		defer close(ch)
		metricsClient, podsClient, err := client.Clients(c.KubeConfig, c.KubeContext)
		if err != nil {
			ch <- WatchResponse{error: err}
			return
		}
		repo := NewNodeRepository()
		p := func() {
			data, rErr := c.apiRequest(ctx, repo, podsClient, metricsClient)
			if rErr != nil {
				ch <- WatchResponse{error: rErr}
				return
			}
			ch <- WatchResponse{data: data}
		}

		p()

		ticker := time.NewTicker(time.Duration(c.WatchPeriod) * time.Second) //nolint:gosec // it is ok
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p()
			case <-ctx.Done():
				return
			}
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
	klog.InitFlags(nil)
	defer klog.Flush()
	return flag.Set("v", strconv.Itoa(int(c.KLogLevel))) //nolint:gosec // it is safe
}

func (c Config) Process(successProcessor SuccessProcessor) error {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer done()
	if err := c.prepare(); err != nil {
		return err
	}
	resources, err := c.Request(ctx)
	if err != nil {
		return fmt.Errorf("cannot get k8s resources: %w", err)
	}
	successProcessor.Success(resources)
	return nil
}

func (c Config) ProcessWatch(successProcessor SuccessProcessor, errorProcessor ErrorProcessor) error {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer done()
	if err := c.prepare(); err != nil {
		return err
	}
	for resources := range c.Watch(ctx) {
		if resources.error != nil {
			errorProcessor.Error(resources.error)
		} else {
			successProcessor.Success(resources.data)
		}
	}
	return nil
}

type SuccessProcessor interface {
	Success(NodeResourceList)
}

type ErrorProcessor interface {
	Error(error)
}

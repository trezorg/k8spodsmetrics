package metricsresources

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/trezorg/k8spodsmetrics/internal/alert"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/pkg/client"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/klog/v2"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type Config struct {
	KubeConfig   string
	KubeContext  string
	Namespace    string
	Label        string
	LogLevel     string
	Output       string
	Sorting      string
	Alert        string
	KLogLevel    uint
	WatchPeriod  uint
	Reverse      bool
	WatchMetrics bool
}

type WatchResponse struct {
	error error
	data  PodMetricsResourceList
}

func (c Config) request(ctx context.Context, metricsClient metricsv1beta1.MetricsV1beta1Interface, podsClient corev1.CoreV1Interface) (PodMetricsResourceList, error) {
	logger.Debug("Getting metrics...")
	var podMetricsResourceList PodMetricsResourceList
	c_errors := make([]error, 2)
	var podsList pods.PodResourceList
	var metricsList podmetrics.PodMetricList
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsList, c_errors[0] = podmetrics.Metrics(ctx, metricsClient, podmetrics.MetricFilter{
			Namespace:     c.Namespace,
			LabelSelector: c.Label,
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		podsList, c_errors[1] = pods.Pods(ctx, podsClient, pods.PodFilter{
			Namespace:     c.Namespace,
			LabelSelector: c.Label,
		}, "")
	}()

	wg.Wait()

	var rErr error

	for _, err := range c_errors {
		if err != nil {
			rErr = errors.Join(rErr, err)
		}
	}

	if rErr != nil {
		return podMetricsResourceList, rErr
	}

	podMetricsResourceList = merge(podsList, metricsList)
	podMetricsResourceList = podMetricsResourceList.filterByAlert(alert.Alert(c.Alert))
	podMetricsResourceList.sort(c.Sorting, c.Reverse)
	return podMetricsResourceList, nil
}

func (c Config) Request(ctx context.Context) (PodMetricsResourceList, error) {
	var err error
	logger.Debug("Preparing client...")
	metricsClient, podsClient, err := client.Clients(c.KubeConfig, c.KubeContext)
	if err != nil {
		return nil, err
	}
	return c.request(ctx, metricsClient, podsClient)
}

func (c Config) Watch(ctx context.Context) chan WatchResponse {
	ch := make(chan WatchResponse, 1)
	logger.Debug("Preparing client...")

	go func() {
		defer close(ch)
		metricsClient, podsClient, err := client.Clients(c.KubeConfig, c.KubeContext)
		if err != nil {
			ch <- WatchResponse{error: err}
			return
		}
		p := func() {
			data, err := c.request(ctx, metricsClient, podsClient)
			if err != nil {
				ch <- WatchResponse{error: err}
				return
			}
			ch <- WatchResponse{data: data}
		}

		p()

		ticker := time.NewTicker(time.Duration(c.WatchPeriod) * time.Second)
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
	logger.InitLogger(c.LogLevel)
	if c.KubeConfig == "" {
		var err error
		c.KubeConfig, err = client.FindKubeConfig()
		if err != nil {
			return err
		}
	}
	klog.InitFlags(nil)
	defer klog.Flush()
	if err := flag.Set("v", strconv.Itoa(int(c.KLogLevel))); err != nil {
		return err
	}
	return nil
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
	Success(PodMetricsResourceList)
}

type ErrorProcessor interface {
	Error(error)
}

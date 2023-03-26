package metricsresources

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/pkg/client"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type Config struct {
	KubeConfig   string
	KubeContext  string
	Namespace    string
	Label        string
	LogLevel     string
	KLogLevel    uint
	OnlyAlert    bool
	WatchMetrics bool
	WatchPeriod  uint
}

func (config Config) request(ctx context.Context, metricsClient metricsv1beta1.MetricsV1beta1Interface, podsClient corev1.CoreV1Interface) (PodMetricsResourceList, error) {
	logger.Debug("Getting metrics...")
	var podMetricsResourceList PodMetricsResourceList
	errors := make([]error, 2)
	var podsList pods.PodResourceList
	var metricsList podmetrics.PodMetricList
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsList, errors[0] = podmetrics.Metrics(ctx, metricsClient, podmetrics.MetricFilter{
			Namespace:     config.Namespace,
			LabelSelector: config.Label,
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		podsList, errors[1] = pods.Pods(ctx, podsClient, pods.PodFilter{
			Namespace:     config.Namespace,
			LabelSelector: config.Label,
		}, "")
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
		return podMetricsResourceList, err
	}

	podMetricsResourceList = merge(podsList, metricsList)
	if config.OnlyAlert {
		podMetricsResourceList = podMetricsResourceList.filterAlerts()
	}
	return podMetricsResourceList, nil
}

func (config Config) Request(ctx context.Context) (PodMetricsResourceList, error) {
	var err error
	logger.Debug("Preparing client...")
	metricsClient, podsClient, err := client.Clients(config.KubeConfig, config.KubeContext)
	if err != nil {
		return nil, err
	}
	return config.request(ctx, metricsClient, podsClient)
}

func (config Config) Watch(ctx context.Context, successFunc SuccessFunction, errorFunc ErrorFunction) {
	var err error
	logger.Debug("Preparing client...")
	metricsClient, podsClient, err := client.Clients(config.KubeConfig, config.KubeContext)
	if err != nil {
		errorFunc(nil)
		return
	}
	ticker := time.NewTicker(time.Duration(config.WatchPeriod) * time.Second)
	defer ticker.Stop()

	p := func() {
		r, err := config.request(ctx, metricsClient, podsClient)
		if err != nil {
			errorFunc(err)
			return
		}
		successFunc(r)
	}

	p()

	for {
		select {
		case <-ticker.C:
			p()
		case <-ctx.Done():
			return
		}
	}

}

package serviceorchestration

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"log/slog"

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

type ClientsFactory func(string, string) (metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface, error)

type RequestFunc[T any] func(context.Context, metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface) (T, error)

type WatchResponse[T any] struct {
	Error error
	Data  T
}

func RequestWithClients[T any](
	ctx context.Context,
	kubeConfig string,
	kubeContext string,
	clientsFactory ClientsFactory,
	request RequestFunc[T],
) (T, error) {
	var zero T

	slog.Debug("Preparing client...")
	metricsClient, coreClient, err := clientsFactory(kubeConfig, kubeContext)
	if err != nil {
		return zero, err
	}

	return request(ctx, metricsClient, coreClient)
}

func WatchWithClients[T any](
	ctx context.Context,
	kubeConfig string,
	kubeContext string,
	watchPeriodSeconds uint,
	clientsFactory ClientsFactory,
	request RequestFunc[T],
) chan WatchResponse[T] {
	ch := make(chan WatchResponse[T], 1)
	slog.Debug("Preparing client...")

	go func() {
		defer close(ch)

		watchPeriod, err := time.ParseDuration(strconv.FormatUint(uint64(watchPeriodSeconds), 10) + "s")
		if err != nil {
			ch <- WatchResponse[T]{Error: err}
			return
		}

		metricsClient, coreClient, err := clientsFactory(kubeConfig, kubeContext)
		if err != nil {
			ch <- WatchResponse[T]{Error: err}
			return
		}

		produce := func() {
			data, requestErr := request(ctx, metricsClient, coreClient)
			if requestErr != nil {
				ch <- WatchResponse[T]{Error: requestErr}
				return
			}

			ch <- WatchResponse[T]{Data: data}
		}

		produce()

		ticker := time.NewTicker(watchPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				produce()
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch
}

func RunWithPreparedContext(prepare func() error, run func(context.Context) error) error {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer done()

	if err := prepare(); err != nil {
		return err
	}

	return run(ctx)
}

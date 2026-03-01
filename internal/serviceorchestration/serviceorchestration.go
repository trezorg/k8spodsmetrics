package serviceorchestration

import (
	"context"
	"errors"
	"fmt"
	"math"
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
type RepoRequestFunc[T any, R any] func(context.Context, R, metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface) (T, error)

type WatchResponse[T any] struct {
	Error error
	Data  T
}

var (
	ErrSignalCanceled = errors.New("request canceled by signal")
	ErrRequestTimeout = errors.New("request timed out")
)

func RequestWithClients[T any](
	ctx context.Context,
	kubeConfig string,
	kubeContext string,
	timeout uint,
	clientsFactory ClientsFactory,
	request RequestFunc[T],
) (T, error) {
	var zero T

	if timeout > 0 {
		if timeout > uint(math.MaxInt64/int64(time.Second)) {
			return zero, fmt.Errorf("timeout is too large: %d", timeout)
		}
		timeoutDuration := time.Duration(int64(timeout)) * time.Second
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeoutCause(ctx, timeoutDuration, ErrRequestTimeout)
		defer cancel()
	}

	slog.Debug("Preparing client...")
	metricsClient, coreClient, err := clientsFactory(kubeConfig, kubeContext)
	if err != nil {
		return zero, err
	}

	return request(ctx, metricsClient, coreClient)
}

func RequestWithRepo[T any, R any](
	ctx context.Context,
	kubeConfig string,
	kubeContext string,
	timeout uint,
	clientsFactory ClientsFactory,
	repoFactory func() R,
	request RepoRequestFunc[T, R],
) (T, error) {
	repo := repoFactory()
	requestWithRepo := func(
		requestContext context.Context,
		metricsClient metricsv1beta1.MetricsV1beta1Interface,
		coreClient corev1.CoreV1Interface,
	) (T, error) {
		return request(requestContext, repo, metricsClient, coreClient)
	}

	return RequestWithClients(ctx, kubeConfig, kubeContext, timeout, clientsFactory, requestWithRepo)
}

func WatchWithClients[T any](
	ctx context.Context,
	kubeConfig string,
	kubeContext string,
	watchPeriodSeconds uint,
	timeout uint,
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

		var timeoutDuration time.Duration
		if timeout > 0 {
			if timeout > uint(math.MaxInt64/int64(time.Second)) {
				ch <- WatchResponse[T]{Error: fmt.Errorf("timeout is too large: %d", timeout)}
				return
			}
			timeoutDuration = time.Duration(int64(timeout)) * time.Second
		}

		produce := func() {
			requestCtx := ctx
			if timeout > 0 {
				var cancel context.CancelFunc
				requestCtx, cancel = context.WithTimeoutCause(ctx, timeoutDuration, ErrRequestTimeout)
				defer cancel()
			}
			data, requestErr := request(requestCtx, metricsClient, coreClient)
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

func WatchWithRepo[T any, R any](
	ctx context.Context,
	kubeConfig string,
	kubeContext string,
	watchPeriodSeconds uint,
	timeout uint,
	clientsFactory ClientsFactory,
	repoFactory func() R,
	request RepoRequestFunc[T, R],
) chan WatchResponse[T] {
	repo := repoFactory()
	requestWithRepo := func(
		requestContext context.Context,
		metricsClient metricsv1beta1.MetricsV1beta1Interface,
		coreClient corev1.CoreV1Interface,
	) (T, error) {
		return request(requestContext, repo, metricsClient, coreClient)
	}

	return WatchWithClients(ctx, kubeConfig, kubeContext, watchPeriodSeconds, timeout, clientsFactory, requestWithRepo)
}

func RunWithPreparedContext(prepare func() error, run func(context.Context) error) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	if err := prepare(); err != nil {
		return err
	}

	return run(ctx)
}

func ProcessRequest[T any](
	prepare func() error,
	request func(context.Context) (T, error),
	successProcessor func(T),
) error {
	return RunWithPreparedContext(prepare, func(ctx context.Context) error {
		resources, err := request(ctx)
		if err != nil {
			return fmt.Errorf("cannot get k8s resources: %w", err)
		}
		successProcessor(resources)
		return nil
	})
}

func ProcessWatch[T any](
	prepare func() error,
	watch func(context.Context) chan WatchResponse[T],
	successProcessor func(T),
	errorProcessor func(error),
) error {
	return RunWithPreparedContext(prepare, func(ctx context.Context) error {
		for resources := range watch(ctx) {
			if resources.Error != nil {
				errorProcessor(resources.Error)
			} else {
				successProcessor(resources.Data)
			}
		}
		return nil
	})
}

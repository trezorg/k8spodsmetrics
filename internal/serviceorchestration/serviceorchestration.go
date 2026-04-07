package serviceorchestration

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
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

func durationFromSeconds(seconds uint, field string) (time.Duration, error) {
	if seconds > uint(math.MaxInt64/int64(time.Second)) {
		return 0, fmt.Errorf("%s is too large: %d", field, seconds)
	}

	return time.Duration(seconds) * time.Second, nil
}

func withSignalCause(parent context.Context, signals chan os.Signal) (context.Context, context.CancelCauseFunc) {
	ctx, cancel := context.WithCancelCause(parent)

	go func() {
		select {
		case <-ctx.Done():
		case <-signals:
			signal.Stop(signals)
			cancel(ErrSignalCanceled)
		}
	}()

	return ctx, cancel
}

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
		timeoutDuration, err := durationFromSeconds(timeout, "timeout")
		if err != nil {
			return zero, err
		}
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
) <-chan WatchResponse[T] {
	ch := make(chan WatchResponse[T], 1)
	slog.Debug("Preparing client...")

	go func() {
		defer close(ch)

		watchPeriod, err := durationFromSeconds(watchPeriodSeconds, "watch period")
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
			timeoutDuration, err = durationFromSeconds(timeout, "timeout")
			if err != nil {
				ch <- WatchResponse[T]{Error: err}
				return
			}
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
) <-chan WatchResponse[T] {
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
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, shutdownSignals...)
	defer signal.Stop(signals)

	ctx, cancel := withSignalCause(context.Background(), signals)
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
	watch func(context.Context) <-chan WatchResponse[T],
	successProcessor func(T),
	errorProcessor func(error),
) error {
	return RunWithPreparedContext(prepare, func(ctx context.Context) error {
		lastErrorFingerprint := ""
		for resources := range watch(ctx) {
			if resources.Error != nil {
				fingerprint := watchErrorFingerprint(resources.Error)
				if fingerprint != lastErrorFingerprint {
					errorProcessor(resources.Error)
					lastErrorFingerprint = fingerprint
				}
			} else {
				lastErrorFingerprint = ""
				successProcessor(resources.Data)
			}
		}
		return nil
	})
}

func watchErrorFingerprint(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%T:%s", err, err.Error())
}

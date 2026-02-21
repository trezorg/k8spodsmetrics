package serviceorchestration

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

func TestRequestWithClients(t *testing.T) {
	t.Run("returns clients factory error", func(t *testing.T) {
		expectedErr := errors.New("cannot create clients")
		_, err := RequestWithClients(
			context.Background(),
			"config",
			"context",
			func(string, string) (metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface, error) {
				return nil, nil, expectedErr
			},
			func(context.Context, metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface) (int, error) {
				return 0, nil
			},
		)

		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("returns request result", func(t *testing.T) {
		result, err := RequestWithClients(
			context.Background(),
			"config",
			"context",
			func(string, string) (metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface, error) {
				return nil, nil, nil
			},
			func(context.Context, metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface) (int, error) {
				return 42, nil
			},
		)

		require.NoError(t, err)
		require.Equal(t, 42, result)
	})
}

func TestWatchWithClients(t *testing.T) {
	t.Run("returns clients factory error", func(t *testing.T) {
		expectedErr := errors.New("cannot create clients")
		responses := WatchWithClients(
			context.Background(),
			"config",
			"context",
			1,
			func(string, string) (metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface, error) {
				return nil, nil, expectedErr
			},
			func(context.Context, metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface) (int, error) {
				return 0, nil
			},
		)

		response, ok := <-responses
		require.True(t, ok)
		require.ErrorIs(t, response.Error, expectedErr)
		require.Zero(t, response.Data)

		_, isOpen := <-responses
		require.False(t, isOpen)
	})

	t.Run("publishes first and ticker responses", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		calls := 0
		responses := WatchWithClients(
			ctx,
			"config",
			"context",
			1,
			func(string, string) (metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface, error) {
				return nil, nil, nil
			},
			func(context.Context, metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface) (int, error) {
				calls++
				if calls == 2 {
					cancel()
				}
				return calls, nil
			},
		)

		values := []int{}
		for response := range responses {
			require.NoError(t, response.Error)
			values = append(values, response.Data)
		}

		require.Equal(t, []int{1, 2}, values)
	})
}

func TestRunWithPreparedContext(t *testing.T) {
	t.Run("returns prepare error", func(t *testing.T) {
		expectedErr := errors.New("prepare error")
		runCalled := false

		err := RunWithPreparedContext(
			func() error { return expectedErr },
			func(context.Context) error {
				runCalled = true
				return nil
			},
		)

		require.ErrorIs(t, err, expectedErr)
		require.False(t, runCalled)
	})

	t.Run("executes run when prepared", func(t *testing.T) {
		err := RunWithPreparedContext(
			func() error { return nil },
			func(ctx context.Context) error {
				require.NotNil(t, ctx)
				return nil
			},
		)

		require.NoError(t, err)
	})
}

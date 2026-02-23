package metricsresources

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

type mockSuccessProcessor struct {
	called bool
}

func (m *mockSuccessProcessor) Success(list metricsresources.PodMetricsResourceList) {
	m.called = true
}

type mockErrorProcessor struct {
	called bool
}

func (m *mockErrorProcessor) Error(err error) {
	m.called = true
}

func TestNewScreenSuccessWriter(t *testing.T) {
	t.Run("creates success writer", func(t *testing.T) {
		mock := &mockSuccessProcessor{}
		writer := NewScreenSuccessWriter(mock)
		require.NotNil(t, writer)

		list := metricsresources.PodMetricsResourceList{
			{
				PodResource: pods.PodResource{
					NamespaceName: pods.NamespaceName{
						Name:      "test-pod",
						Namespace: "default",
					},
				},
				PodMetric: podmetrics.PodMetric{
					Name:       "test-pod",
					Namespace:  "default",
					Containers: []podmetrics.ContainerMetric{},
				},
			},
		}

		writer.Success(list)
		require.True(t, mock.called)
	})
}

func TestNewScreenErrorWriter(t *testing.T) {
	t.Run("creates error writer", func(t *testing.T) {
		mock := &mockErrorProcessor{}
		writer := NewScreenErrorWriter(mock)
		require.NotNil(t, writer)

		err := errors.New("test error")
		writer.Error(err)
		require.True(t, mock.called)
	})
}

func TestScreenSuccessWriter_Success(t *testing.T) {
	t.Run("calls underlying writer", func(t *testing.T) {
		mock := &mockSuccessProcessor{}
		writer := NewScreenSuccessWriter(mock)
		writer.Success(metricsresources.PodMetricsResourceList{})
		require.True(t, mock.called)
	})
}

func TestScreenErrorWriter_Error(t *testing.T) {
	t.Run("calls underlying writer", func(t *testing.T) {
		mock := &mockErrorProcessor{}
		writer := NewScreenErrorWriter(mock)
		err := errors.New("test error")
		writer.Error(err)
		require.True(t, mock.called)
	})
}

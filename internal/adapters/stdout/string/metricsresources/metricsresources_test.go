package metricsresources

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func TestPrint(t *testing.T) {
	t.Run("prints string", func(t *testing.T) {
		list := metricsresources.PodMetricsResourceList{
			{
				PodResource: pods.PodResource{
					NamespaceName: pods.NamespaceName{
						Name:      "test-pod",
						Namespace: "default",
					},
					NodeName: "node-1",
				},
				PodMetric: podmetrics.PodMetric{
					Name:       "test-pod",
					Namespace:  "default",
					Containers: []podmetrics.ContainerMetric{},
				},
			},
		}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		Print(list)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		require.NotEmpty(t, output)
		require.Contains(t, output, "Name:")
		require.Contains(t, output, "test-pod")
		require.Contains(t, output, "Namespace:")
		require.Contains(t, output, "default")
	})

	t.Run("prints empty list", func(t *testing.T) {
		list := metricsresources.PodMetricsResourceList{}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		Print(list)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		require.NotEmpty(t, output)
	})
}

func TestString_Success(t *testing.T) {
	t.Run("calls Print", func(t *testing.T) {
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

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		formatter := String(Print)
		formatter.Success(list)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		require.NotEmpty(t, output)
	})
}

func TestString_Error(t *testing.T) {
	t.Run("logs error without panicking", func(t *testing.T) {
		formatter := String(Print)
		err := errors.New("test error")
		require.NotPanics(t, func() {
			formatter.Error(err)
		})
	})
}

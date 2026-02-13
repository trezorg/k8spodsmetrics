package metricsresources

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func TestPrint(t *testing.T) {
	t.Run("prints valid JSON", func(t *testing.T) {
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

		var decoded metricsresources.PodMetricsResourceOutputEnvelope
		err := json.Unmarshal([]byte(output), &decoded)
		require.NoError(t, err)
		require.Len(t, decoded.Items, 1)
		require.Equal(t, "test-pod", decoded.Items[0].Name)
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

		var decoded metricsresources.PodMetricsResourceOutputEnvelope
		err := json.Unmarshal([]byte(output), &decoded)
		require.NoError(t, err)
		require.Empty(t, decoded.Items)
	})
}

func TestJSON_Success(t *testing.T) {
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

		formatter := JSON(Print)
		formatter.Success(list)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		var decoded metricsresources.PodMetricsResourceOutputEnvelope
		err := json.Unmarshal([]byte(output), &decoded)
		require.NoError(t, err)
		require.Len(t, decoded.Items, 1)
	})
}

func TestJSON_Error(t *testing.T) {
	t.Run("logs error without panicking", func(t *testing.T) {
		formatter := JSON(Print)
		err := errors.New("test error")
		require.NotPanics(t, func() {
			formatter.Error(err)
		})
	})
}

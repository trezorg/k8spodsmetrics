package metricsresources

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func posResourceList(name, namespace, container string) pods.PodResourceList {
	return []pods.PodResource{
		{
			NamespaceName: pods.NamespaceName{
				Name:      name,
				Namespace: namespace,
			},
			Containers: []pods.ContainerResource{
				{
					Name: container,
					Limits: pods.Resource{
						CPU:    1,
						Memory: 1024,
					},
					Requests: pods.Resource{
						CPU:    1,
						Memory: 1024,
					},
				},
			},
		},
	}
}

func posMetricsList(name, namespace, container string) podmetrics.PodMetricList {
	return []podmetrics.PodMetric{
		{
			Name:      name,
			Namespace: namespace,
			Containers: []podmetrics.ContainerMetric{
				{
					Name: container,
					Metric: podmetrics.Metric{
						CPU:    2000,
						Memory: 2048,
					},
				},
			},
		},
	}
}

func TestMergeSameNamespaceAndName(t *testing.T) {
	logger.InitDefaultLogger()
	podResourceList := posResourceList("foo", "bar", "foo-container")
	podMetricList := posMetricsList("foo", "bar", "foo-container")
	pods := merge(podResourceList, podMetricList)
	require.Len(t, pods, 1)
	for _, pod := range pods {
		require.Len(t, pod.PodResource.Containers, 1)
		require.Contains(t, pod.String(), "/", pod.String())
	}
	require.Contains(t, pods.String(), "/", pods.String())
}

func TestMergeDifferentNamespaceAndName(t *testing.T) {
	logger.InitDefaultLogger()
	podResourceList := posResourceList("foo", "bar", "foo-container")
	podMetricList := posMetricsList("foo1", "bar", "foo-container")
	pods := merge(podResourceList, podMetricList)
	require.Len(t, pods, 1)
	for _, pod := range pods {
		require.NotContains(t, pod.String(), "/", pod.String())
	}
	require.NotContains(t, pods.String(), "/", pods.String())
}

func TestStringify(t *testing.T) {
	logger.InitDefaultLogger()
	podResourceList := posResourceList("foo", "bar", "foo-container")
	podMetricList := []podmetrics.PodMetric{}
	pods := merge(podResourceList, podMetricList)
	require.Len(t, pods, 1)
	text := pods[0].String()
	require.Greater(t, len(text), 0)
	require.NotContains(t, text, "/", text)
	require.NotContains(t, pods.String(), "/", pods.String())
}

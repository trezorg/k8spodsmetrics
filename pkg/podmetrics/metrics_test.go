package podmetrics

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktesting "k8s.io/client-go/testing"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsfake "k8s.io/metrics/pkg/client/clientset/versioned/fake"
)

func TestMetricFilter(t *testing.T) {
	t.Run("default filter", func(t *testing.T) {
		filter := MetricFilter{}
		require.Empty(t, filter.Namespaces)
		require.Empty(t, filter.LabelSelector)
		require.Empty(t, filter.FieldSelector)
	})

	t.Run("with namespaces", func(t *testing.T) {
		filter := MetricFilter{Namespaces: []string{"test-ns"}}
		require.Equal(t, []string{"test-ns"}, filter.Namespaces)
	})

	t.Run("with multiple namespaces", func(t *testing.T) {
		filter := MetricFilter{Namespaces: []string{"ns1", "ns2", "ns3"}}
		require.Equal(t, []string{"ns1", "ns2", "ns3"}, filter.Namespaces)
	})

	t.Run("with label selector", func(t *testing.T) {
		filter := MetricFilter{LabelSelector: "app=test"}
		require.Equal(t, "app=test", filter.LabelSelector)
	})

	t.Run("with field selector", func(t *testing.T) {
		filter := MetricFilter{FieldSelector: "spec.nodeName=node1"}
		require.Equal(t, "spec.nodeName=node1", filter.FieldSelector)
	})
}

func TestMetric(t *testing.T) {
	t.Run("default metric", func(t *testing.T) {
		metric := Metric{}
		require.Equal(t, int64(0), metric.CPU)
		require.Equal(t, int64(0), metric.Memory)
	})

	t.Run("with values", func(t *testing.T) {
		metric := Metric{
			CPU:              1000,
			Memory:           1024 * 1024,
			Storage:          2048 * 1024,
			StorageEphemeral: 512 * 1024,
		}
		require.Equal(t, int64(1000), metric.CPU)
		require.Equal(t, int64(1024*1024), metric.Memory)
		require.Equal(t, int64(2048*1024), metric.Storage)
		require.Equal(t, int64(512*1024), metric.StorageEphemeral)
	})
}

func TestContainerMetric(t *testing.T) {
	t.Run("default metric", func(t *testing.T) {
		metric := ContainerMetric{}
		require.Empty(t, metric.Name)
		require.Equal(t, int64(0), metric.CPU)
	})

	t.Run("with values", func(t *testing.T) {
		metric := ContainerMetric{
			Name:   "container1",
			Metric: Metric{CPU: 1000, Memory: 1024 * 1024},
		}
		require.Equal(t, "container1", metric.Name)
		require.Equal(t, int64(1000), metric.CPU)
		require.Equal(t, int64(1024*1024), metric.Memory)
	})
}

func TestPodMetric(t *testing.T) {
	t.Run("default metric", func(t *testing.T) {
		metric := PodMetric{}
		require.Empty(t, metric.Namespace)
		require.Empty(t, metric.Name)
		require.Empty(t, metric.Containers)
	})

	t.Run("with values", func(t *testing.T) {
		metric := PodMetric{
			Namespace: "test-ns",
			Name:      "test-pod",
			Containers: []ContainerMetric{
				{Name: "c1", Metric: Metric{CPU: 100, Memory: 1024}},
			},
		}
		require.Equal(t, "test-ns", metric.Namespace)
		require.Equal(t, "test-pod", metric.Name)
		require.Len(t, metric.Containers, 1)
	})
}

func TestPodMetricList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		list := PodMetricList{}
		require.Empty(t, list)
	})

	t.Run("with metrics", func(t *testing.T) {
		list := PodMetricList{
			{Namespace: "ns1", Name: "pod1"},
			{Namespace: "ns2", Name: "pod2"},
		}
		require.Len(t, list, 2)
	})
}

func TestListMetricsFollowsPagination(t *testing.T) {
	ctx := t.Context()
	client := metricsfake.NewSimpleClientset()
	type listOptionsGetter interface {
		GetListOptions() metav1.ListOptions
	}

	client.PrependReactor("list", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		if listAction.GetListOptions().Continue != "" {
			return false, nil, nil
		}

		return true, &metricsv1beta1.PodMetricsList{
			ListMeta: metav1.ListMeta{Continue: "page-2"},
			Items: []metricsv1beta1.PodMetrics{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "pod-1", Namespace: "default"},
					Containers: []metricsv1beta1.ContainerMetrics{
						{
							Name: "app",
							Usage: v1.ResourceList{
								v1.ResourceCPU: resource.MustParse("100m"),
							},
						},
					},
				},
			},
		}, nil
	})
	client.PrependReactor("list", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		if listAction.GetListOptions().Continue != "page-2" {
			return false, nil, nil
		}

		return true, &metricsv1beta1.PodMetricsList{
			Items: []metricsv1beta1.PodMetrics{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "pod-2", Namespace: "default"},
					Containers: []metricsv1beta1.ContainerMetrics{
						{
							Name: "app",
							Usage: v1.ResourceList{
								v1.ResourceCPU: resource.MustParse("200m"),
							},
						},
					},
				},
			},
		}, nil
	})

	result, err := listMetrics(ctx, client.MetricsV1beta1(), MetricFilter{}, "default")
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, "pod-1", result[0].Name)
	require.Equal(t, "pod-2", result[1].Name)
}

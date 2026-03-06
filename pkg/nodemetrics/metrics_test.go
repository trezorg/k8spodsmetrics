package nodemetrics

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

func TestNodeMetric(t *testing.T) {
	t.Run("default metric", func(t *testing.T) {
		metric := NodeMetric{}
		require.Empty(t, metric.Name)
		require.Equal(t, int64(0), metric.CPU)
		require.Equal(t, int64(0), metric.Memory)
		require.Equal(t, int64(0), metric.Storage)
		require.Equal(t, int64(0), metric.StorageEphemeral)
	})

	t.Run("with values", func(t *testing.T) {
		metric := NodeMetric{
			Name:             "node-1",
			CPU:              2000,
			Memory:           4 * 1024 * 1024 * 1024,
			Storage:          10 * 1024 * 1024 * 1024,
			StorageEphemeral: 5 * 1024 * 1024 * 1024,
		}
		require.Equal(t, "node-1", metric.Name)
		require.Equal(t, int64(2000), metric.CPU)
		require.Equal(t, int64(4*1024*1024*1024), metric.Memory)
		require.Equal(t, int64(10*1024*1024*1024), metric.Storage)
		require.Equal(t, int64(5*1024*1024*1024), metric.StorageEphemeral)
	})
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		list := List{}
		require.Empty(t, list)
	})

	t.Run("with metrics", func(t *testing.T) {
		list := List{
			{Name: "node-1", CPU: 2000, Memory: 4 * 1024 * 1024 * 1024},
			{Name: "node-2", CPU: 1000, Memory: 2 * 1024 * 1024 * 1024},
		}
		require.Len(t, list, 2)
		require.Equal(t, "node-1", list[0].Name)
		require.Equal(t, "node-2", list[1].Name)
	})
}

func TestMetricsFilter(t *testing.T) {
	t.Run("default filter", func(t *testing.T) {
		filter := MetricsFilter{}
		require.Empty(t, filter.LabelSelector)
		require.Empty(t, filter.FieldSelector)
	})

	t.Run("with label selector", func(t *testing.T) {
		filter := MetricsFilter{LabelSelector: "env=prod"}
		require.Equal(t, "env=prod", filter.LabelSelector)
	})

	t.Run("with field selector", func(t *testing.T) {
		filter := MetricsFilter{FieldSelector: "spec.unschedulable=true"}
		require.Equal(t, "spec.unschedulable=true", filter.FieldSelector)
	})
}

func TestMetricsFollowsPagination(t *testing.T) {
	ctx := t.Context()
	client := metricsfake.NewSimpleClientset()
	type listOptionsGetter interface {
		GetListOptions() metav1.ListOptions
	}

	client.PrependReactor("list", "nodes", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		if listAction.GetListOptions().Continue != "" {
			return false, nil, nil
		}

		return true, &metricsv1beta1.NodeMetricsList{
			ListMeta: metav1.ListMeta{Continue: "page-2"},
			Items: []metricsv1beta1.NodeMetrics{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "node-1"},
					Usage: v1.ResourceList{
						v1.ResourceCPU: resource.MustParse("100m"),
					},
				},
			},
		}, nil
	})
	client.PrependReactor("list", "nodes", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		if listAction.GetListOptions().Continue != "page-2" {
			return false, nil, nil
		}

		return true, &metricsv1beta1.NodeMetricsList{
			Items: []metricsv1beta1.NodeMetrics{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "node-2"},
					Usage: v1.ResourceList{
						v1.ResourceCPU: resource.MustParse("200m"),
					},
				},
			},
		}, nil
	})

	result, err := Metrics(ctx, client.MetricsV1beta1(), MetricsFilter{}, "")
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, "node-1", result[0].Name)
	require.Equal(t, "node-2", result[1].Name)
}

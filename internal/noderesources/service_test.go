package noderesources

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	corefake "k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsfake "k8s.io/metrics/pkg/client/clientset/versioned/fake"
)

func TestConfig(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		config := Config{}
		require.Empty(t, config.KubeConfig)
		require.Empty(t, config.Name)
	})

	t.Run("with values", func(t *testing.T) {
		config := Config{
			KubeConfig:  "/path/to/config",
			KubeContext: "test-context",
			Name:        "node1",
			Label:       "node-role.kubernetes.io/worker",
			Sorting:     "name",
			Reverse:     true,
			Alert:       "memory",
			WatchPeriod: 10,
		}
		require.Equal(t, "/path/to/config", config.KubeConfig)
		require.Equal(t, "test-context", config.KubeContext)
		require.Equal(t, "node1", config.Name)
	})
}

func TestWatchResponse(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		resp := WatchResponse{Error: errors.New("test error")}
		require.Error(t, resp.Error)
		require.Empty(t, resp.Data)
	})

	t.Run("with data", func(t *testing.T) {
		data := NodeResourceList{{Name: "node1"}}
		resp := WatchResponse{Data: data}
		require.NoError(t, resp.Error)
		require.Len(t, resp.Data, 1)
	})
}

func TestConfigValidate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := Config{Sorting: "name", Alert: "none"}
		require.NoError(t, cfg.Validate())
	})

	t.Run("invalid sorting", func(t *testing.T) {
		cfg := Config{Sorting: "invalid", Alert: "none"}
		require.ErrorContains(t, cfg.Validate(), "sorting should be one of")
	})

	t.Run("invalid alert", func(t *testing.T) {
		cfg := Config{Sorting: "name", Alert: "invalid"}
		require.ErrorContains(t, cfg.Validate(), "alert should be one of")
	})
}

func TestConfigValidateWatch(t *testing.T) {
	t.Run("zero watch period", func(t *testing.T) {
		cfg := Config{Sorting: "name", Alert: "none", WatchPeriod: 0}
		require.ErrorContains(t, cfg.ValidateWatch(), "watch period must be greater than 0")
	})
}

type noopSuccessProcessor struct{}

func (noopSuccessProcessor) Success(NodeResourceList) {}

type noopErrorProcessor struct{}

func (noopErrorProcessor) Error(error) {}

func TestProcessValidationError(t *testing.T) {
	cfg := Config{KubeConfig: "dummy", Sorting: "invalid", Alert: "none"}

	err := cfg.Process(noopSuccessProcessor{})
	require.ErrorContains(t, err, "sorting should be one of")
}

func TestProcessWatchValidationError(t *testing.T) {
	cfg := Config{KubeConfig: "dummy", Sorting: "name", Alert: "none", WatchPeriod: 0}

	err := cfg.ProcessWatch(noopSuccessProcessor{}, noopErrorProcessor{})
	require.ErrorContains(t, err, "watch period must be greater than 0")
}

func TestFetchNodeMetricsFollowsPagination(t *testing.T) {
	ctx := t.Context()
	coreClient := corefake.NewSimpleClientset()
	metricsClient := metricsfake.NewSimpleClientset()

	type listOptionsGetter interface {
		GetListOptions() metav1.ListOptions
	}

	coreClient.PrependReactor("list", "nodes", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		require.Equal(t, "", listAction.GetListOptions().Continue)

		return true, &v1.NodeList{
			ListMeta: metav1.ListMeta{Continue: "page-2"},
			Items: []v1.Node{{
				ObjectMeta: metav1.ObjectMeta{Name: "node-1"},
				Status: v1.NodeStatus{
					Capacity:    v1.ResourceList{v1.ResourceCPU: resource.MustParse("4")},
					Allocatable: v1.ResourceList{v1.ResourceCPU: resource.MustParse("3")},
				},
			}},
		}, nil
	})
	coreClient.PrependReactor("list", "nodes", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		if listAction.GetListOptions().Continue != "page-2" {
			return false, nil, nil
		}

		return true, &v1.NodeList{
			Items: []v1.Node{{
				ObjectMeta: metav1.ObjectMeta{Name: "node-2"},
				Status: v1.NodeStatus{
					Capacity:    v1.ResourceList{v1.ResourceCPU: resource.MustParse("2")},
					Allocatable: v1.ResourceList{v1.ResourceCPU: resource.MustParse("1500m")},
				},
			}},
		}, nil
	})

	coreClient.PrependReactor("list", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		require.Equal(t, "", listAction.GetListOptions().Continue)

		return true, &v1.PodList{
			ListMeta: metav1.ListMeta{Continue: "page-2"},
			Items: []v1.Pod{{
				ObjectMeta: metav1.ObjectMeta{Name: "pod-1", Namespace: "default"},
				Spec: v1.PodSpec{
					NodeName: "node-1",
					Containers: []v1.Container{{
						Name: "app",
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("100m")},
						},
					}},
				},
			}},
		}, nil
	})
	coreClient.PrependReactor("list", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		if listAction.GetListOptions().Continue != "page-2" {
			return false, nil, nil
		}

		return true, &v1.PodList{
			Items: []v1.Pod{{
				ObjectMeta: metav1.ObjectMeta{Name: "pod-2", Namespace: "default"},
				Spec: v1.PodSpec{
					NodeName: "node-2",
					Containers: []v1.Container{{
						Name: "app",
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("200m")},
						},
					}},
				},
			}},
		}, nil
	})

	metricsClient.PrependReactor("list", "nodes", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		require.Equal(t, "", listAction.GetListOptions().Continue)

		return true, &metricsv1beta1.NodeMetricsList{
			ListMeta: metav1.ListMeta{Continue: "page-2"},
			Items: []metricsv1beta1.NodeMetrics{{
				ObjectMeta: metav1.ObjectMeta{Name: "node-1"},
				Usage:      v1.ResourceList{v1.ResourceCPU: resource.MustParse("300m")},
			}},
		}, nil
	})
	metricsClient.PrependReactor("list", "nodes", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		if listAction.GetListOptions().Continue != "page-2" {
			return false, nil, nil
		}

		return true, &metricsv1beta1.NodeMetricsList{
			Items: []metricsv1beta1.NodeMetrics{{
				ObjectMeta: metav1.ObjectMeta{Name: "node-2"},
				Usage:      v1.ResourceList{v1.ResourceCPU: resource.MustParse("400m")},
			}},
		}, nil
	})

	result, err := FetchNodeMetrics(ctx, NewNodeRepository(), coreClient.CoreV1(), metricsClient.MetricsV1beta1(), FetchConfig{})
	require.NoError(t, err)
	require.Len(t, result, 2)

	byName := make(map[string]NodeResource, len(result))
	for _, item := range result {
		byName[item.Name] = item
	}

	require.Equal(t, int64(100), byName["node-1"].CPURequest)
	require.Equal(t, int64(300), byName["node-1"].UsedCPU)
	require.Equal(t, int64(200), byName["node-2"].CPURequest)
	require.Equal(t, int64(400), byName["node-2"].UsedCPU)
}

package metricsresources

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
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
		require.Empty(t, config.Namespaces)
	})

	t.Run("with values", func(t *testing.T) {
		config := Config{
			KubeConfig:    "/path/to/config",
			KubeContext:   "test-context",
			Namespaces:    []string{"test-ns"},
			Label:         "app=test",
			FieldSelector: "spec.nodeName=node1",
			Nodes:         []string{"node1", "node2"},
			Sorting:       "name",
			Alert:         "memory",
			WatchPeriod:   10,
			Reverse:       true,
		}
		require.Equal(t, "/path/to/config", config.KubeConfig)
		require.Equal(t, "test-context", config.KubeContext)
		require.Equal(t, []string{"test-ns"}, config.Namespaces)
	})
}

func TestWatchResponse(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		resp := WatchResponse{Error: errors.New("test error")}
		require.Error(t, resp.Error)
		require.Empty(t, resp.Data)
	})

	t.Run("with data", func(t *testing.T) {
		data := PodMetricsResourceList{{PodResource: pods.PodResource{NamespaceName: pods.NamespaceName{Name: "test"}}}}
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

func (noopSuccessProcessor) Success(PodMetricsResourceList) {}

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

func TestMerge(t *testing.T) {
	t.Run("empty lists", func(t *testing.T) {
		result := merge(pods.PodResourceList{}, podmetrics.PodMetricList{})
		require.Empty(t, result)
	})

	t.Run("pod resource without metrics", func(t *testing.T) {
		podResources := pods.PodResourceList{
			{NamespaceName: pods.NamespaceName{Namespace: "ns1", Name: "pod1"}},
		}
		result := merge(podResources, podmetrics.PodMetricList{})
		require.Len(t, result, 1)
		require.Equal(t, "pod1", result[0].NamespaceName.Name)
		require.Equal(t, "ns1", result[0].NamespaceName.Namespace)
	})

	t.Run("matching pod and metrics", func(t *testing.T) {
		podResources := pods.PodResourceList{
			{NamespaceName: pods.NamespaceName{Namespace: "ns1", Name: "pod1"}},
		}
		podMetrics := podmetrics.PodMetricList{
			{Namespace: "ns1", Name: "pod1"},
		}
		result := merge(podResources, podMetrics)
		require.Len(t, result, 1)
		require.Equal(t, "pod1", result[0].NamespaceName.Name)
		require.Equal(t, "ns1", result[0].NamespaceName.Namespace)
	})

	t.Run("mismatched namespace/name", func(t *testing.T) {
		podResources := pods.PodResourceList{
			{NamespaceName: pods.NamespaceName{Namespace: "ns1", Name: "pod1"}},
		}
		podMetrics := podmetrics.PodMetricList{
			{Namespace: "ns2", Name: "pod2"},
		}
		result := merge(podResources, podMetrics)
		require.Len(t, result, 1)

		output := result.toOutput()
		require.Len(t, output.Items, 1)
		require.Equal(t, "pod1", output.Items[0].Name)
		require.Equal(t, "ns1", output.Items[0].Namespace)
	})

	t.Run("mismatched namespace/name does not emit warn log", func(t *testing.T) {
		oldLogger := slog.Default()
		defer slog.SetDefault(oldLogger)

		var logBuffer bytes.Buffer
		handler := slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{Level: slog.LevelDebug})
		slog.SetDefault(slog.New(handler))

		podResources := pods.PodResourceList{
			{NamespaceName: pods.NamespaceName{Namespace: "ns1", Name: "pod1"}},
		}
		podMetrics := podmetrics.PodMetricList{
			{Namespace: "ns2", Name: "pod2"},
		}

		result := merge(podResources, podMetrics)
		require.Len(t, result, 1)

		logs := logBuffer.String()
		require.NotContains(t, logs, `"level":"WARN"`)
		require.Contains(t, logs, "Skipped unmatched pod metrics")
	})
}

func TestFetchPodMetricsFollowsPagination(t *testing.T) {
	ctx := t.Context()
	coreClient := corefake.NewSimpleClientset()
	metricsClient := metricsfake.NewSimpleClientset()

	type listOptionsGetter interface {
		GetListOptions() metav1.ListOptions
	}

	coreClient.PrependReactor("list", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		require.Equal(t, "", listAction.GetListOptions().Continue)

		return true, &v1.PodList{
			ListMeta: metav1.ListMeta{Continue: "page-2"},
			Items: []v1.Pod{{
				ObjectMeta: metav1.ObjectMeta{Name: "pod-1", Namespace: "default"},
				Spec: v1.PodSpec{
					NodeName: "worker-1",
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
					NodeName: "worker-2",
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

	metricsClient.PrependReactor("list", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		require.Equal(t, "", listAction.GetListOptions().Continue)

		return true, &metricsv1beta1.PodMetricsList{
			ListMeta: metav1.ListMeta{Continue: "page-2"},
			Items: []metricsv1beta1.PodMetrics{{
				ObjectMeta: metav1.ObjectMeta{Name: "pod-1", Namespace: "default"},
				Containers: []metricsv1beta1.ContainerMetrics{{
					Name:  "app",
					Usage: v1.ResourceList{v1.ResourceCPU: resource.MustParse("150m")},
				}},
			}},
		}, nil
	})
	metricsClient.PrependReactor("list", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
		listAction, ok := action.(listOptionsGetter)
		require.True(t, ok)
		if listAction.GetListOptions().Continue != "page-2" {
			return false, nil, nil
		}

		return true, &metricsv1beta1.PodMetricsList{
			Items: []metricsv1beta1.PodMetrics{{
				ObjectMeta: metav1.ObjectMeta{Name: "pod-2", Namespace: "default"},
				Containers: []metricsv1beta1.ContainerMetrics{{
					Name:  "app",
					Usage: v1.ResourceList{v1.ResourceCPU: resource.MustParse("250m")},
				}},
			}},
		}, nil
	})

	result, err := FetchPodMetrics(ctx, NewPodRepository(), metricsClient.MetricsV1beta1(), coreClient.CoreV1(), FetchConfig{
		Namespaces: []string{"default"},
	})
	require.NoError(t, err)
	require.Len(t, result, 2)

	byName := make(map[string]PodMetricsResource, len(result))
	for _, item := range result {
		byName[item.Name] = item
	}

	require.Equal(t, int64(100), byName["pod-1"].PodResource.Containers[0].Requests.CPU)
	require.Equal(t, int64(150), byName["pod-1"].PodMetric.Containers[0].CPU)
	require.Equal(t, "worker-1", byName["pod-1"].NodeName)
	require.Equal(t, int64(200), byName["pod-2"].PodResource.Containers[0].Requests.CPU)
	require.Equal(t, int64(250), byName["pod-2"].PodMetric.Containers[0].CPU)
	require.Equal(t, "worker-2", byName["pod-2"].NodeName)
}

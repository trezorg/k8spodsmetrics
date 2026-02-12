package pods

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBuildFieldSelector(t *testing.T) {
	testCases := []struct {
		name     string
		filter   PodFilter
		expected string
	}{
		{
			name:     "no node name",
			filter:   PodFilter{FieldSelector: "metadata.name=foo"},
			expected: "metadata.name=foo",
		},
		{
			name:     "only node name",
			filter:   PodFilter{NodeName: "worker-1"},
			expected: "spec.nodeName=worker-1",
		},
		{
			name:     "field selector with node",
			filter:   PodFilter{FieldSelector: "status.phase=Running", NodeName: "worker-1"},
			expected: "status.phase=Running,spec.nodeName=worker-1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := buildFieldSelector(tc.filter)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestExtractContainerResources(t *testing.T) {
	t.Run("empty resources", func(t *testing.T) {
		container := v1.Container{
			Name: "test-container",
		}
		result := extractContainerResources(container)
		require.Equal(t, "test-container", result.Name)
		require.Equal(t, Resource{}, result.Limits)
		require.Equal(t, Resource{}, result.Requests)
	})

	t.Run("with cpu and memory limits", func(t *testing.T) {
		container := v1.Container{
			Name: "app-container",
			Resources: v1.ResourceRequirements{
				Limits: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("500m"),
					v1.ResourceMemory: resource.MustParse("1Gi"),
				},
				Requests: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("250m"),
					v1.ResourceMemory: resource.MustParse("512Mi"),
				},
			},
		}
		result := extractContainerResources(container)
		require.Equal(t, "app-container", result.Name)
		require.Equal(t, int64(500), result.Limits.CPU)
		require.Equal(t, int64(1024*1024*1024), result.Limits.Memory)
		require.Equal(t, int64(250), result.Requests.CPU)
		require.Equal(t, int64(512*1024*1024), result.Requests.Memory)
	})

	t.Run("with storage resources", func(t *testing.T) {
		container := v1.Container{
			Name: "storage-container",
			Resources: v1.ResourceRequirements{
				Limits: v1.ResourceList{
					v1.ResourceStorage:          resource.MustParse("10Gi"),
					v1.ResourceEphemeralStorage: resource.MustParse("5Gi"),
				},
			},
		}
		result := extractContainerResources(container)
		require.Equal(t, int64(10*1024*1024*1024), result.Limits.Storage)
		require.Equal(t, int64(5*1024*1024*1024), result.Limits.StorageEphemeral)
	})
}

func TestConvertPodToResource(t *testing.T) {
	t.Run("single container", func(t *testing.T) {
		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "test-ns",
			},
			Spec: v1.PodSpec{
				NodeName: "worker-1",
				Containers: []v1.Container{
					{
						Name: "main",
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceCPU: resource.MustParse("100m"),
							},
						},
					},
				},
			},
		}
		result := convertPodToResource(pod)
		require.Equal(t, "test-pod", result.Name)
		require.Equal(t, "test-ns", result.Namespace)
		require.Equal(t, "worker-1", result.NodeName)
		require.Len(t, result.Containers, 1)
		require.Equal(t, "main", result.Containers[0].Name)
		require.Equal(t, int64(100), result.Containers[0].Requests.CPU)
	})

	t.Run("multiple containers sorted", func(t *testing.T) {
		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "multi-pod",
				Namespace: "default",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Name: "z-container"},
					{Name: "a-container"},
					{Name: "m-container"},
				},
			},
		}
		result := convertPodToResource(pod)
		require.Len(t, result.Containers, 3)
		require.Equal(t, "a-container", result.Containers[0].Name)
		require.Equal(t, "m-container", result.Containers[1].Name)
		require.Equal(t, "z-container", result.Containers[2].Name)
	})
}

func TestResourceStructs(t *testing.T) {
	t.Run("Resource defaults", func(t *testing.T) {
		var r Resource
		require.Equal(t, int64(0), r.CPU)
		require.Equal(t, int64(0), r.Memory)
		require.Equal(t, int64(0), r.Storage)
		require.Equal(t, int64(0), r.StorageEphemeral)
	})

	t.Run("ContainerResource defaults", func(t *testing.T) {
		var cr ContainerResource
		require.Empty(t, cr.Name)
	})

	t.Run("PodFilter defaults", func(t *testing.T) {
		var pf PodFilter
		require.Empty(t, pf.Namespace)
		require.Empty(t, pf.LabelSelector)
		require.Empty(t, pf.FieldSelector)
		require.Empty(t, pf.NodeName)
	})

	t.Run("NamespaceName", func(t *testing.T) {
		nn := NamespaceName{Namespace: "prod", Name: "api"}
		require.Equal(t, "prod", nn.Namespace)
		require.Equal(t, "api", nn.Name)
	})

	t.Run("PodResourceList", func(t *testing.T) {
		list := PodResourceList{
			{NamespaceName: NamespaceName{Name: "pod1"}},
			{NamespaceName: NamespaceName{Name: "pod2"}},
		}
		require.Len(t, list, 2)
	})
}

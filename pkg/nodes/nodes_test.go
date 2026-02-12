package nodes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNodes(t *testing.T) {
	ctx := context.Background()

	t.Run("list all nodes", func(t *testing.T) {
		nodes := []v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node-1",
				},
				Status: v1.NodeStatus{
					Capacity: v1.ResourceList{
						v1.ResourceCPU:              resource.MustParse("4"),
						v1.ResourceMemory:           resource.MustParse("8Gi"),
						v1.ResourceStorage:          resource.MustParse("100Gi"),
						v1.ResourceEphemeralStorage: resource.MustParse("20Gi"),
					},
					Allocatable: v1.ResourceList{
						v1.ResourceCPU:              resource.MustParse("3.9"),
						v1.ResourceMemory:           resource.MustParse("7.8Gi"),
						v1.ResourceStorage:          resource.MustParse("90Gi"),
						v1.ResourceEphemeralStorage: resource.MustParse("18Gi"),
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node-2",
				},
				Status: v1.NodeStatus{
					Capacity: v1.ResourceList{
						v1.ResourceCPU:              resource.MustParse("2"),
						v1.ResourceMemory:           resource.MustParse("4Gi"),
						v1.ResourceStorage:          resource.MustParse("50Gi"),
						v1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
					},
					Allocatable: v1.ResourceList{
						v1.ResourceCPU:              resource.MustParse("1.9"),
						v1.ResourceMemory:           resource.MustParse("3.9Gi"),
						v1.ResourceStorage:          resource.MustParse("45Gi"),
						v1.ResourceEphemeralStorage: resource.MustParse("9Gi"),
					},
				},
			},
		}

		client := fake.NewSimpleClientset(&v1.NodeList{Items: nodes})

		result, err := Nodes(ctx, client.CoreV1(), NodeFilter{}, "")
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, "node-1", result[0].Name)
		require.Equal(t, "node-2", result[1].Name)

		require.Equal(t, int64(4000), result[0].CPU)
		require.Equal(t, int64(3900), result[0].AllocatableCPU)
		require.Equal(t, int64(100), result[0].UsedCPU)

		require.Equal(t, int64(2000), result[1].CPU)
		require.Equal(t, int64(1900), result[1].AllocatableCPU)
		require.Equal(t, int64(100), result[1].UsedCPU)
	})

	t.Run("get specific node", func(t *testing.T) {
		node := &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-1",
			},
			Status: v1.NodeStatus{
				Capacity: v1.ResourceList{
					v1.ResourceCPU:              resource.MustParse("4"),
					v1.ResourceMemory:           resource.MustParse("8Gi"),
					v1.ResourceStorage:          resource.MustParse("100Gi"),
					v1.ResourceEphemeralStorage: resource.MustParse("20Gi"),
				},
				Allocatable: v1.ResourceList{
					v1.ResourceCPU:              resource.MustParse("3.9"),
					v1.ResourceMemory:           resource.MustParse("7.8Gi"),
					v1.ResourceStorage:          resource.MustParse("90Gi"),
					v1.ResourceEphemeralStorage: resource.MustParse("18Gi"),
				},
			},
		}

		client := fake.NewSimpleClientset(node)

		result, err := Nodes(ctx, client.CoreV1(), NodeFilter{}, "node-1")
		require.NoError(t, err)
		require.Len(t, result, 1)
		require.Equal(t, "node-1", result[0].Name)
	})

	t.Run("list with label selector", func(t *testing.T) {
		nodes := []v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node-1",
					Labels: map[string]string{
						"env": "prod",
					},
				},
				Status: v1.NodeStatus{
					Capacity: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("4"),
						v1.ResourceMemory: resource.MustParse("8Gi"),
					},
					Allocatable: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("3.9"),
						v1.ResourceMemory: resource.MustParse("7.8Gi"),
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node-2",
					Labels: map[string]string{
						"env": "dev",
					},
				},
				Status: v1.NodeStatus{
					Capacity: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("2"),
						v1.ResourceMemory: resource.MustParse("4Gi"),
					},
					Allocatable: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("1.9"),
						v1.ResourceMemory: resource.MustParse("3.9Gi"),
					},
				},
			},
		}

		client := fake.NewSimpleClientset(&v1.NodeList{Items: nodes})

		result, err := Nodes(ctx, client.CoreV1(), NodeFilter{LabelSelector: "env=prod"}, "")
		require.NoError(t, err)
		require.Len(t, result, 1)
		require.Equal(t, "node-1", result[0].Name)
	})

	t.Run("empty list", func(t *testing.T) {
		client := fake.NewSimpleClientset()

		result, err := Nodes(ctx, client.CoreV1(), NodeFilter{}, "")
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("node not found", func(t *testing.T) {
		client := fake.NewSimpleClientset()

		_, err := Nodes(ctx, client.CoreV1(), NodeFilter{}, "nonexistent-node")
		require.Error(t, err)
	})
}

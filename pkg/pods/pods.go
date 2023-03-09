package pods

import (
	"context"
	"sort"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Resource struct {
	CPU    int64
	Memory int64
}
type ContainerResource struct {
	Name     string
	Limits   Resource
	Requests Resource
}

type NamespaceName struct {
	Namespace string
	Name      string
}

type PodResource struct {
	NamespaceName
	Containers []ContainerResource
	NodeName   string
}

type PodResourceList []PodResource
type PodFilter struct {
	Namespace     string
	LabelSelector string
	FieldSelector string
}

// Pods get pods for MetricFilter
func Pods(ctx context.Context, corev1 corev1.CoreV1Interface, filter PodFilter) (PodResourceList, error) {
	var result PodResourceList
	pods, err := corev1.Pods(filter.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: filter.LabelSelector,
		FieldSelector: filter.FieldSelector,
	})
	if err != nil {
		return result, err
	}
	for _, pod := range pods.Items {
		podResource := PodResource{
			NamespaceName: NamespaceName{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
			NodeName: pod.Spec.NodeName,
		}
		for _, container := range pod.Spec.Containers {
			limits := container.Resources.Limits
			requests := container.Resources.Requests
			containerResource := ContainerResource{
				Name: container.Name,
			}
			if cpuLimit, ok := limits[v1.ResourceCPU]; ok {
				containerResource.Limits.CPU = cpuLimit.MilliValue()
			}
			if memoryLimit, ok := limits[v1.ResourceMemory]; ok {
				if memory, ok := memoryLimit.AsInt64(); ok {
					containerResource.Limits.Memory = memory
				}
			}
			if cpuRequest, ok := requests[v1.ResourceCPU]; ok {
				containerResource.Requests.CPU = cpuRequest.MilliValue()
			}
			if memoryRequest, ok := requests[v1.ResourceMemory]; ok {
				if memory, ok := memoryRequest.AsInt64(); ok {
					containerResource.Requests.Memory = memory
				}
			}
			podResource.Containers = append(podResource.Containers, containerResource)

		}
		sort.Slice(podResource.Containers, func(i, j int) bool {
			return podResource.Containers[i].Name < podResource.Containers[j].Name
		})
		result = append(result, podResource)
	}
	return result, nil
}

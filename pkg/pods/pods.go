package pods

import (
	"context"
	"sort"

	v1 "k8s.io/api/core/v1" //nolint:revive // it is used
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Resource struct {
	CPU              int64 `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory           int64 `json:"memory,omitempty" yaml:"memory,omitempty"`
	Storage          int64 `json:"storage,omitempty" yaml:"storage,omitempty"`
	StorageEphemeral int64 `json:"storage_ephemeral,omitempty" yaml:"storage_ephemeral,omitempty"`
}

type ContainerResource struct {
	Name     string   `json:"name,omitempty" yaml:"name,omitempty"`
	Limits   Resource `json:"limits,omitempty" yaml:"limits,omitempty"`
	Requests Resource `json:"requests,omitempty" yaml:"requests,omitempty"`
}

type NamespaceName struct {
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
}

type PodResource struct {
	NamespaceName `json:"namespace_name,omitempty" yaml:"namespace_name,omitempty"`
	NodeName      string              `json:"node_name,omitempty" yaml:"node_name,omitempty"`
	Containers    []ContainerResource `json:"containers,omitempty" yaml:"containers,omitempty"`
}

type PodResourceList []PodResource

type PodFilter struct {
	Namespace     string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	LabelSelector string `json:"label_selector,omitempty" yaml:"label_selector,omitempty"`
	FieldSelector string `json:"field_selector,omitempty" yaml:"field_selector,omitempty"`
}

func extractContainerResources(container v1.Container) ContainerResource {
	limits := container.Resources.Limits
	requests := container.Resources.Requests
	containerResource := ContainerResource{
		Name: container.Name,
	}

	// Extract limits
	if cpuLimit, ok := limits[v1.ResourceCPU]; ok {
		containerResource.Limits.CPU = cpuLimit.MilliValue()
	}
	if memoryLimit, ok := limits[v1.ResourceMemory]; ok {
		if memory, ok := memoryLimit.AsInt64(); ok {
			containerResource.Limits.Memory = memory
		}
	}
	if storageLimit, ok := limits[v1.ResourceStorage]; ok {
		if storage, ok := storageLimit.AsInt64(); ok {
			containerResource.Limits.Storage = storage
		}
	}
	if storageEphemeralLimit, ok := limits[v1.ResourceEphemeralStorage]; ok {
		if storage, ok := storageEphemeralLimit.AsInt64(); ok {
			containerResource.Limits.StorageEphemeral = storage
		}
	}

	// Extract requests
	if cpuRequest, ok := requests[v1.ResourceCPU]; ok {
		containerResource.Requests.CPU = cpuRequest.MilliValue()
	}
	if memoryRequest, ok := requests[v1.ResourceMemory]; ok {
		if memory, ok := memoryRequest.AsInt64(); ok {
			containerResource.Requests.Memory = memory
		}
	}
	if storageRequest, ok := requests[v1.ResourceStorage]; ok {
		if storage, ok := storageRequest.AsInt64(); ok {
			containerResource.Requests.Storage = storage
		}
	}
	if storageEphemeralRequest, ok := requests[v1.ResourceEphemeralStorage]; ok {
		if storage, ok := storageEphemeralRequest.AsInt64(); ok {
			containerResource.Requests.StorageEphemeral = storage
		}
	}

	return containerResource
}

func convertPodToResource(pod v1.Pod) PodResource {
	podResource := PodResource{
		NamespaceName: NamespaceName{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		NodeName: pod.Spec.NodeName,
	}

	containers := make([]ContainerResource, 0, len(pod.Spec.Containers))
	for _, container := range pod.Spec.Containers {
		containers = append(containers, extractContainerResources(container))
	}

	sort.Slice(containers, func(i, j int) bool {
		return containers[i].Name < containers[j].Name
	})

	podResource.Containers = containers
	return podResource
}

func Pods(
	ctx context.Context,
	coreV1Ifc corev1.CoreV1Interface,
	filter PodFilter,
	nodeName string,
) (PodResourceList, error) {
	pods, err := coreV1Ifc.Pods(filter.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: filter.LabelSelector,
		FieldSelector: filter.FieldSelector,
	})
	if err != nil {
		return nil, err
	}

	result := make(PodResourceList, 0, len(pods.Items))
	for _, pod := range pods.Items {
		if nodeName != "" && nodeName != pod.Spec.NodeName {
			continue
		}
		result = append(result, convertPodToResource(pod))
	}
	return result, nil
}

package pods

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"sync"

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
	Limits   Resource `json:"limits" yaml:"limits"`
	Requests Resource `json:"requests" yaml:"requests"`
}

type NamespaceName struct {
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
}

type PodResource struct {
	NamespaceName `json:"namespace_name" yaml:"namespace_name"`
	NodeName      string              `json:"node_name,omitempty" yaml:"node_name,omitempty"`
	Containers    []ContainerResource `json:"containers,omitempty" yaml:"containers,omitempty"`
}

type PodResourceList []PodResource

type PodFilter struct {
	Namespaces    []string `json:"namespaces,omitempty" yaml:"namespaces,omitempty"`
	LabelSelector string   `json:"label_selector,omitempty" yaml:"label_selector,omitempty"`
	FieldSelector string   `json:"field_selector,omitempty" yaml:"field_selector,omitempty"`
	NodeName      string   `json:"node,omitempty" yaml:"node,omitempty"`
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
	nodeNames ...string,
) (PodResourceList, error) {
	nodeNames = slices.DeleteFunc(nodeNames, func(n string) bool { return n == "" })
	filter.Namespaces = slices.DeleteFunc(filter.Namespaces, func(n string) bool { return n == "" })

	// If no namespaces specified, query all namespaces (empty string)
	if len(filter.Namespaces) == 0 {
		return podsForNamespace(ctx, coreV1Ifc, filter, nodeNames, "")
	}

	// Single namespace
	if len(filter.Namespaces) == 1 {
		return podsForNamespace(ctx, coreV1Ifc, filter, nodeNames, filter.Namespaces[0])
	}

	// Multiple namespaces: query each in parallel
	var wg sync.WaitGroup

	var errs error
	rErrors := make([]error, len(filter.Namespaces))
	pods := make([]PodResourceList, len(filter.Namespaces))

	for idx, ns := range filter.Namespaces {
		wg.Go(func() {
			pods[idx], rErrors[idx] = podsForNamespace(ctx, coreV1Ifc, filter, nodeNames, ns)
		})
	}

	wg.Wait()

	for _, err := range rErrors {
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}

	resultLen := 0
	for _, p := range pods {
		resultLen += len(p)
	}
	result := make(PodResourceList, 0, resultLen)
	for _, p := range pods {
		result = append(result, p...)
	}
	return result, nil
}

func podsForNamespace(
	ctx context.Context,
	coreV1Ifc corev1.CoreV1Interface,
	filter PodFilter,
	nodeNames []string,
	namespace string,
) (PodResourceList, error) {
	if len(nodeNames) == 0 {
		return listPods(ctx, coreV1Ifc, filter, namespace)
	}

	var wg sync.WaitGroup

	var errs error
	rErrors := make([]error, len(nodeNames))
	pods := make([]PodResourceList, len(nodeNames))

	for idx, nodeName := range nodeNames {
		wg.Go(func() {
			nodeFilter := filter
			nodeFilter.NodeName = nodeName
			pods[idx], rErrors[idx] = listPods(ctx, coreV1Ifc, nodeFilter, namespace)
		})
	}

	wg.Wait()

	for _, err := range rErrors {
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return nil, errs
	}

	resultLen := 0
	for _, p := range pods {
		resultLen += len(p)
	}
	result := make(PodResourceList, 0, resultLen)
	for _, p := range pods {
		result = append(result, p...)
	}
	return result, nil
}

func listPods(
	ctx context.Context,
	coreV1Ifc corev1.CoreV1Interface,
	filter PodFilter,
	namespace string,
) (PodResourceList, error) {
	opts := metav1.ListOptions{
		LabelSelector: filter.LabelSelector,
		FieldSelector: buildFieldSelector(filter),
	}
	pods, err := coreV1Ifc.Pods(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	result := make(PodResourceList, 0, len(pods.Items))
	for _, pod := range pods.Items {
		result = append(result, convertPodToResource(pod))
	}
	return result, nil
}

func buildFieldSelector(filter PodFilter) string {
	fieldSelector := filter.FieldSelector
	if filter.NodeName == "" {
		return fieldSelector
	}

	nodeSelector := fmt.Sprintf("spec.nodeName=%s", filter.NodeName)
	if fieldSelector == "" {
		return nodeSelector
	}

	return fieldSelector + "," + nodeSelector
}

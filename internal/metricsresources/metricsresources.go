package metricsresources

import (
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

const (
	Requests ResourceType = iota
	Limits
	unset = int64(-1)
)

type (
	ResourceType uint

	PodMetricsResource struct {
		pods.PodResource
		podmetrics.PodMetric
	}

	PodMetricsResourceList []PodMetricsResource

	MetricsResource struct {
		CPURequest              int64 `json:"cpu_request,omitempty" yaml:"cpu_request,omitempty"`
		MemoryRequest           int64 `json:"memory_request,omitempty" yaml:"memory_request,omitempty"`
		CPUUsed                 int64 `json:"cpu_used,omitempty" yaml:"cpu_used,omitempty"`
		MemoryUsed              int64 `json:"memory_used,omitempty" yaml:"memory_used,omitempty"`
		StorageRequest          int64 `json:"storage_request,omitempty" yaml:"storage_request,omitempty"`
		StorageEphemeralRequest int64 `json:"storage_ephemeral_request,omitempty" yaml:"storage_ephemeral_request,omitempty"`
		StorageUsed             int64 `json:"storage_used,omitempty" yaml:"storage_used,omitempty"`
		StorageEphemeralUsed    int64 `json:"storage_ephemeral_used,omitempty" yaml:"storage_ephemeral_used,omitempty"`
	}

	Resource struct {
		CPU    int64 `json:"cpu,omitempty" yaml:"cpu,omitempty"`
		Memory int64 `json:"memory,omitempty" yaml:"memory,omitempty"`
	}

	ContainerMetricsResource struct {
		Name     string          `json:"name,omitempty" yaml:"name,omitempty"`
		Limits   MetricsResource `json:"limits" yaml:"limits"`
		Requests MetricsResource `json:"requests" yaml:"requests"`
	}

	ContainerMetricsResourceOutput struct {
		Name     string   `json:"name,omitempty" yaml:"name"`
		Limits   Resource `json:"limits" yaml:"limits"`
		Requests Resource `json:"requests" yaml:"requests"`
		Used     Resource `json:"used" yaml:"used"`
	}

	ContainerMetricsResources        []ContainerMetricsResource
	ContainerMetricsResourcesOutputs []ContainerMetricsResourceOutput

	PodMetricsResourceOutput struct {
		Name       string                           `json:"name,omitempty" yaml:"name,omitempty"`
		Namespace  string                           `json:"namespace,omitempty" yaml:"namespace,omitempty"`
		Node       string                           `json:"node,omitempty" yaml:"node,omitempty"`
		Containers ContainerMetricsResourcesOutputs `json:"containers,omitempty" yaml:"containers,omitempty"`
	}
	PodMetricsResourceListOutput []PodMetricsResourceOutput

	PodMetricsResourceOutputEnvelope struct {
		Items PodMetricsResourceListOutput `json:"items,omitempty" yaml:"items,omitempty"`
	}

	containerMetricsPredicate   func(c ContainerMetricsResources) bool
	podResourceMetricsPredicate func(c PodMetricsResource) bool
)

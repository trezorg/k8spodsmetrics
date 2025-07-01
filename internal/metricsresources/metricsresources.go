package metricsresources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"text/template"

	escapes "github.com/snugfox/ansi-escapes"
	alerts "github.com/trezorg/k8spodsmetrics/internal/alert"
	"github.com/trezorg/k8spodsmetrics/internal/humanize"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

const (
	Requests ResourceType = iota
	Limits
	unset = int64(-1)
)

var (
	metricsPodTemplate = template.Must(template.New("metricPod").Parse(`Name:		{{.Name}}
Namespace:	{{.Namespace}}
Node:		{{.NodeName}}
Containers:
  {{ range $index, $container := .ContainersMetrics -}}
  Name:         {{ $container.Name }}
  Requests: 	{{ (index $container.Requests).StringWithColor "yellow" }}
  Limits:    	{{ (index $container.Limits).StringWithColor "red" }}
  {{ end -}}`))
	metricsPodsTemplate = template.Must(template.New("metricPods").Parse(`{{ range $index, $pod := . -}}
{{ $pod }}
{{ end -}}`))
)

type (
	ResourceType uint

	PodMetricsResource struct {
		pods.PodResource
		podmetrics.PodMetric
	}

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
		Limits   MetricsResource `json:"limits,omitempty" yaml:"limits,omitempty"`
		Requests MetricsResource `json:"requests,omitempty" yaml:"requests,omitempty"`
	}

	ContainerMetricsResourceOutput struct {
		Name     string   `json:"name,omitempty" yaml:"name,omitempty"`
		Limits   Resource `json:"limits,omitempty" yaml:"limits,omitempty"`
		Requests Resource `json:"requests,omitempty" yaml:"requests,omitempty"`
		Used     Resource `json:"used,omitempty" yaml:"used,omitempty"`
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

func (c ContainerMetricsResource) IsMemoryAlerted() bool {
	return c.Limits.MemoryAlert() || c.Requests.MemoryAlert()
}

func (c ContainerMetricsResource) IsMemoryRequestAlerted() bool {
	return c.Requests.MemoryAlert()
}

func (c ContainerMetricsResource) IsMemoryLimitAlerted() bool {
	return c.Limits.MemoryAlert()
}

func (c ContainerMetricsResource) IsCPUAlerted() bool {
	return c.Limits.CPUAlert() || c.Requests.CPUAlert()
}

func (c ContainerMetricsResource) IsCPURequestAlerted() bool {
	return c.Requests.CPUAlert()
}

func (c ContainerMetricsResource) IsCPULimitAlerted() bool {
	return c.Limits.CPUAlert()
}

func (c ContainerMetricsResource) IsAlerted() bool {
	return c.IsMemoryAlerted() || c.IsCPUAlerted()
}

func (c ContainerMetricsResource) MemoryUsed() string {
	if c.Limits.MemoryAlert() {
		return c.Limits.MemoryUsedString(escapes.TextColorRed)
	}
	return c.Requests.MemoryUsedString(escapes.TextColorYellow)
}

func (c ContainerMetricsResource) CPUUsed() string {
	if c.Limits.CPUAlert() {
		return c.Limits.CPUUsedString(escapes.TextColorRed)
	}
	return c.Requests.CPUUsedString(escapes.TextColorRed)
}

func (c ContainerMetricsResource) StorageUsed() string {
	return c.Requests.StorageString()
}

func (c ContainerMetricsResource) StorageEphemeralUsed() string {
	return c.Requests.StorageEphemeralString()
}

func (c ContainerMetricsResource) toOutput() ContainerMetricsResourceOutput {
	return ContainerMetricsResourceOutput{
		Name: c.Name,
		Limits: Resource{
			CPU:    c.Limits.CPURequest,
			Memory: c.Limits.MemoryRequest,
		},
		Requests: Resource{
			CPU:    c.Requests.CPURequest,
			Memory: c.Requests.MemoryRequest,
		},
		Used: Resource{
			CPU:    c.Requests.CPUUsed,
			Memory: c.Requests.MemoryUsed,
		},
	}
}

func (c ContainerMetricsResources) IsMemoryAlerted() bool {
	for _, container := range c {
		if container.IsMemoryAlerted() {
			return true
		}
	}
	return false
}

func (c ContainerMetricsResources) IsMemoryRequestAlerted() bool {
	for _, container := range c {
		if container.IsMemoryRequestAlerted() {
			return true
		}
	}
	return false
}

func (c ContainerMetricsResources) IsMemoryLimitAlerted() bool {
	for _, container := range c {
		if container.IsMemoryLimitAlerted() {
			return true
		}
	}
	return false
}

func (c ContainerMetricsResources) IsCPUAlerted() bool {
	for _, container := range c {
		if container.IsCPUAlerted() {
			return true
		}
	}
	return false
}

func (c ContainerMetricsResources) IsCPURequestAlerted() bool {
	for _, container := range c {
		if container.IsCPURequestAlerted() {
			return true
		}
	}
	return false
}

func (c ContainerMetricsResources) IsCPULimitAlerted() bool {
	for _, container := range c {
		if container.IsCPULimitAlerted() {
			return true
		}
	}
	return false
}

func (c ContainerMetricsResources) IsAlerted() bool {
	for _, container := range c {
		if container.IsAlerted() {
			return true
		}
	}
	return false
}

func (c ContainerMetricsResources) toOutput() ContainerMetricsResourcesOutputs {
	result := make(ContainerMetricsResourcesOutputs, 0, len(c))
	for _, container := range c {
		result = append(result, container.toOutput())
	}
	return result
}

func (m MetricsResource) CPUAlert() bool {
	return m.CPURequest > 0 && m.CPURequest <= m.CPUUsed
}

func (m MetricsResource) MemoryAlert() bool {
	return m.MemoryRequest > 0 && m.MemoryRequest <= m.MemoryUsed
}

func (m MetricsResource) CPU(alertColor string) string {
	if m.CPUUsed == unset {
		return fmt.Sprintf("%d", m.CPURequest)
	}
	cpuStartColor := ""
	cpuEndColor := ""
	if m.CPUAlert() {
		cpuStartColor = alertColor
		cpuEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%d/%s%d%s",
		m.CPURequest,
		cpuStartColor,
		m.CPUUsed,
		cpuEndColor,
	)
}

func (m MetricsResource) CPURequestString() string {
	return fmt.Sprintf("%d", m.CPURequest)
}

func (m MetricsResource) CPUUsedString(alertColor string) string {
	if m.CPUUsed == unset {
		return ""
	}
	cpuStartColor := ""
	cpuEndColor := ""
	if m.CPUAlert() {
		cpuStartColor = alertColor
		cpuEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuStartColor,
		m.CPUUsed,
		cpuEndColor,
	)
}

func (m MetricsResource) Memory(alertColor string) string {
	if m.MemoryUsed == unset {
		return humanize.Bytes(m.MemoryRequest)
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if m.MemoryAlert() {
		memoryStartColor = alertColor
		memoryEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s/%s%s%s",
		humanize.Bytes(m.MemoryRequest),
		memoryStartColor,
		humanize.Bytes(m.MemoryUsed),
		memoryEndColor,
	)
}

func (m MetricsResource) MemoryRequestString() string {
	return humanize.Bytes(m.MemoryRequest)
}

func (m MetricsResource) MemoryUsedString(alertColor string) string {
	if m.MemoryUsed == unset {
		return ""
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if m.MemoryAlert() {
		memoryStartColor = alertColor
		memoryEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryStartColor,
		humanize.Bytes(m.MemoryUsed),
		memoryEndColor,
	)
}

func (m MetricsResource) StringWithColor(alertColor string) string {
	cpuStartColor := ""
	cpuEndColor := ""
	if m.CPUAlert() {
		cpuStartColor = alertColor
		cpuEndColor = escapes.ColorReset
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if m.MemoryAlert() {
		memoryStartColor = alertColor
		memoryEndColor = escapes.ColorReset
	}
	if m.MemoryUsed == unset && m.CPUUsed == unset {
		return fmt.Sprintf(
			"CPU=%d, Memory=%s",
			m.CPURequest,
			humanize.Bytes(m.MemoryRequest),
		)
	}
	return fmt.Sprintf(
		"CPU=%d/%s%d%s, Memory=%s/%s%s%s",
		m.CPURequest,
		cpuStartColor,
		m.CPUUsed,
		cpuEndColor,
		humanize.Bytes(m.MemoryRequest),
		memoryStartColor,
		humanize.Bytes(m.MemoryUsed),
		memoryEndColor,
	)
}

func (m MetricsResource) String() string {
	return m.StringWithColor(escapes.TextColorRed)
}

func (r PodMetricsResource) ContainersMetrics() ContainerMetricsResources {
	containerMetricsResources := make(ContainerMetricsResources, 0, len(r.PodMetric.Containers))
	for i, container := range r.PodResource.Containers {
		containerMetricsResource := ContainerMetricsResource{
			Name: container.Name,
		}
		var cpuMetric, memoryMetric, storageMetric, storageEphemeralMetric int64
		if i < len(r.PodMetric.Containers) {
			metrics := r.PodMetric.Containers[i]
			cpuMetric = metrics.CPU
			memoryMetric = metrics.Memory
			storageMetric = metrics.Storage
			storageEphemeralMetric = metrics.StorageEphemeral
		} else {
			cpuMetric = unset
			memoryMetric = unset
			storageMetric = unset
			storageEphemeralMetric = unset
		}
		containerMetricsResource.Requests = MetricsResource{
			CPURequest:              container.Requests.CPU,
			MemoryRequest:           container.Requests.Memory,
			CPUUsed:                 cpuMetric,
			MemoryUsed:              memoryMetric,
			StorageRequest:          container.Requests.Storage,
			StorageEphemeralRequest: container.Requests.StorageEphemeral,
			StorageUsed:             storageMetric,
			StorageEphemeralUsed:    storageEphemeralMetric,
		}
		containerMetricsResource.Limits = MetricsResource{
			CPURequest:              container.Limits.CPU,
			MemoryRequest:           container.Limits.Memory,
			CPUUsed:                 cpuMetric,
			MemoryUsed:              memoryMetric,
			StorageRequest:          container.Limits.Storage,
			StorageEphemeralRequest: container.Limits.StorageEphemeral,
			StorageUsed:             storageMetric,
			StorageEphemeralUsed:    storageEphemeralMetric,
		}
		containerMetricsResources = append(containerMetricsResources, containerMetricsResource)
	}
	return containerMetricsResources
}

func (r PodMetricsResource) String() string {
	var buffer bytes.Buffer
	if err := metricsPodTemplate.Execute(&buffer, r); err != nil {
		panic(err)
	}
	return buffer.String()
}

type PodMetricsResourceList []PodMetricsResource

func (r PodMetricsResourceList) String() string {
	var buffer bytes.Buffer
	if err := metricsPodsTemplate.Execute(&buffer, r); err != nil {
		panic(err)
	}
	return buffer.String()
}

func (r PodMetricsResourceList) filterBy(predicate containerMetricsPredicate) PodMetricsResourceList {
	var result PodMetricsResourceList
	for _, pod := range r {
		if predicate(pod.ContainersMetrics()) {
			result = append(result, pod)
		}
	}
	return result
}

func (r PodMetricsResourceList) filterByPodResource(predicate podResourceMetricsPredicate) PodMetricsResourceList {
	var result PodMetricsResourceList
	for _, pod := range r {
		if predicate(pod) {
			result = append(result, pod)
		}
	}
	return result
}

func (r PodMetricsResourceList) filterNodes(nodes []string) PodMetricsResourceList {
	if len(nodes) == 0 {
		return r
	}
	return r.filterByPodResource(func(r PodMetricsResource) bool {
		return slices.Contains(nodes, r.NodeName)
	})
}

func (m MetricsResource) StorageString() string {
	return humanize.Bytes(m.StorageUsed)
}

func (m MetricsResource) StorageEphemeralString() string {
	return humanize.Bytes(m.StorageEphemeralUsed)
}

func (m MetricsResource) StorageRequestString() string {
	return humanize.Bytes(m.StorageRequest)
}

func (m MetricsResource) StorageEphemeralRequestString() string {
	return humanize.Bytes(m.StorageEphemeralRequest)
}

func (r PodMetricsResourceList) filterByAlert(alert alerts.Alert) PodMetricsResourceList {
	switch alert {
	case alerts.Any:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsAlerted() })
	case alerts.Memory:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsMemoryAlerted() })
	case alerts.MemoryRequest:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsMemoryRequestAlerted() })
	case alerts.MemoryLimit:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsMemoryLimitAlerted() })
	case alerts.CPU:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsCPUAlerted() })
	case alerts.CPURequest:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsCPURequestAlerted() })
	case alerts.CPULimit:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsCPULimitAlerted() })
	case alerts.None, alerts.Storage, alerts.StorageEphemeral:
		return r
	default:
		return r
	}
}

func (r PodMetricsResource) toOutput() PodMetricsResourceOutput {
	containers := r.ContainersMetrics()
	return PodMetricsResourceOutput{
		r.Name,
		r.Namespace,
		r.NodeName,
		containers.toOutput(),
	}
}

func (r PodMetricsResourceList) toOutput() PodMetricsResourceOutputEnvelope {
	items := make([]PodMetricsResourceOutput, 0, len(r))
	for _, item := range r {
		items = append(items, item.toOutput())
	}
	return PodMetricsResourceOutputEnvelope{
		Items: items,
	}
}

func (r PodMetricsResource) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(r.toOutput(), "", "    ")
}

func (r PodMetricsResourceList) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(r.toOutput(), "", "    ")
}

func (r PodMetricsResource) MarshalYAML() (any, error) {
	node := yaml.Node{}
	err := node.Encode(r.toOutput())
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (r PodMetricsResourceList) MarshalYAML() (any, error) {
	node := yaml.Node{}
	err := node.Encode(r.toOutput())
	if err != nil {
		return nil, err
	}
	return node, nil
}

func merge(podResourceList pods.PodResourceList, podMetricList podmetrics.PodMetricList) PodMetricsResourceList {
	podsMap := make(map[pods.NamespaceName]*PodMetricsResource)
	for _, pr := range podResourceList {
		podsMap[pods.NamespaceName{Namespace: pr.Namespace, Name: pr.Name}] = &PodMetricsResource{PodResource: pr}
	}
	for _, pm := range podMetricList {
		podMetricsResource, ok := podsMap[pods.NamespaceName{Namespace: pm.Namespace, Name: pm.Name}]
		if !ok {
			logger.Warn("Cannot substitute namespace and name", slog.String("namespace", pm.Namespace), slog.String("name", pm.Name))
			continue
		}
		podMetricsResource.PodMetric = pm
	}
	podMetricsResourceList := make(PodMetricsResourceList, 0, len(podsMap))
	for _, podMetricsResource := range podsMap {
		podMetricsResourceList = append(podMetricsResourceList, *podMetricsResource)
	}
	return podMetricsResourceList
}

package metricsresources

import (
	"bytes"
	"encoding/json"
	"fmt"

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
  Requests: 	{{ index $container.Requests }}
  Limits:    	{{ index $container.Limits }}
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
		CPURequest    int64  `json:"cpu_request,omitempty" yaml:"cpu_request,omitempty"`
		MemoryRequest int64  `json:"memory_request,omitempty" yaml:"memory_request,omitempty"`
		CPUUsed       int64  `json:"cpu_used,omitempty" yaml:"cpu_used,omitempty"`
		MemoryUsed    int64  `json:"memory_used,omitempty" yaml:"memory_used,omitempty"`
		AlertColor    string `json:"-" yaml:"-"`
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
	containerMetricsPredicate func(c ContainerMetricsResources) bool
)

func (c ContainerMetricsResource) IsMemoryAlerted() bool {
	return c.Limits.MemoryAlert() || c.Requests.MemoryAlert()
}

func (c ContainerMetricsResource) IsCPUAlerted() bool {
	return c.Limits.CPUAlert() || c.Requests.CPUAlert()
}

func (c ContainerMetricsResource) IsAlerted() bool {
	return c.IsMemoryAlerted() || c.IsCPUAlerted()
}

func (c ContainerMetricsResource) MemoryUsed() string {
	if c.Limits.MemoryAlert() {
		return c.Limits.MemoryUsedString()
	}
	return c.Requests.MemoryUsedString()
}

func (c ContainerMetricsResource) CPUUsed() string {
	if c.Limits.CPUAlert() {
		return c.Limits.CPUUsedString()
	}
	return c.Requests.CPUUsedString()
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

func (c ContainerMetricsResources) IsCPUAlerted() bool {
	for _, container := range c {
		if container.IsCPUAlerted() {
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

func (metricResource MetricsResource) CPU() string {
	if metricResource.CPUUsed == unset {
		return fmt.Sprintf("%d", metricResource.CPURequest)
	}
	cpuStartColor := ""
	cpuEndColor := ""
	if metricResource.CPUAlert() {
		cpuStartColor = metricResource.AlertColor
		cpuEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%d/%s%d%s",
		metricResource.CPURequest,
		cpuStartColor,
		metricResource.CPUUsed,
		cpuEndColor,
	)
}

func (metricResource MetricsResource) CPURequestString() string {
	return fmt.Sprintf("%d", metricResource.CPURequest)
}

func (metricResource MetricsResource) CPUUsedString() string {
	if metricResource.CPUUsed == unset {
		return ""
	}
	cpuStartColor := ""
	cpuEndColor := ""
	if metricResource.CPUAlert() {
		cpuStartColor = metricResource.AlertColor
		cpuEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuStartColor,
		metricResource.CPUUsed,
		cpuEndColor,
	)
}

func (metricResource MetricsResource) Memory() string {
	if metricResource.MemoryUsed == unset {
		return humanize.Bytes(metricResource.MemoryRequest)
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if metricResource.MemoryAlert() {
		memoryStartColor = metricResource.AlertColor
		memoryEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s/%s%s%s",
		humanize.Bytes(metricResource.MemoryRequest),
		memoryStartColor,
		humanize.Bytes(metricResource.MemoryUsed),
		memoryEndColor,
	)
}

func (metricResource MetricsResource) MemoryRequestString() string {
	return humanize.Bytes(metricResource.MemoryRequest)
}

func (metricResource MetricsResource) MemoryUsedString() string {
	if metricResource.MemoryUsed == unset {
		return ""
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if metricResource.MemoryAlert() {
		memoryStartColor = metricResource.AlertColor
		memoryEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryStartColor,
		humanize.Bytes(metricResource.MemoryUsed),
		memoryEndColor,
	)
}

func (metricResource MetricsResource) String() string {
	cpuStartColor := ""
	cpuEndColor := ""
	if metricResource.CPUAlert() {
		cpuStartColor = metricResource.AlertColor
		cpuEndColor = escapes.ColorReset
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if metricResource.MemoryAlert() {
		memoryStartColor = metricResource.AlertColor
		memoryEndColor = escapes.ColorReset
	}
	if metricResource.MemoryUsed == unset && metricResource.CPUUsed == unset {
		return fmt.Sprintf(
			"CPU=%d, Memory=%s",
			metricResource.CPURequest,
			humanize.Bytes(metricResource.MemoryRequest),
		)
	}
	return fmt.Sprintf(
		"CPU=%d/%s%d%s, Memory=%s/%s%s%s",
		metricResource.CPURequest,
		cpuStartColor,
		metricResource.CPUUsed,
		cpuEndColor,
		humanize.Bytes(metricResource.MemoryRequest),
		memoryStartColor,
		humanize.Bytes(metricResource.MemoryUsed),
		memoryEndColor,
	)
}

func (r PodMetricsResource) ContainersMetrics() ContainerMetricsResources {
	var containerMetricsResources ContainerMetricsResources
	for i, container := range r.PodResource.Containers {
		containerMetricsResource := ContainerMetricsResource{
			Name: container.Name,
		}
		cpuMetric, memoryMetric := int64(0), int64(0)
		if i < len(r.PodMetric.Containers) {
			metrics := r.PodMetric.Containers[i]
			cpuMetric = metrics.CPU
			memoryMetric = metrics.Memory
		} else {
			cpuMetric = unset
			memoryMetric = unset
		}
		containerMetricsResource.Requests = MetricsResource{
			CPURequest:    container.Requests.CPU,
			MemoryRequest: container.Requests.Memory,
			CPUUsed:       cpuMetric,
			MemoryUsed:    memoryMetric,
			AlertColor:    escapes.TextColorYellow,
		}
		containerMetricsResource.Limits = MetricsResource{
			CPURequest:    container.Limits.CPU,
			MemoryRequest: container.Limits.Memory,
			CPUUsed:       cpuMetric,
			MemoryUsed:    memoryMetric,
			AlertColor:    escapes.TextColorRed,
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

func (r PodMetricsResourceList) filterByAlert(alert alerts.Alert) PodMetricsResourceList {
	switch alert {
	case alerts.Any:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsAlerted() })
	case alerts.Memory:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsMemoryAlerted() })
	case alerts.CPU:
		return r.filterBy(func(c ContainerMetricsResources) bool { return c.IsCPUAlerted() })
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

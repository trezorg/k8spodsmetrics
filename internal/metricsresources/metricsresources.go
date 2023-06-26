package metricsresources

import (
	"bytes"
	"fmt"
	"sort"

	"text/template"

	escapes "github.com/snugfox/ansi-escapes"
	"github.com/trezorg/k8spodsmetrics/internal/humanize"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	"golang.org/x/exp/slog"
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
  Requests: 	{{ index $container.MetricsResources 0 }}
  Limits:    	{{ index $container.MetricsResources 1 }}
  {{ end -}}`))
	metricsPodsTemplate = template.Must(template.New("metricPods").Parse(`{{ range $index, $pod := . -}}
{{ $pod }}
{{ end -}}`))
)

type (
	ResourceType       uint
	PodMetricsResource struct {
		pods.PodResource
		podmetrics.PodMetric
	}

	MetricsResource struct {
		CPURequest    int64
		MemoryRequest int64
		CPUMetric     int64
		MemoryMetric  int64
		AlertColor    string
	}
	ContainerMetricsResource struct {
		Name             string
		MetricsResources []MetricsResource
	}
	ContainerMetricsResources []ContainerMetricsResource
)

func (c ContainerMetricsResource) IsAlerted() bool {
	for _, m := range c.MetricsResources {
		if m.CPUAlert() || m.MemoryAlert() {
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

func (m MetricsResource) CPUAlert() bool {
	return m.CPURequest > 0 && m.CPURequest <= m.CPUMetric
}
func (m MetricsResource) MemoryAlert() bool {
	return m.MemoryRequest > 0 && m.MemoryRequest <= m.MemoryMetric
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
	if metricResource.MemoryMetric == unset && metricResource.CPUMetric == unset {
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
		metricResource.CPUMetric,
		cpuEndColor,
		humanize.Bytes(metricResource.MemoryRequest),
		memoryStartColor,
		humanize.Bytes(metricResource.MemoryMetric),
		memoryEndColor,
	)
}

func (r PodMetricsResource) ContainersMetrics() ContainerMetricsResources {
	var containerMetricsResources ContainerMetricsResources
	alertColors := []string{escapes.TextColorYellow, escapes.TextColorRed}
	for i, container := range r.PodResource.Containers {
		containerMetricsResource := ContainerMetricsResource{
			Name:             container.Name,
			MetricsResources: make([]MetricsResource, 0, 2),
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
		for i, resource := range []pods.Resource{container.Requests, container.Limits} {
			alertColor := alertColors[i]
			containerMetricsResource.MetricsResources = append(containerMetricsResource.MetricsResources, MetricsResource{
				CPURequest:    resource.CPU,
				MemoryRequest: resource.Memory,
				CPUMetric:     cpuMetric,
				MemoryMetric:  memoryMetric,
				AlertColor:    alertColor,
			})
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

func (rList PodMetricsResourceList) String() string {
	var buffer bytes.Buffer
	if err := metricsPodsTemplate.Execute(&buffer, rList); err != nil {
		panic(err)
	}
	return buffer.String()
}

func (rList PodMetricsResourceList) filterAlerts() PodMetricsResourceList {
	var result PodMetricsResourceList
	for _, pod := range rList {
		if pod.ContainersMetrics().IsAlerted() {
			result = append(result, pod)
		}
	}
	return result
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
	sort.Slice(podMetricsResourceList, func(i, j int) bool {
		if podMetricsResourceList[i].Namespace < podMetricsResourceList[j].Namespace {
			return true
		}
		if podMetricsResourceList[i].Namespace > podMetricsResourceList[j].Namespace {
			return false
		}
		return podMetricsResourceList[i].Name < podMetricsResourceList[j].Name
	})
	return podMetricsResourceList
}

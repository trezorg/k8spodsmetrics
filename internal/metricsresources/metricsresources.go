package metricsresources

import (
	"bytes"
	"fmt"
	"sort"

	"text/template"

	escapes "github.com/snugfox/ansi-escapes"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	"golang.org/x/exp/slog"
)

type ResourceType uint

const (
	Requests ResourceType = iota
	Limits
	unset = int64(-1)
)

var (
	metricsResourcesTemplate = template.Must(template.New("metricResources").Parse(
		`Namespace: {{.Namespace}}. Name: {{.Name}}
Containers:
  {{ range $index, $container := .PodResource.Containers -}}
  Name:         {{ $container.Name }}
  Requests: 	{{ $.ContainerMetricsText $index 0 }}
  Limits:    	{{ $.ContainerMetricsText $index 1 -}}
  {{ end -}}`))
)

type PodMetricsResource struct {
	pods.PodResource
	podmetrics.PodMetric
}

type MetricsResource struct {
	CPURequest    int64
	MemoryRequest int64
	CPUMetric     int64
	MemoryMetric  int64
}

func resourceText(metricResource MetricsResource, alertColor string) string {
	cpuStartColor := ""
	cpuEndColor := ""
	if metricResource.CPURequest <= metricResource.CPUMetric {
		cpuStartColor = alertColor
		cpuEndColor = escapes.ColorReset
	}
	memoryStartColor := ""
	memoryEndColor := ""
	if metricResource.MemoryRequest <= metricResource.MemoryMetric {
		memoryStartColor = alertColor
		memoryEndColor = escapes.ColorReset
	}
	if metricResource.MemoryMetric == unset && metricResource.CPUMetric == unset {
		return fmt.Sprintf(
			"CPU=%d, Memory=%d",
			metricResource.CPURequest,
			metricResource.MemoryRequest,
		)
	}
	return fmt.Sprintf(
		"CPU=%d/%s%d%s, Memory=%d/%s%d%s",
		metricResource.CPURequest,
		cpuStartColor,
		metricResource.CPUMetric,
		cpuEndColor,
		metricResource.MemoryRequest,
		memoryStartColor,
		metricResource.MemoryMetric,
		memoryEndColor,
	)
}

func (r PodMetricsResource) ContainerMetrics(i int, resourceType ResourceType) *MetricsResource {
	if i >= len(r.PodResource.Containers) {
		return nil
	}
	container := r.PodResource.Containers[i]
	var resources pods.Resource
	switch resourceType {
	case Requests:
		resources = container.Requests
	case Limits:
		resources = container.Limits
	default:
		return nil
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
	return &MetricsResource{
		CPURequest:    resources.CPU,
		MemoryRequest: resources.Memory,
		CPUMetric:     cpuMetric,
		MemoryMetric:  memoryMetric,
	}
}

func (r PodMetricsResource) ContainerMetricsText(i int, resourceType ResourceType) string {
	containerMetrics := r.ContainerMetrics(i, resourceType)
	if containerMetrics == nil {
		return ""
	}
	alertColor := escapes.TextColorYellow
	if resourceType == Limits {
		alertColor = escapes.TextColorRed
	}
	return resourceText(*containerMetrics, alertColor)

}

func (r PodMetricsResource) String() string {
	var buffer bytes.Buffer
	if err := metricsResourcesTemplate.Execute(&buffer, r); err != nil {
		return ""
	}
	return buffer.String()
}

type PodMetricsResourceList []PodMetricsResource

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

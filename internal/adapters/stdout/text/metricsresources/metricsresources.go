package metricsresources

import (
	"bytes"
	"log/slog"
	"text/template"

	stdoutcommon "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/common"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

var metricsPodTemplate = template.Must(template.New("metricPod").Parse(`Name:		{{.PodResource.Name}}
Namespace:	{{.PodResource.Namespace}}
Node:		{{.NodeName}}
Containers:
  {{ range $index, $container := .ContainersMetrics -}}
  Name:         {{ $container.Name }}
  Requests: 	{{ (index $container.Requests).StringWithColor "yellow" }}
  Limits:    	{{ (index $container.Limits).StringWithColor "red" }}
  {{ end -}}`))

type Text func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	var buffer bytes.Buffer
	for _, pod := range list {
		if err := metricsPodTemplate.Execute(&buffer, pod); err != nil {
			panic(err)
		}
		if err := buffer.WriteByte('\n'); err != nil {
			panic(err)
		}
	}
	stdoutcommon.WriteStringLine(buffer.String())
}

func (j Text) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (Text) Error(err error) {
	slog.Error("text metrics resources output failed", "error", err)
}

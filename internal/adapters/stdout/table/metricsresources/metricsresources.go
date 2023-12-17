package metricsresources

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

type Table func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendHeader(table.Row{
		"Pod Name",
		"Namespace",
		"Node Name",
		"Container Name",
		"CPU",
		"CPU",
		"CPU",
		"Memory",
		"Memory",
		"Memory",
	}, rowConfigAutoMerge)
	t.AppendHeader(table.Row{"", "", "", "", "Request", "Limit", "Used", "Request", "Limit", "Used"}, rowConfigAutoMerge)
	for _, resource := range list {
		metrics := resource.ContainersMetrics()
		metric := metrics[0]
		t.AppendRow(table.Row{
			resource.Name,
			resource.Namespace,
			resource.NodeName,
			metrics[0].Name,
			metric.Requests.CPURequestString(),
			metric.Limits.CPURequestString(),
			metric.CPUUsed(),
			metric.Requests.MemoryRequestString(),
			metric.Limits.MemoryRequestString(),
			metric.MemoryUsed(),
		})
		for _, metric := range metrics[1:] {
			t.AppendRow(table.Row{
				"",
				"",
				"",
				metrics[0].Name,
				metric.Requests.CPURequestString(),
				metric.Limits.CPURequestString(),
				metric.CPUUsed(),
				metric.Requests.MemoryRequestString(),
				metric.Limits.MemoryRequestString(),
				metric.MemoryUsed(),
			})
		}
		t.AppendSeparator()
	}
	t.Render()
}

func (s Table) Success(list metricsresources.PodMetricsResourceList) {
	s(list)
}

func (s Table) Error(err error) {
	logger.Error("", err)
}

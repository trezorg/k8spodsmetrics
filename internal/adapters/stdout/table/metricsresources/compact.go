package metricsresources

import (
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	formatmetricsresources "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/format/metricsresources"
	servicemetricsresources "github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/resources"
)

const (
	compactNamespaceColumn   = 1
	compactPodColumn         = 2
	compactNodeColumn        = 3
	compactFirstMetricColumn = 4
	compactSecondMetricCol   = 5
	compactThirdMetricCol    = 6
	maxCompactColumns        = 7
)

func ToCompactTable(outputResources resources.Resources) Table {
	return Table(func(list servicemetricsresources.PodMetricsResourceList) {
		PrintCompactTo(os.Stdout, list, outputResources)
	})
}

func ToCompactWriter(outputResources resources.Resources) func(io.Writer, servicemetricsresources.PodMetricsResourceList) {
	return func(w io.Writer, list servicemetricsresources.PodMetricsResourceList) {
		PrintCompactTo(w, list, outputResources)
	}
}

func PrintCompactTo(w io.Writer, list servicemetricsresources.PodMetricsResourceList, outputResources resources.Resources) {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	configureCompactTable(t)
	t.AppendHeader(compactHeaderRow(outputResources))

	total := servicemetricsresources.ContainerMetricsResource{}
	rendered := 0
	for _, resource := range list {
		containers := resource.ContainersMetrics()
		if len(containers) == 0 {
			continue
		}
		aggregated := aggregatePodContainers(resource)
		t.AppendRow(compactPodRow(resource, aggregated, outputResources))
		accumulatePodTotal(&total, aggregated)
		rendered++
	}

	if rendered > 1 {
		t.AppendFooter(compactTotalRow(total, outputResources))
	}

	t.Render()
}

func compactHeaderRow(outputResources resources.Resources) table.Row {
	row := table.Row{"NAMESPACE", "POD", "NODE"}
	if outputResources.IsCPU() {
		row = append(row, "CPU(req/used/lim)")
	}
	if outputResources.IsMemory() {
		row = append(row, "MEM(req/used/lim)")
	}
	if outputResources.IsStorage() {
		row = append(row, "STO(req/used/lim)", "EPH(req/used/lim)")
	}
	return row
}

func compactPodRow(
	resource servicemetricsresources.PodMetricsResource,
	aggregated servicemetricsresources.ContainerMetricsResource,
	outputResources resources.Resources,
) table.Row {
	formatter := formatmetricsresources.NewContainer(aggregated)
	row := table.Row{
		resource.PodResource.Namespace,
		resource.PodResource.Name,
		resource.NodeName,
	}
	if outputResources.IsCPU() {
		row = append(row, formatter.CPUCompactString())
	}
	if outputResources.IsMemory() {
		row = append(row, formatter.MemoryCompactString())
	}
	if outputResources.IsStorage() {
		row = append(row, formatter.StorageCompactString(), formatter.StorageEphemeralCompactString())
	}
	return row
}

func compactTotalRow(total servicemetricsresources.ContainerMetricsResource, outputResources resources.Resources) table.Row {
	formatter := formatmetricsresources.NewContainer(total)
	row := table.Row{"TOTAL", "", ""}
	if outputResources.IsCPU() {
		row = append(row, formatter.CPUCompactString())
	}
	if outputResources.IsMemory() {
		row = append(row, formatter.MemoryCompactString())
	}
	if outputResources.IsStorage() {
		row = append(row, formatter.StorageCompactString(), formatter.StorageEphemeralCompactString())
	}
	return row
}

func aggregatePodContainers(resource servicemetricsresources.PodMetricsResource) servicemetricsresources.ContainerMetricsResource {
	var aggregated servicemetricsresources.ContainerMetricsResource
	for _, container := range resource.ContainersMetrics() {
		aggregated.Requests.CPURequest += container.Requests.CPURequest
		aggregated.Requests.MemoryRequest += container.Requests.MemoryRequest
		aggregated.Requests.StorageRequest += container.Requests.StorageRequest
		aggregated.Requests.StorageEphemeralRequest += container.Requests.StorageEphemeralRequest
		aggregated.Limits.CPURequest += container.Limits.CPURequest
		aggregated.Limits.MemoryRequest += container.Limits.MemoryRequest
		aggregated.Limits.StorageRequest += container.Limits.StorageRequest
		aggregated.Limits.StorageEphemeralRequest += container.Limits.StorageEphemeralRequest
		aggregated.Requests.CPUUsed += container.Requests.CPUUsed
		aggregated.Requests.MemoryUsed += container.Requests.MemoryUsed
		aggregated.Requests.StorageUsed += container.Requests.StorageUsed
		aggregated.Requests.StorageEphemeralUsed += container.Requests.StorageEphemeralUsed
		aggregated.Limits.CPUUsed += container.Limits.CPUUsed
		aggregated.Limits.MemoryUsed += container.Limits.MemoryUsed
		aggregated.Limits.StorageUsed += container.Limits.StorageUsed
		aggregated.Limits.StorageEphemeralUsed += container.Limits.StorageEphemeralUsed
	}
	return aggregated
}

func accumulatePodTotal(total *servicemetricsresources.ContainerMetricsResource, aggregated servicemetricsresources.ContainerMetricsResource) {
	total.Requests.CPURequest += aggregated.Requests.CPURequest
	total.Requests.MemoryRequest += aggregated.Requests.MemoryRequest
	total.Requests.StorageRequest += aggregated.Requests.StorageRequest
	total.Requests.StorageEphemeralRequest += aggregated.Requests.StorageEphemeralRequest
	total.Limits.CPURequest += aggregated.Limits.CPURequest
	total.Limits.MemoryRequest += aggregated.Limits.MemoryRequest
	total.Limits.StorageRequest += aggregated.Limits.StorageRequest
	total.Limits.StorageEphemeralRequest += aggregated.Limits.StorageEphemeralRequest
	total.Requests.CPUUsed += aggregated.Requests.CPUUsed
	total.Requests.MemoryUsed += aggregated.Requests.MemoryUsed
	total.Requests.StorageUsed += aggregated.Requests.StorageUsed
	total.Requests.StorageEphemeralUsed += aggregated.Requests.StorageEphemeralUsed
	total.Limits.CPUUsed += aggregated.Limits.CPUUsed
	total.Limits.MemoryUsed += aggregated.Limits.MemoryUsed
	total.Limits.StorageUsed += aggregated.Limits.StorageUsed
	total.Limits.StorageEphemeralUsed += aggregated.Limits.StorageEphemeralUsed
}

func configureCompactTable(t table.Writer) {
	applyTableStyle(t)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: compactNamespaceColumn, Align: text.AlignLeft},
		{Number: compactPodColumn, Align: text.AlignLeft},
		{Number: compactNodeColumn, Align: text.AlignLeft},
		{Number: compactFirstMetricColumn, Align: text.AlignRight},
		{Number: compactSecondMetricCol, Align: text.AlignRight},
		{Number: compactThirdMetricCol, Align: text.AlignRight},
		{Number: maxCompactColumns, Align: text.AlignRight},
	})
}

func applyTableStyle(t table.Writer) {
	t.SetStyle(table.StyleLight)
}

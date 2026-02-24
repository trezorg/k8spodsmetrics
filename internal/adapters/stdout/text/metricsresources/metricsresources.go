package metricsresources

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"

	formatmetricsresources "github.com/trezorg/k8spodsmetrics/internal/adapters/stdout/format/metricsresources"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
)

type Text func(list metricsresources.PodMetricsResourceList)

func Print(list metricsresources.PodMetricsResourceList) {
	PrintTo(os.Stdout, list)
}

func PrintTo(w io.Writer, list metricsresources.PodMetricsResourceList) {
	var buffer bytes.Buffer
	for _, pod := range list {
		_, _ = fmt.Fprintf(&buffer, "Name:\t\t%s\n", pod.PodResource.Name)
		_, _ = fmt.Fprintf(&buffer, "Namespace:\t%s\n", pod.PodResource.Namespace)
		_, _ = fmt.Fprintf(&buffer, "Node:\t\t%s\n", pod.NodeName)
		_, _ = fmt.Fprint(&buffer, "Containers:\n")
		for _, container := range pod.ContainersMetrics() {
			containerFormatter := formatmetricsresources.NewContainer(container)
			_, _ = fmt.Fprintf(&buffer, "  Name:\t\t%s\n", containerFormatter.Name())
			_, _ = fmt.Fprintf(&buffer, "  Requests:\t%s\n", containerFormatter.Requests().StringWithColor("yellow"))
			_, _ = fmt.Fprintf(&buffer, "  Limits:\t%s\n", containerFormatter.Limits().StringWithColor("red"))
		}
		_, _ = fmt.Fprintln(&buffer)
	}
	_, _ = io.WriteString(w, buffer.String())
	_, _ = io.WriteString(w, "\n")
}

func (Text) SuccessTo(w io.Writer, list metricsresources.PodMetricsResourceList) {
	PrintTo(w, list)
}

func (j Text) Success(list metricsresources.PodMetricsResourceList) {
	j(list)
}

func (Text) Error(err error) {
	slog.Error("text metrics resources output failed", "error", err)
}

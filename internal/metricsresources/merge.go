package metricsresources

import (
	"log/slog"

	"github.com/trezorg/k8spodsmetrics/pkg/podmetrics"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
)

func merge(podResourceList pods.PodResourceList, podMetricList podmetrics.PodMetricList) PodMetricsResourceList {
	podsMap := make(map[pods.NamespaceName]*PodMetricsResource)
	for _, pr := range podResourceList {
		podsMap[pods.NamespaceName{Namespace: pr.Namespace, Name: pr.Name}] = &PodMetricsResource{PodResource: pr}
	}
	unmatchedMetrics := 0
	for _, pm := range podMetricList {
		podMetricsResource, ok := podsMap[pods.NamespaceName{Namespace: pm.Namespace, Name: pm.Name}]
		if !ok {
			unmatchedMetrics++
			continue
		}
		podMetricsResource.PodMetric = pm
	}
	if unmatchedMetrics > 0 {
		slog.Debug("Skipped unmatched pod metrics", slog.Int("count", unmatchedMetrics))
	}
	podMetricsResourceList := make(PodMetricsResourceList, 0, len(podsMap))
	for _, podMetricsResource := range podsMap {
		podMetricsResourceList = append(podMetricsResourceList, *podMetricsResource)
	}
	return podMetricsResourceList
}

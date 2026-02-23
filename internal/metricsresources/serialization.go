package metricsresources

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

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

func (c ContainerMetricsResources) toOutput() ContainerMetricsResourcesOutputs {
	result := make(ContainerMetricsResourcesOutputs, 0, len(c))
	for _, container := range c {
		result = append(result, container.toOutput())
	}
	return result
}

func (r PodMetricsResource) toOutput() PodMetricsResourceOutput {
	containers := r.ContainersMetrics()
	return PodMetricsResourceOutput{
		r.PodResource.Name,
		r.PodResource.Namespace,
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

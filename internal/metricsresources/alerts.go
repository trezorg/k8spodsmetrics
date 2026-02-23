package metricsresources

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

func (m MetricsResource) CPUAlert() bool {
	return m.CPURequest > 0 && m.CPURequest <= m.CPUUsed
}

func (m MetricsResource) MemoryAlert() bool {
	return m.MemoryRequest > 0 && m.MemoryRequest <= m.MemoryUsed
}

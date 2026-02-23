package noderesources

func (n NodeResource) IsAlerted() bool {
	return n.IsCPUAlerted() || n.IsMemoryAlerted()
}

func (n NodeResource) IsMemoryAlerted() bool {
	return n.Memory <= n.MemoryLimit || n.Memory <= n.MemoryRequest
}

func (n NodeResource) IsMemoryRequestAlerted() bool {
	return n.Memory <= n.MemoryRequest
}

func (n NodeResource) IsMemoryLimitAlerted() bool {
	return n.Memory <= n.MemoryLimit
}

func (n NodeResource) IsCPUAlerted() bool {
	return n.CPU <= n.CPULimit || n.CPU <= n.CPURequest
}

func (n NodeResource) IsCPURequestAlerted() bool {
	return n.CPU <= n.CPURequest
}

func (n NodeResource) IsCPULimitAlerted() bool {
	return n.CPU <= n.CPULimit
}

func (n NodeResource) IsStorageAlerted() bool {
	if n.Storage <= 0 {
		return false
	}
	return (float64(n.UsedStorage)/float64(n.Storage))*100 > storageUsedPercentAlert
}

func (n NodeResource) IsStorageEphemeralAlerted() bool {
	if n.StorageEphemeral <= 0 {
		return false
	}
	return (float64(n.UsedStorageEphemeral)/float64(n.StorageEphemeral))*100 > storageEphemeralPercentAlert
}

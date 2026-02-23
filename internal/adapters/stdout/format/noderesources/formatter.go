package noderesources

import (
	"fmt"

	escapes "github.com/snugfox/ansi-escapes"
	"github.com/trezorg/k8spodsmetrics/internal/humanize"
	servicenoderesources "github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type Formatter struct {
	resource servicenoderesources.NodeResource
}

func New(resource servicenoderesources.NodeResource) Formatter {
	return Formatter{resource: resource}
}

func (f Formatter) MemoryTemplate() string {
	memoryRequestStartColor := ""
	memoryRequestEndColor := ""
	memoryLimitStartColor := ""
	memoryLimitEndColor := ""
	if f.resource.Memory <= f.resource.MemoryRequest {
		memoryRequestStartColor = escapes.TextColorYellow
		memoryRequestEndColor = escapes.ColorReset
	}
	if f.resource.Memory <= f.resource.MemoryLimit {
		memoryLimitStartColor = escapes.TextColorRed
		memoryLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"Node=%s/%s, Requests=%s%s%s, Limits=%s%s%s",
		humanize.Bytes(f.resource.Memory),
		humanize.Bytes(f.resource.UsedMemory),
		memoryRequestStartColor,
		humanize.Bytes(f.resource.MemoryRequest),
		memoryRequestEndColor,
		memoryLimitStartColor,
		humanize.Bytes(f.resource.MemoryLimit),
		memoryLimitEndColor,
	)
}

func (f Formatter) MemoryRequestString() string {
	memoryRequestStartColor := ""
	memoryRequestEndColor := ""
	if f.resource.Memory <= f.resource.MemoryRequest {
		memoryRequestStartColor = escapes.TextColorYellow
		memoryRequestEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryRequestStartColor,
		humanize.Bytes(f.resource.MemoryRequest),
		memoryRequestEndColor,
	)
}

func (f Formatter) MemoryLimitString() string {
	memoryLimitStartColor := ""
	memoryLimitEndColor := ""
	if f.resource.Memory <= f.resource.MemoryLimit {
		memoryLimitStartColor = escapes.TextColorRed
		memoryLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryLimitStartColor,
		humanize.Bytes(f.resource.MemoryLimit),
		memoryLimitEndColor,
	)
}

func (f Formatter) MemoryAvailableString() string {
	memoryAvailableStartColor := ""
	memoryAvailableEndColor := ""
	if f.resource.AvailableMemory == 0 {
		memoryAvailableStartColor = escapes.TextColorRed
		memoryAvailableEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryAvailableStartColor,
		humanize.Bytes(f.resource.AvailableMemory),
		memoryAvailableEndColor,
	)
}

func (f Formatter) MemoryFreeString() string {
	memoryFreeStartColor := ""
	memoryFreeEndColor := ""
	if f.resource.FreeMemory == 0 {
		memoryFreeStartColor = escapes.TextColorRed
		memoryFreeEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryFreeStartColor,
		humanize.Bytes(f.resource.FreeMemory),
		memoryFreeEndColor,
	)
}

func (f Formatter) MemoryNodeString() string {
	return humanize.Bytes(f.resource.Memory)
}

func (f Formatter) MemoryNodeUsedString() string {
	return humanize.Bytes(f.resource.UsedMemory)
}

func (f Formatter) MemoryNodeAllocatableString() string {
	return humanize.Bytes(f.resource.AllocatableMemory)
}

func (f Formatter) MemoryNodeAlocatableString() string {
	return f.MemoryNodeAllocatableString()
}

func (f Formatter) CPUTemplate() string {
	cpuRequestStartColor := ""
	cpuRequestEndColor := ""
	cpuLimitStartColor := ""
	cpuLimitEndColor := ""
	if f.resource.CPU <= f.resource.CPURequest {
		cpuRequestStartColor = escapes.TextColorYellow
		cpuRequestEndColor = escapes.ColorReset
	}
	if f.resource.CPU <= f.resource.CPULimit {
		cpuLimitStartColor = escapes.TextColorRed
		cpuLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"Node=%d/%d, Requests=%s%d%s, Limits=%s%d%s",
		f.resource.CPU,
		f.resource.UsedCPU,
		cpuRequestStartColor,
		f.resource.CPURequest,
		cpuRequestEndColor,
		cpuLimitStartColor,
		f.resource.CPULimit,
		cpuLimitEndColor,
	)
}

func (f Formatter) CPURequestString() string {
	cpuRequestStartColor := ""
	cpuRequestEndColor := ""
	if f.resource.CPU <= f.resource.CPURequest {
		cpuRequestStartColor = escapes.TextColorYellow
		cpuRequestEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuRequestStartColor,
		f.resource.CPURequest,
		cpuRequestEndColor,
	)
}

func (f Formatter) CPULimitString() string {
	cpuLimitStartColor := ""
	cpuLimitEndColor := ""
	if f.resource.CPU <= f.resource.CPULimit {
		cpuLimitStartColor = escapes.TextColorRed
		cpuLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuLimitStartColor,
		f.resource.CPULimit,
		cpuLimitEndColor,
	)
}

func (f Formatter) CPUAvailableString() string {
	cpuAvailableStartColor := ""
	cpuAvailableEndColor := ""
	if f.resource.AvailableCPU == 0 {
		cpuAvailableStartColor = escapes.TextColorRed
		cpuAvailableEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuAvailableStartColor,
		f.resource.AvailableCPU,
		cpuAvailableEndColor,
	)
}

func (f Formatter) CPUFreeString() string {
	cpuFreeStartColor := ""
	cpuFreeEndColor := ""
	if f.resource.FreeCPU == 0 {
		cpuFreeStartColor = escapes.TextColorRed
		cpuFreeEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuFreeStartColor,
		f.resource.FreeCPU,
		cpuFreeEndColor,
	)
}

func (f Formatter) StorageString() string {
	return humanize.Bytes(f.resource.Storage)
}

func (f Formatter) StorageAllocatableString() string {
	return humanize.Bytes(f.resource.AllocatableStorage)
}

func (f Formatter) StorageFreeString() string {
	return humanize.Bytes(f.resource.FreeStorage)
}

func (f Formatter) StorageUsedString() string {
	usedStorageStartColor := ""
	usedStorageEndColor := ""
	if f.resource.IsStorageAlerted() {
		usedStorageStartColor = escapes.TextColorRed
		usedStorageEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		usedStorageStartColor,
		humanize.Bytes(f.resource.UsedStorage),
		usedStorageEndColor,
	)
}

func (f Formatter) StorageEphemeralString() string {
	return humanize.Bytes(f.resource.StorageEphemeral)
}

func (f Formatter) StorageAllocatableEphemeralString() string {
	return humanize.Bytes(f.resource.AllocatableStorageEphemeral)
}

func (f Formatter) StorageFreeEphemeralString() string {
	return humanize.Bytes(f.resource.FreeStorageEphemeral)
}

func (f Formatter) StorageUsedEphemeralString() string {
	usedStorageStartColor := ""
	usedStorageEndColor := ""
	if f.resource.IsStorageEphemeralAlerted() {
		usedStorageStartColor = escapes.TextColorRed
		usedStorageEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		usedStorageStartColor,
		humanize.Bytes(f.resource.UsedStorageEphemeral),
		usedStorageEndColor,
	)
}

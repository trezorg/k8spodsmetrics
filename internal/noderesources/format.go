package noderesources

import (
	"fmt"

	escapes "github.com/snugfox/ansi-escapes"
	"github.com/trezorg/k8spodsmetrics/internal/humanize"
)

func (n NodeResource) MemoryTemplate() string {
	memoryRequestStartColor := ""
	memoryRequestEndColor := ""
	memoryLimitStartColor := ""
	memoryLimitEndColor := ""
	if n.Memory <= n.MemoryRequest {
		memoryRequestStartColor = escapes.TextColorYellow
		memoryRequestEndColor = escapes.ColorReset
	}
	if n.Memory <= n.MemoryLimit {
		memoryLimitStartColor = escapes.TextColorRed
		memoryLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"Node=%s/%s, Requests=%s%s%s, Limits=%s%s%s",
		humanize.Bytes(n.Memory),
		humanize.Bytes(n.UsedMemory),
		memoryRequestStartColor,
		humanize.Bytes(n.MemoryRequest),
		memoryRequestEndColor,
		memoryLimitStartColor,
		humanize.Bytes(n.MemoryLimit),
		memoryLimitEndColor,
	)
}

func (n NodeResource) MemoryRequestString() string {
	memoryRequestStartColor := ""
	memoryRequestEndColor := ""
	if n.Memory <= n.MemoryRequest {
		memoryRequestStartColor = escapes.TextColorYellow
		memoryRequestEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryRequestStartColor,
		humanize.Bytes(n.MemoryRequest),
		memoryRequestEndColor,
	)
}

func (n NodeResource) MemoryLimitString() string {
	memoryLimitStartColor := ""
	memoryLimitEndColor := ""
	if n.Memory <= n.MemoryLimit {
		memoryLimitStartColor = escapes.TextColorRed
		memoryLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryLimitStartColor,
		humanize.Bytes(n.MemoryLimit),
		memoryLimitEndColor,
	)
}

func (n NodeResource) MemoryAvailableString() string {
	memoryAvailableStartColor := ""
	memoryAvailableEndColor := ""
	if n.AvailableMemory == 0 {
		memoryAvailableStartColor = escapes.TextColorRed
		memoryAvailableEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryAvailableStartColor,
		humanize.Bytes(n.AvailableMemory),
		memoryAvailableEndColor,
	)
}

func (n NodeResource) MemoryFreeString() string {
	memoryFreeStartColor := ""
	memoryFreeEndColor := ""
	if n.FreeMemory == 0 {
		memoryFreeStartColor = escapes.TextColorRed
		memoryFreeEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		memoryFreeStartColor,
		humanize.Bytes(n.FreeMemory),
		memoryFreeEndColor,
	)
}

func (n NodeResource) MemoryNodeString() string {
	return humanize.Bytes(n.Memory)
}

func (n NodeResource) MemoryNodeUsedString() string {
	return humanize.Bytes(n.UsedMemory)
}

func (n NodeResource) MemoryNodeAllocatableString() string {
	return humanize.Bytes(n.AllocatableMemory)
}

// MemoryNodeAlocatableString is kept for compatibility with the previous typoed method name.
func (n NodeResource) MemoryNodeAlocatableString() string {
	return n.MemoryNodeAllocatableString()
}

func (n NodeResource) CPUTemplate() string {
	cpuRequestStartColor := ""
	cpuRequestEndColor := ""
	cpuLimitStartColor := ""
	cpuLimitEndColor := ""
	if n.CPU <= n.CPURequest {
		cpuRequestStartColor = escapes.TextColorYellow
		cpuRequestEndColor = escapes.ColorReset
	}
	if n.CPU <= n.CPULimit {
		cpuLimitStartColor = escapes.TextColorRed
		cpuLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"Node=%d/%d, Requests=%s%d%s, Limits=%s%d%s",
		n.CPU,
		n.UsedCPU,
		cpuRequestStartColor,
		n.CPURequest,
		cpuRequestEndColor,
		cpuLimitStartColor,
		n.CPULimit,
		cpuLimitEndColor,
	)
}

func (n NodeResource) CPURequestString() string {
	cpuRequestStartColor := ""
	cpuRequestEndColor := ""
	if n.CPU <= n.CPURequest {
		cpuRequestStartColor = escapes.TextColorYellow
		cpuRequestEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuRequestStartColor,
		n.CPURequest,
		cpuRequestEndColor,
	)
}

func (n NodeResource) CPULimitString() string {
	cpuLimitStartColor := ""
	cpuLimitEndColor := ""
	if n.CPU <= n.CPULimit {
		cpuLimitStartColor = escapes.TextColorRed
		cpuLimitEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuLimitStartColor,
		n.CPULimit,
		cpuLimitEndColor,
	)
}

func (n NodeResource) CPUAvailableString() string {
	cpuAvailableStartColor := ""
	cpuAvailableEndColor := ""
	if n.AvailableCPU == 0 {
		cpuAvailableStartColor = escapes.TextColorRed
		cpuAvailableEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuAvailableStartColor,
		n.AvailableCPU,
		cpuAvailableEndColor,
	)
}

func (n NodeResource) CPUFreeString() string {
	cpuFreeStartColor := ""
	cpuFreeEndColor := ""
	if n.FreeCPU == 0 {
		cpuFreeStartColor = escapes.TextColorRed
		cpuFreeEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%d%s",
		cpuFreeStartColor,
		n.FreeCPU,
		cpuFreeEndColor,
	)
}

func (n NodeResource) StorageString() string {
	return humanize.Bytes(n.Storage)
}

func (n NodeResource) StorageAllocatableString() string {
	return humanize.Bytes(n.AllocatableStorage)
}

func (n NodeResource) StorageFreeString() string {
	return humanize.Bytes(n.FreeStorage)
}

func (n NodeResource) StorageUsedString() string {
	usedStorageStartColor := ""
	usedStorageEndColor := ""
	if n.IsStorageAlerted() {
		usedStorageStartColor = escapes.TextColorRed
		usedStorageEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		usedStorageStartColor,
		humanize.Bytes(n.UsedStorage),
		usedStorageEndColor,
	)
}

func (n NodeResource) StorageEphemeralString() string {
	return humanize.Bytes(n.StorageEphemeral)
}

func (n NodeResource) StorageAllocatableEphemeralString() string {
	return humanize.Bytes(n.AllocatableStorageEphemeral)
}

func (n NodeResource) StorageFreeEphemeralString() string {
	return humanize.Bytes(n.FreeStorageEphemeral)
}

func (n NodeResource) StorageUsedEphemeralString() string {
	usedStorageStartColor := ""
	usedStorageEndColor := ""
	if n.IsStorageEphemeralAlerted() {
		usedStorageStartColor = escapes.TextColorRed
		usedStorageEndColor = escapes.ColorReset
	}
	return fmt.Sprintf(
		"%s%s%s",
		usedStorageStartColor,
		humanize.Bytes(n.UsedStorageEphemeral),
		usedStorageEndColor,
	)
}

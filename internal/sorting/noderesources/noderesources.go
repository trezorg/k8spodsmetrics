package noderesources

import (
	"fmt"

	choiceutil "github.com/trezorg/k8spodsmetrics/internal/choices"
)

type Sorting string

const (
	Name                        Sorting = "name"
	RequestCPU                  Sorting = "request_cpu"
	LimitCPU                    Sorting = "limit_cpu"
	UsedCPU                     Sorting = "used_cpu"
	TotalCPU                    Sorting = "total_cpu"
	AvailableCPU                Sorting = "available_cpu"
	FreeCPU                     Sorting = "free_cpu"
	RequestMemory               Sorting = "request_memory"
	LimitMemory                 Sorting = "limit_memory"
	UsedMemory                  Sorting = "used_memory"
	TotalMemory                 Sorting = "total_memory"
	AvailableMemory             Sorting = "available_memory"
	FreeMemory                  Sorting = "free_memory"
	Storage                     Sorting = "storage"
	AllocatableStorage          Sorting = "allocatable_storage"
	UsedStorage                 Sorting = "used_storage"
	FreeStorage                 Sorting = "free_storage"
	StorageEphemeral            Sorting = "storage_ephemeral"
	AllocatableStorageEphemeral Sorting = "allocatable_storage_ephemeral"
	UsedStorageEphemeral        Sorting = "used_storage_ephemeral"
	FreeStorageEphemeral        Sorting = "free_storage_ephemeral"
)

var choices = []Sorting{
	Name,
	RequestCPU,
	LimitCPU,
	UsedCPU,
	TotalCPU,
	AvailableCPU,
	FreeCPU,
	RequestMemory,
	LimitMemory,
	UsedMemory,
	TotalMemory,
	AvailableMemory,
	FreeMemory,
	Storage,
	AllocatableStorage,
	UsedStorage,
	StorageEphemeral,
	AllocatableStorageEphemeral,
	UsedStorageEphemeral,
	FreeStorage,
	FreeStorageEphemeral,
}

func Valid(o Sorting) error {
	if !choiceutil.Valid(o, choices) {
		return fmt.Errorf("sorting should be one of: %s", StringList(", "))
	}
	return nil
}

func StringList(separator string) string {
	return choiceutil.StringList(choices, separator)
}

func StringListDefault() string {
	return StringList(choiceutil.DefaultSeparator)
}

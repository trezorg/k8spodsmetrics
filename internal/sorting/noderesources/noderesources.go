package noderesources

import (
	"fmt"
	"slices"
	"strings"
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
	defaultSeparator            string  = "|"
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
}

func Valid(o Sorting) error {
	if !slices.Contains(choices, o) {
		return fmt.Errorf("sorting should be one of: %s", StringList(", "))
	}
	return nil
}

func StringList(separator string) string {
	builder := strings.Builder{}
	size := 0
	for i := range choices {
		size += len(choices[i])
	}
	size += (len(choices) - 1) * len(separator)
	builder.Grow(size)
	for i := range choices {
		_, _ = builder.WriteString(string(choices[i]))
		if i < len(choices)-1 {
			_, _ = builder.WriteString(separator)
		}
	}
	return builder.String()
}

func StringListDefault() string {
	return StringList(defaultSeparator)
}

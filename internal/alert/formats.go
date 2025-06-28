package alert

import (
	"fmt"
	"slices"
	"strings"
)

type Alert string

const (
	Any              Alert  = "any"
	Memory           Alert  = "memory"
	MemoryRequest    Alert  = "memory_request"
	MemoryLimit      Alert  = "memory_limit"
	CPU              Alert  = "cpu"
	CPURequest       Alert  = "cpu_request"
	CPULimit         Alert  = "cpu_limit"
	Storage          Alert  = "storage"
	StorageEphemeral Alert  = "storage_ephemeral"
	None             Alert  = "none"
	defaultSeparator string = "|"
)

var choices = []Alert{Any, Memory, MemoryLimit, MemoryRequest, CPU, CPULimit, CPURequest, Storage, StorageEphemeral, None}

func Valid(o Alert) error {
	if !slices.Contains(choices, o) {
		return fmt.Errorf("alert should be one of: %#v", choices)
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

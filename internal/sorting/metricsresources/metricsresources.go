package metricsresources

import (
	"fmt"
	"slices"
	"strings"
)

type Sorting string

const (
	Name             Sorting = "name"
	Namespace        Sorting = "namespace"
	RequestCPU       Sorting = "request_cpu"
	LimitCPU         Sorting = "limit_cpu"
	UsedCPU          Sorting = "used_cpu"
	RequestMemory    Sorting = "request_memory"
	LimitMemory      Sorting = "limit_memory"
	UsedMemory       Sorting = "used_memory"
	defaultSeparator string  = "|"
)

var choices = []Sorting{
	Name,
	Namespace,
	RequestCPU,
	LimitCPU,
	UsedCPU,
	RequestMemory,
	LimitMemory,
	UsedMemory,
}

func Valid(o Sorting) error {
	if !slices.Contains(choices, o) {
		return fmt.Errorf("Sorting should be one of: %s", StringList(", "))
	}
	return nil
}

func StringList(separator string) string {
	builder := strings.Builder{}
	size := 0
	for i := 0; i < len(choices); i++ {
		size += len(choices[i])
	}
	size += (len(choices) - 1) * len(separator)
	builder.Grow(size)
	for i := 0; i < len(choices); i++ {
		builder.WriteString(string(choices[i]))
		if i < len(choices)-1 {
			builder.WriteString(separator)
		}
	}
	return builder.String()
}

func StringListDefault() string {
	return StringList(defaultSeparator)
}

package metricsresources

import (
	"fmt"

	choiceutil "github.com/trezorg/k8spodsmetrics/internal/choices"
)

type Sorting string

const (
	Name                 Sorting = "name"
	Namespace            Sorting = "namespace"
	Node                 Sorting = "node"
	RequestCPU           Sorting = "request_cpu"
	LimitCPU             Sorting = "limit_cpu"
	UsedCPU              Sorting = "used_cpu"
	RequestMemory        Sorting = "request_memory"
	LimitMemory          Sorting = "limit_memory"
	UsedMemory           Sorting = "used_memory"
	UsedStorage          Sorting = "used_storage"
	UsedStorageEphemeral Sorting = "used_storage_ephemeral"
	defaultSeparator     string  = "|"
)

var choices = []Sorting{
	Name,
	Namespace,
	Node,
	RequestCPU,
	LimitCPU,
	UsedCPU,
	RequestMemory,
	LimitMemory,
	UsedMemory,
	UsedStorage,
	UsedStorageEphemeral,
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
	return StringList(defaultSeparator)
}

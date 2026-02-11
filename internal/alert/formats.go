package alert

import (
	"fmt"

	choiceutil "github.com/trezorg/k8spodsmetrics/internal/choices"
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
	if !choiceutil.Valid(o, choices) {
		return fmt.Errorf("alert should be one of: %#v", choices)
	}
	return nil
}

func StringList(separator string) string {
	return choiceutil.StringList(choices, separator)
}

func StringListDefault() string {
	return StringList(defaultSeparator)
}

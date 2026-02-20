package resources

import (
	"fmt"
	"slices"

	choiceutil "github.com/trezorg/k8spodsmetrics/internal/choices"
)

type (
	Resource  string
	Resources []Resource
)

const (
	Memory  Resource = "memory"
	CPU     Resource = "cpu"
	Storage Resource = "storage"
	All     Resource = "all"
)

var (
	choices            = []Resource{Memory, CPU, Storage, All}
	stringChoices      = ToStrings(choices...)
	ErrInvalidResource = fmt.Errorf("invalid resource. Should be one of: %#v", stringChoices)
)

func Compact(resources ...Resource) Resources {
	slices.Sort(resources)
	resources = slices.Compact(resources)
	if len(resources) > 0 {
		if slices.Contains(resources, All) {
			return []Resource{All}
		}
	}
	return resources
}

func Valid(resources ...Resource) error {
	for _, r := range resources {
		if !slices.Contains(choices, r) {
			return ErrInvalidResource
		}
	}
	return nil
}

func FromStrings(resources ...string) Resources {
	if len(resources) == 0 {
		return []Resource{All}
	}
	result := make([]Resource, 0, len(resources))
	for _, r := range resources {
		result = append(result, Resource(r))
	}
	return Compact(result...)
}

func ToStrings(resources ...Resource) []string {
	result := make([]string, 0, len(resources))
	for _, r := range resources {
		result = append(result, string(r))
	}
	return result
}

func (r Resources) IsCPU() bool {
	return slices.Contains(r, All) || slices.Contains(r, CPU)
}

func (r Resources) IsMemory() bool {
	return slices.Contains(r, All) || slices.Contains(r, Memory)
}

func (r Resources) IsStorage() bool {
	return slices.Contains(r, All) || slices.Contains(r, Storage)
}

func StringList(separator string) string {
	return choiceutil.StringList(choices, separator)
}

func StringListDefault() string {
	return StringList(choiceutil.DefaultSeparator)
}

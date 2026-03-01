package choices

import (
	"slices"
	"strings"
)

// DefaultSeparator is the standard separator used for formatting choice lists.
const DefaultSeparator = "|"

func Valid[T comparable](value T, choices []T) bool {
	return slices.Contains(choices, value)
}

func StringList[T ~string](choices []T, separator string) string {
	strChoices := make([]string, len(choices))
	for i := range choices {
		strChoices[i] = string(choices[i])
	}
	return strings.Join(strChoices, separator)
}

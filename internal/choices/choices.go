package choices

import (
	"slices"
	"strings"
)

func Valid[T comparable](value T, choices []T) bool {
	return slices.Contains(choices, value)
}

func StringList[T ~string](choices []T, separator string) string {
	builder := strings.Builder{}
	size := 0
	for i := range choices {
		size += len(choices[i])
	}
	size += (len(choices) - 1) * len(separator)
	if size > 0 {
		builder.Grow(size)
	}
	for i := range choices {
		_, _ = builder.WriteString(string(choices[i]))
		if i < len(choices)-1 {
			_, _ = builder.WriteString(separator)
		}
	}
	return builder.String()
}

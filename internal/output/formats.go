package output

import (
	"fmt"
	"slices"
	"strings"
)

type Output string

const (
	Table            Output = "table"
	JSON             Output = "json"
	String           Output = "string"
	Yaml             Output = "yaml"
	defaultSeparator string = "|"
)

var choices = []Output{Table, JSON, String, Yaml}

func Valid(o Output) error {
	if !slices.Contains(choices, o) {
		return fmt.Errorf("output should be one of: %#v", choices)
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

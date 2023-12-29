package output

import (
	"fmt"
	"slices"
	"strings"
)

type Output string

const (
	Table            Output = "table"
	Json             Output = "json"
	String           Output = "string"
	Yaml             Output = "yaml"
	defaultSeparator string = "|"
)

var choices = []Output{Table, Json, String, Yaml}

func Valid(o Output) error {
	if !slices.Contains(choices, o) {
		return fmt.Errorf("Output should be one of: %#v", choices)
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

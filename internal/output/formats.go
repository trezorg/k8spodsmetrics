package output

import (
	"fmt"
	"slices"
)

type Output string

const (
	Table  Output = "table"
	Json   Output = "json"
	String Output = "string"
	Yaml   Output = "yaml"
)

var choices = []Output{Table, Json, String, Yaml}

func Valid(o Output) error {
	if !slices.Contains(choices, o) {
		return fmt.Errorf("Output should be one of: %#v", choices)
	}
	return nil
}

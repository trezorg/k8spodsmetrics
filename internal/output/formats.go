package output

import (
	"fmt"

	choiceutil "github.com/trezorg/k8spodsmetrics/internal/choices"
)

type Output string

const (
	Table Output = "table"
	JSON  Output = "json"
	Text  Output = "text"
	Yaml  Output = "yaml"
)

var choices = []Output{Table, JSON, Text, Yaml}

func Valid(o Output) error {
	if !choiceutil.Valid(o, choices) {
		return fmt.Errorf("output should be one of: %#v", choices)
	}
	return nil
}

func StringList(separator string) string {
	return choiceutil.StringList(choices, separator)
}

func StringListDefault() string {
	return StringList(choiceutil.DefaultSeparator)
}

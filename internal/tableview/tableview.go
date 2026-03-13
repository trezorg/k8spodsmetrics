package tableview

import (
	"fmt"

	choiceutil "github.com/trezorg/k8spodsmetrics/internal/choices"
)

type View string

const (
	Expanded View = "expanded"
	Compact  View = "compact"
)

var choices = []View{Expanded, Compact}

func Valid(v View) error {
	if !choiceutil.Valid(v, choices) {
		return fmt.Errorf("table view should be one of: %#v", choices)
	}
	return nil
}

func StringList(separator string) string {
	return choiceutil.StringList(choices, separator)
}

func StringListDefault() string {
	return StringList(choiceutil.DefaultSeparator)
}

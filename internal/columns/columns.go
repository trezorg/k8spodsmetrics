package columns

import (
	"fmt"
	"slices"

	choiceutil "github.com/trezorg/k8spodsmetrics/internal/choices"
)

type Column string

const (
	Total       Column = "total"
	Allocatable Column = "allocatable"
	Used        Column = "used"
	Request     Column = "request"
	Limit       Column = "limit"
	Available   Column = "available"
	Free        Column = "free"
)

var nodeColumns = []Column{Total, Allocatable, Used, Request, Limit, Available, Free}
var podColumns = []Column{Request, Limit, Used}

func ValidForNodes(cols []Column) error {
	for _, col := range cols {
		if !slices.Contains(nodeColumns, col) {
			return fmt.Errorf("invalid column '%s' for nodes. Valid columns: %s", col, StringListNodeColumns())
		}
	}
	return nil
}

func ValidForPods(cols []Column) error {
	for _, col := range cols {
		if !slices.Contains(podColumns, col) {
			return fmt.Errorf("invalid column '%s' for pods. Valid columns: %s", col, StringListPodColumns())
		}
	}
	return nil
}

func FromStrings(vals ...string) []Column {
	result := make([]Column, 0, len(vals))
	for _, v := range vals {
		result = append(result, Column(v))
	}
	return result
}

func StringListNodeColumns() string {
	return choiceutil.StringList(nodeColumns, choiceutil.DefaultSeparator)
}

func StringListPodColumns() string {
	return choiceutil.StringList(podColumns, choiceutil.DefaultSeparator)
}

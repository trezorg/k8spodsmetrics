package resources

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func joinResources(resources Resources, separator string) string {
	parts := make([]string, 0, len(resources))
	for _, r := range resources {
		parts = append(parts, string(r))
	}
	return strings.Join(parts, separator)
}

func TestCompact(t *testing.T) {
	testData := []struct {
		in  Resources
		out Resources
	}{
		{
			in:  Resources{Memory, CPU, CPU, Memory},
			out: Resources{CPU, Memory},
		},
		{
			in:  Resources{Memory, CPU, CPU, Memory, All},
			out: Resources{All},
		},
	}

	for _, data := range testData {
		testName := fmt.Sprintf("%s => %s", joinResources(data.in, ","), joinResources(data.out, ","))
		t.Run(testName, func(t *testing.T) {
			require.Equal(t, data.out, Compact(data.in...))
		})
	}
}

func TestIsValid(t *testing.T) {
	testData := []struct {
		in  Resources
		err error
	}{
		{
			in:  Resources{Memory, CPU, CPU, Memory},
			err: nil,
		},
		{
			in:  Resources{Memory, CPU, CPU, Memory, All, Resource("xxx")},
			err: ErrInvalidResource,
		},
	}

	for _, data := range testData {
		testName := joinResources(data.in, ",")
		t.Run(testName, func(t *testing.T) {
			require.Equal(t, data.err, Valid(data.in...))
		})
	}
}

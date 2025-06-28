package humanize

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHumanizeBytes(t *testing.T) {
	type check[T Number] struct {
		val    T
		result string
	}
	checks := []check[int]{
		{10, "10B"},
		{1023, "1023B"},
		{1025, "1KiB"},
		{1024 * 1024, "1MiB"},
		{1024 * 1024 * 6, "6MiB"},
	}

	for i := range checks {
		t.Run(fmt.Sprintf("%v", checks[i]), func(t *testing.T) {
			require.Equal(t, checks[i].result, Bytes(checks[i].val))
		})
	}
}

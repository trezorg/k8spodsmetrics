package humanize

import (
	"fmt"
	"math"
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
		{1024, "1KiB"},
		{1025, "1KiB"},
		{1536, "1.5KiB"},
		{10 * 1024, "10KiB"},
		{1024 * 1024, "1MiB"},
		{1024 * 1024 * 6, "6MiB"},
		{-1024, "-1KiB"},
	}

	for i := range checks {
		t.Run(fmt.Sprintf("int:%v", checks[i]), func(t *testing.T) {
			require.Equal(t, checks[i].result, Bytes(checks[i].val))
		})
	}

	// float64 checks
	checksF := []check[float64]{
		{1536, "1.5KiB"},
		{0, "0B"},
	}
	for i := range checksF {
		t.Run(fmt.Sprintf("float:%v", checksF[i]), func(t *testing.T) {
			require.Equal(t, checksF[i].result, Bytes(checksF[i].val))
		})
	}

	// very large value clamps to highest unit (EiB)
	t.Run("clamp-to-max-unit", func(t *testing.T) {
		// 1024^7 bytes = 1024 EiB (use float64 to avoid integer overflow)
		big := math.Pow(1024, 7)
		require.Equal(t, "1024EiB", Bytes(big))
	})
}

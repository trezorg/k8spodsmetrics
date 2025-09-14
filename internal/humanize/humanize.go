package humanize

import (
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

var sizes = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

const (
	roundFactor       = 10.0
	intRoundThreshold = 10.0
)

// Bytes formats a byte size using IEC (base-1024) units.
// Examples: 0 -> "0B", 1024 -> "1KiB", 1536 -> "1.5KiB", -1024 -> "-1KiB".
func Bytes[V Number](s V) string {
	sf := float64(s)

	if sf == 0 || math.IsNaN(sf) {
		return "0B"
	}

	sign := ""
	if sf < 0 {
		sign = "-"
		sf = -sf
	}

	i := 0
	for sf >= 1024 && i < len(sizes)-1 {
		sf /= 1024
		i++
	}

	// Round to one decimal, then suppress trailing .0
	rounded := math.Round(sf*roundFactor) / roundFactor
	if math.Abs(rounded-math.Round(rounded)) < 1e-9 || rounded >= intRoundThreshold {
		return fmt.Sprintf("%s%d%s", sign, int(math.Round(rounded)), sizes[i])
	}
	return fmt.Sprintf("%s%0.1f%s", sign, rounded, sizes[i])
}

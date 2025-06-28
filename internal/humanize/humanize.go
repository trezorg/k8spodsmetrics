package humanize

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

const (
	epsilon = 0.001
)

var (
	sizes = [...]string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
)

func Bytes[V Number](s V) string {
	div := 1024

	sf := float64(s)

	var i int
	for i = range sizes {
		t := sf / float64(div)
		if t < 1 {
			break
		}
		sf = t
	}
	if sf-float64(int(sf)) < epsilon {
		return fmt.Sprintf("%d%s", int(sf), sizes[i])
	}
	return fmt.Sprintf("%0.1f%s", sf, sizes[i])
}

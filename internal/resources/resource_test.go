package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		testName := fmt.Sprintf("%s => %s", join(data.in, ","), join(data.out, ","))
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
		testName := join(data.in, ",")
		t.Run(testName, func(t *testing.T) {
			require.Equal(t, data.err, Valid(data.in...))
		})
	}
}

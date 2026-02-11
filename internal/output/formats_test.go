package output

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	t.Run("valid outputs", func(t *testing.T) {
		validOutputs := []Output{Table, JSON, String, Yaml}
		for _, out := range validOutputs {
			err := Valid(out)
			require.NoError(t, err)
		}
	})

	t.Run("invalid output", func(t *testing.T) {
		err := Valid(Output("invalid"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "output should be one of")
	})
}

func TestStringList(t *testing.T) {
	t.Run("default separator", func(t *testing.T) {
		list := StringListDefault()
		require.NotEmpty(t, list)
		require.Contains(t, list, "|")
	})

	t.Run("custom separator", func(t *testing.T) {
		list := StringList(",")
		require.NotEmpty(t, list)
		require.Contains(t, list, ",")
	})

	t.Run("all outputs included", func(t *testing.T) {
		list := StringListDefault()
		expectedOutputs := []string{"table", "json", "string", "yaml"}
		for _, out := range expectedOutputs {
			require.Contains(t, list, out)
		}
	})
}

func TestOutputString(t *testing.T) {
	t.Run("table output", func(t *testing.T) {
		require.Equal(t, "table", string(Table))
	})

	t.Run("json output", func(t *testing.T) {
		require.Equal(t, "json", string(JSON))
	})

	t.Run("yaml output", func(t *testing.T) {
		require.Equal(t, "yaml", string(Yaml))
	})

	t.Run("string output", func(t *testing.T) {
		require.Equal(t, "string", string(String))
	})
}

package tableview

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	t.Run("valid table views", func(t *testing.T) {
		for _, view := range []View{Expanded, Compact} {
			require.NoError(t, Valid(view))
		}
	})

	t.Run("invalid table view", func(t *testing.T) {
		err := Valid(View("invalid"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "table view should be one of")
	})
}

func TestStringListDefault(t *testing.T) {
	list := StringListDefault()
	require.Contains(t, list, string(Expanded))
	require.Contains(t, list, string(Compact))
	require.Contains(t, list, "|")
}

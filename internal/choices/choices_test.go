package choices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	t.Run("valid choice", func(t *testing.T) {
		choices := []string{"a", "b", "c"}
		require.True(t, Valid("a", choices))
		require.True(t, Valid("b", choices))
		require.True(t, Valid("c", choices))
	})

	t.Run("invalid choice", func(t *testing.T) {
		choices := []string{"a", "b", "c"}
		require.False(t, Valid("d", choices))
	})

	t.Run("empty choices", func(t *testing.T) {
		choices := []string{}
		require.False(t, Valid("a", choices))
	})

	t.Run("with integers", func(t *testing.T) {
		choices := []int{1, 2, 3}
		require.True(t, Valid(1, choices))
		require.False(t, Valid(4, choices))
	})
}

func TestStringList(t *testing.T) {
	t.Run("single element", func(t *testing.T) {
		choices := []string{"a"}
		result := StringList(choices, ",")
		require.Equal(t, "a", result)
	})

	t.Run("multiple elements", func(t *testing.T) {
		choices := []string{"a", "b", "c"}
		result := StringList(choices, ",")
		require.Equal(t, "a,b,c", result)
	})

	t.Run("empty list", func(t *testing.T) {
		choices := []string{}
		result := StringList(choices, ",")
		require.Equal(t, "", result)
	})

	t.Run("custom separator", func(t *testing.T) {
		choices := []string{"a", "b", "c"}
		result := StringList(choices, "|")
		require.Equal(t, "a|b|c", result)
	})
}

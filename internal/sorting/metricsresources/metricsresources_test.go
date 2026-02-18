package metricsresources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	testCases := []struct {
		name        string
		sorting     Sorting
		expectError bool
	}{
		{"valid name", Name, false},
		{"valid namespace", Namespace, false},
		{"valid node", Node, false},
		{"valid request_cpu", RequestCPU, false},
		{"valid limit_cpu", LimitCPU, false},
		{"valid used_cpu", UsedCPU, false},
		{"valid request_memory", RequestMemory, false},
		{"valid limit_memory", LimitMemory, false},
		{"valid used_memory", UsedMemory, false},
		{"valid used_storage", UsedStorage, false},
		{"valid used_storage_ephemeral", UsedStorageEphemeral, false},
		{"invalid empty", Sorting(""), true},
		{"invalid unknown", Sorting("unknown_field"), true},
		{"invalid case", Sorting("NAME"), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := Valid(tc.sorting)
			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), "sorting should be one of")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStringList(t *testing.T) {
	t.Run("comma separator", func(t *testing.T) {
		result := StringList(", ")
		require.Contains(t, result, "name")
		require.Contains(t, result, "namespace")
		require.Contains(t, result, "request_cpu")
		require.Contains(t, result, "used_memory")
	})

	t.Run("pipe separator", func(t *testing.T) {
		result := StringList("|")
		require.Contains(t, result, "name")
		require.Contains(t, result, "|")
	})
}

func TestStringListDefault(t *testing.T) {
	result := StringListDefault()
	require.Contains(t, result, "name")
	require.Contains(t, result, "|")
}

func TestSortingConstants(t *testing.T) {
	require.Equal(t, Sorting("name"), Name)
	require.Equal(t, Sorting("namespace"), Namespace)
	require.Equal(t, Sorting("node"), Node)
	require.Equal(t, Sorting("request_cpu"), RequestCPU)
	require.Equal(t, Sorting("limit_cpu"), LimitCPU)
	require.Equal(t, Sorting("used_cpu"), UsedCPU)
	require.Equal(t, Sorting("request_memory"), RequestMemory)
	require.Equal(t, Sorting("limit_memory"), LimitMemory)
	require.Equal(t, Sorting("used_memory"), UsedMemory)
	require.Equal(t, Sorting("used_storage"), UsedStorage)
	require.Equal(t, Sorting("used_storage_ephemeral"), UsedStorageEphemeral)
}

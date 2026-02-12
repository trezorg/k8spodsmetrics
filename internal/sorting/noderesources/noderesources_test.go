package noderesources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	validSortings := []Sorting{
		Name, RequestCPU, LimitCPU, UsedCPU, TotalCPU, AvailableCPU, FreeCPU,
		RequestMemory, LimitMemory, UsedMemory, TotalMemory, AvailableMemory, FreeMemory,
		Storage, AllocatableStorage, UsedStorage, FreeStorage,
		StorageEphemeral, AllocatableStorageEphemeral, UsedStorageEphemeral, FreeStorage, FreeStorageEphemeral,
	}

	for _, s := range validSortings {
		t.Run("valid_"+string(s), func(t *testing.T) {
			t.Parallel()
			err := Valid(s)
			require.NoError(t, err)
		})
	}

	invalidSortings := []Sorting{"", "unknown", "NAME", "invalid_field"}
	for _, s := range invalidSortings {
		t.Run("invalid_"+string(s), func(t *testing.T) {
			t.Parallel()
			err := Valid(s)
			require.Error(t, err)
			require.Contains(t, err.Error(), "sorting should be one of")
		})
	}
}

func TestStringList(t *testing.T) {
	t.Run("comma separator", func(t *testing.T) {
		result := StringList(", ")
		require.Contains(t, result, "name")
		require.Contains(t, result, "request_cpu")
		require.Contains(t, result, "total_memory")
		require.Contains(t, result, "free_storage")
	})

	t.Run("custom separator", func(t *testing.T) {
		result := StringList(" | ")
		require.Contains(t, result, "name")
		require.Contains(t, result, " | ")
	})
}

func TestStringListDefault(t *testing.T) {
	result := StringListDefault()
	require.Contains(t, result, "name")
	require.Contains(t, result, "|")
}

func TestSortingConstants(t *testing.T) {
	require.Equal(t, Sorting("name"), Name)
	require.Equal(t, Sorting("request_cpu"), RequestCPU)
	require.Equal(t, Sorting("limit_cpu"), LimitCPU)
	require.Equal(t, Sorting("used_cpu"), UsedCPU)
	require.Equal(t, Sorting("total_cpu"), TotalCPU)
	require.Equal(t, Sorting("available_cpu"), AvailableCPU)
	require.Equal(t, Sorting("free_cpu"), FreeCPU)
	require.Equal(t, Sorting("request_memory"), RequestMemory)
	require.Equal(t, Sorting("limit_memory"), LimitMemory)
	require.Equal(t, Sorting("used_memory"), UsedMemory)
	require.Equal(t, Sorting("total_memory"), TotalMemory)
	require.Equal(t, Sorting("available_memory"), AvailableMemory)
	require.Equal(t, Sorting("free_memory"), FreeMemory)
	require.Equal(t, Sorting("storage"), Storage)
	require.Equal(t, Sorting("allocatable_storage"), AllocatableStorage)
	require.Equal(t, Sorting("used_storage"), UsedStorage)
	require.Equal(t, Sorting("free_storage"), FreeStorage)
}

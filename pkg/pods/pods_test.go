package pods

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildFieldSelector(t *testing.T) {
	testCases := []struct {
		name     string
		filter   PodFilter
		expected string
	}{
		{
			name:     "no node name",
			filter:   PodFilter{FieldSelector: "metadata.name=foo"},
			expected: "metadata.name=foo",
		},
		{
			name:     "only node name",
			filter:   PodFilter{NodeName: "worker-1"},
			expected: "spec.nodeName=worker-1",
		},
		{
			name:     "field selector with node",
			filter:   PodFilter{FieldSelector: "status.phase=Running", NodeName: "worker-1"},
			expected: "status.phase=Running,spec.nodeName=worker-1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := buildFieldSelector(tc.filter)
			require.Equal(t, tc.expected, actual)
		})
	}
}

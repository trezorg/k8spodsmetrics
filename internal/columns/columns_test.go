package columns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidForNodes(t *testing.T) {
	tests := []struct {
		name    string
		columns []Column
		wantErr bool
	}{
		{
			name:    "valid single column",
			columns: []Column{Total},
			wantErr: false,
		},
		{
			name:    "valid multiple columns",
			columns: []Column{Total, Allocatable, Free},
			wantErr: false,
		},
		{
			name:    "all node columns",
			columns: nodeColumns,
			wantErr: false,
		},
		{
			name:    "invalid column",
			columns: []Column{Column("invalid")},
			wantErr: true,
		},
		{
			name:    "pod column invalid for nodes",
			columns: []Column{Request, Limit, Used},
			wantErr: false, // Request, Limit, Used are also valid for nodes
		},
		{
			name:    "empty columns",
			columns: []Column{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidForNodes(tt.columns)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidForPods(t *testing.T) {
	tests := []struct {
		name    string
		columns []Column
		wantErr bool
	}{
		{
			name:    "valid single column",
			columns: []Column{Request},
			wantErr: false,
		},
		{
			name:    "valid multiple columns",
			columns: []Column{Request, Limit, Used},
			wantErr: false,
		},
		{
			name:    "all pod columns",
			columns: podColumns,
			wantErr: false,
		},
		{
			name:    "invalid column",
			columns: []Column{Column("invalid")},
			wantErr: true,
		},
		{
			name:    "node column invalid for pods",
			columns: []Column{Total},
			wantErr: true,
		},
		{
			name:    "available invalid for pods",
			columns: []Column{Available},
			wantErr: true,
		},
		{
			name:    "empty columns",
			columns: []Column{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidForPods(tt.columns)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFromStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []Column
	}{
		{
			name:     "single column",
			input:    []string{"total"},
			expected: []Column{Total},
		},
		{
			name:     "multiple columns",
			input:    []string{"request", "limit", "used"},
			expected: []Column{Request, Limit, Used},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []Column{},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: []Column{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromStrings(tt.input...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStringListNodeColumns(t *testing.T) {
	result := StringListNodeColumns()
	assert.Contains(t, result, "total")
	assert.Contains(t, result, "allocatable")
	assert.Contains(t, result, "used")
	assert.Contains(t, result, "request")
	assert.Contains(t, result, "limit")
	assert.Contains(t, result, "available")
	assert.Contains(t, result, "free")
}

func TestStringListPodColumns(t *testing.T) {
	result := StringListPodColumns()
	assert.Contains(t, result, "request")
	assert.Contains(t, result, "limit")
	assert.Contains(t, result, "used")
	assert.NotContains(t, result, "total")
	assert.NotContains(t, result, "available")
}

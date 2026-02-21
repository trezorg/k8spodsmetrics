package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLevelFromString(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "debug", input: "DEBUG"},
		{name: "info lower case", input: "info"},
		{name: "empty", input: ""},
		{name: "warn alias", input: "warning"},
		{name: "error", input: "error"},
		{name: "unknown", input: "trace", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := levelFromString(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestInitLogger(t *testing.T) {
	err := initLogger("debug")
	require.NoError(t, err)

	err = initLogger("unknown")
	require.ErrorContains(t, err, "unknown log level")
}

package noderesources

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func nodeResourceList(name string) NodeResourceList {
	return []NodeResource{
		{
			Name:          name,
			CPU:           1024,
			Memory:        1024,
			CPURequest:    512,
			MemoryRequest: 512,
			CPULimit:      512,
			MemoryLimit:   512,
		},
	}
}

func TestStringify(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})))
	nodeResourceList := nodeResourceList("foo")
	text := nodeResourceList.String()
	require.Greater(t, len(text), 0)
	require.Contains(t, text, "/", text)
}

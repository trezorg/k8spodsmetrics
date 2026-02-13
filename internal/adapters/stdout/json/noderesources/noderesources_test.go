package noderesources

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

func TestPrint(t *testing.T) {
	t.Run("prints valid JSON", func(t *testing.T) {
		list := noderesources.NodeResourceList{
			{
				Name:              "node-1",
				CPU:               4000,
				AllocatableCPU:    3900,
				Memory:            8 * 1024 * 1024 * 1024,
				AllocatableMemory: 7 * 1024 * 1024 * 1024,
			},
		}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		Print(list)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		var decoded noderesources.NodeResourceListEnvelop
		err := json.Unmarshal([]byte(output), &decoded)
		require.NoError(t, err)
		require.Len(t, decoded.Items, 1)
		require.Equal(t, "node-1", decoded.Items[0].Name)
	})

	t.Run("prints empty list", func(t *testing.T) {
		list := noderesources.NodeResourceList{}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		Print(list)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		var decoded noderesources.NodeResourceListEnvelop
		err := json.Unmarshal([]byte(output), &decoded)
		require.NoError(t, err)
		require.Empty(t, decoded.Items)
	})
}

func TestJSON_Success(t *testing.T) {
	t.Run("calls Print", func(t *testing.T) {
		list := noderesources.NodeResourceList{
			{
				Name:           "node-1",
				CPU:            4000,
				AllocatableCPU: 3900,
			},
		}

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		formatter := JSON(Print)
		formatter.Success(list)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		var decoded noderesources.NodeResourceListEnvelop
		err := json.Unmarshal([]byte(output), &decoded)
		require.NoError(t, err)
		require.Len(t, decoded.Items, 1)
	})
}

func TestJSON_Error(t *testing.T) {
	t.Run("logs error without panicking", func(t *testing.T) {
		formatter := JSON(Print)
		err := errors.New("test error")
		require.NotPanics(t, func() {
			formatter.Error(err)
		})
	})
}

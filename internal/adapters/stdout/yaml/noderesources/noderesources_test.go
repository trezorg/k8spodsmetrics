package noderesources

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

func TestPrint(t *testing.T) {
	t.Run("prints valid YAML", func(t *testing.T) {
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

		require.NotEmpty(t, output)
		require.Contains(t, output, "items:")
		require.Contains(t, output, "- name: node-1")
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

		require.NotEmpty(t, output)
	})
}

func TestYaml_Success(t *testing.T) {
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

		formatter := Yaml(Print)
		formatter.Success(list)

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		require.NotEmpty(t, output)
		require.Contains(t, output, "items:")
	})
}

func TestYaml_Error(t *testing.T) {
	t.Run("logs error without panicking", func(t *testing.T) {
		formatter := Yaml(Print)
		err := errors.New("test error")
		require.NotPanics(t, func() {
			formatter.Error(err)
		})
	})
}

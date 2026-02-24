package noderesources

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trezorg/k8spodsmetrics/internal/noderesources"
)

type mockErrorProcessor struct {
	called bool
}

func (m *mockErrorProcessor) Error(err error) {
	m.called = true
}

func TestNewScreenSuccessWriter(t *testing.T) {
	t.Run("creates success writer", func(t *testing.T) {
		called := false
		writer := NewScreenSuccessWriter(func(_ io.Writer, _ noderesources.NodeResourceList) {
			called = true
		})
		require.NotNil(t, writer)

		list := noderesources.NodeResourceList{
			{
				Name: "node-1",
				CPU:  4000,
			},
		}

		writer.Success(list)
		require.True(t, called)
	})
}

func TestNewScreenErrorWriter(t *testing.T) {
	t.Run("creates error writer", func(t *testing.T) {
		mock := &mockErrorProcessor{}
		writer := NewScreenErrorWriter(mock)
		require.NotNil(t, writer)

		err := errors.New("test error")
		writer.Error(err)
		require.True(t, mock.called)
	})
}

func TestScreenSuccessWriter_Success(t *testing.T) {
	t.Run("calls underlying writer", func(t *testing.T) {
		called := false
		writer := NewScreenSuccessWriter(func(_ io.Writer, _ noderesources.NodeResourceList) {
			called = true
		})
		writer.Success(noderesources.NodeResourceList{})
		require.True(t, called)
	})
}

func TestScreenErrorWriter_Error(t *testing.T) {
	t.Run("calls underlying writer", func(t *testing.T) {
		mock := &mockErrorProcessor{}
		writer := NewScreenErrorWriter(mock)
		err := errors.New("test error")
		writer.Error(err)
		require.True(t, mock.called)
	})
}

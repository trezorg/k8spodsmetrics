package stdin

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trezorg/k8spodsmetrics/internal/columns"
	"github.com/trezorg/k8spodsmetrics/internal/config"
	"github.com/trezorg/k8spodsmetrics/internal/output"
)

func TestParseColumnsForOutput(t *testing.T) {
	t.Run("non table output skips parsing and validation", func(t *testing.T) {
		parseCalled := false
		validateCalled := false

		cols, err := parseColumnsForOutput(
			output.JSON,
			[]string{"invalid"},
			func(_ []string) []columns.Column {
				parseCalled = true
				return []columns.Column{columns.Column("invalid")}
			},
			func(_ []columns.Column) error {
				validateCalled = true
				return errors.New("should not be called")
			},
		)

		require.NoError(t, err)
		require.Nil(t, cols)
		require.False(t, parseCalled)
		require.False(t, validateCalled)
	})

	t.Run("table output parses and validates columns", func(t *testing.T) {
		parseCalled := false
		validateCalled := false

		cols, err := parseColumnsForOutput(
			output.Table,
			[]string{"used"},
			func(_ []string) []columns.Column {
				parseCalled = true
				return []columns.Column{columns.Used}
			},
			func(parsed []columns.Column) error {
				validateCalled = true
				require.Equal(t, []columns.Column{columns.Used}, parsed)
				return nil
			},
		)

		require.NoError(t, err)
		require.Equal(t, []columns.Column{columns.Used}, cols)
		require.True(t, parseCalled)
		require.True(t, validateCalled)
	})

	t.Run("table output returns validation error", func(t *testing.T) {
		expectedErr := errors.New("invalid columns")

		cols, err := parseColumnsForOutput(
			output.Table,
			[]string{"invalid"},
			func(_ []string) []columns.Column {
				return []columns.Column{columns.Column("invalid")}
			},
			func(_ []columns.Column) error {
				return expectedErr
			},
		)

		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, cols)
	})
}

func TestApplyCommonConfig(t *testing.T) {
	t.Run("uses file bool when flag is not explicitly set", func(t *testing.T) {
		cfg := &commonConfig{WatchMetrics: false}
		fileCfg := &config.Config{Common: config.Common{WatchMetrics: true}}

		merged := applyCommonConfig(cfg, fileCfg, false)
		require.True(t, merged.WatchMetrics)
	})

	t.Run("keeps cli bool when flag is explicitly set", func(t *testing.T) {
		cfg := &commonConfig{WatchMetrics: false}
		fileCfg := &config.Config{Common: config.Common{WatchMetrics: true}}

		merged := applyCommonConfig(cfg, fileCfg, true)
		require.False(t, merged.WatchMetrics)
	})
}

func TestApplyPodsConfig(t *testing.T) {
	t.Run("uses file reverse when flag is not explicitly set", func(t *testing.T) {
		cfg := &podConfig{Reverse: false}
		fileCfg := &config.Config{Pods: config.Pods{Reverse: true}}

		merged := applyPodsConfig(cfg, fileCfg, false)
		require.True(t, merged.Reverse)
	})

	t.Run("keeps cli reverse when flag is explicitly set", func(t *testing.T) {
		cfg := &podConfig{Reverse: false}
		fileCfg := &config.Config{Pods: config.Pods{Reverse: true}}

		merged := applyPodsConfig(cfg, fileCfg, true)
		require.False(t, merged.Reverse)
	})
}

func TestApplySummaryConfig(t *testing.T) {
	t.Run("uses file reverse when flag is not explicitly set", func(t *testing.T) {
		cfg := &summaryConfig{Reverse: false}
		fileCfg := &config.Config{Summary: config.Summary{Reverse: true}}

		merged := applySummaryConfig(cfg, fileCfg, false)
		require.True(t, merged.Reverse)
	})

	t.Run("keeps cli reverse when flag is explicitly set", func(t *testing.T) {
		cfg := &summaryConfig{Reverse: false}
		fileCfg := &config.Config{Summary: config.Summary{Reverse: true}}

		merged := applySummaryConfig(cfg, fileCfg, true)
		require.False(t, merged.Reverse)
	})
}

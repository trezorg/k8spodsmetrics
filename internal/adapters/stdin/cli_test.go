package stdin

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trezorg/k8spodsmetrics/internal/columns"
	"github.com/trezorg/k8spodsmetrics/internal/config"
	"github.com/trezorg/k8spodsmetrics/internal/output"
	"github.com/trezorg/k8spodsmetrics/internal/tableview"
	"github.com/urfave/cli/v2"
)

func TestCommonFlagsAlertNaming(t *testing.T) {
	flags := commonFlags(&commonConfig{})

	var alertFlag *cli.StringFlag
	for _, flag := range flags {
		f, ok := flag.(*cli.StringFlag)
		if !ok {
			continue
		}
		if f.Name == "alert" {
			alertFlag = f
			break
		}
	}

	require.NotNil(t, alertFlag)
	require.Contains(t, alertFlag.Aliases, "alerts")
	require.Contains(t, alertFlag.Aliases, "a")
}

func TestCommonFlagsTableViewNaming(t *testing.T) {
	flags := commonFlags(&commonConfig{})

	var tableViewFlag *cli.StringFlag
	for _, flag := range flags {
		f, ok := flag.(*cli.StringFlag)
		if !ok {
			continue
		}
		if f.Name == "table-view" {
			tableViewFlag = f
			break
		}
	}

	require.NotNil(t, tableViewFlag)
	require.Equal(t, string(tableview.Compact), tableViewFlag.Value)
	require.Contains(t, tableViewFlag.Usage, string(tableview.Compact))
}

func TestPodsFlagsResourcesNaming(t *testing.T) {
	flags := podsFlags()

	var resourcesFlag *cli.StringSliceFlag
	for _, flag := range flags {
		f, ok := flag.(*cli.StringSliceFlag)
		if !ok {
			continue
		}
		if f.Name == "resources" {
			resourcesFlag = f
			break
		}
	}

	require.NotNil(t, resourcesFlag)
	require.Contains(t, resourcesFlag.Aliases, "resource")
	require.Contains(t, resourcesFlag.Aliases, "res")
}

func TestSummaryFlagsResourcesNaming(t *testing.T) {
	flags := summaryFlags()

	var resourcesFlag *cli.StringSliceFlag
	for _, flag := range flags {
		f, ok := flag.(*cli.StringSliceFlag)
		if !ok {
			continue
		}
		if f.Name == "resources" {
			resourcesFlag = f
			break
		}
	}

	require.NotNil(t, resourcesFlag)
	require.Contains(t, resourcesFlag.Aliases, "resource")
	require.Contains(t, resourcesFlag.Aliases, "res")
}

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

func TestValidateTableViewColumns(t *testing.T) {
	t.Run("expanded allows columns", func(t *testing.T) {
		require.NoError(t, validateTableViewColumns(tableview.Expanded, []string{"used"}))
	})

	t.Run("default compact rejects columns", func(t *testing.T) {
		err := validateTableViewColumns(tableview.Compact, []string{"used"})
		require.ErrorContains(t, err, "--columns is only supported")
	})

	t.Run("compact rejects columns", func(t *testing.T) {
		err := validateTableViewColumns(tableview.Compact, []string{"used"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "--columns is only supported")
	})

	t.Run("compact allows empty columns", func(t *testing.T) {
		require.NoError(t, validateTableViewColumns(tableview.Compact, nil))
	})
}

func TestApplyCommonConfig(t *testing.T) {
	t.Run("uses file bool when flag is not explicitly set", func(t *testing.T) {
		cfg := &commonConfig{WatchMetrics: false}
		fileCfg := &config.Config{Common: config.Common{WatchMetrics: true}}

		merged := applyCommonConfig(cfg, fileCfg, false, false)
		require.True(t, merged.WatchMetrics)
	})

	t.Run("keeps cli bool when flag is explicitly set", func(t *testing.T) {
		cfg := &commonConfig{WatchMetrics: false}
		fileCfg := &config.Config{Common: config.Common{WatchMetrics: true}}

		merged := applyCommonConfig(cfg, fileCfg, true, false)
		require.False(t, merged.WatchMetrics)
	})

	t.Run("uses file timeout when flag is not explicitly set", func(t *testing.T) {
		cfg := &commonConfig{Timeout: defaultTimeoutSeconds}
		fileCfg := &config.Config{Common: config.Common{Timeout: 45}}

		merged := applyCommonConfig(cfg, fileCfg, false, false)
		require.Equal(t, uint(45), merged.Timeout)
	})

	t.Run("keeps cli timeout when flag is explicitly set", func(t *testing.T) {
		cfg := &commonConfig{Timeout: 12}
		fileCfg := &config.Config{Common: config.Common{Timeout: 45}}

		merged := applyCommonConfig(cfg, fileCfg, false, true)
		require.Equal(t, uint(12), merged.Timeout)
	})

	t.Run("uses default timeout when unset in cli and file", func(t *testing.T) {
		cfg := &commonConfig{Timeout: defaultTimeoutSeconds}

		merged := applyCommonConfig(cfg, nil, false, false)
		require.Equal(t, uint(defaultTimeoutSeconds), merged.Timeout)
	})

	t.Run("uses file table view when flag is not explicitly set", func(t *testing.T) {
		cfg := &commonConfig{TableView: ""}
		fileCfg := &config.Config{Common: config.Common{TableView: string(tableview.Compact)}}

		merged := applyCommonConfig(cfg, fileCfg, false, false)
		require.Equal(t, string(tableview.Compact), merged.TableView)
	})
}

func TestResolveCommonConfig(t *testing.T) {
	t.Run("uses file-backed defaults when flags are omitted", func(t *testing.T) {
		cfg := commonConfig{
			Output:      string(output.Table),
			TableView:   string(tableview.Expanded),
			Alert:       "none",
			WatchPeriod: defaultWatchPeriodSeconds,
			Columns:     []string{"used"},
			fileConfig: &config.Config{Common: config.Common{
				Output:      string(output.JSON),
				TableView:   string(tableview.Compact),
				Alert:       "cpu",
				WatchPeriod: 12,
				Columns:     []string{"limit"},
			}},
		}

		resolved := resolveCommonConfig(cfg, actionFlags{})

		require.Equal(t, string(output.JSON), resolved.Output)
		require.Equal(t, string(tableview.Compact), resolved.TableView)
		require.Equal(t, "cpu", resolved.Alert)
		require.Equal(t, uint(12), resolved.WatchPeriod)
		require.Equal(t, []string{"limit"}, resolved.Columns)
	})

	t.Run("keeps cli values when flags are explicitly set", func(t *testing.T) {
		cfg := commonConfig{
			Output:      string(output.Table),
			TableView:   string(tableview.Expanded),
			Alert:       "none",
			WatchPeriod: defaultWatchPeriodSeconds,
			Columns:     []string{"used"},
			fileConfig: &config.Config{Common: config.Common{
				Output:      string(output.JSON),
				TableView:   string(tableview.Compact),
				Alert:       "cpu",
				WatchPeriod: 12,
				Columns:     []string{"limit"},
			}},
		}

		resolved := resolveCommonConfig(cfg, actionFlags{
			outputSet:      true,
			tableViewSet:   true,
			alertSet:       true,
			watchPeriodSet: true,
			columnsSet:     true,
		})

		require.Equal(t, string(output.Table), resolved.Output)
		require.Equal(t, string(tableview.Expanded), resolved.TableView)
		require.Equal(t, "none", resolved.Alert)
		require.Equal(t, uint(defaultWatchPeriodSeconds), resolved.WatchPeriod)
		require.Equal(t, []string{"used"}, resolved.Columns)
	})

	t.Run("defaults table view to compact when unset", func(t *testing.T) {
		resolved := resolveCommonConfig(commonConfig{}, actionFlags{})
		require.Equal(t, string(tableview.Compact), resolved.TableView)
	})

	t.Run("cli columns imply expanded when table view is not explicitly set", func(t *testing.T) {
		resolved := resolveCommonConfig(commonConfig{Columns: []string{"used"}}, actionFlags{columnsSet: true})
		require.Equal(t, string(tableview.Expanded), resolved.TableView)
	})

	t.Run("file columns imply expanded when table view is not explicitly set", func(t *testing.T) {
		resolved := resolveCommonConfig(commonConfig{
			fileConfig: &config.Config{Common: config.Common{Columns: []string{"used"}}},
		}, actionFlags{})
		require.Equal(t, string(tableview.Expanded), resolved.TableView)
	})

	t.Run("explicit compact table view is preserved when columns are set", func(t *testing.T) {
		resolved := resolveCommonConfig(commonConfig{
			TableView: string(tableview.Compact),
			Columns:   []string{"used"},
		}, actionFlags{tableViewSet: true, columnsSet: true})
		require.Equal(t, string(tableview.Compact), resolved.TableView)
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

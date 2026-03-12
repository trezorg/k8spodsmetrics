package stdin

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trezorg/k8spodsmetrics/internal/config"
	"github.com/urfave/cli/v2"
)

func TestCommonConfigValidate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := commonConfig{
			Output:      "table",
			Alert:       "none",
			WatchPeriod: 5,
		}

		require.NoError(t, cfg.Validate())
	})

	t.Run("zero watch period when watch enabled", func(t *testing.T) {
		cfg := commonConfig{
			Output:       "table",
			Alert:        "none",
			WatchMetrics: true,
			WatchPeriod:  0,
		}

		require.ErrorContains(t, cfg.Validate(), "watch period must be greater than 0")
	})

	t.Run("zero watch period when watch disabled", func(t *testing.T) {
		cfg := commonConfig{
			Output:       "table",
			Alert:        "none",
			WatchMetrics: false,
			WatchPeriod:  0,
		}

		require.NoError(t, cfg.Validate())
	})

	t.Run("invalid output", func(t *testing.T) {
		cfg := commonConfig{
			Output:      "invalid",
			Alert:       "none",
			WatchPeriod: 5,
		}

		require.ErrorContains(t, cfg.Validate(), "output should be one of")
	})

	t.Run("invalid alert", func(t *testing.T) {
		cfg := commonConfig{
			Output:      "table",
			Alert:       "invalid",
			WatchPeriod: 5,
		}

		require.ErrorContains(t, cfg.Validate(), "alert should be one of")
	})
}

func TestPodConfigValidate(t *testing.T) {
	t.Run("invalid sorting", func(t *testing.T) {
		cfg := podConfig{
			Sorting: "invalid",
			Resources: []string{
				"all",
			},
			commonConfig: commonConfig{
				Output:      "table",
				Alert:       "none",
				WatchPeriod: 5,
			},
		}

		require.ErrorContains(t, cfg.Validate(), "sorting should be one of")
	})

	t.Run("invalid resources", func(t *testing.T) {
		cfg := podConfig{
			Sorting: "namespace",
			Resources: []string{
				"invalid",
			},
			commonConfig: commonConfig{
				Output:      "table",
				Alert:       "none",
				WatchPeriod: 5,
			},
		}

		require.ErrorContains(t, cfg.Validate(), "invalid resource")
	})
}

func TestSummaryConfigValidate(t *testing.T) {
	t.Run("invalid sorting", func(t *testing.T) {
		cfg := summaryConfig{
			Sorting: "invalid",
			Resources: []string{
				"all",
			},
			commonConfig: commonConfig{
				Output:      "table",
				Alert:       "none",
				WatchPeriod: 5,
			},
		}

		require.ErrorContains(t, cfg.Validate(), "sorting should be one of")
	})

	t.Run("invalid resources", func(t *testing.T) {
		cfg := summaryConfig{
			Sorting: "name",
			Resources: []string{
				"invalid",
			},
			commonConfig: commonConfig{
				Output:      "table",
				Alert:       "none",
				WatchPeriod: 5,
			},
		}

		require.ErrorContains(t, cfg.Validate(), "invalid resource")
	})
}

func TestValidateRejectsInvalidMergedFileConfig(t *testing.T) {
	base := commonConfig{
		Output:      "",
		Alert:       "none",
		WatchPeriod: 5,
	}

	base.fileConfig = &config.Config{
		Common: config.Common{Output: "invalid"},
		Summary: config.Summary{
			Sorting:   "name",
			Resources: []string{"all"},
		},
	}

	resolved := resolveSummaryActionConfig(newSummaryTestContext(t), base)

	require.ErrorContains(t, resolved.Validate(), "output should be one of")
}

func newSummaryTestContext(t *testing.T, args ...string) *cli.Context {
	t.Helper()

	set := flag.NewFlagSet("summary", flag.ContinueOnError)
	cfg := commonConfig{}
	for _, cliFlag := range commonFlags(&cfg) {
		require.NoError(t, cliFlag.Apply(set))
	}
	for _, cliFlag := range summaryFlags() {
		require.NoError(t, cliFlag.Apply(set))
	}
	require.NoError(t, set.Parse(args))

	return cli.NewContext(cli.NewApp(), set, nil)
}

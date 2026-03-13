package stdin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAppRunConfigRegression(t *testing.T) {
	t.Setenv("KUBECONFIG", filepath.Join(t.TempDir(), "missing-kubeconfig-from-env"))

	t.Run("invalid config path surfaces through app", func(t *testing.T) {
		missingConfigPath := filepath.Join(t.TempDir(), "missing.yaml")

		err := runApp(t, "--config", missingConfigPath, "summary")

		require.ErrorContains(t, err, "failed to read config file")
	})

	t.Run("invalid config content surfaces through app", func(t *testing.T) {
		configPath := writeConfigFile(t, "common:\n  output: [\n")

		err := runApp(t, "--config", configPath, "summary")

		require.ErrorContains(t, err, "failed to parse config file")
	})

	t.Run("file output is used when output flag is omitted", func(t *testing.T) {
		configPath := writeConfigFile(t, "common:\n  output: invalid\nsummary:\n  sorting: name\n")

		err := runApp(t, "--config", configPath, "summary")

		require.ErrorContains(t, err, "output should be one of")
	})

	t.Run("invalid table view surfaces through app", func(t *testing.T) {
		configPath := writeConfigFile(t, "common:\n  table-view: invalid\nsummary:\n  sorting: name\n")

		err := runApp(t, "--config", configPath, "summary")

		require.ErrorContains(t, err, "table view should be one of")
	})

	t.Run("file table view is used when flag is omitted", func(t *testing.T) {
		configPath := writeConfigFile(t, "common:\n  output: table\n  table-view: invalid\nsummary:\n  sorting: name\n")

		err := runApp(t, "--config", configPath, "summary")

		require.ErrorContains(t, err, "table view should be one of")
	})

	t.Run("file sorting is used when sorting flag is omitted", func(t *testing.T) {
		configPath := writeConfigFile(t, "summary:\n  sorting: invalid\n")

		err := runApp(t, "--config", configPath, "summary")

		require.ErrorContains(t, err, "sorting should be one of")
	})

	t.Run("file resources are used when resources flag is omitted", func(t *testing.T) {
		configPath := writeConfigFile(t, "summary:\n  resources:\n    - invalid\n")

		err := runApp(t, "--config", configPath, "summary")

		require.ErrorContains(t, err, "invalid resource")
	})

	t.Run("file columns are used when columns flag is omitted", func(t *testing.T) {
		configPath := writeConfigFile(t, "common:\n  output: table\n  columns:\n    - invalid\nsummary:\n  sorting: name\n")

		err := runApp(t, "--config", configPath, "summary")

		require.ErrorContains(t, err, "invalid column")
	})

	t.Run("explicit CLI values override invalid file sorting and output", func(t *testing.T) {
		configPath := writeConfigFile(t, "common:\n  output: invalid\nsummary:\n  sorting: invalid\n")
		missingKubeconfigPath := filepath.Join(t.TempDir(), "missing-kubeconfig")

		err := runApp(
			t,
			"--config", configPath,
			"--output", "table",
			"--kubeconfig", missingKubeconfigPath,
			"summary",
			"--sorting", "name",
		)

		require.Error(t, err)
		require.ErrorContains(t, err, missingKubeconfigPath)
		require.NotContains(t, err.Error(), "output should be one of")
		require.NotContains(t, err.Error(), "sorting should be one of")
	})

	t.Run("explicit table view overrides invalid file table view", func(t *testing.T) {
		configPath := writeConfigFile(t, "common:\n  output: table\n  table-view: invalid\nsummary:\n  sorting: name\n")
		missingKubeconfigPath := filepath.Join(t.TempDir(), "missing-kubeconfig")

		err := runApp(
			t,
			"--config", configPath,
			"--output", "table",
			"--table-view", "compact",
			"--kubeconfig", missingKubeconfigPath,
			"summary",
			"--sorting", "name",
		)

		require.Error(t, err)
		require.ErrorContains(t, err, missingKubeconfigPath)
		require.NotContains(t, err.Error(), "table view should be one of")
	})

	t.Run("explicit watch false beats file watch true and reaches kubeconfig error", func(t *testing.T) {
		configPath := writeConfigFile(t, "common:\n  watch: true\n")
		missingKubeconfigPath := filepath.Join(t.TempDir(), "missing-kubeconfig")

		err := runApp(
			t,
			"--config", configPath,
			"--kubeconfig", missingKubeconfigPath,
			"--watch-period", "0",
			"--watch=false",
			"summary",
		)

		require.Error(t, err)
		require.ErrorContains(t, err, missingKubeconfigPath)
		require.NotContains(t, err.Error(), "watch period must be greater than 0")
	})

	t.Run("columns are rejected for compact table view", func(t *testing.T) {
		err := runApp(
			t,
			"--output", "table",
			"--table-view", "compact",
			"--columns", "used",
			"summary",
		)

		require.ErrorContains(t, err, "--columns is only supported with --table-view expanded")
	})

	t.Run("columns imply expanded when table view is not explicitly set", func(t *testing.T) {
		missingKubeconfigPath := filepath.Join(t.TempDir(), "missing-kubeconfig")
		err := runApp(
			t,
			"--output", "table",
			"--columns", "used",
			"--kubeconfig", missingKubeconfigPath,
			"summary",
		)

		require.Error(t, err)
		require.ErrorContains(t, err, missingKubeconfigPath)
		require.NotContains(t, err.Error(), "--columns is only supported with --table-view expanded")
	})
}

func runApp(t *testing.T, args ...string) error {
	t.Helper()

	app := NewApp("test")
	return app.Run(append([]string{"k8spodsmetrics"}, args...))
}

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(content), 0o644))
	return configPath
}

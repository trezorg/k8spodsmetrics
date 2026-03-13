package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestReadmeConfigExample(t *testing.T) {
	readmeConfig := extractReadmeConfigYAML(t)

	t.Run("valid example", func(t *testing.T) {
		cfg := decodeStrictConfig(t, readmeConfig)

		require.Equal(t, "compact", cfg.Common.TableView)
		require.Equal(t, []string{"request", "limit", "used"}, cfg.Common.Columns)
		require.True(t, cfg.Common.WatchMetrics)
		require.Equal(t, uint(45), cfg.Common.Timeout)

		require.Equal(t, StringOrSlice{"default"}, cfg.Pods.Namespaces)
		require.Equal(t, []string{"node1", "node2"}, cfg.Pods.Nodes)
		require.Equal(t, []string{"cpu", "memory"}, cfg.Pods.Resources)

		require.Equal(t, "used_cpu", cfg.Summary.Sorting)
		require.Equal(t, []string{"all"}, cfg.Summary.Resources)
	})

	t.Run("invalid key", func(t *testing.T) {
		invalidConfig := strings.Replace(readmeConfig, "common:\n", "common:\n  unsupported-key: true\n", 1)

		_, err := strictDecodeConfig(invalidConfig)
		require.Error(t, err)
		require.ErrorContains(t, err, "unsupported-key")
	})
}

func extractReadmeConfigYAML(t *testing.T) string {
	t.Helper()

	readmePath := filepath.Join("..", "..", "README.md")
	readmeBytes, err := os.ReadFile(readmePath)
	require.NoError(t, err)

	readmeContent := string(readmeBytes)
	anchorIndex := strings.Index(readmeContent, "Example configuration file:")
	require.NotEqual(t, -1, anchorIndex, "README config example anchor not found")

	configSection := readmeContent[anchorIndex:]
	blockStart := strings.Index(configSection, "```yaml")
	require.NotEqual(t, -1, blockStart, "README YAML block start not found")
	blockStart += len("```yaml")

	blockEnd := strings.Index(configSection[blockStart:], "```")
	require.NotEqual(t, -1, blockEnd, "README YAML block end not found")

	return strings.TrimSpace(configSection[blockStart : blockStart+blockEnd])
}

func decodeStrictConfig(t *testing.T, src string) Config {
	t.Helper()

	cfg, err := strictDecodeConfig(src)
	require.NoError(t, err)

	return cfg
}

func strictDecodeConfig(src string) (Config, error) {
	decoder := yaml.NewDecoder(strings.NewReader(src))
	decoder.KnownFields(true)

	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("strict decode config: %w", err)
	}

	return cfg, nil
}

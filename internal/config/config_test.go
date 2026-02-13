package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("loads valid config", func(t *testing.T) {
		yamlContent := `
common:
  kubeconfig: /path/to/kubeconfig
  context: my-context
  output: json
  alert: cpu
  kloglevel: 2
  watch-period: 10
  watch: true
pods:
  namespace: default
  label: app=nginx
  field-selector: status.phase=Running
  nodes:
    - node1
    - node2
  sorting: name
  reverse: true
  resources:
    - cpu
    - memory
summary:
  name: node-name
  label: kubernetes.io/role=master
  sorting: used_cpu
  reverse: false
  resources:
    - all
`
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")
		err := os.WriteFile(configPath, []byte(yamlContent), 0644)
		require.NoError(t, err)

		cfg, err := Load(configPath)
		require.NoError(t, err)
		require.Equal(t, "/path/to/kubeconfig", cfg.Common.KubeConfig)
		require.Equal(t, "my-context", cfg.Common.KubeContext)
		require.Equal(t, "json", cfg.Common.Output)
		require.Equal(t, "cpu", cfg.Common.Alert)
		require.Equal(t, uint(2), cfg.Common.KLogLevel)
		require.Equal(t, uint(10), cfg.Common.WatchPeriod)
		require.True(t, cfg.Common.WatchMetrics)

		require.Equal(t, StringOrSlice{"default"}, cfg.Pods.Namespaces)
		require.Equal(t, "app=nginx", cfg.Pods.Label)
		require.Equal(t, "status.phase=Running", cfg.Pods.FieldSelector)
		require.Equal(t, []string{"node1", "node2"}, cfg.Pods.Nodes)
		require.Equal(t, "name", cfg.Pods.Sorting)
		require.True(t, cfg.Pods.Reverse)
		require.Equal(t, []string{"cpu", "memory"}, cfg.Pods.Resources)

		require.Equal(t, "node-name", cfg.Summary.Name)
		require.Equal(t, "kubernetes.io/role=master", cfg.Summary.Label)
		require.Equal(t, "used_cpu", cfg.Summary.Sorting)
		require.False(t, cfg.Summary.Reverse)
		require.Equal(t, []string{"all"}, cfg.Summary.Resources)
	})

	t.Run("loads config with multiple namespaces", func(t *testing.T) {
		yamlContent := `
pods:
  namespace:
    - ns1
    - ns2
    - ns3
`
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")
		err := os.WriteFile(configPath, []byte(yamlContent), 0644)
		require.NoError(t, err)

		cfg, err := Load(configPath)
		require.NoError(t, err)
		require.Equal(t, StringOrSlice{"ns1", "ns2", "ns3"}, cfg.Pods.Namespaces)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := Load("/nonexistent/path/config.yaml")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to read config file")
	})

	t.Run("invalid yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")
		err := os.WriteFile(configPath, []byte("invalid: yaml: content:\n  - ["), 0644)
		require.NoError(t, err)

		_, err = Load(configPath)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse config file")
	})

	t.Run("empty config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")
		err := os.WriteFile(configPath, []byte("{}"), 0644)
		require.NoError(t, err)

		cfg, err := Load(configPath)
		require.NoError(t, err)
		require.Equal(t, "", cfg.Common.KubeConfig)
		require.Empty(t, cfg.Pods.Namespaces)
	})
}

func TestMergeCommon(t *testing.T) {
	t.Run("merges empty values from file", func(t *testing.T) {
		fileConfig := &Config{
			Common: Common{
				KubeConfig:   "/path/to/kubeconfig",
				KubeContext:  "my-context",
				Output:       "json",
				Alert:        "cpu",
				KLogLevel:    2,
				WatchPeriod:  10,
				WatchMetrics: true,
			},
		}
		common := &Common{}

		fileConfig.MergeCommon(common)
		require.Equal(t, "/path/to/kubeconfig", common.KubeConfig)
		require.Equal(t, "my-context", common.KubeContext)
		require.Equal(t, "json", common.Output)
		require.Equal(t, "cpu", common.Alert)
		require.Equal(t, uint(2), common.KLogLevel)
		require.Equal(t, uint(10), common.WatchPeriod)
		require.True(t, common.WatchMetrics)
	})

	t.Run("cli string and numeric values take precedence", func(t *testing.T) {
		fileConfig := &Config{
			Common: Common{
				KubeConfig:  "/file/kubeconfig",
				KubeContext: "file-context",
				Output:      "yaml",
				Alert:       "memory",
				KLogLevel:   1,
				WatchPeriod: 5,
			},
		}
		common := &Common{
			KubeConfig:  "/cli/kubeconfig",
			KubeContext: "cli-context",
			Output:      "json",
			Alert:       "cpu",
			KLogLevel:   3,
			WatchPeriod: 15,
		}

		fileConfig.MergeCommon(common)
		require.Equal(t, "/cli/kubeconfig", common.KubeConfig)
		require.Equal(t, "cli-context", common.KubeContext)
		require.Equal(t, "json", common.Output)
		require.Equal(t, "cpu", common.Alert)
		require.Equal(t, uint(3), common.KLogLevel)
		require.Equal(t, uint(15), common.WatchPeriod)
	})

	// Note: Boolean fields have a limitation - CLI default false cannot override file's true.
	// This is by design since CLI boolean flags cannot distinguish "not set" from "explicitly false".
	// If file has watch: true, the merged value will be true even if CLI doesn't pass --watch.
	t.Run("file boolean true overrides cli default false", func(t *testing.T) {
		fileConfig := &Config{
			Common: Common{
				WatchMetrics: true,
			},
		}
		common := &Common{
			WatchMetrics: false, // CLI default
		}

		fileConfig.MergeCommon(common)
		require.True(t, common.WatchMetrics) // File's true overrides CLI's default false
	})

	t.Run("partial merge", func(t *testing.T) {
		fileConfig := &Config{
			Common: Common{
				KubeConfig:  "/file/kubeconfig",
				KubeContext: "file-context",
				Output:      "",
				Alert:       "",
			},
		}
		common := &Common{
			KubeConfig: "",
			Output:     "json",
		}

		fileConfig.MergeCommon(common)
		require.Equal(t, "/file/kubeconfig", common.KubeConfig)
		require.Equal(t, "file-context", common.KubeContext)
		require.Equal(t, "json", common.Output)
	})
}

func TestMergePods(t *testing.T) {
	t.Run("merges empty values from file", func(t *testing.T) {
		fileConfig := &Config{
			Pods: Pods{
				Namespaces:    StringOrSlice{"default"},
				Label:         "app=nginx",
				FieldSelector: "status.phase=Running",
				Nodes:         []string{"node1", "node2"},
				Sorting:       "name",
				Reverse:       true,
				Resources:     []string{"cpu", "memory"},
			},
		}
		pods := &Pods{}

		fileConfig.MergePods(pods)
		require.Equal(t, StringOrSlice{"default"}, pods.Namespaces)
		require.Equal(t, "app=nginx", pods.Label)
		require.Equal(t, "status.phase=Running", pods.FieldSelector)
		require.Equal(t, []string{"node1", "node2"}, pods.Nodes)
		require.Equal(t, "name", pods.Sorting)
		require.True(t, pods.Reverse)
		require.Equal(t, []string{"cpu", "memory"}, pods.Resources)
	})

	t.Run("cli string and slice values take precedence", func(t *testing.T) {
		fileConfig := &Config{
			Pods: Pods{
				Namespaces:    StringOrSlice{"file-ns"},
				Label:         "file=label",
				FieldSelector: "file=selector",
				Nodes:         []string{"file-node"},
				Sorting:       "namespace",
				Resources:     []string{"file-res"},
			},
		}
		pods := &Pods{
			Namespaces:    StringOrSlice{"cli-ns"},
			Label:         "cli=label",
			FieldSelector: "cli=selector",
			Nodes:         []string{"cli-node"},
			Sorting:       "name",
			Resources:     []string{"cli-res"},
		}

		fileConfig.MergePods(pods)
		require.Equal(t, StringOrSlice{"cli-ns"}, pods.Namespaces)
		require.Equal(t, "cli=label", pods.Label)
		require.Equal(t, "cli=selector", pods.FieldSelector)
		require.Equal(t, []string{"cli-node"}, pods.Nodes)
		require.Equal(t, "name", pods.Sorting)
		require.Equal(t, []string{"cli-res"}, pods.Resources)
	})

	// Note: Boolean fields have a limitation - CLI default false cannot override file's true.
	t.Run("file boolean true overrides cli default false", func(t *testing.T) {
		fileConfig := &Config{
			Pods: Pods{
				Reverse: true,
			},
		}
		pods := &Pods{
			Reverse: false, // CLI default
		}

		fileConfig.MergePods(pods)
		require.True(t, pods.Reverse) // File's true overrides CLI's default false
	})
}

func TestMergeSummary(t *testing.T) {
	t.Run("merges empty values from file", func(t *testing.T) {
		fileConfig := &Config{
			Summary: Summary{
				Name:      "node-name",
				Label:     "kubernetes.io/role=master",
				Sorting:   "used_cpu",
				Reverse:   true,
				Resources: []string{"cpu", "memory"},
			},
		}
		summary := &Summary{}

		fileConfig.MergeSummary(summary)
		require.Equal(t, "node-name", summary.Name)
		require.Equal(t, "kubernetes.io/role=master", summary.Label)
		require.Equal(t, "used_cpu", summary.Sorting)
		require.True(t, summary.Reverse)
		require.Equal(t, []string{"cpu", "memory"}, summary.Resources)
	})

	t.Run("cli string and slice values take precedence", func(t *testing.T) {
		fileConfig := &Config{
			Summary: Summary{
				Name:      "file-node",
				Label:     "file=label",
				Sorting:   "name",
				Resources: []string{"file-res"},
			},
		}
		summary := &Summary{
			Name:      "cli-node",
			Label:     "cli=label",
			Sorting:   "used_cpu",
			Resources: []string{"cli-res"},
		}

		fileConfig.MergeSummary(summary)
		require.Equal(t, "cli-node", summary.Name)
		require.Equal(t, "cli=label", summary.Label)
		require.Equal(t, "used_cpu", summary.Sorting)
		require.Equal(t, []string{"cli-res"}, summary.Resources)
	})

	// Note: Boolean fields have a limitation - CLI default false cannot override file's true.
	t.Run("file boolean true overrides cli default false", func(t *testing.T) {
		fileConfig := &Config{
			Summary: Summary{
				Reverse: true,
			},
		}
		summary := &Summary{
			Reverse: false, // CLI default
		}

		fileConfig.MergeSummary(summary)
		require.True(t, summary.Reverse) // File's true overrides CLI's default false
	})
}

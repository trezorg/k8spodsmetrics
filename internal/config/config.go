// Package config provides YAML configuration file support for the CLI.
//
// Configuration File Format (YAML):
//
//	common:
//	  kubeconfig: /path/to/kubeconfig
//	  context: my-context
//	  output: json|yaml|table|string
//	  alert: cpu|memory
//	  watch-period: 10
//	  watch: true
//	pods:
//	  namespace: default          # Single namespace (string)
//	  # OR
//	  namespace:                  # Multiple namespaces (list)
//	    - ns1
//	    - ns2
//	    - ns3
//	  label: app=nginx
//	  field-selector: status.phase=Running
//	  nodes:
//	    - node1
//	    - node2
//	  sorting: name|namespace|cpu|memory
//	  reverse: true
//	  resources:
//	    - cpu
//	    - memory
//	summary:
//	  name: node-name
//	  label: kubernetes.io/role=master
//	  sorting: used_cpu|used_memory|name
//	  reverse: false
//	  resources:
//	    - all
//
// Merge Behavior:
//   - CLI flags take precedence over file config values
//   - Empty/zero values from CLI are replaced with file config values
//   - Boolean limitation: CLI default false cannot override file's true
//     (use --watch=false explicitly if supported by CLI library)
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Common holds shared configuration options applicable to all commands.
type Common struct {
	KubeConfig   string `yaml:"kubeconfig"`
	KubeContext  string `yaml:"context"`
	Output       string `yaml:"output"`
	Alert        string `yaml:"alert"`
	WatchPeriod  uint   `yaml:"watch-period"`
	WatchMetrics bool   `yaml:"watch"`
}

// StringOrSlice is a custom type that can unmarshal from either a string or a slice of strings in YAML.
type StringOrSlice []string

// UnmarshalYAML implements yaml.Unmarshaler to accept either a string or a slice.
func (s *StringOrSlice) UnmarshalYAML(node *yaml.Node) error {
	var single string
	if err := node.Decode(&single); err == nil {
		*s = []string{single}
		return nil
	}

	var multi []string
	if err := node.Decode(&multi); err != nil {
		return fmt.Errorf("expected string or array of strings, got: %s", node.ShortTag())
	}
	*s = multi
	return nil
}

// Pods holds configuration specific to the pods command.
type Pods struct {
	Namespaces    StringOrSlice `yaml:"namespace"`
	Label         string        `yaml:"label"`
	FieldSelector string        `yaml:"field-selector"`
	Nodes         []string      `yaml:"nodes"`
	Sorting       string        `yaml:"sorting"`
	Reverse       bool          `yaml:"reverse"`
	Resources     []string      `yaml:"resources"`
}

// Summary holds configuration specific to the summary command.
type Summary struct {
	Name      string   `yaml:"name"`
	Label     string   `yaml:"label"`
	Sorting   string   `yaml:"sorting"`
	Reverse   bool     `yaml:"reverse"`
	Resources []string `yaml:"resources"`
}

// Config represents the complete configuration file structure.
type Config struct {
	Common  Common  `yaml:"common"`
	Pods    Pods    `yaml:"pods"`
	Summary Summary `yaml:"summary"`
}

// Load reads and parses a YAML configuration file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// MergeCommon merges file config values into the provided Common struct.
// Only empty/zero values in the target are replaced with file config values.
// Note: For boolean WatchMetrics, file's true will override target's false.
func (c *Config) MergeCommon(common *Common) {
	if common.KubeConfig == "" && c.Common.KubeConfig != "" {
		common.KubeConfig = c.Common.KubeConfig
	}
	if common.KubeContext == "" && c.Common.KubeContext != "" {
		common.KubeContext = c.Common.KubeContext
	}
	if common.Output == "" && c.Common.Output != "" {
		common.Output = c.Common.Output
	}
	if common.Alert == "" && c.Common.Alert != "" {
		common.Alert = c.Common.Alert
	}
	if common.WatchPeriod == 0 && c.Common.WatchPeriod != 0 {
		common.WatchPeriod = c.Common.WatchPeriod
	}
	if !common.WatchMetrics && c.Common.WatchMetrics {
		common.WatchMetrics = c.Common.WatchMetrics
	}
}

// MergePods merges file config values into the provided Pods struct.
// Only empty/zero values in the target are replaced with file config values.
// Note: For boolean Reverse, file's true will override target's false.
func (c *Config) MergePods(pods *Pods) {
	if len(pods.Namespaces) == 0 && len(c.Pods.Namespaces) > 0 {
		pods.Namespaces = c.Pods.Namespaces
	}
	if pods.Label == "" && c.Pods.Label != "" {
		pods.Label = c.Pods.Label
	}
	if pods.FieldSelector == "" && c.Pods.FieldSelector != "" {
		pods.FieldSelector = c.Pods.FieldSelector
	}
	if len(pods.Nodes) == 0 && len(c.Pods.Nodes) > 0 {
		pods.Nodes = c.Pods.Nodes
	}
	if pods.Sorting == "" && c.Pods.Sorting != "" {
		pods.Sorting = c.Pods.Sorting
	}
	if !pods.Reverse && c.Pods.Reverse {
		pods.Reverse = c.Pods.Reverse
	}
	if len(pods.Resources) == 0 && len(c.Pods.Resources) > 0 {
		pods.Resources = c.Pods.Resources
	}
}

// MergeSummary merges file config values into the provided Summary struct.
// Only empty/zero values in the target are replaced with file config values.
// Note: For boolean Reverse, file's true will override target's false.
func (c *Config) MergeSummary(summary *Summary) {
	if summary.Name == "" && c.Summary.Name != "" {
		summary.Name = c.Summary.Name
	}
	if summary.Label == "" && c.Summary.Label != "" {
		summary.Label = c.Summary.Label
	}
	if summary.Sorting == "" && c.Summary.Sorting != "" {
		summary.Sorting = c.Summary.Sorting
	}
	if !summary.Reverse && c.Summary.Reverse {
		summary.Reverse = c.Summary.Reverse
	}
	if len(summary.Resources) == 0 && len(c.Summary.Resources) > 0 {
		summary.Resources = c.Summary.Resources
	}
}

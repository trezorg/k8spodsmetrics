package metricsresources

const (
	defaultWatchPeriod = 5
)

type ConfigBuilder struct {
	config *Config
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{
			WatchPeriod: defaultWatchPeriod,
		},
	}
}

func (b *ConfigBuilder) WithKubeConfig(kubeconfig string) *ConfigBuilder {
	b.config.KubeConfig = kubeconfig
	return b
}

func (b *ConfigBuilder) WithKubeContext(context string) *ConfigBuilder {
	b.config.KubeContext = context
	return b
}

func (b *ConfigBuilder) WithNamespaces(namespaces []string) *ConfigBuilder {
	b.config.Namespaces = namespaces
	return b
}

func (b *ConfigBuilder) WithLabel(label string) *ConfigBuilder {
	b.config.Label = label
	return b
}

func (b *ConfigBuilder) WithFieldSelector(selector string) *ConfigBuilder {
	b.config.FieldSelector = selector
	return b
}

func (b *ConfigBuilder) WithNodes(nodes []string) *ConfigBuilder {
	b.config.Nodes = nodes
	return b
}

func (b *ConfigBuilder) WithOutput(output string) *ConfigBuilder {
	b.config.Output = output
	return b
}

func (b *ConfigBuilder) WithSorting(sorting string) *ConfigBuilder {
	b.config.Sorting = sorting
	return b
}

func (b *ConfigBuilder) WithResources(resources []string) *ConfigBuilder {
	b.config.Resources = resources
	return b
}

func (b *ConfigBuilder) WithAlert(alert string) *ConfigBuilder {
	b.config.Alert = alert
	return b
}

func (b *ConfigBuilder) WithWatchPeriod(period uint) *ConfigBuilder {
	b.config.WatchPeriod = period
	return b
}

func (b *ConfigBuilder) WithReverse(reverse bool) *ConfigBuilder {
	b.config.Reverse = reverse
	return b
}

func (b *ConfigBuilder) WithWatchMetrics(watch bool) *ConfigBuilder {
	b.config.WatchMetrics = watch
	return b
}

func (b *ConfigBuilder) Build() Config {
	return *b.config
}

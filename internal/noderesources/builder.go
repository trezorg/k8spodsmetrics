package noderesources

const (
	defaultWatchPeriod = 5
	defaultKLogLevel   = 3
)

type ConfigBuilder struct {
	config *Config
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{
			WatchPeriod: defaultWatchPeriod,
			KLogLevel:   defaultKLogLevel,
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

func (b *ConfigBuilder) WithLabel(label string) *ConfigBuilder {
	b.config.Label = label
	return b
}

func (b *ConfigBuilder) WithName(name string) *ConfigBuilder {
	b.config.Name = name
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

func (b *ConfigBuilder) WithKLogLevel(level uint) *ConfigBuilder {
	b.config.KLogLevel = level
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

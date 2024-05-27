package configuration

type ConfigLoader interface {
	LoadConfig() (*Config, error)
}

package config

type Provider interface {
	Config() *Config
	Parse() error
}

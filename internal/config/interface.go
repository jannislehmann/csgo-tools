package config

// UseCase defines the config service functions.
type UseCase interface {
	GetConfig() *Config
	IsDebug() bool
	IsTrace() bool
}

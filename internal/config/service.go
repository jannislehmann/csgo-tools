package config

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Service struct {
	config *Config
}

func NewService() *Service {
	service := Service{}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AddConfigPath("./configs")
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("csgo")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	service.config = &Config{}
	if err := viper.Unmarshal(&service.config); err != nil {
		log.Panic(err)
	}

	service.setLoggingLevel()

	return &service
}

// Config holds the application configuration.
type Config struct {
	DemosDir string          `mapstructure:"demosDir"`
	Auth     *AuthConfig     `mapstructure:"auth"`
	Steam    *SteamConfig    `mapstructure:"steam"`
	Database *DatabaseConfig `mapstructure:"database"`
	Debug    string          `mapstructure:"debug"`
	Parser   *ParserConfig   `mapstructure:"parser"`
}

// AuthConfig contains the host url for the authentication callback.
type AuthConfig struct {
	Host   string `mapstructure:"host"`
	Secret string `mapstructure:"secret"`
}

// SteamConfig holds the configuration about the steam account to use for communicating with the GameCoordinator.
type SteamConfig struct {
	SteamAPIKey     string `mapstructure:"apiKey"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	TwoFactorSecret string `mapstructure:"twoFactorSecret"`
}

// DatabaseConfig holds database connection information.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type ParserConfig struct {
	WorkerCount string `mapstructure:"workerCount"`
}

// GetConfig returns the application configuration.
func (s *Service) GetConfig() *Config {
	return s.config
}

// IsDebug returns whether the application is in debug mode.
func (s *Service) IsDebug() bool {
	return s.GetConfig().Debug == "true" || s.IsTrace()
}

// IsTrace returns whether the application should do extended debugging.
func (s *Service) IsTrace() bool {
	return s.GetConfig().Debug == "trace"
}

// SetLoggingLevel sets the logging level in relation to the level set in the config file.
func (s *Service) setLoggingLevel() {
	if s.IsTrace() {
		log.SetLevel(log.TraceLevel)
		log.SetReportCaller(true)
	} else if s.IsDebug() {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

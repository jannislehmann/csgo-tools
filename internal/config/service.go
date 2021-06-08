package config

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	config *Config
}

func NewService() *Service {
	service := Service{}

	file := "./configs/config.json"
	configFile, err := os.Open(file)
	if err != nil {
		configFile.Close()
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	service.config = &Config{}

	if err = jsonParser.Decode(service.GetConfig()); err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	service.setLoggingLevel()

	return &service
}

// Config holds the application configuration.
type Config struct {
	DemosDir string          `json:"demosDir"`
	Steam    *SteamConfig    `json:"steam"`
	Database *DatabaseConfig `json:"database"`
	Debug    string          `json:"debug"`
}

// SteamConfig holds the configuration about the steam account to use for communicating with the GameCoordinator.
type SteamConfig struct {
	SteamAPIKey     string `json:"apiKey"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	TwoFactorSecret string `json:"twoFactorSecret"`
}

// DatabaseConfig holds database connection information.
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
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

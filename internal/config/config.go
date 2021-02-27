package config

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

var config Config

// Config holds the application configuration.
type Config struct {
	DemosDir string          `json:"demosDir"`
	Steam    *SteamConfig    `json:"steam"`
	CSGO     []*CSGOConfig   `json:"csgo"`
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

// CSGOConfig holds the accounts to watch.
type CSGOConfig struct {
	HistoryAPIKey  string `json:"matchHistoryAuthenticationCode"`
	KnownMatchCode string `json:"knownMatchCode"`
	SteamID        string `json:"steamId"` // should be uint64
}

// DatabaseConfig holds database connection information.
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func init() {
	file := "./configs/config.json"
	configFile, err := os.Open(file)
	if err != nil {
		configFile.Close()
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	config = Config{}
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()
}

// GetConfiguration returns the application configuration.
func GetConfiguration() *Config {
	return &config
}

// IsDebug returns whether the application is in debug mode.
func (c *Config) IsDebug() bool {
	return config.Debug == "true"
}

// IsTrace returns whether the application should do extended debugging.
func (c *Config) IsTrace() bool {
	return config.Debug == "trace"
}

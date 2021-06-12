package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/match"
	"github.com/Cludch/csgo-tools/internal/gamecoordinator"
	"github.com/Cludch/csgo-tools/internal/steam_client"
	"github.com/Philipp15b/go-steam/v2"
)

var configService *config.Service
var matchService *match.Service
var steamService *steam_client.Service
var gamecoordinatorService *gamecoordinator.Service

func setup() {
	err := steam.InitializeSteamDirectory()
	if err != nil {
		log.Error(err)
	}

	configService = config.NewService()
	db := entity.NewService(configService)

	matchService = match.NewService(match.NewRepositoryMongo(db))
	gamecoordinatorService = gamecoordinator.NewService(matchService)
	steamService = steam_client.NewService(gamecoordinatorService)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
}

func main() {
	setup()

	configData := configService.GetConfig()
	steamService.Connect(configData.Steam.Username, configData.Steam.Password, configData.Steam.TwoFactorSecret)
}

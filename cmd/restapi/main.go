package main

import (
	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/match"
	"github.com/Cludch/csgo-tools/internal/domain/player"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var configService *config.Service
var matchService *match.Service
var playerService *player.Service

// Sets up the global variables (config, db) and the logger.
func setup() {
	configService = config.NewService()
	db := entity.NewService(configService)
	matchService = match.NewService(match.NewRepositoryMongo(db))
	playerService = player.NewService(player.NewRepositoryMongo(db))

	if !configService.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
}

func main() {
	setup()

	log.Info("starting rest api")

	router := gin.Default()
	matchController := match.NewController(matchService)
	playerController := player.NewController(playerService)

	router.GET("/match", matchController.GetMatches)
	router.GET("/match/:id", matchController.GetMatchDetails)
	router.GET("/player/", playerController.GetPlayers)
	router.GET("/player/:id", playerController.GetPlayerDetails)
	router.GET("/player/:id/stats", playerController.GetPlayerAverageStats)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}

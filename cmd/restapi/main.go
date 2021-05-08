package main

import (
	"errors"
	"net/http"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/demoparser"
	"github.com/Cludch/csgo-tools/internal/entity"
	"github.com/Cludch/csgo-tools/pkg/restapi"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var configData *config.Config
var db *gorm.DB

// Sets up the global variables (config, db) and the logger.
func init() {
	db = entity.GetDatabase()
	configData = config.GetConfiguration()
	demoparser.ConfigData = configData
	demoparser.DB = db

	if !configData.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}

	configData.SetLoggingLevel()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
}

func main() {
	log.Info("starting rest api")

	router := gin.Default()

	router.GET("/match", getMatches)
	router.GET("/match/:id", getMatchDetails)
	router.GET("/player/:id", getPlayerStats)
	router.GET("/player/:id/stats", getAveragePlayerStats)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}

// getMatches returns all matches from the database containing the metadata.
func getMatches(c *gin.Context) {
	var matches []demoparser.MatchResult
	db.Find(&matches)
	c.JSON(http.StatusOK, matches)
}

func getMatchDetails(c *gin.Context) {
	var match demoparser.MatchResult

	// Check whether match exists and load it and the teams.
	err := db.Preload("Teams").Preload("Teams.Players").Select("match_results.match_id", "map", "time", "duration").Where("match_results.match_id = ?", c.Param("id")).First(&match).Error
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, match)
}

func getPlayerStats(c *gin.Context) {
	var playerResults []demoparser.PlayerResult

	if err := db.Where("steam_id = ?", c.Param("id")).Find(&playerResults).Error; err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, playerResults)
}

// Takes all player results and creates an average from those stats.
func getAveragePlayerStats(c *gin.Context) {
	var playerResults []demoparser.PlayerResult

	if err := db.Where("steam_id = ?", c.Param("id")).Find(&playerResults).Error; err != nil {
		handleError(c, err)
		return
	}

	playerStats := &restapi.PlayerGameStats{}

	assists, kills, entryKills, headshots, deaths, mvps := 0, 0, 0, 0, 0, 0

	for _, playerResult := range playerResults {
		playerStats.Games++
		playerStats.SteamID = playerResult.SteamID

		assists += int(playerResult.Assists)
		kills += int(playerResult.Kills)
		entryKills += int(playerResult.EntryKills)
		headshots += int(playerResult.Headshots)
		deaths += int(playerResult.Deaths)
		mvps += int(playerResult.MVPs)

		playerStats.Won1v3 += int(playerResult.Won1v3)
		playerStats.Won1v4 += int(playerResult.Won1v4)
		playerStats.Won1v5 += int(playerResult.Won1v5)
		playerStats.RoundsWith3K += int(playerResult.RoundsWith3K)
		playerStats.RoundsWith4K += int(playerResult.RoundsWith4K)
		playerStats.RoundsWith5K += int(playerResult.RoundsWith5K)
	}

	playerStats.AssistsPerGame += assists / playerStats.Games
	playerStats.KillsPerGame += kills / playerStats.Games
	playerStats.EntryKillsPerGame += entryKills / playerStats.Games
	playerStats.HeadshotsPerGame += headshots / playerStats.Games
	playerStats.DeathsPerGame += deaths / playerStats.Games
	playerStats.MVPsPerGame += mvps / playerStats.Games

	c.JSON(http.StatusOK, playerStats)
}

func handleError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found!"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
	}
}

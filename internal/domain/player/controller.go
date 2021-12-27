package player

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service UseCase
}

func NewController(s UseCase) *Controller {
	return &Controller{
		service: s,
	}
}

func (c *Controller) GetPlayers(g *gin.Context) {
	players, _ := c.service.GetAll()
	playerList := &PlayerList{Players: make([]*PlayerListEntry, len(players))}

	for i, player := range players {
		results := player.Results
		lenResults := len(results)

		lastPlayerName := results[lenResults-1]

		playerList.Players[i] = &PlayerListEntry{
			ID: player.ID, Games: lenResults, Name: lastPlayerName.Name,
		}

	}

	g.JSON(http.StatusOK, playerList)
}

func (c *Controller) GetPlayerDetails(g *gin.Context) {
	id, _ := strconv.ParseUint(g.Param("id"), 10, 64)
	player, _ := c.service.GetPlayer(id)
	g.JSON(http.StatusOK, player)
}

func (c *Controller) GetPlayerAverageStats(g *gin.Context) {
	id, _ := strconv.ParseUint(g.Param("id"), 10, 64)
	player, _ := c.service.GetPlayer(id)

	playerStats := &PlayerGameStats{}

	for _, playerResult := range player.Results {
		playerStats.Games++
		playerStats.SteamID = playerResult.SteamID
		playerStats.Name = playerResult.Name

		playerStats.AssistsPerRound += int(playerResult.Assists) / int(playerResult.MatchRounds)
		playerStats.KillsPerRound += int(playerResult.Kills) / int(playerResult.MatchRounds)
		playerStats.EntryKillsPerRound += int(playerResult.EntryKills) / int(playerResult.MatchRounds)
		playerStats.OpeningDuelAttempsPerRound += int(playerResult.OpeningDuelAttempts) / int(playerResult.MatchRounds)
		playerStats.HeadshotsPerRound += int(playerResult.Headshots) / int(playerResult.MatchRounds)
		playerStats.DeathsPerRound += int(playerResult.Deaths) / int(playerResult.MatchRounds)
		playerStats.MVPsPerRound += int(playerResult.MVPs) / int(playerResult.MatchRounds)
		playerStats.DamagePerRound += playerResult.DamageDealt / int(playerResult.MatchRounds)

		playerStats.Won1v3 += int(playerResult.Won1v3)
		playerStats.Won1v4 += int(playerResult.Won1v4)
		playerStats.Won1v5 += int(playerResult.Won1v5)
		playerStats.RoundsWith3K += int(playerResult.RoundsWith3K)
		playerStats.RoundsWith4K += int(playerResult.RoundsWith4K)
		playerStats.RoundsWith5K += int(playerResult.RoundsWith5K)
	}

	g.JSON(http.StatusOK, playerStats)
}

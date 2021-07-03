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

func (c *Controller) GetPlayer(g *gin.Context) {
	id, _ := strconv.ParseUint(g.Param("id"), 10, 64)
	player, _ := c.service.GetPlayer(id)
	g.JSON(http.StatusOK, player)
}

func (c *Controller) GetPlayerAverageStats(g *gin.Context) {
	id, _ := strconv.ParseUint(g.Param("id"), 10, 64)
	player, _ := c.service.GetPlayer(id)

	playerStats := &PlayerGameStats{}

	assists, kills, entryKills, headshots, deaths, mvps := 0, 0, 0, 0, 0, 0

	for _, playerResult := range player.Results {
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
	g.JSON(http.StatusOK, playerStats)
}

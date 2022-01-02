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

		lastPlayerName := results[lenResults-1].Name
		lastWinCount := results[lenResults-1].WinCount
		lastRank := results[lenResults-1].RankNew

		playerList.Players[i] = &PlayerListEntry{
			ID: player.ID, Games: lenResults, Name: lastPlayerName, Wins: lastWinCount, Rank: lastRank,
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

	// Keep track of the amount of rounds and all assists etc.
	var matchRounds float32
	assists, kills, entryKills, openingDuelAttempts, headshots, deaths, mvps, damageDealt := 0, 0, 0, 0, 0, 0, 0, 0

	for _, playerResult := range player.Results {
		playerStats.Games++
		playerStats.SteamID = playerResult.SteamID
		playerStats.Name = playerResult.Name

		matchRounds += float32(playerResult.MatchRounds)

		assists += int(playerResult.Assists)
		kills += int(playerResult.Kills)
		entryKills += int(playerResult.EntryKills)
		openingDuelAttempts += int(playerResult.OpeningDuelAttempts)
		headshots += int(playerResult.Headshots)
		deaths += int(playerResult.Deaths)
		mvps += int(playerResult.MVPs)
		damageDealt += int(playerResult.DamageDealt)

		playerStats.Won1v3 += int(playerResult.Won1v3)
		playerStats.Won1v4 += int(playerResult.Won1v4)
		playerStats.Won1v5 += int(playerResult.Won1v5)
		playerStats.RoundsWith3K += int(playerResult.RoundsWith3K)
		playerStats.RoundsWith4K += int(playerResult.RoundsWith4K)
		playerStats.RoundsWith5K += int(playerResult.RoundsWith5K)
	}

	playerStats.AssistsPerRound += float32(assists) / matchRounds
	playerStats.KillsPerRound += float32(kills) / matchRounds
	playerStats.EntryKillsPerRound += float32(entryKills) / matchRounds
	playerStats.OpeningDuelAttempsPerRound += float32(openingDuelAttempts) / matchRounds
	playerStats.HeadshotsPerRound += float32(headshots) / matchRounds
	playerStats.DeathsPerRound += float32(deaths) / matchRounds
	playerStats.MVPsPerRound += float32(mvps) / matchRounds
	playerStats.DamagePerRound += float32(damageDealt) / matchRounds

	g.JSON(http.StatusOK, playerStats)
}

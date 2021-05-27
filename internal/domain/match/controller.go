package match

import (
	"net/http"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/restapi"
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

// GetMatches returns all matches from the database.
func (c *Controller) GetMatches(g *gin.Context) {
	matches, _ := c.service.GetAll()
	g.JSON(http.StatusOK, matches)
}

func (c *Controller) GetMatchDetails(g *gin.Context) {
	id, err := entity.StringToID(g.Param("id"))
	if err != nil {
		handleError(g, err)
	}

	match, _ := c.service.GetMatch(id)
	g.JSON(http.StatusOK, match)
}

func (c *Controller) GetPlayer(g *gin.Context) {
	// TODO: Get player(?)
	// matches, _ := c.service.GetValveMatchesMissingDownloadUrl()
	// g.JSON(http.StatusOK, matches)
}

func (c *Controller) GetPlayerAverageStats(g *gin.Context) {
	// TODO: Fill playerResults
	// id := g.Param("id")
	var playerResults []PlayerResult

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
	g.JSON(http.StatusOK, playerStats)
}

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
}

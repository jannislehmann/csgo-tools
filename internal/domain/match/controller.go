package match

import (
	"net/http"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
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
	matches, _ := c.service.GetAllParsed()
	matchList := &MatchList{Matches: make([]*MatchListEntry, len(matches))}

	clanPlayerIds := getClanPlayersIds()

	for i, match := range matches {
		clanTeam := common.TeamUnassigned

		// Search for the team with a member of the clan
		for _, team := range match.Result.Teams {
			for _, teamPlayer := range team.Players {
				if _, ok := clanPlayerIds[teamPlayer.SteamID]; ok {
					clanTeam = team.TeamID
				}
			}
		}

		matchList.Matches[i] = &MatchListEntry{
			ID: match.ID, Time: match.Time, Map: match.Result.Map,
			TeamOneScore: match.Result.Teams[0].Wins,
			TeamTwoScore: match.Result.Teams[1].Wins,
			ClanTeam:     clanTeam,
		}
	}

	g.JSON(http.StatusOK, matchList)
}

func (c *Controller) GetMatchDetails(g *gin.Context) {
	id, err := entity.StringToID(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
	}

	match, _ := c.service.GetMatch(id)
	g.JSON(http.StatusOK, match)
}

// Used to determine whether we (the clan) have won the game.
// This makes the api highly subjective and is planned to get changed eventually.
func getClanPlayersIds() map[uint64]bool {
	return map[uint64]bool{
		76561198185324675: true,
		76561198079819126: true,
		76561198070498642: true,
		76561198075069967: true,
		76561198053633135: true,
	}
}

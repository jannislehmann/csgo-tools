package match

import (
	"fmt"
	"net/http"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
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
	matches, _ := c.service.GetAllParsed()
	matchList := &MatchList{Matches: make([]*MatchListEntry, len(matches))}
	resultTemplate := "%d - %d"

	for i, match := range matches {
		matchList.Matches[i] = &MatchListEntry{
			ID: match.ID, Time: match.Time, Map: match.Result.Map,
			Result: fmt.Sprintf(resultTemplate, match.Result.Teams[0].Wins, match.Result.Teams[1].Wins)}
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

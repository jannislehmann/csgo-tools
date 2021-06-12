package match

import (
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
	matches, _ := c.service.GetAll()
	g.JSON(http.StatusOK, matches)
}

func (c *Controller) GetMatchDetails(g *gin.Context) {
	id, err := entity.StringToID(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
	}

	match, _ := c.service.GetMatch(id)
	g.JSON(http.StatusOK, match)
}

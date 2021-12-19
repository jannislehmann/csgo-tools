package auth

import (
	"fmt"
	"net/http"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/user"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"

	log "github.com/sirupsen/logrus"
)

type Controller struct {
	service     UseCase
	userService user.UseCase
}

func NewController(s UseCase, u user.UseCase) *Controller {
	return &Controller{
		service:     s,
		userService: u,
	}
}

func (c *Controller) Auth(g *gin.Context) {
	q := g.Request.URL.Query()
	q.Add("provider", "steam")
	g.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(g.Writer, g.Request)
}

func (c *Controller) Callback(g *gin.Context) {
	q := g.Request.URL.Query()
	q.Add("provider", "steam")
	g.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(g.Writer, g.Request)
	if err != nil {
		_ = g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	token, err := c.service.HandleAuth(user)

	if err != nil {
		_ = g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	const msg = "auth: user with id %s signed in"
	log.Debugf(msg, user.UserID)

	json := gin.H{
		"type":  "Bearer",
		"token": token,
	}
	g.JSON(http.StatusOK, json)
}

func (c *Controller) AuthorizeRequest(g *gin.Context) {
	const BEARER_SCHEMA = "Bearer "
	authHeader := g.GetHeader("Authorization")
	tokenString := authHeader[len(BEARER_SCHEMA):]
	claims, err := c.service.ValidateToken(tokenString)

	if err == nil {
		g.Set("userId", claims.Id)
		g.Next()
	} else {
		log.Error(err)
		g.AbortWithStatus(http.StatusUnauthorized)
	}
}

func (c *Controller) GetUserDetails(g *gin.Context) {
	userId, exists := g.Get("userId")
	if !exists {
		log.Errorf("unable to extract userId from authentication middleware for %v", userId)
		g.AbortWithStatus(http.StatusUnauthorized)
	}

	parsedUserId, _ := entity.StringToID(fmt.Sprint(userId))
	user, err := c.userService.GetUser(parsedUserId)
	if err != nil {
		log.Errorf("error while fetching user details %e", err)
		g.AbortWithStatus(http.StatusInternalServerError)
	}

	g.JSON(http.StatusOK, user)
}

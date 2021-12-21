package main

import (
	"fmt"
	"net/http"

	"github.com/Cludch/csgo-tools/internal/auth"
	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/user"
	redistore "github.com/boj/redistore"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/steam"
	log "github.com/sirupsen/logrus"
)

var configService *config.Service
var authService *auth.Service
var userService *user.Service

// Sets up the global variables (config, db) and the logger.
func setup() {
	configService = config.NewService()
	db := entity.NewService(configService)

	userService = user.NewService(user.NewRepositoryMongo(db), configService)
	authService = auth.NewService(configService, userService)

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

	log.Info("starting auth service")

	router := gin.Default()

	store, err := redistore.NewRediStore(10, "tcp", ":6379", "csgo_auth_session", []byte("session_secret"))
	// store.Options https://github.com/markbates/goth#security-notes
	if err != nil {
		panic(err)
	}
	defer store.Close()
	gothic.Store = store

	// Register Steam as Goth OpenID 2.0 provider
	goth.UseProviders(
		steam.New(configService.GetConfig().Steam.SteamAPIKey, configService.GetConfig().Auth.Host+"/auth/steam/callback"),
	)

	authController := auth.NewController(authService, userService, store)

	router.GET("/auth/:provider", authController.Auth)
	router.GET("/auth/:provider/callback", authController.Callback)

	// TODO: Middleware
	// Protected API endpoints.
	/*
		authorized := router.Group("/")
		authorized.Use(authController.AuthorizeRequest)
		{
			authorized.GET("/me", authController.GetUserDetails)
		}
	*/

	router.GET("/me", func(g *gin.Context) {
		// TODO: The data in the session is garbage json
		session, err := gothic.GetFromSession("steam", g.Request)
		if err != nil { // TODO: Check perm
			g.AbortWithStatus(http.StatusUnauthorized)
		}

		// FIXME: Just to use the variable and fix the reporting
		fmt.Println(session)

		g.JSON(http.StatusOK, gin.H{})
	})

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}

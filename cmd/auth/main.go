package main

import (
	"github.com/Cludch/csgo-tools/internal/auth"
	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/user"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
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
	authController := auth.NewController(authService, userService)

	// this store is not to be used. We use JWT tokens.
	gothic.Store = sessions.NewCookieStore([]byte(""))

	// Register Steam as Goth OpenID 2.0 provider
	goth.UseProviders(
		steam.New(configService.GetConfig().Steam.SteamAPIKey, configService.GetConfig().Auth.Host+"/auth/steam/callback"),
	)

	router.GET("/auth/:provider", authController.Auth)
	router.GET("/auth/:provider/callback", authController.Callback)

	// Protected API endpoints.
	authorized := router.Group("/")
	authorized.Use(authController.AuthorizeRequest)
	{
		authorized.GET("/me", authController.GetUserDetails)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}

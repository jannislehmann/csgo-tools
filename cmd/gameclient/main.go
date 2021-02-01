package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/entity"
	"github.com/Cludch/csgo-tools/internal/gamecoordinator"
	"github.com/Cludch/csgo-tools/pkg/demo"
	"github.com/Philipp15b/go-steam"
	"github.com/Philipp15b/go-steam/protocol/steamlang"
	"github.com/Philipp15b/go-steam/totp"
	"gorm.io/gorm"
)

var configData *config.Config
var db *gorm.DB
var csgoClient *gamecoordinator.CS

func init() {
	configData = config.GetConfiguration()
	err := steam.InitializeSteamDirectory()

	if err != nil {
		log.Error(err)
	}

	if configData.Debug == "true" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	db = entity.GetDatabase()
	gamecoordinator.DB = db
}

func main() {
	demos := demo.ScanDemosDir(config.GetConfiguration().DemosDir)
	for _, matchID := range demos {
		entity.CreateDownloadedMatchFromMatchID(matchID)
	}

	totpInstance := totp.NewTotp(configData.Steam.TwoFactorSecret)

	myLoginInfo := new(steam.LogOnDetails)
	myLoginInfo.Username = configData.Steam.Username
	myLoginInfo.Password = configData.Steam.Password
	twoFactorCode, err := totpInstance.GenerateCode()

	if err != nil {
		log.Error(err)
	}

	myLoginInfo.TwoFactorCode = twoFactorCode

	client := steam.NewClient()
	client.Connect()
	for event := range client.Events() {
		switch e := event.(type) {
		case *steam.ConnectedEvent:
			log.Info("connected to steam. Logging in...")
			client.Auth.LogOn(myLoginInfo)
		case *steam.LoggedOnEvent:
			log.Info("logged on")
			client.Social.SetPersonaState(steamlang.EPersonaState_Online)
			csgoClient = gamecoordinator.NewCSGO(client)
			csgoClient.SetPlaying(true)
			csgoClient.ShakeHands()
		case *gamecoordinator.GCReadyEvent:
			csgoClient.HandleGCReady(e)
		case steam.FatalErrorEvent:
			log.Panic(e)
		}
	}
}

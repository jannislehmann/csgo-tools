package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/entity"
	"github.com/Cludch/csgo-tools/pkg/valveapi"
	"gorm.io/gorm"
)

var configData *config.Config
var db *gorm.DB

// Sets up the global variables (config, db) and the logger
func init() {
	db = entity.GetDatabase()
	configData = config.GetConfiguration()

	if configData.Debug == "true" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
}

func main() {
	// Add accounts from config to database if not existing
	// This also adds the first known share code
	entity.AddConfigUsers(configData.CSGO)

	var csgoUsers []entity.CSGOUser

	// Create a loop that checks for new share codes each minute
	t := time.NewTicker(time.Minute)
	for {
		result := db.Preload("ShareCode").Find(&csgoUsers)

		if err := result.Error; err != nil {
			panic(err)
		}

		// Iterate all csgo users and request the next share code for the latest share code
		for _, csgoUser := range csgoUsers {
			if csgoUser.Disabled {
				continue
			}

			steamID := csgoUser.SteamID
			shareCode, err := valveapi.GetNextMatch(configData.Steam.SteamAPIKey, steamID, csgoUser.MatchHistoryAuthenticationCode, csgoUser.ShareCode.Encoded)

			// Disable user
			if err != nil {
				db.Model(&csgoUser).Update("Disabled", true)
				log.Warnf("disabled csgo user %d due to an error in fetching the share code", steamID)
				log.Error(err)
				continue
			}

			// No new match
			if shareCode == "" {
				continue
			}

			log.Infof("found match share code %v", shareCode)

			// Create share code
			sc := entity.CreateShareCodeFromEncoded(shareCode)
			// Create match
			entity.CreateMatch(sc)
			// Update csgo user
			csgoUser.UpdateLatestShareCode(sc)
		}
		<-t.C
	}
}

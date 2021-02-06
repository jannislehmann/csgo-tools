package main

import (
	"os"
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
	var nonDownloadedMatches []entity.Match

	// Create a loop that checks for new download urls
	t := time.NewTicker(time.Minute)
	for {
		result := db.Find(&nonDownloadedMatches, "download_url != '' AND downloaded = false")

		if err := result.Error; err != nil {
			panic(err)
		}

		// Iterate all matches and download them
		for _, match := range nonDownloadedMatches {
			// Download match
			err := valveapi.DownloadDemo(match.DownloadURL, configData.DemosDir, match.MatchTime)
			if err != nil {
				if os.IsTimeout(err) {
					log.Error("Lost connection", err)
					continue
				}
				log.Error(err)
			}

			// Mark as downloaded
			db.Model(&match).Update("Downloaded", true)
		}
		<-t.C
	}
}

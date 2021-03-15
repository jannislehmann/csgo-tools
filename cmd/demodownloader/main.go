package main

import (
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/entity"
	"github.com/Cludch/csgo-tools/pkg/valveapi"
	"gorm.io/gorm"
)

var configData *config.Config
var db *gorm.DB

// Sets up the global variables (config, db) and the logger.
func init() {
	db = entity.GetDatabase()
	configData = config.GetConfiguration()

	configData.SetLoggingLevel()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
}

func main() {
	var nonDownloadedMatches []entity.Match

	// Create a loop that checks for new download urls.
	t := time.NewTicker(time.Minute)
	for {
		result := db.Find(&nonDownloadedMatches, "download_url != '' AND downloaded = false")

		if err := result.Error; err != nil {
			panic(err)
		}

		// Iterate all matches and download them.
		for _, match := range nonDownloadedMatches {
			// Download match.
			url := match.DownloadURL
			err := valveapi.DownloadDemo(url, configData.DemosDir, match.MatchTime)
			if err != nil {
				if os.IsTimeout(err) {
					log.Error("Lost connection", err)
					continue
				}
				log.Error(err)
			}

			fileName := strings.Split(path.Base(url), ".")[0] + ".dem"
			// Mark as downloaded and save file name.
			db.Model(&match).Updates(entity.Match{Filename: fileName, Downloaded: true})
		}
		<-t.C
	}
}

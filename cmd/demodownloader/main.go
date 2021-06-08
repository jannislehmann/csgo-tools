package main

import (
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/match"
	"github.com/Cludch/csgo-tools/pkg/util"
)

var configService *config.Service
var matchService *match.Service

// Sets up the global variables (config, db) and the logger.
func setup() {
	configService = config.NewService()
	db := entity.NewService(configService)
	matchService = match.NewService(match.NewRepositoryMongo(db))

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
}

func main() {
	setup()

	// Create a loop that checks for new download urls.
	t := time.NewTicker(time.Minute)
	for {
		nonDownloadedMatches, err := matchService.GetDownloadableMatches()
		if err != nil {
			log.Fatal(err)
		}

		// Iterate all matches and download them.
		for _, m := range nonDownloadedMatches {
			var filename = ""
			status := m.Status

			// Download match.
			url := m.DownloadURL
			if err := util.DownloadDemo(url, configService.GetConfig().DemosDir, m.Time); err != nil {
				if os.IsTimeout(err) {
					log.Error("Lost connection", err)
					continue
				} else if util.IsDemoNotFoundError(err) {
					status = match.Unavailable
				}

				log.Error(err)
			} else {
				filename = strings.Split(path.Base(url), ".")[0] + ".dem"
				status = match.Downloaded

				const msg = "downloaded demo %v"
				log.Infof(msg, m.Filename)
			}

			// Mark as downloaded and save file name.
			if err := matchService.SetDownloaded(m, status, filename); err != nil {
				log.Error(err)
			}
		}
		<-t.C
	}
}

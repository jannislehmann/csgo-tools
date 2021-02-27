package main

import (
	"time"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/demoparser"
	"github.com/Cludch/csgo-tools/internal/entity"
	"github.com/Cludch/csgo-tools/pkg/demo"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var configData *config.Config
var db *gorm.DB

// Sets up the global variables (config, db) and the logger
func init() {
	db = entity.GetDatabase()
	configData = config.GetConfiguration()
	demoparser.ConfigData = configData

	if configData.IsDebug() {
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
	log.Info("starting demoparser")

	var nonParsedMatches []entity.Match

	// Create a loop that checks for unparsed demos
	t := time.NewTicker(time.Hour)
	for {
		// Get non-parsed matches from the db
		result := db.Find(&nonParsedMatches, "parsed = false")

		if err := result.Error; err != nil {
			log.Panic(err)
		}

		for _, match := range nonParsedMatches {
			fileName := match.Filename
			if fileName == "" {
				continue
			}
			parser := &demoparser.DemoParser{}
			demoFile := &demo.File{MatchID: match.MatchID, Filename: fileName}
			parser.Parse(configData.DemosDir, demoFile)
		}

		<-t.C
	}
}
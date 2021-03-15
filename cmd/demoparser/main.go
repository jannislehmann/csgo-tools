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

// Sets up the global variables (config, db) and the logger.
func init() {
	db = entity.GetDatabase()
	configData = config.GetConfiguration()
	demoparser.ConfigData = configData
	demoparser.DB = db

	configData.SetLoggingLevel()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
}

func main() {
	log.Info("starting demoparser")

	demos := demo.ScanDemosDir(config.GetConfiguration().DemosDir)
	for _, match := range demos {
		entity.CreateDownloadedMatchFromMatchID(match.MatchID, match.Filename, match.MatchTime)
	}

	var nonParsedMatches []entity.Match

	const numJobs = 5
	matchQueue := make(chan entity.Match, numJobs)

	// Start numJobs-times parallel workers.
	for w := 1; w <= numJobs; w++ {
		go worker(matchQueue)
	}

	// Create a loop that checks for unparsed demos.
	t := time.NewTicker(time.Hour)
	for {
		// Get non-parsed matches from the db.
		result := db.Find(&nonParsedMatches, "parsed = false")

		if err := result.Error; err != nil {
			log.Panic(err)
		}

		// Enqueue found matches.
		for _, match := range nonParsedMatches {
			matchQueue <- match
		}

		<-t.C
	}
}

// Takes a match from the channel, parses and persists it.
func worker(matches <-chan entity.Match) {
	for match := range matches {
		fileName := match.Filename
		if fileName == "" {
			return
		}

		parser := &demoparser.DemoParser{}
		demoFile := &demo.File{MatchID: match.ID, MatchTime: match.CreatedAt, Filename: fileName}

		err := parser.Parse(configData.DemosDir, demoFile)

		if err != nil {
			log.Error(err)
			return
		}

		result := parser.Match.Process()
		persistErr := result.Persist()

		if persistErr == nil && !configData.IsDebug() {
			db.Model(&match).Update("Parsed", true)
		}

		log.Infof("Finished parsing %d", demoFile.MatchID)
	}
}

package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/demoparser"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/match"
	"github.com/Cludch/csgo-tools/internal/domain/player"
	"github.com/Cludch/csgo-tools/pkg/demo"
	log "github.com/sirupsen/logrus"
)

const ParserVersion = 13

var configService *config.Service
var matchService *match.Service
var playerService *player.Service

// Sets up the global variables (config, db) and the logger.
func setup() {
	configService = config.NewService()
	db := entity.NewService(configService)
	matchService = match.NewService(match.NewRepositoryMongo(db))
	playerService = player.NewService(player.NewRepositoryMongo(db))

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
}

func main() {
	setup()

	// Scan for new local files
	demos, _ := demo.ScanDemosDir(configService.GetConfig().DemosDir)
	for _, demo := range demos {
		m, err := matchService.CreateMatchFromManualUpload(demo.Filename, demo.MatchTime)
		if err != nil {
			msg := "unable to create manual uploaded demo for file %s"
			log.Warn(msg, demo.Filename)
		} else if m != nil {
			msg := "found demo file %s and created manual upload entity"
			log.Infof(msg, m.Filename)
		}
	}

	log.Info("starting demoparser")

	numJobs, _ := strconv.ParseInt(configService.GetConfig().Parser.WorkerCount, 10, 32)
	matchQueue := make(chan *match.Match, numJobs)

	msg := "using %d workers"
	log.Infof(msg, numJobs)

	// Start numJobs-times parallel workers.
	for w := int64(1); w <= numJobs; w++ {
		go worker(matchQueue)
	}

	// Create a loop that checks for unparsed demos.
	t := time.NewTicker(time.Hour)
	for {
		// Get non-parsed matches from the db.
		nonParsedMatches, err := matchService.GetParseableMatches(ParserVersion)

		if err != nil {
			log.Fatal(err)
		}

		// Enqueue found matches.
		for _, match := range nonParsedMatches {
			matchQueue <- match
		}

		<-t.C
	}
}

// Takes a match from the channel, parses and persists it.
func worker(matches <-chan *match.Match) {
	for m := range matches {
		filename := m.Filename
		if filename == "" {
			return
		}

		parser := demoparser.NewService(configService)
		demoFile := &demo.Demo{ID: m.ID, MatchTime: m.CreatedAt, Filename: filename}

		// Check if file exists. File may have gotten deleted after being parsed the first time.
		if _, err := os.Stat(filepath.Join(configService.GetConfig().DemosDir, demoFile.Filename)); errors.Is(err, os.ErrNotExist) {
			// Set demo as unavailable.
			if err := matchService.SetStatusAndFilename(m, match.Unavailable, demoFile.Filename); err != nil {
				log.Warnf("Demo file %v for match with id %v is no longer available.", demoFile.Filename, demoFile.ID)
			}
		}

		if err := parser.Parse(configService.GetConfig().DemosDir, demoFile); err != nil {
			log.Error(err)
			continue
		}

		if !parser.GameOver {
			log.Errorf("Game %v did not finish before parsing ended. The file might be incomplete.", demoFile.Filename)
			continue
		}

		result := match.CreateResult(parser.Match)
		if err := matchService.UpdateResult(m, result, ParserVersion); err != nil {
			log.Error(err)
			continue
		}

		for _, t := range m.Result.Teams {
			for _, playerResult := range t.Players {
				player, err := playerService.GetPlayer(playerResult.SteamID)
				if err != nil {
					const msg = "main: unable to query player: %s"
					log.Errorf(msg, err)
					continue
				}

				playerResult.MatchID = m.ID
				playerResult.Map = m.Result.Map
				playerResult.Time = m.Result.Time
				playerResult.MatchRounds = byte(len(m.Result.Rounds))
				playerResult.ScoreOwnTeam = t.Wins

				// This gets the team index in the array by turning the index around.
				// There could be a smarter way, but this is a fast one.
				enemyTeamId := (t.TeamID + 1) % 2
				playerResult.ScoreEnemyTeam = m.Result.Teams[enemyTeamId].Wins
				if err := playerService.AddResult(player, playerResult); err != nil {
					log.Error(err)
				}
			}
		}

		const msg = "demoparser: finished parsing %s"
		log.Infof(msg, filename)
	}
}

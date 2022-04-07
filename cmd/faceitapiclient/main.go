package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/match"
	"github.com/Cludch/csgo-tools/internal/domain/user"
	"github.com/Cludch/csgo-tools/pkg/faceitapi"
)

var configService *config.Service
var matchService *match.Service
var userService *user.Service

// Sets up the global variables (config, db) and the logger.
func setup() {
	configService = config.NewService()
	db := entity.NewService(configService)

	matchService = match.NewService(match.NewRepositoryMongo(db))
	userService = user.NewService(user.NewRepositoryMongo(db), configService)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
}

func main() {
	setup()

	faceitAPIKey := configService.GetConfig().Faceit.FaceitAPIKey

	// Create a loop that checks for new share codes each minute.
	t := time.NewTicker(time.Minute)
	for {
		users, err := userService.GetUsersWithFaceitId()

		if err != nil {
			log.Fatal(err)
		}

		// Iterate all faceit users and request player match history and the match details
		// This is by no means efficient
		for _, u := range users {
			playerMatchHistory, err := faceitapi.GetPlayerMatchHistory(faceitAPIKey, u.Faceit.ID)
			if err != nil {
				log.Error(err)
			}

			if playerMatchHistory.Result == nil {
				continue
			}

			for _, matchHistory := range *playerMatchHistory.Result {
				matchId := matchHistory.MatchId
				matchDetails, err := faceitapi.GetMatchDetails(faceitAPIKey, matchId)

				if err != nil {
					log.Error(err)
				}

				if matchDetails == nil {
					continue
				}

				downloadUrl := matchDetails.DemoUrl[0]
				startTime := time.Unix(matchDetails.StartTime, 0)
				if _, err = matchService.CreateDownloadableMatchFromFaceitId(matchId, downloadUrl, startTime); err != nil {
					const msg = "unable to create match downloadable faceit match for id %s and url %s: %s"
					log.Errorf(msg, matchId, downloadUrl, err)
					continue
				}

				log.Infof("created downloadable faceit match for id %s", matchId)
			}
		}

		<-t.C
	}
}

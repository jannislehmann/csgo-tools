package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/match"
	"github.com/Cludch/csgo-tools/internal/domain/user"
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

	// Create a loop that checks for new share codes each minute.
	t := time.NewTicker(time.Minute)
	for {
		users, err := userService.GetUsersWithAuthenticationCode()

		if err != nil {
			log.Fatal(err)
		}

		// Iterate all csgo users and request the next share code for the latest share code.
		for _, u := range users {
			sc, err := userService.QueryLatestShareCode(u)
			if err != nil {
				log.Error(err)
			}

			if sc == nil {
				continue
			}

			if _, err = matchService.CreateMatchFromSharecode(sc); err != nil {
				log.Errorf("unable to create match from sharecode %s", err)
				continue
			}

			if err = userService.UpdateLatestShareCode(u, sc); err != nil {
				log.Errorf("unable to update user latest share code %s", err)
				continue
			}

		}

		<-t.C
	}
}

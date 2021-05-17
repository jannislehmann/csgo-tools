package gamecoordinator

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Cludch/csgo-tools/internal/entity"
	csgo "github.com/Philipp15b/go-steam/v2/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v2/protocol/gamecoordinator"
)

// DB is required when the package is initiated
var DB *gorm.DB

// Channel is used to request demos one after another.
var matchResponse chan bool

// HandleMatchList handles a gc message containing matches and tries to download those.
func (c *CS) HandleMatchList(packet *gamecoordinator.GCPacket) error {
	matchList := new(csgo.CMsgGCCStrike15V2_MatchList)
	packet.ReadProtoMsg(matchList)

	for _, match := range matchList.GetMatches() {
		for _, round := range match.GetRoundstatsall() {
			// Demo link is only linked in the last round and in this case the reserveration id is set.
			// Reservation id is the outcome id.
			if round.GetReservationid() == 0 {
				continue
			}

			matchTime := match.GetMatchtime()
			id := round.GetReservationid()
			url := round.GetMap()

			var match entity.Match

			DB.Find(&match, "id = ?", id)

			if match.DownloadURL != "" || match.Downloaded {
				continue
			}

			if match.ID == 0 {
				// Match is from the match history and does not exist in db -> create.
				match.ID = id
				DB.Create(&match)
				log.Debugf("created match %d", match.ID)
			}

			match.MatchTime = time.Unix(int64(matchTime), 0)
			match.DownloadURL = url

			DB.Save(&match)
			log.Debugf("saved match details for %d", match.ID)

		}
	}

	matchResponse <- true
	return nil
}

// HandleGCReady starts a daemon and takes non-downloaded share codes from the database.
func (c *CS) HandleGCReady(e *GCReadyEvent) {
	matchResponse = make(chan bool)

	// Download all recents games from the logged in account
	c.GetRecentGames()
	<-matchResponse

	// Request demos for non-processed share codes from the database
	var matches []entity.Match
	t := time.NewTicker(time.Minute)
	for {
		result := DB.Preload("ShareCode").Find(&matches, "download_url = '' AND downloaded = false")

		if err := result.Error; err != nil {
			panic(err)
		}

		// Request demo after another and timeout a request after 5 seconds.
		for _, m := range matches {
			sc := m.ShareCode.Encoded
			c.RequestMatch(sc)

			select {
			case <-matchResponse:
				// Continue with the next match.
			case <-time.After(15 * time.Second):
				matchResponse <- false
			}
		}

		<-t.C
	}
}

// HandleClientWelcome creates a ready event and tries sends a command to download recent games.
func (c *CS) HandleClientWelcome(packet *gamecoordinator.GCPacket) error {
	log.Info("connected to csgo gc")
	if c.isConnected {
		return nil
	}

	c.isConnected = true
	c.emit(&GCReadyEvent{})

	return nil
}

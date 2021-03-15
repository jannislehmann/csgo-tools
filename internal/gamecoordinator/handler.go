package gamecoordinator

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Cludch/csgo-tools/internal/entity"
	"github.com/Cludch/csgo-tools/internal/gamecoordinator/protocol"
	"github.com/Philipp15b/go-steam/protocol/gamecoordinator"
)

// DB is required when the package is initiated
var DB *gorm.DB

// HandleMatchList handles a gc message containing matches and tries to download those.
func (c *CS) HandleMatchList(packet *gamecoordinator.GCPacket) error {
	matchList := new(protocol.CMsgGCCStrike15V2_MatchList)
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

	return nil
}

// HandleGCReady starts a daemon and takes non-downloaded share codes from the database.
func (c *CS) HandleGCReady(e *GCReadyEvent) {
	// Download all recents games from the logged in account
	c.GetRecentGames()

	// Request demos for non-processed share codes from the database
	var matches []entity.Match
	t := time.NewTicker(time.Minute)
	for {
		result := DB.Preload("ShareCode").Find(&matches, "download_url = '' AND downloaded = false")

		if err := result.Error; err != nil {
			panic(err)
		}

		for _, match := range matches {
			sc := match.ShareCode.Encoded
			c.RequestMatch(sc)
			time.Sleep(time.Second * 5)
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

package gamecoordinator

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Cludch/csgo-tools/internal/pkg/database"
	"github.com/Cludch/csgo-tools/internal/pkg/gamecoordinator/protocol"
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
			if round.GetReservationid() == 0 {
				continue
			}

			matchTime := match.GetMatchtime()
			matchID := match.GetMatchid()
			url := round.GetMap()

			var match database.Match
			DB.Find(&match, "match_id = ?", matchID)

			match.MatchTime = time.Unix(int64(matchTime), 0)
			match.DownloadURL = url
			DB.Save(&match)
		}
	}

	return nil
}

// HandleGCReady starts a daemon and takes non-downloaded share codes from the database.
func (c *CS) HandleGCReady(e *GCReadyEvent) {
	// Download all recents games from the logged in account
	c.GetRecentGames()

	// Request demos for non-processed share codes from the database
	var matches []database.Match
	t := time.NewTicker(time.Minute * 5)
	for {
		result := DB.Preload("ShareCode").Find(&matches, "download_url = '' AND downloaded = false")

		if err := result.Error; err != nil {
			panic(err)
		}

		for _, match := range matches {
			sc := match.ShareCode.Encoded
			log.Debugf("requesting match details for %v", sc)
			go c.RequestMatch(sc)
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

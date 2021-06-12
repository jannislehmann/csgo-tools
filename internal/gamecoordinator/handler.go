package gamecoordinator

import (
	"time"

	log "github.com/sirupsen/logrus"

	csgo "github.com/Philipp15b/go-steam/v2/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v2/protocol/gamecoordinator"
)

// Channel is used to request demos one after another.
var matchResponse chan bool = make(chan bool, 1)

// HandleMatchList handles a gc message containing matches and tries to download those.
func (s *Service) HandleMatchList(packet *gamecoordinator.GCPacket) {
	matchList := new(csgo.CMsgGCCStrike15V2_MatchList)
	packet.ReadProtoMsg(matchList)

	for _, matchEntry := range matchList.GetMatches() {
		for _, round := range matchEntry.GetRoundstatsall() {
			// Demo link is only linked in the last round and in this case the reserveration id is set.
			// Reservation id is the outcome id.
			if round.GetReservationid() == 0 {
				continue
			}

			id := round.GetReservationid()
			time := time.Unix(int64(*matchEntry.Matchtime), 0)
			err := s.matchService.UpdateDownloadInformationForOutcomeId(id, time, round.GetMap())
			if err != nil {
				const msg = "gamecoordinator: %s"
				log.Errorf(msg, err)
			} else {
				const msg = "gamecoordinator: saved match details for %d"
				log.Debugf(msg, id)
			}
		}
	}

	matchResponse <- true
}

// HandleGCReady starts a daemon and takes non-downloaded share codes from the database.
func (s *Service) HandleGCReady(e *GCReadyEvent) {
	// Request demos for non-processed share codes from the database
	t := time.NewTicker(time.Minute * 5)
	for {
		matches, err := s.matchService.GetValveMatchesMissingDownloadUrl()
		if err != nil {
			log.Error(err)
		}

		// Request demo after another and timeout a request after 5 seconds.
		ch := make(chan bool, 1)
		for _, m := range matches {
			matchResponse = make(chan bool)
			sc := m.ShareCode
			go s.RequestMatch(sc)
			select {
			case <-matchResponse:
				const msg = "gamecoordinator: received response for %s"
				log.Debugf(msg, sc.Encoded)
				ch <- true
			case <-time.After(15 * time.Second):
				const msg = "gamecoordinator: failed to receive response for %s"
				log.Debugf(msg, sc.Encoded)
				ch <- false
			}

			<-ch
		}

		<-t.C
	}
}

// HandleClientWelcome creates a ready event and tries sends a command to download recent games.
func (s *Service) HandleClientWelcome(packet *gamecoordinator.GCPacket) {
	if !s.gc.isConnected {
		log.Info("connected to csgo gc")
		s.gc.isConnected = true
		s.emit(&GCReadyEvent{})
	}
}

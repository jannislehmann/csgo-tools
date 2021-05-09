package gamecoordinator

import (
	"github.com/Cludch/csgo-tools/internal/entity"
	csgo "github.com/Philipp15b/go-steam/v2/csgo/protocol/protobuf"
	"github.com/golang/protobuf/proto" //nolint //thinks break if we use the new package

	log "github.com/sirupsen/logrus"
)

// GetRecentGames requests the players match history.
func (c *CS) GetRecentGames() {
	newAccID := c.client.SteamId().ToUint64() - 76561197960265728
	c.Write(uint32(csgo.ECsgoGCMsg_k_EMsgGCCStrike15_v2_MatchListRequestRecentUserGames), &csgo.CMsgGCCStrike15V2_MatchListRequestRecentUserGames{
		Accountid: proto.Uint32(uint32(newAccID)),
	})
}

// RequestMatch requests the match information for a share code
func (c *CS) RequestMatch(shareCode string) {
	// Decode share code
	sc := entity.DecodeShareCode(shareCode)
	if sc == nil {
		return
	}

	log.Debugf("requesting match details for %v %d", sc.Encoded, sc.OutcomeID)

	// Request match info
	c.Write(uint32(csgo.ECsgoGCMsg_k_EMsgGCCStrike15_v2_MatchListRequestFullGameInfo), &csgo.CMsgGCCStrike15V2_MatchListRequestFullGameInfo{
		Matchid:   proto.Uint64(uint64(sc.MatchID)),
		Outcomeid: proto.Uint64(uint64(sc.OutcomeID)),
		Token:     proto.Uint32(uint32(sc.Token)),
	})
}

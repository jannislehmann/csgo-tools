package gamecoordinator

import (
	"github.com/Cludch/csgo-tools/internal/entity"
	"github.com/Cludch/csgo-tools/internal/gamecoordinator/protocol"
	"github.com/golang/protobuf/proto" //nolint //thinks break if we use the new package
)

// GetRecentGames requests the players match history.
func (c *CS) GetRecentGames() {
	newAccID := c.client.SteamId().ToUint64() - 76561197960265728
	c.Write(uint32(protocol.ECsgoGCMsg_k_EMsgGCCStrike15_v2_MatchListRequestRecentUserGames), &protocol.CMsgGCCStrike15V2_MatchListRequestRecentUserGames{
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

	// Request match info
	c.Write(uint32(protocol.ECsgoGCMsg_k_EMsgGCCStrike15_v2_MatchListRequestFullGameInfo), &protocol.CMsgGCCStrike15V2_MatchListRequestFullGameInfo{
		Matchid:   proto.Uint64(uint64(sc.MatchID)),
		Outcomeid: proto.Uint64(uint64(sc.OutcomeID)),
		Token:     proto.Uint32(uint32(sc.Token)),
	})
}

package player

import (
	"time"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
)

type Player struct {
	ID        uint64          `json:"steamId" bson:"_id"`
	CreatedAt time.Time       `json:"-" bson:"createdAt"`
	FaceitID  entity.ID       `json:"faceitId" bson:"faceitId,omitempty"`
	Results   []*PlayerResult `json:"results" bson:"results"`
}

// PlayerResult holds different performance metrics from one game.
type PlayerResult struct {
	MatchID      entity.ID `json:"matchId" bson:"matchId"`
	SteamID      uint64    `json:"steamId" bson:"steamId"`
	Name         string    `json:"name" bson:"name"`
	Kills        byte      `json:"kills" bson:"kills"`
	EntryKills   byte      `json:"entryKills" bson:"entryKills"`
	Headshots    byte      `json:"headshots" bson:"headshots"`
	Assists      byte      `json:"assists" bson:"assists"`
	Deaths       byte      `json:"deaths" bson:"deaths"`
	MVPs         byte      `json:"mvps" bson:"mvps"`
	Won1v3       byte      `json:"won1v3" bson:"won1v3"`
	Won1v4       byte      `json:"won1v4" bson:"won1v4"`
	Won1v5       byte      `json:"won1v5" bson:"won1v5"`
	RoundsWith3K byte      `json:"3k" bson:"3k"`
	RoundsWith4K byte      `json:"4k" bson:"4k"`
	RoundsWith5K byte      `json:"5k" bson:"5k"`
}

func NewPlayer(id uint64) (*Player, error) {
	m := &Player{
		ID:        id,
		CreatedAt: time.Now(),
		Results:   []*PlayerResult{},
	}

	return m, nil
}

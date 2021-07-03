package match

import (
	"time"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
)

type MatchList struct {
	Matches []*MatchListEntry `json:"matches"`
}

type MatchListEntry struct {
	ID     entity.ID `json:"id"`
	Time   time.Time `json:"time"`
	Map    string    `json:"map"`
	Result string    `json:"result"`
}

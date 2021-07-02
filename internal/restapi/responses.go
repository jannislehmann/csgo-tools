package restapi

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

// PlayerGameStats describes average stats across all matches.
type PlayerGameStats struct {
	Games             int    `json:"games"`
	SteamID           uint64 `json:"steamId"`
	KillsPerGame      int    `json:"killsPerGame"`
	EntryKillsPerGame int    `json:"entryKillsPerGame"`
	HeadshotsPerGame  int    `json:"headshotsPerGame"`
	AssistsPerGame    int    `json:"assistsPerGame"`
	DeathsPerGame     int    `json:"deathsPerGame"`
	MVPsPerGame       int    `json:"mvpsPerGame"`
	Won1v3            int    `json:"won1v3"`
	Won1v4            int    `json:"won1v4"`
	Won1v5            int    `json:"won1v5"`
	RoundsWith3K      int    `json:"roundsWith3k"`
	RoundsWith4K      int    `json:"roundsWith4k"`
	RoundsWith5K      int    `json:"roundsWith5k"`
}

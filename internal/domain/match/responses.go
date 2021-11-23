package match

import (
	"time"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
)

type MatchList struct {
	Matches []*MatchListEntry `json:"matches"`
}

type MatchListEntry struct {
	ID           entity.ID   `json:"id"`
	Time         time.Time   `json:"time"`
	Map          string      `json:"map"`
	TeamOneScore byte        `json:"teamOneScore"`
	TeamTwoScore byte        `json:"teamTwoScore"`
	ClanTeam     common.Team `json:"clanTeam"`
}

package demoparser

import (
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
)

// GetTeamIndex returns 0 for T, 1 for CT and 2 for everything else.
func GetTeamIndex(team common.Team, sidesSwitched bool) byte {
	if team == common.TeamTerrorists {
		if !sidesSwitched {
			return 0
		}
		return 1
	} else if team == common.TeamCounterTerrorists {
		if !sidesSwitched {
			return 1
		}
		return 0
	}

	// Could also return an error here but we do not expect this to happen.
	return 2
}

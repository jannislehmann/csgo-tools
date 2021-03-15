package demoparser

import (
	"errors"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/pkg/demo"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	log "github.com/sirupsen/logrus"
)

// ConfigData holds the configuration instance for the package.
var ConfigData *config.Config

// DemoParser holds the instance of one demo consisting of the file handle and the parsed data.
type DemoParser struct {
	parser        demoinfocs.Parser
	Match         *MatchData
	CurrentRound  byte
	RoundStart    time.Duration
	RoundOngoing  bool
	SidesSwitched bool
}

// MatchData holds information about the match itself.
type MatchData struct {
	ID       uint64
	Map      string
	Header   *common.DemoHeader
	Players  []*Player
	Teams    [2]*Team
	Duration time.Duration
	Time     time.Time
	Rounds   []*Round
}

// Team represents a team and links to it's players.
type Team struct {
	StartedAs common.Team
	State     *common.TeamState
	Players   []*Player
}

// Player represents one player either as T or CT.
type Player struct {
	SteamID uint64
	Name    string
	Team    *Team
}

// Round contains information about one round.
type Round struct {
	Duration  time.Duration
	Kills     []*Kill
	Winner    *Team
	WinReason events.RoundEndReason
	MVP       *Player
}

// Kill holds information about a kill that happenend during the match.
type Kill struct {
	Time          time.Duration
	Victim        *Player
	Killer        *Player
	Assister      *Player
	Weapon        common.EquipmentType
	IsHeadshot    bool
	AssistedFlash bool
	AttackerBlind bool
	NoScope       bool
	ThroughSmoke  bool
	ThroughWall   bool
}

// Parse takes a demo file and starts parsing by registering all required event handlers.
func (p *DemoParser) Parse(dir string, demoFile *demo.File) error {
	p.Match = &MatchData{ID: demoFile.MatchID, Time: demoFile.MatchTime}

	f, err := os.Open(path.Join(dir, demoFile.Filename))

	if err != nil {
		return err
	}

	log.Infof("Starting demo parsing of match %d", p.Match.ID)

	p.parser = demoinfocs.NewParser(f)
	defer p.parser.Close()
	defer f.Close()

	// Parsing the header within an event handler crashes.
	header, _ := p.parser.ParseHeader()
	p.Match.Header = &header

	// Register all handler
	p.parser.RegisterEventHandler(p.handleMatchStart)
	p.parser.RegisterEventHandler(p.handleGamePhaseChanged)
	p.parser.RegisterEventHandler(p.handleKill)
	p.parser.RegisterEventHandler(p.handleMVP)
	p.parser.RegisterEventHandler(p.handleRoundStart)
	p.parser.RegisterEventHandler(p.handleRoundEnd)

	return p.parser.ParseToEnd()
}

func (p *DemoParser) getPlayer(player *common.Player) (*Player, error) {
	for _, localPlayer := range p.Match.Players {
		if player.SteamID64 == localPlayer.SteamID {
			return localPlayer, nil
		}
	}

	for _, gamePlayer := range p.parser.GameState().Participants().Playing() {
		if player.SteamID64 == gamePlayer.SteamID64 {
			return p.AddPlayer(player), nil
		}
	}

	return nil, errors.New("Player not found in local match struct " + strconv.FormatUint(player.SteamID64, 10))
}

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

// AddPlayer adds a player to the game and returns the pointer.
func (p *DemoParser) AddPlayer(player *common.Player) *Player {
	teamID := GetTeamIndex(player.Team, p.SidesSwitched)
	teams := p.Match.Teams
	teamPlayers := teams[teamID].Players

	customPlayer := &Player{SteamID: player.SteamID64, Name: player.Name, Team: teams[teamID]}

	teams[teamID].Players = append(teamPlayers, customPlayer)
	p.Match.Players = append(p.Match.Players, customPlayer)

	return customPlayer
}

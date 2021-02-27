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

// https://github.com/megaclan3000/megaclan3000/blob/master/internal/demoparser/parser.go

// DemoParser holds the instance of one demo consisting of the file handle and the parsed data.
type DemoParser struct {
	parser       demoinfocs.Parser
	Match        *MatchData
	CurrentRound int
	RoundStart   time.Duration
	RoundOngoing bool
	IsFirstHalf  bool
}

// MatchData holds information about the match itself.
type MatchData struct {
	ID            uint64
	Map           string
	Header        *common.DemoHeader
	Teams         map[int]*Team
	Players       []*Player
	MatchDuration time.Duration
	Rounds        []*Round
}

type Team struct {
	TeamState *common.TeamState
	Players   []*Player
	StartedAs common.Team
}

type Player struct {
	SteamID uint64
	Name    string
	Team    *Team
}

type Round struct {
	Duration  time.Duration
	Kills     []*Kill
	Winner    *Team
	WinReason events.RoundEndReason
	MVP       *Player
}

// Kill holds information about a kill that happenend during the match.
type Kill struct {
	Time       time.Duration
	IsHeadshot bool
	Victim     *Player
	Killer     *Player
	Assister   *Player
	Weapon     common.EquipmentType
}

func (p *DemoParser) Parse(dir string, demoFile *demo.File) error {
	p.Match = &MatchData{ID: demoFile.MatchID}

	f, err := os.Open(path.Join(dir, demoFile.Filename))

	if err != nil {
		return err
	}

	log.Infof("Starting demo parsing of match %d", p.Match.ID)

	p.parser = demoinfocs.NewParser(f)
	defer p.parser.Close()
	defer f.Close()

	// Parsing the header within an event handler crashes.
	header, err := p.parser.ParseHeader()
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

	return nil, errors.New("Player not found in local match struct " + strconv.FormatUint(player.SteamID64, 10))
}

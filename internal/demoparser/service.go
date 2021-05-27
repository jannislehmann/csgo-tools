package demoparser

import (
	"errors"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/pkg/demo"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	log "github.com/sirupsen/logrus"
)

// ConfigData holds the configuration instance for the package.
var ConfigData *config.Config

// DemoParser holds the instance of one demo consisting of the file handle and the parsed data.
type Service struct {
	configurationService config.UseCase
	parser               demoinfocs.Parser
	Match                *MatchData
	CurrentRound         byte
	RoundStart           time.Duration
	RoundOngoing         bool
	SidesSwitched        bool
}

func NewService(c config.UseCase) *Service {
	return &Service{
		configurationService: c,
	}
}

// MatchData holds information about the match itself.
type MatchData struct {
	ID       entity.ID
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
	IsDuringRound bool
	IsHeadshot    bool
	AssistedFlash bool
	AttackerBlind bool
	NoScope       bool
	ThroughSmoke  bool
	ThroughWall   bool
}

// Parse takes a demo file and starts parsing by registering all required event handlers.
func (s *Service) Parse(dir string, demoFile *demo.Demo) error {
	s.Match = &MatchData{ID: demoFile.ID, Time: demoFile.MatchTime}

	f, err := os.Open(path.Join(dir, demoFile.Filename))

	if err != nil {
		return err
	}

	log.Infof("Starting demo parsing of match %d", s.Match.ID)

	s.parser = demoinfocs.NewParser(f)
	defer s.parser.Close()
	defer f.Close()

	// Parsing the header within an event handler crashes.
	header, _ := s.parser.ParseHeader()
	s.Match.Header = &header

	// Register all handler
	s.parser.RegisterEventHandler(s.handleMatchStart)
	s.parser.RegisterEventHandler(s.handleGamePhaseChanged)
	s.parser.RegisterEventHandler(s.handleKill)
	s.parser.RegisterEventHandler(s.handleMVP)
	s.parser.RegisterEventHandler(s.handleRoundStart)
	s.parser.RegisterEventHandler(s.handleRoundEnd)

	return s.parser.ParseToEnd()
}

func (s *Service) getPlayer(player *common.Player) (*Player, error) {
	if player.IsBot {
		return nil, errors.New("Player is a bot")
	}

	for _, localPlayer := range s.Match.Players {
		if player.SteamID64 == localPlayer.SteamID {
			return localPlayer, nil
		}
	}

	for _, gamePlayer := range s.parser.GameState().Participants().Playing() {
		if player.SteamID64 == gamePlayer.SteamID64 {
			return s.AddPlayer(player), nil
		}
	}

	return nil, errors.New("Player not found in local match struct " + strconv.FormatUint(player.SteamID64, 10))
}

// AddPlayer adds a player to the game and returns the pointer.
func (s *Service) AddPlayer(player *common.Player) *Player {
	teamID := GetTeamIndex(player.Team, s.SidesSwitched)
	teams := s.Match.Teams
	teamPlayers := teams[teamID].Players

	customPlayer := &Player{SteamID: player.SteamID64, Name: player.Name, Team: teams[teamID]}

	teams[teamID].Players = append(teamPlayers, customPlayer)
	s.Match.Players = append(s.Match.Players, customPlayer)

	return customPlayer
}

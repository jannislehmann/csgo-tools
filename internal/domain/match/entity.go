package match

import (
	"time"

	"github.com/Cludch/csgo-tools/internal/demoparser"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/pkg/share_code"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	log "github.com/sirupsen/logrus"
)

// Source describes the source of a demo / match.
type Source string

const (
	MatchMaking Source = "MatchMaking"
	Faceit      Source = "Faceit"
	Manual      Source = "Manual"
)

// Status describes the current status / step of the match and its processing.
type Status string

const (
	Created      Status = "Created"
	Downloaded   Status = "Downloaded"
	Expired      Status = "Expired"
	Unavailable  Status = "Unavailable"
	Error        Status = "Error"
	Downloadable Status = "Downloadable"
	Parsed       Status = "Parsed"
)

// Match holds the central information about a csgo match from different data sources.
type Match struct {
	ID            entity.ID                 `json:"id" bson:"_id,omitempty"`
	CreatedAt     time.Time                 `json:"-" bson:"createdAt"`
	Source        Source                    `json:"source" bson:"source"`
	Status        Status                    `json:"status" bson:"status"`
	Time          time.Time                 `json:"time" bson:"time,omitempty"`
	Filename      string                    `json:"filename" bson:"filename,omitempty"`
	DownloadURL   string                    `json:"url" bson:"url,omitempty"`
	ShareCode     *share_code.ShareCodeData `json:"shareCode" bson:"shareCode,omitempty"`
	FaceitMatchId string                    `json:"faceitMatchId" bson:"faceitMatchId,omitempty"`
	Result        *MatchResult              `json:"result" bson:"result,omitempty"`
}

// MatchResult holds meta data and the teams of one match.
type MatchResult struct {
	ParserVersion byte          `json:"parserVersion" bson:"parserVersion"`
	Map           string        `json:"map" bson:"map"`
	Time          time.Time     `json:"time" bson:"time"`
	Duration      time.Duration `json:"duration" bson:"duration"`
	// 0 = T / 1 = CT
	Teams []*TeamResult `json:"teams"`
}

// TeamResult describes the players and wins for one team.
type TeamResult struct {
	// TeamID describes the side the team started as.
	TeamID          common.Team     `json:"id" bson:"id"`
	Players         []*PlayerResult `json:"players" bson:"players"`
	Wins            byte            `json:"wins" bson:"wins"`
	PistolRoundWins byte            `json:"pistolRoundWins" bson:"pistolRoundWins"`
}

// PlayerResult holds different performance metrics from one game.
type PlayerResult struct {
	SteamID      uint64 `json:"steamId" bson:"steamId"`
	Name         string `json:"name" bson:"name"`
	Kills        byte   `json:"kills" bson:"kills"`
	EntryKills   byte   `json:"entryKills" bson:"entryKills"`
	Headshots    byte   `json:"headshots" bson:"headshots"`
	Assists      byte   `json:"assists" bson:"assists"`
	Deaths       byte   `json:"deaths" bson:"deaths"`
	MVPs         byte   `json:"mvps" bson:"mvps"`
	Won1v3       byte   `json:"won1v3" bson:"won1v3"`
	Won1v4       byte   `json:"won1v4" bson:"won1v4"`
	Won1v5       byte   `json:"won1v5" bson:"won1v5"`
	RoundsWith3K byte   `json:"3k" bson:"3k"`
	RoundsWith4K byte   `json:"4k" bson:"4k"`
	RoundsWith5K byte   `json:"5k" bson:"5k"`
}

func NewMatch(source Source) (*Match, error) {
	m := &Match{
		ID:        entity.NewID(),
		CreatedAt: time.Now(),
		Source:    source,
		Status:    Created,
	}

	return m, nil
}

// Process processes the match data and creates more performance-based results per player in order to persist these in the database.
func CreateResult(m *demoparser.MatchData) *MatchResult {
	// Create result.
	result := &MatchResult{Map: m.Map, Duration: m.Duration, Time: m.Time, Teams: make([]*TeamResult, 2)}

	// Create teams.
	for _, team := range m.Teams {
		// Could also use team.State.ID - 2 as they return the same as the enum.
		result.Teams[demoparser.GetTeamIndex(team.StartedAs, false)] = &TeamResult{TeamID: team.StartedAs}
	}

	// Create players.
	for _, player := range m.Players {
		if player.SteamID == 0 {
			log.Debugf("steamid 0 for %s in %d", player.Name, m.ID)
		}

		// Get starting team and append player.
		team := result.Teams[demoparser.GetTeamIndex(player.Team.StartedAs, false)]
		team.Players = append(team.Players, &PlayerResult{SteamID: player.SteamID, Name: player.Name})
	}

	result.processRounds(m.Rounds)

	return result
}

func (m *MatchResult) processRounds(rounds []*demoparser.Round) {
	for index, round := range rounds {

		// MVP can be nil when the round ended because one team surrendered.
		mvp := m.getPlayer(round.MVP)
		if mvp != nil {
			mvp.MVPs++
		}

		winner := m.getTeam(round.Winner.StartedAs)
		winner.Wins++

		// Pistol round wins.
		roundNumber := index + 1
		if roundNumber == 1 || roundNumber == 16 {
			winner.PistolRoundWins++
		}

		playerKills := make(map[*PlayerResult]byte)

		// Process in round function in order to calculate all round information like amount of kills / round.
		for _, kill := range round.Kills {
			// Victim may be null, if it was a bot.
			if kill.Victim != nil {
				m.getPlayer(kill.Victim).Deaths++
			}

			// Killer may not be set if the player died e.g. through fall damage.
			if kill.Killer != nil {
				killer := m.getPlayer(kill.Killer)
				killer.Kills++
				if kill.IsHeadshot {
					killer.Headshots++
				}

				if _, found := playerKills[killer]; !found {
					playerKills[killer] = 0
				}
				playerKills[killer]++
			}

			// Assister may not be set.
			if kill.Assister != nil {
				assister := m.getPlayer(kill.Assister)
				assister.Assists++
			}
		}

		// Increase players 3/4/5 Kills per round.
		for player, kills := range playerKills {
			if kills <= 2 {
				continue
			}

			switch kills {
			case 3:
				player.RoundsWith3K++
			case 4:
				player.RoundsWith4K++
			case 5:
				player.RoundsWith5K++
			}
		}
	}
}

func (m *MatchResult) getTeam(team common.Team) *TeamResult {
	return m.Teams[demoparser.GetTeamIndex(team, false)]
}

func (m *MatchResult) getPlayer(player *demoparser.Player) *PlayerResult {
	if player == nil {
		return nil
	}

	for _, team := range m.Teams {
		for _, resultPlayer := range team.Players {
			if resultPlayer == nil {
				return nil
			}

			if resultPlayer.SteamID == player.SteamID {
				return resultPlayer
			}
		}
	}

	return nil
}

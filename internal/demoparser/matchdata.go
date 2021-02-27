package demoparser

import (
	"time"

	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	log "github.com/sirupsen/logrus"
)

// TODO: Persist all the result entities

// MatchResult holds meta data and the teams of one match.
type MatchResult struct {
	ID       uint64
	Map      string
	Time     time.Time
	Duration time.Duration
	Teams    [2]*TeamResult // 0 = T / 1 = CT
}

// TeamResult describes the players and wins for one team.
type TeamResult struct {
	StartedAs       common.Team
	Players         []*PlayerResult
	Wins            uint8
	PistolRoundWins uint8
}

// PlayerResult holds different performance metrics from one game.
type PlayerResult struct {
	SteamID    uint64
	Name       string
	Kills      uint8
	EntryKills uint8
	Headshots  uint8
	Assists    uint8
	Deaths     uint8
	MVPs       uint8
	// 1v3,4,5
	// 3k,4k,5k
	// Weapon stats
}

// Process processes the match data and creates more performance-based results per player in order to persist these in the database.
func (m *MatchData) Process() *MatchResult {
	// Create result
	result := &MatchResult{ID: m.ID, Map: m.Map, Duration: m.Duration, Time: m.Time}

	// Create teams
	for _, team := range m.Teams {
		// Could also use team.State.ID - 2 as they return the same as the enum
		result.Teams[getTeamIndex(team.StartedAs)] = &TeamResult{StartedAs: team.StartedAs}
	}

	// Create players
	for _, player := range m.Players {
		// Get starting team and append player
		team := result.Teams[getTeamIndex(player.Team.StartedAs)]
		team.Players = append(team.Players, &PlayerResult{SteamID: player.SteamID, Name: player.Name})
	}

	result.processRounds(m.Rounds)

	return result
}

func (m *MatchResult) processRounds(rounds []*Round) {
	for index, round := range rounds {
		mvp := m.getPlayer(round.MVP)
		mvp.MVPs++

		winner := m.getTeam(round.Winner.StartedAs)
		winner.Wins++

		// Pistol round wins
		roundNumber := index + 1
		if roundNumber == 1 || roundNumber == 16 {
			winner.PistolRoundWins++
		}

		// Process in round function in order to calculate all round information like amount of kills / round
		for _, kill := range round.Kills {
			m.getPlayer(kill.Victim).Deaths++
			// Killer may not be set if the player died e.g. through fall damage
			if kill.Killer != nil {
				killer := m.getPlayer(kill.Killer)
				killer.Kills++
				if kill.IsHeadshot {
					killer.Headshots++
				}
			}
			// Assister may not be set
			if kill.Assister != nil {
				m.getPlayer(kill.Assister).Assists++
			}
		}
	}
}

// Print prints the game's result.
func (m *MatchResult) Print() {
	log.Infof("Match result was %d-%d", m.Teams[0].Wins, m.Teams[1].Wins)
	for _, team := range m.Teams {
		for _, player := range team.Players {
			// TODO: Sort by kills
			headshotAccuracy := uint8(0)
			if player.Headshots > 0 && player.Kills > 0 {
				headshotAccuracy = uint8(float32(player.Headshots) / float32(player.Kills) * 100)
			}
			log.Infof("Player %v had a K/D/A of %d/%d/%d and %d headshots. HS accuracy was %d percent. Entry kills %d and %d MVPs", player.Name, player.Kills, player.Deaths,
				player.Assists, player.Headshots, headshotAccuracy, player.EntryKills, player.MVPs)
		}
	}
}

func getTeamIndex(team common.Team) uint8 {
	return GetTeamIndex(team, false)
}

func (m *MatchResult) getTeam(team common.Team) *TeamResult {
	return m.Teams[getTeamIndex(team)]
}

func (m *MatchResult) getPlayer(player *Player) *PlayerResult {
	for _, team := range m.Teams {
		for _, resultPlayer := range team.Players {
			if resultPlayer.SteamID == player.SteamID {
				return resultPlayer
			}
		}
	}

	return nil
}

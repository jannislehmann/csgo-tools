package demoparser

import (
	"sort"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// DB is required when the package is initiated.
var DB *gorm.DB

// MatchResult holds meta data and the teams of one match.
type MatchResult struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uint64         `gorm:"primaryKey"`
	Map       string
	Time      time.Time
	Duration  time.Duration
	// 0 = T / 1 = CT
	Teams []*TeamResult `gorm:"foreignKey:ID"`
}

// TeamResult describes the players and wins for one team.
type TeamResult struct {
	gorm.Model
	MatchID         uint64 `gorm:"primaryKey"`
	StartedAs       common.Team
	Players         []*PlayerResult `gorm:"foreignKey:SteamID"`
	Wins            byte
	PistolRoundWins byte
}

// PlayerResult holds different performance metrics from one game.
type PlayerResult struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	MatchID      uint64         `gorm:"primaryKey"`
	SteamID      uint64         `gorm:"primaryKey"`
	Name         string
	Kills        byte
	EntryKills   byte
	Headshots    byte
	Assists      byte
	Deaths       byte
	MVPs         byte
	Won1v3       byte
	Won1v4       byte
	Won1v5       byte
	RoundsWith3K byte
	RoundsWith4K byte
	RoundsWith5K byte
}

// Process processes the match data and creates more performance-based results per player in order to persist these in the database.
func (m *MatchData) Process() *MatchResult {
	// Create result.
	result := &MatchResult{ID: m.ID, Map: m.Map, Duration: m.Duration, Time: m.Time, Teams: make([]*TeamResult, 2)}

	// Create teams.
	for _, team := range m.Teams {
		// Could also use team.State.ID - 2 as they return the same as the enum.
		result.Teams[getTeamIndex(team.StartedAs)] = &TeamResult{StartedAs: team.StartedAs}
	}

	// Create players
	for _, player := range m.Players {
		// Get starting team and append player.
		team := result.Teams[getTeamIndex(player.Team.StartedAs)]
		team.Players = append(team.Players, &PlayerResult{SteamID: player.SteamID, Name: player.Name})
	}

	result.processRounds(m.Rounds)

	return result
}

func (m *MatchResult) processRounds(rounds []*Round) {
	for index, round := range rounds {

		// MVP can be nil when the round ended because one team surrendered.
		mvp := m.getPlayer(round.MVP)
		if mvp != nil {
			mvp.MVPs++
		}

		winner := m.getTeam(round.Winner.StartedAs)
		winner.Wins++

		// Pistol round wins
		roundNumber := index + 1
		if roundNumber == 1 || roundNumber == 16 {
			winner.PistolRoundWins++
		}

		playerKills := make(map[*PlayerResult]byte)

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
				playerKills[killer]++
			}

			// Assister may not be set
			if kill.Assister != nil {
				m.getPlayer(kill.Assister).Assists++
			}
		}

		// Increase players 3/4/5 Kills per round.
		for player, kills := range playerKills {
			if kills >= 2 {
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

// Print prints the game's result.
func (m *MatchResult) Print() {
	log.Infof("Match result was %d-%d", m.Teams[0].Wins, m.Teams[1].Wins)
	for _, team := range m.Teams {
		players := team.Players

		sort.Slice(players[:], func(i, j int) bool {
			return players[i].Kills < players[j].Kills
		})

		for _, player := range players {
			headshotAccuracy := byte(0)
			if player.Headshots > 0 && player.Kills > 0 {
				headshotAccuracy = byte(float32(player.Headshots) / float32(player.Kills) * 100)
			}
			log.Infof("Player %v had a K/D/A of %d/%d/%d and %d headshots. HS accuracy was %d percent. Entry kills %d and %d MVPs", player.Name, player.Kills, player.Deaths,
				player.Assists, player.Headshots, headshotAccuracy, player.EntryKills, player.MVPs)
		}
	}
}

// Persist persists the match results in the database.
func (m *MatchResult) Persist() error {
	// TODO: Persist
	return DB.Create(&m).Error
}

func getTeamIndex(team common.Team) byte {
	return GetTeamIndex(team, false)
}

func (m *MatchResult) getTeam(team common.Team) *TeamResult {
	return m.Teams[getTeamIndex(team)]
}

func (m *MatchResult) getPlayer(player *Player) *PlayerResult {
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

package match

import (
	"time"

	"github.com/Cludch/csgo-tools/internal/demoparser"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/player"
	"github.com/Cludch/csgo-tools/pkg/share_code"
	"github.com/go-playground/validator"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	log "github.com/sirupsen/logrus"
)

var validate = validator.New()

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
	Source        Source                    `json:"source" bson:"source" validate:"required"`
	Status        Status                    `json:"status" bson:"status" validate:"required"`
	Time          time.Time                 `json:"time" bson:"time,omitempty"`
	Filename      string                    `json:"filename" bson:"filename,omitempty"`
	DownloadURL   string                    `json:"url" bson:"url,omitempty"`
	ShareCode     *share_code.ShareCodeData `json:"shareCode" bson:"shareCode,omitempty"`
	FaceitMatchId string                    `json:"faceitMatchId" bson:"faceitMatchId,omitempty"`
	Result        *MatchResult              `json:"result" bson:"result,omitempty" validation:"dive"`
}

// MatchResult holds meta data and the teams of one match.
type MatchResult struct {
	ParserVersion byte          `json:"parserVersion" bson:"parserVersion" validate:"required,gt=0"`
	Map           string        `json:"map" bson:"map" validate:"required"`
	Time          time.Time     `json:"time" bson:"time" validate:"required"`
	Duration      time.Duration `json:"duration" bson:"duration" validate:"required"`
	// 0 = T / 1 = CT
	Teams  []*TeamResult  `json:"teams" validate:"required,dive"`
	Rounds []*RoundResult `json:"rounds" validate:"required,dive"`
}

// TeamResult describes the players and wins for one team.
type TeamResult struct {
	// TeamID describes the side the team started as.
	TeamID          common.Team            `json:"id" bson:"id" validate:"required,gte=2,lte=3"`
	Players         []*player.PlayerResult `json:"players" bson:"players" validate:"required,dive"`
	Wins            byte                   `json:"wins" bson:"wins"`
	PistolRoundWins byte                   `json:"pistolRoundWins" bson:"pistolRoundWins"`
}

// RoundResult contains information about a single round.
type RoundResult struct {
	RoundNumber  byte          `json:"roundNumber" bson:"roundNumber" validate:"required"`
	Duration     time.Duration `json:"duration" bson:"duration" validate:"required"`
	Kills        []*KillResult `json:"kills" bson:"kills"  validate:"required,dive"`
	MVPPlayerID  uint64        `json:"mvp" bson:"mvpPlayerId"`
	WinnerTeamID common.Team   `json:"winnerTeamId" bson:"winnerTeamId" validate:"required,gte=2,lte=3"`
}

// KillResult contains information about a kill.
type KillResult struct {
	Tick            time.Duration        `json:"tick" bson:"tick" validate:"required"`
	VictimID        uint64               `json:"victimId" bson:"victimId"`
	KillerID        uint64               `json:"killerId" bson:"killerId"`
	AssisterID      uint64               `json:"assisterId" bson:"assisterId"`
	Weapon          common.EquipmentType `json:"weapon" bson:"" validate:"required"`
	IsDuringRound   bool                 `json:"isDuringRound" bson:"isDuringRound"`
	IsHeadshot      bool                 `json:"isHeadshot" bson:"isHeadshot"`
	IsFlashAssist   bool                 `json:"isFlashAssist" bson:"isFlashAssist"`
	IsAttackerBlind bool                 `json:"isAttackerBlind" bson:"isAttackerBlind"`
	IsNoScope       bool                 `json:"isNoScope" bson:"isNoScope"`
	IsThroughSmoke  bool                 `json:"isThroughSmoke" bson:"isThroughSmoke"`
	IsThroughWall   bool                 `json:"isThroughWall" bson:"isThroughWall"`
}

func NewMatch(source Source) (*Match, error) {
	m := &Match{
		ID:        entity.NewID(),
		CreatedAt: time.Now(),
		Source:    source,
		Status:    Created,
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

// Process processes the match data and creates more performance-based results per player in order to persist these in the database.
func CreateResult(m *demoparser.MatchData) *MatchResult {
	// Create result.
	result := &MatchResult{Map: m.Map, Duration: m.Duration, Time: m.Time, Teams: make([]*TeamResult, 2), Rounds: make([]*RoundResult, len(m.Rounds))}

	// Create teams.
	for _, team := range m.Teams {
		// Could also use team.State.ID - 2 as they return the same as the enum.
		result.Teams[demoparser.GetTeamIndex(team.StartedAs, false)] = &TeamResult{TeamID: team.StartedAs}
	}

	// Create players.
	for _, p := range m.Players {
		if p.SteamID == 0 {
			const msg = "match: steamid 0 for %s in %d"
			log.Debugf(msg, p.Name, m.ID)
		}

		// Get starting team and append player.
		team := result.Teams[demoparser.GetTeamIndex(p.Team.StartedAs, false)]
		team.Players = append(team.Players, &player.PlayerResult{SteamID: p.SteamID, Name: p.Name, WinCount: p.WinCount, RankOld: p.RankOld, RankNew: p.RankNew})
	}

	result.processRounds(m.Rounds)

	return result
}

// CreateRoundResult takes a parsed round and returns a persistable RoundResult.
func CreateRoundResult(r *demoparser.Round) *RoundResult {
	return &RoundResult{Duration: r.Duration, Kills: make([]*KillResult, len(r.Kills))}
}

// CreateKillResult takes a parsed kill and returns a persistable KillResult.
func CreateKillResult(k *demoparser.Kill) *KillResult {
	killResult := &KillResult{Tick: k.Tick,
		Weapon: k.Weapon, IsDuringRound: k.IsDuringRound, IsHeadshot: k.IsHeadshot, IsFlashAssist: k.IsFlashAssist, IsAttackerBlind: k.IsAttackerBlind,
		IsNoScope: k.IsNoScope, IsThroughSmoke: k.IsThroughSmoke, IsThroughWall: k.IsThroughWall,
	}

	// --- Check whether the victim, killer and assister are set before trying to access its property.

	if k.Victim != nil {
		killResult.VictimID = k.Victim.SteamID
	}

	if k.Killer != nil {
		killResult.KillerID = k.Killer.SteamID
	}

	if k.Assister != nil {
		killResult.AssisterID = k.Assister.SteamID
	}

	return killResult
}

func (m *MatchResult) processRounds(rounds []*demoparser.Round) {
	for index, round := range rounds {
		// Create round result and log
		roundResult := CreateRoundResult(round)
		roundResult.RoundNumber = byte(index + 1)
		m.Rounds[index] = roundResult

		// MVP can be nil when the round ended because one team surrendered.
		mvp := m.getPlayer(round.MVP)
		if mvp != nil {
			mvp.MVPs++
			roundResult.MVPPlayerID = mvp.SteamID
		}

		// Get winner
		winner := m.getTeam(round.Winner.StartedAs)
		winner.Wins++
		roundResult.WinnerTeamID = winner.TeamID

		// Pistol round wins.
		roundNumber := index + 1
		if roundNumber == 1 || roundNumber == 16 {
			winner.PistolRoundWins++
		}

		// Iterage damage and add that to the damage dealt by each player.
		for _, damage := range round.Damage {
			// Attacker might be null according to open issues of the demoparser.
			attacker := damage.Attacker
			if attacker != nil {
				player := m.getPlayer(attacker)
				player.DamageDealt += damage.HealthDamageTaken
			}
		}

		playerKills := make(map[*player.PlayerResult]byte)

		// Process in round function in order to calculate all round information like amount of kills / round.
		for index, kill := range round.Kills {
			// Create and add killresult
			killResult := CreateKillResult(kill)
			roundResult.Kills[index] = killResult

			// Victim may be null, if it was a bot.
			if kill.Victim != nil {
				victim := m.getPlayer(kill.Victim)
				victim.Deaths++

				// First death of each round is an attempted opening duel.
				if index == 0 {
					victim.OpeningDuelAttempts++
				}
			}

			// Killer may not be set if the player died e.g. through fall damage.
			if kill.Killer != nil {
				killer := m.getPlayer(kill.Killer)
				killer.Kills++

				// First kill of each round is an entry kill.
				if index == 0 {
					killer.OpeningDuelAttempts++
					killer.EntryKills++
				}

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

func (m *MatchResult) getPlayer(player *demoparser.Player) *player.PlayerResult {
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

func (m *Match) Validate() error {
	err := validate.Struct(m)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

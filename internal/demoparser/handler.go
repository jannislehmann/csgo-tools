package demoparser

import (
	"fmt"

	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"

	log "github.com/sirupsen/logrus"
)

// Inits the players and teams.
func (p *DemoParser) handleMatchStart(events.MatchStart) {
	p.Match.Map = p.Match.Header.MapName
	p.IsFirstHalf = true

	if ConfigData.IsDebug() {
		p.debug(fmt.Sprintf("Game started on map %v", p.Match.Map))
	}

	gameState := p.parser.GameState()

	// Create teams.
	ct := gameState.TeamCounterTerrorists()
	t := gameState.TeamTerrorists()

	p.Match.Teams = make(map[int]*(Team), 2)
	teams := p.Match.Teams

	teams[ct.ID()] = &Team{State: ct, StartedAs: common.TeamCounterTerrorists}
	teams[t.ID()] = &Team{State: t, StartedAs: common.TeamTerrorists}

	// Create players and map them to the teams.
	for _, player := range gameState.Participants().Playing() {
		if player.IsBot {
			continue
		}

		teamID := player.TeamState.ID()
		teamPlayers := p.Match.Teams[teamID].Players

		customPlayer := &Player{SteamID: player.SteamID64, Name: player.Name, Team: teams[teamID]}

		p.Match.Players = append(p.Match.Players, customPlayer)
		teamPlayers = append(teamPlayers, customPlayer)
	}
}

func (p *DemoParser) handleGamePhaseChanged(e events.GamePhaseChanged) {
	switch e.NewGamePhase {
	case common.GamePhaseInit:
		p.IsFirstHalf = true
	case common.GamePhaseTeamSideSwitch:
		p.IsFirstHalf = false
	case common.GamePhaseGameEnded:
		p.Match.Duration = p.parser.CurrentTime()
	}
}

func (p *DemoParser) handleRoundStart(e events.RoundStart) {
	p.CurrentRound++
	p.RoundOngoing = true
	p.RoundStart = p.parser.CurrentTime()
	p.Match.Rounds = append(p.Match.Rounds, &Round{})

	if ConfigData.IsDebug() {
		p.debug(fmt.Sprintf("Starting round %d", p.CurrentRound))
	}
}

func (p *DemoParser) handleMVP(e events.RoundMVPAnnouncement) {
	player, err := p.getPlayer(e.Player)
	if err != nil {
		log.Panic(err)
	}

	if ConfigData.IsDebug() {
		p.debug(fmt.Sprintf("MVP for round %d is %v", p.CurrentRound, player.Name))
	}

	p.Match.Rounds[p.CurrentRound-1].MVP = player
}

func (p *DemoParser) handleRoundEnd(e events.RoundEnd) {
	p.RoundOngoing = false
	round := p.Match.Rounds[p.CurrentRound-1]

	if ConfigData.IsDebug() {
		p.debug(fmt.Sprintf("Ending round %d with winner %v", p.CurrentRound, e.Message))
	}

	round.Winner = p.Match.Teams[e.WinnerState.ID()]
	round.WinReason = e.Reason
	round.Duration = p.RoundStart - p.parser.CurrentTime()
}

func (p *DemoParser) handleKill(e events.Kill) {
	// Ignore warm-up kills
	if p.CurrentRound == 0 {
		return
	}

	killer, err := p.getPlayer(e.Killer)
	if err != nil {
		log.Panic(err)
	}

	victim, err := p.getPlayer(e.Victim)
	if err != nil {
		log.Panic(err)
	}

	round := p.Match.Rounds[p.CurrentRound-1]
	kill := &Kill{Time: p.parser.CurrentTime(), Weapon: e.Weapon.Type, IsHeadshot: e.IsHeadshot, Killer: killer, Victim: victim}
	round.Kills = append(round.Kills, kill)

	// Add optional assister
	if e.Assister != nil {
		assister, err := p.getPlayer(e.Assister)
		if err != nil {
			log.Panic(err)
		}
		kill.Assister = assister
	}
}

func (p *DemoParser) debug(message string) {
	if ConfigData.IsTrace() {
		log.WithFields(log.Fields{
			"Match": p.Match.ID,
			"Round": p.CurrentRound,
		}).Trace(message)
	} else {
		log.Debug(message)
	}
}

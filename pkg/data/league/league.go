package league

import (
	"git.sr.ht/~patrickpichler/wuzzlmoasta/pkg/data/game"
	"github.com/google/uuid"
)

type LeagueField func(*League)

func SetLeagueState(state LeagueState) LeagueField {
	return func(l *League) {
		l.State = state
	}
}

func SetLeagueName(name string) LeagueField {
	return func(l *League) {
		l.Name = name
	}
}

func SetLeagueDescription(description string) LeagueField {
	return func(l *League) {
		l.Description = description
	}
}

func AddPlayer(player uuid.UUID) LeagueField {
	return func(l *League) {
		l.Players = append(l.Players, player)
	}
}

func RemovePlayer(player uuid.UUID) LeagueField {
	return func(l *League) {
		for i, u := range l.Players {
			if u == player {
				l.Players[i] = l.Players[len(l.Players)-1]
				l.Players = l.Players[:len(l.Players)-1]
			}
		}
	}
}

func SetPlayers(players ...uuid.UUID) LeagueField {
	return func(l *League) {
		l.Players = players
	}
}

type League struct {
	Id          uuid.UUID
	State       LeagueState
	Name        string
	Description string
	Games       []game.Game
	Players     []uuid.UUID
}

type LeagueState int8

const (
	Planned LeagueState = 0
	Running             = 1
	Done                = 2
)

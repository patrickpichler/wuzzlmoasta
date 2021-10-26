package game

import (
	"time"

	"github.com/google/uuid"
)

type GameField func(*Game)

func Date(date time.Time) GameField {
	return func(g *Game) {
		g.Date = date
	}
}

func State(state GameState) GameField {
	return func(g *Game) {
		g.State = state
	}
}

func Team1(team1 Team) GameField {
	return func(g *Game) {
		g.Team1 = team1
	}
}

func Team2(team2 Team) GameField {
	return func(g *Game) {
		g.Team2 = team2
	}
}

func Score(score GameScore) GameField {
	return func(g *Game) {
		g.Score = score
	}
}

type Game struct {
	Id       uuid.UUID
	LeagueId uuid.UUID
	Date     time.Time
	State    GameState
	Team1    Team
	Team2    Team
	Score    GameScore
}

type GameScore struct {
	Team1 uint
	Team2 uint
}

type GameState int8

const (
	Scheduled GameState = 0
	Done                = 1
	Aborted             = 2
	Canceled            = 3
)

type Team struct {
	Players []uuid.UUID
}

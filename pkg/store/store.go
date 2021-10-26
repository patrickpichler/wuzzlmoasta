package store

import (
	"errors"
	"time"

	"git.sr.ht/~patrickpichler/wuzzlmoasta/pkg/data/game"
	"git.sr.ht/~patrickpichler/wuzzlmoasta/pkg/data/league"
	"github.com/google/uuid"
)

var (
	LeagueNotFound = errors.New("league not found")
	GameNotFound   = errors.New("game not found")
)

type Store interface {
	StartNewLeague(name, description string, players []uuid.UUID) (league.League, error)
	ListLeagues() ([]league.League, error)
	GetLeagueById(id uuid.UUID) (league.League, error)
	UpdateLeague(id uuid.UUID, fields ...league.LeagueField) (league.League, error)
	ListGamesInLeague(id uuid.UUID) ([]game.Game, error)
	AddDoneGameToLeague(leagueId uuid.UUID, date time.Time, team1, team2 game.Team, score game.GameScore) (game.Game, error)
	UpdateGame(id uuid.UUID, fields ...game.GameField) error
	ListGames() ([]game.Game, error)
}

func New() Store {
	return &inMemoryStore{}
}

package store

import (
	"time"

	"git.sr.ht/~patrickpichler/wuzzlmoasta/pkg/data/game"
	"git.sr.ht/~patrickpichler/wuzzlmoasta/pkg/data/league"
	"github.com/google/uuid"
)

type inMemoryStore struct {
	leagues []league.League
	games   []game.Game
}

var (
	_ Store = (*inMemoryStore)(nil)
)

func (s *inMemoryStore) StartNewLeague(name, description string, players []uuid.UUID) (league.League, error) {
	league := league.League{
		Id:          uuid.New(),
		State:       league.Running,
		Name:        name,
		Description: description,
		Players:     append([]uuid.UUID(nil), players...),
	}

	s.leagues = append(s.leagues, league)

	return league, nil
}

func (s *inMemoryStore) ListLeagues() ([]league.League, error) {
	result := make([]league.League, len(s.leagues))

	copy(result, s.leagues)

	return result, nil
}

func (s *inMemoryStore) GetLeagueById(id uuid.UUID) (league.League, error) {
	for _, l := range s.leagues {
		if l.Id == id {
			return l, nil
		}
	}

	return league.League{}, LeagueNotFound
}

func (s *inMemoryStore) UpdateLeague(id uuid.UUID, fields ...league.LeagueField) (league.League, error) {
	index := -1

	for i, l := range s.leagues {
		if l.Id == id {
			index = i
			break
		}
	}

	if index >= 0 {
		for _, f := range fields {
			f(&s.leagues[index])
		}
	}

	return league.League{}, LeagueNotFound
}

func (s *inMemoryStore) ListGamesInLeague(id uuid.UUID) ([]game.Game, error) {
	league, err := s.GetLeagueById(id)

	if err != nil {
		return nil, err
	}

	var result []game.Game

	for _, g := range s.games {
		if g.LeagueId == league.Id {
			result = append(result, g)
		}
	}

	return result, nil
}

func (s *inMemoryStore) AddDoneGameToLeague(leagueId uuid.UUID, date time.Time, team1, team2 game.Team, score game.GameScore) (game.Game, error) {
	g := game.Game{
		Id:       uuid.New(),
		LeagueId: leagueId,
		Date:     date,
		State:    game.Done,
		Team1:    team1,
		Team2:    team2,
		Score:    score,
	}

	s.games = append(s.games, g)

	return g, nil
}

func (s *inMemoryStore) UpdateGame(id uuid.UUID, fields ...game.GameField) error {
	for i, g := range s.games {
		if g.Id == id {
			for _, f := range fields {
				f(&s.games[i])
			}

			return nil
		}
	}

	return GameNotFound
}

func (s *inMemoryStore) ListGames() ([]game.Game, error) {
	result := make([]game.Game, len(s.games))

	copy(result, s.games)

	return result, nil
}

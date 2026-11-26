// Package memory keeps the game's data in memory.
package memory

import (
	"context"
	"slices"
	"sync"

	"github.com/gcalin/whomstdle/internal/game/domain"
)

type gameKey struct {
	player domain.PlayerID
	puzzle domain.Puzzle
}

// Games keeps every game in memory. Safe for concurrent use.
type Games struct {
	mu    sync.Mutex // handlers run concurrently
	games map[gameKey]domain.Game
}

func NewGames() *Games { return &Games{games: map[gameKey]domain.Game{}} }

// Get returns the player's game for puzzle p; ok is false when there is none yet. It never stores anything.
func (s *Games) Get(ctx context.Context, player domain.PlayerID, p domain.Puzzle) (domain.Game, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	g, ok := s.games[gameKey{player, p}]
	return clone(g), ok, nil
}

// Update follows the contract of app.GameStore.Update (see "The design → Interfaces").
func (s *Games) Update(ctx context.Context, player domain.PlayerID, p domain.Puzzle, change func(*domain.Game) error) (domain.Game, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := gameKey{player, p}
	g := clone(s.games[k])
	if err := change(&g); err != nil {
		return domain.Game{}, err
	}
	s.games[k] = g
	return clone(g), nil
}

// Len returns how many games are stored. Tests use it.
func (s *Games) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.games)
}

// clone returns a copy of g that shares no slice with it.
func clone(g domain.Game) domain.Game {
	g.Asked, g.Tried = slices.Clone(g.Asked), slices.Clone(g.Tried)
	return g
}

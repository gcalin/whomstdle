package app

import (
	"context"

	"github.com/gcalin/whomstdle/internal/game/domain"
)

// GameStore keeps every player's games.
type GameStore interface {
	// Get returns the player's game for puzzle p; ok is false when there is none yet.
	Get(ctx context.Context, player domain.PlayerID, p domain.Puzzle) (g domain.Game, ok bool, err error)
	// Update runs change on the player's game for puzzle p (a new, empty game when there is none yet),
	// as one atomic step. If change returns nil, the result is stored and a copy of it is returned.
	// If change returns an error, nothing is stored (not even a new game) and that error is returned.
	Update(ctx context.Context, player domain.PlayerID, p domain.Puzzle, change func(*domain.Game) error) (domain.Game, error)
}

// Catalog knows the people and which person each puzzle uses.
type Catalog interface {
	Person(ctx context.Context, id domain.PersonID) (p domain.Person, ok bool, err error)
	People(ctx context.Context) ([]domain.Person, error)
	PersonFor(ctx context.Context, p domain.Puzzle) (domain.Person, error)
}

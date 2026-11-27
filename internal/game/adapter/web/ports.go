package web

import (
	"context"

	"github.com/gcalin/whomstdle/internal/game/app"
	"github.com/gcalin/whomstdle/internal/game/domain"
)

// Game is everything a player can do. app.Service implements it.
type Game interface {
	Today(ctx context.Context, player domain.PlayerID) (app.Board, error)
	Ask(ctx context.Context, player domain.PlayerID, p domain.Puzzle, q domain.Question) (app.Board, error)
	Guess(ctx context.Context, player domain.PlayerID, p domain.Puzzle, id domain.PersonID) (app.Board, error)
	GiveUp(ctx context.Context, player domain.PlayerID, p domain.Puzzle) (app.Board, error)
	Choices(ctx context.Context) (app.Choices, error)
}

// Package app runs the game's use cases: show today, ask, guess, give up.
package app

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gcalin/whomstdle/internal/game/domain"
)

// Service runs the game's use cases.
type Service struct {
	games   GameStore
	catalog Catalog
	now     func() time.Time
}

func New(games GameStore, catalog Catalog, now func() time.Time) *Service {
	return &Service{games: games, catalog: catalog, now: now}
}

// Today returns the player's board for today's puzzle. It never stores anything.
func (s *Service) Today(ctx context.Context, player domain.PlayerID) (Board, error) {
	p := domain.Current(s.now())
	g, _, err := s.games.Get(ctx, player, p)
	if err != nil {
		return Board{}, fmt.Errorf("today: %w", err)
	}
	b, err := s.board(ctx, p, g)
	if err != nil {
		return Board{}, fmt.Errorf("today: %w", err)
	}
	return b, nil
}

func (s *Service) Ask(ctx context.Context, player domain.PlayerID, p domain.Puzzle, q domain.Question) (Board, error) {
	if !q.Valid() {
		return Board{}, ErrUnknownQuestion
	}
	if !domain.Open(p, s.now()) {
		return Board{}, ErrClosed
	}
	return s.update(ctx, "ask", player, p, func(g *domain.Game) error { g.Ask(q); return nil })
}

func (s *Service) Guess(ctx context.Context, player domain.PlayerID, p domain.Puzzle, id domain.PersonID) (Board, error) {
	_, ok, err := s.catalog.Person(ctx, id)
	if err != nil {
		return Board{}, fmt.Errorf("guess: %w", err)
	}
	if !ok {
		return Board{}, ErrUnknownPerson
	}
	if !domain.Open(p, s.now()) {
		return Board{}, ErrClosed
	}
	answer, err := s.catalog.PersonFor(ctx, p)
	if err != nil {
		return Board{}, fmt.Errorf("guess: %w", err)
	}
	return s.update(ctx, "guess", player, p, func(g *domain.Game) error { g.Guess(id, answer.ID); return nil })
}

func (s *Service) GiveUp(ctx context.Context, player domain.PlayerID, p domain.Puzzle) (Board, error) {
	if !domain.Open(p, s.now()) {
		return Board{}, ErrClosed
	}
	return s.update(ctx, "give up", player, p, func(g *domain.Game) error { g.GiveUp(); return nil })
}

// Choices returns every question, and every person sorted by name.
func (s *Service) Choices(ctx context.Context) (Choices, error) {
	people, err := s.catalog.People(ctx)
	if err != nil {
		return Choices{}, fmt.Errorf("choices: %w", err)
	}
	c := Choices{Questions: slices.Clone(domain.Questions)}
	for _, p := range people {
		c.People = append(c.People, personName(p))
	}
	slices.SortFunc(c.People, func(a, b PersonName) int { return strings.Compare(a.Name, b.Name) })
	return c, nil
}

// update runs change on the player's game for puzzle p and returns the new board.
// Errors are wrapped with op, the name of the action.
func (s *Service) update(ctx context.Context, op string, player domain.PlayerID, p domain.Puzzle, change func(*domain.Game) error) (Board, error) {
	g, err := s.games.Update(ctx, player, p, change)
	if err != nil {
		return Board{}, fmt.Errorf("%s: %w", op, err)
	}
	b, err := s.board(ctx, p, g)
	if err != nil {
		return Board{}, fmt.Errorf("%s: %w", op, err)
	}
	return b, nil
}

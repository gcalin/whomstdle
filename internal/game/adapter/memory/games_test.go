package memory

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"

	"github.com/gcalin/whomstdle/internal/game/domain"
)

var ctx = context.Background()

func ask(q domain.Question) func(*domain.Game) error {
	return func(g *domain.Game) error { g.Ask(q); return nil }
}

func TestGetEmpty(t *testing.T) {
	s := NewGames()
	if _, ok, _ := s.Get(ctx, "a", 1); ok || s.Len() != 0 {
		t.Fatalf("Get on an empty store: ok = %v, Len = %d, want false and 0", ok, s.Len())
	}
}

func TestUpdateThenGet(t *testing.T) {
	s := NewGames()
	if _, err := s.Update(ctx, "a", 1, ask(domain.Alive)); err != nil {
		t.Fatal(err)
	}
	g, ok, _ := s.Get(ctx, "a", 1)
	if !ok || !slices.Equal(g.Asked, []domain.Question{domain.Alive}) {
		t.Fatalf("Get after Update: ok = %v, Asked = %v", ok, g.Asked)
	}
}

func TestFailedChangeStoresNothing(t *testing.T) {
	s := NewGames()
	boom := errors.New("boom")
	fail := func(g *domain.Game) error { g.Ask(domain.Gender); return boom }
	if _, err := s.Update(ctx, "a", 1, fail); !errors.Is(err, boom) || s.Len() != 0 {
		t.Fatalf("a failed change on a new game: err = %v, Len = %d, want boom and 0", err, s.Len())
	}
	if _, err := s.Update(ctx, "a", 1, ask(domain.Alive)); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(ctx, "a", 1, fail); !errors.Is(err, boom) {
		t.Fatalf("err = %v, want boom", err)
	}
	if g, _, _ := s.Get(ctx, "a", 1); !slices.Equal(g.Asked, []domain.Question{domain.Alive}) {
		t.Fatalf("a failed change on an existing game changed it: Asked = %v", g.Asked)
	}
}

func TestCopiesDoNotLeak(t *testing.T) {
	s := NewGames()
	g, _ := s.Update(ctx, "a", 1, ask(domain.Alive))
	g.Asked = append(g.Asked, domain.Gender)
	g.Asked[0] = domain.Hint
	if got, _, _ := s.Get(ctx, "a", 1); !slices.Equal(got.Asked, []domain.Question{domain.Alive}) {
		t.Fatalf("changing a returned copy changed the stored game: Asked = %v", got.Asked)
	}
}

func TestConcurrentUpdates(t *testing.T) {
	s := NewGames()
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Go(func() {
			_, _ = s.Update(ctx, "a", 1, func(g *domain.Game) error { g.Guess(domain.PersonID(100+i), 1); return nil })
		})
	}
	wg.Wait()
	if g, _, _ := s.Get(ctx, "a", 1); len(g.Tried) != 20 {
		t.Fatalf("got %d wrong guesses, want 20", len(g.Tried))
	}
}

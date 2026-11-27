package app

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/gcalin/whomstdle/internal/game/domain"
)

var ctx = context.Background()

type gameKey struct {
	player domain.PlayerID
	puzzle domain.Puzzle
}

// fakeStore keeps games in a map and counts Update calls. When err is set, every call returns it.
type fakeStore struct {
	games   map[gameKey]domain.Game
	updates int
	err     error
}

func (f *fakeStore) Get(_ context.Context, player domain.PlayerID, p domain.Puzzle) (domain.Game, bool, error) {
	if f.err != nil {
		return domain.Game{}, false, f.err
	}
	g, ok := f.games[gameKey{player, p}]
	return g, ok, nil
}

func (f *fakeStore) Update(_ context.Context, player domain.PlayerID, p domain.Puzzle, change func(*domain.Game) error) (domain.Game, error) {
	f.updates++
	if f.err != nil {
		return domain.Game{}, f.err
	}
	g := f.games[gameKey{player, p}]
	g.Asked, g.Tried = slices.Clone(g.Asked), slices.Clone(g.Tried)
	if err := change(&g); err != nil {
		return domain.Game{}, err
	}
	f.games[gameKey{player, p}] = g
	return g, nil
}

// answers gives person every question with the same answer.
func answers(a string) map[domain.Question]string {
	m := map[domain.Question]string{}
	for _, q := range domain.Questions {
		m[q] = a
	}
	return m
}

// fakeCatalog has two people. Puzzle 1 uses person 1 and puzzle 2 uses person 2.
type fakeCatalog struct{}

var people = []domain.Person{
	{ID: 1, Name: "Zed", Answers: answers("one")},
	{ID: 2, Name: "Amy", Answers: answers("two")},
}

func (fakeCatalog) Person(_ context.Context, id domain.PersonID) (domain.Person, bool, error) {
	for _, p := range people {
		if p.ID == id {
			return p, true, nil
		}
	}
	return domain.Person{}, false, nil
}

func (fakeCatalog) People(context.Context) ([]domain.Person, error) {
	return slices.Clone(people), nil
}

func (fakeCatalog) PersonFor(_ context.Context, p domain.Puzzle) (domain.Person, error) {
	return people[(p-1)%2], nil
}

// newService returns a Service with fakes, and a clock that reads *at.
func newService(at *time.Time) (*Service, *fakeStore) {
	store := &fakeStore{games: map[gameKey]domain.Game{}}
	return New(store, fakeCatalog{}, func() time.Time { return *at }), store
}

var day1 = domain.Start(1).Add(12 * time.Hour)

func TestNameHiddenWhilePlaying(t *testing.T) {
	at := day1
	s, _ := newService(&at)
	for _, q := range domain.Questions {
		b, err := s.Ask(ctx, "a", 1, q)
		if err != nil || b.Name != "" {
			t.Fatalf("Ask(%d): Name = %q, err = %v, want empty and nil", q, b.Name, err)
		}
	}
	b, _ := s.Guess(ctx, "a", 1, 1)
	if b.Name != "Zed" || len(b.Answers) != 10 {
		t.Errorf("after solving: Name = %q with %d answers, want Zed and 10", b.Name, len(b.Answers))
	}
	b, _ = s.GiveUp(ctx, "b", 1)
	if b.Name != "Zed" || len(b.Answers) != 10 {
		t.Errorf("after giving up: Name = %q with %d answers, want Zed and 10", b.Name, len(b.Answers))
	}
}

func TestAnswersInOrderAsked(t *testing.T) {
	at := day1
	s, _ := newService(&at)
	if _, err := s.Ask(ctx, "a", 1, domain.Field); err != nil {
		t.Fatal(err)
	}
	b, _ := s.Ask(ctx, "a", 1, domain.Alive)
	want := []Answer{{domain.Field, "one"}, {domain.Alive, "one"}}
	if !slices.Equal(b.Answers, want) || b.Asked != 2 {
		t.Fatalf("Answers = %v, Asked = %d, want %v and 2", b.Answers, b.Asked, want)
	}
	b, _ = s.GiveUp(ctx, "a", 1)
	if b.Asked != 2 || b.Answers[0].Question != domain.Questions[0] {
		t.Errorf("after the end: Asked = %d, first answer %v, want 2 and all questions in order", b.Asked, b.Answers[0])
	}
}

func TestSecondsLeft(t *testing.T) {
	at := day1
	s, _ := newService(&at)
	b, _ := s.Ask(ctx, "a", 1, domain.Alive)
	if b.SecondsLeft != 0 {
		t.Errorf("playing: SecondsLeft = %d, want 0", b.SecondsLeft)
	}
	b, _ = s.GiveUp(ctx, "a", 1)
	if b.SecondsLeft != 43200 {
		t.Errorf("finished: SecondsLeft = %d, want 43200", b.SecondsLeft)
	}
}

func TestHintNotCounted(t *testing.T) {
	at := day1
	s, _ := newService(&at)
	var b Board
	for _, q := range []domain.Question{domain.Alive, domain.Gender, domain.Field, domain.Continent, domain.BirthDecade, domain.Hint} {
		b, _ = s.Ask(ctx, "a", 1, q)
	}
	if !b.HintUsed || b.Asked != 5 || len(b.Answers) != 6 {
		t.Errorf("HintUsed = %v, Asked = %d, %d answers, want true, 5 and 6", b.HintUsed, b.Asked, len(b.Answers))
	}
}

func TestTodayStoresNothing(t *testing.T) {
	at := day1
	s, store := newService(&at)
	if b, err := s.Today(ctx, "a"); err != nil || b.Puzzle != 1 || store.updates != 0 {
		t.Fatalf("Today: Puzzle = %d, err = %v, %d updates, want 1, nil, 0", b.Puzzle, err, store.updates)
	}
}

func TestRefused(t *testing.T) {
	at := day1
	s, store := newService(&at)
	if _, err := s.Ask(ctx, "a", 1, 10); !errors.Is(err, ErrUnknownQuestion) {
		t.Errorf("unknown question: err = %v", err)
	}
	if _, err := s.Guess(ctx, "a", 1, 7); !errors.Is(err, ErrUnknownPerson) {
		t.Errorf("unknown person: err = %v", err)
	}
	if _, err := s.Ask(ctx, "a", 2, domain.Alive); !errors.Is(err, ErrClosed) {
		t.Errorf("Ask on a closed puzzle: err = %v", err)
	}
	if _, err := s.Guess(ctx, "a", 2, 1); !errors.Is(err, ErrClosed) {
		t.Errorf("Guess on a closed puzzle: err = %v", err)
	}
	if _, err := s.GiveUp(ctx, "a", 2); !errors.Is(err, ErrClosed) {
		t.Errorf("GiveUp on a closed puzzle: err = %v", err)
	}
	if store.updates != 0 {
		t.Errorf("refused requests called Update %d times", store.updates)
	}
}

func TestMidnight(t *testing.T) {
	at := domain.Start(2).Add(-2 * time.Minute)
	s, _ := newService(&at)
	if _, err := s.Ask(ctx, "a", 1, domain.Alive); err != nil {
		t.Fatal(err)
	}
	at = domain.Start(2).Add(3 * time.Minute)
	b, err := s.Ask(ctx, "a", 1, domain.Gender)
	if err != nil || b.Answers[1].Text != "one" {
		t.Fatalf("after midnight, puzzle 1: answers %v, err = %v, want person 1's", b.Answers, err)
	}
	at = domain.Start(2).Add(61 * time.Minute)
	if _, err := s.Ask(ctx, "a", 1, domain.Field); !errors.Is(err, ErrClosed) {
		t.Errorf("after the grace hour: err = %v, want ErrClosed", err)
	}
}

func TestChoicesSorted(t *testing.T) {
	at := day1
	s, _ := newService(&at)
	c, err := s.Choices(ctx)
	if err != nil || len(c.People) != 2 || c.People[0].Name != "Amy" || c.People[1].Name != "Zed" {
		t.Fatalf("Choices.People = %v, err = %v, want Amy then Zed", c.People, err)
	}
	if !slices.Equal(c.Questions, domain.Questions) {
		t.Errorf("Choices.Questions = %v", c.Questions)
	}
}

func TestStoreErrors(t *testing.T) {
	at := day1
	s, store := newService(&at)
	store.err = errors.New("boom")
	_, askErr := s.Ask(ctx, "a", 1, domain.Alive)
	_, guessErr := s.Guess(ctx, "a", 1, 1)
	_, giveUpErr := s.GiveUp(ctx, "a", 1)
	_, todayErr := s.Today(ctx, "a")
	for name, err := range map[string]error{"Ask": askErr, "Guess": guessErr, "GiveUp": giveUpErr, "Today": todayErr} {
		if !errors.Is(err, store.err) {
			t.Errorf("%s: err = %v, want the store's error", name, err)
		}
	}
}

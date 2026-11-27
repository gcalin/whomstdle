package app

import (
	"context"
	"time"

	"github.com/gcalin/whomstdle/internal/game/domain"
)

// Board is one player's game for one puzzle, as the player may see it.
type Board struct {
	Puzzle   domain.Puzzle
	Status   domain.Status
	Asked    int // how many questions the player asked, not counting the hint
	HintUsed bool
	Answers  []Answer     // playing: only the questions asked, in the order asked. Finished: all 10, in question order.
	Tried    []PersonName // wrong guesses, in the order made
	Score    int
	Name     string // the person's name; empty while playing (game rule 6)
	// SecondsLeft is the time until the next puzzle starts, set only when finished.
	// It is 0 or less for a game finished in the grace hour: the next puzzle has already started.
	SecondsLeft int
}

// Answer is one question and the person's answer to it.
type Answer struct {
	Question domain.Question
	Text     string
}

// PersonName is a person without their answers: safe to send to the page.
type PersonName struct {
	ID      domain.PersonID
	Name    string
	Aliases []string
}

// Choices is what the page's search boxes offer.
type Choices struct {
	Questions []domain.Question
	People    []PersonName // sorted by name: the order must not hint at the answer
}

func personName(p domain.Person) PersonName {
	return PersonName{ID: p.ID, Name: p.Name, Aliases: p.Aliases}
}

// board builds what the player may see of game g for puzzle p.
func (s *Service) board(ctx context.Context, p domain.Puzzle, g domain.Game) (Board, error) {
	person, err := s.catalog.PersonFor(ctx, p)
	if err != nil {
		return Board{}, err
	}
	b := Board{Puzzle: p, Status: g.Status, Asked: len(g.Asked), HintUsed: g.HintUsed(), Score: g.Score()}
	if b.HintUsed {
		b.Asked--
	}
	for _, id := range g.Tried {
		tried, _, err := s.catalog.Person(ctx, id) // a tried ID is always a known person
		if err != nil {
			return Board{}, err
		}
		b.Tried = append(b.Tried, personName(tried))
	}
	shown := g.Asked
	if g.Finished() {
		shown = domain.Questions
		b.Name = person.Name
		b.SecondsLeft = int(domain.Start(p+1).Sub(s.now()) / time.Second)
	}
	for _, q := range shown {
		b.Answers = append(b.Answers, Answer{Question: q, Text: person.Answer(q)})
	}
	return b, nil
}

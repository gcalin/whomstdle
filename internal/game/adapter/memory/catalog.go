package memory

import (
	"context"
	"slices"

	"github.com/gcalin/whomstdle/internal/game/domain"
)

// Catalog looks up people and the person for each puzzle.
type Catalog struct {
	people   []domain.Person
	schedule []domain.PersonID // schedule[n-1] is the person for puzzle n
}

// NewCatalog returns the built-in people and schedule.
func NewCatalog() *Catalog { return &Catalog{people: people, schedule: schedule} }

func (c *Catalog) Person(ctx context.Context, id domain.PersonID) (domain.Person, bool, error) {
	for _, p := range c.people {
		if p.ID == id {
			return p, true, nil
		}
	}
	return domain.Person{}, false, nil
}

// People returns a copy of the list.
func (c *Catalog) People(ctx context.Context) ([]domain.Person, error) {
	return slices.Clone(c.people), nil
}

func (c *Catalog) PersonFor(ctx context.Context, p domain.Puzzle) (domain.Person, error) {
	id := c.schedule[int(p-1)%len(c.schedule)] // wraps if the schedule runs out; a test warns a week before
	person, _, _ := c.Person(ctx, id)
	return person, nil
}

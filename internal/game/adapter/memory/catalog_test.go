package memory

import (
	"context"
	"testing"
	"time"

	"github.com/gcalin/whomstdle/internal/game/domain"
)

func TestIDs(t *testing.T) {
	seen := map[domain.PersonID]bool{}
	for _, p := range people {
		if p.ID < 1 || seen[p.ID] {
			t.Errorf("%s: ID %d is below 1 or used twice", p.Name, p.ID)
		}
		seen[p.ID] = true
	}
}

func TestSchedule(t *testing.T) {
	if schedule[0] != 1 {
		t.Errorf("puzzle 1 has started and must stay ID 1, got %d", schedule[0])
	}
	c := NewCatalog()
	for i, id := range schedule {
		if _, ok, _ := c.Person(context.Background(), id); !ok {
			t.Errorf("puzzle %d: no person with ID %d", i+1, id)
		}
		if i > 0 && schedule[i-1] == id {
			t.Errorf("puzzles %d and %d have the same person", i, i+1)
		}
	}
	// Uses the real clock on purpose: fails a week before the schedule runs out.
	if left := len(schedule) - int(domain.Current(time.Now())); left < 7 {
		t.Errorf("only %d puzzles left: the schedule runs out after %s, add entries to schedule",
			left, domain.Start(domain.Puzzle(len(schedule))).Format("2006-01-02"))
	}
}

// A missing answer would show the player an empty reply, so every person must answer every question.
func TestEveryPersonAnswersEveryQuestion(t *testing.T) {
	for _, p := range people {
		for _, q := range domain.Questions {
			if p.Answer(q) == "" {
				t.Errorf("%s has no answer to %q", p.Name, q.Text())
			}
		}
	}
}

func TestPersonForIgnoresOrder(t *testing.T) {
	ctx := context.Background()
	before, _ := NewCatalog().PersonFor(ctx, 1)
	c := &Catalog{people: append([]domain.Person{{ID: 99, Name: "New"}}, people...), schedule: schedule}
	if got, _ := c.PersonFor(ctx, 1); got.Name != before.Name {
		t.Errorf("adding a person changed puzzle 1 from %s to %s", before.Name, got.Name)
	}
}

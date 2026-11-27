package domain

import (
	"slices"
	"testing"
)

func TestAsk(t *testing.T) {
	var g Game
	g.Ask(Alive)
	g.Ask(Gender)
	g.Ask(Alive)
	if want := []Question{Alive, Gender}; !slices.Equal(g.Asked, want) {
		t.Errorf("Asked = %v, want %v: a repeated question changes nothing", g.Asked, want)
	}
}

func TestAskHintLocked(t *testing.T) {
	var g Game
	for _, q := range []Question{Alive, Gender, Field, Continent} {
		g.Ask(q)
	}
	g.Ask(Hint)
	if g.HintUsed() || len(g.Asked) != 4 {
		t.Errorf("hint after 4 questions: Asked = %v, want the hint ignored", g.Asked)
	}
	g.Ask(BirthDecade)
	g.Ask(Hint)
	if !g.HintUsed() || len(g.Asked) != 6 {
		t.Errorf("hint after 5 questions: Asked = %v, want the hint added", g.Asked)
	}
}

func TestGuess(t *testing.T) {
	var g Game
	g.Guess(2, 1)
	g.Guess(2, 1)
	if want := []PersonID{2}; !slices.Equal(g.Tried, want) || g.Finished() {
		t.Errorf("after the same wrong guess twice: Tried = %v, Status = %v, want %v and Playing", g.Tried, g.Status, want)
	}
	g.Guess(1, 1)
	if g.Status != Solved {
		t.Errorf("a right guess: Status = %v, want Solved", g.Status)
	}

	var first Game
	first.Guess(1, 1)
	if first.Status != Solved || first.Score() != 0 {
		t.Errorf("a right first guess: Status = %v, Score = %d, want Solved and 0", first.Status, first.Score())
	}
}

func TestGiveUp(t *testing.T) {
	var g Game
	g.GiveUp()
	if g.Status != GaveUp {
		t.Errorf("Status = %v, want GaveUp", g.Status)
	}
}

func TestFinished(t *testing.T) {
	for _, c := range []struct {
		s    Status
		want bool
	}{{Playing, false}, {Solved, true}, {GaveUp, true}} {
		if got := (Game{Status: c.s}).Finished(); got != c.want {
			t.Errorf("Finished() with Status %v = %v, want %v", c.s, got, c.want)
		}
	}
}

// After the end, no action changes the game.
func TestFinishedGameIgnoresActions(t *testing.T) {
	for _, s := range []Status{Solved, GaveUp} {
		before := Game{Asked: []Question{Alive}, Tried: []PersonID{2}, Status: s}
		for name, act := range map[string]func(*Game){
			"Ask":         func(g *Game) { g.Ask(Gender) },
			"wrong Guess": func(g *Game) { g.Guess(3, 1) },
			"right Guess": func(g *Game) { g.Guess(1, 1) },
			"GiveUp":      func(g *Game) { g.GiveUp() },
		} {
			g := Game{Asked: slices.Clone(before.Asked), Tried: slices.Clone(before.Tried), Status: s}
			act(&g)
			if !slices.Equal(g.Asked, before.Asked) || !slices.Equal(g.Tried, before.Tried) || g.Status != s {
				t.Errorf("%s after Status %v changed the game: %+v", name, s, g)
			}
		}
	}
}

func TestScore(t *testing.T) {
	for _, c := range []struct {
		name string
		g    Game
		want int
	}{
		{"new game", Game{}, 0},
		{"questions only", Game{Asked: []Question{Alive, Gender}}, 2},
		{"wrong guesses only", Game{Tried: []PersonID{2, 3}}, 2},
		{"questions and wrong guesses", Game{Asked: []Question{Alive}, Tried: []PersonID{2, 3}}, 3},
		{"hint costs 3", Game{Asked: []Question{Alive, Hint}}, 4},
		{"solved game counts the same", Game{Asked: []Question{Alive}, Tried: []PersonID{2}, Status: Solved}, 2},
	} {
		if got := c.g.Score(); got != c.want {
			t.Errorf("%s: Score() = %d, want %d", c.name, got, c.want)
		}
	}
}

func TestAnswer(t *testing.T) {
	p := Person{Answers: map[Question]string{Alive: "No"}}
	if got := p.Answer(Alive); got != "No" {
		t.Errorf("Answer(Alive) = %q, want %q", got, "No")
	}
	if got := p.Answer(Gender); got != "" {
		t.Errorf("Answer(Gender) with no answer stored = %q, want empty", got)
	}
}

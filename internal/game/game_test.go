package game

import "testing"

func TestScore(t *testing.T) {
	for _, c := range []struct {
		name string
		g    Game
		want int
	}{
		{"new game", Game{}, 0},
		{"questions only", Game{Asked: []Question{Alive, Gender}}, 2},
		{"wrong guesses only", Game{Tried: []int{2, 3}}, 2},
		{"questions and wrong guesses", Game{Asked: []Question{Alive}, Tried: []int{2, 3}}, 3},
		{"solved game counts the same", Game{Asked: []Question{Alive}, Tried: []int{2}, Status: Solved}, 2},
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

// A missing answer would show the player an empty reply, so every person must answer every question.
func TestEveryPersonAnswersEveryQuestion(t *testing.T) {
	for _, p := range People {
		for _, q := range Questions {
			if p.Answer(q) == "" {
				t.Errorf("%s has no answer to %q", p.Name, q.Text())
			}
		}
	}
}

package game

import "testing"

// fixture returns a small, known game used across the tests. Keeping the test
// data separate from cmd/game's real cast means these tests never break when
// you add or edit people.
func fixture() *Game {
	einstein := Person{
		Name:    "Albert Einstein",
		Aliases: []string{"Einstein"},
		Answers: map[QuestionID]string{
			"alive":  "no",
			"gender": "man",
		},
	}
	questions := []Question{
		{ID: "alive", Text: "Still alive?"},
		{ID: "gender", Text: "Man or woman?"},
	}
	return New(einstein, questions)
}

// From the plan: "asking about a dead person answers no to still alive?".
func TestAskDeadPersonIsNotAlive(t *testing.T) {
	g := fixture()
	got, ok := g.Ask("alive")
	if !ok {
		t.Fatal("Ask(alive): ok = false, want true")
	}
	if got != "no" {
		t.Errorf("Ask(alive) = %q, want %q", got, "no")
	}
}

func TestAskUnknownQuestion(t *testing.T) {
	g := fixture()
	if _, ok := g.Ask("height"); ok {
		t.Error("Ask(height): ok = true for an unknown question, want false")
	}
}

// From the plan: "a correct guess ends the game". Also checks that matching is
// case-insensitive and works on an alias.
func TestCorrectGuessEndsGame(t *testing.T) {
	g := fixture()
	if g.Over() {
		t.Fatal("a new game should not be over")
	}
	if !g.Guess("einstein") {
		t.Fatal(`Guess("einstein") = false, want true (alias, different case)`)
	}
	if !g.Over() {
		t.Error("game should be over after a correct guess")
	}
	if !g.Won() {
		t.Error("game should be won after a correct guess")
	}
}

func TestWrongGuessContinues(t *testing.T) {
	g := fixture()
	if g.Guess("Marie Curie") {
		t.Fatal(`Guess("Marie Curie") = true, want false`)
	}
	if g.Over() {
		t.Error("game should continue after a wrong guess")
	}
	if g.WrongGuesses() != 1 {
		t.Errorf("WrongGuesses() = %d, want 1", g.WrongGuesses())
	}
}

func TestAskCountsQuestions(t *testing.T) {
	g := fixture()
	g.Ask("alive")
	g.Ask("gender")
	if g.Asked() != 2 {
		t.Errorf("Asked() = %d, want 2", g.Asked())
	}
}

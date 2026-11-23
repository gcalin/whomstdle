// Package game holds the core logic of the daily celebrity guessing game, with
// no terminal, web, or database code. Phase 0 drives it from a terminal loop
// (see cmd/game); Phase 1 will wrap the same logic in net/http. Keeping the
// rules here, free of I/O, is what makes them easy to unit-test.
package game

import "strings"

// QuestionID identifies one of the fixed questions a player can ask.
type QuestionID string

// Question is one fixed question the player may ask about the secret person.
type Question struct {
	ID   QuestionID
	Text string
}

// Person is a celebrity the game can pick as the secret answer. Answers maps a
// QuestionID to this person's answer for it, e.g. "no", "woman", "Dutch".
type Person struct {
	Name    string
	Aliases []string
	Answers map[QuestionID]string
}

// matches reports whether guess names this person, compared case-insensitively
// against the canonical name and every alias.
func (p Person) matches(guess string) bool {
	guess = strings.TrimSpace(strings.ToLower(guess))
	if guess == "" {
		return false
	}
	if strings.ToLower(p.Name) == guess {
		return true
	}
	for _, alias := range p.Aliases {
		if strings.ToLower(alias) == guess {
			return true
		}
	}
	return false
}

// Game is a single play session against one secret person. The zero value is
// not usable; create one with New.
type Game struct {
	secret    Person
	questions []Question
	asked     int
	wrong     int
	over      bool
	won       bool
}

// New starts a game with the given secret person and the set of askable questions.
func New(secret Person, questions []Question) *Game {
	return &Game{secret: secret, questions: questions}
}

// Questions returns the questions the player may ask.
func (g *Game) Questions() []Question { return g.questions }

// Ask returns the secret person's answer to question id. ok is false when there
// is no such question. A question is counted only while the game is still going.
func (g *Game) Ask(id QuestionID) (answer string, ok bool) {
	answer, ok = g.secret.Answers[id]
	if !ok {
		return "", false
	}
	if !g.over {
		g.asked++
	}
	return answer, true
}

// Guess checks a guessed name. A correct guess wins and ends the game; a wrong
// guess is counted and the game continues. correct reports whether it was right.
func (g *Game) Guess(name string) (correct bool) {
	if g.over {
		return false
	}
	if g.secret.matches(name) {
		g.over = true
		g.won = true
		return true
	}
	g.wrong++
	return false
}

// Over reports whether the game has finished.
func (g *Game) Over() bool { return g.over }

// Won reports whether the game ended with a correct guess.
func (g *Game) Won() bool { return g.won }

// Asked returns how many questions have been asked.
func (g *Game) Asked() int { return g.asked }

// WrongGuesses returns how many incorrect guesses have been made.
func (g *Game) WrongGuesses() int { return g.wrong }

// Secret returns the secret person. Only call this once the game is Over — the
// terminal and (later) web layers must never reveal it mid-game.
func (g *Game) Secret() Person { return g.secret }

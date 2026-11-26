package game

import "slices"

// Status says whether a game is still being played, and how it ended.
type Status int

const (
	Playing Status = iota // the zero value, so a new game is being played
	Solved
	GaveUp
)

// Game is one player's game for one puzzle.
type Game struct {
	Asked  []Question // in the order asked, no repeats
	Tried  []PersonID // person IDs of wrong guesses, no repeats
	Status Status
}

// Score: 1 point per question and per wrong guess. Step 7 adds the hint cost.
func (g Game) Score() int { return len(g.Asked) + len(g.Tried) }

// Finished reports whether the game has ended.
func (g Game) Finished() bool { return g.Status != Playing }

// Ask records question q. It does nothing when the game is finished or q was already asked.
func (g *Game) Ask(q Question) {
	if !g.Finished() && !slices.Contains(g.Asked, q) {
		g.Asked = append(g.Asked, q)
	}
}

// Guess records a guess of person guess, when the answer is person answer.
// Right: the game is Solved. Wrong: guess is added to Tried once. Finished: nothing changes.
func (g *Game) Guess(guess, answer PersonID) {
	switch {
	case g.Finished():
	case guess == answer:
		g.Status = Solved
	case !slices.Contains(g.Tried, guess):
		g.Tried = append(g.Tried, guess)
	}
}

// GiveUp ends the game as GaveUp. It does nothing when the game is finished.
func (g *Game) GiveUp() {
	if !g.Finished() {
		g.Status = GaveUp
	}
}

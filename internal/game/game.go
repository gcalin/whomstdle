package game

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

package game

// Game is one player's game for one puzzle.
type Game struct {
	Asked []Question // in the order asked, no repeats
	Tried []int      // person IDs of wrong guesses, no repeats
}

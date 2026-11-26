package web

import (
	"slices"
	"sync"

	"github.com/gcalin/whomstdle/internal/game/domain"
)

// memStore keeps every game in memory.
type memStore struct {
	mu    sync.Mutex // handlers run concurrently
	games map[gameKey]*domain.Game
}

type gameKey struct {
	player domain.PlayerID
	puzzle domain.Puzzle
}

// update runs f on the player's game for the puzzle while holding the lock.
// If there is no game yet: create it first when create is true, otherwise return ok == false.
// It returns a copy of the game after f ran. f must not write the response.
func (s *memStore) update(player domain.PlayerID, puzzle domain.Puzzle, create bool, f func(*domain.Game)) (g domain.Game, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := gameKey{player, puzzle}
	stored := s.games[k]
	if stored == nil {
		if !create {
			return domain.Game{}, false
		}
		stored = &domain.Game{}
		s.games[k] = stored
	}
	f(stored)
	c := *stored
	c.Asked, c.Tried = slices.Clone(c.Asked), slices.Clone(c.Tried)
	return c, true
}

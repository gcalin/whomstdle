package memory

import "github.com/gcalin/whomstdle/internal/game/domain"

// schedule[n-1] is the person ID for puzzle n.
// Only add to the end. Never change an entry for a puzzle that has started.
// Do not schedule people in ID order, or the IDs in the page give away the answer.
var schedule = []domain.PersonID{1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 3, 2}

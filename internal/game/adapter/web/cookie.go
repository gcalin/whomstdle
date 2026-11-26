package web

import (
	"crypto/rand"
	"net/http"

	"github.com/gcalin/whomstdle/internal/game/domain"
)

// playerID returns the player cookie, and creates it when it is missing.
func playerID(w http.ResponseWriter, r *http.Request) domain.PlayerID {
	if c, err := r.Cookie("player"); err == nil {
		return domain.PlayerID(c.Value)
	}
	id := rand.Text()
	http.SetCookie(w, &http.Cookie{Name: "player", Value: id, Path: "/", MaxAge: 365 * 24 * 3600, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	return domain.PlayerID(id)
}

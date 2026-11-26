// Package web serves the game in the browser: html/template pages plus htmx fragments.
package web

import (
	"crypto/rand"
	"embed"
	"html/template"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gcalin/whomstdle/internal/game"
)

//go:embed templates/*.html
var templateFS embed.FS

var tmpl = template.Must(template.ParseFS(templateFS, "templates/*.html"))

// pageData is sent once inside the page; the browser searches it as the player types.
var pageData = func() map[string]any {
	var qs, people []map[string]any
	for _, q := range game.Questions {
		qs = append(qs, map[string]any{"id": int(q), "text": q.Text(), "aliases": q.Aliases()})
	}
	// Sorted by name: the order in game.People must not hint at the answer.
	sorted := slices.Clone(game.People)
	slices.SortFunc(sorted, func(a, b game.Person) int { return strings.Compare(a.Name, b.Name) })
	for _, p := range sorted {
		people = append(people, map[string]any{"id": p.ID, "name": p.Name, "aliases": p.Aliases})
	}
	return map[string]any{"questions": qs, "people": people}
}()

type server struct {
	now   func() time.Time
	store memStore
	mux   *http.ServeMux
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }

// playerID returns the player cookie, and creates it when it is missing.
func playerID(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie("player"); err == nil {
		return c.Value
	}
	id := rand.Text()
	http.SetCookie(w, &http.Cookie{Name: "player", Value: id, Path: "/", MaxAge: 365 * 24 * 3600, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	return id
}

// playData is what play.html draws: the whole game area of game g for puzzle n.
// ok colours the message as good news.
func playData(n int, g game.Game, msg string, ok bool) map[string]any {
	person := game.PersonFor(n)
	var answers []map[string]any
	for _, q := range slices.Backward(g.Asked) {
		answers = append(answers, map[string]any{"ID": int(q), "Question": q.Text(), "Answer": person.Answer(q)})
	}
	var tried []*game.Person
	for _, id := range g.Tried {
		p, _ := game.PersonByID(id)
		tried = append(tried, p)
	}
	return map[string]any{"Answers": answers, "Tried": tried, "Msg": msg, "OK": ok}
}

// NewHandler returns the game's HTTP routes. now supplies the current time so tests can fix the day.
func NewHandler(now func() time.Time) http.Handler { return newServer(now) }

func newServer(now func() time.Time) *server {
	mux := http.NewServeMux()
	s := &server{now: now, store: memStore{games: map[gameKey]*game.Game{}}, mux: mux}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		n := game.Current(now())
		g, _ := s.store.update(playerID(w, r), n, false, func(*game.Game) {})
		tmpl.ExecuteTemplate(w, "index.html", map[string]any{"Play": playData(n, g, "", false), "Data": pageData})
	})

	mux.HandleFunc("POST /ask/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 0 || id >= len(game.Questions) {
			http.NotFound(w, r)
			return
		}
		q := game.Questions[id]
		n := game.Current(now())
		g, _ := s.store.update(playerID(w, r), n, true, func(g *game.Game) {
			if !slices.Contains(g.Asked, q) {
				g.Asked = append(g.Asked, q)
			}
		})
		tmpl.ExecuteTemplate(w, "play.html", playData(n, g, "", false))
	})

	mux.HandleFunc("POST /guess", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.FormValue("person"))
		p, ok := game.PersonByID(id)
		if err != nil || !ok {
			http.Error(w, "unknown person", http.StatusBadRequest)
			return
		}
		n := game.Current(now())
		correct := p.ID == game.PersonFor(n).ID
		msg := p.Name + " is not who I'm thinking of. Try again!"
		if correct {
			msg = "🎉 That's right! " + p.Name + " is today's person!"
		}
		g, _ := s.store.update(playerID(w, r), n, true, func(g *game.Game) {
			if !correct && !slices.Contains(g.Tried, p.ID) {
				g.Tried = append(g.Tried, p.ID)
			}
		})
		tmpl.ExecuteTemplate(w, "play.html", playData(n, g, msg, correct))
	})

	return s
}

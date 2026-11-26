// Package web serves the game in the browser: html/template pages plus htmx fragments.
package web

import (
	"crypto/rand"
	"embed"
	"html/template"
	"log"
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

// render writes template name to w. A template error is logged, because part of the page may already be sent.
func render(w http.ResponseWriter, name string, data any) {
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("render %s: %v", name, err)
	}
}

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
		people = append(people, map[string]any{"id": int(p.ID), "name": p.Name, "aliases": p.Aliases})
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
func playerID(w http.ResponseWriter, r *http.Request) game.PlayerID {
	if c, err := r.Cookie("player"); err == nil {
		return game.PlayerID(c.Value)
	}
	id := rand.Text()
	http.SetCookie(w, &http.Cookie{Name: "player", Value: id, Path: "/", MaxAge: 365 * 24 * 3600, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	return game.PlayerID(id)
}

// playData is what play.html draws: the whole game area of game g for puzzle n.
// msg is shown only while the game is being played.
func playData(n game.Puzzle, g game.Game, msg string) map[string]any {
	person := game.PersonFor(n)
	answer := func(q game.Question) map[string]any {
		return map[string]any{"ID": int(q), "Question": q.Text(), "Answer": person.Answer(q)}
	}
	if g.Finished() {
		var all []map[string]any
		for _, q := range game.Questions {
			all = append(all, answer(q))
		}
		// No score after giving up: giving up at once would score 0, the best score.
		summary := count(len(g.Asked), "question", "questions") + ", " + count(len(g.Tried), "miss", "misses")
		if g.Status == game.Solved {
			summary = "Score " + strconv.Itoa(g.Score()) + ": " + summary
		}
		return map[string]any{"Finished": true, "Won": g.Status == game.Solved, "Name": person.Name, "Answers": all, "Summary": summary}
	}
	var answers []map[string]any
	for _, q := range slices.Backward(g.Asked) {
		answers = append(answers, answer(q))
	}
	var tried []*game.Person
	for _, id := range g.Tried {
		p, _ := game.PersonByID(id)
		tried = append(tried, p)
	}
	return map[string]any{"Answers": answers, "Tried": tried, "Msg": msg}
}

// count returns, for example, "1 miss" or "2 misses".
func count(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return strconv.Itoa(n) + " " + many
}

// puzzle reads the puzzle number the page sent. When it is not a number (400) or the puzzle
// no longer accepts actions ("ended" contents), it writes the reply and returns ok == false.
func puzzle(w http.ResponseWriter, r *http.Request, now time.Time) (p game.Puzzle, ok bool) {
	n, err := strconv.Atoi(r.FormValue("puzzle"))
	if err != nil {
		http.Error(w, "bad puzzle", http.StatusBadRequest)
		return 0, false
	}
	p = game.Puzzle(n)
	if !game.Open(p, now) {
		render(w, "ended.html", nil)
		return 0, false
	}
	return p, true
}

// NewHandler returns the game's HTTP routes. now supplies the current time so tests can fix the day.
func NewHandler(now func() time.Time) http.Handler { return newServer(now) }

func newServer(now func() time.Time) *server {
	mux := http.NewServeMux()
	s := &server{now: now, store: memStore{games: map[gameKey]*game.Game{}}, mux: mux}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		n := game.Current(now())
		g, _ := s.store.update(playerID(w, r), n, false, func(*game.Game) {})
		render(w, "index.html", map[string]any{"N": n, "Play": playData(n, g, ""), "Data": pageData})
	})

	mux.HandleFunc("POST /ask/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 0 || id >= len(game.Questions) {
			http.NotFound(w, r)
			return
		}
		q := game.Questions[id]
		n, ok := puzzle(w, r, now())
		if !ok {
			return
		}
		g, _ := s.store.update(playerID(w, r), n, true, func(g *game.Game) { g.Ask(q) })
		render(w, "play.html", playData(n, g, ""))
	})

	mux.HandleFunc("POST /guess", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.FormValue("person"))
		p, ok := game.PersonByID(game.PersonID(id))
		if err != nil || !ok {
			http.Error(w, "unknown person", http.StatusBadRequest)
			return
		}
		n, ok := puzzle(w, r, now())
		if !ok {
			return
		}
		g, _ := s.store.update(playerID(w, r), n, true, func(g *game.Game) { g.Guess(p.ID, game.PersonFor(n).ID) })
		// A right guess or a finished game shows the result card, which has no message.
		render(w, "play.html", playData(n, g, p.Name+" is not who I'm thinking of. Try again!"))
	})

	mux.HandleFunc("POST /giveup", func(w http.ResponseWriter, r *http.Request) {
		n, ok := puzzle(w, r, now())
		if !ok {
			return
		}
		g, _ := s.store.update(playerID(w, r), n, true, func(g *game.Game) { g.GiveUp() })
		render(w, "play.html", playData(n, g, ""))
	})

	return s
}

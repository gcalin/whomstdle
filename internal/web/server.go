// Package web serves the game in the browser: html/template pages plus htmx fragments.
package web

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"html/template"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
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

// state is one player's progress for one day.
type state struct {
	Asked []game.Question
	Tried []string
}

type server struct {
	now   func() time.Time
	mu    sync.Mutex // guards games: handlers run concurrently
	games map[string]*state
}

// play runs f on the caller's game for today. Handlers must not write the response inside f.
func (s *server) play(w http.ResponseWriter, r *http.Request, f func(*state)) {
	c, err := r.Cookie("player")
	if err != nil {
		b := make([]byte, 16)
		rand.Read(b)
		c = &http.Cookie{Name: "player", Value: hex.EncodeToString(b), Path: "/", MaxAge: 365 * 24 * 3600, HttpOnly: true, SameSite: http.SameSiteLaxMode}
		http.SetCookie(w, c)
	}
	key := c.Value + s.now().UTC().Format("2006-01-02")

	s.mu.Lock()
	defer s.mu.Unlock()
	st := s.games[key]
	if st == nil {
		st = &state{}
		s.games[key] = st
	}
	f(st)
}

// NewHandler returns the game's HTTP routes. now supplies the current time so tests can fix the day.
func NewHandler(now func() time.Time) http.Handler {
	s := &server{now: now, games: map[string]*state{}}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		d := struct {
			Answers []map[string]any
			Tried   []string
			Data    any
		}{Data: pageData}
		s.play(w, r, func(st *state) {
			for _, q := range slices.Backward(st.Asked) {
				d.Answers = append(d.Answers, map[string]any{"ID": int(q), "Question": q.Text(), "Answer": game.PersonFor(game.Current(now())).Answer(q)})
			}
			d.Tried = slices.Clone(st.Tried)
		})
		tmpl.ExecuteTemplate(w, "index.html", d)
	})

	mux.HandleFunc("POST /ask/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 0 || id >= len(game.Questions) {
			http.NotFound(w, r)
			return
		}
		q := game.Questions[id]
		fresh := false
		s.play(w, r, func(st *state) {
			if fresh = !slices.Contains(st.Asked, q); fresh {
				st.Asked = append(st.Asked, q)
			}
		})
		if !fresh {
			w.WriteHeader(http.StatusNoContent) // already asked: htmx leaves the page alone
			return
		}
		tmpl.ExecuteTemplate(w, "answer.html", map[string]any{
			"ID":       id,
			"Question": q.Text(),
			"Answer":   game.PersonFor(game.Current(now())).Answer(q),
		})
	})

	mux.HandleFunc("POST /guess", func(w http.ResponseWriter, r *http.Request) {
		guess := r.FormValue("guess")
		target := game.PersonFor(game.Current(now()))
		p, known := game.FindPerson(guess)
		correct := known && p == target
		wrong := ""
		if known && !correct {
			s.play(w, r, func(st *state) {
				if !slices.Contains(st.Tried, p.Name) {
					st.Tried = append(st.Tried, p.Name)
					wrong = p.Name
				}
			})
		}
		tmpl.ExecuteTemplate(w, "guess.html", map[string]any{
			"Guess": guess, "Known": known, "Correct": correct, "Name": target.Name, "Wrong": wrong,
		})
	})

	return mux
}

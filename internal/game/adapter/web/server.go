// Package web serves the game in the browser: html/template pages plus htmx fragments.
package web

import (
	"embed"
	"errors"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/gcalin/whomstdle/internal/game/app"
	"github.com/gcalin/whomstdle/internal/game/domain"
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

type server struct {
	svc *app.Service
	mux *http.ServeMux
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }

// parsePuzzle reads the puzzle number the page sent. When it is not a number, it replies 400 and returns ok == false.
func parsePuzzle(w http.ResponseWriter, r *http.Request) (domain.Puzzle, bool) {
	n, err := strconv.Atoi(r.FormValue("puzzle"))
	if err != nil {
		http.Error(w, "bad puzzle", http.StatusBadRequest)
		return 0, false
	}
	return domain.Puzzle(n), true
}

// replyError writes the reply for an error from the app.
func replyError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, app.ErrUnknownQuestion):
		http.NotFound(w, r)
	case errors.Is(err, app.ErrUnknownPerson):
		http.Error(w, "unknown person", http.StatusBadRequest)
	case errors.Is(err, app.ErrClosed):
		render(w, "ended.html", nil)
	default:
		log.Printf("%s %s: %v", r.Method, r.URL.Path, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

// replyPlay renders the game area for board b, or the reply for err.
func replyPlay(w http.ResponseWriter, r *http.Request, b app.Board, err error, msg string) {
	if err != nil {
		replyError(w, r, err)
		return
	}
	render(w, "play.html", newPlayView(b, msg))
}

// NewHandler returns the game's HTTP routes.
func NewHandler(svc *app.Service) http.Handler {
	mux := http.NewServeMux()
	s := &server{svc: svc, mux: mux}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		b, err := s.svc.Today(r.Context(), playerID(w, r))
		if err != nil {
			replyError(w, r, err)
			return
		}
		c, err := s.svc.Choices(r.Context())
		if err != nil {
			replyError(w, r, err)
			return
		}
		render(w, "index.html", pageView{N: b.Puzzle, Play: newPlayView(b, ""), Data: newSearchData(c)})
	})

	mux.HandleFunc("POST /ask/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		p, ok := parsePuzzle(w, r)
		if !ok {
			return
		}
		b, err := s.svc.Ask(r.Context(), playerID(w, r), p, domain.Question(id))
		replyPlay(w, r, b, err, "")
	})

	mux.HandleFunc("POST /guess", func(w http.ResponseWriter, r *http.Request) {
		n, err := strconv.Atoi(r.FormValue("person"))
		if err != nil {
			http.Error(w, "unknown person", http.StatusBadRequest)
			return
		}
		p, ok := parsePuzzle(w, r)
		if !ok {
			return
		}
		id := domain.PersonID(n)
		b, err := s.svc.Guess(r.Context(), playerID(w, r), p, id)
		// A right guess or a finished game shows the result card, which has no message.
		msg := ""
		if i := slices.IndexFunc(b.Tried, func(t app.PersonName) bool { return t.ID == id }); i >= 0 && b.Status == domain.Playing {
			msg = b.Tried[i].Name + " is not who I'm thinking of. Try again!"
		}
		replyPlay(w, r, b, err, msg)
	})

	mux.HandleFunc("POST /giveup", func(w http.ResponseWriter, r *http.Request) {
		p, ok := parsePuzzle(w, r)
		if !ok {
			return
		}
		b, err := s.svc.GiveUp(r.Context(), playerID(w, r), p)
		replyPlay(w, r, b, err, "")
	})

	return s
}

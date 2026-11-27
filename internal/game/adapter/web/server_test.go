package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gcalin/whomstdle/internal/game/app"
	"github.com/gcalin/whomstdle/internal/game/domain"
)

// fakeGame returns board (or choices) and err, and records the last call.
type fakeGame struct {
	board   app.Board
	choices app.Choices
	err     error

	called string // the method name, empty when nothing was called
	player domain.PlayerID
	puzzle domain.Puzzle
	arg    int // the question or person
}

func (f *fakeGame) Today(_ context.Context, player domain.PlayerID) (app.Board, error) {
	f.called, f.player = "Today", player
	return f.board, f.err
}

func (f *fakeGame) Ask(_ context.Context, player domain.PlayerID, p domain.Puzzle, q domain.Question) (app.Board, error) {
	f.called, f.player, f.puzzle, f.arg = "Ask", player, p, int(q)
	return f.board, f.err
}

func (f *fakeGame) Guess(_ context.Context, player domain.PlayerID, p domain.Puzzle, id domain.PersonID) (app.Board, error) {
	f.called, f.player, f.puzzle, f.arg = "Guess", player, p, int(id)
	return f.board, f.err
}

func (f *fakeGame) GiveUp(_ context.Context, player domain.PlayerID, p domain.Puzzle) (app.Board, error) {
	f.called, f.player, f.puzzle = "GiveUp", player, p
	return f.board, f.err
}

func (f *fakeGame) Choices(context.Context) (app.Choices, error) {
	return f.choices, f.err
}

// send sends one request to a handler for f, with the player cookie "p1".
func send(f *fakeGame, method, path string, form url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "player", Value: "p1"})
	rec := httptest.NewRecorder()
	NewHandler(f).ServeHTTP(rec, req)
	return rec
}

func TestIndexSetsCookie(t *testing.T) {
	f := &fakeGame{choices: app.Choices{Questions: []domain.Question{0}, People: []app.PersonName{{ID: 7, Name: "Test Person"}}}}
	rec := httptest.NewRecorder()
	NewHandler(f).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if len(rec.Result().Cookies()) != 1 {
		t.Error("GET / should set the player cookie")
	}
	if body := rec.Body.String(); !strings.Contains(body, "Test Person") || !strings.Contains(body, "Are they alive?") {
		t.Errorf("page is missing the choices: %q", body)
	}
}

func TestCalls(t *testing.T) {
	for _, c := range []struct {
		path string
		form url.Values
		want string
		arg  int
	}{
		{"/ask/1", url.Values{"puzzle": {"1"}}, "Ask", 1},
		{"/guess", url.Values{"puzzle": {"1"}, "person": {"3"}}, "Guess", 3},
		{"/giveup", url.Values{"puzzle": {"1"}}, "GiveUp", 0},
	} {
		f := &fakeGame{}
		send(f, "POST", c.path, c.form)
		if f.called != c.want || f.arg != c.arg || f.puzzle != 1 || f.player != "p1" {
			t.Errorf("POST %s: called %s(%v, %d, %d), want %s(p1, 1, %d)", c.path, f.called, f.player, f.puzzle, f.arg, c.want, c.arg)
		}
	}
}

func TestBadValues(t *testing.T) {
	for _, c := range []struct {
		path string
		form url.Values
		code int
	}{
		{"/ask/abc", url.Values{"puzzle": {"1"}}, http.StatusNotFound},
		{"/ask/1", url.Values{"puzzle": {"abc"}}, http.StatusBadRequest},
		{"/guess", url.Values{"puzzle": {"1"}, "person": {"abc"}}, http.StatusBadRequest},
	} {
		f := &fakeGame{}
		if rec := send(f, "POST", c.path, c.form); rec.Code != c.code || f.called != "" {
			t.Errorf("POST %s %v: got %d and called %q, want %d and no call", c.path, c.form, rec.Code, f.called, c.code)
		}
	}
}

func TestErrors(t *testing.T) {
	for _, c := range []struct {
		err  error
		code int
	}{
		{app.ErrUnknownQuestion, http.StatusNotFound},
		{app.ErrUnknownPerson, http.StatusBadRequest},
		{app.ErrClosed, http.StatusOK},
		{errors.New("disk full"), http.StatusInternalServerError},
	} {
		rec := send(&fakeGame{err: c.err}, "POST", "/ask/1", url.Values{"puzzle": {"1"}})
		if rec.Code != c.code {
			t.Errorf("%v: got %d, want %d", c.err, rec.Code, c.code)
		}
		if c.err == app.ErrClosed && !strings.Contains(rec.Body.String(), "This puzzle has ended") {
			t.Errorf("closed: %q", rec.Body.String())
		}
	}
}

func TestRender(t *testing.T) {
	f := &fakeGame{board: app.Board{Status: domain.Solved, Name: "Test Person"}}
	if body := send(f, "POST", "/giveup", url.Values{"puzzle": {"1"}}).Body.String(); !strings.Contains(body, "It's Test Person.") {
		t.Errorf("finished board: %q", body)
	}
	f = &fakeGame{board: app.Board{Status: domain.Solved, SecondsLeft: 43200}}
	if body := send(f, "POST", "/giveup", url.Values{"puzzle": {"1"}}).Body.String(); !strings.Contains(body, `data-left="43200"`) {
		t.Errorf("countdown: %q", body)
	}
	f = &fakeGame{board: app.Board{Tried: []app.PersonName{{ID: 3, Name: "Wrong Person"}}}}
	if body := send(f, "POST", "/guess", url.Values{"puzzle": {"1"}, "person": {"3"}}).Body.String(); !strings.Contains(body, "Wrong Person is not who I") {
		t.Errorf("wrong guess: %q", body)
	}
}

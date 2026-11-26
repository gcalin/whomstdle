package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gcalin/whomstdle/internal/game"
)

var day = time.Date(2026, 11, 25, 12, 0, 0, 0, time.UTC)

func dayClock() time.Time { return day }

var p1 = url.Values{"puzzle": {"1"}}

// client returns a function that sends requests to h as one player (it keeps the cookie).
func client(h http.Handler) func(method, path string, form url.Values) *httptest.ResponseRecorder {
	var cookies []*http.Cookie
	return func(method, path string, form url.Values) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if cs := rec.Result().Cookies(); len(cs) > 0 {
			cookies = cs
		}
		return rec
	}
}

func TestIndex(t *testing.T) {
	rec := client(NewHandler(dayClock))("GET", "/", nil)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Why are they famous?") || !strings.Contains(rec.Body.String(), "Michael Jackson") {
		t.Fatalf("page should embed the question aliases and names: %d", rec.Code)
	}
}

func TestAsk(t *testing.T) {
	do := client(NewHandler(dayClock))
	want := game.PersonFor(game.Current(day)).Answer(game.Gender)
	if body := do("POST", "/ask/1", p1).Body.String(); !strings.Contains(body, want) {
		t.Fatalf("%q missing %q", body, want)
	}
	if rec := do("POST", "/ask/1", p1); rec.Code != 200 || strings.Count(rec.Body.String(), `<li data-id="1"><span>`) != 1 {
		t.Fatalf("repeat question: got %d, want 200 with the answer once: %q", rec.Code, rec.Body.String())
	}
	if rec := do("POST", "/ask/99", p1); rec.Code != http.StatusNotFound {
		t.Fatalf("got %d, want 404", rec.Code)
	}
	if body := do("GET", "/", nil).Body.String(); !strings.Contains(body, want) {
		t.Fatalf("reload lost the answer: %q", body)
	}
}

func TestGuess(t *testing.T) {
	do := client(NewHandler(dayClock))
	target := strconv.Itoa(game.PersonFor(game.Current(day)).ID) // puzzle 1: Andrei Boboc, ID 1
	tried := `<li data-id="3">Michael Jackson</li>`
	if body := do("POST", "/guess", url.Values{"person": {"3"}, "puzzle": {"1"}}).Body.String(); strings.Contains(body, `class="ok"`) || !strings.Contains(body, tried) {
		t.Fatalf("wrong guess: %q", body)
	}
	if body := do("POST", "/ask/2", p1).Body.String(); !strings.Contains(body, tried) {
		t.Fatalf("ask reply lost the tried names: %q", body)
	}
	if body := do("POST", "/guess", url.Values{"person": {target}, "puzzle": {"1"}}).Body.String(); !strings.Contains(body, `class="ok"`) || !strings.Contains(body, `<li data-id="2"><span>`) || !strings.Contains(body, tried) {
		t.Fatalf("correct guess reply should show the message, the answers and the tried names: %q", body)
	}
	if body := do("GET", "/", nil).Body.String(); !strings.Contains(body, tried) {
		t.Fatalf("tried guess lost on reload: %q", body)
	}
	for _, bad := range []string{"abc", "999"} {
		if rec := do("POST", "/guess", url.Values{"person": {bad}, "puzzle": {"1"}}); rec.Code != http.StatusBadRequest {
			t.Fatalf("person=%s: got %d, want 400", bad, rec.Code)
		}
	}
}

func TestPageDataPeople(t *testing.T) {
	people := pageData["people"].([]map[string]any)
	for i, p := range people {
		if p["id"].(int) < 1 {
			t.Errorf("%s has no id", p["name"])
		}
		if i > 0 && people[i-1]["name"].(string) > p["name"].(string) {
			t.Errorf("people not sorted by name: %s before %s", people[i-1]["name"], p["name"])
		}
	}
}

func TestNoGameUntilAction(t *testing.T) {
	s := newServer(dayClock)
	do := client(s)
	do("GET", "/", nil)
	do("POST", "/ask/99", p1)
	do("POST", "/guess", url.Values{"person": {"abc"}, "puzzle": {"1"}})
	if len(s.store.games) != 0 {
		t.Fatalf("page load or bad request created a game: %d games", len(s.store.games))
	}
	do("POST", "/ask/1", p1)
	if len(s.store.games) != 1 {
		t.Fatalf("asking should create one game, got %d", len(s.store.games))
	}
}

func TestMidnight(t *testing.T) {
	at := time.Date(2026, 11, 25, 23, 58, 0, 0, time.UTC) // puzzle 1
	s := newServer(func() time.Time { return at })
	do := client(s)
	alive := `<li data-id="0"><span>Are they alive?</span><strong>Yes</strong></li>` // puzzle 2's answer is "No"
	do("POST", "/ask/0", p1)

	at = at.Add(5 * time.Minute) // 00:03, puzzle 2
	if body := do("POST", "/ask/1", p1).Body.String(); !strings.Contains(body, alive) {
		t.Fatalf("the game lost its person or answers after midnight: %q", body)
	}
	if body := do("POST", "/guess", url.Values{"person": {"1"}, "puzzle": {"1"}}).Body.String(); !strings.Contains(body, `class="ok"`) {
		t.Fatalf("puzzle 1's person should still win after midnight: %q", body)
	}

	at = at.Add(58 * time.Minute) // 01:01, the grace hour is over
	rec := do("POST", "/ask/2", p1)
	if body := rec.Body.String(); rec.Code != 200 || !strings.Contains(body, "This puzzle has ended") || strings.Contains(body, "Andrei Boboc") {
		t.Fatalf("after the grace hour: got %d %q", rec.Code, body)
	}
	games := len(s.store.games)
	client(s)("POST", "/ask/0", p1) // a new player
	if len(s.store.games) != games {
		t.Fatal("a refused first action created a game")
	}
	if rec := do("POST", "/ask/0", url.Values{"puzzle": {"abc"}}); rec.Code != http.StatusBadRequest {
		t.Fatalf("puzzle=abc: got %d, want 400", rec.Code)
	}
}

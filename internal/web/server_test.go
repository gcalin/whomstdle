package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gcalin/whomstdle/internal/game"
)

var day = time.Date(2026, 11, 25, 12, 0, 0, 0, time.UTC)

func dayClock() time.Time { return day }

var p1 = url.Values{"puzzle": {"1"}}

// guessForm is the form of a guess for puzzle 1.
func guessForm(person string) url.Values { return url.Values{"person": {person}, "puzzle": {"1"}} }

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
	target := strconv.Itoa(int(game.PersonFor(game.Current(day)).ID)) // puzzle 1: Andrei Boboc, ID 1
	tried := `<li data-id="3">Michael Jackson</li>`
	if body := do("POST", "/guess", url.Values{"person": {"3"}, "puzzle": {"1"}}).Body.String(); strings.Contains(body, `class="ok"`) || !strings.Contains(body, tried) {
		t.Fatalf("wrong guess: %q", body)
	}
	if body := do("POST", "/ask/2", p1).Body.String(); !strings.Contains(body, tried) {
		t.Fatalf("ask reply lost the tried names: %q", body)
	}
	if body := do("GET", "/", nil).Body.String(); !strings.Contains(body, tried) {
		t.Fatalf("tried guess lost on reload: %q", body)
	}
	if body := do("POST", "/guess", url.Values{"person": {target}, "puzzle": {"1"}}).Body.String(); !strings.Contains(body, `class="ok"`) {
		t.Fatalf("correct guess should win: %q", body)
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
	do("POST", "/giveup", url.Values{"puzzle": {"abc"}})
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

	if body := client(s)("POST", "/giveup", p1).Body.String(); !strings.Contains(body, "Andrei Boboc") {
		t.Fatalf("giving up yesterday's puzzle should show yesterday's person: %q", body)
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

func TestSolve(t *testing.T) {
	do := client(NewHandler(dayClock))
	for _, q := range []string{"0", "1", "2", "3"} {
		do("POST", "/ask/"+q, p1)
	}
	do("POST", "/guess", guessForm("3"))
	card := do("POST", "/guess", guessForm("1")).Body.String()
	for _, want := range []string{"You got it!", "Andrei Boboc", "Score 5: 4 questions, 1 miss"} {
		if !strings.Contains(card, want) {
			t.Errorf("card is missing %q", want)
		}
	}
	if n := strings.Count(card, "<li data-id="); n != len(game.Questions) {
		t.Errorf("card shows %d answers, want all %d", n, len(game.Questions))
	}
	for _, gone := range []string{`id="q"`, `id="g"`, "/giveup"} {
		if strings.Contains(card, gone) {
			t.Errorf("card should not contain %q", gone)
		}
	}
	for _, a := range []struct {
		path string
		form url.Values
	}{{"/ask/5", p1}, {"/guess", guessForm("2")}, {"/guess", guessForm("1")}, {"/giveup", p1}} {
		if body := do("POST", a.path, a.form).Body.String(); body != card {
			t.Errorf("POST %s after the win changed the reply: %q", a.path, body)
		}
	}
	if body := do("GET", "/", nil).Body.String(); !strings.Contains(body, card) {
		t.Errorf("reload does not show the card: %q", body)
	}
	if body := client(NewHandler(dayClock))("POST", "/guess", guessForm("1")).Body.String(); !strings.Contains(body, "Score 0: 0 questions, 0 misses") {
		t.Errorf("a right first guess should score 0: %q", body)
	}
}

func TestGiveUp(t *testing.T) {
	h := NewHandler(dayClock)
	if body := client(h)("POST", "/giveup", p1).Body.String(); !strings.Contains(body, "You gave up.") || !strings.Contains(body, "0 questions, 0 misses") {
		t.Errorf("give up as the first action: %q", body)
	}
	do := client(h)
	do("POST", "/ask/1", p1)
	card := do("POST", "/giveup", p1).Body.String()
	for _, want := range []string{"You gave up.", "Andrei Boboc", "1 question, 0 misses"} {
		if !strings.Contains(card, want) {
			t.Errorf("card is missing %q", want)
		}
	}
	if strings.Contains(card, "Score") || strings.Contains(card, "You got it!") {
		t.Errorf("a give up has no score and is not a win: %q", card)
	}
	if body := do("POST", "/guess", guessForm("1")).Body.String(); body != card {
		t.Errorf("a right guess after giving up changed the reply: %q", body)
	}
	if body := do("GET", "/", nil).Body.String(); !strings.Contains(body, card) {
		t.Errorf("reload does not show the card: %q", body)
	}
}

func TestNoNameWhilePlaying(t *testing.T) {
	do := client(NewHandler(dayClock))
	name := game.PersonFor(1).Name
	replies := []string{do("POST", "/guess", guessForm("2")).Body.String(), do("POST", "/guess", guessForm("3")).Body.String()}
	for id := range game.Questions { // every question, the hint too
		replies = append(replies, do("POST", "/ask/"+strconv.Itoa(id), p1).Body.String())
	}
	for _, body := range replies {
		if strings.Contains(body, name) {
			t.Fatalf("reply for a game still being played gives away %s: %q", name, body)
		}
	}
}

// Right guesses and give ups at the same moment: the first one ends the game and every reply shows that same ending.
func TestFinishAtTheSameTime(t *testing.T) {
	h := NewHandler(dayClock)
	cookie := client(h)("POST", "/ask/0", p1).Result().Cookies()[0]
	bodies := make([]string, 20)
	var wg sync.WaitGroup
	for i := range bodies {
		wg.Go(func() {
			path, form := "/guess", guessForm("1")
			if i%2 == 1 {
				path, form = "/giveup", p1
			}
			req := httptest.NewRequest("POST", path, strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.AddCookie(cookie)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			bodies[i] = rec.Body.String()
		})
	}
	wg.Wait()
	if !strings.Contains(bodies[0], "Andrei Boboc") {
		t.Fatalf("expected a result card: %q", bodies[0])
	}
	for _, body := range bodies {
		if body != bodies[0] {
			t.Fatalf("replies disagree on how the game ended:\n%q\n%q", bodies[0], body)
		}
	}
}

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
	if body := do("POST", "/ask/1", nil).Body.String(); !strings.Contains(body, want) {
		t.Fatalf("%q missing %q", body, want)
	}
	if rec := do("POST", "/ask/1", nil); rec.Code != http.StatusNoContent {
		t.Fatalf("repeat question: got %d, want 204", rec.Code)
	}
	if rec := do("POST", "/ask/99", nil); rec.Code != http.StatusNotFound {
		t.Fatalf("got %d, want 404", rec.Code)
	}
	if body := do("GET", "/", nil).Body.String(); !strings.Contains(body, want) {
		t.Fatalf("reload lost the answer: %q", body)
	}
}

func TestGuess(t *testing.T) {
	do := client(NewHandler(dayClock))
	target := strconv.Itoa(game.PersonFor(game.Current(day)).ID) // puzzle 1: Andrei Boboc, ID 1
	if body := do("POST", "/guess", url.Values{"person": {target}}).Body.String(); !strings.Contains(body, "That's right") {
		t.Fatalf("correct guess: %q", body)
	}
	if body := do("POST", "/guess", url.Values{"person": {"3"}}).Body.String(); strings.Contains(body, "That's right") {
		t.Fatalf("wrong guess accepted: %q", body)
	}
	if body := do("GET", "/", nil).Body.String(); !strings.Contains(body, `<li data-id="3">Michael Jackson</li>`) {
		t.Fatalf("tried guess lost on reload: %q", body)
	}
	for _, bad := range []string{"abc", "999"} {
		if rec := do("POST", "/guess", url.Values{"person": {bad}}); rec.Code != http.StatusBadRequest {
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
	do("POST", "/ask/99", nil)
	do("POST", "/guess", url.Values{"person": {"abc"}})
	if len(s.store.games) != 0 {
		t.Fatalf("page load or bad request created a game: %d games", len(s.store.games))
	}
	do("POST", "/ask/1", nil)
	if len(s.store.games) != 1 {
		t.Fatalf("asking should create one game, got %d", len(s.store.games))
	}
}

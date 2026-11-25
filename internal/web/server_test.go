package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gcalin/whomstdle/internal/game"
)

var day = time.Date(2026, 11, 25, 12, 0, 0, 0, time.UTC)

// client returns a function that sends requests as one player (it keeps the cookie).
func client() func(method, path string, form url.Values) *httptest.ResponseRecorder {
	h := NewHandler(func() time.Time { return day })
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
	if rec := client()("GET", "/", nil); rec.Code != 200 || !strings.Contains(rec.Body.String(), "hx-get") {
		t.Fatalf("unexpected response: %d", rec.Code)
	}
}

func TestAsk(t *testing.T) {
	do := client()
	want := game.TodaysPerson(day).Answer(game.Gender)
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

func TestQuestions(t *testing.T) {
	do := client()
	if body := do("GET", "/questions?q=where+are+they", nil).Body.String(); !strings.Contains(body, "Which continent") {
		t.Fatalf("got %q", body)
	}
	do("POST", "/ask/3", nil)
	if body := do("GET", "/questions?q=where+are+they", nil).Body.String(); strings.Contains(body, "Which continent") {
		t.Fatalf("asked question still listed: %q", body)
	}
}

func TestGuess(t *testing.T) {
	do := client()
	target := game.TodaysPerson(day).Name
	if body := do("POST", "/guess", url.Values{"guess": {target}}).Body.String(); !strings.Contains(body, "That's right") {
		t.Fatalf("correct guess: %q", body)
	}
	if body := do("POST", "/guess", url.Values{"guess": {"Nobody"}}).Body.String(); strings.Contains(body, "That's right") {
		t.Fatalf("wrong guess accepted: %q", body)
	}
}

func TestSearch(t *testing.T) {
	do := client()
	if body := do("GET", "/search?q=mich", nil).Body.String(); !strings.Contains(body, "Michael Jackson") {
		t.Fatalf("got %q", body)
	}
	// Michael Jackson is not today's person, so this guess is wrong and gets remembered.
	do("POST", "/guess", url.Values{"guess": {"Michael Jackson"}})
	if body := do("GET", "/search?q=mich", nil).Body.String(); strings.Contains(body, "Michael Jackson") {
		t.Fatalf("tried guess still listed: %q", body)
	}
}

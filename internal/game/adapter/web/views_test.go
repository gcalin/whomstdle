package web

import (
	"strings"
	"testing"

	"github.com/gcalin/whomstdle/internal/game/app"
	"github.com/gcalin/whomstdle/internal/game/domain"
)

func TestNewPlayView(t *testing.T) {
	playing := app.Board{Answers: []app.Answer{{Question: 0, Text: "Yes"}, {Question: 1, Text: "Man"}}}
	if v := newPlayView(playing, ""); v.Answers[0].ID != 1 || v.Answers[1].ID != 0 {
		t.Errorf("playing: answers %v, want question 1 then 0", v.Answers)
	}
	solved := app.Board{Status: domain.Solved, Asked: 4, Tried: []app.PersonName{{ID: 3, Name: "Michael Jackson"}}, Score: 5}
	if got, want := newPlayView(solved, "").Summary, "Score 5: 4 questions, 1 miss"; got != want {
		t.Errorf("solved: summary %q, want %q", got, want)
	}
	if got := newPlayView(app.Board{Status: domain.GaveUp, Score: 3}, "").Summary; strings.Contains(got, "Score") {
		t.Errorf("gave up: summary %q should have no score", got)
	}
}

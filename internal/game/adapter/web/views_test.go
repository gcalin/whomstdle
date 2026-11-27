package web

import (
	"slices"
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
	if v := newPlayView(app.Board{Asked: 3}, ""); v.HintLock != "2 more questions" || v.Hint != "" {
		t.Errorf("3 asked: hint lock %q, want \"2 more questions\"", v.HintLock)
	}
	if v := newPlayView(app.Board{Asked: 5}, ""); v.HintLock != "" {
		t.Errorf("5 asked: hint lock %q, want none", v.HintLock)
	}
	hintAsked := app.Board{Asked: 5, HintUsed: true, Answers: []app.Answer{{Question: 0, Text: "Yes"}, {Question: domain.Hint, Text: "The king of pop"}}}
	if v := newPlayView(hintAsked, ""); v.Hint != "The king of pop" || len(v.Answers) != 1 {
		t.Errorf("hint used: hint %q and %d answers, want the hint text and 1 answer", v.Hint, len(v.Answers))
	}
	solved := app.Board{Status: domain.Solved, Asked: 4, Tried: []app.PersonName{{ID: 3, Name: "Michael Jackson"}}, Score: 5}
	if got, want := newPlayView(solved, "").Summary, "Score 5: 4 questions, 1 miss"; got != want {
		t.Errorf("solved: summary %q, want %q", got, want)
	}
	hinted := app.Board{Status: domain.Solved, Asked: 5, HintUsed: true, Tried: []app.PersonName{{ID: 3, Name: "Michael Jackson"}}, Score: 9}
	if got, want := newPlayView(hinted, "").Summary, "Score 9: 5 questions, hint +3, 1 miss"; got != want {
		t.Errorf("hinted: summary %q, want %q", got, want)
	}
	if got := newPlayView(app.Board{Status: domain.GaveUp, Score: 3}, "").Summary; strings.Contains(got, "Score") {
		t.Errorf("gave up: summary %q should have no score", got)
	}
}

func TestNewSearchDataLeavesOutHint(t *testing.T) {
	d := newSearchData(app.Choices{Questions: domain.Questions})
	if len(d.Questions) != len(domain.Questions)-1 || slices.ContainsFunc(d.Questions, func(q questionData) bool { return q.ID == domain.Hint }) {
		t.Errorf("questions %v, want all but the hint", d.Questions)
	}
}

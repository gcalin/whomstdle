package web

import (
	"slices"
	"strconv"

	"github.com/gcalin/whomstdle/internal/game/app"
	"github.com/gcalin/whomstdle/internal/game/domain"
)

// pageView is what index.html draws.
type pageView struct {
	N    domain.Puzzle
	Play playView
	Data searchData
}

// playView is what play.html draws: the whole game area. Msg is shown only while playing.
type playView struct {
	Finished bool
	Won      bool
	Name     string
	Answers  []answerView
	Summary  string
	Tried    []app.PersonName
	Msg      string
	// Hint is the hint's text once used. Until then HintLock says how many questions are still needed,
	// e.g. "2 more questions", or is empty when the hint can be used. Both are set only while playing.
	Hint     string
	HintLock string
	// SecondsLeft is the time until the next puzzle; used only on the result card.
	SecondsLeft int
}

type answerView struct {
	ID       domain.Question
	Question string // the question's text
	Answer   string
}

// searchData is sent once inside the page as JSON; the page script searches it as the player types.
type searchData struct {
	Questions []questionData `json:"questions"`
	People    []personData   `json:"people"`
}

type questionData struct {
	ID      domain.Question `json:"id"`
	Text    string          `json:"text"`
	Aliases []string        `json:"aliases"`
}

type personData struct {
	ID      domain.PersonID `json:"id"`
	Name    string          `json:"name"`
	Aliases []string        `json:"aliases"`
}

func newPlayView(b app.Board, msg string) playView {
	v := playView{Finished: b.Status != domain.Playing, Won: b.Status == domain.Solved, Name: b.Name, Tried: b.Tried, Msg: msg, SecondsLeft: b.SecondsLeft}
	for _, a := range b.Answers {
		if a.Question == domain.Hint && !v.Finished {
			v.Hint = a.Text // shown on the hint line, not in the list
			continue
		}
		v.Answers = append(v.Answers, answerView{ID: a.Question, Question: a.Question.Text(), Answer: a.Text})
	}
	if !v.Finished {
		slices.Reverse(v.Answers) // newest first
		if left := domain.HintAfter - b.Asked; left > 0 {
			v.HintLock = count(left, "more question", "more questions")
		}
		return v
	}
	// No score after giving up: giving up at once would score 0, the best score.
	v.Summary = count(b.Asked, "question", "questions")
	if b.HintUsed {
		v.Summary += ", hint +" + strconv.Itoa(domain.HintCost)
	}
	v.Summary += ", " + count(len(b.Tried), "miss", "misses")
	if v.Won {
		v.Summary = "Score " + strconv.Itoa(b.Score) + ": " + v.Summary
	}
	return v
}

func newSearchData(c app.Choices) searchData {
	var d searchData
	for _, q := range c.Questions {
		if q == domain.Hint {
			continue // the hint has its own button, not a place in the question search
		}
		d.Questions = append(d.Questions, questionData{ID: q, Text: q.Text(), Aliases: q.Aliases()})
	}
	for _, p := range c.People {
		d.People = append(d.People, personData{ID: p.ID, Name: p.Name, Aliases: p.Aliases})
	}
	return d
}

// count returns, for example, "1 miss" or "2 misses".
func count(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return strconv.Itoa(n) + " " + many
}

// HintID and HintCost let play.html draw the hint button.
func (playView) HintID() domain.Question { return domain.Hint }
func (playView) HintCost() int           { return domain.HintCost }

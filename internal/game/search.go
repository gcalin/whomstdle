package game

import (
	"slices"
	"strings"
	"unicode"
)

func normalize(s string) string {
	return strings.ToLower(strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			return r
		}
		return -1
	}, s))
}

// match reports whether every word of query appears in at least one of texts.
func match(query string, texts ...string) bool {
	words := strings.Fields(normalize(query))
	if len(words) == 0 {
		return false
	}
	for _, t := range texts {
		t = normalize(t)
		all := true
		for _, w := range words {
			all = all && strings.Contains(t, w)
		}
		if all {
			return true
		}
	}
	return false
}

// SearchQuestions returns the askable questions matching q, except those already asked.
func SearchQuestions(q string, asked []Question) []Question {
	var found []Question
	for _, question := range Questions {
		if question < Recap && !slices.Contains(asked, question) && match(q, append([]string{question.Text()}, aliases[question]...)...) {
			found = append(found, question)
		}
	}
	return found
}

// Search returns the names of people matching q, except those already tried.
func Search(q string, tried []string) []string {
	var names []string
	for _, p := range People {
		if !slices.Contains(tried, p.Name) && match(q, append([]string{p.Name}, p.Aliases...)...) {
			names = append(names, p.Name)
		}
	}
	return names
}

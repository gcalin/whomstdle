// Command game is the Phase 0 terminal version of the daily celebrity guessing
// game: no web, no database. It wires a hard-coded cast and question set to the
// core rules in internal/game and reads commands from stdin.
//
// Run it with:
//
//	go run ./cmd/game
//
// Commands once it is running: "ask <id>", "guess <name>", "help", "quit".
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gcalin/whomstdle/internal/game"
)

// questions are the fixed questions a player may ask. Expand this set as part of
// Phase 0/1; Phase 2 moves it into Postgres.
var questions = []game.Question{
	{ID: "alive", Text: "Still alive?"},
	{ID: "gender", Text: "Man or woman?"},
	{ID: "nationality", Text: "Nationality?"},
	{ID: "field", Text: "Field (film, music, sport, politics, science)?"},
	{ID: "before1960", Text: "Born before 1960?"},
	{ID: "knownfor", Text: "Known for?"},
}

// people is the hard-coded cast. The plan suggests starting with ~5 here and
// growing to 20+ in Phase 1. Every person must answer every question above.
var people = []game.Person{
	{
		Name:    "Albert Einstein",
		Aliases: []string{"Einstein"},
		Answers: map[game.QuestionID]string{
			"alive":       "no",
			"gender":      "man",
			"nationality": "German",
			"field":       "science",
			"before1960":  "yes",
			"knownfor":    "the theory of relativity",
		},
	},
	{
		Name:    "Marie Curie",
		Aliases: []string{"Curie", "Maria Skłodowska-Curie"},
		Answers: map[game.QuestionID]string{
			"alive":       "no",
			"gender":      "woman",
			"nationality": "Polish",
			"field":       "science",
			"before1960":  "yes",
			"knownfor":    "research on radioactivity",
		},
	},
	{
		Name:    "Nelson Mandela",
		Aliases: []string{"Mandela", "Madiba"},
		Answers: map[game.QuestionID]string{
			"alive":       "no",
			"gender":      "man",
			"nationality": "South African",
			"field":       "politics",
			"before1960":  "yes",
			"knownfor":    "ending apartheid",
		},
	},
	{
		Name:    "Serena Williams",
		Aliases: []string{"Serena"},
		Answers: map[game.QuestionID]string{
			"alive":       "yes",
			"gender":      "woman",
			"nationality": "American",
			"field":       "sport",
			"before1960":  "no",
			"knownfor":    "23 Grand Slam singles titles",
		},
	},
	{
		Name:    "David Bowie",
		Aliases: []string{"Bowie", "Ziggy Stardust"},
		Answers: map[game.QuestionID]string{
			"alive":       "no",
			"gender":      "man",
			"nationality": "British",
			"field":       "music",
			"before1960":  "no",
			"knownfor":    "reinventing rock and pop",
		},
	},
}

func main() {
	// Phase 0 picks a fixed secret. Phase 1 will choose today's person
	// deterministically from the date (daysSinceLaunch % len(people)).
	g := game.New(people[0], questions)

	in := bufio.NewScanner(os.Stdin)
	fmt.Println("Whomst'd'le — guess the secret person.")
	printHelp(g)

	for !g.Over() {
		fmt.Print("> ")
		if !in.Scan() {
			fmt.Println()
			return // EOF (Ctrl-D)
		}
		line := strings.TrimSpace(in.Text())

		switch {
		case line == "":
			continue
		case line == "quit" || line == "exit":
			fmt.Printf("The answer was %s. Bye.\n", g.Secret().Name)
			return
		case line == "help" || line == "?":
			printHelp(g)
		case strings.HasPrefix(line, "ask "):
			id := game.QuestionID(strings.TrimSpace(strings.TrimPrefix(line, "ask ")))
			if answer, ok := g.Ask(id); ok {
				fmt.Printf("  %s\n", answer)
			} else {
				fmt.Printf("  no such question %q — type 'help' to list them\n", id)
			}
		case strings.HasPrefix(line, "guess "):
			name := strings.TrimSpace(strings.TrimPrefix(line, "guess "))
			if g.Guess(name) {
				fmt.Printf("  Correct! It was %s. You asked %d question(s) and made %d wrong guess(es).\n",
					g.Secret().Name, g.Asked(), g.WrongGuesses())
			} else {
				fmt.Println("  Wrong — keep going.")
			}
		default:
			fmt.Println("  Commands: ask <id>, guess <name>, help, quit")
		}
	}
}

// printHelp lists the askable question IDs and the available commands.
func printHelp(g *game.Game) {
	fmt.Println("Ask about the secret person, then guess who it is.")
	fmt.Println("Questions you can ask (use the id):")
	for _, q := range g.Questions() {
		fmt.Printf("  ask %-12s %s\n", q.ID, q.Text)
	}
	fmt.Println("Guess with:  guess <name>    Quit with:  quit")
}

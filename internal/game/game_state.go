// Package game holds the core rules of the guessing game, plus the terminal
// read/print loop that drives them.
package game

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

type GameState struct {
	alreadyAsked map[Question]int
	numGuesses   int
	targetPerson *Person
	reader       *bufio.Reader
}

func NewGameState() *GameState {
	return &GameState{
		alreadyAsked: map[Question]int{},
		targetPerson: &People[rand.IntN(len(People))],
		reader:       bufio.NewReader(os.Stdin),
	}
}

func (g *GameState) ProcessQuestion() (Question, int, error) {
	s, err := g.reader.ReadString('\n')

	if err != nil {
		fmt.Println("Oops! Something went wrong!")
		os.Exit(1)
	}

	i, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return Alive, -1, errors.New("that's not a number, try again")
	}

	if 1 <= i && i <= len(Questions) {
		return Questions[i-1], i, nil
	}
	return Alive, -1, fmt.Errorf("try questions between %d and %d", 1, len(Questions))
}

func (g *GameState) ProcessAnswer(q Question, qID int) {
	switch q {
	case Recap:
		for prevQ, i := range g.alreadyAsked {
			fmt.Printf("%d. %s\n", i, g.targetPerson.Answer(prevQ))
		}
	case MakeAGuess:
		g.numGuesses++
		fmt.Println("Type your guess and hit Enter!")
		answer, err := g.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ooops! Something went very wrong!")
			os.Exit(1)
		}
		answer = strings.TrimSpace(answer)
		if person, ok := FindPerson(answer); ok {
			if person == g.targetPerson {
				fmt.Printf("That's right! %s is today's person! You guessed after %d questions and %d guesses.\n", g.targetPerson.Name, len(g.alreadyAsked), g.numGuesses)
				os.Exit(0)
			} else {
				fmt.Printf("%s is not who I'm thinking of! Try again!\n", answer)
			}
		} else {
			fmt.Printf("I don't even know who %s is!\n", answer)
		}
	default:
		fmt.Println(g.targetPerson.Answer(q))
		g.alreadyAsked[q] = qID
	}
}

func (g *GameState) WaitForEnter() {
	g.reader.ReadString('\n')
}

func (g *GameState) DisplayQuestions() {
	for i, q := range Questions {
		if _, ok := g.alreadyAsked[q]; ok {
			continue
		}
		fmt.Printf("%d. %s\n", i+1, q.Text())
	}
}

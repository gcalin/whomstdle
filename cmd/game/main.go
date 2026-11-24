// Command game is the terminal version of the guessing game.
package main

import (
	"fmt"

	"github.com/gcalin/whomstdle/internal/game"
)

func main() {
	state := game.NewGameState()

	fmt.Println("Whomst are we thinking of? Press Enter to start!")
	state.WaitForEnter()
	fmt.Println()

	fmt.Println("I'm thinking of a person. Ask me questions to figure out who they are! Type the number of the question you want to ask and then hit Enter.")

	for {
		fmt.Println("---")

		state.DisplayQuestions()

		fmt.Println("Make your choice...")

		q, id, err := state.ProcessQuestion()

		if err != nil {
			fmt.Println(err)
			continue
		}

		state.ProcessAnswer(q, id)
	}
}

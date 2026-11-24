package game

type Question int

const (
	Alive Question = iota
	KnownFor
	Weight
	Gender
	Recap
	MakeAGuess
)

var Questions = []Question{Alive, KnownFor, Weight, Gender, Recap, MakeAGuess}

var questionText = map[Question]string{
	Alive:      "Are they alive as of today?",
	KnownFor:   "What are they known for?",
	Weight:     "How much do they weigh?",
	Gender:     "What is their gender?",
	Recap:      "Recap the questions",
	MakeAGuess: "Make a guess!",
}

var answerPrefix = map[Question]string{
	Alive:      "",
	KnownFor:   "They are best known for ",
	Weight:     "They weigh right about ",
	Gender:     "They are ",
	Recap:      "",
	MakeAGuess: "",
}

func (q Question) Text() string {
	return questionText[q]
}

func (q Question) Prefix() string {
	return answerPrefix[q]
}

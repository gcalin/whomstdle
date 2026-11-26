package domain

type Question int

const (
	Alive Question = iota
	Gender
	Field
	Continent
	BirthDecade
	ActiveEra
	Languages
	PublicRole
	Awards
	Hint
)

var Questions = []Question{Alive, Gender, Field, Continent, BirthDecade, ActiveEra, Languages, PublicRole, Awards, Hint}

var questionText = map[Question]string{
	Alive:       "Are they alive?",
	Gender:      "What is their gender?",
	Field:       "What are they known for?",
	Continent:   "Which continent are they from?",
	BirthDecade: "When were they born?",
	ActiveEra:   "In which decades were they most famous?",
	Languages:   "Which languages did they speak?",
	PublicRole:  "Do they perform in front of an audience?",
	Awards:      "Have they won a major international award?",
	Hint:        "Give me a one-sentence hint",
}

func (q Question) Aliases() []string {
	return aliases[q]
}

func (q Question) Text() string {
	return questionText[q]
}

var aliases = map[Question][]string{
	Alive:       {"Are they still alive?", "Are they dead?", "Have they died?", "Are they living?", "Did they pass away?"},
	Gender:      {"Are they a man or a woman?", "Are they male or female?", "Which sex are they?"},
	Field:       {"Why are they famous?", "What do they do?", "What is their job?", "What is their profession?", "What is their field?", "What made them famous?"},
	Continent:   {"Where are they from?", "What part of the world are they from?"},
	BirthDecade: {"Which decade were they born in?", "What decade were they born?"},
	ActiveEra:   {"When were they famous?", "When were they active?", "When did they peak?", "Which era are they from?"},
	Languages:   {"What language do they speak?", "What is their mother tongue?", "What is their native language?"},
	PublicRole:  {"Are they a performer?", "Do they perform live?"},
	Awards:      {"Have they won any awards?", "Are they award-winning?"},
	Hint:        {"Can I get a hint?", "Give me a clue", "Tell me something about them", "I'm stuck", "Help"},
}

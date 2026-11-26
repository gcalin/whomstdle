package game

type Person struct {
	ID      int
	Name    string
	Aliases []string
	Answers map[Question]string
}

func (p *Person) Answer(q Question) string {
	return p.Answers[q]
}

var People = []Person{
	{
		ID:      1,
		Name:    "Andrei Boboc",
		Aliases: []string{"Boboc"},
		Answers: map[Question]string{
			Alive: "Yes", Gender: "Man", Field: "Science & tech", Continent: "Europe",
			BirthDecade: "1990s", ActiveEra: "2020s", Languages: "Romanian, English",
			PublicRole: "No", Awards: "No", Hint: "AI engineer with unpredictable bowel movements",
		},
	},
	{
		ID:      2,
		Name:    "Nihal",
		Aliases: []string{"Mihai Daniel Filologie", "Mihai Iosub"},
		Answers: map[Question]string{
			Alive: "Yes", Gender: "Man", Field: "Other", Continent: "Europe",
			BirthDecade: "1990s", ActiveEra: "2020s", Languages: "Romanian, English",
			PublicRole: "No", Awards: "No", Hint: "Communist thought leader and mobile gamer",
		},
	},
	{
		ID:   3,
		Name: "Michael Jackson",
		Answers: map[Question]string{
			Alive: "No", Gender: "Man", Field: "Music", Continent: "North America",
			BirthDecade: "1950s", ActiveEra: "1970s, 1980s, 1990s", Languages: "English",
			PublicRole: "Yes", Awards: "Yes", Hint: "The king of pop",
		},
	},
}

// PersonByID returns the person with the given ID.
func PersonByID(id int) (*Person, bool) {
	for i := range People {
		if People[i].ID == id {
			return &People[i], true
		}
	}
	return nil, false
}

// byName indexes People by canonical name and every alias, so a player's
// typed guess can be matched regardless of which name they use. It's built
// once from People, which stays the single source of truth.
var byName = buildIndex()

func buildIndex() map[string]*Person {
	index := map[string]*Person{}
	for i := range People {
		p := &People[i]
		index[p.Name] = p
		for _, alias := range p.Aliases {
			index[alias] = p
		}
	}
	return index
}

func FindPerson(name string) (*Person, bool) {
	p, ok := byName[name]
	return p, ok
}

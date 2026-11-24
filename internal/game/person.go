package game

type Person struct {
	Name    string
	Aliases []string
	Answers map[Question]string
}

func (p *Person) Answer(q Question) string {
	return q.Prefix() + p.Answers[q] + "."
}

var People = []Person{
	{
		Name:    "Andrei Boboc",
		Aliases: []string{"Boboc"},
		Answers: map[Question]string{
			Alive:    "Yes",
			KnownFor: "AI Engineering and unpredictable bowel movements",
			Weight:   "200kg",
			Gender:   "Male",
		},
	},
	{
		Name:    "Nihal",
		Aliases: []string{"Mihai Daniel Filologie", "Mihai Iosub"},
		Answers: map[Question]string{
			Alive:    "Yes",
			KnownFor: "being a communist thought leader & mobile gamer",
			Weight:   "110kg",
			Gender:   "Male",
		},
	},
	{
		Name:    "Michael Jackson",
		Aliases: nil,
		Answers: map[Question]string{
			Alive:    "Yes",
			KnownFor: "being the king of pop music",
			Weight:   "45kg",
			Gender:   "Male",
		},
	},
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

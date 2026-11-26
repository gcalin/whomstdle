package memory

import "github.com/gcalin/whomstdle/internal/game/domain"

var people = []domain.Person{
	{
		ID:      1,
		Name:    "Andrei Boboc",
		Aliases: []string{"Boboc"},
		Answers: map[domain.Question]string{
			domain.Alive: "Yes", domain.Gender: "Man", domain.Field: "Science & tech", domain.Continent: "Europe",
			domain.BirthDecade: "1990s", domain.ActiveEra: "2020s", domain.Languages: "Romanian, English",
			domain.PublicRole: "No", domain.Awards: "No", domain.Hint: "AI engineer with unpredictable bowel movements",
		},
	},
	{
		ID:      2,
		Name:    "Nihal",
		Aliases: []string{"Mihai Daniel Filologie", "Mihai Iosub"},
		Answers: map[domain.Question]string{
			domain.Alive: "Yes", domain.Gender: "Man", domain.Field: "Other", domain.Continent: "Europe",
			domain.BirthDecade: "1990s", domain.ActiveEra: "2020s", domain.Languages: "Romanian, English",
			domain.PublicRole: "No", domain.Awards: "No", domain.Hint: "Communist thought leader and mobile gamer",
		},
	},
	{
		ID:   3,
		Name: "Michael Jackson",
		Answers: map[domain.Question]string{
			domain.Alive: "No", domain.Gender: "Man", domain.Field: "Music", domain.Continent: "North America",
			domain.BirthDecade: "1950s", domain.ActiveEra: "1970s, 1980s, 1990s", domain.Languages: "English",
			domain.PublicRole: "Yes", domain.Awards: "Yes", domain.Hint: "The king of pop",
		},
	},
}

package domain

// PersonID is a person's fixed ID. Never change or reuse one.
type PersonID int

type Person struct {
	ID      PersonID
	Name    string
	Aliases []string
	Answers map[Question]string
}

func (p *Person) Answer(q Question) string {
	return p.Answers[q]
}

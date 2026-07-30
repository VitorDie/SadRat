package domain

type JobEntity struct {
	ID      string
	Command string
	Args    []string
}
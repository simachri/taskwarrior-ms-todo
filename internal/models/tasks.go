package models

type Task struct {
	ID    *string
	Title *string
	// Format is yyyy-MM-DDThh:mm:ss, example: 2022-08-02T00:00:00.0000000
	CompletedAt *string
}

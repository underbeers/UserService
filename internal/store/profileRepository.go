package store

import "github.com/google/uuid"

type ProfileRepository struct {
	store *Store
}

type Profiler interface {
}

type Profile struct {
	ID         uuid.UUID `db:"id"`
	FirstName  string    `db:"first_name"`
	SecondName string    `db:"second_name"`
	SurName    string    `db:"sur_name"`
	Status     string    `db:"status"`
}

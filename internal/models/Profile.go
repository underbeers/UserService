package models

import "github.com/google/uuid"

type Profile struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	SurName   string    `db:"sur_name"`
	Status    int       `db:"status"`
}
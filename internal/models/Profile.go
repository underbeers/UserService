package models

import (
	"github.com/google/uuid"
	"time"
)

type Profile struct {
	ID               uuid.UUID `db:"id"`
	FirstName        string    `db:"first_name"`
	SurName          string    `db:"sur_name"`
	Status           int       `db:"status"`
	ImageLink        *string   `db:"image_link"`
	DateRegistration time.Time `db:"date_registration"`
	Description      string    `db:"description"`
}

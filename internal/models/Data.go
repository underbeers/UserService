package models

import "github.com/google/uuid"

type Data struct {
	ID              int       `db:"id"`
	ProfileID       uuid.UUID `db:"id_profile"`
	PasswordEncoded string    `db:"password_encoded"`
	PasswordSalt    string    `db:"password_salt"`
}

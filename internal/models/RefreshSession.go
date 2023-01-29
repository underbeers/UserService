package models

import (
	"github.com/google/uuid"
	"time"
)

type RefreshSession struct {
	ID             int       `db:"id"`
	ProfileID      uuid.UUID `db:"id_profile"`
	RefreshTokenID uuid.UUID `db:"id_refresh_token"`
	IssuedAt       time.Time `db:"issued_at"`
	ExpiresIn      time.Time `db:"expires_in"`
}

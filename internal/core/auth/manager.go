package auth

import (
	"github.com/google/uuid"
	"time"
)

type Maker interface {
	CreateToken(profileID uuid.UUID, duration time.Duration,
		durationRtoken time.Duration, rToken uuid.UUID) (string, uuid.UUID, error)
	VerifyToken(token string) (*Payload, error)
	ParseExpiredToken(token string) (*Payload, error)
}

package auth

import (
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	TokenID   uuid.UUID `json:"tokenID"`
	ProfileID uuid.UUID `json:"profileID"`
	RtokenID  uuid.UUID `json:"rtokenID"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
	ExpiredIn time.Time `json:"expiredIn"`
}

func NewPayload(profileID uuid.UUID, duration time.Duration, durationRToken time.Duration,
	rTokenID uuid.UUID) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, genErr.NewError(err, core.ErrUnavailableResource, "msg", "error while creating payload")
	}
	payload := Payload{
		TokenID:   tokenID,
		ProfileID: profileID,
		RtokenID:  rTokenID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		ExpiredIn: time.Now().Add(durationRToken),
	}

	return &payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return genErr.NewError(nil, core.ErrTokenExpired)
	}

	return nil
}

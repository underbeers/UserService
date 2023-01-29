package auth

import (
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"git.friends.com/PetLand/UserService/v2/internal/store"
	"github.com/google/uuid"
)

func StoreSession(s *store.Store, token string) error {
	payload, err := ParseToken(token)
	if err != nil {
		return genErr.NewError(err, core.ErrInvalidToken, "msg", "error while ParsingToken")
	}

	rs := models.RefreshSession{
		ProfileID:      payload.ProfileID,
		RefreshTokenID: payload.TokenID,
		IssuedAt:       payload.IssuedAt,
		ExpiresIn:      payload.ExpiredAt,
	}

	err = s.TokenSession().Create(&rs)
	if err != nil {
		return genErr.NewError(err, core.ErrSession, "msg", "failed to create session")
	}

	return nil
}

func DeleteSessions(s *store.Store, uid uuid.UUID) error {
	err := s.TokenSession().DeleteUserSessions(uid)
	if err != nil {
		return genErr.NewError(err, core.ErrSession, "msg", "failed to delete session")
	}

	return nil
}

func GetSessionByRTokenID(s *store.Store, uid uuid.UUID) (*models.RefreshSession, error) {
	result, err := s.TokenSession().GetSessionByRTokenID(uid)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrSession, "msg", "can't get session by refresh tokenID")
	}

	return result, nil
}

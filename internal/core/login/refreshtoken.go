package login

import (
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/core/auth"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"git.friends.com/PetLand/UserService/v2/internal/store"
	"time"
)

const ErrCreatingAToken = "error while creating access token"

func Refresh(refToken string, s *store.Store) (*models.Tokens, error) {
	payload, err := auth.ParseToken(refToken)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrInvalidToken)
	}

	currentSession, err := s.TokenSession().GetSessionByRTokenID(payload.RtokenID)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrSession, "msg", "can't get session by refresh token ID")
	}

	tokenExpired := time.Now().After(currentSession.ExpiresIn)
	if tokenExpired {
		return nil, genErr.NewError(nil, core.ErrTokenExpired)
	}

	maker, err := auth.NewJWTMaker()
	if err != nil {
		return nil, genErr.NewError(err, core.ErrInvalidToken, "msg", core.ErrCreatingNewJWTMaker.Error())
	}
	aToken, _, err := auth.Maker.CreateToken(maker, payload.ProfileID, time.Minute*accessLifeTime, time.Minute*refreshLifeTime, payload.RtokenID) //nolint:lll
	if err != nil {
		return nil, genErr.NewError(err, core.ErrInvalidToken, "msg", ErrCreatingAToken)
	}
	pair := models.Tokens{
		AccessToken:  aToken,
		RefreshToken: "",
	}

	return &pair, nil
}

func RefreshAccess(refToken string, s *store.Store) (*models.Tokens, error) {
	// Refreshing automatic
	maker, err := auth.NewJWTMaker()
	if err != nil {
		return nil, genErr.NewError(err, core.ErrCreatingNewJWTMaker, "msg", core.ErrCreatingNewJWTMaker.Error())
	}
	payload, _ := maker.ParseExpiredToken(refToken)
	aToken, _, err := auth.Maker.CreateToken(maker, payload.ProfileID, time.Minute*accessLifeTime, time.Minute*refreshLifeTime, payload.RtokenID) //nolint:lll
	if err != nil {
		return nil, genErr.NewError(err, core.ErrInvalidToken, "msg", ErrCreatingAToken)
	}
	pair := models.Tokens{
		AccessToken:  aToken,
		RefreshToken: "",
	}

	return &pair, nil
}

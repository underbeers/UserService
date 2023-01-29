package login

import (
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/core/auth"
	"git.friends.com/PetLand/UserService/v2/internal/core/register"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"git.friends.com/PetLand/UserService/v2/internal/store"
	"github.com/google/uuid"
	"time"
)

const (
	accessLifeTime  = 10 // access token lifetime in minutes
	refreshLifeTime = 30 // refresh token lifetime in minutes
)

func Login(req *models.Login, s *store.Store) (*models.Tokens, error) {
	userUUID, err := GetUserUUID(req, s)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrRepository, "error while getting user UUID")
	}

	user, err := s.UserData().GetByUserID(userUUID)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrRepository, "UUID", userUUID)
	}

	passwordValid := register.ComparePassword(user.PasswordEncoded, req.Password, []byte(user.PasswordSalt))
	if !passwordValid {
		return nil, genErr.NewError(nil, core.ErrBadCredentials)
	}

	pair, err := GenerateTokenPair(user.ProfileID)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrInvalidToken, "msg", "error while generating token pair")
	}

	err = auth.StoreSession(s, pair.RefreshToken)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrSession, "msg", "error while store session")
	}

	return pair, nil
}

func GetUserUUID(req *models.Login, s *store.Store) (uuid.UUID, error) {
	user, err := s.Contacts().GetByEmail(req.Login)
	if err != nil {
		return uuid.Nil, genErr.NewError(err, core.ErrRepository, "error while trying login by email")
	}

	return user.ProfileID, nil
}

func GenerateTokenPair(userUUID uuid.UUID) (*models.Tokens, error) {
	maker, err := auth.NewJWTMaker()
	if err != nil {
		return nil, genErr.NewError(err, core.ErrInvalidToken, "msg", core.ErrCreatingNewJWTMaker)
	}

	rToken, rTokenID, err := auth.Maker.CreateToken(maker, userUUID,
		time.Minute*refreshLifeTime, time.Minute*accessLifeTime, uuid.Nil)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrInvalidToken, "msg", "error while creating refresh token")
	}

	aToken, _, err := auth.Maker.CreateToken(maker, userUUID,
		time.Minute*accessLifeTime, time.Minute*refreshLifeTime, rTokenID)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrInvalidToken, "msg", "error while creating access token")
	}

	return &models.Tokens{AccessToken: aToken, RefreshToken: rToken}, nil
}

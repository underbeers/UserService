package auth

import (
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"os"
	"time"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker() (Maker, error) {
	return &JWTMaker{os.Getenv("SECRET_JWT")}, nil
}

func (maker *JWTMaker) CreateToken(profileID uuid.UUID, duration time.Duration,
	durationRToken time.Duration, rTokenID uuid.UUID) (string, uuid.UUID, error) {
	payload, err := NewPayload(profileID, duration, durationRToken, rTokenID)
	if err != nil {
		return "", uuid.Nil, genErr.NewError(err, core.ErrInvalidToken, "msg", "can't create payload for the token.")
	}
	if rTokenID == uuid.Nil {
		payload.RtokenID = payload.TokenID
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", uuid.Nil, genErr.NewError(err, core.ErrInvalidToken, "msg", "can't sign jwt token.")
	}

	return token, payload.RtokenID, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, genErr.NewError(nil, core.ErrInvalidToken, "msg", "SigningMethodHMAC")
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, _ := err.(*jwt.ValidationError)
		if verr.Errors == jwt.ValidationErrorClaimsInvalid {
			_, ok := jwtToken.Claims.(*Payload)
			if !ok {
				return nil, genErr.NewError(err, core.ErrTokenExpired)
			}

			return nil, genErr.NewError(err, core.ErrInvalidToken)
		}

		return nil, genErr.NewError(err, core.ErrInvalidToken)
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, genErr.NewError(nil, core.ErrInvalidToken)
	}

	return payload, nil
}

func ParseToken(token string) (*Payload, error) {
	maker, err := NewJWTMaker()
	if err != nil {
		return nil, genErr.NewError(err, core.ErrInvalidToken, "msg", "error while store refresh session")
	}

	payload, err := maker.VerifyToken(token)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrInvalidToken, "msg", "error while VerifyToken")
	}

	return payload, nil
}

func (maker *JWTMaker) ParseExpiredToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, genErr.NewError(nil, core.ErrInvalidToken)
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, _ := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, genErr.NewError(nil, core.ErrInvalidToken)
	}

	return payload, nil
}

package signup

import (
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateUser(u *models.UserEx) error {
	validate := validator.New()
	type fields struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=6,max=255"`
	}
	f := &fields{
		Email:    u.Contacts.Email,
		Password: u.Data.PasswordEncoded,
	}

	err := validate.Struct(f)
	if err != nil {
		return genErr.NewError(err, core.ErrInvalidData)
	}

	return nil
}

func ValidateCharset(pwd string) error {
	charset := []string{"[A-Z]+", "[a-z]+", "[0-9]+"}
	for _, v := range charset {
		res, err := regexp.MatchString(v, pwd)
		if err != nil {
			return genErr.NewError(err, core.ErrValidationFailure)
		}
		if !res {
			return core.ErrInvalidPassword
		}
	}

	return nil
}

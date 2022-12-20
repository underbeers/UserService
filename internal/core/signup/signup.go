package signup

import (
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"git.friends.com/PetLand/UserService/v2/internal/store"
)

func SignUp(s *store.Store, u *models.UserEx) error {
	tx, err := s.BeginTransaction()
	if err != nil {
		return genErr.NewError(err, core.ErrRepository)
	}
	err = s.Profile().CreateNewTx(tx, &u.Profile)
	if err != nil {
		return genErr.NewError(err, core.ErrRepository)
	}
	u.Data.ProfileID = u.ID
	u.Contacts.ProfileID = u.ID
	err = s.UserData().CreateTx(tx, u.Data)
	if err != nil {
		return genErr.NewError(err, core.ErrRepository)
	}
	err = s.Contacts().CreateTx(tx, &u.Contacts.Contacts)
	if err != nil {
		return genErr.NewError(err, core.ErrRepository)
	}

	err = s.CommitTransaction(tx)
	if err != nil {
		return genErr.NewError(err, core.ErrRepository)
	}

	return nil
}

func CheckIfSigned(s *store.Store, u *models.UserEx) (bool, error) {
	check, err := s.Contacts().CheckIfSigned(&u.Contacts.Contacts)
	if err != nil {
		return false, genErr.NewError(err, core.ErrRepository)
	}
	return check, nil
}

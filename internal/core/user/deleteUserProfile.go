package user

import (
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/store"
	"github.com/google/uuid"
)

func DeleteUserProfile(profileID string, s *store.Store) error {
	id, err := uuid.Parse(profileID)
	if err != nil {
		return genErr.NewError(err, core.ErrParseUUID)
	}

	err = s.UserData().Delete(id)
	if err != nil {
		return genErr.NewError(err, core.ErrSession, "msg", "failed to delete user data")
	}
	err = s.Contacts().Delete(id)
	if err != nil {
		return genErr.NewError(err, core.ErrSession, "msg", "failed to delete user data")
	}
	err = s.Profile().Delete(id)
	if err != nil {
		return genErr.NewError(err, core.ErrSession, "msg", "failed to delete user profile")
	}

	return nil
}

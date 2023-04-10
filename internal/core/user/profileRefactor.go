package user

import (
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/core/register"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"git.friends.com/PetLand/UserService/v2/internal/store"
	"github.com/google/uuid"
)

func ChangePassword(data *models.Data, newPassword string, store *store.Store) error {
	data.PasswordEncoded = newPassword

	if err := register.EncryptPassword(data); err != nil {
		return genErr.NewError(err, core.ErrEncryptPassword)
	}

	err := store.UserData().ChangePassword(data.ProfileID, data.PasswordEncoded, data.PasswordSalt)
	if err != nil {
		return genErr.NewError(err, core.ErrRepository, "msg", "failed to change password")
	}

	return nil
}

func ChangeChatID(userID uuid.UUID, chatID string, store *store.Store) error {
	if err := store.Contacts().ChangeChatID(userID, chatID); err != nil {
		return genErr.NewError(err, core.ErrRepository)
	}

	return nil
}

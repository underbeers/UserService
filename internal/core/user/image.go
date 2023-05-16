package user

import (
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/store"
	"github.com/google/uuid"
)

func SetImage(userID uuid.UUID, imageLink string, store *store.Store) error {
	if err := store.Profile().SetImage(userID, imageLink); err != nil {
		return genErr.NewError(err, core.ErrRepository)
	}

	return nil
}

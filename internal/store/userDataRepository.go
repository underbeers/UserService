package store

import (
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserDataRepository struct {
	store *Store
}

type UserDater interface {
	Create(d *models.Data) error
	CreateTx(tx *sqlx.Tx, d *models.Data) error
}

func (r *UserDataRepository) Create(d *models.Data) error {
	return r.CreateTx(nil, d)
}

func (r *UserDataRepository) CreateTx(tx *sqlx.Tx, d *models.Data) error {
	if _, err := r.store.db.Exec(
		tx,
		"INSERT INTO  public.user_data (id_profile, password_encoded, password_salt) VALUES ($1, $2, $3);",
		d.ProfileID,
		d.PasswordEncoded,
		d.PasswordSalt,
	); err != nil {
		return r.store.Rollback(tx, err)
	}

	return nil
}

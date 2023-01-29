package store

import (
	"database/sql"
	"errors"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserDataRepository struct {
	store *Store
}

type UserDater interface {
	Create(d *models.Data) error
	CreateTx(tx *sqlx.Tx, d *models.Data) error
	GetByUserID(id uuid.UUID) (*models.Data, error)
}

func (r *UserDataRepository) Create(d *models.Data) error {
	return r.CreateTx(nil, d)
}

func (r *UserDataRepository) CreateTx(tx *sqlx.Tx, d *models.Data) error {
	if _, err := r.store.db.Exec(
		tx,
		"INSERT INTO user_service.public.user_data (id_profile, password_encoded, password_salt) VALUES ($1, $2, $3);",
		d.ProfileID,
		d.PasswordEncoded,
		d.PasswordSalt,
	); err != nil {
		return r.store.Rollback(tx, err)
	}

	return nil
}

func (r *UserDataRepository) GetByUserID(id uuid.UUID) (*models.Data, error) {
	return r.getByUserIDTx(nil, id)
}

func (r *UserDataRepository) getByUserIDTx(tx *sqlx.Tx, id uuid.UUID) (*models.Data, error) {
	data := &models.Data{}
	row := r.store.db.QueryRow(
		tx, `SELECT id, id_profile, password_encoded, password_salt 
FROM user_service.public.user_data WHERE id_profile = $1`, id)

	err := row.StructScan(data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = genErr.NewError(err, ErrNotFound, "UUID", id)
		}
		err = genErr.NewError(err, ErrScanStructFailed)

		return nil, r.store.Rollback(tx, err)
	}

	return data, nil
}

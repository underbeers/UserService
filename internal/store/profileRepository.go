package store

import (
	"errors"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProfileRepository struct {
	store *Store
}

type Profiler interface {
	CreateNewTx(tx *sqlx.Tx, c *models.Profile) error
	CreateNew(c *models.Profile) error
	GetByUserID(id uuid.UUID) (*models.Profile, error)
	GetByUserIDTx(tx *sqlx.Tx, id uuid.UUID) (*models.Profile, error)
	Delete(id uuid.UUID) error
	DeleteTx(tx *sqlx.Tx, id uuid.UUID) error
}

type Profile struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	SurName   string    `db:"sur_name"`
	Status    string    `db:"status"`
}

func (r *ProfileRepository) CreateNew(c *models.Profile) error {
	return r.CreateNewTx(nil, c)
}

func (r *ProfileRepository) CreateNewTx(tx *sqlx.Tx, c *models.Profile) error {
	c.ID = uuid.New()
	if _, err := r.store.db.Exec(
		tx,
		`INSERT INTO user_service.public.user_profile (id, first_name, sur_name, status) VALUES ($1, $2, $3, $4);`,
		c.ID,
		c.FirstName,
		c.SurName,
		c.Status,
	); err != nil {
		return r.store.Rollback(tx, err)
	}

	return nil
}

func (r *ProfileRepository) GetByUserID(id uuid.UUID) (*models.Profile, error) {
	return r.GetByUserIDTx(nil, id)
}

func (r *ProfileRepository) GetByUserIDTx(tx *sqlx.Tx, id uuid.UUID) (*models.Profile, error) {
	profile := &models.Profile{}
	row := r.store.db.QueryRow(tx,
		`SELECT id, first_name, sur_name, status FROM user_service.public.user_profile WHERE id = $1`, id)

	err := row.StructScan(profile)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			err = genErr.NewError(err, ErrNotFound, "id", id)

			return nil, err
		}
		err = genErr.NewError(err, ErrScanStructFailed)

		return nil, r.store.Rollback(tx, err)
	}

	return profile, nil
}

func (r *ProfileRepository) Delete(id uuid.UUID) error {
	return r.DeleteTx(nil, id)
}

func (r *ProfileRepository) DeleteTx(tx *sqlx.Tx, id uuid.UUID) error {
	if err := r.store.db.QueryRow(tx, `DELETE FROM user_service.public.user_profile WHERE id=$1`, id).Err(); err != nil {
		return r.store.Rollback(tx, err)
	}

	return nil
}

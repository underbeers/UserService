package store

import (
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

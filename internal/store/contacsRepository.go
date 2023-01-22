package store

import (
	"database/sql"
	"errors"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"github.com/jmoiron/sqlx"
)

type ContactsRepository struct {
	store *Store
}

type Contacter interface {
	Create(c *models.Contacts) error
	CreateTx(tx *sqlx.Tx, c *models.Contacts) error
	CheckIfSigned(c *models.Contacts) (bool, error)
	GetByEmail(email string) (*models.Contacts, error)
	GetByEmailTx(tx *sqlx.Tx, email string) (*models.Contacts, error)
}

func (r *ContactsRepository) Create(c *models.Contacts) error {
	return r.CreateTx(nil, c)
}

func (r *ContactsRepository) CreateTx(tx *sqlx.Tx, c *models.Contacts) error {
	if err := r.store.db.QueryRow(
		tx,
		`INSERT INTO user_service.public.user_contacts (id_profile, email, mobile_phone, 
			email_subscription, show_phone) VALUES ($1, $2, $3, $4, $5);`,
		c.ProfileID,
		c.Email,
		c.MobilePhone,
		c.EmailSubscription,
		c.ShowPhone,
	).Err(); err != nil {
		return r.store.Rollback(tx, err)
	}

	return nil
}

func (r *ContactsRepository) CheckIfSigned(c *models.Contacts) (bool, error) {
	var signed bool
	row := r.store.db.QueryRow(nil,
		`SELECT EXISTS (SELECT * FROM user_service.public.user_contacts WHERE email = $1)`,
		c.Email)
	if err := row.Scan(&signed); err != nil {
		err = genErr.NewError(err, ErrScanStructFailed)

		return true, err
	}

	if signed {
		return true, nil
	}

	return false, nil
}

func (r *ContactsRepository) GetByEmail(email string) (*models.Contacts, error) {
	return r.GetByEmailTx(nil, email)
}

func (r *ContactsRepository) GetByEmailTx(tx *sqlx.Tx, email string) (*models.Contacts, error) {
	contacts := &models.Contacts{}
	row := r.store.db.QueryRow(tx,
		`SELECT id, id_profile, email, mobile_phone, show_phone 
			FROM user_service.public.user_contacts WHERE email = $1`, email)

	err := row.StructScan(contacts)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = genErr.NewError(err, ErrNotFound, "email", email)
		}
		err = genErr.NewError(err, ErrScanStructFailed)

		return nil, r.store.Rollback(tx, err)
	}

	return contacts, nil
}

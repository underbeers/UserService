package store

import (
	"database/sql"
	"errors"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"github.com/google/uuid"
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
	GetByUserProfileID(id uuid.UUID) (*models.Contacts, error)
	GetByUserProfileIDTx(tx *sqlx.Tx, id uuid.UUID) (*models.Contacts, error)
	GetByHashID(hash string) (*models.Contacts, error)
	GetByHashIDTx(tx *sqlx.Tx, hash string) (*models.Contacts, error)
	InsertHashID(profileID uuid.UUID, hash string) error
	InsertHashIDTx(tx *sqlx.Tx, profileID uuid.UUID, hash string) error
	Delete(profileID uuid.UUID) error
	DeleteTx(tx *sqlx.Tx, profileID uuid.UUID) error
}

func (r *ContactsRepository) Create(c *models.Contacts) error {
	return r.CreateTx(nil, c)
}

func (r *ContactsRepository) CreateTx(tx *sqlx.Tx, c *models.Contacts) error {
	if err := r.store.db.QueryRow(
		tx,
		`INSERT INTO user_contacts (id_profile, email, mobile_phone, 
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
			FROM user_contacts WHERE email = $1`, email)

	err := row.StructScan(contacts)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = genErr.NewError(err, ErrNotFound, "email", email)

			return nil, err
		}
		err = genErr.NewError(err, ErrScanStructFailed)

		return nil, r.store.Rollback(tx, err)
	}

	return contacts, nil
}

func (r *ContactsRepository) GetByUserProfileID(id uuid.UUID) (*models.Contacts, error) {
	return r.GetByUserProfileIDTx(nil, id)
}

func (r *ContactsRepository) GetByUserProfileIDTx(tx *sqlx.Tx, id uuid.UUID) (*models.Contacts, error) {
	contacts := &models.Contacts{}
	row := r.store.db.QueryRow(tx,
		`SELECT id, id_profile, email, mobile_phone, show_phone 
			FROM user_contacts WHERE id_profile = $1`, id)

	err := row.StructScan(contacts)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = genErr.NewError(err, ErrNotFound, "id", id)

			return nil, err
		}
		err = genErr.NewError(err, ErrScanStructFailed)

		return nil, r.store.Rollback(tx, err)
	}

	return contacts, nil
}

func (r *ContactsRepository) GetByHashID(hash string) (*models.Contacts, error) {
	return r.GetByHashIDTx(nil, hash)
}

func (r *ContactsRepository) GetByHashIDTx(tx *sqlx.Tx, hash string) (*models.Contacts, error) {
	contacts := &models.Contacts{}
	row := r.store.db.QueryRow(tx,
		`SELECT id, id_profile, email, mobile_phone, show_phone, hash_id 
			FROM user_contacts WHERE hash_id = $1`, hash)

	err := row.StructScan(contacts)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = genErr.NewError(err, ErrNotFound, "hash_id", hash)

			return nil, err
		}
		err = genErr.NewError(err, ErrScanStructFailed)

		return nil, r.store.Rollback(tx, err)
	}

	return contacts, nil
}

func (r *ContactsRepository) InsertHashID(profileID uuid.UUID, hash string) error {
	return r.InsertHashIDTx(nil, profileID, hash)
}

func (r *ContactsRepository) InsertHashIDTx(tx *sqlx.Tx, profileID uuid.UUID, hash string) error {
	if err := r.store.db.QueryRow(
		tx,
		`UPDATE user_contacts SET hash_id = $1 WHERE id_profile = $2`, hash, profileID).Err(); err != nil {
		return r.store.Rollback(tx, err)
	}

	return nil
}

func (r *ContactsRepository) Delete(profileID uuid.UUID) error {
	return r.DeleteTx(nil, profileID)
}

func (r *ContactsRepository) DeleteTx(tx *sqlx.Tx, profileID uuid.UUID) error {
	if err := r.store.db.QueryRow(tx,
		`DELETE FROM user_contacts WHERE id_profile=$1`, profileID,
	).Err(); err != nil {
		return r.store.Rollback(tx, err)
	}

	return nil
}

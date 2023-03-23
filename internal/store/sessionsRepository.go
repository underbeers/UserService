package store

import (
	"database/sql"
	"errors"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	store *Store
}

type Sessioner interface {
	Create(rs *models.RefreshSession) error
	CreateTx(tx *sqlx.Tx, rs *models.RefreshSession) error
	GetSessionByRTokenID(uid uuid.UUID) (*models.RefreshSession, error)
	DeleteUserSessions(uid uuid.UUID) error
	DeleteUserSessionsTx(tx *sqlx.Tx, uid uuid.UUID) error
}

func (s *SessionRepository) Create(rs *models.RefreshSession) error {
	return s.CreateTx(nil, rs)
}

func (s *SessionRepository) CreateTx(tx *sqlx.Tx, rs *models.RefreshSession) error {
	if _, err := s.store.db.Exec(
		tx,
		`INSERT INTO refresh_sessions
        (id_profile, id_refresh_token, issued_at, expires_in)
        VALUES ($1, $2, $3, $4);`,
		rs.ProfileID,
		rs.RefreshTokenID,
		rs.IssuedAt,
		rs.ExpiresIn,
	); err != nil {
		return s.store.Rollback(tx, err)
	}

	return nil
}

func (s *SessionRepository) GetSessionByRTokenID(uid uuid.UUID) (*models.RefreshSession, error) {
	return s.getSessionByTokenIDTx(nil, uid)
}

func (s *SessionRepository) getSessionByTokenIDTx(tx *sqlx.Tx, uid uuid.UUID) (*models.RefreshSession, error) {
	session := &models.RefreshSession{}

	row := s.store.db.QueryRow(tx, `SELECT id_profile, id_refresh_token, issued_at, expires_in 
		FROM refresh_sessions WHERE id_refresh_token=$1;`, uid)

	err := row.StructScan(session)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = genErr.NewError(err, ErrNotFound, "UUID", uid)
		}
		err = genErr.NewError(err, ErrScanStructFailed)

		return nil, s.store.Rollback(tx, err)
	}

	return session, nil
}

func (s *SessionRepository) DeleteUserSessions(uid uuid.UUID) error {
	return s.DeleteUserSessionsTx(nil, uid)
}

func (s *SessionRepository) DeleteUserSessionsTx(tx *sqlx.Tx, uid uuid.UUID) error {
	_, err := s.store.db.Exec(tx,
		`DELETE FROM refresh_sessions WHERE id_profile = $1`, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = genErr.NewError(err, ErrNotFound, "UUID", uid)
		}
		err = genErr.NewError(err, ErrScanStructFailed)

		return s.store.Rollback(tx, err)
	}

	return nil
}

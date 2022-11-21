package store

import (
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Store struct {
	db                     *DB
	logger                 *zap.SugaredLogger
	Itx                    TX
	profileRepository      Profiler
	contactsRepository     Contacter
	userDataRepository     UserDater
	expertDataRepository   ExpertDater
	verificationRepository Verificationer
}

type TX interface {
	Rollback(tx *sqlx.Tx, err error) error
	BeginTransaction() (*sqlx.Tx, error)
	CommitTransaction(tx *sqlx.Tx) error
}

func New(db *sqlx.DB, l *zap.SugaredLogger) *Store {
	return &Store{
		db:     &DB{db},
		logger: l,
	}
}

func (s *Store) UserData() UserDater {
	if s.userDataRepository != nil {
		return s.userDataRepository
	}

	s.userDataRepository = &UserDataRepository{
		store: s,
	}

	return s.userDataRepository
}

func (s *Store) ExpertData() ExpertDater {
	if s.expertDataRepository != nil {
		return s.expertDataRepository
	}

	s.expertDataRepository = &ExpertDataRepository{
		store: s,
	}

	return s.expertDataRepository
}

func (s *Store) Profile() Profiler {
	if s.profileRepository != nil {
		return s.profileRepository
	}

	s.profileRepository = &ProfileRepository{
		store: s,
	}

	return s.profileRepository
}

func (s *Store) Contacts() Contacter {
	if s.contactsRepository != nil {
		return s.contactsRepository
	}

	s.contactsRepository = &ContactsRepository{
		store: s,
	}

	return s.contactsRepository
}

func (s *Store) Verification() Verificationer {
	if s.verificationRepository != nil {
		return s.verificationRepository
	}

	s.verificationRepository = &VerificationRepository{
		store: s,
	}

	return s.verificationRepository
}

func (s *Store) BeginTransaction() (*sqlx.Tx, error) {
	if s.Itx != nil {
		return s.Itx.BeginTransaction()
	} else {
		tx, err := s.db.sqlxDB.Beginx()
		if err != nil {
			s.logger.Error(ErrTransactionCreationFailed)

			return nil, genErr.NewError(err, ErrTransactionCreationFailed)
		}

		return tx, nil
	}
}

func (s *Store) CommitTransaction(tx *sqlx.Tx) error {
	if s.Itx != nil {
		return s.Itx.CommitTransaction(tx)
	} else {
		if err := tx.Commit(); err != nil {
			// s.logger.Errorf(wErr.WrapError(ErrTransactionCommitFailed.Error(), err).Error())
			s.logger.Error(err)
			if err = s.Rollback(tx, err); err != nil {
				return err
			}
		}

		return nil
	}
}

func (s *Store) Rollback(tx *sqlx.Tx, err error) error {
	if s.Itx != nil {
		return s.Itx.Rollback(tx, err)
	} else {
		if tx == nil {
			return err
		}

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return genErr.NewError(err, ErrTransactionCreationFailed, nil)
		}

		return genErr.NewError(err, ErrTransactionFailed, nil)
	}
}

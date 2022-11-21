package store

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	sqlxDB *sqlx.DB
}

func (w *DB) Query(tx *sqlx.Tx, query string, args ...any) (*sqlx.Rows, error) {
	if tx == nil {
		return w.sqlxDB.Queryx(query, args...) //nolint:wrapcheck
	} else {
		return tx.Queryx(query, args...) //nolint:wrapcheck
	}
}

func (w *DB) QueryRow(tx *sqlx.Tx, query string, args ...any) *sqlx.Row {
	if tx == nil {
		return w.sqlxDB.QueryRowx(query, args...)
	} else {
		return tx.QueryRowx(query, args...)
	}
}

func (w *DB) Exec(tx *sqlx.Tx, query string, args ...any) (sql.Result, error) {
	if tx == nil {
		return w.sqlxDB.Exec(query, args...) //nolint:wrapcheck
	} else {
		return tx.Exec(query, args...) //nolint:wrapcheck
	}
}

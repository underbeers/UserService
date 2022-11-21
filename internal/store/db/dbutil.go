package db

import (
	"database/sql"
	"git.friends.com/PetLand/UserService/v2/internal/config"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

func NewDB(config *config.Config, logger *zap.Logger) (*sqlx.DB, error) {
	// parse connection string
	dbConf, err := pgx.ParseConfig(config.DB.GetConnectionString())
	if err != nil {
		return nil, genErr.NewError(err, genErr.New("UnavailableResource"), "msg", "failed to parse config")
	}

	dbConf.Logger = zapadapter.NewLogger(logger)
	dbConf.Host = config.DB.Host

	//register pgx conn
	dsn := stdlib.RegisterConnConfig(dbConf)

	sql.Register("wrapper", stdlib.GetDefaultDriver())
	wdb, err := sql.Open("wrapper", dsn)
	if err != nil {
		return nil, genErr.NewError(err, genErr.New("UnavailableResource"), "msg", "failed to connect to database")
	}

	const (
		maxOpenConns    = 50
		maxIdleConns    = 50
		connMaxLifetime = 3
		connMaxIdleTime = 5
	)
	db := sqlx.NewDb(wdb, "pgx")
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Minute)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, genErr.NewError(err, genErr.New("Unavailable resource"), "msg", "failed to connect to database")
	}

	return db, nil
}

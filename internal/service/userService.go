package service

import (
	"context"
	"fmt"
	"git.friends.com/PetLand/UserService/v2/internal/config"
	"git.friends.com/PetLand/UserService/v2/internal/store/db"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
)

type service struct {
	Logger *zap.SugaredLogger
	conf   *config.Config
	router *mux.Router
	//store *store.Store
}

const (
	baseURL   = "/api/v1/"
	requestID = "X-request-ID"
	attempts  = 3 // The number of attempts we're trying to say Hello() to ApiGateway
	timeout   = 5 // Timeout in Seconds, how long we wait to reconnect to ApiGateway
)

func NewService(cfg *config.Config) *service {
	logger := NewLogger(cfg.DebugMode)
	srv := &service{Logger: logger, conf: cfg, router: mux.NewRouter()}
	srv.registerHandlers()

	return srv
}

func (srv *service) Start() error {
	database, err := db.NewDB(srv.conf, srv.Logger.Desugar())
	if err != nil {
		srv.Logger.Fatal("Can't initialize connection to database", zap.Error(err))

		return fmt.Errorf("failed to initialize connection to database: %w", err)
	}
	checkMigrationVersion(srv, database)
	//srv.store = store.New(database, srv.Logger)
	/*var errorCnt int
	var er error
	for errorCnt < attempts {
		time.Sleep(time.Second * timeout)
		srv.Logger.Infof("Attempt to connect to APIGateway %d of %d", errorCnt+1, attempts)
		if err = pingAPIGateway(srv); err != nil {
			errorCnt++
			er = err
		} else {
			break
		}
	}
	if errorCnt >= attempts {
		return genErr.NewError(er, ErrConnectAPIGateWay, msg, "failed to send info to the APIGateway")
	}
	//err = HelloAPIGateway(srv)
	if err != nil {
		return err
	}
	srv.Logger.Info("Start to listen to", zap.String("port", srv.conf.Listen.Port))
	*/

	return fmt.Errorf("failed to listen and serve: %w", http.ListenAndServe(":"+srv.conf.Listen.Port, srv.router))
}

func (srv *service) registerHandlers() {
	srv.router.Use(srv.getRequestID)
	srv.registerClientHandlers()
}

func (srv *service) DebugModeRunning() bool {
	return srv.conf.DebugMode
}

func (srv *service) getRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type xRequestIDKey string
		xRequestID := xRequestIDKey(requestID)
		id := r.Header.Get(requestID)
		ctx := context.WithValue(r.Context(), xRequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkMigrationVersion(srv *service, db *sqlx.DB) {
	type MigrationSchema struct {
		Version int  `db:"version"`
		Dirty   bool `db:"dirty"`
	}
	/*result := db.QueryRowx("SELECT version, dirty FROM schema_migrations;")
	if result.Err() != nil {
		srv.Logger.Fatal("failed to connect to database")
	}
	migration := &MigrationSchema{}
	err := result.StructScan(migration)
	if err != nil {
		srv.Logger.Fatal("failed to scan schema migrations:", zap.Error(err))
	}
	if migration.Dirty {
		srv.Logger.Fatal("schema migration is in dirty mode")
	}
	if srv.conf.VersionDB != migration.Version {
		srv.Logger.Fatalf("Mismatched db versions. Expected: %d, got: %d", srv.conf.VersionDB, migration.Version)
	}
	*/
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"git.friends.com/PetLand/UserService/v2/internal/config"
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"git.friends.com/PetLand/UserService/v2/internal/store"
	"git.friends.com/PetLand/UserService/v2/internal/store/db"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type service struct {
	Logger *zap.SugaredLogger
	conf   *config.Config
	router *mux.Router
	store  *store.Store
}

const (
	baseURL   = "/api/v1/"
	requestID = "X-request-ID"
	attempts  = 3 // The number of attempts we're trying to say Hello() to ApiGateway
	timeout   = 5 // Timeout in Seconds, how long we wait to reconnect to ApiGateway
)

func NewService(cfg *config.Config) *service {
	logger := NewLogger()
	srv := &service{Logger: logger, conf: cfg, router: mux.NewRouter()}
	srv.registerHandlers()

	return srv
}

func (srv *service) Start() error {

	srv.Logger.Infof("Start to listen to port %s", srv.conf.Listen.Port)
	srv.Logger.Info("Config variables: ")
	srv.Logger.Infof("GATEWAY_PORT: %v", srv.conf.Gateway.Port)
	srv.Logger.Infof("GATEWAY_IP: %v", srv.conf.Gateway.IP)

	database, err := db.NewDB(srv.conf, srv.Logger.Desugar())
	if err != nil {
		srv.Logger.Fatal("Can't initialize connection to database", zap.Error(err))

		return fmt.Errorf("failed to initialize connection to database: %w", err)
	}
	checkMigrationVersion(srv, database)
	srv.store = store.New(database, srv.Logger)
	var errorCnt int
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
	err = HelloAPIGateway(srv)
	if err != nil {
		return err
	}
	srv.Logger.Info("Start to listen to", zap.String("port", srv.conf.Listen.Port))

	return fmt.Errorf("failed to listen and serve: %w", http.ListenAndServe(":"+srv.conf.Listen.Port, srv.router))
}

func (srv *service) registerHandlers() {
	srv.router.Use(srv.getRequestID)
	srv.registerClientHandlers()
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
	result := db.QueryRowx("SELECT version, dirty FROM schema_migrations;")
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
}

func GetServiceInfo(srv *service) *config.Service {
	handles, err := getHandles(srv)
	if err != nil {
		srv.Logger.Fatalf("failed to getHandles, %v", err)
	}

	instance := config.Service{
		Name:      "user",
		Label:     "pl_user_service",
		IP:        srv.conf.Listen.IP,
		Port:      srv.conf.Listen.Port,
		Endpoints: nil,
	}
	unprotected, err := getUnprotected()
	if err != nil {
		srv.Logger.Fatalf("failed to getUnprotected, %v", err)
	}

	for k, v := range handles {
		// skip endpoint-info
		if k == "endpoint-info/" {
			continue
		}
		endpoint := models.Endpoint{
			URL:       k,
			Protected: true,
			Methods:   v,
		}
		if unprotected[k] {
			endpoint.Protected = false
		}
		instance.Endpoints = append(instance.Endpoints, endpoint)
	}

	return &instance
}

func getHandles(srv *service) (map[string][]string, error) {
	data := make(map[string][]string)
	err := srv.router.Walk(
		func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, _ := route.GetPathTemplate()
			n, _ := route.GetMethods()
			path = strings.Split(path, "/api/v1/")[1]
			d, ok := data[path]
			if ok {
				n = append(n, d...)
				data[path] = n

				return nil
			}
			data[path] = n

			return nil
		})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getUnprotected() (map[string]bool, error) {
	// Read's list of unprotected endpoints
	lst, err := os.OpenFile("service.json", os.O_RDONLY, 0o600) //nolint:gomnd
	if err != nil {
		return nil, genErr.NewError(err, ErrOpenFile)
	}
	reader, err := io.ReadAll(lst)
	if err != nil {
		return nil, genErr.NewError(err, ErrReadFile)
	}
	data := struct {
		URLS []string `json:"urls"`
	}{}
	err = json.Unmarshal(reader, &data)
	if err != nil {
		return nil, genErr.NewError(err, core.ErrMarshalUnmarshal)
	}
	result := make(map[string]bool)
	for _, k := range data.URLS {
		result[k] = true
	}

	return result, nil
}

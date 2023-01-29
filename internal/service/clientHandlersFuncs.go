package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"git.friends.com/PetLand/UserService/v2/internal/config"
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/core/login"
	"git.friends.com/PetLand/UserService/v2/internal/core/register"
	"git.friends.com/PetLand/UserService/v2/internal/core/signup"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"github.com/google/uuid"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	protocol   = "http"
	msg        = "msg"
	keyValPair = 2
	userIDAuth = "UserID"
)

func (srv *service) registerClientHandlers() {
	srv.router.HandleFunc(baseURL+"helloMessage/", srv.handleHelloMessage()).Methods(http.MethodGet)
	srv.router.HandleFunc(baseURL+"registration/new/", srv.handleCreteNewUser()).Methods(http.MethodPost, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"login/", srv.handleLoginUser()).Methods(http.MethodPost, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"login/token/", srv.handleRefreshToken()).Methods(http.MethodGet, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"user/info/", srv.handleUserInfo()).Methods(http.MethodGet, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"endpoint-info/", srv.handleInfo()).Methods(http.MethodGet)
}

func (srv *service) handleHelloMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.respond(w, http.StatusOK, "Hello, it's work!")
	}
}

func (srv *service) handleCreteNewUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			FirstName   string `json:"firstName"`
			SurName     string `json:"surName"`
			Email       string `json:"email"`
			MobilePhone string `json:"mobilePhone"`
			Password    string `json:"password"`
		}

		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())

			return
		}

		u := &models.UserEx{
			Profile: models.Profile{
				FirstName: req.FirstName,
				SurName:   req.SurName,
				Status:    0,
			},
			Data: &models.Data{
				PasswordEncoded: req.Password,
			},
			Contacts: &models.ContactsEX{
				Contacts: models.Contacts{
					Email:             req.Email,
					MobilePhone:       req.MobilePhone,
					EmailSubscription: false,
					ShowPhone:         true,
				},
			},
		}

		signed, err := signup.CheckIfSigned(srv.store, u)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}
		if signed {
			srv.respond(w, http.StatusConflict, "User with this email already exists")

			return
		}

		if err := signup.ValidateUser(u); err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}

		if err := register.EncryptPassword(u.Data); err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}

		if err := signup.SignUp(srv.store, u); err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}

		srv.respond(w, http.StatusCreated, u.ID)
	}
}

func (srv *service) handleLoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &models.Login{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		tokens, err := login.Login(req, srv.store)
		if err != nil {
			switch {
			case errors.Is(err, core.ErrInvalidToken):
				srv.warning(w, http.StatusBadRequest, genErr.NewError(err, core.ErrInvalidToken))

				return
			}
			srv.warning(w, http.StatusBadRequest, genErr.NewError(err, core.ErrBadCredentials))

			return
		}

		tokenJSON, err := json.Marshal(models.AccessToken{AccessToken: tokens.AccessToken})
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())
		}
		_, err = w.Write(tokenJSON)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())
		}
	}
}

func (srv *service) handleRefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokens *models.Tokens
		expired := r.Header.Get("ExpiredIn")
		if len(expired) > 0 {
			exp, err := strconv.ParseInt(expired, 10, 64)
			if err != nil {
				srv.error(w, http.StatusInternalServerError, err, r.Context())
			}
			if time.Now().Before(time.Unix(exp, 0)) {
				providedToken, err := getAuthHeader(r.Header.Get("Authorization"))
				if err != nil {
					srv.warning(w, http.StatusBadRequest, genErr.NewError(nil, err))

					return
				}
				tokens, err = login.RefreshAccess(providedToken, srv.store)
				if err != nil {
					srv.error(w, http.StatusInternalServerError, err, r.Context())
				}

				return
			}
		}

		providedToken, err := getAuthHeader(r.Header.Get("Authorization"))
		if err != nil {
			srv.warning(w, http.StatusBadRequest, genErr.NewError(nil, err))

			return
		}

		tokens, err = login.Refresh(providedToken, srv.store)
		if err != nil {
			switch {
			case errors.Is(err, core.ErrTokenExpired):
				srv.warning(w, http.StatusUnauthorized, genErr.NewError(err, core.ErrTokenExpired))

				return
			case errors.Is(err, core.ErrTokenHeaderMismatch):
				srv.warning(w, http.StatusBadRequest, genErr.NewError(err, core.ErrTokenHeaderMismatch))

				return
			case errors.Is(err, core.ErrInvalidToken):
				srv.warning(w, http.StatusBadRequest, genErr.NewError(err, core.ErrInvalidToken))

				return
			}

			srv.warning(w, http.StatusBadRequest, genErr.NewError(err, core.ErrInvalidToken))

			return
		}

		err = writeJSONBody(w, tokens)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())
		}
	}
}

func (srv *service) handleUserInfo() http.HandlerFunc {
	type Response struct {
		FirstName   string
		SurName     string
		MobilePhone string
		Email       string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(userIDAuth)
		if len(userID) == 0 {
			srv.warning(w, http.StatusUnauthorized, ErrInvalidHeader)

			return
		}
		id, err := uuid.Parse(userID)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, ErrParams, r.Context())
		}
		contacts, err := srv.store.Contacts().GetByUserProfileID(id)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}
		user, err := srv.store.Profile().GetByUserID(id)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}

		resp := &Response{
			FirstName:   user.FirstName,
			SurName:     user.SurName,
			MobilePhone: contacts.MobilePhone,
			Email:       contacts.Email,
		}
		w.Header().Add("Content-Type", "application/json")
		userInfoJSON, err := json.Marshal(resp)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}
		_, err = w.Write(userInfoJSON)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}
	}
}

func getAuthHeader(header string) (string, error) {
	if len(header) == 0 {
		return "", genErr.NewError(nil, ErrNoHeader, msg, ErrAuthHeaderMissing)
	}
	providedHeader := strings.Split(header, " ")
	if len(providedHeader) != keyValPair {
		return "", genErr.NewError(nil, ErrInvalidHeader, msg, ErrInvalidHeader)
	}

	return providedHeader[1], nil
}

func writeJSONBody(w http.ResponseWriter, tokens *models.Tokens) error {
	w.Header().Add("Content-Type", "application/json")
	tokenJSON, err := json.Marshal(models.AccessToken{AccessToken: tokens.AccessToken})
	if err != nil {
		return genErr.NewError(err, ErrMarshalUnmarshal, "msg", "error while Marshal AccessToken")
	}
	_, err = w.Write(tokenJSON)
	if err != nil {
		return genErr.NewError(err, ErrWriteBody, "msg", "error while writing tokens")
	}

	return nil
}

func (srv *service) handleInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		serviceInfo := GetServiceInfo(srv)
		payload, err := json.Marshal(serviceInfo)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())
		}

		_, err = w.Write(payload)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())
		}
	}
}

func HelloAPIGateway(srv *service) error {
	var domain string

	cfg := config.ReadConfig()
	if srv.conf.DebugMode {
		domain = cfg.Gateway.IP
	} else {
		domain = cfg.Gateway.Label
	}
	gatewayURL, err := url.Parse(
		protocol + "://" + domain + ":" + cfg.Gateway.Port + baseURL + "hello/")
	if err != nil {
		return genErr.NewError(err, ErrConnectAPIGateWay, msg, "can't parse ur for endpoint 'hello/'")
	}

	//endpoints := config.ReadServicesList()

	info := &models.Hello{
		Name:      "user",
		Label:     "pl_user_service",
		IP:        cfg.Listen.IP,
		Port:      cfg.Listen.Port,
		Endpoints: nil,
	}
	jsonStr, err := json.Marshal(info)
	if err != nil {
		return genErr.NewError(err, ErrMarshal)
	}

	go knock(gatewayURL.String(), jsonStr)

	return nil
}

func knock(url string, payload []byte) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload)) //nolint: gosec, noctx
	if resp == nil {
		// FIXME:Super dirty. Need to handle error
		log.Println("can't say Hello to Gateway", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		log.Println("knock() Post Error", err)
	}
	if resp.StatusCode == http.StatusOK {
		log.Println("Successfully greet ApiGateway")
	}
}

func pingAPIGateway(srv *service) error {
	gwURL, err := gatewayURL(srv)
	if err != nil {
		return genErr.NewError(err, ErrConnectAPIGateWay, msg, "pingAPIGateway gwURL generation")
	}
	resp, err := http.Get(gwURL.String()) //nolint: noctx

	if resp == nil {
		return genErr.NewError(err, ErrConnectAPIGateWay, msg, "can't ping APIGateway")
	}
	if err != nil {
		return genErr.NewError(err, ErrConnectAPIGateWay, msg, "error pingAPIGateway http.Get(gwURL)")
	}
	if err := resp.Body.Close(); err != nil {
		log.Println(ErrCloseResponseBody.Error())
	}

	return nil
}

func gatewayURL(srv *service) (*url.URL, error) {
	var domain string
	if srv.conf.DebugMode {
		domain = srv.conf.Gateway.IP
	} else {
		domain = srv.conf.Gateway.Label
	}
	gwURL, err := url.Parse(
		protocol + "://" + domain + ":" + srv.conf.Gateway.Port + baseURL + "hello/")
	srv.Logger.Info(gwURL)
	if err != nil {
		return nil, genErr.NewError(err, ErrConnectAPIGateWay)
	}

	return gwURL, nil
}

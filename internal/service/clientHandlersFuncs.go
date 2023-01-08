package service

import (
	"encoding/json"
	"git.friends.com/PetLand/UserService/v2/internal/core/register"
	"git.friends.com/PetLand/UserService/v2/internal/core/signup"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"log"
	"net/http"
	"net/url"
)

const (
	protocol = "http"
	msg      = "msg"
)

func (srv *service) registerClientHandlers() {
	srv.router.HandleFunc(baseURL+"helloMessage/", srv.handleHelloMessage()).Methods(http.MethodGet)
	srv.router.HandleFunc(baseURL+"registration/new/", srv.handleCreteNewUser()).Methods(http.MethodPost)
	srv.router.HandleFunc(baseURL+"login/", srv.handleLoginUser()).Methods(http.MethodPost)
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

func (srv *service) handleCheckRegistration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (srv *service) handleLoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		return
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
	if err != nil {
		return nil, genErr.NewError(err, ErrConnectAPIGateWay)
	}

	return gwURL, nil
}

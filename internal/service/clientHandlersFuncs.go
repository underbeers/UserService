package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/core/login"
	"git.friends.com/PetLand/UserService/v2/internal/core/register"
	"git.friends.com/PetLand/UserService/v2/internal/core/signup"
	"git.friends.com/PetLand/UserService/v2/internal/core/user"
	"git.friends.com/PetLand/UserService/v2/internal/core/utils"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"git.friends.com/PetLand/UserService/v2/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
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
	srv.router.HandleFunc(baseURL+"helloMessage", srv.handleHelloMessage()).Methods(http.MethodGet)
	srv.router.HandleFunc(baseURL+"registration/new", srv.handleCreteNewUser()).Methods(http.MethodPost, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"login", srv.handleLoginUser()).Methods(http.MethodPost, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"login/token", srv.handleRefreshToken()).Methods(http.MethodGet, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"user/info", srv.handleUserInfo()).Methods(http.MethodGet, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"user/infoByID", srv.handleUserInfoByID()).Methods(http.MethodGet, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"user/delete", srv.handleDeleteProfile()).Methods(http.MethodDelete, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"user/password/change", srv.handleChangePassword()).Methods(http.MethodPatch, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"endpoint-info/", srv.handleInfo()).Methods(http.MethodGet)
	srv.router.HandleFunc(baseURL+"email/code", srv.handleSendEmail()).Methods(http.MethodPost, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"password/refresh", srv.handleForgotPassword()).Methods(http.MethodPost, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"password/reset", srv.handleResetPassword()).Methods(http.MethodPatch, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"user/chat/update", srv.handleUpdateChatID()).Methods(http.MethodPatch, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"user/{userID}/chatID", srv.handleGetChatID()).Methods(http.MethodGet, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"user/image/set", srv.handleSetUserImage()).Methods(http.MethodPost, http.MethodOptions)
	srv.router.HandleFunc(baseURL+"user/description/update", srv.handleUpdateDescription()).Methods(http.MethodPatch, http.MethodOptions)
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
					EmailSubscription: false,
					SessionID:         "",
					ChatID:            "",
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

func (srv *service) handleSendEmail() http.HandlerFunc {
	type Request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
			return
		}
		err := utils.SendEmail(req.Email, req.Code)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())
			return
		}
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

func (srv *service) handleGetChatID() http.HandlerFunc {
	type Resp struct {
		ChatID string `json:"chatID"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]

		id, err := uuid.Parse(userID)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, ErrParams, r.Context())
			return
		}
		contacts, err := srv.store.Contacts().GetByUserProfileID(id)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}

		resp := &Resp{
			ChatID: contacts.ChatID,
		}
		userChatID, err := json.Marshal(resp)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}
		_, err = w.Write(userChatID)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}
		srv.respond(w, http.StatusOK, nil)
	}
}

func (srv *service) handleSetUserImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Data struct {
			Original string `json:"original"`
		}

		type Request struct {
			StatusCode  int    `json:"statusCode"`
			Data        Data   `json:"data"`
			AccessToken string `json:"accessToken"`
		}

		req := &Request{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}
		if err := json.Unmarshal(body, req); err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}

		userID := req.AccessToken

		id, err := uuid.Parse(userID)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, ErrParams, r.Context())
		}

		err = user.SetImage(id, req.Data.Original, srv.store)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, ErrParams, r.Context())

			return
		}

		srv.respond(w, http.StatusOK, nil)
	}
}

func (srv *service) handleUserInfo() http.HandlerFunc {
	type Response struct {
		UserID           uuid.UUID `json:"userID"`
		FirstName        string    `json:"firstName"`
		SurName          string    `json:"surName"`
		Email            string    `json:"email"`
		ChatID           string    `json:"chatID"`
		SessionID        string    `json:"sessionID"`
		ImageLink        string    `json:"imageLink"`
		DateRegistration time.Time `json:"date_registration"`
		Description      string    `json:"description"`
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
		profile, err := srv.store.Profile().GetByUserID(id)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}

		temp := ""
		if profile.ImageLink != nil {
			temp = *profile.ImageLink
		}

		resp := &Response{
			UserID:           id,
			FirstName:        profile.FirstName,
			SurName:          profile.SurName,
			Email:            contacts.Email,
			ChatID:           contacts.ChatID,
			SessionID:        contacts.SessionID,
			ImageLink:        temp,
			DateRegistration: profile.DateRegistration,
			Description:      profile.Description,
		}
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
		srv.respond(w, http.StatusOK, nil)
	}
}

func (srv *service) handleUserInfoByID() http.HandlerFunc {

	type Response struct {
		UserID           uuid.UUID `json:"userID"`
		FirstName        string    `json:"firstName"`
		SurName          string    `json:"surName"`
		Email            string    `json:"email"`
		ChatID           string    `json:"chatID"`
		ImageLink        string    `json:"imageLink"`
		DateRegistration time.Time `json:"date_registration"`
		Description      string    `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		resp := &Response{}
		if query.Has("userID") {
			userID := query.Get("userID")
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
			profile, err := srv.store.Profile().GetByUserID(id)
			if err != nil {
				srv.error(w, http.StatusInternalServerError, err, r.Context())

				return
			}

			temp := ""
			if profile.ImageLink != nil {
				temp = *profile.ImageLink
			}

			resp = &Response{
				UserID:           id,
				FirstName:        profile.FirstName,
				SurName:          profile.SurName,
				Email:            contacts.Email,
				ChatID:           contacts.ChatID,
				ImageLink:        temp,
				DateRegistration: profile.DateRegistration,
				Description:      profile.Description,
			}

		}

		if query.Has("chatID") {
			chatID := query.Get("chatID")
			if len(chatID) == 0 {
				srv.warning(w, http.StatusUnauthorized, ErrInvalidHeader)

				return
			}

			contacts, err := srv.store.Contacts().GetByUserChatID(chatID)
			if err != nil {
				srv.error(w, http.StatusInternalServerError, err, r.Context())

				return
			}

			profile, err := srv.store.Profile().GetByUserID(contacts.ProfileID)
			if err != nil {
				srv.error(w, http.StatusInternalServerError, err, r.Context())

				return
			}

			temp := ""
			if profile.ImageLink != nil {
				temp = *profile.ImageLink
			}

			resp = &Response{
				UserID:           contacts.ProfileID,
				FirstName:        profile.FirstName,
				SurName:          profile.SurName,
				Email:            contacts.Email,
				ChatID:           contacts.ChatID,
				ImageLink:        temp,
				DateRegistration: profile.DateRegistration,
				Description:      profile.Description,
			}

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
		srv.respond(w, http.StatusOK, nil)
	}
}

func (srv *service) handleDeleteProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(userIDAuth)
		if len(userID) == 0 {
			srv.warning(w, http.StatusUnauthorized, ErrInvalidHeader)

			return
		}
		err := user.DeleteUserProfile(userID, srv.store)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())
			return
		}

		srv.respond(w, http.StatusNoContent, nil)
	}
}

func (srv *service) handleChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			OldPassword string `json:"oldPassword"`
			NewPassword string `json:"newPassword"`
		}

		userID := r.Header.Get(userIDAuth)
		if len(userID) == 0 {
			srv.warning(w, http.StatusUnauthorized, ErrInvalidHeader)

			return
		}

		profileID, err := uuid.Parse(userID)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, core.ErrParseUUID, r.Context())

			return
		}

		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())

			return
		}

		data, err := srv.store.UserData().GetByUserID(profileID)
		if err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
		}

		passwordValid := register.ComparePassword(data.PasswordEncoded, req.OldPassword, []byte(data.PasswordSalt))
		if !passwordValid {
			srv.error(w, http.StatusBadRequest, err, r.Context())
		}

		if err := user.ChangePassword(data, req.NewPassword, srv.store); err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}
	}
}

func (srv *service) handleUpdateChatID() http.HandlerFunc {
	type Req struct {
		ChatID    string `json:"chatID"`
		SessionID string `json:"sessionID"`
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
		req := &Req{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
			return
		}
		err = user.ChangeChatInfo(id, req.ChatID, req.SessionID, srv.store)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())
			return
		}

		srv.respond(w, http.StatusOK, nil)
	}
}

func (srv *service) handleUpdateDescription() http.HandlerFunc {
	type Req struct {
		Description string `json:"description"`
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
		req := &Req{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
			return
		}
		err = user.ChangeDescriptionInfo(id, req.Description, srv.store)
		if err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())
			return
		}

		srv.respond(w, http.StatusOK, nil)
	}
}

func (srv *service) handleForgotPassword() http.HandlerFunc {
	type Request struct {
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
			return
		}
		contacts, err := srv.store.Contacts().GetByEmail(req.Email)
		if err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
			return
		}

		if err := user.ForgotPassword(contacts, req.Email, srv.store); err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}

	}
}

func (srv *service) handleResetPassword() http.HandlerFunc {
	type Request struct {
		HashID      string `json:"hashID"`
		NewPassword string `json:"newPassword"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
			return
		}

		profileID, err := srv.store.Contacts().GetByHashID(req.HashID)
		if err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
		}

		//data, err := srv.store.Contacts().GetByUserProfileID(req.ProfileID)
		data, err := srv.store.UserData().GetByUserID(profileID.ProfileID)
		if err != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
		}

		if err := user.ChangePassword(data, req.NewPassword, srv.store); err != nil {
			srv.error(w, http.StatusInternalServerError, err, r.Context())

			return
		}

		err2 := srv.store.Contacts().InsertHashID(profileID.ProfileID, "NULL")
		if err2 != nil {
			srv.error(w, http.StatusBadRequest, err, r.Context())
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

	cfg := srv.conf
	domain = cfg.Gateway.IP
	gatewayURL, err := url.Parse(
		protocol + "://" + domain + ":" + cfg.Gateway.Port + baseURL + "hello")
	if err != nil {
		return genErr.NewError(err, ErrConnectAPIGateWay, msg, "can't parse ur for endpoint 'hello'")
	}

	srv.Logger.Infof("Connection gateway url %s", gatewayURL.String())

	info := &models.Hello{
		Name:      "user",
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
	domain = srv.conf.Gateway.IP
	gwURL, err := url.Parse(
		protocol + "://" + domain + ":" + srv.conf.Gateway.Port + baseURL + "hello?service=user&port=" + srv.conf.Listen.Port)
	if err != nil {
		return nil, genErr.NewError(err, ErrConnectAPIGateWay)
	}
	srv.Logger.Infof("Connection gateway url %s", gwURL.String())

	return gwURL, nil
}

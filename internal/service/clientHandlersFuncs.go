package service

import (
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"log"
	"net/http"
	"net/url"
)

const (
	protocol = "http"
	msg      = "msg"
	GET      = "GET"
	POST     = "POST"
	PATCH    = "PATCH"
	DELETE   = "DELETE"
)

func (srv *service) registerClientHandlers() {
	srv.router.HandleFunc(baseURL+"helloMessage/", srv.handleHelloMessage()).Methods(GET)
}

func (srv *service) handleHelloMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.respond(w, http.StatusOK, "Hello, it's work!")
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
	if srv.conf.IsLocal {
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

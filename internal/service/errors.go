package service

import "git.friends.com/PetLand/UserService/v2/internal/genErr"

var (
	ErrCantWriteFile     = genErr.New("can't write file")
	ErrCloseResponseBody = genErr.New("can't close response body")
	ErrParams            = genErr.New("error params")
	ErrConnectAPIGateWay = genErr.New("can't connect to API GateWay")
	ErrParseUUID         = genErr.New("could not parse uuid")
	ErrMarshal           = genErr.New("could not marshal data")
	ErrAuthHeaderMissing = genErr.New("'Authorization' header missing")
	ErrInvalidHeader     = genErr.New("header is invalid")
)

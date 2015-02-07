package transfer

import (
	"net/http"
	"strings"

	"github.com/elos/data"
	"github.com/elos/models/user"
)

type Authenticator func(data.Store, *http.Request) (data.Client, bool, error)

var Auth = func(extract func(*http.Request) (string, string)) Authenticator {
	return func(s data.Store, r *http.Request) (data.Client, bool, error) {
		id, key := extract(r)

		if id == "" || key == "" {
			return nil, false, nil
		}
		return user.Authenticate(s, id, key)
	}
}

var HTTPCredentials = func(r *http.Request) (string, string) {
	return "", ""
}

var SocketCredentials = func(r *http.Request) (string, string) {
	tokens := strings.Split(r.Header.Get("Sec-WebSocket-Protocol"), "-")
	// Query Parameter Method of Authentication
	/*
		id := r.FormValue("id")
		key := r.FormValue("key")
	*/
	if len(tokens) != 2 {
		return "", ""
	} else {
		return tokens[0], tokens[1]
	}
}

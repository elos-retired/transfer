package transfer

import (
	"net/http"
	"strings"

	"github.com/elos/data"
	"github.com/elos/models/user"
)

const AuthHeader = "Elos-Auth"
const AuthDelimeter = "-"

type Authenticator func(data.Store, *http.Request) (data.Client, bool, error)

var Auth = func(credentials Credentialer) Authenticator {
	return func(s data.Store, r *http.Request) (data.Client, bool, error) {
		id, key, ok := credentials(r)
		if !ok {
			return nil, false, ErrCredentialsMalformed
		}

		return user.Authenticate(s, id, key)
	}
}

type Credentialer func(*http.Request) (string, string, bool)

var HTTPCredentialer = func(r *http.Request) (string, string, bool) {
	tokens := strings.Split(r.Header.Get(AuthHeader), AuthDelimeter)

	if len(tokens) != 2 {
		return "", "", false
	} else {
		return tokens[0], tokens[1], true
	}
}

var SocketCredentialer = func(r *http.Request) (string, string, bool) {
	tokens := strings.Split(r.Header.Get(WebSocketProtocolHeader), AuthDelimeter)
	// Query Parameter Method of Authentication
	/*
		id := r.FormValue("id")
		key := r.FormValue("key")
	*/
	if len(tokens) != 2 {
		return "", "", false
	} else {
		return tokens[0], tokens[1], true
	}
}

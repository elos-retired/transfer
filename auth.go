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

var Auth = func(credentialer Credentialer) Authenticator {
	return func(s data.Store, r *http.Request) (data.Client, bool, error) {
		id, key, ok := credentialer(r)
		if !ok {
			return nil, false, ErrMalformedCredentials
		}

		return Authenticate(s, id, key)
	}
}

var Authenticate = user.Authenticate

type Credentialer func(*http.Request) (string, string, bool)

var HTTPCredentialer = credentialer(httpTokens)
var SocketCredentialer = credentialer(socketTokens)

func credentialer(t tokenizer) Credentialer {
	return func(r *http.Request) (id string, key string, ok bool) {
		tokens := t(r)

		if len(tokens) < 2 || len(tokens) > 2 {
			return
		}

		id = tokens[0]
		key = tokens[1]

		if id != "" && key != "" {
			ok = true
		}

		return
	}
}

type tokenizer func(*http.Request) []string

func httpTokens(r *http.Request) []string {
	return strings.Split(r.Header.Get(AuthHeader), AuthDelimeter)
}

func socketTokens(r *http.Request) []string {
	return strings.Split(r.Header.Get(WebSocketProtocolHeader), AuthDelimeter)
}

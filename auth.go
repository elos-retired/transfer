package transfer

import (
	"net/http"
	"strings"

	"github.com/elos/data"
	"github.com/elos/models/user"
)

type RequestAuthenticator func(data.Store, *http.Request) (data.Identifiable, bool, error)

func AuthenticateRequest(s data.Store, r *http.Request) (data.Identifiable, bool, error) {
	// Use the WebSocket protocol header to identify and authenticate the user
	id, key := ExtractCredentials(r)

	if id == "" || key == "" {
		return nil, false, nil
	}

	return user.Authenticate(s, id, key)
}

func ExtractCredentials(r *http.Request) (string, string) {
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

package transfer

import (
	"github.com/elos/data"
	gorilla "github.com/gorilla/websocket"
	"net/http"
	"sync"
)

const WebSocketProtocolHeader = "Sec-WebSocket-Protocol"

// A good default upgrader from gorilla/socket
var DefaultWebSocketUpgrader WebSocketUpgrader = NewGorillaUpgrader(1024, 1024, true)

// the utility a route will use to upgrade a request to a websocket
type WebSocketUpgrader interface {
	Upgrade(http.ResponseWriter, *http.Request, data.Identifiable) (SocketConnection, error)
}

// NullUpgrader {{{

type NullUpgrader struct {
	Upgraded   map[*http.Request]bool
	Connection SocketConnection
	Error      error
	m          sync.Mutex
}

func NewNullUpgrader(c SocketConnection) *NullUpgrader {
	return (&NullUpgrader{Connection: c}).Reset()
}

func (u *NullUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, a data.Identifiable) (SocketConnection, error) {
	u.m.Lock()
	defer u.m.Unlock()

	if u.Error != nil {
		return nil, u.Error
	}
	u.Upgraded[r] = true
	return u.Connection, nil
}

func (u *NullUpgrader) Reset() *NullUpgrader {
	u.m.Lock()
	defer u.m.Unlock()

	u.Upgraded = make(map[*http.Request]bool)
	u.Error = nil
	return u
}

func (u *NullUpgrader) SetError(e error) {
	u.m.Lock()
	defer u.m.Unlock()

	u.Error = e
}

// }}}

// Gorilla Upgrader {{{
// wrapper for gorillla upgrader
type GorillaUpgrader struct {
	Upgrader *gorilla.Upgrader
}

func (u *GorillaUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, a data.Identifiable) (SocketConnection, error) {
	wc, err := u.Upgrader.Upgrade(w, r, ExtractProtocolHeader(r))

	if err != nil {
		return NewNullConnection(a), err
	}

	return NewGorillaConnection(wc, a), nil
}

func NewGorillaUpgrader(rbs int, wbs int, checkO bool) *GorillaUpgrader {
	u := &gorilla.Upgrader{
		ReadBufferSize:  rbs,
		WriteBufferSize: wbs,
		CheckOrigin: func(r *http.Request) bool {
			return checkO
		},
	}
	return &GorillaUpgrader{
		Upgrader: u,
	}
}

// Gorilla Upgrader }}}

func ExtractProtocolHeader(r *http.Request) http.Header {
	header := http.Header{}
	header.Add(WebSocketProtocolHeader, r.Header.Get(WebSocketProtocolHeader))
	return header
}

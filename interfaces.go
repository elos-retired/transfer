package transfer

import (
	"errors"
	"github.com/elos/data"
)

var ConnectionClosedError = errors.New("SocketConnection is closed")

type JReader interface {
	ReadJSON(interface{}) error
}

type JWriter interface {
	WriteJSON(interface{}) error
}

type AnonConnection interface {
	JWriter
}

type Connection interface {
	AnonConnection
	Agent() data.Identifiable
}

type AnonSocketConnection interface {
	AnonConnection
	JReader
	Close() error
}

type SocketConnection interface {
	AnonSocketConnection
	Agent() data.Identifiable
}

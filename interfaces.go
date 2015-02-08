package transfer

import "github.com/elos/data"

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
	data.Client
}

type AnonSocketConnection interface {
	AnonConnection
	JReader
	Close() error
}

type SocketConnection interface {
	AnonSocketConnection
	data.Client
}

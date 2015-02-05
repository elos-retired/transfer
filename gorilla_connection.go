package transfer

import (
	"github.com/elos/data"
)

type GorillaConnection struct {
	conn  AnonSocketConnection
	agent data.Identifiable
}

func NewGorillaConnection(c AnonSocketConnection, a data.Identifiable) SocketConnection {
	return &GorillaConnection{
		conn:  c,
		agent: a,
	}
}

func (c *GorillaConnection) WriteJSON(v interface{}) error {
	return c.conn.WriteJSON(v)
}

func (c *GorillaConnection) ReadJSON(v interface{}) error {
	return c.conn.ReadJSON(v)
}

func (c *GorillaConnection) Close() error {
	return c.conn.Close()
}

func (c *GorillaConnection) Agent() data.Identifiable {
	return c.agent
}

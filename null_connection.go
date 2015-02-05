package transfer

import (
	"github.com/elos/data"
	"sync"
)

/*
	NullConnection implements Connection - mostly for testing
		Writes: Record of interfaces written
		Reads: Record f interfaces read
		Closed: Defaults to false, becomes true if Close() is called
		Error: Error to return, defaults to nil, and if nil, no error return
		agent: Data agent
		m: Mutex for thread-safety
*/
type NullConnection struct {
	Writes map[interface{}]bool
	Reads  map[interface{}]bool
	Closed bool
	Error  error
	agent  data.Identifiable
	m      sync.Mutex

	LastWrite interface{}
}

// Allocates and returns a new *NullConnection
func NewNullConnection(a data.Identifiable) *NullConnection {
	return (&NullConnection{agent: a}).Reset()
}

/*
	Resets:
	- Writes -> emtpy,
	- Reads -> empty,
	- Closed -> false,
	- Error -> nil
*/
func (c *NullConnection) Reset() *NullConnection {
	c.m.Lock()
	defer c.m.Unlock()

	c.Writes = make(map[interface{}]bool)
	c.Reads = make(map[interface{}]bool)
	c.Closed = false
	c.Error = nil
	return c
}

func (c *NullConnection) WriteJSON(v interface{}) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.Error != nil {
		return c.Error
	}

	if c.Closed {
		return ConnectionClosedError
	}

	c.LastWrite = v

	c.Writes[v] = true
	return nil
}

func (c *NullConnection) ReadJSON(v interface{}) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.Error != nil {
		return c.Error
	}

	if c.Closed {
		return ConnectionClosedError
	}

	c.Reads[v] = true
	return nil
}

func (c *NullConnection) Close() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.Error != nil {
		return c.Error
	}

	c.Closed = true
	return nil
}

func (c *NullConnection) Agent() data.Identifiable {
	c.m.Lock()
	defer c.m.Unlock()

	return c.agent
}

func (c *NullConnection) SetError(e error) {
	c.m.Lock()
	defer c.m.Unlock()

	c.Error = e
}

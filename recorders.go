package transfer

import (
	"net/http"
	"sync"

	"github.com/elos/data"
)

// RecorderConnection {{{

/*
	RecorderConnection implements Connection - mostly for testing
		Writes: Record of interfaces written
		Reads: Record f interfaces read
		Closed: Defaults to false, becomes true if Close() is called
		Error: Error to return, defaults to nil, and if nil, no error return
		agent: Data agent
		m: Mutex for thread-safety
*/
type RecorderConnection struct {
	data.Client
	sync.Mutex

	LastWrite interface{}
	Writes    map[interface{}]bool
	Reads     map[interface{}]bool
	Closed    bool
	Error     error
}

// Allocates and returns a new *RecorderConnection
func NewRecorderConnection(c data.Client) *RecorderConnection {
	return (&RecorderConnection{Client: c}).Reset()
}

func (c *RecorderConnection) SetError(e error) {
	c.Lock()
	defer c.Unlock()

	c.Error = e
}

/*
	Resets:
	- Writes -> empty,
	- Reads -> empty,
	- Closed -> false,
	- Error -> nil
*/
func (c *RecorderConnection) Reset() *RecorderConnection {
	c.Lock()
	defer c.Unlock()

	c.Writes = make(map[interface{}]bool)
	c.Reads = make(map[interface{}]bool)
	c.Closed = false
	c.Error = nil

	return c
}

func (c *RecorderConnection) WriteJSON(v interface{}) error {
	c.Lock()
	defer c.Unlock()

	if c.Error != nil {
		return c.Error
	} else if c.Closed {
		return ErrConnectionClosed
	} else {
		c.LastWrite = v
		c.Writes[v] = true
		return nil
	}
}

func (c *RecorderConnection) ReadJSON(v interface{}) error {
	c.Lock()
	defer c.Unlock()

	if c.Error != nil {
		return c.Error
	} else if c.Closed {
		return ErrConnectionClosed
	} else {
		c.Reads[v] = true
		return nil
	}
}

func (c *RecorderConnection) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.Error != nil {
		return c.Error
	}

	c.Closed = true
	return nil
}

// RecorderConnection }}}

// RecorderUpgrader {{{

type RecorderUpgrader struct {
	Upgraded   map[*http.Request]bool
	Connection SocketConnection
	Error      error
	m          sync.Mutex
}

func NewRecorderUpgrader(c SocketConnection) *RecorderUpgrader {
	return (&RecorderUpgrader{Connection: c}).Reset()
}

func (u *RecorderUpgrader) Reset() *RecorderUpgrader {
	u.m.Lock()
	defer u.m.Unlock()

	u.Upgraded = make(map[*http.Request]bool)
	u.Error = nil
	return u
}

func (u *RecorderUpgrader) SetError(e error) {
	u.m.Lock()
	defer u.m.Unlock()

	u.Error = e
}

func (u *RecorderUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, a data.Identifiable) (SocketConnection, error) {
	u.m.Lock()
	defer u.m.Unlock()

	if u.Error != nil {
		return nil, u.Error
	}

	u.Upgraded[r] = true
	return u.Connection, nil
}

// TestUprader }}}

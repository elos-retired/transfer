package transfer

import (
	"github.com/elos/data"
)

type Action string

const (
	POST   Action = "POST"
	GET    Action = "GET"
	DELETE Action = "DELETE"
	SYNC   Action = "SYNC"
	ECHO   Action = "ECHO"
)

// Inbound
type Envelope struct {
	Connection `json:"-"`
	Action     Action                     `json:"action"`
	Data       map[data.Kind]data.AttrMap `json:"data"`
}

func New(c Connection, a Action, k data.Kind, attrs data.AttrMap) *Envelope {
	return &Envelope{
		Connection: c,
		Action:     a,
		Data: map[data.Kind]data.AttrMap{
			k: attrs,
		},
	}
}

func NewEnvelope(a Action, data map[data.Kind]data.AttrMap) *Envelope {
	return &Envelope{
		Action: a,
		Data:   data,
	}
}

// Outbound
type Package struct {
	Action Action       `json:"action"`
	Data   data.KindMap `json:"data"`
}

func NewPackage(a Action, data map[data.Kind]data.Record) *Package {
	return &Package{
		Action: a,
		Data:   data,
	}
}

package transfer

import (
	"github.com/elos/data"
)

func EchoHandler(e *Envelope, db data.DB) {
	e.Connection.WriteJSON(e)
}

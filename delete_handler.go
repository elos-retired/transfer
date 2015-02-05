package transfer

import (
	"github.com/elos/data"
)

/*
	Takes a well-formed envelope, a database and a connection
	and attempts to remove that record from the database.

	Successful removal prompts a direct data.DELETE response

	Unsuccessful removal prompts a direct POST response
	containing the record in question
*/
func DeleteHandler(e *Envelope, s data.Store) {
	var (
		kind data.Kind
		info data.AttrMap
	)

	for kind, info = range e.Data {
		m, err := s.Unmarshal(kind, info)
		if err != nil {
			return
		}

		e.Connection.WriteJSON(NewPackage(DELETE, Map(m)))
	}
}

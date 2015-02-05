package transfer

import "github.com/elos/data"

// GET {{{

func get(e *Envelope, s data.Store) {
	// kind is s.Kind
	// info is map[string]interface{}
	for kind, info := range e.Data {
		m, _ := s.Unmarshal(kind, info)

		err := s.PopulateByID(m)

		if err != nil {
			if err == data.ErrNotFound {
				e.Connection.WriteJSON(APIError{404, 404, "Not Found", "Bad id?"})
				return
			}
			// Otherwise we don't know
			e.Connection.WriteJSON(APIError{400, 400, "Oh shit", ""})
			return
		}

		e.Connection.WriteJSON(m)
	}
}

// GET }}}

// POST {{{

func post(e *Envelope, s data.Store) {
	var (
		kind  data.Kind
		attrs data.AttrMap
	)

	for kind, attrs = range e.Data {
		m, err := s.Unmarshal(kind, attrs)

		if err != nil {
			e.Connection.WriteJSON(APIError{400, 400, "Error", "error"})
			return
		}

		if !m.ID().Valid() {
			m.SetID(s.NewID())
		}

		err = s.Save(m)

		if err != nil {
			e.Connection.WriteJSON(APIError{400, 400, "Error saving the model", "Check yoself"})
			return
		}

		e.Connection.WriteJSON(m)
	}
}

// POST }}}

// DELETE {{{

/*
	Takes a well-formed envelope, a database and a connection
	and attempts to remove that record from the database.

	Successful removal prompts a direct data.DELETE response

	Unsuccessful removal prompts a direct POST response
	containing the record in question
*/
func del(e *Envelope, s data.Store) {
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

// DELETE }}}

// ECHO {{{

func echo(e *Envelope, db data.DB) {
	e.Connection.WriteJSON(e)
}

// ECHO }}}

// SYNC {{{

func syn(e *Envelope, db data.DB) {
	// not implemented
}

// SYNC }}}

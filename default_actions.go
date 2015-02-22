package transfer

import "github.com/elos/data"

// GET {{{

var GetRoute = func(e *Envelope, a data.Access) {
	for kind, info := range e.Data {
		m, err := a.Unmarshal(kind, info)
		if err != nil {
			e.WriteJSON(ErrGeneric)
			return
		}

		if err = a.PopulateByID(m); err != nil {
			switch err {
			case data.ErrAccessDenial:
				e.WriteJSON(ErrAccess)
			case data.ErrNotFound:
				e.WriteJSON(ErrNotFound)
			default:
				e.WriteJSON(ErrGeneric)
			}

			return
		}

		e.WriteJSON(NewPackage(POST, data.Map(m)))
	}
}

// GET }}}

// POST {{{

var PostRoute = func(e *Envelope, a data.Access) {
	for kind, attrs := range e.Data {
		m, err := a.Unmarshal(kind, attrs)
		if err != nil {
			e.WriteJSON(ErrGeneric)
			return
		}

		if !m.ID().Valid() {
			m.SetID(a.NewID())
		}

		if err = a.Save(m); err != nil {
			switch err {
			case data.ErrAccessDenial:
				e.WriteJSON(ErrAccess)
			default:
				e.WriteJSON(ErrGeneric)
			}
			return
		}

		e.WriteJSON(NewPackage(POST, data.Map(m)))
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
var DeleteRoute = func(e *Envelope, a data.Access) {
	for kind, info := range e.Data {
		m, err := a.Unmarshal(kind, info)
		if err != nil {
			e.WriteJSON(ErrGeneric)
			return
		}

		if err = a.Delete(m); err != nil {
			switch err {
			case data.ErrAccessDenial:
				e.WriteJSON(ErrAccess)
			default:
				e.WriteJSON(ErrGeneric)
			}

			return
		}

		e.WriteJSON(NewPackage(DELETE, data.Map(m)))
	}
}

// DELETE }}}

// ECHO {{{

var EchoRoute = func(e *Envelope, a data.Access) {
	e.WriteJSON(e)
}

// ECHO }}}

// SYNC {{{

var SyncRoute = func(e *Envelope, a data.Access) {
	// not implemented
}

// SYNC }}}

// NoActon {{{

var NoActionRoute = func(e *Envelope) {
	e.WriteJSON(ErrBadMethod)
}

// NoActon }}}

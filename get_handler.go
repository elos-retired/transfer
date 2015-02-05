package transfer

import (
	"github.com/elos/data"
	"github.com/elos/stack/util"
)

func GetHandler(e *Envelope, s data.Store) {
	// kind is s.Kind
	// info is map[string]interface{}
	for kind, info := range e.Data {
		m, _ := s.Unmarshal(kind, info)

		err := s.PopulateByID(m)

		if err != nil {
			if err == data.ErrNotFound {
				e.Connection.WriteJSON(util.ApiError{404, 404, "Not Found", "Bad id?"})
				return
			}
			// Otherwise we don't know
			e.Connection.WriteJSON(util.ApiError{400, 400, "Oh shit", ""})
			return
		}

		e.Connection.WriteJSON(m)
	}
}

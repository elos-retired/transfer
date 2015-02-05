package transfer

import (
	"github.com/elos/data"
	"github.com/elos/stack/util"
)

func PostHandler(e *Envelope, s data.Store) {
	var (
		kind  data.Kind
		attrs data.AttrMap
	)

	for kind, attrs = range e.Data {
		m, err := s.Unmarshal(kind, attrs)

		if err != nil {
			e.Connection.WriteJSON(util.ApiError{400, 400, "Error", "error"})
			return
		}

		if !m.ID().Valid() {
			m.SetID(s.NewID())
		}

		err = s.Save(m)

		if err != nil {
			e.Connection.WriteJSON(util.ApiError{400, 400, "Error saving the model", "Check yoself"})
			return
		}

		e.Connection.WriteJSON(m)
	}
}

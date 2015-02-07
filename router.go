package transfer

import (
	"github.com/elos/data"
)

type Router interface {
	Route(e *Envelope, s data.Store)
}

type ActionRouter struct {
	actions map[Action]Router
}

func NewActionRouter() *ActionRouter {
	return &ActionRouter{
		actions: make(map[Action]Router),
	}
}

func (r *ActionRouter) Route(e *Envelope, s data.Store) {
	// Overrides
	router, ok := r.actions[e.Action]
	if ok {
		go router.Route(e, s)
		return
	}

	// Defaults
	switch e.Action {
	case POST:
		go post(e, s)
	case GET:
		go get(e, s)
	case DELETE:
		go del(e, s)
	case SYNC:
		go synchronize(e, s)
	case ECHO:
		go echo(e, s)
	default:
		go e.Connection.WriteJSON(NewInvalidMethodError())
	}
}

var DefaultRouter = NewActionRouter()

func Route(e *Envelope, s data.Store) {
	DefaultRouter.Route(e, s)
}

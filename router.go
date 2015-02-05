package transfer

import (
	"github.com/elos/data"
	"github.com/elos/stack/util"
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
		go PostHandler(e, s)
	case GET:
		go GetHandler(e, s)
	case DELETE:
		go DeleteHandler(e, s)
	case SYNC:
		go SyncHandler(e, s)
	case ECHO:
		go EchoHandler(e, s)
	default:
		go e.Connection.WriteJSON(util.NewInvalidMethodError())
	}
}

var DefaultRouter = NewActionRouter()

func Route(e *Envelope, s data.Store) {
	DefaultRouter.Route(e, s)
}

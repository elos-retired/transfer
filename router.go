package transfer

import (
	"github.com/elos/data"
)

type Router interface {
	Route(e *Envelope, s *data.Access)
}

type ActionRouter struct {
	actions map[Action]Router
}

func NewActionRouter() *ActionRouter {
	return &ActionRouter{
		actions: make(map[Action]Router),
	}
}

func (r *ActionRouter) Route(e *Envelope, a *data.Access) {
	// Overrides
	router, ok := r.actions[e.Action]
	if ok {
		go router.Route(e, a)
		return
	}

	// Defaults
	switch e.Action {
	case POST:
		go PostRoute(e, a)
	case GET:
		go GetRoute(e, a)
	case DELETE:
		go DeleteRoute(e, a)
	case SYNC:
		go SyncRoute(e, a)
	case ECHO:
		go EchoRoute(e, a)
	default:
		go NoActionRoute(e)
	}
}

var DefaultRouter = NewActionRouter()

func Route(e *Envelope, a *data.Access) {
	DefaultRouter.Route(e, a)
}

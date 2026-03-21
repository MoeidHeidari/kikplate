package routes

import (
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/lib"
)

type BadgeRoutes struct {
	handler handlers.BadgeHandler
	rh      lib.RequestHandler
}

func NewBadgeRoutes(handler handlers.BadgeHandler, rh lib.RequestHandler) BadgeRoutes {
	return BadgeRoutes{handler: handler, rh: rh}
}

func (r BadgeRoutes) Setup() {
	r.rh.Mux.Get("/badges", r.handler.List)
}

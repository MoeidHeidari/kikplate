package routes

import (
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/lib"
)

type ConfigRoutes struct {
	handler handlers.ConfigHandler
	rh      lib.RequestHandler
}

func NewConfigRoutes(handler handlers.ConfigHandler, rh lib.RequestHandler) ConfigRoutes {
	return ConfigRoutes{handler: handler, rh: rh}
}

func (r ConfigRoutes) Setup() {
	r.rh.Mux.Get("/config", r.handler.GetConfig)
}

package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/lib"
)

type OrganizationRoutes struct {
	handler handlers.OrganizationHandler
	rh      lib.RequestHandler
}

func NewOrganizationRoutes(handler handlers.OrganizationHandler, rh lib.RequestHandler) OrganizationRoutes {
	return OrganizationRoutes{handler: handler, rh: rh}
}

func (r OrganizationRoutes) Setup() {
	r.rh.Mux.Route("/organizations", func(m chi.Router) {
		m.Get("/", r.handler.ListPublic)
		m.Post("/", r.handler.Create)

		m.Get("/me", r.handler.ListMine)
		m.Get("/by-name/{name}", r.handler.GetByName)
		m.Get("/{id}", r.handler.GetByID)
		m.Put("/{id}", r.handler.Update)
		m.Delete("/{id}", r.handler.Delete)
	})
}

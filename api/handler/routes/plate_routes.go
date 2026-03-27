package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/lib"
)

type PlateRoutes struct {
	handler handlers.PlateHandler
	rh      lib.RequestHandler
}

func NewPlateRoutes(handler handlers.PlateHandler, rh lib.RequestHandler) PlateRoutes {
	return PlateRoutes{handler: handler, rh: rh}
}

func (r PlateRoutes) Setup() {
	r.rh.Mux.Route("/plates", func(m chi.Router) {
		m.Get("/", r.handler.List)
		m.Get("/stats", r.handler.Stats)
		m.Get("/filters", r.handler.FilterOptions)
		m.Get("/bookmarked", r.handler.ListBookmarked)
		m.Post("/repository", r.handler.SubmitRepository)

		m.Get("/{slug}", r.handler.GetBySlug)
		m.Patch("/{id}", r.handler.Update)
		m.Post("/{id}/verify", r.handler.VerifyRepository)
		m.Put("/{id}/organization", r.handler.MoveToOrganization)
		m.Delete("/{id}", r.handler.Archive)
		m.Delete("/{id}/remove", r.handler.Remove)
		m.Post("/{id}/bookmark", r.handler.SetBookmark)
		m.Post("/{id}/reviews", r.handler.SubmitReview)
		m.Put("/{id}/tags", r.handler.ReplaceTags)
		m.Post("/{id}/approve", r.handler.Approve)
		m.Post("/{id}/reject", r.handler.Reject)
		m.Post("/{id}/badges", r.handler.GrantBadge)
		m.Delete("/{id}/badges/{slug}", r.handler.RevokeBadge)
	})
}

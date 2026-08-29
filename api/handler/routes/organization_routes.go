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
		m.Get("/invitations/me", r.handler.ListMyInvitations)
		m.Post("/invitations/{invitationId}/accept", r.handler.AcceptInvitation)
		m.Post("/invitations/{invitationId}/decline", r.handler.DeclineInvitation)
		m.Get("/by-name/{name}", r.handler.GetByName)
		m.Get("/{id}", r.handler.GetByID)
		m.Get("/{id}/github/connect", r.handler.ConnectGitHub)
		m.Get("/{id}/members", r.handler.ListMembers)
		m.Get("/{id}/invitations", r.handler.ListInvitations)
		m.Delete("/{id}/members/{accountId}", r.handler.RemoveMember)
		m.Delete("/{id}/invitations/{invitationId}", r.handler.RevokeInvitation)
		m.Post("/{id}/invitations", r.handler.InviteMember)
		m.Post("/{id}/leave", r.handler.Leave)
		m.Put("/{id}", r.handler.Update)
		m.Delete("/{id}", r.handler.Delete)
	})
}

package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/lib"
)

type AuthRoutes struct {
	logger  lib.Logger
	env     lib.Env
	handler handlers.AuthHandler
	mux     lib.RequestHandler
}

func NewAuthRoutes(
	logger lib.Logger,
	env lib.Env,
	handler handlers.AuthHandler,
	mux lib.RequestHandler,
) AuthRoutes {
	return AuthRoutes{
		logger:  logger,
		env:     env,
		handler: handler,
		mux:     mux,
	}
}

func (r AuthRoutes) Setup() {
	r.mux.Mux.Route("/auth", func(router chi.Router) {
		router.Post("/register", r.handler.Register)
		router.Get("/verify-email", r.handler.VerifyEmail)
		router.Post("/login", r.handler.LoginLocal)
		router.Get("/{provider}/redirect", r.handler.OAuthRedirect)
		router.Get("/{provider}/callback", r.handler.OAuthCallback)
	})

	// protected routes
	r.mux.Mux.Group(func(router chi.Router) {
		router.Use(middleware.RequireAuth)
		router.Get("/me", r.handler.Me)
	})
}
